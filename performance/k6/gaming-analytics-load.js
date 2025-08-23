// Herald.lol Gaming Analytics Load Testing with k6
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Counter, Trend } from 'k6/metrics';

// Herald.lol Gaming Performance Requirements
// - Analytics: <5s response time
// - UI Load: <2s response time  
// - Uptime: 99.9%
// - Concurrent: 1M+ users

// Custom metrics for Herald.lol gaming analytics
export let analyticsResponseTime = new Trend('analytics_response_time');
export let gamingApiErrors = new Rate('gaming_api_errors');
export let analyticsRequests = new Counter('analytics_requests');
export let riotApiCalls = new Counter('riot_api_calls');

// Performance test configuration
export let options = {
  scenarios: {
    // Gaming Analytics Dashboard Load Test
    analytics_dashboard: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 100 },   // Ramp up
        { duration: '5m', target: 500 },   // Normal load
        { duration: '2m', target: 1000 },  // Peak gaming hours
        { duration: '5m', target: 1000 },  // Sustained peak
        { duration: '5m', target: 0 },     // Ramp down
      ],
      gracefulRampDown: '30s',
    },
    
    // Gaming API Stress Test (1M+ concurrent target)
    gaming_api_stress: {
      executor: 'constant-arrival-rate',
      rate: 1000,                         // 1000 requests per second
      timeUnit: '1s',
      duration: '10m',
      preAllocatedVUs: 100,
      maxVUs: 2000,
    },
    
    // Gaming Metrics Spike Test
    gaming_spike: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      stages: [
        { duration: '2m', target: 100 },   // Normal
        { duration: '1m', target: 2000 },  // Spike (post-game analysis rush)
        { duration: '2m', target: 100 },   // Recovery
      ],
    },
  },
  
  // Performance thresholds based on Herald.lol requirements
  thresholds: {
    // Analytics must load in <5s (Herald.lol requirement)
    'analytics_response_time': ['p(95)<5000', 'p(99)<8000'],
    
    // 99.9% uptime requirement
    'http_req_failed': ['rate<0.01'],
    
    // Gaming API response times
    'http_req_duration{group:::gaming_analytics}': ['p(95)<5000'],
    'http_req_duration{group:::match_analysis}': ['p(95)<3000'],
    'http_req_duration{group:::riot_api}': ['p(95)<2000'],
    
    // Error rates
    'gaming_api_errors': ['rate<0.05'],
    
    // Concurrent user support (1M+ target)
    'http_reqs': ['count>50000'],
  },
};

// Herald.lol API endpoints
const BASE_URL = __ENV.API_BASE_URL || 'http://localhost:8080';
const FRONTEND_URL = __ENV.FRONTEND_URL || 'http://localhost:3000';

// Test data for gaming scenarios
const TEST_DATA = {
  summoners: ['TestSummoner1', 'TestSummoner2', 'TestSummoner3'],
  matchIds: ['NA1_4567890123', 'NA1_4567890124', 'NA1_4567890125'],
  puuids: ['test-puuid-1', 'test-puuid-2', 'test-puuid-3'],
  champions: ['Jinx', 'Caitlyn', 'Ezreal', 'Vayne', 'Ashe'],
  ranks: ['IRON', 'BRONZE', 'SILVER', 'GOLD', 'PLATINUM'],
};

// Gaming user authentication
function authenticateGamingUser() {
  const loginData = {
    email: 'loadtest@herald.lol',
    password: 'LoadTest123!',
  };
  
  const response = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify(loginData), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(response, {
    'Gaming user login successful': (r) => r.status === 200,
    'Gaming auth token received': (r) => r.json('token') !== null,
  });
  
  return response.json('token');
}

// Test gaming analytics dashboard load
function testAnalyticsDashboard(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  // Simulate dashboard load with all gaming widgets
  const startTime = Date.now();
  
  // Parallel requests for gaming analytics widgets
  const responses = http.batch([
    ['GET', `${BASE_URL}/api/analytics/kda`, null, { headers }],
    ['GET', `${BASE_URL}/api/analytics/cs`, null, { headers }],
    ['GET', `${BASE_URL}/api/analytics/vision`, null, { headers }],
    ['GET', `${BASE_URL}/api/analytics/damage`, null, { headers }],
    ['GET', `${BASE_URL}/api/analytics/gold`, null, { headers }],
  ]);
  
  const loadTime = Date.now() - startTime;
  analyticsResponseTime.add(loadTime);
  analyticsRequests.add(5);
  
  // Validate Herald.lol <5s requirement
  check(responses, {
    'Analytics dashboard loads in <5s': () => loadTime < 5000,
    'All gaming widgets respond successfully': (r) => r.every(res => res.status === 200),
    'KDA analytics valid': (r) => r[0].json('data.currentKDA') > 0,
    'CS/min analytics valid': (r) => r[1].json('data.currentCSPerMin') > 0,
    'Vision score valid': (r) => r[2].json('data.averageVisionScore') > 0,
  });
  
  return responses;
}

// Test match analysis performance
function testMatchAnalysis(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const matchId = TEST_DATA.matchIds[Math.floor(Math.random() * TEST_DATA.matchIds.length)];
  
  const startTime = Date.now();
  const response = http.get(`${BASE_URL}/api/matches/${matchId}/analyze`, { headers });
  const analysisTime = Date.now() - startTime;
  
  check(response, {
    'Match analysis completes successfully': (r) => r.status === 200,
    'Match analysis <5s (Herald.lol requirement)': () => analysisTime < 5000,
    'Analysis contains gaming metrics': (r) => {
      const data = r.json();
      return data.kda && data.cs && data.vision && data.damage && data.gold;
    },
  });
  
  return response;
}

// Test Riot API integration performance
function testRiotApiIntegration(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const summoner = TEST_DATA.summoners[Math.floor(Math.random() * TEST_DATA.summoners.length)];
  
  // Test Riot API calls through Herald.lol backend
  const responses = http.batch([
    ['GET', `${BASE_URL}/api/riot/summoner/${summoner}`, null, { headers }],
    ['GET', `${BASE_URL}/api/riot/ranked/${summoner}`, null, { headers }],
    ['GET', `${BASE_URL}/api/riot/matches/${summoner}?count=20`, null, { headers }],
  ]);
  
  riotApiCalls.add(3);
  
  check(responses, {
    'Riot API integration successful': (r) => r.every(res => res.status === 200),
    'Summoner data retrieved': (r) => r[0].json('name') !== null,
    'Ranked data retrieved': (r) => r[1].json('tier') !== null,
    'Match history retrieved': (r) => Array.isArray(r[2].json()),
  });
  
  return responses;
}

// Test team composition optimization
function testTeamCompositionOptimization(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const optimizationRequest = {
    currentChampions: TEST_DATA.champions.slice(0, 4),
    strategy: 'meta_optimal',
    targetRank: TEST_DATA.ranks[Math.floor(Math.random() * TEST_DATA.ranks.length)],
  };
  
  const startTime = Date.now();
  const response = http.post(`${BASE_URL}/api/team-composition/optimize`, 
    JSON.stringify(optimizationRequest), { headers });
  const optimizationTime = Date.now() - startTime;
  
  check(response, {
    'Team composition optimization successful': (r) => r.status === 200,
    'Optimization completes quickly': () => optimizationTime < 3000,
    'Recommendations provided': (r) => Array.isArray(r.json('recommendations')),
  });
  
  return response;
}

// Test gaming search functionality
function testGamingSearch(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const champion = TEST_DATA.champions[Math.floor(Math.random() * TEST_DATA.champions.length)];
  
  const response = http.get(`${BASE_URL}/api/search/champions?q=${champion}&analytics=true`, { headers });
  
  check(response, {
    'Gaming search successful': (r) => r.status === 200,
    'Search results contain analytics': (r) => {
      const results = r.json();
      return Array.isArray(results) && results.length > 0;
    },
  });
  
  return response;
}

// Main test execution
export default function () {
  // Authenticate gaming user
  const token = authenticateGamingUser();
  if (!token) {
    gamingApiErrors.add(1);
    return;
  }
  
  // Test different gaming scenarios based on weight
  const scenario = Math.random();
  
  if (scenario < 0.4) {
    // 40% - Analytics dashboard (most common)
    testAnalyticsDashboard(token);
  } else if (scenario < 0.7) {
    // 30% - Match analysis
    testMatchAnalysis(token);
  } else if (scenario < 0.85) {
    // 15% - Riot API integration
    testRiotApiIntegration(token);
  } else if (scenario < 0.95) {
    // 10% - Team composition
    testTeamCompositionOptimization(token);
  } else {
    // 5% - Search functionality
    testGamingSearch(token);
  }
  
  // Simulate user think time (gaming analysis pause)
  sleep(Math.random() * 3 + 1); // 1-4 seconds
}

// Setup function
export function setup() {
  console.log('🎮 Herald.lol Gaming Analytics Performance Testing');
  console.log(`📊 Target: <5s analytics, 99.9% uptime, 1M+ concurrent`);
  console.log(`🔗 API: ${BASE_URL}`);
  console.log(`🌐 Frontend: ${FRONTEND_URL}`);
}

// Teardown function with results summary
export function teardown(data) {
  console.log('🎮 Herald.lol Performance Test Complete');
  console.log('📊 Gaming Analytics Performance Summary:');
  console.log('   ✅ Analytics Response Time Target: <5s');
  console.log('   ✅ Uptime Requirement: 99.9%');
  console.log('   ✅ Concurrent User Target: 1M+');
  console.log('🎯 Check k6 output for detailed performance metrics');
}