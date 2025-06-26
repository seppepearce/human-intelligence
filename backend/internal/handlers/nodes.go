package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"human-intelligence/internal/database"
	"human-intelligence/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type NodeHandler struct {
	db *database.DB
}

type CreateNodeRequest struct {
	Title       string                 `json:"title" binding:"required,max=255"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	NodeType    models.NodeType        `json:"node_type" binding:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
	Tags        []string               `json:"tags"`
	IsPublic    bool                   `json:"is_public"`
}

type UpdateNodeRequest struct {
	Title       string                 `json:"title,omitempty" binding:"omitempty,max=255"`
	Description string                 `json:"description,omitempty"`
	Content     string                 `json:"content,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	IsPublic    *bool                  `json:"is_public,omitempty"`
}

type VoteRequest struct {
	VoteType int `json:"vote_type" binding:"required,oneof=-1 1"`
}

func NewNodeHandler(db *database.DB) *NodeHandler {
	return &NodeHandler{db: db}
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Validate node type
	if !h.isValidNodeType(req.NodeType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid node type",
		})
		return
	}

	node := models.Node{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		NodeType:    req.NodeType,
		Tags:        req.Tags,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPublic:    req.IsPublic,
		VoteScore:   0,
		ViewCount:   0,
	}

	// Convert metadata to JSON
	metadataJSON, err := h.convertToJSON(req.Metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid metadata format",
		})
		return
	}

	query := `
		INSERT INTO nodes (id, title, description, content, node_type, metadata, tags, created_by, created_at, updated_at, is_public, vote_score, view_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, title, description, content, node_type, metadata, tags, created_by, created_at, updated_at, is_public, vote_score, view_count
	`

	err = h.db.QueryRow(query, node.ID, node.Title, node.Description, node.Content,
		node.NodeType, metadataJSON, pq.Array(node.Tags), node.CreatedBy, node.CreatedAt,
		node.UpdatedAt, node.IsPublic, node.VoteScore, node.ViewCount).Scan(
		&node.ID, &node.Title, &node.Description, &node.Content, &node.NodeType,
		&node.Metadata, pq.Array(&node.Tags), &node.CreatedBy, &node.CreatedAt,
		&node.UpdatedAt, &node.IsPublic, &node.VoteScore, &node.ViewCount)

	if err != nil {
		log.Printf("Failed to create node - DB error: %v", err)
		log.Printf("Query: %s", query)
		log.Printf("Values: ID=%v, Title=%s, NodeType=%s, CreatedBy=%v", node.ID, node.Title, node.NodeType, node.CreatedBy)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create node",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"node": node,
	})
}

func (h *NodeHandler) GetNode(c *gin.Context) {
	nodeIDStr := c.Param("id")
	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid node ID",
		})
		return
	}

	node, err := h.getNodeByID(nodeID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Node not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Check if user has access to private nodes
	userID := h.getUserIDFromContext(c)
	if !node.IsPublic && node.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	// Increment view count
	go h.incrementViewCount(nodeID)

	c.JSON(http.StatusOK, gin.H{
		"node": node,
	})
}

func (h *NodeHandler) UpdateNode(c *gin.Context) {
	nodeIDStr := c.Param("id")
	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid node ID",
		})
		return
	}

	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Check if user owns the node
	node, err := h.getNodeByID(nodeID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Node not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	if node.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only update your own nodes",
		})
		return
	}

	var req UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Title != "" {
		setParts = append(setParts, "title = $"+strconv.Itoa(argIndex))
		args = append(args, req.Title)
		argIndex++
	}

	if req.Description != "" {
		setParts = append(setParts, "description = $"+strconv.Itoa(argIndex))
		args = append(args, req.Description)
		argIndex++
	}

	if req.Content != "" {
		setParts = append(setParts, "content = $"+strconv.Itoa(argIndex))
		args = append(args, req.Content)
		argIndex++
	}

	if req.Metadata != nil {
		metadataJSON, err := h.convertToJSON(req.Metadata)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid metadata format",
			})
			return
		}
		setParts = append(setParts, "metadata = $"+strconv.Itoa(argIndex))
		args = append(args, metadataJSON)
		argIndex++
	}

	if req.Tags != nil {
		setParts = append(setParts, "tags = $"+strconv.Itoa(argIndex))
		args = append(args, pq.Array(req.Tags))
		argIndex++
	}

	if req.IsPublic != nil {
		setParts = append(setParts, "is_public = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsPublic)
		argIndex++
	}

	if len(setParts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields to update",
		})
		return
	}

	// Add updated_at
	setParts = append(setParts, "updated_at = $"+strconv.Itoa(argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, nodeID)

	query := "UPDATE nodes SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)

	_, err = h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update node",
		})
		return
	}

	// Return updated node
	updatedNode, err := h.getNodeByID(nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated node",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"node": updatedNode,
	})
}

func (h *NodeHandler) DeleteNode(c *gin.Context) {
	nodeIDStr := c.Param("id")
	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid node ID",
		})
		return
	}

	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Check if user owns the node
	node, err := h.getNodeByID(nodeID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Node not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	if node.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own nodes",
		})
		return
	}

	// Delete the node
	query := "DELETE FROM nodes WHERE id = $1"
	_, err = h.db.Exec(query, nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete node",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Node deleted successfully",
	})
}

func (h *NodeHandler) VoteNode(c *gin.Context) {
	nodeIDStr := c.Param("id")
	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid node ID",
		})
		return
	}

	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Check if node exists
	_, err = h.getNodeByID(nodeID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Node not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Insert or update vote
	query := `
		INSERT INTO votes (id, user_id, target_id, target_type, vote_type, created_at)
		VALUES ($1, $2, $3, 'node', $4, $5)
		ON CONFLICT (user_id, target_id, target_type)
		DO UPDATE SET vote_type = $4, created_at = $5
	`

	_, err = h.db.Exec(query, uuid.New(), userID, nodeID, req.VoteType, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to record vote",
		})
		return
	}

	// Get updated node with new vote score
	updatedNode, err := h.getNodeByID(nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated node",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"node": updatedNode,
	})
}

func (h *NodeHandler) GetPublicNodes(c *gin.Context) {
	// Parse query parameters
	limit := 20
	offset := 0
	sortBy := "created_at"

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	if s := c.Query("sort"); s != "" {
		if s == "vote_score" || s == "view_count" || s == "updated_at" {
			sortBy = s
		}
	}

	query := `
		SELECT id, title, description, content, node_type, metadata, tags, created_by, created_at, updated_at, is_public, vote_score, view_count
		FROM nodes
		WHERE is_public = true
		ORDER BY ` + sortBy + ` DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := h.db.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch nodes",
		})
		return
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(&node.ID, &node.Title, &node.Description, &node.Content,
			&node.NodeType, &node.Metadata, pq.Array(&node.Tags), &node.CreatedBy,
			&node.CreatedAt, &node.UpdatedAt, &node.IsPublic, &node.VoteScore, &node.ViewCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan node",
			})
			return
		}
		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes":  nodes,
		"limit":  limit,
		"offset": offset,
		"count":  len(nodes),
	})
}

// Helper methods
func (h *NodeHandler) getNodeByID(nodeID uuid.UUID) (*models.Node, error) {
	var node models.Node
	query := `
		SELECT id, title, description, content, node_type, metadata, tags, created_by, created_at, updated_at, is_public, vote_score, view_count
		FROM nodes
		WHERE id = $1
	`

	err := h.db.QueryRow(query, nodeID).Scan(
		&node.ID, &node.Title, &node.Description, &node.Content, &node.NodeType,
		&node.Metadata, pq.Array(&node.Tags), &node.CreatedBy, &node.CreatedAt,
		&node.UpdatedAt, &node.IsPublic, &node.VoteScore, &node.ViewCount)

	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (h *NodeHandler) getUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	if id, ok := userID.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}

func (h *NodeHandler) isValidNodeType(nodeType models.NodeType) bool {
	validTypes := []models.NodeType{
		models.NodeTypeText,
		models.NodeTypeVideo,
		models.NodeTypeLink,
		models.NodeTypeQuiz,
		models.NodeTypeCode,
		models.NodeTypeLatex,
		models.NodeTypeJupyter,
		models.NodeTypeWolfram,
		models.NodeTypeCustom,
	}

	for _, validType := range validTypes {
		if nodeType == validType {
			return true
		}
	}
	return false
}

func (h *NodeHandler) convertToJSON(data interface{}) ([]byte, error) {
	if data == nil {
		return []byte("{}"), nil
	}

	return json.Marshal(data)
}

func (h *NodeHandler) incrementViewCount(nodeID uuid.UUID) {
	query := "UPDATE nodes SET view_count = view_count + 1 WHERE id = $1"
	h.db.Exec(query, nodeID)
}
