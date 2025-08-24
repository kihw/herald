# Herald.lol Gaming Analytics Platform - Backend

🎮 **L'infrastructure gaming analytics la plus performante pour League of Legends**

## 🚀 Architecture

### Stack Principal
- **Backend:** Go 1.23+ + Gin Web Framework + gRPC
- **Database:** PostgreSQL + Redis Cluster
- **Infrastructure:** Kubernetes + Docker
- **Deployment:** Blue-Green Strategy + Auto-scaling

### Performance Targets 🎯
- **⚡ Analytics Response:** < 5 secondes
- **🚀 API Latency:** < 1 seconde  
- **📊 Uptime:** 99.9%
- **👥 Concurrent Users:** 1M+ supported

## 🛠️ Quick Start

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- Kubernetes cluster
- PostgreSQL & Redis

### Local Development
```bash
# Clone and setup
git clone https://github.com/herald-lol/herald
cd herald/backend

# Install dependencies
go mod download

# Setup environment
cp .env.example .env

# Run services
docker-compose up -d postgres redis

# Run the server
go run main.go
```

### Production Deployment

#### Blue-Green Deployment
```bash
# Deploy to green environment with auto-switch
./scripts/blue-green-deploy.sh -i herald/gaming-analytics:v1.2.3 -t green -a

# Monitor deployment
./scripts/deployment-monitor.sh --realtime

# Rollback if needed
./scripts/blue-green-deploy.sh --rollback
```

#### Deploy to Kubernetes
```bash
# Apply base infrastructure
kubectl apply -f k8s/blue-green/namespace.yaml
kubectl apply -f k8s/blue-green/configmap.yaml

# Deploy blue environment
kubectl apply -f k8s/blue-green/herald-blue-deployment.yaml

# Deploy services and ingress
kubectl apply -f k8s/blue-green/herald-service.yaml
kubectl apply -f k8s/blue-green/ingress.yaml

# Enable auto-scaling
kubectl apply -f k8s/blue-green/hpa.yaml
```

## 🎮 Gaming Analytics Services

### Core Services
- **Analytics Service:** KDA, CS/min, Vision Score analysis
- **Match Processing:** Real-time match data processing
- **Riot API Integration:** Rate-limited Riot Games API client
- **Real-time Service:** WebSocket streaming for live data

### gRPC Services
```bash
# Generate protobuf code
protoc --go_out=. --go-grpc_out=. api/proto/analytics.proto
protoc --go_out=. --go-grpc_out=. api/proto/match.proto
protoc --go_out=. --go-grpc_out=. api/proto/riot.proto

# Run gRPC server
go run cmd/grpc-server/main.go
```

## 📊 Monitoring & Observability

### Health Checks
```bash
# Check deployment status
./scripts/blue-green-deploy.sh --check

# Monitor performance (5 minute window)
./scripts/deployment-monitor.sh --monitor 300

# Real-time monitoring
./scripts/deployment-monitor.sh --realtime
```

### Metrics Endpoints
- **Health:** `GET /health`
- **Ready:** `GET /ready`
- **Metrics:** `GET /metrics` (Prometheus format)

### Performance Monitoring
- Response time tracking (<5s target)
- gRPC server health monitoring  
- Resource usage optimization
- Auto-scaling based on gaming metrics

## 🔧 Development

### Project Structure
```
backend/
├── api/proto/           # gRPC protobuf definitions
├── cmd/                 # Application entrypoints
├── internal/            # Private application code
│   ├── analytics/       # Gaming analytics engine
│   ├── auth/           # Authentication & authorization
│   ├── grpc/           # gRPC server implementations  
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   ├── services/       # Business logic services
│   └── websocket/      # Real-time WebSocket handlers
├── k8s/                # Kubernetes manifests
├── scripts/            # Deployment scripts
└── tests/              # Test files
```

### Testing
```bash
# Run all tests
go test ./... -v -cover

# Run gaming analytics tests
go test ./internal/analytics/... -v

# Benchmark performance tests
go test -bench=. -benchmem ./internal/analytics/
```

### Code Quality
```bash
# Lint code
golangci-lint run

# Format code
go fmt ./...

# Security scan
gosec ./...
```

## 🔐 Security & Compliance

### Gaming Data Protection
- GDPR compliance for EU players
- Riot Games ToS compliance
- API key security with vault storage
- Player data anonymization

### Infrastructure Security
- Zero Trust architecture
- End-to-end encryption (AES-256 + TLS 1.3)
- OAuth 2.0 + MFA authentication
- Complete audit trail

## 🚀 Blue-Green Deployment Strategy

### Deployment Process
1. **Prepare:** Deploy new version to inactive environment
2. **Validate:** Comprehensive health checks + performance tests
3. **Switch:** Zero-downtime traffic routing
4. **Monitor:** Real-time performance monitoring
5. **Rollback:** Instant rollback capability if issues detected

### Gaming-Optimized Health Checks
- **Response Time:** < 5s analytics target
- **gRPC Health:** Server connectivity validation
- **Gaming Metrics:** Performance-specific checks
- **Load Testing:** Concurrent user simulation

## 📈 Scaling Configuration

### Auto-scaling Triggers
- CPU utilization > 70%
- Memory utilization > 80%
- Analytics response time > 4s
- Requests per second > 50/pod

### Gaming Performance Optimization
- Connection pooling for Riot API
- Redis caching for frequent queries
- Horizontal scaling for concurrent users
- Performance monitoring with alerting

## 🎯 Gaming Analytics Features

### League of Legends Analytics
- **KDA Analysis:** Comprehensive kill/death/assist metrics
- **CS/min Tracking:** Creep score optimization insights
- **Vision Score:** Map control and warding analysis
- **Champion Mastery:** Performance by champion
- **Rank Progression:** Climbing analysis and predictions

### Real-time Capabilities
- Live match streaming
- Real-time performance alerts
- Dynamic coaching recommendations
- Instant match result processing

## 🔗 API Endpoints

### Gaming Analytics
```bash
# Player KDA analysis
GET /analytics/kda/{playerID}?timeRange=30d&champion=Jinx

# CS per minute analysis  
GET /analytics/cs/{playerID}?position=ADC&timeRange=7d

# Performance comparison
GET /analytics/compare/{playerID}?timeRange=30d

# Match processing
POST /matches/process
```

### Real-time WebSocket
```javascript
// Connect to real-time gaming updates
const socket = new WebSocket('wss://api.herald.lol/ws');

// Subscribe to match updates
socket.send(JSON.stringify({
  action: 'watch_match',
  data: { match_id: 'NA1_1234567890' }
}));
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-gaming-feature`)
3. Commit changes (`git commit -m 'Add amazing gaming feature'`)
4. Push to branch (`git push origin feature/amazing-gaming-feature`)
5. Create Pull Request

### Development Guidelines
- Follow Go best practices and idioms
- Maintain <5s response time for analytics
- Include comprehensive tests
- Update documentation
- Gaming-first approach to features

## 📄 License

MIT License - see [LICENSE](LICENSE) file

---

## 🎮 Herald.lol Mission

**Démocratiser l'accès aux outils d'analyse gaming professionnels pour tous les joueurs de League of Legends.**

🌟 **Vision:** Devenir la référence mondiale pour l'analytics gaming multi-jeux en unifiant l'écosystème Riot Games.

---

**Built with ❤️ for the gaming community** | **Performance-first architecture** | **Cloud-native scalability**