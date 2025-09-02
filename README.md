# Consumer RabbitMQ to MongoDB

Este projeto é um **consumer** que lê mensagens de uma fila do RabbitMQ e insere em lote no MongoDB. Ele foi desenvolvido em Go e suporta envio em **batches** com controle de tempo.

> **Observação sobre execução contínua vs agendada:**  
> Existe uma branch separada `goroutine-consumer
` onde o consumer utiliza **goroutines** para processar continuamente as mensagens, mantendo o sistema sempre ativo.  
> Para este exemplo específico, o objetivo é rodar o consumer **uma vez por dia** usando um cron (ou GitHub Actions), portanto o programa é executado, processa as mensagens disponíveis e termina.


---

## Funcionalidades

- Conexão com RabbitMQ para consumir mensagens de uma fila específica.
- Armazenamento de mensagens no MongoDB em lotes (`batch`) para otimização.
- Configuração de variáveis de ambiente para filas, banco de dados e coleções.
- Timer para enviar mensagens incompletas caso o batch não atinja o tamanho configurado em tempo definido (ex: 3 segundos).
- Logs detalhados para acompanhamento do consumo e inserção.

---

## Estrutura do Projeto

```
consumer-rabbitmq-to-db/
├── src/
│ ├── config/
│ │ ├── db/
│ │ │ └── mongo/ # Inicialização e função para pegar collections
│ │ ├── env/ # Carregamento de variáveis de ambiente
│ │ └── rabbitmq/ # Conexão e consumo do RabbitMQ
│ ├── exec/
│ │ ├── consumer/ # Lógica de consumo e buffer
│ │ ├── model/ # Modelos de mensagens
│ │ └── repository/ # Repository para inserir batches no MongoDB
├── main.go # Entrada principal do programa
├── go.mod
└── go.sum
```


---

## Requisitos

- Go >= 1.21
- RabbitMQ
- MongoDB
- Variáveis de ambiente configuradas

---

## Variáveis de Ambiente

- `RABBITMQ_URI` → URI de conexão com RabbitMQ  
- `RABBITMQ_QUEUE_LOG` → Nome da fila a ser consumida  
- `MONGO_URI` → URI de conexão com MongoDB  
- `MONGO_DB_NAME` → Nome do banco de dados MongoDB  
- `MONGO_COLLECTION_LOGS` → Nome da coleção onde as mensagens serão inseridas  

Exemplo `.env`:

```env
RABBITMQ_URI=amqp://user:pass@localhost:5672/
RABBITMQ_QUEUE_LOG=membros_log
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=mydatabase
MONGO_COLLECTION_LOGS=logs
```

## Estrutura da Mensagem
- O consumer espera mensagens no seguinte formato JSON:

```
{
  "id": "123",
  "message": "Mensagem de teste",
  "timestamp": "2025-09-01T21:35:49-03:00"
}
```

## Configuração de Batch e Timer
- `batchSize` → Número de mensagens para inserir por vez (ex: 10)
- Timer → Se o `batch` não atingir batchSize em 3 segundos, ele envia as mensagens atuais.
