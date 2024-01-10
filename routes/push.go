package routes

import (
	"gitleet/services"
	"github.com/gofiber/fiber/v2"
)

func PushRoute(app fiber.Router, h *Handlers) {
	pushRoute := app.Group("/push")
	pushRoute.Get("/:userId", h.PushLatest)
}

func (h *Handlers) PushLatest(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	submissions, err := h.Services.LeetcodeService.GetRecentSubmission(userId)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for _, s := range submissions {
		submissionDetail, err := h.Services.LeetcodeService.GetCode(s.Id)
		if err != nil {
			h.Logger.Println(err)
			return ctx.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = services.PushToGithub(submissionDetail)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return ctx.SendStatus(200)
}
