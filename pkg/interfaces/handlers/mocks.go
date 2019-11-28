package handlers

import (
	"github.com/stretchr/testify/mock"
)

// MockInputRequest is a mock class
type MockInputRequest struct {
	mock.Mock
}

// Set is a mocked method
func (m *MockInputRequest) Set(input interface{}) InputRequest {
	args := m.Called(input)
	return args.Get(0).(InputRequest)
}

// FromJSONBody is a mocked method
func (m *MockInputRequest) FromJSONBody() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

// FromRawBody is a mocked method
func (m *MockInputRequest) FromRawBody() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

// FromPath is a mocked method
func (m *MockInputRequest) FromPath() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

// FromQuery is a mocked method
func (m *MockInputRequest) FromQuery() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

// FromHeaders is a mocked method
func (m *MockInputRequest) FromHeaders() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

// FromCookies is a mocked method
func (m *MockInputRequest) FromCookies() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

// FromForm is a mocked method
func (m *MockInputRequest) FromForm() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}
