#!/bin/bash

# Human Intelligence Neo4j Knowledge Platform Setup Script
# This script helps you set up the development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.neo4j.yml"
PROJECT_NAME="human-intelligence"
NEO4J_PASSWORD="hi_password"
NEO4J_USER="neo4j"

# Function to print colored output
print_color() {
    printf "${1}${2}${NC}\n"
}

print_header() {
    echo ""
    print_color $PURPLE "════════════════════════════════════════════════════════════════"
    print_color $PURPLE "  🧠 Human Intelligence - Neo4j Knowledge Platform Setup"
    print_color $PURPLE "════════════════════════════════════════════════════════════════"
    echo ""
}

print_step() {
    print_color $CYAN "→ $1"
}

print_success() {
    print_color $GREEN "✅ $1"
}

print_warning() {
    print_color $YELLOW "⚠️  $1"
}

print_error() {
    print_color $RED "❌ $1"
}

print_info() {
    print_color $BLUE "ℹ️  $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_step "Checking prerequisites..."

    local missing_deps=()

    # Check Docker
    if ! command_exists docker; then
        missing_deps+=("docker")
    fi

    # Check Docker Compose
    if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
        missing_deps+=("docker-compose")
    fi

    # Check Go (optional for development)
    if ! command_exists go; then
        print_warning "Go is not installed (required for backend development)"
    else
        print_info "Go version: $(go version | cut -d' ' -f3)"
    fi

    # Check Node.js (optional for development)
    if ! command_exists node; then
        print_warning "Node.js is not installed (required for frontend development)"
    else
        print_info "Node.js version: $(node --version)"
    fi

    # Check Python (optional for AI service development)
    if ! command_exists python3; then
        print_warning "Python 3 is not installed (required for AI service development)"
    else
        print_info "Python version: $(python3 --version)"
    fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        echo ""
        print_info "Please install the missing dependencies:"
        echo ""
        print_color $YELLOW "  Docker: https://docs.docker.com/get-docker/"
        print_color $YELLOW "  Docker Compose: https://docs.docker.com/compose/install/"
        echo ""
        exit 1
    fi

    print_success "All required dependencies are installed"
}

# Function to setup environment files
setup_environment() {
    print_step "Setting up environment files..."

    # Backend environment
    if [ ! -f "backend/.env" ]; then
        print_step "Creating backend/.env file..."
        cat > backend/.env << EOF
# Neo4j Configuration
NEO4J_URI=bolt://neo4j:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=hi_password
NEO4J_DATABASE=knowledgegraph

# Application Configuration
GIN_MODE=debug
PORT=8080
APP_VERSION=dev
JWT_SECRET=dev-secret-key-change-in-production
JWT_EXPIRY_HOURS=24

# Redis Configuration
REDIS_ENABLED=true
REDIS_URL=redis://redis:6379

# AI Service Configuration
PYTHON_AI_SERVICE_URL=http://ai-service:8000
EMBEDDINGS_ENABLED=true
RECOMMENDATIONS_ENABLED=true

# Feature Flags
GRAPH_VISUALIZATIONS_ENABLED=true
REAL_TIME_ENABLED=true
SEARCH_ENABLED=true
TAGGING_ENABLED=true
DEBUG_ENABLED=true
SEED_DATABASE=true
EOF
        print_success "Created backend/.env"
    else
        print_info "backend/.env already exists"
    fi

    # Frontend environment
    if [ ! -f "frontend/.env" ]; then
        print_step "Creating frontend/.env file..."
        cat > frontend/.env << EOF
# API Configuration
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_NEO4J_BROWSER_URL=http://localhost:7474
VITE_AI_SERVICE_URL=http://localhost:8000

# Feature Flags
VITE_ENABLE_GRAPH_VISUALIZATIONS=true
VITE_ENABLE_AI_FEATURES=true
VITE_ENABLE_REAL_TIME=true
VITE_ENABLE_ADVANCED_SEARCH=true
VITE_ENABLE_TAGGING=true

# Graph Visualization Settings
VITE_MAX_GRAPH_NODES=1000
VITE_DEFAULT_GRAPH_LAYOUT=force
VITE_ENABLE_3D_GRAPH=false
EOF
        print_success "Created frontend/.env"
    else
        print_info "frontend/.env already exists"
    fi

    # AI Service environment
    if [ ! -f "ai-service/.env" ]; then
        print_step "Creating ai-service/.env file..."
        cat > ai-service/.env << EOF
# Neo4j Configuration
NEO4J_URI=bolt://neo4j:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=hi_password
NEO4J_DATABASE=knowledgegraph

# AI/ML Configuration
EMBEDDINGS_MODEL=sentence-transformers/all-MiniLM-L6-v2
EMBEDDINGS_DIMENSION=384
BATCH_SIZE=32
MIN_SIMILARITY_THRESHOLD=0.7
MAX_RECOMMENDATIONS=10

# Redis Configuration
REDIS_URL=redis://redis:6379
CACHE_TTL=3600

# FastAPI Configuration
FASTAPI_ENV=development
PORT=8000
EOF
        print_success "Created ai-service/.env"
    else
        print_info "ai-service/.env already exists"
    fi
}

# Function to pull Docker images
pull_images() {
    print_step "Pulling Docker images..."

    if [ -f "$COMPOSE_FILE" ]; then
        docker-compose -f "$COMPOSE_FILE" pull
        print_success "Docker images pulled successfully"
    else
        print_error "Docker Compose file not found: $COMPOSE_FILE"
        exit 1
    fi
}

# Function to start services
start_services() {
    print_step "Starting services..."

    # Start core services first
    print_step "Starting Neo4j and Redis..."
    docker-compose -f "$COMPOSE_FILE" up -d neo4j redis

    # Wait for Neo4j to be ready
    print_step "Waiting for Neo4j to be ready..."
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if docker exec hi-neo4j cypher-shell -u neo4j -p hi_password "RETURN 1" >/dev/null 2>&1; then
            print_success "Neo4j is ready"
            break
        fi

        if [ $attempt -eq $max_attempts ]; then
            print_error "Neo4j failed to start after $max_attempts attempts"
            exit 1
        fi

        print_info "Attempt $attempt/$max_attempts - Neo4j not ready yet..."
        sleep 5
        ((attempt++))
    done

    # Start AI service
    print_step "Starting AI service..."
    docker-compose -f "$COMPOSE_FILE" up -d ai-service

    # Wait for AI service to be ready
    print_step "Waiting for AI service to be ready..."
    attempt=1
    while [ $attempt -le 20 ]; do
        if curl -s http://localhost:8000/health >/dev/null 2>&1; then
            print_success "AI service is ready"
            break
        fi

        if [ $attempt -eq 20 ]; then
            print_warning "AI service may not be ready yet, but continuing..."
            break
        fi

        sleep 3
        ((attempt++))
    done

    # Start backend
    print_step "Starting backend..."
    docker-compose -f "$COMPOSE_FILE" up -d backend

    # Start frontend
    print_step "Starting frontend..."
    docker-compose -f "$COMPOSE_FILE" up -d frontend

    print_success "All services started successfully"
}

# Function to initialize database
init_database() {
    print_step "Initializing database schema..."

    # Wait for backend to be ready
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8080/health >/dev/null 2>&1; then
            print_success "Backend is ready"
            break
        fi

        if [ $attempt -eq $max_attempts ]; then
            print_error "Backend failed to start after $max_attempts attempts"
            exit 1
        fi

        print_info "Attempt $attempt/$max_attempts - Backend not ready yet..."
        sleep 5
        ((attempt++))
    done

    # Initialize schema (this should be done automatically by the backend)
    print_info "Schema initialization is handled by the backend service"
    print_success "Database initialized"
}

# Function to seed sample data
seed_data() {
    print_step "Seeding sample data..."

    # Create sample data via Cypher
    docker exec hi-neo4j cypher-shell -u neo4j -p hi_password << 'EOF'
// Create sample users
CREATE (:User {
    id: 'user-1',
    username: 'alice',
    email: 'alice@example.com',
    first_name: 'Alice',
    last_name: 'Smith',
    bio: 'AI researcher and knowledge enthusiast',
    is_active: true,
    is_verified: true,
    created_at: datetime(),
    updated_at: datetime(),
    followers_count: 0,
    following_count: 0,
    nodes_count: 0,
    trees_count: 0,
    reputation_score: 100
});

CREATE (:User {
    id: 'user-2',
    username: 'bob',
    email: 'bob@example.com',
    first_name: 'Bob',
    last_name: 'Johnson',
    bio: 'Software engineer and graph database enthusiast',
    is_active: true,
    is_verified: true,
    created_at: datetime(),
    updated_at: datetime(),
    followers_count: 0,
    following_count: 0,
    nodes_count: 0,
    trees_count: 0,
    reputation_score: 85
});

// Create sample tags
CREATE (:Tag {
    id: 'tag-1',
    name: 'machine-learning',
    color: '#3B82F6',
    created_at: datetime(),
    usage_count: 0
});

CREATE (:Tag {
    id: 'tag-2',
    name: 'graph-databases',
    color: '#10B981',
    created_at: datetime(),
    usage_count: 0
});

CREATE (:Tag {
    id: 'tag-3',
    name: 'neo4j',
    color: '#008CC1',
    created_at: datetime(),
    usage_count: 0
});

// Create sample nodes
CREATE (:Node {
    id: 'node-1',
    title: 'Introduction to Graph Databases',
    content: 'Graph databases are a type of NoSQL database that uses graph structures with nodes, edges, and properties to represent and store data.',
    content_type: 'text',
    description: 'A beginner-friendly introduction to graph databases',
    is_public: true,
    is_published: true,
    created_at: datetime(),
    updated_at: datetime(),
    views_count: 0,
    likes_count: 0,
    shares_count: 0,
    owner_id: 'user-1',
    level: 1,
    position: 1
});

CREATE (:Node {
    id: 'node-2',
    title: 'Neo4j Cypher Basics',
    content: 'Cypher is Neo4j\'s declarative query language. It allows you to express what you want to select, insert, update, or delete from your graph data.',
    content_type: 'text',
    description: 'Learn the basics of Cypher query language',
    is_public: true,
    is_published: true,
    created_at: datetime(),
    updated_at: datetime(),
    views_count: 0,
    likes_count: 0,
    shares_count: 0,
    owner_id: 'user-2',
    level: 2,
    position: 2
});

// Create sample tree
CREATE (:Tree {
    id: 'tree-1',
    name: 'Graph Database Learning Path',
    description: 'A comprehensive learning path for graph databases',
    is_public: true,
    is_template: false,
    created_at: datetime(),
    updated_at: datetime(),
    views_count: 0,
    likes_count: 0,
    forks_count: 0,
    nodes_count: 2,
    owner_id: 'user-1'
});

// Create relationships
MATCH (u1:User {id: 'user-1'}), (n1:Node {id: 'node-1'})
CREATE (u1)-[:CREATED]->(n1);

MATCH (u2:User {id: 'user-2'}), (n2:Node {id: 'node-2'})
CREATE (u2)-[:CREATED]->(n2);

MATCH (t:Tree {id: 'tree-1'}), (n1:Node {id: 'node-1'})
CREATE (t)-[:CONTAINS {position: 1, added_at: datetime()}]->(n1);

MATCH (t:Tree {id: 'tree-1'}), (n2:Node {id: 'node-2'})
CREATE (t)-[:CONTAINS {position: 2, added_at: datetime()}]->(n2);

MATCH (u1:User {id: 'user-1'}), (t:Tree {id: 'tree-1'})
CREATE (u1)-[:CREATED]->(t);

MATCH (n1:Node {id: 'node-1'}), (tag1:Tag {name: 'graph-databases'})
CREATE (n1)-[:TAGGED {tagged_at: datetime()}]->(tag1);

MATCH (n1:Node {id: 'node-1'}), (tag3:Tag {name: 'neo4j'})
CREATE (n1)-[:TAGGED {tagged_at: datetime()}]->(tag3);

MATCH (n2:Node {id: 'node-2'}), (tag3:Tag {name: 'neo4j'})
CREATE (n2)-[:TAGGED {tagged_at: datetime()}]->(tag3);

MATCH (n1:Node {id: 'node-1'}), (n2:Node {id: 'node-2'})
CREATE (n1)-[:RELATED_TO {strength: 0.8, type: 'conceptual'}]->(n2);

MATCH (u1:User {id: 'user-1'}), (u2:User {id: 'user-2'})
CREATE (u1)-[:FOLLOWS {followed_at: datetime()}]->(u2);

RETURN 'Sample data created successfully' as result;
EOF

    print_success "Sample data seeded successfully"
}

# Function to show service status
show_status() {
    print_step "Checking service status..."
    echo ""

    # Check Neo4j
    if docker exec hi-neo4j cypher-shell -u neo4j -p hi_password "RETURN 1" >/dev/null 2>&1; then
        print_success "Neo4j: Running (http://localhost:7474)"
    else
        print_error "Neo4j: Not running"
    fi

    # Check Redis
    if docker exec hi-redis redis-cli ping >/dev/null 2>&1; then
        print_success "Redis: Running"
    else
        print_error "Redis: Not running"
    fi

    # Check Backend
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        print_success "Backend API: Running (http://localhost:8080)"
    else
        print_error "Backend API: Not running"
    fi

    # Check AI Service
    if curl -s http://localhost:8000/health >/dev/null 2>&1; then
        print_success "AI Service: Running (http://localhost:8000)"
    else
        print_error "AI Service: Not running"
    fi

    # Check Frontend
    if curl -s http://localhost:3000 >/dev/null 2>&1; then
        print_success "Frontend: Running (http://localhost:3000)"
    else
        print_error "Frontend: Not running"
    fi
}

# Function to show useful information
show_info() {
    echo ""
    print_color $PURPLE "════════════════════════════════════════════════════════════════"
    print_color $PURPLE "  🎉 Setup Complete! Here's what you can do now:"
    print_color $PURPLE "════════════════════════════════════════════════════════════════"
    echo ""

    print_color $GREEN "🌐 Web Interfaces:"
    print_color $CYAN "  • Frontend:       http://localhost:3000"
    print_color $CYAN "  • Neo4j Browser:  http://localhost:7474 (neo4j/hi_password)"
    print_color $CYAN "  • Backend API:    http://localhost:8080"
    print_color $CYAN "  • AI Service:     http://localhost:8000"
    print_color $CYAN "  • API Docs:       http://localhost:8080/swagger"
    print_color $CYAN "  • AI Docs:        http://localhost:8000/docs"
    echo ""

    print_color $GREEN "🔧 Useful Commands:"
    print_color $YELLOW "  • View logs:      docker-compose -f $COMPOSE_FILE logs -f"
    print_color $YELLOW "  • Stop services:  docker-compose -f $COMPOSE_FILE down"
    print_color $YELLOW "  • Restart:        docker-compose -f $COMPOSE_FILE restart"
    print_color $YELLOW "  • Shell access:   docker exec -it hi-neo4j bash"
    print_color $YELLOW "  • Cypher shell:   docker exec -it hi-neo4j cypher-shell -u neo4j -p hi_password"
    echo ""

    print_color $GREEN "📊 Monitoring:"
    print_color $CYAN "  • Grafana:        http://localhost:3002 (admin/admin)"
    print_color $CYAN "  • Prometheus:     http://localhost:9091"
    echo ""

    print_color $GREEN "🧪 Sample Queries:"
    print_color $YELLOW "  • List all nodes: MATCH (n:Node) RETURN n LIMIT 10"
    print_color $YELLOW "  • User networks:  MATCH (u:User)-[:FOLLOWS]->(f:User) RETURN u, f"
    print_color $YELLOW "  • Tagged content: MATCH (n:Node)-[:TAGGED]->(t:Tag) RETURN n.title, t.name"
    echo ""

    print_color $GREEN "🤖 AI Features:"
    print_color $CYAN "  • Generate embeddings: curl -X POST http://localhost:8000/embeddings/batch-update"
    print_color $CYAN "  • Get recommendations: curl -X POST http://localhost:8000/recommendations/personalized"
    print_color $CYAN "  • Similarity search:   curl -X POST http://localhost:8000/search/similarity"
    echo ""

    print_color $PURPLE "════════════════════════════════════════════════════════════════"
    print_color $PURPLE "  Happy coding! 🚀"
    print_color $PURPLE "════════════════════════════════════════════════════════════════"
}

# Function to handle cleanup
cleanup() {
    print_step "Stopping services..."
    docker-compose -f "$COMPOSE_FILE" down
    print_success "Services stopped"
}

# Function to show help
show_help() {
    echo "Human Intelligence Neo4j Setup Script"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  --full          Full setup (default): check deps, setup env, start services, seed data"
    echo "  --quick         Quick setup: start services only"
    echo "  --env           Setup environment files only"
    echo "  --start         Start services only"
    echo "  --stop          Stop all services"
    echo "  --restart       Restart all services"
    echo "  --status        Show service status"
    echo "  --seed          Seed sample data"
    echo "  --logs          Show service logs"
    echo "  --clean         Clean up everything (remove containers and volumes)"
    echo "  --help          Show this help message"
    echo ""
}

# Main function
main() {
    case "${1:-}" in
        --help)
            show_help
            ;;
        --env)
            print_header
            setup_environment
            ;;
        --quick)
            print_header
            start_services
            show_status
            show_info
            ;;
        --start)
            print_header
            start_services
            show_status
            ;;
        --stop)
            print_header
            cleanup
            ;;
        --restart)
            print_header
            cleanup
            start_services
            show_status
            ;;
        --status)
            print_header
            show_status
            ;;
        --seed)
            print_header
            seed_data
            ;;
        --logs)
            docker-compose -f "$COMPOSE_FILE" logs -f
            ;;
        --clean)
            print_header
            print_warning "This will remove all containers and volumes. Are you sure? (y/N)"
            read -r response
            if [[ "$response" =~ ^[Yy]$ ]]; then
                docker-compose -f "$COMPOSE_FILE" down -v --remove-orphans
                docker system prune -f
                print_success "Cleanup completed"
            else
                print_info "Cleanup cancelled"
            fi
            ;;
        --full|"")
            print_header
            check_prerequisites
            setup_environment
            pull_images
            start_services
            init_database
            seed_data
            show_status
            show_info
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# Trap to handle script interruption
trap cleanup EXIT

# Run main function
main "$@"
