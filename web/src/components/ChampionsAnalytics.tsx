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
  Tabs,
  Tab,
} from '@mui/material';
import {
  ExpandMore,
  EmojiEvents,
  TrendingUp,
  TrendingDown,
  TrendingFlat,
  Star,
  Warning,
  Assessment,
  Sports,
  Timeline,
} from '@mui/icons-material';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, LineChart, Line, ResponsiveContainer, RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar } from 'recharts';

interface ChampionsAnalyticsProps {
  puuid: string;
}

interface ChampionMastery {
  champion_id: number;
  champion_name: string;
  games_played: number;
  performance_metrics: {
    games_played: number;
    wins: number;
    losses: number;
    win_rate: number;
    avg_kills: number;
    avg_deaths: number;
    avg_assists: number;
    avg_kda: number;
    avg_cs_per_min: number;
    avg_gold_per_min: number;
    avg_damage_per_min: number;
    avg_vision_score: number;
    performance_score: number;
    trend_direction: string;
  };
  mastery_score: number;
  best_game: any;
  worst_game: any;
  improvement_suggestions: string[];
  skill_progression: any;
}

interface RolePerformance {
  games_played: number;
  wins: number;
  losses: number;
  win_rate: number;
  avg_kills: number;
  avg_deaths: number;
  avg_assists: number;
  avg_kda: number;
  avg_cs_per_min: number;
  avg_gold_per_min: number;
  avg_damage_per_min: number;
  avg_vision_score: number;
  performance_score: number;
  trend_direction: string;
}

interface ChampionsData {
  user: any;
  role?: string;
  period: string;
  champions?: ChampionMastery[];
  top_champions?: Array<{
    champion_name: string;
    games: number;
    win_rate: number;
    performance_score: number;
    avg_kda: number;
  }>;
  role_performance?: Record<string, RolePerformance>;
}

const ROLE_NAMES: Record<string, string> = {
  'TOP': 'Top',
  'JUNGLE': 'Jungle', 
  'MIDDLE': 'Mid',
  'BOTTOM': 'ADC',
  'UTILITY': 'Support',
  'UNKNOWN': 'Autre'
};

const ROLES = ['TOP', 'JUNGLE', 'MIDDLE', 'BOTTOM', 'UTILITY'];

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`champions-tabpanel-${index}`}
      aria-labelledby={`champions-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

export const ChampionsAnalytics: React.FC<ChampionsAnalyticsProps> = ({ puuid }) => {
  const [period, setPeriod] = useState<string>('month');
  const [selectedRole, setSelectedRole] = useState<string>('all');
  const [tabValue, setTabValue] = useState(0);
  const [data, setData] = useState<ChampionsData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const API = ((window as any).VITE_API_BASE || 'http://localhost:8004');

  useEffect(() => {
    fetchChampionsData();
  }, [puuid, period, selectedRole]);

  const fetchChampionsData = async () => {
    if (!puuid) return;
    
    setLoading(true);
    setError(null);
    
    try {
      let url;
      if (selectedRole !== 'all') {
        // Use role-specific analytics endpoint
        url = `${API}/api/analytics/champions/${selectedRole}?period=${period}`;
      } else {
        // Use general period analytics endpoint
        url = `${API}/api/analytics/period/${period}`;
      }
      
      const response = await fetch(url, {
        credentials: 'include' // Include session cookies
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch champions data: ${response.statusText}`);
      }
      
      const result = await response.json();
      
      if (!result.success) {
        throw new Error(result.message || 'Failed to fetch data');
      }
      
      // Transform data to match component interface
      const championsData: ChampionsData = {
        user: { puuid },
        role: selectedRole !== 'all' ? selectedRole : undefined,
        period: period,
        top_champions: result.data.top_champions || result.data.champions,
        role_performance: result.data.role_performance,
        champions: result.data.champions // For detailed analysis tab
      };
      
      setData(championsData);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'improving':
        return <TrendingUp color="success" fontSize="small" />;
      case 'declining':
        return <TrendingDown color="error" fontSize="small" />;
      default:
        return <TrendingFlat color="action" fontSize="small" />;
    }
  };

  const getMasteryColor = (score: number) => {
    if (score >= 80) return 'success';
    if (score >= 60) return 'primary';
    if (score >= 40) return 'warning';
    return 'error';
  };

  const formatWinRate = (winRate: number) => `${winRate.toFixed(1)}%`;
  const formatKDA = (kda: number) => kda.toFixed(2);
  const formatStat = (stat: number) => stat.toFixed(1);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="h6" sx={{ ml: 2 }}>
          Chargement des analyses de champions...
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

  if (!data) {
    return (
      <Alert severity="info" sx={{ m: 2 }}>
        Aucune donnée de champions disponible.
      </Alert>
    );
  }

  const createRadarData = (metrics: any) => {
    return [
      { stat: 'KDA', value: Math.min(metrics.avg_kda * 25, 100) },
      { stat: 'CS/min', value: Math.min(metrics.avg_cs_per_min * 10, 100) },
      { stat: 'Damage', value: Math.min(metrics.avg_damage_per_min / 20, 100) },
      { stat: 'Vision', value: Math.min(metrics.avg_vision_score * 5, 100) },
      { stat: 'Gold', value: Math.min(metrics.avg_gold_per_min / 5, 100) },
    ];
  };

  return (
    <Box sx={{ p: 2 }}>
      {/* Header */}
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
        <Typography variant="h4" component="h1">
          Analyse des Champions
        </Typography>
        <Box display="flex" gap={2}>
          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel>Rôle</InputLabel>
            <Select
              value={selectedRole}
              label="Rôle"
              onChange={(e) => setSelectedRole(e.target.value)}
            >
              <MenuItem value="all">Tous les rôles</MenuItem>
              {ROLES.map(role => (
                <MenuItem key={role} value={role}>
                  {ROLE_NAMES[role]}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel>Période</InputLabel>
            <Select
              value={period}
              label="Période"
              onChange={(e) => setPeriod(e.target.value)}
            >
              <MenuItem value="week">Cette semaine</MenuItem>
              <MenuItem value="month">Ce mois</MenuItem>
              <MenuItem value="season">Cette saison</MenuItem>
            </Select>
          </FormControl>
        </Box>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={tabValue} onChange={(e, newValue) => setTabValue(newValue)}>
          <Tab label="Vue d'ensemble" />
          <Tab label="Analyse détaillée" />
          <Tab label="Performance par rôle" />
        </Tabs>
      </Box>

      {/* Tab 1: Overview */}
      <TabPanel value={tabValue} index={0}>
        <Grid container spacing={3}>
          {/* Top Champions Summary */}
          <Grid item xs={12}>
            <Card>
              <CardHeader
                title="Top Champions"
                subheader={`${selectedRole !== 'all' ? ROLE_NAMES[selectedRole] : 'Tous rôles'} • ${period}`}
                avatar={<EmojiEvents color="primary" />}
              />
              <CardContent>
                {data.top_champions && data.top_champions.length > 0 ? (
                  <Grid container spacing={2}>
                    {data.top_champions.slice(0, 6).map((champion, index) => (
                      <Grid item xs={12} sm={6} md={4} key={champion.champion_name}>
                        <Card variant="outlined">
                          <CardContent>
                            <Box display="flex" alignItems="center" mb={1}>
                              <Avatar sx={{ bgcolor: 'primary.main', mr: 2, width: 32, height: 32 }}>
                                {index + 1}
                              </Avatar>
                              <Typography variant="h6" noWrap>
                                {champion.champion_name}
                              </Typography>
                            </Box>
                            <Box mb={1}>
                              <Typography variant="body2" color="textSecondary">
                                {champion.games} parties
                              </Typography>
                            </Box>
                            <Box display="flex" justifyContent="space-between" alignItems="center">
                              <Chip
                                label={formatWinRate(champion.win_rate)}
                                color={champion.win_rate >= 60 ? 'success' : champion.win_rate >= 50 ? 'primary' : 'error'}
                                size="small"
                              />
                              <Typography variant="body2">
                                KDA: {formatKDA(champion.avg_kda)}
                              </Typography>
                            </Box>
                            <Box mt={1}>
                              <LinearProgress
                                variant="determinate"
                                value={champion.performance_score}
                                color={champion.performance_score >= 70 ? 'success' : champion.performance_score >= 50 ? 'primary' : 'warning'}
                              />
                              <Typography variant="caption" color="textSecondary">
                                Score: {champion.performance_score.toFixed(0)}/100
                              </Typography>
                            </Box>
                          </CardContent>
                        </Card>
                      </Grid>
                    ))}
                  </Grid>
                ) : (
                  <Typography variant="body1" color="textSecondary" textAlign="center">
                    Aucun champion trouvé pour cette période et ce rôle.
                  </Typography>
                )}
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>

      {/* Tab 2: Detailed Analysis */}
      <TabPanel value={tabValue} index={1}>
        <Grid container spacing={3}>
          {data.champions && data.champions.map((champion, index) => (
            <Grid item xs={12} key={champion.champion_id}>
              <Accordion>
                <AccordionSummary expandIcon={<ExpandMore />}>
                  <Box display="flex" alignItems="center" width="100%">
                    <Avatar sx={{ bgcolor: getMasteryColor(champion.mastery_score) + '.main', mr: 2 }}>
                      <Star />
                    </Avatar>
                    <Box flexGrow={1}>
                      <Typography variant="h6">{champion.champion_name}</Typography>
                      <Typography variant="body2" color="textSecondary">
                        {champion.games_played} parties • Maîtrise: {champion.mastery_score.toFixed(0)}/100
                      </Typography>
                    </Box>
                    <Box display="flex" alignItems="center" gap={1}>
                      <Chip
                        label={formatWinRate(champion.performance_metrics.win_rate)}
                        color={champion.performance_metrics.win_rate >= 60 ? 'success' : 'primary'}
                        size="small"
                      />
                      {getTrendIcon(champion.performance_metrics.trend_direction)}
                    </Box>
                  </Box>
                </AccordionSummary>
                <AccordionDetails>
                  <Grid container spacing={3}>
                    {/* Stats Overview */}
                    <Grid item xs={12} md={8}>
                      <Typography variant="h6" gutterBottom>
                        Statistiques détaillées
                      </Typography>
                      <Grid container spacing={2}>
                        <Grid item xs={6} sm={3}>
                          <Box textAlign="center">
                            <Typography variant="h5" color="primary">
                              {formatKDA(champion.performance_metrics.avg_kda)}
                            </Typography>
                            <Typography variant="body2" color="textSecondary">
                              KDA
                            </Typography>
                          </Box>
                        </Grid>
                        <Grid item xs={6} sm={3}>
                          <Box textAlign="center">
                            <Typography variant="h5" color="secondary">
                              {formatStat(champion.performance_metrics.avg_cs_per_min)}
                            </Typography>
                            <Typography variant="body2" color="textSecondary">
                              CS/min
                            </Typography>
                          </Box>
                        </Grid>
                        <Grid item xs={6} sm={3}>
                          <Box textAlign="center">
                            <Typography variant="h5" color="warning.main">
                              {formatStat(champion.performance_metrics.avg_damage_per_min)}
                            </Typography>
                            <Typography variant="body2" color="textSecondary">
                              Dégâts/min
                            </Typography>
                          </Box>
                        </Grid>
                        <Grid item xs={6} sm={3}>
                          <Box textAlign="center">
                            <Typography variant="h5" color="info.main">
                              {formatStat(champion.performance_metrics.avg_vision_score)}
                            </Typography>
                            <Typography variant="body2" color="textSecondary">
                              Vision
                            </Typography>
                          </Box>
                        </Grid>
                      </Grid>

                      {/* Performance Chart */}
                      <Box mt={3}>
                        <Typography variant="h6" gutterBottom>
                          Radar des performances
                        </Typography>
                        <ResponsiveContainer width="100%" height={250}>
                          <RadarChart data={createRadarData(champion.performance_metrics)}>
                            <PolarGrid />
                            <PolarAngleAxis dataKey="stat" />
                            <PolarRadiusAxis angle={90} domain={[0, 100]} />
                            <Radar
                              name="Performance"
                              dataKey="value"
                              stroke="#8884d8"
                              fill="#8884d8"
                              fillOpacity={0.6}
                            />
                          </RadarChart>
                        </ResponsiveContainer>
                      </Box>
                    </Grid>

                    {/* Suggestions */}
                    <Grid item xs={12} md={4}>
                      <Typography variant="h6" gutterBottom>
                        Suggestions d'amélioration
                      </Typography>
                      <List>
                        {champion.improvement_suggestions.map((suggestion, idx) => (
                          <ListItem key={idx} sx={{ pl: 0 }}>
                            <ListItemIcon>
                              <Assessment color="primary" fontSize="small" />
                            </ListItemIcon>
                            <ListItemText
                              primary={suggestion}
                              primaryTypographyProps={{ variant: 'body2' }}
                            />
                          </ListItem>
                        ))}
                      </List>
                      
                      {/* Mastery Progress */}
                      <Box mt={2}>
                        <Typography variant="body2" color="textSecondary" gutterBottom>
                          Niveau de maîtrise
                        </Typography>
                        <LinearProgress
                          variant="determinate"
                          value={champion.mastery_score}
                          color={getMasteryColor(champion.mastery_score) as any}
                          sx={{ height: 8, borderRadius: 4 }}
                        />
                        <Typography variant="caption" color="textSecondary">
                          {champion.mastery_score.toFixed(0)}/100
                        </Typography>
                      </Box>
                    </Grid>
                  </Grid>
                </AccordionDetails>
              </Accordion>
            </Grid>
          ))}
          
          {(!data.champions || data.champions.length === 0) && (
            <Grid item xs={12}>
              <Alert severity="info">
                Aucune analyse détaillée disponible. Sélectionnez un rôle spécifique pour voir l'analyse par champion.
              </Alert>
            </Grid>
          )}
        </Grid>
      </TabPanel>

      {/* Tab 3: Role Performance */}
      <TabPanel value={tabValue} index={2}>
        <Grid container spacing={3}>
          {data.role_performance && Object.keys(data.role_performance).length > 0 ? (
            Object.entries(data.role_performance).map(([role, performance]) => (
              <Grid item xs={12} md={6} key={role}>
                <Card>
                  <CardHeader
                    title={ROLE_NAMES[role] || role}
                    avatar={<Sports color="primary" />}
                    action={
                      <Box display="flex" alignItems="center">
                        {getTrendIcon(performance.trend_direction)}
                        <Chip
                          label={formatWinRate(performance.win_rate)}
                          color={performance.win_rate >= 60 ? 'success' : performance.win_rate >= 50 ? 'primary' : 'error'}
                          size="small"
                          sx={{ ml: 1 }}
                        />
                      </Box>
                    }
                  />
                  <CardContent>
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          Parties jouées
                        </Typography>
                        <Typography variant="h6">
                          {performance.games_played}
                        </Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          KDA moyen
                        </Typography>
                        <Typography variant="h6">
                          {formatKDA(performance.avg_kda)}
                        </Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          CS/min
                        </Typography>
                        <Typography variant="h6">
                          {formatStat(performance.avg_cs_per_min)}
                        </Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          Score de performance
                        </Typography>
                        <Typography variant="h6">
                          {performance.performance_score.toFixed(0)}/100
                        </Typography>
                      </Grid>
                    </Grid>
                    
                    <Box mt={2}>
                      <LinearProgress
                        variant="determinate"
                        value={performance.performance_score}
                        color={performance.performance_score >= 70 ? 'success' : performance.performance_score >= 50 ? 'primary' : 'warning'}
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            ))
          ) : (
            <Grid item xs={12}>
              <Alert severity="info">
                Aucune donnée de performance par rôle disponible pour cette période.
              </Alert>
            </Grid>
          )}
        </Grid>
      </TabPanel>
    </Box>
  );
};