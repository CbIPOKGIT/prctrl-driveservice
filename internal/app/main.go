package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CbIPOKGIT/prctrl-driveservice/internal/api"
	"github.com/CbIPOKGIT/prctrl-driveservice/internal/server"
	"github.com/gin-gonic/gin"
)

func StartApplication() {
	service := server.New()
	go service.Start()

	go startHttpService()

	handleBreak()
	log.Println("Application stopped")
}

func handleBreak() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
	log.Println("Shutting down...")
}

func startHttpService() {
	server := gin.New()

	apiGroup := server.Group("/api")

	apiGroup.GET("/file", api.LoadFile)
	apiGroup.PUT("/file", api.UploadFile)

	log.Printf("Starting http server on port %s", os.Getenv("API_PORT"))
	if err := server.Run(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
