package usecases

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

type mockProducer struct {
	mock.Mock
}

func (m *mockProducer) Push(topic string, event domain.Event) error {
	args := m.Called(topic, event)
	return args.Error(0)
}

type mockRouter struct {
	mock.Mock
}

func (m *mockRouter) GetTopics(event domain.Event) ([]string, error) {
	args := m.Called(event)
	return args.Get(0).([]string), args.Error(1)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) LogErrorGettingTopics(ev domain.Event, err error) {
	m.Called(ev, err)
}

func (m *mockLogger) LogErrorPushing(ev domain.Event, topic string, err error) {
	m.Called(ev, topic, err)
}

func TestDispatchOK(t *testing.T) {
	mProducer := &mockProducer{}
	mRouter := &mockRouter{}
	mLogger := &mockLogger{}

	mProducer.On("Push", mock.AnythingOfType("string"),
		mock.AnythingOfType("domain.Event")).Return(nil)
	mRouter.On("GetTopics",
		mock.AnythingOfType("domain.Event")).Return([]string{""}, nil)
	interactor := DispatchInteractor{
		Producer: mProducer,
		Router:   mRouter,
		Logger:   mLogger,
	}
	err := interactor.Dispatch(domain.Event{})
	assert.NoError(t, err)

	mProducer.AssertExpectations(t)
	mRouter.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestDispatchPushError(t *testing.T) {
	mProducer := &mockProducer{}
	mRouter := &mockRouter{}
	mLogger := &mockLogger{}

	mProducer.On("Push", mock.AnythingOfType("string"),
		mock.AnythingOfType("domain.Event")).Return(fmt.Errorf("e"))
	mRouter.On("GetTopics",
		mock.AnythingOfType("domain.Event")).Return([]string{""}, nil)
	mLogger.On("LogErrorPushing", mock.AnythingOfType("domain.Event"),
		mock.AnythingOfType("string"), mock.Anything)
	interactor := DispatchInteractor{
		Producer: mProducer,
		Router:   mRouter,
		Logger:   mLogger,
	}
	err := interactor.Dispatch(domain.Event{})
	assert.NoError(t, err)

	mProducer.AssertExpectations(t)
	mRouter.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestDispatchGetTopicsError(t *testing.T) {
	mProducer := &mockProducer{}
	mRouter := &mockRouter{}
	mLogger := &mockLogger{}

	mRouter.On("GetTopics",
		mock.AnythingOfType("domain.Event")).Return(
		[]string{""}, fmt.Errorf("e"))
	mLogger.On("LogErrorGettingTopics", mock.AnythingOfType("domain.Event"),
		mock.Anything)
	interactor := DispatchInteractor{
		Producer: mProducer,
		Router:   mRouter,
		Logger:   mLogger,
	}
	err := interactor.Dispatch(domain.Event{})
	assert.Error(t, err)

	mProducer.AssertExpectations(t)
	mRouter.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
