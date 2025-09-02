package mongo

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollection retorna uma coleção do MongoDB com base no nome fornecido.
// Usa a variável de ambiente MONGO_DB_NAME para determinar o banco de dados.
//
// Parâmetros:
// - name: nome da coleção a ser retornada
//
// Retorna:
// - ponteiro para mongo.Collection correspondente à coleção solicitada
func GetCollection(name string) *mongo.Collection {
	dbName := os.Getenv("MONGO_DB_NAME")
	return MongoClient.Database(dbName).Collection(name)
}
