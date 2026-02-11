package eventing

import (
	"log/slog"
	"strings"

	"app/internal/config"
	"app/internal/eventing/kafka"
	"app/internal/service"
)

func BuildSubscriptionPublisher(cfg config.Config, logger *slog.Logger) service.SubscriptionEventPublisher {
	brokers := make([]string, 0, len(cfg.Kafka.Brokers))
	for _, broker := range cfg.Kafka.Brokers {
		trimmed := strings.TrimSpace(broker)
		if trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}

	if len(brokers) == 0 {
		if logger != nil {
			logger.Info("Kafka brokers not configured, subscription publisher disabled")
		}
		return nil
	}

	return kafka.NewSubscriptionPublisher(brokers, cfg.Kafka.SubscriptionTopic, logger)
}
