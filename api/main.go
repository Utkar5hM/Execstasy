package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/controllers/instances"
	"github.com/Utkar5hM/Execstasy/api/controllers/roles"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/Utkar5hM/Execstasy/api/utils/validations"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	fmt.Println("API_DEBUG: " + os.Getenv("API_DEBUG"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	opt, _ := redis.ParseURL(cfg.REDIS_DB_URL)
	rdb := redis.NewClient(opt)
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("usernameregex", validations.UserNameValidator)
	validate.RegisterValidation("nameregex", validations.NameValidator)
	validate.RegisterValidation("linuxuser", validations.LinuxUserValidator)
	validate.RegisterValidation("scopelinuxuser", validations.ScopeLinuxUserValidator)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.MethodOverride())
	if os.Getenv("API_DEBUG") != "DEBUG" {
		e.Use(middleware.Secure())
	}
	r := e.Group("/api")
	authGroup := r.Group("/users")
	authentication.UseSubroute(authGroup, dbpool, cfg, validate)

	instanceGroup := r.Group("/instances")
	instances.UseSubroute(instanceGroup, dbpool, rdb, cfg, validate)
	OAuthServerGroup := r.Group("/oauth")
	instances.UseOAuthServerSubroute(OAuthServerGroup, dbpool, rdb, cfg, validate)
	rolesGroup := r.Group("/roles")
	roles.UseSubroute(rolesGroup, dbpool, cfg, validate)
	e.Static("/static/", "static")
	e.Logger.Fatal(e.Start(":4000"))

}
