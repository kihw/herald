# Monitoring et Observabilité

## Vue d'Ensemble de l'Observabilité

Herald.lol implémente une **stratégie d'observabilité complète** basée sur les trois piliers fondamentaux : métriques, logs et traces. Cette approche holistique garantit une visibilité totale sur la performance, la santé et l'expérience utilisateur de la plateforme.

## Architecture de Monitoring

### Stack d'Observabilité

#### Prometheus + Grafana Core
- **Prometheus** : Collecte et stockage métriques time-series
- **Grafana** : Visualisation et dashboards interactifs
- **AlertManager** : Gestion et routage des alertes
- **Node Exporter** : Métriques infrastructure système

#### ELK Stack pour Logging
- **Elasticsearch** : Stockage et indexation logs
- **Logstash** : Processing et enrichissement logs
- **Kibana** : Analyse et visualisation logs
- **Filebeat** : Collecte logs distribuée

#### Distributed Tracing
- **Jaeger** : Tracing distribué pour microservices
- **OpenTelemetry** : Instrumentation standardisée
- **Zipkin** : Tracing alternatif pour certains composants

## Métriques et KPIs

### Métriques Infrastructure

#### System Metrics
```yaml
# Prometheus Configuration - Infrastructure
- job_name: 'node-exporter'
  static_configs:
  - targets: ['node-exporter:9100']
  scrape_interval: 15s
  metrics_path: /metrics
  
  metric_relabel_configs:
  - source_labels: [__name__]
    regex: 'node_(cpu|memory|disk|network)_.*'
    target_label: __tmp_keep
    replacement: 'true'
```

#### Métriques Système Critiques
- **CPU Utilization** : Utilisation CPU par core et global
- **Memory Usage** : Mémoire utilisée/disponible avec swap
- **Disk I/O** : IOPS, latence, throughput disques
- **Network Traffic** : Bande passante, paquets, erreurs
- **Load Average** : Charge système 1m/5m/15m

### Métriques Application

#### Backend API Metrics
```go
// Prometheus metrics dans le code Go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "herald_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status_code"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "herald_http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    databaseConnections = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "herald_database_connections",
            Help: "Number of database connections",
        },
        []string{"state"},
    )
    
    riotApiCalls = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "herald_riot_api_calls_total",
            Help: "Total Riot API calls",
        },
        []string{"endpoint", "region", "status"},
    )
)
```

#### Métriques Business
- **User Registrations** : Nouvelles inscriptions par période
- **Data Synchronizations** : Syncs réussies/échouées par utilisateur
- **API Rate Limits** : Utilisation quotas Riot API
- **Export Requests** : Demandes d'export par format
- **Feature Usage** : Utilisation fonctionnalités par utilisateur

### Métriques Frontend

#### Real User Monitoring (RUM)
```typescript
// RUM Metrics Collection
interface PerformanceMetrics {
  // Core Web Vitals
  firstContentfulPaint: number;
  largestContentfulPaint: number;
  firstInputDelay: number;
  cumulativeLayoutShift: number;
  
  // Custom Business Metrics
  timeToInteractive: number;
  routeChangeTime: number;
  apiResponseTime: number;
  errorRate: number;
}

class MetricsCollector {
  collectWebVitals() {
    // Collection automatique Web Vitals
    getCLS(this.sendMetric.bind(this));
    getFCP(this.sendMetric.bind(this));
    getFID(this.sendMetric.bind(this));
    getLCP(this.sendMetric.bind(this));
  }
  
  trackUserInteraction(action: string, context: any) {
    // Tracking interactions utilisateur
    this.sendEvent({
      type: 'user_interaction',
      action,
      context,
      timestamp: Date.now()
    });
  }
}
```

## Dashboards et Visualisation

### Infrastructure Dashboards

#### System Overview Dashboard
```json
{
  "dashboard": {
    "title": "Herald.lol - Infrastructure Overview",
    "panels": [
      {
        "title": "CPU Usage",
        "type": "stat",
        "targets": [
          {
            "expr": "100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)",
            "legendFormat": "CPU Usage %"
          }
        ]
      },
      {
        "title": "Memory Usage",
        "type": "stat", 
        "targets": [
          {
            "expr": "(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100",
            "legendFormat": "Memory Usage %"
          }
        ]
      },
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(herald_http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      }
    ]
  }
}
```

#### Kubernetes Cluster Dashboard
- **Pod Status** : État pods par namespace
- **Resource Utilization** : CPU/Memory par node
- **Network Traffic** : Trafic inter-pods
- **Storage Usage** : Utilisation PVC et volumes

### Application Performance Dashboards

#### API Performance Dashboard
```yaml
# Grafana Dashboard Configuration
panels:
  - title: "API Response Times"
    type: heatmap
    datasource: prometheus
    targets:
      - expr: rate(herald_http_request_duration_seconds_bucket[5m])
        legendFormat: "{{le}}"
    
  - title: "Error Rate by Endpoint"
    type: graph
    targets:
      - expr: rate(herald_http_requests_total{status_code=~"5.."}[5m]) / rate(herald_http_requests_total[5m])
        legendFormat: "{{endpoint}}"
        
  - title: "Database Performance"
    type: graph
    targets:
      - expr: rate(herald_database_queries_total[5m])
        legendFormat: "Queries/sec"
      - expr: herald_database_connections{state="active"}
        legendFormat: "Active Connections"
```

#### Business Metrics Dashboard
- **User Activity** : Utilisateurs actifs par période
- **Gaming Analytics Usage** : Utilisation fonctionnalités analytics
- **Data Pipeline Health** : Santé pipeline synchronisation
- **Revenue Metrics** : Métriques business et conversion

### Gaming-Specific Dashboards

#### Riot API Integration Health
```promql
# Queries Prometheus pour Riot API
# Taux de succès par région
rate(herald_riot_api_calls_total{status="success"}[5m]) / rate(herald_riot_api_calls_total[5m])

# Latence moyenne par endpoint
rate(herald_riot_api_duration_seconds_sum[5m]) / rate(herald_riot_api_duration_seconds_count[5m])

# Rate limiting status
herald_riot_api_rate_limit_remaining / herald_riot_api_rate_limit_total
```

## Alerting et Notification

### Règles d'Alerte Critiques

#### Infrastructure Alerts
```yaml
# prometheus-alerts.yml
groups:
- name: infrastructure.critical
  rules:
  - alert: HostDown
    expr: up == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Host {{ $labels.instance }} is down"
      description: "Host has been down for more than 5 minutes"
      
  - alert: HighCPUUsage
    expr: 100 - (avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 85
    for: 10m
    labels:
      severity: warning
    annotations:
      summary: "High CPU usage on {{ $labels.instance }}"
      
  - alert: HighMemoryUsage
    expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 90
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High memory usage on {{ $labels.instance }}"
```

#### Application Alerts
```yaml
- name: application.critical
  rules:
  - alert: HighErrorRate
    expr: rate(herald_http_requests_total{status_code=~"5.."}[5m]) / rate(herald_http_requests_total[5m]) > 0.1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      
  - alert: DatabaseConnectionsHigh
    expr: herald_database_connections{state="active"} > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High number of database connections"
      
  - alert: RiotAPIRateLimitApproaching
    expr: herald_riot_api_rate_limit_remaining / herald_riot_api_rate_limit_total < 0.1
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "Riot API rate limit approaching"
```

### Notification Channels

#### Multi-Channel Alerting
```yaml
# alertmanager.yml
route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'default'
  routes:
  - match:
      severity: critical
    receiver: 'critical-alerts'
  - match:
      severity: warning
    receiver: 'warning-alerts'

receivers:
- name: 'critical-alerts'
  slack_configs:
  - api_url: 'SLACK_WEBHOOK_URL'
    channel: '#alerts-critical'
    title: 'Critical Alert - Herald.lol'
    text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'
  pagerduty_configs:
  - service_key: 'PAGERDUTY_SERVICE_KEY'
    description: '{{ .GroupLabels.alertname }}'
    
- name: 'warning-alerts'
  slack_configs:
  - api_url: 'SLACK_WEBHOOK_URL'
    channel: '#alerts-warning'
    title: 'Warning Alert - Herald.lol'
```

## Logging et Audit

### Structured Logging

#### Backend Logging Structure
```go
// Structured logging avec logrus
type LogEntry struct {
    Level       string    `json:"level"`
    Timestamp   time.Time `json:"timestamp"`
    Message     string    `json:"message"`
    Component   string    `json:"component"`
    UserID      string    `json:"user_id,omitempty"`
    RequestID   string    `json:"request_id,omitempty"`
    Endpoint    string    `json:"endpoint,omitempty"`
    Duration    int64     `json:"duration_ms,omitempty"`
    Error       string    `json:"error,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func LogAPIRequest(c *gin.Context, duration time.Duration) {
    log.WithFields(logrus.Fields{
        "component":   "api",
        "method":      c.Request.Method,
        "endpoint":    c.Request.URL.Path,
        "status_code": c.Writer.Status(),
        "duration_ms": duration.Milliseconds(),
        "user_id":     getUserID(c),
        "request_id":  getRequestID(c),
        "ip_address":  c.ClientIP(),
    }).Info("API request completed")
}
```

#### Frontend Error Logging
```typescript
// Error boundary avec logging
class ErrorLogger {
  static logError(error: Error, context: any) {
    const errorLog = {
      level: 'error',
      timestamp: new Date().toISOString(),
      message: error.message,
      stack: error.stack,
      url: window.location.href,
      userAgent: navigator.userAgent,
      context: context,
      sessionId: getSessionId(),
      userId: getCurrentUserId()
    };
    
    // Send to logging service
    fetch('/api/logs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(errorLog)
    });
  }
}
```

### Security Audit Logging

#### Audit Events Tracking
```go
type AuditEvent struct {
    ID           string    `json:"id"`
    Timestamp    time.Time `json:"timestamp"`
    EventType    string    `json:"event_type"`
    UserID       string    `json:"user_id"`
    IPAddress    string    `json:"ip_address"`
    UserAgent    string    `json:"user_agent"`
    Resource     string    `json:"resource"`
    Action       string    `json:"action"`
    Result       string    `json:"result"`
    Details      map[string]interface{} `json:"details"`
}

func AuditLog(eventType, userID, action, resource, result string, details map[string]interface{}) {
    event := AuditEvent{
        ID:        generateUUID(),
        Timestamp: time.Now(),
        EventType: eventType,
        UserID:    userID,
        Action:    action,
        Resource:  resource,
        Result:    result,
        Details:   details,
    }
    
    // Log to audit system
    auditLogger.WithFields(logrus.Fields(event)).Info("Audit event")
}
```

## Performance Monitoring

### Application Performance Monitoring

#### Go Application Profiling
```go
// Performance profiling intégré
import _ "net/http/pprof"

func init() {
    if os.Getenv("ENABLE_PROFILING") == "true" {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
}

// Custom metrics collection
func instrumentHandler(handler gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        handler(c)
        
        duration := time.Since(start)
        httpRequestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration.Seconds())
        
        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
    }
}
```

#### Database Performance Monitoring
```go
// Database query monitoring
type QueryMonitor struct {
    db *sql.DB
}

func (qm *QueryMonitor) QueryWithMetrics(query string, args ...interface{}) (*sql.Rows, error) {
    start := time.Now()
    
    rows, err := qm.db.Query(query, args...)
    
    duration := time.Since(start)
    dbQueryDuration.WithLabelValues(
        getQueryType(query),
        getTableName(query),
    ).Observe(duration.Seconds())
    
    if err != nil {
        dbQueryErrors.WithLabelValues(
            getQueryType(query),
            "error",
        ).Inc()
    }
    
    return rows, err
}
```

### SLA et Objectives

#### Service Level Indicators (SLIs)
- **Availability** : 99.9% uptime (8.76 heures downtime/an)
- **Latency** : P95 < 500ms pour API endpoints
- **Error Rate** : < 0.1% pour requests utilisateur
- **Throughput** : Support 10,000 req/sec en peak

#### Service Level Objectives (SLOs)
```yaml
slos:
  api_availability:
    target: 99.9%
    measurement_window: 30d
    
  api_latency:
    target: 95% des requêtes < 500ms
    measurement_window: 24h
    
  data_sync_success:
    target: 99.5% sync réussies
    measurement_window: 7d
    
  riot_api_integration:
    target: 99% disponibilité
    measurement_window: 24h
```

Cette stratégie de monitoring complète garantit une observabilité totale et une réactivité optimale pour maintenir les plus hauts standards de performance et fiabilité de Herald.lol.