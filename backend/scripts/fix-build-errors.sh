#!/bin/bash

# Herald.lol Gaming Analytics - Quick Build Error Fix
# Fix common compilation errors for gaming platform

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Gaming colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
GOLD='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] [FIX]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] [SUCCESS]${NC} $1"
}

log_gaming() {
    echo -e "${GOLD}[$(date +'%Y-%m-%d %H:%M:%S')] [HERALD-GAMING]${NC} $1"
}

# Fix summoner unused variables
fix_summoner_errors() {
    log_info "🔧 Fixing summoner service errors..."
    
    cd "$PROJECT_ROOT"
    
    # Fix unused variable in helpers.go
    if [ -f "internal/summoner/helpers.go" ]; then
        sed -i 's/analyticsRequest := analytics.AnalyticsRequest{/_ = analytics.AnalyticsRequest{/g' internal/summoner/helpers.go
        log_success "  ✅ Fixed unused analyticsRequest variable"
    fi
    
    # Fix unused imports
    if [ -f "internal/summoner/models.go" ]; then
        # Comment out unused riot import
        sed -i 's|"github.com/herald-lol/herald/backend/internal/riot"|// "github.com/herald-lol/herald/backend/internal/riot"|g' internal/summoner/models.go
        log_success "  ✅ Fixed unused riot import in models.go"
    fi
    
    if [ -f "internal/summoner/service.go" ]; then
        # Comment out unused json import  
        sed -i 's|"encoding/json"|// "encoding/json"|g' internal/summoner/service.go
        log_success "  ✅ Fixed unused json import in service.go"
    fi
}

# Add missing types for coaching service
fix_coaching_types() {
    log_info "🔧 Adding missing coaching service types..."
    
    cd "$PROJECT_ROOT"
    
    # Add missing types to coaching_types.go
    cat >> internal/services/coaching_types.go << 'EOF'

// Additional coaching types for Herald.lol gaming platform
type ChampionCoachingTip struct {
    Champion   string   `json:"champion"`
    Role       string   `json:"role"`
    Tips       []string `json:"tips"`
    Priority   int      `json:"priority"`
    Difficulty string   `json:"difficulty"`
}

type MetaAdaptationGuidance struct {
    MetaPatch      string   `json:"meta_patch"`
    Adaptations    []string `json:"adaptations"`
    ChampionShifts []string `json:"champion_shifts"`
    StrategyChanges []string `json:"strategy_changes"`
}

type PerformanceGoal struct {
    Metric     string  `json:"metric"`
    Current    float64 `json:"current"`
    Target     float64 `json:"target"`
    Timeline   string  `json:"timeline"`
    Difficulty string  `json:"difficulty"`
}

type CoachingSchedule struct {
    SessionType string    `json:"session_type"`
    Duration    int       `json:"duration"`
    Frequency   string    `json:"frequency"`
    Goals       []string  `json:"goals"`
    StartTime   time.Time `json:"start_time"`
}

type ProgressTracking struct {
    PlayerID     string             `json:"player_id"`
    Goals        []PerformanceGoal  `json:"goals"`
    Achievements []string           `json:"achievements"`
    LastUpdated  time.Time          `json:"last_updated"`
    Progress     map[string]float64 `json:"progress"`
}
EOF

    log_success "  ✅ Added missing coaching types"
}

# Add missing analytics types
fix_analytics_types() {
    log_info "🔧 Adding missing analytics types..."
    
    cd "$PROJECT_ROOT"
    
    # Add missing AnalyticsResult to analytics models
    if ! grep -q "AnalyticsResult" internal/analytics/models.go; then
        cat >> internal/analytics/models.go << 'EOF'

// AnalyticsResult represents the result of gaming analytics processing
type AnalyticsResult struct {
    PlayerID           string             `json:"player_id"`
    MatchID           string             `json:"match_id"`
    ProcessingTime    time.Duration      `json:"processing_time"`
    PerformanceScore  float64           `json:"performance_score"`
    Metrics          map[string]float64 `json:"metrics"`
    Recommendations  []string           `json:"recommendations"`
    Errors           []string           `json:"errors"`
    Success          bool               `json:"success"`
}
EOF
        log_success "  ✅ Added AnalyticsResult type"
    fi
}

# Add missing service interfaces
fix_service_interfaces() {
    log_info "🔧 Adding missing service interfaces..."
    
    cd "$PROJECT_ROOT"
    
    # Create missing service interfaces file
    cat > internal/services/interfaces.go << 'EOF'
package services

import (
    "context"
    "github.com/herald-lol/herald/backend/internal/models"
)

// MatchService interface for match-related operations
type MatchService interface {
    GetMatch(ctx context.Context, matchID string) (*models.Match, error)
    GetPlayerMatches(ctx context.Context, playerID string, limit int) ([]*models.Match, error)
    ProcessMatch(ctx context.Context, matchID string) error
}

// UserService interface for user-related operations  
type UserService interface {
    GetUser(ctx context.Context, userID string) (*models.User, error)
    GetUserByRiotPUUID(ctx context.Context, puuid string) (*models.User, error)
    UpdateUser(ctx context.Context, user *models.User) error
}

// AnalyticsService interface for analytics operations
type AnalyticsService interface {
    ProcessPlayerPerformance(ctx context.Context, playerID string, matchID string) error
    GetPlayerStats(ctx context.Context, playerID string) (map[string]interface{}, error)
    GetPerformanceMetrics(ctx context.Context, playerID string) (map[string]float64, error)
}
EOF

    log_success "  ✅ Added missing service interfaces"
}

# Main build fix
fix_build_errors() {
    log_gaming "🎮 Fixing Herald.lol gaming platform build errors..."
    
    fix_summoner_errors
    fix_coaching_types  
    fix_analytics_types
    fix_service_interfaces
    
    # Test build
    log_info "🧪 Testing Go build..."
    cd "$PROJECT_ROOT"
    
    if go build -o /tmp/herald-test ./cmd/server >/dev/null 2>&1; then
        log_success "✅ Herald.lol gaming platform builds successfully!"
        rm -f /tmp/herald-test
        return 0
    else
        log_info "⚠️ Still some build errors - check output:"
        go build -o /tmp/herald-test ./cmd/server 2>&1 | head -10
        return 1
    fi
}

# Main function
main() {
    log_gaming "🎮 Starting Herald.lol build error fixes"
    
    fix_build_errors
    
    log_success "🎮 Herald.lol gaming platform build fixes completed!"
}

# Handle interruption gracefully
trap 'echo ""; log_info "🎮 Herald.lol build fix interrupted"; exit 0' SIGINT

# Run main function
main "$@"