# Containerisation Docker et Orchestration Kubernetes

## Vue d'Ensemble de la Containerisation

Herald.lol implémente une **stratégie de containerisation native** qui optimise le déploiement, la scalabilité et la maintenance à travers tous les environnements. Cette approche container-first garantit la portabilité et la cohérence des déploiements.

## Architecture Docker Multi-Stage

### Dockerfile Backend Optimisé

#### Build Stage Multi-Étapes
```dockerfile
# Stage 1: Dependencies et Build
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Cache des dépendances pour builds incrémentaux
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o herald-server .

# Stage 2: Runtime Minimal
FROM alpine:3.18 AS runtime

# Installation certificats SSL et timezone
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Création utilisateur non-privilégié
RUN addgroup -g 1001 herald && \
    adduser -D -s /bin/sh -u 1001 -G herald herald

# Copie binaire depuis builder stage
COPY --from=builder /app/herald-server .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

# Configuration sécurisée
USER herald
EXPOSE 8000

# Health check intégré
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

CMD ["./herald-server"]
```

#### Optimisations Docker Avancées
- **Multi-Stage Builds** : Réduction taille image finale 80%
- **Build Cache Optimization** : Cache layer intelligent pour CI/CD
- **Security Scanning** : Scan sécurité automatisé avec Snyk
- **Distroless Images** : Images sans OS pour sécurité maximale

### Dockerfile Frontend React

#### Build Production Optimisé
```dockerfile
# Stage 1: Build Dependencies
FROM node:20-alpine AS deps

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# Stage 2: Build Application
FROM node:20-alpine AS builder

WORKDIR /app
COPY . .
COPY --from=deps /app/node_modules ./node_modules
RUN npm run build

# Stage 3: Nginx Production
FROM nginx:1.25-alpine AS runtime

# Configuration Nginx optimisée
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /app/dist /usr/share/nginx/html

# Configuration sécurisée
RUN addgroup -g 1001 herald && \
    adduser -D -s /bin/sh -u 1001 -G herald herald

# Permissions et configuration
RUN chown -R herald:herald /usr/share/nginx/html && \
    chown -R herald:herald /var/cache/nginx

USER herald
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
```

### Configuration Nginx Production

#### Nginx Haute Performance
```nginx
user herald;
worker_processes auto;
worker_rlimit_nofile 65535;

events {
    worker_connections 4096;
    use epoll;
    multi_accept on;
}

http {
    # Performance et Sécurité
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    keepalive_requests 1000;
    
    # Compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types
        application/javascript
        application/json
        text/css
        text/xml;
    
    # Headers Sécurité
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # Backend Proxy
    upstream herald_backend {
        server herald-backend:8000 max_fails=3 fail_timeout=30s;
        keepalive 32;
    }
    
    server {
        listen 80;
        server_name herald.lol;
        root /usr/share/nginx/html;
        index index.html;
        
        # API Proxy
        location /api/ {
            proxy_pass http://herald_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
        }
        
        # Static Assets avec Cache
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
        
        # SPA Fallback
        location / {
            try_files $uri $uri/ /index.html;
        }
    }
}
```

## Orchestration Kubernetes

### Architecture Kubernetes Multi-Environment

#### Namespace Configuration
```yaml
# k8s/namespaces/herald-production.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: herald-production
  labels:
    environment: production
    app: herald
    cost-center: gaming-analytics
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: herald-quota
  namespace: herald-production
spec:
  hard:
    requests.cpu: "20"
    requests.memory: 40Gi
    limits.cpu: "40"
    limits.memory: 80Gi
    persistentvolumeclaims: "10"
    services: "10"
    secrets: "20"
---
apiVersion: v1
kind: LimitRange
metadata:
  name: herald-limits
  namespace: herald-production
spec:
  limits:
  - default:
      cpu: "2"
      memory: "4Gi"
    defaultRequest:
      cpu: "500m"
      memory: "1Gi"
    type: Container
```

### Déploiement Backend Kubernetes

#### Deployment Backend avec Auto-Scaling
```yaml
# k8s/backend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: herald-backend
  namespace: herald-production
  labels:
    app: herald-backend
    version: v2.1.0
spec:
  replicas: 6
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 2
      maxUnavailable: 1
  selector:
    matchLabels:
      app: herald-backend
  template:
    metadata:
      labels:
        app: herald-backend
        version: v2.1.0
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8000"
        prometheus.io/path: "/metrics"
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
      - name: herald-backend
        image: herald/backend:v2.1.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8000
          name: http
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: herald-secrets
              key: database-url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: herald-secrets
              key: redis-url
        - name: ENVIRONMENT
          value: "production"
        resources:
          requests:
            memory: "2Gi"
            cpu: "1"
          limits:
            memory: "4Gi"
            cpu: "2"
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sh", "-c", "sleep 15"]
      terminationGracePeriodSeconds: 30
---
apiVersion: v1
kind: Service
metadata:
  name: herald-backend
  namespace: herald-production
  labels:
    app: herald-backend
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8000
    name: http
  selector:
    app: herald-backend
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: herald-backend-hpa
  namespace: herald-production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: herald-backend
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

### Déploiement Frontend Kubernetes

#### Frontend avec CDN Integration
```yaml
# k8s/frontend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: herald-frontend
  namespace: herald-production
spec:
  replicas: 4
  selector:
    matchLabels:
      app: herald-frontend
  template:
    metadata:
      labels:
        app: herald-frontend
    spec:
      containers:
      - name: herald-frontend
        image: herald/frontend:v2.1.0
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: herald-frontend
  namespace: herald-production
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: herald-frontend
```

### Ingress et Load Balancing

#### Ingress Controller avec SSL
```yaml
# k8s/ingress/production-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: herald-ingress
  namespace: herald-production
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/use-regex: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "300"
spec:
  tls:
  - hosts:
    - herald.lol
    - www.herald.lol
    - api.herald.lol
    secretName: herald-tls
  rules:
  - host: herald.lol
    http:
      paths:
      - path: /api/(.*)
        pathType: Prefix
        backend:
          service:
            name: herald-backend
            port:
              number: 80
      - path: /
        pathType: Prefix
        backend:
          service:
            name: herald-frontend
            port:
              number: 80
  - host: api.herald.lol
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: herald-backend
            port:
              number: 80
```

## Storage et Persistence

### Persistent Volumes Configuration

#### PostgreSQL Cluster Storage
```yaml
# k8s/storage/postgres-pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: postgres-pv
  labels:
    type: ssd
    environment: production
spec:
  capacity:
    storage: 1Ti
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: fast-ssd
  awsElasticBlockStore:
    volumeID: vol-0a1b2c3d4e5f6g7h8
    fsType: ext4
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: herald-production
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 500Gi
  storageClassName: fast-ssd
```

#### Redis Cluster Persistence
```yaml
# k8s/storage/redis-cluster.yaml
apiVersion: redis.redis.opstreelabs.in/v1beta1
kind: RedisCluster
metadata:
  name: herald-redis
  namespace: herald-production
spec:
  clusterSize: 6
  kubernetesConfig:
    image: redis:7.2-alpine
    imagePullPolicy: Always
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 1
        memory: 2Gi
  storage:
    volumeClaimTemplate:
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 100Gi
        storageClassName: fast-ssd
  redisConfig:
    maxmemory: 1536mb
    maxmemory-policy: allkeys-lru
    save: "900 1 300 10 60 10000"
    appendonly: "yes"
    appendfsync: everysec
```

## Monitoring et Logging Kubernetes

### Prometheus Monitoring Stack

#### ServiceMonitor Configuration
```yaml
# k8s/monitoring/service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: herald-backend-monitor
  namespace: herald-production
  labels:
    app: herald-backend
spec:
  selector:
    matchLabels:
      app: herald-backend
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
    scrapeTimeout: 10s
---
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: herald-alerts
  namespace: herald-production
spec:
  groups:
  - name: herald.production
    rules:
    - alert: HeraldBackendDown
      expr: up{job="herald-backend"} == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Herald backend is down"
        description: "Herald backend has been down for more than 5 minutes"
    
    - alert: HeraldHighLatency
      expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High latency detected"
        description: "95th percentile latency is above 1 second"
```

### Centralized Logging

#### Fluent Bit DaemonSet
```yaml
# k8s/logging/fluent-bit.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluent-bit
  namespace: herald-production
spec:
  selector:
    matchLabels:
      app: fluent-bit
  template:
    metadata:
      labels:
        app: fluent-bit
    spec:
      containers:
      - name: fluent-bit
        image: fluent/fluent-bit:2.1
        ports:
        - containerPort: 2020
        env:
        - name: FLUENT_ELASTICSEARCH_HOST
          value: "elasticsearch.logging.svc.cluster.local"
        - name: FLUENT_ELASTICSEARCH_PORT
          value: "9200"
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
        - name: fluent-bit-config
          mountPath: /fluent-bit/etc/
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
      - name: fluent-bit-config
        configMap:
          name: fluent-bit-config
```

Cette architecture Docker et Kubernetes robuste garantit la haute disponibilité, la scalabilité et la maintenance optimale de Herald.lol en production.