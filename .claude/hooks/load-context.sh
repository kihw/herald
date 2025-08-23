#!/bin/bash
# Herald.lol Context Loader
# Loads gaming platform context at session start

echo "🎮 Loading Herald.lol Gaming Analytics Platform Context..."

# Check if we're in a Herald.lol project
if [ ! -f "go.mod" ] && [ ! -f "package.json" ]; then
    echo "📁 Herald.lol project structure not detected"
    exit 0
fi

# Load gaming platform context
echo "🎯 Herald.lol Gaming Analytics Platform"
echo "📊 Stack: Go + Gin + React + TypeScript + PostgreSQL + Redis"
echo "🎮 Focus: League of Legends & TFT Analytics"
echo "⚡ Performance Target: <5s analysis, 99.9% uptime, 1M+ concurrent users"
echo ""

# Check current git branch for gaming workflows
if git rev-parse --git-dir > /dev/null 2>&1; then
    CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")
    echo "🌿 Git Branch: $CURRENT_BRANCH"
    
    # Gaming-specific branch hints
    if [[ $CURRENT_BRANCH == *"analytics"* ]]; then
        echo "📈 Analytics branch - Focus on gaming metrics (KDA, CS/min, Vision Score)"
    elif [[ $CURRENT_BRANCH == *"riot"* ]]; then
        echo "🎮 Riot API branch - Remember rate limits and ToS compliance"
    elif [[ $CURRENT_BRANCH == *"performance"* ]]; then
        echo "⚡ Performance branch - Target <5s post-game analysis"
    elif [[ $CURRENT_BRANCH == *"ui"* ]] || [[ $CURRENT_BRANCH == *"frontend"* ]]; then
        echo "🎨 Frontend branch - Gaming UX with LoL theme integration"
    fi
    echo ""
fi

# Check for gaming platform dependencies
echo "🔍 Checking Herald.lol dependencies..."

# Go backend checks
if [ -f "go.mod" ]; then
    echo "✅ Go backend detected"
    if grep -q "gin-gonic/gin" go.mod 2>/dev/null; then
        echo "  🌐 Gin web framework found"
    fi
    if grep -q "lib/pq\|pgx" go.mod 2>/dev/null; then
        echo "  🗄️  PostgreSQL driver found"
    fi
    if grep -q "redis" go.mod 2>/dev/null; then
        echo "  🔄 Redis client found"
    fi
fi

# Frontend checks
if [ -f "package.json" ]; then
    echo "✅ Frontend project detected"
    if grep -q "react" package.json 2>/dev/null; then
        echo "  ⚛️  React found"
    fi
    if grep -q "typescript" package.json 2>/dev/null; then
        echo "  📘 TypeScript found"
    fi
    if grep -q "@mui/material\|@material-ui" package.json 2>/dev/null; then
        echo "  🎨 Material-UI found"
    fi
    if grep -q "chart.js\|recharts\|d3" package.json 2>/dev/null; then
        echo "  📊 Charts library found - Perfect for gaming analytics!"
    fi
fi

# Check for Docker/Kubernetes
if [ -f "Dockerfile" ]; then
    echo "🐳 Docker configuration found"
fi

if [ -d "k8s" ] || [ -f "kubernetes.yaml" ]; then
    echo "☸️  Kubernetes manifests found"
fi

# Gaming platform specific reminders
echo ""
echo "🎮 Herald.lol Development Reminders:"
echo "  • Riot Games API rate limits: Personal (100/2min), Production (varies)"
echo "  • Gaming metrics priority: KDA, CS/min, Vision Score, Damage Share"
echo "  • Performance targets: <5s analysis, <2s UI load"
echo "  • Scalability goal: 100k+ MAU, 1M+ concurrent support"
echo "  • Security: GDPR compliance + gaming data protection"
echo ""

echo "🚀 Ready for Herald.lol gaming analytics development!"

# Output success for Claude context injection
exit 0