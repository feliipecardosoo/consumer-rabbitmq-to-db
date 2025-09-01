package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient é a instância global do cliente MongoDB que pode ser utilizada em toda a aplicação.
var MongoClient *mongo.Client

// InitMongoConnection inicializa a conexão com o MongoDB.
// Busca a URI do MongoDB na variável de ambiente "MONGO_URI".
// Caso não esteja configurada, encerra a aplicação com log.Fatal.
// Cria um contexto com timeout de 10 segundos para a conexão e verifica se o MongoDB está acessível (Ping).
// Em caso de sucesso, MongoClient estará pronto para uso.
func InitMongoConnection() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI não configurado")
	}

	// Configura opções de cliente
	clientOptions := options.Client().ApplyURI(uri)

	// Contexto com timeout para conexão
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	// Testa a conexão
	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Não foi possível pingar o MongoDB: %v", err)
	}

	fmt.Println("Conectado ao MongoDB com sucesso!")
}
