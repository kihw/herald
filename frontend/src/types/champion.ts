// Champion Analytics Types for Herald.lol Frontend

export interface ChampionAnalysis {
  id: string;
  playerId: string;
  champion: string;
  position: string;
  timeRange: string;
  
  // Core Performance Metrics
  masteryLevel: number;
  masteryPoints: number;
  playRate: number;
  winRate: number;
  totalGames: number;
  recentForm: string;
  
  // Performance Scoring
  overallRating: number;
  mechanicsScore: number;
  gameKnowledgeScore: number;
  consistencyScore: number;
  adaptabilityScore: number;
  
  // Champion-Specific Metrics
  championStats: ChampionSpecificStats;
  powerSpikes: ChampionPowerSpike[];
  itemBuilds: ItemBuildAnalysis;
  skillOrder: SkillOrderAnalysis;
  runeOptimization: RuneOptimizationData;
  
  // Matchup Analysis
  matchupPerformance: MatchupAnalysisData;
  strengthMatchups: MatchupData[];
  weaknessMatchups: MatchupData[];
  
  // Game Phase Performance
  lanePhasePerformance: GamePhaseData;
  midGamePerformance: GamePhaseData;
  lateGamePerformance: GamePhaseData;
  
  // Team Fighting Analysis
  teamFightRole: string;
  teamFightRating: number;
  teamFightStats: TeamFightAnalysisData;
  
  // Comparative Analysis
  roleBenchmark: ChampionBenchmarkData;
  rankBenchmark: ChampionBenchmarkData;
  globalBenchmark: ChampionBenchmarkData;
  
  // Performance Trends
  trendData: ChampionTrendPoint[];
  trendDirection: 'improving' | 'declining' | 'stable';
  trendConfidence: number;
  
  // Strengths and Weaknesses
  coreStrengths: PerformanceInsight[];
  improvementAreas: PerformanceInsight[];
  
  // Recommendations
  playStyleRecommendations: PlayStyleRecommendation[];
  trainingRecommendations: TrainingRecommendation[];
  
  // Advanced Analytics
  carryPotential: number;
  clutchFactor: number;
  learningCurve: LearningCurveData;
  metaAlignment: number;
  
  // Metadata
  generatedAt: string;
  lastUpdated: string;
}

export interface ChampionSpecificStats {
  // Core Stats
  averageKDA: number;
  averageKills: number;
  averageDeaths: number;
  averageAssists: number;
  
  // Champion-Specific Metrics
  averageCS: number;
  csPerMinute: number;
  goldPerMinute: number;
  damagePerMinute: number;
  visionScore: number;
  
  // Game Impact
  killParticipation: number;
  damageShare: number;
  goldShare: number;
  
  // Efficiency Metrics
  goldEfficiency: number;
  damageEfficiency: number;
  objectiveControl: number;
  
  // Champion-Specific Abilities
  skillAccuracy: Record<string, number>;
  ultimateUsage: UltimateUsageStats;
  passiveUtilization: number;
}

export interface UltimateUsageStats {
  averageUsesPerGame: number;
  accuracyRate: number;
  impactScore: number;
  timingOptimality: number;
  teamFightUsage: number;
  clutchUsage: number;
}

export interface ChampionPowerSpike {
  level: number;
  itemThreshold: string;
  powerRating: number;
  winRateIncrease: number;
  optimalTiming: number; // Game time in seconds
  utilizationRate: number;
}

export interface ItemBuildAnalysis {
  mostSuccessfulBuild: ItemBuildPath[];
  adaptabilityScore: number;
  buildVariety: number;
  counterBuildRate: number;
  coreItemTiming: Record<string, number>;
  situationalItems: SituationalItemData[];
}

export interface ItemBuildPath {
  items: string[];
  winRate: number;
  playRate: number;
  averageTiming: number[]; // Time to complete each item
  situations: string[]; // When this build is optimal
}

export interface SituationalItemData {
  itemId: number;
  itemName: string;
  usageRate: number;
  successRate: number;
  situations: string[];
  triggers: string[];
}

export interface SkillOrderAnalysis {
  mostCommonOrder: string;
  optimalOrder: string;
  adaptabilityRate: number;
  skillMaxOrder: SkillMaxData[];
  situationalOrders: SituationalSkillData[];
}

export interface SkillMaxData {
  skill: string;
  maxOrder: number;
  winRate: number;
  situations: string[];
}

export interface SituationalSkillData {
  situation: string;
  skillOrder: string;
  usageRate: number;
  successRate: number;
}

export interface RuneOptimizationData {
  primaryTree: string;
  secondaryTree: string;
  keystoneOptimality: number;
  runeAdaptation: number;
  mostSuccessfulSetup: RuneSetupData;
  situationalSetups: SituationalRuneData[];
}

export interface RuneSetupData {
  primaryTree: string;
  primaryRunes: string[];
  secondaryTree: string;
  secondaryRunes: string[];
  statShards: string[];
  winRate: number;
  playRate: number;
}

export interface SituationalRuneData {
  situation: string;
  runeSetup: RuneSetupData;
  usageRate: number;
  successRate: number;
}

export interface MatchupAnalysisData {
  totalMatchups: number;
  favorableMatchups: number;
  unfavorableMatchups: number;
  matchupAdaptability: number;
  lanePhaseWinRate: number;
  scalingAdvantage: number;
}

export interface MatchupData {
  opponentChampion: string;
  gamesPlayed: number;
  winRate: number;
  lanePhasePerformance: number;
  scalingComparison: number;
  averageCSAdvantage: number;
  averageGoldAdvantage: number;
  keyStrategies?: string[];
  commonMistakes?: string[];
}

export interface GamePhaseData {
  phaseRating: number;
  phaseWinRate: number;
  averagePerformance: number;
  consistencyScore: number;
  keyMetrics: Record<string, number>;
  strengthAreas: string[];
  weaknessAreas: string[];
}

export interface TeamFightAnalysisData {
  participationRate: number;
  survivalRate: number;
  damageContribution: number;
  ccContribution: number;
  positioningScore: number;
  engageTiming: number;
  targetPriority: number;
  ultimateEfficiency: number;
}

export interface ChampionBenchmarkData {
  percentile: number;
  averageRating: number;
  topPercentileRating: number;
  sampleSize: number;
  comparisonType: 'role' | 'rank' | 'global' | 'champion';
  filterValue: string;
}

export interface ChampionTrendPoint {
  date: string;
  overallRating: number;
  winRate: number;
  kda: number;
  csPerMinute: number;
  damageShare: number;
  matchId?: string;
  gameLength?: number;
}

export interface PerformanceInsight {
  category: string;
  title: string;
  description: string;
  metricValue: number;
  benchmarkValue: number;
  confidence: number;
  impact: 'high' | 'medium' | 'low';
}

export interface PlayStyleRecommendation {
  category: string;
  title: string;
  description: string;
  priority: 'high' | 'medium' | 'low';
  difficulty: 'easy' | 'medium' | 'hard';
  expectedImprovement: number;
  keyFocus: string[];
}

export interface TrainingRecommendation {
  type: string;
  title: string;
  description: string;
  duration: string;
  frequency: string;
  skillsImproved: string[];
  expectedTimeline: string;
}

export interface LearningCurveData {
  currentStage: string;
  progressScore: number;
  masteryTrajectory: 'accelerating' | 'steady_improvement' | 'plateauing' | 'declining';
  plateauRisk: number;
  nextMilestone: string;
  estimatedTimeToMastery: number; // In games
}

// Service Response Types
export interface ChampionMasteryRanking {
  champion: string;
  masteryPoints: number;
  winRate: number;
  playRate: number;
  rating: number;
  rank: number;
}

export interface ChampionComparisonData {
  playerId: string;
  timeRange: string;
  champions: ChampionComparisonEntry[];
}

export interface ChampionComparisonEntry {
  champion: string;
  winRate: number;
  playRate: number;
  averageKDA: number;
  damagePerMinute: number;
  csPerMinute: number;
  overallRating: number;
}

export interface ChampionMatchupData {
  player_id: string;
  champion: string;
  time_range: string;
  matchup_performance: MatchupAnalysisData;
  strength_matchups: MatchupData[];
  weakness_matchups: MatchupData[];
  specific_matchup?: MatchupData | { message: string };
}

export interface ChampionBuildData {
  player_id: string;
  champion: string;
  time_range: string;
  item_builds: ItemBuildAnalysis;
  skill_order: SkillOrderAnalysis;
  rune_optimization: RuneOptimizationData;
  focus_data?: ItemBuildAnalysis | SkillOrderAnalysis | RuneOptimizationData;
}

export interface ChampionCoachingData {
  player_id: string;
  champion: string;
  time_range: string;
  current_performance: {
    overall_rating: number;
    mechanics_score: number;
    game_knowledge_score: number;
    consistency_score: number;
  };
  core_strengths: PerformanceInsight[];
  improvement_areas: PerformanceInsight[];
  playstyle_recommendations: PlayStyleRecommendation[];
  training_recommendations: TrainingRecommendation[];
  learning_curve: LearningCurveData;
}

// Chart data interfaces for champion analytics
export interface ChampionRadarChartData {
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

export interface ChampionTrendChartData {
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

export interface ChampionBarChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    backgroundColor: string[];
    borderColor: string[];
    borderWidth: number;
  }>;
}

// Champion mastery and skill level types
export type ChampionMasteryTier = 'Bronze' | 'Silver' | 'Gold' | 'Platinum' | 'Diamond' | 'Master';

export type ChampionSkillLevel = 'Beginner' | 'Intermediate' | 'Advanced' | 'Expert' | 'Master';

export type ChampionRole = 'TOP' | 'JUNGLE' | 'MID' | 'ADC' | 'SUPPORT';

export type GamePhase = 'early' | 'mid' | 'late';

export type TrendDirection = 'improving' | 'declining' | 'stable';

export type ImpactLevel = 'high' | 'medium' | 'low';

export type Priority = 'high' | 'medium' | 'low';

export type Difficulty = 'easy' | 'medium' | 'hard';

export type TeamFightRole = 
  | 'Primary Damage Dealer'
  | 'Secondary Damage Dealer'
  | 'Tank/Engage'
  | 'Peeler/Support'
  | 'Assassin/Flanker'
  | 'Utility/Control';

// Champion-specific ability types
export type ChampionAbility = 'Q' | 'W' | 'E' | 'R' | 'Passive';

export type SkillOrder = string; // e.g., "Q-E-W-Q-Q-R"

// Item and build types
export interface ChampionItemBuild {
  coreItems: string[];
  situationalItems: string[];
  startingItems: string[];
  boots: string;
  buildPath: string[];
}

// Rune types
export interface ChampionRunePage {
  primaryTree: string;
  keystone: string;
  primaryRunes: string[];
  secondaryTree: string;
  secondaryRunes: string[];
  statShards: string[];
}

// Meta and tier list types
export type MetaTier = 'S+' | 'S' | 'A+' | 'A' | 'B+' | 'B' | 'C+' | 'C' | 'D';

export interface ChampionMetaData {
  currentTier: MetaTier;
  tierHistory: Array<{
    patch: string;
    tier: MetaTier;
    reason?: string;
  }>;
  metaAlignment: number;
  popularityTrend: 'rising' | 'stable' | 'declining';
}

// Synergy and team composition types
export interface ChampionSynergy {
  champion: string;
  synergyScore: number;
  synergyType: 'engage' | 'peel' | 'damage' | 'utility' | 'scaling';
  description: string;
}

export interface TeamComposition {
  type: 'engage' | 'poke' | 'split_push' | 'team_fight' | 'pick' | 'siege';
  effectiveness: number;
  championRoles: Record<string, string>;
  winRate: number;
  gameCount: number;
}

// Learning and improvement types
export interface ChampionLearningResource {
  type: 'guide' | 'video' | 'vod' | 'tool';
  title: string;
  creator: string;
  skillLevel: ChampionSkillLevel;
  focusAreas: string[];
  rating: number;
  url: string;
  duration?: string;
}

export interface ChampionMilestone {
  milestone: string;
  description: string;
  currentProgress: number;
  targetValue: number;
  estimatedGames: number;
  keyMetrics: string[];
}

export interface ChampionPracticeExercise {
  name: string;
  description: string;
  category: string;
  difficulty: Difficulty;
  duration: string;
  frequency: string;
  skillsImproved: string[];
}

// Performance comparison types
export interface ChampionPerformanceComparison {
  metric: string;
  playerValue: number;
  roleAverage: number;
  rankAverage: number;
  globalAverage: number;
  percentiles: {
    role: number;
    rank: number;
    global: number;
  };
}