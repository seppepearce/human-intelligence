package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/human-intelligence/ai-knowledge-trees/backend/internal/models"
	_ "modernc.org/sqlite"
)

// DB wraps the database connection and provides methods for data access
type DB struct {
	conn *sql.DB
}

// New creates a new database connection
func New() (*DB, error) {
	var db *sql.DB
	var err error

	// Check if we're using Turso or local SQLite
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_TOKEN")

	if tursoURL != "" && tursoToken != "" {
		// Connect to Turso
		connStr := fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken)
		db, err = sql.Open("sqlite", connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to turso: %w", err)
		}
		log.Println("Connected to Turso database")
	} else {
		// Use local SQLite for development
		dbPath := os.Getenv("SQLITE_PATH")
		if dbPath == "" {
			dbPath = "./ai_knowledge_trees.db"
		}
		db, err = sql.Open("sqlite", dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
		}
		log.Printf("Connected to local SQLite database: %s", dbPath)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	dbInstance := &DB{conn: db}

	// Initialize schema
	if err = dbInstance.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return dbInstance, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// initSchema creates the database tables if they don't exist
func (db *DB) initSchema() error {
	schema := `
	-- Trees: Knowledge collections
	CREATE TABLE IF NOT EXISTS trees (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		owner_id TEXT NOT NULL,
		is_public BOOLEAN DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Nodes: AI-enhanced knowledge units
	CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY,
		tree_id TEXT NOT NULL REFERENCES trees(id) ON DELETE CASCADE,
		parent_id TEXT REFERENCES nodes(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		content TEXT DEFAULT '',

		-- Tree structure
		depth INTEGER DEFAULT 0,
		position INTEGER DEFAULT 0,

		-- AI enhancements
		embedding BLOB,              -- Vector representation
		concepts JSON,               -- AI-extracted concepts
		difficulty INTEGER,          -- 1-10 scale
		ai_suggestions TEXT,         -- Improvement recommendations

		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_trees_owner_id ON trees(owner_id);
	CREATE INDEX IF NOT EXISTS idx_trees_is_public ON trees(is_public);
	CREATE INDEX IF NOT EXISTS idx_trees_created_at ON trees(created_at);

	CREATE INDEX IF NOT EXISTS idx_nodes_tree_id ON nodes(tree_id);
	CREATE INDEX IF NOT EXISTS idx_nodes_parent_id ON nodes(parent_id);
	CREATE INDEX IF NOT EXISTS idx_nodes_depth ON nodes(depth);
	CREATE INDEX IF NOT EXISTS idx_nodes_difficulty ON nodes(difficulty);

	-- Triggers for updated_at
	CREATE TRIGGER IF NOT EXISTS update_trees_timestamp
		AFTER UPDATE ON trees
		FOR EACH ROW
		BEGIN
			UPDATE trees SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
		END;

	CREATE TRIGGER IF NOT EXISTS update_nodes_timestamp
		AFTER UPDATE ON nodes
		FOR EACH ROW
		BEGIN
			UPDATE nodes SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
		END;
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// === Tree Operations ===

// CreateTree inserts a new tree into the database
func (db *DB) CreateTree(tree *models.Tree) error {
	query := `
		INSERT INTO trees (id, title, description, owner_id, is_public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query, tree.ID, tree.Title, tree.Description, tree.OwnerID, tree.IsPublic, tree.CreatedAt, tree.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}
	return nil
}

// GetTree retrieves a tree by ID
func (db *DB) GetTree(id string) (*models.Tree, error) {
	query := `
		SELECT id, title, description, owner_id, is_public, created_at, updated_at
		FROM trees
		WHERE id = ?
	`
	tree := &models.Tree{}
	err := db.conn.QueryRow(query, id).Scan(
		&tree.ID, &tree.Title, &tree.Description, &tree.OwnerID, &tree.IsPublic, &tree.CreatedAt, &tree.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}
	return tree, nil
}

// GetTreeWithNodes retrieves a tree with all its nodes
func (db *DB) GetTreeWithNodes(id string) (*models.Tree, error) {
	tree, err := db.GetTree(id)
	if err != nil || tree == nil {
		return tree, err
	}

	nodes, err := db.GetNodesByTreeID(id)
	if err != nil {
		return nil, err
	}

	tree.Nodes = nodes
	tree.Stats = db.calculateTreeStats(nodes)
	return tree, nil
}

// UpdateTree updates an existing tree
func (db *DB) UpdateTree(id string, req *models.UpdateTreeRequest) error {
	setParts := []string{}
	args := []interface{}{}

	if req.Title != nil {
		setParts = append(setParts, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.IsPublic != nil {
		setParts = append(setParts, "is_public = ?")
		args = append(args, *req.IsPublic)
	}

	if len(setParts) == 0 {
		return nil // Nothing to update
	}

	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE trees SET %s WHERE id = ?", strings.Join(setParts, ", "))
	_, err := db.conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update tree: %w", err)
	}
	return nil
}

// DeleteTree deletes a tree and all its nodes
func (db *DB) DeleteTree(id string) error {
	_, err := db.conn.Exec("DELETE FROM trees WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete tree: %w", err)
	}
	return nil
}

// SearchTrees searches for trees based on criteria
func (db *DB) SearchTrees(req *models.TreeSearchRequest) ([]models.Tree, int, error) {
	whereParts := []string{}
	args := []interface{}{}

	if req.Query != "" {
		whereParts = append(whereParts, "(title LIKE ? OR description LIKE ?)")
		searchTerm := "%" + req.Query + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if req.OwnerID != "" {
		whereParts = append(whereParts, "owner_id = ?")
		args = append(args, req.OwnerID)
	}

	if req.IsPublic != nil {
		whereParts = append(whereParts, "is_public = ?")
		args = append(args, *req.IsPublic)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Count total results
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM trees %s", whereClause)
	var total int
	err := db.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count trees: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, title, description, owner_id, is_public, created_at, updated_at
		FROM trees %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, req.SortBy, req.SortOrder)

	args = append(args, req.Limit, req.Offset)
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search trees: %w", err)
	}
	defer rows.Close()

	trees := []models.Tree{}
	for rows.Next() {
		tree := models.Tree{}
		err := rows.Scan(&tree.ID, &tree.Title, &tree.Description, &tree.OwnerID, &tree.IsPublic, &tree.CreatedAt, &tree.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tree: %w", err)
		}
		trees = append(trees, tree)
	}

	return trees, total, nil
}

// === Node Operations ===

// CreateNode inserts a new node into the database
func (db *DB) CreateNode(node *models.Node) error {
	query := `
		INSERT INTO nodes (id, tree_id, parent_id, title, content, depth, position, embedding, concepts, difficulty, ai_suggestions, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query,
		node.ID, node.TreeID, node.ParentID, node.Title, node.Content,
		node.Depth, node.Position, node.Embedding, node.Concepts,
		node.Difficulty, node.AISuggestions, node.CreatedAt, node.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}
	return nil
}

// GetNode retrieves a node by ID
func (db *DB) GetNode(id string) (*models.Node, error) {
	query := `
		SELECT id, tree_id, parent_id, title, content, depth, position, embedding, concepts, difficulty, ai_suggestions, created_at, updated_at
		FROM nodes
		WHERE id = ?
	`
	node := &models.Node{}

	// Use nullable types for fields that can be NULL
	var embedding sql.NullString
	var concepts sql.NullString
	var difficulty sql.NullInt64
	var aiSuggestions sql.NullString

	err := db.conn.QueryRow(query, id).Scan(
		&node.ID, &node.TreeID, &node.ParentID, &node.Title, &node.Content,
		&node.Depth, &node.Position, &embedding, &concepts,
		&difficulty, &aiSuggestions, &node.CreatedAt, &node.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// Handle NULL values
	if embedding.Valid {
		node.Embedding = []byte(embedding.String)
	}
	if concepts.Valid {
		node.Concepts = []byte(concepts.String)
	}
	if difficulty.Valid {
		difficultyInt := int(difficulty.Int64)
		node.Difficulty = &difficultyInt
	}
	if aiSuggestions.Valid {
		node.AISuggestions = &aiSuggestions.String
	}
	return node, nil
}

// GetNodesByTreeID retrieves all nodes for a tree, organized hierarchically
func (db *DB) GetNodesByTreeID(treeID string) ([]models.Node, error) {
	query := `
		SELECT id, tree_id, parent_id, title, content, depth, position, embedding, concepts, difficulty, ai_suggestions, created_at, updated_at
		FROM nodes
		WHERE tree_id = ?
		ORDER BY depth, position
	`
	rows, err := db.conn.Query(query, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}
	defer rows.Close()

	nodes := []models.Node{}
	for rows.Next() {
		node := models.Node{}

		// Use nullable types for fields that can be NULL
		var embedding sql.NullString
		var concepts sql.NullString
		var difficulty sql.NullInt64
		var aiSuggestions sql.NullString

		err := rows.Scan(
			&node.ID, &node.TreeID, &node.ParentID, &node.Title, &node.Content,
			&node.Depth, &node.Position, &embedding, &concepts,
			&difficulty, &aiSuggestions, &node.CreatedAt, &node.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}

		// Handle NULL values
		if embedding.Valid {
			node.Embedding = []byte(embedding.String)
		}
		if concepts.Valid {
			node.Concepts = []byte(concepts.String)
		}
		if difficulty.Valid {
			difficultyInt := int(difficulty.Int64)
			node.Difficulty = &difficultyInt
		}
		if aiSuggestions.Valid {
			node.AISuggestions = &aiSuggestions.String
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// UpdateNode updates an existing node
func (db *DB) UpdateNode(id string, req *models.UpdateNodeRequest) error {
	setParts := []string{}
	args := []interface{}{}

	if req.Title != nil {
		setParts = append(setParts, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Content != nil {
		setParts = append(setParts, "content = ?")
		args = append(args, *req.Content)
	}
	if req.Position != nil {
		setParts = append(setParts, "position = ?")
		args = append(args, *req.Position)
	}

	if len(setParts) == 0 {
		return nil // Nothing to update
	}

	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE nodes SET %s WHERE id = ?", strings.Join(setParts, ", "))
	_, err := db.conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}
	return nil
}

// UpdateNodeAI updates AI-related fields for a node
func (db *DB) UpdateNodeAI(id string, embedding []byte, concepts json.RawMessage, difficulty *int, suggestions *string) error {
	query := `
		UPDATE nodes
		SET embedding = ?, concepts = ?, difficulty = ?, ai_suggestions = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.conn.Exec(query, embedding, concepts, difficulty, suggestions, id)
	if err != nil {
		return fmt.Errorf("failed to update node AI data: %w", err)
	}
	return nil
}

// DeleteNode deletes a node and all its children
func (db *DB) DeleteNode(id string) error {
	_, err := db.conn.Exec("DELETE FROM nodes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}
	return nil
}

// === Vector/Semantic Search Operations ===

// SemanticSearch performs vector similarity search across nodes
func (db *DB) SemanticSearch(queryEmbedding []byte, treeID *string, threshold float64, limit int) ([]models.SemanticSearchResult, error) {
	whereParts := []string{"embedding IS NOT NULL"}
	args := []interface{}{}

	if treeID != nil {
		whereParts = append(whereParts, "n.tree_id = ?")
		args = append(args, *treeID)
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Note: This is a simplified version. For production, you'd want to use
	// a proper vector similarity function. SQLite doesn't have built-in vector
	// similarity, so this would need to be implemented differently or use an extension.
	query := fmt.Sprintf(`
		SELECT
			n.id, n.tree_id, n.parent_id, n.title, n.content, n.depth, n.position,
			n.embedding, n.concepts, n.difficulty, n.ai_suggestions, n.created_at, n.updated_at,
			t.id, t.title, t.description, t.owner_id, t.is_public, t.created_at, t.updated_at,
			1.0 as similarity
		FROM nodes n
		JOIN trees t ON n.tree_id = t.id
		WHERE %s
		ORDER BY similarity DESC
		LIMIT ?
	`, whereClause)

	args = append(args, limit)
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to perform semantic search: %w", err)
	}
	defer rows.Close()

	results := []models.SemanticSearchResult{}
	for rows.Next() {
		var result models.SemanticSearchResult
		err := rows.Scan(
			&result.Node.ID, &result.Node.TreeID, &result.Node.ParentID, &result.Node.Title, &result.Node.Content,
			&result.Node.Depth, &result.Node.Position, &result.Node.Embedding, &result.Node.Concepts,
			&result.Node.Difficulty, &result.Node.AISuggestions, &result.Node.CreatedAt, &result.Node.UpdatedAt,
			&result.Tree.ID, &result.Tree.Title, &result.Tree.Description, &result.Tree.OwnerID,
			&result.Tree.IsPublic, &result.Tree.CreatedAt, &result.Tree.UpdatedAt,
			&result.Similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		result.Relevance = result.GetRelevanceLevel()
		results = append(results, result)
	}

	return results, nil
}

// === Helper Functions ===

// calculateTreeStats computes statistics for a tree based on its nodes
func (db *DB) calculateTreeStats(nodes []models.Node) *models.TreeStats {
	if len(nodes) == 0 {
		return &models.TreeStats{}
	}

	stats := &models.TreeStats{
		NodeCount: len(nodes),
	}

	totalDifficulty := 0
	difficultyCount := 0
	conceptSet := make(map[string]bool)

	for _, node := range nodes {
		// Calculate max depth
		if node.Depth > stats.MaxDepth {
			stats.MaxDepth = node.Depth
		}

		// Calculate average difficulty
		if node.Difficulty != nil {
			totalDifficulty += *node.Difficulty
			difficultyCount++
		}

		// Count unique concepts
		concepts, err := node.GetConceptsAsSlice()
		if err == nil {
			for _, concept := range concepts {
				conceptSet[concept.Term] = true
			}
		}
	}

	if difficultyCount > 0 {
		stats.AvgDifficulty = float64(totalDifficulty) / float64(difficultyCount)
	}

	stats.ConceptCount = len(conceptSet)
	return stats
}

// Health checks the database connection
func (db *DB) Health() error {
	return db.conn.Ping()
}
