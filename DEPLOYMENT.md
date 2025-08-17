# 🚀 Deployment Guide

This guide covers deployment options for the LoL Match Exporter in production environments.

## 📋 Prerequisites

- Python 3.8+
- Node.js 16+
- Valid Riot Games API Key
- 512MB+ RAM (2GB+ recommended)
- 1GB+ disk space

## 🐳 Docker Deployment (Recommended)

### 1. Create Dockerfile

```dockerfile
# Multi-stage build for frontend
FROM node:18-alpine as frontend-builder
WORKDIR /app/web

# Install dependencies needed for native modules
RUN apk add --no-cache python3 make g++

COPY web/package*.json ./

# Clear npm cache and install with force to handle optional dependencies
RUN npm cache clean --force && \
    rm -rf node_modules package-lock.json && \
    npm install --no-optional

COPY web/ ./
RUN npm run build

# Python backend
FROM python:3.11-slim
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy backend code
COPY *.py ./
COPY --from=frontend-builder /app/web/dist ./web/dist

# Create jobs directory
RUN mkdir -p jobs

# Expose port
EXPOSE 8000

# Run server
CMD ["python", "server.py"]
```

### 2. Create docker-compose.yml

```yaml
version: '3.8'
services:
  lol-exporter:
    build: .
    ports:
      - "8000:8000"
    environment:
      - RIOT_API_KEY=${RIOT_API_KEY}
      - EXPORTER_API_KEY=${EXPORTER_API_KEY}
    volumes:
      - ./jobs:/app/jobs
      - ./logs:/app/logs
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Optional: Add reverse proxy
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - lol-exporter
    restart: unless-stopped
```

### 3. Deploy with Docker Compose

```bash
# Set environment variables
echo "RIOT_API_KEY=your_api_key_here" > .env
echo "EXPORTER_API_KEY=your_server_key_here" >> .env

# Build and start
docker-compose up -d

# View logs
docker-compose logs -f lol-exporter
```

## ☁️ Cloud Deployment Options

### AWS Deployment

#### Option 1: AWS App Runner
```bash
# Create apprunner.yaml
version: 1.0
runtime: python3
build:
  commands:
    build:
      - pip install -r requirements.txt
      - cd web && npm install && npm run build
run:
  runtime-version: 3.11
  command: python server.py
  network:
    port: 8000
    env: PORT
  env:
    - name: RIOT_API_KEY
      value: your_api_key
```

#### Option 2: ECS with Fargate
```json
{
  "family": "lol-exporter",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "lol-exporter",
      "image": "your-account.dkr.ecr.region.amazonaws.com/lol-exporter:latest",
      "portMappings": [
        {
          "containerPort": 8000,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "RIOT_API_KEY",
          "value": "your_api_key"
        }
      ]
    }
  ]
}
```

### Google Cloud Platform

#### Cloud Run Deployment
```bash
# Build and push image
gcloud builds submit --tag gcr.io/PROJECT-ID/lol-exporter

# Deploy to Cloud Run
gcloud run deploy lol-exporter \
  --image gcr.io/PROJECT-ID/lol-exporter \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars RIOT_API_KEY=your_api_key
```

### Azure Container Instances

```bash
# Create resource group
az group create --name lol-exporter-rg --location eastus

# Deploy container
az container create \
  --resource-group lol-exporter-rg \
  --name lol-exporter \
  --image your-registry/lol-exporter:latest \
  --ports 8000 \
  --environment-variables RIOT_API_KEY=your_api_key \
  --cpu 1 \
  --memory 2
```

## 🖥️ Traditional Server Deployment

### 1. System Setup (Ubuntu 20.04+)

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y python3 python3-pip nodejs npm nginx

# Install Python dependencies
pip3 install -r requirements.txt

# Build frontend
cd web
npm install
npm run build
cd ..
```

### 2. Systemd Service

Create `/etc/systemd/system/lol-exporter.service`:

```ini
[Unit]
Description=LoL Match Exporter
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/lol-exporter
Environment=RIOT_API_KEY=your_api_key
Environment=EXPORTER_API_KEY=your_server_key
ExecStart=/usr/bin/python3 server.py
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl enable lol-exporter
sudo systemctl start lol-exporter
sudo systemctl status lol-exporter
```

### 3. Nginx Configuration

Create `/etc/nginx/sites-available/lol-exporter`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL Configuration
    ssl_certificate /etc/ssl/certs/your-domain.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # Proxy to FastAPI
    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # SSE support
        proxy_buffering off;
        proxy_cache off;
    }

    # Static files
    location /static/ {
        alias /opt/lol-exporter/web/dist/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # File downloads
    location /export/ {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Large file support
        client_max_body_size 100M;
        proxy_read_timeout 300s;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/lol-exporter /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🔒 Security Considerations

### 1. Environment Variables
```bash
# Never commit API keys to version control
# Use environment variables or secret management
export RIOT_API_KEY="RGAPI-xxxxxxxx"
export EXPORTER_API_KEY="secure-random-string"
```

### 2. Rate Limiting
```python
# Add to server.py for additional protection
from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address

limiter = Limiter(key_func=get_remote_address)
app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)

@app.post("/export")
@limiter.limit("5/minute")
async def start_export(request: Request, req: ExportRequest):
    # ... existing code
```

### 3. CORS Configuration
```python
# Restrict CORS in production
app.add_middleware(
    CORSMiddleware,
    allow_origins=["https://your-domain.com"],  # Specific domain
    allow_credentials=True,
    allow_methods=["GET", "POST", "DELETE"],
    allow_headers=["*"],
)
```

### 4. HTTPS Enforcement
```python
# Add HTTPS redirect middleware
from fastapi.middleware.httpsredirect import HTTPSRedirectMiddleware

if os.getenv("ENVIRONMENT") == "production":
    app.add_middleware(HTTPSRedirectMiddleware)
```

## 📊 Monitoring & Logging

### 1. Application Metrics
```python
# Add to server.py
from prometheus_client import Counter, Histogram, generate_latest
import time

REQUEST_COUNT = Counter('http_requests_total', 'Total HTTP requests', ['method', 'endpoint'])
REQUEST_LATENCY = Histogram('http_request_duration_seconds', 'HTTP request latency')

@app.middleware("http")
async def metrics_middleware(request: Request, call_next):
    start_time = time.time()
    response = await call_next(request)
    duration = time.time() - start_time
    
    REQUEST_COUNT.labels(method=request.method, endpoint=request.url.path).inc()
    REQUEST_LATENCY.observe(duration)
    
    return response

@app.get("/metrics")
async def get_metrics():
    return Response(generate_latest(), media_type="text/plain")
```

### 2. Structured Logging
```python
import structlog

logger = structlog.get_logger()

@app.post("/export")
async def start_export(req: ExportRequest):
    logger.info("Export started", 
                riot_id=req.riotId, 
                platform=req.platform, 
                count=req.count)
```

### 3. Health Checks
```python
@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "timestamp": datetime.utcnow().isoformat(),
        "version": "2.1.0",
        "active_jobs": len(jobs),
        "cache_size": cache.get_stats()["size"] if 'cache' in globals() else 0
    }
```

## 🔧 Performance Tuning

### 1. Production Server Configuration
```python
# Use Gunicorn for production
# gunicorn_config.py
bind = "0.0.0.0:8000"
workers = 4
worker_class = "uvicorn.workers.UvicornWorker"
worker_connections = 1000
max_requests = 1000
max_requests_jitter = 100
timeout = 300
keepalive = 5
```

### 2. Database for Job Persistence (Optional)
```python
# Use Redis or PostgreSQL for job state
import redis

redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

class JobState(BaseModel):
    def save(self):
        redis_client.hset(f"job:{self.id}", mapping=self.dict())
    
    @classmethod
    def load(cls, job_id: str):
        data = redis_client.hgetall(f"job:{job_id}")
        return cls(**data) if data else None
```

## 📈 Scaling Considerations

### Horizontal Scaling
- Use load balancer (AWS ALB, nginx)
- Shared storage for job files (S3, NFS)
- Redis for job state synchronization
- Queue system for background processing (Celery, RQ)

### Vertical Scaling
- Increase CPU/memory based on concurrent exports
- Monitor rate limiter effectiveness
- Cache hit rates and API response times

## 🚨 Troubleshooting

### Common Issues
1. **High memory usage**: Reduce cache size, implement job cleanup
2. **Rate limiting**: Increase backoff multiplier, reduce concurrent requests
3. **Large exports failing**: Implement streaming responses, increase timeouts
4. **SSL certificate**: Use Let's Encrypt for free certificates

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=DEBUG
python server.py
```

This deployment guide ensures your LoL Match Exporter runs reliably in production! 🚀