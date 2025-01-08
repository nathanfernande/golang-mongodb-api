package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nathanfernande/go-mongodb/controllers"
)

func UserRoute(app *fiber.App) {
    //todas as rotas relacionadas aos usuarios estarão aqui
	app.Post("/user", controllers.CreateUser)
    app.Get("/user/:userId", controllers.GetAUser)
    app.Put("/user/:userId", controllers.EditAUser)
    app.Delete("/user/:userId", controllers.DeleteAUser)
    app.Get("/users", controllers.GetAllUsers)
}