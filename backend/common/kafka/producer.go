package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger logger.Logger
}

func NewProducer(brokers []string, log logger.Logger) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}

	return &Producer{
		writer: w,
		logger: log,
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, message interface{}) error {
	value, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("failed to marshal message", zap.Error(err))
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
		Topic: topic,
		Time:  time.Now(),
	})

	if err != nil {
		p.logger.Error("failed to write message to kafka",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return err
	}

	p.logger.Info("message published to kafka",
		zap.String("topic", topic),
		zap.String("key", key),
	)

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
