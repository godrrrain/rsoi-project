package main

import (
	"lab2/src/gateway-service/handler"
	"lab2/src/jobqueue"
	"lab2/src/middleware"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	jwtMiddleware := middleware.NewJWTMiddleware("http://idp-service:8090")

	router.GET("/api/v1/libraries", handler.GetLibrariesByCity)
	router.GET("/api/v1/libraries/:uid/books/", handler.GetBooksByLibraryUid)

	router.GET("/api/v1/rating/", jwtMiddleware.Middleware(), handler.GetRating)
	router.GET("/api/v1/reservations", jwtMiddleware.Middleware(), handler.GetReservations)
	router.GET("/api/v1/reservations/all", jwtMiddleware.Middleware(), handler.GetReservationsAll)
	router.POST("/api/v1/reservations", jwtMiddleware.Middleware(), handler.CreateReservation)
	router.POST("/api/v1/reservations/:uid/return", jwtMiddleware.Middleware(), handler.ReturnBook)

	router.GET("/manage/health", handler.GetHealth)

	router.Run()
}
