package authentication

import (
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) IsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		sql, _, _ := goqu.From("users").Select("role").Where(goqu.Ex{"id": claims.Id}).ToSQL()
		var role string
		err := h.DB.QueryRow(c.Request().Context(), sql).Scan(&role)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Access Denied", err))
		}
		if role != "admin" {
			return c.JSON(http.StatusForbidden, helper.ErrorMessage("Access Denied", "You do not have permission to perform this action"))
		}
		return next(c)
	}
}

func IsLoggedIn(jwt_secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey:  []byte(jwt_secret),
		TokenLookup: "header:Authorization:Bearer ,cookie:jwt",
	})
}
