package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func generateRandomState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// fallback or panic as appropriate for your app
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func setSecureStateCookie(c echo.Context, state string) {
	cookie := new(http.Cookie)
	cookie.Name = "oauth_state"
	cookie.Value = state
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(5 * time.Minute)

	c.SetCookie(cookie)
}

func getSecureStateCookie(c echo.Context) string {
	cookie, err := c.Cookie("oauth_state")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (h *AuthHandler) GitLabLogin(c echo.Context) error {
	challenge := c.QueryParam("challenge")
	gitlabConfig := h.Config.GitlabLoginConfig
	state := generateRandomState()
	setSecureStateCookie(c, state)
	url := gitlabConfig.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.AccessTypeOnline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GitLabCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	stateFromCookie := getSecureStateCookie(c)
	if state != stateFromCookie {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid state parameter", nil))
	}
	redirectURL := fmt.Sprintf("%s/users/login/callback?provider=gitlab&code=%s&state=%s", h.Config.FRONTEND_URL, code, state)
	return c.Redirect(http.StatusFound, redirectURL)
}

type oauthExchanceData struct {
	CodeVerifier string `json:"code_verifier"`
	Code         string `json:"code"`
	State        string `json:"state"`
}

func (h *AuthHandler) GitLabExchange(c echo.Context) error {
	var data oauthExchanceData
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}
	gitlabConfig := h.Config.GitlabLoginConfig
	token, err := gitlabConfig.Exchange(c.Request().Context(), data.Code, oauth2.SetAuthURLParam("code_verifier", data.CodeVerifier))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to exchange code for token\ncode: "+
			data.Code+"\nverifier: "+data.CodeVerifier, err))
	}

	client := gitlabConfig.Client(c.Request().Context(), token)
	resp, err := client.Get(h.Config.GitlabLoginEndpoint + "/api/v4/user")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to fetch user information", err))
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return err
	}

	sql_fetched_user := &User{}
	sql, _, _ := goqu.From("users").Where(goqu.C("email").Eq(userInfo.Email)).Select(goqu.COUNT("*")).ToSQL()
	var count int
	row := h.DB.QueryRow(c.Request().Context(), sql)
	if err = row.Scan(&count); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to check user existence", err))
	}

	if count == 0 {
		sql, _, _ = goqu.Insert("users").
			Cols("username", "name", "email", "role").
			Vals(goqu.Vals{userInfo.Username, userInfo.Name, userInfo.Email, "user"}).
			Returning("name", "username", "email", "role", "id").ToSQL()
		row = h.DB.QueryRow(c.Request().Context(), sql)
		if err = row.Scan(&sql_fetched_user.Name, &sql_fetched_user.Username, &sql_fetched_user.Email, &sql_fetched_user.Role, &sql_fetched_user.Id); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to fetch user after creation", err))
		}
	} else {
		sql, _, _ = goqu.From("users").
			Where(goqu.C("email").Eq(userInfo.Email)).Select("name", "username", "email", "role", "id").ToSQL()
		row = h.DB.QueryRow(c.Request().Context(), sql)
		if err = row.Scan(&sql_fetched_user.Name, &sql_fetched_user.Username, &sql_fetched_user.Email, &sql_fetched_user.Role, &sql_fetched_user.Id); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to fetch user", err))
		}
	}
	expiry := time.Now().Add(time.Hour * 72)
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

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Login successful",
		"status":       "success",
		"access_token": tokenString,
		"expiry":       expiry,
	})
}
