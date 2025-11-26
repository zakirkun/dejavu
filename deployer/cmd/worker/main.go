package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dejavu/deployer/internal/worker"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize worker
	w, err := worker.New()
	if err != nil {
		log.Fatal("Failed to initialize worker:", err)
	}
	defer w.Close()

	log.Println("🚀 Deployer Worker started, waiting for deployments...")

	// Start worker
	if err := w.Start(); err != nil {
		log.Fatal("Failed to start worker:", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker...")
}

