import { apiClient } from './apiClient';
import type { 
  WardAnalysis, 
  WardTrendPoint, 
  WardPlacementPattern,
  MapControlData,
  WardTypeAnalysis,
  WardOptimizationData
} from '../types/ward';

class WardService {
  private baseUrl = '/api/v1/ward';

  /**
   * Get comprehensive ward placement and map control analysis
   */
  async getWardAnalysis(
    playerId: string,
    timeRange: string,
    champion?: string,
    position?: string,
    gameMode?: string
  ): Promise<WardAnalysis> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(champion && { champion }),
      ...(position && { position }),
      ...(gameMode && { game_mode: gameMode })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/analysis?${params}`);
    return response.data;
  }

  /**
   * Get map control analysis with optional zone filtering
   */
  async getMapControl(
    playerId: string,
    timeRange: string,
    zone?: string,
    metric?: string
  ): Promise<MapControlData> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(zone && { zone }),
      ...(metric && { metric })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/map-control?${params}`);
    return response.data;
  }

  /**
   * Get ward placement pattern analysis
   */
  async getWardPlacementPatterns(
    playerId: string,
    timeRange?: string,
    patternType?: string
  ): Promise<{
    player_id: string;
    time_range: string;
    placement_patterns: WardPlacementPattern;
    optimal_placements: Array<{
      location: string;
      frequency: number;
      effectiveness: number;
    }>;
    placement_timing: {
      earlyGameFrequency: number;
      midGameFrequency: number;
      lateGameFrequency: number;
      optimalTimings: Array<{
        gameTime: number;
        context: string;
        effectiveness: number;
      }>;
    };
    placement_optimization: WardOptimizationData;
  }> {
    const params = new URLSearchParams({
      ...(timeRange && { time_range: timeRange }),
      ...(patternType && { pattern_type: patternType })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/patterns?${params}`);
    return response.data;
  }

  /**
   * Get ward type specific analysis
   */
  async getWardTypeAnalysis(
    playerId: string,
    wardType?: string,
    timeRange?: string
  ): Promise<{
    player_id: string;
    time_range: string;
    yellow_wards_analysis: WardTypeAnalysis;
    control_wards_analysis: WardTypeAnalysis;
    blue_ward_analysis: WardTypeAnalysis;
    focus_analysis?: WardTypeAnalysis;
  }> {
    const params = new URLSearchParams({
      ...(wardType && { ward_type: wardType }),
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/types?${params}`);
    return response.data;
  }

  /**
   * Get ward performance trends
   */
  async getWardTrends(
    playerId: string,
    metric: string,
    period?: string,
    days?: number
  ): Promise<WardTrendPoint[]> {
    const params = new URLSearchParams({
      metric,
      ...(period && { period }),
      ...(days && { days: days.toString() })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/trends?${params}`);
    return response.data;
  }

  /**
   * Get ward clearing analysis
   */
  async getWardClearing(
    playerId: string,
    timeRange?: string,
    clearingType?: string
  ): Promise<{
    player_id: string;
    time_range: string;
    ward_clearing_patterns: {
      proactiveClearing: number;
      reactiveClearing: number;
      opportunisticClearing: number;
      clearingEfficiency: number;
      averageClearTime: number;
    };
    counter_warding_score: number;
    clearing_optimization: WardOptimizationData;
    vision_denied_score: number;
  }> {
    const params = new URLSearchParams({
      ...(timeRange && { time_range: timeRange }),
      ...(clearingType && { clearing_type: clearingType })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/clearing?${params}`);
    return response.data;
  }

  /**
   * Get ward placement and clearing optimization suggestions
   */
  async getWardOptimization(
    playerId: string,
    timeRange?: string,
    optimizationType?: string
  ): Promise<{
    player_id: string;
    time_range: string;
    current_performance: {
      map_control_score: number;
      ward_efficiency: number;
      counter_warding_score: number;
    };
    placement_optimization: WardOptimizationData;
    clearing_optimization: WardOptimizationData;
    recommendations: Array<{
      type: string;
      title: string;
      description: string;
      priority: string;
      expectedImprovement: number;
    }>;
    improvement_potential: {
      expected_control_gain: number;
      expected_safety_gain: number;
      expected_denial_gain: number;
    };
  }> {
    const params = new URLSearchParams({
      ...(timeRange && { time_range: timeRange }),
      ...(optimizationType && { optimization_type: optimizationType })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/optimization?${params}`);
    return response.data;
  }

  /**
   * Get objective-specific vision control analysis
   */
  async getObjectiveControl(
    playerId: string,
    objective?: string,
    timeRange?: string
  ): Promise<{
    player_id: string;
    time_range: string;
    objective_setup: {
      DragonSetupScore: number;
      BaronSetupScore: number;
      HeraldSetupScore: number;
      ElderSetupScore: number;
      averageSetupTime: number;
      setupConsistency: number;
    };
    safety_provided: {
      overallSafetyScore: number;
      escapeRouteCoverage: number;
      flankProtection: number;
      engageWarning: number;
    };
    strategic_coverage: {
      overallScore: number;
      dragonPitCoverage: number;
      baronPitCoverage: number;
      keyChokepoints: number;
      rotationPaths: number;
    };
    focus_score?: number;
    focus_coverage?: number;
  }> {
    const params = new URLSearchParams({
      ...(objective && { objective }),
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/objectives?${params}`);
    return response.data;
  }

  /**
   * Get ward heatmap data for visualization
   */
  async getWardHeatmap(
    playerId: string,
    timeRange: string,
    wardType?: string
  ): Promise<{
    heatmapData: Array<{
      x: number;
      y: number;
      intensity: number;
      wardType: string;
      effectiveness: number;
    }>;
    mapBounds: {
      minX: number;
      maxX: number;
      minY: number;
      maxY: number;
    };
  }> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(wardType && { ward_type: wardType })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/heatmap?${params}`);
    return response.data;
  }

  /**
   * Get ward comparison with other players
   */
  async compareWardPerformance(
    playerId: string,
    compareWithPlayerIds: string[],
    timeRange: string,
    metric?: string
  ): Promise<{
    comparison_data: Array<{
      player_id: string;
      player_name: string;
      map_control_score: number;
      ward_efficiency: number;
      wards_placed_per_game: number;
      wards_killed_per_game: number;
      counter_warding_score: number;
      percentile: number;
    }>;
    ranking: Array<{
      player_id: string;
      rank: number;
      score: number;
    }>;
  }> {
    const params = new URLSearchParams({
      time_range: timeRange,
      compare_with: compareWithPlayerIds.join(','),
      ...(metric && { metric })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/compare?${params}`);
    return response.data;
  }

  /**
   * Get ward coaching recommendations
   */
  async getWardCoachingRecommendations(
    playerId: string,
    timeRange: string,
    focusArea?: string
  ): Promise<{
    recommendations: Array<{
      category: string;
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
    }>;
    learningPath: Array<{
      step: number;
      title: string;
      description: string;
      estimatedTime: string;
      prerequisites: string[];
    }>;
    progressMetrics: Array<{
      metric: string;
      currentValue: number;
      targetValue: number;
      timeframe: string;
    }>;
  }> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(focusArea && { focus_area: focusArea })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/coaching?${params}`);
    return response.data;
  }
}

export const wardService = new WardService();