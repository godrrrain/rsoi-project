package groupconsumer

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func InitConsumerGroup() (sarama.ConsumerGroup, error) {
	broker := "kafka:9092"
	groupID := "statistics-service"

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	var consumerGroup sarama.ConsumerGroup
	var err error

	maxRetries := 5
	retryInterval := 2 * time.Second

	for i := range maxRetries {
		consumerGroup, err = sarama.NewConsumerGroup([]string{broker}, groupID, config)
		if err == nil {
			return consumerGroup, nil
		}

		log.Printf("Waiting for Kafka... attempt %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
}
