#!/bin/bash

# Neo4j Learning Trees Platform - Development CLI
# Simple script for MVP development

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}🌳 Neo4j Learning Trees - Development CLI${NC}\n"
}

print_status() {
    echo -e "${GREEN}Development Environment Status:${NC}\n"

    # Check Neo4j
    if docker-compose ps neo4j | grep -q "Up.*healthy"; then
        echo -e "✅ Neo4j running on port 7687/7474"
    else
        echo -e "❌ Neo4j not running"
    fi

    # Check Redis
    if docker-compose ps redis | grep -q "Up.*healthy"; then
        echo -e "✅ Redis running on port 6379"
    else
        echo -e "❌ Redis not running"
    fi

    # Check Backend
    if curl -s http://localhost:8085/health >/dev/null 2>&1; then
        echo -e "✅ Backend running on port 8085"
    else
        echo -e "❌ Backend not running"
    fi

    echo ""
}

start_services() {
    echo -e "${YELLOW}Starting Neo4j and Redis...${NC}"
    docker-compose up -d neo4j redis

    echo -e "${YELLOW}Waiting for services to be ready...${NC}"
    sleep 10

    echo -e "${YELLOW}Building backend...${NC}"
    cd backend
    go build -o server cmd/server/main.go

    echo -e "${GREEN}Services started! Use 'backend' command to start the server.${NC}"
    echo ""
    echo "Quick links:"
    echo "  • Neo4j Browser: http://localhost:7474 (neo4j/hi_password)"
    echo "  • Test Interface: open test-tree-api.html"
    echo ""
}

start_backend() {
    cd backend
    if [ ! -f "server" ]; then
        echo -e "${YELLOW}Building backend...${NC}"
        go build -o server cmd/server/main.go
    fi

    echo -e "${GREEN}Starting backend server on port 8085...${NC}"
    export PORT=8085
    export NEO4J_URI=bolt://localhost:7687
    export NEO4J_USERNAME=neo4j
    export NEO4J_PASSWORD=hi_password
    export NEO4J_DATABASE=knowledgegraph
    export GIN_MODE=debug

    ./server
}

stop_services() {
    echo -e "${YELLOW}Stopping all services...${NC}"
    docker-compose down
    pkill -f "server" 2>/dev/null || true
    echo -e "${GREEN}All services stopped.${NC}"
}

build_backend() {
    echo -e "${YELLOW}Building Go backend...${NC}"
    cd backend
    go mod tidy
    go build -o server cmd/server/main.go
    echo -e "${GREEN}Backend built successfully!${NC}"
}

test_api() {
    echo -e "${YELLOW}Testing API endpoints...${NC}"

    # Health check
    if curl -s http://localhost:8085/health | grep -q "healthy"; then
        echo -e "✅ Health check passed"
    else
        echo -e "❌ Health check failed"
        exit 1
    fi

    # Create test tree
    echo -e "${YELLOW}Creating test tree...${NC}"
    TREE_RESPONSE=$(curl -s -X POST http://localhost:8085/api/v1/trees \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Test Learning Tree",
            "description": "A test tree for development",
            "is_public": true,
            "root_node": {
                "title": "Getting Started",
                "content": "This is the root of our learning journey",
                "content_type": "text"
            }
        }')

    if echo "$TREE_RESPONSE" | grep -q "tree"; then
        echo -e "✅ Tree creation test passed"
        echo -e "${GREEN}API tests completed successfully!${NC}"
    else
        echo -e "❌ Tree creation test failed"
        echo "$TREE_RESPONSE"
        exit 1
    fi
}

clean_all() {
    echo -e "${YELLOW}Cleaning up development environment...${NC}"
    docker-compose down -v  # Remove volumes too
    cd backend
    rm -f server
    echo -e "${GREEN}Environment cleaned!${NC}"
}

show_help() {
    print_header
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  start     Start Neo4j and Redis services"
    echo "  backend   Start the Go backend server"
    echo "  stop      Stop all services"
    echo "  status    Show service status"
    echo "  build     Build the backend"
    echo "  test      Test API endpoints"
    echo "  clean     Clean up everything"
    echo "  help      Show this help message"
    echo ""
    echo "Quick start:"
    echo "  $0 start    # Start services"
    echo "  $0 backend  # Start backend (in another terminal)"
    echo ""
    echo "Or use: ./quick-start.sh for automated setup"
}

# Main command handling
case "${1:-help}" in
    "start")
        print_header
        start_services
        ;;
    "backend")
        print_header
        start_backend
        ;;
    "stop")
        print_header
        stop_services
        ;;
    "status")
        print_header
        print_status
        ;;
    "build")
        print_header
        build_backend
        ;;
    "test")
        print_header
        test_api
        ;;
    "clean")
        print_header
        clean_all
        ;;
    "help"|*)
        show_help
        ;;
esac
