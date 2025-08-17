import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Grid,
  Typography,
  Card,
  CardContent,
  IconButton,
  Tooltip,
  Alert,
  LinearProgress,
  Chip,
  Switch,
  FormControlLabel,
} from '@mui/material';
import {
  PlayArrow,
  Pause,
  Refresh,
  Timeline,
  TrendingUp,
  TrendingDown,
  ShowChart,
} from '@mui/icons-material';
import { InteractiveCharts, ChartSeries, ChartDataPoint } from './InteractiveCharts';

interface PerformanceMetric {
  id: string;
  name: string;
  value: number;
  trend: 'up' | 'down' | 'stable';
  change: number;
  unit: string;
  color: string;
}

interface RealTimePerformanceChartsProps {
  puuid?: string;
  refreshInterval?: number;
  showLiveData?: boolean;
}

export const RealTimePerformanceCharts: React.FC<RealTimePerformanceChartsProps> = ({
  puuid,
  refreshInterval = 30000, // 30 seconds
  showLiveData = true,
}) => {
  const [isLive, setIsLive] = useState(showLiveData);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date());
  
  // Données de performance en temps réel simulées
  const [performanceData, setPerformanceData] = useState<ChartSeries[]>([
    {
      name: 'MMR Évolution',
      data: generateTimeSeriesData('mmr', 30),
      color: '#1976d2',
    },
    {
      name: 'Winrate',
      data: generateTimeSeriesData('winrate', 30),
      color: '#388e3c',
    },
    {
      name: 'KDA Moyen',
      data: generateTimeSeriesData('kda', 30),
      color: '#f57c00',
    },
  ]);

  const [metrics, setMetrics] = useState<PerformanceMetric[]>([
    {
      id: 'mmr',
      name: 'MMR Estimé',
      value: 1247,
      trend: 'up',
      change: 15,
      unit: 'LP',
      color: '#1976d2',
    },
    {
      id: 'winrate',
      name: 'Winrate (7j)',
      value: 68,
      trend: 'up',
      change: 3.2,
      unit: '%',
      color: '#388e3c',
    },
    {
      id: 'kda',
      name: 'KDA Moyen',
      value: 2.4,
      trend: 'down',
      change: -0.1,
      unit: '',
      color: '#f57c00',
    },
    {
      id: 'cs_min',
      name: 'CS/min',
      value: 6.8,
      trend: 'stable',
      change: 0.0,
      unit: '',
      color: '#7b1fa2',
    },
  ]);

  // Génération de données temporelles simulées
  function generateTimeSeriesData(type: string, points: number): ChartDataPoint[] {
    const now = new Date();
    const data: ChartDataPoint[] = [];
    
    for (let i = points - 1; i >= 0; i--) {
      const timestamp = new Date(now.getTime() - i * 60000); // Chaque minute
      let value = 0;
      
      switch (type) {
        case 'mmr':
          value = 1200 + Math.random() * 100 + Math.sin(i / 5) * 20;
          break;
        case 'winrate':
          value = 60 + Math.random() * 20 + Math.cos(i / 3) * 10;
          break;
        case 'kda':
          value = 2 + Math.random() * 1.5 + Math.sin(i / 4) * 0.5;
          break;
        default:
          value = Math.random() * 100;
      }
      
      data.push({
        x: timestamp.toISOString(),
        y: Math.round(value * 100) / 100,
        label: `${type} at ${timestamp.toLocaleTimeString()}`,
        metadata: { type, timestamp },
      });
    }
    
    return data;
  }

  // Simulation de mise à jour des données en temps réel
  const updatePerformanceData = useCallback(async () => {
    if (!isLive) return;
    
    setLoading(true);
    try {
      // Simulation d'un appel API
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const now = new Date();
      
      // Mise à jour des séries de données
      setPerformanceData(prev => prev.map(series => {
        const newPoint: ChartDataPoint = {
          x: now.toISOString(),
          y: series.data[series.data.length - 1].y + (Math.random() - 0.5) * 10,
          label: `${series.name} at ${now.toLocaleTimeString()}`,
          metadata: { type: series.name, timestamp: now },
        };
        
        // Garder seulement les 30 derniers points
        const newData = [...series.data.slice(-29), newPoint];
        return { ...series, data: newData };
      }));
      
      // Mise à jour des métriques
      setMetrics(prev => prev.map(metric => ({
        ...metric,
        value: metric.value + (Math.random() - 0.5) * (metric.value * 0.02),
        change: (Math.random() - 0.5) * 5,
        trend: Math.random() > 0.5 ? 'up' : Math.random() > 0.5 ? 'down' : 'stable',
      })));
      
      setLastUpdate(now);
      setError(null);
    } catch (err) {
      setError('Erreur lors de la mise à jour des données');
    } finally {
      setLoading(false);
    }
  }, [isLive]);

  // Effet pour les mises à jour automatiques
  useEffect(() => {
    if (!isLive) return;
    
    const interval = setInterval(updatePerformanceData, refreshInterval);
    return () => clearInterval(interval);
  }, [isLive, refreshInterval, updatePerformanceData]);

  // Gestion des clics sur les points de données
  const handleDataPointClick = (point: ChartDataPoint, series: ChartSeries) => {
    console.log('Point cliqué:', { point, series });
    // Ici on pourrait ouvrir un modal avec plus de détails
  };

  // Gestion de l'export
  const handleExport = (format: 'png' | 'svg' | 'pdf') => {
    console.log('Export en format:', format);
    // Implémentation de l'export
  };

  return (
    <Box>
      {/* En-tête avec contrôles */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h5" fontWeight="bold">
          Performance en temps réel
        </Typography>
        <Box display="flex" alignItems="center" gap={2}>
          <FormControlLabel
            control={
              <Switch
                checked={isLive}
                onChange={(e) => setIsLive(e.target.checked)}
                color="primary"
              />
            }
            label="Temps réel"
          />
          <Tooltip title={isLive ? "Pause" : "Reprendre"}>
            <IconButton
              onClick={() => setIsLive(!isLive)}
              color={isLive ? "primary" : "default"}
            >
              {isLive ? <Pause /> : <PlayArrow />}
            </IconButton>
          </Tooltip>
          <Tooltip title="Actualiser maintenant">
            <IconButton onClick={updatePerformanceData} disabled={loading}>
              <Refresh />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      {/* Indicateur de statut */}
      <Box mb={2}>
        {loading && <LinearProgress />}
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box display="flex" alignItems="center" gap={1}>
            <Chip
              icon={isLive ? <Timeline /> : <Pause />}
              label={isLive ? "En direct" : "Pausé"}
              color={isLive ? "success" : "default"}
              size="small"
            />
            <Typography variant="caption" color="text.secondary">
              Dernière mise à jour: {lastUpdate.toLocaleTimeString()}
            </Typography>
          </Box>
        </Box>
      </Box>

      {/* Métriques rapides */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        {metrics.map((metric) => (
          <Grid item xs={12} sm={6} md={3} key={metric.id}>
            <Card>
              <CardContent sx={{ textAlign: 'center', py: 2 }}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  {metric.name}
                </Typography>
                <Typography variant="h4" fontWeight="bold" color={metric.color}>
                  {metric.value.toFixed(metric.unit === '' ? 1 : 0)}{metric.unit}
                </Typography>
                <Box display="flex" alignItems="center" justifyContent="center" gap={0.5}>
                  {metric.trend === 'up' && (
                    <TrendingUp fontSize="small" color="success" />
                  )}
                  {metric.trend === 'down' && (
                    <TrendingDown fontSize="small" color="error" />
                  )}
                  {metric.trend === 'stable' && (
                    <ShowChart fontSize="small" color="action" />
                  )}
                  <Typography
                    variant="caption"
                    color={
                      metric.trend === 'up'
                        ? 'success.main'
                        : metric.trend === 'down'
                        ? 'error.main'
                        : 'text.secondary'
                    }
                  >
                    {metric.change > 0 ? '+' : ''}{metric.change.toFixed(1)}
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Graphiques principaux */}
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <InteractiveCharts
            title="Évolution de la performance"
            series={performanceData}
            chartType="line"
            height={400}
            showControls={true}
            realTimeUpdate={isLive}
            onDataPointClick={handleDataPointClick}
            onExport={handleExport}
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <InteractiveCharts
            title="Distribution MMR"
            series={[
              {
                name: 'MMR Distribution',
                data: [
                  { x: 'Bronze', y: 15 },
                  { x: 'Silver', y: 25 },
                  { x: 'Gold', y: 30 },
                  { x: 'Platinum', y: 20 },
                  { x: 'Diamond', y: 10 },
                ],
              }
            ]}
            chartType="doughnut"
            height={300}
            showControls={false}
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <InteractiveCharts
            title="Performance par champion"
            series={[
              {
                name: 'Champions Performance',
                data: [
                  { x: 'Jinx', y: 85 },
                  { x: 'Yasuo', y: 75 },
                  { x: 'Thresh', y: 90 },
                  { x: 'Leona', y: 80 },
                  { x: 'Zed', y: 70 },
                ],
              }
            ]}
            chartType="bar"
            height={300}
            showControls={false}
          />
        </Grid>

        <Grid item xs={12}>
          <InteractiveCharts
            title="Heatmap Performance Champion/Rôle"
            series={[]}
            chartType="heatmap"
            height={400}
            showControls={false}
          />
        </Grid>
        
        <Grid item xs={12}>
          <InteractiveCharts
            title="Timeline des matches"
            series={performanceData}
            chartType="d3-timeline"
            height={300}
            showControls={false}
            onDataPointClick={handleDataPointClick}
          />
        </Grid>
      </Grid>
    </Box>
  );
};