package config

import (
	"os"
	"strings"

	// this will automatically load your .env file:

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DATABASE_URL        string
	JWT_SECRET          string
	GoogleLoginConfig   oauth2.Config
	REDIS_DB_URL        string
	GitlabLoginConfig   oauth2.Config
	GitlabLoginEndpoint string
	FRONTEND_URL        string
	BACKEND_URL         string
	EnableGoogleOauth   bool
	EnableGitlabOauth   bool
}

type Handler struct {
	DB        *pgxpool.Pool
	Config    *Config
	RDB       *redis.Client
	Validator *validator.Validate
}

func LoadConfig() (*Config, error) {

	GITLAB_BASEURL := os.Getenv("GITLAB_BASEURL")
	FRONTEND_URL := os.Getenv("FRONTEND_URL")
	BACKEND_URL := os.Getenv("BACKEND_URL")
	if BACKEND_URL == "" {
		BACKEND_URL = "http://localhost:4000"
	}
	if FRONTEND_URL == "" {
		FRONTEND_URL = "http://localhost:4000"
	}
	if GITLAB_BASEURL == "" {
		GITLAB_BASEURL = "https://gitlab.com"
	}
	cfg := &Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		GoogleLoginConfig: oauth2.Config{
			RedirectURL:  BACKEND_URL + "/api/users/oauth/google/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint: google.Endpoint,
		},
		REDIS_DB_URL: os.Getenv("REDIS_URL"),
		GitlabLoginConfig: oauth2.Config{
			RedirectURL:  BACKEND_URL + "/api/users/oauth/gitlab/callback",
			ClientID:     os.Getenv("GITLAB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITLAB_CLIENT_SECRET"),
			Scopes:       []string{"read_user"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  GITLAB_BASEURL + "/oauth/authorize",
				TokenURL: GITLAB_BASEURL + "/oauth/token",
			},
		},
		GitlabLoginEndpoint: GITLAB_BASEURL,
		FRONTEND_URL:        FRONTEND_URL,
		EnableGoogleOauth:   strings.ToLower(os.Getenv("ENABLE_GOOGLE_OAUTH")) == "true",
		EnableGitlabOauth:   strings.ToLower(os.Getenv("ENABLE_GITLAB_OAUTH")) == "true",
	}

	return cfg, nil
}
