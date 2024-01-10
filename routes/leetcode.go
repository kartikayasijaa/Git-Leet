package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LeetcodeRoute(app fiber.Router, h *Handlers) {
	leetcodeRoute := app.Group("/leetcode")
	leetcodeRoute.Patch("/update", h.UpdateLeetcode)
}

func (h *Handlers) UpdateLeetcode(ctx *fiber.Ctx) error {
	type UpateLeetcodeResponse struct {
		LeetcodeUsername string `json:"leetcode_username"`
		UserID           string `json:"userId"`
	}

	res := new(UpateLeetcodeResponse)

	if err := ctx.BodyParser(&res); err != nil {
		h.Logger.Println(err.Error())
		return ctx.Status(400).JSON(&fiber.Map{
			"error": "Invalid request body",
		})
	}

	totalSubmissions, err := h.Services.LeetcodeService.GetTotalSubmission(res.LeetcodeUsername)
	if err != nil {
		h.Logger.Println(err.Error())
		return ctx.Status(400).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if res.LeetcodeUsername == "" || res.UserID == "" {
		fmt.Println("errpr")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Leetcode username and userID are required fields",
		})
	}

	err = h.Services.DBService.UpdateLeetcode(res.LeetcodeUsername, res.UserID, totalSubmissions)
	if err != nil {
		h.Logger.Println(err.Error())
		return ctx.Status(400).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message":           "Updated Successfully",
		"leetcode_username": res.LeetcodeUsername,
	})

}
