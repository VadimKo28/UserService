package kafka

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"strconv"

	"app/internal/domain/subscription"

	kafkago "github.com/segmentio/kafka-go"
)

type SubscriptionPublisher struct {
	writer *kafkago.Writer
}

func NewSubscriptionPublisher(brokers []string, topic string, logger *slog.Logger) *SubscriptionPublisher {
	errorLogger := log.New(os.Stderr, "kafka writer: ", log.LstdFlags)
	if logger != nil {
		errorLogger = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	}

	return &SubscriptionPublisher{
		writer: &kafkago.Writer{
			Addr:        kafkago.TCP(brokers...),
			Topic:       topic,
			Balancer:    &kafkago.LeastBytes{},
			Async:       true,
			ErrorLogger: errorLogger,
		},
	}
}

func (p *SubscriptionPublisher) Publish(ctx context.Context, event subscription.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(strconv.Itoa(event.Payload.UserID)),
		Value: payload,
		Time:  event.OccurredAt,
	})
}

func (p *SubscriptionPublisher) Close() error {
	return p.writer.Close()
}
