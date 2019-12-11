package loggers

import (
	"testing"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

func TestDispatchInteractorLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeDispatchInteractorlogger(m)
	l.LogErrorGettingTopics(domain.Event{}, nil)
	l.LogErrorPushing(domain.Event{}, "", nil)
	m.AssertExpectations(t)
}
