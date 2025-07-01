# 🌳 Neo4j Learning Trees Platform

A knowledge-sharing platform built on Neo4j that organizes learning into hierarchical tree structures. Users can create learning trees with branching nodes, representing natural learning progressions from simple to complex concepts.

## 🚀 Quick Start

```bash
# Start the platform
./quick-start.sh
```

This will:
1. Start Neo4j and Redis via Docker Compose
2. Build the Go backend
3. Start the server on port 8085

## 🏗️ Architecture

- **Backend**: Go with Gin framework
- **Database**: Neo4j 5.15 (graph database)
- **Cache**: Redis 7
- **Development**: Docker Compose for services

## 📊 API Endpoints

### Trees (Hierarchical Learning Paths)
- `POST /api/v1/trees` - Create tree with root node
- `GET /api/v1/trees/{treeId}` - Get complete tree structure
- `PUT /api/v1/trees/{treeId}` - Update tree metadata
- `DELETE /api/v1/trees/{treeId}` - Delete entire tree
- `POST /api/v1/trees/{treeId}/nodes` - Add child node to tree
- `GET /api/v1/trees/{treeId}/nodes` - Get all nodes in tree

### Nodes (Learning Content)
- `POST /api/v1/nodes` - Create standalone node
- `GET /api/v1/nodes/{nodeId}` - Get node details
- `GET /api/v1/nodes?q=search` - Search nodes

### Users & Tags
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{userId}` - Get user details
- `POST /api/v1/tags` - Create tag
- `GET /api/v1/tags/popular` - Get popular tags
- `POST /api/v1/nodes/{nodeId}/tags/{tagName}` - Tag a node

### Visualization
- `GET /api/v1/graph/visualization/{userId}` - Get graph data for D3.js

## 🌱 Example: Creating a Learning Tree

```bash
# 1. Create a learning tree with root node
curl -X POST http://localhost:8085/api/v1/trees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Learning C Programming",
    "description": "Complete journey through C programming",
    "is_public": true,
    "root_node": {
      "title": "Hello World",
      "content": "First program - printing to screen",
      "content_type": "text"
    }
  }'

# 2. Add child nodes to create branches
curl -X POST http://localhost:8085/api/v1/trees/{treeId}/nodes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Variables & Data Types",
    "content": "Learning about int, char, float...",
    "parent_node_id": "{rootNodeId}",
    "position": 1
  }'
```

## 🎯 Core Concepts

### Trees (Terrestrial View)
- **Learning Paths**: Focused, hierarchical learning journeys
- **Root Nodes**: Starting points (e.g., "Hello World")
- **Branching**: Natural progression from simple to complex
- **Hierarchical**: Parent-child relationships between concepts

### Tags (Cosmic View)
- **Semantic Connections**: Link related trees across domains
- **Community Formation**: Users discover content through tags
- **Cross-Pollination**: Bridge different learning paths

## 🧪 Testing

### Web Interface
Open `test-tree-api.html` in your browser for a complete testing interface with:
- Tree creation and management
- Node addition with hierarchy
- Visual tree structure display
- Real-time API testing

### Health Check
```bash
curl http://localhost:8085/health
```

### Manual Testing
```bash
# Create user
curl -X POST http://localhost:8085/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username": "learner", "email": "test@example.com"}'

# Search nodes
curl "http://localhost:8085/api/v1/nodes?q=programming"

# Get popular tags
curl http://localhost:8085/api/v1/tags/popular
```

## 🔧 Development

### Project Structure
```
backend/
├── cmd/server/          # Main application
├── internal/
│   ├── database/        # Neo4j service
│   ├── handlers/        # HTTP handlers
│   └── models/          # Data models
├── go.mod
└── go.sum

docker-compose.yml       # Services (Neo4j, Redis)
quick-start.sh          # Development setup
```

### Key Features Implemented
- ✅ Hierarchical tree creation with root nodes
- ✅ Parent-child node relationships
- ✅ Tree structure retrieval with statistics
- ✅ Node branching and positioning
- ✅ Tag system for semantic connections
- ✅ Search functionality across content
- ✅ Graph visualization data endpoints
- ✅ User management
- ✅ Health monitoring

### Graph Schema
```cypher
# Core relationships
(:User)-[:CREATED]->(:Tree)
(:Tree)-[:CONTAINS]->(:Node)
(:Node)-[:PARENT_OF]->(:Node)
(:Node)-[:CHILD_OF]->(:Node)
(:Node)-[:TAGGED]->(:Tag)
```

## 🚧 Next Steps

1. **Enhanced UI**: Build React/Svelte frontend
2. **Tree Forking**: Copy and extend existing trees
3. **Collaboration**: Multi-user tree editing
4. **AI Integration**: Smart content suggestions
5. **Advanced Search**: Full-text and semantic search

## 🔗 Services

- **Backend**: http://localhost:8085
- **Neo4j Browser**: http://localhost:7474 (neo4j/hi_password)
- **Redis**: localhost:6379

## 📝 License

MIT License - see LICENSE file for details.