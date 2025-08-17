import { useMemo, useCallback } from 'react';

// Type definitions for optimized data structures
export interface OptimizedMatch {
  id: string;
  date: number;
  champion: string;
  win: boolean;
  kda: number;
  duration: number;
  queue: string;
  lane: string;
}

export interface ChampionStats {
  champion: string;
  games: number;
  wins: number;
  winRate: number;
  avgKDA: number;
  totalKills: number;
  totalDeaths: number;
  totalAssists: number;
  avgDuration: number;
  lastPlayed: number;
}

export interface PerformanceInsights {
  recentTrend: 'improving' | 'declining' | 'stable';
  bestChampions: string[];
  recommendedChampions: string[];
  weakPoints: string[];
  strengths: string[];
}

/**
 * Optimized data processing utilities with memoization and efficient algorithms
 */
export class DataProcessor {
  // Cache for computed results
  private static cache = new Map<string, any>();
  
  // Cache TTL (5 minutes)
  private static CACHE_TTL = 5 * 60 * 1000;

  /**
   * Get cached result or compute if not available/expired
   */
  private static getCached<T>(key: string, computeFn: () => T): T {
    const cached = this.cache.get(key);
    
    if (cached && Date.now() - cached.timestamp < this.CACHE_TTL) {
      return cached.data;
    }
    
    const result = computeFn();
    this.cache.set(key, {
      data: result,
      timestamp: Date.now()
    });
    
    return result;
  }

  /**
   * Process raw match data into optimized format
   */
  static processMatches(rawData: any[]): OptimizedMatch[] {
    const cacheKey = `processed_matches_${rawData.length}_${rawData[0]?.matchId || 'empty'}`;
    
    return this.getCached(cacheKey, () => {
      return rawData.map(match => ({
        id: match.matchId || match.id,
        date: new Date(match.date || match.gameCreation).getTime(),
        champion: match.champion || match.championName,
        win: Boolean(match.win),
        kda: this.calculateKDA(match.kills, match.deaths, match.assists),
        duration: match.duration_s || match.gameDuration || 0,
        queue: match.queue || match.queueId?.toString() || 'Unknown',
        lane: match.lane || match.teamPosition || 'Unknown'
      }));
    });
  }

  /**
   * Calculate KDA with death handling
   */
  private static calculateKDA(kills: number, deaths: number, assists: number): number {
    const k = kills || 0;
    const d = deaths || 0;
    const a = assists || 0;
    
    return d === 0 ? k + a : Number(((k + a) / d).toFixed(2));
  }

  /**
   * Compute champion statistics with efficient grouping
   */
  static computeChampionStats(matches: OptimizedMatch[]): ChampionStats[] {
    const cacheKey = `champion_stats_${matches.length}_${matches[0]?.date || 0}`;
    
    return this.getCached(cacheKey, () => {
      const statsMap = new Map<string, {
        games: number;
        wins: number;
        totalKDA: number;
        totalDuration: number;
        kills: number;
        deaths: number;
        assists: number;
        lastPlayed: number;
      }>();

      // Single pass through matches for efficiency
      for (const match of matches) {
        const current = statsMap.get(match.champion) || {
          games: 0,
          wins: 0,
          totalKDA: 0,
          totalDuration: 0,
          kills: 0,
          deaths: 0,
          assists: 0,
          lastPlayed: 0
        };

        current.games += 1;
        current.wins += match.win ? 1 : 0;
        current.totalKDA += match.kda;
        current.totalDuration += match.duration;
        current.lastPlayed = Math.max(current.lastPlayed, match.date);

        statsMap.set(match.champion, current);
      }

      // Convert to array with computed averages
      return Array.from(statsMap.entries()).map(([champion, stats]) => ({
        champion,
        games: stats.games,
        wins: stats.wins,
        winRate: Number((stats.wins / stats.games * 100).toFixed(1)),
        avgKDA: Number((stats.totalKDA / stats.games).toFixed(2)),
        totalKills: stats.kills,
        totalDeaths: stats.deaths,
        totalAssists: stats.assists,
        avgDuration: Math.round(stats.totalDuration / stats.games),
        lastPlayed: stats.lastPlayed
      }));
    });
  }

  /**
   * Generate performance insights using AI-like analysis
   */
  static generateInsights(matches: OptimizedMatch[], championStats: ChampionStats[]): PerformanceInsights {
    const cacheKey = `insights_${matches.length}_${championStats.length}`;
    
    return this.getCached(cacheKey, () => {
      // Analyze recent trend (last 10 matches)
      const recentMatches = matches.slice(-10);
      const recentWinRate = recentMatches.filter(m => m.win).length / recentMatches.length;
      const overallWinRate = matches.filter(m => m.win).length / matches.length;
      
      let recentTrend: 'improving' | 'declining' | 'stable' = 'stable';
      if (recentWinRate > overallWinRate + 0.1) recentTrend = 'improving';
      else if (recentWinRate < overallWinRate - 0.1) recentTrend = 'declining';

      // Best champions (high win rate + reasonable games)
      const bestChampions = championStats
        .filter(c => c.games >= 3)
        .sort((a, b) => b.winRate - a.winRate)
        .slice(0, 3)
        .map(c => c.champion);

      // Recommended champions (good performance, not played recently)
      const daysSinceEpoch = Date.now() / (1000 * 60 * 60 * 24);
      const recommendedChampions = championStats
        .filter(c => c.winRate > 60 && (daysSinceEpoch - c.lastPlayed / (1000 * 60 * 60 * 24)) > 7)
        .slice(0, 2)
        .map(c => c.champion);

      // Analyze weak points
      const weakPoints: string[] = [];
      const avgKDA = championStats.reduce((sum, c) => sum + c.avgKDA, 0) / championStats.length;
      if (avgKDA < 1.5) weakPoints.push('KDA could be improved');
      
      const avgWinRate = championStats.reduce((sum, c) => sum + c.winRate, 0) / championStats.length;
      if (avgWinRate < 50) weakPoints.push('Overall win rate needs improvement');

      // Identify strengths
      const strengths: string[] = [];
      if (avgKDA > 2.0) strengths.push('Excellent KDA performance');
      if (avgWinRate > 60) strengths.push('Strong win rate across champions');
      if (recentTrend === 'improving') strengths.push('Improving recent performance');

      return {
        recentTrend,
        bestChampions,
        recommendedChampions,
        weakPoints,
        strengths
      };
    });
  }

  /**
   * Filter matches by date range efficiently
   */
  static filterMatchesByDateRange(
    matches: OptimizedMatch[],
    startDate: Date,
    endDate: Date
  ): OptimizedMatch[] {
    const start = startDate.getTime();
    const end = endDate.getTime();
    
    // Use binary search if matches are sorted by date
    return matches.filter(match => match.date >= start && match.date <= end);
  }

  /**
   * Get performance summary for time periods
   */
  static getPerformanceSummary(matches: OptimizedMatch[]) {
    const now = Date.now();
    const day = 24 * 60 * 60 * 1000;
    
    const timeRanges = {
      last7Days: matches.filter(m => now - m.date <= 7 * day),
      last30Days: matches.filter(m => now - m.date <= 30 * day),
      last90Days: matches.filter(m => now - m.date <= 90 * day)
    };

    return Object.entries(timeRanges).reduce((summary, [period, periodMatches]) => {
      const wins = periodMatches.filter(m => m.win).length;
      const games = periodMatches.length;
      
      summary[period] = {
        games,
        wins,
        winRate: games > 0 ? Number((wins / games * 100).toFixed(1)) : 0,
        avgKDA: games > 0 ? Number((periodMatches.reduce((sum, m) => sum + m.kda, 0) / games).toFixed(2)) : 0
      };
      
      return summary;
    }, {} as Record<string, any>);
  }

  /**
   * Clear cache manually
   */
  static clearCache(): void {
    this.cache.clear();
  }

  /**
   * Get cache statistics
   */
  static getCacheStats(): { size: number; keys: string[] } {
    return {
      size: this.cache.size,
      keys: Array.from(this.cache.keys())
    };
  }
}

/**
 * React hook for optimized data processing
 */
export const useOptimizedData = (rawData: any[]) => {
  const matches = useMemo(
    () => DataProcessor.processMatches(rawData),
    [rawData]
  );

  const championStats = useMemo(
    () => DataProcessor.computeChampionStats(matches),
    [matches]
  );

  const insights = useMemo(
    () => DataProcessor.generateInsights(matches, championStats),
    [matches, championStats]
  );

  const performanceSummary = useMemo(
    () => DataProcessor.getPerformanceSummary(matches),
    [matches]
  );

  const filterByDateRange = useCallback(
    (startDate: Date, endDate: Date) => 
      DataProcessor.filterMatchesByDateRange(matches, startDate, endDate),
    [matches]
  );

  return {
    matches,
    championStats,
    insights,
    performanceSummary,
    filterByDateRange,
    clearCache: DataProcessor.clearCache,
    cacheStats: DataProcessor.getCacheStats()
  };
};

export default DataProcessor;
