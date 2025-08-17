#!/bin/bash

# Script de test pour vérifier la configuration Jenkins
echo "=== Test de configuration Jenkins ==="

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Variables
JENKINS_URL=${JENKINS_URL:-"http://localhost:8080"}
PROJECT_NAME="lol-match-exporter-freestyle"

echo -e "${YELLOW}Testing Jenkins configuration...${NC}"

# Test 1: Vérifier Go
echo -n "Testing Go installation... "
if command -v go &> /dev/null; then
    echo -e "${GREEN}✓ Go found: $(go version)${NC}"
else
    echo -e "${RED}✗ Go not found${NC}"
    exit 1
fi

# Test 2: Vérifier Docker
echo -n "Testing Docker installation... "
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓ Docker found: $(docker --version)${NC}"
else
    echo -e "${RED}✗ Docker not found${NC}"
    exit 1
fi

# Test 3: Vérifier SSH vers le serveur de production
echo -n "Testing SSH connection to production server... "
if ssh -o ConnectTimeout=5 -o BatchMode=yes debian@51.178.17.78 "echo 'SSH OK'" &> /dev/null; then
    echo -e "${GREEN}✓ SSH connection OK${NC}"
else
    echo -e "${RED}✗ SSH connection failed${NC}"
    echo "Please check SSH keys and server availability"
fi

# Test 4: Vérifier les fichiers de configuration
echo "Testing configuration files:"

files=("Dockerfile.simple-fullstack" "docker-compose.complete.yml" "nginx/nginx-fullstack.conf")
for file in "${files[@]}"; do
    echo -n "  Checking $file... "
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
    fi
done

# Test 5: Build de test
echo -n "Testing Go build... "
if go build -o test-main . &> /dev/null; then
    echo -e "${GREEN}✓ Go build successful${NC}"
    rm -f test-main
else
    echo -e "${RED}✗ Go build failed${NC}"
fi

# Test 6: Test Docker build (rapide)
echo -n "Testing Docker build... "
if docker build -f Dockerfile.simple-fullstack -t jenkins-test:latest . &> /dev/null; then
    echo -e "${GREEN}✓ Docker build successful${NC}"
    docker rmi jenkins-test:latest &> /dev/null
else
    echo -e "${RED}✗ Docker build failed${NC}"
fi

# Test 7: Vérifier Jenkins (si accessible)
if curl -s "$JENKINS_URL" &> /dev/null; then
    echo -e "${GREEN}✓ Jenkins accessible at $JENKINS_URL${NC}"
    
    # Vérifier si le projet existe
    if curl -s "$JENKINS_URL/job/$PROJECT_NAME/" &> /dev/null; then
        echo -e "${GREEN}✓ Project $PROJECT_NAME exists in Jenkins${NC}"
    else
        echo -e "${YELLOW}! Project $PROJECT_NAME not found in Jenkins${NC}"
        echo "  Create it manually or check the name"
    fi
else
    echo -e "${YELLOW}! Jenkins not accessible at $JENKINS_URL${NC}"
    echo "  Please check Jenkins URL and ensure it's running"
fi

echo ""
echo -e "${GREEN}=== Configuration test completed ===${NC}"
echo ""
echo "Next steps:"
echo "1. Create the freestyle project in Jenkins"
echo "2. Configure build steps as described in the guide"
echo "3. Set up SSH credentials for deployment"
echo "4. Configure webhooks for automatic builds"
