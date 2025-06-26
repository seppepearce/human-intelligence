#!/bin/bash

echo "🏛️ Starting AI Knowledge Trees Development Environment..."

# Check if backend is built
if [ ! -f "backend/bin/server" ]; then
    echo "Building backend..."
    cd backend && go build -o bin/server . && cd ..
fi

# Start backend in background
echo "Starting backend on port 8081..."
cd backend && ./bin/server &
BACKEND_PID=$!
cd ..

# Wait a moment for backend to start
sleep 2

# Start frontend
echo "Starting frontend on port 5173..."
cd frontend && npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "🎉 Development servers started!"
echo "Frontend: http://localhost:5173"
echo "Backend:  http://localhost:8081"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for user interrupt
trap 'kill $BACKEND_PID $FRONTEND_PID 2>/dev/null' EXIT
wait
