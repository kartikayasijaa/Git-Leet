package routes

import (
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
			"error" : err.Error(),
		})
	} 
	
	code, err := services.GetCode(submissions[0].Id)
	
	if err != nil {
		h.Logger.Println(err)
		return ctx.Status(400).JSON(fiber.Map{
			"error" : err.Error(),
		})
	} 

	return ctx.Status(200).JSON(code)	
}