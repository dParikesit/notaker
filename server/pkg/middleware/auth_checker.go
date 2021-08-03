package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)
func AuthChecker(c*fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Session error",
		})
		return err
  }
	fmt.Println(sess.Keys())
	fmt.Println(sess.Fresh())
	if sess.Get("email")==nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Unauthorized",
		})
		return err
  }
	return c.Next()
}