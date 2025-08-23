// Skill Progression Types for Herald.lol Frontend

export interface SkillProgressionAnalysis {
  id: string;
  summonerId: string;
  analysisType: 'overall' | 'champion' | 'role' | 'skill';
  timeRange: TimeRange;
  skillCategories: SkillCategoryProgress[];
  overallProgress: OverallProgressData;
  rankProgression: RankProgressionData;
  championMastery: ChampionMasteryProgress[];
  coreSkills: CoreSkillsAnalysis;
  learningCurve: LearningCurveData;
  milestones: SkillMilestone[];
  predictions: ProgressionPredictions;
  recommendations: ProgressionRecommendation[];
  confidence: number;
  createdAt: string;
}

export interface TimeRange {
  startDate: string;
  endDate: string;
  periodType: 'week' | 'month' | 'season' | 'year' | 'custom';
  periodCount: number;
}

export interface SkillCategoryProgress {
  category: 'mechanical' | 'tactical' | 'strategic' | 'mental' | 'champion_specific';
  currentRating: number; // 0-100
  previousRating: number;
  progressRate: number; // improvement per week
  trend: ProgressTrend;
  subcategories: SkillSubcategory[];
  benchmarks: SkillBenchmarks;
  improvementTips: string[];
  practiceAreas: PracticeArea[];
}

export interface ProgressTrend {
  direction: 'improving' | 'stable' | 'declining' | 'inconsistent';
  strength: number; // 0-100 how strong the trend is
  duration: string; // how long trend has been active
  consistency: number; // 0-100 how consistent the trend is
  recentChanges: TrendPoint[];
  predictedNext: number;
}

export interface TrendPoint {
  date: string;
  value: number;
  change: number;
  events: string[]; // patch changes, meta shifts, etc.
}

export interface SkillSubcategory {
  name: string;
  currentValue: number;
  previousValue: number;
  progress: number;
  weight: number; // importance to overall category
  metrics: SkillMetric[];
  targetValue: number;
  estimatedTime: string;
}

export interface SkillMetric {
  metricName: string;
  currentValue: number;
  targetValue: number;
  percentile: number; // compared to similar rank players
  improvement: number;
  unit: string;
  description: string;
}

export interface SkillBenchmarks {
  rankBenchmarks: Record<string, number>; // Iron->Challenger expected values
  roleAverage: number;
  globalAverage: number;
  topPercentile: number; // top 10%
  yourRank: string;
  nextRankTarget: number;
}

export interface PracticeArea {
  area: string;
  priority: 'high' | 'medium' | 'low';
  currentLevel: number;
  targetLevel: number;
  timeEstimate: string;
  practiceMethod: string[];
  difficulty: 'easy' | 'moderate' | 'hard';
  impactRating: number; // how much improvement this will give
}

export interface OverallProgressData {
  overallRating: number; // 0-100
  previousRating: number;
  progressVelocity: number; // points per week
  skillTier: string; // Bronze, Silver, Gold, etc.
  nextTierProgress: number; // 0-100
  strengthAreas: string[];
  weaknessAreas: string[];
  mostImproved: string[];
  needsWork: string[];
  learningEfficiency: number; // how quickly they improve
  motivationFactors: MotivationFactor[];
  learningStyle: LearningStyleAnalysis;
}

export interface MotivationFactor {
  factor: string;
  impact: 'positive' | 'negative' | 'neutral';
  strength: number; // 0-100
  suggestions: string[];
}

export interface LearningStyleAnalysis {
  primaryStyle: 'visual' | 'kinesthetic' | 'analytical' | 'social';
  secondaryStyle: string;
  learningSpeed: 'fast' | 'moderate' | 'steady' | 'slow';
  retentionRate: number; // how well they retain skills
  optimalMethods: string[];
  avoidMethods: string[];
}

export interface RankProgressionData {
  currentRank: string;
  currentLP: number;
  peakRank: string;
  startSeasonRank: string;
  rankHistory: RankHistoryPoint[];
  promotionAttempts: number;
  demotionRisk: number; // 0-100
  promotionChance: number; // 0-100
  rankStability: number; // how stable is current rank
  expectedRank: string; // based on current skill
  rankingFactors: RankingFactor[];
  mmrEstimate: MMREstimateData;
}

export interface RankHistoryPoint {
  date: string;
  rank: string;
  lp: number;
  change: number;
  matchResult: 'win' | 'loss';
  performance: number; // 0-100
}

export interface RankingFactor {
  factor: string;
  impact: number; // -100 to 100
  description: string;
  improvement: string;
}

export interface MMREstimateData {
  estimatedMMR: number;
  confidence: number;
  rankMMRRange: string;
  mmrTrend: 'increasing' | 'stable' | 'decreasing';
  gainLossPattern: string;
}

export interface ChampionMasteryProgress {
  champion: string;
  role: string;
  masteryLevel: number;
  masteryPoints: number;
  gamesPlayed: number;
  winRate: number;
  performanceRating: number; // 0-100
  progressionStage: 'learning' | 'improving' | 'mastering' | 'expert';
  skillAreas: ChampionSkillArea[];
  mechanics: ChampionMechanicsData;
  gameKnowledge: ChampionKnowledgeData;
  decisionMaking: DecisionMakingData;
  improvementRate: number; // per week
  timeInvested: string;
  nextMilestone: ChampionMilestone;
}

export interface ChampionSkillArea {
  skillName: string;
  currentRating: number;
  targetRating: number;
  difficulty: string;
  priority: string;
  practiceTime: string;
  confidence: number;
}

export interface ChampionMechanicsData {
  overall: number;
  combos: MechanicSkillData;
  positioning: MechanicSkillData;
  teamfighting: MechanicSkillData;
  laning: MechanicSkillData;
  skillShots: MechanicSkillData;
  animation: MechanicSkillData; // canceling, weaving
  itemization: MechanicSkillData;
  advancedTechniques: AdvancedTechnique[];
}

export interface MechanicSkillData {
  rating: number;
  consistency: number;
  improvement: number;
  examples: string[];
  tips: string[];
}

export interface AdvancedTechnique {
  technique: string;
  mastery: number; // 0-100
  difficulty: string;
  impact: number;
  tutorial: string;
}

export interface ChampionKnowledgeData {
  powerSpikes: number;
  matchups: number;
  itemBuilds: number;
  runeSelection: number;
  waveManagement: number;
  roaming: number;
  objectives: number;
}

export interface DecisionMakingData {
  lanePhase: number;
  midGame: number;
  lateGame: number;
  teamfights: number;
  soloPlays: number;
  riskManagement: number;
  adaptability: number;
}

export interface ChampionMilestone {
  milestone: string;
  description: string;
  requirements: string[];
  reward: string;
  estimatedTime: string;
  progress: number; // 0-100
}

export interface CoreSkillsAnalysis {
  mechanical: CoreSkillData;
  gameKnowledge: CoreSkillData;
  strategic: CoreSkillData;
  mental: CoreSkillData;
  communication: CoreSkillData;
  adaptability: CoreSkillData;
  leadership: CoreSkillData;
}

export interface CoreSkillData {
  rating: number;
  progress: number;
  trend: ProgressTrend;
  components: SkillComponent[];
  percentile: number;
  nextLevel: string;
  blockers: string[];
  catalysts: string[];
}

export interface SkillComponent {
  name: string;
  value: number;
  weight: number;
  improvement: number;
  target: number;
}

export interface LearningCurveData {
  curveType: 'linear' | 'exponential' | 'plateau' | 'sigmoid';
  learningPhase: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  plateau: PlateauAnalysis;
  breakthroughs: Breakthrough[];
  optimalPractice: OptimalPracticeData;
  efficiencyMetrics: EfficiencyMetrics;
}

export interface PlateauAnalysis {
  inPlateau: boolean;
  plateauDuration: string;
  plateauLevel: number;
  breakoutTips: string[];
  breakoutChance: number;
}

export interface Breakthrough {
  date: string;
  skillArea: string;
  impactSize: number;
  trigger: string;
  description: string;
  lessons: string[];
}

export interface OptimalPracticeData {
  hoursPerWeek: number;
  sessionLength: string;
  practiceRatio: PracticeRatio;
  focusAreas: string[];
  practiceSchedule: PracticeSession[];
}

export interface PracticeRatio {
  rankedGames: number; // %
  normalGames: number; // %
  practiceTools: number; // %
  vodReview: number; // %
  theoryCraft: number; // %
}

export interface PracticeSession {
  type: string;
  duration: string;
  focus: string[];
  goals: string[];
  frequency: string;
}

export interface EfficiencyMetrics {
  learningRate: number; // skill points per hour
  retentionRate: number; // how well skills are retained
  transferRate: number; // how well skills transfer between champions
  focusQuality: number; // quality of practice sessions
  improvementFactor: number; // multiplier on learning
}

export interface SkillMilestone {
  id: string;
  name: string;
  category: string;
  description: string;
  achieved: boolean;
  achievementDate?: string;
  progress: number; // 0-100
  requirements: Requirement[];
  reward: MilestoneReward;
  difficulty: 'easy' | 'moderate' | 'hard' | 'extreme';
  estimatedTime: string;
  tips: string[];
}

export interface Requirement {
  type: string;
  description: string;
  target: number;
  current: number;
  met: boolean;
}

export interface MilestoneReward {
  type: 'badge' | 'title' | 'unlock' | 'bonus';
  name: string;
  description: string;
  value: string;
}

export interface ProgressionPredictions {
  rankPrediction: RankPrediction;
  skillPredictions: SkillPrediction[];
  timeToGoals: TimeToGoal[];
  potentialAnalysis: PotentialAnalysis;
  scenarios: ProgressionScenario[];
}

export interface RankPrediction {
  predictedRank: string;
  confidence: number;
  timeFrame: string;
  keyFactors: string[];
  requirements: string[];
  likelihood: number; // 0-100
}

export interface SkillPrediction {
  skillArea: string;
  currentRating: number;
  predictedRating: number;
  timeFrame: string;
  confidence: number;
  assumptions: string[];
}

export interface TimeToGoal {
  goal: string;
  estimatedTime: string;
  confidence: number;
  milestones: string[];
  blockers: string[];
  accelerators: string[];
}

export interface PotentialAnalysis {
  overallPotential: number; // 0-100
  peakPrediction: string; // predicted peak rank
  limitingFactors: string[];
  strengthAreas: string[];
  untappedPotential: UntappedArea[];
  talentAssessment: TalentAssessment;
}

export interface UntappedArea {
  area: string;
  potential: number;
  difficulty: string;
  impact: number;
  timeFrame: string;
}

export interface TalentAssessment {
  naturalTalent: number;
  workEthic: number;
  learningSpeed: number;
  consistency: number;
  adaptability: number;
  competitiveDrive: number;
  talentProfile: 'prodigy' | 'grinder' | 'balanced' | 'late_bloomer';
  recommendations: string[];
}

export interface ProgressionScenario {
  scenarioName: string;
  description: string;
  probability: number;
  timeFrame: string;
  requirements: string[];
  expectedOutcome: string;
  keyActions: string[];
}

export interface ProgressionRecommendation {
  id: string;
  type: 'practice' | 'champion' | 'playstyle' | 'mindset';
  priority: 'critical' | 'high' | 'medium' | 'low';
  title: string;
  description: string;
  impactRating: number; // 0-100
  difficulty: 'easy' | 'moderate' | 'hard';
  timeCommitment: string;
  actionSteps: ActionStep[];
  success: SuccessMetrics;
  related: string[]; // other recommendation IDs
  status?: 'active' | 'completed' | 'dismissed' | 'paused';
  progress?: number;
}

export interface ActionStep {
  step: string;
  description: string;
  duration: string;
  resources: string[];
  completed: boolean;
}

export interface SuccessMetrics {
  measurableGoals: MeasurableGoal[];
  timeline: string;
  successRate: number;
}

export interface MeasurableGoal {
  metric: string;
  target: number;
  current: number;
  unit: string;
  deadline: string;
}

// Request Types
export interface SkillProgressionRequest {
  summonerId: string;
  analysisType?: 'overall' | 'champion' | 'role' | 'skill';
  timeRange?: {
    startDate?: string;
    endDate?: string;
    periodType?: 'week' | 'month' | 'season' | 'year' | 'custom';
    periodCount?: number;
  };
  focusAreas?: string[];
}

export interface SkillCategoryTrackingRequest {
  summonerId: string;
  category: string;
  rating: number;
  percentile?: number;
  improvement?: number;
  confidence?: number;
}

export interface RankChangeRequest {
  summonerId: string;
  newRank: string;
  newLP: number;
  change: number;
  matchResult: 'win' | 'loss';
  performance: number;
}

export interface PracticeSessionRequest {
  summonerId: string;
  sessionType: 'cs_drill' | 'mechanics' | 'vod_review' | 'theory' | 'ranked_practice';
  focusAreas: string[];
  goals: string[];
  duration?: number; // planned duration in minutes
}

export interface SkillGoalRequest {
  summonerId: string;
  goalType: 'rank' | 'skill_rating' | 'champion_mastery' | 'custom';
  target: string;
  priority?: 'high' | 'medium' | 'low';
  deadline?: string;
  description?: string;
  strategy?: string[];
}

export interface BreakthroughRequest {
  summonerId: string;
  skillArea: string;
  impactSize: number;
  trigger: string;
  description: string;
  lessons: string[];
  beforeRating: number;
  afterRating: number;
}

// Response Types
export interface ProgressionOverviewResponse {
  summonerId: string;
  overallProgress: OverallProgressData;
  topStrengths: string[];
  topWeaknesses: string[];
  recentMilestones: SkillMilestone[];
  nextGoals: TimeToGoal[];
  confidence: number;
  lastUpdated: string;
}

export interface SkillCategoriesResponse {
  summonerId: string;
  skillCategories: SkillCategoryProgress[];
  confidence: number;
}

export interface RankHistoryResponse {
  summonerId: string;
  season: string;
  gameMode: string;
  history: RankHistoryPoint[];
  totalGames: number;
}

export interface MilestonesResponse {
  summonerId: string;
  milestones: SkillMilestone[];
  status?: string;
  category?: string;
}

export interface RecommendationsResponse {
  summonerId: string;
  recommendations: ProgressionRecommendation[];
  totalCount: number;
  filters: {
    priority?: string;
    status: string;
    limit: number;
  };
}

export interface SkillBenchmarksResponse {
  skillArea: string;
  rank: string;
  role: string;
  benchmarks: Array<{
    metricName: string;
    expectedValue: number;
    minValue: number;
    maxValue: number;
    unit: string;
    sampleSize: number;
  }>;
}

export interface ProgressionTrendsResponse {
  summonerId: string;
  timeRange: number;
  trends: Array<{
    category: string;
    direction: string;
    strength: number;
    velocity: number;
  }>;
  overall: {
    direction: string;
    strength: number;
    velocity: number;
  };
}

// Chart and Visualization Data
export interface SkillRadarData {
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

export interface ProgressionLineChart {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    borderColor: string;
    backgroundColor: string;
    fill: boolean;
    tension: number;
  }>;
}

export interface RankProgressionChart {
  timeline: string[];
  ranks: string[];
  lpChanges: number[];
  performances: number[];
  winStreak: boolean[];
}

export interface LearningCurveChart {
  timePoints: string[];
  skillRatings: number[];
  trendLine: number[];
  plateauPeriods: Array<{
    start: string;
    end: string;
    level: number;
  }>;
  breakthroughs: Array<{
    date: string;
    impact: number;
    description: string;
  }>;
}

// Utility Types
export type SkillCategory = 'mechanical' | 'tactical' | 'strategic' | 'mental' | 'champion_specific';
export type LearningPhase = 'beginner' | 'intermediate' | 'advanced' | 'expert';
export type TrendDirection = 'improving' | 'stable' | 'declining' | 'inconsistent';
export type Priority = 'critical' | 'high' | 'medium' | 'low';
export type Difficulty = 'easy' | 'moderate' | 'hard' | 'extreme';
export type RecommendationType = 'practice' | 'champion' | 'playstyle' | 'mindset';
export type GoalType = 'rank' | 'skill_rating' | 'champion_mastery' | 'custom';
export type TalentProfile = 'prodigy' | 'grinder' | 'balanced' | 'late_bloomer';

// Filters and Options
export interface SkillProgressionFilters {
  categories?: SkillCategory[];
  timeRange?: TimeRange;
  minRating?: number;
  maxRating?: number;
  trendsOnly?: boolean;
  includeHistory?: boolean;
}

export interface MilestoneFilters {
  achieved?: boolean;
  category?: string;
  difficulty?: Difficulty[];
  progress?: {
    min: number;
    max: number;
  };
}

export interface RecommendationFilters {
  type?: RecommendationType[];
  priority?: Priority[];
  status?: string[];
  difficulty?: Difficulty[];
  impactRating?: {
    min: number;
    max: number;
  };
}

// Export namespace for easier imports
export namespace SkillProgression {
  export type Analysis = SkillProgressionAnalysis;
  export type Request = SkillProgressionRequest;
  export type Overview = ProgressionOverviewResponse;
  export type Categories = SkillCategoriesResponse;
  export type Recommendations = RecommendationsResponse;
  export type Milestones = MilestonesResponse;
  export type Trends = ProgressionTrendsResponse;
  export type Filters = SkillProgressionFilters;
}