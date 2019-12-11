package handlers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) LogNewMessage(msg string) {
	m.Called(msg)
}

func (m *mockLogger) LogErrorDispatching(ev domain.Event, err error) {
	m.Called(ev, err)
}

func (m *mockLogger) LogSuccess(ev domain.Event) {
	m.Called(ev)
}

func (m *mockLogger) LogErrorDecodingInput(msg []byte, err error) {
	m.Called(msg, err)
}

type mockConsumer struct {
	mock.Mock
}

func (m *mockConsumer) GetMessages() chan []byte {
	args := m.Called()
	return args.Get(0).(chan []byte)
}

func (m *mockConsumer) Listen() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockConsumer) Close() error {
	args := m.Called()
	return args.Error(0)
}

type mockDispatchInteractor struct {
	mock.Mock
}

func (m *mockDispatchInteractor) Dispatch(event domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestConsumeOK(t *testing.T) {
	mLogger := &mockLogger{}
	mConsumer := &mockConsumer{}
	mInteractor := &mockDispatchInteractor{}
	handler := NewDispatchEventHandler(mConsumer,
		mInteractor, mLogger)
	mLogger.On("LogNewMessage", mock.Anything)
	mInteractor.On("Dispatch",
		mock.AnythingOfType("domain.Event")).Return(nil)
	mLogger.On("LogSuccess", mock.Anything)
	messageCh := make(chan []byte, 1)
	messageCh <- []byte(`{"type":"bump","date":"2019-12-04 16:48:38",
	"content":{"action_id":"6","ad_id":"6283212",
	"list_time":"2019-12-04 16:48:38.638176"}}`)
	mConsumer.On("GetMessages").Return(messageCh)
	close(messageCh)
	handler.Consume()
	mInteractor.AssertExpectations(t)
	mConsumer.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestConsumeUnmarshalError(t *testing.T) {
	mLogger := &mockLogger{}
	mConsumer := &mockConsumer{}
	mInteractor := &mockDispatchInteractor{}
	handler := NewDispatchEventHandler(mConsumer,
		mInteractor, mLogger)
	mLogger.On("LogNewMessage", mock.Anything)
	mLogger.On("LogErrorDecodingInput",
		mock.AnythingOfType("[]uint8"),
		mock.Anything)
	messageCh := make(chan []byte, 1)
	messageCh <- []byte(`XXXXXXXXX`)
	mConsumer.On("GetMessages").Return(messageCh)
	close(messageCh)
	handler.Consume()
	mInteractor.AssertExpectations(t)
	mConsumer.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestConsumeMissingTypeError(t *testing.T) {
	mLogger := &mockLogger{}
	mConsumer := &mockConsumer{}
	mInteractor := &mockDispatchInteractor{}
	handler := NewDispatchEventHandler(mConsumer,
		mInteractor, mLogger)
	mLogger.On("LogNewMessage", mock.Anything)
	mLogger.On("LogErrorDecodingInput",
		mock.AnythingOfType("[]uint8"),
		mock.Anything)
	messageCh := make(chan []byte, 1)
	messageCh <- []byte(`{"type":"","date":"2019-12-04 16:48:38",
	"content":{"action_id":"6","ad_id":"6283212",
	"list_time":"2019-12-04 16:48:38.638176"}}`)
	mConsumer.On("GetMessages").Return(messageCh)
	close(messageCh)
	handler.Consume()
	mInteractor.AssertExpectations(t)
	mConsumer.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestConsumeDateError(t *testing.T) {
	mLogger := &mockLogger{}
	mConsumer := &mockConsumer{}
	mInteractor := &mockDispatchInteractor{}
	handler := NewDispatchEventHandler(mConsumer,
		mInteractor, mLogger)
	mLogger.On("LogNewMessage", mock.Anything)
	mLogger.On("LogErrorDecodingInput",
		mock.AnythingOfType("[]uint8"),
		mock.Anything)
	messageCh := make(chan []byte, 1)
	messageCh <- []byte(`{"type":"bump","date":"XXXXX",
	"content":{"action_id":"6","ad_id":"6283212",
	"list_time":"2019-12-04 16:48:38.638176"}}`)
	mConsumer.On("GetMessages").Return(messageCh)
	close(messageCh)
	handler.Consume()
	mInteractor.AssertExpectations(t)
	mConsumer.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestConsumeDispatchError(t *testing.T) {
	mLogger := &mockLogger{}
	mConsumer := &mockConsumer{}
	mInteractor := &mockDispatchInteractor{}
	handler := NewDispatchEventHandler(mConsumer,
		mInteractor, mLogger)
	mLogger.On("LogNewMessage", mock.Anything)
	mInteractor.On("Dispatch",
		mock.AnythingOfType("domain.Event")).
		Return(fmt.Errorf("err"))
	mLogger.On("LogErrorDispatching",
		mock.AnythingOfType("domain.Event"),
		mock.Anything)
	messageCh := make(chan []byte, 1)
	messageCh <- []byte(`{"type":"bump","date":"2019-12-04 16:48:38",
	"content":{"action_id":"6","ad_id":"6283212",
	"list_time":"2019-12-04 16:48:38.638176"}}`)
	mConsumer.On("GetMessages").Return(messageCh)
	close(messageCh)
	handler.Consume()
	mInteractor.AssertExpectations(t)
	mConsumer.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
