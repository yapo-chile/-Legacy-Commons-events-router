package loggers

import (
	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/handlers"
)

type dispatchEventsHandlerLogger struct {
	logger Logger
}

// LogNewMessage logs new incoming message
func (l *dispatchEventsHandlerLogger) LogNewMessage(m string) {
	l.logger.Info("< new kafka message: %s", m)
}

// LogErrorDispatching logs error on dispatching event
func (l *dispatchEventsHandlerLogger) LogErrorDispatching(ev domain.Event, err error) {
	l.logger.Error("< error dispatching event %s: %v", ev.Type, err)
}

// LogErrorDecodingInput logs error on decoding input from queue
func (l *dispatchEventsHandlerLogger) LogErrorDecodingInput(message []byte, err error) {
	l.logger.Error("< error decoding input from kafka message %s: %v", string(message), err)
}

// LogSuccess logs info on success dispatching event
func (l *dispatchEventsHandlerLogger) LogSuccess(ev domain.Event) {
	l.logger.Info("> %s event sent", ev.Type)
}

// MakeDispatchEventHandlerlogger sets up a DispatchEventLogger instrumented
// via the provided logger
func MakeDispatchEventHandlerlogger(logger Logger) handlers.DispatchEventLogger {
	return &dispatchEventsHandlerLogger{
		logger: logger,
	}
}
