package authentication

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (h *AuthHandler) GitLabLogin(c echo.Context) error {
	url := h.Config.GitlabLoginConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GitLabCallback(c echo.Context) error {
	code := c.QueryParam("code")
	token, err := h.Config.GitlabLoginConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		return err
	}

	client := h.Config.GitlabLoginConfig.Client(c.Request().Context(), token)
	resp, err := client.Get(h.Config.GitlabLoginEndpoint + "/api/v4/user")
	if err != nil {
		return err
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

	claims := &JwtCustomClaims{
		sql_fetched_user.Name,
		sql_fetched_user.Username,
		sql_fetched_user.Role,
		sql_fetched_user.Id,
		sql_fetched_user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.Config.JWT_SECRET))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}
