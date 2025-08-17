#!/bin/bash
# Commands to run on the server for deployment

# Create deployment directory
sudo mkdir -p /opt/lol-match-exporter
cd /opt/lol-match-exporter

# Extract the archive
sudo tar -xzf /tmp/lol-exporter-deployment.tar.gz -C /opt/lol-match-exporter --strip-components=0

# Set proper permissions
sudo chown -R debian:debian /opt/lol-match-exporter
chmod +x /opt/lol-match-exporter/deploy-herald.sh

# Copy production environment file
cp /opt/lol-match-exporter/.env.herald /opt/lol-match-exporter/.env

# Install Docker and Docker Compose if not present
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker debian
fi

if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create necessary directories
sudo mkdir -p /opt/lol-match-exporter/data
sudo mkdir -p /opt/lol-match-exporter/logs
sudo mkdir -p /opt/lol-match-exporter/exports

# Set ownership
sudo chown -R debian:debian /opt/lol-match-exporter

echo "✅ Server preparation completed!"
echo "Next: Run 'docker-compose -f docker-compose.production.yml up -d' to start the application"
