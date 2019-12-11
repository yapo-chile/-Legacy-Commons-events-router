package repository

import (
	"encoding/json"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
	"github.mpi-internal.com/Yapo/events-router/pkg/usecases"
)

// producer allows to push events to queue
type producer struct {
	handler KafkaProducer
}

// MakeProducer creates new instance of Producer
func MakeProducer(handler KafkaProducer) usecases.Producer {
	return &producer{
		handler: handler,
	}
}

type kafkaMessage struct {
	Type    string      `json:"type"`
	Date    string      `json:"date"`
	Content interface{} `json:"content"`
}

// Push pushes given event to given topic
func (p *producer) Push(topic string, event domain.Event) error {
	message := kafkaMessage{
		Type:    event.Type,
		Date:    event.Date.Format("2006-01-02 15:04:05"), // TODO: implement unix time
		Content: event.Content,
	}
	bytes, _ := json.Marshal(message) // nolint
	return p.handler.SendMessage(topic, bytes)
}
