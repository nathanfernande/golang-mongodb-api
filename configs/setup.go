package configs

import (
    "context" // Permite controlar o tempo limite (timeout) e o cancelamento em chamadas de funções.
    "fmt"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options" //Parte do driver oficial do MongoDB para Go, usado para conexão e configuração.
)

//retorna um ponteiro para uma instância de mongo.Client. Essa função é responsável por configurar e estabelecer a conexão com o MongoDB.
func ConnectDB() *mongo.Client  {
	//Cria um novo cliente MongoDB usando mongo.NewClient e define a URI de conexão com ApplyURI, que é fornecida pela função EnvMongoURI.
    client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))

	//Se ocorrer um erro durante a criação do cliente, o programa exibe o erro e encerra a execução com log.Fatal.
    if err != nil {
        log.Fatal(err)
    }

	//Cria um contexto com um tempo limite de 10 segundos. Esse contexto é usado para controlar operações de conexão.
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Conecta ao MongoDB usando o cliente e o contexto configurado.
    err = client.Connect(ctx)

	//caso ocorra um erro na conexão, ele é registrado e o programa é encerrado com log.Fatal.
    if err != nil {
        log.Fatal(err)
    }

    //Envia um comando ping ao MongoDB para verificar se a conexão está ativa.
    err = client.Ping(ctx, nil)
	//Se o ping falhar, o programa encerra com o erro.
    if err != nil {
        log.Fatal(err)
    }
	//Caso contrário, exibe a mensagem "Connected to MongoDB" no console.
    fmt.Println("Connected to MongoDB")

	//Retorna a instância do cliente conectado ao MongoDB.
    return client
}

//Cria uma variável global chamada DB que armazena o cliente retornado pela função ConnectDB
var DB *mongo.Client = ConnectDB()

//retorna uma coleção (mongo.Collection) de um banco de dados.
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	//Especifica o banco de dados chamado golangAPI
	//Obtém uma coleção específica pelo nome fornecido no parâmetro collectionName.
    collection := client.Database("golangAPI").Collection(collectionName)
    return collection
}

//Configura uma conexão com o MongoDB utilizando uma URI obtida via EnvMongoURI.
//Implementa um timeout para operações de conexão.
//Faz um teste de conectividade com o MongoDB usando Ping.
//Expõe a função GetCollection para obter coleções específicas de um banco de dados.
//É organizado para ser reutilizável em diferentes partes do projeto.