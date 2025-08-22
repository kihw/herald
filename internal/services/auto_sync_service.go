package services

import (
	"database/sql"
	"log"
	"time"
)

// AutoSyncService gère les synchronisations automatiques
type AutoSyncService struct {
	db           *sql.DB
	matchService *MatchService
	ticker       *time.Ticker
	stopChan     chan bool
	running      bool
}

// NewAutoSyncService crée un nouveau service de synchronisation automatique
func NewAutoSyncService(db *sql.DB, matchService *MatchService) *AutoSyncService {
	return &AutoSyncService{
		db:           db,
		matchService: matchService,
		stopChan:     make(chan bool),
		running:      false,
	}
}

// Start démarre le service de synchronisation automatique
func (as *AutoSyncService) Start() {
	if as.running {
		return
	}

	as.running = true
	as.ticker = time.NewTicker(1 * time.Hour) // Vérifier toutes les heures

	go func() {
		log.Println("🔄 Auto-sync service started")
		
		for {
			select {
			case <-as.ticker.C:
				as.checkAndSyncUsers()
			case <-as.stopChan:
				log.Println("🛑 Auto-sync service stopped")
				return
			}
		}
	}()
}

// Stop arrête le service de synchronisation automatique
func (as *AutoSyncService) Stop() {
	if !as.running {
		return
	}

	as.running = false
	if as.ticker != nil {
		as.ticker.Stop()
	}
	as.stopChan <- true
}

// checkAndSyncUsers vérifie quels utilisateurs ont besoin d'une synchronisation
func (as *AutoSyncService) checkAndSyncUsers() {
	log.Println("🔍 Checking users for auto-sync...")

	query := `
		SELECT 
			u.id, 
			u.riot_id, 
			u.riot_tag,
			u.last_sync,
			COALESCE(us.auto_sync_enabled, 1) as auto_sync_enabled,
			COALESCE(us.sync_frequency_hours, 24) as sync_frequency_hours
		FROM users u
		LEFT JOIN user_settings us ON u.id = us.user_id
		WHERE u.is_validated = 1
		AND COALESCE(us.auto_sync_enabled, 1) = 1
	`

	rows, err := as.db.Query(query)
	if err != nil {
		log.Printf("❌ Error querying users for auto-sync: %v", err)
		return
	}
	defer rows.Close()

	usersToSync := 0
	for rows.Next() {
		var userID int
		var riotID, riotTag string
		var lastSync sql.NullTime
		var autoSyncEnabled bool
		var syncFrequencyHours int

		err := rows.Scan(&userID, &riotID, &riotTag, &lastSync, &autoSyncEnabled, &syncFrequencyHours)
		if err != nil {
			log.Printf("❌ Error scanning user row: %v", err)
			continue
		}

		if !autoSyncEnabled {
			continue
		}

		// Vérifier si une synchronisation est nécessaire
		if as.needsSync(lastSync, syncFrequencyHours) {
			go as.syncUser(userID, riotID, riotTag)
			usersToSync++
		}
	}

	if usersToSync > 0 {
		log.Printf("🔄 Initiated auto-sync for %d users", usersToSync)
	} else {
		log.Println("✅ No users need auto-sync at this time")
	}
}

// needsSync détermine si un utilisateur a besoin d'une synchronisation
func (as *AutoSyncService) needsSync(lastSync sql.NullTime, frequencyHours int) bool {
	if !lastSync.Valid {
		return true // Jamais synchronisé
	}

	nextSyncTime := lastSync.Time.Add(time.Duration(frequencyHours) * time.Hour)
	return time.Now().After(nextSyncTime)
}

// syncUser effectue la synchronisation pour un utilisateur
func (as *AutoSyncService) syncUser(userID int, riotID, riotTag string) {
	log.Printf("🔄 Starting auto-sync for user %s#%s (ID: %d)", riotID, riotTag, userID)

	// Vérifier s'il y a déjà une synchronisation en cours
	if as.hasPendingSync(userID) {
		log.Printf("⏳ User %d already has a pending sync, skipping", userID)
		return
	}

	// Lancer la synchronisation via le match service
	_, err := as.matchService.SyncUserMatches(userID, 20) // Synchroniser 20 matches récents
	if err != nil {
		log.Printf("❌ Auto-sync failed for user %d: %v", userID, err)
		return
	}

	log.Printf("✅ Auto-sync initiated successfully for user %d", userID)
}

// hasPendingSync vérifie si un utilisateur a déjà une synchronisation en cours
func (as *AutoSyncService) hasPendingSync(userID int) bool {
	query := `
		SELECT COUNT(*) 
		FROM sync_jobs 
		WHERE user_id = ? 
		AND status IN ('pending', 'running')
		AND started_at > datetime('now', '-1 hour')
	`

	var count int
	err := as.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		log.Printf("❌ Error checking pending sync for user %d: %v", userID, err)
		return false
	}

	return count > 0
}

// GetSyncStats retourne des statistiques sur les synchronisations automatiques
func (as *AutoSyncService) GetSyncStats() map[string]interface{} {
	stats := map[string]interface{}{
		"service_running": as.running,
		"last_check":      time.Now().Format("2006-01-02 15:04:05"),
	}

	// Compter les utilisateurs avec auto-sync activé
	query := `
		SELECT COUNT(*) 
		FROM users u
		LEFT JOIN user_settings us ON u.id = us.user_id
		WHERE u.is_validated = 1
		AND COALESCE(us.auto_sync_enabled, 1) = 1
	`

	var autoSyncUsers int
	err := as.db.QueryRow(query).Scan(&autoSyncUsers)
	if err == nil {
		stats["auto_sync_users"] = autoSyncUsers
	}

	// Compter les synchronisations récentes
	recentSyncQuery := `
		SELECT COUNT(*) 
		FROM sync_jobs 
		WHERE started_at > datetime('now', '-24 hours')
		AND job_type = 'match_sync'
	`

	var recentSyncs int
	err = as.db.QueryRow(recentSyncQuery).Scan(&recentSyncs)
	if err == nil {
		stats["syncs_last_24h"] = recentSyncs
	}

	return stats
}

// IsRunning retourne true si le service est en cours d'exécution
func (as *AutoSyncService) IsRunning() bool {
	return as.running
}

// ForceCheck force une vérification immédiate des utilisateurs
func (as *AutoSyncService) ForceCheck() {
	if !as.running {
		log.Println("❌ Auto-sync service is not running")
		return
	}

	log.Println("🔄 Forcing auto-sync check...")
	go as.checkAndSyncUsers()
}