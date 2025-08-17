import React, { useMemo, useState } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Button,
  ToggleButton,
  ToggleButtonGroup,
  Divider,
} from '@mui/material';
import {
  TrendingUp,
  SportsEsports,
  Person,
  Timeline,
  ViewModule,
  ViewList,
} from '@mui/icons-material';
import { PerformanceHeatmap, HeatmapDataPoint } from '../components/heatmap/PerformanceHeatmap';
import { Row } from '../types';

export interface HeatmapViewProps {
  data: Row[];
  loading?: boolean;
  error?: string | null;
}

type HeatmapType = 'champion-role' | 'role-champion' | 'champion-time' | 'performance-time';
type MetricType = 'winrate' | 'kda' | 'games' | 'performance';

export const HeatmapView: React.FC<HeatmapViewProps> = ({
  data,
  loading = false,
  error = null,
}) => {
  const [heatmapType, setHeatmapType] = useState<HeatmapType>('champion-role');
  const [selectedMetric, setSelectedMetric] = useState<MetricType>('winrate');
  const [viewMode, setViewMode] = useState<'single' | 'grid'>('single');

  // Transformation des données selon le type de heatmap
  const heatmapData = useMemo((): HeatmapDataPoint[] => {
    if (!data.length) return [];

    const result: HeatmapDataPoint[] = [];

    switch (heatmapType) {
      case 'champion-role': {
        // Grouper par champion et rôle
        const grouped: Record<string, Record<string, Row[]>> = {};
        
        data.forEach(match => {
          if (!grouped[match.champion]) {
            grouped[match.champion] = {};
          }
          if (!grouped[match.champion][match.lane]) {
            grouped[match.champion][match.lane] = [];
          }
          grouped[match.champion][match.lane].push(match);
        });

        Object.entries(grouped).forEach(([champion, roles]) => {
          Object.entries(roles).forEach(([role, matches]) => {
            const totalGames = matches.length;
            const wins = matches.filter(m => m.win).length;
            const winrate = (wins / totalGames) * 100;
            const avgKDA = matches.reduce((sum, m) => {
              const kda = m.deaths > 0 ? (m.kills + m.assists) / m.deaths : m.kills + m.assists;
              return sum + kda;
            }, 0) / totalGames;

            let value: number;
            switch (selectedMetric) {
              case 'winrate':
                value = winrate;
                break;
              case 'kda':
                value = avgKDA;
                break;
              case 'games':
                value = totalGames;
                break;
              case 'performance':
                value = (winrate / 100) * avgKDA * Math.log(totalGames + 1);
                break;
              default:
                value = winrate;
            }

            result.push({
              x: champion,
              y: role,
              value,
              count: totalGames,
              metadata: {
                totalGames,
                wins,
                losses: totalGames - wins,
                avgKDA,
                avgKills: matches.reduce((sum, m) => sum + m.kills, 0) / totalGames,
                avgDeaths: matches.reduce((sum, m) => sum + m.deaths, 0) / totalGames,
                avgAssists: matches.reduce((sum, m) => sum + m.assists, 0) / totalGames,
                avgCS: matches.reduce((sum, m) => sum + (m.cs_per_min || 0), 0) / totalGames,
                lastPlayed: new Date(Math.max(...matches.map(m => new Date(m.game_creation).getTime()))),
              },
            });
          });
        });
        break;
      }

      case 'role-champion': {
        // Identique à champion-role mais avec x et y inversés
        const grouped: Record<string, Record<string, Row[]>> = {};
        
        data.forEach(match => {
          if (!grouped[match.lane]) {
            grouped[match.lane] = {};
          }
          if (!grouped[match.lane][match.champion]) {
            grouped[match.lane][match.champion] = [];
          }
          grouped[match.lane][match.champion].push(match);
        });

        Object.entries(grouped).forEach(([role, champions]) => {
          Object.entries(champions).forEach(([champion, matches]) => {
            const totalGames = matches.length;
            const wins = matches.filter(m => m.win).length;
            const winrate = (wins / totalGames) * 100;
            const avgKDA = matches.reduce((sum, m) => {
              const kda = m.deaths > 0 ? (m.kills + m.assists) / m.deaths : m.kills + m.assists;
              return sum + kda;
            }, 0) / totalGames;

            let value: number;
            switch (selectedMetric) {
              case 'winrate':
                value = winrate;
                break;
              case 'kda':
                value = avgKDA;
                break;
              case 'games':
                value = totalGames;
                break;
              case 'performance':
                value = (winrate / 100) * avgKDA * Math.log(totalGames + 1);
                break;
              default:
                value = winrate;
            }

            result.push({
              x: role,
              y: champion,
              value,
              count: totalGames,
              metadata: {
                totalGames,
                wins,
                losses: totalGames - wins,
                avgKDA,
                avgKills: matches.reduce((sum, m) => sum + m.kills, 0) / totalGames,
                avgDeaths: matches.reduce((sum, m) => sum + m.deaths, 0) / totalGames,
                avgAssists: matches.reduce((sum, m) => sum + m.assists, 0) / totalGames,
                avgCS: matches.reduce((sum, m) => sum + (m.cs_per_min || 0), 0) / totalGames,
                lastPlayed: new Date(Math.max(...matches.map(m => new Date(m.game_creation).getTime()))),
              },
            });
          });
        });
        break;
      }

      case 'champion-time': {
        // Grouper par champion et période temporelle (mois)
        const grouped: Record<string, Record<string, Row[]>> = {};
        
        data.forEach(match => {
          const month = new Date(match.game_creation).toISOString().slice(0, 7); // YYYY-MM
          
          if (!grouped[match.champion]) {
            grouped[match.champion] = {};
          }
          if (!grouped[match.champion][month]) {
            grouped[match.champion][month] = [];
          }
          grouped[match.champion][month].push(match);
        });

        Object.entries(grouped).forEach(([champion, months]) => {
          Object.entries(months).forEach(([month, matches]) => {
            const totalGames = matches.length;
            const wins = matches.filter(m => m.win).length;
            const winrate = (wins / totalGames) * 100;

            result.push({
              x: champion,
              y: month,
              value: selectedMetric === 'games' ? totalGames : winrate,
              count: totalGames,
              metadata: {
                totalGames,
                wins,
                losses: totalGames - wins,
                avgKDA: matches.reduce((sum, m) => {
                  const kda = m.deaths > 0 ? (m.kills + m.assists) / m.deaths : m.kills + m.assists;
                  return sum + kda;
                }, 0) / totalGames,
                avgKills: matches.reduce((sum, m) => sum + m.kills, 0) / totalGames,
                avgDeaths: matches.reduce((sum, m) => sum + m.deaths, 0) / totalGames,
                avgAssists: matches.reduce((sum, m) => sum + m.assists, 0) / totalGames,
                avgCS: matches.reduce((sum, m) => sum + (m.cs_per_min || 0), 0) / totalGames,
              },
            });
          });
        });
        break;
      }
    }

    return result;
  }, [data, heatmapType, selectedMetric]);

  // Statistiques générales
  const stats = useMemo(() => {
    if (!data.length) return null;

    const totalMatches = data.length;
    const uniqueChampions = new Set(data.map(m => m.champion)).size;
    const uniqueRoles = new Set(data.map(m => m.lane)).size;
    const totalWins = data.filter(m => m.win).length;
    const overallWinrate = (totalWins / totalMatches) * 100;

    // Top 3 combinaisons champion-rôle
    const combinations: Record<string, { count: number; wins: number }> = {};
    data.forEach(match => {
      const key = `${match.champion}-${match.lane}`;
      if (!combinations[key]) {
        combinations[key] = { count: 0, wins: 0 };
      }
      combinations[key].count++;
      if (match.win) combinations[key].wins++;
    });

    const topCombinations = Object.entries(combinations)
      .map(([key, stats]) => ({
        combination: key,
        count: stats.count,
        winrate: (stats.wins / stats.count) * 100,
      }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 3);

    return {
      totalMatches,
      uniqueChampions,
      uniqueRoles,
      overallWinrate,
      topCombinations,
    };
  }, [data]);

  const handleCellClick = (dataPoint: HeatmapDataPoint) => {
    console.log('Heatmap cell clicked:', dataPoint);
    // Ici on pourrait ouvrir un modal avec plus de détails
  };

  if (loading) {
    return (
      <Box>
        <Typography variant="h4" fontWeight="bold" gutterBottom>
          Heatmap de Performance
        </Typography>
        <Alert severity="info">
          Chargement des données en cours...
        </Alert>
      </Box>
    );
  }

  if (error) {
    return (
      <Box>
        <Typography variant="h4" fontWeight="bold" gutterBottom>
          Heatmap de Performance
        </Typography>
        <Alert severity="error">
          {error}
        </Alert>
      </Box>
    );
  }

  if (data.length === 0) {
    return (
      <Box>
        <Typography variant="h4" fontWeight="bold" gutterBottom>
          Heatmap de Performance
        </Typography>
        <Alert severity="info">
          Aucune donnée disponible. Importez des matches pour voir les heatmaps de performance.
        </Alert>
      </Box>
    );
  }

  const getHeatmapTitle = () => {
    switch (heatmapType) {
      case 'champion-role':
        return 'Performance par Champion et Rôle';
      case 'role-champion':
        return 'Performance par Rôle et Champion';
      case 'champion-time':
        return 'Évolution Temporelle par Champion';
      case 'performance-time':
        return 'Évolution de la Performance';
      default:
        return 'Heatmap de Performance';
    }
  };

  const getAxisLabels = () => {
    switch (heatmapType) {
      case 'champion-role':
        return { x: 'Champions', y: 'Rôles' };
      case 'role-champion':
        return { x: 'Rôles', y: 'Champions' };
      case 'champion-time':
        return { x: 'Champions', y: 'Mois' };
      case 'performance-time':
        return { x: 'Période', y: 'Métrique' };
      default:
        return { x: 'X', y: 'Y' };
    }
  };

  const axisLabels = getAxisLabels();

  return (
    <Box>
      {/* En-tête */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Heatmap de Performance
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Visualisation matricielle des performances par différentes dimensions
          </Typography>
        </Box>
        <ToggleButtonGroup
          value={viewMode}
          exclusive
          onChange={(_, newMode) => newMode && setViewMode(newMode)}
          size="small"
        >
          <ToggleButton value="single">
            <ViewModule />
          </ToggleButton>
          <ToggleButton value="grid">
            <ViewList />
          </ToggleButton>
        </ToggleButtonGroup>
      </Box>

      {/* Statistiques générales */}
      {stats && (
        <Grid container spacing={2} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Typography variant="h4" color="primary">
                  {stats.totalMatches}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Total matches
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Typography variant="h4" color="secondary">
                  {stats.uniqueChampions}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Champions joués
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Typography variant="h4" color="info.main">
                  {stats.uniqueRoles}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Rôles joués
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Typography 
                  variant="h4" 
                  color={stats.overallWinrate >= 50 ? 'success.main' : 'error.main'}
                >
                  {stats.overallWinrate.toFixed(1)}%
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Winrate global
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}

      {/* Contrôles */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Type de vue</InputLabel>
                <Select
                  value={heatmapType}
                  label="Type de vue"
                  onChange={(e) => setHeatmapType(e.target.value as HeatmapType)}
                >
                  <MenuItem value="champion-role">Champion × Rôle</MenuItem>
                  <MenuItem value="role-champion">Rôle × Champion</MenuItem>
                  <MenuItem value="champion-time">Champion × Temps</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Métrique</InputLabel>
                <Select
                  value={selectedMetric}
                  label="Métrique"
                  onChange={(e) => setSelectedMetric(e.target.value as MetricType)}
                >
                  <MenuItem value="winrate">Winrate (%)</MenuItem>
                  <MenuItem value="kda">KDA Moyen</MenuItem>
                  <MenuItem value="games">Nombre de parties</MenuItem>
                  <MenuItem value="performance">Score de performance</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <Box display="flex" flexWrap="wrap" gap={1}>
                <Chip 
                  icon={<SportsEsports />} 
                  label={`${heatmapData.length} points de données`} 
                  variant="outlined" 
                />
                <Chip 
                  icon={<Timeline />} 
                  label={getHeatmapTitle()} 
                  color="primary" 
                  variant="outlined" 
                />
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Heatmap principale */}
      <PerformanceHeatmap
        data={heatmapData}
        title={getHeatmapTitle()}
        xAxisLabel={axisLabels.x}
        yAxisLabel={axisLabels.y}
        metric={selectedMetric}
        colorScheme="blues"
        showLabels={true}
        showTooltip={true}
        showLegend={true}
        interactive={true}
        height={500}
        width={800}
        onCellClick={handleCellClick}
      />

      {/* Top combinations */}
      {stats && stats.topCombinations.length > 0 && (
        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Top combinaisons Champion-Rôle
            </Typography>
            <Grid container spacing={2}>
              {stats.topCombinations.map((combo, index) => (
                <Grid item xs={12} md={4} key={combo.combination}>
                  <Box 
                    sx={{ 
                      p: 2, 
                      border: 1, 
                      borderColor: 'divider', 
                      borderRadius: 1,
                      textAlign: 'center',
                    }}
                  >
                    <Typography variant="h6">
                      #{index + 1}
                    </Typography>
                    <Typography variant="body1" fontWeight="bold">
                      {combo.combination.replace('-', ' • ')}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {combo.count} parties • {combo.winrate.toFixed(1)}% WR
                    </Typography>
                  </Box>
                </Grid>
              ))}
            </Grid>
          </CardContent>
        </Card>
      )}
    </Box>
  );
};