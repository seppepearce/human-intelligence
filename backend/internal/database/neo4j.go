package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"human-intelligence/internal/models"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jService handles Neo4j database operations
type Neo4jService struct {
	driver   neo4j.DriverWithContext
	database string
	config   *Neo4jConfig
}

// Neo4jConfig holds Neo4j connection configuration
type Neo4jConfig struct {
	URI                          string
	Username                     string
	Password                     string
	Database                     string
	MaxConnectionPoolSize        int
	ConnectionTimeout            time.Duration
	MaxTransactionRetries        int
	InitialRetryDelay            time.Duration
	MaxRetryDelay                time.Duration
	RetryDelayMultiplier         float64
	ConnectionAcquisitionTimeout time.Duration
}

// NewNeo4jService creates a new Neo4j service instance
func NewNeo4jService(config *Neo4jConfig) (*Neo4jService, error) {
	// Configure authentication
	auth := neo4j.BasicAuth(config.Username, config.Password, "")

	// Configure driver
	driverConfig := func(conf *neo4j.Config) {
		conf.MaxConnectionPoolSize = config.MaxConnectionPoolSize
		conf.ConnectionAcquisitionTimeout = config.ConnectionAcquisitionTimeout
		conf.MaxTransactionRetryTime = config.MaxRetryDelay
	}

	// Create driver
	driver, err := neo4j.NewDriverWithContext(config.URI, auth, driverConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	service := &Neo4jService{
		driver:   driver,
		database: config.Database,
		config:   config,
	}

	// Verify connectivity
	if err := service.VerifyConnectivity(context.Background()); err != nil {
		driver.Close(context.Background())
		return nil, fmt.Errorf("failed to verify Neo4j connectivity: %w", err)
	}

	// Initialize schema
	if err := service.InitializeSchema(context.Background()); err != nil {
		log.Printf("Warning: failed to initialize schema: %v", err)
	}

	return service, nil
}

// Close closes the Neo4j driver
func (s *Neo4jService) Close(ctx context.Context) error {
	return s.driver.Close(ctx)
}

// VerifyConnectivity verifies the connection to Neo4j
func (s *Neo4jService) VerifyConnectivity(ctx context.Context) error {
	return s.driver.VerifyConnectivity(ctx)
}

// InitializeSchema creates indexes and constraints
func (s *Neo4jService) InitializeSchema(ctx context.Context) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	// Create constraints and indexes
	constraints := []string{
		// User constraints
		"CREATE CONSTRAINT user_id_unique IF NOT EXISTS FOR (u:User) REQUIRE u.id IS UNIQUE",
		"CREATE CONSTRAINT user_username_unique IF NOT EXISTS FOR (u:User) REQUIRE u.username IS UNIQUE",
		"CREATE CONSTRAINT user_email_unique IF NOT EXISTS FOR (u:User) REQUIRE u.email IS UNIQUE",

		// Node constraints
		"CREATE CONSTRAINT node_id_unique IF NOT EXISTS FOR (n:Node) REQUIRE n.id IS UNIQUE",

		// Tree constraints
		"CREATE CONSTRAINT tree_id_unique IF NOT EXISTS FOR (t:Tree) REQUIRE t.id IS UNIQUE",

		// Tag constraints
		"CREATE CONSTRAINT tag_id_unique IF NOT EXISTS FOR (tag:Tag) REQUIRE tag.id IS UNIQUE",
		"CREATE CONSTRAINT tag_name_unique IF NOT EXISTS FOR (tag:Tag) REQUIRE tag.name IS UNIQUE",
	}

	indexes := []string{
		// User indexes
		"CREATE INDEX user_created_at IF NOT EXISTS FOR (u:User) ON (u.created_at)",
		"CREATE INDEX user_is_active IF NOT EXISTS FOR (u:User) ON (u.is_active)",

		// Node indexes
		"CREATE INDEX node_title IF NOT EXISTS FOR (n:Node) ON (n.title)",
		"CREATE INDEX node_content_type IF NOT EXISTS FOR (n:Node) ON (n.content_type)",
		"CREATE INDEX node_is_public IF NOT EXISTS FOR (n:Node) ON (n.is_public)",
		"CREATE INDEX node_is_published IF NOT EXISTS FOR (n:Node) ON (n.is_published)",
		"CREATE INDEX node_created_at IF NOT EXISTS FOR (n:Node) ON (n.created_at)",
		"CREATE INDEX node_owner_id IF NOT EXISTS FOR (n:Node) ON (n.owner_id)",
		"CREATE FULLTEXT INDEX node_search IF NOT EXISTS FOR (n:Node) ON EACH [n.title, n.content, n.description]",

		// Tree indexes
		"CREATE INDEX tree_name IF NOT EXISTS FOR (t:Tree) ON (t.name)",
		"CREATE INDEX tree_is_public IF NOT EXISTS FOR (t:Tree) ON (t.is_public)",
		"CREATE INDEX tree_created_at IF NOT EXISTS FOR (t:Tree) ON (t.created_at)",
		"CREATE INDEX tree_owner_id IF NOT EXISTS FOR (t:Tree) ON (t.owner_id)",
		"CREATE FULLTEXT INDEX tree_search IF NOT EXISTS FOR (t:Tree) ON EACH [t.name, t.description]",

		// Tag indexes
		"CREATE INDEX tag_usage_count IF NOT EXISTS FOR (tag:Tag) ON (tag.usage_count)",
	}

	// Execute constraints
	for _, constraint := range constraints {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			_, err := tx.Run(ctx, constraint, nil)
			return nil, err
		})
		if err != nil {
			log.Printf("Warning: failed to create constraint: %s, error: %v", constraint, err)
		}
	}

	// Execute indexes
	for _, index := range indexes {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			_, err := tx.Run(ctx, index, nil)
			return nil, err
		})
		if err != nil {
			log.Printf("Warning: failed to create index: %s, error: %v", index, err)
		}
	}

	return nil
}

// ExecuteQuery executes a read query
func (s *Neo4jService) ExecuteQuery(ctx context.Context, query string, params map[string]interface{}) ([]*neo4j.Record, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Collect(ctx)
	})

	if err != nil {
		return nil, err
	}

	return result.([]*neo4j.Record), nil
}

// ExecuteWrite executes a write query
func (s *Neo4jService) ExecuteWrite(ctx context.Context, query string, params map[string]interface{}) (*neo4j.ResultSummary, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	if err != nil {
		return nil, err
	}

	return result.(*neo4j.ResultSummary), nil
}

// User operations

// CreateUser creates a new user
func (s *Neo4jService) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		CREATE (u:User {
			id: $id,
			username: $username,
			email: $email,
			password_hash: $password_hash,
			first_name: $first_name,
			last_name: $last_name,
			avatar: $avatar,
			bio: $bio,
			is_active: $is_active,
			is_verified: $is_verified,
			created_at: datetime(),
			updated_at: datetime(),
			preferences: $preferences,
			followers_count: 0,
			following_count: 0,
			nodes_count: 0,
			trees_count: 0,
			reputation_score: 0
		})
		RETURN u
	`

	params := map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"avatar":        user.Avatar,
		"bio":           user.Bio,
		"is_active":     user.IsActive,
		"is_verified":   user.IsVerified,
		"preferences":   user.Preferences,
	}
	_, err := s.ExecuteWrite(ctx, query, params)
	return err
}

// GetUserByID retrieves a user by ID
func (s *Neo4jService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		MATCH (u:User {id: $id})
		RETURN u
	`

	params := map[string]interface{}{
		"id": id,
	}

	records, err := s.ExecuteQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, models.ErrRecordNotFound
	}

	user := &models.User{}
	if err := user.FromRecord(records[0]); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *Neo4jService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		MATCH (u:User {username: $username})
		RETURN u
	`

	params := map[string]interface{}{
		"username": username,
	}

	records, err := s.ExecuteQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, models.ErrRecordNotFound
	}

	user := &models.User{}
	if err := user.FromRecord(records[0]); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user
func (s *Neo4jService) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		MATCH (u:User {id: $id})
		SET u += $props
		RETURN u
	`

	props := user.GetProperties()
	delete(props, "id") // Don't update ID

	params := map[string]interface{}{
		"id":    user.ID,
		"props": props,
	}

	_, err := s.ExecuteWrite(ctx, query, params)
	return err
}

// Node operations

// CreateNode creates a new node
func (s *Neo4jService) CreateNode(ctx context.Context, node *models.Node) error {
	query := `
		CREATE (n:Node {
			id: $id,
			title: $title,
			content: $content,
			content_type: $content_type,
			description: $description,
			is_public: $is_public,
			is_published: $is_published,
			created_at: datetime(),
			updated_at: datetime(),
			views_count: 0,
			likes_count: 0,
			shares_count: 0,
			metadata: $metadata,
			owner_id: $owner_id,
			parent_node_id: $parent_node_id,
			level: $level,
			position: $position
		})
		RETURN n
	`

	params := map[string]interface{}{
		"id":             node.ID,
		"title":          node.Title,
		"content":        node.Content,
		"content_type":   node.ContentType,
		"description":    node.Description,
		"is_public":      node.IsPublic,
		"is_published":   node.IsPublished,
		"metadata":       node.Metadata,
		"owner_id":       node.OwnerID,
		"parent_node_id": node.ParentNodeID,
		"level":          node.Level,
		"position":       node.Position,
	}

	_, err := s.ExecuteWrite(ctx, query, params)
	return err
}

// GetNodeByID retrieves a node by ID
func (s *Neo4jService) GetNodeByID(ctx context.Context, id string) (*models.Node, error) {
	query := `
		MATCH (n:Node {id: $id})
		OPTIONAL MATCH (n)<-[:CREATED]-(owner:User)
		OPTIONAL MATCH (n)-[:TAGGED]->(tag:Tag)
		RETURN n, owner, collect(tag.name) as tags
	`

	params := map[string]interface{}{
		"id": id,
	}

	records, err := s.ExecuteQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, models.ErrRecordNotFound
	}

	record := records[0]
	node := &models.Node{}

	// Get node data
	if nodeData, ok := record.Get("n"); ok {
		nodeValue := nodeData.(neo4j.Node)
		if err := node.FromNode(nodeValue); err != nil {
			return nil, err
		}
	}

	// Get owner data
	if ownerData, ok := record.Get("owner"); ok && ownerData != nil {
		owner := &models.User{}
		ownerValue := ownerData.(neo4j.Node)
		if err := owner.FromNode(ownerValue); err == nil {
			node.Owner = owner
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

	return node, nil
}

// SearchNodes searches for nodes based on query parameters
func (s *Neo4jService) SearchNodes(ctx context.Context, query *models.GraphQuery) (*models.GraphResponse, error) {
	cypher := `
		MATCH (n:Node)
		OPTIONAL MATCH (n)<-[:CREATED]-(owner:User)
		OPTIONAL MATCH (n)-[:TAGGED]->(tag:Tag)
		WHERE 1=1
	`
	params := make(map[string]interface{})

	// Add search conditions
	if query.SearchString != "" {
		cypher += ` AND (n.title CONTAINS $search OR n.content CONTAINS $search OR n.description CONTAINS $search)`
		params["search"] = query.SearchString
	}

	if query.IsPublic != nil {
		cypher += ` AND n.is_public = $is_public`
		params["is_public"] = *query.IsPublic
	}

	if query.OwnerID != "" {
		cypher += ` AND n.owner_id = $owner_id`
		params["owner_id"] = query.OwnerID
	}

	if query.ContentType != "" {
		cypher += ` AND n.content_type = $content_type`
		params["content_type"] = query.ContentType
	}

	// Add tag filtering
	if len(query.Tags) > 0 {
		cypher += ` AND tag.name IN $tags`
		params["tags"] = query.Tags
	}

	cypher += ` RETURN n, owner, collect(tag.name) as tags`

	// Add ordering
	if query.SortBy != "" {
		order := "ASC"
		if query.SortOrder == "desc" {
			order = "DESC"
		}
		cypher += fmt.Sprintf(` ORDER BY n.%s %s`, query.SortBy, order)
	} else {
		cypher += ` ORDER BY n.created_at DESC`
	}

	// Add pagination
	if query.Limit > 0 {
		cypher += ` SKIP $offset LIMIT $limit`
		params["offset"] = query.Offset
		params["limit"] = query.Limit
	}

	records, err := s.ExecuteQuery(ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	nodes := make([]*models.Node, len(records))
	for i, record := range records {
		node := &models.Node{}

		// Parse node
		if nodeData, ok := record.Get("n"); ok {
			nodeValue := nodeData.(neo4j.Node)
			if err := node.FromNode(nodeValue); err != nil {
				return nil, err
			}
		}

		// Parse owner
		if ownerData, ok := record.Get("owner"); ok && ownerData != nil {
			owner := &models.User{}
			ownerValue := ownerData.(neo4j.Node)
			if err := owner.FromNode(ownerValue); err == nil {
				node.Owner = owner
			}
		}

		// Parse tags
		if tagsData, ok := record.Get("tags"); ok {
			if tags, ok := tagsData.([]interface{}); ok {
				node.Tags = make([]string, len(tags))
				for j, tag := range tags {
					if tagStr, ok := tag.(string); ok {
						node.Tags[j] = tagStr
					}
				}
			}
		}

		nodes[i] = node
	}

	// Get total count
	countQuery := `
		MATCH (n:Node)
		WHERE 1=1
	`
	countParams := make(map[string]interface{})

	if query.SearchString != "" {
		countQuery += ` AND (n.title CONTAINS $search OR n.content CONTAINS $search OR n.description CONTAINS $search)`
		countParams["search"] = query.SearchString
	}

	if query.IsPublic != nil {
		countQuery += ` AND n.is_public = $is_public`
		countParams["is_public"] = *query.IsPublic
	}

	if query.OwnerID != "" {
		countQuery += ` AND n.owner_id = $owner_id`
		countParams["owner_id"] = query.OwnerID
	}

	if query.ContentType != "" {
		countQuery += ` AND n.content_type = $content_type`
		countParams["content_type"] = query.ContentType
	}

	if len(query.Tags) > 0 {
		countQuery += `
			AND EXISTS {
				MATCH (n)-[:TAGGED]->(tag:Tag)
				WHERE tag.name IN $tags
			}
		`
		countParams["tags"] = query.Tags
	}

	countQuery += ` RETURN count(n) as total`

	countRecords, err := s.ExecuteQuery(ctx, countQuery, countParams)
	if err != nil {
		return nil, err
	}

	var total int64
	if len(countRecords) > 0 {
		if totalData, ok := countRecords[0].Get("total"); ok {
			total = totalData.(int64)
		}
	}

	response := &models.GraphResponse{
		Nodes:    nodes,
		Total:    total,
		Page:     query.Offset/query.Limit + 1,
		PageSize: query.Limit,
		HasMore:  int64(query.Offset+query.Limit) < total,
	}

	return response, nil
}

// Tree operations

// CreateTree creates a new tree
func (s *Neo4jService) CreateTree(ctx context.Context, tree *models.Tree) error {
	query := `
		CREATE (t:Tree {
			id: $id,
			name: $name,
			description: $description,
			is_public: $is_public,
			is_template: $is_template,
			created_at: datetime(),
			updated_at: datetime(),
			views_count: 0,
			likes_count: 0,
			forks_count: 0,
			nodes_count: $nodes_count,
			owner_id: $owner_id,
			parent_tree_id: $parent_tree_id
		})
		RETURN t
	`

	params := map[string]interface{}{
		"id":             tree.ID,
		"name":           tree.Name,
		"description":    tree.Description,
		"is_public":      tree.IsPublic,
		"is_template":    tree.IsTemplate,
		"nodes_count":    tree.NodesCount,
		"owner_id":       tree.OwnerID,
		"parent_tree_id": tree.ParentTreeID,
	}

	_, err := s.ExecuteWrite(ctx, query, params)
	return err
}

// AddNodeToTree adds a node to a tree
func (s *Neo4jService) AddNodeToTree(ctx context.Context, treeID, nodeID string, position int) error {
	query := `
		MATCH (t:Tree {id: $tree_id})
		MATCH (n:Node {id: $node_id})
		CREATE (t)-[:CONTAINS {position: $position, added_at: datetime()}]->(n)
		SET t.nodes_count = t.nodes_count + 1
		RETURN t, n
	`

	params := map[string]interface{}{
		"tree_id":  treeID,
		"node_id":  nodeID,
		"position": position,
	}

	_, err := s.ExecuteWrite(ctx, query, params)
	return err
}

// Tag operations

// CreateTag creates a new tag
func (s *Neo4jService) CreateTag(ctx context.Context, tag *models.Tag) error {
	query := `
		CREATE (tag:Tag {
			id: $id,
			name: $name,
			color: $color,
			created_at: datetime(),
			usage_count: 0
		})
		RETURN tag
	`

	params := map[string]interface{}{
		"id":    tag.ID,
		"name":  tag.Name,
		"color": tag.Color,
	}
	_, err := s.ExecuteWrite(ctx, query, params)
	return err
}

// TagNode tags a node with a tag
func (s *Neo4jService) TagNode(ctx context.Context, nodeID, tagName string) error {
	query := `
		MATCH (n:Node {id: $node_id})
		MERGE (tag:Tag {name: $tag_name})
		ON CREATE SET tag.id = randomUUID(), tag.created_at = datetime(), tag.usage_count = 0, tag.color = '#' + substring(randomUUID(), 0, 6)
		CREATE (n)-[:TAGGED {tagged_at: datetime()}]->(tag)
		SET tag.usage_count = tag.usage_count + 1
		RETURN n, tag
	`

	params := map[string]interface{}{
		"node_id":  nodeID,
		"tag_name": tagName,
	}

	_, err := s.ExecuteWrite(ctx, query, params)
	return err
}

// GetPopularTags gets the most popular tags
func (s *Neo4jService) GetPopularTags(ctx context.Context, limit int) ([]*models.Tag, error) {
	query := `
		MATCH (tag:Tag)
		RETURN tag
		ORDER BY tag.usage_count DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"limit": limit,
	}

	records, err := s.ExecuteQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	tags := make([]*models.Tag, len(records))
	for i, record := range records {
		tag := &models.Tag{}
		if err := tag.FromRecord(record); err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	return tags, nil
}

// GetVisualizationData gets data for graph visualization
func (s *Neo4jService) GetVisualizationData(ctx context.Context, userID string, depth int, limit int) (*models.VisualizationData, error) {
	query := `
		MATCH (u:User {id: $user_id})
		OPTIONAL MATCH (u)-[:CREATED]->(n:Node)
		OPTIONAL MATCH (u)-[:CREATED]->(t:Tree)
		OPTIONAL MATCH (t)-[:CONTAINS]->(tn:Node)
		OPTIONAL MATCH (n)-[:TAGGED]->(tag:Tag)
		OPTIONAL MATCH (tn)-[:TAGGED]->(ttag:Tag)
		WITH u, collect(DISTINCT n) + collect(DISTINCT t) + collect(DISTINCT tn) + collect(DISTINCT tag) + collect(DISTINCT ttag) as entities
		UNWIND entities as entity
		RETURN entity
		LIMIT $limit
	`

	params := map[string]interface{}{
		"user_id": userID,
		"depth":   depth,
		"limit":   limit,
	}

	records, err := s.ExecuteQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := make([]*models.VisualizationNode, 0)
	links := make([]*models.VisualizationLink, 0)

	// Process records to create visualization data
	for _, record := range records {
		if entityData, ok := record.Get("entity"); ok && entityData != nil {
			entity := entityData.(neo4j.Node)
			labels := entity.Labels

			vizNode := &models.VisualizationNode{
				ID:    fmt.Sprintf("%v", entity.Props["id"]),
				Label: fmt.Sprintf("%v", entity.Props["title"]),
				Size:  10,
				Color: "#6366f1",
			}

			if len(labels) > 0 {
				vizNode.Type = labels[0]
				switch labels[0] {
				case "User":
					vizNode.Color = "#ef4444"
					vizNode.Size = 20
					if name, ok := entity.Props["username"]; ok {
						vizNode.Label = fmt.Sprintf("%v", name)
					}
				case "Node":
					vizNode.Color = "#3b82f6"
					vizNode.Size = 15
				case "Tree":
					vizNode.Color = "#10b981"
					vizNode.Size = 18
					if name, ok := entity.Props["name"]; ok {
						vizNode.Label = fmt.Sprintf("%v", name)
					}
				case "Tag":
					vizNode.Color = "#f59e0b"
					vizNode.Size = 8
					if name, ok := entity.Props["name"]; ok {
						vizNode.Label = fmt.Sprintf("%v", name)
					}
				}
			}

			nodes = append(nodes, vizNode)
		}
	}

	// Get relationships
	relQuery := `
		MATCH (a)-[r]->(b)
		WHERE a.id IN $node_ids AND b.id IN $node_ids
		RETURN a.id as source, b.id as target, type(r) as relationship_type
	`

	nodeIDs := make([]string, len(nodes))
	for i, node := range nodes {
		nodeIDs[i] = node.ID
	}

	relParams := map[string]interface{}{
		"node_ids": nodeIDs,
	}

	relRecords, err := s.ExecuteQuery(ctx, relQuery, relParams)
	if err == nil {
		for _, record := range relRecords {
			if source, ok := record.Get("source"); ok {
				if target, ok := record.Get("target"); ok {
					if relType, ok := record.Get("relationship_type"); ok {
						link := &models.VisualizationLink{
							Source:   fmt.Sprintf("%v", source),
							Target:   fmt.Sprintf("%v", target),
							Type:     fmt.Sprintf("%v", relType),
							Weight:   1.0,
							Distance: 50.0,
							Color:    "#94a3b8",
						}
						links = append(links, link)
					}
				}
			}
		}
	}

	return &models.VisualizationData{
		Nodes: nodes,
		Links: links,
	}, nil
}

// Health check
func (s *Neo4jService) HealthCheck(ctx context.Context) error {
	query := "RETURN 1 as health"
	_, err := s.ExecuteQuery(ctx, query, nil)
	return err
}
