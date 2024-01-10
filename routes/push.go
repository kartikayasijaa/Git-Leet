package routes

import (
	"fmt"
	"gitleet/services"

	"github.com/gofiber/fiber/v2"
)

func PushRoute(app fiber.Router, h *Handlers) {
	pushRoute := app.Group("/push")
	pushRoute.Get("/:userId", h.AuthMiddleware, h.PushLatest)
}

func (h *Handlers) PushLatest(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")
	token := ctx.Locals("token").(string)
	if len(userId) <= 0 {
		return ctx.Status(400).JSON(&fiber.Map{
			"error": "User Id can not be null",
		})
	}

	// Start a transaction
	tx := h.Services.DBService.DB.Begin()

	// Defer the transaction rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := h.Services.DBService.GetUser(userId)
	if err != nil {
		h.Logger.Println(err.Error())
		tx.Rollback()
		return ctx.Status(400).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	submissions, err := h.Services.LeetcodeService.GetRecentSubmission(user.LeetcodeUsername, user.LeetcodePrevSubmission)
	if err != nil {
		tx.Rollback()
		return ctx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for _, s := range submissions {
		submissionDetail, err := h.Services.LeetcodeService.GetCode(s.Id)
		if err != nil {
			h.Logger.Println(err)
			tx.Rollback()
			return ctx.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = services.PushToGithub(token, user.GithubUsername, user.GithubRepoName, user.GithubRepoBranch, submissionDetail)
		if err != nil {
			h.Logger.Println(err.Error())
			tx.Rollback()
			return ctx.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	newSubmission := len(submissions) + int(user.LeetcodePrevSubmission)
	if err = h.Services.DBService.UpdatePrevSubmission(userId, int32(newSubmission)); err != nil {
		tx.Rollback()
		return ctx.Status(400).JSON(&fiber.Map{
			"error": "error updating previous Submissions",
		})
	}

	// Commit the transaction if everything succeeded
	tx.Commit()

	return ctx.Status(200).JSON(&fiber.Map{
		"message": fmt.Sprintf("Committed %d new Submissions", len(submissions)),
	})
}
