package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"budget-tracker/internal/api"
	"budget-tracker/internal/api/handlers"
	"budget-tracker/internal/repository"
	"budget-tracker/internal/services/ai"
)

func main() {
	log.Println("Starting Budget Tracker API server...")

	// Initialize database
	dbConfig := repository.NewConfigFromEnv()
	db, err := repository.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize AI client (optional - receipt processing won't work without it)
	var aiClient *ai.Client
	aiClient, err = ai.NewClientFromEnv()
	if err != nil {
		log.Printf("Warning: AI client not initialized: %v", err)
		log.Println("Receipt processing will be unavailable")
		aiClient = nil
	} else {
		log.Println("AI client initialized successfully")
	}

	// Initialize repositories
	budgetRepo := repository.NewBudgetRepository(db)
	expectedExpenseRepo := repository.NewExpectedExpenseRepository(db)
	actualExpenseRepo := repository.NewActualExpenseRepository(db)

	// Initialize handlers
	budgetHandler := handlers.NewBudgetHandler(budgetRepo)
	expectedExpenseHandler := handlers.NewExpectedExpenseHandler(expectedExpenseRepo)
	actualExpenseHandler := handlers.NewActualExpenseHandler(actualExpenseRepo)
	receiptHandler := handlers.NewReceiptHandler(aiClient, expectedExpenseRepo, actualExpenseRepo)
	notificationHandler := handlers.NewNotificationHandler(budgetRepo, expectedExpenseRepo, actualExpenseRepo)

	// Create router with all handlers
	h := &api.Handlers{
		Budget:          budgetHandler,
		ExpectedExpense: expectedExpenseHandler,
		ActualExpense:   actualExpenseHandler,
		Receipt:         receiptHandler,
		Notification:    notificationHandler,
	}
	router := api.NewRouter(h)

	// Apply middleware
	handler := api.Chain(
		router,
		api.Recovery,
		api.Logger,
		api.CORS(api.DefaultCORSConfig()),
	)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second, // Longer timeout for AI processing
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
