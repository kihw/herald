// Test setup for Herald.lol Frontend
import '@testing-library/jest-dom';
import { cleanup } from '@testing-library/react';
import { afterEach, beforeAll, afterAll, vi } from 'vitest';
import { server } from './mocks/server';

// Establish API mocking before all tests
beforeAll(() => {
  server.listen();
});

// Reset any request handlers that we may add during the tests,
// so they don't affect other tests
afterEach(() => {
  server.resetHandlers();
  cleanup();
});

// Clean up after the tests are finished
afterAll(() => {
  server.close();
});

// Mock matchMedia for components that use it
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// Mock ResizeObserver for components that use it
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}));

// Mock IntersectionObserver for components that use it
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}));

// Gaming analytics test helpers
export const mockMatchData = {
  matchId: 'NA1_4567890123',
  summonerId: 'test-summoner-id',
  champion: 'Jinx',
  role: 'ADC',
  kills: 12,
  deaths: 3,
  assists: 18,
  cs: 245,
  gameDuration: 1834000, // 30:34 in milliseconds
  visionScore: 32,
  totalDamageDealt: 45678,
  goldEarned: 16234,
  items: [3031, 3006, 3087, 3036, 3035, 3142],
  kda: 10.0,
  csPerMin: 8.0,
  damageShare: 0.32,
  visionControl: 0.78,
  goldEfficiency: 0.89
};

export const mockUserProfile = {
  id: 'test-user-id',
  summonerId: 'test-summoner-id',
  summonerName: 'TestSummoner',
  puuid: 'test-puuid',
  region: 'na1',
  tier: 'GOLD',
  rank: 'II',
  leaguePoints: 67,
  preferences: {
    theme: 'dark',
    language: 'en',
    notifications: true
  }
};

export const mockAnalyticsData = {
  overallStats: {
    averageKDA: 2.1,
    averageCS: 6.8,
    winRate: 0.64,
    totalGames: 123,
    averageVisionScore: 28
  },
  recentPerformance: {
    last10Games: {
      wins: 7,
      losses: 3,
      averageKDA: 2.4,
      trending: 'up'
    }
  },
  championMastery: [
    { champion: 'Jinx', points: 156789, level: 7, gamesPlayed: 34 },
    { champion: 'Caitlyn', points: 89456, level: 6, gamesPlayed: 28 },
    { champion: 'Ashe', points: 67234, level: 5, gamesPlayed: 22 }
  ]
};

// Gaming calculation test utilities
export const validateGamingMetrics = (analysis: any) => {
  // Ensure gaming calculations are reasonable
  expect(analysis.kda).toBeGreaterThanOrEqual(0);
  expect(analysis.csPerMin).toBeGreaterThanOrEqual(0);
  expect(analysis.visionScore).toBeGreaterThanOrEqual(0);
  expect(analysis.damageShare).toBeGreaterThanOrEqual(0);
  expect(analysis.damageShare).toBeLessThanOrEqual(1);
  expect(analysis.goldEfficiency).toBeGreaterThanOrEqual(0);
  expect(analysis.goldEfficiency).toBeLessThanOrEqual(2); // Can exceed 1 for exceptional performance
};

// Performance test helper
export const measurePerformance = async (fn: () => Promise<any>, maxMs: number = 5000) => {
  const start = Date.now();
  const result = await fn();
  const duration = Date.now() - start;
  
  expect(duration).toBeLessThan(maxMs);
  return { result, duration };
};

// Mock Riot API responses for consistent testing
export const mockRiotApiResponses = {
  summonerByName: {
    id: 'test-summoner-id',
    accountId: 'test-account-id',
    puuid: 'test-puuid',
    name: 'TestSummoner',
    profileIconId: 1234,
    revisionDate: Date.now(),
    summonerLevel: 156
  },
  
  matchIds: ['NA1_4567890123', 'NA1_4567890124', 'NA1_4567890125'],
  
  matchDetail: {
    metadata: {
      matchId: 'NA1_4567890123',
      participants: ['test-puuid']
    },
    info: {
      gameDuration: 1834,
      gameMode: 'CLASSIC',
      gameType: 'MATCHED_GAME',
      participants: [{
        puuid: 'test-puuid',
        championName: 'Jinx',
        teamPosition: 'BOTTOM',
        kills: 12,
        deaths: 3,
        assists: 18,
        totalMinionsKilled: 187,
        neutralMinionsKilled: 58,
        visionScore: 32,
        totalDamageDealtToChampions: 45678,
        goldEarned: 16234,
        item0: 3031,
        item1: 3006,
        item2: 3087,
        item3: 3036,
        item4: 3035,
        item5: 3142,
        win: true
      }]
    }
  }
};

// Test data validation helpers
export const isValidGamingData = (data: any): boolean => {
  if (!data || typeof data !== 'object') return false;
  
  // Basic gaming data structure validation
  const requiredFields = ['matchId', 'summonerId', 'champion', 'kills', 'deaths', 'assists'];
  return requiredFields.every(field => data.hasOwnProperty(field) && data[field] !== undefined);
};

export const isValidAnalyticsResponse = (response: any): boolean => {
  if (!response || typeof response !== 'object') return false;
  
  // Analytics response structure validation
  return response.hasOwnProperty('data') && 
         response.hasOwnProperty('confidence') &&
         typeof response.confidence === 'number' &&
         response.confidence >= 0 && response.confidence <= 1;
};