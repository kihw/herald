# Gaming Analytics Developer Guide

This comprehensive guide covers everything you need to know about integrating Herald.lol's gaming analytics into your applications, with detailed examples and best practices for League of Legends and TFT analytics.

## 🎮 Gaming Analytics Overview

Herald.lol's gaming analytics provide deep insights into player performance, match statistics, and competitive metrics. Our analytics engine processes millions of gaming data points to deliver actionable insights that help players improve their gameplay.

### Core Gaming Metrics

#### 1. KDA Analysis (Kill/Death/Assist)
```javascript
// Example KDA data structure
{
  "kda": {
    "kills": 7.2,           // Average kills per game
    "deaths": 5.1,          // Average deaths per game  
    "assists": 9.4,         // Average assists per game
    "ratio": 3.25,          // (Kills + Assists) / Deaths
    "kill_participation": 0.68, // Percentage of team kills participated in
    "first_blood_rate": 0.15,   // Percentage of games with first blood
    "trends": {
      "7d_change": 0.23,    // KDA improvement over 7 days
      "30d_change": 0.15    // KDA improvement over 30 days
    }
  }
}
```

#### 2. CS/min (Creep Score per Minute)
```javascript
{
  "cs_per_minute": {
    "overall": 6.8,         // Overall CS/min average
    "laning_phase": 7.2,    // CS/min during laning phase (0-15 min)
    "mid_game": 6.1,        // CS/min during mid game (15-25 min)
    "late_game": 5.9,       // CS/min during late game (25+ min)
    "by_role": {
      "adc": 8.2,           // Role-specific CS/min
      "mid": 7.1,
      "top": 6.8
    },
    "percentile_rank": 75,  // Percentile compared to same rank
    "improvement_trend": 0.12 // Recent improvement percentage
  }
}
```

#### 3. Vision Score
```javascript
{
  "vision_score": {
    "average": 42.3,        // Average vision score per game
    "wards_placed": 12.5,   // Average wards placed per game
    "wards_cleared": 8.2,   // Average wards cleared per game
    "vision_per_minute": 1.8, // Vision score per minute
    "control_wards": 3.1,   // Control wards per game
    "trinket_efficiency": 0.85, // How well trinket is utilized
    "map_control": {
      "early_game": 0.35,   // Map control percentage early game
      "mid_game": 0.42,     // Map control percentage mid game
      "late_game": 0.38     // Map control percentage late game
    }
  }
}
```

#### 4. Damage Share Analysis
```javascript
{
  "damage_share": {
    "total_damage": 0.28,   // Share of team's total damage
    "champion_damage": 0.31, // Share of champion damage
    "objective_damage": 0.25, // Share of objective damage
    "damage_per_minute": 587.3, // Damage dealt per minute
    "damage_efficiency": 1.15, // Damage per gold ratio
    "damage_taken_share": 0.22, // Share of damage taken
    "effective_hp": 2847.5  // Effective health considering resistances
  }
}
```

#### 5. Gold Efficiency
```javascript
{
  "gold_efficiency": {
    "overall": 1.15,        // Overall gold efficiency ratio
    "gold_per_minute": 423.7, // Average gold per minute
    "cs_gold": 0.68,        // Percentage from CS
    "kill_gold": 0.18,      // Percentage from kills/assists
    "objective_gold": 0.14, // Percentage from objectives
    "item_efficiency": 0.92, // Item build efficiency
    "economy_ranking": 82   // Percentile in economic performance
  }
}
```

## 🔍 Advanced Analytics Features

### Champion-Specific Analytics

Get detailed performance data for specific champions:

```javascript
const championAnalytics = await client.gaming.getChampionAnalytics('NA', 'HeraldChampion', {
  champion: 'Jinx',
  timeRange: '30d',
  gameMode: 'RANKED_SOLO_5x5'
});

console.log(championAnalytics.data);
// {
//   champion_name: "Jinx",
//   games_played: 23,
//   win_rate: 0.652,
//   kda: { ratio: 3.8, kills: 8.1, deaths: 4.2, assists: 7.9 },
//   cs_per_minute: 8.4,
//   damage_share: 0.34,
//   item_builds: [
//     {
//       build: ["Kraken Slayer", "Phantom Dancer", "Infinity Edge"],
//       win_rate: 0.75,
//       games: 12
//     }
//   ],
//   skill_orders: [
//     {
//       order: "Q>E>W>Q>Q>R",
//       win_rate: 0.68,
//       games: 15
//     }
//   ]
// }
```

### Performance Trends

Track performance improvements over time:

```javascript
const trends = await client.gaming.getPerformanceTrends('NA', 'HeraldChampion', {
  metrics: ['kda', 'cs_per_minute', 'vision_score'],
  timeRange: '90d',
  granularity: 'daily'
});

// Visualize trends data
trends.data.forEach(day => {
  console.log(`${day.date}: KDA ${day.kda.ratio} (${day.kda.trend})`);
});
```

### Comparative Analytics

Compare performance against other players:

```javascript
const comparison = await client.gaming.comparePerformance('NA', 'HeraldChampion', {
  compareTo: 'rank_average', // or 'region_average', 'global_average'
  rank: 'GOLD',
  tier: 'II',
  metrics: ['kda', 'cs_per_minute', 'vision_score', 'damage_share']
});

console.log(comparison.data);
// {
//   player_performance: { kda: 3.25, cs_per_minute: 6.8, ... },
//   comparison_baseline: { kda: 2.15, cs_per_minute: 5.9, ... },
//   percentile_rankings: { kda: 85, cs_per_minute: 72, ... },
//   strengths: ["KDA", "Vision Control"],
//   weaknesses: ["CS Efficiency", "Objective Control"]
// }
```

## 🎯 Gaming Insights Engine

Herald.lol's AI-powered insights engine analyzes gaming data to provide actionable recommendations:

### Automated Insights

```javascript
const insights = await client.gaming.getGamingInsights('NA', 'HeraldChampion', {
  timeRange: '30d',
  analysisDepth: 'detailed' // basic, detailed, comprehensive
});

console.log(insights.gaming_insights);
// [
//   {
//     type: "improvement",
//     message: "Your CS/min has improved by 12% this week",
//     impact: "high",
//     confidence: 0.89,
//     action_items: [
//       "Continue focusing on last-hitting in early game",
//       "Practice farming while harassing opponent"
//     ],
//     supporting_data: {
//       current_cs_min: 6.8,
//       previous_cs_min: 6.1,
//       improvement_percentage: 0.12
//     }
//   },
//   {
//     type: "weakness",
//     message: "Vision score below average for your rank",
//     impact: "medium",
//     confidence: 0.76,
//     action_items: [
//       "Place more control wards in river bushes",
//       "Clear enemy wards when roaming"
//     ]
//   }
// ]
```

### Performance Coaching

Get AI-powered coaching recommendations:

```javascript
const coaching = await client.gaming.getCoachingRecommendations('NA', 'HeraldChampion', {
  focus_areas: ['laning', 'teamfight', 'macro'],
  skill_level: 'intermediate',
  main_role: 'adc'
});

console.log(coaching.recommendations);
// {
//   laning: {
//     priority: "high",
//     specific_advice: [
//       "Work on CS under tower - you're missing 15% of tower CS",
//       "Trade more efficiently - you're taking 23% more poke damage"
//     ],
//     drills: [
//       {
//         name: "CS Under Tower Practice",
//         description: "Practice last-hitting under tower for 10 minutes daily",
//         difficulty: "beginner"
//       }
//     ]
//   },
//   teamfight: {
//     priority: "medium",
//     specific_advice: [
//       "Positioning: Stay 550+ units from enemy frontline",
//       "Target selection: Focus highest threat within safe range"
//     ]
//   }
// }
```

## 📊 Real-time Match Analysis

### Live Match Tracking

Monitor ongoing matches in real-time:

```javascript
const liveMatch = await client.gaming.getCurrentMatch('NA', 'HeraldChampion');

if (liveMatch.active) {
  console.log(liveMatch.data);
  // {
  //   match_id: "NA1_4567890123",
  //   game_duration: 847, // seconds
  //   participants: [...],
  //   real_time_stats: {
  //     gold_advantage: 1250,
  //     experience_advantage: 890,
  //     objectives: {
  //       towers: { blue: 2, red: 1 },
  //       dragons: { blue: 1, red: 0 },
  //       kills: { blue: 8, red: 5 }
  //     }
  //   },
  //   predictions: {
  //     win_probability: { blue: 0.67, red: 0.33 },
  //     confidence: 0.74
  //   }
  // }
}
```

### Post-Match Analysis

Get detailed analysis immediately after match completion:

```javascript
const matchAnalysis = await client.gaming.analyzeMatch('NA', 'NA1_4567890123', {
  player_focus: 'HeraldChampion',
  analysis_type: 'comprehensive'
});

console.log(matchAnalysis.player_analysis);
// {
//   overall_performance: "above_average",
//   grade: "B+",
//   key_moments: [
//     {
//       timestamp: 1247,
//       event: "solo_kill",
//       impact: "high",
//       description: "Excellent trade timing led to solo kill"
//     }
//   ],
//   improvement_areas: [
//     "Ward placement frequency",
//     "Objective priority"
//   ],
//   match_impact: {
//     damage_contribution: 0.31,
//     vision_contribution: 0.28,
//     objective_contribution: 0.22
//   }
// }
```

## 🏆 Team Analytics

### Team Performance Metrics

Analyze team performance and synergy:

```javascript
const teamAnalytics = await client.teams.getTeamAnalytics('team_herald_champions', {
  timeRange: '30d',
  includeIndividualStats: true
});

console.log(teamAnalytics.data);
// {
//   team_stats: {
//     win_rate: 0.68,
//     average_game_duration: 1847,
//     kda: { kills: 15.2, deaths: 11.8, assists: 31.4 },
//     objective_control: {
//       first_dragon: 0.72,
//       first_baron: 0.58,
//       first_tower: 0.81
//     }
//   },
//   synergy_analysis: {
//     bot_lane_synergy: 0.84,    // ADC/Support synergy score
//     jungle_coordination: 0.79,  // Jungle/Lane coordination
//     team_fight_coordination: 0.72
//   },
//   draft_analysis: {
//     champion_pool_diversity: 0.67,
//     meta_adaptation: 0.73,
//     ban_efficiency: 0.58
//   }
// }
```

### Player Role Analysis

Understand how each player contributes to team success:

```javascript
const roleAnalysis = await client.teams.getPlayerRoleAnalysis('team_herald_champions', {
  player: 'HeraldChampion',
  role: 'adc'
});

console.log(roleAnalysis.role_performance);
// {
//   role_grade: "A-",
//   role_efficiency: 0.82,
//   team_synergy: {
//     with_support: 0.89,
//     with_jungle: 0.71,
//     with_team: 0.75
//   },
//   carry_potential: 0.78,
//   consistency: 0.84,
//   clutch_performance: 0.67
// }
```

## 📈 Data Export and Integration

### Bulk Analytics Export

Export comprehensive analytics data for offline analysis:

```javascript
const exportRequest = await client.gaming.exportSummonerAnalytics('NA', 'HeraldChampion', {
  format: 'json',
  timeRange: '90d',
  includeMatches: true,
  includeTimeline: true,
  mfaToken: 'your-mfa-token'
});

// Poll for export completion
const exportStatus = await client.exports.getStatus(exportRequest.export_id);
if (exportStatus.status === 'completed') {
  const downloadUrl = exportStatus.download_url;
  // Download and process the data
}
```

### Streaming Data Integration

Set up real-time data streams for continuous monitoring:

```javascript
const stream = client.gaming.createAnalyticsStream('NA', 'HeraldChampion', {
  metrics: ['kda', 'cs_per_minute', 'vision_score'],
  updateInterval: 'match_completion',
  includeInsights: true
});

stream.on('analytics_update', (data) => {
  console.log('New analytics data:', data);
  // Process real-time analytics updates
});

stream.on('insights_update', (insights) => {
  console.log('New gaming insights:', insights);
  // Handle AI-generated insights
});
```

## 🎮 Gaming Platform Best Practices

### Performance Optimization

1. **Caching Strategy**
```javascript
// Cache analytics data for better performance
const cache = new Map();
const cacheKey = `analytics_${region}_${summonerName}_${timeRange}`;

let analytics = cache.get(cacheKey);
if (!analytics) {
  analytics = await client.gaming.getSummonerAnalytics(region, summonerName, { timeRange });
  cache.set(cacheKey, analytics);
  // Set expiration for cache
  setTimeout(() => cache.delete(cacheKey), 300000); // 5 minutes
}
```

2. **Batch Requests**
```javascript
// Get analytics for multiple players efficiently
const players = ['Player1', 'Player2', 'Player3'];
const batchResults = await client.gaming.batchGetSummonerAnalytics('NA', players, {
  timeRange: '7d',
  metrics: ['kda', 'cs_per_minute']
});
```

3. **Error Handling**
```javascript
async function getGamingAnalytics(region, summonerName) {
  try {
    const analytics = await client.gaming.getSummonerAnalytics(region, summonerName);
    return analytics;
  } catch (error) {
    if (error.status === 404) {
      console.log('Summoner not found');
      return null;
    } else if (error.status === 429) {
      console.log('Rate limited, retrying in 60 seconds');
      await new Promise(resolve => setTimeout(resolve, 60000));
      return getGamingAnalytics(region, summonerName);
    } else {
      console.error('Analytics error:', error);
      throw error;
    }
  }
}
```

### Rate Limit Management

```javascript
class RateLimitManager {
  constructor() {
    this.requestQueue = [];
    this.processing = false;
  }

  async makeRequest(requestFn) {
    return new Promise((resolve, reject) => {
      this.requestQueue.push({ requestFn, resolve, reject });
      this.processQueue();
    });
  }

  async processQueue() {
    if (this.processing || this.requestQueue.length === 0) return;
    
    this.processing = true;
    
    while (this.requestQueue.length > 0) {
      const { requestFn, resolve, reject } = this.requestQueue.shift();
      
      try {
        const result = await requestFn();
        resolve(result);
        
        // Respect rate limits
        await new Promise(resolve => setTimeout(resolve, 1000)); // 1 second between requests
      } catch (error) {
        if (error.status === 429) {
          // Rate limited - put request back and wait
          this.requestQueue.unshift({ requestFn, resolve, reject });
          const retryAfter = error.headers['retry-after'] * 1000 || 60000;
          await new Promise(resolve => setTimeout(resolve, retryAfter));
        } else {
          reject(error);
        }
      }
    }
    
    this.processing = false;
  }
}
```

## 🔧 Advanced Use Cases

### Building Gaming Dashboards

Create comprehensive gaming dashboards with real-time updates:

```javascript
class GamingDashboard {
  constructor(region, summonerName) {
    this.region = region;
    this.summonerName = summonerName;
    this.client = new HeraldClient({ apiKey: process.env.HERALD_API_KEY });
  }

  async initialize() {
    // Get comprehensive player data
    const [analytics, recentMatches, trends] = await Promise.all([
      this.client.gaming.getSummonerAnalytics(this.region, this.summonerName),
      this.client.gaming.getRecentMatches(this.region, this.summonerName, { limit: 10 }),
      this.client.gaming.getPerformanceTrends(this.region, this.summonerName, { timeRange: '30d' })
    ]);

    return {
      analytics: analytics.data,
      recentMatches: recentMatches.matches,
      trends: trends.data
    };
  }

  async getRealtimeUpdates() {
    // Check for live match
    const liveMatch = await this.client.gaming.getCurrentMatch(this.region, this.summonerName);
    
    if (liveMatch.active) {
      return {
        type: 'live_match',
        data: liveMatch.data
      };
    }

    // Check for new completed matches
    const latestMatch = await this.client.gaming.getLatestMatch(this.region, this.summonerName);
    return {
      type: 'match_completed',
      data: latestMatch
    };
  }
}
```

### Automated Performance Tracking

Set up automated performance tracking and alerts:

```javascript
class PerformanceTracker {
  constructor(players) {
    this.players = players;
    this.client = new HeraldClient({ apiKey: process.env.HERALD_API_KEY });
    this.thresholds = {
      kda_improvement: 0.1,  // 10% improvement
      rank_change: true,     // Any rank change
      win_streak: 5          // 5+ win streak
    };
  }

  async trackPerformance() {
    for (const player of this.players) {
      const { region, summonerName } = player;
      
      try {
        // Get current analytics
        const current = await this.client.gaming.getSummonerAnalytics(region, summonerName, {
          timeRange: '7d'
        });

        // Compare with previous week
        const previous = await this.client.gaming.getSummonerAnalytics(region, summonerName, {
          timeRange: '7d',
          offset: '7d'
        });

        // Check for significant improvements
        const kdaImprovement = (current.data.kda.ratio - previous.data.kda.ratio) / previous.data.kda.ratio;
        
        if (kdaImprovement > this.thresholds.kda_improvement) {
          await this.sendAlert({
            player: summonerName,
            type: 'improvement',
            message: `KDA improved by ${(kdaImprovement * 100).toFixed(1)}%`,
            data: { current: current.data.kda.ratio, previous: previous.data.kda.ratio }
          });
        }

        // Check win streak
        const recentMatches = await this.client.gaming.getRecentMatches(region, summonerName, { limit: 10 });
        const winStreak = this.calculateWinStreak(recentMatches.matches);
        
        if (winStreak >= this.thresholds.win_streak) {
          await this.sendAlert({
            player: summonerName,
            type: 'win_streak',
            message: `On a ${winStreak} game win streak!`,
            data: { streak: winStreak }
          });
        }

      } catch (error) {
        console.error(`Error tracking ${summonerName}:`, error);
      }
    }
  }

  calculateWinStreak(matches) {
    let streak = 0;
    for (const match of matches) {
      if (match.result === 'win') {
        streak++;
      } else {
        break;
      }
    }
    return streak;
  }

  async sendAlert(alert) {
    // Send to Discord, Slack, email, etc.
    console.log('Performance Alert:', alert);
  }
}
```

## 📚 SDK Integration Examples

### React Dashboard Component

```jsx
import React, { useState, useEffect } from 'react';
import { HeraldClient } from 'herald-lol-sdk';

const GamingAnalyticsDashboard = ({ region, summonerName }) => {
  const [analytics, setAnalytics] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const client = new HeraldClient({
    apiKey: process.env.REACT_APP_HERALD_API_KEY
  });

  useEffect(() => {
    loadAnalytics();
  }, [region, summonerName]);

  const loadAnalytics = async () => {
    try {
      setLoading(true);
      const response = await client.gaming.getSummonerAnalytics(region, summonerName, {
        timeRange: '30d'
      });
      setAnalytics(response.data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="loading">Loading gaming analytics...</div>;
  if (error) return <div className="error">Error: {error}</div>;
  if (!analytics) return <div>No data available</div>;

  return (
    <div className="gaming-dashboard">
      <h2>Gaming Analytics for {summonerName}</h2>
      
      <div className="metrics-grid">
        <div className="metric-card">
          <h3>KDA</h3>
          <div className="kda-display">
            {analytics.kda.kills.toFixed(1)} / {analytics.kda.deaths.toFixed(1)} / {analytics.kda.assists.toFixed(1)}
          </div>
          <div className="kda-ratio">Ratio: {analytics.kda.ratio.toFixed(2)}</div>
        </div>

        <div className="metric-card">
          <h3>CS/min</h3>
          <div className="cs-display">{analytics.cs_per_minute.toFixed(1)}</div>
          <div className="percentile">
            {analytics.cs_percentile}th percentile for your rank
          </div>
        </div>

        <div className="metric-card">
          <h3>Vision Score</h3>
          <div className="vision-display">{analytics.vision_score.toFixed(1)}</div>
          <div className="vision-details">
            {analytics.vision_wards_placed.toFixed(1)} wards placed per game
          </div>
        </div>

        <div className="metric-card">
          <h3>Damage Share</h3>
          <div className="damage-display">
            {(analytics.damage_share * 100).toFixed(1)}%
          </div>
          <div className="damage-details">
            {analytics.damage_per_minute.toFixed(0)} DPM
          </div>
        </div>
      </div>

      <div className="insights-section">
        <h3>Gaming Insights</h3>
        {analytics.gaming_insights.map((insight, index) => (
          <div key={index} className={`insight insight-${insight.type}`}>
            <div className="insight-message">{insight.message}</div>
            <div className="insight-impact">Impact: {insight.impact}</div>
            {insight.action_items && (
              <ul className="action-items">
                {insight.action_items.map((item, i) => (
                  <li key={i}>{item}</li>
                ))}
              </ul>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default GamingAnalyticsDashboard;
```

This comprehensive guide provides everything developers need to integrate Herald.lol's gaming analytics into their applications, from basic API calls to advanced real-time analytics and AI-powered insights.