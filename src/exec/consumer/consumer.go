package consumer

import (
	"consumer-rabbitmq-to-db/src/config/rabbitmq"
	"consumer-rabbitmq-to-db/src/exec/repository"
)

// ConsumerExec define os métodos que um consumer deve implementar.
// Atualmente, possui apenas um método para consumir mensagens de uma fila RabbitMQ
// e enviar os dados em batch para o MongoDB.
type ConsumerExec interface {
	// ConsumerLogFila consome mensagens da fila especificada no RabbitMQ
	// e envia em lotes (batch) para o repositório MongoDB.
	//
	// Parâmetros:
	// - rm: instância de RabbitMQ usada para consumir mensagens
	// - repo: repositório que fará a persistência das mensagens
	// - batchSize: número máximo de mensagens a serem enviadas em cada batch
	ConsumerLogFila(rm *rabbitmq.RabbitMQ, repo repository.LogRepository, batchSize int)
}
