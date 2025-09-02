package repository

import (
	"consumer-rabbitmq-to-db/src/config/db/mongo"
	"consumer-rabbitmq-to-db/src/exec/model"
	"context"
	"fmt"
	"log"
	"os"
)

// logRepository é a implementação concreta da interface LogRepository
// responsável por persistir mensagens no MongoDB.
type logRepository struct{}

// NewlogRepository cria e retorna uma nova instância de logRepository.
// Retorna a interface LogRepository para permitir abstração e testes.
func NewlogRepository() LogRepository {
	return &logRepository{}
}

// InsertBatch insere um lote de mensagens no MongoDB.
// messages: slice de model.Message a ser inserido.
// O método converte cada model.Message para model.MessageDB antes de inserir,
// garantindo que os campos fiquem no formato esperado pelo MongoDB.
// Caso a variável de ambiente MONGO_COLLECTION_LOGS não esteja definida,
// ou a coleção não possa ser obtida, retorna erro.
// Retorna erro caso a inserção falhe.
func (r *logRepository) InsertBatch(messages []model.Message) error {
	collectionName := os.Getenv("MONGO_COLLECTION_LOGS")
	if collectionName == "" {
		log.Println("[InsertBatch] Variável de ambiente MONGO_COLLECTION_LOGS não definida")
		return fmt.Errorf("MONGO_COLLECTION_LOGS não definida")
	}

	collection := mongo.GetCollection(collectionName)
	if collection == nil {
		log.Printf("[InsertBatch] Não foi possível obter a coleção %s\n", collectionName)
		return fmt.Errorf("não foi possível obter a coleção %s", collectionName)
	}

	if len(messages) == 0 {
		// Nada a inserir, retorna nil
		return nil
	}

	// Converte cada mensagem para MessageDB e cria slice de interface{}
	docs := make([]interface{}, len(messages))
	for i, m := range messages {
		docs[i] = m.ToMessageDB()
	}

	// Insere todas as mensagens de uma vez no MongoDB
	_, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		log.Println("[InsertBatch] Erro ao inserir batch no MongoDB:", err)
		return err
	}

	log.Printf("[InsertBatch] Inserido batch de %d mensagens\n", len(messages))
	return nil
}
