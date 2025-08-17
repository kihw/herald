#!/bin/bash

echo "🚀 Starting LoL Match Exporter..."
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}⚠️  Go is not installed. Please install Go to run the backend.${NC}"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${YELLOW}⚠️  Node.js is not installed. Please install Node.js to run the frontend.${NC}"
    exit 1
fi

echo -e "${BLUE}📦 Installing dependencies...${NC}"

# Install Go dependencies
echo "Installing Go modules..."
go mod tidy

# Install Node.js dependencies
echo "Installing Node.js dependencies..."
cd web && npm install && cd ..

echo ""
echo -e "${GREEN}✅ Dependencies installed successfully!${NC}"
echo ""

# Build the development server
echo -e "${BLUE}🔨 Building development server...${NC}"
go build -o dev-server ./cmd/dev-server

# Build the frontend
echo -e "${BLUE}🔨 Building frontend...${NC}"
cd web && npm run build && cd ..

echo ""
echo -e "${GREEN}✅ Build completed successfully!${NC}"
echo ""

# Start the servers
echo -e "${BLUE}🚀 Starting servers...${NC}"
echo ""
echo -e "${GREEN}Backend API:${NC} http://localhost:8001"
echo -e "${GREEN}Frontend:${NC}    http://localhost:5173"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all servers${NC}"
echo ""

# Function to handle cleanup
cleanup() {
    echo ""
    echo -e "${YELLOW}🛑 Stopping all servers...${NC}"
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    echo -e "${GREEN}✅ All servers stopped. Goodbye!${NC}"
    exit 0
}

# Set trap to catch Ctrl+C
trap cleanup SIGINT

# Start backend in background
echo -e "${BLUE}Starting backend server...${NC}"
./dev-server &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Start frontend in background
echo -e "${BLUE}Starting frontend server...${NC}"
cd web && npm run dev &
FRONTEND_PID=$!
cd ..

# Wait for both processes
wait $BACKEND_PID $FRONTEND_PID
