package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"lab2/src/rating-service/storage"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	storage storage.Storage
}

type RatingResponse struct {
	Stars int `json:"stars"`
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetRating(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
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

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	var reqRating RatingResponse

	err := json.NewDecoder(c.Request.Body).Decode(&reqRating)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	err = h.storage.UpdateRating(context.Background(), username, reqRating.Stars)
	if err != nil {
		fmt.Printf("failed to update raing %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "rating updated",
	})
}

func (h *Handler) GetHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
