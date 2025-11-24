package outbound

type KafkaProducerInterface interface {
	ProduceMessage(topic string, data []byte) error
	Close() error
}
