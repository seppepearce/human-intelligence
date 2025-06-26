package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"human-intelligence/internal/database"
	"human-intelligence/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PathHandler struct {
	db *database.DB
}

type CreatePathRequest struct {
	Name        string             `json:"name" binding:"required,max=255"`
	Description string             `json:"description"`
	AccessLevel models.AccessLevel `json:"access_level"`
	Tags        []string           `json:"tags"`
	NodeIDs     []uuid.UUID        `json:"node_ids"`
}

type UpdatePathRequest struct {
	Name        string             `json:"name,omitempty" binding:"omitempty,max=255"`
	Description string             `json:"description,omitempty"`
	AccessLevel models.AccessLevel `json:"access_level,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
}

type AddNodeRequest struct {
	NodeID        uuid.UUID   `json:"node_id" binding:"required"`
	Position      int         `json:"position"`
	IsRequired    bool        `json:"is_required"`
	Prerequisites []uuid.UUID `json:"prerequisites"`
}

type CreateTLDRRequest struct {
	Content string `json:"content" binding:"required,max=140"`
}

func NewPathHandler(db *database.DB) *PathHandler {
	return &PathHandler{db: db}
}

func (h *PathHandler) CreatePath(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req CreatePathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Validate access level
	if req.AccessLevel == "" {
		req.AccessLevel = models.AccessPublic
	}

	path := models.LearningPath{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		AccessLevel:     req.AccessLevel,
		CreatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		VoteScore:       0,
		ForkCount:       0,
		CompletionCount: 0,
		Tags:            req.Tags,
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start transaction",
		})
		return
	}
	defer tx.Rollback()

	// Insert path
	query := `
		INSERT INTO learning_paths (id, name, description, access_level, created_by, created_at, updated_at, vote_score, fork_count, completion_count, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(query, path.ID, path.Name, path.Description, path.AccessLevel,
		path.CreatedBy, path.CreatedAt, path.UpdatedAt, path.VoteScore,
		path.ForkCount, path.CompletionCount, path.Tags)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create path",
		})
		return
	}

	// Add nodes to path if provided
	if len(req.NodeIDs) > 0 {
		for i, nodeID := range req.NodeIDs {
			nodeQuery := `
				INSERT INTO path_nodes (id, path_id, node_id, position, is_required)
				VALUES ($1, $2, $3, $4, $5)
			`
			_, err = tx.Exec(nodeQuery, uuid.New(), path.ID, nodeID, i+1, true)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to add nodes to path",
				})
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	// Fetch the created path with nodes
	createdPath, err := h.getPathByID(path.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch created path",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"path": createdPath,
	})
}

func (h *PathHandler) GetPath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
		})
		return
	}

	userID := h.getUserIDFromContext(c)
	path, err := h.getPathByID(pathID, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Path not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Check access permissions
	if !h.hasPathAccess(path, userID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": path,
	})
}

func (h *PathHandler) UpdatePath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
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

	// Check if user owns the path
	path, err := h.getPathByID(pathID, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Path not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	if path.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only update your own paths",
		})
		return
	}

	var req UpdatePathRequest
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

	if req.Name != "" {
		setParts = append(setParts, "name = $"+strconv.Itoa(argIndex))
		args = append(args, req.Name)
		argIndex++
	}

	if req.Description != "" {
		setParts = append(setParts, "description = $"+strconv.Itoa(argIndex))
		args = append(args, req.Description)
		argIndex++
	}

	if req.AccessLevel != "" {
		setParts = append(setParts, "access_level = $"+strconv.Itoa(argIndex))
		args = append(args, req.AccessLevel)
		argIndex++
	}

	if req.Tags != nil {
		setParts = append(setParts, "tags = $"+strconv.Itoa(argIndex))
		args = append(args, req.Tags)
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
	args = append(args, pathID)

	query := "UPDATE learning_paths SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)

	_, err = h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update path",
		})
		return
	}

	// Return updated path
	updatedPath, err := h.getPathByID(pathID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated path",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": updatedPath,
	})
}

func (h *PathHandler) DeletePath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
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

	// Check if user owns the path
	path, err := h.getPathByID(pathID, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Path not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	if path.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own paths",
		})
		return
	}

	// Delete the path (cascading deletes will handle related records)
	query := "DELETE FROM learning_paths WHERE id = $1"
	_, err = h.db.Exec(query, pathID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete path",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Path deleted successfully",
	})
}

func (h *PathHandler) ForkPath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
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

	// Get the original path
	originalPath, err := h.getPathByID(pathID, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Path not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Check if user has access to fork this path
	if !h.hasPathAccess(originalPath, userID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	// Create forked path
	forkedPath := models.LearningPath{
		ID:              uuid.New(),
		Name:            "Fork of " + originalPath.Name,
		Description:     originalPath.Description,
		AccessLevel:     models.AccessPublic, // Forks are public by default
		CreatedBy:       userID,
		ParentPathID:    &pathID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		VoteScore:       0,
		ForkCount:       0,
		CompletionCount: 0,
		Tags:            originalPath.Tags,
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start transaction",
		})
		return
	}
	defer tx.Rollback()

	// Insert forked path
	query := `
		INSERT INTO learning_paths (id, name, description, access_level, created_by, parent_path_id, created_at, updated_at, vote_score, fork_count, completion_count, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = tx.Exec(query, forkedPath.ID, forkedPath.Name, forkedPath.Description,
		forkedPath.AccessLevel, forkedPath.CreatedBy, forkedPath.ParentPathID,
		forkedPath.CreatedAt, forkedPath.UpdatedAt, forkedPath.VoteScore,
		forkedPath.ForkCount, forkedPath.CompletionCount, forkedPath.Tags)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create forked path",
		})
		return
	}

	// Copy nodes from original path
	copyNodesQuery := `
		INSERT INTO path_nodes (id, path_id, node_id, position, is_required, prerequisites)
		SELECT uuid_generate_v4(), $1, node_id, position, is_required, prerequisites
		FROM path_nodes
		WHERE path_id = $2
		ORDER BY position
	`

	_, err = tx.Exec(copyNodesQuery, forkedPath.ID, pathID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to copy nodes to forked path",
		})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	// Fetch the created path with nodes
	createdPath, err := h.getPathByID(forkedPath.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch forked path",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"path": createdPath,
	})
}

func (h *PathHandler) AddNodeToPath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
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

	var req AddNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Check if user owns the path
	path, err := h.getPathByID(pathID, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Path not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	if path.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only modify your own paths",
		})
		return
	}

	// If no position specified, add to end
	if req.Position == 0 {
		var maxPosition int
		posQuery := "SELECT COALESCE(MAX(position), 0) FROM path_nodes WHERE path_id = $1"
		h.db.QueryRow(posQuery, pathID).Scan(&maxPosition)
		req.Position = maxPosition + 1
	}

	// Add node to path
	query := `
		INSERT INTO path_nodes (id, path_id, node_id, position, is_required, prerequisites)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = h.db.Exec(query, uuid.New(), pathID, req.NodeID, req.Position, req.IsRequired, req.Prerequisites)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add node to path",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Node added to path successfully",
	})
}

func (h *PathHandler) RemoveNodeFromPath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
		})
		return
	}

	nodeIDStr := c.Param("nodeId")
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

	// Check if user owns the path
	path, err := h.getPathByID(pathID, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Path not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	if path.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only modify your own paths",
		})
		return
	}

	// Remove node from path
	query := "DELETE FROM path_nodes WHERE path_id = $1 AND node_id = $2"
	result, err := h.db.Exec(query, pathID, nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove node from path",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Node not found in path",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Node removed from path successfully",
	})
}

func (h *PathHandler) CompletePath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
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

	// Update or insert user progress
	query := `
		INSERT INTO user_progress (id, user_id, path_id, completed_at, last_active_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, path_id)
		DO UPDATE SET completed_at = $4, last_active_at = $5
	`

	now := time.Now()
	_, err = h.db.Exec(query, uuid.New(), userID, pathID, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to mark path as completed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Path marked as completed",
	})
}

func (h *PathHandler) CreateTLDR(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
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

	var req CreateTLDRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Create TLDR
	tldr := models.TLDR{
		ID:        uuid.New(),
		PathID:    pathID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		VoteScore: 0,
	}

	query := `
		INSERT INTO tldrs (id, path_id, user_id, content, created_at, vote_score)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (path_id, user_id)
		DO UPDATE SET content = $4, created_at = $5
		RETURNING id, path_id, user_id, content, created_at, vote_score
	`

	err = h.db.QueryRow(query, tldr.ID, tldr.PathID, tldr.UserID, tldr.Content,
		tldr.CreatedAt, tldr.VoteScore).Scan(&tldr.ID, &tldr.PathID, &tldr.UserID,
		&tldr.Content, &tldr.CreatedAt, &tldr.VoteScore)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create TLDR",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tldr": tldr,
	})
}

func (h *PathHandler) VotePath(c *gin.Context) {
	pathIDStr := c.Param("id")
	pathID, err := uuid.Parse(pathIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid path ID",
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

	// Insert or update vote
	query := `
		INSERT INTO votes (id, user_id, target_id, target_type, vote_type, created_at)
		VALUES ($1, $2, $3, 'path', $4, $5)
		ON CONFLICT (user_id, target_id, target_type)
		DO UPDATE SET vote_type = $4, created_at = $5
	`

	_, err = h.db.Exec(query, uuid.New(), userID, pathID, req.VoteType, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to record vote",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote recorded successfully",
	})
}

func (h *PathHandler) GetPublicPaths(c *gin.Context) {
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
		if s == "vote_score" || s == "completion_count" || s == "fork_count" || s == "updated_at" {
			sortBy = s
		}
	}

	query := `
		SELECT id, name, description, access_level, created_by, parent_path_id, created_at, updated_at, vote_score, fork_count, completion_count, tags
		FROM learning_paths
		WHERE access_level = 'public'
		ORDER BY ` + sortBy + ` DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := h.db.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch paths",
		})
		return
	}
	defer rows.Close()

	var paths []models.LearningPath
	for rows.Next() {
		var path models.LearningPath
		err := rows.Scan(&path.ID, &path.Name, &path.Description, &path.AccessLevel,
			&path.CreatedBy, &path.ParentPathID, &path.CreatedAt, &path.UpdatedAt,
			&path.VoteScore, &path.ForkCount, &path.CompletionCount, &path.Tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan path",
			})
			return
		}
		paths = append(paths, path)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"paths":  paths,
		"limit":  limit,
		"offset": offset,
		"count":  len(paths),
	})
}

// Helper methods
func (h *PathHandler) getPathByID(pathID uuid.UUID, userID uuid.UUID) (*models.LearningPath, error) {
	var path models.LearningPath
	query := `
		SELECT id, name, description, access_level, created_by, parent_path_id, created_at, updated_at, vote_score, fork_count, completion_count, tags
		FROM learning_paths
		WHERE id = $1
	`

	err := h.db.QueryRow(query, pathID).Scan(
		&path.ID, &path.Name, &path.Description, &path.AccessLevel, &path.CreatedBy,
		&path.ParentPathID, &path.CreatedAt, &path.UpdatedAt, &path.VoteScore,
		&path.ForkCount, &path.CompletionCount, &path.Tags)

	if err != nil {
		return nil, err
	}

	// Load nodes
	nodes, err := h.getPathNodes(pathID)
	if err == nil {
		path.Nodes = nodes
	}

	return &path, nil
}

func (h *PathHandler) getPathNodes(pathID uuid.UUID) ([]models.PathNode, error) {
	query := `
		SELECT pn.id, pn.path_id, pn.node_id, pn.position, pn.is_required, pn.prerequisites,
		       n.id, n.title, n.description, n.content, n.node_type, n.metadata, n.tags, n.created_by, n.created_at, n.updated_at, n.is_public, n.vote_score, n.view_count
		FROM path_nodes pn
		JOIN nodes n ON pn.node_id = n.id
		WHERE pn.path_id = $1
		ORDER BY pn.position
	`

	rows, err := h.db.Query(query, pathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pathNodes []models.PathNode
	for rows.Next() {
		var pathNode models.PathNode
		var node models.Node

		err := rows.Scan(&pathNode.ID, &pathNode.PathID, &pathNode.NodeID,
			&pathNode.Position, &pathNode.IsRequired, &pathNode.Prerequisites,
			&node.ID, &node.Title, &node.Description, &node.Content, &node.NodeType,
			&node.Metadata, &node.Tags, &node.CreatedBy, &node.CreatedAt,
			&node.UpdatedAt, &node.IsPublic, &node.VoteScore, &node.ViewCount)

		if err != nil {
			return nil, err
		}

		pathNode.Node = &node
		pathNodes = append(pathNodes, pathNode)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pathNodes, nil
}

func (h *PathHandler) hasPathAccess(path *models.LearningPath, userID uuid.UUID) bool {
	// Public paths are accessible to everyone
	if path.AccessLevel == models.AccessPublic {
		return true
	}

	// Owner can always access their own paths
	if path.CreatedBy == userID {
		return true
	}

	// For private/group paths, check access control
	if path.AccessLevel == models.AccessPrivate || path.AccessLevel == models.AccessGroup {
		// TODO: Implement access control check
		// For now, deny access to non-owners
		return false
	}

	return false
}

func (h *PathHandler) getUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	if id, ok := userID.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}
