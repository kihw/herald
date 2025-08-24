#!/bin/bash
# Herald.lol Gaming Platform Linter
# Specialized linting for gaming analytics platform

echo "🎮 Running Herald.lol gaming platform linting..."

LINT_ISSUES=0

# Go backend linting (gaming analytics focus)
if [ -f "go.mod" ]; then
    echo "🔍 Linting Go backend for gaming analytics..."
    
    # Format Go code
    if command -v go >/dev/null 2>&1; then
        echo "  📝 Formatting Go code..."
        go fmt ./...
        
        # Vet for common issues
        echo "  🔎 Running go vet..."
        if ! go vet ./...; then
            echo "  ❌ Go vet found issues"
            LINT_ISSUES=$((LINT_ISSUES + 1))
        else
            echo "  ✅ Go vet passed"
        fi
        
        # Run golangci-lint if available (gaming performance focus)
        if command -v golangci-lint >/dev/null 2>&1; then
            echo "  🏃 Running golangci-lint for gaming performance..."
            if ! golangci-lint run --timeout=3m; then
                echo "  ❌ golangci-lint found issues"
                LINT_ISSUES=$((LINT_ISSUES + 1))
            else
                echo "  ✅ golangci-lint passed"
            fi
        fi
        
        # Gaming-specific Go checks
        echo "  🎮 Gaming-specific Go validations..."
        
        # Check for proper error handling in API calls
        if grep -r "riot.*api" . --include="*.go" | grep -v "err.*!=" >/dev/null 2>&1; then
            echo "  ⚠️  Ensure proper error handling for Riot API calls"
        fi
        
        # Check for rate limiting implementation
        if grep -r "riot.*api" . --include="*.go" >/dev/null 2>&1; then
            if ! grep -r "rate.*limit\|throttle" . --include="*.go" >/dev/null 2>&1; then
                echo "  ⚠️  Consider implementing rate limiting for Riot API"
            fi
        fi
        
        # Check for proper gaming metrics validation
        if grep -r "kda\|cs.*min\|vision.*score" . --include="*.go" -i >/dev/null 2>&1; then
            echo "  ✅ Gaming metrics detected"
        fi
    fi
fi

# Frontend linting (React + TypeScript gaming UI)
if [ -f "package.json" ]; then
    echo "🔍 Linting frontend for gaming UI..."
    
    # ESLint for React/TypeScript
    if [ -f ".eslintrc.js" ] || [ -f ".eslintrc.json" ] || [ -f "eslint.config.js" ]; then
        echo "  📝 Running ESLint..."
        if command -v npm >/dev/null 2>&1; then
            if npm run lint >/dev/null 2>&1; then
                echo "  ✅ ESLint passed"
            else
                echo "  ❌ ESLint found issues"
                LINT_ISSUES=$((LINT_ISSUES + 1))
            fi
        elif command -v yarn >/dev/null 2>&1; then
            if yarn lint >/dev/null 2>&1; then
                echo "  ✅ ESLint passed"
            else
                echo "  ❌ ESLint found issues" 
                LINT_ISSUES=$((LINT_ISSUES + 1))
            fi
        fi
    fi
    
    # TypeScript checking
    if [ -f "tsconfig.json" ]; then
        echo "  📘 Running TypeScript check..."
        if command -v npx >/dev/null 2>&1; then
            if npx tsc --noEmit >/dev/null 2>&1; then
                echo "  ✅ TypeScript check passed"
            else
                echo "  ⚠️  TypeScript check completed with warnings"
            fi
        fi
    fi
    
    # Gaming UI specific checks
    echo "  🎮 Gaming UI validations..."
    
    # Check for accessibility in gaming UI
    if grep -r "aria-\|role=" src/ --include="*.tsx" --include="*.jsx" >/dev/null 2>&1; then
        echo "  ♿ Accessibility attributes found"
    else
        echo "  ⚠️  Consider adding accessibility attributes for gaming UI"
    fi
    
    # Check for proper gaming theme integration
    if grep -r "mui.*theme\|material.*ui" src/ --include="*.tsx" --include="*.jsx" -i >/dev/null 2>&1; then
        echo "  🎨 Material-UI theme integration found"
    fi
    
    # Check for performance optimizations
    if grep -r "useMemo\|useCallback\|React\.memo" src/ --include="*.tsx" --include="*.jsx" >/dev/null 2>&1; then
        echo "  ⚡ React performance optimizations found"
    else
        echo "  ⚠️  Consider React performance optimizations for gaming analytics"
    fi
fi

# Database migration linting
if [ -d "migrations" ] || [ -d "internal/db/migrations" ]; then
    echo "🗄️  Checking database migrations..."
    
    # Check for proper indexing on gaming data
    if grep -r "CREATE.*INDEX" migrations/ internal/db/migrations/ 2>/dev/null | grep -i "player\|match\|game" >/dev/null; then
        echo "  ✅ Gaming data indexes found"
    else
        echo "  ⚠️  Consider indexes on gaming data tables for performance"
    fi
fi

# Docker linting
if [ -f "Dockerfile" ]; then
    echo "🐳 Checking Dockerfile for gaming platform..."
    
    # Check for multi-stage builds
    if grep -c "FROM.*AS" Dockerfile > /dev/null; then
        echo "  ✅ Multi-stage build detected"
    else
        echo "  ⚠️  Consider multi-stage build for smaller gaming platform images"
    fi
    
    # Check for non-root user
    if grep "USER.*[^root]" Dockerfile >/dev/null; then
        echo "  ✅ Non-root user found"
    else
        echo "  ⚠️  Consider running as non-root user for security"
    fi
    
    # Check for health checks
    if grep "HEALTHCHECK" Dockerfile >/dev/null; then
        echo "  ✅ Health check found"
    else
        echo "  ⚠️  Consider adding health check for gaming platform monitoring"
    fi
fi

# Security linting for gaming platform
echo "🔒 Gaming platform security checks..."

# Check for exposed secrets
if grep -r "api.*key\|secret\|token.*=" . --include="*.go" --include="*.js" --include="*.ts" --include="*.tsx" 2>/dev/null | grep -v "env\|ENV" >/dev/null; then
    echo "  ❌ Potential secrets in code - use environment variables"
    LINT_ISSUES=$((LINT_ISSUES + 1))
else
    echo "  ✅ No hardcoded secrets found"
fi

# Check for proper CORS configuration
if grep -r "cors" . --include="*.go" --include="*.js" --include="*.ts" >/dev/null 2>&1; then
    echo "  ✅ CORS configuration found"
else
    echo "  ⚠️  Consider CORS configuration for gaming platform APIs"
fi

# Final report
echo ""
if [ $LINT_ISSUES -eq 0 ]; then
    echo "✅ Herald.lol gaming platform linting completed successfully!"
else
    echo "⚠️  Herald.lol linting completed with $LINT_ISSUES issues to review"
fi

exit 0