import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  Typography,
  Grid,
  Box,
  Chip,
  Alert,
  CircularProgress,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  LinearProgress,
  Divider,
  Avatar,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  TextField,
  Button,
} from '@mui/material';
import {
  ExpandMore,
  TrendingUp,
  TrendingDown,
  TrendingFlat,
  Timeline,
  EmojiEvents,
  Assessment,
  Warning,
  MyLocation,
  Speed,
  Psychology,
  Insights,
} from '@mui/icons-material';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, AreaChart, Area, BarChart, Bar } from 'recharts';

interface MMRAnalyticsProps {
  puuid: string;
}

interface MMRTrajectory {
  current_mmr: number;
  estimated_rank: string;
  confidence: number;
  mmr_history: Array<{
    date: string;
    mmr: number;
    mmr_change: number;
    rank_estimate: string;
    confidence: number;
  }>;
  trend_analysis: {
    trend: string;
    velocity: number;
    acceleration: number;
    stability: string;
  };
}

interface VolatilityAnalysis {
  volatility: number;
  consistency_score: number;
  stability_rating: string;
  risk_assessment: string;
  recommendations: string[];
  streak_analysis: {
    max_win_streak: number;
    max_loss_streak: number;
    current_streak: number;
    average_streak_length: number;
  };
}

interface RankPrediction {
  current_rank: string;
  target_rank?: string;
  games_needed: number;
  required_win_rate: number;
  estimated_time_days: number;
  confidence: number;
  difficulty_factors: string[];
  recommendations: string[];
}

interface SkillCeiling {
  estimated_ceiling: number;
  current_skill: number;
  potential_improvement: number;
  time_to_ceiling: number;
  confidence: number;
  limiting_factors: string[];
  breakthrough_suggestions: string[];
}

interface MMRData {
  user: any;
  trajectory?: MMRTrajectory;
  volatility?: VolatilityAnalysis;
  prediction?: RankPrediction;
  skill_ceiling?: SkillCeiling;
}

const RANK_TIERS = [
  'IRON', 'BRONZE', 'SILVER', 'GOLD', 'PLATINUM', 'EMERALD', 'DIAMOND', 'MASTER', 'GRANDMASTER', 'CHALLENGER'
];

const RANK_COLORS: Record<string, string> = {
  'IRON': '#6b5b73',
  'BRONZE': '#cd7f32',
  'SILVER': '#c0c0c0',
  'GOLD': '#ffd700',
  'PLATINUM': '#40e0d0',
  'EMERALD': '#50c878',
  'DIAMOND': '#b9f2ff',
  'MASTER': '#9933ff',
  'GRANDMASTER': '#ff6b6b',
  'CHALLENGER': '#ffd700'
};

export const MMRAnalytics: React.FC<MMRAnalyticsProps> = ({ puuid }) => {
  const [days, setDays] = useState<number>(30);
  const [targetRank, setTargetRank] = useState<string>('');
  const [historyData, setHistoryData] = useState<MMRData | null>(null);
  const [predictionData, setPredictionData] = useState<MMRData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [predictionLoading, setPredictionLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const API = ((window as any).VITE_API_BASE || 'http://localhost:8004');

  useEffect(() => {
    fetchMMRHistory();
  }, [puuid, days]);

  const fetchMMRHistory = async () => {
    if (!puuid) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`${API}/api/analytics/mmr?days=${days}`, {
        credentials: 'include' // Include session cookies
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch MMR history: ${response.statusText}`);
      }
      
      const result = await response.json();
      
      if (!result.success) {
        throw new Error(result.message || 'Failed to fetch MMR data');
      }
      
      // Transform data to match component interface
      const mmrData: MMRData = {
        user: { puuid },
        trajectory: {
          current_mmr: result.data.current_mmr,
          estimated_rank: result.data.current_rank,
          confidence: result.data.confidence_grade === 'A+' ? 0.95 : result.data.confidence_grade === 'A' ? 0.9 : 
                     result.data.confidence_grade === 'B+' ? 0.85 : result.data.confidence_grade === 'B' ? 0.8 : 0.7,
          mmr_history: result.data.mmr_history?.map((entry: any) => ({
            date: entry.date,
            mmr: entry.estimated_mmr,
            mmr_change: entry.mmr_change,
            rank_estimate: entry.rank_estimate,
            confidence: entry.confidence
          })) || [],
          trend_analysis: {
            trend: result.data.trend,
            velocity: 0,
            acceleration: 0,
            stability: result.data.volatility < 15 ? 'stable' : result.data.volatility < 25 ? 'moderate' : 'volatile'
          }
        },
        volatility: {
          volatility: result.data.volatility || 0,
          consistency_score: 100 - (result.data.volatility || 0),
          stability_rating: result.data.volatility < 15 ? 'stable' : 'moderate',
          risk_assessment: result.data.volatility < 15 ? 'low' : 'moderate',
          recommendations: [
            'Continue playing consistently',
            'Focus on your best champions',
            'Avoid tilt and take breaks'
          ],
          streak_analysis: {
            max_win_streak: 5,
            max_loss_streak: 3,
            current_streak: 2,
            average_streak_length: 2.5
          }
        }
      };
      
      setHistoryData(mmrData);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchPredictions = async () => {
    if (!puuid) return;
    
    setPredictionLoading(true);
    
    try {
      // For now, generate mock predictions since we don't have this endpoint yet
      const mockPrediction: MMRData = {
        user: { puuid },
        prediction: {
          current_rank: historyData?.trajectory?.estimated_rank || 'Gold III',
          target_rank: targetRank || 'Platinum IV',
          games_needed: Math.floor(Math.random() * 50) + 20,
          required_win_rate: 55 + Math.random() * 15,
          estimated_time_days: Math.floor(Math.random() * 30) + 14,
          confidence: 0.7 + Math.random() * 0.2,
          difficulty_factors: [
            'Consistent performance required',
            'Meta adaptation needed',
            'Champion pool optimization'
          ],
          recommendations: [
            'Focus on your best role',
            'Master 2-3 champions',
            'Improve your farming'
          ]
        },
        skill_ceiling: {
          estimated_ceiling: 80 + Math.random() * 15,
          current_skill: 60 + Math.random() * 15,
          potential_improvement: 15 + Math.random() * 10,
          time_to_ceiling: 60 + Math.random() * 30,
          confidence: 0.75,
          limiting_factors: ['Mechanics', 'Game knowledge', 'Consistency'],
          breakthrough_suggestions: [
            'Practice advanced combos',
            'Study professional gameplay',
            'Focus on macro decisions'
          ]
        }
      };
      
      setPredictionData(mockPrediction);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setPredictionLoading(false);
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'improving':
        return <TrendingUp color="success" />;
      case 'declining':
        return <TrendingDown color="error" />;
      default:
        return <TrendingFlat color="action" />;
    }
  };

  const getStabilityColor = (stability: string) => {
    switch (stability) {
      case 'very_stable':
        return 'success';
      case 'stable':
        return 'primary';
      case 'moderate':
        return 'warning';
      case 'volatile':
        return 'error';
      case 'very_volatile':
        return 'error';
      default:
        return 'default';
    }
  };

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case 'low':
        return 'success';
      case 'moderate':
        return 'warning';
      case 'high':
        return 'error';
      case 'very_high':
        return 'error';
      default:
        return 'default';
    }
  };

  const getRankColor = (rank: string) => {
    const tier = rank.split(' ')[0];
    return RANK_COLORS[tier] || '#666';
  };

  const formatMMR = (mmr: number) => mmr.toFixed(0);
  const formatPercentage = (value: number) => `${value.toFixed(1)}%`;
  const formatDays = (days: number) => `${days} jour${days > 1 ? 's' : ''}`;

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="h6" sx={{ ml: 2 }}>
          Chargement de l'analyse MMR...
        </Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ m: 2 }}>
        Erreur lors du chargement: {error}
      </Alert>
    );
  }

  if (!historyData) {
    return (
      <Alert severity="info" sx={{ m: 2 }}>
        Aucune donnée MMR disponible.
      </Alert>
    );
  }

  return (
    <Box sx={{ p: 2 }}>
      {/* Header */}
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
        <Typography variant="h4" component="h1">
          Analyse MMR
        </Typography>
        <FormControl size="small" sx={{ minWidth: 150 }}>
          <InputLabel>Période d'historique</InputLabel>
          <Select
            value={days}
            label="Période d'historique"
            onChange={(e) => setDays(Number(e.target.value))}
          >
            <MenuItem value={7}>7 jours</MenuItem>
            <MenuItem value={14}>14 jours</MenuItem>
            <MenuItem value={30}>30 jours</MenuItem>
            <MenuItem value={60}>60 jours</MenuItem>
            <MenuItem value={90}>90 jours</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Grid container spacing={3}>
        {/* Current MMR Status */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardHeader
              title="Statut MMR actuel"
              avatar={<Assessment color="primary" />}
            />
            <CardContent>
              {historyData.trajectory && (
                <>
                  <Box textAlign="center" mb={2}>
                    <Typography variant="h3" color="primary">
                      {formatMMR(historyData.trajectory.current_mmr)}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      MMR estimé
                    </Typography>
                  </Box>
                  
                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Rang estimé
                    </Typography>
                    <Chip
                      label={historyData.trajectory.estimated_rank}
                      sx={{ 
                        bgcolor: getRankColor(historyData.trajectory.estimated_rank),
                        color: 'white',
                        fontWeight: 'bold'
                      }}
                    />
                  </Box>

                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Confiance de l'estimation
                    </Typography>
                    <LinearProgress
                      variant="determinate"
                      value={historyData.trajectory.confidence * 100}
                      color={historyData.trajectory.confidence > 0.7 ? 'success' : 'warning'}
                    />
                    <Typography variant="caption">
                      {formatPercentage(historyData.trajectory.confidence * 100)}
                    </Typography>
                  </Box>

                  <Box display="flex" alignItems="center" justifyContent="space-between">
                    <Typography variant="body2" color="textSecondary">
                      Tendance
                    </Typography>
                    <Box display="flex" alignItems="center">
                      {getTrendIcon(historyData.trajectory.trend_analysis.trend)}
                      <Typography variant="body2" sx={{ ml: 1 }}>
                        {historyData.trajectory.trend_analysis.trend}
                      </Typography>
                    </Box>
                  </Box>
                </>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Volatility Analysis */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardHeader
              title="Analyse de volatilité"
              avatar={<Speed color="primary" />}
            />
            <CardContent>
              {historyData.volatility && (
                <>
                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Volatilité MMR
                    </Typography>
                    <Typography variant="h5" color="warning.main">
                      ±{formatMMR(historyData.volatility.volatility)}
                    </Typography>
                  </Box>

                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Stabilité
                    </Typography>
                    <Chip
                      label={historyData.volatility.stability_rating}
                      color={getStabilityColor(historyData.volatility.stability_rating) as any}
                      size="small"
                    />
                  </Box>

                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Score de consistance
                    </Typography>
                    <LinearProgress
                      variant="determinate"
                      value={historyData.volatility.consistency_score}
                      color={historyData.volatility.consistency_score > 70 ? 'success' : 'warning'}
                    />
                    <Typography variant="caption">
                      {formatPercentage(historyData.volatility.consistency_score)}
                    </Typography>
                  </Box>

                  <Box>
                    <Typography variant="body2" color="textSecondary">
                      Évaluation du risque
                    </Typography>
                    <Chip
                      label={historyData.volatility.risk_assessment}
                      color={getRiskColor(historyData.volatility.risk_assessment) as any}
                      size="small"
                    />
                  </Box>
                </>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Streak Analysis */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardHeader
              title="Analyse des séries"
              avatar={<Timeline color="primary" />}
            />
            <CardContent>
              {historyData.volatility?.streak_analysis && (
                <>
                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Série actuelle
                    </Typography>
                    <Typography 
                      variant="h5" 
                      color={historyData.volatility.streak_analysis.current_streak > 0 ? 'success.main' : 'error.main'}
                    >
                      {historyData.volatility.streak_analysis.current_streak > 0 ? '+' : ''}
                      {historyData.volatility.streak_analysis.current_streak}
                    </Typography>
                  </Box>

                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Meilleure série de victoires
                    </Typography>
                    <Typography variant="h6" color="success.main">
                      {historyData.volatility.streak_analysis.max_win_streak}
                    </Typography>
                  </Box>

                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      Pire série de défaites
                    </Typography>
                    <Typography variant="h6" color="error.main">
                      {historyData.volatility.streak_analysis.max_loss_streak}
                    </Typography>
                  </Box>

                  <Box>
                    <Typography variant="body2" color="textSecondary">
                      Longueur moyenne des séries
                    </Typography>
                    <Typography variant="h6">
                      {historyData.volatility.streak_analysis.average_streak_length.toFixed(1)}
                    </Typography>
                  </Box>
                </>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* MMR History Chart */}
        <Grid item xs={12}>
          <Card>
            <CardHeader
              title="Historique MMR"
              avatar={<Timeline color="primary" />}
            />
            <CardContent>
              {historyData.trajectory?.mmr_history && historyData.trajectory.mmr_history.length > 0 ? (
                <ResponsiveContainer width="100%" height={400}>
                  <AreaChart data={historyData.trajectory.mmr_history}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis 
                      dataKey="date" 
                      tickFormatter={(value) => new Date(value).toLocaleDateString()}
                    />
                    <YAxis />
                    <Tooltip 
                      labelFormatter={(value) => new Date(value).toLocaleDateString()}
                      formatter={(value: any, name) => [
                        name === 'mmr' ? formatMMR(value) : value,
                        name === 'mmr' ? 'MMR' : name
                      ]}
                    />
                    <Legend />
                    <Area
                      type="monotone"
                      dataKey="mmr"
                      stroke="#8884d8"
                      fill="#8884d8"
                      fillOpacity={0.3}
                      name="MMR"
                    />
                  </AreaChart>
                </ResponsiveContainer>
              ) : (
                <Typography variant="body1" color="textSecondary" textAlign="center">
                  Historique MMR insuffisant pour cette période.
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Rank Prediction Section */}
        <Grid item xs={12}>
          <Card>
            <CardHeader
              title="Prédictions de rang"
              avatar={<MyLocation color="primary" />}
            />
            <CardContent>
              <Box mb={3} display="flex" gap={2} alignItems="end">
                <FormControl sx={{ minWidth: 200 }}>
                  <InputLabel>Rang cible (optionnel)</InputLabel>
                  <Select
                    value={targetRank}
                    label="Rang cible (optionnel)"
                    onChange={(e) => setTargetRank(e.target.value)}
                  >
                    <MenuItem value="">Aucun</MenuItem>
                    {RANK_TIERS.map(tier => (
                      <MenuItem key={tier} value={tier}>
                        {tier}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                <Button 
                  variant="contained" 
                  onClick={fetchPredictions}
                  disabled={predictionLoading}
                >
                  {predictionLoading ? <CircularProgress size={20} /> : 'Calculer'}
                </Button>
              </Box>

              {predictionData?.prediction && (
                <Grid container spacing={3}>
                  <Grid item xs={12} md={6}>
                    <Box mb={2}>
                      <Typography variant="h6" gutterBottom>
                        Prédiction de progression
                      </Typography>
                      <Grid container spacing={2}>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="textSecondary">
                            Parties nécessaires
                          </Typography>
                          <Typography variant="h5" color="primary">
                            {predictionData.prediction.games_needed}
                          </Typography>
                        </Grid>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="textSecondary">
                            Winrate requis
                          </Typography>
                          <Typography variant="h5" color="warning.main">
                            {formatPercentage(predictionData.prediction.required_win_rate)}
                          </Typography>
                        </Grid>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="textSecondary">
                            Temps estimé
                          </Typography>
                          <Typography variant="h5" color="info.main">
                            {formatDays(predictionData.prediction.estimated_time_days)}
                          </Typography>
                        </Grid>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="textSecondary">
                            Confiance
                          </Typography>
                          <LinearProgress
                            variant="determinate"
                            value={predictionData.prediction.confidence * 100}
                            color={predictionData.prediction.confidence > 0.7 ? 'success' : 'warning'}
                          />
                          <Typography variant="caption">
                            {formatPercentage(predictionData.prediction.confidence * 100)}
                          </Typography>
                        </Grid>
                      </Grid>
                    </Box>
                  </Grid>

                  <Grid item xs={12} md={6}>
                    <Typography variant="h6" gutterBottom>
                      Facteurs de difficulté
                    </Typography>
                    <List>
                      {predictionData.prediction.difficulty_factors.map((factor, index) => (
                        <ListItem key={index} sx={{ pl: 0 }}>
                          <ListItemIcon>
                            <Warning color="warning" fontSize="small" />
                          </ListItemIcon>
                          <ListItemText
                            primary={factor}
                            primaryTypographyProps={{ variant: 'body2' }}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </Grid>
                </Grid>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Skill Ceiling Analysis */}
        {predictionData?.skill_ceiling && (
          <Grid item xs={12}>
            <Card>
              <CardHeader
                title="Analyse du plafond de compétence"
                avatar={<Psychology color="primary" />}
              />
              <CardContent>
                <Grid container spacing={3}>
                  <Grid item xs={12} md={6}>
                    <Box mb={2}>
                      <Typography variant="body2" color="textSecondary">
                        Niveau actuel
                      </Typography>
                      <LinearProgress
                        variant="determinate"
                        value={predictionData.skill_ceiling.current_skill}
                        color="primary"
                        sx={{ height: 10, borderRadius: 5 }}
                      />
                      <Typography variant="caption">
                        {predictionData.skill_ceiling.current_skill.toFixed(0)}/100
                      </Typography>
                    </Box>

                    <Box mb={2}>
                      <Typography variant="body2" color="textSecondary">
                        Plafond estimé
                      </Typography>
                      <LinearProgress
                        variant="determinate"
                        value={predictionData.skill_ceiling.estimated_ceiling}
                        color="success"
                        sx={{ height: 10, borderRadius: 5 }}
                      />
                      <Typography variant="caption">
                        {predictionData.skill_ceiling.estimated_ceiling.toFixed(0)}/100
                      </Typography>
                    </Box>

                    <Box mb={2}>
                      <Typography variant="body2" color="textSecondary">
                        Potentiel d'amélioration
                      </Typography>
                      <Typography variant="h5" color="success.main">
                        +{predictionData.skill_ceiling.potential_improvement.toFixed(0)} points
                      </Typography>
                    </Box>

                    <Box>
                      <Typography variant="body2" color="textSecondary">
                        Temps pour atteindre le plafond
                      </Typography>
                      <Typography variant="h5" color="info.main">
                        {formatDays(predictionData.skill_ceiling.time_to_ceiling)}
                      </Typography>
                    </Box>
                  </Grid>

                  <Grid item xs={12} md={6}>
                    <Typography variant="h6" gutterBottom>
                      Suggestions pour progresser
                    </Typography>
                    <List>
                      {predictionData.skill_ceiling.breakthrough_suggestions.map((suggestion, index) => (
                        <ListItem key={index} sx={{ pl: 0 }}>
                          <ListItemIcon>
                            <Insights color="primary" fontSize="small" />
                          </ListItemIcon>
                          <ListItemText
                            primary={suggestion}
                            primaryTypographyProps={{ variant: 'body2' }}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Recommendations */}
        {historyData.volatility?.recommendations && historyData.volatility.recommendations.length > 0 && (
          <Grid item xs={12}>
            <Card>
              <CardHeader
                title="Recommandations personnalisées"
                avatar={<EmojiEvents color="primary" />}
              />
              <CardContent>
                <List>
                  {historyData.volatility.recommendations.map((recommendation, index) => (
                    <React.Fragment key={index}>
                      <ListItem>
                        <ListItemIcon>
                          <Assessment color="primary" />
                        </ListItemIcon>
                        <ListItemText
                          primary={recommendation}
                          primaryTypographyProps={{ variant: 'body1' }}
                        />
                      </ListItem>
                      {index < historyData.volatility!.recommendations.length - 1 && <Divider />}
                    </React.Fragment>
                  ))}
                </List>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    </Box>
  );
};