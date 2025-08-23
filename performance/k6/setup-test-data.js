// Herald.lol Performance Test Data Setup
// Generates test data for k6 performance testing scenarios

const fs = require('fs');
const path = require('path');

// Herald.lol gaming test data
const GAMING_TEST_DATA = {
  // Test users for performance testing
  testUsers: Array.from({ length: 1000 }, (_, i) => ({
    email: `perftest${i}@herald.lol`,
    username: `PerfTest${i}`,
    password: 'PerfTest123!',
    displayName: `Performance Tester ${i}`,
    role: 'user',
    isVerified: true,
    createdAt: new Date().toISOString(),
  })),

  // Summoner data for Riot API testing
  testSummoners: [
    'PerfTestSummoner1',
    'PerfTestSummoner2', 
    'PerfTestSummoner3',
    'PerfTestSummoner4',
    'PerfTestSummoner5',
    'LoadTestPlayer1',
    'LoadTestPlayer2',
    'StressTestUser1',
    'StressTestUser2',
    'SpikeTestGamer1',
  ],

  // Match IDs for analysis testing
  testMatches: Array.from({ length: 100 }, (_, i) => 
    `NA1_${4567890000 + i}`
  ),

  // PUUIDs for player identification
  testPuuids: Array.from({ length: 100 }, (_, i) => 
    `test-puuid-${String(i).padStart(5, '0')}`
  ),

  // Champions for testing
  champions: [
    'Jinx', 'Caitlyn', 'Ezreal', 'Vayne', 'Ashe',
    'Jhin', 'MissFortune', 'Lucian', 'Draven', 'Sivir',
    'Leona', 'Thresh', 'Braum', 'Morgana', 'Lulu',
    'Yasuo', 'Zed', 'Akali', 'Katarina', 'Talon',
    'Lee Sin', 'Graves', 'Kindred', 'Kha\'Zix', 'Rengar',
  ],

  // Gaming analytics test data
  analyticsTestData: {
    kda: Array.from({ length: 50 }, () => ({
      kills: Math.floor(Math.random() * 20),
      deaths: Math.floor(Math.random() * 10) + 1,
      assists: Math.floor(Math.random() * 25),
      matchId: `NA1_${Math.floor(Math.random() * 1000000000)}`,
      timestamp: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
    })),

    csData: Array.from({ length: 50 }, () => ({
      totalCS: Math.floor(Math.random() * 300) + 100,
      gameDuration: Math.floor(Math.random() * 1800) + 1200, // 20-50 minutes
      matchId: `NA1_${Math.floor(Math.random() * 1000000000)}`,
      timestamp: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
    })),

    visionData: Array.from({ length: 50 }, () => ({
      visionScore: Math.floor(Math.random() * 100) + 20,
      wardsPlaced: Math.floor(Math.random() * 20) + 5,
      wardsCleared: Math.floor(Math.random() * 15),
      matchId: `NA1_${Math.floor(Math.random() * 1000000000)}`,
      timestamp: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
    })),
  },

  // Load test scenarios
  loadScenarios: [
    {
      name: 'gaming_dashboard_load',
      description: 'Standard gaming dashboard load',
      concurrent_users: 1000,
      duration: '5m',
      ramp_up: '2m',
    },
    {
      name: 'post_game_analysis_rush',
      description: 'Post-game analysis rush hour',
      concurrent_users: 5000,
      duration: '10m', 
      ramp_up: '1m',
    },
    {
      name: 'tournament_end_spike',
      description: 'Tournament end analysis spike',
      concurrent_users: 10000,
      duration: '5m',
      ramp_up: '30s',
    },
    {
      name: 'rank_reset_surge',
      description: 'Season start rank checking surge',
      concurrent_users: 8000,
      duration: '15m',
      ramp_up: '2m',
    },
  ],

  // Performance thresholds for Herald.lol
  performanceThresholds: {
    analytics_response_time: 5000, // 5 seconds max
    ui_load_time: 2000,            // 2 seconds max
    uptime_target: 0.999,          // 99.9% uptime
    error_rate_max: 0.01,          // 1% max error rate
    concurrent_users_target: 1000000, // 1M concurrent target
  },

  // Riot API test configurations
  riotApiConfig: {
    rate_limits: {
      personal: { requests: 100, per_seconds: 120 },
      production: { requests: 3000, per_seconds: 10 },
    },
    endpoints: [
      '/lol/summoner/v4/summoners/by-name/{summonerName}',
      '/lol/league/v4/entries/by-summoner/{encryptedSummonerId}',
      '/lol/match/v5/matches/by-puuid/{puuid}/ids',
      '/lol/match/v5/matches/{matchId}',
      '/lol/spectator/v4/active-games/by-summoner/{encryptedSummonerId}',
    ],
  },
};

// Function to create test data files
function createTestDataFiles() {
  const dataDir = path.join(__dirname, 'test-data');
  
  // Create data directory
  if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
  }

  // Write test users file
  fs.writeFileSync(
    path.join(dataDir, 'test-users.json'),
    JSON.stringify(GAMING_TEST_DATA.testUsers, null, 2)
  );

  // Write summoner test data
  fs.writeFileSync(
    path.join(dataDir, 'test-summoners.json'),
    JSON.stringify(GAMING_TEST_DATA.testSummoners, null, 2)
  );

  // Write match test data
  fs.writeFileSync(
    path.join(dataDir, 'test-matches.json'),
    JSON.stringify(GAMING_TEST_DATA.testMatches, null, 2)
  );

  // Write analytics test data
  fs.writeFileSync(
    path.join(dataDir, 'analytics-test-data.json'),
    JSON.stringify(GAMING_TEST_DATA.analyticsTestData, null, 2)
  );

  // Write load scenarios
  fs.writeFileSync(
    path.join(dataDir, 'load-scenarios.json'),
    JSON.stringify(GAMING_TEST_DATA.loadScenarios, null, 2)
  );

  // Write performance thresholds
  fs.writeFileSync(
    path.join(dataDir, 'performance-thresholds.json'),
    JSON.stringify(GAMING_TEST_DATA.performanceThresholds, null, 2)
  );

  console.log('✅ Herald.lol test data files created successfully!');
  console.log(`📁 Test data directory: ${dataDir}`);
  console.log('📊 Test data files:');
  console.log('   - test-users.json (1000 test users)');
  console.log('   - test-summoners.json (10 test summoners)');
  console.log('   - test-matches.json (100 test matches)');
  console.log('   - analytics-test-data.json (gaming analytics data)');
  console.log('   - load-scenarios.json (performance test scenarios)');
  console.log('   - performance-thresholds.json (Herald.lol requirements)');
}

// Function to generate k6 data files
function generateK6DataFiles() {
  const k6DataDir = path.join(__dirname, 'k6-data');
  
  // Create k6 data directory
  if (!fs.existsSync(k6DataDir)) {
    fs.mkdirSync(k6DataDir, { recursive: true });
  }

  // Create CSV files for k6 data sources
  
  // Test users CSV
  const usersCSV = [
    'email,username,password',
    ...GAMING_TEST_DATA.testUsers.slice(0, 100).map(user => 
      `${user.email},${user.username},${user.password}`
    )
  ].join('\n');
  
  fs.writeFileSync(path.join(k6DataDir, 'test-users.csv'), usersCSV);

  // Summoners CSV
  const summonersCSV = [
    'summonerName',
    ...GAMING_TEST_DATA.testSummoners.map(summoner => summoner)
  ].join('\n');
  
  fs.writeFileSync(path.join(k6DataDir, 'summoners.csv'), summonersCSV);

  // Matches CSV
  const matchesCSV = [
    'matchId',
    ...GAMING_TEST_DATA.testMatches.slice(0, 50).map(match => match)
  ].join('\n');
  
  fs.writeFileSync(path.join(k6DataDir, 'matches.csv'), matchesCSV);

  console.log('✅ k6 data files created successfully!');
  console.log(`📁 k6 data directory: ${k6DataDir}`);
  console.log('📊 k6 data files:');
  console.log('   - test-users.csv (100 test users for k6)');
  console.log('   - summoners.csv (10 summoners for testing)');
  console.log('   - matches.csv (50 matches for analysis testing)');
}

// Main execution
function main() {
  console.log('🎮 Herald.lol Performance Test Data Setup');
  console.log('========================================');
  console.log('');

  try {
    createTestDataFiles();
    console.log('');
    generateK6DataFiles();
    console.log('');
    console.log('🎯 Test data setup complete!');
    console.log('🚀 Ready for Herald.lol performance testing with k6');
  } catch (error) {
    console.error('❌ Error setting up test data:', error);
    process.exit(1);
  }
}

// Run if called directly
if (require.main === module) {
  main();
}

module.exports = {
  GAMING_TEST_DATA,
  createTestDataFiles,
  generateK6DataFiles,
};