// Team Composition Service for Herald.lol Frontend
import axiosInstance from './api';
import {
  TeamCompositionOptimization,
  CompositionOptimizationRequest,
  CompositionAnalysis,
  MetaComposition,
  PlayerComfortData,
  BanStrategy,
  DraftOptimization,
  DraftState,
  PlayerOptimizationData,
  OptimizationStrategy,
  OptimizationConstraints,
  OptimizationPreferences,
  CompositionRecommendation,
  ChampionRole,
  WinCondition,
  SynergyAnalysis,
  CounterAnalysis,
  ScalingProfile
} from '../types/teamComposition';

class TeamCompositionService {
  private readonly baseURL = '/team-composition';

  // Team Composition Optimization
  async optimizeTeamComposition(request: CompositionOptimizationRequest): Promise<TeamCompositionOptimization> {
    const response = await axiosInstance.post(`${this.baseURL}/optimize`, request);
    return response.data;
  }

  async analyzeComposition(
    blueTeam: Array<{ champion: string; role: string }>,
    redTeam?: Array<{ champion: string; role: string }>,
    gameMode: string = 'ranked',
    patch?: string
  ): Promise<CompositionAnalysis> {
    const response = await axiosInstance.post(`${this.baseURL}/analyze`, {
      blueTeam,
      redTeam,
      gameMode,
      patch
    });
    return response.data;
  }

  async getCompositionSuggestions(
    summonerId: string,
    role?: string,
    gameMode: string = 'ranked',
    strategy: OptimizationStrategy = 'balanced',
    limit: number = 10
  ): Promise<CompositionRecommendation[]> {
    const params = new URLSearchParams({
      gameMode,
      strategy,
      limit: limit.toString()
    });
    
    if (role) {
      params.append('role', role);
    }

    const response = await axiosInstance.get(`${this.baseURL}/suggestions/${summonerId}?${params}`);
    return response.data;
  }

  async validateComposition(
    composition: Array<{ champion: string; role: string }>,
    gameMode: string = 'ranked'
  ): Promise<{
    isValid: boolean;
    issues: string[];
    suggestions: string[];
    score: number;
  }> {
    const response = await axiosInstance.post(`${this.baseURL}/validate`, {
      composition,
      gameMode
    });
    return response.data;
  }

  async compareCompositions(
    compositions: Array<{
      name: string;
      champions: Array<{ champion: string; role: string }>;
    }>,
    gameMode: string = 'ranked',
    criteria: string[] = ['synergy', 'scaling', 'teamfight']
  ): Promise<{
    comparisons: Array<{
      name: string;
      score: number;
      ranking: number;
      strengths: string[];
      weaknesses: string[];
    }>;
    winner: string;
    analysis: string;
  }> {
    const response = await axiosInstance.post(`${this.baseURL}/compare`, {
      compositions,
      gameMode,
      criteria
    });
    return response.data;
  }

  // Meta Compositions
  async getMetaCompositions(
    gameMode: string = 'ranked',
    rank: string = 'all',
    region: string = 'global',
    patch?: string,
    limit: number = 20
  ): Promise<MetaComposition[]> {
    const params = new URLSearchParams({
      gameMode,
      rank,
      region,
      limit: limit.toString()
    });
    
    if (patch) {
      params.append('patch', patch);
    }

    const response = await axiosInstance.get(`${this.baseURL}/meta-compositions?${params}`);
    return response.data;
  }

  // Synergy and Counter Analysis
  async analyzeTeamSynergy(
    champions: Array<{ champion: string; role: string }>,
    synergyType?: string
  ): Promise<SynergyAnalysis> {
    const response = await axiosInstance.post(`${this.baseURL}/synergy-analysis`, {
      champions,
      synergyType
    });
    return response.data;
  }

  async analyzeCounters(
    enemyComposition: Array<{ champion: string; role: string }>,
    availableChampions: string[],
    targetRole?: string,
    counterType?: string
  ): Promise<CounterAnalysis> {
    const response = await axiosInstance.post(`${this.baseURL}/counter-analysis`, {
      enemyComposition,
      availableChampions,
      targetRole,
      counterType
    });
    return response.data;
  }

  // Role-Specific Recommendations
  async getRoleRecommendations(
    summonerId: string,
    role: string,
    existingChampions: string[] = [],
    gameMode: string = 'ranked',
    strategy: OptimizationStrategy = 'balanced',
    limit: number = 15
  ): Promise<{
    recommendations: Array<{
      champion: string;
      score: number;
      reasoning: string[];
      synergy: number;
      metaFit: number;
      comfort: number;
      risk: 'low' | 'medium' | 'high';
    }>;
    analysis: {
      roleNeeds: string[];
      teamGaps: string[];
      suggestions: string[];
    };
  }> {
    const params = new URLSearchParams({
      gameMode,
      strategy,
      limit: limit.toString()
    });

    existingChampions.forEach(champion => {
      params.append('existing_champions', champion);
    });

    const response = await axiosInstance.get(`${this.baseURL}/role-recommendations/${summonerId}/${role}?${params}`);
    return response.data;
  }

  // Draft Optimization
  async optimizeDraftPicks(
    draftState: DraftState,
    playerData: PlayerOptimizationData[],
    gameMode: string = 'ranked',
    strategy: OptimizationStrategy = 'balanced'
  ): Promise<DraftOptimization> {
    const response = await axiosInstance.post(`${this.baseURL}/draft-optimization`, {
      draftState,
      playerData,
      gameMode,
      strategy
    });
    return response.data;
  }

  // Player Comfort and Champion Pools
  async getPlayerComfortPicks(
    summonerId: string,
    role?: string,
    gameMode: string = 'ranked',
    recentGames: number = 50,
    limit: number = 20
  ): Promise<PlayerComfortData> {
    const params = new URLSearchParams({
      gameMode,
      recentGames: recentGames.toString(),
      limit: limit.toString()
    });
    
    if (role) {
      params.append('role', role);
    }

    const response = await axiosInstance.get(`${this.baseURL}/player-comfort/${summonerId}?${params}`);
    return response.data;
  }

  async getChampionPools(
    summonerId: string,
    role?: string,
    gameMode: string = 'ranked',
    poolType: 'comfort' | 'meta' | 'flex' | 'wide' = 'comfort',
    recentGames: number = 100
  ): Promise<{
    championPools: Record<string, PlayerComfortData>;
    analysis: {
      poolDepth: number;
      metaAlignment: number;
      flexibility: number;
      gaps: string[];
      recommendations: string[];
    };
  }> {
    const params = new URLSearchParams({
      gameMode,
      poolType,
      recentGames: recentGames.toString()
    });
    
    if (role) {
      params.append('role', role);
    }

    const response = await axiosInstance.get(`${this.baseURL}/champion-pools/${summonerId}?${params}`);
    return response.data;
  }

  // Strategic Analysis
  async analyzeWinConditions(
    teamComposition: Array<{ champion: string; role: string }>,
    enemyComposition?: Array<{ champion: string; role: string }>,
    gameMode: string = 'ranked'
  ): Promise<{
    winConditions: WinCondition[];
    analysis: {
      primaryPath: string;
      alternativePaths: string[];
      risks: string[];
      opportunities: string[];
    };
    timeline: {
      early: string[];
      mid: string[];
      late: string[];
    };
  }> {
    const response = await axiosInstance.post(`${this.baseURL}/win-condition-analysis`, {
      teamComposition,
      enemyComposition,
      gameMode
    });
    return response.data;
  }

  async analyzeScaling(
    teamComposition: Array<{ champion: string; role: string }>,
    compareAgainst?: Array<{ champion: string; role: string }>,
    gameMode: string = 'ranked'
  ): Promise<{
    scalingProfile: ScalingProfile;
    analysis: {
      advantage: string[];
      disadvantage: string[];
      recommendations: string[];
    };
    comparison?: {
      earlyAdvantage: number;
      midAdvantage: number;
      lateAdvantage: number;
      overallAdvantage: number;
    };
  }> {
    const response = await axiosInstance.post(`${this.baseURL}/scaling-analysis`, {
      teamComposition,
      compareAgainst,
      gameMode
    });
    return response.data;
  }

  // Ban Strategy
  async getBanStrategy(
    playerData: Array<{ summonerId: string; role: string }>,
    enemyData: Array<{ summonerId?: string; role: string }> = [],
    banPhase: 'first_ban' | 'second_ban' = 'first_ban',
    existingBans: string[] = [],
    gameMode: string = 'ranked',
    strategy: 'target_player' | 'protect_comp' | 'meta_deny' = 'balanced'
  ): Promise<BanStrategy> {
    const response = await axiosInstance.post(`${this.baseURL}/ban-strategy`, {
      playerData,
      enemyData,
      banPhase,
      existingBans,
      gameMode,
      strategy
    });
    return response.data;
  }

  // Utility Methods
  async getOptimizationStrategies(): Promise<Array<{
    value: OptimizationStrategy;
    label: string;
    description: string;
    useCase: string[];
  }>> {
    // Static data for optimization strategies
    return [
      {
        value: 'meta_optimal',
        label: 'Meta Optimal',
        description: 'Focus on current meta strength and proven compositions',
        useCase: ['Ranked climbing', 'High-stakes matches', 'Following pro meta']
      },
      {
        value: 'synergy_focused',
        label: 'Synergy Focused',
        description: 'Maximize team synergy and combo potential',
        useCase: ['Team coordination', 'Learning team fighting', 'Combo practice']
      },
      {
        value: 'balanced',
        label: 'Balanced',
        description: 'Balance all factors including meta, synergy, and comfort',
        useCase: ['General gameplay', 'Versatile strategy', 'All-around improvement']
      },
      {
        value: 'comfort_picks',
        label: 'Comfort Picks',
        description: 'Prioritize player comfort and champion mastery',
        useCase: ['Learning fundamentals', 'Consistency focus', 'New players']
      },
      {
        value: 'counter_focused',
        label: 'Counter Focused',
        description: 'Focus on countering enemy composition',
        useCase: ['Draft phase', 'Known enemy picks', 'Reactive strategy']
      },
      {
        value: 'scaling_focused',
        label: 'Late Game Focused',
        description: 'Focus on late game team compositions',
        useCase: ['Scaling practice', 'Late game team fights', 'Patience strategy']
      },
      {
        value: 'early_focused',
        label: 'Early Game Focused',
        description: 'Focus on early game pressure and snowball',
        useCase: ['Quick games', 'Aggressive playstyle', 'Lane dominance']
      }
    ];
  }

  async getPlayStyles(): Promise<Array<{
    value: string;
    label: string;
    description: string;
    champions: string[];
    difficulty: 'beginner' | 'intermediate' | 'advanced';
  }>> {
    // Static data for play styles
    return [
      {
        value: 'team_fight',
        label: 'Team Fight',
        description: 'Focus on 5v5 team fighting around objectives',
        champions: ['Malphite', 'Orianna', 'Jinx', 'Braum'],
        difficulty: 'beginner'
      },
      {
        value: 'split_push',
        label: 'Split Push',
        description: 'Create pressure through splitting and side lane control',
        champions: ['Fiora', 'Tryndamere', 'Jax', 'Camille'],
        difficulty: 'advanced'
      },
      {
        value: 'poke_siege',
        label: 'Poke & Siege',
        description: 'Control range and siege objectives with poke damage',
        champions: ['Jayce', 'Xerath', 'Caitlyn', 'Karma'],
        difficulty: 'intermediate'
      },
      {
        value: 'pick_comp',
        label: 'Pick Composition',
        description: 'Focus on catching isolated enemies and creating picks',
        champions: ['Thresh', 'LeBlanc', 'Rengar', 'Pyke'],
        difficulty: 'advanced'
      },
      {
        value: 'protect_adc',
        label: 'Protect the ADC',
        description: 'Focus on protecting and enabling the ADC',
        champions: ['Lulu', 'Janna', 'Kog\'Maw', 'Maokai'],
        difficulty: 'intermediate'
      },
      {
        value: 'engage_comp',
        label: 'Engage Composition',
        description: 'Focus on hard engage and team fight initiation',
        champions: ['Malphite', 'Amumu', 'Leona', 'Jarvan IV'],
        difficulty: 'beginner'
      }
    ];
  }

  // Performance Analysis
  async analyzeCompositionPerformance(
    compositions: string[],
    timeframe: '7d' | '30d' | '90d' = '30d',
    rank?: string,
    region: string = 'global'
  ): Promise<{
    performance: Array<{
      composition: string;
      winRate: number;
      playRate: number;
      trend: 'rising' | 'stable' | 'declining';
      avgGameLength: number;
      strengthPhases: string[];
    }>;
    meta: {
      topPerformers: string[];
      emerging: string[];
      declining: string[];
      stable: string[];
    };
  }> {
    const params = new URLSearchParams({
      timeframe,
      region
    });
    
    if (rank) {
      params.append('rank', rank);
    }

    compositions.forEach(comp => {
      params.append('compositions', comp);
    });

    const response = await axiosInstance.get(`${this.baseURL}/performance-analysis?${params}`);
    return response.data;
  }

  // Export utilities for chart data transformation
  createSynergyRadarData(synergy: SynergyAnalysis) {
    return {
      labels: ['Combo Potential', 'Engage Synergy', 'Protection', 'Damage Amp', 'Utility Chain'],
      data: [
        synergy.comboPotential?.reduce((sum, combo) => sum + combo.damage, 0) / (synergy.comboPotential?.length || 1) || 0,
        synergy.chainSynergies?.filter(chain => chain.type === 'engage').length * 20 || 0,
        synergy.chainSynergies?.filter(chain => chain.type === 'protection').length * 25 || 0,
        synergy.chainSynergies?.filter(chain => chain.type === 'amplification').length * 20 || 0,
        synergy.overall || 0
      ],
      backgroundColor: 'rgba(54, 162, 235, 0.2)',
      borderColor: 'rgba(54, 162, 235, 1)',
      pointBackgroundColor: 'rgba(54, 162, 235, 1)'
    };
  }

  createScalingCurveData(profile: ScalingProfile) {
    return {
      timeline: profile.scalingCurve.map(point => point.minute),
      teamPower: profile.scalingCurve.map(point => point.powerLevel),
      powerSpikes: profile.criticalTimings.map(timing => ({
        time: parseInt(timing.timing.split(' ')[0]) || 0, // Extract minute from timing string
        power: timing.importance === 'critical' ? 90 : timing.importance === 'important' ? 70 : 50,
        description: timing.description
      }))
    };
  }
}

export default new TeamCompositionService();