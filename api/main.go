package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/controllers/instances"
	"github.com/Utkar5hM/Execstasy/api/controllers/roles"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/Utkar5hM/Execstasy/api/utils/validations"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	username := claims.Username
	message := fmt.Sprintf("Welcome %s!\nYour Role: %s", username, claims.Role)
	return c.String(http.StatusOK, message)
	// return c.String(http.StatusOK, "Welcome!")
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connection to Database first. :)
	dbpool, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_DB_URL,
		Password: cfg.REDIS_PASSWORD, // no password set
		DB:       cfg.REDIS_DB,       // use default DB
	})
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("usernameregex", validations.UserNameValidator)
	validate.RegisterValidation("nameregex", validations.NameValidator)
	validate.RegisterValidation("linuxuser", validations.LinuxUserValidator)
	validate.RegisterValidation("scopelinuxuser", validations.ScopeLinuxUserValidator)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.MethodOverride())

	// e.Use(middleware.Secure())
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
