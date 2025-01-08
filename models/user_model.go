package models

import "go.mongodb.org/mongo-driver/bson/primitive" //Importa o pacote primitive da biblioteca oficial do driver MongoDB para Go. Este pacote fornece tipos básicos usados pelo MongoDB, como ObjectID

//Define uma struct chamada User que representa um modelo de usuário. Esta estrutura é usada para mapear dados entre a aplicação e o banco de dados MongoDB.

type User struct {
    Id       primitive.ObjectID `json:"id,omitempty"` //omitempty: Omite o campo no JSON se ele estiver vazio
    Name     string             `json:"name,omitempty" validate:"required"` //validate: "required" é uma validação que garante que o campo não esteja vazio
    Location string             `json:"location,omitempty" validate:"required"`
    Title    string             `json:"title,omitempty" validate:"required"`
}
