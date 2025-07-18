package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eliothedeman/smol/control"
	"github.com/eliothedeman/smol/unit"
)

func main() {
	fmt.Println("smol - neural network system")
	log.Println("Starting smol...")

	registry := unit.NewRegistry()
	lifecycle := control.NewLifecycle()

	registry.Register("lifecycle", lifecycle)

	if err := registry.Start(); err != nil {
		log.Fatalf("Failed to start registry: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	go func() {
		lifecycle.AddTaskWithContext(func(ctx context.Context) {
			<-ctx.Done()
			log.Println("Shutting down lifecycle...")
			lifecycle.Shutdown()
		})
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	registry.Stop()
	log.Println("Shutdown complete")
}
