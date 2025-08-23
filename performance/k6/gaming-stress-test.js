// Herald.lol Gaming Stress Testing - Extreme Load Scenarios
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Counter, Trend, Gauge } from 'k6/metrics';

// Custom metrics for extreme gaming load
export let extremeLoadResponseTime = new Trend('extreme_load_response_time');
export let concurrentGamingUsers = new Gauge('concurrent_gaming_users');
export let gamingSystemFailures = new Rate('gaming_system_failures');
export let riotApiRateLimit = new Counter('riot_api_rate_limit_hits');

// Extreme stress test configuration for Herald.lol
export let options = {
  scenarios: {
    // Post-game analysis rush (extreme concurrent load)
    post_game_rush: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 1000 },   // Rapid ramp-up
        { duration: '2m', target: 5000 },   // Major tournament end
        { duration: '3m', target: 10000 },  // Peak analysis load
        { duration: '2m', target: 15000 },  // Extreme stress
        { duration: '2m', target: 5000 },   // Recovery
        { duration: '2m', target: 0 },      // Cool down
      ],
    },
    
    // Gaming API endurance test
    gaming_api_endurance: {
      executor: 'constant-vus',
      vus: 2000,
      duration: '30m',
    },
    
    // Analytics burst load (ranking updates)
    analytics_burst: {
      executor: 'ramping-arrival-rate',
      startRate: 100,
      stages: [
        { duration: '30s', target: 500 },   // Warm up
        { duration: '1m', target: 2000 },   // High load
        { duration: '30s', target: 5000 },  // Burst
        { duration: '1m', target: 1000 },   // Settle
        { duration: '30s', target: 100 },   // Cool down
      ],
    },
  },
  
  // Extreme stress thresholds
  thresholds: {
    // Must maintain <5s even under extreme load
    'extreme_load_response_time': ['p(50)<5000', 'p(95)<10000'],
    
    // System should not fail completely
    'gaming_system_failures': ['rate<0.1'], // Allow 10% degradation under extreme load
    
    // Response times under stress
    'http_req_duration': ['p(99)<15000'], // 15s max under extreme conditions
    
    // Error rates
    'http_req_failed': ['rate<0.05'],
  },
};

const BASE_URL = __ENV.API_BASE_URL || 'http://localhost:8080';

// Heavy gaming test data
const STRESS_TEST_DATA = {
  heavyAnalyticsRequests: [
    'analytics/kda/detailed?matches=100',
    'analytics/cs/timeline?duration=30d',
    'analytics/vision/heatmap?resolution=high',
    'analytics/damage/breakdown?detailed=true',
    'analytics/gold/efficiency?deep_analysis=true',
  ],
  
  complexQueries: [
    'team-composition/optimize?strategy=genetic&iterations=1000',
    'counter-picks/analyze?depth=5&meta_analysis=true',
    'skill-progression/predict?timeframe=1y',
    'coaching/insights/generate?comprehensive=true',
  ],
  
  batchOperations: [
    'matches/batch-analyze?count=50',
    'analytics/multi-metric?metrics=all',
    'riot/bulk-update?summoners=100',
  ],
};

// Stress test user pool
const STRESS_USERS = Array.from({ length: 100 }, (_, i) => ({
  email: `stresstest${i}@herald.lol`,
  password: 'StressTest123!',
  token: null,
}));

// Authenticate stress test user
function authenticateStressUser() {
  const user = STRESS_USERS[Math.floor(Math.random() * STRESS_USERS.length)];
  
  if (user.token) {
    return user.token; // Reuse token
  }
  
  const response = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify({
    email: user.email,
    password: user.password,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  if (response.status === 200) {
    user.token = response.json('token');
    return user.token;
  }
  
  gamingSystemFailures.add(1);
  return null;
}

// Heavy analytics load test
function runHeavyAnalyticsLoad(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const startTime = Date.now();
  
  // Parallel heavy analytics requests
  const heavyRequests = STRESS_TEST_DATA.heavyAnalyticsRequests.map(endpoint => 
    ['GET', `${BASE_URL}/api/${endpoint}`, null, { headers }]
  );
  
  const responses = http.batch(heavyRequests);
  const totalTime = Date.now() - startTime;
  
  extremeLoadResponseTime.add(totalTime);
  concurrentGamingUsers.add(__VU);
  
  check(responses, {
    'Heavy analytics survive stress load': (r) => r.every(res => res.status < 500),
    'At least 50% succeed under stress': (r) => {
      const successCount = r.filter(res => res.status === 200).length;
      return (successCount / r.length) >= 0.5;
    },
  });
  
  return responses;
}

// Complex query stress test
function runComplexQueryStress(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const complexQuery = STRESS_TEST_DATA.complexQueries[
    Math.floor(Math.random() * STRESS_TEST_DATA.complexQueries.length)
  ];
  
  const startTime = Date.now();
  const response = http.post(`${BASE_URL}/api/${complexQuery}`, '{}', { headers });
  const queryTime = Date.now() - startTime;
  
  extremeLoadResponseTime.add(queryTime);
  
  const success = check(response, {
    'Complex query handles stress': (r) => r.status < 500,
    'Query completes within reasonable time': () => queryTime < 30000, // 30s max
  });
  
  if (!success) {
    gamingSystemFailures.add(1);
  }
  
  return response;
}

// Batch operation stress test
function runBatchOperationStress(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const batchOperation = STRESS_TEST_DATA.batchOperations[
    Math.floor(Math.random() * STRESS_TEST_DATA.batchOperations.length)
  ];
  
  const startTime = Date.now();
  const response = http.post(`${BASE_URL}/api/${batchOperation}`, '{}', { headers });
  const batchTime = Date.now() - startTime;
  
  extremeLoadResponseTime.add(batchTime);
  
  check(response, {
    'Batch operation survives stress': (r) => r.status < 500,
    'Batch operation reasonable time': () => batchTime < 60000, // 1 minute max
  });
  
  return response;
}

// Database stress test
function runDatabaseStress(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  // Simulate database-heavy operations
  const dbStressRequests = [
    `analytics/historical/query?depth=365d&metrics=all`,
    `matches/search/advanced?complex_filter=true`,
    `rankings/leaderboard/global?limit=10000`,
  ];
  
  const responses = http.batch(
    dbStressRequests.map(endpoint => 
      ['GET', `${BASE_URL}/api/${endpoint}`, null, { headers }]
    )
  );
  
  check(responses, {
    'Database survives heavy load': (r) => r.some(res => res.status === 200),
    'No complete database failure': (r) => !r.every(res => res.status >= 500),
  });
  
  return responses;
}

// Riot API rate limit stress
function runRiotApiStress(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  // Rapid Riot API calls to test rate limiting
  const riotRequests = Array.from({ length: 10 }, (_, i) => 
    ['GET', `${BASE_URL}/api/riot/summoner/stresstest${i}`, null, { headers }]
  );
  
  const responses = http.batch(riotRequests);
  
  // Count rate limit hits
  const rateLimitHits = responses.filter(r => r.status === 429).length;
  riotApiRateLimit.add(rateLimitHits);
  
  check(responses, {
    'Riot API rate limiting works': (r) => r.some(res => res.status === 200),
    'Rate limits properly enforced': () => rateLimitHits > 0 && rateLimitHits < 10,
  });
  
  return responses;
}

// Memory stress test
function runMemoryStress(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  // Request large datasets to stress memory
  const largeDataRequests = [
    'analytics/export/full?format=detailed',
    'matches/history/complete?summoner=test&limit=1000',
    'statistics/comprehensive?timeframe=all',
  ];
  
  const responses = http.batch(
    largeDataRequests.map(endpoint => 
      ['GET', `${BASE_URL}/api/${endpoint}`, null, { headers }]
    )
  );
  
  check(responses, {
    'Memory-intensive requests handled': (r) => r.some(res => res.status === 200),
    'System does not crash under memory pressure': (r) => !r.every(res => res.status >= 500),
  });
  
  return responses;
}

// Main stress test execution
export default function () {
  const token = authenticateStressUser();
  if (!token) {
    return;
  }
  
  // Update concurrent user count
  concurrentGamingUsers.add(__VU);
  
  // Execute different stress scenarios
  const stressScenario = Math.random();
  
  if (stressScenario < 0.3) {
    // 30% - Heavy analytics load
    runHeavyAnalyticsLoad(token);
  } else if (stressScenario < 0.5) {
    // 20% - Complex queries
    runComplexQueryStress(token);
  } else if (stressScenario < 0.65) {
    // 15% - Batch operations
    runBatchOperationStress(token);
  } else if (stressScenario < 0.8) {
    // 15% - Database stress
    runDatabaseStress(token);
  } else if (stressScenario < 0.9) {
    // 10% - Riot API stress
    runRiotApiStress(token);
  } else {
    // 10% - Memory stress
    runMemoryStress(token);
  }
  
  // Minimal sleep under stress conditions
  sleep(0.1 + Math.random() * 0.5); // 0.1-0.6 seconds
}

export function setup() {
  console.log('🔥 Herald.lol EXTREME STRESS TESTING');
  console.log('⚠️  WARNING: This test applies extreme load conditions');
  console.log('🎯 Testing system limits and failure modes');
  console.log('📊 Expected some degradation under extreme conditions');
}

export function teardown() {
  console.log('🔥 Herald.lol Extreme Stress Test Complete');
  console.log('📊 System Resilience Summary:');
  console.log('   🎯 Target: Graceful degradation under extreme load');
  console.log('   ✅ No complete system failures allowed');
  console.log('   ⚡ Maintained partial functionality under stress');
}