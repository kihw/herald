package middleware

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Herald.lol Gaming Analytics - Advanced DDoS Protection
// Multi-layered DDoS protection system for gaming platform

// DDoSProtector provides comprehensive DDoS protection
type DDoSProtector struct {
	redisClient *redis.Client
	config      *GamingRateLimitConfig
	patterns    *AttackPatternDetector
	mitigation  *MitigationEngine
}

// AttackPatternDetector detects various attack patterns
type AttackPatternDetector struct {
	redisClient *redis.Client
}

// MitigationEngine handles attack mitigation strategies
type MitigationEngine struct {
	redisClient *redis.Client
	config      *GamingRateLimitConfig
}

// BlockInfo contains information about blocked clients
type BlockInfo struct {
	Reason       string     `json:"reason"`
	BlockedAt    time.Time  `json:"blocked_at"`
	BlockedUntil *time.Time `json:"blocked_until,omitempty"`
	AttackType   string     `json:"attack_type"`
	Severity     string     `json:"severity"`
	RequestCount int        `json:"request_count"`
	IPAddress    string     `json:"ip_address"`
}

// AttackSignature represents detected attack patterns
type AttackSignature struct {
	Type       string    `json:"type"`
	Confidence float64   `json:"confidence"`
	Evidence   []string  `json:"evidence"`
	Detected   time.Time `json:"detected"`
	Severity   string    `json:"severity"`
}

// NewDDoSProtector creates new DDoS protection system
func NewDDoSProtector(redisClient *redis.Client, config *GamingRateLimitConfig) *DDoSProtector {
	return &DDoSProtector{
		redisClient: redisClient,
		config:      config,
		patterns: &AttackPatternDetector{
			redisClient: redisClient,
		},
		mitigation: &MitigationEngine{
			redisClient: redisClient,
			config:      config,
		},
	}
}

// IsBlocked checks if a client is currently blocked
func (ddos *DDoSProtector) IsBlocked(ctx context.Context, clientID string) (bool, *BlockInfo) {
	key := fmt.Sprintf("herald:ddos:blocked:%s", clientID)

	result, err := ddos.redisClient.HGetAll(ctx, key).Result()
	if err != nil || len(result) == 0 {
		return false, nil
	}

	blockInfo := &BlockInfo{
		Reason:     result["reason"],
		AttackType: result["attack_type"],
		Severity:   result["severity"],
		IPAddress:  result["ip_address"],
	}

	// Parse timestamps
	if blockedAtStr, exists := result["blocked_at"]; exists {
		if timestamp, err := strconv.ParseInt(blockedAtStr, 10, 64); err == nil {
			blockInfo.BlockedAt = time.Unix(timestamp, 0)
		}
	}

	if blockedUntilStr, exists := result["blocked_until"]; exists && blockedUntilStr != "" {
		if timestamp, err := strconv.ParseInt(blockedUntilStr, 10, 64); err == nil {
			blockedUntil := time.Unix(timestamp, 0)
			blockInfo.BlockedUntil = &blockedUntil

			// Check if block has expired
			if time.Now().After(blockedUntil) {
				ddos.unblockClient(ctx, clientID)
				return false, nil
			}
		}
	}

	if requestCountStr, exists := result["request_count"]; exists {
		if count, err := strconv.Atoi(requestCountStr); err == nil {
			blockInfo.RequestCount = count
		}
	}

	return true, blockInfo
}

// RecordRequest records a request for DDoS analysis
func (ddos *DDoSProtector) RecordRequest(ctx context.Context, clientID, ipAddress string) {
	now := time.Now()

	// Record request in time-series for pattern analysis
	requestKey := fmt.Sprintf("herald:ddos:requests:%s", clientID)
	ddos.redisClient.ZAdd(ctx, requestKey, &redis.Z{
		Score:  float64(now.Unix()),
		Member: fmt.Sprintf("%d-%s", now.UnixNano(), ipAddress),
	})
	ddos.redisClient.Expire(ctx, requestKey, 10*time.Minute)

	// Record IP-based requests for IP blocking
	ipKey := fmt.Sprintf("herald:ddos:ip_requests:%s", ipAddress)
	ddos.redisClient.ZAdd(ctx, ipKey, &redis.Z{
		Score:  float64(now.Unix()),
		Member: fmt.Sprintf("%d-%s", now.UnixNano(), clientID),
	})
	ddos.redisClient.Expire(ctx, ipKey, 10*time.Minute)

	// Analyze patterns and trigger mitigation if needed
	go ddos.analyzeAndMitigate(ctx, clientID, ipAddress)
}

// analyzeAndMitigate performs real-time analysis and mitigation
func (ddos *DDoSProtector) analyzeAndMitigate(ctx context.Context, clientID, ipAddress string) {
	// Check for various attack patterns
	signatures := ddos.patterns.DetectAttackPatterns(ctx, clientID, ipAddress)

	for _, signature := range signatures {
		if signature.Confidence > 0.7 { // High confidence threshold
			// Trigger mitigation
			ddos.mitigation.TriggerMitigation(ctx, clientID, ipAddress, signature)

			// Log attack detection
			ddos.logAttackDetection(ctx, clientID, ipAddress, signature)
		}
	}
}

// DetectAttackPatterns detects various DDoS attack patterns
func (apd *AttackPatternDetector) DetectAttackPatterns(ctx context.Context, clientID, ipAddress string) []*AttackSignature {
	var signatures []*AttackSignature

	// 1. Volumetric attack detection
	if volumetricSig := apd.detectVolumetricAttack(ctx, clientID, ipAddress); volumetricSig != nil {
		signatures = append(signatures, volumetricSig)
	}

	// 2. Burst attack detection
	if burstSig := apd.detectBurstAttack(ctx, clientID, ipAddress); burstSig != nil {
		signatures = append(signatures, burstSig)
	}

	// 3. Slowloris attack detection
	if slowlorisSig := apd.detectSlowlorisAttack(ctx, clientID, ipAddress); slowlorisSig != nil {
		signatures = append(signatures, slowlorisSig)
	}

	// 4. Distributed attack detection
	if distributedSig := apd.detectDistributedAttack(ctx, ipAddress); distributedSig != nil {
		signatures = append(signatures, distributedSig)
	}

	// 5. Gaming-specific attack patterns
	if gamingSig := apd.detectGamingAttackPatterns(ctx, clientID, ipAddress); gamingSig != nil {
		signatures = append(signatures, gamingSig)
	}

	return signatures
}

// detectVolumetricAttack detects high-volume attacks
func (apd *AttackPatternDetector) detectVolumetricAttack(ctx context.Context, clientID, ipAddress string) *AttackSignature {
	now := time.Now()
	windowStart := now.Add(-time.Minute)

	// Check requests in the last minute
	requestKey := fmt.Sprintf("herald:ddos:requests:%s", clientID)
	count, err := apd.redisClient.ZCount(ctx, requestKey, strconv.FormatInt(windowStart.Unix(), 10), "+inf").Result()
	if err != nil {
		return nil
	}

	// Volumetric attack threshold
	threshold := int64(1000) // 1000 requests per minute

	if count > threshold {
		confidence := float64(count) / float64(threshold)
		if confidence > 3.0 {
			confidence = 1.0 // Cap at 100%
		} else {
			confidence = confidence / 3.0
		}

		return &AttackSignature{
			Type:       "volumetric",
			Confidence: confidence,
			Evidence: []string{
				fmt.Sprintf("High request volume: %d requests/minute", count),
				fmt.Sprintf("Threshold exceeded by %.1fx", float64(count)/float64(threshold)),
			},
			Detected: now,
			Severity: apd.getSeverity(confidence),
		}
	}

	return nil
}

// detectBurstAttack detects burst attack patterns
func (apd *AttackPatternDetector) detectBurstAttack(ctx context.Context, clientID, ipAddress string) *AttackSignature {
	now := time.Now()

	// Check for burst patterns in 5-second windows
	var burstCounts []int64
	for i := 0; i < 12; i++ { // Check last 12 5-second windows (1 minute)
		windowEnd := now.Add(-time.Duration(i) * 5 * time.Second)
		windowStart := windowEnd.Add(-5 * time.Second)

		requestKey := fmt.Sprintf("herald:ddos:requests:%s", clientID)
		count, _ := apd.redisClient.ZCount(ctx, requestKey,
			strconv.FormatInt(windowStart.Unix(), 10),
			strconv.FormatInt(windowEnd.Unix(), 10)).Result()

		burstCounts = append(burstCounts, count)
	}

	// Analyze burst patterns
	var highBursts int
	var totalRequests int64
	burstThreshold := int64(100) // 100 requests in 5 seconds

	for _, count := range burstCounts {
		totalRequests += count
		if count > burstThreshold {
			highBursts++
		}
	}

	// Detect burst attack if multiple high bursts with gaps
	if highBursts >= 3 && totalRequests > 500 {
		confidence := float64(highBursts) / 12.0 * 2.0 // Scale by burst frequency
		if confidence > 1.0 {
			confidence = 1.0
		}

		return &AttackSignature{
			Type:       "burst",
			Confidence: confidence,
			Evidence: []string{
				fmt.Sprintf("Multiple burst windows: %d high-volume windows", highBursts),
				fmt.Sprintf("Total requests: %d", totalRequests),
				fmt.Sprintf("Burst pattern detected over 1-minute window"),
			},
			Detected: now,
			Severity: apd.getSeverity(confidence),
		}
	}

	return nil
}

// detectSlowlorisAttack detects slowloris-style attacks
func (apd *AttackPatternDetector) detectSlowlorisAttack(ctx context.Context, clientID, ipAddress string) *AttackSignature {
	// Check for many long-running connections from same IP
	connectionKey := fmt.Sprintf("herald:ddos:connections:%s", ipAddress)

	// This would typically track connection timestamps and durations
	// For now, we'll use a simplified approach
	now := time.Now()
	windowStart := now.Add(-10 * time.Minute)

	connectionCount, err := apd.redisClient.ZCount(ctx, connectionKey,
		strconv.FormatInt(windowStart.Unix(), 10), "+inf").Result()
	if err != nil {
		return nil
	}

	// Check for sustained, low-rate connections
	if connectionCount > 100 { // Many slow connections
		confidence := 0.6 // Medium confidence for slowloris detection

		return &AttackSignature{
			Type:       "slowloris",
			Confidence: confidence,
			Evidence: []string{
				fmt.Sprintf("Sustained connections: %d over 10 minutes", connectionCount),
				"Pattern consistent with slowloris attack",
			},
			Detected: now,
			Severity: "medium",
		}
	}

	return nil
}

// detectDistributedAttack detects distributed attacks across IPs
func (apd *AttackPatternDetector) detectDistributedAttack(ctx context.Context, ipAddress string) *AttackSignature {
	now := time.Now()

	// Check for correlated attacks from multiple IPs in same subnet
	subnet := apd.getSubnet(ipAddress)
	if subnet == "" {
		return nil
	}

	// Count active IPs in subnet attacking
	subnetKey := fmt.Sprintf("herald:ddos:subnet:%s", subnet)
	activeIPs, err := apd.redisClient.ZCard(ctx, subnetKey).Result()
	if err != nil {
		return nil
	}

	// Distributed attack threshold
	if activeIPs > 10 { // More than 10 IPs from same subnet
		confidence := float64(activeIPs) / 50.0 // Scale by IP count
		if confidence > 1.0 {
			confidence = 1.0
		}

		return &AttackSignature{
			Type:       "distributed",
			Confidence: confidence,
			Evidence: []string{
				fmt.Sprintf("Multiple attacking IPs from subnet: %d", activeIPs),
				fmt.Sprintf("Subnet: %s", subnet),
				"Coordinated distributed attack pattern",
			},
			Detected: now,
			Severity: apd.getSeverity(confidence),
		}
	}

	return nil
}

// detectGamingAttackPatterns detects gaming-specific attack patterns
func (apd *AttackPatternDetector) detectGamingAttackPatterns(ctx context.Context, clientID, ipAddress string) *AttackSignature {
	now := time.Now()

	// Check for gaming API abuse patterns
	patterns := []string{
		"herald:ddos:gaming:analytics_spam",
		"herald:ddos:gaming:riot_api_abuse",
		"herald:ddos:gaming:export_spam",
	}

	var totalAbuse int64
	var evidenceList []string

	for _, pattern := range patterns {
		key := fmt.Sprintf("%s:%s", pattern, clientID)
		count, _ := apd.redisClient.Get(ctx, key).Int64()
		if count > 50 { // Threshold for gaming API abuse
			totalAbuse += count
			evidenceList = append(evidenceList, fmt.Sprintf("%s: %d requests", strings.Split(pattern, ":")[3], count))
		}
	}

	if totalAbuse > 100 {
		confidence := float64(totalAbuse) / 500.0 // Scale by abuse level
		if confidence > 1.0 {
			confidence = 1.0
		}

		return &AttackSignature{
			Type:       "gaming_api_abuse",
			Confidence: confidence,
			Evidence:   evidenceList,
			Detected:   now,
			Severity:   apd.getSeverity(confidence),
		}
	}

	return nil
}

// TriggerMitigation triggers appropriate mitigation strategies
func (me *MitigationEngine) TriggerMitigation(ctx context.Context, clientID, ipAddress string, signature *AttackSignature) {
	now := time.Now()

	// Determine mitigation strategy based on attack type and severity
	var blockDuration time.Duration
	var reason string

	switch signature.Type {
	case "volumetric":
		blockDuration = me.config.BlockDuration
		if signature.Severity == "critical" {
			blockDuration = blockDuration * 3
		}
		reason = "Volumetric DDoS attack detected"

	case "burst":
		blockDuration = me.config.BlockDuration / 2
		reason = "Burst attack pattern detected"

	case "slowloris":
		blockDuration = me.config.BlockDuration * 2
		reason = "Slowloris attack detected"

	case "distributed":
		blockDuration = me.config.BlockDuration * 4
		reason = "Distributed DDoS attack detected"

	case "gaming_api_abuse":
		blockDuration = me.config.BlockDuration
		reason = "Gaming API abuse detected"

	default:
		blockDuration = me.config.BlockDuration
		reason = "Suspicious activity detected"
	}

	// Apply mitigation
	me.blockClient(ctx, clientID, ipAddress, signature, blockDuration, reason)

	// Apply additional mitigations for severe attacks
	if signature.Severity == "critical" {
		me.applyAdvancedMitigation(ctx, ipAddress, signature)
	}
}

// blockClient blocks a client for specified duration
func (me *MitigationEngine) blockClient(ctx context.Context, clientID, ipAddress string, signature *AttackSignature, duration time.Duration, reason string) {
	key := fmt.Sprintf("herald:ddos:blocked:%s", clientID)
	blockedUntil := time.Now().Add(duration)

	blockInfo := map[string]interface{}{
		"reason":        reason,
		"blocked_at":    time.Now().Unix(),
		"blocked_until": blockedUntil.Unix(),
		"attack_type":   signature.Type,
		"severity":      signature.Severity,
		"confidence":    signature.Confidence,
		"ip_address":    ipAddress,
	}

	me.redisClient.HMSet(ctx, key, blockInfo)
	me.redisClient.Expire(ctx, key, duration+time.Hour) // Expire slightly after block duration

	// Also block by IP for severe attacks
	if signature.Severity == "critical" || signature.Severity == "high" {
		ipKey := fmt.Sprintf("herald:ddos:blocked_ip:%s", ipAddress)
		me.redisClient.HMSet(ctx, ipKey, blockInfo)
		me.redisClient.Expire(ctx, ipKey, duration+time.Hour)
	}
}

// applyAdvancedMitigation applies advanced mitigation for severe attacks
func (me *MitigationEngine) applyAdvancedMitigation(ctx context.Context, ipAddress string, signature *AttackSignature) {
	// Block entire subnet for distributed attacks
	if signature.Type == "distributed" {
		subnet := me.getSubnet(ipAddress)
		if subnet != "" {
			subnetKey := fmt.Sprintf("herald:ddos:blocked_subnet:%s", subnet)
			me.redisClient.Set(ctx, subnetKey, time.Now().Unix(), me.config.BlockDuration*2)
		}
	}

	// Implement CAPTCHA challenges
	challengeKey := fmt.Sprintf("herald:ddos:challenge:%s", ipAddress)
	me.redisClient.Set(ctx, challengeKey, "required", me.config.BlockDuration)

	// Notify administrators for critical attacks
	if signature.Severity == "critical" {
		me.notifyAdministrators(ctx, ipAddress, signature)
	}
}

// Helper methods

func (ddos *DDoSProtector) unblockClient(ctx context.Context, clientID string) {
	key := fmt.Sprintf("herald:ddos:blocked:%s", clientID)
	ddos.redisClient.Del(ctx, key)
}

func (ddos *DDoSProtector) logAttackDetection(ctx context.Context, clientID, ipAddress string, signature *AttackSignature) {
	logEntry := map[string]interface{}{
		"client_id":   clientID,
		"ip_address":  ipAddress,
		"attack_type": signature.Type,
		"confidence":  signature.Confidence,
		"severity":    signature.Severity,
		"evidence":    signature.Evidence,
		"detected_at": signature.Detected.Unix(),
		"platform":    "herald-lol",
	}

	logKey := fmt.Sprintf("herald:ddos:attack_log:%s", time.Now().Format("2006-01-02"))
	ddos.redisClient.LPush(ctx, logKey, logEntry)
	ddos.redisClient.Expire(ctx, logKey, 30*24*time.Hour) // Keep logs for 30 days
}

func (apd *AttackPatternDetector) getSeverity(confidence float64) string {
	if confidence >= 0.9 {
		return "critical"
	} else if confidence >= 0.7 {
		return "high"
	} else if confidence >= 0.5 {
		return "medium"
	}
	return "low"
}

func (apd *AttackPatternDetector) getSubnet(ipAddress string) string {
	parts := strings.Split(ipAddress, ".")
	if len(parts) >= 3 {
		return fmt.Sprintf("%s.%s.%s.0/24", parts[0], parts[1], parts[2])
	}
	return ""
}

func (me *MitigationEngine) getSubnet(ipAddress string) string {
	parts := strings.Split(ipAddress, ".")
	if len(parts) >= 3 {
		return fmt.Sprintf("%s.%s.%s.0/24", parts[0], parts[1], parts[2])
	}
	return ""
}

func (me *MitigationEngine) notifyAdministrators(ctx context.Context, ipAddress string, signature *AttackSignature) {
	// Send alert to administrators
	alertKey := "herald:ddos:critical_alerts"
	alert := map[string]interface{}{
		"ip_address":  ipAddress,
		"attack_type": signature.Type,
		"confidence":  signature.Confidence,
		"severity":    signature.Severity,
		"timestamp":   time.Now().Unix(),
		"platform":    "herald-lol",
	}

	me.redisClient.LPush(ctx, alertKey, alert)
	me.redisClient.LTrim(ctx, alertKey, 0, 99) // Keep last 100 alerts
}
