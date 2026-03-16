package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer wraps kafka-go writer for publishing events.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer.
func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	return &Producer{writer: w}
}

// Publish sends an event to Kafka.
func (p *Producer) Publish(ctx context.Context, evt Event) error {
	evt.Timestamp = time.Now().Unix()

	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(evt.Type),
		Value: data,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("📤 Event published: %s (aggregate_id=%d)", evt.Type, evt.AggregateID)
	return nil
}

// Close closes the Kafka producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Consumer wraps kafka-go reader for consuming events.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: r}
}

// Consume reads events from Kafka and passes them to the handler.
func (c *Consumer) Consume(ctx context.Context, handler func(Event) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var evt Event
			if err := json.Unmarshal(msg.Value, &evt); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}

			if err := handler(evt); err != nil {
				log.Printf("Error handling event %s: %v", evt.Type, err)
			}
		}
	}
}

// Close closes the Kafka consumer.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
