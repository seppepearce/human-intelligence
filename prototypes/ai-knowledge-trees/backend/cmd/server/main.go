package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/ai"
	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/db"
	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/handlers"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	log.Println("🗄️  Initializing database...")
	database, err := db.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize LocalAI client
	log.Println("🧠 Initializing LocalAI client...")
	localAI := ai.NewLocalAI()

	// Test LocalAI connection
	if err := localAI.Health(); err != nil {
		log.Printf("⚠️  LocalAI health check failed: %v", err)
		log.Println("🔄 AI features will use fallback mode until LocalAI is available")
	} else {
		log.Println("✅ LocalAI connection established")
	}

	// Initialize HTTP handlers
	handler := handlers.New(database, localAI)

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173", // SvelteKit dev server
			"http://localhost:3000", // Alternative dev port
			"http://localhost:8080", // Production port
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Allow all origins in development
	if gin.Mode() == gin.DebugMode {
		corsConfig.AllowOrigins = []string{"*"}
		corsConfig.AllowCredentials = false
	}

	router.Use(cors.New(corsConfig))

	// Custom middleware for request logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		latency := time.Since(start)

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%s] %s %s %d %d %v",
			clientIP,
			method,
			path,
			statusCode,
			bodySize,
			latency,
		)
	})

	// Register routes
	handler.RegisterRoutes(router)

	// Add a root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "🏛️ AI Knowledge Trees API",
			"description": "Where ancient wisdom meets modern intelligence",
			"version":     "1.0.0",
			"status":      "running",
			"endpoints": gin.H{
				"health": "/api/v1/health",
				"trees":  "/api/v1/trees",
				"nodes":  "/api/v1/nodes",
				"search": "/api/v1/search",
				"ai":     "/api/v1/ai",
				"docs":   "https://github.com/human-intelligence/ai-knowledge-trees",
			},
			"philosophy": "The unexamined life is not worth living. - Socrates",
		})
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		log.Printf("📖 API documentation: http://localhost:%s/", port)
		log.Printf("❤️  Health check: http://localhost:%s/api/v1/health", port)
		log.Printf("🌳 Tree management: http://localhost:%s/api/v1/trees", port)

		if localAI.IsAvailable() {
			log.Printf("🧠 AI features: http://localhost:%s/api/v1/ai", port)
		} else {
			log.Printf("⚠️  AI features running in fallback mode")
		}

		log.Println("🏛️  \"Wisdom begins in wonder.\" - Socrates")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server gracefully stopped")
	log.Println("🏛️  May wisdom guide your next session.")
}
