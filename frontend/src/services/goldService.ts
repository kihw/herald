import { BaseService } from './baseService';

export interface GoldAnalysis {
  player_id: string;
  champion?: string;
  position?: string;
  time_range: string;
  average_gold_earned: number;
  average_gold_per_minute: number;
  gold_efficiency_score: number;
  economy_rating: string;
  gold_sources: GoldSourcesData;
  item_efficiency: ItemEfficiencyData;
  spending_patterns: SpendingPatternsData;
  early_game_gold: GoldPhaseData;
  mid_game_gold: GoldPhaseData;
  late_game_gold: GoldPhaseData;
  role_benchmark: GoldBenchmark;
  rank_benchmark: GoldBenchmark;
  global_benchmark: GoldBenchmark;
  gold_advantage_win_rate: number;
  gold_disadvantage_win_rate: number;
  gold_impact_score: number;
  trend_direction: string;
  trend_slope: number;
  trend_confidence: number;
  trend_data: GoldTrendPoint[];
  income_optimization: IncomeOptimizationData;
  spending_optimization: SpendingOptimizationData;
  strength_areas: string[];
  improvement_areas: string[];
  recommendations: GoldRecommendation[];
  recent_matches: MatchGoldData[];
}

export interface GoldSourcesData {
  farming_gold: number;
  farming_percent: number;
  kills_gold: number;
  kills_percent: number;
  assists_gold: number;
  assists_percent: number;
  objective_gold: number;
  objective_percent: number;
  passive_gold: number;
  passive_percent: number;
  items_gold: number;
  items_percent: number;
  cs_gold_per_minute: number;
  kill_gold_efficiency: number;
  objective_gold_share: number;
}

export interface ItemEfficiencyData {
  average_items_completed: number;
  item_completion_speed: number;
  gold_spent_on_items: number;
  item_value_efficiency: number;
  damage_items_percent: number;
  defensive_items_percent: number;
  utility_items_percent: number;
  first_item_timing: number;
  core_items_timing: number;
  six_items_timing: number;
  optimal_item_order: boolean;
  counter_build_efficiency: number;
  component_utilization: number;
}

export interface SpendingPatternsData {
  control_wards_percent: number;
  consumables_percent: number;
  back_timing: BackTimingData;
  gold_efficiency_by_phase: PhaseGoldEfficiency[];
  average_shopping_time: number;
  optimal_back_timing: number;
  emergency_backs: number;
}

export interface BackTimingData {
  average_back_timing: number;
  optimal_backs: number;
  suboptimal_backs: number;
  forced_backs: number;
  gold_per_back: number;
}

export interface PhaseGoldEfficiency {
  phase: string;
  gold_per_minute: number;
  spending_efficiency: number;
  income_efficiency: number;
  economy_score: number;
}

export interface GoldPhaseData {
  phase: string;
  average_gold_per_minute: number;
  gold_advantage: number;
  farming_efficiency: number;
  kill_participation: number;
  objective_participation: number;
  spending_score: number;
  efficiency_rating: string;
}

export interface GoldBenchmark {
  category: string;
  average_gold_per_minute: number;
  top_10_percent: number;
  top_25_percent: number;
  median: number;
  player_percentile: number;
  efficiency_average: number;
}

export interface IncomeOptimizationData {
  cs_improvement_potential: number;
  kp_improvement_potential: number;
  objective_improvement_potential: number;
  early_farming_suggestions: string[];
  mid_game_position_suggestions: string[];
  late_game_focus_suggestions: string[];
  expected_gpm_increase: number;
  expected_win_rate_increase: number;
}

export interface SpendingOptimizationData {
  item_order_optimization: string[];
  back_timing_optimization: string[];
  gold_allocation_suggestions: string[];
  component_buying_tips: string[];
  power_spike_timing: string[];
  early_game_priorities: string[];
  mid_game_priorities: string[];
  late_game_priorities: string[];
}

export interface GoldTrendPoint {
  date: string;
  gold_per_minute: number;
  gold_efficiency: number;
  farming_efficiency: number;
  spending_efficiency: number;
  moving_average: number;
}

export interface GoldRecommendation {
  priority: 'high' | 'medium' | 'low';
  category: string;
  title: string;
  description: string;
  impact: string;
  game_phase: string[];
  expected_gpm_increase: number;
  implementation_difficulty: string;
}

export interface MatchGoldData {
  match_id: string;
  champion: string;
  position: string;
  total_gold_earned: number;
  gold_per_minute: number;
  gold_efficiency_score: number;
  farming_gold: number;
  kill_gold: number;
  objective_gold: number;
  items_completed: number;
  control_wards_spent: number;
  game_duration: number;
  result: string;
  date: string;
  gold_advantage_at_15: number;
  team_gold_share: number;
}

export class GoldService extends BaseService {
  /**
   * Get comprehensive gold efficiency analysis for a player
   */
  async getGoldAnalysis(
    playerId: string, 
    timeRange: '7d' | '30d' | '90d' = '30d',
    champion?: string,
    position?: string,
    gameMode?: string
  ): Promise<GoldAnalysis> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(champion && { champion }),
      ...(position && { position }),
      ...(gameMode && { game_mode: gameMode })
    });

    return this.get(`/gold/${playerId}/analysis?${params}`);
  }

  /**
   * Get gold income source breakdown
   */
  async getGoldSources(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    sourceType?: 'farming' | 'kills' | 'objectives' | 'passive' | 'items'
  ): Promise<GoldSourcesData> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(sourceType && { source_type: sourceType })
    });

    return this.get(`/gold/${playerId}/sources?${params}`);
  }

  /**
   * Get item purchase and utilization efficiency
   */
  async getItemEfficiency(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    itemCategory?: 'damage' | 'defensive' | 'utility'
  ): Promise<ItemEfficiencyData> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(itemCategory && { item_category: itemCategory })
    });

    return this.get(`/gold/${playerId}/item-efficiency?${params}`);
  }

  /**
   * Get gold spending behavior analysis
   */
  async getSpendingPatterns(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    spendingType?: 'items' | 'wards' | 'consumables'
  ): Promise<SpendingPatternsData> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(spendingType && { spending_type: spendingType })
    });

    return this.get(`/gold/${playerId}/spending?${params}`);
  }

  /**
   * Get gold efficiency trends over time
   */
  async getGoldTrends(
    playerId: string,
    metric: 'gold_per_minute' | 'efficiency' | 'farming_efficiency' | 'spending_efficiency',
    period: 'daily' | 'weekly' = 'daily',
    days: number = 30
  ): Promise<GoldTrendPoint[]> {
    const params = new URLSearchParams({
      metric,
      period,
      days: days.toString()
    });

    return this.get(`/gold/${playerId}/trends?${params}`);
  }

  /**
   * Compare gold performance with benchmarks
   */
  async getGoldComparison(
    playerId: string,
    benchmarkType: 'role' | 'rank' | 'global',
    timeRange: '7d' | '30d' | '90d' = '30d'
  ): Promise<{
    benchmark_type: string;
    player_metrics: {
      gold_per_minute: number;
      gold_efficiency: number;
      economy_rating: string;
      gold_impact_score: number;
    };
    benchmark_metrics: GoldBenchmark;
    percentile: number;
  }> {
    const params = new URLSearchParams({
      benchmark_type: benchmarkType,
      time_range: timeRange
    });

    return this.get(`/gold/${playerId}/comparison?${params}`);
  }

  /**
   * Get gold efficiency optimization suggestions
   */
  async getGoldOptimization(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    focusArea?: 'income' | 'spending' | 'efficiency',
    difficulty?: 'easy' | 'medium' | 'hard'
  ): Promise<{
    player_id: string;
    time_range: string;
    current_performance: {
      gold_per_minute: number;
      gold_efficiency: number;
      economy_rating: string;
    };
    income_optimization: IncomeOptimizationData;
    spending_optimization: SpendingOptimizationData;
    recommendations: GoldRecommendation[];
    improvement_potential: {
      expected_gpm_increase: number;
      expected_wr_increase: number;
    };
  }> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(focusArea && { focus_area: focusArea }),
      ...(difficulty && { difficulty })
    });

    return this.get(`/gold/${playerId}/optimization?${params}`);
  }

  /**
   * Get gold efficiency by game phase
   */
  async getGoldPhaseAnalysis(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    phase?: 'early' | 'mid' | 'late'
  ): Promise<{
    player_id: string;
    time_range: string;
    early_game: GoldPhaseData;
    mid_game: GoldPhaseData;
    late_game: GoldPhaseData;
    phase_data?: GoldPhaseData;
  }> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(phase && { phase })
    });

    return this.get(`/gold/${playerId}/phases?${params}`);
  }

  /**
   * Get gold analytics summary (multiple metrics in one call)
   */
  async getGoldSummary(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d'
  ): Promise<{
    analysis: GoldAnalysis;
    sources: GoldSourcesData;
    efficiency: ItemEfficiencyData;
    spending: SpendingPatternsData;
    trends: GoldTrendPoint[];
    benchmarks: {
      role: GoldBenchmark;
      rank: GoldBenchmark;
      global: GoldBenchmark;
    };
  }> {
    const [analysis, sources, efficiency, spending, trends, roleBenchmark, rankBenchmark, globalBenchmark] = await Promise.all([
      this.getGoldAnalysis(playerId, timeRange),
      this.getGoldSources(playerId, timeRange),
      this.getItemEfficiency(playerId, timeRange),
      this.getSpendingPatterns(playerId, timeRange),
      this.getGoldTrends(playerId, 'gold_per_minute', 'daily'),
      this.getGoldComparison(playerId, 'role', timeRange),
      this.getGoldComparison(playerId, 'rank', timeRange),
      this.getGoldComparison(playerId, 'global', timeRange)
    ]);

    return {
      analysis,
      sources,
      efficiency,
      spending,
      trends,
      benchmarks: {
        role: roleBenchmark.benchmark_metrics,
        rank: rankBenchmark.benchmark_metrics,
        global: globalBenchmark.benchmark_metrics
      }
    };
  }

  /**
   * Get gold insights for multiple time periods (trend analysis)
   */
  async getGoldInsights(playerId: string): Promise<{
    current_period: GoldAnalysis;
    previous_period: GoldAnalysis;
    improvement_areas: string[];
    strength_areas: string[];
    trend_direction: 'improving' | 'declining' | 'stable';
    optimization_priority: 'income' | 'spending' | 'efficiency';
  }> {
    const [current, previous] = await Promise.all([
      this.getGoldAnalysis(playerId, '30d'),
      this.getGoldAnalysis(playerId, '90d')
    ]);

    // Simple trend analysis
    const trendDirection = current.average_gold_per_minute > previous.average_gold_per_minute ? 'improving' :
                          current.average_gold_per_minute < previous.average_gold_per_minute * 0.95 ? 'declining' : 'stable';

    const improvementAreas = [];
    const strengthAreas = [];

    // Analyze areas for improvement
    if (current.role_benchmark.player_percentile < 50) {
      improvementAreas.push('gold_generation');
    }
    if (current.gold_efficiency_score < 70) {
      improvementAreas.push('gold_efficiency');
    }
    if (current.gold_sources.farming_percent < 45) {
      improvementAreas.push('farming_focus');
    }

    // Analyze strengths
    if (current.role_benchmark.player_percentile > 75) {
      strengthAreas.push('gold_generation');
    }
    if (current.gold_efficiency_score > 80) {
      strengthAreas.push('gold_efficiency');
    }
    if (current.gold_sources.farming_percent > 55) {
      strengthAreas.push('farming_consistency');
    }

    // Determine optimization priority
    let optimizationPriority: 'income' | 'spending' | 'efficiency' = 'income';
    if (current.gold_efficiency_score < 60) {
      optimizationPriority = 'efficiency';
    } else if (current.average_gold_per_minute < current.role_benchmark.average_gold_per_minute * 0.9) {
      optimizationPriority = 'income';
    } else {
      optimizationPriority = 'spending';
    }

    return {
      current_period: current,
      previous_period: previous,
      improvement_areas: improvementAreas,
      strength_areas: strengthAreas,
      trend_direction: trendDirection,
      optimization_priority: optimizationPriority
    };
  }

  /**
   * Get gold coaching recommendations based on role and playstyle
   */
  async getGoldCoaching(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    playstyle?: 'aggressive' | 'passive' | 'balanced'
  ): Promise<{
    current_rating: string;
    target_rating: string;
    coaching_plan: {
      week_1: string[];
      week_2: string[];
      week_3: string[];
      week_4: string[];
    };
    practice_drills: string[];
    expected_improvement: {
      gpm_increase: number;
      efficiency_increase: number;
      timeline_weeks: number;
    };
  }> {
    const analysis = await this.getGoldAnalysis(playerId, timeRange);
    const optimization = await this.getGoldOptimization(playerId, timeRange);
    
    // Generate coaching plan based on analysis
    const currentRating = analysis.economy_rating;
    const targetRating = currentRating === 'poor' ? 'average' :
                        currentRating === 'average' ? 'good' :
                        currentRating === 'good' ? 'excellent' : 'excellent';
    
    const coachingPlan = {
      week_1: [
        'Focus on last-hitting fundamentals',
        'Practice optimal back timing',
        'Track gold sources after each game'
      ],
      week_2: [
        'Improve jungle camp timing',
        'Work on item component optimization',
        'Analyze farming patterns by game phase'
      ],
      week_3: [
        'Practice team fight positioning for gold efficiency',
        'Optimize control ward spending',
        'Focus on objective participation'
      ],
      week_4: [
        'Integrate all improvements',
        'Focus on consistency',
        'Analyze progress and adjust strategy'
      ]
    };
    
    const practiceDrills = [
      'CS training in practice tool (10 min/day)',
      'Item build path optimization review',
      'Back timing decision tree practice',
      'Objective priority ranking exercises'
    ];
    
    const expectedImprovement = {
      gpm_increase: optimization.improvement_potential.expected_gpm_increase * 0.7, // Conservative estimate
      efficiency_increase: 15.0, // Expected efficiency score improvement
      timeline_weeks: 4
    };
    
    return {
      current_rating: currentRating,
      target_rating: targetRating,
      coaching_plan: coachingPlan,
      practice_drills: practiceDrills,
      expected_improvement: expectedImprovement
    };
  }

  /**
   * Get champion-specific gold efficiency analysis
   */
  async getChampionGoldAnalysis(
    playerId: string,
    champion: string,
    timeRange: '7d' | '30d' | '90d' = '30d'
  ): Promise<{
    champion: string;
    games_played: number;
    gold_performance: GoldAnalysis;
    champion_benchmarks: GoldBenchmark;
    optimal_build_paths: string[];
    power_spike_analysis: {
      first_item: number;
      core_items: number;
      late_game: number;
    };
    recommendations: string[];
  }> {
    const analysis = await this.getGoldAnalysis(playerId, timeRange, champion);
    
    // Champion-specific analysis would be enhanced with champion data
    const optimalBuildPaths = [
      'Standard scaling build for consistent gold income',
      'Early game focused build for lane dominance',
      'Utility build for team-oriented gold efficiency'
    ];
    
    const powerSpikeAnalysis = {
      first_item: analysis.item_efficiency.first_item_timing,
      core_items: analysis.item_efficiency.core_items_timing,
      late_game: analysis.item_efficiency.six_items_timing
    };
    
    const recommendations = analysis.recommendations.slice(0, 5).map(rec => rec.title);
    
    return {
      champion,
      games_played: analysis.recent_matches.length,
      gold_performance: analysis,
      champion_benchmarks: analysis.role_benchmark, // Would be champion-specific in real implementation
      optimal_build_paths: optimalBuildPaths,
      power_spike_analysis: powerSpikeAnalysis,
      recommendations: recommendations
    };
  }
}

// Export singleton instance
export const goldService = new GoldService();