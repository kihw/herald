package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"lol-match-exporter/internal/auth"
	"lol-match-exporter/internal/handlers"
	"lol-match-exporter/internal/models"
	"lol-match-exporter/internal/services"
)

type ExportJob struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	LogLines    []string   `json:"log_lines"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Error       string     `json:"error,omitempty"`
	ZipPath     string     `json:"zip_path,omitempty"`
	cmd         *exec.Cmd
	ctx         context.Context
	cancel      context.CancelFunc
	logChan     chan string
	clientChans []chan string
	mutex       sync.RWMutex
}

type ExportRequest struct {
	Username  string `json:"username"`
	TagLine   string `json:"tagline"`
	RiotId    string `json:"riotId"`
	GameCount int    `json:"gameCount"`
	Count     int    `json:"count"`
	APIKey    string `json:"apiKey"`
}

type UserSession struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Server struct {
	jobs             map[string]*ExportJob
	mutex            sync.RWMutex
	exportService    *services.ExportService
	templateService  *services.TemplateService
	riotAPIService   *services.RiotAPIService
	googleOAuth      *auth.GoogleOAuthService
	db               *sql.DB
	groupHandler     *handlers.GroupHandler
	matchService     *services.MatchService
	analyticsService *services.SimpleAnalyticsService
	autoSyncService  *services.AutoSyncService
	sessions         map[string]*UserSession
	sessionMutex    sync.RWMutex
}

func NewServer() *Server {
	// Initialiser la base de données SQLite
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Récupérer la clé API Riot depuis les variables d'environnement
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	var riotService *services.RiotAPIService
	if riotAPIKey != "" {
		riotService = services.NewRiotAPIService(riotAPIKey)
	}

	// Initialiser un service d'analytics simple pour SQLite
	analyticsService := services.NewSimpleAnalyticsService(db)
	
	// Initialiser le service de matches
	matchService := services.NewMatchService(db, riotService)
	
	// Initialiser le service de synchronisation automatique
	autoSyncService := services.NewAutoSyncService(db, matchService)

	return &Server{
		jobs:             make(map[string]*ExportJob),
		exportService:    services.NewExportService("./exports"),
		templateService:  services.NewTemplateService(),
		riotAPIService:   riotService,
		googleOAuth:      auth.NewGoogleOAuthService(),
		db:               db,
		groupHandler:     handlers.NewGroupHandler(db, nil),
		matchService:     matchService,
		analyticsService: analyticsService,
		autoSyncService:  autoSyncService,
		sessions:         make(map[string]*UserSession),
	}
}

func generateJobID() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateSessionID génère un ID de session unique
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// createUserSession crée une nouvelle session pour un utilisateur
func (s *Server) createUserSession(userID int) (string, error) {
	sessionID := generateSessionID()
	now := time.Now()
	session := &UserSession{
		SessionID: sessionID,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(7 * 24 * time.Hour), // 7 jours
	}
	
	s.sessionMutex.Lock()
	s.sessions[sessionID] = session
	s.sessionMutex.Unlock()
	
	return sessionID, nil
}

// getUserFromSession récupère un utilisateur à partir d'un ID de session
func (s *Server) getUserFromSession(sessionID string) (*models.User, error) {
	s.sessionMutex.RLock()
	session, exists := s.sessions[sessionID]
	s.sessionMutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	
	if time.Now().After(session.ExpiresAt) {
		// Session expirée
		s.sessionMutex.Lock()
		delete(s.sessions, sessionID)
		s.sessionMutex.Unlock()
		return nil, fmt.Errorf("session expired")
	}
	
	// Récupérer l'utilisateur depuis la base de données
	query := `
		SELECT id, riot_id, riot_tag, riot_puuid, summoner_id, account_id, 
			   profile_icon_id, summoner_level, region, is_validated, 
			   created_at, updated_at, last_sync
		FROM users 
		WHERE id = ?
	`
	
	var user models.User
	var summonerID, accountID sql.NullString
	var lastSync sql.NullTime
	
	err := s.db.QueryRow(query, session.UserID).Scan(
		&user.ID, &user.RiotID, &user.RiotTag, &user.RiotPUUID,
		&summonerID, &accountID, &user.ProfileIconID, &user.SummonerLevel,
		&user.Region, &user.IsValidated, &user.CreatedAt, &user.UpdatedAt,
		&lastSync,
	)
	
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	// Gérer les valeurs nullables
	if summonerID.Valid {
		user.SummonerID = &summonerID.String
	}
	if accountID.Valid {
		user.AccountID = &accountID.String
	}
	if lastSync.Valid {
		user.LastSync = &lastSync.Time
	}
	
	return &user, nil
}

// requireAuth middleware pour vérifier l'authentification
func (s *Server) requireAuth(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		c.Abort()
		return
	}
	
	user, err := s.getUserFromSession(sessionID)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid session"})
		c.Abort()
		return
	}
	
	// Ajouter l'utilisateur au contexte
	c.Set("user", user)
	c.Next()
}

func (s *Server) startExport(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Parse RiotId if provided (format: username#tagline)
	if req.RiotId != "" {
		parts := strings.Split(req.RiotId, "#")
		if len(parts) != 2 {
			c.JSON(400, gin.H{"error": "Invalid RiotId format. Expected: GameName#Tag"})
			return
		}
		req.Username = parts[0]
		req.TagLine = parts[1]
	}

	// Validate required fields
	if req.Username == "" || req.TagLine == "" {
		c.JSON(400, gin.H{"error": "Username and tagline are required (provide riotId or separate username/tagline)"})
		return
	}

	// Use count if provided, otherwise use gameCount, default to 1000
	if req.Count > 0 {
		req.GameCount = req.Count
	} else if req.GameCount == 0 {
		req.GameCount = 1000
	}

	jobID := generateJobID()
	ctx, cancel := context.WithCancel(context.Background())

	job := &ExportJob{
		ID:          jobID,
		Status:      "running",
		Progress:    0,
		LogLines:    []string{},
		StartTime:   time.Now(),
		ctx:         ctx,
		cancel:      cancel,
		logChan:     make(chan string, 100),
		clientChans: []chan string{},
	}

	s.mutex.Lock()
	s.jobs[jobID] = job
	s.mutex.Unlock()

	log.Printf("[SERVER] Starting job %s for %s#%s", jobID, req.Username, req.TagLine)

	// Start the export process
	go s.runExportProcess(job, req)

	// Start log processing
	go s.processLogs(job)

	c.JSON(200, gin.H{"job_id": jobID})
}

func (s *Server) runExportProcess(job *ExportJob, req ExportRequest) {
	defer func() {
		job.cancel()
		close(job.logChan)
	}()

	// Log de début
	select {
	case job.logChan <- fmt.Sprintf("[INFO] Démarrage de l'export pour %s#%s", req.Username, req.TagLine):
	case <-job.ctx.Done():
		return
	}

	job.mutex.Lock()
	job.Progress = 10
	job.mutex.Unlock()

	// Récupérer les données via l'API Riot ou générer des données de test
	var matches []services.MatchData
	var err error

	// Pour le moment, utilisons toujours des données de test pour la démonstration
	// if s.riotAPIService != nil {
	if false && s.riotAPIService != nil {
		select {
		case job.logChan <- "[INFO] Récupération des données via l'API Riot Games...":
		case <-job.ctx.Done():
			return
		}

		matches, err = s.riotAPIService.GetMatchesData(req.Username, req.TagLine, req.GameCount, nil)
		if err != nil {
			s.setJobError(job, fmt.Sprintf("Erreur API Riot: %s", err.Error()))
			select {
			case job.logChan <- fmt.Sprintf("[ERROR] Échec de récupération des données: %s", err.Error()):
			case <-job.ctx.Done():
			}
			return
		}
	} else {
		select {
		case job.logChan <- fmt.Sprintf("[INFO] Génération de %d matchs de test pour %s#%s", req.GameCount, req.Username, req.TagLine):
		case <-job.ctx.Done():
			return
		}
		matches = s.generateTestMatches()
	}

	job.mutex.Lock()
	job.Progress = 50
	job.mutex.Unlock()

	select {
	case job.logChan <- fmt.Sprintf("[INFO] %d matchs récupérés, génération du fichier d'export...", len(matches)):
	case <-job.ctx.Done():
		return
	}

	// Créer les options d'export par défaut
	options := services.ExportOptions{
		Format:   services.FormatCSV,
		Filename: fmt.Sprintf("%s_%s_matches", req.Username, req.TagLine),
		Filter: services.ExportFilter{
			RecentFirst: true,
		},
		Compression: true,
		Metadata:    true,
	}

	// Exporter les données
	exportPath, err := s.exportService.ExportMatches(matches, options)

	job.mutex.Lock()
	defer job.mutex.Unlock()

	if err != nil {
		if job.ctx.Err() == context.Canceled {
			job.Status = "cancelled"
		} else {
			job.Status = "failed"
			job.Error = err.Error()
		}
		select {
		case job.logChan <- fmt.Sprintf("[ERROR] Échec de l'export: %s", err.Error()):
		case <-job.ctx.Done():
		}
	} else {
		job.Status = "completed"
		job.Progress = 100
		job.ZipPath = exportPath

		select {
		case job.logChan <- fmt.Sprintf("[SUCCESS] Export terminé: %s", exportPath):
		case <-job.ctx.Done():
		}
	}

	now := time.Now()
	job.EndTime = &now
}

func (s *Server) readOutput(job *ExportJob, reader io.Reader, source string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// Send to log channel
		select {
		case job.logChan <- fmt.Sprintf("[%s] %s", source, line):
		case <-job.ctx.Done():
			return
		}

		// Parse progress if possible
		s.parseProgress(job, line)
	}
}

func (s *Server) parseProgress(job *ExportJob, line string) {
	// Simple progress parsing - you can enhance this based on your Python script output
	if strings.Contains(line, "Progress:") {
		// Extract progress percentage from line
		// This is a simplified example
		job.mutex.Lock()
		if strings.Contains(line, "25%") {
			job.Progress = 25
		} else if strings.Contains(line, "50%") {
			job.Progress = 50
		} else if strings.Contains(line, "75%") {
			job.Progress = 75
		}
		job.mutex.Unlock()
	}
}

func (s *Server) processLogs(job *ExportJob) {
	for {
		select {
		case line, ok := <-job.logChan:
			if !ok {
				return
			}

			job.mutex.Lock()
			job.LogLines = append(job.LogLines, line)
			// Keep only last 1000 lines
			if len(job.LogLines) > 1000 {
				job.LogLines = job.LogLines[len(job.LogLines)-1000:]
			}

			// Send to all connected clients
			for _, clientChan := range job.clientChans {
				select {
				case clientChan <- line:
				default:
					// Client channel is full, skip
				}
			}
			job.mutex.Unlock()

		case <-job.ctx.Done():
			return
		}
	}
}

func (s *Server) setJobError(job *ExportJob, errorMsg string) {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	job.Status = "failed"
	job.Error = errorMsg
	now := time.Now()
	job.EndTime = &now

	log.Printf("[SERVER] Job %s failed: %s", job.ID, errorMsg)
}

func (s *Server) getJobStatus(c *gin.Context) {
	jobID := c.Param("job_id")

	s.mutex.RLock()
	job, exists := s.jobs[jobID]
	s.mutex.RUnlock()

	if !exists {
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	job.mutex.RLock()
	response := gin.H{
		"id":         job.ID,
		"status":     job.Status,
		"progress":   job.Progress,
		"log_lines":  job.LogLines,
		"start_time": job.StartTime,
	}

	if job.EndTime != nil {
		response["end_time"] = job.EndTime
	}

	if job.Error != "" {
		response["error"] = job.Error
	}

	if job.ZipPath != "" {
		response["zip_path"] = job.ZipPath
	}
	job.mutex.RUnlock()

	c.JSON(200, response)
}

func (s *Server) streamLogs(c *gin.Context) {
	jobID := c.Param("job_id")

	s.mutex.RLock()
	job, exists := s.jobs[jobID]
	s.mutex.RUnlock()

	if !exists {
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create client channel
	clientChan := make(chan string, 100)

	job.mutex.Lock()
	job.clientChans = append(job.clientChans, clientChan)

	// Send existing logs
	for _, line := range job.LogLines {
		select {
		case clientChan <- line:
		default:
			break
		}
	}
	job.mutex.Unlock()

	// Stream logs
	clientGone := c.Request.Context().Done()
	ticker := time.NewTicker(30 * time.Second) // Heartbeat
	defer ticker.Stop()

	for {
		select {
		case line := <-clientChan:
			data, _ := json.Marshal(gin.H{"log": line})
			c.SSEvent("log", string(data))
			c.Writer.Flush()

		case <-ticker.C:
			// Send heartbeat
			c.SSEvent("heartbeat", time.Now().Unix())
			c.Writer.Flush()

		case <-clientGone:
			// Remove client channel
			job.mutex.Lock()
			for i, ch := range job.clientChans {
				if ch == clientChan {
					job.clientChans = append(job.clientChans[:i], job.clientChans[i+1:]...)
					break
				}
			}
			job.mutex.Unlock()
			close(clientChan)
			return

		case <-job.ctx.Done():
			return
		}
	}
}

func (s *Server) downloadZip(c *gin.Context) {
	jobID := c.Param("job_id")
	filename := c.Param("filename") // Optionnel pour compatibilité frontend

	s.mutex.RLock()
	job, exists := s.jobs[jobID]
	s.mutex.RUnlock()

	if !exists {
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}

	job.mutex.RLock()
	zipPath := job.ZipPath
	job.mutex.RUnlock()

	if zipPath == "" {
		c.JSON(404, gin.H{"error": "No zip file available"})
		return
	}

	// Si un filename spécifique est demandé, vérifier qu'il correspond
	if filename != "" {
		expectedFilename := filepath.Base(zipPath)
		if filename != expectedFilename {
			c.JSON(404, gin.H{"error": fmt.Sprintf("File %s not found, available: %s", filename, expectedFilename)})
			return
		}
	}

	// Convertir le chemin relatif en absolu si nécessaire
	if !filepath.IsAbs(zipPath) {
		zipPath = filepath.Join(".", zipPath)
	}
	
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": fmt.Sprintf("Zip file not found at path: %s", zipPath)})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(zipPath)))
	c.File(zipPath)
}

func (s *Server) findGeneratedZip(username, tagline string) string {
	// Look for zip files matching the pattern
	pattern := fmt.Sprintf("*%s*%s*.zip", username, tagline)
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}

	// Return the most recent file
	var newest string
	var newestTime time.Time

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if info.ModTime().After(newestTime) {
			newest = match
			newestTime = info.ModTime()
		}
	}

	return newest
}

func (s *Server) health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func (s *Server) cleanupOldJobs() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-1 * time.Hour)

		s.mutex.Lock()
		for jobID, job := range s.jobs {
			job.mutex.RLock()
			shouldDelete := job.EndTime != nil && job.EndTime.Before(cutoff)
			job.mutex.RUnlock()

			if shouldDelete {
				job.cancel()
				delete(s.jobs, jobID)
				log.Printf("[SERVER] Cleaned up old job %s", jobID)
			}
		}
		s.mutex.Unlock()
	}
}

func main() {
	server := NewServer()

	// Start cleanup routine
	go server.cleanupOldJobs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Serve static files
	r.Use(static.Serve("/", static.LocalFile("./web/dist", false)))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", server.health)

		// Auth endpoints
		api.GET("/auth/session", server.getSession)
		api.GET("/auth/regions", server.getSupportedRegions)
		api.POST("/auth/validate", server.validateAuth)
		api.POST("/auth/logout", server.logout)
		
		// Google OAuth endpoints (mock for now)
		api.POST("/auth/google/init", server.initGoogleOAuth)
		api.GET("/auth/google/callback", server.googleOAuthCallback)

		api.POST("/export", server.startExport)
		api.GET("/jobs/:job_id", server.getJobStatus)
		api.GET("/jobs/:job_id/logs", server.streamLogs)
		api.GET("/jobs/:job_id/download", server.downloadZip)

		// Advanced export endpoints
		api.POST("/export/advanced", server.startAdvancedExport)
		api.GET("/export/formats", server.getSupportedFormats)
		api.POST("/export/validate", server.validateExportOptions)
		api.GET("/export/history", server.getExportHistory)
		
		// Endpoints pour compatibilité frontend
		api.GET("/export/:job_id/status", server.getJobStatus)
		api.GET("/export/:job_id/events", server.streamLogs)
		api.GET("/export/:job_id/download/:filename", server.downloadZip)

		// Template endpoints
		api.GET("/templates", server.getAllTemplates)
		api.GET("/templates/:id", server.getTemplate)
		api.POST("/templates", server.createTemplate)
		api.PUT("/templates/:id", server.updateTemplate)
		api.DELETE("/templates/:id", server.deleteTemplate)
		api.POST("/export/from-template", server.exportFromTemplate)
		
		// Dashboard endpoint
		api.GET("/dashboard", server.getDashboard)
		
		// Match synchronization endpoints
		api.POST("/sync/matches", server.syncMatches)
		api.GET("/sync/status/:jobId", server.getSyncStatus)
		api.GET("/matches", server.getUserMatches)
		
		// Analytics endpoints
		api.GET("/analytics/period/:period", server.getPeriodAnalytics)
		api.GET("/analytics/recommendations", server.getRecommendations)
		api.GET("/analytics/trends", server.getPerformanceTrends)
		api.POST("/analytics/refresh", server.refreshAnalytics)
		
		// Settings endpoints
		api.GET("/settings", server.getUserSettings)
		api.PUT("/settings", server.updateUserSettings)

		// Group endpoints
		groups := api.Group("/groups")
		{
			groups.POST("/", server.groupHandler.CreateGroup)
			groups.GET("/search", server.groupHandler.SearchGroups)
			groups.GET("/my", server.groupHandler.GetUserGroups)
			groups.POST("/join", server.groupHandler.JoinGroup)
			
			groups.GET("/:id", server.groupHandler.GetGroup)
			groups.GET("/:id/members", server.groupHandler.GetGroupMembers)
			groups.POST("/:id/invite", server.groupHandler.InviteToGroup)
			groups.DELETE("/:id/members", server.groupHandler.RemoveMember)
			groups.GET("/:id/stats", server.groupHandler.GetGroupStats)
			groups.PUT("/:id/settings", server.groupHandler.UpdateGroupSettings)
			
			// Comparison endpoints
			groups.POST("/:id/comparisons", server.groupHandler.CreateComparison)
			groups.GET("/:id/comparisons", server.groupHandler.GetGroupComparisons)
			groups.GET("/:id/comparisons/:comparisonId", server.groupHandler.GetComparison)
			groups.POST("/:id/comparisons/:comparisonId/regenerate", server.groupHandler.RegenerateComparison)
		}
	}

	// Legacy compatibility routes (redirect to API)
	r.POST("/export", server.startExport)
	r.GET("/jobs/:job_id", server.getJobStatus)
	r.GET("/jobs/:job_id/logs", server.streamLogs)
	r.GET("/jobs/:job_id/download", server.downloadZip)

	// Fallback for SPA
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.File("./web/dist/index.html")
		} else {
			c.JSON(404, gin.H{"error": "Not found"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Démarrer le service de synchronisation automatique
	server.autoSyncService.Start()
	log.Println("🔄 Auto-sync service started")

	log.Printf("[SERVER] Starting Go backend on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// AdvancedExportRequest définit la structure de requête pour l'export avancé
type AdvancedExportRequest struct {
	Username string                 `json:"username" binding:"required"`
	TagLine  string                 `json:"tagline" binding:"required"`
	Options  services.ExportOptions `json:"options" binding:"required"`
}

// startAdvancedExport démarre un export avancé avec des options personnalisées
func (s *Server) startAdvancedExport(c *gin.Context) {
	var req AdvancedExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Valider les options d'export
	if err := s.exportService.ValidateOptions(req.Options); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jobID := generateJobID()
	ctx, cancel := context.WithCancel(context.Background())

	job := &ExportJob{
		ID:          jobID,
		Status:      "pending",
		Progress:    0,
		LogLines:    []string{},
		StartTime:   time.Now(),
		ctx:         ctx,
		cancel:      cancel,
		logChan:     make(chan string, 100),
		clientChans: []chan string{},
	}

	s.mutex.Lock()
	s.jobs[jobID] = job
	s.mutex.Unlock()

	// Démarre l'export en arrière-plan
	go s.runAdvancedExport(job, req)

	c.JSON(200, gin.H{"job_id": jobID})
}

// runAdvancedExport exécute l'export avancé
func (s *Server) runAdvancedExport(job *ExportJob, req AdvancedExportRequest) {
	job.mutex.Lock()
	job.Status = "running"
	job.Progress = 10
	job.mutex.Unlock()

	// Log de début
	select {
	case job.logChan <- fmt.Sprintf("[INFO] Démarrage de l'export avancé pour %s#%s", req.Username, req.TagLine):
	case <-job.ctx.Done():
		return
	}

	// Récupérer les données de matchs
	var matches []services.MatchData
	var err error

	if s.riotAPIService != nil {
		// Utiliser la vraie API Riot
		select {
		case job.logChan <- "[INFO] Récupération des données via l'API Riot Games...":
		case <-job.ctx.Done():
			return
		}

		// Extraire les queues du filtre si spécifiées
		var queueIDs []int
		if req.Options.Filter.Queues != nil {
			queueIDs = req.Options.Filter.Queues
		}

		matches, err = s.riotAPIService.GetMatchesData(req.Username, req.TagLine, 50, queueIDs)
		if err != nil {
			job.mutex.Lock()
			job.Status = "failed"
			job.Error = fmt.Sprintf("Erreur API Riot: %s", err.Error())
			job.mutex.Unlock()

			select {
			case job.logChan <- fmt.Sprintf("[ERROR] Échec de récupération des données: %s", err.Error()):
			case <-job.ctx.Done():
			}
			return
		}
	} else {
		// Utiliser des données de test
		select {
		case job.logChan <- "[INFO] Utilisation de données de test (clé API Riot non configurée)":
		case <-job.ctx.Done():
			return
		}
		matches = s.generateTestMatches()
	}

	job.mutex.Lock()
	job.Progress = 50
	job.mutex.Unlock()

	select {
	case job.logChan <- fmt.Sprintf("[INFO] Données récupérées, %d matchs trouvés", len(matches)):
	case <-job.ctx.Done():
		return
	}

	// Exporter les données
	exportPath, err := s.exportService.ExportMatches(matches, req.Options)

	job.mutex.Lock()
	defer job.mutex.Unlock()

	if err != nil {
		if job.ctx.Err() == context.Canceled {
			job.Status = "cancelled"
		} else {
			job.Status = "failed"
			job.Error = err.Error()
		}
		select {
		case job.logChan <- fmt.Sprintf("[ERROR] Échec de l'export: %s", err.Error()):
		case <-job.ctx.Done():
		}
	} else {
		job.Status = "completed"
		job.Progress = 100
		job.ZipPath = exportPath

		select {
		case job.logChan <- fmt.Sprintf("[SUCCESS] Export terminé: %s", exportPath):
		case <-job.ctx.Done():
		}
	}

	now := time.Now()
	job.EndTime = &now
}

// generateTestMatches génère des données de test pour l'export
func (s *Server) generateTestMatches() []services.MatchData {
	matches := make([]services.MatchData, 10)

	for i := 0; i < 10; i++ {
		lp := 1200 + i*50
		mmr := 1400 + i*60

		matches[i] = services.MatchData{
			MatchID:      fmt.Sprintf("EUW1_%d", 5000000000+i),
			GameCreation: time.Now().AddDate(0, 0, -i),
			GameDuration: 1800 + i*120,
			QueueID:      420, // Ranked Solo/Duo
			GameMode:     "CLASSIC",
			GameType:     "MATCHED_GAME",
			ChampionID:   1 + i,
			ChampionName: fmt.Sprintf("Champion%d", i+1),
			Role:         "MIDDLE",
			Lane:         "MIDDLE",
			Win:          i%2 == 0,
			Kills:        5 + i,
			Deaths:       2 + i/2,
			Assists:      8 + i*2,
			KDA:          float64(5+i+8+i*2) / float64(max(1, 2+i/2)),
			CS:           150 + i*10,
			Gold:         12000 + i*500,
			Damage:       25000 + i*1000,
			Vision:       15 + i,
			Items:        []int{3006, 3020, 3031, 3046, 3072, 3139},
			Summoners:    []int{4, 7},
			Rank:         "GOLD",
			LP:           &lp,
			MMR:          &mmr,
		}
	}

	return matches
}

// getSupportedFormats retourne les formats d'export supportés
func (s *Server) getSupportedFormats(c *gin.Context) {
	formats := s.exportService.GetSupportedFormats()
	c.JSON(200, gin.H{"formats": formats})
}

// validateExportOptions valide les options d'export
func (s *Server) validateExportOptions(c *gin.Context) {
	var options services.ExportOptions
	if err := c.ShouldBindJSON(&options); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := s.exportService.ValidateOptions(options); err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "valid": false})
		return
	}

	c.JSON(200, gin.H{"valid": true})
}

// getExportHistory retourne l'historique des exports
func (s *Server) getExportHistory(c *gin.Context) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var history []gin.H
	for _, job := range s.jobs {
		job.mutex.RLock()
		historyItem := gin.H{
			"id":         job.ID,
			"status":     job.Status,
			"start_time": job.StartTime,
			"end_time":   job.EndTime,
		}
		if job.ZipPath != "" {
			historyItem["download_available"] = true
		}
		job.mutex.RUnlock()
		history = append(history, historyItem)
	}

	c.JSON(200, gin.H{"history": history})
}

// getAllTemplates retourne tous les templates disponibles
func (s *Server) getAllTemplates(c *gin.Context) {
	templates := s.templateService.GetAllTemplates()
	c.JSON(200, gin.H{"templates": templates})
}

// getTemplate récupère un template spécifique
func (s *Server) getTemplate(c *gin.Context) {
	templateID := c.Param("id")
	template, exists := s.templateService.GetTemplate(templateID)

	if !exists {
		c.JSON(404, gin.H{"error": "Template non trouvé"})
		return
	}

	c.JSON(200, template)
}

// createTemplate crée un nouveau template
func (s *Server) createTemplate(c *gin.Context) {
	var template services.ExportTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Valider le template
	if err := s.templateService.ValidateTemplate(template); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Créer le template
	if err := s.templateService.CreateTemplate(template); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, template)
}

// updateTemplate met à jour un template existant
func (s *Server) updateTemplate(c *gin.Context) {
	templateID := c.Param("id")

	var template services.ExportTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Valider le template
	if err := s.templateService.ValidateTemplate(template); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Mettre à jour le template
	if err := s.templateService.UpdateTemplate(templateID, template); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, template)
}

// deleteTemplate supprime un template
func (s *Server) deleteTemplate(c *gin.Context) {
	templateID := c.Param("id")

	if err := s.templateService.DeleteTemplate(templateID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Template supprimé avec succès"})
}

// TemplateExportRequest définit la structure pour l'export depuis un template
type TemplateExportRequest struct {
	Username   string                  `json:"username" binding:"required"`
	TagLine    string                  `json:"tagline" binding:"required"`
	TemplateID string                  `json:"template_id" binding:"required"`
	Overrides  *services.ExportOptions `json:"overrides,omitempty"`
}

// exportFromTemplate démarre un export en utilisant un template
func (s *Server) exportFromTemplate(c *gin.Context) {
	var req TemplateExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Récupérer le template
	template, exists := s.templateService.GetTemplate(req.TemplateID)
	if !exists {
		c.JSON(404, gin.H{"error": "Template non trouvé"})
		return
	}

	// Appliquer les overrides si fournis
	var options services.ExportOptions
	if req.Overrides != nil {
		var err error
		options, err = s.templateService.ApplyTemplate(req.TemplateID, *req.Overrides)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		options = template.Options
	}

	// Créer la requête d'export avancé
	advancedReq := AdvancedExportRequest{
		Username: req.Username,
		TagLine:  req.TagLine,
		Options:  options,
	}

	// Valider les options
	if err := s.exportService.ValidateOptions(options); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jobID := generateJobID()
	ctx, cancel := context.WithCancel(context.Background())

	job := &ExportJob{
		ID:          jobID,
		Status:      "pending",
		Progress:    0,
		LogLines:    []string{},
		StartTime:   time.Now(),
		ctx:         ctx,
		cancel:      cancel,
		logChan:     make(chan string, 100),
		clientChans: []chan string{},
	}

	s.mutex.Lock()
	s.jobs[jobID] = job
	s.mutex.Unlock()

	// Démarre l'export en arrière-plan
	go s.runAdvancedExport(job, advancedReq)

	c.JSON(200, gin.H{
		"job_id":        jobID,
		"template_used": template.Name,
	})
}

// getSession retourne les informations de session
func (s *Server) getSession(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(200, gin.H{
			"authenticated": false,
			"user":          nil,
		})
		return
	}
	
	user, err := s.getUserFromSession(sessionID)
	if err != nil {
		c.JSON(200, gin.H{
			"authenticated": false,
			"user":          nil,
		})
		return
	}
	
	c.JSON(200, gin.H{
		"authenticated": true,
		"user":          user,
	})
}

// syncMatches lance la synchronisation des matches d'un utilisateur
func (s *Server) syncMatches(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)

	var request struct {
		Count int `json:"count"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		request.Count = 20 // Default
	}

	if request.Count <= 0 || request.Count > 100 {
		request.Count = 20
	}

	syncJob, err := s.matchService.SyncUserMatches(userObj.ID, request.Count)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start match synchronization: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"job_id": syncJob.ID,
		"status": syncJob.Status,
		"message": "Match synchronization started",
	})
}

// getSyncStatus récupère le statut d'un job de synchronisation
func (s *Server) getSyncStatus(c *gin.Context) {
	jobIDStr := c.Param("jobId")
	if jobIDStr == "" {
		c.JSON(400, gin.H{"error": "Job ID is required"})
		return
	}

	// Simple query to get sync job status
	query := `
		SELECT id, user_id, job_type, status, started_at, completed_at,
			   matches_processed, matches_new, matches_updated, error_message
		FROM sync_jobs WHERE id = ?
	`

	var job models.SyncJob
	var startedAt, completedAt sql.NullTime
	var matchesProcessed, matchesNew, matchesUpdated sql.NullInt64
	var errorMessage sql.NullString

	err := s.db.QueryRow(query, jobIDStr).Scan(
		&job.ID, &job.UserID, &job.JobType, &job.Status,
		&startedAt, &completedAt, &matchesProcessed,
		&matchesNew, &matchesUpdated, &errorMessage,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Sync job not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	if startedAt.Valid {
		job.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	if matchesProcessed.Valid {
		job.MatchesProcessed = int(matchesProcessed.Int64)
	}
	if matchesNew.Valid {
		job.MatchesNew = int(matchesNew.Int64)
	}
	if matchesUpdated.Valid {
		job.MatchesUpdated = int(matchesUpdated.Int64)
	}
	if errorMessage.Valid {
		job.ErrorMessage = &errorMessage.String
	}

	c.JSON(200, job)
}

// getUserMatches récupère les matches d'un utilisateur
func (s *Server) getUserMatches(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 20
	offset := 0

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	matches, err := s.matchService.GetUserMatches(userObj.ID, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get matches: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"matches": matches,
		"count":   len(matches),
		"limit":   limit,
		"offset":  offset,
	})
}

// getDashboard retourne les données du dashboard
func (s *Server) getDashboard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userObj := user.(*models.User)

	// Récupérer les statistiques de l'utilisateur
	stats, err := s.matchService.GetUserStats(userObj.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user stats: " + err.Error()})
		return
	}

	c.JSON(200, stats)
}

// logout déconnecte l'utilisateur en supprimant sa session
func (s *Server) logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err == nil {
		// Supprimer la session
		s.sessionMutex.Lock()
		delete(s.sessions, sessionID)
		s.sessionMutex.Unlock()
	}
	
	// Supprimer le cookie
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	
	c.JSON(200, gin.H{
		"message": "Logged out successfully",
	})
}

// getSupportedRegions retourne les régions supportées
func (s *Server) getSupportedRegions(c *gin.Context) {
	regions := []gin.H{
		{"code": "euw1", "name": "Europe West", "id": "EUW1"},
		{"code": "na1", "name": "North America", "id": "NA1"},
		{"code": "eun1", "name": "Europe Nordic & East", "id": "EUN1"},
		{"code": "kr", "name": "Korea", "id": "KR"},
		{"code": "jp1", "name": "Japan", "id": "JP1"},
		{"code": "br1", "name": "Brazil", "id": "BR1"},
		{"code": "la1", "name": "Latin America North", "id": "LA1"},
		{"code": "la2", "name": "Latin America South", "id": "LA2"},
		{"code": "oc1", "name": "Oceania", "id": "OC1"},
		{"code": "tr1", "name": "Turkey", "id": "TR1"},
		{"code": "ru", "name": "Russia", "id": "RU"},
		{"code": "ph2", "name": "Philippines", "id": "PH2"},
		{"code": "sg2", "name": "Singapore", "id": "SG2"},
		{"code": "th2", "name": "Thailand", "id": "TH2"},
		{"code": "tw2", "name": "Taiwan", "id": "TW2"},
		{"code": "vn2", "name": "Vietnam", "id": "VN2"},
	}

	c.JSON(200, gin.H{
		"regions": regions,
		"default": "euw1",
	})
}

// ValidateAuthRequest définit la structure de requête pour la validation d'authentification
type ValidateAuthRequest struct {
	RiotID  string `json:"riot_id" binding:"required"`
	RiotTag string `json:"riot_tag" binding:"required"`
	Region  string `json:"region" binding:"required"`
}

// validateAuth valide les informations d'authentification/compte avec l'API Riot
func (s *Server) validateAuth(c *gin.Context) {
	var req ValidateAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"valid": false,
			"error_message": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validation des champs requis
	if req.RiotID == "" || req.RiotTag == "" || req.Region == "" {
		c.JSON(400, gin.H{
			"valid": false,
			"error_message": "Riot ID, tag, and region are required",
		})
		return
	}

	// Validation de la région
	validRegions := []string{"euw1", "na1", "eun1", "kr", "jp1", "br1", "la1", "la2", "oc1", "tr1", "ru", "ph2", "sg2", "th2", "tw2", "vn2"}
	isValidRegion := false
	for _, region := range validRegions {
		if req.Region == region {
			isValidRegion = true
			break
		}
	}
	if !isValidRegion {
		c.JSON(400, gin.H{
			"valid": false,
			"error_message": "Invalid region: " + req.Region,
		})
		return
	}

	// Vérifier si l'API Riot est disponible
	if s.riotAPIService == nil {
		c.JSON(500, gin.H{
			"valid": false,
			"error_message": "Riot API service not available",
		})
		return
	}

	// Valider le compte avec l'API Riot
	account, err := s.riotAPIService.GetAccountByRiotID(req.RiotID, req.RiotTag)
	if err != nil {
		c.JSON(404, gin.H{
			"valid": false,
			"error_message": "Account not found or Riot API error: " + err.Error(),
		})
		return
	}

	// Récupérer les informations d'invocateur
	summoner, err := s.riotAPIService.GetSummonerByPUUID(account.PUUID, req.Region)
	if err != nil {
		c.JSON(404, gin.H{
			"valid": false,
			"error_message": "Summoner not found in region " + req.Region + ": " + err.Error(),
		})
		return
	}

	// Créer le service utilisateur
	userService := services.NewUserService(s.db)

	// Vérifier si l'utilisateur existe déjà
	existingUser, err := userService.GetUserByRiotID(req.RiotID, req.RiotTag)
	if err != nil {
		c.JSON(500, gin.H{
			"valid": false,
			"error_message": "Database error: " + err.Error(),
		})
		return
	}

	var user *models.User
	if existingUser != nil {
		// Utilisateur existant - mettre à jour ses informations
		err = userService.UpdateUserSummonerInfo(
			existingUser.ID,
			summoner.ID,
			summoner.AccountID,
			summoner.ProfileIconID,
			summoner.SummonerLevel,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"valid": false,
				"error_message": "Failed to update user info: " + err.Error(),
			})
			return
		}
		
		// Récupérer l'utilisateur mis à jour
		user, err = userService.GetUserByRiotID(req.RiotID, req.RiotTag)
		if err != nil {
			c.JSON(500, gin.H{
				"valid": false,
				"error_message": "Failed to retrieve updated user: " + err.Error(),
			})
			return
		}
	} else {
		// Nouvel utilisateur - le créer
		user, err = userService.CreateUser(
			req.RiotID,
			req.RiotTag,
			account.PUUID,
			req.Region,
			summoner.ProfileIconID,
			summoner.SummonerLevel,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"valid": false,
				"error_message": "Failed to create user: " + err.Error(),
			})
			return
		}
		
		// Mettre à jour avec les IDs d'invocateur
		err = userService.UpdateUserSummonerInfo(
			user.ID,
			summoner.ID,
			summoner.AccountID,
			summoner.ProfileIconID,
			summoner.SummonerLevel,
		)
		if err != nil {
			log.Printf("Warning: Failed to update summoner info for new user: %v", err)
		}
	}

	// Créer une session pour l'utilisateur
	sessionID, err := s.createUserSession(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"valid": false,
			"error_message": "Failed to create session: " + err.Error(),
		})
		return
	}

	// Définir le cookie de session
	c.SetCookie("session_id", sessionID, 86400*7, "/", "", false, true) // 7 jours

	c.JSON(200, gin.H{
		"valid": true,
		"user": user,
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// initDatabase initialise la base de données SQLite
func initDatabase() (*sql.DB, error) {
	// Créer le dossier data s'il n'existe pas
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	
	// Ouvrir la base de données SQLite
	db, err := sql.Open("sqlite", "./data/herald.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Tester la connexion
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	// Exécuter les migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	
	log.Println("[DATABASE] SQLite database initialized successfully")
	return db, nil
}

// runMigrations exécute les migrations de base de données
func runMigrations(db *sql.DB) error {
	// Créer la table de suivi des migrations
	migrationTrackingSQL := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version TEXT UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	
	if _, err := db.Exec(migrationTrackingSQL); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	// Liste des migrations à appliquer
	migrations := []struct {
		version string
		file    string
	}{
		{"001", "./internal/db/migrations/001_users_matches.sql"},
	}
	
	for _, migration := range migrations {
		// Vérifier si la migration a déjà été appliquée
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", migration.version).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.version, err)
		}
		
		if count > 0 {
			log.Printf("[MIGRATION] Skipping migration %s (already applied)", migration.version)
			continue
		}
		
		// Lire le fichier de migration
		migrationSQL, err := os.ReadFile(migration.file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration.file, err)
		}
		
		log.Printf("[MIGRATION] Applying migration %s", migration.version)
		
		// Exécuter la migration dans une transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", migration.version, err)
		}
		
		// Exécuter le SQL de migration
		if _, err := tx.Exec(string(migrationSQL)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.version, err)
		}
		
		// Marquer la migration comme appliquée
		if _, err := tx.Exec("INSERT INTO migrations (version) VALUES (?)", migration.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.version, err)
		}
		
		// Valider la transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.version, err)
		}
		
		log.Printf("[MIGRATION] Successfully applied migration %s", migration.version)
	}
	
	return nil
}

// initGoogleOAuth initie le flow OAuth Google
func (s *Server) initGoogleOAuth(c *gin.Context) {
	if !s.googleOAuth.IsConfigured() {
		c.JSON(503, gin.H{
			"error":   "OAuth Google non configuré",
			"message": "Veuillez configurer GOOGLE_CLIENT_ID et GOOGLE_CLIENT_SECRET",
			"mock_url": "https://accounts.google.com/oauth/authorize?client_id=mock&redirect_uri=https://herald.lol/auth/google/callback&response_type=code&scope=email+profile",
		})
		return
	}

	state := s.googleOAuth.GenerateState()
	authURL := s.googleOAuth.GetAuthURL(state)

	// En production, stockez le state de manière sécurisée (Redis, base de données, etc.)
	// Pour cette démo, on fait confiance au validateur côté callback

	c.JSON(200, gin.H{
		"auth_url": authURL,
		"state":    state,
		"message":  "Redirection vers Google OAuth",
	})
}

// googleOAuthCallback gère le callback OAuth Google
func (s *Server) googleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// Gérer les erreurs OAuth
	if errorParam != "" {
		log.Printf("OAuth error: %s", errorParam)
		c.Redirect(302, "/?oauth_error="+errorParam)
		return
	}

	if code == "" {
		log.Printf("Missing OAuth code")
		c.Redirect(302, "/?oauth_error=missing_code")
		return
	}

	// Valider le state
	if !s.googleOAuth.ValidateState(state) {
		log.Printf("Invalid OAuth state: %s", state)
		c.Redirect(302, "/?oauth_error=invalid_state")
		return
	}

	if !s.googleOAuth.IsConfigured() {
		// Mode mock pour développement
		c.Redirect(302, "/?oauth_success=mock&user=mock_user&email=user@example.com")
		return
	}

	// Échanger le code contre un token
	ctx := c.Request.Context()
	tokenResp, err := s.googleOAuth.ExchangeCodeForToken(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		c.Redirect(302, "/?oauth_error=token_exchange_failed")
		return
	}

	// Récupérer les informations utilisateur
	userInfo, err := s.googleOAuth.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.Redirect(302, "/?oauth_error=userinfo_failed")
		return
	}

	// En production, vous stockeriez ces informations en session/base de données
	log.Printf("User authenticated: %s (%s)", userInfo.Name, userInfo.Email)

	// Rediriger vers le frontend avec les informations utilisateur
	redirectURL := fmt.Sprintf("/?oauth_success=true&user=%s&email=%s&picture=%s",
		url.QueryEscape(userInfo.Name),
		url.QueryEscape(userInfo.Email),
		url.QueryEscape(userInfo.Picture))

	c.Redirect(302, redirectURL)
}

// Analytics endpoint handlers

// getPeriodAnalytics récupère les analytics pour une période donnée
func (s *Server) getPeriodAnalytics(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userObj := user.(*models.User)
	
	period := c.Param("period")
	if period == "" {
		period = "week" // Valeur par défaut
	}
	
	// Vérifier si la période est valide
	validPeriods := map[string]bool{
		"today": true, "week": true, "month": true, "season": true,
	}
	if !validPeriods[period] {
		c.JSON(400, gin.H{"error": "Invalid period. Must be one of: today, week, month, season"})
		return
	}
	
	stats, err := s.analyticsService.GetPeriodStats(userObj.ID, period)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get period analytics: " + err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"success": true,
		"data":    stats,
	})
}

// getRecommendations récupère les recommandations IA
func (s *Server) getRecommendations(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userObj := user.(*models.User)
	
	recommendations, err := s.analyticsService.GetRecommendations(userObj.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get recommendations: " + err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"success": true,
		"data":    recommendations,
	})
}

// getPerformanceTrends récupère les tendances de performance
func (s *Server) getPerformanceTrends(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userObj := user.(*models.User)
	
	trends, err := s.analyticsService.GetPerformanceTrends(userObj.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get performance trends: " + err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"success": true,
		"data":    trends,
	})
}

// refreshAnalytics force le rafraîchissement des analytics
func (s *Server) refreshAnalytics(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userObj := user.(*models.User)
	
	var request struct {
		Period string `json:"period"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		request.Period = "week" // Valeur par défaut
	}
	
	// Mettre à jour les statistiques
	err := s.analyticsService.UpdateChampionStats(userObj.ID, request.Period)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to refresh analytics: " + err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"success": true,
		"message": "Analytics refreshed successfully",
	})
}

// Settings endpoint handlers

// getUserSettings récupère les paramètres de l'utilisateur
func (s *Server) getUserSettings(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userObj := user.(*models.User)
	
	// Récupérer les settings de l'utilisateur depuis la base de données
	settings, err := s.getUserSettingsFromDB(userObj.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user settings: " + err.Error()})
		return
	}
	
	c.JSON(200, settings)
}

// updateUserSettings met à jour les paramètres de l'utilisateur
func (s *Server) updateUserSettings(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userObj := user.(*models.User)
	
	var settings UserSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(400, gin.H{"error": "Invalid settings data: " + err.Error()})
		return
	}
	
	// Valider les settings
	if err := s.validateUserSettings(&settings); err != nil {
		c.JSON(400, gin.H{"error": "Invalid settings: " + err.Error()})
		return
	}
	
	// Sauvegarder en base de données
	err := s.saveUserSettings(userObj.ID, &settings)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save settings: " + err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"success": true,
		"message": "Settings updated successfully",
	})
}

// UserSettings structure pour les paramètres utilisateur
type UserSettings struct {
	IncludeTimeline    bool `json:"include_timeline"`
	IncludeAllData     bool `json:"include_all_data"`
	LightMode          bool `json:"light_mode"`
	AutoSyncEnabled    bool `json:"auto_sync_enabled"`
	SyncFrequencyHours int  `json:"sync_frequency_hours"`
}

// getUserSettingsFromDB récupère les settings depuis la base de données
func (s *Server) getUserSettingsFromDB(userID int) (*UserSettings, error) {
	query := `
		SELECT 
			include_timeline, 
			include_all_data, 
			light_mode, 
			auto_sync_enabled, 
			sync_frequency_hours
		FROM user_settings 
		WHERE user_id = ?
	`
	
	var settings UserSettings
	err := s.db.QueryRow(query, userID).Scan(
		&settings.IncludeTimeline,
		&settings.IncludeAllData,
		&settings.LightMode,
		&settings.AutoSyncEnabled,
		&settings.SyncFrequencyHours,
	)
	
	if err == sql.ErrNoRows {
		// Créer des settings par défaut si aucun n'existe
		defaultSettings := &UserSettings{
			IncludeTimeline:    true,
			IncludeAllData:     true,
			LightMode:          true,
			AutoSyncEnabled:    true,
			SyncFrequencyHours: 24,
		}
		
		// Sauvegarder les settings par défaut
		if saveErr := s.saveUserSettings(userID, defaultSettings); saveErr != nil {
			return nil, fmt.Errorf("failed to create default settings: %w", saveErr)
		}
		
		return defaultSettings, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to query user settings: %w", err)
	}
	
	return &settings, nil
}

// saveUserSettings sauvegarde les settings en base de données
func (s *Server) saveUserSettings(userID int, settings *UserSettings) error {
	query := `
		INSERT INTO user_settings (
			user_id, 
			include_timeline, 
			include_all_data, 
			light_mode, 
			auto_sync_enabled, 
			sync_frequency_hours,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(user_id) DO UPDATE SET
			include_timeline = excluded.include_timeline,
			include_all_data = excluded.include_all_data,
			light_mode = excluded.light_mode,
			auto_sync_enabled = excluded.auto_sync_enabled,
			sync_frequency_hours = excluded.sync_frequency_hours,
			updated_at = datetime('now')
	`
	
	_, err := s.db.Exec(query,
		userID,
		settings.IncludeTimeline,
		settings.IncludeAllData,
		settings.LightMode,
		settings.AutoSyncEnabled,
		settings.SyncFrequencyHours,
	)
	
	if err != nil {
		return fmt.Errorf("failed to save user settings: %w", err)
	}
	
	return nil
}

// validateUserSettings valide les paramètres utilisateur
func (s *Server) validateUserSettings(settings *UserSettings) error {
	// Valider la fréquence de sync
	validFrequencies := map[int]bool{
		1: true, 6: true, 12: true, 24: true, 72: true, 168: true,
	}
	
	if !validFrequencies[settings.SyncFrequencyHours] {
		return fmt.Errorf("invalid sync frequency: %d hours", settings.SyncFrequencyHours)
	}
	
	return nil
}
