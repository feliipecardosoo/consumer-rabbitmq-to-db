package main

import (
	"consumer-rabbitmq-to-db/src/config/db/mongo"
	"consumer-rabbitmq-to-db/src/config/env"
	"consumer-rabbitmq-to-db/src/config/rabbitmq"
	"consumer-rabbitmq-to-db/src/exec/consumer"
	"consumer-rabbitmq-to-db/src/exec/repository"
	"log"
)

func main() {
	// Carrega variáveis de ambiente
	env.LoadEnv()

	// Conexão com RabbitMQ
	rm := rabbitmq.InitPublisher()
	defer func() {
		if err := rm.Close(); err != nil {
			log.Println("Erro ao fechar RabbitMQ:", err)
		}
	}()

	// Conexão MongoDB
	mongo.InitMongoConnection()

	// Cria instância do repository
	repo := repository.NewlogRepository()

	// Cria instância do consumer
	cons := consumer.NewConsumerExec()

	// Define batch
	batchSize := 10

	cons.ConsumerLogFila(rm, repo, batchSize)
	log.Println("Consumer finalizado.")
}
