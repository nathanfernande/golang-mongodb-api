package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathanfernande/go-mongodb/configs"
	"github.com/nathanfernande/go-mongodb/routes"
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
