package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	config.Handler
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	challenge := c.QueryParam("challenge")
	state := generateRandomState()
	setSecureStateCookie(c, state)
	url := h.Config.GoogleLoginConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	stateFromCookie := getSecureStateCookie(c)
	if state != stateFromCookie {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid state parameter", nil))
	}
	redirectURL := fmt.Sprintf("%s/users/login/callback?provider=google&code=%s&state=%s", h.Config.FRONTEND_URL, code, state)
	return c.Redirect(http.StatusFound, redirectURL)
}

func (h *AuthHandler) GoogleExchange(c echo.Context) error {
	var data oauthExchanceData
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}
	token, err := h.Config.GoogleLoginConfig.Exchange(
		c.Request().Context(),
		data.Code,
		oauth2.SetAuthURLParam("code_verifier", data.CodeVerifier))
	if err != nil {
		return err
	}
	client := h.Config.GoogleLoginConfig.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return err
	}

	// Check if user exists in the database
	sql_fetched_user := &User{}
	sql, _, _ := goqu.From("users").Where(goqu.C("email").Eq(userInfo.Email)).Select(goqu.COUNT("*")).ToSQL()
	var count int
	row := h.DB.QueryRow(c.Request().Context(), sql)
	err = row.Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to check user existence", err))
	}
	if count == 0 {
		sql, _, _ = goqu.Insert("users").
			Cols("username", "name", "email", "role").
			Vals(goqu.Vals{userInfo.Email, userInfo.Name, userInfo.Email, "user"}).
			Returning("name", "username", "email", "role", "id").ToSQL()
		row = h.DB.QueryRow(c.Request().Context(), sql)
		err = row.Scan(&sql_fetched_user.Name, &sql_fetched_user.Username, &sql_fetched_user.Email, &sql_fetched_user.Role, &sql_fetched_user.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to fetch user after creation", err))
		}
	} else {
		sql, _, _ = goqu.From("users").
			Where(goqu.C("email").Eq(userInfo.Email)).Select("name", "username", "email", "role", "id").ToSQL()
		row = h.DB.QueryRow(c.Request().Context(), sql)
		err = row.Scan(&sql_fetched_user.Name, &sql_fetched_user.Username, &sql_fetched_user.Email, &sql_fetched_user.Role, &sql_fetched_user.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to fetch user", err))
		}
	}

	expiry := time.Now().Add(time.Hour * 72)
	// Generate JWT token
	claims := &JwtCustomClaims{
		sql_fetched_user.Name,
		sql_fetched_user.Username,
		sql_fetched_user.Role,
		sql_fetched_user.Id,
		sql_fetched_user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.Config.JWT_SECRET))
	if err != nil {
		return err
	}

	// Set JWT token as cookie
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Login successful",
		"status":       "success",
		"access_token": tokenString,
		"expiry":       expiry,
	})
}
