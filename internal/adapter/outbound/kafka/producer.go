package kafka

import (
	"github.com/IBM/sarama"
	"github.com/labstack/gommon/log"
)

type IKafka interface {
	ProduceMessage(string, []byte) error
}

type Kafka struct {
	brokers []string
	config  *sarama.Config
}

func NewKafkaProducer(brokers []string, config *sarama.Config) IKafka {
	return &Kafka{
		brokers: brokers,
		config:  config,
	}
}

func (k *Kafka) ProduceMessage(topic string, data []byte) error {
	producer, err := sarama.NewSyncProducer(k.brokers, k.config)
	if err != nil {
		log.Errorf("[Kafka-1] Failed to create producer: %v", err)
		return err
	}

	defer func(producer sarama.SyncProducer) {
		err = producer.Close()
		if err != nil {
			log.Errorf("[Kafka-2] Failed to close producer: %v", err)
			return
		}
	}(producer)

	message := &sarama.ProducerMessage{
		Topic:   topic,
		Headers: nil,
		Value:   sarama.ByteEncoder(data),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Errorf("[Kafka-3] Failed to produce message to kafka: %v", err)
		return err
	}

	log.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
