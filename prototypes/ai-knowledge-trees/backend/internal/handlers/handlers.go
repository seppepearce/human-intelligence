package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/ai"
	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/db"
	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/models"
)

// Handler represents the HTTP handler with its dependencies
type Handler struct {
	db      *db.DB
	localAI *ai.LocalAI
}

// New creates a new handler instance
func New(database *db.DB, localAI *ai.LocalAI) *Handler {
	return &Handler{
		db:      database,
		localAI: localAI,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", h.HealthCheck)

		// Tree routes
		trees := v1.Group("/trees")
		{
			trees.GET("", h.GetTrees)
			trees.POST("", h.CreateTree)
			trees.GET("/:id", h.GetTree)
			trees.PUT("/:id", h.UpdateTree)
			trees.DELETE("/:id", h.DeleteTree)
			trees.GET("/:id/nodes", h.GetTreeNodes)
		}

		// Node routes
		nodes := v1.Group("/nodes")
		{
			nodes.POST("", h.CreateNode)
			nodes.GET("/:id", h.GetNode)
			nodes.PUT("/:id", h.UpdateNode)
			nodes.DELETE("/:id", h.DeleteNode)
			nodes.POST("/:id/analyze", h.AnalyzeNode)
		}

		// Search routes
		search := v1.Group("/search")
		{
			search.GET("/trees", h.SearchTrees)
			search.POST("/semantic", h.SemanticSearch)
		}

		// AI routes
		aiRoutes := v1.Group("/ai")
		{
			aiRoutes.POST("/analyze", h.AnalyzeContent)
			aiRoutes.POST("/suggest-connections", h.SuggestConnections)
			aiRoutes.GET("/health", h.AIHealthCheck)
		}
	}
}

// === Health Check Handlers ===

// HealthCheck checks the overall system health
func (h *Handler) HealthCheck(c *gin.Context) {
	health := map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	// Check database connection
	if err := h.db.Health(); err != nil {
		health["database"] = "unhealthy: " + err.Error()
		health["status"] = "degraded"
	} else {
		health["database"] = "healthy"
	}

	// Check LocalAI connection
	if h.localAI.IsAvailable() {
		health["localai"] = "healthy"
	} else {
		health["localai"] = "unavailable (using fallback mode)"
	}

	statusCode := http.StatusOK
	if health["status"] == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, models.NewSuccessResponse(health))
}

// AIHealthCheck specifically checks LocalAI health
func (h *Handler) AIHealthCheck(c *gin.Context) {
	if err := h.localAI.Health(); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.NewErrorResponse(
			"AI_UNAVAILABLE",
			"LocalAI service is not available",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(map[string]string{
		"status":  "healthy",
		"service": "LocalAI",
	}))
}

// === Tree Handlers ===

// CreateTree creates a new knowledge tree
func (h *Handler) CreateTree(c *gin.Context) {
	var req models.CreateTreeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// For now, use a default owner ID - in production, get from auth
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		ownerID = "default-user"
	}

	tree := models.NewTree(req.Title, req.Description, ownerID, req.IsPublic)

	if err := h.db.CreateTree(tree); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to create tree",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(tree))
}

// GetTrees retrieves trees with optional filtering
func (h *Handler) GetTrees(c *gin.Context) {
	req := &models.TreeSearchRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_QUERY",
			"Invalid query parameters",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"VALIDATION_ERROR",
			"Invalid search parameters",
			err.Error(),
		))
		return
	}

	trees, total, err := h.db.SearchTrees(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to search trees",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(trees, total, req.Limit, req.Offset))
}

// GetTree retrieves a specific tree with its nodes
func (h *Handler) GetTree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Tree ID is required",
			"",
		))
		return
	}

	tree, err := h.db.GetTreeWithNodes(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to retrieve tree",
			err.Error(),
		))
		return
	}

	if tree == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"TREE_NOT_FOUND",
			"Tree not found",
			"",
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(tree))
}

// UpdateTree updates an existing tree
func (h *Handler) UpdateTree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Tree ID is required",
			"",
		))
		return
	}

	var req models.UpdateTreeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Check if tree exists
	existingTree, err := h.db.GetTree(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to check tree existence",
			err.Error(),
		))
		return
	}

	if existingTree == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"TREE_NOT_FOUND",
			"Tree not found",
			"",
		))
		return
	}

	if err := h.db.UpdateTree(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to update tree",
			err.Error(),
		))
		return
	}

	// Return updated tree
	updatedTree, _ := h.db.GetTree(id)
	c.JSON(http.StatusOK, models.NewSuccessResponse(updatedTree))
}

// DeleteTree deletes a tree and all its nodes
func (h *Handler) DeleteTree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Tree ID is required",
			"",
		))
		return
	}

	// Check if tree exists
	existingTree, err := h.db.GetTree(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to check tree existence",
			err.Error(),
		))
		return
	}

	if existingTree == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"TREE_NOT_FOUND",
			"Tree not found",
			"",
		))
		return
	}

	if err := h.db.DeleteTree(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to delete tree",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(map[string]string{
		"message": "Tree deleted successfully",
		"id":      id,
	}))
}

// GetTreeNodes retrieves all nodes for a specific tree
func (h *Handler) GetTreeNodes(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Tree ID is required",
			"",
		))
		return
	}

	nodes, err := h.db.GetNodesByTreeID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to retrieve nodes",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(nodes))
}

// === Node Handlers ===

// CreateNode creates a new node in a tree
func (h *Handler) CreateNode(c *gin.Context) {
	var req models.CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Validate tree exists
	tree, err := h.db.GetTree(req.TreeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to validate tree",
			err.Error(),
		))
		return
	}

	if tree == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"TREE_NOT_FOUND",
			"Tree not found",
			"",
		))
		return
	}

	// Calculate depth
	depth := 0
	if req.ParentID != nil {
		parentNode, err := h.db.GetNode(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
				"DATABASE_ERROR",
				"Failed to validate parent node",
				err.Error(),
			))
			return
		}
		if parentNode == nil {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(
				"PARENT_NOT_FOUND",
				"Parent node not found",
				"",
			))
			return
		}
		depth = parentNode.Depth + 1
	}

	node := models.NewNode(req.TreeID, req.ParentID, req.Title, req.Content, depth, req.Position)

	// Create the node
	if err := h.db.CreateNode(node); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to create node",
			err.Error(),
		))
		return
	}

	// Asynchronously analyze the node with AI
	go h.analyzeNodeAsync(node.ID, node.Title, node.Content)

	c.JSON(http.StatusCreated, models.NewSuccessResponse(node))
}

// GetNode retrieves a specific node
func (h *Handler) GetNode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Node ID is required",
			"",
		))
		return
	}

	node, err := h.db.GetNode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to retrieve node",
			err.Error(),
		))
		return
	}

	if node == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"NODE_NOT_FOUND",
			"Node not found",
			"",
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(node))
}

// UpdateNode updates an existing node
func (h *Handler) UpdateNode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Node ID is required",
			"",
		))
		return
	}

	var req models.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Check if node exists
	existingNode, err := h.db.GetNode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to check node existence",
			err.Error(),
		))
		return
	}

	if existingNode == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"NODE_NOT_FOUND",
			"Node not found",
			"",
		))
		return
	}

	if err := h.db.UpdateNode(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to update node",
			err.Error(),
		))
		return
	}

	// If content was updated, re-analyze with AI
	if req.Title != nil || req.Content != nil {
		updatedNode, _ := h.db.GetNode(id)
		if updatedNode != nil {
			go h.analyzeNodeAsync(updatedNode.ID, updatedNode.Title, updatedNode.Content)
		}
	}

	// Return updated node
	updatedNode, _ := h.db.GetNode(id)
	c.JSON(http.StatusOK, models.NewSuccessResponse(updatedNode))
}

// DeleteNode deletes a node and all its children
func (h *Handler) DeleteNode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Node ID is required",
			"",
		))
		return
	}

	// Check if node exists
	existingNode, err := h.db.GetNode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to check node existence",
			err.Error(),
		))
		return
	}

	if existingNode == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"NODE_NOT_FOUND",
			"Node not found",
			"",
		))
		return
	}

	if err := h.db.DeleteNode(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to delete node",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(map[string]string{
		"message": "Node deleted successfully",
		"id":      id,
	}))
}

// AnalyzeNode manually triggers AI analysis for a node
func (h *Handler) AnalyzeNode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"MISSING_ID",
			"Node ID is required",
			"",
		))
		return
	}

	node, err := h.db.GetNode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to retrieve node",
			err.Error(),
		))
		return
	}

	if node == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"NODE_NOT_FOUND",
			"Node not found",
			"",
		))
		return
	}

	// Perform AI analysis
	analysis, err := h.performNodeAnalysis(node.ID, node.Title, node.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"AI_ERROR",
			"Failed to analyze node",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(analysis))
}

// === Search Handlers ===

// SearchTrees performs text-based search on trees
func (h *Handler) SearchTrees(c *gin.Context) {
	req := &models.TreeSearchRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_QUERY",
			"Invalid query parameters",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"VALIDATION_ERROR",
			"Invalid search parameters",
			err.Error(),
		))
		return
	}

	trees, total, err := h.db.SearchTrees(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to search trees",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(trees, total, req.Limit, req.Offset))
}

// SemanticSearch performs vector-based semantic search
func (h *Handler) SemanticSearch(c *gin.Context) {
	var req models.SemanticSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"VALIDATION_ERROR",
			"Invalid search parameters",
			err.Error(),
		))
		return
	}

	// Generate embedding for the query
	queryEmbedding, err := h.localAI.GenerateEmbedding(req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"AI_ERROR",
			"Failed to generate query embedding",
			err.Error(),
		))
		return
	}

	// Convert to bytes for database storage
	embeddingBytes, err := json.Marshal(queryEmbedding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"ENCODING_ERROR",
			"Failed to encode embedding",
			err.Error(),
		))
		return
	}

	// Perform semantic search
	results, err := h.db.SemanticSearch(embeddingBytes, req.TreeID, req.Threshold, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to perform semantic search",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(results))
}

// === AI Handlers ===

// AnalyzeContent performs AI analysis on arbitrary content
func (h *Handler) AnalyzeContent(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	analysis, err := h.localAI.AnalyzeContent(req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"AI_ERROR",
			"Failed to analyze content",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(analysis))
}

// SuggestConnections suggests connections between nodes
func (h *Handler) SuggestConnections(c *gin.Context) {
	var req struct {
		NodeID string  `json:"node_id" binding:"required"`
		TreeID *string `json:"tree_id,omitempty"`
		Limit  int     `json:"limit,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	// Get the source node
	sourceNode, err := h.db.GetNode(req.NodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to retrieve source node",
			err.Error(),
		))
		return
	}

	if sourceNode == nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"NODE_NOT_FOUND",
			"Source node not found",
			"",
		))
		return
	}

	// Get candidate nodes (from same tree or all trees)
	var candidateNodes []models.Node
	if req.TreeID != nil {
		candidateNodes, err = h.db.GetNodesByTreeID(*req.TreeID)
	} else {
		candidateNodes, err = h.db.GetNodesByTreeID(sourceNode.TreeID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"DATABASE_ERROR",
			"Failed to retrieve candidate nodes",
			err.Error(),
		))
		return
	}

	// Filter out the source node itself
	filteredCandidates := []models.Node{}
	for _, node := range candidateNodes {
		if node.ID != sourceNode.ID {
			filteredCandidates = append(filteredCandidates, node)
		}
	}

	// Limit candidates to avoid token limits
	if len(filteredCandidates) > req.Limit {
		filteredCandidates = filteredCandidates[:req.Limit]
	}

	// Get AI suggestions
	suggestions, err := h.localAI.SuggestConnections(sourceNode, filteredCandidates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"AI_ERROR",
			"Failed to generate connection suggestions",
			err.Error(),
		))
		return
	}

	result := map[string]interface{}{
		"source_node":     sourceNode,
		"suggestions":     suggestions,
		"candidate_count": len(filteredCandidates),
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(result))
}

// === Helper Functions ===

// analyzeNodeAsync performs AI analysis on a node asynchronously
func (h *Handler) analyzeNodeAsync(nodeID, title, content string) {
	analysis, err := h.performNodeAnalysis(nodeID, title, content)
	if err != nil {
		// Log error but don't fail the request
		return
	}

	// Update the node with AI analysis results
	conceptsJSON, _ := json.Marshal(analysis.Concepts)
	suggestions := ""
	if len(analysis.Suggestions) > 0 {
		suggestionsJSON, _ := json.Marshal(analysis.Suggestions)
		suggestions = string(suggestionsJSON)
	}

	h.db.UpdateNodeAI(nodeID, nil, conceptsJSON, &analysis.Difficulty, &suggestions)
}

// performNodeAnalysis performs comprehensive AI analysis on a node
func (h *Handler) performNodeAnalysis(nodeID, title, content string) (*models.AIAnalysis, error) {
	// Generate embedding
	embedding, err := h.localAI.GenerateEmbedding(title + " " + content)
	if err != nil {
		return nil, err
	}

	// Store embedding in database
	embeddingBytes, _ := json.Marshal(embedding)
	h.db.UpdateNodeAI(nodeID, embeddingBytes, nil, nil, nil)

	// Perform content analysis
	analysis, err := h.localAI.AnalyzeContent(title, content)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}
