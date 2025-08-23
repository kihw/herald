// Match Prediction Types for Herald.lol Frontend

export interface MatchPrediction {
  id: string;
  predictionType: 'pre_game' | 'draft' | 'live';
  gameMode: string;
  
  // Teams
  blueTeam: TeamPredictionData;
  redTeam: TeamPredictionData;
  
  // Prediction Results
  winProbability: ProbabilityData;
  gameAnalysis: GameFlowPrediction;
  
  // Detailed Predictions
  playerPerformance: PlayerMatchPrediction[];
  teamFightAnalysis: TeamFightPredictions;
  objectiveControl: ObjectivePredictions;
  
  // Draft Analysis
  draftAnalysis?: DraftAnalysisData;
  
  // Meta Context
  metaContext: MetaPredictionContext;
  
  // Confidence and Validation
  predictionConfidence: PredictionConfidenceData;
  
  // Metadata
  createdAt: string;
  predictionValidUntil: string;
  actualResult?: MatchResult;
}

export interface TeamPredictionData {
  teamId: string;
  players: PlayerPredictionSummary[];
  
  // Team Strength Analysis
  overallStrength: number; // 0-100
  teamSynergy: number; // 0-100
  experienceLevel: number; // 0-100
  recentForm: number; // 0-100
  
  // Composition Analysis
  compositionType: string; // team_fight, split_push, poke, etc.
  compositionScore: number; // 0-100
  scalingCurve: TeamScalingData;
  
  // Strengths and Weaknesses
  keyStrengths: string[];
  keyWeaknesses: string[];
  winConditions: WinCondition[];
  
  // Predicted Performance
  predictedKDA: number;
  predictedGold: number;
  predictedDamage: number;
}

export interface PlayerPredictionSummary {
  summonerId: string;
  summonerName: string;
  role: string;
  champion: string;
  rank: string;
  skillRating: number; // 0-100
  recentPerformance: number; // 0-100
  championMastery: number; // 0-100
  roleEfficiency: number; // 0-100
}

export interface TeamScalingData {
  earlyGame: number; // 0-15 min strength
  midGame: number; // 15-30 min strength
  lateGame: number; // 30+ min strength
  powerSpikes: number[]; // minute marks of major power spikes
}

export interface WinCondition {
  condition: string; // e.g., "Control early game", "Scale to late game"
  probability: number; // 0-100 chance this condition leads to win
  requirements: string[]; // what needs to happen
  counters: string[]; // how enemy can prevent this
}

export interface ProbabilityData {
  blueWinProbability: number; // 0-100
  redWinProbability: number; // 0-100
  probabilityFactors: ProbabilityFactor[];
  confidenceInterval: number; // 0-100
  modelAccuracy: number; // 0-100 based on historical predictions
}

export interface ProbabilityFactor {
  factor: string; // e.g., "Team composition", "Player skill gap"
  impact: number; // -50 to +50 (percentage points)
  confidence: number; // 0-100
  description: string; // explanation of the factor
}

export interface GameFlowPrediction {
  predictedGameLength: number; // minutes
  gamePhaseAnalysis: GamePhaseData[];
  keyMoments: KeyMomentPrediction[];
  victoryScenarios: VictoryScenario[];
  riskFactors: RiskFactor[];
}

export interface GamePhaseData {
  phase: 'early' | 'mid' | 'late';
  duration: string; // time range
  blueAdvantage: number; // -100 to +100
  redAdvantage: number; // -100 to +100
  keyObjectives: string[];
  criticalEvents: string[];
}

export interface KeyMomentPrediction {
  timestamp: number; // minute mark
  event: string; // e.g., "First dragon", "Baron spawn"
  importance: number; // 0-100
  prediction: string; // what's likely to happen
  consequences: string; // impact on game state
}

export interface VictoryScenario {
  team: 'blue' | 'red';
  scenario: string; // description
  probability: number; // 0-100
  timeline: string; // when this might happen
  triggers: string[]; // events that enable this scenario
}

export interface RiskFactor {
  risk: string; // description of the risk
  team: 'blue' | 'red'; // which team is at risk
  severity: 'low' | 'medium' | 'high' | 'critical';
  probability: number; // 0-100 chance this risk occurs
  mitigation: string[]; // how to avoid/minimize the risk
}

export interface PlayerMatchPrediction {
  summonerId: string;
  summonerName: string;
  role: string;
  champion: string;
  
  // Performance Predictions
  predictedKDA: KDAPrediction;
  predictedCS: CSPrediction;
  predictedDamage: DamagePrediction;
  predictedVision: VisionPrediction;
  predictedGold: GoldPrediction;
  
  // Impact Predictions
  carryPotential: number; // 0-100
  teamFightImpact: number; // 0-100
  laningPerformance: number; // 0-100
  objectiveImpact: number; // 0-100
  
  // Matchup Analysis
  laningMatchup: MatchupAnalysis;
  counterThreats: ThreatAnalysis[];
  synergyPartners: SynergyAnalysis[];
  
  // Confidence
  predictionConfidence: number; // 0-100
}

export interface KDAPrediction {
  kills: NumberRange;
  deaths: NumberRange;
  assists: NumberRange;
  kdaRange: NumberRange;
}

export interface CSPrediction {
  totalCS: IntRange;
  csPerMinute: NumberRange;
  csAt15Min: IntRange;
}

export interface DamagePrediction {
  totalDamage: IntRange;
  damagePerMinute: IntRange;
  damageShare: NumberRange; // % of team damage
  damageToChamps: IntRange;
}

export interface VisionPrediction {
  visionScore: NumberRange;
  wardsPlaced: IntRange;
  wardsDestroyed: IntRange;
  visionDenied: IntRange;
}

export interface GoldPrediction {
  totalGold: IntRange;
  goldPerMinute: IntRange;
  goldAt15Min: IntRange;
  goldEfficiency: NumberRange;
}

export interface NumberRange {
  min: number;
  expected: number;
  max: number;
}

export interface IntRange {
  min: number;
  expected: number;
  max: number;
}

export interface MatchupAnalysis {
  opponent: string; // enemy champion
  matchupRating: 'favorable' | 'even' | 'unfavorable';
  advantageScore: number; // -100 to +100
  keyFactors: string[];
  playstyleTips: string[];
  powerSpikes: string[];
}

export interface ThreatAnalysis {
  threatChampion: string;
  threatLevel: 'low' | 'medium' | 'high' | 'extreme';
  threatType: string; // burst, sustain, crowd_control, etc.
  counters: string[]; // how to play against this threat
  itemCounters: string[]; // items that help against threat
}

export interface SynergyAnalysis {
  partnerChampion: string;
  synergyRating: 'excellent' | 'good' | 'average' | 'poor';
  synergyType: string; // engage, protect, combo, etc.
  comboPotential: number; // 0-100
  playAroundTips: string[];
}

export interface TeamFightPredictions {
  teamFightStrength: TeamFightComparison;
  engageOptions: EngageOption[];
  teamFightScenarios: TeamFightScenario[];
  positioningAnalysis: PositioningPrediction;
}

export interface TeamFightComparison {
  blueTeamStrength: number; // 0-100
  redTeamStrength: number; // 0-100
  blueFightAdvantage: number; // -100 to +100
  keyAdvantages: string[];
  keyWeaknesses: string[];
}

export interface EngageOption {
  team: 'blue' | 'red';
  engageMethod: string; // e.g., "Malphite ultimate"
  engageStrength: number; // 0-100
  successRate: number; // 0-100
  counterPlay: string[];
  optimalTiming: string[];
}

export interface TeamFightScenario {
  scenarioName: string; // e.g., "5v5 at Baron"
  blueWinChance: number; // 0-100
  redWinChance: number; // 0-100
  keyFactors: string[];
  optimalStrategy: string;
}

export interface PositioningPrediction {
  frontlineStrength: number; // 0-100
  backlineProtection: number; // 0-100
  flankPotential: number; // 0-100
  positioningTips: string[];
}

export interface ObjectivePredictions {
  dragonControl: ObjectiveControlData;
  baronControl: ObjectiveControlData;
  riftHeraldControl: ObjectiveControlData;
  towerControl: TowerControlPrediction;
}

export interface ObjectiveControlData {
  blueControlChance: number; // 0-100
  redControlChance: number; // 0-100
  contestedRate: number; // 0-100
  controlFactors: string[];
  optimalStrategy: string;
}

export interface TowerControlPrediction {
  earlySieging: number; // 0-100
  midGamePressure: number; // 0-100
  lateGamePush: number; // 0-100
  splitPushPotential: number; // 0-100
  siegingAdvantages: string[];
}

export interface DraftAnalysisData {
  draftPhase: string; // pick_ban, completed
  blueDraftRating: number; // 0-100
  redDraftRating: number; // 0-100
  draftAdvantage: number; // -100 to +100 (positive = blue advantage)
  banAnalysis: BanAnalysisItem[];
  pickAnalysis: PickAnalysisItem[];
  compositionFit: CompositionAnalysis;
  flexPickAdvantage: FlexPickData;
}

export interface BanAnalysisItem {
  bannedChampion: string;
  banEffectiveness: number; // 0-100
  targetPlayer: string; // which player this ban targets
  impactReason: string; // why this ban is effective
  alternatives: string[]; // other champions that could have been banned
}

export interface PickAnalysisItem {
  pickedChampion: string;
  pickStrength: number; // 0-100
  pickReasoning: string; // why this pick is good/bad
  metaFit: number; // 0-100
  counterPotential: number; // 0-100
  synergyRating: number; // 0-100
  alternativePicks: string[];
}

export interface CompositionAnalysis {
  blueCompType: string; // team_fight, split_push, poke, etc.
  redCompType: string;
  blueCompStrength: number; // 0-100
  redCompStrength: number; // 0-100
  compMatchup: string; // how compositions interact
  winConditions: string[];
}

export interface FlexPickData {
  hasFlexPicks: boolean;
  flexAdvantage: number; // 0-100
  flexChampions: string[];
  flexStrategy: string; // how flex picks are used
}

export interface MetaPredictionContext {
  currentPatch: string;
  metaRelevance: number; // 0-100
  championTiers: ChampionTierContext[];
  metaTrends: MetaTrendContext[];
  patchImpact: PatchImpactData;
}

export interface ChampionTierContext {
  champion: string;
  tier: string; // S+, S, A+, A, B+, B, C+, C, D
  winRate: number; // 0-100
  pickRate: number; // 0-100
  banRate: number; // 0-100
  metaImpact: number; // 0-100
}

export interface MetaTrendContext {
  trendType: 'rising' | 'stable' | 'declining';
  champions: string[];
  impact: number; // 0-100
  description: string;
}

export interface PatchImpactData {
  patchAge: number; // days since patch
  stabilityRating: number; // 0-100
  majorChanges: string[];
  affectedChampions: string[];
}

export interface PredictionConfidenceData {
  overallConfidence: number; // 0-100
  dataQuality: number; // 0-100
  sampleSize: number;
  modelAccuracy: ModelAccuracyData;
  uncertaintyFactors: UncertaintyFactor[];
  confidenceBreakdown: ConfidenceBreakdown;
}

export interface ModelAccuracyData {
  winPredictionAccuracy: number; // 0-100
  playerPredictionAccuracy: number; // 0-100
  gameFlowAccuracy: number; // 0-100
  lastCalibration: string;
}

export interface UncertaintyFactor {
  factor: string; // e.g., "New patch", "Limited data"
  impact: 'low' | 'medium' | 'high';
  description: string;
  mitigation: string; // how uncertainty is handled
}

export interface ConfidenceBreakdown {
  winProbabilityConfidence: number; // 0-100
  playerPerformanceConfidence: number; // 0-100
  gameFlowConfidence: number; // 0-100
  draftAnalysisConfidence: number; // 0-100
  objectiveConfidence: number; // 0-100
}

export interface MatchResult {
  winningTeam: 'blue' | 'red';
  gameLength: number; // minutes
  actualPlayerData: ActualPlayerPerformance[];
  validationScore: number; // 0-100 how accurate prediction was
  resultDate: string;
}

export interface ActualPlayerPerformance {
  summonerId: string;
  actualKDA: number;
  actualCS: number;
  actualDamage: number;
  actualVision: number;
  actualGold: number;
  performanceRating: number; // 0-100
}

// Service Request Types
export interface MatchPredictionRequest {
  predictionType: 'pre_game' | 'draft' | 'live';
  gameMode: string;
  blueTeam: PlayerMatchData[];
  redTeam: PlayerMatchData[];
  draftData?: DraftData;
  metaData?: MetaMatchData;
}

export interface PlayerMatchData {
  summonerId: string;
  summonerName: string;
  role: string;
  champion?: string; // may not be selected yet
  rank: string;
  recentGames: number; // games to analyze
}

export interface DraftData {
  phase: string; // pick_ban, completed
  blueBans: string[];
  redBans: string[];
  bluePicks: string[];
  redPicks: string[];
  pickOrder: string[]; // order of picks/bans
  timeRemaining: number; // seconds left in draft
}

export interface MetaMatchData {
  patch: string;
  region: string;
  rankTier: string;
  gameType: string;
  customOptions: Record<string, string>;
}

// Pre-game Analysis Types
export interface PreGameAnalysis {
  blueTeamStrength: number;
  redTeamStrength: number;
  predictedOutcome: {
    blueWinProbability: number;
    redWinProbability: number;
    confidence: number;
  };
  keyFactors: KeyFactor[];
  teamAnalysis: {
    blueTeam: TeamAnalysisSummary;
    redTeam: TeamAnalysisSummary;
  };
  recommendations: {
    blueTeam: string[];
    redTeam: string[];
  };
}

export interface KeyFactor {
  factor: string;
  impact: string;
  description: string;
}

export interface TeamAnalysisSummary {
  strengths: string[];
  weaknesses: string[];
  winConditions: string[];
}

// Live Match Prediction Types
export interface LiveMatchPrediction {
  summonerId: string;
  inGame: boolean;
  gameInfo: {
    gameId: string;
    gameMode: string;
    gameStartTime: string;
    gameDuration: number; // seconds
  };
  currentPrediction: {
    winProbability: {
      blueTeam: number;
      redTeam: number;
    };
    gameState: {
      phase: string;
      predictedLength: number;
      blueTeamGoldAdvantage: number;
      blueTeamKillAdvantage: number;
    };
  };
  playerPerformance: {
    currentKDA: string;
    predictedFinalKDA: string;
    cs: number;
    predictedFinalCS: number;
    gold: number;
    performanceRating: number;
  };
  keyEvents: LiveGameEvent[];
  nextPredictions: {
    nextDragon?: ObjectiveSpawn;
    nextMajorEvent: string;
  };
}

export interface LiveGameEvent {
  timestamp: number;
  event: string;
  impact: string;
}

export interface ObjectiveSpawn {
  spawnTime: number;
  controlProbability: {
    blueTeam: number;
    redTeam: number;
  };
}

// Chart Data Types for Visualization
export interface PredictionChartData {
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

export interface WinProbabilityChart {
  labels: string[];
  blueTeamData: number[];
  redTeamData: number[];
  confidence: number[];
}

export interface GameFlowChart {
  timeline: number[];
  blueAdvantage: number[];
  redAdvantage: number[];
  keyEvents: Array<{
    time: number;
    event: string;
    impact: number;
  }>;
}

// Utility Types
export type PredictionType = 'pre_game' | 'draft' | 'live';
export type GameMode = 'ranked' | 'normal' | 'tournament' | 'custom';
export type TeamSide = 'blue' | 'red';
export type GamePhase = 'early' | 'mid' | 'late';
export type MatchupRating = 'favorable' | 'even' | 'unfavorable';
export type ThreatLevel = 'low' | 'medium' | 'high' | 'extreme';
export type SynergyRating = 'excellent' | 'good' | 'average' | 'poor';
export type ConfidenceLevel = 'low' | 'medium' | 'high' | 'very_high';

// Historical Analysis Types
export interface PredictionHistory {
  summonerId: string;
  predictions: HistoricalPrediction[];
  summary: {
    totalPredictions: number;
    averageAccuracy: number;
    winPredictionAccuracy: number;
    lossPredictionAccuracy: number;
    mostAccurateRole: string;
    leastAccurateRole: string;
  };
  trends: {
    accuracyTrend: 'improving' | 'stable' | 'declining';
    recentAccuracy: number;
    bestPredictionStreak: number;
    currentStreak: number;
  };
}

export interface HistoricalPrediction {
  predictionId: string;
  createdAt: string;
  gameMode: string;
  predictedWinProbability: number;
  actualResult: 'win' | 'loss';
  accuracy: number;
  champion: string;
  role: string;
  gameLength: number;
}

// Model Performance Types
export interface ModelPerformance {
  modelInfo: {
    version: string;
    lastUpdated: string;
    trainingDataSize: number;
    featuresCount: number;
  };
  performanceMetrics: {
    accuracy: number;
    precision: number;
    recall: number;
    f1Score: number;
    aucRoc: number;
    logLoss: number;
  };
  predictionDistribution: {
    highConfidence: { count: number; accuracy: number };
    mediumConfidence: { count: number; accuracy: number };
    lowConfidence: { count: number; accuracy: number };
  };
  featureImportance: Array<{
    feature: string;
    importance: number;
  }>;
  validationResults: {
    crossValidationScore: number;
    validationSetAccuracy: number;
    overfittingScore: number;
    generalizationScore: number;
  };
  recentPerformance: {
    last7Days: {
      predictions: number;
      accuracy: number;
      trend: string;
    };
    last30Days: {
      predictions: number;
      accuracy: number;
      trend: string;
    };
  };
}