package utils

import (
	// "os"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetCookie - returns a cookie
func GetCookie(data string, name string) *fiber.Cookie {
	isProd := false
	if os.Getenv("IS_PROD") == "YES" {
		isProd = true
	}
	cookie := new(fiber.Cookie)
	cookie.HTTPOnly = true
	cookie.Name = name
	cookie.Value = data
	cookie.Secure = isProd
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30)
	cookie.SameSite = "None"
	return cookie
}
