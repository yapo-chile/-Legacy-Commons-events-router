package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

type mockKafkaProducer struct {
	mock.Mock
}

func (m *mockKafkaProducer) SendMessage(topic string, bytes []byte) error {
	args := m.Called(topic, bytes)
	return args.Error(0)
}

func (m *mockKafkaProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestPushOK(t *testing.T) {
	mProducer := &mockKafkaProducer{}
	mProducer.On("SendMessage", mock.AnythingOfType("string"),
		mock.AnythingOfType("[]uint8")).Return(nil)
	repo := MakeProducer(mProducer)
	err := repo.Push("some-topic", domain.Event{})
	assert.NoError(t, err)
	mProducer.AssertExpectations(t)
}
