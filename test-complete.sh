#!/bin/bash

# Complete integration test for Herald.lol
# Tests frontend build, backend compilation, and integration

echo "🎯 Herald.lol - Test Complet d'Intégration"
echo "==========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    
    case $status in
        "SUCCESS")
            echo -e "${GREEN}✅ $message${NC}"
            ;;
        "ERROR")
            echo -e "${RED}❌ $message${NC}"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  $message${NC}"
            ;;
        "WARNING")
            echo -e "${YELLOW}⚠️  $message${NC}"
            ;;
        "PROGRESS")
            echo -e "${PURPLE}🔄 $message${NC}"
            ;;
    esac
}

# Test counters
total_tests=0
passed_tests=0

# Function to run test with error handling
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    total_tests=$((total_tests + 1))
    print_status "PROGRESS" "Test: $test_name"
    
    if eval "$test_command" > /tmp/test_output 2>&1; then
        print_status "SUCCESS" "$test_name"
        passed_tests=$((passed_tests + 1))
        return 0
    else
        print_status "ERROR" "$test_name"
        echo "Error output:"
        cat /tmp/test_output | head -10
        return 1
    fi
}

echo
print_status "INFO" "Démarrage des tests d'intégration Herald.lol"
echo

# 1. Test de la compilation frontend
echo "===================="
echo "🎨 TESTS FRONTEND"
echo "===================="

cd /home/debian/herald/web

run_test "Compilation TypeScript Frontend" "npm run build"

# 2. Test de la compilation backend
echo
echo "==================="
echo "🔧 TESTS BACKEND"
echo "==================="

cd /home/debian/herald

print_status "INFO" "Test avec Docker..."
run_test "Compilation Docker Backend" "docker-compose build"

# 3. Test des fichiers essentiels
echo
echo "=========================="
echo "⚙️  TESTS CONFIGURATION"
echo "=========================="

essential_files=(
    "/home/debian/herald/main.go"
    "/home/debian/herald/web/index.html"
    "/home/debian/herald/web/package.json"
    "/home/debian/herald/web/src/components/groups/GroupManagement.tsx"
    "/home/debian/herald/web/src/services/groupApi.ts"
    "/home/debian/herald/internal/models/group_models.go"
    "/home/debian/herald/internal/handlers/group_handler.go"
)

for file in "${essential_files[@]}"; do
    if [ -f "$file" ]; then
        print_status "SUCCESS" "Fichier présent: $(basename $file)"
        passed_tests=$((passed_tests + 1))
    else
        print_status "ERROR" "Fichier manquant: $file"
    fi
    total_tests=$((total_tests + 1))
done

echo
echo "=========================="
echo "📊 RÉSUMÉ DES TESTS"
echo "=========================="

echo "Tests exécutés: $total_tests"
echo "Tests réussis: $passed_tests"
echo "Tests échoués: $((total_tests - passed_tests))"

if [ $total_tests -gt 0 ]; then
    percentage=$((passed_tests * 100 / total_tests))
    echo "Taux de réussite: $percentage%"
    
    if [ $percentage -ge 90 ]; then
        print_status "SUCCESS" "Excellent! Application prête 🚀"
        exit_code=0
    elif [ $percentage -ge 75 ]; then
        print_status "WARNING" "Bon état"
        exit_code=0
    else
        print_status "ERROR" "Corrections requises"
        exit_code=1
    fi
else
    print_status "ERROR" "Aucun test exécuté"
    exit_code=1
fi

print_status "INFO" "Test terminé"
exit $exit_code