import React, { useState, useMemo } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Typography,
  Grid,
  Avatar,
  Chip,
  IconButton,
  Tooltip,
  Button,
  Divider,
  LinearProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Switch,
  FormControlLabel,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Alert,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material';
import {
  Compare,
  TrendingUp,
  TrendingDown,
  EmojiEvents,
  SportsEsports,
  Timeline,
  SwapHoriz,
  Close,
  Add,
  ExpandMore,
  Visibility,
  VisibilityOff,
} from '@mui/icons-material';
import { InteractiveCharts, ChartSeries } from '../charts/InteractiveCharts';

export interface PlayerStats {
  id: string;
  name: string;
  rank: string;
  avatar?: string;
  totalGames: number;
  winRate: number;
  avgKDA: number;
  avgKills: number;
  avgDeaths: number;
  avgAssists: number;
  avgCS: number;
  avgVisionScore: number;
  avgGameDuration: number;
  avgGoldPerMinute: number;
  avgDamagePerMinute: number;
  topChampions: Array<{
    champion: string;
    games: number;
    winRate: number;
    avgKDA: number;
  }>;
  performanceByRole: Record<string, {
    games: number;
    winRate: number;
    avgKDA: number;
  }>;
  monthlyProgress: Array<{
    month: string;
    winRate: number;
    avgKDA: number;
    games: number;
  }>;
}

export interface PerformanceComparatorProps {
  players: PlayerStats[];
  onPlayerAdd?: () => void;
  onPlayerRemove?: (playerId: string) => void;
  maxPlayers?: number;
  showCharts?: boolean;
  showDetailedStats?: boolean;
}

const METRIC_CONFIGS = [
  { key: 'winRate', label: 'Winrate', unit: '%', format: (v: number) => `${v.toFixed(1)}%`, color: 'success' },
  { key: 'avgKDA', label: 'KDA Moyen', unit: '', format: (v: number) => v.toFixed(2), color: 'primary' },
  { key: 'avgKills', label: 'Kills/Game', unit: '', format: (v: number) => v.toFixed(1), color: 'error' },
  { key: 'avgDeaths', label: 'Deaths/Game', unit: '', format: (v: number) => v.toFixed(1), color: 'warning' },
  { key: 'avgAssists', label: 'Assists/Game', unit: '', format: (v: number) => v.toFixed(1), color: 'info' },
  { key: 'avgCS', label: 'CS/min', unit: '', format: (v: number) => v.toFixed(1), color: 'secondary' },
  { key: 'avgVisionScore', label: 'Vision Score', unit: '', format: (v: number) => v.toFixed(0), color: 'primary' },
  { key: 'avgGoldPerMinute', label: 'Gold/min', unit: '', format: (v: number) => v.toFixed(0), color: 'warning' },
  { key: 'avgDamagePerMinute', label: 'Damage/min', unit: '', format: (v: number) => v.toFixed(0), color: 'error' },
] as const;

export const PerformanceComparator: React.FC<PerformanceComparatorProps> = ({
  players,
  onPlayerAdd,
  onPlayerRemove,
  maxPlayers = 4,
  showCharts = true,
  showDetailedStats = true,
}) => {
  const [selectedMetrics, setSelectedMetrics] = useState<string[]>([
    'winRate', 'avgKDA', 'avgKills', 'avgDeaths', 'avgAssists', 'avgCS'
  ]);
  const [showPercentageComparison, setShowPercentageComparison] = useState(true);
  const [comparisonMode, setComparisonMode] = useState<'absolute' | 'relative'>('absolute');
  const [expandedSections, setExpandedSections] = useState<string[]>(['overview']);

  // Calcul des métriques de comparaison
  const comparisonData = useMemo(() => {
    if (players.length < 2) return null;

    const metrics = METRIC_CONFIGS.filter(m => selectedMetrics.includes(m.key));
    const comparisons: Record<string, any> = {};

    metrics.forEach(metric => {
      const values = players.map(p => (p as any)[metric.key] as number);
      const max = Math.max(...values);
      const min = Math.min(...values);
      const avg = values.reduce((sum, v) => sum + v, 0) / values.length;

      comparisons[metric.key] = {
        values,
        max,
        min,
        avg,
        range: max - min,
        leader: players[values.indexOf(max)],
        laggard: players[values.indexOf(min)],
      };
    });

    return comparisons;
  }, [players, selectedMetrics]);

  // Données pour les graphiques
  const chartData = useMemo((): ChartSeries[] => {
    if (!showCharts || players.length === 0) return [];

    // Graphique de progression mensuelle
    const progressionSeries = players.map((player, index) => ({
      name: player.name,
      data: player.monthlyProgress.map(month => ({
        x: month.month,
        y: month.winRate,
        metadata: { player: player.name, month: month.month, games: month.games, kda: month.avgKDA },
      })),
      color: ['#1976d2', '#388e3c', '#f57c00', '#d32f2f'][index % 4],
    }));

    return progressionSeries;
  }, [players, showCharts]);

  // Données pour le graphique radar de comparaison
  const radarData = useMemo((): ChartSeries[] => {
    if (players.length === 0) return [];

    const normalizedMetrics = ['winRate', 'avgKDA', 'avgCS', 'avgVisionScore'];
    
    return players.map((player, index) => ({
      name: player.name,
      data: normalizedMetrics.map(metric => {
        let value = (player as any)[metric];
        // Normalisation pour le radar (0-100)
        switch (metric) {
          case 'winRate':
            value = value; // Déjà en pourcentage
            break;
          case 'avgKDA':
            value = Math.min(value * 20, 100); // KDA * 20, max 100
            break;
          case 'avgCS':
            value = Math.min(value * 10, 100); // CS * 10, max 100
            break;
          case 'avgVisionScore':
            value = Math.min(value / 2, 100); // Vision / 2, max 100
            break;
        }
        return {
          x: metric.replace('avg', '').replace('Rate', ''),
          y: value,
        };
      }),
      color: ['#1976d2', '#388e3c', '#f57c00', '#d32f2f'][index % 4],
    }));
  }, [players]);

  // Gestion des sections dépliées
  const handleAccordionChange = (section: string) => (
    event: React.SyntheticEvent,
    isExpanded: boolean
  ) => {
    setExpandedSections(prev => 
      isExpanded 
        ? [...prev, section]
        : prev.filter(s => s !== section)
    );
  };

  // Rendu d'une métrique de comparaison
  const renderMetricComparison = (metricKey: string) => {
    const metric = METRIC_CONFIGS.find(m => m.key === metricKey);
    const comparison = comparisonData?.[metricKey];
    
    if (!metric || !comparison) return null;

    return (
      <TableRow key={metricKey}>
        <TableCell component="th" scope="row">
          <Typography variant="body2" fontWeight="medium">
            {metric.label}
          </Typography>
        </TableCell>
        {players.map((player, index) => {
          const value = (player as any)[metricKey] as number;
          const isLeader = comparison.leader.id === player.id;
          const isLaggard = comparison.laggard.id === player.id;
          const percentage = comparison.max > 0 ? (value / comparison.max) * 100 : 0;
          
          return (
            <TableCell key={player.id} align="center">
              <Box>
                <Typography 
                  variant="body2" 
                  fontWeight={isLeader ? 'bold' : 'normal'}
                  color={isLeader ? 'success.main' : isLaggard ? 'error.main' : 'text.primary'}
                >
                  {metric.format(value)}
                </Typography>
                {showPercentageComparison && comparison.range > 0 && (
                  <LinearProgress
                    variant="determinate"
                    value={percentage}
                    sx={{ 
                      mt: 0.5, 
                      height: 4,
                      bgcolor: 'grey.200',
                      '& .MuiLinearProgress-bar': {
                        bgcolor: isLeader ? 'success.main' : isLaggard ? 'error.main' : 'primary.main'
                      }
                    }}
                  />
                )}
              </Box>
            </TableCell>
          );
        })}
      </TableRow>
    );
  };

  if (players.length === 0) {
    return (
      <Card>
        <CardContent sx={{ textAlign: 'center', py: 6 }}>
          <Compare sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            Aucun joueur à comparer
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Ajoutez au moins 2 joueurs pour commencer la comparaison
          </Typography>
          {onPlayerAdd && (
            <Button variant="contained" startIcon={<Add />} onClick={onPlayerAdd}>
              Ajouter un joueur
            </Button>
          )}
        </CardContent>
      </Card>
    );
  }

  if (players.length === 1) {
    return (
      <Alert severity="info" sx={{ mb: 2 }}>
        Ajoutez au moins un autre joueur pour comparer les performances.
      </Alert>
    );
  }

  return (
    <Box>
      {/* En-tête avec contrôles */}
      <Card sx={{ mb: 3 }}>
        <CardHeader
          title={
            <Box display="flex" alignItems="center" gap={1}>
              <Compare color="primary" />
              <Typography variant="h6">
                Comparaison de performance ({players.length} joueurs)
              </Typography>
            </Box>
          }
          action={
            <Box display="flex" gap={1}>
              {onPlayerAdd && players.length < maxPlayers && (
                <Button
                  variant="outlined"
                  size="small"
                  startIcon={<Add />}
                  onClick={onPlayerAdd}
                >
                  Ajouter
                </Button>
              )}
              <FormControl size="small" sx={{ minWidth: 120 }}>
                <InputLabel>Mode</InputLabel>
                <Select
                  value={comparisonMode}
                  label="Mode"
                  onChange={(e) => setComparisonMode(e.target.value as 'absolute' | 'relative')}
                >
                  <MenuItem value="absolute">Absolue</MenuItem>
                  <MenuItem value="relative">Relative</MenuItem>
                </Select>
              </FormControl>
            </Box>
          }
        />
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <FormControlLabel
                control={
                  <Switch
                    checked={showPercentageComparison}
                    onChange={(e) => setShowPercentageComparison(e.target.checked)}
                  />
                }
                label="Afficher les barres de progression"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <Box display="flex" flexWrap="wrap" gap={1}>
                <Typography variant="body2" color="text.secondary" sx={{ mr: 1 }}>
                  Métriques:
                </Typography>
                {METRIC_CONFIGS.map(metric => (
                  <Chip
                    key={metric.key}
                    label={metric.label}
                    size="small"
                    clickable
                    color={selectedMetrics.includes(metric.key) ? 'primary' : 'default'}
                    onClick={() => {
                      setSelectedMetrics(prev =>
                        prev.includes(metric.key)
                          ? prev.filter(m => m !== metric.key)
                          : [...prev, metric.key]
                      );
                    }}
                  />
                ))}
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Cartes des joueurs */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        {players.map((player) => (
          <Grid item xs={12} sm={6} md={12 / Math.min(players.length, 4)} key={player.id}>
            <Card variant="outlined">
              <CardContent>
                <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
                  <Box display="flex" alignItems="center" gap={2}>
                    <Avatar src={player.avatar} sx={{ bgcolor: 'primary.main' }}>
                      {player.name[0]}
                    </Avatar>
                    <Box>
                      <Typography variant="h6">{player.name}</Typography>
                      <Chip label={player.rank} size="small" color="primary" />
                    </Box>
                  </Box>
                  {onPlayerRemove && (
                    <IconButton
                      size="small"
                      onClick={() => onPlayerRemove(player.id)}
                      color="error"
                    >
                      <Close />
                    </IconButton>
                  )}
                </Box>
                
                <Grid container spacing={1}>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Parties
                    </Typography>
                    <Typography variant="h6">{player.totalGames}</Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Winrate
                    </Typography>
                    <Typography variant="h6" color="success.main">
                      {player.winRate.toFixed(1)}%
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      KDA
                    </Typography>
                    <Typography variant="h6">{player.avgKDA.toFixed(2)}</Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      CS/min
                    </Typography>
                    <Typography variant="h6">{player.avgCS.toFixed(1)}</Typography>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Tableau de comparaison */}
      <Accordion 
        expanded={expandedSections.includes('comparison')}
        onChange={handleAccordionChange('comparison')}
      >
        <AccordionSummary expandIcon={<ExpandMore />}>
          <Typography variant="h6">Comparaison détaillée</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Métrique</TableCell>
                  {players.map(player => (
                    <TableCell key={player.id} align="center">
                      <Box display="flex" alignItems="center" justifyContent="center" gap={1}>
                        <Avatar src={player.avatar} sx={{ width: 24, height: 24 }}>
                          {player.name[0]}
                        </Avatar>
                        <Typography variant="body2">{player.name}</Typography>
                      </Box>
                    </TableCell>
                  ))}
                </TableRow>
              </TableHead>
              <TableBody>
                {selectedMetrics.map(metric => renderMetricComparison(metric))}
              </TableBody>
            </Table>
          </TableContainer>
        </AccordionDetails>
      </Accordion>

      {/* Graphiques de comparaison */}
      {showCharts && (
        <Accordion 
          expanded={expandedSections.includes('charts')}
          onChange={handleAccordionChange('charts')}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Typography variant="h6">Graphiques de comparaison</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <InteractiveCharts
                  title="Progression mensuelle du winrate"
                  series={chartData}
                  chartType="line"
                  height={300}
                  showControls={false}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <InteractiveCharts
                  title="Comparaison radar des performances"
                  series={radarData}
                  chartType="radar"
                  height={300}
                  showControls={false}
                />
              </Grid>
            </Grid>
          </AccordionDetails>
        </Accordion>
      )}
    </Box>
  );
};