import { http, HttpResponse } from 'msw';
import { AuthResponse, User } from '@/types';

// Mock user data for testing
const mockUser: User = {
  id: '1',
  email: 'test@herald.lol',
  username: 'testuser',
  display_name: 'Test User',
  avatar: '',
  bio: '',
  timezone: 'UTC',
  language: 'en',
  is_active: true,
  is_premium: false,
  last_login: new Date().toISOString(),
  login_count: 1,
  total_matches: 0,
  last_sync_at: new Date().toISOString(),
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  preferences: {
    id: '1',
    user_id: '1',
    theme: 'dark',
    compact_mode: false,
    show_detailed_stats: true,
    default_timeframe: '7d',
    email_notifications: true,
    push_notifications: true,
    match_notifications: true,
    rank_change_notifications: true,
    auto_sync_matches: true,
    sync_interval: 300,
    include_normal_games: true,
    include_aram_games: true,
    public_profile: true,
    show_in_leaderboards: true,
    allow_data_export: true,
    receive_ai_coaching: true,
    skill_level: 'intermediate',
    preferred_coaching_style: 'balanced',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
};

const mockAuthResponse: AuthResponse = {
  token: 'mock-jwt-token',
  refresh_token: 'mock-refresh-token',
  user: mockUser,
  expires_in: 86400,
};

export const handlers = [
  // Auth endpoints
  http.post('/api/v1/auth/login', () => {
    return HttpResponse.json(mockAuthResponse);
  }),

  http.post('/api/v1/auth/register', () => {
    return HttpResponse.json(mockAuthResponse);
  }),

  http.get('/api/v1/auth/profile', () => {
    return HttpResponse.json(mockUser);
  }),

  http.post('/api/v1/auth/logout', () => {
    return HttpResponse.json({ success: true });
  }),

  http.post('/api/v1/auth/refresh', () => {
    return HttpResponse.json(mockAuthResponse);
  }),

  // Health check
  http.get('/api/v1/health', () => {
    return HttpResponse.json({
      status: 'ok',
      timestamp: new Date().toISOString(),
      version: '1.0.0',
    });
  }),

  // Riot API endpoints
  http.post('/api/v1/riot/link', () => {
    return HttpResponse.json({
      id: '1',
      user_id: '1',
      puuid: 'mock-puuid',
      summoner_name: 'TestSummoner',
      tag_line: 'NA1',
      region: 'na1',
      is_verified: true,
    });
  }),

  http.get('/api/v1/riot/accounts', () => {
    return HttpResponse.json([]);
  }),

  // Match endpoints
  http.get('/api/v1/matches', () => {
    return HttpResponse.json([]);
  }),

  // Analytics endpoints
  http.get('/api/v1/analytics/players/:accountId/stats', () => {
    return HttpResponse.json({
      total_matches: 10,
      wins: 6,
      losses: 4,
      win_rate: 0.6,
      avg_kda: 2.5,
    });
  }),

  // Gaming Analytics Mock Endpoints
  http.get('/api/v1/analytics/kda/:summonerId', () => {
    return HttpResponse.json({
      data: {
        currentKDA: 2.34,
        previousKDA: 2.1,
        trend: 'improving',
        recentMatches: [
          { matchId: 'NA1_1', kda: 2.8, date: '2024-01-20' },
          { matchId: 'NA1_2', kda: 2.1, date: '2024-01-19' },
        ],
        percentile: 65,
        rankComparison: 'above_average'
      },
      confidence: 0.92
    });
  }),

  http.get('/api/v1/analytics/cs/:summonerId', () => {
    return HttpResponse.json({
      data: {
        currentCSPerMin: 7.2,
        previousCSPerMin: 6.8,
        trend: 'improving',
        benchmarks: {
          iron: 4.0,
          bronze: 5.2,
          silver: 6.1,
          gold: 6.8,
          platinum: 7.5,
          diamond: 8.2
        },
        efficiency: 0.85
      },
      confidence: 0.88
    });
  }),

  http.get('/api/v1/analytics/vision/:summonerId', () => {
    return HttpResponse.json({
      data: {
        averageVisionScore: 28.4,
        wardPlacement: 'good',
        visionControl: 0.72,
        heatmapData: [
          { x: 100, y: 200, intensity: 0.8 },
          { x: 150, y: 300, intensity: 0.6 }
        ]
      },
      confidence: 0.79
    });
  }),

  http.get('/api/v1/analytics/damage/:summonerId', () => {
    return HttpResponse.json({
      data: {
        damageShare: 0.32,
        damagePerMinute: 850,
        efficiency: 0.89,
        teamContribution: 'high',
        breakdown: {
          physical: 0.65,
          magic: 0.25,
          true: 0.1
        }
      },
      confidence: 0.91
    });
  }),

  http.get('/api/v1/analytics/gold/:summonerId', () => {
    return HttpResponse.json({
      data: {
        goldPerMinute: 425,
        goldEfficiency: 0.87,
        economicRating: 'good',
        trends: [
          { minute: 5, gold: 2100 },
          { minute: 10, gold: 4250 },
          { minute: 15, gold: 6380 }
        ]
      },
      confidence: 0.84
    });
  }),

  http.get('/api/v1/analytics/champions/:summonerId', () => {
    return HttpResponse.json({
      data: {
        masteryData: [
          { champion: 'Jinx', level: 7, points: 156789, winRate: 0.68 },
          { champion: 'Caitlyn', level: 6, points: 89456, winRate: 0.72 }
        ],
        recommendations: ['Focus on late-game positioning', 'Improve early trades']
      },
      confidence: 0.86
    });
  }),

  http.get('/api/v1/analytics/meta', () => {
    return HttpResponse.json({
      data: {
        tierList: [
          { champion: 'Jinx', tier: 'S', winRate: 0.52, pickRate: 0.14 },
          { champion: 'Caitlyn', tier: 'A', winRate: 0.49, pickRate: 0.18 }
        ],
        trends: 'marksmen_meta',
        patchImpact: 'moderate'
      },
      confidence: 0.95
    });
  }),

  http.post('/api/v1/analytics/predict/match', () => {
    return HttpResponse.json({
      data: {
        winProbability: 0.67,
        keyFactors: ['champion_comp', 'recent_form', 'skill_matchup'],
        predictions: {
          gameLength: 1800,
          playerPerformance: 'above_average'
        }
      },
      confidence: 0.73
    });
  }),

  http.post('/api/v1/team-composition/optimize', () => {
    return HttpResponse.json({
      data: {
        recommendations: [
          {
            composition: ['Jinx', 'Thresh', 'Graves', 'Orianna', 'Malphite'],
            synergy: 0.89,
            winRate: 0.64,
            strategy: 'teamfight_focused'
          }
        ]
      },
      confidence: 0.81
    });
  }),

  http.post('/api/v1/counter-picks/analyze', () => {
    return HttpResponse.json({
      data: {
        counters: [
          { champion: 'Caitlyn', advantage: 0.15, confidence: 0.87 },
          { champion: 'Draven', advantage: 0.12, confidence: 0.82 }
        ],
        reasoning: 'Range advantage and early game pressure'
      },
      confidence: 0.85
    });
  }),
];

// Error handlers
export const errorHandlers = [
  http.post('/api/v1/auth/login', () => {
    return HttpResponse.json(
      { error: 'Invalid credentials', message: 'Email or password is incorrect' },
      { status: 401 }
    );
  }),

  http.get('/api/v1/auth/profile', () => {
    return HttpResponse.json(
      { error: 'Unauthorized', message: 'Token is invalid' },
      { status: 401 }
    );
  }),
];