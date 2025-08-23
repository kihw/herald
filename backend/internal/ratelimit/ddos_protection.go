package ratelimit

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Herald.lol Gaming Analytics - DDoS Protection
// Advanced DDoS protection with pattern detection and geographical filtering

// DDoSProtector provides DDoS protection functionality
type DDoSProtector struct {
	redis          *redis.Client
	config         *DDoSProtectionConfig
	patternMatcher *SuspiciousPatternMatcher
	geoFilter      *GeographicalFilter
}

// SuspiciousActivity represents detected suspicious activity
type SuspiciousActivity struct {
	IP               string                 `json:"ip"`
	UserAgent        string                 `json:"user_agent"`
	RequestCount     int                    `json:"request_count"`
	TimeWindow       time.Duration          `json:"time_window"`
	SuspiciousScore  int                    `json:"suspicious_score"`
	DetectedPatterns []string               `json:"detected_patterns"`
	FirstSeen        time.Time              `json:"first_seen"`
	LastSeen         time.Time              `json:"last_seen"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// DDoSProtectionResult contains protection decision
type DDoSProtectionResult struct {
	Action           string              `json:"action"`           // allow, block, challenge
	Reason           string              `json:"reason"`
	SuspiciousScore  int                 `json:"suspicious_score"`
	BlockDuration    time.Duration       `json:"block_duration,omitempty"`
	Activity         *SuspiciousActivity `json:"activity,omitempty"`
	RequiresChallenge bool               `json:"requires_challenge"`
}

// NewDDoSProtector creates new DDoS protector
func NewDDoSProtector(redis *redis.Client, config *DDoSProtectionConfig) *DDoSProtector {
	return &DDoSProtector{
		redis:          redis,
		config:         config,
		patternMatcher: NewSuspiciousPatternMatcher(config.SuspiciousPatterns),
		geoFilter:      NewGeographicalFilter(config.GeoBlocking, config.BlockedCountries),
	}
}

// CheckForDDoS performs comprehensive DDoS detection
func (d *DDoSProtector) CheckForDDoS(ctx context.Context, request *DDoSRequest) (*DDoSProtectionResult, error) {
	if !d.config.Enabled {
		return &DDoSProtectionResult{Action: "allow"}, nil
	}
	
	// 1. Check if IP is already blocked
	if blocked, reason, err := d.isIPBlocked(ctx, request.ClientIP); err != nil {
		return nil, err
	} else if blocked {
		return &DDoSProtectionResult{
			Action:        "block",
			Reason:        fmt.Sprintf("IP blocked: %s", reason),
			BlockDuration: d.config.BlockDuration,
		}, nil
	}
	
	// 2. Check geographical restrictions
	if d.config.GeoBlocking {
		if blocked, err := d.geoFilter.IsBlocked(request.ClientIP, request.Country); err != nil {
			return nil, err
		} else if blocked {
			return d.blockIP(ctx, request.ClientIP, "geographical_restriction", 24*time.Hour)
		}
	}
	
	// 3. Analyze request patterns
	activity, err := d.analyzeActivity(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze activity: %w", err)
	}
	
	// 4. Calculate suspicious score
	suspiciousScore := d.calculateSuspiciousScore(request, activity)
	
	// 5. Make protection decision
	if suspiciousScore > 80 {
		return d.blockIP(ctx, request.ClientIP, "high_suspicious_score", d.config.BlockDuration)
	} else if suspiciousScore > 50 {
		return &DDoSProtectionResult{
			Action:           "challenge",
			Reason:           "moderate_suspicious_score",
			SuspiciousScore:  suspiciousScore,
			Activity:         activity,
			RequiresChallenge: true,
		}, nil
	}
	
	// 6. Check request rate threshold
	if activity.RequestCount > d.config.RequestThreshold {
		return d.blockIP(ctx, request.ClientIP, "request_threshold_exceeded", d.config.BlockDuration)
	}
	
	return &DDoSProtectionResult{
		Action:          "allow",
		SuspiciousScore: suspiciousScore,
		Activity:        activity,
	}, nil
}

// analyzeActivity analyzes request activity patterns
func (d *DDoSProtector) analyzeActivity(ctx context.Context, request *DDoSRequest) (*SuspiciousActivity, error) {
	now := time.Now()
	windowKey := fmt.Sprintf("activity:%s:%d", request.ClientIP, now.Unix()/int64(d.config.WindowDuration.Seconds()))
	
	// Get or create activity record
	activity, err := d.getActivity(ctx, request.ClientIP, windowKey)
	if err != nil {
		return nil, err
	}
	
	// Update activity
	activity.RequestCount++
	activity.LastSeen = now
	if activity.FirstSeen.IsZero() {
		activity.FirstSeen = now
	}
	
	// Store updated activity
	if err := d.storeActivity(ctx, windowKey, activity); err != nil {
		return nil, err
	}
	
	// Detect suspicious patterns
	patterns := d.patternMatcher.AnalyzeRequest(request)
	activity.DetectedPatterns = append(activity.DetectedPatterns, patterns...)
	
	return activity, nil
}

// calculateSuspiciousScore calculates overall suspicious score
func (d *DDoSProtector) calculateSuspiciousScore(request *DDoSRequest, activity *SuspiciousActivity) int {
	score := 0
	
	// High request rate
	if activity.RequestCount > d.config.RequestThreshold {
		score += 40
	} else if activity.RequestCount > d.config.RequestThreshold/2 {
		score += 20
	}
	
	// Suspicious user agent patterns
	if d.patternMatcher.IsSuspiciousUserAgent(request.UserAgent) {
		score += 25
	}
	
	// Suspicious patterns in URLs or headers
	if len(activity.DetectedPatterns) > 0 {
		score += len(activity.DetectedPatterns) * 10
	}
	
	// No referrer (possible bot)
	if request.Referrer == "" && !strings.Contains(request.Path, "/api/") {
		score += 10
	}
	
	// Unusual request timing (too fast)
	if !activity.FirstSeen.IsZero() && !activity.LastSeen.IsZero() {
		duration := activity.LastSeen.Sub(activity.FirstSeen)
		if duration < 10*time.Second && activity.RequestCount > 50 {
			score += 30
		}
	}
	
	// Known attack patterns in path
	if d.containsAttackPatterns(request.Path) {
		score += 50
	}
	
	// High number of different endpoints accessed
	if d.isEndpointScan(request.ClientIP, request.Path) {
		score += 20
	}
	
	// Cap at 100
	if score > 100 {
		score = 100
	}
	
	return score
}

// blockIP blocks an IP address
func (d *DDoSProtector) blockIP(ctx context.Context, ip, reason string, duration time.Duration) (*DDoSProtectionResult, error) {
	blockKey := fmt.Sprintf("blocked:ip:%s", ip)
	
	// Store block reason and metadata
	blockInfo := map[string]interface{}{
		"reason":     reason,
		"blocked_at": time.Now().Unix(),
		"duration":   int(duration.Seconds()),
		"expires_at": time.Now().Add(duration).Unix(),
	}
	
	// Set block
	err := d.redis.Set(ctx, blockKey, reason, duration).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to block IP: %w", err)
	}
	
	// Log block event
	logKey := fmt.Sprintf("ddos_blocks:%d", time.Now().Unix())
	d.redis.HSet(ctx, logKey, ip, blockInfo)
	d.redis.Expire(ctx, logKey, 24*time.Hour)
	
	// Update global DDoS stats
	d.updateDDoSStats(ctx, "blocks", 1)
	
	return &DDoSProtectionResult{
		Action:        "block",
		Reason:        reason,
		BlockDuration: duration,
	}, nil
}

// isIPBlocked checks if IP is currently blocked
func (d *DDoSProtector) isIPBlocked(ctx context.Context, ip string) (bool, string, error) {
	blockKey := fmt.Sprintf("blocked:ip:%s", ip)
	reason, err := d.redis.Get(ctx, blockKey).Result()
	if err == redis.Nil {
		return false, "", nil
	}
	if err != nil {
		return false, "", fmt.Errorf("failed to check IP block status: %w", err)
	}
	return true, reason, nil
}

// getActivity retrieves or creates activity record
func (d *DDoSProtector) getActivity(ctx context.Context, ip, windowKey string) (*SuspiciousActivity, error) {
	// Try to get existing activity
	activityData, err := d.redis.HGetAll(ctx, windowKey).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get activity data: %w", err)
	}
	
	if len(activityData) == 0 {
		// Create new activity record
		return &SuspiciousActivity{
			IP:               ip,
			RequestCount:     0,
			DetectedPatterns: []string{},
			Metadata:         make(map[string]interface{}),
		}, nil
	}
	
	// Parse existing activity
	activity := &SuspiciousActivity{
		IP:               ip,
		DetectedPatterns: []string{},
		Metadata:         make(map[string]interface{}),
	}
	
	if val, ok := activityData["request_count"]; ok {
		if count, err := redis.ParseInt(val); err == nil {
			activity.RequestCount = int(count)
		}
	}
	
	if val, ok := activityData["user_agent"]; ok {
		activity.UserAgent = val
	}
	
	return activity, nil
}

// storeActivity stores activity record
func (d *DDoSProtector) storeActivity(ctx context.Context, windowKey string, activity *SuspiciousActivity) error {
	activityData := map[string]interface{}{
		"request_count": activity.RequestCount,
		"user_agent":    activity.UserAgent,
		"first_seen":    activity.FirstSeen.Unix(),
		"last_seen":     activity.LastSeen.Unix(),
	}
	
	err := d.redis.HSet(ctx, windowKey, activityData).Err()
	if err != nil {
		return fmt.Errorf("failed to store activity: %w", err)
	}
	
	// Set expiration
	d.redis.Expire(ctx, windowKey, d.config.WindowDuration*2)
	
	return nil
}

// containsAttackPatterns checks for common attack patterns
func (d *DDoSProtector) containsAttackPatterns(path string) bool {
	attackPatterns := []string{
		"../", "..\\", ".env", "wp-admin", "wp-login", "phpmyadmin",
		"admin/", "login.php", "config.php", "shell", "cmd",
		"<script", "javascript:", "onload=", "eval(",
		"union select", "drop table", "insert into",
		"etc/passwd", "proc/self", "/dev/null",
	}
	
	lowerPath := strings.ToLower(path)
	for _, pattern := range attackPatterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}
	
	return false
}

// isEndpointScan checks if IP is scanning for endpoints
func (d *DDoSProtector) isEndpointScan(ip, path string) bool {
	ctx := context.Background()
	scanKey := fmt.Sprintf("scan:%s:endpoints", ip)
	
	// Add current path to set
	d.redis.SAdd(ctx, scanKey, path)
	d.redis.Expire(ctx, scanKey, 5*time.Minute)
	
	// Count unique endpoints accessed
	count, err := d.redis.SCard(ctx, scanKey).Result()
	if err != nil {
		return false
	}
	
	// If accessing many different endpoints quickly, it's likely a scan
	return count > 20
}

// updateDDoSStats updates global DDoS statistics
func (d *DDoSProtector) updateDDoSStats(ctx context.Context, metric string, value int64) {
	now := time.Now()
	statsKey := fmt.Sprintf("ddos_stats:%s:%d", metric, now.Unix()/3600) // Hourly stats
	
	d.redis.IncrBy(ctx, statsKey, value)
	d.redis.Expire(ctx, statsKey, 24*time.Hour)
}

// DDoSRequest contains request information for DDoS analysis
type DDoSRequest struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Referrer  string `json:"referrer"`
	Country   string `json:"country,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// DDoSProtectionMiddleware creates Gin middleware for DDoS protection
func (d *DDoSProtector) DDoSProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &DDoSRequest{
			ClientIP:  c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			Path:      c.Request.URL.Path,
			Method:    c.Request.Method,
			Referrer:  c.GetHeader("Referer"),
			Country:   c.GetHeader("CF-IPCountry"), // Cloudflare country header
		}
		
		result, err := d.CheckForDDoS(c.Request.Context(), request)
		if err != nil {
			c.JSON(500, gin.H{
				"error":           "DDoS protection check failed",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		switch result.Action {
		case "block":
			c.Header("X-DDoS-Protection", "blocked")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":           "Request blocked by DDoS protection",
				"reason":          result.Reason,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
			
		case "challenge":
			c.Header("X-DDoS-Protection", "challenge")
			// In a real implementation, this would redirect to a CAPTCHA or similar
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":             "Challenge required",
				"reason":            result.Reason,
				"suspicious_score":  result.SuspiciousScore,
				"challenge_required": true,
				"gaming_platform":   "herald-lol",
			})
			c.Abort()
			return
			
		case "allow":
			c.Header("X-DDoS-Protection", "passed")
			if result.SuspiciousScore > 0 {
				c.Header("X-Suspicious-Score", fmt.Sprintf("%d", result.SuspiciousScore))
			}
		}
		
		c.Next()
	}
}

// GetDDoSStats retrieves DDoS protection statistics
func (d *DDoSProtector) GetDDoSStats(ctx context.Context, hours int) (*DDoSStats, error) {
	now := time.Now()
	stats := &DDoSStats{
		TimeRange: fmt.Sprintf("%d hours", hours),
	}
	
	// Collect hourly stats
	for i := 0; i < hours; i++ {
		hour := now.Add(-time.Duration(i) * time.Hour)
		hourKey := hour.Unix() / 3600
		
		// Get blocks for this hour
		blocksKey := fmt.Sprintf("ddos_stats:blocks:%d", hourKey)
		blocks, err := d.redis.Get(ctx, blocksKey).Int64()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		
		stats.TotalBlocks += blocks
	}
	
	// Get currently blocked IPs count
	blockedPattern := "blocked:ip:*"
	blockedIPs, err := d.redis.Keys(ctx, blockedPattern).Result()
	if err == nil {
		stats.CurrentlyBlocked = len(blockedIPs)
	}
	
	return stats, nil
}

// DDoSStats contains DDoS protection statistics
type DDoSStats struct {
	TimeRange        string `json:"time_range"`
	TotalBlocks      int64  `json:"total_blocks"`
	CurrentlyBlocked int    `json:"currently_blocked"`
}

// SuspiciousPatternMatcher analyzes requests for suspicious patterns
type SuspiciousPatternMatcher struct {
	patterns []string
}

// NewSuspiciousPatternMatcher creates new pattern matcher
func NewSuspiciousPatternMatcher(patterns []string) *SuspiciousPatternMatcher {
	return &SuspiciousPatternMatcher{
		patterns: patterns,
	}
}

// AnalyzeRequest analyzes request for suspicious patterns
func (s *SuspiciousPatternMatcher) AnalyzeRequest(request *DDoSRequest) []string {
	var detected []string
	
	// Check user agent
	if s.IsSuspiciousUserAgent(request.UserAgent) {
		detected = append(detected, "suspicious_user_agent")
	}
	
	// Check for automated tools
	lowerUA := strings.ToLower(request.UserAgent)
	automatedTools := []string{"bot", "crawler", "spider", "scraper", "curl", "wget", "python", "requests"}
	for _, tool := range automatedTools {
		if strings.Contains(lowerUA, tool) {
			detected = append(detected, fmt.Sprintf("automated_tool_%s", tool))
		}
	}
	
	return detected
}

// IsSuspiciousUserAgent checks if user agent is suspicious
func (s *SuspiciousPatternMatcher) IsSuspiciousUserAgent(userAgent string) bool {
	if userAgent == "" {
		return true
	}
	
	suspiciousPatterns := []string{
		"bot", "crawler", "spider", "scraper", "scanner",
		"hack", "attack", "exploit", "injection",
		"masscan", "nmap", "nikto", "sqlmap",
	}
	
	lowerUA := strings.ToLower(userAgent)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerUA, pattern) {
			return true
		}
	}
	
	return false
}

// GeographicalFilter handles geographical filtering
type GeographicalFilter struct {
	enabled          bool
	blockedCountries []string
}

// NewGeographicalFilter creates new geographical filter
func NewGeographicalFilter(enabled bool, blockedCountries []string) *GeographicalFilter {
	return &GeographicalFilter{
		enabled:          enabled,
		blockedCountries: blockedCountries,
	}
}

// IsBlocked checks if IP/country should be blocked
func (g *GeographicalFilter) IsBlocked(ip, country string) (bool, error) {
	if !g.enabled {
		return false, nil
	}
	
	if country == "" {
		// Try to determine country from IP (simplified)
		country = g.getCountryFromIP(ip)
	}
	
	for _, blocked := range g.blockedCountries {
		if country == blocked {
			return true, nil
		}
	}
	
	return false, nil
}

// getCountryFromIP determines country from IP (simplified implementation)
func (g *GeographicalFilter) getCountryFromIP(ip string) string {
	// In a real implementation, this would use a GeoIP database
	// For now, return empty string
	return ""
}

// UnblockIP removes IP from blocked list
func (d *DDoSProtector) UnblockIP(ctx context.Context, ip string) error {
	blockKey := fmt.Sprintf("blocked:ip:%s", ip)
	return d.redis.Del(ctx, blockKey).Err()
}

// GetBlockedIPs returns list of currently blocked IPs
func (d *DDoSProtector) GetBlockedIPs(ctx context.Context) ([]string, error) {
	pattern := "blocked:ip:*"
	keys, err := d.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	
	var blockedIPs []string
	for _, key := range keys {
		// Extract IP from key: blocked:ip:192.168.1.1
		parts := strings.Split(key, ":")
		if len(parts) >= 3 {
			ip := strings.Join(parts[2:], ":") // Handle IPv6
			blockedIPs = append(blockedIPs, ip)
		}
	}
	
	return blockedIPs, nil
}