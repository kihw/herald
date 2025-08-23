#!/bin/bash

echo "🧪 Testing Herald.lol setup..."

# Check if required files exist
echo "📁 Checking file structure..."
required_files=(
    "docker-compose.dev.yml"
    "Makefile"
    ".env"
    "backend/go.mod"
    "backend/go.sum"
    "backend/cmd/server/main.go"
    "frontend/package.json"
    "frontend/vite.config.ts"
    "nginx/dev.conf"
)

missing_files=0
for file in "${required_files[@]}"; do
    if [[ -f "$file" ]]; then
        echo "✅ $file"
    else
        echo "❌ $file"
        missing_files=$((missing_files + 1))
    fi
done

if [[ $missing_files -gt 0 ]]; then
    echo "❌ $missing_files required files are missing"
    exit 1
fi

echo ""
echo "🐳 Testing Docker Compose configuration..."
if docker-compose -f docker-compose.dev.yml config > /dev/null 2>&1; then
    echo "✅ Docker Compose configuration is valid"
else
    echo "❌ Docker Compose configuration has errors"
    exit 1
fi

echo ""
echo "🔧 Environment Configuration:"
echo "   - Database: $(grep 'DB_HOST=' .env | cut -d'=' -f2)"
echo "   - Redis: $(grep 'REDIS_HOST=' .env | cut -d'=' -f2)" 
echo "   - Environment: $(grep 'ENV=' .env | cut -d'=' -f2)"

echo ""
echo "✅ Herald.lol setup test completed successfully!"
echo ""
echo "Next steps:"
echo "1. Run 'make build' to build the services"
echo "2. Run 'make dev' to start the development environment"
echo "3. Visit http://localhost to access the application"