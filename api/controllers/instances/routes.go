package instances

import (
	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	instanceGroup := g.Group("")
	instanceGroup.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	useInstanceRoutes(instanceGroup, db, cfg)
}

func UseOAuthServerSubroute(g *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	h := &instanceHandler{config.Handler{DB: db, Config: cfg, RDB: rdb}}
	g.POST("/device_authorization", h.deviceAuthorization)
	g.POST("/token", h.token)
	verify := g.Group("")
	verify.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	verify.POST("", h.VerifyUserCode)
}

func useInstanceRoutes(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {

	h := &instanceHandler{config.Handler{DB: db, Config: cfg}}
	g.GET("", h.getInstances)
	g.GET("/view/:id", h.getInstance)
	g.POST("", h.createInstance)
	g.PUT("/edit/:id", h.editInstance)
	g.DELETE("/:id", h.deleteInstance)

	g.GET("/users/:id", h.getInstanceUsers)
	g.POST("/users/:id", h.addUserInstanceAccess)
	g.DELETE("/users/:id", h.deleteInstanceUsers)

	g.GET("/roles/:id", h.getInstanceRoles)
	g.POST("/roles/:id", h.addInstanceRoles)
	g.DELETE("/roles/:id", h.deleteInstanceRoles)
	// g.GET("/hostUsers/:id", h.getInstanceHostUsers)

	g.Use(h.isAdminOrCreatorMiddleware)
	g.POST("/hostUsers/:id", h.addInstanceHostUser)
	g.DELETE("/hostUsers/:id", h.deleteInstanceHostUser)
	g.PUT("/status/:id", h.setStatusInstance)
}
