#!/bin/bash

# Human Intelligence - Development Wrapper Script
# Convenient entry point for all development tasks

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo ""
    printf "${PURPLE}🧠 Human Intelligence - Development CLI${NC}\n"
    echo ""
}

print_info() {
    printf "${BLUE}ℹ️  $1${NC}\n"
}

show_help() {
    print_header
    echo "Usage: ./dev.sh <command> [options]"
    echo ""
    echo "Commands:"
    echo "  start                 Start full development environment (Neo4j + Backend + Frontend)"
    echo "  backend               Start only backend (assumes Neo4j running)"
    echo "  stop                  Stop all development processes"
    echo "  services              Start only database services (Neo4j + Redis)"
    echo "  status                Show status of development environment"
    echo "  logs                  Show development logs"
    echo "  test                  Open test page in browser"
    echo "  neo4j                 Open Neo4j browser"
    echo ""
    echo "Options (passed to underlying scripts):"
    echo "  --backend-port PORT   Set backend port (default: 8085)"
    echo "  --frontend-port PORT  Set frontend port (default: 3000)"
    echo "  --no-open            Don't open browser automatically"
    echo "  --help               Show help for specific command"
    echo ""
    echo "Examples:"
    echo "  ./dev.sh start                    # Start everything"
    echo "  ./dev.sh backend                  # Just backend"
    echo "  ./dev.sh start --no-open         # Start without opening browser"
    echo "  ./dev.sh stop                     # Stop everything"
    echo ""
    echo "Quick links:"
    echo "  Frontend:      http://localhost:3000"
    echo "  Backend API:   http://localhost:8085"
    echo "  Neo4j Browser: http://localhost:7474"
    echo "  Test Page:     file://$(pwd)/test-platform.html"
    echo ""
}

case "${1:-help}" in
    start|dev)
        shift
        print_info "Starting full development environment..."
        exec ./scripts/dev/start-dev.sh "$@"
        ;;

    backend|be)
        shift
        print_info "Starting backend only..."
        exec ./scripts/dev/start-backend.sh "$@"
        ;;

    stop|kill)
        shift
        print_info "Stopping development environment..."
        exec ./scripts/dev/stop-dev.sh "$@"
        ;;

    services|db)
        shift
        print_info "Starting database services only..."
        exec ./scripts/dev/start-dev.sh --services-only "$@"
        ;;

    status)
        print_header
        echo "Development Environment Status:"
        echo ""

        # Check processes
        if lsof -Pi :8085 -sTCP:LISTEN -t >/dev/null 2>&1; then
            printf "${GREEN}✅ Backend running on port 8085${NC}\n"
        else
            echo "❌ Backend not running"
        fi

        if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
            printf "${GREEN}✅ Frontend running on port 3000${NC}\n"
        else
            echo "❌ Frontend not running"
        fi

        # Check Docker services
        if command -v docker >/dev/null 2>&1; then
            if docker ps --format "table {{.Names}}" | grep -q "hi-neo4j"; then
                printf "${GREEN}✅ Neo4j running${NC}\n"
            else
                echo "❌ Neo4j not running"
            fi

            if docker ps --format "table {{.Names}}" | grep -q "hi-redis"; then
                printf "${GREEN}✅ Redis running${NC}\n"
            else
                echo "❌ Redis not running"
            fi
        fi

        echo ""
        ;;

    logs)
        print_header
        echo "Available logs:"
        echo ""

        if [ -f "backend.log" ]; then
            echo "Backend logs:"
            echo "  tail -f backend.log"
            echo ""
        fi

        if [ -f "frontend.log" ]; then
            echo "Frontend logs:"
            echo "  tail -f frontend.log"
            echo ""
        fi

        echo "Docker logs:"
        echo "  docker-compose -f docker-compose.neo4j.yml logs -f"
        echo ""
        ;;

    test)
        if command -v xdg-open >/dev/null 2>&1; then
            xdg-open "file://$(pwd)/test-platform.html"
        elif command -v open >/dev/null 2>&1; then
            open "file://$(pwd)/test-platform.html"
        else
            echo "Open manually: file://$(pwd)/test-platform.html"
        fi
        ;;

    neo4j|n4j)
        if command -v xdg-open >/dev/null 2>&1; then
            xdg-open "http://localhost:7474"
        elif command -v open >/dev/null 2>&1; then
            open "http://localhost:7474"
        else
            echo "Open manually: http://localhost:7474"
            echo "Login: neo4j / hi_password"
        fi
        ;;

    help|--help|-h)
        show_help
        ;;

    *)
        echo "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
