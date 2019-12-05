package usecases

import (
	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

// GomsRepository interface that represents all the methods available to
// interact with events-router microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}

// Producer represents a handler that interacts with the producer infrastructure
type Producer interface {
	Push(topic string, event domain.Event) error
}

// Router allows get topics from remote configuration for each incoming event
type Router interface {
	GetTopics(event domain.Event) ([]string, error)
}
