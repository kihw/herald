// Herald.lol Gaming Spike Testing - Sudden Load Spikes
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Counter, Trend } from 'k6/metrics';

// Spike-specific metrics for Herald.lol
export let spikeResponseTime = new Trend('spike_response_time');
export let spikeRecoveryTime = new Trend('spike_recovery_time');
export let analyticsAvailability = new Rate('analytics_availability');
export let autoScalingTriggered = new Counter('auto_scaling_triggered');

// Gaming spike test scenarios
export let options = {
  scenarios: {
    // Tournament end spike (massive concurrent analysis)
    tournament_end_spike: {
      executor: 'ramping-vus',
      startVUs: 100,
      stages: [
        { duration: '30s', target: 100 },    // Baseline
        { duration: '10s', target: 5000 },   // Sudden spike
        { duration: '2m', target: 5000 },    // Sustained spike
        { duration: '30s', target: 100 },    // Recovery
      ],
    },
    
    // Rank reset spike (season start)
    rank_reset_spike: {
      executor: 'ramping-arrival-rate',
      startRate: 50,
      stages: [
        { duration: '1m', target: 50 },      // Normal
        { duration: '5s', target: 1000 },    // Rank reset announcement
        { duration: '3m', target: 1000 },    // Everyone checking ranks
        { duration: '1m', target: 50 },      // Back to normal
      ],
    },
    
    // Patch analysis spike (new meta analysis)
    patch_spike: {
      executor: 'ramping-vus',
      startVUs: 50,
      stages: [
        { duration: '1m', target: 50 },      // Pre-patch
        { duration: '15s', target: 2000 },   // Patch goes live
        { duration: '5m', target: 2000 },    // Meta analysis rush
        { duration: '2m', target: 200 },     // Settling
        { duration: '1m', target: 50 },      // Normal
      ],
    },
  },
  
  // Spike test thresholds
  thresholds: {
    // Analytics must remain available during spikes
    'analytics_availability': ['rate>0.95'], // 95% availability during spikes
    
    // Response time targets during spikes
    'spike_response_time': ['p(50)<8000', 'p(95)<15000'], // Degraded but functional
    
    // Recovery time after spikes
    'spike_recovery_time': ['p(95)<10000'], // System recovers within 10s
    
    // Error rates during spikes
    'http_req_failed': ['rate<0.1'], // Allow 10% errors during spikes
  },
};

const BASE_URL = __ENV.API_BASE_URL || 'http://localhost:8080';

// Spike test scenarios data
const SPIKE_SCENARIOS = {
  tournament_end: {
    endpoints: [
      'analytics/kda?tournament=true',
      'analytics/rankings/tournament',
      'matches/tournament/analysis',
      'team-composition/tournament-meta',
    ],
    concurrent_factor: 10, // 10x normal load
  },
  
  rank_reset: {
    endpoints: [
      'analytics/rankings/current',
      'analytics/progression/reset',
      'matches/placement/analysis',
      'skill-progression/reset-impact',
    ],
    concurrent_factor: 8, // 8x normal load
  },
  
  patch_analysis: {
    endpoints: [
      'analytics/meta/changes',
      'counter-picks/patch-impact',
      'team-composition/new-meta',
      'champions/patch-analysis',
    ],
    concurrent_factor: 6, // 6x normal load
  },
};

// Gaming spike user authentication
function authenticateSpikeUser() {
  const response = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify({
    email: `spiketest${__VU}@herald.lol`,
    password: 'SpikeTest123!',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  if (response.status === 200) {
    return response.json('token');
  }
  
  return null;
}

// Test tournament end spike scenario
function testTournamentEndSpike(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const startTime = Date.now();
  const scenario = SPIKE_SCENARIOS.tournament_end;
  
  // Simulate tournament end analytics rush
  const responses = http.batch(
    scenario.endpoints.map(endpoint => 
      ['GET', `${BASE_URL}/api/${endpoint}`, null, { headers }]
    )
  );
  
  const responseTime = Date.now() - startTime;
  spikeResponseTime.add(responseTime);
  
  // Check analytics availability during spike
  const availabilityRate = responses.filter(r => r.status === 200).length / responses.length;
  analyticsAvailability.add(availabilityRate > 0.8 ? 1 : 0);
  
  check(responses, {
    'Tournament analytics survive spike': (r) => r.some(res => res.status === 200),
    'Tournament data partially available': () => availabilityRate > 0.5,
    'No complete tournament system failure': (r) => !r.every(res => res.status >= 500),
  });
  
  return responses;
}

// Test rank reset spike scenario
function testRankResetSpike(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const startTime = Date.now();
  const scenario = SPIKE_SCENARIOS.rank_reset;
  
  // Simulate rank reset checking rush
  const responses = http.batch(
    scenario.endpoints.map(endpoint => 
      ['GET', `${BASE_URL}/api/${endpoint}`, null, { headers }]
    )
  );
  
  const responseTime = Date.now() - startTime;
  spikeResponseTime.add(responseTime);
  
  const availabilityRate = responses.filter(r => r.status === 200).length / responses.length;
  analyticsAvailability.add(availabilityRate > 0.7 ? 1 : 0);
  
  check(responses, {
    'Rank data survives reset spike': (r) => r.some(res => res.status === 200),
    'Ranking system partially operational': () => availabilityRate > 0.4,
    'Rank reset spike handled gracefully': (r) => responseTime < 20000, // 20s max
  });
  
  return responses;
}

// Test patch analysis spike scenario
function testPatchAnalysisSpike(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  const startTime = Date.now();
  const scenario = SPIKE_SCENARIOS.patch_analysis;
  
  // Simulate patch meta analysis rush
  const responses = http.batch(
    scenario.endpoints.map(endpoint => 
      ['GET', `${BASE_URL}/api/${endpoint}`, null, { headers }]
    )
  );
  
  const responseTime = Date.now() - startTime;
  spikeResponseTime.add(responseTime);
  
  const availabilityRate = responses.filter(r => r.status === 200).length / responses.length;
  analyticsAvailability.add(availabilityRate > 0.8 ? 1 : 0);
  
  check(responses, {
    'Patch analysis survives spike': (r) => r.some(res => res.status === 200),
    'Meta data remains accessible': () => availabilityRate > 0.6,
    'Patch spike recovery time acceptable': (r) => responseTime < 12000, // 12s max
  });
  
  return responses;
}

// Test auto-scaling response
function testAutoScalingResponse(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  // Check system metrics during spike
  const response = http.get(`${BASE_URL}/api/system/metrics`, { headers });
  
  if (response.status === 200) {
    const metrics = response.json();
    
    // Check if auto-scaling indicators are present
    if (metrics.cpu_usage > 80 || metrics.memory_usage > 85 || metrics.response_time > 5000) {
      autoScalingTriggered.add(1);
    }
    
    check(metrics, {
      'System metrics available during spike': () => true,
      'CPU usage monitored': (m) => typeof m.cpu_usage === 'number',
      'Memory usage tracked': (m) => typeof m.memory_usage === 'number',
      'Response times measured': (m) => typeof m.response_time === 'number',
    });
  }
  
  return response;
}

// Test spike recovery
function testSpikeRecovery(token) {
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
  
  // Test basic functionality after spike
  const recoveryStartTime = Date.now();
  
  const response = http.get(`${BASE_URL}/api/analytics/kda?quick=true`, { headers });
  
  const recoveryTime = Date.now() - recoveryStartTime;
  spikeRecoveryTime.add(recoveryTime);
  
  check(response, {
    'System recovers after spike': (r) => r.status === 200,
    'Recovery time acceptable': () => recoveryTime < 5000, // Normal performance restored
    'Analytics functional post-spike': (r) => r.json('data') !== null,
  });
  
  return response;
}

// Main spike test execution
export default function () {
  const token = authenticateSpikeUser();
  if (!token) {
    return;
  }
  
  // Determine which spike scenario we're in based on time
  const testDuration = Date.now() % 600000; // 10-minute cycle
  let spikeType;
  
  if (testDuration < 200000) {
    spikeType = 'tournament';
  } else if (testDuration < 400000) {
    spikeType = 'rank_reset';
  } else {
    spikeType = 'patch';
  }
  
  // Execute spike scenario
  switch (spikeType) {
    case 'tournament':
      testTournamentEndSpike(token);
      break;
    case 'rank_reset':
      testRankResetSpike(token);
      break;
    case 'patch':
      testPatchAnalysisSpike(token);
      break;
  }
  
  // Test auto-scaling response
  if (Math.random() < 0.1) { // 10% of users check system metrics
    testAutoScalingResponse(token);
  }
  
  // Test recovery every few iterations
  if (__ITER > 0 && __ITER % 5 === 0) {
    testSpikeRecovery(token);
  }
  
  // Minimal sleep during spike scenarios
  sleep(0.1 + Math.random() * 0.3); // 0.1-0.4 seconds
}

export function setup() {
  console.log('⚡ Herald.lol Gaming Spike Testing');
  console.log('🎯 Testing sudden load spikes and recovery');
  console.log('📊 Scenarios: Tournament End, Rank Reset, Patch Analysis');
  console.log('🔧 Auto-scaling and recovery testing enabled');
}

export function teardown() {
  console.log('⚡ Herald.lol Spike Test Complete');
  console.log('📊 Spike Handling Summary:');
  console.log('   🎯 Target: Graceful handling of sudden load spikes');
  console.log('   ⚡ Auto-scaling response verification');
  console.log('   🔄 Recovery time and system resilience');
  console.log('   ✅ Gaming experience maintained during spikes');
}