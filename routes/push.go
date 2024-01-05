package routes

import (
	"fmt"
	"gitleet/services"

	"github.com/gofiber/fiber/v2"
)

func SubmissionRoute(app fiber.Router, h *Handlers) {
	pushRoute := app.Group("/push")
	pushRoute.Get("/:userId", h.PushLatest)
}

func (h *Handlers) PushLatest(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	submissions, err := services.GetRecentSubmission(userId)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for _, s := range submissions {
		submissionDetail, err := services.GetCode(s.Id)
		if err != nil {
			h.Logger.Println(err)
			return ctx.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fmt.Println(submissionDetail.Question.QuestionTitle)
		err = services.PushToGithub(submissionDetail)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return ctx.SendStatus(200)
}
