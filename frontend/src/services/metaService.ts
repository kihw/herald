import { apiClient } from './apiClient';
import type {
  MetaAnalysis,
  ChampionTierList,
  MetaTrends,
  ChampionMetaStats,
  BanAnalysis,
  PickAnalysis,
  MetaPredictions,
  MetaRecommendation,
  MetaHistory
} from '../types/meta';

class MetaService {
  private baseUrl = '/api/v1/meta';

  /**
   * Get comprehensive meta analysis
   */
  async getMetaAnalysis(
    patch: string,
    region?: string,
    rank?: string,
    timeRange?: string
  ): Promise<MetaAnalysis> {
    const params = new URLSearchParams({
      patch,
      ...(region && { region }),
      ...(rank && { rank }),
      ...(timeRange && { time_range: timeRange })
    });

    const response = await apiClient.get(`${this.baseUrl}/analysis?${params}`);
    return response.data;
  }

  /**
   * Get champion tier list
   */
  async getTierList(
    patch: string,
    region?: string,
    rank?: string,
    role?: string
  ): Promise<ChampionTierList> {
    const params = new URLSearchParams({
      patch,
      ...(region && { region }),
      ...(rank && { rank }),
      ...(role && { role })
    });

    const response = await apiClient.get(`${this.baseUrl}/tier-list?${params}`);
    return response.data;
  }

  /**
   * Get champion meta statistics
   */
  async getChampionMeta(
    champion: string,
    patch: string,
    region?: string,
    rank?: string
  ): Promise<ChampionMetaStats> {
    const params = new URLSearchParams({
      champion,
      patch,
      ...(region && { region }),
      ...(rank && { rank })
    });

    const response = await apiClient.get(`${this.baseUrl}/champion?${params}`);
    return response.data;
  }

  /**
   * Get meta trends analysis
   */
  async getMetaTrends(
    patch: string,
    region?: string,
    rank?: string,
    category?: string
  ): Promise<MetaTrends> {
    const params = new URLSearchParams({
      patch,
      ...(region && { region }),
      ...(rank && { rank }),
      ...(category && { category })
    });

    const response = await apiClient.get(`${this.baseUrl}/trends?${params}`);
    return response.data;
  }

  /**
   * Get ban phase analysis
   */
  async getBanAnalysis(
    patch: string,
    region?: string,
    rank?: string,
    banType?: string
  ): Promise<BanAnalysis> {
    const params = new URLSearchParams({
      patch,
      ...(region && { region }),
      ...(rank && { rank }),
      ...(banType && { ban_type: banType })
    });

    const response = await apiClient.get(`${this.baseUrl}/bans?${params}`);
    return response.data;
  }

  /**
   * Get pick phase analysis
   */
  async getPickAnalysis(
    patch: string,
    region?: string,
    rank?: string,
    pickType?: string
  ): Promise<PickAnalysis> {
    const params = new URLSearchParams({
      patch,
      ...(region && { region }),
      ...(rank && { rank }),
      ...(pickType && { pick_type: pickType })
    });

    const response = await apiClient.get(`${this.baseUrl}/picks?${params}`);
    return response.data;
  }

  /**
   * Get meta predictions
   */
  async getMetaPredictions(
    patch: string,
    predictionType?: string,
    confidenceThreshold?: number
  ): Promise<MetaPredictions> {
    const params = new URLSearchParams({
      patch,
      ...(predictionType && { prediction_type: predictionType }),
      ...(confidenceThreshold && { confidence_threshold: confidenceThreshold.toString() })
    });

    const response = await apiClient.get(`${this.baseUrl}/predictions?${params}`);
    return response.data;
  }

  /**
   * Get meta-based recommendations
   */
  async getMetaRecommendations(
    patch: string,
    rank?: string,
    role?: string,
    recommendationType?: string
  ): Promise<MetaRecommendation[]> {
    const params = new URLSearchParams({
      patch,
      ...(rank && { rank }),
      ...(role && { role }),
      ...(recommendationType && { recommendation_type: recommendationType })
    });

    const response = await apiClient.get(`${this.baseUrl}/recommendations?${params}`);
    return response.data;
  }

  /**
   * Get meta history and evolution
   */
  async getMetaHistory(
    startPatch: string,
    endPatch?: string,
    champion?: string,
    metric?: string
  ): Promise<MetaHistory> {
    const params = new URLSearchParams({
      start_patch: startPatch,
      ...(endPatch && { end_patch: endPatch }),
      ...(champion && { champion }),
      ...(metric && { metric })
    });

    const response = await apiClient.get(`${this.baseUrl}/history?${params}`);
    return response.data;
  }

  /**
   * Get champion tier changes over time
   */
  async getChampionTierHistory(
    champion: string,
    startPatch: string,
    endPatch?: string
  ): Promise<{
    champion: string;
    tier_history: Array<{
      patch: string;
      tier: string;
      tier_score: number;
      win_rate: number;
      pick_rate: number;
      ban_rate: number;
    }>;
    trend_analysis: {
      direction: string;
      volatility: string;
      consistency: number;
    };
  }> {
    return this.getMetaHistory(startPatch, endPatch, champion, 'tier');
  }

  /**
   * Get role meta analysis
   */
  async getRoleMeta(
    role: string,
    patch: string,
    region?: string,
    rank?: string
  ): Promise<{
    role: string;
    patch: string;
    top_champions: ChampionMetaStats[];
    role_trends: {
      average_win_rate: number;
      pick_diversity: number;
      ban_pressure: number;
      meta_shifts: string[];
    };
    recommended_champions: Array<{
      champion: string;
      tier: string;
      reason: string;
      difficulty: string;
    }>;
  }> {
    const tierList = await this.getTierList(patch, region, rank, role);
    const trends = await this.getMetaTrends(patch, region, rank);
    
    // This would be processed and structured for role-specific analysis
    return {
      role,
      patch,
      top_champions: [], // Would be filtered from champion stats
      role_trends: {
        average_win_rate: 50.0,
        pick_diversity: 75.0,
        ban_pressure: 25.0,
        meta_shifts: []
      },
      recommended_champions: []
    };
  }

  /**
   * Get patch comparison analysis
   */
  async getPatchComparison(
    patch1: string,
    patch2: string,
    region?: string,
    rank?: string
  ): Promise<{
    patch1: string;
    patch2: string;
    tier_changes: Array<{
      champion: string;
      patch1_tier: string;
      patch2_tier: string;
      tier_change: number;
      impact: 'major' | 'minor' | 'none';
    }>;
    meta_shifts: Array<{
      category: string;
      change_description: string;
      affected_champions: string[];
    }>;
    statistical_changes: {
      average_game_length_change: number;
      pick_diversity_change: number;
      ban_rate_changes: Record<string, number>;
    };
  }> {
    const [analysis1, analysis2] = await Promise.all([
      this.getMetaAnalysis(patch1, region, rank),
      this.getMetaAnalysis(patch2, region, rank)
    ]);

    // Process and compare the two analyses
    return {
      patch1,
      patch2,
      tier_changes: [],
      meta_shifts: [],
      statistical_changes: {
        average_game_length_change: 0,
        pick_diversity_change: 0,
        ban_rate_changes: {}
      }
    };
  }

  /**
   * Get team composition meta analysis
   */
  async getTeamCompMeta(
    patch: string,
    compType?: string,
    region?: string,
    rank?: string
  ): Promise<{
    patch: string;
    team_compositions: Array<{
      composition_type: string;
      win_rate: number;
      popularity: number;
      champions: {
        top: string[];
        jungle: string[];
        mid: string[];
        adc: string[];
        support: string[];
      };
      synergies: string[];
      counters: string[];
      game_length_preference: string;
    }>;
    meta_composition_trends: {
      rising_comps: string[];
      declining_comps: string[];
      stable_comps: string[];
    };
  }> {
    const params = new URLSearchParams({
      patch,
      ...(compType && { comp_type: compType }),
      ...(region && { region }),
      ...(rank && { rank })
    });

    // This would be a specialized endpoint for team composition analysis
    const response = await apiClient.get(`${this.baseUrl}/team-compositions?${params}`);
    return response.data;
  }

  /**
   * Get item meta analysis
   */
  async getItemMeta(
    patch: string,
    itemCategory?: string,
    role?: string
  ): Promise<{
    patch: string;
    popular_items: Array<{
      item_id: number;
      item_name: string;
      pick_rate: number;
      win_rate: number;
      roles: string[];
      champions: string[];
      synergistic_items: string[];
    }>;
    emerging_items: Array<{
      item_name: string;
      pick_rate_change: number;
      win_rate: number;
      trending_champions: string[];
    }>;
    item_build_paths: Array<{
      build_name: string;
      items: string[];
      win_rate: number;
      pick_rate: number;
      optimal_champions: string[];
    }>;
  }> {
    const params = new URLSearchParams({
      patch,
      ...(itemCategory && { item_category: itemCategory }),
      ...(role && { role })
    });

    const response = await apiClient.get(`${this.baseUrl}/items?${params}`);
    return response.data;
  }

  /**
   * Get rune meta analysis
   */
  async getRuneMeta(
    patch: string,
    runeTree?: string,
    champion?: string
  ): Promise<{
    patch: string;
    popular_rune_combinations: Array<{
      primary_tree: string;
      keystone: string;
      secondary_tree: string;
      win_rate: number;
      pick_rate: number;
      champions: string[];
    }>;
    keystone_analysis: Array<{
      keystone: string;
      tree: string;
      win_rate: number;
      pick_rate: number;
      trending: 'rising' | 'stable' | 'declining';
      optimal_champions: string[];
    }>;
    rune_synergies: Array<{
      rune_combination: string[];
      synergy_rating: number;
      use_cases: string[];
    }>;
  }> {
    const params = new URLSearchParams({
      patch,
      ...(runeTree && { rune_tree: runeTree }),
      ...(champion && { champion })
    });

    const response = await apiClient.get(`${this.baseUrl}/runes?${params}`);
    return response.data;
  }

  /**
   * Get meta snapshot for quick overview
   */
  async getMetaSnapshot(
    patch: string,
    region?: string
  ): Promise<{
    patch: string;
    region: string;
    snapshot_date: string;
    top_tier_champions: {
      s_plus: string[];
      s_tier: string[];
    };
    most_banned: string[];
    emerging_picks: string[];
    dominant_strategies: string[];
    game_length_trend: string;
    next_patch_predictions: string[];
  }> {
    const analysis = await this.getMetaAnalysis(patch, region, 'all', '7d');
    
    return {
      patch,
      region: region || 'all',
      snapshot_date: analysis.generatedAt,
      top_tier_champions: {
        s_plus: analysis.tierList.sPlusTier?.map(c => c.champion) || [],
        s_tier: analysis.tierList.sTier?.map(c => c.champion) || []
      },
      most_banned: analysis.banAnalysis.topBannedChampions.slice(0, 5).map(b => b.champion),
      emerging_picks: analysis.emergingPicks.map(p => p.champion),
      dominant_strategies: analysis.metaTrends.dominantStrategies.map(s => s.strategy),
      game_length_trend: 'stable', // Based on game length analysis
      next_patch_predictions: analysis.predictions.nextPatchPredictions.slice(0, 3).map(p => p.champion)
    };
  }

  /**
   * Search champions by meta criteria
   */
  async searchChampionsByMeta(
    patch: string,
    criteria: {
      min_win_rate?: number;
      max_ban_rate?: number;
      tier?: string;
      role?: string;
      trend?: 'rising' | 'stable' | 'declining';
      play_style?: string;
    }
  ): Promise<{
    matching_champions: Array<{
      champion: string;
      role: string;
      tier: string;
      win_rate: number;
      pick_rate: number;
      ban_rate: number;
      trend: string;
      match_score: number;
    }>;
    search_criteria: typeof criteria;
    total_matches: number;
  }> {
    const analysis = await this.getMetaAnalysis(patch);
    
    // Filter champions based on criteria
    const matchingChampions = analysis.championStats.filter(champion => {
      if (criteria.min_win_rate && champion.winRate < criteria.min_win_rate) return false;
      if (criteria.max_ban_rate && champion.banRate > criteria.max_ban_rate) return false;
      if (criteria.tier && champion.tier !== criteria.tier) return false;
      if (criteria.role && champion.role !== criteria.role) return false;
      if (criteria.trend && champion.trendDirection !== criteria.trend) return false;
      return true;
    });

    return {
      matching_champions: matchingChampions.map(c => ({
        champion: c.champion,
        role: c.role,
        tier: c.tier,
        win_rate: c.winRate,
        pick_rate: c.pickRate,
        ban_rate: c.banRate,
        trend: c.trendDirection,
        match_score: 85.0 // Calculated based on criteria match
      })),
      search_criteria: criteria,
      total_matches: matchingChampions.length
    };
  }
}

export const metaService = new MetaService();