#!/bin/bash

echo "🧪 Testing AI Knowledge Trees System..."

# Test backend build
echo "Testing backend build..."
cd backend
if go build -o bin/server .; then
    echo "✅ Backend builds successfully"
else
    echo "❌ Backend build failed"
    exit 1
fi
cd ..

# Test frontend build
echo "Testing frontend build..."
cd frontend
if npm run build; then
    echo "✅ Frontend builds successfully"
else
    echo "❌ Frontend build failed"
    exit 1
fi
cd ..

echo "🎉 All tests passed!"
