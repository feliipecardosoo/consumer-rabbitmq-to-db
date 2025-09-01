package main

import (
	"consumer-rabbitmq-to-db/src/config/db/mongo"
	"consumer-rabbitmq-to-db/src/config/env"
	"consumer-rabbitmq-to-db/src/config/rabbitmq"
	"log"
)

func main() {
	// Carrega variáveis de ambiente
	env.LoadEnv()

	// Conexão com RabbitMQ
	rm := rabbitmq.InitPublisher() // usa singleton
	defer func() {
		if err := rm.Close(); err != nil {
			log.Println("Erro ao fechar RabbitMQ:", err)
		}
	}()

	// Conexão MongoDB
	mongo.InitMongoConnection()
}
