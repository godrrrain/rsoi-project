package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"lab2/src/rating-service/storage"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	storage  storage.Storage
	producer sarama.SyncProducer
}

type RatingResponse struct {
	Stars int `json:"stars"`
}

type UpdateRatingRequest struct {
	Stars    int    `json:"stars"`
	Username string `json:"username"`
}

func NewHandler(storage storage.Storage, producer sarama.SyncProducer) *Handler {
	return &Handler{storage: storage, producer: producer}
}

func (h *Handler) GetRating(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "empty username",
		})
		return
	}

	rating, err := h.storage.GetRating(context.Background(), username)

	if err != nil {
		fmt.Printf("failed to get rating %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RatingResponse{
		Stars: rating.Stars,
	})
}

func (h *Handler) UpdateRating(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "empty username",
		})
		return
	}

	var reqRating UpdateRatingRequest

	err := json.NewDecoder(c.Request.Body).Decode(&reqRating)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if reqRating.Username != "" {
		username = reqRating.Username
	}

	err = h.storage.UpdateRating(context.Background(), username, reqRating.Stars)
	if err != nil {
		fmt.Printf("failed to update raing %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	h.sendEvent("Рейтинг обновился")

	c.JSON(http.StatusOK, MessageResponse{
		Message: "rating updated",
	})
}

func (h *Handler) GetHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *Handler) sendEvent(event string) {
	if h.producer == nil {
		return
	}

	const topic = "events"

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(event),
	}

	_, _, err := h.producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("failed to send kafka event: %s\n", err.Error())
	}
}
