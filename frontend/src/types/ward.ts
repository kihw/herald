// Ward Analytics Types for Herald.lol Frontend

export interface WardAnalysis {
  id: string;
  playerId: string;
  champion?: string;
  position?: string;
  timeRange: string;
  
  // Core Ward Metrics
  mapControlScore: number;
  territoryControlled: number;
  wardEfficiency: number;
  wardPlacementScore: number;
  counterWardingScore: number;
  visionDeniedScore: number;
  
  // Strategic Coverage
  strategicCoverage: StrategicCoverageData;
  
  // Placement Patterns
  placementPatterns: WardPlacementPatterns;
  optimalPlacements: Array<{
    location: string;
    frequency: number;
    effectiveness: number;
    coordinates: { x: number; y: number };
  }>;
  
  // Zone Control
  zoneControl: Record<string, ZoneControlData>;
  riverControl: RiverControlData;
  jungleControl: JungleControlData;
  
  // Ward Types Analysis
  yellowWardsAnalysis: WardTypeAnalysis;
  controlWardsAnalysis: WardTypeAnalysis;
  blueWardAnalysis: WardTypeAnalysis;
  
  // Timing Analysis
  placementTiming: PlacementTimingData;
  
  // Clearing Patterns
  wardClearingPatterns: WardClearingPatterns;
  
  // Objective Control
  objectiveSetup: ObjectiveSetupData;
  safetyProvided: SafetyProvidedData;
  
  // Optimization
  placementOptimization: WardOptimizationData;
  clearingOptimization: WardOptimizationData;
  
  // Performance Trends
  trendData: WardTrendPoint[];
  trendDirection: 'improving' | 'declining' | 'stable';
  trendConfidence: number;
  
  // Recommendations
  recommendations: Array<{
    type: 'strength' | 'improvement' | 'critical';
    title: string;
    description: string;
    priority: 'high' | 'medium' | 'low';
    expectedImpact: number;
  }>;
  
  // Benchmarking
  roleBenchmark: BenchmarkData;
  rankBenchmark: BenchmarkData;
  globalBenchmark: BenchmarkData;
  
  // Metadata
  generatedAt: string;
  lastUpdated: string;
}

export interface StrategicCoverageData {
  overallScore: number;
  dragonPitCoverage: number;
  baronPitCoverage: number;
  keyChokepoints: number;
  rotationPaths: number;
  escapeRoutes: number;
  engagePoints: number;
  safetyZones: number;
}

export interface WardPlacementPatterns {
  earlyGamePattern: PlacementPattern;
  midGamePattern: PlacementPattern;
  lateGamePattern: PlacementPattern;
  defensivePattern: PlacementPattern;
  offensivePattern: PlacementPattern;
  consistencyScore: number;
  adaptabilityScore: number;
}

export interface PlacementPattern {
  frequency: number;
  locations: Array<{
    zone: string;
    percentage: number;
    effectiveness: number;
  }>;
  timing: {
    averageDelay: number;
    optimalityScore: number;
  };
}

export interface ZoneControlData {
  controlScore: number;
  wardsPlaced: number;
  wardsKilled: number;
  visionTime: number;
  contestedTime: number;
  dominanceScore: number;
  strategicValue: number;
}

export interface RiverControlData extends ZoneControlData {
  scuttleControl: number;
  crossingPrevention: number;
  ganksDetected: number;
  rotationTracking: number;
}

export interface JungleControlData extends ZoneControlData {
  invadeDetection: number;
  counterJungling: number;
  objectiveApproaches: number;
  junglerTracking: number;
}

export interface WardTypeAnalysis {
  totalPlaced: number;
  totalKilled: number;
  averageLifespan: number;
  averageValue: number;
  efficiencyScore: number;
  optimalUsage: number;
  wastedPlacements: number;
  
  // Specific to ward type
  totalUsed?: number; // For blue trinket
  accuracyRate?: number; // For blue trinket
  informationGained?: number; // For blue trinket
}

export interface PlacementTimingData {
  earlyGameFrequency: number;
  midGameFrequency: number;
  lateGameFrequency: number;
  
  optimalTimings: Array<{
    gameTime: number;
    context: string;
    effectiveness: number;
  }>;
  
  timingConsistency: number;
  reactiveSpeed: number;
  proactiveScore: number;
}

export interface WardClearingPatterns {
  proactiveClearing: number;
  reactiveClearing: number;
  opportunisticClearing: number;
  clearingEfficiency: number;
  averageClearTime: number;
  priorityTargeting: number;
  teamClearing: number;
}

export interface ObjectiveSetupData {
  DragonSetupScore: number;
  BaronSetupScore: number;
  HeraldSetupScore: number;
  ElderSetupScore: number;
  
  averageSetupTime: number;
  setupConsistency: number;
  visionDensity: number;
  clearingPreparation: number;
  teamCoordination: number;
}

export interface SafetyProvidedData {
  overallSafetyScore: number;
  escapeRouteCoverage: number;
  flankProtection: number;
  engageWarning: number;
  roamingSafety: number;
  lanePhaseProtection: number;
}

export interface WardOptimizationData {
  currentQualityScore: number;
  expectedControlGain: number;
  expectedSafetyGain: number;
  expectedDenialGain: number;
  
  suggestions: Array<{
    area: string;
    reason: string;
    priority: number;
    expectedImprovement: number;
  }>;
  
  currentCounterWardingEfficiency?: number;
}

export interface WardTrendPoint {
  date: string;
  value: number;
  gamePhase?: string;
  matchId?: string;
  champion?: string;
  position?: string;
}

export interface BenchmarkData {
  percentile: number;
  averageScore: number;
  topPercentileScore: number;
  sampleSize: number;
  comparisonType: 'role' | 'rank' | 'global' | 'champion';
  filterValue: string;
}

export interface MapControlData {
  player_id: string;
  time_range: string;
  map_control_score: number;
  territory_controlled: number;
  strategic_coverage: StrategicCoverageData;
  zone_control: Record<string, ZoneControlData>;
  river_control: RiverControlData;
  jungle_control: JungleControlData;
  focus_data?: any;
}

// Ward Placement Pattern specific types
export interface WardPlacementPattern {
  earlyGamePattern: PlacementPattern;
  midGamePattern: PlacementPattern;
  lateGamePattern: PlacementPattern;
  defensivePattern: PlacementPattern;
  offensivePattern: PlacementPattern;
  consistencyScore: number;
  adaptabilityScore: number;
}

// Additional utility types for frontend components
export interface WardHeatmapPoint {
  x: number;
  y: number;
  intensity: number;
  wardType: 'YELLOW' | 'CONTROL' | 'BLUE_TRINKET';
  effectiveness: number;
  gamePhase: 'early' | 'mid' | 'late';
}

export interface WardComparisonData {
  player_id: string;
  player_name: string;
  map_control_score: number;
  ward_efficiency: number;
  wards_placed_per_game: number;
  wards_killed_per_game: number;
  counter_warding_score: number;
  percentile: number;
  rank?: number;
}

export interface WardCoachingRecommendation {
  category: 'placement' | 'clearing' | 'timing' | 'efficiency';
  title: string;
  description: string;
  priority: 'high' | 'medium' | 'low';
  difficulty: 'easy' | 'medium' | 'hard';
  expectedImpact: number;
  
  practiceExercises: Array<{
    name: string;
    description: string;
    duration: string;
  }>;
}

export interface WardLearningPath {
  step: number;
  title: string;
  description: string;
  estimatedTime: string;
  prerequisites: string[];
  
  progressMetrics: Array<{
    metric: string;
    currentValue: number;
    targetValue: number;
    timeframe: string;
  }>;
}

// Chart data interfaces for ward analytics
export interface WardRadarChartData {
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

export interface WardTrendChartData {
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

export interface WardTypeDistributionData {
  labels: string[];
  datasets: Array<{
    data: number[];
    backgroundColor: string[];
    hoverBackgroundColor: string[];
  }>;
}

// Ward efficiency rating types
export type WardEfficiencyRating = {
  label: 'Excellent' | 'Good' | 'Average' | 'Below Average' | 'Poor';
  color: string;
  threshold: number;
};

export type WardZoneType = 'river' | 'jungle' | 'dragon' | 'baron' | 'objectives' | 'lane' | 'neutral';

export type WardMetricType = 'wards_placed' | 'wards_killed' | 'map_control' | 'efficiency' | 'control' | 'coverage';

export type WardPatternType = 'placement' | 'timing' | 'optimization';

export type WardClearingType = 'proactive' | 'reactive' | 'opportunistic';

export type WardOptimizationType = 'placement' | 'clearing' | 'timing';

export type ObjectiveType = 'dragon' | 'baron' | 'herald' | 'elder';

export type WardGamePhase = 'early' | 'mid' | 'late';

export type WardType = 'YELLOW' | 'CONTROL' | 'BLUE_TRINKET';