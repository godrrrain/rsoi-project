package main

import (
	"context"
	"lab2/src/statistics-service/groupconsumer"
	"lab2/src/statistics-service/handler"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const topic = "events"

func main() {
	consumerGroup, err := groupconsumer.InitConsumerGroup()
	if err != nil {
		log.Fatalf("failed to create kafka consumer group")
	}

	log.Println("Successfully connected to Kafka")

	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Printf("failed to close consumer group: %v", err)
		}
	}()

	consumer := handler.NewEventConsumer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumer); err != nil {
				log.Printf("error consuming kafka topic: %v", err)
				time.Sleep(50 * time.Millisecond)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("consumer group error: %v", err)
		}
	}()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/api/v1/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, consumer.Stats())
	})
	router.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.Run(":8040")
}
