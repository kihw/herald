// Herald.lol Gaming Analytics Platform - Load Testing Script
// K6 performance testing for gaming workloads

import http from 'k6/http';
import ws from 'k6/ws';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Gaming-specific metrics
const gamingAnalyticsResponseTime = new Trend('gaming_analytics_response_time');
const gamingApiErrors = new Rate('gaming_api_error_rate');
const riotApiCalls = new Counter('riot_api_calls');
const matchProcessingTime = new Trend('match_processing_time');
const realtimeConnectionSuccess = new Rate('realtime_connection_success');

// Test configuration for gaming workloads
export let options = {
  // Gaming analytics load testing profile
  stages: [
    // Ramp up to simulate game start surge
    { duration: '30s', target: 100 }, // Game lobby creation
    { duration: '1m', target: 500 },  // Match start surge
    { duration: '2m', target: 1000 }, // Peak gaming hours
    { duration: '5m', target: 1500 }, // Sustained gaming load
    { duration: '2m', target: 2000 }, // Maximum concurrent users
    { duration: '3m', target: 1000 }, // Sustained high load
    { duration: '1m', target: 500 },  // Cool down period
    { duration: '30s', target: 0 },   // Complete cool down
  ],
  
  // Gaming performance thresholds
  thresholds: {
    // API response time: <500ms for gaming analytics
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    
    // Gaming analytics specific: <5s for post-game analysis
    gaming_analytics_response_time: ['p(95)<5000', 'p(99)<8000'],
    
    // Gaming API error rate: <1%
    gaming_api_error_rate: ['rate<0.01'],
    
    // Match processing: <2s average
    match_processing_time: ['p(95)<2000'],
    
    // Real-time connections: >99% success
    realtime_connection_success: ['rate>0.99'],
    
    // Overall API success rate
    http_req_failed: ['rate<0.05'],
  },
  
  // Test scenarios
  scenarios: {
    // Gaming API load testing
    gaming_api_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 500 },
        { duration: '5m', target: 500 },
        { duration: '1m', target: 0 },
      ],
      exec: 'testGamingAPI',
    },
    
    // Real-time WebSocket connections
    realtime_websocket: {
      executor: 'constant-vus',
      vus: 100,
      duration: '8m',
      exec: 'testRealtimeConnections',
    },
    
    // Gaming analytics processing
    analytics_processing: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 200 },
        { duration: '3m', target: 200 },
        { duration: '1m', target: 0 },
      ],
      exec: 'testGamingAnalytics',
    },
  },
  
  // Tags for gaming metrics
  tags: {
    test_type: 'gaming_load_test',
    service: 'herald_gaming_platform',
    region: 'us-east-1',
  },
};

// Base URL configuration
const BASE_URL = __ENV.BASE_URL || 'https://api.herald.lol';
const WS_URL = __ENV.WS_URL || 'wss://ws.herald.lol';

// Gaming test data
const GAMING_TEST_DATA = {
  summonerNames: ['Faker', 'Caps', 'Bjergsen', 'Doublelift', 'Rekkles'],
  regions: ['na1', 'euw1', 'kr'],
  championIds: [1, 2, 3, 4, 5, 17, 18, 19, 20, 21],
  queueIds: [420, 440, 450], // Ranked Solo, Ranked Flex, ARAM
  
  // Sample match data for testing
  sampleMatch: {
    matchId: 'NA1_4567890123',
    gameId: 4567890123,
    platformId: 'NA1',
    gameCreation: Date.now() - 3600000, // 1 hour ago
    gameDuration: 1845, // ~30 minutes
    queueId: 420,
    mapId: 11,
    seasonId: 13,
    gameVersion: '13.24.1',
    participants: [
      {
        participantId: 1,
        teamId: 100,
        championId: 17,
        summonerId: 'test_summoner_1',
        kills: 7,
        deaths: 3,
        assists: 12,
        totalMinionsKilled: 245,
        visionScore: 28,
        goldEarned: 15420,
      }
    ]
  }
};

// Gaming API load testing
export function testGamingAPI() {
  group('Gaming API Load Test', () => {
    const params = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer test-token',
        'X-Gaming-Platform': 'herald',
        'X-Test-Type': 'load-test',
      },
      timeout: '30s',
    };

    group('Summoner Data Retrieval', () => {
      const summonerName = GAMING_TEST_DATA.summonerNames[Math.floor(Math.random() * GAMING_TEST_DATA.summonerNames.length)];
      const region = GAMING_TEST_DATA.regions[Math.floor(Math.random() * GAMING_TEST_DATA.regions.length)];
      
      const response = http.get(`${BASE_URL}/api/v1/gaming/summoner/${region}/${summonerName}`, params);
      
      check(response, {
        'summoner data retrieved successfully': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
        'has summoner data': (r) => {
          try {
            const data = JSON.parse(r.body);
            return data.id && data.name && data.summonerLevel;
          } catch (e) {
            return false;
          }
        },
      });
      
      gamingApiErrors.add(response.status !== 200);
      riotApiCalls.add(1);
    });

    group('Match History Retrieval', () => {
      const response = http.get(`${BASE_URL}/api/v1/gaming/matches/recent?count=10`, params);
      
      check(response, {
        'match history retrieved': (r) => r.status === 200,
        'response time acceptable': (r) => r.timings.duration < 1000,
        'has match data': (r) => {
          try {
            const data = JSON.parse(r.body);
            return Array.isArray(data.matches) && data.matches.length > 0;
          } catch (e) {
            return false;
          }
        },
      });
      
      gamingApiErrors.add(response.status !== 200);
    });

    group('Gaming Leaderboards', () => {
      const response = http.get(`${BASE_URL}/api/v1/gaming/leaderboard/ranked?queue=420`, params);
      
      check(response, {
        'leaderboard data retrieved': (r) => r.status === 200,
        'fast leaderboard response': (r) => r.timings.duration < 300,
        'has ranking data': (r) => {
          try {
            const data = JSON.parse(r.body);
            return Array.isArray(data.leaderboard) && data.leaderboard.length > 0;
          } catch (e) {
            return false;
          }
        },
      });
      
      gamingApiErrors.add(response.status !== 200);
    });
  });

  sleep(Math.random() * 2 + 1); // 1-3 seconds between requests
}

// Gaming analytics processing test
export function testGamingAnalytics() {
  group('Gaming Analytics Processing', () => {
    const params = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer test-token',
        'X-Analytics-Type': 'post-game',
      },
      timeout: '10s', // Gaming analytics target: <5s
    };

    group('Post-Game Analysis', () => {
      const startTime = Date.now();
      
      const response = http.post(
        `${BASE_URL}/api/v1/analytics/match/analyze`,
        JSON.stringify({
          matchId: GAMING_TEST_DATA.sampleMatch.matchId,
          detailed: true,
          includeTimeline: true,
        }),
        params
      );
      
      const processingTime = Date.now() - startTime;
      gamingAnalyticsResponseTime.add(processingTime);
      
      check(response, {
        'analytics processing successful': (r) => r.status === 200,
        'processing time < 5s (gaming target)': (r) => processingTime < 5000,
        'has analytics data': (r) => {
          try {
            const data = JSON.parse(r.body);
            return data.analysis && data.metrics && data.recommendations;
          } catch (e) {
            return false;
          }
        },
      });
      
      gamingApiErrors.add(response.status !== 200);
    });

    group('Performance Metrics Calculation', () => {
      const response = http.post(
        `${BASE_URL}/api/v1/analytics/performance/calculate`,
        JSON.stringify({
          summonerId: 'test_summoner_1',
          gameMode: 'ranked',
          timeframe: '7d',
          metrics: ['kda', 'cs_per_min', 'vision_score', 'damage_share'],
        }),
        params
      );
      
      const processingTime = response.timings.duration;
      gamingAnalyticsResponseTime.add(processingTime);
      
      check(response, {
        'performance metrics calculated': (r) => r.status === 200,
        'metrics processing fast': (r) => r.timings.duration < 2000,
        'has performance data': (r) => {
          try {
            const data = JSON.parse(r.body);
            return data.kda && data.csPerMin && data.visionScore;
          } catch (e) {
            return false;
          }
        },
      });
      
      gamingApiErrors.add(response.status !== 200);
    });

    group('Improvement Recommendations', () => {
      const response = http.get(
        `${BASE_URL}/api/v1/analytics/recommendations?championId=17&role=adc`,
        params
      );
      
      check(response, {
        'recommendations generated': (r) => r.status === 200,
        'recommendation response fast': (r) => r.timings.duration < 1000,
        'has recommendations': (r) => {
          try {
            const data = JSON.parse(r.body);
            return Array.isArray(data.recommendations) && data.recommendations.length > 0;
          } catch (e) {
            return false;
          }
        },
      });
      
      gamingApiErrors.add(response.status !== 200);
    });
  });

  sleep(Math.random() * 3 + 2); // 2-5 seconds between analytics requests
}

// Real-time WebSocket connection testing
export function testRealtimeConnections() {
  group('Real-time WebSocket Gaming', () => {
    const url = `${WS_URL}/gaming/realtime`;
    let connectionSuccess = false;
    
    const response = ws.connect(url, {
      headers: {
        'Authorization': 'Bearer test-token',
        'X-Gaming-Platform': 'herald',
      },
    }, (socket) => {
      connectionSuccess = true;
      
      socket.on('open', () => {
        console.log('🎮 WebSocket connection established');
        
        // Subscribe to live match updates
        socket.send(JSON.stringify({
          type: 'subscribe',
          channels: ['live-matches', 'performance-updates'],
          filters: {
            region: 'na1',
            queueType: 'ranked',
          },
        }));
      });

      socket.on('message', (data) => {
        try {
          const message = JSON.parse(data);
          console.log('📊 Received gaming update:', message.type);
          
          // Validate message structure
          check(message, {
            'valid message format': (msg) => msg.type && msg.data,
            'gaming message types': (msg) => ['match_update', 'performance_alert', 'analytics_complete'].includes(msg.type),
          });
        } catch (e) {
          console.log('❌ Invalid WebSocket message:', e);
        }
      });

      socket.on('error', (e) => {
        console.log('🔥 WebSocket error:', e);
        connectionSuccess = false;
      });

      // Keep connection alive for realistic gaming session duration
      socket.setTimeout(() => {
        console.log('⏰ Closing WebSocket connection after timeout');
        socket.close();
      }, 30000 + Math.random() * 60000); // 30-90 seconds
    });
    
    realtimeConnectionSuccess.add(connectionSuccess);
    
    check(response, {
      'websocket connection established': (r) => r && r.status === 101,
    });
  });
}

// Setup function - run once before all tests
export function setup() {
  console.log('🎮 Starting Herald.lol Gaming Analytics Load Test');
  console.log('📊 Test Configuration:');
  console.log(`   - Base URL: ${BASE_URL}`);
  console.log(`   - WebSocket URL: ${WS_URL}`);
  console.log(`   - Max VUs: 2000`);
  console.log(`   - Duration: 8+ minutes`);
  console.log('🎯 Performance Targets:');
  console.log('   - API Response: <500ms (p95)');
  console.log('   - Gaming Analytics: <5s (p95)');
  console.log('   - Error Rate: <1%');
  console.log('   - WebSocket Success: >99%');
  
  // Verify API is accessible
  const healthCheck = http.get(`${BASE_URL}/health`);
  if (healthCheck.status !== 200) {
    throw new Error(`API health check failed: ${healthCheck.status}`);
  }
  
  console.log('✅ API health check passed - starting load test');
  return { startTime: Date.now() };
}

// Teardown function - run once after all tests
export function teardown(data) {
  const testDuration = (Date.now() - data.startTime) / 1000;
  console.log('🏁 Herald.lol Gaming Load Test Complete');
  console.log(`📈 Test Duration: ${testDuration.toFixed(1)}s`);
  console.log('📊 Check the metrics above for detailed performance results');
  console.log('🎮 Gaming analytics platform load testing finished!');
}

// Handle test options based on environment
export default function() {
  // This is the default function that gets called for each VU iteration
  // The actual test logic is in the scenario functions above
}