# Herald.lol Performance Testing with k6

Performance testing suite for Herald.lol gaming analytics platform using k6.

## 🎯 Herald.lol Performance Requirements

- **Analytics Response Time:** <5 seconds
- **UI Load Time:** <2 seconds  
- **Uptime Target:** 99.9%
- **Concurrent Users:** 1M+ support
- **Gaming Metrics:** KDA, CS/min, Vision Score, Damage Share, Gold Efficiency

## 🚀 Quick Start

### Prerequisites

1. **Install k6:**
   ```bash
   # Linux (Ubuntu/Debian)
   npm run install-k6-linux
   
   # macOS
   npm run install-k6-mac
   
   # Windows
   npm run install-k6-windows
   ```

2. **Start Herald.lol services:**
   ```bash
   # Backend
   cd ../../backend && go run main.go
   
   # Frontend  
   cd ../../frontend && npm run dev
   ```

### Running Tests

#### Interactive Mode
```bash
./run-tests.sh
```

#### Command Line Mode
```bash
# Gaming analytics load test
npm run test:load

# Gaming stress test (extreme load)
npm run test:stress

# Gaming spike test (sudden load spikes)  
npm run test:spike

# Run all tests
npm run test:all

# Quick health check
npm run test:health
```

## 📊 Test Scenarios

### 1. Gaming Analytics Load Test (`gaming-analytics-load.js`)

**Purpose:** Standard load testing for Herald.lol gaming analytics

**Scenarios:**
- Analytics Dashboard Load (100-1000 VUs)
- Gaming API Stress (1000 RPS)
- Gaming Metrics Spike (100-2000 RPS)

**Key Metrics:**
- Analytics response time <5s
- Gaming API error rate <5%
- Concurrent user support 1M+

**Usage:**
```bash
k6 run gaming-analytics-load.js
```

### 2. Gaming Stress Test (`gaming-stress-test.js`)

**Purpose:** Extreme load conditions and system limits

**Scenarios:**
- Post-game Analysis Rush (up to 15k VUs)
- Gaming API Endurance (30min test)
- Analytics Burst Load (5000 RPS spikes)

**Key Metrics:**
- System failure rate <10%
- Graceful degradation under extreme load
- Memory and CPU stress testing

**Usage:**
```bash
k6 run gaming-stress-test.js
```

### 3. Gaming Spike Test (`gaming-spike-test.js`)

**Purpose:** Sudden load spikes and recovery testing

**Scenarios:**
- Tournament End Spike (100 → 5000 VUs in 10s)
- Rank Reset Spike (season start surge)
- Patch Analysis Spike (new meta rush)

**Key Metrics:**
- Analytics availability >95% during spikes
- Recovery time <10s
- Auto-scaling response

**Usage:**
```bash
k6 run gaming-spike-test.js
```

## 🎮 Gaming Test Data

### Test Scenarios
- **Tournament End:** Massive concurrent post-game analysis
- **Rank Reset:** Season start ranking checks
- **Patch Analysis:** New meta analysis rush
- **Daily Peak:** Regular evening gaming hours

### Gaming Metrics Tested
- **KDA Analysis:** Kill/Death/Assist ratio calculations
- **CS/min Analysis:** Creep Score per minute farming efficiency
- **Vision Score:** Map control and ward placement analysis
- **Damage Analysis:** Team fight contribution and damage breakdown
- **Gold Efficiency:** Economic performance and item optimization

### User Scenarios
- Gaming dashboard browsing
- Match analysis requests
- Team composition optimization
- Counter-pick analysis
- Skill progression tracking

## 📈 Performance Thresholds

### Response Times
```javascript
thresholds: {
  // Herald.lol gaming requirements
  'analytics_response_time': ['p(95)<5000'], // <5s analytics
  'http_req_duration{group:::gaming_analytics}': ['p(95)<5000'],
  'http_req_duration{group:::match_analysis}': ['p(95)<3000'],
  'http_req_duration{group:::riot_api}': ['p(95)<2000'],
}
```

### Error Rates
```javascript
thresholds: {
  'http_req_failed': ['rate<0.01'], // 99.9% uptime
  'gaming_api_errors': ['rate<0.05'], // <5% gaming errors
}
```

### Gaming Metrics
```javascript
thresholds: {
  'analytics_availability': ['rate>0.95'], // 95% analytics availability
  'concurrent_gaming_users': ['value>1000'], // 1K+ concurrent support
}
```

## 🔧 Configuration

### Environment Variables
```bash
export API_BASE_URL="http://localhost:8080"
export FRONTEND_URL="http://localhost:3000" 
export TEST_DURATION="5m"
export MAX_VUS="1000"
```

### Gaming Test Configuration
```javascript
// Herald.lol specific settings
const GAMING_CONFIG = {
  analytics_timeout: 5000,    // <5s requirement
  ui_timeout: 2000,           // <2s requirement  
  uptime_target: 0.999,       // 99.9% uptime
  concurrent_target: 1000000, // 1M+ users
};
```

## 📊 Results Analysis

### Output Files
- **JSON Results:** `results/gaming-load_TIMESTAMP.json`
- **HTML Reports:** `results/gaming-load_TIMESTAMP.html`
- **Summary Report:** `results/test_summary_TIMESTAMP.md`

### Key Metrics to Monitor
1. **Analytics Response Time** - Must be <5s
2. **Gaming API Error Rate** - Should be <1%  
3. **Concurrent User Support** - Target 1M+
4. **System Recovery Time** - <10s after spikes
5. **Riot API Integration** - Rate limiting compliance

### Performance Alerts
```javascript
// Herald.lol performance violations
if (analytics_response_time > 5000) {
  alert('Analytics exceeds 5s requirement');
}

if (error_rate > 0.01) {
  alert('System availability below 99.9%');
}
```

## 🎯 Gaming-Specific Testing

### Riot API Integration
- Rate limiting compliance testing
- API key rotation validation
- Match data synchronization performance

### Gaming Analytics Performance
- KDA calculation optimization
- CS/min analysis efficiency
- Vision heatmap generation speed
- Damage breakdown computation
- Gold efficiency algorithms

### Gaming User Patterns
- Post-game analysis rush patterns
- Rank checking surge simulation
- Tournament viewing spikes
- Patch analysis waves

## 🛠️ Troubleshooting

### Common Issues

1. **k6 not found:**
   ```bash
   npm run install-k6
   ```

2. **API connection failed:**
   ```bash
   # Check Herald.lol backend is running
   curl http://localhost:8080/health
   ```

3. **High error rates:**
   - Check server resource usage
   - Verify database connections
   - Monitor Riot API rate limits

4. **Performance degradation:**
   - Review system metrics
   - Check database query performance
   - Validate caching strategies

### Debug Mode
```bash
# Run with debug output
K6_LOG_LEVEL=debug k6 run gaming-analytics-load.js

# Detailed HTTP logging
k6 run --http-debug gaming-analytics-load.js
```

## 🚀 CI/CD Integration

### GitHub Actions Example
```yaml
name: Herald.lol Performance Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  performance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install k6
        run: |
          sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6
      
      - name: Run Herald.lol Performance Tests
        run: |
          cd performance/k6
          ./run-tests.sh all
```

## 📚 Resources

- [k6 Documentation](https://k6.io/docs/)
- [Herald.lol Performance Requirements](../../CLAUDE.md)
- [Gaming Analytics Architecture](../../docs/architecture.md)
- [Riot Games API Documentation](https://developer.riotgames.com/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add new gaming performance scenarios
4. Test with Herald.lol requirements
5. Submit a pull request

## 📄 License

MIT - See [LICENSE](../../LICENSE) for details.

---

**Herald.lol Performance Testing Suite** - Ensuring optimal gaming analytics performance at scale 🎮