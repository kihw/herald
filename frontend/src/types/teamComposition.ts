// Team Composition Types for Herald.lol Frontend

export interface TeamCompositionOptimization {
  id: string;
  strategy: OptimizationStrategy;
  timestamp: string;
  
  // Optimization Results
  recommendations: CompositionRecommendation[];
  analysis: CompositionAnalysis;
  alternatives: AlternativeComposition[];
  
  // Performance Metrics
  optimizationScore: number; // 0-100
  confidence: number; // 0-100
  expectedWinRate: number; // 0-100
  
  // Context
  gameMode: string;
  patch: string;
  metaContext: MetaCompositionContext;
}

export interface CompositionRecommendation {
  id: string;
  name: string;
  description: string;
  
  // Team Setup
  composition: ChampionRole[];
  
  // Scoring
  overallScore: number; // 0-100
  synergyScore: number; // 0-100
  metaScore: number; // 0-100
  comfortScore: number; // 0-100
  balanceScore: number; // 0-100
  
  // Analysis
  strengths: string[];
  weaknesses: string[];
  winConditions: WinCondition[];
  powerSpikes: PowerSpike[];
  
  // Performance Predictions
  expectedPerformance: PerformancePrediction;
  scalingProfile: ScalingProfile;
  teamFightRating: TeamFightRating;
  
  // Player Comfort
  playerComfort: PlayerComfortData[];
  
  // Strategic Analysis
  playStyle: PlayStyleAnalysis;
  objectives: ObjectiveControl;
  
  confidence: number; // 0-100
}

export interface ChampionRole {
  champion: string;
  role: string;
  flexibility: number; // 0-100 (can play multiple roles)
  priority: 'must_pick' | 'preferred' | 'situational' | 'backup';
  reasoning: string[];
}

export interface WinCondition {
  condition: string;
  probability: number; // 0-100
  timeline: string;
  requirements: string[];
  counters: string[];
  keyPlayers: string[];
}

export interface PowerSpike {
  phase: 'early' | 'mid' | 'late';
  timing: string; // e.g., "Level 6", "2 items", "20 minutes"
  strength: number; // 0-100
  description: string;
  champions: string[]; // which champions contribute
}

export interface PerformancePrediction {
  earlyGame: number; // 0-100 predicted strength
  midGame: number;
  lateGame: number;
  teamFighting: number;
  splitPushing: number;
  sieging: number;
  pickPotential: number;
  objectiveControl: number;
}

export interface ScalingProfile {
  overall: 'early' | 'mid' | 'late' | 'balanced';
  earlyGamePower: number; // 0-100
  midGamePower: number;
  lateGamePower: number;
  scalingCurve: ScalingPoint[];
  criticalTimings: CriticalTiming[];
}

export interface ScalingPoint {
  minute: number;
  powerLevel: number; // 0-100
  keyFactors: string[];
}

export interface CriticalTiming {
  timing: string;
  importance: 'critical' | 'important' | 'moderate';
  description: string;
  impact: string;
}

export interface TeamFightRating {
  overall: number; // 0-100
  initiation: number;
  followUp: number;
  protection: number;
  damage: number;
  disruption: number;
  positioning: PositioningRequirements;
  idealScenario: string;
}

export interface PositioningRequirements {
  frontline: string[];
  backline: string[];
  flankers: string[];
  requirements: string[];
  tips: string[];
}

export interface PlayerComfortData {
  summonerId: string;
  summonerName: string;
  role: string;
  champion: string;
  comfortLevel: number; // 0-100
  masteryPoints: number;
  recentGames: number;
  winRate: number;
  performance: number; // 0-100
  confidence: number; // 0-100
}

export interface PlayStyleAnalysis {
  primary: PlayStyleType;
  secondary?: PlayStyleType;
  description: string;
  keyTactics: string[];
  idealGamePlan: string[];
  adaptations: PlayStyleAdaptation[];
}

export type PlayStyleType = 
  | 'team_fight'
  | 'split_push'
  | 'poke_siege'
  | 'pick_comp'
  | 'protect_adc'
  | 'engage_comp'
  | 'scaling_comp'
  | 'early_game'
  | 'control_comp'
  | 'dive_comp';

export interface PlayStyleAdaptation {
  situation: string;
  adaptation: string;
  priority: number; // 1-10
}

export interface ObjectiveControl {
  dragonControl: number; // 0-100
  baronControl: number;
  riftHeraldControl: number;
  towerSiege: number;
  visionControl: number;
  jungleControl: number;
  overallControl: number;
  keyStrengths: string[];
  weaknesses: string[];
}

export interface CompositionAnalysis {
  id: string;
  
  // Overall Assessment
  overallRating: number; // 0-100
  tierRating: string; // S+, S, A+, A, B+, B, C
  metaFit: number; // 0-100
  
  // Component Analysis
  roleBalance: RoleBalanceAnalysis;
  damageProfile: DamageProfile;
  defenseProfile: DefenseProfile;
  utilityProfile: UtilityProfile;
  
  // Team Dynamics
  synergy: SynergyAnalysis;
  conflicts: ConflictAnalysis[];
  flexibility: FlexibilityAnalysis;
  
  // Strategic Assessment
  winConditions: WinCondition[];
  threats: ThreatAnalysis[];
  counters: CounterAnalysis[];
  
  // Performance Projections
  expectedPerformance: PerformancePrediction;
  riskFactors: RiskFactor[];
  
  // Recommendations
  improvements: ImprovementSuggestion[];
  alternatives: string[];
}

export interface RoleBalanceAnalysis {
  tankiness: number; // 0-100
  damage: number;
  utility: number;
  engage: number;
  peel: number;
  waveClear: number;
  balance: 'excellent' | 'good' | 'acceptable' | 'poor';
  gaps: string[];
  overlaps: string[];
}

export interface DamageProfile {
  physicalDamage: number; // 0-100
  magicDamage: number;
  trueDamage: number;
  burstDamage: number;
  sustainedDamage: number;
  aoeCapability: number;
  balance: 'excellent' | 'good' | 'predictable' | 'poor';
  vulnerabilities: string[];
}

export interface DefenseProfile {
  tankiness: number; // 0-100
  mobility: number;
  healing: number;
  shielding: number;
  crowdControl: number;
  disengage: number;
  survivability: 'excellent' | 'good' | 'moderate' | 'fragile';
  weaknesses: string[];
}

export interface UtilityProfile {
  vision: number; // 0-100
  engage: number;
  disengage: number;
  buffs: number;
  debuffs: number;
  zoning: number;
  overall: 'excellent' | 'good' | 'moderate' | 'limited';
  keyUtilities: string[];
}

export interface SynergyAnalysis {
  overall: number; // 0-100
  comboPotential: ComboPotential[];
  chainSynergies: ChainSynergy[];
  antiSynergies: AntiSynergy[];
  improvement: number; // 0-100 potential for improvement
}

export interface ComboPotential {
  champions: string[];
  comboName: string;
  damage: number; // 0-100
  reliability: number; // 0-100
  cooldown: number; // seconds
  range: string;
  description: string;
}

export interface ChainSynergy {
  type: 'engage' | 'follow_up' | 'protection' | 'amplification';
  champions: string[];
  description: string;
  effectiveness: number; // 0-100
}

export interface AntiSynergy {
  champions: string[];
  issue: string;
  severity: 'minor' | 'moderate' | 'major';
  mitigation: string[];
}

export interface ConflictAnalysis {
  type: 'resource' | 'positioning' | 'timing' | 'role_overlap';
  champions: string[];
  description: string;
  impact: 'low' | 'medium' | 'high';
  resolution: string[];
}

export interface FlexibilityAnalysis {
  overall: number; // 0-100
  roleFlexibility: RoleFlexibility[];
  strategicFlexibility: number; // 0-100
  adaptability: number; // 0-100
  pickBanFlexibility: number; // 0-100
}

export interface RoleFlexibility {
  champion: string;
  primaryRole: string;
  alternativeRoles: string[];
  effectiveness: Record<string, number>; // role -> effectiveness (0-100)
}

export interface ThreatAnalysis {
  threat: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  probability: number; // 0-100
  impact: string;
  mitigation: string[];
  responsibleChampions: string[];
}

export interface CounterAnalysis {
  vulnerableTo: VulnerabilityAnalysis[];
  strongAgainst: StrengthAnalysis[];
  neutralMatchups: string[];
  overallThreatLevel: number; // 0-100
}

export interface VulnerabilityAnalysis {
  counterChampions: string[];
  counterStrategy: string;
  severity: 'minor' | 'moderate' | 'major' | 'critical';
  mitigation: string[];
  probability: number; // 0-100 chance of facing this counter
}

export interface StrengthAnalysis {
  targetChampions: string[];
  advantage: string;
  exploitMethod: string[];
  consistency: number; // 0-100
}

export interface RiskFactor {
  factor: string;
  probability: number; // 0-100
  impact: 'low' | 'medium' | 'high' | 'game_ending';
  mitigation: string[];
  warning_signs: string[];
}

export interface ImprovementSuggestion {
  type: 'champion_swap' | 'strategic_adjustment' | 'playstyle_change';
  suggestion: string;
  reasoning: string;
  impact: number; // 0-100 expected improvement
  difficulty: 'easy' | 'moderate' | 'hard';
  alternatives: string[];
}

export interface AlternativeComposition {
  name: string;
  composition: ChampionRole[];
  score: number; // 0-100
  tradeOffs: TradeOffAnalysis;
  reason: string;
  situationalUse: string;
}

export interface TradeOffAnalysis {
  gains: string[];
  losses: string[];
  netBenefit: number; // -100 to 100
}

export interface MetaCompositionContext {
  patch: string;
  metaRelevance: number; // 0-100
  trendDirection: 'rising' | 'stable' | 'declining';
  popularityRank: number;
  proPlayUsage: number; // 0-100
  soloQueueSuccess: number; // 0-100
  banRate: number; // 0-100
  counterRate: number; // 0-100 how often it gets countered
}

// Optimization Request Types
export interface CompositionOptimizationRequest {
  playerData: PlayerOptimizationData[];
  strategy: OptimizationStrategy;
  constraints: OptimizationConstraints;
  preferences: OptimizationPreferences;
  gameMode: string;
  bannedChampions?: string[];
  requiredChampions?: string[];
}

export interface PlayerOptimizationData {
  summonerId: string;
  summonerName: string;
  role: string;
  championPool: string[];
  comfortLevel: number; // 1-10
  recentGames: number;
  preferredPlaystyle?: PlayStyleType[];
  avoidChampions?: string[];
}

export type OptimizationStrategy = 
  | 'meta_optimal'    // Focus on current meta strength
  | 'synergy_focused' // Maximize team synergy
  | 'balanced'        // Balance all factors
  | 'comfort_picks'   // Prioritize player comfort
  | 'counter_focused' // Focus on countering enemies
  | 'scaling_focused' // Focus on late game
  | 'early_focused';  // Focus on early game

export interface OptimizationConstraints {
  maxNewChampions: number;
  requireADC: boolean;
  requireTank: boolean;
  requireEngage: boolean;
  preferLateGame: boolean;
  preferEarlyGame: boolean;
  mustIncludeRoles: string[];
  maxRiskLevel: 'low' | 'medium' | 'high';
  minSynergyScore: number; // 0-100
}

export interface OptimizationPreferences {
  prioritizeMeta: number; // 0-10
  prioritizeSynergy: number;
  prioritizeComfort: number;
  prioritizeBalance: number;
  prioritizeFlexibility: number;
  avoidHighBanRate: boolean;
  preferProvenCombos: boolean;
  allowExperimental: boolean;
}

// Draft Optimization Types
export interface DraftOptimization {
  currentDraft: DraftState;
  recommendations: DraftRecommendation[];
  analysis: DraftAnalysis;
  predictions: DraftPrediction[];
}

export interface DraftState {
  bluePicks: string[];
  redPicks: string[];
  blueBans: string[];
  redBans: string[];
  currentTurn: 'blue_pick' | 'red_pick' | 'blue_ban' | 'red_ban';
  pickOrder: string[];
  phase: string;
  timeRemaining: number;
}

export interface DraftRecommendation {
  action: 'pick' | 'ban';
  champion: string;
  priority: number; // 1-10
  reasoning: string[];
  alternatives: string[];
  risk: 'low' | 'medium' | 'high';
  impact: number; // 0-100
  situationalFactors: string[];
}

export interface DraftAnalysis {
  blueTeamAdvantage: number; // -100 to 100
  redTeamAdvantage: number;
  keyAdvantages: string[];
  keyRisks: string[];
  flexPickAdvantage: number; // 0-100
  banEffectiveness: number; // 0-100
  counterPickOpportunities: CounterPickOpportunity[];
}

export interface CounterPickOpportunity {
  targetChampion: string;
  counterOptions: string[];
  effectiveness: number; // 0-100
  role: string;
  priority: number; // 1-10
}

export interface DraftPrediction {
  scenario: string;
  probability: number; // 0-100
  outcome: string;
  preparation: string[];
}

// Meta Composition Types
export interface MetaComposition {
  id: string;
  name: string;
  composition: ChampionRole[];
  
  // Meta Statistics
  winRate: number; // 0-100
  pickRate: number; // 0-100
  banRate: number; // 0-100
  popularity: number; // 0-100
  
  // Performance by Rank
  performanceByRank: Record<string, number>;
  
  // Trend Data
  trend: 'rising' | 'stable' | 'declining';
  trendStrength: number; // 0-100
  
  // Professional Play
  proPlayPresence: number; // 0-100
  proWinRate: number;
  
  // Analysis
  strengths: string[];
  weaknesses: string[];
  counters: string[];
  synergies: string[];
  
  // Context
  patch: string;
  region: string;
  sampleSize: number;
  lastUpdated: string;
}

// Player Comfort and Champion Pool Types
export interface PlayerComfortData {
  summonerId: string;
  role: string;
  
  // Comfort Tiers
  masterTier: ChampionComfort[];      // S-tier comfort
  comfortTier: ChampionComfort[];     // A-tier comfort
  familiarTier: ChampionComfort[];    // B-tier comfort
  learningTier: ChampionComfort[];    // C-tier comfort
  
  // Analysis
  poolDepth: number; // 0-100
  metaAlignment: number; // 0-100
  flexibility: number; // 0-100
  
  // Recommendations
  recommendations: ChampionRecommendation[];
  warnings: string[];
}

export interface ChampionComfort {
  champion: string;
  comfortScore: number; // 0-100
  masteryPoints: number;
  recentGames: number;
  winRate: number;
  averagePerformance: number; // 0-100
  consistency: number; // 0-100
  lastPlayed: string;
  trends: PerformanceTrend;
}

export interface PerformanceTrend {
  direction: 'improving' | 'stable' | 'declining';
  strength: number; // 0-100
  recentWinRate: number;
  recentPerformance: number;
}

export interface ChampionRecommendation {
  champion: string;
  reason: string;
  priority: number; // 1-10
  learningDifficulty: 'easy' | 'moderate' | 'hard';
  metaRelevance: number; // 0-100
  synergyWithPool: number; // 0-100
  expectedTimeToComfort: string;
}

// Ban Strategy Types
export interface BanStrategy {
  recommendations: BanRecommendation[];
  strategy: BanStrategyType;
  analysis: BanAnalysis;
  alternatives: AlternativeBanStrategy[];
}

export interface BanRecommendation {
  champion: string;
  priority: number; // 1-10
  reasoning: string[];
  target: 'specific_player' | 'meta_deny' | 'comp_protection';
  effectiveness: number; // 0-100
  risk: 'low' | 'medium' | 'high';
  alternatives: string[];
}

export type BanStrategyType = 
  | 'target_player'    // Target specific enemy player
  | 'protect_comp'     // Protect our composition
  | 'meta_deny'        // Deny meta champions
  | 'flex_deny'        // Deny flex picks
  | 'balanced';        // Balanced approach

export interface BanAnalysis {
  strategyEffectiveness: number; // 0-100
  riskAssessment: string[];
  expectedImpact: string[];
  counterBanPotential: number; // 0-100
}

export interface AlternativeBanStrategy {
  strategy: BanStrategyType;
  champions: string[];
  reasoning: string;
  effectiveness: number; // 0-100
  situationalUse: string;
}

// Utility Types
export type CompositionTier = 'S+' | 'S' | 'A+' | 'A' | 'B+' | 'B' | 'C+' | 'C' | 'D';
export type GamePhase = 'early' | 'mid' | 'late';
export type RoleType = 'top' | 'jungle' | 'mid' | 'adc' | 'support';
export type TeamSide = 'blue' | 'red';

// Chart Data Types for Visualization
export interface CompositionChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    borderColor: string;
    backgroundColor: string;
    fill?: boolean;
    tension?: number;
  }>;
}

export interface SynergyRadarData {
  labels: string[];
  data: number[];
  backgroundColor: string;
  borderColor: string;
  pointBackgroundColor: string;
}

export interface ScalingCurveData {
  timeline: number[];
  teamPower: number[];
  enemyPower?: number[];
  powerSpikes: Array<{
    time: number;
    power: number;
    description: string;
  }>;
}

export interface WinConditionTreeData {
  nodes: Array<{
    id: string;
    label: string;
    probability: number;
    children: string[];
    requirements: string[];
  }>;
  connections: Array<{
    from: string;
    to: string;
    strength: number;
  }>;
}