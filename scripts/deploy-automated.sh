#!/bin/bash

# Script de déploiement automatisé pour LoL Match Exporter
# Usage: ./deploy.sh [environment] [build_number]

set -e  # Arrêt en cas d'erreur

# Configuration
DEPLOY_HOST="51.178.17.78"
DEPLOY_USER="debian"
APP_NAME="lol-fullstack-app"
DOCKER_COMPOSE_FILE="docker-compose.complete.yml"

# Couleurs pour les logs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction de logging
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
    exit 1
}

# Vérification des prérequis
check_prerequisites() {
    log "Vérification des prérequis..."
    
    # Vérifier Docker
    if ! command -v docker &> /dev/null; then
        error "Docker n'est pas installé"
    fi
    
    # Vérifier la connexion SSH
    if ! ssh -o ConnectTimeout=5 ${DEPLOY_USER}@${DEPLOY_HOST} "echo 'SSH OK'" &> /dev/null; then
        error "Impossible de se connecter en SSH à ${DEPLOY_HOST}"
    fi
    
    # Vérifier Docker sur le serveur distant
    if ! ssh ${DEPLOY_USER}@${DEPLOY_HOST} "docker --version" &> /dev/null; then
        error "Docker n'est pas disponible sur le serveur distant"
    fi
    
    success "Prérequis vérifiés"
}

# Construction de l'image Docker
build_image() {
    local build_number=${1:-"latest"}
    log "Construction de l'image Docker (build: ${build_number})..."
    
    # Build de l'image
    docker build -f Dockerfile.simple-fullstack -t lol-match-exporter-fullstack:${build_number} .
    docker tag lol-match-exporter-fullstack:${build_number} lol-match-exporter-fullstack:latest
    
    success "Image Docker construite"
}

# Tests de base
run_tests() {
    log "Exécution des tests..."
    
    # Test de compilation Go
    if [ -f "main.go" ]; then
        go mod tidy
        go build -o main .
        success "Compilation Go réussie"
    fi
    
    # Tests unitaires
    if go test ./... &> /dev/null; then
        success "Tests unitaires passés"
    else
        warning "Certains tests unitaires ont échoué"
    fi
}

# Déploiement sur le serveur
deploy_to_server() {
    local environment=${1:-"production"}
    local build_number=${2:-"latest"}
    
    log "Déploiement vers ${environment} (build: ${build_number})..."
    
    # Sauvegarder et transférer l'image Docker
    log "Transfert de l'image Docker..."
    docker save lol-match-exporter-fullstack:latest | ssh ${DEPLOY_USER}@${DEPLOY_HOST} 'docker load'
    
    # Transférer les fichiers de configuration
    log "Transfert des fichiers de configuration..."
    scp ${DOCKER_COMPOSE_FILE} ${DEPLOY_USER}@${DEPLOY_HOST}:~/docker-compose.yml
    
    if [ -f "nginx/nginx-fullstack.conf" ]; then
        scp nginx/nginx-fullstack.conf ${DEPLOY_USER}@${DEPLOY_HOST}:~/nginx.conf
    fi
    
    if [ -f "scripts/start-fullstack.sh" ]; then
        scp scripts/start-fullstack.sh ${DEPLOY_USER}@${DEPLOY_HOST}:~/start.sh
    fi
    
    # Exécuter le déploiement sur le serveur distant
    log "Exécution du déploiement distant..."
    ssh ${DEPLOY_USER}@${DEPLOY_HOST} bash << 'EOF'
        set -e
        
        echo "Arrêt de l'ancienne version..."
        docker-compose down --remove-orphans 2>/dev/null || true
        
        echo "Nettoyage des anciens conteneurs..."
        docker container prune -f
        
        echo "Démarrage de la nouvelle version..."
        if [ -f "start.sh" ]; then
            chmod +x start.sh
            ./start.sh
        else
            docker-compose up -d
        fi
        
        echo "Attente du démarrage..."
        sleep 15
        
        echo "Vérification du conteneur..."
        if docker ps | grep lol-fullstack-app > /dev/null; then
            echo "Conteneur démarré avec succès"
        else
            echo "Erreur: le conteneur n'a pas démarré"
            docker logs lol-fullstack-app 2>/dev/null || true
            exit 1
        fi
EOF
    
    success "Déploiement terminé"
}

# Tests de santé post-déploiement
health_check() {
    log "Vérification de la santé de l'application..."
    
    # Attendre que l'application soit prête
    sleep 10
    
    # Test de l'endpoint de santé
    for i in {1..5}; do
        if curl -f -s https://herald.lol/api/health > /dev/null 2>&1; then
            success "Endpoint de santé OK"
            break
        else
            warning "Tentative ${i}/5 - endpoint de santé non disponible"
            sleep 5
        fi
        
        if [ $i -eq 5 ]; then
            error "Endpoint de santé non accessible après 5 tentatives"
        fi
    done
    
    # Test des endpoints critiques
    local endpoints=(
        "https://herald.lol/api/auth/session"
        "https://herald.lol/api/auth/regions"
    )
    
    for endpoint in "${endpoints[@]}"; do
        if curl -f -s "$endpoint" > /dev/null 2>&1; then
            success "Endpoint ${endpoint} OK"
        else
            warning "Endpoint ${endpoint} non accessible"
        fi
    done
}

# Rollback en cas de problème
rollback() {
    log "Rollback vers la version précédente..."
    
    ssh ${DEPLOY_USER}@${DEPLOY_HOST} bash << 'EOF'
        set -e
        
        echo "Recherche de la version précédente..."
        LAST_STABLE=$(docker images --format "table {{.Repository}}:{{.Tag}}" | \
                     grep lol-match-exporter-fullstack | \
                     grep -v latest | \
                     head -1)
        
        if [ -n "$LAST_STABLE" ]; then
            echo "Rollback vers $LAST_STABLE"
            
            docker-compose down
            docker tag $LAST_STABLE lol-match-exporter-fullstack:latest
            
            if [ -f "start.sh" ]; then
                ./start.sh
            else
                docker-compose up -d
            fi
            
            echo "Rollback terminé vers $LAST_STABLE"
        else
            echo "Aucune version précédente trouvée"
            exit 1
        fi
EOF
    
    success "Rollback effectué"
}

# Fonction principale
main() {
    local environment=${1:-"production"}
    local build_number=${2:-$(date +%Y%m%d-%H%M%S)}
    
    log "=== Déploiement LoL Match Exporter ==="
    log "Environnement: ${environment}"
    log "Build: ${build_number}"
    log "Serveur: ${DEPLOY_HOST}"
    
    # Étapes du déploiement
    check_prerequisites
    run_tests
    build_image "$build_number"
    deploy_to_server "$environment" "$build_number"
    health_check
    
    success "=== Déploiement réussi ==="
}

# Gestion des arguments
case "${1:-}" in
    "rollback")
        rollback
        ;;
    "health")
        health_check
        ;;
    "build")
        build_image "${2:-latest}"
        ;;
    *)
        main "$@"
        ;;
esac
