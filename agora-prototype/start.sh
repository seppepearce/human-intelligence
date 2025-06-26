#!/bin/bash

# 🏛️ Digital Agora Prototype Startup Script
# Starts both the Go WebSocket server and SvelteKit frontend

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${CYAN}🏛️  Agora:${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅ Success:${NC} $1"
}

print_error() {
    echo -e "${RED}❌ Error:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️  Warning:${NC} $1"
}

# Function to cleanup processes on exit
cleanup() {
    print_status "Shutting down Digital Agora..."
    if [ ! -z "$SERVER_PID" ]; then
        kill $SERVER_PID 2>/dev/null || true
        print_status "Go server stopped"
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        print_status "Frontend server stopped"
    fi
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

print_status "Starting Digital Agora Prototype..."
echo -e "${PURPLE}===============================================${NC}"
echo -e "${CYAN}🏛️  Welcome to the Digital Agora${NC}"
echo -e "${CYAN}   Ancient Wisdom meets Modern Tools${NC}"
echo -e "${PURPLE}===============================================${NC}"
echo

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    print_error "Please run this script from the agora-prototype directory"
    exit 1
fi

# Check dependencies
print_status "Checking system dependencies..."

# Check Node.js
if ! command -v node &> /dev/null; then
    print_error "Node.js is required but not installed"
    exit 1
fi

# Check Go
if ! command -v go &> /dev/null; then
    print_error "Go is required but not installed"
    exit 1
fi

print_success "System dependencies verified"

# Install npm dependencies if needed
if [ ! -d "node_modules" ]; then
    print_status "Installing frontend dependencies..."
    npm install
    print_success "Frontend dependencies installed"
else
    print_status "Frontend dependencies already installed"
fi

# Start Go WebSocket server
print_status "Starting WebSocket server..."
cd server

# Download Go dependencies if needed
if [ ! -f "go.sum" ]; then
    print_status "Downloading Go dependencies..."
    go mod tidy
    print_success "Go dependencies downloaded"
fi

# Build and start Go server
go build -o agora-server main.go
./agora-server &
SERVER_PID=$!

cd ..

# Wait a moment for server to start
sleep 2

# Check if server started successfully
if ! kill -0 $SERVER_PID 2>/dev/null; then
    print_error "Failed to start WebSocket server"
    exit 1
fi

print_success "WebSocket server started (PID: $SERVER_PID)"

# Start SvelteKit frontend
print_status "Starting frontend development server..."
npm run dev &
FRONTEND_PID=$!

# Wait a moment for frontend to start
sleep 3

# Check if frontend started successfully
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    print_error "Failed to start frontend server"
    cleanup
    exit 1
fi

print_success "Frontend server started (PID: $FRONTEND_PID)"

echo
echo -e "${GREEN}🎉 Digital Agora is now running!${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}Frontend:${NC}    http://localhost:3001"
echo -e "${CYAN}WebSocket:${NC}   ws://localhost:8083/ws"
echo -e "${CYAN}API:${NC}         http://localhost:8083/api"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo
print_status "Services are running. Monitoring for changes..."

# Keep script running and monitor processes
while true; do
    # Check if server process is still running
    if ! kill -0 $SERVER_PID 2>/dev/null; then
        print_error "WebSocket server stopped unexpectedly"
        cleanup
        exit 1
    fi

    # Check if frontend process is still running
    if ! kill -0 $FRONTEND_PID 2>/dev/null; then
        print_error "Frontend server stopped unexpectedly"
        cleanup
        exit 1
    fi

    sleep 5
done
