# 🌳 Human Intelligence - Learning Trees

> **Where knowledge grows naturally** - A Neo4j-powered platform for cultivating learning through interconnected trees of knowledge

[![License: MIT](https://img.shields.io/badge/License-MIT-neon.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Neo4j](https://img.shields.io/badge/Neo4j-5.15+-008CC1.svg)](https://neo4j.com/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6.svg)](https://typescriptlang.org/)

## 🌱 Vision: Growing Knowledge Gardens

Human Intelligence reimagines learning as **growing a garden of knowledge trees**. Each learner cultivates their own collection of trees, where every tree represents a focused learning journey.

### 🔍 Two Views of Knowledge

**🌳 Terrestrial View (Trees & Nodes)**
- **Trees**: Focused learning paths - trains of thought that branch and grow
- **Nodes**: Individual learning steps, experiments, insights, or discoveries
- **Branching**: Natural divergence when exploration leads to new directions
- **Growth**: Trees expand organically as learning progresses

**🌌 Cosmic View (Tags & Forests)**  
- **Tags**: Semantic connectors that group related trees across the knowledge space
- **Forests**: Communities that form around shared tags and interests
- **Connections**: Trees relate to each other through overlapping tags and concepts

### 📚 Learning Journey Example

```
🌳 "Learning C Programming" Tree:
Root → Hello World (print to screen)
  ├── Variables & Data Types
  ├── Control Flow (if/else, loops)
  │   ├── FizzBuzz Implementation
  │   └── Simple Calculator
  ├── Functions & Scope
  ├── Pointers & Memory
  │   ├── Malloc/Free Experiments
  │   ├── Memory Leak Debugging
  │   └── Custom Allocator Research
  └── Building a Todo CLI App

🏷️ Tags: #programming, #c, #systems, #memory-management, #cli-tools
```

## 🛠️ Tech Stack

### Core Architecture
- **Database**: Neo4j 5.15+ (Native graph database for knowledge relationships)
- **Backend**: Go 1.21+ with Gin framework (Fast, concurrent API server)
- **Frontend**: SvelteKit + TypeScript (Reactive UI for knowledge visualization)
- **Visualization**: D3.js (Interactive graph rendering)
- **Cache**: Redis (Session management and performance)

### Graph Data Model
```cypher
// Core entities and relationships
(:User)-[:CREATED]->(:Tree)-[:CONTAINS]->(:Node)
(:Node)-[:BRANCHES_TO]->(:Node)
(:Tree)-[:TAGGED_WITH]->(:Tag)
(:User)-[:FOLLOWS]->(:User)
(:User)-[:EXPLORES]->(:Tree)
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose

### 🌟 Development Scripts

We've built convenient development tools for a smooth experience:

```bash
# Start everything (Neo4j + Backend + Frontend)
./dev.sh start

# Start just the backend
./dev.sh backend

# Check what's running
./dev.sh status

# Stop everything
./dev.sh stop

# Quick manual start
./quick-start.sh
```

### 🎯 Step-by-Step Setup

1. **Clone and Navigate**
   ```bash
   git clone https://github.com/your-org/human-intelligence.git
   cd human-intelligence
   ```

2. **Start the Knowledge Platform**
   ```bash
   ./dev.sh start
   ```

3. **Access Your Garden**
   - **Frontend**: http://localhost:3000
   - **API**: http://localhost:8085
   - **Neo4j Browser**: http://localhost:7474 (neo4j/hi_password)
   - **Test Page**: `open test-platform.html`

## 🌿 MVP Features

### ✅ Foundation Complete
- **Graph Database**: Neo4j with optimized schema for learning trees
- **User Management**: Create accounts and manage learning profiles
- **REST API**: Complete CRUD operations for all entities
- **Graph Visualization**: D3.js powered knowledge maps
- **Development Tools**: Streamlined development workflow

### 🎯 Next: Core Learning Features
- **Tree Creation**: Start new learning journeys
- **Node Management**: Add insights, experiments, and discoveries
- **Branching Logic**: Handle natural learning divergence
- **Tag System**: Semantic organization and discovery
- **Tree Exploration**: Browse and fork others' learning paths

### 🔮 Future Growth
- **Forest Communities**: Collaborative learning groups
- **Learning Analytics**: Progress tracking and insights
- **AI Suggestions**: Intelligent path recommendations
- **Real-time Collaboration**: Live learning sessions
- **Mobile Experience**: Native apps for learning on-the-go

## 🔧 Development

### API Endpoints

```bash
# Trees (Learning Paths)
POST   /api/v1/trees              # Create new learning tree
GET    /api/v1/trees/:id          # Get tree with all nodes
PUT    /api/v1/trees/:id          # Update tree metadata
DELETE /api/v1/trees/:id          # Delete tree

# Nodes (Learning Steps)
POST   /api/v1/nodes              # Create new learning node
GET    /api/v1/nodes/:id          # Get node details
PUT    /api/v1/nodes/:id          # Update node content
DELETE /api/v1/nodes/:id          # Delete node
POST   /api/v1/trees/:treeId/nodes/:nodeId  # Add node to tree

# Tags (Semantic Connectors)
POST   /api/v1/tags               # Create new tag
GET    /api/v1/tags/popular       # Get popular tags
POST   /api/v1/nodes/:nodeId/tags/:tagName  # Tag a node

# Discovery & Visualization
GET    /api/v1/search?q=:query    # Search trees and nodes
GET    /api/v1/graph/visualization/:userId  # Get graph data for D3.js
```

### Graph Queries

```cypher
-- Find a user's learning forest
MATCH (u:User {id: $userId})-[:CREATED]->(t:Tree)
OPTIONAL MATCH (t)-[:CONTAINS]->(n:Node)
RETURN u, t, collect(n) as nodes

-- Discover related trees through tags
MATCH (t1:Tree)-[:TAGGED_WITH]->(tag:Tag)<-[:TAGGED_WITH]-(t2:Tree)
WHERE t1.id = $treeId AND t1 <> t2
RETURN t2, tag, count(*) as shared_tags
ORDER BY shared_tags DESC

-- Find branching paths in a tree
MATCH (tree:Tree)-[:CONTAINS]->(root:Node)
WHERE tree.id = $treeId AND NOT (root)<-[:BRANCHES_TO]-()
MATCH path = (root)-[:BRANCHES_TO*]->(leaf:Node)
WHERE NOT (leaf)-[:BRANCHES_TO]->()
RETURN path
```

### Development Commands

```bash
# Start development environment
./dev.sh start                    # Full stack
./dev.sh backend                  # Backend only
./dev.sh services                 # Databases only

# Check status
./dev.sh status                   # Show running services
./dev.sh logs                     # View logs

# Quick actions
./dev.sh test                     # Open test page
./dev.sh neo4j                    # Open Neo4j browser
./dev.sh stop                     # Stop everything
```

## 🧪 Testing Your Garden

Open the test interface to explore the platform:

```bash
# Open test page
./dev.sh test

# Or manually
open test-platform.html
```

Test key features:
- ✅ Create users and trees
- ✅ Add nodes to learning paths  
- ✅ Apply semantic tags
- ✅ Visualize knowledge graphs
- ✅ Search across the forest

## 🗂️ Project Structure

```
human-intelligence/
├── backend/                 # Go API server
│   ├── cmd/server/         # Application entry point
│   ├── internal/models/    # Graph data models
│   └── internal/database/  # Neo4j integration
├── frontend/               # SvelteKit interface
│   ├── src/lib/components/ # UI components
│   └── src/routes/         # Application pages
├── scripts/dev/            # Development tools
│   ├── start-dev.sh       # Full stack startup
│   ├── start-backend.sh   # Backend only
│   └── stop-dev.sh        # Graceful shutdown
├── dev.sh                 # Main development CLI
├── quick-start.sh         # Simple manual start
└── test-platform.html    # Development test interface
```

## 🤝 Contributing

We're building the future of collaborative learning! 

### Development Workflow
1. **Start your environment**: `./dev.sh start`
2. **Make your changes**: Edit code with hot reloading
3. **Test your garden**: Use the test interface
4. **Check the graph**: Explore relationships in Neo4j Browser
5. **Submit your growth**: Create a pull request

### Code Philosophy
- **Graph-First**: Design with relationships in mind
- **Learning-Centered**: Every feature should enhance the learning experience
- **Natural Growth**: Features should feel organic, not forced
- **Community Focused**: Enable collaboration without compromising personal learning

## 🎯 Vision: The Learning Future

We're creating a platform where:
- **Learning is Visual**: See your knowledge grow like a living garden
- **Connections Emerge**: Discover relationships between different areas of study
- **Communities Form**: Find others on similar learning journeys
- **Knowledge Persists**: Build a permanent record of your intellectual growth
- **Wisdom Spreads**: Share insights that help others learn faster

**Join us in growing the forest of human knowledge.** 🌳🌲🌴

## 📜 License

MIT License - see [LICENSE](LICENSE) for details.

---

<div align="center">

**Made with 🧠 for learners, by learners**

*Human Intelligence © 2024*

</div>