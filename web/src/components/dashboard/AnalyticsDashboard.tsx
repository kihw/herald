import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  Tabs,
  Tab,
  Alert,
  IconButton,
  Tooltip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Card,
  CardContent,
  CardHeader,
  Grid,
  Chip,
  LinearProgress,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Avatar,
  Button
} from '@mui/material';
import {
  Analytics,
  Psychology,
  EmojiEvents,
  TrendingUp,
  TrendingDown,
  TrendingFlat,
  Assessment,
  Refresh,
  Insights,
  Timeline,
  Speed
} from '@mui/icons-material';

import { ChampionsAnalytics } from '../ChampionsAnalytics';
import { MMRAnalytics } from '../MMRAnalytics';
import ToastNotification from '../common/ToastNotification';
import LoadingSkeleton from '../common/LoadingSkeleton';
import { useAuth } from '../../context/AuthContext';

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
      id={`analytics-tabpanel-${index}`}
      aria-labelledby={`analytics-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
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
  role_performance: Record<string, any>;
  recent_trend: string;
  suggestions: string[];
}

interface Recommendation {
  type: string;
  title: string;
  description: string;
  priority: number;
  confidence: number;
  expected_improvement: string;
  action_items: string[];
  role?: string;
  time_period: string;
}

const PERIODS = [
  { value: 'today', label: "Aujourd'hui" },
  { value: 'week', label: 'Cette semaine' },
  { value: 'month', label: 'Ce mois' },
  { value: 'season', label: 'Cette saison' }
];

const ROLE_NAMES: Record<string, string> = {
  'TOP': 'Top',
  'JUNGLE': 'Jungle', 
  'MIDDLE': 'Mid',
  'BOTTOM': 'ADC',
  'UTILITY': 'Support'
};

const AnalyticsDashboard: React.FC = () => {
  const { state } = useAuth();
  const { user } = state;
  const [tabValue, setTabValue] = useState(0);
  const [selectedPeriod, setSelectedPeriod] = useState('week');
  
  // Analytics data state
  const [periodStats, setPeriodStats] = useState<PeriodStats | null>(null);
  const [recommendations, setRecommendations] = useState<Recommendation[]>([]);
  const [performanceTrends, setPerformanceTrends] = useState<any>(null);
  
  // Loading states
  const [statsLoading, setStatsLoading] = useState(true);
  const [recommendationsLoading, setRecommendationsLoading] = useState(true);
  const [trendsLoading, setTrendsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  
  // Error state
  const [error, setError] = useState<string | null>(null);
  
  // Notification state
  const [notification, setNotification] = useState({
    open: false,
    message: '',
    title: '',
    severity: 'success' as 'success' | 'error' | 'warning' | 'info',
  });

  const API = (import.meta.env?.VITE_API_BASE || 'http://localhost:8004');

  useEffect(() => {
    if (user) {
      loadAnalyticsData();
    }
  }, [user, selectedPeriod]);

  const loadAnalyticsData = async () => {
    await Promise.all([
      fetchPeriodStats(),
      fetchRecommendations(),
      fetchPerformanceTrends()
    ]);
  };

  const fetchPeriodStats = async () => {
    setStatsLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`${API}/api/analytics/period/${selectedPeriod}`, {
        credentials: 'include'
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch period stats: ${response.statusText}`);
      }
      
      const result = await response.json();
      
      if (!result.success) {
        throw new Error(result.message || 'Failed to fetch period stats');
      }
      
      setPeriodStats(result.data);
    } catch (err: any) {
      setError(err.message);
      showNotification('error', 'Erreur', 'Impossible de charger les statistiques');
    } finally {
      setStatsLoading(false);
    }
  };

  const fetchRecommendations = async () => {
    setRecommendationsLoading(true);
    
    try {
      const response = await fetch(`${API}/api/analytics/recommendations`, {
        credentials: 'include'
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch recommendations: ${response.statusText}`);
      }
      
      const result = await response.json();
      
      if (!result.success) {
        throw new Error(result.message || 'Failed to fetch recommendations');
      }
      
      setRecommendations(result.data);
    } catch (err: any) {
      console.error('Failed to fetch recommendations:', err);
      // Don't show error for recommendations as it's not critical
    } finally {
      setRecommendationsLoading(false);
    }
  };

  const fetchPerformanceTrends = async () => {
    setTrendsLoading(true);
    
    try {
      const response = await fetch(`${API}/api/analytics/trends`, {
        credentials: 'include'
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch trends: ${response.statusText}`);
      }
      
      const result = await response.json();
      
      if (!result.success) {
        throw new Error(result.message || 'Failed to fetch trends');
      }
      
      setPerformanceTrends(result.data);
    } catch (err: any) {
      console.error('Failed to fetch trends:', err);
    } finally {
      setTrendsLoading(false);
    }
  };

  const refreshAnalytics = async () => {
    setRefreshing(true);
    
    try {
      // Call refresh endpoint
      const response = await fetch(`${API}/api/analytics/refresh`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ period: selectedPeriod })
      });
      
      if (!response.ok) {
        throw new Error(`Failed to refresh analytics: ${response.statusText}`);
      }
      
      const result = await response.json();
      
      if (!result.success) {
        throw new Error(result.message || 'Failed to refresh analytics');
      }
      
      // Reload all data
      await loadAnalyticsData();
      
      showNotification('success', 'Succès', 'Analytics rafraîchies avec succès');
    } catch (err: any) {
      showNotification('error', 'Erreur', 'Impossible de rafraîchir les analytics');
    } finally {
      setRefreshing(false);
    }
  };

  const showNotification = (severity: 'success' | 'error' | 'warning' | 'info', title: string, message: string) => {
    setNotification({ open: true, severity, title, message });
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

  const formatWinRate = (winRate: number) => `${winRate.toFixed(1)}%`;
  const formatKDA = (kda: number) => kda.toFixed(2);

  if (!user) {
    return (
      <Alert severity="warning" sx={{ m: 2 }}>
        Veuillez vous connecter pour accéder aux analytics.
      </Alert>
    );
  }

  return (
    <Container maxWidth="xl" sx={{ py: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
        <Typography variant="h4" component="h1" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Analytics color="primary" />
          Dashboard Analytics
        </Typography>
        
        <Box display="flex" gap={2} alignItems="center">
          <FormControl size="small" sx={{ minWidth: 150 }}>
            <InputLabel>Période</InputLabel>
            <Select
              value={selectedPeriod}
              label="Période"
              onChange={(e) => setSelectedPeriod(e.target.value)}
            >
              {PERIODS.map(period => (
                <MenuItem key={period.value} value={period.value}>
                  {period.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          
          <Tooltip title="Rafraîchir les analytics">
            <IconButton 
              onClick={refreshAnalytics} 
              disabled={refreshing}
              color="primary"
            >
              <Refresh className={refreshing ? 'rotating' : ''} />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={tabValue} onChange={(e, newValue) => setTabValue(newValue)}>
          <Tab label="Vue d'ensemble" icon={<Assessment />} />
          <Tab label="Analyse des Champions" icon={<EmojiEvents />} />
          <Tab label="MMR & Progression" icon={<Timeline />} />
          <Tab label="Recommandations IA" icon={<Psychology />} />
        </Tabs>
      </Box>

      {/* Tab 1: Overview */}
      <TabPanel value={tabValue} index={0}>
        {statsLoading ? (
          <LoadingSkeleton />
        ) : error ? (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        ) : periodStats ? (
          <Grid container spacing={3}>
            {/* Main Stats Cards */}
            <Grid item xs={12} md={3}>
              <Card>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h3" color="primary">
                    {periodStats.total_games}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Parties jouées
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Card>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h3" color={periodStats.win_rate >= 60 ? 'success.main' : 'text.primary'}>
                    {formatWinRate(periodStats.win_rate)}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Taux de victoire
                  </Typography>
                  <Box display="flex" justifyContent="center" mt={1}>
                    {getTrendIcon(periodStats.recent_trend)}
                  </Box>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Card>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h3" color="secondary.main">
                    {formatKDA(periodStats.avg_kda)}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    KDA moyen
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Card>
                <CardContent sx={{ textAlign: 'center' }}>
                  <Chip
                    label={ROLE_NAMES[periodStats.best_role] || periodStats.best_role}
                    color="success"
                    sx={{ fontSize: '1.1rem', py: 2 }}
                  />
                  <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
                    Meilleur rôle
                  </Typography>
                </CardContent>
              </Card>
            </Grid>

            {/* Top Champions */}
            <Grid item xs={12} md={6}>
              <Card>
                <CardHeader
                  title="Top Champions"
                  subheader={PERIODS.find(p => p.value === selectedPeriod)?.label}
                  avatar={<EmojiEvents color="primary" />}
                />
                <CardContent>
                  {periodStats.top_champions && periodStats.top_champions.length > 0 ? (
                    <List>
                      {periodStats.top_champions.slice(0, 5).map((champion, index) => (
                        <ListItem key={champion.champion_name}>
                          <ListItemIcon>
                            <Avatar sx={{ bgcolor: 'primary.main', width: 32, height: 32 }}>
                              {index + 1}
                            </Avatar>
                          </ListItemIcon>
                          <ListItemText
                            primary={champion.champion_name}
                            secondary={`${champion.games} parties • ${formatWinRate(champion.win_rate)} • KDA: ${formatKDA(champion.avg_kda)}`}
                          />
                          <Box sx={{ minWidth: 80 }}>
                            <LinearProgress
                              variant="determinate"
                              value={champion.performance_score || 0}
                              color={(champion.performance_score || 0) >= 70 ? 'success' : (champion.performance_score || 0) >= 50 ? 'primary' : 'warning'}
                            />
                            <Typography variant="caption" color="textSecondary">
                              {champion.performance_score?.toFixed(0) || 0}/100
                            </Typography>
                          </Box>
                        </ListItem>
                      ))}
                    </List>
                  ) : (
                    <Typography variant="body2" color="textSecondary" textAlign="center">
                      Aucun champion trouvé pour cette période.
                    </Typography>
                  )}
                </CardContent>
              </Card>
            </Grid>

            {/* Performance Trends */}
            <Grid item xs={12} md={6}>
              <Card>
                <CardHeader
                  title="Tendances de Performance"
                  avatar={<Speed color="primary" />}
                />
                <CardContent>
                  {trendsLoading ? (
                    <Box display="flex" justifyContent="center" py={3}>
                      <LoadingSkeleton />
                    </Box>
                  ) : performanceTrends ? (
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          Tendance quotidienne
                        </Typography>
                        <Box display="flex" alignItems="center">
                          {getTrendIcon(performanceTrends.daily_trend?.trend)}
                          <Typography variant="h6" sx={{ ml: 1 }}>
                            {formatWinRate(performanceTrends.daily_trend?.win_rate || 0)}
                          </Typography>
                        </Box>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          Tendance hebdomadaire
                        </Typography>
                        <Box display="flex" alignItems="center">
                          {getTrendIcon(performanceTrends.weekly_trend?.trend)}
                          <Typography variant="h6" sx={{ ml: 1 }}>
                            {formatWinRate(performanceTrends.weekly_trend?.win_rate || 0)}
                          </Typography>
                        </Box>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          Score de consistance
                        </Typography>
                        <Typography variant="h6" color="info.main">
                          {performanceTrends.consistency_score?.toFixed(1) || 0}/100
                        </Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="body2" color="textSecondary">
                          Pic de performance
                        </Typography>
                        <Typography variant="h6" color="success.main">
                          {performanceTrends.peak_performance?.performance?.toFixed(1) || 0}%
                        </Typography>
                      </Grid>
                    </Grid>
                  ) : (
                    <Typography variant="body2" color="textSecondary" textAlign="center">
                      Données de tendances non disponibles.
                    </Typography>
                  )}
                </CardContent>
              </Card>
            </Grid>

            {/* Suggestions */}
            <Grid item xs={12}>
              <Card>
                <CardHeader
                  title="Suggestions d'Amélioration"
                  subheader={`Basées sur vos performances ${PERIODS.find(p => p.value === selectedPeriod)?.label.toLowerCase()}`}
                  avatar={<Insights color="primary" />}
                />
                <CardContent>
                  {periodStats.suggestions && periodStats.suggestions.length > 0 ? (
                    <List>
                      {periodStats.suggestions.map((suggestion, index) => (
                        <ListItem key={index}>
                          <ListItemIcon>
                            <Assessment color="primary" fontSize="small" />
                          </ListItemIcon>
                          <ListItemText primary={suggestion} />
                        </ListItem>
                      ))}
                    </List>
                  ) : (
                    <Typography variant="body2" color="textSecondary" textAlign="center">
                      Aucune suggestion disponible pour le moment.
                    </Typography>
                  )}
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        ) : (
          <Alert severity="info">
            Aucune donnée disponible pour cette période.
          </Alert>
        )}
      </TabPanel>

      {/* Tab 2: Champions Analytics */}
      <TabPanel value={tabValue} index={1}>
        <ChampionsAnalytics puuid={user.riot_puuid || user.id.toString()} />
      </TabPanel>

      {/* Tab 3: MMR Analytics */}
      <TabPanel value={tabValue} index={2}>
        <MMRAnalytics puuid={user.riot_puuid || user.id.toString()} />
      </TabPanel>

      {/* Tab 4: AI Recommendations */}
      <TabPanel value={tabValue} index={3}>
        {recommendationsLoading ? (
          <LoadingSkeleton />
        ) : (
          <Grid container spacing={3}>
            {recommendations.length > 0 ? (
              recommendations.map((rec, index) => (
                <Grid item xs={12} md={6} key={index}>
                  <Card>
                    <CardHeader
                      title={rec.title}
                      subheader={`Priorité: ${rec.priority} • Confiance: ${(rec.confidence * 100).toFixed(0)}%`}
                      avatar={
                        <Avatar sx={{ 
                          bgcolor: rec.priority === 1 ? 'error.main' : rec.priority === 2 ? 'warning.main' : 'info.main' 
                        }}>
                          {rec.priority}
                        </Avatar>
                      }
                    />
                    <CardContent>
                      <Typography variant="body1" gutterBottom>
                        {rec.description}
                      </Typography>
                      
                      <Box sx={{ mt: 2, mb: 2 }}>
                        <Chip
                          label={rec.expected_improvement}
                          color="success"
                          size="small"
                        />
                        <Chip
                          label={rec.time_period}
                          color="info"
                          size="small"
                          sx={{ ml: 1 }}
                        />
                        {rec.role && (
                          <Chip
                            label={ROLE_NAMES[rec.role] || rec.role}
                            color="primary"
                            size="small"
                            sx={{ ml: 1 }}
                          />
                        )}
                      </Box>

                      <Typography variant="subtitle2" gutterBottom>
                        Actions recommandées:
                      </Typography>
                      <List dense>
                        {rec.action_items.map((action, actionIndex) => (
                          <ListItem key={actionIndex} sx={{ pl: 0 }}>
                            <ListItemIcon>
                              <Assessment color="primary" fontSize="small" />
                            </ListItemIcon>
                            <ListItemText 
                              primary={action}
                              primaryTypographyProps={{ variant: 'body2' }}
                            />
                          </ListItem>
                        ))}
                      </List>
                    </CardContent>
                  </Card>
                </Grid>
              ))
            ) : (
              <Grid item xs={12}>
                <Alert severity="info">
                  Aucune recommandation disponible pour le moment. Jouez quelques parties pour obtenir des suggestions personnalisées !
                </Alert>
              </Grid>
            )}
          </Grid>
        )}
      </TabPanel>

      {/* Toast Notification */}
      <ToastNotification
        open={notification.open}
        message={notification.message}
        title={notification.title}
        severity={notification.severity}
        onClose={() => setNotification({ ...notification, open: false })}
      />

      <style jsx>{`
        .rotating {
          animation: rotate 1s linear infinite;
        }
        
        @keyframes rotate {
          from {
            transform: rotate(0deg);
          }
          to {
            transform: rotate(360deg);
          }
        }
      `}</style>
    </Container>
  );
};

export default AnalyticsDashboard;