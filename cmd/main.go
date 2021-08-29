package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"log"
	"net/http"
	"os"
	"os/signal"
	"simple-go-auth/pkg/auth"
	"simple-go-auth/pkg/handler"
	"simple-go-auth/pkg/middleware"
	"simple-go-auth/pkg/redis"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	appAddr := ":" + os.Getenv("PORT")

	// redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	// TODO: config in json instead
	redisInfo := redis.RedisClientInfo{redis_host, redis_port, redis_password}

	redisClient := redis.NewRedisClient(redisInfo)

	var authService = auth.NewAuthService(redisClient)
	var token = auth.NewToken(os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET"))
	var handlers = handlers.NewHandlers(authService, token)

	var router = gin.Default()
	router.POST("/login", handlers.Login)
	router.POST("/todo", middleware.TokenAuthMiddleware(), handlers.CreateTodo)
	router.POST("/logout", middleware.TokenAuthMiddleware(), handlers.Logout)
	router.POST("/refresh", handlers.Refresh)

	srv := &http.Server{
		Addr:    appAddr,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %s\n", err)
		}
	}()

	//Wait for interrupt signal to gracefully shutdown the server
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigs
	log.Println("Shuting down the server ...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Failed to gracefully shutdown the server:", err)
	}
	log.Println("Finished shutdown")
}
