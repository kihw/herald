#!/bin/bash

# Herald.lol Visual Regression Testing Runner
# Execute Playwright visual tests for gaming analytics components

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
FRONTEND_URL="${FRONTEND_URL:-http://localhost:3000}"
API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
RESULTS_DIR="./test-results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo -e "${BLUE}🎮 Herald.lol Visual Regression Testing${NC}"
echo -e "${BLUE}====================================${NC}"
echo ""
echo -e "Frontend URL: ${GREEN}$FRONTEND_URL${NC}"
echo -e "API Base URL: ${GREEN}$API_BASE_URL${NC}"
echo -e "Results Directory: ${GREEN}$RESULTS_DIR${NC}"
echo ""

# Function to check system health
check_system_health() {
    echo -e "${BLUE}🏥 Checking Herald.lol system health...${NC}"
    
    # Check API health
    if curl -f -s "$API_BASE_URL/health" >/dev/null; then
        echo -e "${GREEN}✅ Herald.lol API health check passed${NC}"
    else
        echo -e "${YELLOW}⚠️  Herald.lol API health check failed - continuing with mocked data${NC}"
    fi
    
    # Check frontend
    if curl -f -s "$FRONTEND_URL" >/dev/null; then
        echo -e "${GREEN}✅ Herald.lol frontend health check passed${NC}"
    else
        echo -e "${RED}❌ Herald.lol frontend health check failed${NC}"
        echo -e "${YELLOW}⚠️  Make sure 'npm run dev' is running${NC}"
        return 1
    fi
    
    echo ""
}

# Function to display test menu
show_menu() {
    echo -e "${BLUE}🎯 Select Herald.lol Visual Regression Test:${NC}"
    echo ""
    echo -e "1) ${GREEN}Gaming Analytics Dashboard${NC} - Visual test analytics widgets"
    echo -e "2) ${GREEN}Match Analysis Interface${NC} - Visual test match analysis"
    echo -e "3) ${GREEN}Gaming Components${NC} - Visual test individual components"
    echo -e "4) ${BLUE}Run All Visual Tests${NC} - Complete visual regression suite"
    echo -e "5) ${YELLOW}Update Visual Baselines${NC} - Update screenshot baselines"
    echo -e "6) ${GREEN}Visual Test Report${NC} - View test results report"
    echo -e "7) ${BLUE}Interactive Debug Mode${NC} - Debug visual tests interactively"
    echo -e "8) Exit"
    echo ""
    echo -n "Enter your choice (1-8): "
}

# Function to run specific visual test
run_visual_test() {
    local test_file="$1"
    local test_name="$2"
    
    echo -e "${YELLOW}📸 Running $test_name visual tests...${NC}"
    
    # Set environment variables
    export FRONTEND_URL="$FRONTEND_URL"
    export API_BASE_URL="$API_BASE_URL"
    
    # Run Playwright test
    if npx playwright test --config=playwright.config.ts "$test_file"; then
        echo -e "${GREEN}✅ $test_name visual tests completed successfully${NC}"
        return 0
    else
        echo -e "${RED}❌ $test_name visual tests failed${NC}"
        return 1
    fi
}

# Function to run all visual tests
run_all_visual_tests() {
    echo -e "${BLUE}🚀 Running complete Herald.lol visual regression suite...${NC}"
    echo ""
    
    local failed_tests=0
    
    # Gaming Analytics Dashboard
    echo -e "${YELLOW}📊 Testing Gaming Analytics Dashboard...${NC}"
    if ! run_visual_test "visual-tests/gaming-analytics-dashboard.visual.test.ts" "Gaming Analytics Dashboard"; then
        failed_tests=$((failed_tests + 1))
    fi
    echo ""
    
    # Match Analysis Interface
    echo -e "${YELLOW}🎮 Testing Match Analysis Interface...${NC}"
    if ! run_visual_test "visual-tests/match-analysis.visual.test.ts" "Match Analysis Interface"; then
        failed_tests=$((failed_tests + 1))
    fi
    echo ""
    
    # Gaming Components
    echo -e "${YELLOW}🎯 Testing Gaming Components...${NC}"
    if ! run_visual_test "visual-tests/gaming-components.visual.test.ts" "Gaming Components"; then
        failed_tests=$((failed_tests + 1))
    fi
    echo ""
    
    # Summary
    if [ $failed_tests -eq 0 ]; then
        echo -e "${GREEN}🎉 All Herald.lol visual tests passed!${NC}"
        echo -e "${GREEN}✅ Gaming Analytics Dashboard: PASSED${NC}"
        echo -e "${GREEN}✅ Match Analysis Interface: PASSED${NC}"
        echo -e "${GREEN}✅ Gaming Components: PASSED${NC}"
    else
        echo -e "${RED}❌ $failed_tests visual test suite(s) failed${NC}"
        echo -e "${YELLOW}📊 Check test results for details${NC}"
    fi
    
    return $failed_tests
}

# Function to update visual baselines
update_visual_baselines() {
    echo -e "${YELLOW}📸 Updating Herald.lol visual baselines...${NC}"
    echo -e "${YELLOW}⚠️  This will update all screenshot baselines${NC}"
    echo -n "Are you sure? (y/N): "
    read -r confirm
    
    if [[ $confirm =~ ^[Yy]$ ]]; then
        export FRONTEND_URL="$FRONTEND_URL"
        export API_BASE_URL="$API_BASE_URL"
        
        if npx playwright test --config=playwright.config.ts --update-snapshots; then
            echo -e "${GREEN}✅ Visual baselines updated successfully${NC}"
            echo -e "${BLUE}ℹ️  Review changes with 'git diff' before committing${NC}"
        else
            echo -e "${RED}❌ Failed to update visual baselines${NC}"
        fi
    else
        echo -e "${BLUE}ℹ️  Visual baseline update cancelled${NC}"
    fi
}

# Function to show visual test report
show_visual_report() {
    echo -e "${BLUE}📊 Herald.lol Visual Test Report${NC}"
    
    # Check if report exists
    if [ -f "$RESULTS_DIR/index.html" ]; then
        echo -e "${GREEN}✅ Opening HTML report...${NC}"
        
        # Try to open in browser
        if command -v open >/dev/null 2>&1; then
            open "$RESULTS_DIR/index.html"
        elif command -v xdg-open >/dev/null 2>&1; then
            xdg-open "$RESULTS_DIR/index.html"
        else
            echo -e "${BLUE}📊 Report location: $(pwd)/$RESULTS_DIR/index.html${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  No visual test report found${NC}"
        echo -e "${BLUE}ℹ️  Run visual tests first to generate report${NC}"
    fi
    
    # Show summary if available
    if [ -f "$RESULTS_DIR/visual-test-summary.md" ]; then
        echo -e "${BLUE}📋 Visual Test Summary:${NC}"
        head -20 "$RESULTS_DIR/visual-test-summary.md"
    fi
}

# Function for interactive debug mode
interactive_debug() {
    echo -e "${BLUE}🔍 Herald.lol Visual Test Interactive Debug Mode${NC}"
    echo -e "${BLUE}===============================================${NC}"
    
    export FRONTEND_URL="$FRONTEND_URL"
    export API_BASE_URL="$API_BASE_URL"
    
    echo -e "${YELLOW}🎮 Starting interactive gaming visual tests...${NC}"
    echo -e "${BLUE}ℹ️  Use browser developer tools to inspect elements${NC}"
    echo -e "${BLUE}ℹ️  Tests will pause for manual inspection${NC}"
    
    npx playwright test --config=playwright.config.ts --ui
}

# Function to install prerequisites
install_prerequisites() {
    echo -e "${BLUE}📦 Installing Herald.lol visual testing prerequisites...${NC}"
    
    # Install Playwright
    if ! command -v npx >/dev/null 2>&1; then
        echo -e "${RED}❌ npx not found. Please install Node.js and npm${NC}"
        exit 1
    fi
    
    # Install Playwright browsers
    echo -e "${YELLOW}🌐 Installing Playwright browsers...${NC}"
    npx playwright install
    
    # Install system dependencies
    echo -e "${YELLOW}🔧 Installing system dependencies...${NC}"
    npx playwright install-deps
    
    echo -e "${GREEN}✅ Prerequisites installed successfully${NC}"
}

# Function to clean up old results
cleanup_results() {
    echo -e "${YELLOW}🧹 Cleaning up old visual test results...${NC}"
    
    if [ -d "$RESULTS_DIR" ]; then
        # Archive old results
        if [ -f "$RESULTS_DIR/visual-results.json" ]; then
            mv "$RESULTS_DIR" "${RESULTS_DIR}_${TIMESTAMP}"
            echo -e "${GREEN}✅ Old results archived to ${RESULTS_DIR}_${TIMESTAMP}${NC}"
        fi
    fi
    
    # Create fresh results directory
    mkdir -p "$RESULTS_DIR"
}

# Main execution
main() {
    # Check if running in CI
    if [ "$CI" = "true" ]; then
        echo -e "${BLUE}🤖 Running in CI mode${NC}"
        check_system_health
        cleanup_results
        run_all_visual_tests
        return $?
    fi
    
    # Check for command line arguments
    if [ $# -gt 0 ]; then
        case "$1" in
            "dashboard")
                check_system_health
                run_visual_test "visual-tests/gaming-analytics-dashboard.visual.test.ts" "Gaming Analytics Dashboard"
                ;;
            "match")
                check_system_health
                run_visual_test "visual-tests/match-analysis.visual.test.ts" "Match Analysis Interface"
                ;;
            "components")
                check_system_health
                run_visual_test "visual-tests/gaming-components.visual.test.ts" "Gaming Components"
                ;;
            "all")
                check_system_health
                cleanup_results
                run_all_visual_tests
                ;;
            "update")
                check_system_health
                update_visual_baselines
                ;;
            "report")
                show_visual_report
                ;;
            "debug")
                check_system_health
                interactive_debug
                ;;
            "install")
                install_prerequisites
                ;;
            "clean")
                cleanup_results
                ;;
            *)
                echo -e "${RED}Unknown command: $1${NC}"
                echo -e "${BLUE}Usage: $0 [dashboard|match|components|all|update|report|debug|install|clean]${NC}"
                exit 1
                ;;
        esac
        return $?
    fi
    
    # Interactive mode
    check_system_health || return 1
    
    while true; do
        show_menu
        read -r choice
        
        case $choice in
            1)
                run_visual_test "visual-tests/gaming-analytics-dashboard.visual.test.ts" "Gaming Analytics Dashboard"
                ;;
            2)
                run_visual_test "visual-tests/match-analysis.visual.test.ts" "Match Analysis Interface"
                ;;
            3)
                run_visual_test "visual-tests/gaming-components.visual.test.ts" "Gaming Components"
                ;;
            4)
                cleanup_results
                run_all_visual_tests
                ;;
            5)
                update_visual_baselines
                ;;
            6)
                show_visual_report
                ;;
            7)
                interactive_debug
                ;;
            8)
                echo -e "${GREEN}🎮 Herald.lol Visual Regression Testing Complete!${NC}"
                break
                ;;
            *)
                echo -e "${RED}Invalid choice. Please select 1-8.${NC}"
                ;;
        esac
        
        echo ""
    done
}

# Run main function
main "$@"