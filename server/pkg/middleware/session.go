package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func Init() error {
	Store = session.New(session.Config{
		CookieSecure: true,
		CookieHTTPOnly: true,
		CookieSameSite: "strict",
	})
	return nil
}