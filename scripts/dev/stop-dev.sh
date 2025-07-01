#!/bin/bash

# Human Intelligence - Development Environment Stop Script
# This script stops all development processes and services

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

print_header() {
    echo ""
    printf "${PURPLE}════════════════════════════════════════════════════════════════${NC}\n"
    printf "${PURPLE}  🧠 Human Intelligence - Stop Development Environment${NC}\n"
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

# Kill processes on specific ports
kill_port_processes() {
    local port=$1
    local service_name=$2

    print_step "Stopping $service_name processes on port $port..."

    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        local pids=$(lsof -ti:$port)
        if [ -n "$pids" ]; then
            echo "$pids" | xargs kill -TERM 2>/dev/null || true
            sleep 2

            # Force kill if still running
            if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
                echo "$pids" | xargs kill -KILL 2>/dev/null || true
                print_warning "$service_name processes force killed"
            else
                print_success "$service_name stopped gracefully"
            fi
        fi
    else
        print_info "No $service_name processes found on port $port"
    fi
}

# Stop Docker services
stop_docker_services() {
    print_step "Stopping Docker services..."

    if command -v docker-compose &> /dev/null && [ -f "docker-compose.yml" ]; then
        if docker ps --format "table {{.Names}}" | grep -q "hi-"; then
            docker-compose down
            print_success "Docker services stopped"
        else
            print_info "No Docker services running"
        fi
    else
        print_info "Docker Compose not available or config not found"
    fi
}

# Kill any remaining backend processes
stop_backend_processes() {
    print_step "Stopping backend processes..."

    # Kill by port
    kill_port_processes $BACKEND_PORT "Backend"

    # Kill by process name
    if pgrep -f "server" >/dev/null; then
        pkill -f "server" || true
        print_success "Backend processes stopped"
    fi

    # Clean up log file
    if [ -f "backend.log" ]; then
        rm -f backend.log
        print_info "Backend log file cleaned up"
    fi
}

# Stop frontend processes
stop_frontend_processes() {
    print_step "Stopping frontend processes..."

    # Kill by port
    kill_port_processes $FRONTEND_PORT "Frontend"

    # Kill common frontend dev server processes
    local frontend_processes=("vite" "webpack-dev-server" "next" "svelte-kit")

    for process in "${frontend_processes[@]}"; do
        if pgrep -f "$process" >/dev/null; then
            pkill -f "$process" || true
        fi
    done

    # Clean up log file
    if [ -f "frontend.log" ]; then
        rm -f frontend.log
        print_info "Frontend log file cleaned up"
    fi

    print_success "Frontend processes stopped"
}

# Clean up any remaining processes
cleanup_remaining() {
    print_step "Cleaning up remaining processes..."

    # Kill any processes that might be related to our project
    local project_processes=("human-intelligence" "hi-" "server")

    for process in "${project_processes[@]}"; do
        if pgrep -f "$process" >/dev/null 2>&1; then
            print_info "Found remaining $process processes, cleaning up..."
            pkill -f "$process" 2>/dev/null || true
        fi
    done

    # Clean up any remaining log files
    for log_file in *.log; do
        if [ -f "$log_file" ]; then
            rm -f "$log_file"
        fi
    done

    print_success "Cleanup completed"
}

# Show final status
show_final_status() {
    echo ""
    print_success "🎉 Development environment stopped!"
    echo ""
    print_info "All processes have been terminated:"
    print_info "  ✓ Backend server (port $BACKEND_PORT)"
    print_info "  ✓ Frontend server (port $FRONTEND_PORT)"
    print_info "  ✓ Neo4j database"
    print_info "  ✓ Redis cache"
    print_info "  ✓ Log files cleaned up"
    echo ""
    print_info "To restart development environment:"
    print_info "  ./scripts/dev/start-dev.sh"
    echo ""
}

# Show help
show_help() {
    echo "Human Intelligence Development Stop Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --backend-port PORT    Backend port to stop (default: 8085)"
    echo "  --frontend-port PORT   Frontend port to stop (default: 3000)"
    echo "  --services-only        Only stop Docker services"
    echo "  --processes-only       Only stop application processes"
    echo "  --force                Force kill all processes immediately"
    echo "  --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                     # Stop everything"
    echo "  $0 --services-only     # Only stop Docker services"
    echo "  $0 --force             # Force kill everything"
    echo ""
}

# Main function
main() {
    local SERVICES_ONLY=false
    local PROCESSES_ONLY=false
    local FORCE=false

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
            --services-only)
                SERVICES_ONLY=true
                shift
                ;;
            --processes-only)
                PROCESSES_ONLY=true
                shift
                ;;
            --force)
                FORCE=true
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

    if [ "$FORCE" = true ]; then
        print_warning "Force mode enabled - killing all processes immediately"
    fi

    if [ "$SERVICES_ONLY" = false ]; then
        stop_backend_processes
        stop_frontend_processes

        if [ "$FORCE" = false ]; then
            cleanup_remaining
        fi
    fi

    if [ "$PROCESSES_ONLY" = false ]; then
        stop_docker_services
    fi

    if [ "$FORCE" = true ]; then
        print_step "Force killing any remaining processes..."
        # Nuclear option - kill everything related to our ports
        fuser -k ${BACKEND_PORT}/tcp 2>/dev/null || true
        fuser -k ${FRONTEND_PORT}/tcp 2>/dev/null || true
        fuser -k 7474/tcp 2>/dev/null || true  # Neo4j HTTP
        fuser -k 7687/tcp 2>/dev/null || true  # Neo4j Bolt
        fuser -k 6379/tcp 2>/dev/null || true  # Redis
        print_success "Force kill completed"
    fi

    show_final_status
}

# Run main function
main "$@"
