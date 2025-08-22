import React, { useState, useEffect } from 'react';
import { getApiUrl } from '../utils/api-config';
import { OAuthCallback } from './auth/OAuthCallback';
import { GoogleAuth } from './auth/GoogleAuth';
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
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  TrendingFlat,
  EmojiEvents,
  Warning,
  Lightbulb,
  Timeline,
  Sports,
  Analytics,
} from '@mui/icons-material';
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, LineChart, Line, ResponsiveContainer } from 'recharts';

interface DashboardProps {
  puuid: string;
}

interface PeriodStats {
  period: string;
  total_games: number;
  win_rate: number;
  avg_kda: number;
  best_role: string;
  worst_role: string;
  top_champions: Array<{
    champion_name: string;
    games: number;
    win_rate: number;
    performance_score: number;
    avg_kda: number;
  }>;
  role_performance: Record<string, {
    games_played: number;
    win_rate: number;
    performance_score: number;
    avg_kda: number;
  }>;
  recent_trend: string;
  suggestions: string[];
}

interface PerformanceTrends {
  daily_trend: { trend: string; games: number; win_rate: number };
  weekly_trend: { trend: string; games: number; win_rate: number };
  monthly_trend: { trend: string; games: number; win_rate: number };
  seasonal_trend: { trend: string; games: number; win_rate: number };
  improvement_velocity: number;
  consistency_score: number;
  peak_performance: { period: string; performance: number; games: number };
}

interface Suggestion {
  type: string;
  title: string;
  description: string;
  priority: number;
}

interface DashboardData {
  user: any;
  period: string;
  stats: PeriodStats;
  trends: PerformanceTrends;
  suggestions: Suggestion[];
}

const ROLE_NAMES: Record<string, string> = {
  'TOP': 'Top',
  'JUNGLE': 'Jungle',
  'MIDDLE': 'Mid',
  'BOTTOM': 'ADC',
  'UTILITY': 'Support',
  'UNKNOWN': 'Autre'
};

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8'];

export const Dashboard: React.FC<DashboardProps> = ({ puuid }) => {
  const [period, setPeriod] = useState<string>('week');
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [user, setUser] = useState<any>(null);
  const [authChecked, setAuthChecked] = useState<boolean>(false);

  const API = getApiUrl('');

  // Vérifier l'authentification au chargement
  useEffect(() => {
    const checkAuth = () => {
      const storedUser = localStorage.getItem('herald_user');
      if (storedUser) {
        try {
          const userData = JSON.parse(storedUser);
          setUser(userData);
        } catch (e) {
          localStorage.removeItem('herald_user');
        }
      }
      setAuthChecked(true);
    };
    
    checkAuth();
  }, []);

  const handleAuthSuccess = (userData: any) => {
    setUser(userData);
    setAuthChecked(true);
  };

  const handleAuthError = (error: string) => {
    console.error('Auth error:', error);
    setAuthChecked(true);
  };

  const handleLogout = () => {
    localStorage.removeItem('herald_user');
    setUser(null);
  };

  useEffect(() => {
    fetchDashboardData();
  }, [puuid, period]);

  const fetchDashboardData = async () => {
    if (!puuid) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`${API}/api/users/${puuid}/dashboard/${period}`);
      if (!response.ok) {
        throw new Error(`Failed to fetch dashboard data: ${response.statusText}`);
      }
      
      const dashboardData = await response.json();
      setData(dashboardData);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
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

  const getTrendColor = (trend: string) => {
    switch (trend) {
      case 'improving':
        return 'success';
      case 'declining':
        return 'error';
      default:
        return 'default';
    }
  };

  const getPriorityIcon = (priority: number) => {
    switch (priority) {
      case 1:
        return <Warning color="error" />;
      case 2:
        return <Analytics color="warning" />;
      default:
        return <Lightbulb color="info" />;
    }
  };

  const formatWinRate = (winRate: number) => `${winRate.toFixed(1)}%`;
  const formatKDA = (kda: number) => kda.toFixed(2);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="h6" sx={{ ml: 2 }}>
          Chargement du dashboard...
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
        Aucune donnée disponible. Effectuez d'abord un export pour analyser vos performances.
      </Alert>
    );
  }

  const roleChartData = Object.entries(data.stats.role_performance).map(([role, perf]) => ({
    name: ROLE_NAMES[role] || role,
    games: perf.games_played,
    winRate: perf.win_rate,
    performance: perf.performance_score,
    kda: perf.avg_kda,
  }));

  const trendsData = [
    { period: 'Aujourd\'hui', ...data.trends.daily_trend },
    { period: 'Semaine', ...data.trends.weekly_trend },
    { period: 'Mois', ...data.trends.monthly_trend },
    { period: 'Saison', ...data.trends.seasonal_trend },
  ];

  // Gérer le callback OAuth
  if (!authChecked) {
    return <OAuthCallback onAuthSuccess={handleAuthSuccess} onAuthError={handleAuthError} />;
  }

  // Afficher l'authentification si l'utilisateur n'est pas connecté
  if (!user) {
    return (
      <Box sx={{ p: 2 }}>
        <Card sx={{ maxWidth: 500, mx: 'auto', mt: 4 }}>
          <CardContent>
            <Typography variant="h5" component="h1" gutterBottom textAlign="center">
              Herald.lol - Dashboard
            </Typography>
            <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 3 }}>
              Connectez-vous pour accéder à vos statistiques League of Legends
            </Typography>
            <GoogleAuth onSuccess={handleAuthSuccess} onError={handleAuthError} />
          </CardContent>
        </Card>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 2 }}>
      {/* Auth Status & Header */}
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Box>
          <Typography variant="h4" component="h1">
            Dashboard Analytics
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Connecté en tant que {user.name}
          </Typography>
        </Box>
        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
          <Chip 
            label={user.email} 
            variant="outlined" 
            size="small"
            avatar={user.picture ? <img src={user.picture} alt="avatar" style={{width: 24, height: 24, borderRadius: '50%'}} /> : undefined}
          />
          <Chip 
            label="Déconnexion" 
            onClick={handleLogout} 
            color="secondary" 
            variant="outlined" 
            size="small" 
          />
        </Box>
      </Box>
      
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h6" component="h2">
          Statistiques de performance
        </Typography>
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Période</InputLabel>
          <Select
            value={period}
            label="Période"
            onChange={(e) => setPeriod(e.target.value)}
          >
            <MenuItem value="today">Aujourd'hui</MenuItem>
            <MenuItem value="week">Cette semaine</MenuItem>
            <MenuItem value="month">Ce mois</MenuItem>
            <MenuItem value="season">Cette saison</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Grid container spacing={3}>
        {/* Performance Overview */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardHeader
              title="Vue d'ensemble des performances"
              subheader={`Période: ${data.period}`}
              avatar={<Analytics color="primary" />}
            />
            <CardContent>
              <Grid container spacing={2}>
                <Grid item xs={6} sm={3}>
                  <Box textAlign="center">
                    <Typography variant="h4" color="primary">
                      {data.stats.total_games}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Parties jouées
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Box textAlign="center">
                    <Typography variant="h4" color={data.stats.win_rate >= 50 ? 'success.main' : 'error.main'}>
                      {formatWinRate(data.stats.win_rate)}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Winrate
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Box textAlign="center">
                    <Typography variant="h4" color="secondary">
                      {formatKDA(data.stats.avg_kda)}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      KDA moyen
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Box textAlign="center">
                    <Box display="flex" alignItems="center" justifyContent="center">
                      {getTrendIcon(data.stats.recent_trend)}
                      <Chip
                        label={data.stats.recent_trend}
                        color={getTrendColor(data.stats.recent_trend) as any}
                        size="small"
                        sx={{ ml: 1 }}
                      />
                    </Box>
                    <Typography variant="body2" color="textSecondary">
                      Tendance récente
                    </Typography>
                  </Box>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Quick Stats */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardHeader
              title="Statistiques rapides"
              avatar={<Timeline color="primary" />}
            />
            <CardContent>
              <Box mb={2}>
                <Typography variant="body2" color="textSecondary">
                  Meilleur rôle
                </Typography>
                <Typography variant="h6" color="success.main">
                  {ROLE_NAMES[data.stats.best_role] || data.stats.best_role}
                </Typography>
              </Box>
              <Box mb={2}>
                <Typography variant="body2" color="textSecondary">
                  Rôle à améliorer
                </Typography>
                <Typography variant="h6" color="warning.main">
                  {ROLE_NAMES[data.stats.worst_role] || data.stats.worst_role}
                </Typography>
              </Box>
              <Box mb={2}>
                <Typography variant="body2" color="textSecondary">
                  Score de consistance
                </Typography>
                <Box display="flex" alignItems="center">
                  <LinearProgress
                    variant="determinate"
                    value={data.trends.consistency_score}
                    sx={{ flexGrow: 1, mr: 1 }}
                  />
                  <Typography variant="body2">
                    {data.trends.consistency_score.toFixed(0)}%
                  </Typography>
                </Box>
              </Box>
              <Box>
                <Typography variant="body2" color="textSecondary">
                  Vélocité d'amélioration
                </Typography>
                <Typography 
                  variant="h6" 
                  color={data.trends.improvement_velocity > 0 ? 'success.main' : 'error.main'}
                >
                  {data.trends.improvement_velocity > 0 ? '+' : ''}{data.trends.improvement_velocity.toFixed(1)}%
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Performance by Role */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader
              title="Performance par rôle"
              avatar={<Sports color="primary" />}
            />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={roleChartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip 
                    formatter={(value, name) => {
                      if (name === 'winRate') return [`${value}%`, 'Winrate'];
                      if (name === 'kda') return [value, 'KDA'];
                      if (name === 'performance') return [value, 'Score'];
                      return [value, name];
                    }}
                  />
                  <Legend />
                  <Bar dataKey="winRate" fill="#8884d8" name="Winrate %" />
                  <Bar dataKey="performance" fill="#82ca9d" name="Score de performance" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Trends Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader
              title="Tendances multi-temporelles"
              avatar={<Timeline color="primary" />}
            />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={trendsData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="period" />
                  <YAxis />
                  <Tooltip formatter={(value) => [`${value}%`, 'Winrate']} />
                  <Legend />
                  <Line 
                    type="monotone" 
                    dataKey="win_rate" 
                    stroke="#8884d8" 
                    name="Winrate %"
                    strokeWidth={2}
                  />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Top Champions */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader
              title="Top Champions"
              avatar={<EmojiEvents color="primary" />}
            />
            <CardContent>
              <List>
                {data.stats.top_champions.slice(0, 5).map((champion, index) => (
                  <React.Fragment key={champion.champion_name}>
                    <ListItem>
                      <ListItemIcon>
                        <Typography variant="h6" color="primary">
                          #{index + 1}
                        </Typography>
                      </ListItemIcon>
                      <ListItemText
                        primary={champion.champion_name}
                        secondary={
                          <Box>
                            <Typography variant="body2">
                              {champion.games} parties • {formatWinRate(champion.win_rate)} WR
                            </Typography>
                            <Typography variant="body2" color="textSecondary">
                              KDA: {formatKDA(champion.avg_kda)} • Score: {champion.performance_score.toFixed(0)}
                            </Typography>
                          </Box>
                        }
                      />
                      <Chip
                        label={formatWinRate(champion.win_rate)}
                        color={champion.win_rate >= 60 ? 'success' : champion.win_rate >= 50 ? 'primary' : 'error'}
                        size="small"
                      />
                    </ListItem>
                    {index < 4 && <Divider />}
                  </React.Fragment>
                ))}
              </List>
            </CardContent>
          </Card>
        </Grid>

        {/* AI Suggestions */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader
              title="Suggestions IA"
              avatar={<Lightbulb color="primary" />}
            />
            <CardContent>
              <List>
                {data.suggestions.slice(0, 5).map((suggestion, index) => (
                  <React.Fragment key={index}>
                    <ListItem>
                      <ListItemIcon>
                        {getPriorityIcon(suggestion.priority)}
                      </ListItemIcon>
                      <ListItemText
                        primary={suggestion.title}
                        secondary={suggestion.description}
                      />
                      <Chip
                        label={suggestion.type}
                        size="small"
                        variant="outlined"
                      />
                    </ListItem>
                    {index < Math.min(data.suggestions.length - 1, 4) && <Divider />}
                  </React.Fragment>
                ))}
              </List>
              {data.suggestions.length === 0 && (
                <Typography variant="body2" color="textSecondary" textAlign="center">
                  Jouez plus de parties pour recevoir des suggestions personnalisées
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Peak Performance */}
        {data.trends.peak_performance.period !== 'insufficient_data' && (
          <Grid item xs={12}>
            <Card>
              <CardHeader
                title="Performance de pointe"
                avatar={<EmojiEvents color="primary" />}
              />
              <CardContent>
                <Alert severity="success">
                  <Typography variant="body1">
                    <strong>Meilleure période:</strong> {data.trends.peak_performance.period}
                  </Typography>
                  <Typography variant="body2">
                    Winrate de {formatWinRate(data.trends.peak_performance.performance)} sur {data.trends.peak_performance.games} parties.
                    Analysez cette période pour reproduire ces performances !
                  </Typography>
                </Alert>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    </Box>
  );
};