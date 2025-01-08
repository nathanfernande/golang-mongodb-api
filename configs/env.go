package configs

import (
    "log"
    "os"
    "github.com/joho/godotenv" // Pacote ques ajuda a carregar variáveis de ambiente a partir de um arquivo .env
)

//Cria uma função EnvMongoURI que verifica se a variável de ambiente está carregada corretamente e retorna a variável de ambiente.
func EnvMongoURI() string {
	//Esse método lê o arquivo .env e disponibiliza suas variáveis como variáveis de ambiente.
    err := godotenv.Load()

	//Verifica se ocorreu um erro ao tentar carregar o arquivo .env:
	//se err não for nil (houve um erro), o programa exibe a mensagem "Error loading .env file" no console e encerra a execução com log.Fatal
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	//Retorna o valor da variável de ambiente MONGOURI, que foi definida no arquivo .env ou no sistema operacional. O método os.Getenv busca o valor da variável de ambiente pelo nome fornecido ("MONGOURI").
    return os.Getenv("MONGOURI")
}
//Carregar variáveis de ambiente a partir de um arquivo .env (geralmente usado para armazenar configurações sensíveis como strings de conexão, chaves de API, etc.).
//Garantir que a variável MONGOURI, usada para conexão com o MongoDB, esteja acessível no programa.
//Se o arquivo .env não puder ser carregado, o programa encerra com uma mensagem de erro.