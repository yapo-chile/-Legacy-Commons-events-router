package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
)

type mockRConfig struct {
	mock.Mock
}

func (m *mockRConfig) Get(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestGetTopicsOK(t *testing.T) {
	mRconfig := &mockRConfig{}
	mRconfig.On("Get",
		mock.AnythingOfType("string")).
		Return(`["test"]`)
	repo := MakeRouter(mRconfig)
	res, err := repo.GetTopics(domain.Event{})
	assert.NoError(t, err)
	expected := []string{"test"}
	assert.Equal(t, expected, res)
	mRconfig.AssertExpectations(t)
}

func TestGetTopicsErrorUnmarshal(t *testing.T) {
	mRconfig := &mockRConfig{}
	mRconfig.On("Get",
		mock.AnythingOfType("string")).
		Return(`bad response`)
	repo := MakeRouter(mRconfig)
	res, err := repo.GetTopics(domain.Event{})
	assert.Error(t, err)
	expected := []string{}
	assert.Equal(t, expected, res)
	mRconfig.AssertExpectations(t)
}

func TestGetTopicsConfigNotFound(t *testing.T) {
	mRconfig := &mockRConfig{}
	mRconfig.On("Get",
		mock.AnythingOfType("string")).Return(``)
	repo := MakeRouter(mRconfig)
	res, err := repo.GetTopics(domain.Event{})
	assert.Error(t, err)
	expected := []string{}
	assert.Equal(t, expected, res)
	mRconfig.AssertExpectations(t)
}
