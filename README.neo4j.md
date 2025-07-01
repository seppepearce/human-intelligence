# 🧠 Human Intelligence - Neo4j Knowledge Platform

> **Where AI augments Human Intelligence** - A graph-powered knowledge sharing platform built on Neo4j

[![License: MIT](https://img.shields.io/badge/License-MIT-neon.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Node Version](https://img.shields.io/badge/Node-18+-339933.svg)](https://nodejs.org/)
[![Neo4j](https://img.shields.io/badge/Neo4j-5.15+-008CC1.svg)](https://neo4j.com/)
[![Python](https://img.shields.io/badge/Python-3.11+-3776AB.svg)](https://python.org/)

## 🌟 Vision

Human Intelligence is a **graph-powered collaborative learning platform** that leverages Neo4j's graph database to create rich, interconnected knowledge networks. Think of it as the intersection of GitHub's collaboration model, Stack Overflow's knowledge sharing, and a modern graph-based learning platform.

### 🎯 Core Concept
- **Nodes**: Atomic knowledge units (text, video, code, interactive content)
- **Trees**: Structured learning paths with graph-based relationships
- **Graph Relationships**: Rich connections between users, content, and concepts
- **AI-Powered Discovery**: Semantic search and recommendations using graph embeddings
- **Visual Knowledge Maps**: Interactive D3.js visualizations of your knowledge graph

## 🚀 Key Features

### 🔥 MVP Features
- ✅ **Graph-Native Architecture**: Built on Neo4j from the ground up
- ✅ **Interactive Graph Visualizations**: D3.js powered knowledge maps
- ✅ **Smart Tagging System**: Graph-based tag relationships and clustering
- ✅ **Advanced Search**: Full-text, semantic, and graph traversal search
- ✅ **User Authentication**: Secure JWT-based auth with user relationship tracking
- ✅ **Real-time Collaboration**: WebSocket-powered live updates
- ✅ **AI-Powered Recommendations**: Machine learning content suggestions
- ✅ **Responsive Design**: Vaporwave-inspired UI with graph aesthetics

### 🔮 Advanced Features
- 🤖 **Semantic Embeddings**: Vector similarity search in graph context
- 📊 **Graph Analytics**: Community detection, centrality analysis, path finding
- 🔄 **Knowledge Evolution**: Track how concepts and relationships change over time
- 🌐 **Multi-dimensional Relationships**: Rich relationship types with properties
- 📱 **Graph Mobile Experience**: Touch-optimized graph navigation
- 🎯 **Personalized Learning Paths**: AI-generated routes through your knowledge graph

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+ with Gin framework
- **Database**: Neo4j 5.15+ (Graph Database)
- **Cache**: Redis 7+ for session management and caching
- **AI/ML**: Python 3.11+ with FastAPI for embeddings and recommendations
- **Real-time**: Gorilla WebSocket for live updates
- **Authentication**: JWT with bcrypt password hashing

### Frontend
- **Framework**: SvelteKit 4.0+ with TypeScript
- **Styling**: Tailwind CSS with custom graph-themed components
- **Visualization**: D3.js for interactive graph rendering
- **Alternative Graph**: Cytoscape.js for complex network layouts
- **Real-time**: WebSocket client with automatic reconnection
- **State Management**: Svelte stores with graph-aware reactivity

### AI/ML Service
- **Framework**: FastAPI with async support
- **Embeddings**: Sentence Transformers (all-MiniLM-L6-v2)
- **ML Libraries**: scikit-learn, numpy for similarity calculations
- **Vector Operations**: Cosine similarity, clustering algorithms
- **Caching**: Redis for embedding cache and model predictions

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Monitoring**: Prometheus + Grafana with Neo4j metrics
- **Graph Browser**: Neo4j Browser for development
- **CI/CD**: GitHub Actions with graph database testing

## 🏗️ Graph Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   SvelteKit     │    │   Go Backend    │    │     Neo4j       │
│   Frontend      │◄──►│   API Server    │◄──►│ Graph Database  │
│   + D3.js       │    │   + WebSocket   │    │   + Cypher      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Graph Viz     │    │   Python AI     │    │   Redis Cache   │
│   D3.js + Cyto  │    │   Embeddings    │    │   + Sessions    │
│                 │    │   + ML Models   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Graph Data Model

```cypher
// Core Entities
(:User)-[:CREATED]->(:Node)-[:TAGGED]->(:Tag)
(:User)-[:CREATED]->(:Tree)-[:CONTAINS]->(:Node)
(:User)-[:FOLLOWS]->(:User)
(:User)-[:LIKES]->(:Node)
(:Node)-[:RELATED_TO]->(:Node)
(:Tree)-[:FORKED_FROM]->(:Tree)

// Rich Relationships with Properties
(:Node)-[:CONNECTED_TO {strength: 0.8, type: "conceptual"}]->(:Node)
(:User)-[:VIEWED {timestamp: datetime(), duration: 120}]->(:Node)
```

## 🚦 Quick Start

### Prerequisites
- **Go 1.21+**
- **Node.js 18+** 
- **Python 3.11+**
- **Neo4j 5.15+** (or use Docker)
- **Redis 7+** (or use Docker)
- **Docker & Docker Compose** (recommended)

### 🐳 Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-org/human-intelligence.git
cd human-intelligence

# Start with Neo4j stack
docker-compose -f docker-compose.neo4j.yml up -d

# Services will be available at:
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Neo4j Browser: http://localhost:7474 (neo4j/hi_password)
# AI Service: http://localhost:8000
# Grafana: http://localhost:3002 (admin/admin)
```

### 🔧 Option 2: Manual Setup

#### Neo4j Setup
```bash
# Download and start Neo4j
wget https://neo4j.com/artifact.php?name=neo4j-community-5.15.0-unix.tar.gz
tar -xzf neo4j-community-5.15.0-unix.tar.gz
cd neo4j-community-5.15.0

# Configure
echo "dbms.default_database=knowledge_graph" >> conf/neo4j.conf
echo "dbms.security.auth_enabled=true" >> conf/neo4j.conf

# Start Neo4j
bin/neo4j start

# Set password
bin/neo4j-admin set-initial-password hi_password
```

#### Backend Setup
```bash
cd backend

# Copy environment config
cp .env.example .env
# Edit .env with your Neo4j credentials

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

#### AI Service Setup
```bash
cd ai-service

# Create virtual environment
python -m venv venv
source venv/bin/activate  # or `venv\Scripts\activate` on Windows

# Install dependencies
pip install -r requirements.txt

# Start the service
python main.py
```

#### Frontend Setup
```bash
cd frontend

# Install dependencies
pnpm install  # or npm install

# Start development server
pnpm dev
```

## 📊 Graph Visualization Features

### Interactive Knowledge Maps
- **Force-directed layouts** showing content relationships
- **Hierarchical tree views** for structured learning paths
- **Tag clustering** with community detection
- **User network analysis** showing collaboration patterns
- **Concept evolution** tracking changes over time

### Visualization Controls
```typescript
// Example D3.js integration
const graphData = await fetch('/api/v1/graph/visualization')
  .then(res => res.json());

const simulation = d3.forceSimulation(graphData.nodes)
  .force("link", d3.forceLink(graphData.links).id(d => d.id))
  .force("charge", d3.forceManyBody().strength(-300))
  .force("center", d3.forceCenter(width / 2, height / 2));
```

## 🔍 Advanced Search Capabilities

### Graph-Powered Search
```http
# Full-text search with graph context
GET /api/v1/search?q=machine+learning&expand_graph=true

# Semantic similarity search
POST /api/v1/search/semantic
{
  "query": "neural networks",
  "include_related": true,
  "max_depth": 2
}

# Graph traversal queries
POST /api/v1/graph/traverse
{
  "start_node": "user-123",
  "relationship": "CREATED",
  "depth": 3,
  "filters": {"content_type": "video"}
}
```

## 📚 API Documentation

### Core Endpoints

#### Authentication
```http
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
```

#### Nodes (Knowledge Units)
```http
GET    /api/v1/nodes              # List with graph context
POST   /api/v1/nodes              # Create node + relationships
GET    /api/v1/nodes/:id          # Get with relationships
PUT    /api/v1/nodes/:id          # Update + relationship changes
DELETE /api/v1/nodes/:id          # Delete + cleanup relationships
GET    /api/v1/nodes/:id/related  # Get related nodes via graph
```

#### Trees (Learning Paths)
```http
GET    /api/v1/trees                    # List with graph metrics
POST   /api/v1/trees                    # Create tree structure
GET    /api/v1/trees/:id                # Get with node relationships
POST   /api/v1/trees/:id/fork           # Fork with graph inheritance
GET    /api/v1/trees/:id/analytics      # Graph analytics
```

#### Graph Operations
```http
GET    /api/v1/graph/visualization/:id  # Get visualization data
POST   /api/v1/graph/search             # Graph search queries
GET    /api/v1/graph/analytics          # Network analysis
POST   /api/v1/graph/relationships      # Create/update relationships
```

#### AI-Powered Features
```http
POST /api/v1/ai/embeddings              # Generate embeddings
POST /api/v1/ai/recommendations         # Get personalized suggestions
POST /api/v1/ai/similarity              # Find similar content
GET  /api/v1/ai/analytics               # AI usage statistics
```

### WebSocket Events
```javascript
// Real-time graph updates
{
  "type": "graph_update",
  "data": {
    "action": "node_created",
    "node_id": "node-456",
    "relationships": [
      {"type": "CREATED", "from": "user-123", "to": "node-456"}
    ]
  }
}

// Live collaboration
{
  "type": "collaboration",
  "data": {
    "user": "user-789",
    "action": "editing_node",
    "node_id": "node-456",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

## 🔧 Development

### Neo4j Schema Management
```bash
# Initialize constraints and indexes
go run cmd/server/main.go init-schema

# Seed with sample data
go run cmd/server/main.go seed-data

# Backup graph data
docker exec hi-neo4j neo4j-admin database dump --database=knowledge_graph
```

### Graph Query Development
```cypher
-- Find learning paths for a topic
MATCH path = (start:Node {title: "Machine Learning"})-[:CONNECTED_TO*1..3]->(end:Node)
WHERE end.content_type = "tutorial"
RETURN path, length(path) as depth
ORDER BY depth ASC

-- Analyze user collaboration networks
MATCH (u1:User)-[:FOLLOWS]->(u2:User)
WHERE u1.id = $user_id
OPTIONAL MATCH (u1)-[:CREATED]->(n:Node)<-[:LIKES]-(u2)
RETURN u2, count(n) as shared_interests
ORDER BY shared_interests DESC
```

### AI Service Development
```bash
cd ai-service

# Update embeddings for all nodes
curl -X POST http://localhost:8000/embeddings/batch-update

# Get similarity recommendations
curl -X POST http://localhost:8000/search/similarity \
  -H "Content-Type: application/json" \
  -d '{"query_text": "graph databases", "limit": 5}'

# Check embedding statistics
curl http://localhost:8000/analytics/embeddings-stats
```

### Frontend Graph Components
```svelte
<!-- GraphVisualization.svelte -->
<script>
  import { onMount } from 'svelte';
  import * as d3 from 'd3';
  
  export let graphData;
  export let width = 800;
  export let height = 600;
  
  let svg;
  
  onMount(() => {
    const simulation = d3.forceSimulation(graphData.nodes)
      .force('link', d3.forceLink(graphData.links).id(d => d.id))
      .force('charge', d3.forceManyBody().strength(-300))
      .force('center', d3.forceCenter(width / 2, height / 2));
    
    // Render graph...
  });
</script>

<div class="graph-container">
  <svg bind:this={svg} {width} {height}></svg>
</div>
```

## 🚢 Deployment

### Docker Production
```bash
# Build production images
docker-compose -f docker-compose.neo4j.yml build

# Deploy with all services
docker-compose -f docker-compose.neo4j.yml --profile full up -d

# Deploy minimal stack
docker-compose -f docker-compose.neo4j.yml up -d neo4j backend frontend ai-service
```

### Environment Configuration
```bash
# Production environment variables
NEO4J_URI=bolt://neo4j:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=your-secure-password
NEO4J_DATABASE=knowledge_graph

# AI Service configuration  
EMBEDDINGS_MODEL=sentence-transformers/all-MiniLM-L6-v2
EMBEDDINGS_DIMENSION=384
MIN_SIMILARITY_THRESHOLD=0.7

# Frontend configuration
VITE_API_URL=https://api.yourdomain.com
VITE_WS_URL=wss://api.yourdomain.com
VITE_NEO4J_BROWSER_URL=https://neo4j.yourdomain.com
```

### Neo4j Production Setup
```bash
# Production Neo4j configuration
echo "dbms.memory.heap.initial_size=2G" >> neo4j.conf
echo "dbms.memory.heap.max_size=4G" >> neo4j.conf
echo "dbms.memory.pagecache.size=2G" >> neo4j.conf
echo "dbms.security.auth_enabled=true" >> neo4j.conf
echo "dbms.connector.bolt.listen_address=0.0.0.0:7687" >> neo4j.conf
```

## 📊 Monitoring & Analytics

### Neo4j Metrics
```cypher
-- Database health check
CALL dbms.queryJmx("org.neo4j:instance=kernel#0,name=Store file sizes")

-- Performance monitoring
CALL dbms.queryJmx("org.neo4j:instance=kernel#0,name=Transactions")
```

### Grafana Dashboards
- **Graph Database Metrics**: Node/relationship counts, query performance
- **User Activity**: Creation patterns, collaboration networks
- **AI Service Performance**: Embedding generation, recommendation accuracy
- **Application Health**: API response times, WebSocket connections

## 🧪 Testing

### Graph Database Tests
```bash
# Backend tests with Neo4j testcontainers
cd backend
go test ./... -tags=integration

# Test specific graph operations
go test ./internal/database -v -run TestNeo4jService
```

### AI Service Tests
```bash
cd ai-service
pytest tests/ -v
pytest tests/test_embeddings.py -v
pytest tests/test_recommendations.py -v
```

### Frontend Graph Tests
```bash
cd frontend
npm run test
npm run test:e2e
npm run test:graph-components
```

## 🤝 Contributing

### Development Workflow
1. **Fork** the repository
2. **Create feature branch**: `git checkout -b feature/graph-analytics`
3. **Add graph tests**: Ensure your changes work with Neo4j
4. **Update documentation**: Include Cypher queries and graph examples
5. **Test thoroughly**: Run integration tests with Neo4j
6. **Submit PR**: Include graph schema changes if applicable

### Code Style
- **Go**: Use `gofmt`, `golint`, follow graph service patterns
- **TypeScript/Svelte**: Use Prettier, ESLint, type graph data properly
- **Python**: Use `black`, `flake8`, type hints for AI functions
- **Cypher**: Follow Neo4j style guide, use parameterized queries
- **Graph Schema**: Document all node labels and relationship types

## 📖 Learning Resources

### Neo4j & Graph Databases
- [Neo4j Graph Academy](https://graphacademy.neo4j.com/)
- [Cypher Query Language](https://neo4j.com/developer/cypher/)
- [Graph Data Science](https://neo4j.com/docs/graph-data-science/)

### D3.js Graph Visualization
- [D3.js Force Layout](https://github.com/d3/d3-force)
- [Observable Graph Examples](https://observablehq.com/@d3/force-directed-graph)
- [Cytoscape.js Documentation](https://js.cytoscape.org/)

## 🔗 Links

- **Demo**: [https://hi-graph.demo.com](https://hi-graph.demo.com)
- **Neo4j Browser**: [http://localhost:7474](http://localhost:7474)
- **API Documentation**: [http://localhost:8080/swagger](http://localhost:8080/swagger)
- **AI Service Docs**: [http://localhost:8000/docs](http://localhost:8000/docs)
- **Discord**: [https://discord.gg/hi-graph](https://discord.gg/hi-graph)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎉 Acknowledgments

- **Neo4j Community** for the excellent graph database
- **D3.js Community** for powerful visualization tools
- **SvelteKit Team** for the fantastic frontend framework
- **Sentence Transformers** for semantic embeddings
- **Vaporwave Artists** for the aesthetic inspiration

---

<div align="center">

**Made with 🧠 by humans, augmented by AI, powered by graphs**

*Human Intelligence Neo4j Platform © 2024*

</div>