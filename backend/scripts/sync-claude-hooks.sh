#!/bin/bash

# Herald.lol Gaming Analytics - Claude Hooks Synchronization
# Sync hooks between main project and backend directories

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$BACKEND_DIR")"

# Gaming colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
GOLD='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] [SYNC]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] [SUCCESS]${NC} $1"
}

log_gaming() {
    echo -e "${GOLD}[$(date +'%Y-%m-%d %H:%M:%S')] [HERALD-GAMING]${NC} $1"
}

# Sync Claude hooks
sync_hooks() {
    log_gaming "🎮 Synchronizing Herald.lol Claude hooks..."
    
    cd "$BACKEND_DIR"
    
    # Create hooks directory if not exists
    mkdir -p .claude/hooks
    
    # Check if main project has hooks
    if [ -d "$PROJECT_ROOT/.claude/hooks" ]; then
        log_info "📁 Found Claude hooks in main project directory"
        
        # Copy hooks to backend
        log_info "📋 Copying hooks to backend directory..."
        cp -r "$PROJECT_ROOT/.claude/hooks/"* .claude/hooks/
        
        # Make sure they're executable
        chmod +x .claude/hooks/*.sh
        chmod +x .claude/hooks/*.py
        
        log_success "✅ Herald.lol Claude hooks synchronized"
        
        # List synchronized hooks
        echo ""
        log_info "🔧 Available Herald.lol gaming hooks:"
        for hook in .claude/hooks/*; do
            if [ -f "$hook" ]; then
                hook_name=$(basename "$hook")
                echo "  🎮 $hook_name"
            fi
        done
        
    else
        log_info "⚠️ No Claude hooks found in main project directory"
        return 1
    fi
    
    echo ""
    log_success "🎮 Herald.lol Claude hooks ready for gaming development!"
}

# Test hooks
test_hooks() {
    log_gaming "🧪 Testing Herald.lol Claude hooks..."
    
    cd "$BACKEND_DIR"
    
    if [ ! -d ".claude/hooks" ]; then
        log_info "❌ No hooks directory found - run sync first"
        return 1
    fi
    
    # Test load-context hook
    if [ -f ".claude/hooks/load-context.sh" ]; then
        log_info "🧪 Testing load-context hook..."
        ./.claude/hooks/load-context.sh >/dev/null 2>&1 && \
            log_success "✅ load-context hook working" || \
            log_info "⚠️ load-context hook has issues"
    fi
    
    # Test deploy-check hook
    if [ -f ".claude/hooks/deploy-check.sh" ]; then
        log_info "🧪 Testing deploy-check hook..."
        ./.claude/hooks/deploy-check.sh >/dev/null 2>&1 && \
            log_success "✅ deploy-check hook working" || \
            log_info "⚠️ deploy-check hook has issues"
    fi
    
    log_success "🎮 Herald.lol Claude hooks testing completed"
}

# Show hook status
show_status() {
    log_gaming "📊 Herald.lol Claude Hooks Status"
    echo "=================================="
    
    cd "$BACKEND_DIR"
    
    if [ -d ".claude/hooks" ]; then
        echo ""
        log_info "🔧 Backend Hooks Directory: .claude/hooks/"
        ls -la .claude/hooks/ | sed 's/^/  /'
    else
        echo ""
        log_info "❌ No backend hooks directory found"
    fi
    
    if [ -d "$PROJECT_ROOT/.claude/hooks" ]; then
        echo ""
        log_info "🔧 Main Project Hooks Directory: $PROJECT_ROOT/.claude/hooks/"
        ls -la "$PROJECT_ROOT/.claude/hooks/" | sed 's/^/  /'
    else
        echo ""
        log_info "❌ No main project hooks directory found"
    fi
}

# Usage
usage() {
    cat << EOF
🎮 Herald.lol Claude Hooks Synchronization

Usage: $0 [COMMAND]

COMMANDS:
    sync      Sync hooks from main project to backend (default)
    test      Test all Herald.lol gaming hooks
    status    Show hooks status
    -h, --help    Show this help message

EXAMPLES:
    # Sync Herald.lol gaming hooks
    $0 sync

    # Test gaming hooks
    $0 test
    
    # Check hooks status
    $0 status

🎯 Gaming Focus: Herald.lol Claude Code integration
⚡ Performance: <5s analytics, 1M+ users
EOF
}

# Main function
main() {
    case "${1:-sync}" in
        sync)
            sync_hooks
            ;;
        test)
            test_hooks
            ;;
        status)
            show_status
            ;;
        -h|--help)
            usage
            ;;
        *)
            log_info "❌ Unknown command: $1"
            echo ""
            usage
            exit 1
            ;;
    esac
}

# Handle interruption gracefully
trap 'echo ""; log_info "🎮 Herald.lol hooks sync interrupted"; exit 0' SIGINT

# Run main function
main "$@"