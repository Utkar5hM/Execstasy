package authentication

import (
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &AuthHandler{config.Handler{DB: db, Config: cfg}}
	g.POST("/oauth/google/login", h.GoogleLogin)
	g.GET("/oauth/google/login", h.GoogleLogin)
	g.GET("/oauth/google/callback", h.GoogleCallback)
	protectedGroup := g.Group("")
	protectedGroup.Use(IsLoggedIn(cfg.JWT_SECRET))
	protectedGroup.GET("", h.getUsers)
	// protectedGroup.GET("/search", h.searchUsers)
	protectedGroup.GET("/me", h.getMe)
	protectedGroup.PUT("/me", h.updateMe)
	protectedGroup.PUT("/role", h.IsAdminMiddleware(h.updateRole))
}
