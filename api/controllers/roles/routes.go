package roles

import (
	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config, v *validator.Validate) {
	h := &roleHandler{config.Handler{DB: db, Config: cfg, Validator: v}}
	ah := &authentication.AuthHandler{Handler: config.Handler{DB: db, Config: cfg, Validator: v}}

	l := g.Group("")
	l.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	l.GET("", h.getRoles)
	l.GET("/view/:id", h.getRole)
	l.GET("/me", h.getMyRoles)
	l.GET("/access/:id", h.hasRoleEditAccess)
	l.GET("/users/:id", h.getRoleUsers)

	a := g.Group("")
	a.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	a.Use(h.hasRoleEditAccessMiddleware)

	a.POST("/users/:id", h.AddUserToRole)
	a.DELETE("/users/:id", h.DeleteUserFromRole)

	l.PUT("/edit/:id", ah.IsAdminMiddleware(h.editRole))
	l.POST("", ah.IsAdminMiddleware(h.createRole))
	l.DELETE("/:id", ah.IsAdminMiddleware(h.deleteRole))
}
