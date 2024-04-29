package main

import (
	"lab2/src/gateway-service/handler"
	"lab2/src/jobqueue"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

func main() {

	var st gobreaker.Settings
	st.Name = "Library Circuit Breaker"
	libraryCb := gobreaker.NewCircuitBreaker(st)

	st.Name = "Rating Circuit Breaker"
	ratingCb := gobreaker.NewCircuitBreaker(st)

	st.Name = "Reservation Circuit Breaker"
	reservationCb := gobreaker.NewCircuitBreaker(st)

	jobScheduler := jobqueue.NewJobScheduler(10 * time.Second)
	jobScheduler.Start()

	handler := handler.NewHandler(libraryCb, ratingCb, reservationCb, jobScheduler)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/libraries", handler.GetLibrariesByCity)
	router.GET("/api/v1/libraries/:uid/books/", handler.GetBooksByLibraryUid)
	router.GET("/api/v1/rating/", handler.GetRating)
	router.GET("/api/v1/reservations", handler.GetReservations)
	router.POST("/api/v1/reservations", handler.CreateReservation)
	router.POST("/api/v1/reservations/:uid/return", handler.ReturnBook)

	router.GET("/manage/health", handler.GetHealth)

	router.Run()
}
