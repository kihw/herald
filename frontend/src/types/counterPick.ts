// Counter Pick Types for Herald.lol Frontend

export interface CounterPickAnalysis {
  id: string;
  targetChampion: string;
  targetRole: string;
  counterPicks: CounterPickSuggestion[];
  laneCounters: LaneCounterData[];
  teamFightCounters: TeamFightCounterData[];
  itemCounters: ItemCounterData[];
  playStyleCounters: PlayStyleCounterData[];
  metaContext: CounterMetaContext;
  confidence: number;
  createdAt: string;
}

export interface CounterPickSuggestion {
  champion: string;
  counterStrength: number; // 0-100
  winRateAdvantage: number;
  laneAdvantage: number;
  teamFightAdvantage: number;
  scalingAdvantage: number;
  counterReasons: string[];
  playingTips: string[];
  itemRecommendations: string[];
  powerSpikes: PowerSpikeData[];
  weaknesses: CounterWeakness[];
  matchupDifficulty: 'easy' | 'moderate' | 'hard';
  metaFit: number;
  playerComfort?: number;
  banPriority: number;
  flexibility: number;
  safetyRating: number;
}

export interface PowerSpikeData {
  timing: string; // e.g., "Level 6", "First item", "20 minutes"
  strength: number; // 0-100
  duration: string;
  counterPlay: string[];
}

export interface CounterWeakness {
  weakness: string;
  severity: 'minor' | 'moderate' | 'major';
  exploitHow: string[];
  timing: string;
}

export interface LaneCounterData {
  phase: 'early' | 'mid' | 'late';
  advantage: number; // -100 to 100
  keyFactors: string[];
  playingTips: string[];
  wardingTips: string[];
  tradingPatterns: string[];
  allInPotential: number;
  roamingPotential: number;
  scalingComparison: string;
}

export interface TeamFightCounterData {
  counterType: 'engage' | 'disengage' | 'peel' | 'burst' | 'sustain' | 'crowd_control';
  effectiveness: number; // 0-100
  positioning: string[];
  comboCounters: string[];
  teamCoordination: string[];
  objectiveControl: string[];
}

export interface ItemCounterData {
  itemName: string;
  counterType: 'defensive' | 'offensive' | 'utility';
  effectiveness: number; // 0-100
  buildPriority: number; // 1-6
  situationalUse: string;
  costEffectiveness: number;
}

export interface PlayStyleCounterData {
  targetPlayStyle: string;
  counterStrategy: string;
  keyPrinciples: string[];
  timing: string[];
  teamSupport: string[];
  riskLevel: 'low' | 'medium' | 'high';
}

export interface CounterMetaContext {
  patch: string;
  targetPickRate: number;
  targetBanRate: number;
  counterPickRate?: number;
  metaTrend: 'rising' | 'stable' | 'declining';
  proPlayUsage: number;
}

// Multi-Target Counter Analysis
export interface MultiTargetCounterAnalysis {
  id: string;
  targetChampions: TargetChampionData[];
  universalCounters: UniversalCounterSuggestion[];
  specificCounters: SpecificCounterSuggestion[];
  teamCounters: TeamCounterStrategy[];
  banRecommendations: BanRecommendation[];
  overallStrategy: CounterStrategy;
  confidence: number;
  createdAt: string;
}

export interface TargetChampionData {
  champion: string;
  role: string;
  threatLevel: 'low' | 'medium' | 'high' | 'critical';
  priority: number; // 0-100
  reasons: string[];
}

export interface UniversalCounterSuggestion {
  champion: string;
  countersTargets: string[];
  averageStrength: number;
  versatility: number;
  recommendReasons: string[];
}

export interface SpecificCounterSuggestion {
  champion: string;
  primaryTarget: string;
  secondaryTargets: string[];
  counterStrength: number;
  specialization: string;
}

export interface TeamCounterStrategy {
  strategy: string;
  requiredChampions: string[];
  effectiveness: number;
  complexity: 'simple' | 'moderate' | 'complex';
  description: string;
  execution: string[];
}

export interface BanRecommendation {
  champion: string;
  priority: number; // 0-100
  reasoning: string;
  impact: string;
  alternatives: string[];
}

export interface CounterStrategy {
  primary: string;
  secondary?: string;
  approach: string;
  keyPrinciples: string[];
  timeline: string[];
}

// Request Types
export interface CounterPickRequest {
  targetChampion: string;
  targetRole: string;
  gameMode?: string;
  playerChampionPool?: string[];
  playerRank?: string;
  preferences?: CounterPickPreferences;
}

export interface CounterPickPreferences {
  prioritizeLane: boolean;
  prioritizeTeamfight: boolean;
  prioritizeMeta: boolean;
  prioritizeComfort: boolean;
}

export interface MultiTargetCounterRequest {
  targetChampions: Array<{
    champion: string;
    role: string;
    threatLevel?: 'low' | 'medium' | 'high' | 'critical';
    priority?: number;
  }>;
  playerChampionPool?: string[];
  gameMode?: string;
  strategy?: 'universal' | 'specific' | 'hybrid';
}

export interface TeamCounterRequest {
  enemyTeam: Array<{
    champion: string;
    role: string;
  }>;
  ourTeam?: Array<{
    champion?: string;
    role: string;
  }>;
  gameMode?: string;
}

// Matchup Analysis
export interface MatchupAnalysis {
  champion1: string;
  champion2: string;
  role: string;
  champion1VsChampion2: CounterPickSuggestion;
  champion2VsChampion1: CounterPickSuggestion;
  lanePhases: LaneCounterData[];
  teamFightData: TeamFightCounterData[];
  itemCounters: ItemCounterData[];
  overallAssessment: MatchupAssessment;
}

export interface MatchupAssessment {
  favoredChampion: string;
  confidenceLevel: number;
  keyFactors: string[];
  gameplan: {
    early: string[];
    mid: string[];
    late: string[];
  };
  criticalMoments: string[];
}

// User Favorites and History
export interface CounterPickFavorite {
  id: number;
  champion: string;
  targetChampion: string;
  role: string;
  notes: string;
  personalRating: number; // 1-10
  timesUsed: number;
  successRate: number;
  createdAt: string;
  updatedAt: string;
}

export interface CounterPickHistory {
  id: number;
  champion: string;
  targetChampion: string;
  role: string;
  gameMode: string;
  matchId: string;
  result: 'win' | 'loss';
  performance: number; // 0-100
  predictedStrength: number;
  actualStrength: number;
  accuracy: number;
  createdAt: string;
}

export interface CounterPickMetrics {
  champion: string;
  targetChampion: string;
  role: string;
  gameMode: string;
  sampleSize: number;
  winRate: number;
  laneWinRate: number;
  kda: number;
  damageShare: number;
  goldDifferential: number;
  csAdvantage: number;
  visionScore: number;
  counterStrength: number;
  confidence: number;
}

// Performance Analysis
export interface CounterPerformance {
  summonerId: string;
  champion: string;
  performance: {
    overallWinRate: number;
    counterWinRate: number;
    averageKDA: number;
    gamesPlayed: number;
    bestMatchups: string[];
    worstMatchups: string[];
    improvement: {
      trend: 'improving' | 'stable' | 'declining';
      recentPerformance: number;
      areas: string[];
    };
  };
}

// Meta Counter Data
export interface MetaCounterData {
  targetChampion: string;
  targetRole: string;
  topCounters: CounterPickSuggestion[];
  metaContext: CounterMetaContext;
  popularity: {
    pickRate: number;
    banRate: number;
    counterRate: number;
  };
  trends: {
    rising: string[];
    stable: string[];
    declining: string[];
  };
}

// Ban Strategy Types
export interface CounterBanStrategy {
  banStrategy: BanRecommendation[];
  universalCounters: UniversalCounterSuggestion[];
  specificCounters: SpecificCounterSuggestion[];
  recommendation: CounterStrategy;
  confidence: number;
}

export interface BanStrategyRequest {
  threateningChampions: string[];
  playerChampionPool?: string[];
  banPhase?: 'first' | 'second';
  gameMode?: string;
}

// Threat Assessment
export interface ThreatAssessment {
  champions: string[];
  threatAssessment: {
    overall: 'low' | 'medium' | 'high' | 'critical';
    early: 'low' | 'medium' | 'high' | 'critical';
    mid: 'low' | 'medium' | 'high' | 'critical';
    late: 'low' | 'medium' | 'high' | 'critical';
    priorities: string[];
  };
  recommendations: {
    bans: string[];
    counters: string[];
    strategies: string[];
  };
}

// Champion Pool Analysis
export interface ChampionPoolCounterAnalysis {
  summonerId: string;
  championPool: Array<{
    champion: string;
    role: string;
    comfortLevel: number;
    counters: string[];
    vulnerabilities: string[];
    strengths: string[];
    recommendations: string[];
  }>;
  gaps: {
    missingCounters: string[];
    weakMatchups: string[];
    suggestions: string[];
  };
  overall: {
    versatility: number;
    counterCoverage: number;
    metaAlignment: number;
    recommendations: string[];
  };
}

// Coaching and Tips
export interface CounterPickTips {
  champion: string;
  targetChampion: string;
  role: string;
  tips: {
    laning: CounterTip[];
    teamfighting: CounterTip[];
    itemization: CounterTip[];
    general: CounterTip[];
  };
  warnings: CounterWarning[];
  keyTimings: KeyTiming[];
}

export interface CounterTip {
  category: string;
  tip: string;
  importance: 'low' | 'medium' | 'high' | 'critical';
  phase: 'early' | 'mid' | 'late' | 'all';
  situational?: string;
}

export interface CounterWarning {
  warning: string;
  severity: 'minor' | 'moderate' | 'major' | 'critical';
  timing: string;
  mitigation: string[];
}

export interface KeyTiming {
  timing: string;
  description: string;
  opportunity: string;
  risk: string;
}

// Chart and Visualization Data
export interface CounterPickChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    backgroundColor: string;
    borderColor: string;
    fill?: boolean;
    tension?: number;
  }>;
}

export interface WinRateComparisonChart {
  champions: string[];
  winRates: number[];
  laneWinRates: number[];
  teamFightWinRates: number[];
  confidence: number[];
}

export interface CounterStrengthRadar {
  labels: string[];
  data: number[];
  backgroundColor: string;
  borderColor: string;
  pointBackgroundColor: string;
  maxValue: number;
}

export interface LanePhaseChart {
  phases: Array<{
    phase: string;
    advantage: number;
    confidence: number;
    keyFactors: string[];
  }>;
}

// Feedback and Learning
export interface CounterPickFeedback {
  champion: string;
  targetChampion: string;
  matchId: string;
  result: 'win' | 'loss';
  performance: number; // 1-10
  feedback: string;
  categories: {
    laning: number;
    teamfighting: number;
    itemization: number;
    execution: number;
  };
  improvements: string[];
}

// Service Response Types
export interface CounterPickResponse {
  targetChampion: string;
  targetRole: string;
  gameMode: string;
  suggestions: CounterPickSuggestion[];
  metaContext: CounterMetaContext;
  confidence: number;
}

export interface PersonalizedCounterResponse {
  summonerId: string;
  targetChampion: string;
  targetRole: string;
  personalizedCounters: CounterPickSuggestion[];
  recommendations: {
    comfort: string;
    meta: string;
    winrate: string;
  };
  confidence: number;
}

// Filter and Sort Options
export interface CounterPickFilters {
  minCounterStrength?: number;
  maxBanRate?: number;
  minMetaFit?: number;
  difficulties?: ('easy' | 'moderate' | 'hard')[];
  roles?: string[];
  counterTypes?: ('lane' | 'teamfight' | 'scaling' | 'item')[];
  playerChampionsOnly?: boolean;
}

export interface CounterPickSortOptions {
  sortBy: 'counterStrength' | 'winRate' | 'metaFit' | 'safetyRating' | 'playerComfort';
  sortOrder: 'asc' | 'desc';
}

// Utility Types
export type CounterPickDifficulty = 'easy' | 'moderate' | 'hard';
export type CounterType = 'lane' | 'teamfight' | 'scaling' | 'item' | 'playstyle';
export type GamePhase = 'early' | 'mid' | 'late';
export type ThreatLevel = 'low' | 'medium' | 'high' | 'critical';
export type CounterConfidence = 'low' | 'medium' | 'high' | 'very_high';

// API Error Types
export interface CounterPickError {
  code: string;
  message: string;
  details?: string;
  suggestions?: string[];
}

// Export all types as a namespace for easier imports
export namespace CounterPick {
  export type Analysis = CounterPickAnalysis;
  export type Suggestion = CounterPickSuggestion;
  export type Request = CounterPickRequest;
  export type Preferences = CounterPickPreferences;
  export type Response = CounterPickResponse;
  export type Filters = CounterPickFilters;
  export type SortOptions = CounterPickSortOptions;
}