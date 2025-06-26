# 🧠 Human Intelligence (HI!)

> **Where AI augments HI** - A factorial multiplier to human intelligence

[![License: MIT](https://img.shields.io/badge/License-MIT-neon.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Node Version](https://img.shields.io/badge/Node-18+-339933.svg)](https://nodejs.org/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-4.0+-FF3E00.svg)](https://kit.svelte.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://postgresql.org/)

## 🌟 Vision

Human Intelligence is a **collaborative learning platform** that revolutionizes how knowledge is shared and acquired. Think of it as the love child of GitHub, Stack Overflow, and a modern forum - but specifically designed for learning with AI augmentation.

### Core Concept
- **Nodes**: Atomic learning units (text, video, code, quizzes, etc.)
- **Learning Paths**: Git-like sequences of nodes that can be forked, merged, and evolved
- **TLDRs**: 140-character summaries users create after completing paths
- **Live Forum**: Real-time visualization of learning activity across the platform
- **AI Integration**: LocalAI for free tier, premium APIs for paid features

## 🎨 Aesthetic

**Zoomer shitpost meets boomer messageboard with hacker flair**

- 🌈 **Vaporwave** color palette (neon pink, cyan, purple)
- 🏛️ **Greek marble** textures for that classical knowledge vibe
- 💻 **Terminal green** for the hacker aesthetic
- ⚡ **Real-time animations** showing learning activity
- 🎮 **Retro grid** backgrounds and neon glows

## 🚀 Features

### MVP Features
- ✅ User authentication and profiles
- ✅ Create and manage learning nodes
- ✅ Build learning paths with git-like branching
- ✅ Fork and merge learning paths
- ✅ TLDR summaries and community voting
- ✅ Real-time activity feed
- ✅ Basic text search
- ✅ Responsive vaporwave UI

### Planned Features
- 🔮 **AI-Powered Search**: Semantic search with embeddings
- 🤖 **AI Recommendations**: Personalized learning suggestions
- 🧩 **Plugin System**: Custom node types (LaTeX, Jupyter, Wolfram, etc.)
- 👥 **Collaboration Tools**: Real-time editing and discussion
- 🏫 **Classroom Mode**: Teacher tools and grading
- 📊 **Analytics**: Learning progress and insights
- 🌐 **Internationalization**: Multi-language support
- 📱 **Mobile App**: Native iOS/Android apps

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP) + Gorilla WebSocket
- **Database**: PostgreSQL 15+ with JSONB
- **Search**: Elasticsearch/Meilisearch (configurable)
- **Cache**: Redis (optional)
- **AI**: LocalAI (free) + OpenAI/Anthropic (premium)

### Frontend
- **Framework**: SvelteKit 4.0+
- **Styling**: Tailwind CSS with custom vaporwave theme
- **Real-time**: WebSocket client
- **Charts**: D3.js for learning path visualizations
- **Animations**: Custom CSS animations + Svelte transitions

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Deployment**: Kubernetes manifests included
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana dashboards

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   SvelteKit     │    │   Go Backend    │    │   PostgreSQL    │
│   Frontend      │◄──►│   API Server    │◄──►│   Database      │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WebSocket     │    │   LocalAI       │    │   Redis Cache   │
│   Real-time     │    │   AI Services   │    │   (Optional)    │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚦 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose (recommended)

### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-org/human-intelligence.git
cd human-intelligence

# Start all services
docker-compose up -d

# The app will be available at:
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Database: localhost:5432
```

### Option 2: Manual Setup

#### Backend Setup

```bash
cd backend

# Copy environment config
cp .env.example .env
# Edit .env with your database credentials

# Install dependencies
go mod download

# Run database migrations
go run cmd/server/main.go migrate

# Start the server
go run cmd/server/main.go
```

#### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

#### Database Setup

```sql
-- Create database
CREATE DATABASE human_intelligence;

-- Create user (optional)
CREATE USER hi_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE human_intelligence TO hi_user;
```

## 🔧 Development

### Backend Development

```bash
# Run with hot reload
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Generate API docs
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go

# Database migrations
go run cmd/server/main.go migrate
go run cmd/server/main.go rollback
```

### Frontend Development

```bash
# Development server
npm run dev

# Type checking
npm run check

# Linting
npm run lint

# Build for production
npm run build

# Preview production build
npm run preview
```

### Database Migrations

```bash
# Create new migration
touch backend/migrations/002_add_feature.sql

# Apply migrations
go run cmd/server/main.go migrate

# Check migration status
go run cmd/server/main.go migrate-status
```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
go test -v ./internal/handlers/
go test -race ./...
```

### Frontend Tests
```bash
cd frontend
npm run test
npm run test:watch
npm run test:ui
```

### E2E Tests
```bash
cd frontend
npm run test:e2e
```

## 📚 API Documentation

### Authentication
```http
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
```

### Nodes
```http
GET    /api/v1/nodes          # List public nodes
POST   /api/v1/nodes          # Create node
GET    /api/v1/nodes/:id      # Get node
PUT    /api/v1/nodes/:id      # Update node
DELETE /api/v1/nodes/:id      # Delete node
POST   /api/v1/nodes/:id/vote # Vote on node
```

### Learning Paths
```http
GET    /api/v1/paths               # List public paths
POST   /api/v1/paths               # Create path
GET    /api/v1/paths/:id           # Get path
PUT    /api/v1/paths/:id           # Update path
DELETE /api/v1/paths/:id           # Delete path
POST   /api/v1/paths/:id/fork      # Fork path
POST   /api/v1/paths/:id/complete  # Mark as completed
POST   /api/v1/paths/:id/tldr      # Create TLDR
```

### Search
```http
GET  /api/v1/search?q=query        # Text search
POST /api/v1/search/semantic       # Semantic search
GET  /api/v1/search/suggestions    # AI recommendations
```

### WebSocket Events
```javascript
// Real-time activity updates
{
  "type": "activity",
  "data": {
    "user": "username",
    "action": "node_completed",
    "target": "node_id",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

## 🎯 Environment Variables

### Backend (.env)
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=human_intelligence

# JWT
JWT_SECRET=your-secret-key

# AI Services
LOCALAI_ENABLED=true
LOCALAI_BASE_URL=http://localhost:8080
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=...

# Features
REDIS_ENABLED=false
SEARCH_ENGINE=simple
PLUGINS_ENABLED=true
```

### Frontend (.env)
```bash
# API Configuration
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080

# Feature Flags
VITE_ENABLE_AI_FEATURES=true
VITE_ENABLE_REAL_TIME=true
```

## 🚢 Deployment

### Docker Production
```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes
```bash
# Apply manifests
kubectl apply -f k8s/

# Check deployment
kubectl get pods -n human-intelligence
```

### Manual Deployment
```bash
# Backend
cd backend
CGO_ENABLED=0 GOOS=linux go build -o hi-server cmd/server/main.go

# Frontend
cd frontend
npm run build
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Run the test suite (`go test ./...` and `npm test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style
- **Go**: Use `gofmt` and `golint`
- **JavaScript/Svelte**: Use Prettier and ESLint
- **CSS**: Follow BEM methodology
- **Commits**: Use conventional commits

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎉 Acknowledgments

- Inspired by the collaborative nature of GitHub
- Learning from the community wisdom of Stack Overflow
- Aesthetic influenced by the vaporwave and synthwave movements
- Built for the next generation of learners

## 🔗 Links

- **Demo**: [https://hi.demo.com](https://hi.demo.com)
- **Documentation**: [https://docs.hi.com](https://docs.hi.com)
- **Discord**: [https://discord.gg/hi](https://discord.gg/hi)
- **Twitter**: [@HumanIntelHQ](https://twitter.com/HumanIntelHQ)

---

<div align="center">

**Made with 🧠 by humans, augmented by AI**

*Human Intelligence © 2024*

</div>