// Improvement Recommendations Types for Herald.lol Frontend

export interface ImprovementRecommendation {
  id: string;
  summonerId: string;
  category: string;
  priority: 'critical' | 'high' | 'medium' | 'low';
  title: string;
  description: string;
  impactScore: number; // 0-100
  difficultyLevel: 'easy' | 'medium' | 'hard' | 'expert';
  timeToSeeResults: number; // days
  estimatedROI: number; // expected improvement %
  
  // Detailed Action Plan
  actionPlan: ImprovementActionPlan;
  
  // Progress Tracking
  progressTracking: ProgressTrackingData;
  
  // Context and Reasoning
  recommendationContext: RecommendationContext;
  
  // Metadata
  createdAt: string;
  updatedAt: string;
  validUntil: string;
  status: 'active' | 'completed' | 'dismissed' | 'expired';
}

export interface ImprovementActionPlan {
  primaryObjective: string;
  secondaryObjectives: string[];
  actionSteps: ActionStep[];
  practiceExercises: PracticeExercise[];
  resources: LearningResource[];
  milestones: ImprovementMilestone[];
  successMetrics: string[];
}

export interface ActionStep {
  stepNumber: number;
  title: string;
  description: string;
  duration: string; // e.g., "15 minutes daily"
  frequency: string; // e.g., "daily", "3x per week"
  prerequisites: string[];
  tools: string[]; // practice tool, replay analysis, etc.
}

export interface PracticeExercise {
  name: string;
  description: string;
  duration: number; // minutes
  difficulty: string;
  focus: string[]; // what skills this targets
  instructions: string[];
  variations: string[]; // different ways to practice
}

export interface LearningResource {
  type: 'video' | 'guide' | 'tool' | 'coach';
  title: string;
  url?: string;
  description: string;
  duration?: string;
  difficulty: string;
}

export interface ImprovementMilestone {
  milestoneNumber: number;
  title: string;
  description: string;
  targetMetrics: string[];
  timeFrame: string;
  rewardValue: number; // potential LP/rank improvement
}

export interface ProgressTrackingData {
  overallProgress: number; // 0-100
  completedSteps: number[];
  currentMilestone: number;
  milestoneProgress: MilestoneProgress[];
  weeklyProgress: WeeklyProgressData[];
  performanceImpact: PerformanceImpactData;
  lastProgressUpdate: string;
}

export interface MilestoneProgress {
  milestoneNumber: number;
  progress: number; // 0-100
  startedAt: string;
  completedAt?: string;
  currentMetrics: Record<string, number>;
  targetMetrics: Record<string, number>;
}

export interface WeeklyProgressData {
  week: number;
  startDate: string;
  practiceTimeMinutes: number;
  gamesPlayed: number;
  skillImprovement: Record<string, number>; // skill -> improvement delta
  rankProgress: number; // LP change
  consistencyScore: number; // 0-100
}

export interface PerformanceImpactData {
  baselineMetrics: Record<string, number>;
  currentMetrics: Record<string, number>;
  improvementDeltas: Record<string, number>;
  roiActual: number; // actual improvement %
  roiPredicted: number; // predicted improvement %
  confidenceInterval: number; // prediction accuracy
}

export interface RecommendationContext {
  triggeringFactors: string[]; // what led to this recommendation
  dataSources: string[]; // recent games, long-term trends, etc.
  analysisDepth: 'surface' | 'moderate' | 'deep';
  confidenceScore: number; // 0-100
  alternativeOptions: AlternativeRecommendation[];
  personalizationFactors: PersonalizationData;
}

export interface AlternativeRecommendation {
  title: string;
  description: string;
  impactScore: number;
  difficulty: string;
  reason: string; // why this wasn't chosen as primary
}

export interface PersonalizationData {
  playStyle: string;
  learningStyle: 'visual' | 'hands_on' | 'analytical';
  timeCommitment: 'casual' | 'moderate' | 'intensive';
  currentRank: string;
  mainRole: string;
  championPool: string[];
  weakestAreas: string[];
  strengthAreas: string[];
  recentPerformance: RecentPerfSummary;
  goals: string[];
  preferences: UserPreferences;
}

export interface RecentPerfSummary {
  games: number;
  winRate: number;
  averageKDA: number;
  consistencyRating: number;
  improvementTrend: 'improving' | 'stable' | 'declining';
  problemAreas: string[];
}

export interface UserPreferences {
  preferredDifficulty: string;
  focusAreas: string[];
  avoidanceAreas: string[];
  practiceStyle: 'structured' | 'flexible' | 'mixed';
  feedbackFrequency: 'daily' | 'weekly' | 'monthly';
}

// Player Analysis Types
export interface PlayerAnalysisResult {
  summonerId: string;
  analysisDate: string;
  overallRating: number; // 0-100
  skillBreakdown: Record<string, number>; // skill -> rating
  improvementPotential: Record<string, number>; // skill -> potential gain
  criticalWeaknesses: CriticalWeakness[];
  underutilizedStrengths: UnderutilizedStrength[];
  personalizationData: PersonalizationData;
  recentTrends: RecentTrendAnalysis;
  competitiveBenchmark: CompetitiveBenchmark;
}

export interface CriticalWeakness {
  area: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  impactOnWinRate: number; // estimated WR improvement if fixed
  frequency: number; // how often this issue occurs
  rootCauses: string[];
  quickWins: string[]; // easy improvements
}

export interface UnderutilizedStrength {
  strength: string;
  currentUsage: number; // 0-100
  optimalUsage: number; // 0-100
  leverageStrategy: string[];
  potentialGain: number; // estimated improvement
}

export interface RecentTrendAnalysis {
  trendPeriodDays: number;
  overallTrend: 'improving' | 'stable' | 'declining';
  skillTrends: Record<string, string>; // skill -> trend
  consistencyTrend: string;
  performanceVolatility: number;
  recentBreakthroughs: RecentBreakthrough[];
  recentStruggles: RecentStruggle[];
}

export interface RecentBreakthrough {
  area: string;
  description: string;
  impact: number;
  date: string;
  sustainability: 'high' | 'medium' | 'low';
}

export interface RecentStruggle {
  area: string;
  description: string;
  frequency: number; // how often it happens
  severity: string;
  pattern: 'situational' | 'consistent' | 'random';
  trend: 'worsening' | 'stable' | 'improving';
}

export interface CompetitiveBenchmark {
  rankTier: string;
  regionalPercentile: number;
  skillPercentiles: Record<string, number>; // skill -> percentile
  strongerThanPeers: string[];
  weakerThanPeers: string[];
  competitiveAdvantages: string[];
  competitiveDisadvantages: string[];
}

// Quick Wins Types
export interface QuickWin {
  title: string;
  description: string;
  expectedImpact: string;
  timeToImplement: string;
  difficulty: 'very_easy' | 'easy' | 'medium';
  roiScore: number; // 0-100
  instructions: string[];
  habitFormation?: HabitFormationGuide;
  timingGuide?: Record<string, string>;
}

export interface HabitFormationGuide {
  trigger: string;
  action: string;
  reward: string;
  tracking: string;
}

// Coaching Plan Types
export interface CoachingPlan {
  summonerId: string;
  planDetails: CoachingPlanDetails;
  weeklyBreakdown: WeeklyPlan[];
  dailyRoutine: DailyRoutine;
  progressTracking: ProgressTrackingSetup;
}

export interface CoachingPlanDetails {
  durationDays: number;
  intensityLevel: 'light' | 'moderate' | 'intensive';
  primaryGoal: string;
  secondaryGoals: string[];
  estimatedOutcome: string;
}

export interface WeeklyPlan {
  week: number;
  focusTheme: string;
  primarySkills: string[];
  dailyTasks: DailyTasks;
  milestone: WeeklyMilestone;
}

export interface DailyTasks {
  practiceTime: number; // minutes
  gamesMinimum: number;
  specificFocus: string[];
  successMetrics: string[];
}

export interface WeeklyMilestone {
  title: string;
  description: string;
  target: string;
}

export interface DailyRoutine {
  warmUp: RoutineSection;
  focusedPractice: RoutineSection;
  rankedGames: RankedGameRoutine;
  reviewSession: RoutineSection;
}

export interface RoutineSection {
  durationMinutes: number;
  activities: string[];
}

export interface RankedGameRoutine {
  minimumGames: number;
  maximumGames: number;
  focusMindset: string[];
}

export interface ProgressTrackingSetup {
  dailyMetrics: string[];
  weeklyReview: string[];
  successIndicators: string[];
}

// Improvement Insights Types
export interface ImprovementInsight {
  type: 'opportunity' | 'trend' | 'warning' | 'achievement';
  priority: 'critical' | 'high' | 'medium' | 'low';
  title: string;
  description: string;
  impact: 'low' | 'medium' | 'high';
  difficulty?: string;
  timeframe?: string;
  actionSteps?: string[];
  trend?: 'positive' | 'negative' | 'neutral';
  continueDoing?: string[];
  suggestedFocus?: string[];
}

export interface OverallTrajectory {
  direction: 'improving' | 'stable' | 'declining';
  consistency: 'stable' | 'volatile';
  keyFocusAreas: string[];
  strengthAreas: string[];
  predictedOutcome: string;
}

// Overall Progress Types
export interface OverallProgress {
  summonerId: string;
  overallProgress: ProgressSummary;
  recentAchievements: Achievement[];
  skillProgress: Record<string, SkillProgressData>;
  weeklySummary: WeeklySummary;
}

export interface ProgressSummary {
  improvementScore: number;
  activeRecommendations: number;
  completedRecommendations: number;
  totalRecommendations: number;
  streakDays: number;
  consistencyRating: number;
}

export interface Achievement {
  achievement: string;
  description: string;
  date: string;
  impact: string;
}

export interface SkillProgressData {
  baseline: number;
  current: number;
  target: number;
  progress: number; // percentage to target
  trend: 'improving' | 'stable' | 'declining';
}

export interface WeeklySummary {
  weekNumber: number;
  practiceMinutes: number;
  gamesPlayed: number;
  recommendationsWorkedOn: number;
  improvementAreasFocused: string[];
  winRateChange: number;
  rankProgress: number; // LP gained
}

// Service Request Types
export interface RecommendationOptions {
  focusCategory?: string; // mechanical, macro, mental, etc.
  difficultyFilter?: string; // easy, medium, hard, expert
  timeConstraint?: number; // max time per day in minutes
  maxRecommendations: number;
  includeAlternatives: boolean;
  priorityAreas?: string[];
}

export interface RecommendationResponse {
  summonerId: string;
  recommendations: ImprovementRecommendation[];
  options: RecommendationOptions;
  generatedAt: string;
}

export interface ProgressUpdateRequest {
  overallProgress: number;
  completedSteps: number[];
  currentMilestone: number;
  performanceMetrics: Record<string, number>;
  notes?: string;
}

// Chart Data Types for Visualization
export interface ImprovementChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    borderColor: string;
    backgroundColor: string;
    fill: boolean;
    tension?: number;
  }>;
}

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

export interface ProgressBarData {
  skill: string;
  current: number;
  target: number;
  baseline: number;
  progress: number;
}

// Utility Types
export type ImprovementCategory = 
  | 'mechanical_skill'
  | 'game_knowledge' 
  | 'map_awareness'
  | 'team_fighting'
  | 'laning'
  | 'objective_control'
  | 'vision_control'
  | 'positioning'
  | 'decision_making'
  | 'mental_resilience'
  | 'communication'
  | 'champion_mastery';

export type PriorityLevel = 'critical' | 'high' | 'medium' | 'low';
export type DifficultyLevel = 'very_easy' | 'easy' | 'medium' | 'hard' | 'expert';
export type ImpactLevel = 'low' | 'medium' | 'high';
export type TrendDirection = 'improving' | 'stable' | 'declining';
export type LearningStyle = 'visual' | 'hands_on' | 'analytical' | 'mixed';
export type TimeCommitment = 'casual' | 'moderate' | 'intensive';
export type PracticeStyle = 'structured' | 'flexible' | 'mixed';

// Advanced Types for Gamification
export interface ImprovementBadge {
  id: string;
  name: string;
  description: string;
  iconUrl: string;
  rarity: 'common' | 'rare' | 'epic' | 'legendary';
  category: string;
  earnedAt?: string;
  progress?: number; // 0-100 for badges in progress
  requirements: string[];
}

export interface ImprovementStreak {
  type: 'practice' | 'consistency' | 'improvement' | 'completion';
  currentStreak: number;
  longestStreak: number;
  streakValue: string; // what constitutes the streak
  lastActivityDate: string;
  nextMilestone: number;
}

export interface ImprovementLeaderboard {
  timeframe: 'daily' | 'weekly' | 'monthly' | 'all_time';
  category: string;
  userRank: number;
  totalParticipants: number;
  topEntries: LeaderboardEntry[];
  userEntry: LeaderboardEntry;
}

export interface LeaderboardEntry {
  rank: number;
  summonerName: string;
  value: number;
  trend: 'up' | 'down' | 'same';
  badge?: string;
}