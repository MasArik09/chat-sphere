package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"chatsphere/internal/auth"
	"chatsphere/internal/conversations"
	"chatsphere/internal/database"
	"chatsphere/internal/messages"
	"chatsphere/internal/users"
	"chatsphere/pkg/config"
)

func main() {
	log.Println("Starting ChatSphere API Server...")

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set up router
	router := gin.Default()

	// Configure CORS middleware (standard for cross-origin local dev)
	router.Use(corsMiddleware())

	// Initialize Repositories, Services, and Handlers
	userRepo := users.NewPostgresUserRepository(db)
	authService := auth.NewAuthService(userRepo, cfg)
	authHandler := auth.NewAuthHandler(authService, userRepo)

	conversationRepo := conversations.NewPostgresConversationRepository(db)
	conversationService := conversations.NewConversationService(conversationRepo)
	conversationHandler := conversations.NewConversationHandler(conversationService)

	messageRepo := messages.NewPostgresMessageRepository(db)
	messageService := messages.NewMessageService(messageRepo, conversationRepo)
	messageHandler := messages.NewMessageHandler(messageService)

	// API Routing Groups
	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", auth.AuthMiddleware(cfg), authHandler.Me)
		}

		convGroup := v1.Group("/conversations", auth.AuthMiddleware(cfg))
		{
			convGroup.POST("", conversationHandler.Create)
			convGroup.GET("", conversationHandler.List)
			convGroup.GET("/:id", conversationHandler.Detail)
			convGroup.POST("/:id/participants", conversationHandler.AddParticipant)
			convGroup.DELETE("/:id/participants/:userId", conversationHandler.RemoveParticipant)

			// Message routes nested under conversations
			convGroup.POST("/:id/messages", messageHandler.Send)
			convGroup.GET("/:id/messages", messageHandler.List)
		}
	}

	// Register health check endpoint
	router.GET("/health", healthCheckHandler(db))

	// Create HTTP Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Listening and serving HTTP on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func healthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ping database to ensure connectivity
		if err := db.Ping(); err != nil {
			log.Printf("Health check failed: database ping error: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection unhealthy",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
