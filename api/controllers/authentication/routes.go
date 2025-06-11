package authentication

import (
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &authHandler{config.Handler{DB: db, Config: cfg}}
	g.POST("/oauth/google/login", h.GoogleLogin)
	g.GET("/oauth/google/login", h.GoogleLogin)
	g.GET("/oauth/google/callback", h.GoogleCallback)
	protectedGroup := g.Group("")
	protectedGroup.Use(IsLoggedIn(cfg.JWT_SECRET))
	protectedGroup.GET("", h.getUsers)

}

func IsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*JwtCustomClaims)
		if user.Role != "admin" {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func IsAdmin(c echo.Context) bool {
	user := c.Get("user").(*JwtCustomClaims)
	if user.Role != "admin" {
		return true
	}
	return false
}
