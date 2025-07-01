package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"human-intelligence/internal/database"
	"human-intelligence/internal/handlers"
	"human-intelligence/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	neo4j       *database.Neo4jService
	treeHandler *handlers.TreeHandler
	router      *gin.Engine
	server      *http.Server
}

func main() {
	// Initialize server
	server := &Server{}

	// Setup Neo4j connection
	if err := server.initNeo4j(); err != nil {
		log.Fatalf("Failed to initialize Neo4j: %v", err)
	}
	defer server.neo4j.Close(context.Background())

	// Setup routes
	server.setupRouter()
	server.setupRoutes()

	// Start server
	port := getEnv("PORT", "8080")
	server.server = &http.Server{
		Addr:    ":" + port,
		Handler: server.router,
	}

	// Start in goroutine
	go func() {
		log.Printf("🚀 Neo4j Knowledge Platform starting on port %s", port)
		if err := server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	server.gracefulShutdown()
}

func (s *Server) initNeo4j() error {
	config := &database.Neo4jConfig{
		URI:                          getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Username:                     getEnv("NEO4J_USERNAME", "neo4j"),
		Password:                     getEnv("NEO4J_PASSWORD", "hi_password"),
		Database:                     getEnv("NEO4J_DATABASE", "knowledgegraph"),
		MaxConnectionPoolSize:        50,
		ConnectionTimeout:            30 * time.Second,
		MaxTransactionRetries:        3,
		InitialRetryDelay:            time.Second,
		MaxRetryDelay:                30 * time.Second,
		RetryDelayMultiplier:         2.0,
		ConnectionAcquisitionTimeout: 60 * time.Second,
	}

	neo4jService, err := database.NewNeo4jService(config)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j service: %w", err)
	}

	s.neo4j = neo4jService
	s.treeHandler = handlers.NewTreeHandler(neo4jService)
	log.Println("✅ Connected to Neo4j successfully")
	return nil
}

func (s *Server) setupRouter() {
	s.router = gin.Default()

	// CORS middleware - permissive for development
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}))
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Users
		api.POST("/users", s.createUser)
		api.GET("/users/:id", s.getUser)
		api.GET("/users/username/:username", s.getUserByUsername)

		// Nodes
		api.POST("/nodes", s.createNode)
		api.GET("/nodes/:id", s.getNode)
		api.GET("/nodes", s.searchNodes)

		// Enhanced Trees with hierarchical structure
		api.POST("/trees", s.treeHandler.CreateTree)
		api.GET("/trees/:treeId", s.treeHandler.GetTreeStructure)
		api.PUT("/trees/:treeId", s.treeHandler.UpdateTreeMetadata)
		api.DELETE("/trees/:treeId", s.treeHandler.DeleteTree)
		api.POST("/trees/:treeId/nodes", s.treeHandler.AddNodeToTree)
		api.GET("/trees/:treeId/nodes", s.treeHandler.GetTreeNodes)

		// Legacy tree endpoints (backward compatibility)
		api.POST("/trees/:treeId/nodes/:nodeId", s.addNodeToTree)

		// Tags
		api.POST("/tags", s.createTag)
		api.POST("/nodes/:nodeId/tags/:tagName", s.tagNode)
		api.GET("/tags/popular", s.getPopularTags)

		// Graph visualization
		api.GET("/graph/visualization/:userId", s.getVisualizationData)

		// Search
		api.GET("/search", s.search)

		// Debug endpoint
		api.POST("/debug/simple-node", s.createSimpleNode)
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	ctx := context.Background()
	if err := s.neo4j.HealthCheck(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"database":  "neo4j",
	})
}

// User endpoints
func (s *Server) createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = generateID()
	user.IsActive = true

	ctx := context.Background()
	if err := s.neo4j.CreateUser(ctx, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (s *Server) getUser(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	user, err := s.neo4j.GetUserByID(ctx, id)
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) getUserByUsername(c *gin.Context) {
	username := c.Param("username")
	ctx := context.Background()

	user, err := s.neo4j.GetUserByUsername(ctx, username)
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Node endpoints
func (s *Server) createNode(c *gin.Context) {
	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node.ID = generateID()

	ctx := context.Background()
	if err := s.neo4j.CreateNode(ctx, &node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, node)
}

func (s *Server) getNode(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	node, err := s.neo4j.GetNodeByID(ctx, id)
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, node)
}

func (s *Server) searchNodes(c *gin.Context) {
	query := &models.GraphQuery{
		SearchString: c.Query("q"),
		Limit:        getIntParam(c, "limit", 20),
		Offset:       getIntParam(c, "offset", 0),
	}

	if tags := c.QueryArray("tags"); len(tags) > 0 {
		query.Tags = tags
	}

	ctx := context.Background()
	response, err := s.neo4j.SearchNodes(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Legacy tree endpoints for backward compatibility
func (s *Server) createTreeLegacy(c *gin.Context) {
	var tree models.Tree
	if err := c.ShouldBindJSON(&tree); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tree.ID = generateID()

	ctx := context.Background()
	if err := s.neo4j.CreateTree(ctx, &tree); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tree)
}

func (s *Server) addNodeToTree(c *gin.Context) {
	treeID := c.Param("treeId")
	nodeID := c.Param("nodeId")
	position := getIntParam(c, "position", 0)

	ctx := context.Background()
	if err := s.neo4j.AddNodeToTree(ctx, treeID, nodeID, position); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node added to tree successfully"})
}

// Tag endpoints
func (s *Server) createTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag.ID = generateID()

	ctx := context.Background()
	if err := s.neo4j.CreateTag(ctx, &tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

func (s *Server) tagNode(c *gin.Context) {
	nodeID := c.Param("nodeId")
	tagName := c.Param("tagName")

	ctx := context.Background()
	if err := s.neo4j.TagNode(ctx, nodeID, tagName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node tagged successfully"})
}

func (s *Server) getPopularTags(c *gin.Context) {
	limit := getIntParam(c, "limit", 10)

	ctx := context.Background()
	tags, err := s.neo4j.GetPopularTags(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// Graph visualization
func (s *Server) getVisualizationData(c *gin.Context) {
	userID := c.Param("userId")
	depth := getIntParam(c, "depth", 2)
	limit := getIntParam(c, "limit", 100)

	ctx := context.Background()
	data, err := s.neo4j.GetVisualizationData(ctx, userID, depth, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// Search endpoint
func (s *Server) search(c *gin.Context) {
	searchString := c.Query("q")
	if searchString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	query := &models.GraphQuery{
		SearchString: searchString,
		Limit:        getIntParam(c, "limit", 20),
		Offset:       getIntParam(c, "offset", 0),
	}

	ctx := context.Background()
	response, err := s.neo4j.SearchNodes(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

// Utility functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntParam(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (s *Server) createSimpleNode(c *gin.Context) {
	ctx := context.Background()

	// Create a simple node directly with Neo4j query
	query := `
		CREATE (n:Node {
			id: $id,
			title: $title,
			content: $content,
			content_type: 'text',
			is_public: true,
			created_at: datetime(),
			updated_at: datetime(),
			owner_id: 'debug-user'
		})
		RETURN n.id as id, n.title as title, n.created_at as created_at
	`

	params := map[string]interface{}{
		"id":      fmt.Sprintf("debug-node-%d", time.Now().UnixNano()),
		"title":   "Debug Node",
		"content": "This is a debug node created directly",
	}

	records, err := s.neo4j.ExecuteQuery(ctx, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(records) > 0 {
		record := records[0]
		result := map[string]interface{}{}

		if id, ok := record.Get("id"); ok {
			result["id"] = id
		}
		if title, ok := record.Get("title"); ok {
			result["title"] = title
		}
		if createdAt, ok := record.Get("created_at"); ok {
			result["created_at"] = createdAt
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Debug node created successfully",
			"node":    result,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No records returned"})
	}
}



func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
