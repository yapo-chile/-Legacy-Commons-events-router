package loggers

import (
	"testing"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

func TestDispatchEventsHandlerLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeDispatchEventHandlerlogger(m)
	l.LogNewMessage("")
	l.LogErrorDispatching(domain.Event{}, nil)
	l.LogErrorDecodingInput([]byte{}, nil)
	l.LogSuccess(domain.Event{})
	m.AssertExpectations(t)
}
