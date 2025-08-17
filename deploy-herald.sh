#!/bin/bash

# Deployment script for herald.lol (51.178.17.78)
# Run this script from your local machine

set -e

SERVER_IP="51.178.17.78"
SERVER_USER="root"  # Adjust if needed
REMOTE_DIR="/opt/lol-match-exporter"
LOCAL_DIR="$(pwd)"

echo "🚀 Deploying LoL Match Exporter to herald.lol..."

# 1. Create deployment archive
echo "📦 Creating deployment archive..."
tar --exclude='*.exe' \
    --exclude='*.db' \
    --exclude='node_modules' \
    --exclude='web/dist' \
    --exclude='web/node_modules' \
    --exclude='data' \
    --exclude='exports' \
    --exclude='logs' \
    --exclude='.git' \
    --exclude='*.log' \
    --exclude='*.tmp' \
    -czf lol-exporter-deployment.tar.gz .

echo "✅ Archive created: lol-exporter-deployment.tar.gz"

# 2. Transfer to server
echo "📤 Transferring files to $SERVER_IP..."
scp lol-exporter-deployment.tar.gz $SERVER_USER@$SERVER_IP:/tmp/

# 3. Deploy on server
echo "🏗️ Deploying on server..."
ssh $SERVER_USER@$SERVER_IP << 'ENDSSH'
    # Update system
    apt-get update
    
    # Install Docker if not present
    if ! command -v docker &> /dev/null; then
        echo "🐳 Installing Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sh get-docker.sh
        systemctl enable docker
        systemctl start docker
    fi
    
    # Install Docker Compose if not present
    if ! command -v docker-compose &> /dev/null; then
        echo "🐳 Installing Docker Compose..."
        curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
    fi
    
    # Create deployment directory
    mkdir -p /opt/lol-match-exporter
    cd /opt/lol-match-exporter
    
    # Stop existing containers if running
    docker-compose -f docker-compose.production.yml down 2>/dev/null || true
    
    # Extract new deployment
    tar -xzf /tmp/lol-exporter-deployment.tar.gz
    
    # Create required directories
    mkdir -p data exports logs logs/nginx
    
    # Copy environment file
    cp .env.herald .env
    
    # Set permissions
    chown -R 1000:1000 data exports logs
    
    # Build and start containers
    echo "🚀 Starting containers..."
    docker-compose -f docker-compose.production.yml up -d --build
    
    # Wait for services to start
    sleep 10
    
    # Check health
    echo "🏥 Checking health..."
    curl -f http://localhost/health || echo "⚠️ Health check failed"
    
    # Show status
    docker-compose -f docker-compose.production.yml ps
    
    # Clean up
    rm /tmp/lol-exporter-deployment.tar.gz
ENDSSH

# 4. Test deployment
echo "🧪 Testing deployment..."
sleep 5
curl -f http://herald.lol/health && echo "✅ herald.lol is responding!" || echo "❌ Health check failed"

# Clean up local archive
rm lol-exporter-deployment.tar.gz

echo "🎉 Deployment completed!"
echo "🌐 Application should be available at: http://herald.lol"
echo "📊 Health endpoint: http://herald.lol/health"
echo "📚 API docs: http://herald.lol/docs"