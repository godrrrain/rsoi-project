package handler

import (
	"sync"

	"github.com/IBM/sarama"
)

type EventConsumer struct {
	stats map[string]int
	mu    sync.RWMutex
}

type StatisticsResponse map[string]int

func NewEventConsumer() *EventConsumer {
	return &EventConsumer{stats: make(map[string]int)}
}

func (c *EventConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *EventConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *EventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		event := string(message.Value)
		c.mu.Lock()
		c.stats[event]++
		c.mu.Unlock()
		session.MarkMessage(message, "")
	}
	return nil
}

func (c *EventConsumer) Stats() StatisticsResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(StatisticsResponse, len(c.stats))
	for k, v := range c.stats {
		result[k] = v
	}
	return result
}
