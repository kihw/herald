// Predictive Analytics Types for Herald.lol Frontend

export interface PredictiveAnalysis {
  id: string;
  summonerId: string;
  analysisType: string;
  generatedAt: string;
  
  // Core Predictions
  performancePrediction: PerformancePredictionData;
  rankProgression: RankProgressionPrediction;
  skillDevelopment: SkillDevelopmentForecast;
  
  // Champion & Meta
  championRecommendations: ChampionRecommendationData[];
  metaAdaptation: MetaAdaptationForecast;
  
  // Team Analysis
  teamPerformance?: TeamPerformancePrediction;
  teamSynergy?: TeamSynergyAnalysis;
  
  // Career Analytics
  careerTrajectory: CareerTrajectoryForecast;
  playerPotential: PlayerPotentialAssessment;
  
  // Model Insights
  modelConfidence: ModelConfidenceData;
  actionableInsights: ActionableInsight[];
  
  // Metadata
  lastUpdated: string;
  validUntil: string;
  dataQuality: number;
}

export interface PerformancePredictionData {
  horizon: 'short_term' | 'medium_term' | 'long_term';
  
  // Next Game Predictions
  nextGameWinProbability: number;
  expectedKDA: {
    kills: number;
    deaths: number;
    assists: number;
    kdaRatio: number;
  };
  expectedPerformanceScore: number;
  
  // Short-term Forecasts (1-10 games)
  shortTermForecasts: {
    winRateEstimate: number;
    lpGainEstimate: number;
    performanceConsistency: number;
    improvementAreas: string[];
  };
  
  // Medium-term Forecasts (1-4 weeks)
  mediumTermForecasts: {
    skillRatingChange: number;
    rankProgressionLikelihood: number;
    masteryImprovements: Array<{
      skill: string;
      currentLevel: number;
      predictedLevel: number;
      timeToAchieve: number;
    }>;
  };
  
  // Long-term Forecasts (1-12 months)
  longTermForecasts: {
    peakRankPrediction: string;
    skillCeilingEstimate: number;
    competitivePotential: 'recreational' | 'semi_competitive' | 'competitive' | 'professional';
    careerMilestones: Array<{
      milestone: string;
      estimatedTimeframe: string;
      probability: number;
    }>;
  };
  
  // Performance Factors
  performanceFactors: {
    championPool: {
      current: string[];
      recommended: string[];
      diversity: number;
      effectiveness: number;
    };
    playstyleConsistency: number;
    adaptabilityScore: number;
    mentalResilience: number;
    learningRate: number;
  };
  
  confidence: number;
  assumptions: string[];
}

export interface RankProgressionPrediction {
  currentRank: string;
  currentLP: number;
  targetRank?: string;
  
  // Progression Timeline
  progressionScenarios: Array<{
    scenario: 'optimistic' | 'realistic' | 'conservative';
    timeToTarget: number; // days
    gamesRequired: number;
    winRateRequired: number;
    probability: number;
  }>;
  
  // Milestone Predictions
  nearTermMilestones: Array<{
    rank: string;
    estimatedDate: string;
    probability: number;
    requirements: {
      averageWinRate: number;
      gamesPerWeek: number;
      performanceConsistency: number;
    };
  }>;
  
  // Factors Influencing Progression
  progressionFactors: {
    currentTrend: 'climbing' | 'stable' | 'declining';
    winStreakPotential: number;
    tiltRecoveryRate: number;
    peakPerformanceFrequency: number;
    championMastery: number;
  };
  
  // Barriers and Accelerators
  barriers: Array<{
    factor: string;
    impact: 'low' | 'medium' | 'high';
    description: string;
    mitigation: string[];
  }>;
  
  accelerators: Array<{
    factor: string;
    impact: 'low' | 'medium' | 'high';
    description: string;
    activation: string[];
  }>;
  
  confidence: number;
  lastUpdated: string;
}

export interface SkillDevelopmentForecast {
  forecastPeriod: number; // days
  skillCategories: string[];
  
  // Individual Skill Forecasts
  skillForecasts: Array<{
    category: string;
    currentLevel: number;
    predictedLevel: number;
    improvementRate: number;
    learningCurve: 'linear' | 'exponential' | 'logarithmic' | 'plateau';
    
    // Milestones
    milestones: Array<{
      level: number;
      estimatedDate: string;
      requirements: string[];
      benefits: string[];
    }>;
    
    // Development Plan
    developmentPlan: {
      focusAreas: string[];
      practiceRecommendations: string[];
      expectedChallenges: string[];
      successMetrics: string[];
    };
  }>;
  
  // Overall Development Trajectory
  overallTrajectory: {
    currentSkillLevel: number;
    predictedSkillLevel: number;
    developmentRate: 'accelerating' | 'steady' | 'slowing';
    skillBalance: number; // how balanced skills are
    specialization: string[]; // areas of emerging specialization
  };
  
  // Learning Analytics
  learningAnalytics: {
    learningStyle: 'visual' | 'practical' | 'analytical' | 'social';
    optimalPracticeFrequency: number;
    retentionRate: number;
    transferLearningPotential: number;
  };
  
  confidence: number;
  recommendations: string[];
}

export interface ChampionRecommendationData {
  champion: string;
  role: string;
  recommendationType: 'meta_strong' | 'personal_fit' | 'skill_development' | 'counter_pick';
  
  // Recommendation Scores
  overallScore: number;
  fitScore: number;
  metaScore: number;
  learningScore: number;
  
  // Performance Predictions
  predictedWinRate: number;
  predictedKDA: number;
  predictedPerformanceScore: number;
  improvementPotential: number;
  
  // Learning Analysis
  learningDifficulty: 'easy' | 'medium' | 'hard' | 'expert';
  timeToMastery: number; // games
  skillTransferability: number;
  
  // Reasoning
  reasons: Array<{
    factor: string;
    importance: 'high' | 'medium' | 'low';
    description: string;
  }>;
  
  // Context
  bestSituations: string[];
  avoidSituations: string[];
  synergisticChampions: string[];
  counterChampions: string[];
  
  // Meta Context
  currentMetaRating: 'S+' | 'S' | 'A+' | 'A' | 'B+' | 'B' | 'C+' | 'C' | 'D';
  metaTrend: 'rising' | 'stable' | 'declining';
  patchStability: 'stable' | 'volatile';
  
  confidence: number;
}

export interface MetaAdaptationForecast {
  currentPatch: string;
  targetPatch: string;
  adaptationStrategy: 'proactive' | 'gradual' | 'reactive';
  
  // Adaptation Plan
  adaptationPlan: {
    championTransitions: Array<{
      from: string;
      to: string;
      reason: string;
      priority: 'high' | 'medium' | 'low';
      timeline: string;
    }>;
    
    playstyleAdjustments: Array<{
      aspect: string;
      currentApproach: string;
      recommendedApproach: string;
      difficulty: 'easy' | 'medium' | 'hard';
    }>;
    
    buildAdaptations: Array<{
      champion: string;
      currentBuild: string[];
      recommendedBuild: string[];
      situational: boolean;
    }>;
  };
  
  // Meta Impact Analysis
  metaImpact: {
    playerImpact: 'positive' | 'neutral' | 'negative';
    affectedChampions: string[];
    opportunityChampions: string[];
    riskChampions: string[];
  };
  
  // Adaptation Predictions
  adaptationPredictions: {
    adaptationSuccess: number;
    performanceImpact: number;
    timeToAdapt: number; // games
    difficultyRating: number;
  };
  
  confidence: number;
  recommendations: string[];
}

export interface TeamPerformancePrediction {
  teamId: string;
  predictionScope: 'next_game' | 'next_series' | 'tournament';
  
  // Team Performance Metrics
  predictedWinRate: number;
  expectedPerformanceScore: number;
  teamSynergyScore: number;
  strategicFlexibility: number;
  
  // Individual Predictions
  memberPredictions: Array<{
    summonerId: string;
    role: string;
    predictedPerformance: number;
    carryPotential: number;
    consistencyRating: number;
    pressureHandling: number;
  }>;
  
  // Strategic Analysis
  strategicPredictions: {
    preferredStyle: string;
    adaptabilityScore: number;
    championPoolDepth: number;
    draftFlexibility: number;
  };
  
  // Matchup Analysis
  matchupFactors: Array<{
    factor: string;
    advantage: number; // -100 to 100
    importance: number;
    description: string;
  }>;
  
  confidence: number;
}

export interface TeamSynergyAnalysis {
  teamId: string;
  analysisDepth: 'basic' | 'comprehensive' | 'detailed';
  
  // Synergy Scores
  overallSynergy: number;
  communicationSynergy: number;
  playstyleSynergy: number;
  championPoolSynergy: number;
  strategicSynergy: number;
  
  // Synergy Breakdown
  synergyBreakdown: Array<{
    category: string;
    score: number;
    strengths: string[];
    weaknesses: string[];
    recommendations: string[];
  }>;
  
  // Pair Analysis
  memberPairAnalysis: Array<{
    member1: string;
    member2: string;
    synergyScore: number;
    workingRelationship: 'excellent' | 'good' | 'average' | 'needs_work';
    strengths: string[];
    improvementAreas: string[];
  }>;
  
  // Team Dynamics
  teamDynamics: {
    leadership: string; // member ID
    shotCaller: string;
    carryPlayers: string[];
    supportPlayers: string[];
    flexibilityRating: number;
  };
  
  recommendations: Array<{
    type: 'composition' | 'communication' | 'strategy' | 'practice';
    priority: 'high' | 'medium' | 'low';
    description: string;
    expectedImprovement: number;
  }>;
}

export interface CareerTrajectoryForecast {
  summonerId: string;
  forecastHorizon: 'short_term' | 'medium_term' | 'long_term';
  careerGoals?: string;
  
  // Trajectory Scenarios
  trajectoryScenarios: Array<{
    scenario: 'conservative' | 'realistic' | 'optimistic';
    probability: number;
    milestones: Array<{
      achievement: string;
      timeframe: string;
      requirements: string[];
    }>;
    peakPerformance: {
      estimatedRank: string;
      timeframe: string;
      sustainability: number;
    };
  }>;
  
  // Career Phases
  careerPhases: Array<{
    phase: string;
    duration: string;
    focus: string[];
    expectedGrowth: number;
    challenges: string[];
    opportunities: string[];
  }>;
  
  // Development Path
  developmentPath: {
    currentStage: string;
    nextStage: string;
    transitionRequirements: string[];
    estimatedTransitionTime: number;
    successProbability: number;
  };
  
  // Competitive Potential
  competitivePotential: {
    level: 'casual' | 'enthusiast' | 'semi_pro' | 'professional' | 'elite';
    confidence: number;
    limitingFactors: string[];
    acceleratingFactors: string[];
    timelineToReach: string;
  };
  
  confidence: number;
  lastUpdated: string;
}

export interface PlayerPotentialAssessment {
  summonerId: string;
  assessmentType: 'quick' | 'comprehensive' | 'detailed';
  
  // Overall Potential
  overallPotential: {
    rating: number; // 0-100
    tier: 'limited' | 'moderate' | 'high' | 'exceptional' | 'elite';
    confidence: number;
  };
  
  // Potential Breakdown
  potentialBreakdown: {
    mechanicalSkill: PotentialCategory;
    gameKnowledge: PotentialCategory;
    decisionMaking: PotentialCategory;
    teamwork: PotentialCategory;
    adaptability: PotentialCategory;
    mentalToughness: PotentialCategory;
    learning: PotentialCategory;
    leadership: PotentialCategory;
  };
  
  // Ceiling Analysis
  ceilingAnalysis: {
    estimatedCeiling: string; // rank
    ceilingConfidence: number;
    timeToReachCeiling: number; // months
    ceilingLimiters: string[];
    ceilingEnablers: string[];
  };
  
  // Development Recommendations
  developmentRecommendations: Array<{
    area: string;
    currentLevel: number;
    potentialLevel: number;
    priority: 'critical' | 'high' | 'medium' | 'low';
    developmentPlan: string[];
    timeframe: string;
    difficulty: 'easy' | 'medium' | 'hard' | 'expert';
  }>;
  
  // Comparative Analysis
  comparativeAnalysis: {
    percentileRating: number;
    similarPlayerProfiles: string[];
    strengthsRelativeToRank: string[];
    weaknessesRelativeToRank: string[];
  };
  
  confidence: number;
}

export interface PotentialCategory {
  current: number;
  ceiling: number;
  improvementRate: number;
  timeToRealizePotential: number;
  keyLimiters: string[];
  keyEnablers: string[];
}

export interface ModelConfidenceData {
  overallConfidence: number;
  dataQuality: number;
  sampleSize: number;
  recencyWeight: number;
  
  // Model Metrics
  modelAccuracy: {
    performancePrediction: number;
    rankPrediction: number;
    skillDevelopment: number;
    championRecommendation: number;
  };
  
  // Uncertainty Factors
  uncertaintyFactors: Array<{
    factor: string;
    impact: 'low' | 'medium' | 'high';
    description: string;
  }>;
  
  // Data Limitations
  dataLimitations: string[];
  assumptions: string[];
  
  lastCalibration: string;
}

export interface ActionableInsight {
  type: 'improvement' | 'opportunity' | 'warning' | 'recommendation';
  priority: 'critical' | 'high' | 'medium' | 'low';
  category: string;
  
  title: string;
  description: string;
  impact: 'high' | 'medium' | 'low';
  difficulty: 'easy' | 'medium' | 'hard';
  timeframe: 'immediate' | 'short_term' | 'medium_term' | 'long_term';
  
  actionSteps: string[];
  successMetrics: string[];
  expectedOutcome: string;
  
  confidence: number;
  validUntil: string;
}

// Service Response Types
export interface PredictiveResponse<T> {
  data: T;
  metadata: {
    generatedAt: string;
    modelVersion: string;
    confidence: number;
    validUntil: string;
  };
  status: 'success' | 'partial' | 'limited';
  warnings?: string[];
}

// Chart data interfaces for predictive analytics
export interface PredictiveChartData {
  labels: string[];
  datasets: Array<{
    label: string;
    data: number[];
    borderColor: string;
    backgroundColor: string;
    fill: boolean;
    tension: number;
    confidence?: number[];
  }>;
}

export interface PredictionInterval {
  lower: number;
  prediction: number;
  upper: number;
  confidence: number;
}

// Utility types for predictive analytics
export type PredictionHorizon = 'short_term' | 'medium_term' | 'long_term';
export type PredictionType = 'performance' | 'rank' | 'skill' | 'meta' | 'team' | 'career';
export type ConfidenceLevel = 'low' | 'medium' | 'high' | 'very_high';
export type ImpactLevel = 'low' | 'medium' | 'high' | 'critical';
export type PriorityLevel = 'low' | 'medium' | 'high' | 'critical';
export type DifficultyLevel = 'easy' | 'medium' | 'hard' | 'expert';
export type TimeframeType = 'immediate' | 'short_term' | 'medium_term' | 'long_term';

// Advanced prediction types
export interface PredictionScenario {
  name: string;
  probability: number;
  outcomes: Array<{
    metric: string;
    value: number;
    confidence: number;
  }>;
  assumptions: string[];
  timeline: string;
}

export interface ModelPerformanceMetrics {
  accuracy: number;
  precision: number;
  recall: number;
  f1Score: number;
  lastEvaluated: string;
  evaluationSample: number;
}