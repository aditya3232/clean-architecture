package service

import (
	"clean-architecture/config"
	"clean-architecture/internal/domain/entity"
	"clean-architecture/internal/port/outbound"
	"clean-architecture/utils"
	"context"
	"encoding/json"
	"time"
)

type KafkaServiceInterface interface {
	PublishMessage(ctx context.Context, req entity.PublishMessage) error
}

type kafkaService struct {
	cfg   *config.Config
	kafka outbound.KafkaProducerInterface
}

func NewKafkaService(cfg *config.Config, kafka outbound.KafkaProducerInterface) KafkaServiceInterface {
	return &kafkaService{
		cfg:   cfg,
		kafka: kafka,
	}
}

func (s *kafkaService) produceToKafka(req entity.PublishMessage) error {
	event := entity.KafkaEvent{
		Name: "message_published",
	}

	metadata := entity.KafkaMetaData{
		Sender:    "clean_architecture_service",
		SendingAt: time.Now().Format(time.RFC3339),
	}

	notifType := "EMAIL"
	if req.QueueName == utils.PUSH_NOTIF {
		notifType = "PUSH"
	}

	body := entity.KafkaBody{
		Type: "JSON",
		Data: &entity.KafkaData{
			ReceiverEmail:    req.Email,
			Message:          req.Message,
			ReceiverId:       req.UserId,
			Subject:          req.Subject,
			NotificationType: notifType,
		},
	}

	kafkaMessage := entity.KafkaMessage{
		Event:    event,
		Metadata: metadata,
		Body:     body,
	}

	topic := s.cfg.Kafka.Topic
	kafkaMessageJSON, _ := json.Marshal(kafkaMessage)
	err := s.kafka.ProduceMessage(topic, kafkaMessageJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *kafkaService) PublishMessage(ctx context.Context, req entity.PublishMessage) error {
	err := s.produceToKafka(req)
	if err != nil {
		return err
	}
	return nil
}
