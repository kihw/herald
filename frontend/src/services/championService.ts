import { apiClient } from './apiClient';
import type {
  ChampionAnalysis,
  ChampionTrendPoint,
  ChampionMasteryRanking,
  ChampionComparisonData,
  ChampionPowerSpike,
  ChampionMatchupData,
  ChampionBuildData,
  ChampionCoachingData
} from '../types/champion';

class ChampionService {
  private baseUrl = '/api/v1/champion';

  /**
   * Get comprehensive champion-specific performance analysis
   */
  async getChampionAnalysis(
    playerId: string,
    timeRange: string,
    champion: string,
    position?: string
  ): Promise<ChampionAnalysis> {
    const params = new URLSearchParams({
      time_range: timeRange,
      champion,
      ...(position && { position })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/analysis?${params}`);
    return response.data;
  }

  /**
   * Get champion mastery ranking
   */
  async getChampionMastery(
    playerId: string,
    timeRange: string,
    limit?: number
  ): Promise<ChampionMasteryRanking[]> {
    const params = new URLSearchParams({
      time_range: timeRange,
      ...(limit && { limit: limit.toString() })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/mastery?${params}`);
    return response.data;
  }

  /**
   * Compare performance across multiple champions
   */
  async getChampionComparison(
    playerId: string,
    champions: string[],
    timeRange: string
  ): Promise<ChampionComparisonData> {
    const params = new URLSearchParams({
      champions: champions.join(','),
      time_range: timeRange
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/comparison?${params}`);
    return response.data;
  }

  /**
   * Get champion power spike analysis
   */
  async getChampionPowerSpikes(
    playerId: string,
    champion: string,
    timeRange?: string
  ): Promise<{
    player_id: string;
    champion: string;
    time_range: string;
    power_spikes: ChampionPowerSpike[];
    carry_potential: number;
    scaling_analysis: {
      early_game_rating: number;
      mid_game_rating: number;
      late_game_rating: number;
    };
  }> {
    const params = new URLSearchParams({
      champion,
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/power-spikes?${params}`);
    return response.data;
  }

  /**
   * Get champion matchup analysis
   */
  async getChampionMatchups(
    playerId: string,
    champion: string,
    opponent?: string,
    timeRange?: string
  ): Promise<ChampionMatchupData> {
    const params = new URLSearchParams({
      champion,
      ...(opponent && { opponent }),
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/matchups?${params}`);
    return response.data;
  }

  /**
   * Get champion build analysis
   */
  async getChampionBuilds(
    playerId: string,
    champion: string,
    buildType?: string,
    timeRange?: string
  ): Promise<ChampionBuildData> {
    const params = new URLSearchParams({
      champion,
      ...(buildType && { build_type: buildType }),
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/builds?${params}`);
    return response.data;
  }

  /**
   * Get champion-specific coaching recommendations
   */
  async getChampionCoaching(
    playerId: string,
    champion: string,
    focusArea?: string,
    timeRange?: string
  ): Promise<ChampionCoachingData> {
    const params = new URLSearchParams({
      champion,
      ...(focusArea && { focus_area: focusArea }),
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/coaching?${params}`);
    return response.data;
  }

  /**
   * Get champion performance trends
   */
  async getChampionTrends(
    playerId: string,
    champion: string,
    metric?: string,
    period?: string,
    days?: number
  ): Promise<ChampionTrendPoint[]> {
    const params = new URLSearchParams({
      champion,
      ...(metric && { metric }),
      ...(period && { period }),
      ...(days && { days: days.toString() })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/trends?${params}`);
    return response.data;
  }

  /**
   * Get champion performance by game length
   */
  async getChampionPerformanceByGameLength(
    playerId: string,
    champion: string,
    timeRange: string
  ): Promise<{
    short_games: { average_length: number; win_rate: number; performance_rating: number };
    medium_games: { average_length: number; win_rate: number; performance_rating: number };
    long_games: { average_length: number; win_rate: number; performance_rating: number };
  }> {
    // This would be a custom endpoint for game length analysis
    const params = new URLSearchParams({
      champion,
      time_range: timeRange,
      analysis_type: 'game_length'
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/performance-by-length?${params}`);
    return response.data;
  }

  /**
   * Get champion role flexibility analysis
   */
  async getChampionRoleFlexibility(
    playerId: string,
    champion: string,
    timeRange: string
  ): Promise<{
    primary_role: string;
    secondary_roles: Array<{
      role: string;
      games_played: number;
      win_rate: number;
      performance_rating: number;
    }>;
    flexibility_score: number;
    role_adaptation: number;
  }> {
    const params = new URLSearchParams({
      champion,
      time_range: timeRange
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/role-flexibility?${params}`);
    return response.data;
  }

  /**
   * Get champion meta analysis
   */
  async getChampionMetaAnalysis(
    playerId: string,
    champion: string,
    timeRange: string
  ): Promise<{
    current_meta_tier: string;
    meta_alignment_score: number;
    meta_trends: Array<{
      patch: string;
      tier: string;
      play_rate: number;
      win_rate: number;
    }>;
    optimal_patches: string[];
    meta_recommendations: Array<{
      type: string;
      description: string;
      impact: string;
    }>;
  }> {
    const params = new URLSearchParams({
      champion,
      time_range: timeRange
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/meta-analysis?${params}`);
    return response.data;
  }

  /**
   * Get champion improvement recommendations
   */
  async getChampionImprovementPlan(
    playerId: string,
    champion: string,
    focusArea?: string,
    timeRange?: string
  ): Promise<{
    current_skill_level: string;
    target_skill_level: string;
    improvement_timeline: string;
    milestones: Array<{
      milestone: string;
      description: string;
      estimated_games: number;
      key_metrics: string[];
    }>;
    practice_routine: Array<{
      category: string;
      exercises: Array<{
        name: string;
        description: string;
        frequency: string;
        duration: string;
      }>;
    }>;
  }> {
    const params = new URLSearchParams({
      champion,
      ...(focusArea && { focus_area: focusArea }),
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/improvement-plan?${params}`);
    return response.data;
  }

  /**
   * Get champion synergy analysis with teammates
   */
  async getChampionSynergy(
    playerId: string,
    champion: string,
    timeRange: string
  ): Promise<{
    best_synergies: Array<{
      teammate_champion: string;
      synergy_score: number;
      win_rate_together: number;
      games_played: number;
      key_synergies: string[];
    }>;
    worst_synergies: Array<{
      teammate_champion: string;
      synergy_score: number;
      win_rate_together: number;
      games_played: number;
      common_issues: string[];
    }>;
    team_composition_preferences: {
      engage_comps: number;
      poke_comps: number;
      split_push_comps: number;
      team_fight_comps: number;
    };
  }> {
    const params = new URLSearchParams({
      champion,
      time_range: timeRange
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/synergy?${params}`);
    return response.data;
  }

  /**
   * Get champion ban analysis
   */
  async getChampionBanAnalysis(
    playerId: string,
    champion: string,
    timeRange: string
  ): Promise<{
    ban_rate_against: number;
    ban_priority: number;
    most_banned_by: Array<{
      opponent_champion: string;
      ban_frequency: number;
      reasoning: string[];
    }>;
    threat_level: string;
    counterpick_frequency: number;
  }> {
    const params = new URLSearchParams({
      champion,
      time_range: timeRange
    });

    const response = await apiClient.get(`${this.baseUrl}/${playerId}/ban-analysis?${params}`);
    return response.data;
  }

  /**
   * Get champion learning resources
   */
  async getChampionLearningResources(
    champion: string,
    skillLevel?: string,
    focusArea?: string
  ): Promise<{
    guides: Array<{
      title: string;
      author: string;
      difficulty: string;
      focus_areas: string[];
      rating: number;
      url: string;
    }>;
    video_tutorials: Array<{
      title: string;
      creator: string;
      duration: string;
      topics: string[];
      skill_level: string;
      url: string;
    }>;
    pro_player_resources: Array<{
      player_name: string;
      team: string;
      playstyle: string;
      vod_links: string[];
      notable_matches: string[];
    }>;
    practice_tools: Array<{
      tool_name: string;
      description: string;
      focus_area: string;
      difficulty: string;
    }>;
  }> {
    const params = new URLSearchParams({
      champion,
      ...(skillLevel && { skill_level: skillLevel }),
      ...(focusArea && { focus_area: focusArea })
    });

    const response = await apiClient.get(`${this.baseUrl}/learning-resources?${params}`);
    return response.data;
  }
}

export const championService = new ChampionService();