#!/bin/bash

# Script de monitoring Herald.lol
# Ce script vérifie l'état des services et génère des alertes

LOG_DIR="/home/debian/herald/logs"
ALERT_LOG="$LOG_DIR/alerts.log"
HEALTH_LOG="$LOG_DIR/health.log"

# Créer les répertoires de logs si nécessaires
mkdir -p "$LOG_DIR"

# Fonction de logging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$HEALTH_LOG"
}

alert() {
    echo "[ALERT $(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$ALERT_LOG"
}

# Vérifier le status des containers Docker
check_containers() {
    log "Vérification des containers Docker..."
    
    # Vérifier nginx
    if ! docker ps | grep -q "lol-nginx-proxy"; then
        alert "Container nginx n'est pas en cours d'exécution"
        return 1
    fi
    
    # Vérifier frontend
    if ! docker ps | grep -q "lol-frontend-production"; then
        alert "Container frontend n'est pas en cours d'exécution"
        return 1
    fi
    
    # Vérifier backend
    if ! docker ps | grep -q "lol-exporter-production"; then
        alert "Container backend n'est pas en cours d'exécution"
        return 1
    fi
    
    log "Tous les containers sont opérationnels"
    return 0
}

# Vérifier la connectivité HTTP/HTTPS
check_web_connectivity() {
    log "Vérification de la connectivité web..."
    
    # Vérifier HTTPS
    if ! curl -k -f -s https://herald.lol/ > /dev/null; then
        alert "Site HTTPS non accessible"
        return 1
    fi
    
    # Vérifier l'API
    if ! curl -k -f -s https://herald.lol/api/health > /dev/null 2>&1; then
        log "API health endpoint non accessible (peut être normal)"
    fi
    
    log "Connectivité web OK"
    return 0
}

# Vérifier l'utilisation des ressources
check_resources() {
    log "Vérification des ressources système..."
    
    # Utilisation disque
    DISK_USAGE=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ "$DISK_USAGE" -gt 85 ]; then
        alert "Utilisation disque élevée: ${DISK_USAGE}%"
    fi
    
    # Utilisation mémoire
    MEM_USAGE=$(free | awk '/^Mem:/ {printf "%.0f", $3/$2 * 100}')
    if [ "$MEM_USAGE" -gt 85 ]; then
        alert "Utilisation mémoire élevée: ${MEM_USAGE}%"
    fi
    
    log "Ressources: Disque ${DISK_USAGE}%, Mémoire ${MEM_USAGE}%"
    return 0
}

# Fonction principale
main() {
    log "=== Début du monitoring Herald.lol ==="
    
    local errors=0
    
    check_containers || ((errors++))
    check_web_connectivity || ((errors++))
    check_resources || ((errors++))
    
    if [ $errors -eq 0 ]; then
        log "=== Monitoring terminé: Tout est OK ==="
    else
        alert "=== Monitoring terminé avec $errors erreurs ==="
    fi
    
    return $errors
}

# Exécuter le monitoring
main "$@"