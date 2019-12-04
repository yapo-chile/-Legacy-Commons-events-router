package usecases

import (
	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

// DisptachInteractor allows push events to producer using router
type DisptachInteractor struct {
	Producer Producer
	Router   Router
	Logger   DisptachInteractorLogger
}

// DisptachInteractorLogger logs events in DispatchInteractor
type DisptachInteractorLogger interface {
	LogErrorGettingTopics(ev domain.Event, err error)
	LogErrorPushing(ev domain.Event, topic string, err error)
}

// Dispatch pushes event to related topic defined by router
func (i *DisptachInteractor) Dispatch(event domain.Event) error {
	topics, err := i.Router.GetTopics(event)
	if err != nil {
		i.Logger.LogErrorGettingTopics(event, err)
		return err
	}
	for _, topic := range topics {
		if err := i.Producer.Push(topic, event); err != nil {
			i.Logger.LogErrorPushing(event, topic, err)
		}
	}
	return nil
}
