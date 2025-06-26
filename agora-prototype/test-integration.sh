#!/bin/bash

# 🧪 Integration Test Script - Digital Agora Prototype
# Tests that all components work together properly

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Test configuration
BACKEND_PORT=8082
FRONTEND_PORT=3001
TEST_TIMEOUT=30
SERVER_PID=""
FRONTEND_PID=""

# Function to print colored output
print_test() {
    echo -e "${CYAN}🧪 Test:${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅ Pass:${NC} $1"
}

print_failure() {
    echo -e "${RED}❌ Fail:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️  Warning:${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ️  Info:${NC} $1"
}

# Cleanup function
cleanup() {
    print_info "Cleaning up test processes..."
    if [ ! -z "$SERVER_PID" ]; then
        kill $SERVER_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    # Kill any remaining processes on our ports
    lsof -ti:$BACKEND_PORT | xargs kill -9 2>/dev/null || true
    lsof -ti:$FRONTEND_PORT | xargs kill -9 2>/dev/null || true
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM EXIT

# Function to wait for service to be ready
wait_for_service() {
    local port=$1
    local service_name=$2
    local timeout=$3
    local count=0

    print_info "Waiting for $service_name on port $port..."

    while [ $count -lt $timeout ]; do
        if curl -s "http://localhost:$port" >/dev/null 2>&1; then
            print_success "$service_name is ready"
            return 0
        fi
        sleep 1
        count=$((count + 1))
    done

    print_failure "$service_name failed to start within $timeout seconds"
    return 1
}

# Function to test API endpoint
test_api_endpoint() {
    local endpoint=$1
    local expected_status=$2
    local description=$3

    print_test "Testing $description"

    local response=$(curl -s -w "%{http_code}" "http://localhost:$BACKEND_PORT$endpoint")
    local body="${response%???}"
    local status="${response: -3}"

    if [ "$status" -eq "$expected_status" ]; then
        print_success "API endpoint $endpoint returned $status"
        return 0
    else
        print_failure "API endpoint $endpoint returned $status, expected $expected_status"
        return 1
    fi
}

# Function to test WebSocket connection
test_websocket() {
    print_test "Testing WebSocket connection"

    # Create a simple WebSocket test using Node.js
    cat > /tmp/ws_test.js << 'EOF'
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8082/ws');

ws.on('open', function open() {
    console.log('WebSocket connected');

    // Send a test message
    ws.send(JSON.stringify({
        type: 'presence',
        userId: 'test_user',
        data: {
            id: 'test_user',
            name: 'Test User',
            avatar: '🧪',
            interests: ['Testing']
        }
    }));

    setTimeout(() => {
        ws.close();
        process.exit(0);
    }, 2000);
});

ws.on('message', function message(data) {
    console.log('Received:', data.toString());
});

ws.on('error', function error(err) {
    console.error('WebSocket error:', err);
    process.exit(1);
});

ws.on('close', function close() {
    console.log('WebSocket closed');
});

// Timeout after 10 seconds
setTimeout(() => {
    console.error('WebSocket test timeout');
    process.exit(1);
}, 10000);
EOF

    if node /tmp/ws_test.js >/dev/null 2>&1; then
        print_success "WebSocket connection test passed"
        rm -f /tmp/ws_test.js
        return 0
    else
        print_failure "WebSocket connection test failed"
        rm -f /tmp/ws_test.js
        return 1
    fi
}

# Function to test frontend accessibility
test_frontend() {
    print_test "Testing frontend accessibility"

    local response=$(curl -s -w "%{http_code}" "http://localhost:$FRONTEND_PORT")
    local status="${response: -3}"

    if [ "$status" -eq "200" ]; then
        print_success "Frontend is accessible"
        return 0
    else
        print_failure "Frontend returned status $status"
        return 1
    fi
}

# Main test execution
main() {
    echo -e "${PURPLE}=================================================${NC}"
    echo -e "${CYAN}🧪 Digital Agora Integration Tests${NC}"
    echo -e "${PURPLE}=================================================${NC}"
    echo

    # Test 1: Check dependencies
    print_test "Checking system dependencies"

    if ! command -v node &> /dev/null; then
        print_failure "Node.js not found"
        exit 1
    fi
    print_success "Node.js found: $(node --version)"

    if ! command -v go &> /dev/null; then
        print_failure "Go not found"
        exit 1
    fi
    print_success "Go found: $(go version | cut -d' ' -f3)"

    if ! command -v curl &> /dev/null; then
        print_failure "curl not found"
        exit 1
    fi
    print_success "curl found"

    # Test 2: Check if we're in the right directory
    print_test "Checking project structure"

    if [ ! -f "package.json" ]; then
        print_failure "package.json not found. Please run from agora-prototype directory"
        exit 1
    fi

    if [ ! -d "server" ]; then
        print_failure "server directory not found"
        exit 1
    fi

    if [ ! -d "src" ]; then
        print_failure "src directory not found"
        exit 1
    fi

    print_success "Project structure validated"

    # Test 3: Install dependencies
    print_test "Installing Node.js dependencies"

    if [ ! -d "node_modules" ]; then
        npm install --silent
        if [ $? -eq 0 ]; then
            print_success "Node.js dependencies installed"
        else
            print_failure "Failed to install Node.js dependencies"
            exit 1
        fi
    else
        print_success "Node.js dependencies already installed"
    fi

    # Test 4: Build and start Go server
    print_test "Building and starting Go server"

    cd server

    # Install Go dependencies
    if [ ! -f "go.sum" ]; then
        go mod tidy
    fi

    # Build server
    go build -o agora-server main.go
    if [ $? -ne 0 ]; then
        print_failure "Failed to build Go server"
        exit 1
    fi

    # Start server in background
    ./agora-server > /tmp/agora-server.log 2>&1 &
    SERVER_PID=$!

    cd ..

    # Wait for server to be ready
    if ! wait_for_service $BACKEND_PORT "Go server" $TEST_TIMEOUT; then
        print_failure "Go server failed to start"
        cat /tmp/agora-server.log
        exit 1
    fi

    # Test 5: Test API endpoints
    test_api_endpoint "/api/spaces" 200 "spaces API endpoint"

    # Test 6: Test WebSocket connection
    test_websocket

    # Test 7: Start frontend
    print_test "Starting frontend development server"

    npm run dev > /tmp/agora-frontend.log 2>&1 &
    FRONTEND_PID=$!

    # Wait for frontend to be ready
    if ! wait_for_service $FRONTEND_PORT "Frontend server" $TEST_TIMEOUT; then
        print_failure "Frontend server failed to start"
        cat /tmp/agora-frontend.log
        exit 1
    fi

    # Test 8: Test frontend accessibility
    test_frontend

    # Test 9: Test integration endpoints
    print_test "Testing frontend-backend integration"

    # Check if frontend can load (contains expected content)
    local frontend_content=$(curl -s "http://localhost:$FRONTEND_PORT")
    if echo "$frontend_content" | grep -q "Digital Agora"; then
        print_success "Frontend contains expected content"
    else
        print_failure "Frontend does not contain expected content"
        exit 1
    fi

    # Test 10: Performance check
    print_test "Running basic performance checks"

    # Measure API response time
    local start_time=$(date +%s%N)
    curl -s "http://localhost:$BACKEND_PORT/api/spaces" >/dev/null
    local end_time=$(date +%s%N)
    local response_time=$(( (end_time - start_time) / 1000000 ))

    if [ $response_time -lt 1000 ]; then
        print_success "API response time: ${response_time}ms (excellent)"
    elif [ $response_time -lt 2000 ]; then
        print_success "API response time: ${response_time}ms (good)"
    else
        print_warning "API response time: ${response_time}ms (could be better)"
    fi

    # Test 11: Memory usage check
    print_test "Checking memory usage"

    if command -v ps &> /dev/null; then
        local server_memory=$(ps -o rss= -p $SERVER_PID | xargs)
        local frontend_memory=$(ps -o rss= -p $FRONTEND_PID | xargs)

        if [ ! -z "$server_memory" ] && [ "$server_memory" -lt 100000 ]; then
            print_success "Server memory usage: ${server_memory}KB (efficient)"
        else
            print_warning "Server memory usage: ${server_memory}KB (check for leaks)"
        fi

        if [ ! -z "$frontend_memory" ] && [ "$frontend_memory" -lt 200000 ]; then
            print_success "Frontend memory usage: ${frontend_memory}KB (efficient)"
        else
            print_warning "Frontend memory usage: ${frontend_memory}KB (typical for dev server)"
        fi
    else
        print_info "ps command not available, skipping memory check"
    fi

    # Test Summary
    echo
    echo -e "${PURPLE}=================================================${NC}"
    echo -e "${GREEN}🎉 All Integration Tests Passed!${NC}"
    echo -e "${PURPLE}=================================================${NC}"
    echo
    echo -e "${CYAN}Services Running:${NC}"
    echo -e "${CYAN}Frontend:${NC}  http://localhost:$FRONTEND_PORT"
    echo -e "${CYAN}Backend:${NC}   http://localhost:$BACKEND_PORT"
    echo -e "${CYAN}WebSocket:${NC} ws://localhost:$BACKEND_PORT/ws"
    echo
    echo -e "${YELLOW}Press Ctrl+C to stop services and exit${NC}"
    echo

    # Keep services running for manual testing
    print_info "Services will continue running for manual testing..."
    print_info "Check /tmp/agora-server.log and /tmp/agora-frontend.log for detailed logs"

    # Wait for user to stop
    while true; do
        sleep 5
        # Check if processes are still running
        if ! kill -0 $SERVER_PID 2>/dev/null; then
            print_failure "Server process died unexpectedly"
            break
        fi
        if ! kill -0 $FRONTEND_PID 2>/dev/null; then
            print_failure "Frontend process died unexpectedly"
            break
        fi
    done
}

# Run main function
main "$@"
