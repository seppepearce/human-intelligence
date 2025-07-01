package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"human-intelligence/internal/database"
	"human-intelligence/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TreeHandler struct {
	neo4j *database.Neo4jService
}

// CreateTreeRequest represents the request to create a new learning tree
type CreateTreeRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	IsTemplate  bool   `json:"is_template"`
	RootNode    struct {
		Title       string                 `json:"title" binding:"required,max=255"`
		Content     string                 `json:"content"`
		ContentType string                 `json:"content_type"`
		Description string                 `json:"description"`
		Metadata    map[string]interface{} `json:"metadata"`
	} `json:"root_node" binding:"required"`
}

// TreeNodeRequest represents the request to add a child node to an existing node
type TreeNodeRequest struct {
	Title        string                 `json:"title" binding:"required,max=255"`
	Content      string                 `json:"content"`
	ContentType  string                 `json:"content_type"`
	Description  string                 `json:"description"`
	ParentNodeID string                 `json:"parent_node_id" binding:"required"`
	Position     int                    `json:"position"` // Order among siblings
	Metadata     map[string]interface{} `json:"metadata"`
	Tags         []string               `json:"tags"`
}

// TreeStructureResponse represents the hierarchical tree structure
type TreeStructureResponse struct {
	Tree     *models.Tree      `json:"tree"`
	RootNode *NodeWithChildren `json:"root_node"`
	Stats    TreeStats         `json:"stats"`
}

// NodeWithChildren represents a node with its children for hierarchical display
type NodeWithChildren struct {
	*models.Node
	Children []*NodeWithChildren `json:"children"`
}

// TreeStats provides quick statistics about the tree
type TreeStats struct {
	TotalNodes  int `json:"total_nodes"`
	MaxDepth    int `json:"max_depth"`
	BranchCount int `json:"branch_count"`
}

func NewTreeHandler(neo4j *database.Neo4jService) *TreeHandler {
	return &TreeHandler{neo4j: neo4j}
}

// CreateTree creates a new learning tree with an initial root node
func (h *TreeHandler) CreateTree(c *gin.Context) {
	var req CreateTreeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// For now, we'll use a hardcoded user ID - in production this would come from auth
	userID := "user-alice-2024"

	// Create the tree
	tree := &models.Tree{
		ID:          generateTreeID(),
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		IsTemplate:  req.IsTemplate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerID:     userID,
		NodesCount:  1, // Will have one root node
	}

	// Create the root node
	rootNode := &models.Node{
		ID:          generateNodeID(),
		Title:       req.RootNode.Title,
		Content:     req.RootNode.Content,
		ContentType: getContentType(req.RootNode.ContentType),
		Description: req.RootNode.Description,
		IsPublic:    req.IsPublic,
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerID:     userID,
		Level:       0, // Root level
		Position:    0, // First position
		Metadata:    req.RootNode.Metadata,
	}

	ctx := context.Background()

	// Execute the tree creation with root node in a transaction-like manner
	if err := h.createTreeWithRootNode(ctx, tree, rootNode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create tree with root node",
			"details": err.Error(),
		})
		return
	}

	// Return the created tree with root node
	response := &TreeStructureResponse{
		Tree: tree,
		RootNode: &NodeWithChildren{
			Node:     rootNode,
			Children: []*NodeWithChildren{},
		},
		Stats: TreeStats{
			TotalNodes:  1,
			MaxDepth:    1,
			BranchCount: 0,
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetTreeStructure returns the complete hierarchical structure of a tree
func (h *TreeHandler) GetTreeStructure(c *gin.Context) {
	treeID := c.Param("treeId")
	if treeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tree ID is required"})
		return
	}

	ctx := context.Background()
	structure, err := h.getCompleteTreeStructure(ctx, treeID)
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tree not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get tree structure",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, structure)
}

// AddNodeToTree adds a new child node to an existing node in the tree
func (h *TreeHandler) AddNodeToTree(c *gin.Context) {
	treeID := c.Param("treeId")
	if treeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tree ID is required"})
		return
	}

	var req TreeNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// For now, we'll use a hardcoded user ID
	userID := "user-alice-2024"

	ctx := context.Background()

	// First, verify the parent node exists and belongs to this tree
	parentNode, err := h.neo4j.GetNodeByID(ctx, req.ParentNodeID)
	if err != nil {
		if err == models.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent node not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify parent node"})
		return
	}

	// Create the new child node
	childNode := &models.Node{
		ID:           generateNodeID(),
		Title:        req.Title,
		Content:      req.Content,
		ContentType:  getContentType(req.ContentType),
		Description:  req.Description,
		IsPublic:     parentNode.IsPublic, // Inherit from parent
		IsPublished:  true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		OwnerID:      userID,
		ParentNodeID: &req.ParentNodeID,
		Level:        parentNode.Level + 1,
		Position:     req.Position,
		Metadata:     req.Metadata,
		Tags:         req.Tags,
	}

	// Add the node to the tree with proper relationships
	if err := h.addChildNodeToTree(ctx, treeID, childNode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add node to tree",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Node added successfully",
		"node":    childNode,
	})
}

// GetTreeNodes returns all nodes in a tree with their relationships
func (h *TreeHandler) GetTreeNodes(c *gin.Context) {
	treeID := c.Param("treeId")
	if treeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tree ID is required"})
		return
	}

	ctx := context.Background()
	nodes, err := h.getAllTreeNodes(ctx, treeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get tree nodes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tree_id": treeID,
		"nodes":   nodes,
		"count":   len(nodes),
	})
}

// UpdateTreeMetadata updates tree information (name, description, etc.)
func (h *TreeHandler) UpdateTreeMetadata(c *gin.Context) {
	treeID := c.Param("treeId")
	if treeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tree ID is required"})
		return
	}

	var req struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		IsPublic    *bool  `json:"is_public,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	ctx := context.Background()
	if err := h.updateTreeMetadata(ctx, treeID, req.Name, req.Description, req.IsPublic); err != nil {
		if err == models.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tree not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update tree",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tree updated successfully"})
}

// DeleteTree removes a tree and all its nodes
func (h *TreeHandler) DeleteTree(c *gin.Context) {
	treeID := c.Param("treeId")
	if treeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tree ID is required"})
		return
	}

	ctx := context.Background()
	if err := h.deleteTreeAndNodes(ctx, treeID); err != nil {
		if err == models.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tree not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete tree",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tree deleted successfully"})
}

// Helper methods

func (h *TreeHandler) createTreeWithRootNode(ctx context.Context, tree *models.Tree, rootNode *models.Node) error {
	// Create tree
	if err := h.neo4j.CreateTree(ctx, tree); err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}

	// Create root node
	if err := h.neo4j.CreateNode(ctx, rootNode); err != nil {
		return fmt.Errorf("failed to create root node: %w", err)
	}

	// Connect tree to root node
	if err := h.neo4j.AddNodeToTree(ctx, tree.ID, rootNode.ID, 0); err != nil {
		return fmt.Errorf("failed to connect root node to tree: %w", err)
	}

	return nil
}

func (h *TreeHandler) getCompleteTreeStructure(ctx context.Context, treeID string) (*TreeStructureResponse, error) {
	// Get tree information
	tree, err := h.getTreeByID(ctx, treeID)
	if err != nil {
		return nil, err
	}

	// Get all nodes in the tree
	nodes, err := h.getAllTreeNodes(ctx, treeID)
	if err != nil {
		return nil, err
	}

	// Build hierarchical structure
	rootNode, stats := h.buildHierarchy(nodes)

	return &TreeStructureResponse{
		Tree:     tree,
		RootNode: rootNode,
		Stats:    stats,
	}, nil
}

func (h *TreeHandler) addChildNodeToTree(ctx context.Context, treeID string, childNode *models.Node) error {
	// Create the node
	if err := h.neo4j.CreateNode(ctx, childNode); err != nil {
		return fmt.Errorf("failed to create child node: %w", err)
	}

	// Add node to tree
	if err := h.neo4j.AddNodeToTree(ctx, treeID, childNode.ID, childNode.Position); err != nil {
		return fmt.Errorf("failed to add node to tree: %w", err)
	}

	// Create parent-child relationship
	if childNode.ParentNodeID != nil {
		if err := h.createParentChildRelationship(ctx, *childNode.ParentNodeID, childNode.ID); err != nil {
			return fmt.Errorf("failed to create parent-child relationship: %w", err)
		}
	}

	// Update tree's node count
	if err := h.incrementTreeNodeCount(ctx, treeID); err != nil {
		return fmt.Errorf("failed to update tree node count: %w", err)
	}

	return nil
}

func (h *TreeHandler) buildHierarchy(nodes []*models.Node) (*NodeWithChildren, TreeStats) {
	nodeMap := make(map[string]*NodeWithChildren)
	var rootNode *NodeWithChildren
	stats := TreeStats{}

	// Create NodeWithChildren for each node
	for _, node := range nodes {
		nodeWithChildren := &NodeWithChildren{
			Node:     node,
			Children: []*NodeWithChildren{},
		}
		nodeMap[node.ID] = nodeWithChildren

		if node.ParentNodeID == nil {
			rootNode = nodeWithChildren
		}
	}

	// Build parent-child relationships
	branchCount := 0
	maxDepth := 0

	for _, node := range nodes {
		if node.ParentNodeID != nil {
			if parent, exists := nodeMap[*node.ParentNodeID]; exists {
				parent.Children = append(parent.Children, nodeMap[node.ID])
				if len(parent.Children) > 1 {
					branchCount++
				}
			}
		}
		if node.Level > maxDepth {
			maxDepth = node.Level
		}
	}

	stats.TotalNodes = len(nodes)
	stats.MaxDepth = maxDepth + 1 // Convert 0-based to 1-based
	stats.BranchCount = branchCount

	return rootNode, stats
}

// Neo4j query helpers

func (h *TreeHandler) getTreeByID(ctx context.Context, treeID string) (*models.Tree, error) {
	query := `
		MATCH (t:Tree {id: $tree_id})
		OPTIONAL MATCH (t)<-[:CREATED]-(owner:User)
		RETURN t, owner
	`

	params := map[string]interface{}{
		"tree_id": treeID,
	}

	records, err := h.neo4j.ExecuteQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, models.ErrRecordNotFound
	}

	record := records[0]
	tree := &models.Tree{}

	if treeData, ok := record.Get("t"); ok {
		treeValue := treeData.(neo4j.Node)
		if err := tree.FromNode(treeValue); err != nil {
			return nil, err
		}
	}

	return tree, nil
}

func (h *TreeHandler) getAllTreeNodes(ctx context.Context, treeID string) ([]*models.Node, error) {
	query := `
		MATCH (t:Tree {id: $tree_id})-[:CONTAINS]->(n:Node)
		OPTIONAL MATCH (n)<-[:CREATED]-(owner:User)
		OPTIONAL MATCH (n)-[:TAGGED]->(tag:Tag)
		RETURN n, owner, collect(tag.name) as tags
		ORDER BY n.level, n.position
	`

	params := map[string]interface{}{
		"tree_id": treeID,
	}

	records, err := h.neo4j.ExecuteQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	var nodes []*models.Node
	for _, record := range records {
		node := &models.Node{}

		if nodeData, ok := record.Get("n"); ok {
			nodeValue := nodeData.(neo4j.Node)
			if err := node.FromNode(nodeValue); err != nil {
				continue
			}
		}

		// Get tags
		if tagsData, ok := record.Get("tags"); ok {
			if tags, ok := tagsData.([]interface{}); ok {
				node.Tags = make([]string, len(tags))
				for i, tag := range tags {
					if tagStr, ok := tag.(string); ok {
						node.Tags[i] = tagStr
					}
				}
			}
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (h *TreeHandler) createParentChildRelationship(ctx context.Context, parentID, childID string) error {
	query := `
		MATCH (parent:Node {id: $parent_id})
		MATCH (child:Node {id: $child_id})
		CREATE (parent)-[:PARENT_OF]->(child)
		CREATE (child)-[:CHILD_OF]->(parent)
		RETURN parent, child
	`

	params := map[string]interface{}{
		"parent_id": parentID,
		"child_id":  childID,
	}

	_, err := h.neo4j.ExecuteWrite(ctx, query, params)
	return err
}

func (h *TreeHandler) incrementTreeNodeCount(ctx context.Context, treeID string) error {
	query := `
		MATCH (t:Tree {id: $tree_id})
		SET t.nodes_count = t.nodes_count + 1,
		    t.updated_at = datetime()
		RETURN t
	`

	params := map[string]interface{}{
		"tree_id": treeID,
	}

	_, err := h.neo4j.ExecuteWrite(ctx, query, params)
	return err
}

func (h *TreeHandler) updateTreeMetadata(ctx context.Context, treeID, name, description string, isPublic *bool) error {
	setParts := []string{"updated_at = datetime()"}
	params := map[string]interface{}{
		"tree_id": treeID,
	}

	if name != "" {
		setParts = append(setParts, "name = $name")
		params["name"] = name
	}

	if description != "" {
		setParts = append(setParts, "description = $description")
		params["description"] = description
	}

	if isPublic != nil {
		setParts = append(setParts, "is_public = $is_public")
		params["is_public"] = *isPublic
	}

	query := fmt.Sprintf(`
		MATCH (t:Tree {id: $tree_id})
		SET t.%s
		RETURN t
	`, strings.Join(setParts, ", "))

	_, err := h.neo4j.ExecuteWrite(ctx, query, params)
	return err
}

func (h *TreeHandler) deleteTreeAndNodes(ctx context.Context, treeID string) error {
	query := `
		MATCH (t:Tree {id: $tree_id})
		OPTIONAL MATCH (t)-[:CONTAINS]->(n:Node)
		DETACH DELETE t, n
	`

	params := map[string]interface{}{
		"tree_id": treeID,
	}

	_, err := h.neo4j.ExecuteWrite(ctx, query, params)
	return err
}

// Utility functions

func generateTreeID() string {
	return fmt.Sprintf("tree-%d", time.Now().UnixNano())
}

func generateNodeID() string {
	return fmt.Sprintf("node-%d", time.Now().UnixNano())
}

func getContentType(contentType string) string {
	if contentType == "" {
		return "text"
	}
	validTypes := []string{"text", "markdown", "html", "video", "image", "code", "quiz"}
	for _, valid := range validTypes {
		if contentType == valid {
			return contentType
		}
	}
	return "text"
}

func getIntParam(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
