package routes

import (
	"context"
	"gitleet/services"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

// Controller Handlers
type Handlers struct {
	Logger   *log.Logger
	Db       *gorm.DB
	Context  context.Context
	Services *services.Services
}

func NewHandlers(logger *log.Logger, db *gorm.DB, ctx context.Context, services *services.Services) *Handlers {
	return &Handlers{
		Logger:   logger,
		Db:       db,
		Context:  ctx,
		Services: services,
	}
}

// Auth Controllers
type AuthHandlers struct {
	Handlers
	githubOauthConfig *oauth2.Config
	state             string
}

func NewAuthHandlers(logger *log.Logger, db *gorm.DB, ctx context.Context, services *services.Services) *AuthHandlers {
	authHandlers := &AuthHandlers{
		Handlers: Handlers{
			Logger:    logger,
			Db:        db,
			Context:   ctx,
			Services: services,
		},
	}
	authHandlers.AuthInit()
	return authHandlers
}

func (h *AuthHandlers) AuthInit() {
	h.githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_CALLBACK_URL"),
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}

	h.state = "some_random_state"
}
