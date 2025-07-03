package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// Global flag to track shutdown state
var shuttingDown int32

func health(w http.ResponseWriter, req *http.Request) {
	// Check if we're in shutdown mode
	if atomic.LoadInt32(&shuttingDown) == 1 {
		fmt.Println("health: Shutdown")

		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Shutdown")
		return
	}
	fmt.Println("health: Ok")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok")
}

func main() {
	// Create a new HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)

	server := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server starting on :5000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	fmt.Println("\nReceived interrupt signal, preparing for graceful shutdown...")

	// Immediately set shutdown flag so /health returns 500
	atomic.StoreInt32(&shuttingDown, 1)
	fmt.Println("/health endpoint now returns 500 status")

	// Wait 5 seconds before starting graceful shutdown
	fmt.Println("Waiting 6 seconds before shutting down...")
	time.Sleep(6 * time.Second)

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	fmt.Println("Starting graceful shutdown...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}
