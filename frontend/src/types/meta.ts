// Meta Analytics Types for Herald.lol Frontend

export interface MetaAnalysis {
  id: string;
  patch: string;
  region: string;
  rank: string;
  timeRange: string;
  analysisDate: string;
  
  // Tier Lists
  tierList: ChampionTierList;
  roleTierLists: Record<string, ChampionTierList>;
  
  // Meta Trends
  metaTrends: MetaTrends;
  emergingPicks: EmergingChampionData[];
  deciningPicks: DecliningChampionData[];
  
  // Statistical Analysis
  championStats: ChampionMetaStats[];
  banAnalysis: BanAnalysis;
  pickAnalysis: PickAnalysis;
  
  // Meta Shifts
  metaShifts: MetaShiftData[];
  patchImpact: PatchImpactAnalysis;
  
  // Predictions
  predictions: MetaPredictions;
  
  // Recommendations
  recommendations: MetaRecommendation[];
  
  // Metadata
  generatedAt: string;
  lastUpdated: string;
  dataQuality: number;
}

export interface ChampionTierList {
  sPlusTier: ChampionTierEntry[];
  sTier: ChampionTierEntry[];
  aPlusTier: ChampionTierEntry[];
  aTier: ChampionTierEntry[];
  bPlusTier: ChampionTierEntry[];
  bTier: ChampionTierEntry[];
  cPlusTier: ChampionTierEntry[];
  cTier: ChampionTierEntry[];
  dTier: ChampionTierEntry[];
  
  lastUpdated: string;
  sampleSize: number;
  confidence: number;
}

export interface ChampionTierEntry {
  champion: string;
  tierScore: number;
  winRate: number;
  pickRate: number;
  banRate: number;
  carryPotential: number;
  versatility: number;
  trendDirection: 'rising' | 'stable' | 'falling';
  tierMovement: number; // +2, +1, 0, -1, -2 from previous patch
  recommendedFor: string[]; // ["climbing", "one_trick", "flex_pick"]
}

export interface MetaTrends {
  dominantStrategies: StrategyTrendData[];
  championTypes: ChampionTypeData;
  gameLength: GameLengthTrends;
  objectivePriority: ObjectiveTrends;
  teamCompositions: TeamCompTrend[];
}

export interface StrategyTrendData {
  strategy: string;
  popularity: number;
  winRate: number;
  trend: 'rising' | 'stable' | 'falling';
  champions: string[];
  description: string;
}

export interface ChampionTypeData {
  tanks: ChampionTypeMetrics;
  fighters: ChampionTypeMetrics;
  assassins: ChampionTypeMetrics;
  mages: ChampionTypeMetrics;
  marksmen: ChampionTypeMetrics;
  support: ChampionTypeMetrics;
}

export interface ChampionTypeMetrics {
  pickRate: number;
  winRate: number;
  banRate: number;
  trendDirection: 'rising' | 'stable' | 'declining';
}

export interface GameLengthTrends {
  averageGameLength: number;
  lengthDistribution: Record<string, number>;
  championsByLength: Record<string, string[]>;
  strategiesByLength: Record<string, string[]>;
}

export interface ObjectiveTrends {
  dragonPriority: ObjectivePriorityData;
  baronPriority: ObjectivePriorityData;
  heraldPriority: ObjectivePriorityData;
  towerPriority: ObjectivePriorityData;
}

export interface ObjectivePriorityData {
  priority: number;
  winRateImpact: number;
  controlRate: number;
  contestedRate: number;
}

export interface TeamCompTrend {
  compositionType: string;
  popularity: number;
  winRate: number;
  champions: string[];
  synergies: string[];
  counters: string[];
}

export interface EmergingChampionData {
  champion: string;
  role: string;
  currentTier: string;
  previousTier: string;
  winRateChange: number;
  pickRateChange: number;
  reasonForRise: string[];
  projectedTier: string;
  confidence: number;
}

export interface DecliningChampionData {
  champion: string;
  role: string;
  currentTier: string;
  previousTier: string;
  winRateChange: number;
  pickRateChange: number;
  reasonForDecline: string[];
  projectedTier: string;
  confidence: number;
}

export interface ChampionMetaStats {
  champion: string;
  role: string;
  tier: string;
  tierScore: number;
  
  // Core Statistics
  winRate: number;
  pickRate: number;
  banRate: number;
  presenceRate: number;
  
  // Performance Metrics
  averageKDA: number;
  damageShare: number;
  goldShare: number;
  visionScore: number;
  
  // Meta Analysis
  versatility: number;
  carryPotential: number;
  teamReliance: number;
  scalingCurve: ScalingCurveData;
  
  // Trend Analysis
  trendData: ChampionMetaTrendPoint[];
  trendDirection: 'rising' | 'stable' | 'declining';
  
  // Comparative Analysis
  rankVariance: Record<string, ChampionRankStats>;
  regionVariance: Record<string, ChampionRegionStats>;
  
  // Build and Play Style
  popularBuilds: MetaBuildData[];
  playStyles: PlayStyleData[];
  
  // Matchup Context
  strongAgainst: string[];
  weakAgainst: string[];
  synergizesWith: string[];
}

export interface ScalingCurveData {
  earlyGame: number;
  midGame: number;
  lateGame: number;
  powerSpikes: number[];
}

export interface ChampionMetaTrendPoint {
  date: string;
  winRate: number;
  pickRate: number;
  banRate: number;
  tier: string;
}

export interface ChampionRankStats {
  winRate: number;
  pickRate: number;
  performance: number;
}

export interface ChampionRegionStats {
  winRate: number;
  pickRate: number;
  popularity: number;
}

export interface MetaBuildData {
  buildName: string;
  items: string[];
  runes: string[];
  winRate: number;
  pickRate: number;
  situations: string[];
}

export interface PlayStyleData {
  styleName: string;
  description: string;
  popularity: number;
  effectiveness: number;
  keyFeatures: string[];
}

export interface BanAnalysis {
  topBannedChampions: BanStatsData[];
  banStrategies: BanStrategyData[];
  roleTargeting: Record<string, number>;
  powerBans: PowerBanData[];
}

export interface BanStatsData {
  champion: string;
  banRate: number;
  threatLevel: number;
  banPriority: number;
  reasons: string[];
}

export interface BanStrategyData {
  strategy: string;
  description: string;
  usage: number;
  effectiveness: number;
  targetChampions: string[];
}

export interface PowerBanData {
  champion: string;
  impactOnWinRate: number;
  banValue: number;
  situational: boolean;
}

export interface PickAnalysis {
  blindPickSafe: ChampionPickData[];
  flexPicks: ChampionPickData[];
  counterPicks: CounterPickData[];
  firstPickPriority: ChampionPickData[];
  lastPickOptions: ChampionPickData[];
}

export interface ChampionPickData {
  champion: string;
  pickValue: number;
  versatility: number;
  winRate: number;
  safetyRating: number;
}

export interface CounterPickData {
  champion: string;
  counters: string[];
  counterValue: number;
  effectiveness: number;
}

export interface MetaShiftData {
  shiftType: string;
  description: string;
  catalyst: string;
  affectedChampions: string[];
  impact: number;
  timeline: string;
}

export interface PatchImpactAnalysis {
  patchNumber: string;
  releaseDate: string;
  overallImpact: number;
  championChanges: ChampionPatchChange[];
  itemChanges: ItemPatchChange[];
  systemChanges: SystemPatchChange[];
  metaShiftPrediction: string;
}

export interface ChampionPatchChange {
  champion: string;
  changeType: 'buff' | 'nerf' | 'rework' | 'adjustment';
  severity: number;
  predictedImpact: string;
  changes: string[];
}

export interface ItemPatchChange {
  item: string;
  changeType: string;
  impact: number;
  affectedChampions: string[];
  changes: string[];
}

export interface SystemPatchChange {
  system: string;
  changeType: string;
  impact: number;
  description: string;
  affectedAspects: string[];
}

export interface MetaPredictions {
  nextPatchPredictions: ChampionTierPrediction[];
  emergingChampions: EmergingPrediction[];
  decliningChampions: DecliningPrediction[];
  strategicShifts: StrategicShiftPrediction[];
  confidence: number;
  predictionAccuracy: number;
}

export interface ChampionTierPrediction {
  champion: string;
  currentTier: string;
  predictedTier: string;
  confidence: number;
  reasoningFactors: string[];
}

export interface EmergingPrediction {
  champion: string;
  role: string;
  catalysts: string[];
  timeline: string;
  confidence: number;
}

export interface DecliningPrediction {
  champion: string;
  role: string;
  reasons: string[];
  timeline: string;
  confidence: number;
}

export interface StrategicShiftPrediction {
  strategy: string;
  direction: 'emerging' | 'declining';
  drivers: string[];
  timeline: string;
  impact: number;
}

export interface MetaRecommendation {
  type: string;
  title: string;
  description: string;
  priority: 'high' | 'medium' | 'low';
  targetRank: string;
  champions: string[];
  strategies: string[];
  expected: string;
}

// Service Response Types
export interface MetaHistory {
  startPatch: string;
  endPatch: string;
  champion?: string;
  metric: string;
  dataPoints: Array<{
    patch: string;
    value: number;
    tier?: string;
  }>;
  trendAnalysis: {
    direction: 'rising' | 'stable' | 'declining';
    change: number;
    volatility: 'low' | 'medium' | 'high';
  };
}

// Chart data interfaces for meta analytics
export interface MetaBarChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    backgroundColor: string[];
    borderColor: string[];
    borderWidth: number;
  }>;
}

export interface MetaRadarChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    borderColor: string;
    backgroundColor: string;
    pointBackgroundColor: string;
    pointBorderColor: string;
    pointHoverBackgroundColor: string;
    pointHoverBorderColor: string;
  }>;
}

export interface MetaTrendChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    borderColor: string;
    backgroundColor: string;
    tension: number;
    fill: boolean;
  }>;
}

// Utility types for meta analysis
export type MetaTier = 'S+' | 'S' | 'A+' | 'A' | 'B+' | 'B' | 'C+' | 'C' | 'D';

export type TrendDirection = 'rising' | 'stable' | 'declining' | 'falling';

export type ChampionClass = 
  | 'Tank'
  | 'Fighter'
  | 'Assassin'
  | 'Mage'
  | 'Marksman'
  | 'Support'
  | 'Controller'
  | 'Specialist';

export type GamePhase = 'early' | 'mid' | 'late';

export type MetaPosition = 'TOP' | 'JUNGLE' | 'MID' | 'ADC' | 'SUPPORT';

export type PickPhase = 'first_pick' | 'early_pick' | 'mid_pick' | 'late_pick' | 'last_pick';

export type BanPhase = 'first_ban' | 'second_ban' | 'third_ban' | 'fourth_ban' | 'fifth_ban';

export type CompType = 
  | 'engage'
  | 'poke'
  | 'split_push'
  | 'team_fight'
  | 'pick'
  | 'siege'
  | 'protect'
  | 'dive';

// Meta search and filter types
export interface MetaSearchCriteria {
  minWinRate?: number;
  maxBanRate?: number;
  tier?: MetaTier;
  role?: MetaPosition;
  trend?: TrendDirection;
  playStyle?: string;
  difficulty?: 'easy' | 'medium' | 'hard';
}

export interface MetaFilterOptions {
  patches: string[];
  regions: string[];
  ranks: string[];
  roles: MetaPosition[];
  tiers: MetaTier[];
  trends: TrendDirection[];
}

// Advanced meta analysis types
export interface MetaSnapshot {
  patch: string;
  region: string;
  snapshotDate: string;
  topTierChampions: {
    sPlus: string[];
    sTier: string[];
  };
  mostBanned: string[];
  emergingPicks: string[];
  dominantStrategies: string[];
  gameLengthTrend: string;
  nextPatchPredictions: string[];
}

export interface RoleMeta {
  role: MetaPosition;
  patch: string;
  topChampions: ChampionMetaStats[];
  roleTrends: {
    averageWinRate: number;
    pickDiversity: number;
    banPressure: number;
    metaShifts: string[];
  };
  recommendedChampions: Array<{
    champion: string;
    tier: string;
    reason: string;
    difficulty: string;
  }>;
}

export interface PatchComparison {
  patch1: string;
  patch2: string;
  tierChanges: Array<{
    champion: string;
    patch1Tier: string;
    patch2Tier: string;
    tierChange: number;
    impact: 'major' | 'minor' | 'none';
  }>;
  metaShifts: Array<{
    category: string;
    changeDescription: string;
    affectedChampions: string[];
  }>;
  statisticalChanges: {
    averageGameLengthChange: number;
    pickDiversityChange: number;
    banRateChanges: Record<string, number>;
  };
}

export interface TeamCompMeta {
  patch: string;
  teamCompositions: Array<{
    compositionType: string;
    winRate: number;
    popularity: number;
    champions: {
      top: string[];
      jungle: string[];
      mid: string[];
      adc: string[];
      support: string[];
    };
    synergies: string[];
    counters: string[];
    gameLengthPreference: string;
  }>;
  metaCompositionTrends: {
    risingComps: string[];
    decliningComps: string[];
    stableComps: string[];
  };
}

export interface ItemMeta {
  patch: string;
  popularItems: Array<{
    itemId: number;
    itemName: string;
    pickRate: number;
    winRate: number;
    roles: string[];
    champions: string[];
    synergisticItems: string[];
  }>;
  emergingItems: Array<{
    itemName: string;
    pickRateChange: number;
    winRate: number;
    trendingChampions: string[];
  }>;
  itemBuildPaths: Array<{
    buildName: string;
    items: string[];
    winRate: number;
    pickRate: number;
    optimalChampions: string[];
  }>;
}

export interface RuneMeta {
  patch: string;
  popularRuneCombinations: Array<{
    primaryTree: string;
    keystone: string;
    secondaryTree: string;
    winRate: number;
    pickRate: number;
    champions: string[];
  }>;
  keystoneAnalysis: Array<{
    keystone: string;
    tree: string;
    winRate: number;
    pickRate: number;
    trending: 'rising' | 'stable' | 'declining';
    optimalChampions: string[];
  }>;
  runeSynergies: Array<{
    runeCombination: string[];
    synergyRating: number;
    useCases: string[];
  }>;
}