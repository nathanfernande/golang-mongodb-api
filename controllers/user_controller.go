package controllers

import (
    "context" //Usado para gerenciar o contexto e controlar operações assíncronas, como limites de tempo.
    "github.com/nathanfernande/golang-mongodb-api/configs"
    "github.com/nathanfernande/golang-mongodb-api/models"
    "github.com/nathanfernande/golang-mongodb-api/responses"
    "net/http"
    "time"

    "github.com/go-playground/validator/v10" //Biblioteca para validação de structs, usada para verificar campos obrigatórios.
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

//Declara uma variável que representa a coleção de usuários no MongoDB
//Obtida através da função GetCollection do pacote configs
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
//Inicializa o validador (validator/v10) para validar campos obrigatórios das requisições.
var validate = validator.New()

//Define uma função que cria um novo usuário.
//Parâmetro: c *fiber.Ctx representa o contexto da requisição no Fiber.
//retorna um erro 
func CreateUser(c *fiber.Ctx) error {
	//Cria um contexto com tempo limite de 10 segundos
	//defer cancel(): Garante que os recursos associados ao contexto sejam liberados.
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var user models.User
    defer cancel()

    //c.BodyParser(&user): Tenta converter o corpo da requisição em um objeto User
	//Em caso de erro, retorna uma resposta HTTP 400 com uma mensagem de erro.
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //Usa o validador para verificar se os campos obrigatórios (required) estão preenchidos
	//Se houver falhas, retorna uma resposta HTTP 400.
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }

	//Cria um novo objeto User, gerando um ObjectID único para o campo Id.
    newUser := models.User{
        Id:       primitive.NewObjectID(),
        Name:     user.Name,
        Location: user.Location,
        Title:    user.Title,
    }

	//Insere o novo usuário na coleção do MongoDB.
	//Em caso de erro, retorna uma resposta HTTP 500.
    result, err := userCollection.InsertOne(ctx, newUser)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

	//Retorna uma resposta HTTP 201 com uma mensagem de sucesso e os detalhes do resultado.
    return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})


	//Recebe o corpo da requisição.
	//Valida os dados.
	//Cria um objeto User com os dados recebidos.
	//Insere o objeto na coleção do MongoDB.
	//Retorna uma resposta JSON indicando sucesso ou falha.
}

//Define uma função que puxa os dados de um usuário.
//Parâmetro: c *fiber.Ctx representa o contexto da requisição no Fiber.
//retorna um erro
func GetAUser(c *fiber.Ctx) error {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second): Cria um contexto com timeout de 10 segundos. Se a operação ultrapassar esse tempo, o contexto será cancelado automaticamente. A função cancel é usada para liberar os recursos manualmente, se necessário.
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtém o valor do parâmetro userId da URL. O Fiber armazena parâmetros capturados em rotas (ex.: /:userId) no contexto da requisição, que pode ser acessado com c.Params.
    userId := c.Params("userId")
	//var user models.User: Declara uma variável chamada user do tipo models.User, que armazenará os dados do usuário retornados do banco de dados.
    var user models.User
	//defer cancel(): Garante que a função cancel seja chamada ao sair da função, liberando recursos associados ao contexto ctx.
    defer cancel()

	//Converte o valor de userId (uma string representando o ID do usuário) para um objeto ObjectID do MongoDB. Isso é necessário porque o MongoDB armazena IDs em um formato hexadecimal específico.
    objId, _ := primitive.ObjectIDFromHex(userId)

	//userCollection.FindOne(ctx, bson.M{"_id": objId}): Realiza uma consulta no banco de dados MongoDB para encontrar um documento na coleção userCollection onde o campo _id corresponde a objId.
	//.Decode(&user): Decodifica o documento retornado do banco e o armazena na variável user.
    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

	//Verifica se ocorreu um erro durante a consulta ao banco.
		//Caso haja erro:
			//Define o status HTTP como 500 - Internal Server Error.
			//Retorna uma resposta JSON contendo:
				//Status: O código de status HTTP.
				//Message: A mensagem "error".
				//Data: Um objeto com detalhes do erro.
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

	//Se não houver erro:
		//Define o status HTTP como 200 - OK.
		//Retorna uma resposta JSON contendo:
			//Status: O código de status HTTP.
			//Message: A mensagem "success".
			//Data: Um objeto com os dados do usuário encontrados.
    return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})



	//A função GetAUser recebe um ID do usuário como parâmetro da URL.
	//Converte o ID para o formato adequado (ObjectID).
	//Busca o usuário correspondente no banco de dados MongoDB.
	//Retorna os dados do usuário no formato JSON, ou um erro caso ocorra algum problema.
}

//Define uma função que edita um usuário já existente.
//Parâmetro: c *fiber.Ctx representa o contexto da requisição no Fiber.
//retorna um erro
func EditAUser(c *fiber.Ctx) error {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second): Cria um contexto com timeout de 10 segundos. Após esse período, o contexto será cancelado.
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//userId := c.Params("userId"): Obtém o parâmetro userId da URL.
    userId := c.Params("userId")
	// var user models.User: Declara uma variável user do tipo models.User para armazenar os dados enviados no corpo da requisição.
    var user models.User
	//defer cancel(): Garante que o contexto será liberado ao final da função.
    defer cancel()

	//Converte userId (string) para o tipo ObjectID do MongoDB.
    objId, _ := primitive.ObjectIDFromHex(userId)

    //c.BodyParser(&user): Analisa o corpo da requisição e popula a variável user com os dados recebidos.
	// Se ocorrer um erro (ex.: corpo da requisição inválido), retorna um status 400 - Bad Request com a mensagem de erro.
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //validate.Struct(&user): Usa uma biblioteca de validação (como go-playground/validator) para verificar se os campos obrigatórios do user estão preenchidos.
	//Se os dados forem inválidos, retorna um status 400 - Bad Request com detalhes da validação.
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }

	//Cria um objeto no formato BSON (Binary JSON) contendo os campos que serão atualizados no banco de dados.
    update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

	//userCollection.UpdateOne: Atualiza os dados do usuário no banco de dados.
		//O filtro bson.M{"_id": objId} seleciona o documento com o ID correspondente.
		//O operador "$set" especifica os campos que devem ser atualizados.
    result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	//Se ocorrer algum problema durante a atualização, retorna um status 500 - Internal Server Error.
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    var updatedUser models.User
	//result.MatchedCount == 1: Verifica se algum documento foi encontrado e atualizado.
    if result.MatchedCount == 1 {
		//userCollection.FindOne: Busca o documento atualizado no banco de dados.
		//.Decode(&updatedUser): Decodifica o documento encontrado e armazena os dados em updatedUser.
        err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		//Caso não consiga buscar os dados atualizados, retorna um status 500 - Internal Server Error.
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }
    }

	//Retorna uma resposta com status 200 - OK e os dados do usuário atualizados no formato JSON.
    return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})


	//Obtém o userId da URL e converte para ObjectID.
	//Valida os dados enviados no corpo da requisição.
	//Atualiza os campos no banco de dados.
	//Recupera os dados atualizados do usuário.
	//Retorna uma resposta com os dados atualizados ou um erro, se houver.
}

//Define uma função que deleta um usuário.
//Parâmetro: c *fiber.Ctx representa o contexto da requisição no Fiber.
//retorna um erro
func DeleteAUser(c *fiber.Ctx) error {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second): Cria um contexto com timeout de 10 segundos para a operação de exclusão.
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//userId := c.Params("userId"): Obtém o parâmetro userId da URL.
    userId := c.Params("userId")
	//defer cancel(): Garante que o contexto será cancelado ao sair da função, liberando recursos.
    defer cancel()

	//Converte o userId (string) para um objeto ObjectID, que é o formato utilizado pelo MongoDB para IDs.
    objId, _ := primitive.ObjectIDFromHex(userId)

	//userCollection.DeleteOne(ctx, bson.M{"_id": objId}): Executa a operação de exclusão no banco de dados. O filtro bson.M{"_id": objId} busca o documento com o ID fornecido.
    result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	// Caso ocorra um erro durante a exclusão, retorna uma resposta com status 500 - Internal Server Error e detalhes do erro.
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

	//result.DeletedCount < 1: Verifica se nenhum documento foi excluído. Isso significa que o ID especificado não foi encontrado no banco.
	//Retorna um status 404 - Not Found e uma mensagem indicando que o usuário com o ID fornecido não foi encontrado.
    if result.DeletedCount < 1 {
        return c.Status(http.StatusNotFound).JSON(
            responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
        )
    }

	//Retorna uma resposta com status 200 - OK e uma mensagem indicando que o usuário foi excluído com sucesso.
    return c.Status(http.StatusOK).JSON(
        responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
    )


	//Obtém o ID do usuário a partir da URL e o converte para ObjectID.
	//Executa a exclusão do documento correspondente no MongoDB.
	//Verifica se algum documento foi excluído:
	//Retorna erro 500 se houver problema na operação.
	//Retorna erro 404 se o ID fornecido não corresponder a nenhum usuário.
	//Caso contrário, retorna sucesso com status 200.
}

//Define uma função que puxa os dados de todos os usuários.
//Parâmetro: c *fiber.Ctx representa o contexto da requisição no Fiber.
//retorna um erro
func GetAllUsers(c *fiber.Ctx) error {
	//Cria um contexto (ctx) com um tempo limite de 10 segundos para operações assíncronas (útil para evitar que a função fique bloqueada indefinidamente).
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//Declara um slice (lista) vazio chamado users, que armazenará os dados dos usuários retornados do banco.
	//models.User é um modelo definido previamente.
    var users []models.User
	//defer cancel() garante que o contexto será cancelado ao final da execução da função, liberando recursos.
    defer cancel()

	//Executa uma busca na coleção userCollection (uma referência à coleção de usuários no MongoDB).
	//O bson.M{} representa um filtro vazio, ou seja, todos os documentos (usuários) serão retornados.
	//A busca usa o contexto ctx criado anteriormente.
	//results é um iterador que permite percorrer os documentos retornados.
    results, err := userCollection.Find(ctx, bson.M{})

	//Verifica se ocorreu algum erro na consulta ao banco de dados.
	//Caso positivo, retorna uma resposta HTTP com status 500 (erro interno do servidor) e inclui o erro na resposta no formato JSON.
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //Adiciona um defer para fechar o iterador results automaticamente após o término do uso, liberando os recursos do MongoDB.
    defer results.Close(ctx)

	//O for results.Next(ctx) percorre cada documento retornado pela consulta.
	//Em cada iteração:
		//Declara uma variável singleUser para armazenar os dados de um usuário.
		//Decodifica o documento atual em singleUser usando results.Decode(&singleUser). Caso ocorra um erro durante a decodificação, retorna uma resposta HTTP com status 500.
		//adiciona o usuário decodificado ao slice users com append(users, singleUser).

    for results.Next(ctx) {
        var singleUser models.User
        if err = results.Decode(&singleUser); err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }

        users = append(users, singleUser)
    }


	//Após percorrer todos os documentos, retorna uma resposta HTTP com status 200 (OK).
	//A resposta inclui os usuários encontrados no campo Data, encapsulados em um objeto de resposta (responses.UserResponse).
    return c.Status(http.StatusOK).JSON(
        responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
    )

	
	//Criação de contexto: Um contexto com timeout de 10 segundos é criado para limitar a execução da operação.
	//Consulta no banco: Todos os documentos (usuários) são buscados na coleção.
	//Iteração nos resultados: Os documentos retornados são iterados, decodificados em uma estrutura models.User e adicionados a um slice (users).
	//Tratamento de erros: Caso ocorra erro na consulta ou decodificação, uma resposta de erro HTTP 500 é retornada.
	//Resposta final: Após processar os dados, retorna uma resposta HTTP 200 com a lista de usuários.
}