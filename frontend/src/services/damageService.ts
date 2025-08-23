import { BaseService } from './baseService';

export interface DamageAnalysis {
  player_id: string;
  champion?: string;
  position?: string;
  time_range: string;
  damage_share: number;
  damage_per_minute: number;
  total_damage: number;
  physical_damage_percent: number;
  magic_damage_percent: number;
  true_damage_percent: number;
  carry_potential: number;
  efficiency_rating: string;
  damage_consistency: number;
  team_contribution: TeamContributionData;
  damage_distribution: DamageDistribution;
  game_phase_analysis: GamePhaseAnalysis;
  high_damage_win_rate: number;
  low_damage_win_rate: number;
  role_benchmark: DamageBenchmark;
  rank_benchmark: DamageBenchmark;
  global_benchmark: DamageBenchmark;
  trend_data: DamageTrendPoint[];
  recommendations: DamageRecommendation[];
}

export interface TeamContributionData {
  damage_contribution_score: number;
  kill_participation: number;
  solo_kill_rate: number;
  team_fight_damage_share: number;
  objective_damage_share: number;
  consistent_damage_games: number;
  clutch_performance_score: number;
}

export interface DamageDistribution {
  champion_damage_percent: number;
  structure_damage_percent: number;
  monster_damage_percent: number;
  damage_by_type: {
    physical: number;
    magical: number;
    true_damage: number;
  };
  target_priority: {
    carries: number;
    tanks: number;
    supports: number;
  };
}

export interface GamePhaseAnalysis {
  early_game: {
    damage_per_minute: number;
    damage_share: number;
    carry_potential: number;
    efficiency: string;
  };
  mid_game: {
    damage_per_minute: number;
    damage_share: number;
    carry_potential: number;
    efficiency: string;
  };
  late_game: {
    damage_per_minute: number;
    damage_share: number;
    carry_potential: number;
    efficiency: string;
  };
}

export interface DamageBenchmark {
  category: string;
  average_damage_share: number;
  top_10_percent: number;
  top_25_percent: number;
  median: number;
  player_percentile: number;
}

export interface DamageTrendPoint {
  date: string;
  damage_per_minute: number;
  damage_share: number;
  carry_potential: number;
  moving_average: number;
  efficiency: number;
}

export interface DamageRecommendation {
  priority: 'high' | 'medium' | 'low';
  category: string;
  title: string;
  description: string;
  impact: string;
  game_phase: string[];
}

export interface CarryPotentialAnalysis {
  player_id: string;
  time_range: string;
  carry_potential: number;
  damage_share: number;
  consistency_score: number;
  team_contribution: TeamContributionData;
  game_phase_analysis: {
    early_game_carry: number;
    mid_game_carry: number;
    late_game_carry: number;
  };
  win_rate_correlation: {
    high_damage_games: number;
    low_damage_games: number;
    impact_score: number;
  };
  recommendations: DamageRecommendation[];
}

export class DamageService extends BaseService {
  /**
   * Get comprehensive damage analysis for a player
   */
  async getDamageAnalysis(
    playerId: string, 
    timeRange: '7d' | '30d' | '90d' = '30d',
    champion?: string,
    position?: string,
    gameMode?: string
  ): Promise<DamageAnalysis> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(champion && { champion }),
      ...(position && { position }),
      ...(gameMode && { game_mode: gameMode })
    });

    return this.get(`/damage/${playerId}/analysis?${params}`);
  }

  /**
   * Get team contribution metrics
   */
  async getTeamContribution(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    teamRole?: string,
    season?: string
  ): Promise<TeamContributionData> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(teamRole && { team_role: teamRole }),
      ...(season && { season })
    });

    return this.get(`/damage/${playerId}/contribution?${params}`);
  }

  /**
   * Get damage distribution analysis
   */
  async getDamageDistribution(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    targetType?: 'champions' | 'structures' | 'monsters' | 'objectives'
  ): Promise<DamageDistribution> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(targetType && { target_type: targetType })
    });

    return this.get(`/damage/${playerId}/distribution?${params}`);
  }

  /**
   * Get damage performance trends over time
   */
  async getDamageTrends(
    playerId: string,
    metric: 'damage_per_minute' | 'damage_share' | 'carry_potential' | 'efficiency',
    period: 'daily' | 'weekly' = 'daily',
    days: number = 30
  ): Promise<DamageTrendPoint[]> {
    const params = new URLSearchParams({
      metric,
      period,
      days: days.toString()
    });

    return this.get(`/damage/${playerId}/trends?${params}`);
  }

  /**
   * Compare damage performance with benchmarks
   */
  async getDamageComparison(
    playerId: string,
    benchmarkType: 'role' | 'rank' | 'global',
    timeRange: '7d' | '30d' | '90d' = '30d'
  ): Promise<{
    benchmark_type: string;
    player_metrics: {
      damage_share: number;
      damage_per_minute: number;
      carry_potential: number;
      efficiency_rating: string;
    };
    benchmark_metrics: DamageBenchmark;
    percentile: number;
  }> {
    const params = new URLSearchParams({
      benchmark_type: benchmarkType,
      time_range: timeRange
    });

    return this.get(`/damage/${playerId}/comparison?${params}`);
  }

  /**
   * Get carry potential analysis
   */
  async getCarryPotentialAnalysis(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d',
    winCondition?: 'damage_carry' | 'utility_carry' | 'mixed'
  ): Promise<CarryPotentialAnalysis> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(winCondition && { win_condition: winCondition })
    });

    return this.get(`/damage/${playerId}/carry-potential?${params}`);
  }

  /**
   * Get damage analytics summary (multiple metrics in one call)
   */
  async getDamageSummary(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d'
  ): Promise<{
    analysis: DamageAnalysis;
    trends: DamageTrendPoint[];
    benchmarks: {
      role: DamageBenchmark;
      rank: DamageBenchmark;
      global: DamageBenchmark;
    };
  }> {
    const [analysis, trends, roleBenchmark, rankBenchmark, globalBenchmark] = await Promise.all([
      this.getDamageAnalysis(playerId, timeRange),
      this.getDamageTrends(playerId, 'damage_share', 'daily'),
      this.getDamageComparison(playerId, 'role', timeRange),
      this.getDamageComparison(playerId, 'rank', timeRange),
      this.getDamageComparison(playerId, 'global', timeRange)
    ]);

    return {
      analysis,
      trends,
      benchmarks: {
        role: roleBenchmark.benchmark_metrics,
        rank: rankBenchmark.benchmark_metrics,
        global: globalBenchmark.benchmark_metrics
      }
    };
  }

  /**
   * Get damage insights for multiple time periods (trend analysis)
   */
  async getDamageInsights(playerId: string): Promise<{
    current_period: DamageAnalysis;
    previous_period: DamageAnalysis;
    improvement_areas: string[];
    strength_areas: string[];
    trend_direction: 'improving' | 'declining' | 'stable';
  }> {
    const [current, previous] = await Promise.all([
      this.getDamageAnalysis(playerId, '30d'),
      this.getDamageAnalysis(playerId, '90d') // Get longer period for comparison
    ]);

    // Simple trend analysis
    const trendDirection = current.damage_share > previous.damage_share ? 'improving' :
                          current.damage_share < previous.damage_share * 0.95 ? 'declining' : 'stable';

    const improvementAreas = [];
    const strengthAreas = [];

    // Analyze areas for improvement
    if (current.role_benchmark.player_percentile < 50) {
      improvementAreas.push('damage_output');
    }
    if (current.damage_consistency < 70) {
      improvementAreas.push('consistency');
    }
    if (current.carry_potential < 50) {
      improvementAreas.push('carry_potential');
    }

    // Analyze strengths
    if (current.role_benchmark.player_percentile > 75) {
      strengthAreas.push('damage_output');
    }
    if (current.damage_consistency > 80) {
      strengthAreas.push('consistency');
    }
    if (current.team_contribution.damage_contribution_score > 75) {
      strengthAreas.push('team_contribution');
    }

    return {
      current_period: current,
      previous_period: previous,
      improvement_areas: improvementAreas,
      strength_areas: strengthAreas,
      trend_direction: trendDirection
    };
  }

  /**
   * Get damage efficiency breakdown by game phase
   */
  async getDamageEfficiencyByPhase(
    playerId: string,
    timeRange: '7d' | '30d' | '90d' = '30d'
  ): Promise<{
    early_game: { efficiency: number; damage_per_minute: number; impact: string; };
    mid_game: { efficiency: number; damage_per_minute: number; impact: string; };
    late_game: { efficiency: number; damage_per_minute: number; impact: string; };
    recommendations: string[];
  }> {
    const analysis = await this.getDamageAnalysis(playerId, timeRange);
    
    const getImpact = (efficiency: number, dpm: number): string => {
      if (efficiency > 80 && dpm > 800) return 'High Impact';
      if (efficiency > 60 && dpm > 600) return 'Medium Impact';
      return 'Low Impact';
    };

    const recommendations = [];
    
    // Generate phase-specific recommendations
    if (analysis.game_phase_analysis.early_game.carry_potential < 40) {
      recommendations.push('Focus on early game farming and safe trading');
    }
    if (analysis.game_phase_analysis.mid_game.carry_potential < 50) {
      recommendations.push('Improve mid-game team fighting and objective control');
    }
    if (analysis.game_phase_analysis.late_game.carry_potential < 60) {
      recommendations.push('Work on late-game positioning and target selection');
    }

    return {
      early_game: {
        efficiency: analysis.game_phase_analysis.early_game.carry_potential,
        damage_per_minute: analysis.game_phase_analysis.early_game.damage_per_minute,
        impact: getImpact(analysis.game_phase_analysis.early_game.carry_potential, analysis.game_phase_analysis.early_game.damage_per_minute)
      },
      mid_game: {
        efficiency: analysis.game_phase_analysis.mid_game.carry_potential,
        damage_per_minute: analysis.game_phase_analysis.mid_game.damage_per_minute,
        impact: getImpact(analysis.game_phase_analysis.mid_game.carry_potential, analysis.game_phase_analysis.mid_game.damage_per_minute)
      },
      late_game: {
        efficiency: analysis.game_phase_analysis.late_game.carry_potential,
        damage_per_minute: analysis.game_phase_analysis.late_game.damage_per_minute,
        impact: getImpact(analysis.game_phase_analysis.late_game.carry_potential, analysis.game_phase_analysis.late_game.damage_per_minute)
      },
      recommendations
    };
  }
}

// Export singleton instance
export const damageService = new DamageService();