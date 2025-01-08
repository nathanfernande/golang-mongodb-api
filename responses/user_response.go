package responses

import "github.com/gofiber/fiber/v2"

//Define uma struct chamada UserResponse, usada para padronizar a forma de resposta da api

type UserResponse struct {
    Status  int        `json:"status"`
    Message string     `json:"message"`
    Data    *fiber.Map `json:"data"` // fiber.Map é um atalho para um mapa genérico, semelhante a map[string]interface{}.
}

// {
//    "status": 200,
//    "message": "Operação realizada com sucesso",
//    "data": {
//      "id": 123,
//      "name": "Exemplo"
//    }
//  }