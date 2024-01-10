package routes

import (
	"gitleet/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/github"
)

func GithubRoutes(app fiber.Router, h *Handlers) {
	githubRoute := app.Group("/github")
	githubRoute.Get("/repo", h.AuthMiddleware, h.GetGithubRepo)
	githubRoute.Get("/branch/:username/:repoID", h.AuthMiddleware, h.GetBranch)
	githubRoute.Get("/repo/:username/:term", h.AuthMiddleware, h.SearchRepo)
	githubRoute.Patch("/repo/update", h.AuthMiddleware, h.UpdateRepo)
}

func (h *Handlers) GetGithubRepo(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(string)
	client := services.GetGitHubClient(token)
	repos, _, err := client.Repositories.List(h.Context, "", nil)
	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"error": "Error fetching Repositories",
		})
	}
	return ctx.Status(200).JSON(repos)
}

func (h *Handlers) SearchRepo(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(string)
	username := ctx.Params("username")
	term := ctx.Params("term")
	if len(username) <= 0 || len(term) <= 0 {
		return ctx.Status(400).JSON(&fiber.Map{
			"error": "No search term",
		})
	}

	options := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	searchQuery := "user:" + username + " " + term

	client := services.GetGitHubClient(token)
	result, _, err := client.Search.Repositories(h.Context, searchQuery, options)
	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(result)

}

func (h *Handlers) GetBranch(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(string)
	username := ctx.Params("username")
	repoID := ctx.Params("repoID")

	client := services.GetGitHubClient(token)

	branches, _, err := client.Repositories.ListBranches(h.Context, username, repoID, nil)
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(branches)
}

func (h *Handlers) UpdateRepo(ctx *fiber.Ctx) error {
	type UpdateGithub struct {
		UserID           string `json:"userId"`
		GithubRepoName   string `json:"github_repo"`
		GithubRepoBranch string `json:"github_repo_branch"`
	}

	res := new(UpdateGithub)

	if err := ctx.BodyParser(&res); err != nil {
		h.Logger.Println(err.Error())
		return ctx.Status(400).JSON(&fiber.Map{
			"error": "Invalid request body",
		})
	}
	err := h.Services.DBService.UpdateGithub(res.GithubRepoName, res.GithubRepoBranch, res.UserID)
	if err != nil {
		h.Logger.Println(err.Error())
		return ctx.Status(400).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message":            "Updated Successfully",
		"github_repo":        res.GithubRepoName,
		"github_repo_branch": res.GithubRepoBranch,
	})

}
