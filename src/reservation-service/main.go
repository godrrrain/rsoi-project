package main

import (
	"context"
	"fmt"

	"lab2/src/kafka"
	"lab2/src/middleware"
	"lab2/src/reservation-service/handler"
	"lab2/src/reservation-service/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		"postgres", 5432, "program", "reservations", "test")
	psqlDB, err := storage.NewPgStorage(context.Background(), postgresURL)
	if err != nil {
		fmt.Printf("Postgresql init: %s", err)
	} else {
		fmt.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	producer, err := kafka.NewSyncProducer([]string{"kafka:9092"})
	if err != nil {
		fmt.Printf("kafka producer init: %s", err)
	} else {
		fmt.Println("Connected to Kafka for reservation service")
	}
	defer kafka.CloseProducer(producer)

	handler := handler.NewHandler(psqlDB, producer)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	jwtMiddleware := middleware.NewJWTMiddleware("http://idp-service:8090")

	router.GET("/api/v1/reservations", jwtMiddleware.Middleware(), handler.GetReservations)
	router.GET("/api/v1/reservations/all", jwtMiddleware.Middleware(), handler.GetReservationsAll)
	router.GET("/api/v1/reservations/info/:uid", jwtMiddleware.Middleware(), handler.GetReservationByUid)
	router.GET("/api/v1/reservations/amount", jwtMiddleware.Middleware(), handler.GetRentedReservationAmount)
	router.POST("/api/v1/reservations", jwtMiddleware.Middleware(), handler.CreateReservation)
	router.PUT("/api/v1/reservations/:uid", jwtMiddleware.Middleware(), handler.UpdateReservationStatus)

	router.GET("/manage/health", handler.GetHealth)

	router.Run(":8070")
}
