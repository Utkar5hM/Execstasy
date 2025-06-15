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

	protectedGroup.POST("/users/:id", h.AddUserToRole)
	protectedGroup.DELETE("/users/:id", h.DeleteUserFromRole)

	g.GET("", h.getRoles)
	g.GET("/view/:id", h.getRole)
	protectedGroup.PUT("/edit/:id", h.editRole)
	g.GET("/users/:id", h.getRoleUsers)
	protectedGroup.POST("", h.createRole)
	protectedGroup.DELETE("/:id", h.deleteRole)
}
