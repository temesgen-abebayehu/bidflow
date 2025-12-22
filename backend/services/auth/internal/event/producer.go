package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/temesgen-abebayehu/bidflow/backend/common/kafka"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/domain"
)

type KafkaEventProducer struct {
	producer *kafka.Producer
}

func NewKafkaEventProducer(producer *kafka.Producer) domain.EventProducer {
	return &KafkaEventProducer{producer: producer}
}

func (p *KafkaEventProducer) PublishUserRegistered(ctx context.Context, user *domain.User) error {
	event := UserRegisteredEvent{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      user.Role,
		Timestamp: time.Now(),
	}
	return p.producer.Publish(ctx, TopicUserRegistered, user.ID.String(), event)
}

func (p *KafkaEventProducer) PublishUserVerified(ctx context.Context, userID uuid.UUID) error {
	event := UserVerifiedEvent{
		UserID:    userID,
		Timestamp: time.Now(),
	}
	return p.producer.Publish(ctx, TopicUserVerified, userID.String(), event)
}
