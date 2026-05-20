package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-project/internal/components/account"
	"task-project/internal/components/address"
	"task-project/internal/components/card"
	"task-project/internal/components/customer"
	"task-project/internal/components/person"
	transaction "task-project/internal/components/transactions"
	"task-project/internal/config"
	"task-project/internal/router"
)

func main() {
	fmt.Println("Starting Banking API...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	dbConnection, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	fmt.Println("Successfully connected to database")
	defer dbConnection.Connection.Close()

	// Create repositories
	personRepo := person.NewRepository(dbConnection.Connection)
	customerRepo := customer.NewRepository(dbConnection.Connection)
	addressRepo := address.NewRepository(dbConnection.Connection)
	accountRepo := account.NewRepository(dbConnection.Connection)
	cardRepo := card.NewRepository(dbConnection.Connection)
	transactionRepo := transaction.NewRepository(dbConnection.Connection)

	// Create handlers
	personHandler := person.NewHandler(personRepo)
	customerHandler := customer.NewHandler(customerRepo)
	addressHandler := address.NewHandler(addressRepo)
	accountHandler := account.NewHandler(accountRepo)
	cardHandler := card.NewHandler(cardRepo)
	transactionHandler := transaction.NewHandler(transactionRepo)

	// Setup router
	r := router.SetupRouter(
		cfg,
		personHandler,
		customerHandler,
		addressHandler,
		accountHandler,
		cardHandler,
		transactionHandler,
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Create channel for signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("\nShutting down server...")

	// Create timeout context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped")
	log.Println("Database connection closed")
	log.Println("Banking API exited gracefully")
}
