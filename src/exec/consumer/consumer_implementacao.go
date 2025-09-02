package consumer

import (
	"consumer-rabbitmq-to-db/src/config/rabbitmq"
	"consumer-rabbitmq-to-db/src/exec/model"
	"consumer-rabbitmq-to-db/src/exec/repository"
	"encoding/json"
	"log"
	"os"
	"time"
)

// consumerExec é a implementação concreta de ConsumerExec.
// Responsável por consumir mensagens de uma fila RabbitMQ e enviar batches para o MongoDB.
type consumerExec struct{}

// NewConsumerExec cria e retorna uma nova instância de ConsumerExec.
func NewConsumerExec() ConsumerExec {
	return &consumerExec{}
}

// carregarEnv lê a variável de ambiente RABBITMQ_QUEUE_LOG.
// Retorna o nome da fila a ser consumida.
// Finaliza o programa caso a variável não esteja configurada.
func carregarEnv() string {
	queueName := os.Getenv("RABBITMQ_QUEUE_LOG")
	if queueName == "" {
		log.Fatal("[carregarEnv] RABBITMQ_QUEUE_LOG não configurado")
	}
	return queueName
}

// ConsumerLogFila inicia o consumo contínuo da fila RabbitMQ em uma goroutine.
// rm: instância de RabbitMQ
// repo: instância de LogRepository
// batchSize: tamanho máximo de cada batch
func (consumer *consumerExec) ConsumerLogFila(rm *rabbitmq.RabbitMQ, repo repository.LogRepository, batchSize int) {
	queueName := carregarEnv()

	msgs, err := rm.Consume(queueName)
	if err != nil {
		log.Println("[ConsumerLogFilaOneShot] Erro ao consumir fila:", err)
		return
	}

	log.Printf("[ConsumerLogFilaOneShot] Fila '%s' consumida com sucesso", queueName)

	var buffer []model.Message
	timer := time.NewTimer(3 * time.Second)

	sendBatch := func() {
		if len(buffer) == 0 {
			return
		}
		log.Printf("[ConsumerLogFilaOneShot] Enviando batch de %d mensagens para o Mongo", len(buffer))
		if err := repo.InsertBatch(buffer); err != nil {
			log.Println("[ConsumerLogFilaOneShot] Erro ao inserir batch:", err)
		} else {
			log.Printf("[ConsumerLogFilaOneShot] Batch de %d mensagens inserido com sucesso", len(buffer))
		}
		buffer = buffer[:0]
	}

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				// Canal fechado, envia o que restou e termina
				sendBatch()
				log.Println("[ConsumerLogFilaOneShot] Canal de mensagens fechado. Finalizando...")
				return
			}

			var m model.Message
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				log.Println("[ConsumerLogFilaOneShot] Erro ao decodificar mensagem:", err)
				continue
			}

			buffer = append(buffer, m)
			log.Printf("[ConsumerLogFilaOneShot] Mensagem adicionada ao buffer (tamanho atual: %d)", len(buffer))

			if len(buffer) >= batchSize {
				sendBatch()
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(3 * time.Second)
			}

		case <-timer.C:
			// Timeout de 3 segundos sem novas mensagens
			sendBatch()
			log.Println("[ConsumerLogFilaOneShot] Timeout atingido. Finalizando...")
			return
		}
	}
}
