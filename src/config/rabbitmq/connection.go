package rabbitmq

import (
	"consumer-rabbitmq-to-db/src/config/env"
	"errors"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ encapsula a conexão e o canal com o servidor RabbitMQ.
// Inclui a conexão (Conn) e o canal (Channel) utilizados para enviar e consumir mensagens.
type RabbitMQ struct {
	Conn    *amqp.Connection // Conexão com o servidor RabbitMQ
	Channel *amqp.Channel    // Canal de comunicação com RabbitMQ
}

// QueueLog é a variável global que armazena o nome da fila de logs,
// definida pela variável de ambiente RABBITMQ_QUEUE_LOG.
var QueueLog string

// allowedQueues mantém um mapa das filas válidas que podem ser consumidas.
// Inicializado dinamicamente a partir de QueueLog.
var allowedQueues map[string]bool

// init é executado automaticamente quando o pacote é carregado.
// Faz o seguinte:
// 1. Carrega as variáveis de ambiente chamando env.LoadEnv().
// 2. Lê a variável RABBITMQ_QUEUE_LOG e atribui a QueueLog.
// 3. Encerra a aplicação com log.Fatal se RABBITMQ_QUEUE_LOG não estiver configurada.
// 4. Inicializa allowedQueues com a fila definida em QueueLog.
func init() {
	// Carrega variáveis de ambiente do arquivo .env
	env.LoadEnv()

	// Lê variável de ambiente com o nome da fila
	QueueLog = os.Getenv("RABBITMQ_QUEUE_LOG")
	if QueueLog == "" {
		log.Fatal("RABBITMQ_QUEUE_LOG não configurada")
	}

	// Inicializa allowedQueues depois de ler QueueLog
	allowedQueues = map[string]bool{
		QueueLog: true,
	}
}

// InitRabbitMQ cria uma conexão com o RabbitMQ usando a variável de ambiente RABBITMQ_URI.
// Retorna uma instância de RabbitMQ pronta para consumo e publicação de mensagens.
func InitRabbitMQ() *RabbitMQ {
	uri := os.Getenv("RABBITMQ_URI")
	if uri == "" {
		log.Fatal("RABBITMQ_URI não configurado")
	}

	// Estabelece a conexão
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatal("Erro ao conectar RabbitMQ:", err)
	}

	// Abre um canal de comunicação
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatal("Erro ao abrir canal RabbitMQ:", err)
	}

	log.Println("Conexão RabbitMQ estabelecida com sucesso")
	return &RabbitMQ{Conn: conn, Channel: ch}
}

// Close fecha o canal e a conexão com o RabbitMQ.
// Retorna erro caso falhe ao fechar algum dos recursos.
func (r *RabbitMQ) Close() error {
	if err := r.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := r.Conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

// validateQueue verifica se o nome da fila está entre as filas permitidas.
// Retorna erro caso a fila não seja válida.
func validateQueue(queueName string) error {
	if !allowedQueues[queueName] {
		return errors.New("fila não permitida: " + queueName)
	}
	return nil
}

// Consume inicia o consumo de mensagens de uma fila específica.
// Valida a fila, declara a fila no servidor (caso não exista) e retorna um canal de mensagens.
// queueName: nome da fila que será consumida.
// Retorna: canal de mensagens (<-chan amqp.Delivery) e possível erro.
func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	if err := validateQueue(queueName); err != nil {
		return nil, err
	}

	// Declara a fila no servidor
	_, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Inicia o consumo da fila
	msgs, err := r.Channel.Consume(
		queueName,
		"",    // consumer name (empty = auto-generate)
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	return msgs, nil
}

// RMInstance mantém uma instância singleton do RabbitMQ para publicação
var RMInstance *RabbitMQ

// InitPublisher inicializa a instância singleton do RabbitMQ para envio de mensagens.
// Caso já exista, retorna a instância existente.
func InitPublisher() *RabbitMQ {
	if RMInstance == nil {
		RMInstance = InitRabbitMQ()
	}
	return RMInstance
}
