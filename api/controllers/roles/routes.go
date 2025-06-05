package roles

import (
	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &roleHandler{config.Handler{DB: db, Config: cfg}}
	protectedGroup := g.Group("")
	protectedGroup.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	protectedGroup.Use(h.isAdminMiddleware)
	protectedGroup.POST("", h.createRole)
	protectedGroup.DELETE("", h.deleteRole)
	g.GET("", h.getRoles)
}
