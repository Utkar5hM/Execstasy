package instances

import (
	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config, v *validator.Validate) {
	instanceGroup := g.Group("")
	instanceGroup.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	useInstanceRoutes(instanceGroup, db, cfg, v)
}

func UseOAuthServerSubroute(g *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config, v *validator.Validate) {
	h := &instanceHandler{config.Handler{DB: db, Config: cfg, RDB: rdb, Validator: v}}
	g.POST("/device_authorization", h.deviceAuthorization)
	g.POST("/token", h.token)
	verify := g.Group("")
	verify.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	verify.POST("", h.VerifyUserCode)
}

func useInstanceRoutes(g *echo.Group, db *pgxpool.Pool, cfg *config.Config, v *validator.Validate) {

	h := &instanceHandler{config.Handler{DB: db, Config: cfg, Validator: v}}
	ah := &authentication.AuthHandler{Handler: config.Handler{DB: db, Config: cfg, Validator: v}}
	g.GET("", h.getInstances)
	g.GET("/me", h.getMyInstances)
	g.GET("/view/:id", h.getInstance)

	a := g.Group("")
	a.Use(ah.IsAdminMiddleware)
	a.POST("", h.createInstance)
	a.PUT("/edit/:id", h.editInstance)
	a.DELETE("/:id", h.deleteInstance)

	g.GET("/users/:id", h.getInstanceUsers)
	a.POST("/users/:id", h.addUserInstanceAccess)
	a.DELETE("/users/:id", h.deleteUserInstanceAccess)

	g.GET("/roles/:id", h.getInstanceRoles)
	a.POST("/roles/:id", h.addInstanceRoles)
	a.DELETE("/roles/:id", h.deleteInstanceRoles)
	// g.GET("/hostUsers/:id", h.getInstanceHostUsers)

	a.POST("/hostUsers/:id", h.addInstanceHostUser)
	a.DELETE("/hostUsers/:id", h.deleteInstanceHostUser)
}
