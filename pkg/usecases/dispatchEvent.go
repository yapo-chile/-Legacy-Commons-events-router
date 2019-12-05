package usecases

import (
	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

// DispatchInteractor allows push events to producer using router
type DispatchInteractor struct {
	Producer Producer
	Router   Router
	Logger   DispatchInteractorLogger
}

// DispatchInteractorLogger logs events in DispatchInteractor
type DispatchInteractorLogger interface {
	LogErrorGettingTopics(ev domain.Event, err error)
	LogErrorPushing(ev domain.Event, topic string, err error)
}

// Dispatch pushes event to related topic defined by router
func (i *DispatchInteractor) Dispatch(event domain.Event) error {
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
