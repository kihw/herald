import { apiClient } from './apiClient';
import type {
  ImprovementRecommendation,
  RecommendationOptions,
  RecommendationResponse,
  PlayerAnalysisResult,
  ProgressTrackingData,
  ProgressUpdateRequest,
  ImprovementInsight,
  OverallTrajectory,
  OverallProgress,
  QuickWin,
  CoachingPlan
} from '../types/improvement';

class ImprovementService {
  private baseUrl = '/api/v1/improvement';

  /**
   * Get personalized improvement recommendations
   */
  async getPersonalizedRecommendations(
    summonerId: string,
    options?: Partial<RecommendationOptions>
  ): Promise<RecommendationResponse> {
    const defaultOptions: RecommendationOptions = {
      maxRecommendations: 10,
      includeAlternatives: false,
    };

    const finalOptions = { ...defaultOptions, ...options };
    const params = new URLSearchParams();

    Object.entries(finalOptions).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          params.append(key, value.join(','));
        } else {
          params.append(key, value.toString());
        }
      }
    });

    const response = await apiClient.get(
      `${this.baseUrl}/recommendations/${summonerId}?${params}`
    );
    return response.data;
  }

  /**
   * Get active recommendations for a summoner
   */
  async getActiveRecommendations(summonerId: string): Promise<ImprovementRecommendation[]> {
    const response = await apiClient.get(
      `${this.baseUrl}/recommendations/${summonerId}/active`
    );
    return response.data.recommendations;
  }

  /**
   * Get details for a specific recommendation
   */
  async getRecommendationDetails(recommendationId: string): Promise<{
    recommendationId: string;
    progress: ProgressTrackingData;
  }> {
    const response = await apiClient.get(
      `${this.baseUrl}/recommendation/${recommendationId}`
    );
    return response.data;
  }

  /**
   * Update progress on a recommendation
   */
  async updateRecommendationProgress(
    recommendationId: string,
    progressUpdate: ProgressUpdateRequest
  ): Promise<void> {
    const progressData: ProgressTrackingData = {
      overallProgress: progressUpdate.overallProgress,
      completedSteps: progressUpdate.completedSteps,
      currentMilestone: progressUpdate.currentMilestone,
      milestoneProgress: [],
      weeklyProgress: [],
      performanceImpact: {
        baselineMetrics: {},
        currentMetrics: progressUpdate.performanceMetrics,
        improvementDeltas: {},
        roiActual: 0,
        roiPredicted: 0,
        confidenceInterval: 0.8,
      },
      lastProgressUpdate: new Date().toISOString(),
    };

    await apiClient.put(
      `${this.baseUrl}/recommendation/${recommendationId}/progress`,
      progressData
    );
  }

  /**
   * Mark a recommendation as completed
   */
  async completeRecommendation(recommendationId: string): Promise<void> {
    await apiClient.post(
      `${this.baseUrl}/recommendation/${recommendationId}/complete`
    );
  }

  /**
   * Get comprehensive player analysis
   */
  async getPlayerAnalysis(summonerId: string): Promise<{
    analysis: PlayerAnalysisResult;
  }> {
    const response = await apiClient.get(
      `${this.baseUrl}/analysis/${summonerId}`
    );
    return response.data;
  }

  /**
   * Get improvement insights for a summoner
   */
  async getImprovementInsights(summonerId: string): Promise<{
    summonerId: string;
    insights: ImprovementInsight[];
    overallTrajectory: OverallTrajectory;
  }> {
    const response = await apiClient.get(
      `${this.baseUrl}/insights/${summonerId}`
    );
    return response.data;
  }

  /**
   * Get overall progress tracking for a summoner
   */
  async getOverallProgress(summonerId: string): Promise<OverallProgress> {
    const response = await apiClient.get(
      `${this.baseUrl}/progress/${summonerId}`
    );
    return response.data;
  }

  /**
   * Get quick win recommendations
   */
  async getQuickWins(summonerId: string): Promise<{
    summonerId: string;
    quickWins: QuickWin[];
    implementationPlan: {
      week1: string[];
      week2: string[];
      expectedCumulativeImpact: string;
      successMetrics: string[];
    };
  }> {
    const response = await apiClient.get(
      `${this.baseUrl}/quick-wins/${summonerId}`
    );
    return response.data;
  }

  /**
   * Get comprehensive coaching plan
   */
  async getCoachingPlan(
    summonerId: string,
    duration?: number,
    intensity?: 'light' | 'moderate' | 'intensive'
  ): Promise<CoachingPlan> {
    const params = new URLSearchParams();
    if (duration) params.append('duration', duration.toString());
    if (intensity) params.append('intensity', intensity);

    const response = await apiClient.get(
      `${this.baseUrl}/coaching-plan/${summonerId}?${params}`
    );
    
    return {
      summonerId,
      planDetails: response.data.plan_details,
      weeklyBreakdown: response.data.weekly_breakdown,
      dailyRoutine: response.data.daily_routine,
      progressTracking: response.data.progress_tracking,
    };
  }

  /**
   * Get recommendations by category
   */
  async getRecommendationsByCategory(
    summonerId: string,
    category: string
  ): Promise<ImprovementRecommendation[]> {
    const options: RecommendationOptions = {
      focusCategory: category,
      maxRecommendations: 5,
      includeAlternatives: true,
    };

    const response = await this.getPersonalizedRecommendations(summonerId, options);
    return response.recommendations;
  }

  /**
   * Get recommendations by difficulty level
   */
  async getRecommendationsByDifficulty(
    summonerId: string,
    difficulty: 'easy' | 'medium' | 'hard' | 'expert'
  ): Promise<ImprovementRecommendation[]> {
    const options: RecommendationOptions = {
      difficultyFilter: difficulty,
      maxRecommendations: 8,
      includeAlternatives: false,
    };

    const response = await this.getPersonalizedRecommendations(summonerId, options);
    return response.recommendations;
  }

  /**
   * Get skill-specific improvement plan
   */
  async getSkillImprovementPlan(
    summonerId: string,
    skill: string
  ): Promise<{
    skill: string;
    currentLevel: number;
    targetLevel: number;
    improvementPlan: ImprovementRecommendation[];
    timeEstimate: string;
    milestones: Array<{
      level: number;
      timeframe: string;
      requirements: string[];
    }>;
  }> {
    // This combines multiple API calls to create a comprehensive skill improvement plan
    const [recommendations, analysis] = await Promise.all([
      this.getRecommendationsByCategory(summonerId, skill),
      this.getPlayerAnalysis(summonerId),
    ]);

    const currentLevel = analysis.analysis.skillBreakdown[skill] || 50;
    const improvementPotential = analysis.analysis.improvementPotential[skill] || 10;
    const targetLevel = Math.min(100, currentLevel + improvementPotential);

    return {
      skill,
      currentLevel,
      targetLevel,
      improvementPlan: recommendations,
      timeEstimate: this.estimateImprovementTime(currentLevel, targetLevel),
      milestones: this.generateSkillMilestones(skill, currentLevel, targetLevel),
    };
  }

  /**
   * Get habit formation recommendations
   */
  async getHabitFormationPlan(
    summonerId: string,
    habits: string[]
  ): Promise<{
    summonerId: string;
    targetHabits: string[];
    formationPlan: Array<{
      habit: string;
      trigger: string;
      routine: string;
      reward: string;
      timeToForm: number; // days
      trackingMethod: string;
      tips: string[];
    }>;
    weeklySchedule: Array<{
      week: number;
      focus: string[];
      activities: string[];
      successMetrics: string[];
    }>;
  }> {
    // For now, return a structured habit formation plan based on common gaming habits
    const habitDatabase = {
      vision_control: {
        trigger: 'Before every objective spawn',
        routine: 'Place control ward and sweep area',
        reward: 'Better map control and safety',
        timeToForm: 21,
        trackingMethod: 'Vision score per game',
        tips: [
          'Set a phone reminder 30 seconds before objectives',
          'Always buy control wards on every back',
          'Practice in normal games first',
        ],
      },
      map_awareness: {
        trigger: 'After every CS',
        routine: 'Quick glance at minimap',
        reward: 'Avoid ganks and spot opportunities',
        timeToForm: 14,
        trackingMethod: 'Deaths to ganks per game',
        tips: [
          'Count your map checks per minute',
          'Use a metronome or timer initially',
          'Make it a verbal habit - say "map" after each CS',
        ],
      },
      positioning: {
        trigger: 'When entering team fights',
        routine: 'Check for escape routes and threats',
        reward: 'Survive team fights and deal damage',
        timeToForm: 28,
        trackingMethod: 'Deaths in team fights',
        tips: [
          'Practice positioning in ARAM',
          'Watch replays focusing only on positioning',
          'Use max camera distance',
        ],
      },
    };

    const formationPlan = habits.map(habit => ({
      habit,
      ...(habitDatabase[habit as keyof typeof habitDatabase] || {
        trigger: 'Custom trigger',
        routine: 'Custom routine',
        reward: 'Improved gameplay',
        timeToForm: 21,
        trackingMethod: 'Custom tracking',
        tips: ['Practice consistently', 'Track your progress', 'Be patient with formation'],
      }),
    }));

    const weeklySchedule = [
      {
        week: 1,
        focus: habits.slice(0, 1), // Start with one habit
        activities: ['Establish trigger awareness', 'Practice routine 10x daily'],
        successMetrics: ['80% routine completion', 'Trigger recognition'],
      },
      {
        week: 2,
        focus: habits.slice(0, 1),
        activities: ['Refine routine execution', 'Track performance impact'],
        successMetrics: ['90% routine completion', 'Measurable improvement'],
      },
      {
        week: 3,
        focus: habits.slice(0, 2), // Add second habit if applicable
        activities: ['Maintain first habit', 'Introduce second habit'],
        successMetrics: ['Consistent habit 1', 'Habit 2 trigger awareness'],
      },
      {
        week: 4,
        focus: habits,
        activities: ['Stack all habits', 'Monitor for habit conflicts'],
        successMetrics: ['All habits >70% completion', 'No performance regression'],
      },
    ];

    return {
      summonerId,
      targetHabits: habits,
      formationPlan,
      weeklySchedule,
    };
  }

  /**
   * Generate personalized practice routine
   */
  async generatePracticeRoutine(
    summonerId: string,
    availableTime: number, // minutes per day
    focusAreas: string[]
  ): Promise<{
    summonerId: string;
    dailyTimeAvailable: number;
    routine: {
      warmUp: {
        duration: number;
        activities: Array<{
          name: string;
          duration: number;
          description: string;
          tools: string[];
        }>;
      };
      skillPractice: {
        duration: number;
        activities: Array<{
          skill: string;
          duration: number;
          exercises: string[];
          successCriteria: string[];
        }>;
      };
      gameApplication: {
        duration: number;
        gamesCount: number;
        focus: string[];
        trackingMetrics: string[];
      };
      coolDown: {
        duration: number;
        activities: string[];
      };
    };
    weeklyProgression: Array<{
      week: number;
      adjustments: string[];
      newChallenges: string[];
    }>;
  }> {
    const timeAllocation = {
      warmUp: Math.max(5, Math.floor(availableTime * 0.15)),
      skillPractice: Math.floor(availableTime * 0.25),
      gameApplication: Math.floor(availableTime * 0.55),
      coolDown: Math.max(5, Math.floor(availableTime * 0.05)),
    };

    const skillExercises = {
      mechanical_skill: [
        'Last-hit practice in practice tool',
        'Combo execution drills',
        'Animation canceling practice',
        'Kiting and orbwalking drills',
      ],
      vision_control: [
        'Ward placement timing practice',
        'Vision sweep patterns',
        'Control ward positioning drills',
        'Objective vision setups',
      ],
      positioning: [
        'Team fight positioning scenarios',
        'Laning position optimization',
        'Escape route planning',
        'Threat assessment drills',
      ],
      map_awareness: [
        'Minimap attention training',
        'Enemy tracking exercises',
        'Opportunity recognition drills',
        'Information processing practice',
      ],
    };

    return {
      summonerId,
      dailyTimeAvailable: availableTime,
      routine: {
        warmUp: {
          duration: timeAllocation.warmUp,
          activities: [
            {
              name: 'Mechanics Review',
              duration: Math.ceil(timeAllocation.warmUp * 0.6),
              description: 'Quick mechanical skill refresh',
              tools: ['Practice Tool', 'Custom Game'],
            },
            {
              name: 'Game Plan Setting',
              duration: Math.floor(timeAllocation.warmUp * 0.4),
              description: 'Set focus goals for the session',
              tools: ['Notebook', 'Goal Tracker'],
            },
          ],
        },
        skillPractice: {
          duration: timeAllocation.skillPractice,
          activities: focusAreas.map(skill => ({
            skill,
            duration: Math.floor(timeAllocation.skillPractice / focusAreas.length),
            exercises: skillExercises[skill as keyof typeof skillExercises] || ['General skill practice'],
            successCriteria: [`Improvement in ${skill} metric`, 'Consistent execution', 'Confidence building'],
          })),
        },
        gameApplication: {
          duration: timeAllocation.gameApplication,
          gamesCount: Math.max(1, Math.floor(timeAllocation.gameApplication / 25)), // ~25 min per game
          focus: focusAreas,
          trackingMetrics: focusAreas.map(skill => `${skill}_performance`),
        },
        coolDown: {
          duration: timeAllocation.coolDown,
          activities: [
            'Quick session review',
            'Note key learnings',
            'Plan next session focus',
          ],
        },
      },
      weeklyProgression: [
        {
          week: 1,
          adjustments: ['Establish routine consistency'],
          newChallenges: ['Basic skill execution'],
        },
        {
          week: 2,
          adjustments: ['Increase practice intensity'],
          newChallenges: ['Apply skills under pressure'],
        },
        {
          week: 3,
          adjustments: ['Add situational variations'],
          newChallenges: ['Combine multiple skills'],
        },
        {
          week: 4,
          adjustments: ['Focus on weak points'],
          newChallenges: ['Competitive application'],
        },
      ],
    };
  }

  // Helper methods
  private estimateImprovementTime(currentLevel: number, targetLevel: number): string {
    const improvement = targetLevel - currentLevel;
    const weeks = Math.ceil(improvement / 2); // Assume ~2 points improvement per week
    
    if (weeks <= 4) return `${weeks} weeks`;
    if (weeks <= 8) return `${Math.ceil(weeks / 4)} months`;
    return `${Math.ceil(weeks / 12)} seasons`;
  }

  private generateSkillMilestones(skill: string, currentLevel: number, targetLevel: number) {
    const improvement = targetLevel - currentLevel;
    const milestones = [];
    const stepSize = improvement / 4; // 4 milestones

    for (let i = 1; i <= 4; i++) {
      const milestoneLevel = Math.round(currentLevel + (stepSize * i));
      milestones.push({
        level: milestoneLevel,
        timeframe: `${i * 2} weeks`,
        requirements: [
          `Consistent practice in ${skill}`,
          `Apply learnings in ranked games`,
          `Track improvement metrics`,
        ],
      });
    }

    return milestones;
  }
}

export const improvementService = new ImprovementService();