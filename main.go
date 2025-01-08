package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathanfernande/golang-mongodb-api/configs"
	"github.com/nathanfernande/golang-mongodb-api/routes"
)

func main() {
	//iniciando a aplicação
	app := fiber.New()

	// rodar o banco de dados
	configs.ConnectDB()

	//rotas
	routes.UserRoute(app)

	//inicia o servidos HTTP na porta 6000
	app.Listen(":6000")
}
