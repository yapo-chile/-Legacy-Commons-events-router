package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/repository"
)

// DispatchEventLogger DispatchEvent's type of logs
type DispatchEventLogger interface {
	LogNewMessage(message string)
	LogErrorDecodingInput(message []byte, err error)
	LogErrorDispatching(ev domain.Event, err error)
	LogSuccess(ev domain.Event)
}

// DispatchInteractor allows push events to producer using router
type DispatchInteractor interface {
	Dispatch(event domain.Event) error
}

// DispatchEventHandler struct that represents the transfer from the reader to the message sender
type DispatchEventHandler struct {
	consumer   repository.KafkaConsumer
	interactor DispatchInteractor
	logger     DispatchEventLogger
}

// NewDispatchEventHandler initiallize a DispatchEventhandler
func NewDispatchEventHandler(
	consumer repository.KafkaConsumer,
	interactor DispatchInteractor,
	logger DispatchEventLogger,
) *DispatchEventHandler {
	return &DispatchEventHandler{
		consumer:   consumer,
		interactor: interactor,
		logger:     logger,
	}
}

type dispatchEventHandlerInput struct {
	Type    string      `json:"type"`
	DateStr string      `json:"date"`
	Content interface{} `json:"content"`
}

// Consume process a message and sent to DispatchEvent
// Gets a message on the Messages channel
// Decodes an avro message to a struct
// Parse fields to be compatible with DispatchEvent format
// Send a message to DispatchEvent
func (p *DispatchEventHandler) Consume() {
	for message := range p.consumer.GetMessages() {
		p.logger.LogNewMessage(string(message))
		input := dispatchEventHandlerInput{}
		err := json.Unmarshal(message, &input)
		if err != nil {
			p.logger.LogErrorDecodingInput(message, err)
			continue
		}
		if input.Type == "" {
			p.logger.LogErrorDecodingInput(
				message, fmt.Errorf("missing input type"),
			)
			continue
		}
		date, err := time.Parse("2006-01-02 15:04:05", input.DateStr)
		if err != nil {
			p.logger.LogErrorDecodingInput(message, err)
			continue
		}
		newEvent := domain.Event{
			Type:    input.Type,
			Date:    date,
			Content: input.Content,
		}
		if err := p.interactor.Dispatch(newEvent); err != nil {
			p.logger.LogErrorDispatching(newEvent, err)
			continue
		}
		p.logger.LogSuccess(newEvent)
	}
}
