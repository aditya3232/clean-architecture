package config

import (
	"time"

	"github.com/IBM/sarama"
)

func (cfg Config) NewKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = cfg.Kafka.MaxRetry
	config.Producer.Retry.Backoff = time.Duration(cfg.Kafka.TimeoutInMS) * time.Millisecond
	config.Version = sarama.V2_1_0_0
	return config
}
