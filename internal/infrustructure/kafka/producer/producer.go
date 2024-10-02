package kafka

import "github.com/IBM/sarama"

type Producer interface {
	SendMessage(topic string, message []byte) error
	Close() error
}

type kafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string, config *sarama.Config) (Producer, error) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{producer: producer}, nil
}

func (p *kafkaProducer) SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := p.producer.SendMessage(msg)
	return err
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}
