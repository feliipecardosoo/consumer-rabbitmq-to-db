package repository

import "consumer-rabbitmq-to-db/src/exec/model"

// LogRepository define a interface para inserção de mensagens no repositório.
// A ideia é permitir a inserção de mensagens em lote (batch) no banco de dados.
// Implementações dessa interface podem ser feitas, por exemplo, usando MongoDB ou outro banco.
type LogRepository interface {
	// InsertBatch insere um slice de mensagens no repositório.
	// messages: slice de model.Message que será inserido.
	// Retorna erro caso ocorra falha na inserção.
	InsertBatch([]model.Message) error
}
