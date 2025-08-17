package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

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
	Username  string `json:"username" binding:"required"`
	TagLine   string `json:"tagline" binding:"required"`
	GameCount int    `json:"gameCount"`
	APIKey    string `json:"apiKey"`
}

type Server struct {
	jobs            map[string]*ExportJob
	mutex           sync.RWMutex
	exportService   *services.ExportService
	templateService *services.TemplateService
	riotAPIService  *services.RiotAPIService
}

func NewServer() *Server {
	// Récupérer la clé API Riot depuis les variables d'environnement
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	var riotService *services.RiotAPIService
	if riotAPIKey != "" {
		riotService = services.NewRiotAPIService(riotAPIKey)
	}

	return &Server{
		jobs:            make(map[string]*ExportJob),
		exportService:   services.NewExportService("./exports"),
		templateService: services.NewTemplateService(),
		riotAPIService:  riotService,
	}
}

func generateJobID() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *Server) startExport(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Force count to 1000
	req.GameCount = 1000

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

	// Prepare command
	args := []string{
		"lol_match_exporter.py",
		"--username", req.Username,
		"--tagline", req.TagLine,
		"--count", fmt.Sprintf("%d", req.GameCount),
		"--light-mode",
	}

	if req.APIKey != "" {
		args = append(args, "--api-key", req.APIKey)
	}

	cmd := exec.CommandContext(job.ctx, "python3", args...)
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	job.mutex.Lock()
	job.cmd = cmd
	job.mutex.Unlock()

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s.setJobError(job, fmt.Sprintf("Failed to create stdout pipe: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		s.setJobError(job, fmt.Sprintf("Failed to create stderr pipe: %v", err))
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		s.setJobError(job, fmt.Sprintf("Failed to start export process: %v", err))
		return
	}

	log.Printf("[SERVER] Subprocess created with PID %d", cmd.Process.Pid)

	// Read output
	go s.readOutput(job, stdout, "STDOUT")
	go s.readOutput(job, stderr, "STDERR")

	// Wait for completion
	err = cmd.Wait()

	job.mutex.Lock()
	defer job.mutex.Unlock()

	if err != nil {
		if job.ctx.Err() == context.Canceled {
			job.Status = "cancelled"
		} else {
			job.Status = "failed"
			job.Error = err.Error()
		}
	} else {
		job.Status = "completed"
		job.Progress = 100

		// Look for generated zip file
		zipPath := s.findGeneratedZip(req.Username, req.TagLine)
		if zipPath != "" {
			job.ZipPath = zipPath
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

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Zip file not found"})
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

		// Auth endpoints (mock for frontend compatibility)
		api.GET("/auth/session", server.getSession)
		api.GET("/auth/regions", server.getSupportedRegions)
		api.POST("/auth/validate", server.validateAuth)

		api.POST("/export", server.startExport)
		api.GET("/jobs/:job_id", server.getJobStatus)
		api.GET("/jobs/:job_id/logs", server.streamLogs)
		api.GET("/jobs/:job_id/download", server.downloadZip)

		// Advanced export endpoints
		api.POST("/export/advanced", server.startAdvancedExport)
		api.GET("/export/formats", server.getSupportedFormats)
		api.POST("/export/validate", server.validateExportOptions)
		api.GET("/export/history", server.getExportHistory)

		// Template endpoints
		api.GET("/templates", server.getAllTemplates)
		api.GET("/templates/:id", server.getTemplate)
		api.POST("/templates", server.createTemplate)
		api.PUT("/templates/:id", server.updateTemplate)
		api.DELETE("/templates/:id", server.deleteTemplate)
		api.POST("/export/from-template", server.exportFromTemplate)
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

// getSession retourne les informations de session (mock)
func (s *Server) getSession(c *gin.Context) {
	c.JSON(200, gin.H{
		"authenticated": false,
		"user":          nil,
		"message":       "No authentication required for this demo",
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
	Username string `json:"username" binding:"required"`
	TagLine  string `json:"tagline" binding:"required"`
	Region   string `json:"region"`
}

// validateAuth valide les informations d'authentification/compte
func (s *Server) validateAuth(c *gin.Context) {
	var req ValidateAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Simulation de validation (dans une vraie app, on vérifierait via l'API Riot)
	// Pour le moment, on accepte tous les comptes non vides
	if req.Username == "" || req.TagLine == "" {
		c.JSON(400, gin.H{
			"valid": false,
			"error": "Username and tagline are required",
		})
		return
	}

	// Validation de la région si fournie
	if req.Region != "" {
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
				"error": "Invalid region specified",
			})
			return
		}
	}

	// Si on a une vraie API Riot configurée, on peut faire une validation plus poussée
	if s.riotAPIService != nil {
		// Ici on pourrait vérifier si le compte existe vraiment
		// Pour le moment, on simule juste une réponse positive
	}

	c.JSON(200, gin.H{
		"valid":    true,
		"message":  "Account validation successful",
		"username": req.Username,
		"tagline":  req.TagLine,
		"region":   req.Region,
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
