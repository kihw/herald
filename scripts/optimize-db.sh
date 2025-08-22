#!/bin/bash

# Script d'optimisation des bases de données Herald.lol
# Ce script optimise les bases de données SQLite et nettoie les anciens logs

DB_DIR="/home/debian/herald/data"
LOG_DIR="/home/debian/herald/logs"
EXPORTS_DIR="/home/debian/herald/exports"
BACKUP_DIR="/home/debian/herald/backups"

# Fonction de logging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Créer le répertoire de backup
mkdir -p "$BACKUP_DIR"

# Optimiser les bases de données SQLite
optimize_databases() {
    log "Optimisation des bases de données SQLite..."
    
    for db in "$DB_DIR"/*.db; do
        if [ -f "$db" ]; then
            local db_name=$(basename "$db")
            log "Optimisation de $db_name"
            
            # Backup avant optimisation
            cp "$db" "$BACKUP_DIR/${db_name}.backup.$(date +%Y%m%d)"
            
            # Optimisations SQLite
            sqlite3 "$db" "VACUUM;"
            sqlite3 "$db" "ANALYZE;"
            sqlite3 "$db" "PRAGMA optimize;"
            
            log "Optimisation de $db_name terminée"
        fi
    done
}

# Nettoyer les anciens logs
cleanup_logs() {
    log "Nettoyage des anciens logs..."
    
    # Garder seulement les logs des 7 derniers jours
    find "$LOG_DIR" -name "*.log" -mtime +7 -delete 2>/dev/null || true
    
    # Compresser les logs de plus de 1 jour
    find "$LOG_DIR" -name "*.log" -mtime +1 -exec gzip {} \; 2>/dev/null || true
    
    log "Nettoyage des logs terminé"
}

# Nettoyer les anciens exports
cleanup_exports() {
    log "Nettoyage des anciens exports..."
    
    # Supprimer les exports de plus de 30 jours
    find "$EXPORTS_DIR" -name "*.zip" -mtime +30 -delete 2>/dev/null || true
    find "$EXPORTS_DIR" -name "*.csv" -mtime +30 -delete 2>/dev/null || true
    find "$EXPORTS_DIR" -name "*.xlsx" -mtime +30 -delete 2>/dev/null || true
    
    log "Nettoyage des exports terminé"
}

# Nettoyer les anciens backups
cleanup_backups() {
    log "Nettoyage des anciens backups..."
    
    # Garder seulement les backups des 30 derniers jours
    find "$BACKUP_DIR" -name "*.backup.*" -mtime +30 -delete 2>/dev/null || true
    
    log "Nettoyage des backups terminé"
}

# Afficher les statistiques d'utilisation
show_stats() {
    log "Statistiques d'utilisation:"
    
    # Taille des bases de données
    if [ -d "$DB_DIR" ]; then
        local db_size=$(du -sh "$DB_DIR" 2>/dev/null | cut -f1)
        log "  Bases de données: $db_size"
    fi
    
    # Taille des logs
    if [ -d "$LOG_DIR" ]; then
        local log_size=$(du -sh "$LOG_DIR" 2>/dev/null | cut -f1)
        log "  Logs: $log_size"
    fi
    
    # Taille des exports
    if [ -d "$EXPORTS_DIR" ]; then
        local export_size=$(du -sh "$EXPORTS_DIR" 2>/dev/null | cut -f1)
        log "  Exports: $export_size"
    fi
    
    # Utilisation disque globale
    local disk_usage=$(df -h / | awk 'NR==2 {print $5}')
    log "  Utilisation disque: $disk_usage"
}

# Fonction principale
main() {
    log "=== Début de l'optimisation Herald.lol ==="
    
    show_stats
    
    optimize_databases
    cleanup_logs
    cleanup_exports
    cleanup_backups
    
    log "=== Optimisation terminée ==="
    show_stats
}

# Exécuter l'optimisation
main "$@"