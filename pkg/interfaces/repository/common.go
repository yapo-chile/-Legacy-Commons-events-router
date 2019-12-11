package repository

import "io"

// KafkaConsumer allows kafka reader operations
type KafkaConsumer interface {
	GetMessages() chan []byte
	Listen() error
	io.Closer
}

// KafkaProducer allows send messages to kafka
type KafkaProducer interface {
	SendMessage(topic string, message []byte) error
	io.Closer
}
