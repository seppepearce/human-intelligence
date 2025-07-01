#!/bin/bash

# Human Intelligence - Full Stack Development Startup Script
# This script starts the complete development environment

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
BACKEND_PORT="${BACKEND_PORT:-8085}"
FRONTEND_PORT="${FRONTEND_PORT:-3000}"
NEO4J_URI="${NEO4J_URI:-bolt://localhost:7687}"
NEO4J_USERNAME="${NEO4J_USERNAME:-neo4j}"
NEO4J_PASSWORD="${NEO4J_PASSWORD:-hi_password}"

# Process tracking
BACKEND_PID=""
FRONTEND_PID=""
SERVICES_STARTED=""

print_header() {
    echo ""
    printf "${PURPLE}════════════════════════════════════════════════════════════════${NC}\n"
    printf "${PURPLE}  🧠 Human Intelligence - Full Stack Development Environment${NC}\n"
    printf "${PURPLE}════════════════════════════════════════════════════════════════${NC}\n"
    echo ""
}

print_step() {
    printf "${CYAN}→ $1${NC}\n"
}

print_success() {
    printf "${GREEN}✅ $1${NC}\n"
}

print_warning() {
    printf "${YELLOW}⚠️  $1${NC}\n"
}

print_error() {
    printf "${RED}❌ $1${NC}\n"
}

print_info() {
    printf "${BLUE}ℹ️  $1${NC}\n"
}

# Cleanup function
cleanup() {
    echo ""
    print_info "🛑 Shutting down development environment..."

    if [ -n "$BACKEND_PID" ]; then
        print_info "Stopping backend server (PID: $BACKEND_PID)..."
        kill $BACKEND_PID 2>/dev/null || true
    fi

    if [ -n "$FRONTEND_PID" ]; then
        print_info "Stopping frontend server (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID 2>/dev/null || true
    fi

    # Kill any remaining processes on our ports
    lsof -ti:$BACKEND_PORT | xargs kill -9 2>/dev/null || true
    lsof -ti:$FRONTEND_PORT | xargs kill -9 2>/dev/null || true

    print_success "Development environment stopped"
    exit 0
}

# Set up trap for cleanup
trap cleanup INT TERM

# Check prerequisites
check_prerequisites() {
    print_step "Checking prerequisites..."

    # Check if we're in the right directory
    if [ ! -d "backend" ] || [ ! -d "frontend" ]; then
        print_error "Please run this script from the project root directory"
        exit 1
    fi

    # Check for required tools
    local missing_tools=()

    if ! command -v go &> /dev/null; then
        missing_tools+=("go")
    fi

    if ! command -v node &> /dev/null; then
        missing_tools+=("node")
    fi

    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi

    if ! command -v docker-compose &> /dev/null; then
        missing_tools+=("docker-compose")
    fi

    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        print_info "Please install the missing tools and try again"
        exit 1
    fi

    print_success "All prerequisites met"
}

# Start database services
start_services() {
    print_step "Starting database services..."

    if [ ! -f "docker-compose.yml" ]; then
        print_error "docker-compose.yml not found"
        exit 1
    fi

    # Check if services are already running
    if docker ps --format "{{.Names}}" | grep -q "hi-neo4j"; then
        print_info "Neo4j already running"
    else
        print_info "Starting Neo4j and Redis..."
        docker-compose up -d neo4j redis
        SERVICES_STARTED="true"
        sleep 10  # Give services time to initialize

        # Wait for Neo4j to be ready
        print_info "Waiting for Neo4j to be ready..."
        local max_attempts=30
        local attempt=1

        while [ $attempt -le $max_attempts ]; do
            # Check if container is running first
            if docker ps --format "{{.Names}}" | grep -q "hi-neo4j"; then
                # Try to connect to Neo4j
                if docker-compose exec -T neo4j cypher-shell -u "$NEO4J_USERNAME" -p "$NEO4J_PASSWORD" "RETURN 1" &>/dev/null; then
                    print_success "Neo4j is ready"
                    break
                fi
            fi

            if [ $attempt -eq $max_attempts ]; then
                print_error "Neo4j failed to start after $max_attempts attempts"
                print_info "Checking Neo4j logs..."
                docker-compose logs neo4j | tail -10
                exit 1
            fi

            print_info "Attempt $attempt/$max_attempts - waiting for Neo4j..."
            sleep 5
            ((attempt++))
        done
    fi
}

# Build and start backend
start_backend() {
    print_step "Building and starting backend..."

    cd backend

    # Check if binary exists or needs rebuilding
    if [ ! -f "server" ] || [ "cmd/server/main.go" -nt "server" ]; then
        print_info "Building backend..."
        go mod download
        go build -o server cmd/server/main.go
        print_success "Backend built"
    else
        print_info "Using existing backend binary"
    fi

    # Check if port is available
    if lsof -Pi :$BACKEND_PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        print_warning "Port $BACKEND_PORT is already in use, killing existing process..."
        fuser -k $BACKEND_PORT/tcp 2>/dev/null || true
        sleep 2
    fi

    # Set environment and start backend
    export PORT="$BACKEND_PORT"
    export NEO4J_URI="$NEO4J_URI"
    export NEO4J_USERNAME="$NEO4J_USERNAME"
    export NEO4J_PASSWORD="$NEO4J_PASSWORD"
    export NEO4J_DATABASE="knowledgegraph"
    export GIN_MODE="debug"

    print_info "Starting backend on port $BACKEND_PORT..."
    ./server > ../backend.log 2>&1 &
    BACKEND_PID=$!

    cd ..

    # Wait for backend to start
    sleep 3
    if kill -0 $BACKEND_PID 2>/dev/null; then
        # Test if backend is responding
        if curl -s http://localhost:$BACKEND_PORT/health > /dev/null; then
            print_success "Backend started successfully (PID: $BACKEND_PID)"
        else
            print_warning "Backend started but not responding yet..."
        fi
    else
        print_error "Backend failed to start"
        cat backend.log
        exit 1
    fi
}

# Start frontend
start_frontend() {
    print_step "Starting frontend..."

    cd frontend

    # Install dependencies if needed
    if [ ! -d "node_modules" ] || [ "package.json" -nt "node_modules" ]; then
        print_info "Installing frontend dependencies..."
        if command -v pnpm &> /dev/null; then
            pnpm install
        elif command -v yarn &> /dev/null; then
            yarn install
        else
            npm install
        fi
        print_success "Dependencies installed"
    else
        print_info "Frontend dependencies up to date"
    fi

    # Check if port is available
    if lsof -Pi :$FRONTEND_PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        print_warning "Port $FRONTEND_PORT is already in use, killing existing process..."
        fuser -k $FRONTEND_PORT/tcp 2>/dev/null || true
        sleep 2
    fi

    # Set environment variables
    export VITE_API_URL="http://localhost:$BACKEND_PORT"
    export VITE_WS_URL="ws://localhost:$BACKEND_PORT"
    export VITE_NEO4J_BROWSER_URL="http://localhost:7474"

    print_info "Starting frontend on port $FRONTEND_PORT..."

    # Start frontend dev server
    if command -v pnpm &> /dev/null; then
        pnpm dev --port $FRONTEND_PORT --host > ../frontend.log 2>&1 &
    elif command -v yarn &> /dev/null; then
        yarn dev --port $FRONTEND_PORT --host > ../frontend.log 2>&1 &
    else
        npm run dev -- --port $FRONTEND_PORT --host > ../frontend.log 2>&1 &
    fi

    FRONTEND_PID=$!
    cd ..

    # Wait for frontend to start
    sleep 5
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        print_success "Frontend started successfully (PID: $FRONTEND_PID)"
    else
        print_error "Frontend failed to start"
        cat frontend.log
        exit 1
    fi
}

# Show development info
show_dev_info() {
    echo ""
    print_success "🎉 Development environment is ready!"
    echo ""
    print_info "📱 Frontend:      http://localhost:$FRONTEND_PORT"
    print_info "🔧 Backend API:   http://localhost:$BACKEND_PORT"
    print_info "🗄️  Neo4j Browser: http://localhost:7474 (neo4j/hi_password)"
    print_info "📊 Test Page:     file://$(pwd)/test-platform.html"
    echo ""
    print_info "📜 Logs:"
    print_info "  Backend:  tail -f backend.log"
    print_info "  Frontend: tail -f frontend.log"
    echo ""
    print_info "🔗 Quick Tests:"
    print_info "  curl http://localhost:$BACKEND_PORT/health"
    print_info "  curl http://localhost:$BACKEND_PORT/api/v1/users/user-1"
    echo ""
    print_warning "Press Ctrl+C to stop all services"
    echo ""
}

# Open browser (optional)
open_browser() {
    if command -v xdg-open &> /dev/null; then
        sleep 2
        xdg-open "http://localhost:$FRONTEND_PORT" &>/dev/null &
    elif command -v open &> /dev/null; then
        sleep 2
        open "http://localhost:$FRONTEND_PORT" &>/dev/null &
    fi
}

# Show help
show_help() {
    echo "Human Intelligence Full Stack Development Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --backend-port PORT    Set backend port (default: 8085)"
    echo "  --frontend-port PORT   Set frontend port (default: 3000)"
    echo "  --no-open             Don't open browser automatically"
    echo "  --services-only       Only start database services"
    echo "  --backend-only        Only start backend (assumes services running)"
    echo "  --frontend-only       Only start frontend (assumes backend running)"
    echo "  --help                Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  BACKEND_PORT          Backend server port"
    echo "  FRONTEND_PORT         Frontend server port"
    echo "  NEO4J_URI            Neo4j connection URI"
    echo "  NEO4J_USERNAME       Neo4j username"
    echo "  NEO4J_PASSWORD       Neo4j password"
    echo ""
}

# Main function
main() {
    local NO_OPEN=false
    local SERVICES_ONLY=false
    local BACKEND_ONLY=false
    local FRONTEND_ONLY=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --backend-port)
                BACKEND_PORT="$2"
                shift 2
                ;;
            --frontend-port)
                FRONTEND_PORT="$2"
                shift 2
                ;;
            --no-open)
                NO_OPEN=true
                shift
                ;;
            --services-only)
                SERVICES_ONLY=true
                shift
                ;;
            --backend-only)
                BACKEND_ONLY=true
                shift
                ;;
            --frontend-only)
                FRONTEND_ONLY=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done

    print_header

    check_prerequisites

    if [ "$FRONTEND_ONLY" = false ]; then
        start_services
    fi

    if [ "$SERVICES_ONLY" = false ] && [ "$FRONTEND_ONLY" = false ]; then
        start_backend
    fi

    if [ "$SERVICES_ONLY" = false ] && [ "$BACKEND_ONLY" = false ]; then
        start_frontend
    fi

    show_dev_info

    if [ "$NO_OPEN" = false ] && [ "$SERVICES_ONLY" = false ]; then
        open_browser
    fi

    # Keep script running and wait for interrupt
    while true; do
        sleep 5

        # Check if processes are still running
        if [ -n "$BACKEND_PID" ] && ! kill -0 $BACKEND_PID 2>/dev/null; then
            print_error "Backend process died unexpectedly"
            break
        fi

        if [ -n "$FRONTEND_PID" ] && ! kill -0 $FRONTEND_PID 2>/dev/null; then
            print_error "Frontend process died unexpectedly"
            break
        fi
    done
}

# Run main function
main "$@"
