# Sécurité et Conformité

## Vue d'Ensemble de la Sécurité

Herald.lol implémente une **stratégie de sécurité defense-in-depth** qui protège l'ensemble de la plateforme à tous les niveaux. Cette approche multicouche garantit la protection des données utilisateurs, la conformité réglementaire et la résistance aux menaces modernes.

## Framework de Sécurité

### Security by Design

#### Principes Fondamentaux
- **Zero Trust Architecture** : Aucune confiance implicite, vérification constante
- **Principle of Least Privilege** : Accès minimal nécessaire
- **Defense in Depth** : Couches de sécurité multiples
- **Privacy by Design** : Protection vie privée intégrée dès la conception

#### Security Development Lifecycle
- **Threat Modeling** : Modélisation menaces dès la conception
- **Security Code Review** : Revue code sécurisé automatisée et manuelle
- **Vulnerability Assessment** : Évaluation vulnérabilités continue
- **Penetration Testing** : Tests intrusion réguliers

### Architecture de Sécurité

#### Network Security
```yaml
# Network Segmentation
network_zones:
  dmz:
    - load_balancers
    - cdn_edge_nodes
    - reverse_proxies
  
  application_tier:
    - web_servers
    - api_gateways
    - application_servers
  
  data_tier:
    - databases
    - cache_servers
    - message_queues
    
  management:
    - monitoring_systems
    - backup_systems
    - admin_interfaces

# Firewall Rules
firewall_policies:
  inbound:
    - allow: [80, 443] from: internet to: dmz
    - allow: [8000-8999] from: dmz to: application_tier
    - allow: [5432, 6379] from: application_tier to: data_tier
    - deny: all from: internet to: [application_tier, data_tier]
    
  outbound:
    - allow: [443] from: application_tier to: external_apis
    - allow: [25, 587] from: application_tier to: email_servers
    - deny: all unnecessary outbound traffic
```

#### Application Security

##### Authentication et Authorization
```go
// Multi-Factor Authentication
type AuthService struct {
    jwtManager   *JWTManager
    totpManager  *TOTPManager
    sessionStore *SessionStore
}

func (as *AuthService) AuthenticateUser(credentials *LoginCredentials) (*AuthResult, error) {
    // 1. Validate credentials
    user, err := as.validateCredentials(credentials)
    if err != nil {
        return nil, err
    }
    
    // 2. Check for MFA requirement
    if user.MFAEnabled {
        return &AuthResult{
            Status: "mfa_required",
            TempToken: as.generateTempToken(user.ID),
        }, nil
    }
    
    // 3. Generate JWT tokens
    accessToken, err := as.jwtManager.GenerateAccessToken(user)
    if err != nil {
        return nil, err
    }
    
    refreshToken, err := as.jwtManager.GenerateRefreshToken(user)
    if err != nil {
        return nil, err
    }
    
    // 4. Create secure session
    session, err := as.sessionStore.CreateSession(user.ID, accessToken)
    if err != nil {
        return nil, err
    }
    
    return &AuthResult{
        Status: "success",
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        SessionID: session.ID,
    }, nil
}

// Role-Based Access Control
type RBACManager struct {
    permissions map[string][]string
    roles       map[string][]string
}

func (rbac *RBACManager) CheckPermission(userID, resource, action string) bool {
    userRoles := rbac.getUserRoles(userID)
    
    for _, role := range userRoles {
        rolePermissions := rbac.roles[role]
        for _, permission := range rolePermissions {
            if permission == fmt.Sprintf("%s:%s", resource, action) {
                return true
            }
        }
    }
    
    return false
}
```

##### Input Validation et Sanitization
```go
// Comprehensive input validation
type ValidationRules struct {
    Required    bool
    MinLength   int
    MaxLength   int
    Pattern     *regexp.Regexp
    Sanitize    bool
    XSSProtect  bool
    SQLProtect  bool
}

func ValidateAndSanitize(input string, rules ValidationRules) (string, error) {
    // Required validation
    if rules.Required && strings.TrimSpace(input) == "" {
        return "", errors.New("field is required")
    }
    
    // Length validation
    if len(input) < rules.MinLength || len(input) > rules.MaxLength {
        return "", errors.New("invalid length")
    }
    
    // Pattern validation
    if rules.Pattern != nil && !rules.Pattern.MatchString(input) {
        return "", errors.New("invalid format")
    }
    
    // XSS Protection
    if rules.XSSProtect {
        input = html.EscapeString(input)
    }
    
    // SQL Injection Protection
    if rules.SQLProtect {
        input = strings.ReplaceAll(input, "'", "''")
        input = strings.ReplaceAll(input, ";", "")
    }
    
    // General sanitization
    if rules.Sanitize {
        input = strings.TrimSpace(input)
        input = regexp.MustCompile(`[^\w\s-.]`).ReplaceAllString(input, "")
    }
    
    return input, nil
}
```

## Data Protection et Chiffrement

### Encryption at Rest

#### Database Encryption
```sql
-- PostgreSQL Transparent Data Encryption
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Encrypted sensitive columns
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    riot_puuid VARCHAR(78) UNIQUE NOT NULL,
    riot_id VARCHAR(16) NOT NULL,
    riot_tag VARCHAR(5) NOT NULL,
    email_encrypted BYTEA, -- PGP encrypted
    preferences_encrypted BYTEA, -- AES encrypted
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Encryption functions
CREATE OR REPLACE FUNCTION encrypt_email(email TEXT, key TEXT)
RETURNS BYTEA AS $$
BEGIN
    RETURN pgp_sym_encrypt(email, key);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION decrypt_email(encrypted_data BYTEA, key TEXT)
RETURNS TEXT AS $$
BEGIN
    RETURN pgp_sym_decrypt(encrypted_data, key);
END;
$$ LANGUAGE plpgsql;
```

#### File System Encryption
```bash
# LUKS encryption for data volumes
cryptsetup luksFormat /dev/xvdf
cryptsetup luksOpen /dev/xvdf herald_data
mkfs.ext4 /dev/mapper/herald_data
mount /dev/mapper/herald_data /data

# Backup encryption
gpg --cipher-algo AES256 --compress-algo 1 --s2k-mode 3 \
    --s2k-digest-algo SHA512 --s2k-count 65536 --symmetric \
    --output backup_encrypted.gpg backup_data.tar
```

### Encryption in Transit

#### TLS Configuration
```nginx
# Nginx TLS hardening
ssl_protocols TLSv1.3 TLSv1.2;
ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
ssl_prefer_server_ciphers off;
ssl_ecdh_curve secp384r1;

# HSTS
add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";

# Certificate Transparency
ssl_ct on;
ssl_ct_static_scts /path/to/scts;

# OCSP Stapling
ssl_stapling on;
ssl_stapling_verify on;
ssl_trusted_certificate /path/to/chain.pem;
```

#### Application-Level Encryption
```go
// End-to-end encryption for sensitive data
type EncryptionService struct {
    aesKey []byte
    rsaKey *rsa.PrivateKey
}

func (es *EncryptionService) EncryptSensitiveData(data []byte) ([]byte, error) {
    // Generate random AES key for this data
    aesKey := make([]byte, 32)
    if _, err := rand.Read(aesKey); err != nil {
        return nil, err
    }
    
    // Encrypt data with AES
    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return nil, err
    }
    
    encryptedData := gcm.Seal(nil, nonce, data, nil)
    
    // Encrypt AES key with RSA
    encryptedKey, err := rsa.EncryptOAEP(
        sha256.New(),
        rand.Reader,
        &es.rsaKey.PublicKey,
        aesKey,
        nil,
    )
    if err != nil {
        return nil, err
    }
    
    // Combine encrypted key + nonce + encrypted data
    result := append(encryptedKey, nonce...)
    result = append(result, encryptedData...)
    
    return result, nil
}
```

## Conformité Réglementaire

### GDPR Compliance

#### Data Processing Legal Basis
```go
type DataProcessingBasis string

const (
    Consent           DataProcessingBasis = "consent"
    Contract          DataProcessingBasis = "contract"
    LegalObligation   DataProcessingBasis = "legal_obligation"
    VitalInterests    DataProcessingBasis = "vital_interests"
    PublicTask        DataProcessingBasis = "public_task"
    LegitimateInterest DataProcessingBasis = "legitimate_interest"
)

type DataProcessingRecord struct {
    ID              string                 `json:"id"`
    UserID          string                 `json:"user_id"`
    DataCategory    string                 `json:"data_category"`
    ProcessingBasis DataProcessingBasis    `json:"processing_basis"`
    Purpose         string                 `json:"purpose"`
    Retention       time.Duration          `json:"retention"`
    Timestamp       time.Time              `json:"timestamp"`
    ConsentID       *string                `json:"consent_id,omitempty"`
}

func (dpr *DataProcessingRecord) IsRetentionExpired() bool {
    return time.Since(dpr.Timestamp) > dpr.Retention
}
```

#### Consent Management System
```go
type ConsentManager struct {
    db *sql.DB
}

type ConsentRecord struct {
    ID                string            `json:"id"`
    UserID            string            `json:"user_id"`
    ConsentType       string            `json:"consent_type"`
    Granted           bool              `json:"granted"`
    Timestamp         time.Time         `json:"timestamp"`
    ExpirationDate    *time.Time        `json:"expiration_date,omitempty"`
    WithdrawalDate    *time.Time        `json:"withdrawal_date,omitempty"`
    LegalBasis        string            `json:"legal_basis"`
    ProcessingDetails map[string]string `json:"processing_details"`
}

func (cm *ConsentManager) RecordConsent(userID, consentType string, granted bool, details map[string]string) error {
    consent := ConsentRecord{
        ID:                generateUUID(),
        UserID:            userID,
        ConsentType:       consentType,
        Granted:           granted,
        Timestamp:         time.Now(),
        LegalBasis:        "consent",
        ProcessingDetails: details,
    }
    
    if consentType == "marketing" {
        // Marketing consent expires after 2 years
        expiration := time.Now().AddDate(2, 0, 0)
        consent.ExpirationDate = &expiration
    }
    
    return cm.storeConsent(consent)
}

func (cm *ConsentManager) HasValidConsent(userID, consentType string) bool {
    consent, err := cm.getLatestConsent(userID, consentType)
    if err != nil {
        return false
    }
    
    // Check if consent is granted and not expired
    if !consent.Granted {
        return false
    }
    
    if consent.WithdrawalDate != nil {
        return false
    }
    
    if consent.ExpirationDate != nil && time.Now().After(*consent.ExpirationDate) {
        return false
    }
    
    return true
}
```

#### Data Subject Rights Implementation
```go
type DataSubjectRightsService struct {
    db             *sql.DB
    encryptionSvc  *EncryptionService
    auditLogger    *AuditLogger
}

// Right to Access
func (dsrs *DataSubjectRightsService) ExportUserData(userID string) (*UserDataExport, error) {
    export := &UserDataExport{
        UserID:    userID,
        RequestID: generateUUID(),
        Timestamp: time.Now(),
    }
    
    // Collect all user data
    profile, err := dsrs.getUserProfile(userID)
    if err != nil {
        return nil, err
    }
    export.Profile = profile
    
    matches, err := dsrs.getUserMatches(userID)
    if err != nil {
        return nil, err
    }
    export.Matches = matches
    
    analytics, err := dsrs.getUserAnalytics(userID)
    if err != nil {
        return nil, err
    }
    export.Analytics = analytics
    
    // Log the access request
    dsrs.auditLogger.LogDataAccess(userID, "data_export", export.RequestID)
    
    return export, nil
}

// Right to Erasure
func (dsrs *DataSubjectRightsService) DeleteUserData(userID string, reason string) error {
    tx, err := dsrs.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Mark user as deleted (soft delete for legal compliance)
    _, err = tx.Exec(`
        UPDATE users 
        SET deleted_at = NOW(), 
            deletion_reason = $1,
            email_encrypted = NULL,
            preferences_encrypted = NULL
        WHERE id = $2
    `, reason, userID)
    if err != nil {
        return err
    }
    
    // Anonymize related data
    _, err = tx.Exec(`
        UPDATE matches 
        SET user_id = 'anonymized-' || generate_random_uuid()
        WHERE user_id = $1
    `, userID)
    if err != nil {
        return err
    }
    
    // Log the deletion
    dsrs.auditLogger.LogDataDeletion(userID, reason)
    
    return tx.Commit()
}

// Right to Rectification
func (dsrs *DataSubjectRightsService) UpdateUserData(userID string, updates map[string]interface{}) error {
    // Validate updates
    for field, value := range updates {
        if err := dsrs.validateFieldUpdate(field, value); err != nil {
            return err
        }
    }
    
    // Apply updates
    for field, value := range updates {
        if err := dsrs.updateUserField(userID, field, value); err != nil {
            return err
        }
    }
    
    // Log the rectification
    dsrs.auditLogger.LogDataRectification(userID, updates)
    
    return nil
}
```

### CCPA Compliance

#### Consumer Rights Implementation
```go
type CCPAComplianceService struct {
    dataMapper      *DataMapper
    consentManager  *ConsentManager
    auditLogger     *AuditLogger
}

// Do Not Sell directive
func (ccpa *CCPAComplianceService) SetDoNotSell(userID string, doNotSell bool) error {
    record := ConsentRecord{
        UserID:      userID,
        ConsentType: "data_sale_opt_out",
        Granted:     !doNotSell, // Inverted because it's opt-out
        Timestamp:   time.Now(),
        LegalBasis:  "ccpa_consumer_right",
    }
    
    return ccpa.consentManager.storeConsent(record)
}

// Consumer disclosure requirements
func (ccpa *CCPAComplianceService) GetDataDisclosure(userID string) (*DataDisclosure, error) {
    disclosure := &DataDisclosure{
        UserID:    userID,
        Timestamp: time.Now(),
    }
    
    // Categories of data collected
    disclosure.DataCategories = []string{
        "gaming_performance_data",
        "account_information", 
        "device_information",
        "usage_analytics",
    }
    
    // Sources of data
    disclosure.DataSources = []string{
        "riot_games_api",
        "user_input",
        "device_sensors",
        "third_party_integrations",
    }
    
    // Business purposes
    disclosure.BusinessPurposes = []string{
        "gaming_analytics",
        "service_improvement",
        "fraud_prevention",
        "customer_support",
    }
    
    // Third parties data shared with
    disclosure.ThirdParties = []string{
        "cloud_infrastructure_providers",
        "analytics_service_providers",
        "customer_support_tools",
    }
    
    return disclosure, nil
}
```

## Security Monitoring et Incident Response

### Security Information and Event Management (SIEM)

#### Security Event Correlation
```go
type SecurityEventMonitor struct {
    eventStore    *EventStore
    ruleEngine    *SecurityRuleEngine
    alertManager  *AlertManager
}

type SecurityEvent struct {
    ID           string                 `json:"id"`
    Timestamp    time.Time              `json:"timestamp"`
    EventType    string                 `json:"event_type"`
    Severity     string                 `json:"severity"`
    Source       string                 `json:"source"`
    UserID       *string                `json:"user_id,omitempty"`
    IPAddress    string                 `json:"ip_address"`
    UserAgent    string                 `json:"user_agent"`
    Details      map[string]interface{} `json:"details"`
    ThreatScore  int                    `json:"threat_score"`
}

func (sem *SecurityEventMonitor) ProcessEvent(event SecurityEvent) error {
    // Store the event
    if err := sem.eventStore.Store(event); err != nil {
        return err
    }
    
    // Apply security rules
    threats := sem.ruleEngine.EvaluateEvent(event)
    
    // Generate alerts for high-severity threats
    for _, threat := range threats {
        if threat.Severity >= HIGH_SEVERITY {
            alert := SecurityAlert{
                ThreatID:    threat.ID,
                EventID:     event.ID,
                Severity:    threat.Severity,
                Description: threat.Description,
                Timestamp:   time.Now(),
            }
            
            if err := sem.alertManager.TriggerAlert(alert); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

#### Anomaly Detection
```go
type AnomalyDetector struct {
    mlModel      *MachineLearningModel
    baseline     *BaselineMetrics
    alertManager *AlertManager
}

func (ad *AnomalyDetector) DetectAnomalies(userID string, activity UserActivity) error {
    // Calculate activity metrics
    metrics := ActivityMetrics{
        LoginFrequency:     activity.LoginCount,
        APICallRate:       activity.APICallCount,
        DataVolumeAccess:  activity.DataAccessVolume,
        GeolocationChange: activity.LocationChanges,
        DeviceFingerprint: activity.DeviceFingerprint,
    }
    
    // Compare against user baseline
    baseline := ad.baseline.GetUserBaseline(userID)
    anomalyScore := ad.mlModel.CalculateAnomalyScore(metrics, baseline)
    
    // Trigger alert if anomaly score is high
    if anomalyScore > ANOMALY_THRESHOLD {
        alert := AnomalyAlert{
            UserID:       userID,
            AnomalyScore: anomalyScore,
            Metrics:      metrics,
            Baseline:     baseline,
            Timestamp:    time.Now(),
        }
        
        return ad.alertManager.TriggerAnomalyAlert(alert)
    }
    
    return nil
}
```

### Incident Response Framework

#### Automated Incident Response
```go
type IncidentResponseSystem struct {
    playbooks      map[string]*Playbook
    orchestrator   *ResponseOrchestrator
    notificationMgr *NotificationManager
}

type IncidentPlaybook struct {
    ID          string                 `json:"id"`
    TriggerType string                 `json:"trigger_type"`
    Severity    string                 `json:"severity"`
    Actions     []ResponseAction       `json:"actions"`
    Escalation  *EscalationPolicy      `json:"escalation"`
}

func (irs *IncidentResponseSystem) HandleSecurityIncident(incident SecurityIncident) error {
    // Find appropriate playbook
    playbook := irs.playbooks[incident.Type]
    if playbook == nil {
        return errors.New("no playbook found for incident type")
    }
    
    // Execute response actions
    for _, action := range playbook.Actions {
        if err := irs.orchestrator.ExecuteAction(action, incident); err != nil {
            log.Errorf("Failed to execute action %s: %v", action.ID, err)
            continue
        }
    }
    
    // Handle escalation if needed
    if incident.Severity == "critical" && playbook.Escalation != nil {
        return irs.escalateIncident(incident, playbook.Escalation)
    }
    
    return nil
}

func (irs *IncidentResponseSystem) escalateIncident(incident SecurityIncident, policy *EscalationPolicy) error {
    // Notify security team
    notification := IncidentNotification{
        IncidentID:  incident.ID,
        Severity:    incident.Severity,
        Description: incident.Description,
        Timestamp:   time.Now(),
    }
    
    return irs.notificationMgr.NotifySecurityTeam(notification)
}
```

Cette stratégie de sécurité et conformité robuste garantit que Herald.lol maintient les plus hauts standards de protection des données et de conformité réglementaire tout en offrant une expérience utilisateur fluide et sécurisée.