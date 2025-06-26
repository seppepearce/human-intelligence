package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"human-intelligence/internal/database"
	"human-intelligence/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SearchHandler struct {
	db *database.DB
}

type SearchRequest struct {
	Query       string             `json:"query" binding:"required"`
	Tags        []string           `json:"tags,omitempty"`
	NodeTypes   []models.NodeType  `json:"node_types,omitempty"`
	AccessLevel models.AccessLevel `json:"access_level,omitempty"`
	MinRating   int                `json:"min_rating,omitempty"`
	DateFrom    *time.Time         `json:"date_from,omitempty"`
	DateTo      *time.Time         `json:"date_to,omitempty"`
	SortBy      string             `json:"sort_by,omitempty"`
	Limit       int                `json:"limit,omitempty"`
	Offset      int                `json:"offset,omitempty"`
}

type SemanticSearchRequest struct {
	Query     string   `json:"query" binding:"required"`
	Tags      []string `json:"tags,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Threshold float64  `json:"threshold,omitempty"` // Similarity threshold
}

type SuggestionsRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	Context    string    `json:"context,omitempty"`    // Current learning context
	Interests  []string  `json:"interests,omitempty"`  // User interests/tags
	Difficulty string    `json:"difficulty,omitempty"` // beginner, intermediate, advanced
	Limit      int       `json:"limit,omitempty"`
}

type SearchResponse struct {
	Results     []SearchResult `json:"results"`
	Total       int            `json:"total"`
	Query       string         `json:"query"`
	Took        int            `json:"took"` // Search time in milliseconds
	Suggestions []string       `json:"suggestions,omitempty"`
}

type SearchResult struct {
	Type        string                 `json:"type"` // "node" or "path"
	ID          uuid.UUID              `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	Score       float64                `json:"score"`
	Snippet     string                 `json:"snippet,omitempty"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	VoteScore   int                    `json:"vote_score"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type AIRecommendation struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Reason      string    `json:"reason"`
	Confidence  float64   `json:"confidence"`
	Tags        []string  `json:"tags"`
}

func NewSearchHandler(db *database.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

func (h *SearchHandler) Search(c *gin.Context) {
	start := time.Now()

	// Parse query parameters
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameter 'q' is required",
		})
		return
	}

	// Parse optional parameters
	limit := 20
	offset := 0
	sortBy := "relevance"

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
		if s == "date" || s == "popularity" || s == "rating" {
			sortBy = s
		}
	}

	tags := []string{}
	if t := c.Query("tags"); t != "" {
		tags = strings.Split(t, ",")
	}

	nodeTypes := []string{}
	if nt := c.Query("node_types"); nt != "" {
		nodeTypes = strings.Split(nt, ",")
	}

	// Perform search
	results, total, err := h.performTextSearch(query, tags, nodeTypes, sortBy, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed",
		})
		return
	}

	// Get search suggestions
	suggestions := h.getSearchSuggestions(query)

	took := int(time.Since(start).Milliseconds())

	c.JSON(http.StatusOK, SearchResponse{
		Results:     results,
		Total:       total,
		Query:       query,
		Took:        took,
		Suggestions: suggestions,
	})
}

func (h *SearchHandler) SemanticSearch(c *gin.Context) {
	start := time.Now()

	var req SemanticSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Threshold == 0 {
		req.Threshold = 0.7 // Default similarity threshold
	}

	// Perform semantic search
	results, err := h.performSemanticSearch(req.Query, req.Tags, req.Limit, req.Threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Semantic search failed",
		})
		return
	}

	took := int(time.Since(start).Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"total":   len(results),
		"query":   req.Query,
		"took":    took,
		"method":  "semantic",
	})
}

func (h *SearchHandler) GetSuggestions(c *gin.Context) {
	userID := h.getUserIDFromContext(c)

	// Parse query parameters
	context := c.Query("context")
	interests := []string{}
	if i := c.Query("interests"); i != "" {
		interests = strings.Split(i, ",")
	}
	difficulty := c.DefaultQuery("difficulty", "beginner")
	limit := 10

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	// Get AI recommendations
	recommendations, err := h.getAIRecommendations(userID, context, interests, difficulty, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get recommendations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"context":         context,
		"user_id":         userID,
		"count":           len(recommendations),
	})
}

// Search implementation methods
func (h *SearchHandler) performTextSearch(query string, tags []string, nodeTypes []string, sortBy string, limit int, offset int) ([]SearchResult, int, error) {
	var results []SearchResult
	var total int

	// Build search query - simplified text search for now
	searchQuery := `
		SELECT 'node' as type, id, title, description, tags, created_by, created_at, vote_score, 1.0 as score
		FROM nodes
		WHERE is_public = true AND (
			title ILIKE $1 OR
			description ILIKE $1 OR
			content ILIKE $1
		)
	`

	args := []interface{}{"%" + query + "%"}
	argIndex := 2

	// Add tag filtering
	if len(tags) > 0 {
		searchQuery += " AND tags && $" + strconv.Itoa(argIndex)
		args = append(args, tags)
		argIndex++
	}

	// Add node type filtering
	if len(nodeTypes) > 0 {
		searchQuery += " AND node_type = ANY($" + strconv.Itoa(argIndex) + ")"
		args = append(args, nodeTypes)
		argIndex++
	}

	// Add sorting
	switch sortBy {
	case "date":
		searchQuery += " ORDER BY created_at DESC"
	case "popularity":
		searchQuery += " ORDER BY vote_score DESC"
	case "rating":
		searchQuery += " ORDER BY vote_score DESC"
	default:
		searchQuery += " ORDER BY vote_score DESC" // Default to popularity
	}

	// Add pagination
	searchQuery += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	// Execute search
	rows, err := h.db.Query(searchQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var result SearchResult
		var tags []string
		err := rows.Scan(&result.Type, &result.ID, &result.Title, &result.Description,
			&tags, &result.CreatedBy, &result.CreatedAt, &result.VoteScore, &result.Score)
		if err != nil {
			continue
		}
		result.Tags = tags
		result.Snippet = h.generateSnippet(result.Description, query)
		results = append(results, result)
	}

	// Get total count (simplified)
	total = len(results)

	return results, total, nil
}

func (h *SearchHandler) performSemanticSearch(query string, tags []string, limit int, threshold float64) ([]SearchResult, error) {
	// Placeholder for semantic search implementation
	// In a real implementation, this would:
	// 1. Generate embeddings for the query
	// 2. Search vector database for similar content
	// 3. Return results sorted by similarity score

	var results []SearchResult

	// For now, return a placeholder message
	results = append(results, SearchResult{
		Type:        "message",
		ID:          uuid.New(),
		Title:       "Semantic Search Coming Soon",
		Description: "Advanced semantic search with AI embeddings will be available soon!",
		Tags:        []string{"ai", "search", "coming-soon"},
		Score:       1.0,
		CreatedAt:   time.Now(),
	})

	return results, nil
}

func (h *SearchHandler) getAIRecommendations(userID uuid.UUID, context string, interests []string, difficulty string, limit int) ([]AIRecommendation, error) {
	var recommendations []AIRecommendation

	// Placeholder for AI recommendations
	// In a real implementation, this would:
	// 1. Analyze user's learning history
	// 2. Use AI to suggest relevant content
	// 3. Consider user preferences and difficulty level

	recommendations = append(recommendations, AIRecommendation{
		ID:          uuid.New(),
		Type:        "path",
		Title:       "Introduction to " + strings.Join(interests, " and "),
		Description: "A beginner-friendly learning path tailored to your interests",
		Reason:      "Based on your interests in " + strings.Join(interests, ", "),
		Confidence:  0.85,
		Tags:        interests,
	})

	if context != "" {
		recommendations = append(recommendations, AIRecommendation{
			ID:          uuid.New(),
			Type:        "node",
			Title:       "Deep Dive: " + context,
			Description: "Advanced concepts building on your current learning",
			Reason:      "Complements your current focus on " + context,
			Confidence:  0.78,
			Tags:        []string{context, "advanced"},
		})
	}

	return recommendations, nil
}

func (h *SearchHandler) getSearchSuggestions(query string) []string {
	// Placeholder for search suggestions
	// In a real implementation, this would analyze popular searches,
	// user behavior, and provide intelligent suggestions

	suggestions := []string{
		query + " tutorial",
		query + " beginner guide",
		query + " advanced",
		"learn " + query,
	}

	return suggestions
}

func (h *SearchHandler) generateSnippet(text, query string) string {
	// Simple snippet generation
	if len(text) <= 150 {
		return text
	}

	// Find query term in text and create snippet around it
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	index := strings.Index(lowerText, lowerQuery)
	if index == -1 {
		// Query not found, return first 150 characters
		return text[:147] + "..."
	}

	// Create snippet around the found term
	start := index - 50
	if start < 0 {
		start = 0
	}

	end := index + len(query) + 50
	if end > len(text) {
		end = len(text)
	}

	snippet := text[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}

	return snippet
}

func (h *SearchHandler) getUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	if id, ok := userID.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}
