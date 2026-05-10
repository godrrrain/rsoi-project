package main

import (
	"context"
	"fmt"

	"lab2/src/middleware"
	"lab2/src/rating-service/handler"
	"lab2/src/rating-service/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		"postgres", 5432, "program", "ratings", "test")
	psqlDB, err := storage.NewPgStorage(context.Background(), postgresURL)
	if err != nil {
		fmt.Printf("Postgresql init: %s", err)
	} else {
		fmt.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	handler := handler.NewHandler(psqlDB)

	router := gin.Default()

	router.Use(cors.Default())

	jwtMiddleware := middleware.NewJWTMiddleware("http://idp-service:8090")

	router.GET("/api/v1/rating/", jwtMiddleware.Middleware(), handler.GetRating)
	router.PUT("/api/v1/rating/", jwtMiddleware.Middleware(), handler.UpdateRating)

	router.GET("/manage/health", handler.GetHealth)

	router.Run(":8050")
}
