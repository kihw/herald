#!/bin/bash

# LoL Match Exporter Deployment Script
# Usage: ./deploy.sh [production]

set -e

echo "🚀 Starting LoL Match Exporter deployment..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  Creating .env file..."
    cat > .env << EOF
# Required: Your Riot Games API Key
RIOT_API_KEY=your_api_key_here

# Optional: Server API key for protection
EXPORTER_API_KEY=dev_secret_key
EOF
    echo "⚠️  Please edit .env file with your actual API keys before running again."
    exit 1
fi

# Source environment variables
source .env

if [ -z "$RIOT_API_KEY" ] || [ "$RIOT_API_KEY" = "your_api_key_here" ]; then
    echo "❌ Please set your RIOT_API_KEY in .env file"
    exit 1
fi

# Create necessary directories
mkdir -p jobs logs

# Build and start containers
echo "🔨 Building containers..."
if [ "$1" = "production" ]; then
    echo "🌐 Starting with Nginx reverse proxy..."
    docker-compose --profile production up -d --build
else
    echo "🔧 Starting in development mode..."
    docker-compose up -d --build
fi

# Wait for health check
echo "⏳ Waiting for application to start..."
timeout=60
counter=0

while [ $counter -lt $timeout ]; do
    if curl -f http://localhost:8000/health > /dev/null 2>&1; then
        echo "✅ Application is healthy!"
        break
    fi
    sleep 2
    counter=$((counter + 2))
    echo "Waiting... ($counter/$timeout seconds)"
done

if [ $counter -ge $timeout ]; then
    echo "❌ Application failed to start within $timeout seconds"
    echo "📋 Container logs:"
    docker-compose logs lol-exporter
    exit 1
fi

# Show status
echo ""
echo "🎉 Deployment successful!"
echo ""
echo "📊 Application URLs:"
if [ "$1" = "production" ]; then
    echo "   Frontend: http://localhost"
    echo "   API: http://localhost/api"
else
    echo "   Frontend: http://localhost:8000"
    echo "   API: http://localhost:8000/docs"
fi
echo ""
echo "🔧 Useful commands:"
echo "   View logs: docker-compose logs -f lol-exporter"
echo "   Stop: docker-compose down"
echo "   Rebuild: docker-compose up -d --build"
echo ""
echo "📋 Container status:"
docker-compose ps