package config

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type Handlers struct {
	Logger  *log.Logger
	Db      *gorm.DB
	Context context.Context
}

type AuthHandlers struct {
	Handlers
	githubOauthConfig *oauth2.Config
	state             string
}

func NewHandlers(logger *log.Logger, db *gorm.DB, ctx context.Context) *Handlers {
	return &Handlers{
		Logger:  logger,
		Db:      db,
		Context: ctx,
	}
}

func NewAuthHandlers(logger *log.Logger) *AuthHandlers {
	authHandlers := &AuthHandlers{
		Handlers: Handlers{
			Logger: logger,
		},
	}
	authHandlers.init()
	return authHandlers
}

func (h *AuthHandlers) init() {
	h.githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_CALLBACK_URL"),
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}

	h.state = "some_random_state"
}
