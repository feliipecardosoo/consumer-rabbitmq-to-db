package env

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv carrega as variáveis de ambiente a partir do arquivo `.env` na raiz do projeto.
// Caso ocorra algum erro ao carregar o arquivo, a função encerra a aplicação com log.Fatal.
// Uso típico: chamar LoadEnv() no início da execução da aplicação.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
}
