// Coaching Types for Herald.lol Frontend

export interface CoachingInsight {
  id: string;
  summonerId: string;
  insightType: 'match_analysis' | 'skill_development' | 'strategic' | 'tactical' | 'mental';
  insightData: string; // JSON stored as text
  confidence: number;
  createdAt: string;
  updatedAt: string;
}

export interface CoachingTip {
  id: number;
  tipId: string;
  summonerId: string;
  category: 'mechanical' | 'tactical' | 'strategic' | 'mental' | 'champion_specific';
  type: 'quick_tip' | 'deep_insight' | 'warning' | 'opportunity';
  title: string;
  content: string;
  context: string; // JSON stored as text
  relevance: number; // 0-100 how relevant to player
  actionable: boolean;
  difficulty: 'easy' | 'moderate' | 'hard';
  expected: string; // JSON stored as text (ExpectedOutcome)
  related: string; // JSON array of related tip IDs
  status: 'active' | 'applied' | 'dismissed' | 'archived';
  createdAt: string;
  updatedAt: string;
}

export interface CoachingTipFeedback {
  id: number;
  tipId: number;
  summonerId: string;
  helpful: boolean;
  applied: boolean;
  effective: boolean;
  comments: string;
  rating: number; // 1-10
  createdAt: string;
  updatedAt: string;
}

export interface ImprovementPlan {
  id: string;
  summonerId: string;
  planType: 'skill_development' | 'rank_climb' | 'champion_mastery';
  title: string;
  description: string;
  duration: '4_weeks' | '8_weeks' | '12_weeks' | 'custom';
  status: 'active' | 'completed' | 'paused' | 'abandoned';
  progress: number; // 0-100
  planData: string; // JSON stored as text
  mainObjectives: string; // JSON array
  dailyRoutine: string; // JSON object
  weeklyGoals: string; // JSON array
  checkpoints: string; // JSON array
  successMetrics: string; // JSON array
  startedAt: string;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PracticeRoutine {
  id: string;
  summonerId: string;
  routineName: string;
  routineType: 'daily' | 'weekly' | 'custom';
  duration: string; // total time commitment
  skillFocus: string; // JSON array
  phases: string; // JSON array
  equipment: string; // JSON array
  progression: string; // JSON object
  alternatives: string; // JSON array
  effectiveness: number; // 0-100
  timesUsed: number;
  averageRating: number; // user ratings
  status: 'active' | 'archived' | 'favorite';
  createdAt: string;
  updatedAt: string;
}

export interface PracticeSession {
  id: number;
  routineId?: string;
  summonerId: string;
  sessionType: 'cs_drill' | 'mechanics' | 'vod_review' | 'theory' | 'custom';
  focusAreas: string; // JSON array
  duration: number; // actual duration in minutes
  plannedDuration: number; // planned duration
  quality: number; // 1-10 subjective rating
  goals: string; // JSON array
  achievements: string; // JSON array
  notes: string;
  improvementSeen: boolean;
  followUpNeeded: boolean;
  effectiveness: number; // 0-100 how effective was the session
  startedAt: string;
  completedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface TacticalAdvice {
  id: number;
  adviceId: string;
  summonerId: string;
  category: string; // laning, teamfighting, positioning, vision, etc.
  situation: string;
  problem: string;
  solution: string;
  reasoning: string;
  examples: string; // JSON array of PracticalExample
  difficulty: 'easy' | 'moderate' | 'hard';
  impact: number; // 0-100
  frequency: string; // how often this situation occurs
  urgency: 'critical' | 'high' | 'medium' | 'low';
  related: string; // JSON array of related advice IDs
  applied: boolean; // has user applied this advice
  helpful: boolean; // user feedback
  rating: number; // user rating 1-10
  createdAt: string;
  updatedAt: string;
}

export interface StrategicGuidance {
  id: number;
  guidanceId: string;
  summonerId: string;
  strategyType: 'macro' | 'draft' | 'adaptation' | 'win_conditions';
  title: string;
  overview: string;
  principles: string; // JSON array of StrategicPrinciple
  application: string; // JSON object StrategyApplication
  counters: string; // JSON array of StrategyCounter
  mastery: string; // JSON object MasteryProgression
  advanced: string; // JSON array of AdvancedConcept
  difficulty: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  relevance: number; // 0-100 how relevant to current player level
  mastered: boolean; // has player mastered this
  inProgress: boolean; // currently working on this
  createdAt: string;
  updatedAt: string;
}

export interface MentalCoachingPlan {
  id: number;
  summonerId: string;
  planType: 'tilt_management' | 'confidence_building' | 'focus_training' | 'stress_management';
  currentMentalState: string; // JSON object
  goals: string; // JSON array
  techniques: string; // JSON array
  dailyPractices: string; // JSON array
  tiltTriggers: string; // JSON array
  copingStrategies: string; // JSON array
  progressMetrics: string; // JSON array
  status: 'active' | 'completed' | 'paused';
  progress: number; // 0-100
  effectiveness: number; // user-reported effectiveness
  startedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface PerformanceGoal {
  id: number;
  summonerId: string;
  goalType: 'rank' | 'skill_metric' | 'champion_mastery' | 'habit';
  title: string;
  description: string;
  target: string; // target value
  current: string; // current value
  measurement: string; // how to measure progress
  timeline: string; // target completion date
  priority: 'critical' | 'high' | 'medium' | 'low';
  status: 'active' | 'completed' | 'paused' | 'failed';
  progress: number; // 0-100
  milestones: string; // JSON array
  strategies: string; // JSON array
  blockers: string; // JSON array
  support: string; // JSON array
  achieved: boolean;
  achievementDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CoachingSchedule {
  id: number;
  summonerId: string;
  scheduleType: 'daily' | 'weekly' | 'monthly';
  activities: string; // JSON array
  timeSlots: string; // JSON array
  reminders: string; // JSON array
  flexibility: string; // JSON object
  adherence: number; // 0-100 how well they follow schedule
  adjustments: string; // JSON array of recent adjustments
  status: 'active' | 'paused' | 'archived';
  createdAt: string;
  updatedAt: string;
}

export interface ProgressTracking {
  id: number;
  summonerId: string;
  trackingPeriod: 'daily' | 'weekly' | 'monthly';
  metrics: string; // JSON object
  improvements: string; // JSON array
  setbacks: string; // JSON array
  breakthroughs: string; // JSON array
  overallProgress: number; // 0-100
  motivationLevel: number; // 0-100
  engagementLevel: number; // 0-100
  satisfactionLevel: number; // 0-100
  notes: string;
  recordedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface MatchAnalysisInsight {
  id: number;
  summonerId: string;
  matchId: string;
  insightType: 'tactical' | 'strategic' | 'mechanical' | 'mental';
  category: string; // laning, teamfighting, decision_making, etc.
  title: string;
  description: string;
  timestamp: number; // game timestamp in seconds
  severity: 'critical' | 'major' | 'minor' | 'positive';
  impact: number; // -100 to 100
  actionable: boolean;
  advice: string;
  context: string; // JSON object with game context
  reviewed: boolean; // has player reviewed this
  applied: boolean; // has player applied the advice
  rating: number; // user rating of insight quality
  createdAt: string;
  updatedAt: string;
}

export interface ChampionCoachingTip {
  id: number;
  summonerId: string;
  champion: string;
  role: string;
  tipCategory: 'mechanics' | 'combos' | 'builds' | 'matchups' | 'positioning';
  title: string;
  content: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  mastery: number; // 0-100 current mastery of this tip
  priority: 'high' | 'medium' | 'low';
  practiced: boolean; // has player practiced this
  mastered: boolean; // has player mastered this
  examples: string; // JSON array
  resources: string; // JSON array
  related: string; // JSON array of related tip IDs
  createdAt: string;
  updatedAt: string;
}

// Structured Types for JSON fields
export interface ExpectedOutcome {
  improvement: string;
  timeframe: string;
  metrics: string[];
  confidence: number;
}

export interface PracticalExample {
  situation: string;
  action: string;
  result: string;
  notes?: string;
}

export interface StrategicPrinciple {
  name: string;
  description: string;
  application: string;
  importance: number;
}

export interface StrategyApplication {
  gamePhase: string[];
  conditions: string[];
  execution: string;
  adaptations: string[];
}

export interface StrategyCounter {
  counterStrategy: string;
  identification: string;
  response: string;
  difficulty: string;
}

export interface MasteryProgression {
  currentLevel: 'novice' | 'apprentice' | 'competent' | 'proficient' | 'expert';
  nextLevel: string;
  requirements: string[];
  timeEstimate: string;
}

export interface AdvancedConcept {
  concept: string;
  description: string;
  prerequisites: string[];
  difficulty: string;
  examples: string[];
}

// Request Types
export interface GenerateCoachingInsightsRequest {
  summonerId: string;
  analysisType?: 'comprehensive' | 'tactical' | 'strategic' | 'mental';
  recentMatches?: number; // how many recent matches to analyze
  focusAreas?: string[];
  includePersonalization?: boolean;
}

export interface CoachingTipRequest {
  summonerId: string;
  category?: 'mechanical' | 'tactical' | 'strategic' | 'mental' | 'champion_specific';
  difficulty?: 'easy' | 'moderate' | 'hard';
  limit?: number;
  activeOnly?: boolean;
}

export interface ImprovementPlanRequest {
  summonerId: string;
  planType: 'skill_development' | 'rank_climb' | 'champion_mastery';
  title: string;
  description?: string;
  duration: '4_weeks' | '8_weeks' | '12_weeks' | 'custom';
  mainObjectives: string[];
  customDuration?: number; // in weeks if duration is 'custom'
}

export interface PracticeRoutineRequest {
  summonerId: string;
  routineName: string;
  routineType: 'daily' | 'weekly' | 'custom';
  duration: string;
  skillFocus: string[];
  targetLevel?: 'beginner' | 'intermediate' | 'advanced';
}

export interface PracticeSessionRequest {
  summonerId: string;
  routineId?: string;
  sessionType: 'cs_drill' | 'mechanics' | 'vod_review' | 'theory' | 'custom';
  focusAreas: string[];
  plannedDuration: number; // in minutes
  goals?: string[];
}

export interface TacticalAdviceRequest {
  summonerId: string;
  category?: string;
  situation?: string;
  urgency?: 'critical' | 'high' | 'medium' | 'low';
  limit?: number;
}

export interface StrategicGuidanceRequest {
  summonerId: string;
  strategyType?: 'macro' | 'draft' | 'adaptation' | 'win_conditions';
  difficulty?: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  currentLevel?: string;
}

export interface MentalCoachingRequest {
  summonerId: string;
  planType: 'tilt_management' | 'confidence_building' | 'focus_training' | 'stress_management';
  currentIssues?: string[];
  goals?: string[];
  intensity?: 'light' | 'moderate' | 'intensive';
}

export interface PerformanceGoalRequest {
  summonerId: string;
  goalType: 'rank' | 'skill_metric' | 'champion_mastery' | 'habit';
  title: string;
  description?: string;
  target: string;
  timeline?: string;
  priority?: 'critical' | 'high' | 'medium' | 'low';
}

export interface CoachingScheduleRequest {
  summonerId: string;
  scheduleType: 'daily' | 'weekly' | 'monthly';
  availableTimeSlots: string[];
  preferredActivities: string[];
  flexibilityLevel?: 'low' | 'medium' | 'high';
}

export interface ProgressTrackingRequest {
  summonerId: string;
  trackingPeriod: 'daily' | 'weekly' | 'monthly';
  metrics: Record<string, number>;
  improvements?: string[];
  setbacks?: string[];
  breakthroughs?: string[];
  motivationLevel?: number;
  engagementLevel?: number;
  satisfactionLevel?: number;
  notes?: string;
}

export interface CoachingTipFeedbackRequest {
  tipId: number;
  helpful: boolean;
  applied?: boolean;
  effective?: boolean;
  rating?: number; // 1-10
  comments?: string;
}

export interface MatchAnalysisInsightRequest {
  summonerId: string;
  matchId: string;
  analysisType?: 'tactical' | 'strategic' | 'mechanical' | 'mental' | 'comprehensive';
  focusAreas?: string[];
}

export interface ChampionCoachingTipRequest {
  summonerId: string;
  champion: string;
  role: string;
  tipCategory?: 'mechanics' | 'combos' | 'builds' | 'matchups' | 'positioning';
  difficulty?: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  priority?: 'high' | 'medium' | 'low';
}

// Response Types
export interface CoachingInsightsResponse {
  summonerId: string;
  insights: CoachingInsight[];
  summary: {
    tacticalInsights: number;
    strategicInsights: number;
    mentalInsights: number;
    averageConfidence: number;
  };
  recommendations: string[];
  nextSteps: string[];
  generatedAt: string;
}

export interface CoachingTipsResponse {
  summonerId: string;
  tips: CoachingTip[];
  totalCount: number;
  filters: {
    category?: string;
    difficulty?: string;
    status: string;
    limit: number;
  };
}

export interface ImprovementPlansResponse {
  summonerId: string;
  plans: ImprovementPlan[];
  activePlans: number;
  completedPlans: number;
  totalProgress: number;
}

export interface PracticeRoutinesResponse {
  summonerId: string;
  routines: PracticeRoutine[];
  favoriteRoutines: PracticeRoutine[];
  totalHoursLogged: number;
  averageEffectiveness: number;
}

export interface PracticeSessionsResponse {
  summonerId: string;
  sessions: PracticeSession[];
  totalSessions: number;
  totalPracticeTime: number; // in minutes
  averageQuality: number;
  recentTrend: 'improving' | 'stable' | 'declining';
  statistics: {
    sessionsByType: Record<string, number>;
    averageDuration: number;
    completionRate: number;
    effectivenessScore: number;
  };
}

export interface TacticalAdviceResponse {
  summonerId: string;
  advice: TacticalAdvice[];
  priorityAdvice: TacticalAdvice[];
  appliedCount: number;
  helpfulCount: number;
  categories: string[];
}

export interface StrategicGuidanceResponse {
  summonerId: string;
  guidance: StrategicGuidance[];
  masteredStrategies: number;
  inProgressStrategies: number;
  recommendedNext: StrategicGuidance[];
  skillLevel: 'beginner' | 'intermediate' | 'advanced' | 'expert';
}

export interface MentalCoachingResponse {
  summonerId: string;
  plans: MentalCoachingPlan[];
  currentState: {
    overallMentalHealth: number; // 0-100
    tiltFrequency: string;
    confidenceLevel: number;
    stressLevel: number;
    focusQuality: number;
  };
  recommendations: string[];
  dailyPractices: string[];
}

export interface PerformanceGoalsResponse {
  summonerId: string;
  goals: PerformanceGoal[];
  activeGoals: number;
  completedGoals: number;
  overallProgress: number; // 0-100
  onTrackGoals: number;
  behindGoals: number;
  upcomingMilestones: string[];
}

export interface CoachingScheduleResponse {
  summonerId: string;
  schedule: CoachingSchedule;
  adherenceScore: number; // 0-100
  upcomingActivities: string[];
  missedSessions: number;
  suggestedAdjustments: string[];
}

export interface ProgressTrackingResponse {
  summonerId: string;
  tracking: ProgressTracking[];
  trends: {
    overallProgress: 'improving' | 'stable' | 'declining';
    motivation: 'increasing' | 'stable' | 'decreasing';
    engagement: 'increasing' | 'stable' | 'decreasing';
    satisfaction: 'increasing' | 'stable' | 'decreasing';
  };
  insights: string[];
  concernAreas: string[];
}

export interface MatchAnalysisInsightsResponse {
  summonerId: string;
  matchId: string;
  insights: MatchAnalysisInsight[];
  keyMoments: MatchAnalysisInsight[];
  improvementOpportunities: MatchAnalysisInsight[];
  positiveHighlights: MatchAnalysisInsight[];
  actionableItems: number;
  overallRating: number; // 0-100
}

export interface ChampionCoachingTipsResponse {
  summonerId: string;
  champion: string;
  role: string;
  tips: ChampionCoachingTip[];
  masteryLevel: number; // 0-100 overall mastery
  practiceRecommendations: string[];
  nextMilestones: string[];
  estimatedTimeToMastery: string;
}

export interface CoachingDashboardResponse {
  summonerId: string;
  overview: {
    activePlans: number;
    completedGoals: number;
    practiceHoursThisWeek: number;
    overallProgress: number;
    mentalHealthScore: number;
    coachingEffectiveness: number;
  };
  recentInsights: CoachingInsight[];
  priorityTips: CoachingTip[];
  upcomingMilestones: string[];
  recommendedActions: string[];
  coachingStreak: number; // consecutive days of following coaching
  lastUpdated: string;
}

// Coaching AI Assistant Types
export interface CoachingAIRequest {
  summonerId: string;
  message: string;
  context?: 'general' | 'tactical' | 'strategic' | 'mental' | 'champion_specific';
  includeMatchData?: boolean;
  includePersonalHistory?: boolean;
}

export interface CoachingAIResponse {
  response: string;
  confidence: number;
  sources: string[];
  actionableItems: string[];
  followUpQuestions: string[];
  relatedTopics: string[];
}

// Utility Types
export type CoachingCategory = 'mechanical' | 'tactical' | 'strategic' | 'mental' | 'champion_specific';
export type CoachingTipType = 'quick_tip' | 'deep_insight' | 'warning' | 'opportunity';
export type CoachingDifficulty = 'easy' | 'moderate' | 'hard';
export type CoachingPriority = 'critical' | 'high' | 'medium' | 'low';
export type CoachingStatus = 'active' | 'completed' | 'paused' | 'abandoned' | 'dismissed' | 'archived';
export type PlanType = 'skill_development' | 'rank_climb' | 'champion_mastery';
export type PlanDuration = '4_weeks' | '8_weeks' | '12_weeks' | 'custom';
export type SessionType = 'cs_drill' | 'mechanics' | 'vod_review' | 'theory' | 'custom';
export type RoutineType = 'daily' | 'weekly' | 'custom';
export type GoalType = 'rank' | 'skill_metric' | 'champion_mastery' | 'habit';
export type ScheduleType = 'daily' | 'weekly' | 'monthly';
export type TrackingPeriod = 'daily' | 'weekly' | 'monthly';
export type InsightType = 'tactical' | 'strategic' | 'mechanical' | 'mental';
export type InsightSeverity = 'critical' | 'major' | 'minor' | 'positive';
export type MentalCoachingType = 'tilt_management' | 'confidence_building' | 'focus_training' | 'stress_management';
export type StrategyType = 'macro' | 'draft' | 'adaptation' | 'win_conditions';
export type ChampionTipCategory = 'mechanics' | 'combos' | 'builds' | 'matchups' | 'positioning';
export type CoachingLevel = 'beginner' | 'intermediate' | 'advanced' | 'expert';

// Export namespace for easier imports
export namespace Coaching {
  export type Insight = CoachingInsight;
  export type Tip = CoachingTip;
  export type Plan = ImprovementPlan;
  export type Routine = PracticeRoutine;
  export type Session = PracticeSession;
  export type TacticalAdvice = TacticalAdvice;
  export type StrategicGuidance = StrategicGuidance;
  export type MentalPlan = MentalCoachingPlan;
  export type Goal = PerformanceGoal;
  export type Schedule = CoachingSchedule;
  export type Progress = ProgressTracking;
  export type MatchInsight = MatchAnalysisInsight;
  export type ChampionTip = ChampionCoachingTip;
  export type Dashboard = CoachingDashboardResponse;
  export type AIRequest = CoachingAIRequest;
  export type AIResponse = CoachingAIResponse;
  export type Category = CoachingCategory;
  export type Status = CoachingStatus;
}