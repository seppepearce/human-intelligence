#!/bin/bash

# Human Intelligence - Backend Development Startup Script
# This script builds and runs the Neo4j backend for development

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
BACKEND_DIR="backend"
BINARY_NAME="server-neo4j"
PORT="${PORT:-8085}"
NEO4J_URI="${NEO4J_URI:-bolt://localhost:7687}"
NEO4J_USERNAME="${NEO4J_USERNAME:-neo4j}"
NEO4J_PASSWORD="${NEO4J_PASSWORD:-hi_password}"

print_header() {
    echo ""
    printf "${PURPLE}════════════════════════════════════════════════════════════════${NC}\n"
    printf "${PURPLE}  🧠 Human Intelligence - Backend Development Server${NC}\n"
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

# Check if we're in the right directory
check_directory() {
    if [ ! -d "$BACKEND_DIR" ]; then
        print_error "Backend directory not found. Please run this script from the project root."
        exit 1
    fi

    if [ ! -f "$BACKEND_DIR/cmd/server/main_neo4j.go" ]; then
        print_error "Backend main file not found at $BACKEND_DIR/cmd/server/main_neo4j.go"
        exit 1
    fi
}

# Check if Neo4j is running
check_neo4j() {
    print_step "Checking Neo4j connection..."

    if ! command -v docker &> /dev/null; then
        print_warning "Docker not found. Make sure Neo4j is running manually."
        return
    fi

    # Check if Neo4j container is running
    if docker ps --format "{{.Names}}" | grep -q "hi-neo4j"; then
        print_success "Neo4j container is running"
    else
        print_warning "Neo4j container not found. Starting Neo4j..."
        if [ -f "docker-compose.neo4j.yml" ]; then
            docker-compose -f docker-compose.neo4j.yml up -d neo4j redis
            print_info "Waiting for Neo4j to initialize..."
            sleep 15
        else
            print_error "docker-compose.neo4j.yml not found. Please start Neo4j manually."
            exit 1
        fi
    fi

    # Test Neo4j connection
    max_attempts=20
    attempt=1
    while [ $attempt -le $max_attempts ]; do
        # Check if container is running first
        if docker ps --format "{{.Names}}" | grep -q "hi-neo4j"; then
            # Try to connect to Neo4j using docker-compose exec
            if docker-compose -f docker-compose.neo4j.yml exec -T neo4j cypher-shell -u "$NEO4J_USERNAME" -p "$NEO4J_PASSWORD" "RETURN 1" &>/dev/null; then
                print_success "Neo4j is ready and accepting connections"
                break
            fi
        else
            print_warning "Neo4j container not found, checking if it needs to be started..."
        fi

        if [ $attempt -eq $max_attempts ]; then
            print_error "Neo4j failed to start after $max_attempts attempts"
            print_info "Checking Neo4j logs..."
            docker-compose -f docker-compose.neo4j.yml logs neo4j | tail -10
            exit 1
        fi

        print_info "Attempt $attempt/$max_attempts - waiting for Neo4j..."
        sleep 5
        ((attempt++))
    done
}

# Build the backend
build_backend() {
    print_step "Building backend..."
    cd "$BACKEND_DIR"

    # Check if go.mod exists
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found in backend directory"
        exit 1
    fi

    # Download dependencies if needed
    print_info "Downloading Go dependencies..."
    go mod download

    # Build the binary
    print_info "Compiling backend..."
    if go build -o "$BINARY_NAME" cmd/server/main_neo4j.go; then
        print_success "Backend built successfully"
    else
        print_error "Failed to build backend"
        exit 1
    fi

    cd ..
}

# Check if port is available
check_port() {
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        print_warning "Port $PORT is already in use"
        print_info "Trying to find the process..."

        if command -v lsof &> /dev/null; then
            lsof -i :$PORT
        fi

        read -p "Kill the process and continue? (y/N): " -r
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_info "Killing process on port $PORT..."
            fuser -k $PORT/tcp 2>/dev/null || true
            sleep 2
        else
            print_info "Please choose a different port with: PORT=8086 $0"
            exit 1
        fi
    fi
}

# Start the backend server
start_server() {
    print_step "Starting backend server..."

    cd "$BACKEND_DIR"

    # Set environment variables
    export PORT="$PORT"
    export NEO4J_URI="$NEO4J_URI"
    export NEO4J_USERNAME="$NEO4J_USERNAME"
    export NEO4J_PASSWORD="$NEO4J_PASSWORD"
    export NEO4J_DATABASE="knowledgegraph"
    export GIN_MODE="debug"

    print_info "Configuration:"
    print_info "  Port: $PORT"
    print_info "  Neo4j URI: $NEO4J_URI"
    print_info "  Neo4j Database: knowledgegraph"
    echo ""

    print_success "🚀 Starting Human Intelligence Backend..."
    print_info "API will be available at: http://localhost:$PORT"
    print_info "Health check: http://localhost:$PORT/health"
    print_info "API docs: http://localhost:$PORT/api/v1/"
    echo ""
    print_info "Press Ctrl+C to stop the server"
    echo ""

    # Start the server with trap for graceful shutdown
    trap 'print_info "\n🛑 Shutting down server..."; kill $!; exit 0' INT TERM

    ./"$BINARY_NAME" &
    wait
}

# Show help
show_help() {
    echo "Human Intelligence Backend Development Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --port PORT          Set the server port (default: 8085)"
    echo "  --neo4j-uri URI      Set Neo4j URI (default: bolt://localhost:7687)"
    echo "  --skip-neo4j-check   Skip Neo4j connectivity check"
    echo "  --no-build          Skip building the backend"
    echo "  --help              Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  PORT                 Server port"
    echo "  NEO4J_URI           Neo4j connection URI"
    echo "  NEO4J_USERNAME      Neo4j username (default: neo4j)"
    echo "  NEO4J_PASSWORD      Neo4j password (default: hi_password)"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Start with defaults"
    echo "  $0 --port 8090                       # Start on port 8090"
    echo "  $0 --skip-neo4j-check               # Skip Neo4j check"
    echo "  PORT=8090 $0                         # Start on port 8090 (env var)"
    echo ""
}

# Main function
main() {
    SKIP_NEO4J_CHECK=false
    NO_BUILD=false

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --port)
                PORT="$2"
                shift 2
                ;;
            --neo4j-uri)
                NEO4J_URI="$2"
                shift 2
                ;;
            --skip-neo4j-check)
                SKIP_NEO4J_CHECK=true
                shift
                ;;
            --no-build)
                NO_BUILD=true
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

    check_directory

    if [ "$SKIP_NEO4J_CHECK" = false ]; then
        check_neo4j
    fi

    check_port

    if [ "$NO_BUILD" = false ]; then
        build_backend
    fi

    start_server
}

# Run main function with all arguments
main "$@"
