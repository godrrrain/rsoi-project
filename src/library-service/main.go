package main

import (
	"context"
	"fmt"

	"lab2/src/library-service/handler"
	"lab2/src/library-service/storage"
	"lab2/src/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		"postgres", 5432, "program", "libraries", "test")
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

	router.GET("/api/v1/libraries", handler.GetLibrariesByCity)
	router.GET("/api/v1/libraries/:uid/books/", handler.GetBooksByLibraryUid)
	router.GET("/api/v1/libraries/:uid/", handler.GetLibraryByUid)
	router.GET("/api/v1/books/:uid/", handler.GetBookInfoByUid)

	router.PUT("/api/v1/books/:uid/condition", jwtMiddleware.Middleware(), handler.UpdateBookCondition)
	router.PUT("/api/v1/books/:uid/count/:inc/", jwtMiddleware.Middleware(), handler.UpdateBookCount)

	router.GET("/manage/health", handler.GetHealth)

	router.Run(":8060")
}
