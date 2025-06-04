package events

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// KafkaProducer handles publishing events to Kafka
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishEvent publishes an order event to Kafka
func (p *KafkaProducer) PublishEvent(event *OrderEvent) error {
	data, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.OrderID), // Use order ID as partition key
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf("Published event %s to topic %s, partition %d, offset %d",
		event.ID, p.topic, partition, offset)
	return nil
}

// Close closes the producer
func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}

// KafkaConsumer handles consuming events from Kafka
type KafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
}

// EventHandler defines how to handle received events
type EventHandler func(event *OrderEvent) error

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

// ConsumeEvents consumes events from Kafka and handles them
func (c *KafkaConsumer) ConsumeEvents(handler EventHandler) error {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			event, err := FromJSON(message.Value)
			if err != nil {
				log.Printf("Failed to deserialize event: %v", err)
				continue
			}

			if err := handler(event); err != nil {
				log.Printf("Failed to handle event %s: %v", event.ID, err)
			}

		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v", err)
		}
	}
}

// Close closes the consumer
func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}
