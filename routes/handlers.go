package routes

import (
	"log"
	"os"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Handlers struct {
	Logger *log.Logger
}

type AuthHandlers struct {
	Handlers
	githubOauthConfig *oauth2.Config
	state             string
}

func NewHandlers( logger *log.Logger) *Handlers {
	return &Handlers{
		Logger:      logger,
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