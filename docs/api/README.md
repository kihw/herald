# Herald.lol Gaming Analytics API Documentation

Welcome to the Herald.lol Gaming Analytics API documentation. This comprehensive API provides access to League of Legends and Teamfight Tactics analytics, player statistics, team management, and gaming insights.

## 🎮 Overview

Herald.lol is the premier gaming analytics platform designed to democratize access to professional-grade esports analytics tools. Our API enables developers, analysts, and gaming applications to integrate powerful gaming data and insights.

### Key Features

- **Gaming Analytics**: Advanced KDA analysis, CS/min tracking, Vision Score, Damage Share, Gold Efficiency
- **Real-time Data**: Live match analysis and gaming statistics with <5 second response times
- **Multi-region Support**: Global LoL regions (NA, EUW, KR, JP, etc.)
- **Team Management**: Comprehensive team creation, player management, and performance tracking
- **Riot API Integration**: Seamless access to Riot Games data with intelligent rate limit management
- **Subscription Tiers**: Flexible pricing from free tier to enterprise solutions

### Performance Targets

- **⚡ Analytics Speed**: <5 seconds for gaming analytics
- **🚀 UI Response**: <2 seconds for dashboard loading  
- **📊 Uptime**: 99.9% platform availability
- **👥 Scalability**: Support for 1M+ concurrent users
- **🔄 Real-time**: <1 second latency for live data

## 🚀 Quick Start

### 1. Get Your API Key

Sign up at [herald.lol](https://herald.lol) and generate your API key from the developer dashboard.

```bash
curl -X POST https://api.herald.lol/v3/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-password"
  }'
```

### 2. Make Your First Request

```bash
curl -X GET "https://api.herald.lol/v3/gaming/analytics/summoner/NA/HeraldChampion" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Gaming-API-Key: YOUR_API_KEY"
```

### 3. Get Gaming Analytics

```javascript
const response = await fetch('https://api.herald.lol/v3/gaming/analytics/summoner/NA/HeraldChampion', {
  headers: {
    'Authorization': 'Bearer YOUR_JWT_TOKEN',
    'X-Gaming-API-Key': 'YOUR_API_KEY'
  }
});

const analytics = await response.json();
console.log(analytics.analytics.data.kda); // KDA statistics
console.log(analytics.analytics.data.cs_per_minute); // CS/min performance
```

## 🔐 Authentication

Herald.lol supports multiple authentication methods:

### Bearer Token (Recommended)
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### API Key
```http
X-Gaming-API-Key: hld_1234567890abcdef
```

### OAuth 2.0
Supports Google, Discord, Riot Games, and other providers.

## 📊 Rate Limits

Rate limits are enforced per subscription tier:

| Tier | Requests/Minute | Analytics/Minute | Exports/Day |
|------|-----------------|------------------|-------------|
| **Free** | 60 | 30 | 1 |
| **Premium** | 300 | 100 | 5 |
| **Pro** | 1,200 | 500 | 25 |
| **Enterprise** | 6,000 | 2,000 | Unlimited |

### Gaming-Specific Limits

- **Riot API Proxy**: Additional rate limiting to comply with Riot Games API
- **Analytics Endpoints**: Enhanced limits for complex gaming calculations
- **Export Operations**: MFA-protected with tier-based daily limits

### Rate Limit Headers

All API responses include rate limit information:

```http
X-Gaming-Rate-Limit: 60
X-Gaming-Rate-Remaining: 45
X-Gaming-Rate-Reset: 1705312800
X-Gaming-Rate-Tier: premium
Retry-After: 30
```

## 🎯 Gaming Analytics

### Summoner Analytics

Get comprehensive gaming analytics for any League of Legends player:

```bash
GET /gaming/analytics/summoner/{region}/{summonerName}
```

**Response includes:**
- **KDA Analysis**: Kill/Death/Assist ratios and trends
- **CS/min Performance**: Creep score efficiency tracking
- **Vision Score**: Map control and warding metrics
- **Damage Share**: Team fight contribution analysis
- **Gold Efficiency**: Economic performance metrics
- **Gaming Insights**: AI-powered improvement recommendations

### Advanced Analytics Features

- **Time Range Filtering**: 1d, 7d, 30d, 90d, 1y analysis periods
- **Game Mode Filtering**: Ranked, Normal, ARAM, TFT support
- **Champion-Specific**: Per-champion performance breakdowns
- **Trend Analysis**: Performance improvement tracking
- **Comparative Analytics**: Rank and region benchmarking

### Real-time Match Analysis

```bash
GET /gaming/matches/{region}/{matchId}?include_timeline=true
```

Provides detailed match analysis including:
- Participant statistics and builds
- Timeline events and objectives
- Team fight analysis
- Gold and experience curves
- Damage and healing breakdowns

## 🏆 Team Management

### Creating Gaming Teams

```javascript
const team = await fetch('https://api.herald.lol/v3/gaming/teams', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_TOKEN',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: "Herald Champions",
    tag: "HC",
    region: "NA",
    game_type: "LoL",
    tier: "amateur"
  })
});
```

### Team Analytics

- **Team Performance Metrics**: Win rates, KDA, objective control
- **Player Synergy Analysis**: Team composition effectiveness
- **Scrim and Tournament Tracking**: Competitive performance
- **Role Performance**: Position-specific analytics
- **Champion Pool Analysis**: Draft strategy insights

## 🌍 Multi-Region Support

Herald.lol supports all League of Legends regions:

| Region | Code | Description |
|--------|------|-------------|
| North America | `NA` | NA1 server |
| Europe West | `EUW` | EUW1 server |
| Europe Northeast | `EUNE` | EUN1 server |
| Korea | `KR` | KR server |
| Japan | `JP` | JP1 server |
| Brazil | `BR` | BR1 server |
| Latin America North | `LAN` | LA1 server |
| Latin America South | `LAS` | LA2 server |
| Oceania | `OCE` | OC1 server |
| Russia | `RU` | RU server |
| Turkey | `TR` | TR1 server |

## 📈 Data Export

Export gaming analytics in multiple formats:

```javascript
const exportData = await fetch('https://api.herald.lol/v3/gaming/analytics/summoner/NA/Player/export', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_TOKEN',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    format: "json",
    time_range: "30d",
    include_matches: true,
    mfa_token: "123456" // MFA required for exports
  })
});
```

**Supported Formats:**
- **JSON**: Machine-readable gaming data
- **CSV**: Spreadsheet-friendly format
- **XLSX**: Excel workbook with multiple sheets
- **PDF**: Formatted analytics reports

## 🔒 Security

### Multi-Factor Authentication (MFA)

High-value operations require MFA:
- Data exports
- Team management
- Subscription changes
- Account deletion

**MFA Methods:**
- **TOTP**: Time-based one-time passwords (Google Authenticator, Authy)
- **WebAuthn**: Hardware security keys and biometric authentication
- **Backup Codes**: Emergency recovery codes

### Data Protection

- **GDPR Compliant**: EU data protection compliance
- **Riot ToS Compliant**: Full compliance with Riot Games Terms of Service
- **End-to-End Encryption**: AES-256 encryption for sensitive data
- **Audit Logging**: Complete audit trail for all gaming operations

## 🛠️ SDKs and Libraries

### Official SDKs

- **JavaScript/TypeScript**: `npm install herald-lol-sdk`
- **Python**: `pip install herald-lol`
- **Go**: `go get github.com/herald-lol/go-sdk`
- **Java**: Maven and Gradle support
- **C#**: NuGet package available

### Community SDKs

- **PHP**: Community-maintained Laravel package
- **Ruby**: Gem for Ruby on Rails integration
- **Rust**: Community crate for Rust developers

## 📚 Code Examples

### Getting Started with JavaScript

```javascript
import { HeraldClient } from 'herald-lol-sdk';

const client = new HeraldClient({
  apiKey: 'your-api-key',
  token: 'your-jwt-token'
});

// Get summoner analytics
const analytics = await client.gaming.getSummonerAnalytics('NA', 'HeraldChampion');
console.log(analytics.data.kda.ratio); // 3.25

// Create a gaming team
const team = await client.teams.create({
  name: 'Herald Champions',
  tag: 'HC',
  region: 'NA',
  gameType: 'LoL'
});
```

### Python Example

```python
from herald_lol import HeraldClient

client = HeraldClient(
    api_key="your-api-key",
    token="your-jwt-token"
)

# Get gaming analytics
analytics = client.gaming.get_summoner_analytics("NA", "HeraldChampion")
print(f"KDA Ratio: {analytics.data.kda.ratio}")

# Export analytics data
export = client.gaming.export_summoner_analytics(
    region="NA",
    summoner_name="HeraldChampion",
    format="json",
    time_range="30d",
    mfa_token="123456"
)
```

### Go Example

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/herald-lol/go-sdk/gaming"
)

func main() {
    client := gaming.NewClient("your-api-key", "your-jwt-token")
    
    analytics, err := client.GetSummonerAnalytics(context.Background(), "NA", "HeraldChampion")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("CS/min: %.2f\n", analytics.Data.CSPerMinute)
}
```

## 🚨 Error Handling

Herald.lol uses standard HTTP status codes and provides detailed error information:

### Error Response Format

```json
{
  "error": "Gaming data not found",
  "details": "Summoner not found in specified region",
  "code": "SUMMONER_NOT_FOUND",
  "gaming_platform": "herald-lol",
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_abc123"
}
```

### Common Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `INVALID_PARAMETERS` | Invalid request parameters |
| 401 | `AUTHENTICATION_REQUIRED` | Authentication token required |
| 403 | `INSUFFICIENT_PERMISSIONS` | Insufficient permissions for operation |
| 404 | `RESOURCE_NOT_FOUND` | Gaming data not found |
| 429 | `RATE_LIMIT_EXCEEDED` | API rate limit exceeded |
| 500 | `INTERNAL_ERROR` | Internal server error |

### Gaming-Specific Errors

- `SUMMONER_NOT_FOUND`: Summoner doesn't exist in region
- `MATCH_NOT_FOUND`: Match data not available
- `REGION_UNAVAILABLE`: Region temporarily unavailable
- `RIOT_API_UNAVAILABLE`: Riot Games API is down
- `MFA_REQUIRED`: Multi-factor authentication required
- `SUBSCRIPTION_REQUIRED`: Premium subscription required

## 📖 API Reference

### Complete API Documentation

- **OpenAPI 3.0 Spec**: [openapi.yaml](./openapi.yaml)
- **Interactive Docs**: [https://api.herald.lol/docs](https://api.herald.lol/docs)
- **Postman Collection**: [Download Collection](https://api.herald.lol/postman)

### Endpoint Categories

1. **System** - Health checks and status
2. **Authentication** - User login and registration
3. **Gaming Analytics** - Core analytics functionality
4. **Match Data** - Match details and statistics
5. **Team Management** - Team creation and management
6. **Account** - User account and subscription management

## 🔄 Webhooks

Subscribe to real-time gaming events:

```javascript
// Register webhook for match updates
await client.webhooks.create({
  url: 'https://your-app.com/webhooks/matches',
  events: ['match.completed', 'summoner.rank_changed'],
  summoner_names: ['HeraldChampion'],
  regions: ['NA']
});
```

**Available Events:**
- `match.completed` - New match data available
- `summoner.rank_changed` - Rank promotion/demotion
- `team.match_completed` - Team match results
- `analytics.updated` - Analytics recalculated

## 📊 Status and Monitoring

### Service Status

- **Status Page**: [https://status.herald.lol](https://status.herald.lol)
- **API Health**: `GET /health`
- **Service Metrics**: Real-time performance monitoring

### Performance Metrics

Monitor API performance in real-time:
- Response times per endpoint
- Rate limit utilization
- Error rates by region
- Gaming analytics processing times

## 💼 Enterprise Features

### Dedicated Infrastructure

- **Dedicated API Endpoints**: Custom API domains
- **Enhanced Rate Limits**: Unlimited requests
- **Priority Support**: 24/7 dedicated support
- **Custom Integrations**: White-label solutions

### Advanced Analytics

- **Custom Metrics**: Define your own gaming KPIs
- **Data Warehousing**: Historical data export
- **Real-time Streaming**: Live match data feeds
- **Machine Learning**: Custom AI models

## 🤝 Support

### Getting Help

- **Documentation**: [https://docs.herald.lol](https://docs.herald.lol)
- **API Support**: [api-support@herald.lol](mailto:api-support@herald.lol)
- **Discord Community**: [https://discord.gg/herald-lol](https://discord.gg/herald-lol)
- **GitHub Issues**: [https://github.com/herald-lol/api-issues](https://github.com/herald-lol/api-issues)

### SLA Commitments

- **99.9% Uptime**: Production API availability
- **<5s Response**: Gaming analytics endpoints
- **<2s Response**: Standard API endpoints
- **24/7 Monitoring**: Continuous service monitoring

## 🚀 Roadmap

### Upcoming Features

- **TFT Analytics**: Comprehensive Teamfight Tactics support
- **Mobile SDK**: Native iOS and Android libraries
- **GraphQL API**: Alternative to REST endpoints
- **Real-time Analytics**: Live match analysis
- **Tournament Mode**: Competitive tournament tracking
- **Coach Dashboard**: Advanced coaching tools

### Gaming Integrations

- **League Client**: Direct League of Legends client integration
- **Overwolf**: Gaming overlay applications
- **Blitz.gg**: Third-party tool integrations
- **OP.GG**: Analytics platform partnerships

---

## 📜 Legal

- **Terms of Service**: [https://herald.lol/terms](https://herald.lol/terms)
- **Privacy Policy**: [https://herald.lol/privacy](https://herald.lol/privacy)
- **Riot Games ToS**: Full compliance with Riot Games Terms of Service
- **Data Processing**: GDPR compliant data handling

**© 2024 Herald.lol Gaming Analytics Platform. All rights reserved.**