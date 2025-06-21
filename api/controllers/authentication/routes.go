package authentication

import (
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config, v *validator.Validate) {
	h := &AuthHandler{config.Handler{DB: db, Config: cfg, Validator: v}}
	if cfg.EnableGoogleOauth {
		g.POST("/oauth/google/login", h.GoogleLogin)
		g.GET("/oauth/google/login", h.GoogleLogin)
		g.GET("/oauth/google/callback", h.GoogleCallback)
		g.POST("/oauth/google/exchange", h.GoogleExchange)
	}
	if cfg.EnableGitlabOauth {
		g.GET("/oauth/gitlab/login", h.GitLabLogin)
		g.GET("/oauth/gitlab/callback", h.GitLabCallback)
		g.POST("/oauth/gitlab/exchange", h.GitLabExchange)
	}
	protectedGroup := g.Group("")
	protectedGroup.Use(IsLoggedIn(cfg.JWT_SECRET))
	protectedGroup.GET("", h.getUsers)
	// protectedGroup.GET("/search", h.searchUsers)
	protectedGroup.GET("/me", h.getMe)
	protectedGroup.PUT("/me", h.updateMe)
	protectedGroup.PUT("/role", h.IsAdminMiddleware(h.updateRole))
}
