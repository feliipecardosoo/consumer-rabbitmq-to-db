package model

// MessageDB representa uma mensagem a ser inserida no MongoDB
type MessageDB struct {
	ID        string `bson:"idMembro"`  // Campo "id" no MongoDB
	Message   string `bson:"message"`   // Campo "message" no MongoDB
	Timestamp string `bson:"timestamp"` // Campo "timestamp" no MongoDB
}

// ToMessageDB converte um model.Message em model.MessageDB
func (m Message) ToMessageDB() MessageDB {
	return MessageDB{
		ID:        m.ID,
		Message:   m.Message,
		Timestamp: m.Timestamp,
	}
}
