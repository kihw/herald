#!/bin/bash

# Script to fix import paths for Herald.lol gRPC implementation

echo "🔧 Fixing Herald.lol import paths..."

# Fix github.com/herald-lol/backend/ to github.com/herald-lol/herald/backend/
find . -name "*.go" -type f -exec sed -i 's|github.com/herald-lol/backend/|github.com/herald-lol/herald/backend/|g' {} \;

# Fix github.com/herald/ to github.com/herald-lol/herald/backend/
find . -name "*.go" -type f -exec sed -i 's|github.com/herald/|github.com/herald-lol/herald/backend/|g' {} \;

# Fix herald.lol/internal/ to github.com/herald-lol/herald/backend/internal/
find . -name "*.go" -type f -exec sed -i 's|herald.lol/internal/|github.com/herald-lol/herald/backend/internal/|g' {} \;

# Fix other common patterns
find . -name "*.go" -type f -exec sed -i 's|"github.com/go-redis/redis/v8"|"github.com/redis/go-redis/v9"|g' {} \;

echo "✅ Import paths fixed!"
echo "🧹 Cleaning go modules..."

# Clean module cache
go clean -modcache
go mod tidy

echo "🎮 Herald.lol imports ready for gRPC!"