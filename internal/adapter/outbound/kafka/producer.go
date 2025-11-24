package kafka

import (
	"clean-architecture/internal/port/outbound"

	"github.com/IBM/sarama"
	"github.com/labstack/gommon/log"
)

type Kafka struct {
	producer sarama.SyncProducer
	brokers  []string
	config   *sarama.Config
}

func NewKafkaProducer(brokers []string, config *sarama.Config) (outbound.KafkaProducerInterface, error) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		brokers:  brokers,
		config:   config,
		producer: producer,
	}, nil
}

func (k *Kafka) ProduceMessage(topic string, data []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		log.Errorf("[Kafka-1] Failed to send message: %v", err)
		return err
	}

	log.Infof("[Kafka-2] Sent → topic=%s partition=%d offset=%d", topic, partition, offset)
	return nil
}

func (k *Kafka) Close() error {
	if err := k.producer.Close(); err != nil {
		log.Errorf("[Kafka-3] Failed to close producer: %v", err)
		return err
	}
	return nil
}
