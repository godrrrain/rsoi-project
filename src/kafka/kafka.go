package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func NewSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second
	config.Version = sarama.V2_0_0_0

	var producer sarama.SyncProducer
	var err error

	maxRetries := 5
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		producer, err = sarama.NewSyncProducer(brokers, config)
		if err == nil {
			return producer, nil
		}

		log.Printf("Waiting for Kafka producer... attempt %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to create kafka producer: %w", err)
}

func CloseProducer(producer sarama.SyncProducer) error {
	if producer == nil {
		return nil
	}
	return producer.Close()
}
