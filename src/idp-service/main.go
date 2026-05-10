package main

import (
	"context"
	"fmt"
	"time"

	"lab2/src/idp-service/auth"
	"lab2/src/idp-service/handler"
	"lab2/src/idp-service/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		"postgres", 5432, "program", "idp", "test")

	var db *storage.PgStorage
	var err error

	for i := 0; i < 10; i++ {
		db, err = storage.NewPgStorage(context.Background(), postgresURL)
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to database (attempt %d/10): %s\n", i+1, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		fmt.Printf("Failed to connect to database after 10 attempts: %s\n", err)
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected to PostgreSQL")

	jwtManager, err := auth.NewJWTManager("http://idp-service:8090")
	if err != nil {
		fmt.Printf("Failed to create JWT manager: %s\n", err)
		panic(err)
	}

	h := handler.NewHandler(db, jwtManager, "http://idp-service:8090")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/.well-known/jwks.json", h.GetJWKS)

	router.GET("/oauth2/authorize", h.Authorize)
	router.POST("/oauth2/login", h.Login)
	router.GET("/oauth2/signup", h.Signup)
	router.POST("/oauth2/signup", h.Signup)
	router.POST("/oauth2/token", h.Token)
	router.GET("/oauth2/userinfo", h.UserInfo)

	router.POST("/users", h.CreateUser)
	router.GET("/users/me", h.GetMe)

	router.GET("/manage/health", h.GetHealth)

	fmt.Println("Starting Identity Provider Service on port 8090")
	router.Run(":8090")
}
