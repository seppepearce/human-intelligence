package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"human-intelligence/internal/database"
	"human-intelligence/internal/db/seeds"
	"human-intelligence/internal/handlers"
	"human-intelligence/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Server struct {
	db         *database.DB
	router     *gin.Engine
	httpServer *http.Server
	wsUpgrader websocket.Upgrader
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Check for seeding command
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		runSeeding()
		return
	}

	// Create server instance
	server := &Server{
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
		},
	}

	// Initialize database
	if err := server.initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer server.db.Close()

	// Setup router and routes
	server.setupRouter()
	server.setupRoutes()

	// Create HTTP server
	port := getEnv("PORT", "8080")
	server.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      server.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting Human Intelligence server on port %s", port)
		log.Printf("Environment: %s", getEnv("GIN_MODE", "debug"))

		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	server.gracefulShutdown()
}

func (s *Server) initDatabase() error {
	config := database.DefaultConfig()

	db, err := database.NewDB(config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	s.db = db

	// Run migrations
	migrationsDir := "migrations"
	if err := s.db.RunMigrations(migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (s *Server) setupRouter() {
	// Set Gin mode based on environment
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.New()

	// Middleware
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(corsMiddleware())
	s.router.Use(securityMiddleware())
}

func (s *Server) setupRoutes() {
	// Create handlers
	userHandler := handlers.NewUserHandler(s.db)
	nodeHandler := handlers.NewNodeHandler(s.db)
	pathHandler := handlers.NewPathHandler(s.db)
	searchHandler := handlers.NewSearchHandler(s.db)

	// Health check
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/metrics", s.metrics)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(authMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
				users.GET("/:id", userHandler.GetUser)
			}

			// Node routes
			nodes := protected.Group("/nodes")
			{
				nodes.POST("/", nodeHandler.CreateNode)
				nodes.GET("/:id", nodeHandler.GetNode)
				nodes.PUT("/:id", nodeHandler.UpdateNode)
				nodes.DELETE("/:id", nodeHandler.DeleteNode)
				nodes.POST("/:id/vote", nodeHandler.VoteNode)
			}

			// Learning Path routes
			paths := protected.Group("/paths")
			{
				paths.POST("/", pathHandler.CreatePath)
				paths.GET("/:id", pathHandler.GetPath)
				paths.PUT("/:id", pathHandler.UpdatePath)
				paths.DELETE("/:id", pathHandler.DeletePath)
				paths.POST("/:id/fork", pathHandler.ForkPath)
				paths.POST("/:id/nodes", pathHandler.AddNodeToPath)
				paths.DELETE("/:id/nodes/:nodeId", pathHandler.RemoveNodeFromPath)
				paths.POST("/:id/complete", pathHandler.CompletePath)
				paths.POST("/:id/tldr", pathHandler.CreateTLDR)
				paths.POST("/:id/vote", pathHandler.VotePath)
			}

			// Search routes
			search := protected.Group("/search")
			{
				search.GET("/", searchHandler.Search)
				search.GET("/semantic", searchHandler.SemanticSearch)
				search.GET("/suggestions", searchHandler.GetSuggestions)
			}

			// Activity feed
			protected.GET("/activity", s.getActivityFeed)
			protected.GET("/activity/live", s.liveActivityFeed)
		}

		// Public routes (no auth required)
		public := v1.Group("/public")
		{
			public.GET("/nodes", nodeHandler.GetPublicNodes)
			public.GET("/paths", pathHandler.GetPublicPaths)
			public.GET("/activity", s.getPublicActivity)
		}
	}

	// WebSocket endpoint for real-time features
	s.router.GET("/ws", s.handleWebSocket)

	// Serve static files in development
	if gin.Mode() == gin.DebugMode {
		s.router.Static("/static", "./static")
		// s.router.LoadHTMLGlob("templates/*") // Commented out - no templates needed for API
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	// Check database health
	if err := s.db.Health(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   getEnv("APP_VERSION", "dev"),
	})
}

func (s *Server) metrics(c *gin.Context) {
	stats := s.db.Stats()

	c.JSON(http.StatusOK, gin.H{
		"database": gin.H{
			"open_connections":     stats.OpenConnections,
			"in_use":               stats.InUse,
			"idle":                 stats.Idle,
			"wait_count":           stats.WaitCount,
			"wait_duration":        stats.WaitDuration.String(),
			"max_idle_closed":      stats.MaxIdleClosed,
			"max_idle_time_closed": stats.MaxIdleTimeClosed,
			"max_lifetime_closed":  stats.MaxLifetimeClosed,
		},
	})
}

func (s *Server) getActivityFeed(c *gin.Context) {
	// TODO: Implement activity feed logic
	c.JSON(http.StatusOK, gin.H{
		"activities": []models.ActivityEvent{},
		"message":    "Activity feed implementation pending",
	})
}

func (s *Server) getPublicActivity(c *gin.Context) {
	// TODO: Implement public activity feed
	c.JSON(http.StatusOK, gin.H{
		"activities": []models.ActivityEvent{},
		"message":    "Public activity feed implementation pending",
	})
}

func (s *Server) liveActivityFeed(c *gin.Context) {
	// Upgrade to WebSocket for real-time activity
	conn, err := s.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// TODO: Implement real-time activity streaming
	log.Println("WebSocket connection established for activity feed")

	// Keep connection alive and send periodic updates
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteJSON(gin.H{
				"type":      "heartbeat",
				"timestamp": time.Now().UTC(),
			}); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := s.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	// Handle WebSocket messages
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Echo message back (for now)
		if err := conn.WriteJSON(gin.H{
			"type":      "echo",
			"data":      msg,
			"timestamp": time.Now().UTC(),
		}); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// runSeeding initializes database and runs seeding
func runSeeding() {
	log.Println("🌱 Starting database seeding process...")

	// Initialize database connection
	config := database.DefaultConfig()
	db, err := database.NewDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database for seeding: %v", err)
	}
	defer db.Close()

	// Run migrations first
	migrationsDir := "migrations"
	if err := db.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create seeder and run seeding
	seeder := seeds.NewSeeder(db.DB)
	if err := seeder.SeedAll(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("✅ Database seeding completed successfully!")
}

// Middleware functions
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// In production, implement proper CORS policy
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Parse and validate token
		secret := getJWTSecret()
		log.Printf("JWT Debug - Token: %s...", tokenString[:50]) // First 50 chars
		log.Printf("JWT Debug - Secret: %s", secret)

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("JWT Debug - Parse error: %v", err)
			log.Printf("JWT Debug - Token valid: %v", token != nil && token.Valid)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Store user info in context
		log.Printf("JWT Debug - Successfully authenticated user: %s (ID: %s)", claims.Username, claims.UserID)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// Get JWT secret from environment
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Use a default secret for development
		return "your_super_secret_jwt_key_change_this_in_production"
	}
	return secret
}

// Utility function
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
