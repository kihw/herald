/**
 * Performance and Load Testing for LoL Match Exporter
 * 
 * This file contains comprehensive performance tests using k6 to validate:
 * - API endpoint performance under load
 * - Real-time notification system scalability
 * - Database performance with concurrent users
 * - Memory and CPU usage under stress
 */

import http from 'k6/http';
import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const responseTime = new Trend('response_time');
const wsConnections = new Counter('websocket_connections');
const notificationsReceived = new Counter('notifications_received');

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8001';
const WS_URL = __ENV.WS_URL || 'ws://localhost:8001';

// Test scenarios configuration
export const options = {
  scenarios: {
    // Basic API load test
    api_load: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '2m', target: 10 },  // Ramp up to 10 users
        { duration: '5m', target: 10 },  // Stay at 10 users
        { duration: '2m', target: 20 },  // Ramp up to 20 users
        { duration: '5m', target: 20 },  // Stay at 20 users
        { duration: '2m', target: 0 },   // Ramp down
      ],
      gracefulRampDown: '30s',
      env: { SCENARIO: 'api_load' },
    },
    
    // Stress test for breaking point
    stress_test: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '2m', target: 50 },  // Quick ramp to 50 users
        { duration: '5m', target: 100 }, // Ramp to 100 users
        { duration: '2m', target: 200 }, // Stress with 200 users
        { duration: '5m', target: 200 }, // Maintain stress
        { duration: '3m', target: 0 },   // Recovery
      ],
      env: { SCENARIO: 'stress' },
    },
    
    // Spike test for sudden load
    spike_test: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '30s', target: 5 },   // Normal load
        { duration: '10s', target: 100 }, // Sudden spike
        { duration: '1m', target: 100 },  // Maintain spike
        { duration: '10s', target: 5 },   // Drop back
        { duration: '30s', target: 5 },   // Normal load
      ],
      env: { SCENARIO: 'spike' },
    },
    
    // Real-time notification test
    realtime_notifications: {
      executor: 'constant-vus',
      vus: 20,
      duration: '5m',
      env: { SCENARIO: 'realtime' },
    },
    
    // Soak test for stability
    soak_test: {
      executor: 'constant-vus',
      vus: 10,
      duration: '30m',
      env: { SCENARIO: 'soak' },
    }
  },
  
  thresholds: {
    // API response time thresholds
    'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% under 500ms, 99% under 1s
    'http_req_failed': ['rate<0.05'], // Error rate under 5%
    
    // Custom metric thresholds
    'response_time': ['p(95)<300'],
    'errors': ['rate<0.01'],
    
    // WebSocket thresholds
    'ws_connecting': ['p(95)<200'],
    'ws_msgs_received': ['count>0'],
  }
};

// Test data
const testUsers = [
  { riot_id: 'TestUser1', riot_tag: 'EUW', region: 'euw1' },
  { riot_id: 'TestUser2', riot_tag: 'NA', region: 'na1' },
  { riot_id: 'TestUser3', riot_tag: 'KR', region: 'kr' },
  { riot_id: 'TestUser4', riot_tag: 'JP', region: 'jp1' },
  { riot_id: 'TestUser5', riot_tag: 'BR', region: 'br1' },
];

// Authentication helper
function authenticate() {
  const user = testUsers[Math.floor(Math.random() * testUsers.length)];
  
  const authPayload = JSON.stringify({
    riot_id: user.riot_id,
    riot_tag: user.riot_tag,
    region: user.region
  });
  
  const authResponse = http.post(`${BASE_URL}/api/auth/validate`, authPayload, {
    headers: { 'Content-Type': 'application/json' },
    tags: { name: 'auth' }
  });
  
  check(authResponse, {
    'authentication successful': (r) => r.status === 200,
    'auth response time ok': (r) => r.timings.duration < 1000,
  });
  
  if (authResponse.status === 200) {
    const cookies = http.cookieJar();
    return cookies;
  }
  
  errorRate.add(1);
  return null;
}

// Main test function
export default function () {
  const scenario = __ENV.SCENARIO || 'api_load';
  
  switch (scenario) {
    case 'api_load':
    case 'stress':
    case 'spike':
    case 'soak':
      testAPIEndpoints();
      break;
    case 'realtime':
      testRealtimeNotifications();
      break;
    default:
      testAPIEndpoints();
  }
}

function testAPIEndpoints() {
  // Authenticate
  const cookies = authenticate();
  if (!cookies) return;
  
  // Test health endpoint
  const healthStart = Date.now();
  const healthResponse = http.get(`${BASE_URL}/api/health`, {
    tags: { name: 'health' }
  });
  responseTime.add(Date.now() - healthStart);
  
  check(healthResponse, {
    'health check ok': (r) => r.status === 200,
    'health response has service': (r) => r.json('service') !== undefined,
  });
  
  sleep(0.5);
  
  // Test analytics endpoints
  const analyticsEndpoints = [
    '/api/analytics/health',
    '/api/analytics/period/week',
    '/api/analytics/mmr',
    '/api/analytics/recommendations',
  ];
  
  analyticsEndpoints.forEach(endpoint => {
    const start = Date.now();
    const response = http.get(`${BASE_URL}${endpoint}`, {
      tags: { name: `analytics_${endpoint.split('/').pop()}` }
    });
    responseTime.add(Date.now() - start);
    
    const success = check(response, {
      [`${endpoint} status ok`]: (r) => r.status === 200,
      [`${endpoint} response time ok`]: (r) => r.timings.duration < 2000,
      [`${endpoint} has data`]: (r) => {
        try {
          const data = r.json();
          return data && typeof data === 'object';
        } catch {
          return false;
        }
      }
    });
    
    if (!success) {
      errorRate.add(1);
    }
    
    sleep(0.2);
  });
  
  // Test notification endpoints
  const notificationStart = Date.now();
  const insightsResponse = http.get(`${BASE_URL}/api/notifications/insights`, {
    tags: { name: 'notifications_insights' }
  });
  responseTime.add(Date.now() - notificationStart);
  
  check(insightsResponse, {
    'insights endpoint ok': (r) => r.status === 200,
    'insights has structure': (r) => {
      try {
        const data = r.json();
        return data.insights !== undefined && data.total !== undefined;
      } catch {
        return false;
      }
    }
  });
  
  // Test marking insights as read
  if (insightsResponse.status === 200) {
    const insights = insightsResponse.json().insights;
    if (insights && insights.length > 0) {
      const markReadPayload = JSON.stringify({
        insight_ids: [insights[0].id]
      });
      
      const markReadResponse = http.post(`${BASE_URL}/api/notifications/insights/read`, 
        markReadPayload, {
          headers: { 'Content-Type': 'application/json' },
          tags: { name: 'mark_read' }
        });
      
      check(markReadResponse, {
        'mark as read ok': (r) => r.status === 200,
      });
    }
  }
  
  // Test stats endpoint
  const statsResponse = http.get(`${BASE_URL}/api/notifications/stats`, {
    tags: { name: 'notification_stats' }
  });
  
  check(statsResponse, {
    'stats endpoint ok': (r) => r.status === 200,
    'stats has metrics': (r) => {
      try {
        const data = r.json();
        return data.total_insights !== undefined;
      } catch {
        return false;
      }
    }
  });
  
  sleep(1);
}

function testRealtimeNotifications() {
  // Authenticate first
  const cookies = authenticate();
  if (!cookies) return;
  
  // Test Server-Sent Events connection
  const sseUrl = `${BASE_URL}/api/notifications/stream`;
  
  // Since k6 doesn't support SSE directly, we'll test the endpoint availability
  const sseResponse = http.get(sseUrl, {
    tags: { name: 'sse_connection' }
  });
  
  check(sseResponse, {
    'SSE endpoint accessible': (r) => r.status === 200,
    'SSE content type correct': (r) => 
      r.headers['Content-Type'] && r.headers['Content-Type'].includes('text/event-stream'),
  });
  
  // Test WebSocket alternative (if implemented)
  const wsUrl = `${WS_URL}/ws/notifications`;
  
  ws.connect(wsUrl, {}, function (socket) {
    wsConnections.add(1);
    
    socket.on('open', () => {
      console.log('WebSocket connection established');
      
      // Send authentication
      socket.send(JSON.stringify({
        type: 'auth',
        token: 'test-token'
      }));
    });
    
    socket.on('message', (data) => {
      notificationsReceived.add(1);
      
      try {
        const message = JSON.parse(data);
        check(message, {
          'valid notification format': (msg) => msg.type !== undefined,
        });
      } catch (e) {
        console.log('Invalid message format:', data);
      }
    });
    
    socket.on('error', (e) => {
      console.log('WebSocket error:', e);
      errorRate.add(1);
    });
    
    // Keep connection alive for test duration
    sleep(10);
    
    socket.close();
  });
  
  sleep(1);
}

// Performance monitoring functions
export function handleSummary(data) {
  return {
    'performance-report.json': JSON.stringify(data, null, 2),
    'performance-summary.txt': generateTextSummary(data),
  };
}

function generateTextSummary(data) {
  const summary = [];
  
  summary.push('=== LoL Match Exporter Performance Test Results ===\n');
  
  // Test execution info
  summary.push(`Test Duration: ${data.state.testRunDurationMs}ms`);
  summary.push(`Virtual Users: ${data.metrics.vus?.values?.value || 'N/A'}`);
  summary.push(`Total Requests: ${data.metrics.http_reqs?.values?.count || 'N/A'}\n`);
  
  // HTTP metrics
  if (data.metrics.http_req_duration) {
    const duration = data.metrics.http_req_duration.values;
    summary.push('HTTP Request Duration:');
    summary.push(`  Average: ${duration.avg?.toFixed(2)}ms`);
    summary.push(`  95th percentile: ${duration['p(95)']?.toFixed(2)}ms`);
    summary.push(`  99th percentile: ${duration['p(99)']?.toFixed(2)}ms`);
    summary.push(`  Max: ${duration.max?.toFixed(2)}ms\n`);
  }
  
  // Error rates
  if (data.metrics.http_req_failed) {
    const errorRate = data.metrics.http_req_failed.values.rate * 100;
    summary.push(`Error Rate: ${errorRate.toFixed(2)}%\n`);
  }
  
  // Custom metrics
  if (data.metrics.response_time) {
    summary.push('Custom Response Time Metrics:');
    summary.push(`  95th percentile: ${data.metrics.response_time.values['p(95)']?.toFixed(2)}ms\n`);
  }
  
  // Thresholds
  summary.push('Threshold Results:');
  Object.entries(data.thresholds || {}).forEach(([name, threshold]) => {
    const status = threshold.ok ? '✓ PASS' : '✗ FAIL';
    summary.push(`  ${name}: ${status}`);
  });
  
  // Recommendations
  summary.push('\nRecommendations:');
  if (data.metrics.http_req_duration?.values?.['p(95)'] > 500) {
    summary.push('  - Consider optimizing API response times');
  }
  if (data.metrics.http_req_failed?.values?.rate > 0.01) {
    summary.push('  - Investigate and fix error sources');
  }
  if (data.metrics.vus?.values?.value > 50) {
    summary.push('  - Monitor server resources under high load');
  }
  
  return summary.join('\n');
}

// Setup and teardown
export function setup() {
  console.log('Starting performance test setup...');
  
  // Verify test environment
  const healthResponse = http.get(`${BASE_URL}/api/health`);
  if (healthResponse.status !== 200) {
    throw new Error('Test environment not ready');
  }
  
  console.log('Performance test environment verified');
  return { startTime: Date.now() };
}

export function teardown(data) {
  const duration = Date.now() - data.startTime;
  console.log(`Performance test completed in ${duration}ms`);
}