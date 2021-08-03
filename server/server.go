package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	

	"server/pkg/routes"
	"server/pkg/middleware"
	"server/platform/database"
)

type Data struct {
	Nama  string
	Nomor string
}


func main() {

	if err := database.Connect(); err != nil {
		log.Panic(err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use(logger.New())
	app.Use(recover.New())
	middleware.Init()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Static("/note/images", "./images")
	app.Post("/user/login", routes.Login)
	app.Post("/user/logout", middleware.AuthChecker,routes.Logout)
	app.Post("/note", middleware.AuthChecker,routes.Add)
	app.Get("/note", routes.GetNotes)
	app.Get("/note/:id", routes.GetNote)
	app.Patch("/note/:id", routes.PatchNote)
	app.Delete("/note/:id", routes.DeleteNote)

	app.Listen(":3001")
}
