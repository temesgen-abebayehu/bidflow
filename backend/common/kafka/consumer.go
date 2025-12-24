package kafka

import (
	"context"
	"io"

	"github.com/segmentio/kafka-go"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
	logger logger.Logger
}

type Handler func(ctx context.Context, topic string, key, value []byte) error

func NewConsumer(brokers []string, topics []string, groupID string, log logger.Logger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupTopics: topics,
		GroupID:     groupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})

	return &Consumer{
		reader: r,
		logger: log,
	}
}

func (c *Consumer) Start(ctx context.Context, handler Handler) {
	go func() {
		for {
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if err == io.EOF {
					c.logger.Info("kafka reader closed")
					return
				}
				c.logger.Error("failed to fetch message", zap.Error(err))
				continue
			}

			if err := handler(ctx, m.Topic, m.Key, m.Value); err != nil {
				c.logger.Error("failed to handle message", zap.Error(err))
				// Decide whether to commit or not based on error type
				// For now, we continue, but in production you might want retry logic
			}

			if err := c.reader.CommitMessages(ctx, m); err != nil {
				c.logger.Error("failed to commit message", zap.Error(err))
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
