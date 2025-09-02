package model

// Message representa uma mensagem recebida da fila RabbitMQ.
// É o modelo usado no consumer antes de ser convertido para MessageDB
// para inserção no MongoDB.
type Message struct {
	// ID é o identificador único da mensagem
	ID string `json:"idMembro"`

	// Message é o conteúdo textual da mensagem
	Message string `json:"message"`

	// Timestamp indica o horário em que a mensagem foi gerada ou enviada
	Timestamp string `json:"timestamp"`
}
