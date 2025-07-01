#!/bin/bash

# Simple Quick Start Script for Neo4j Knowledge Platform
# This script starts the basic services step by step

set -e

echo "🧠 Neo4j Knowledge Platform - Quick Start"
echo "=========================================="
echo ""

# Step 1: Start Neo4j and Redis
echo "Step 1: Starting Neo4j and Redis..."
docker-compose -f docker-compose.neo4j.yml up -d neo4j redis

echo "⏳ Waiting 15 seconds for Neo4j to initialize..."
sleep 15

# Step 2: Test Neo4j connection
echo ""
echo "Step 2: Testing Neo4j connection..."
if docker-compose -f docker-compose.neo4j.yml exec -T neo4j cypher-shell -u neo4j -p hi_password "RETURN 'Connected!' as status" 2>/dev/null; then
    echo "✅ Neo4j is ready!"
else
    echo "❌ Neo4j not ready yet. You may need to wait longer or check logs:"
    echo "   docker-compose -f docker-compose.neo4j.yml logs neo4j"
    echo ""
    echo "Continue anyway? (y/n)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Step 3: Build backend
echo ""
echo "Step 3: Building Go backend..."
cd backend
go build -o server-neo4j cmd/server/main_neo4j.go
echo "✅ Backend built!"

# Step 4: Start backend
echo ""
echo "Step 4: Starting backend server..."
echo "🚀 Starting on port 8085..."
echo ""
echo "Backend will start now. Press Ctrl+C to stop."
echo ""
echo "Quick links:"
echo "  • Health check: curl http://localhost:8085/health"
echo "  • Test page: open ../test-platform.html"
echo "  • Neo4j Browser: http://localhost:7474 (neo4j/hi_password)"
echo ""

# Set environment and start
export PORT=8085
export NEO4J_URI=bolt://localhost:7687
export NEO4J_USERNAME=neo4j
export NEO4J_PASSWORD=hi_password
export NEO4J_DATABASE=knowledgegraph
export GIN_MODE=debug

./server-neo4j
