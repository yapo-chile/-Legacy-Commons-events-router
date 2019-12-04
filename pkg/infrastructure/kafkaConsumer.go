package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka" // nolint
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/repository"
)

// KafkaConsumer represents a consumer to kafka queue
type KafkaConsumer struct {
	config    kafka.ConfigMap
	consumer  *kafka.Consumer
	messages  chan []byte
	topics    []string
	logger    loggers.Logger
	connected bool
}

// NewKafkaConsumer Initialize a new KafkaConsumer with fields:
// Host: kafka host
// Port: kafka port
// User: user to connect to kafka queue
// Password: password to connect to kafka queue
// GroupID: unique id to identify a consumer
// SecurityMechanism: mechanism used by kafka (ex. SCRAM-SHA-256)
// SecurityProtocol: protocol to connect with security parameters (ex. sasl_ssl)
// OffsetReset: in which offset the consumer start reading (ex. earliest)
// RebalanceEnable: flag to rebalance topic and partitions to a consumer
// ChannelEnable: flag to read messages and events from the Events channel (true) or Poll(false)
// PartitionEOF: flag to emit when A partition EOF is reached
// TimeOut: range of consumer session timeout
// Topics: array of topics to be readed by consumer
// Messages: channel to store message from kafka queue
func NewKafkaConsumer(
	host string,
	port int,
	groupID, offsetReset string,
	rebalanceEnable, channelEnable, partitionEOF bool,
	timeOut int,
	topics []string,
	logger loggers.Logger) (repository.KafkaConsumer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers":               fmt.Sprintf("%v:%d", host, port),
		"group.id":                        groupID,
		"go.application.rebalance.enable": rebalanceEnable,
		"go.events.channel.enable":        channelEnable,
		"enable.partition.eof":            partitionEOF,
		"auto.offset.reset":               offsetReset,
		"session.timeout.ms":              timeOut,
	}
	consumer := KafkaConsumer{
		topics:   topics,
		config:   conf,
		messages: make(chan []byte),
		logger:   logger,
	}
	err := consumer.connect()
	if err != nil {
		return nil, err
	}
	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		logger.Info("Consumer %+v connected to kafka using config: \n%s\n", consumer, jconf)
	} else {
		logger.Info("Consumer %+v connected to kafka using config: \n%+v\n", consumer, conf)
	}
	return &consumer, nil
}

// Connect initialize a new reader with the given config,
// Subscribe to topics and start to read messages from kafka queue
func (k *KafkaConsumer) connect() error {
	consumer, err := kafka.NewConsumer(&k.config)
	if err != nil {
		return fmt.Errorf("Error on connect %+v", err)
	}
	k.consumer = consumer
	if err := k.consumer.SubscribeTopics(k.topics, nil); err != nil {
		return fmt.Errorf("Error on subscribe to topics %+v", err)
	}
	k.connected = true
	return nil
}

// Listen get every event on the consumer
// when a partition is assignated, assign to the consumer
// when a partition is revoked, unassign to the consumer
// when a message is received in the queue, get the message
// and send to a channel
// when is the EOF of a partition, just try to read again
// when gives an error, write on logger
func (k *KafkaConsumer) Listen() error {
	if !k.connected {
		return fmt.Errorf("Read fails: Not Connected")
	}
	k.logger.Info("Start Read: Waiting messages")
	for ev := range k.consumer.Events() {
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			if e.Partitions != nil {
				if k.consumer.Assign(e.Partitions) != nil {
					k.logger.Error("Error Assign Partitions: %+v", e.Partitions)
				}
			}
		case kafka.RevokedPartitions:
			if k.consumer.Unassign() != nil {
				k.logger.Error("Error Unassign Partitions: %+v", e.Partitions)
			}
		case *kafka.Message:
			k.messages <- e.Value
		case kafka.PartitionEOF:
			k.logger.Info("Waiting From Partition %+v...", e.Partition)
		case kafka.Error:
			// Errors should generally be considered as informational.
			// The client will try to automatically recover
			k.logger.Info("Info: %v", e.Error())
		default:
			k.logger.Info("Event not Handled %+v", e)
		}
	}
	return nil
}

// GetMessages returns the channel where the kafka messages flow
func (k *KafkaConsumer) GetMessages() chan []byte {
	return k.messages
}

// Close ends a connection and close a consumer
func (k *KafkaConsumer) Close() error {
	close(k.messages)
	k.connected = false
	return k.consumer.Close()
}
