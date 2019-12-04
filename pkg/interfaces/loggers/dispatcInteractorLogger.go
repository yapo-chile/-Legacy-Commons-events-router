package loggers

import (
	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
	"github.mpi-internal.com/Yapo/events-router/pkg/usecases"
)

type dispatchInteractorLogger struct {
	logger Logger
}

func (l *dispatchInteractorLogger) LogErrorGettingTopics(ev domain.Event, err error) {
	l.logger.Error("< error getting topic for event %s: %v", ev.Type, err)
}

func (l *dispatchInteractorLogger) LogErrorPushing(ev domain.Event, topic string, err error) {
	l.logger.Error("< error pushing to %s topic, event %s: %v", topic, ev.Type, err)
}

// MakeDispatchInteractorlogger sets up a DispatchInteractorLogger instrumented
// via the provided logger
func MakeDispatchInteractorlogger(logger Logger) usecases.DisptachInteractorLogger {
	return &dispatchInteractorLogger{
		logger: logger,
	}
}
