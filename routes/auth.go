package routes

import (
	"gitleet/services"
	"gitleet/structs"
	"gitleet/utils"
	"net/http"
	"os"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func AuthRoutes(app fiber.Router, h *Handlers) {
	auth := NewAuthHandlers(h.Logger, h.Db, h.Context, h.Services)
	authRoute := app.Group("/auth")
	authRoute.Get("/github", auth.StartGithubOAuth)
	authRoute.Get("/callback", auth.CompleteGithubOAuth)
	authRoute.Get("/refresh", auth.Refresh)
}

func (h *AuthHandlers) StartGithubOAuth(ctx *fiber.Ctx) error {
	url := h.githubOauthConfig.AuthCodeURL(h.state)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (h *AuthHandlers) CompleteGithubOAuth(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if state != "some_random_state" {
		return ctx.Status(400).JSON(&fiber.Map{"error": "Invalid state"})
	}

	token, err := h.githubOauthConfig.Exchange(h.Context, code)
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{"error": err.Error()})
	}
	// get user details
	client := services.GetGitHubClient(token.AccessToken)
	githubUser, _, err := client.Users.Get(h.Context, "")
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{"error": err.Error()})
	}

	_, err = h.Services.DBService.CreateOrUpdateUser(githubUser)
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{"error": err.Error()})
	}

	cookie := utils.GetCookie(token.AccessToken, os.Getenv("REFRESH_COOKIE_NAME"))
	ctx.Cookie(cookie)

	redirectURL := os.Getenv("REDIRECT_FRONTEND_URL")

	return ctx.Redirect(redirectURL, fiber.StatusPermanentRedirect)
}

type RefreshTokenResponse struct {
	User        *structs.User `json:"user"`
	AccessToken string        `json:"access_token"`
}

func (h *AuthHandlers) Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies(os.Getenv("REFRESH_COOKIE_NAME"))
	// Use the RefreshToken method of the Config to create a TokenSource.
	tokenSource := h.githubOauthConfig.TokenSource(h.Context, &oauth2.Token{
		AccessToken: refreshToken,
	})

	// Obtain a new access token using the TokenSource.
	token, err := tokenSource.Token()
	if err != nil {
		h.Logger.Println(err.Error())
		return ctx.Status(500).JSON(&fiber.Map{"error": err.Error()})
	}

	// Update the user details in the database with the new access token.
	client := services.GetGitHubClient(token.AccessToken)
	githubUser, _, err := client.Users.Get(h.Context, "")
	if err != nil {
		h.Logger.Println(err.Error())
		return ctx.Status(500).JSON(&fiber.Map{"error": err.Error()})
	}

	user, err := h.Services.DBService.CreateOrUpdateUser(githubUser)
	if err != nil {
		h.Logger.Println("")
		return ctx.Status(500).JSON(&fiber.Map{"error": err.Error()})
	}

	// // Update the refresh token in the cookie.
	cookie := utils.GetCookie(token.AccessToken, os.Getenv("REFRESH_COOKIE_NAME"))
	ctx.Cookie(cookie)

	response := &RefreshTokenResponse{
		User:        user,
		AccessToken: token.AccessToken,
	}

	return ctx.Status(200).JSON(response)
}

func (h *Handlers) AuthMiddleware(ctx *fiber.Ctx) error {
	accessToken := ctx.Get("Authorization")

	//for now this is enough
	if accessToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"error": "Unauthorized - Access Token is required"})
	}

	ctx.Locals("token", accessToken)

	return ctx.Next()
}
