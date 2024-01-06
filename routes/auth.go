package routes

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/github"
)

func AuthRoutes(app fiber.Router, h *Handlers) {
	auth := NewAuthHandlers(h.Logger, h.Db, h.Context, h.DBServices)
	authRoute := app.Group("/auth")
	authRoute.Get("/github", auth.StartGithubOAuth)
	authRoute.Get("/callback", auth.CompleteGithubOAuth)
}

func (h *AuthHandlers) StartGithubOAuth(ctx *fiber.Ctx) error {
	url := h.githubOauthConfig.AuthCodeURL(h.state)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (h *AuthHandlers) CompleteGithubOAuth(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if state != "some_random_state" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid state"})
	}

	token, err := h.githubOauthConfig.Exchange(ctx.Context(), code)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// get user details
	client := github.NewClient(h.githubOauthConfig.Client(ctx.Context(), token))
	githubUser, _, err := client.Users.Get(ctx.Context(), "")
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = h.DBServices.CreateOrUpdateUser(githubUser)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token.AccessToken,
		HTTPOnly: true,
	})

	redirectURL := os.Getenv("REDIRECT_FRONTEND_URL")

	return ctx.Redirect(redirectURL, fiber.StatusPermanentRedirect)
}
