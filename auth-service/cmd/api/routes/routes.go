package routes

import (
	"backend2/cmd/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", handlers.Register)
	app.Post("/api/login", handlers.Login)
	app.Post("/api/logout", handlers.Logout)

	app.Get("/api/user", handlers.GetUser)
}
