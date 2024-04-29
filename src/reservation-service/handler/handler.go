package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"lab2/src/reservation-service/storage"

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

type RequestCreateReservation struct {
	BookUid    string `json:"bookUid"`
	LibraryUid string `json:"libraryUid"`
	TillDate   string `json:"tillDate"`
}

type RequestUpdateReservation struct {
	Condition string `json:"condition"`
	Date      string `json:"date"`
}

type ReservationResponse struct {
	Reservation_uid string `json:"reservationUid"`
	Username        string `json:"username"`
	Book_uid        string `json:"bookUid"`
	Library_uid     string `json:"libraryUid"`
	Status          string `json:"status"`
	Start_date      string `json:"startDate"`
	Till_date       string `json:"tillDate"`
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetReservations(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	reservations, err := h.storage.GetReservations(context.Background(), username)

	if err != nil {
		fmt.Printf("failed to get reservations %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReservationsToResponse(reservations))
}

func (h *Handler) GetReservationByUid(c *gin.Context) {

	reservation, err := h.storage.GetReservationByUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get reservation %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReservationToResponse(reservation))
}

func (h *Handler) GetRentedReservationAmount(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	reservationAmount, err := h.storage.GetRentedReservationAmount(context.Background(), username)

	if err != nil {
		fmt.Printf("failed to get reservation amount %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reservationAmount)
}

func (h *Handler) CreateReservation(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	var reqCrRes RequestCreateReservation

	err := json.NewDecoder(c.Request.Body).Decode(&reqCrRes)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	reservation, err := h.storage.CreateReservation(context.Background(), username, reqCrRes.BookUid, reqCrRes.LibraryUid, reqCrRes.TillDate)

	if err != nil {
		fmt.Printf("failed to create reservations %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReservationToResponse(reservation))
}

func (h *Handler) UpdateReservationStatus(c *gin.Context) {

	reservation, err := h.storage.GetReservationByUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var reqUpdRes RequestUpdateReservation

	err = json.NewDecoder(c.Request.Body).Decode(&reqUpdRes)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	date, err := time.Parse("2006-01-02", reqUpdRes.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	status := "RETURNED"
	if date.After(reservation.Till_date) {
		status = "EXPIRED"
	}

	err = h.storage.UpdateReservationStatus(context.Background(), c.Param("uid"), status)

	if err != nil {
		fmt.Printf("failed to update reservation %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if status == "EXPIRED" {
		c.JSON(http.StatusNoContent, MessageResponse{
			Message: "status updated",
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "status updated",
	})
}

func ReservationToResponse(reservation storage.Reservation) ReservationResponse {
	return ReservationResponse{
		Reservation_uid: reservation.Reservation_uid,
		Username:        reservation.Username,
		Book_uid:        reservation.Book_uid,
		Library_uid:     reservation.Library_uid,
		Status:          reservation.Status,
		Start_date:      reservation.Start_date.Format("2006-01-02"),
		Till_date:       reservation.Till_date.Format("2006-01-02"),
	}
}

func ReservationsToResponse(reservations []storage.Reservation) []ReservationResponse {
	if reservations == nil {
		return nil
	}

	res := make([]ReservationResponse, len(reservations))

	for index, value := range reservations {
		res[index] = ReservationToResponse(value)
	}

	return res
}

func (h *Handler) GetHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
