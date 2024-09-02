package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CbIPOKGIT/prctrl-driveservice/internal/server"
)

func StartApplication() {
	service := server.New()
	service.Start()

	handleBreak()
	log.Println("Application stopped")
}

func handleBreak() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
	log.Println("Shutting down...")
}
