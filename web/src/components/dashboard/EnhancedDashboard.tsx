import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  CardHeader,
  Tabs,
  Tab,
  Chip,
  Alert,
  CircularProgress,
  IconButton,
  Tooltip,
  Switch,
  FormControlLabel,
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  ShowChart,
  Analytics,
  Timeline,
  Speed,
  Insights,
  Settings,
  Refresh,
} from '@mui/icons-material';
import { RealTimePerformanceCharts } from '../charts/RealTimePerformanceCharts';
import { InteractiveCharts, ChartSeries } from '../charts/InteractiveCharts';

interface EnhancedDashboardProps {
  puuid: string;
}

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
      id={`dashboard-tabpanel-${index}`}
      aria-labelledby={`dashboard-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  );
}

export const EnhancedDashboard: React.FC<EnhancedDashboardProps> = ({ puuid }) => {
  const [currentTab, setCurrentTab] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date());

  // Données simulées pour les différents types de graphiques
  const [performanceOverviewData] = useState<ChartSeries[]>([
    {
      name: 'Score de Performance',
      data: [
        { x: 'Semaine 1', y: 75 },
        { x: 'Semaine 2', y: 82 },
        { x: 'Semaine 3', y: 78 },
        { x: 'Semaine 4', y: 85 },
        { x: 'Semaine 5', y: 88 },
        { x: 'Semaine 6', y: 91 },
      ],
      color: '#1976d2',
    },
    {
      name: 'Winrate (%)',
      data: [
        { x: 'Semaine 1', y: 65 },
        { x: 'Semaine 2', y: 68 },
        { x: 'Semaine 3', y: 62 },
        { x: 'Semaine 4', y: 72 },
        { x: 'Semaine 5', y: 75 },
        { x: 'Semaine 6', y: 78 },
      ],
      color: '#388e3c',
    },
  ]);

  const [championPerformanceData] = useState<ChartSeries[]>([
    {
      name: 'Performance par Champion',
      data: [
        { x: 'Jinx', y: 85, metadata: { games: 12, winrate: 75 } },
        { x: 'Yasuo', y: 72, metadata: { games: 8, winrate: 62 } },
        { x: 'Thresh', y: 91, metadata: { games: 15, winrate: 87 } },
        { x: 'Leona', y: 88, metadata: { games: 10, winrate: 80 } },
        { x: 'Zed', y: 76, metadata: { games: 6, winrate: 67 } },
        { x: 'Caitlyn', y: 83, metadata: { games: 9, winrate: 78 } },
      ],
      color: '#f57c00',
    },
  ]);

  const [roleDistributionData] = useState<ChartSeries[]>([
    {
      name: 'Répartition par Rôle',
      data: [
        { x: 'ADC', y: 35 },
        { x: 'Support', y: 25 },
        { x: 'Mid', y: 20 },
        { x: 'Top', y: 15 },
        { x: 'Jungle', y: 5 },
      ],
      color: '#7b1fa2',
    },
  ]);

  const [rankProgressionData] = useState<ChartSeries[]>([
    {
      name: 'Progression LP',
      data: [
        { x: '2024-01', y: 1150 },
        { x: '2024-02', y: 1180 },
        { x: '2024-03', y: 1165 },
        { x: '2024-04', y: 1195 },
        { x: '2024-05', y: 1220 },
        { x: '2024-06', y: 1245 },
        { x: '2024-07', y: 1260 },
        { x: '2024-08', y: 1275 },
      ],
      color: '#d32f2f',
    },
  ]);

  // Mise à jour manuelle des données
  const handleRefresh = async () => {
    setLoading(true);
    try {
      // Simulation d'un appel API
      await new Promise(resolve => setTimeout(resolve, 2000));
      setLastUpdate(new Date());
      setError(null);
    } catch (err) {
      setError('Erreur lors de la mise à jour des données');
    } finally {
      setLoading(false);
    }
  };

  // Gestion des clics sur les données
  const handleChartClick = (point: any, series: any) => {
    console.log('Données sélectionnées:', { point, series });
    // Ici on pourrait ouvrir un modal détaillé ou naviguer vers une vue spécifique
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  return (
    <Box>
      {/* En-tête du dashboard */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Dashboard Analytics Avancé
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Analyse en temps réel de vos performances League of Legends
          </Typography>
        </Box>
        <Box display="flex" alignItems="center" gap={2}>
          <FormControlLabel
            control={
              <Switch
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                color="primary"
              />
            }
            label="Auto-refresh"
          />
          <Tooltip title="Actualiser les données">
            <IconButton onClick={handleRefresh} disabled={loading}>
              <Refresh />
            </IconButton>
          </Tooltip>
          <Tooltip title="Paramètres">
            <IconButton>
              <Settings />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      {/* Indicateur de statut */}
      {loading && (
        <Box display="flex" alignItems="center" gap={1} mb={2}>
          <CircularProgress size={16} />
          <Typography variant="body2" color="text.secondary">
            Mise à jour en cours...
          </Typography>
        </Box>
      )}

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Chip
          icon={<Insights />}
          label={`PUUID: ${puuid.substring(0, 8)}...`}
          color="primary"
          variant="outlined"
        />
        <Typography variant="caption" color="text.secondary">
          Dernière mise à jour: {lastUpdate.toLocaleString()}
        </Typography>
      </Box>

      {/* Onglets du dashboard */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={currentTab} onChange={handleTabChange} aria-label="dashboard tabs">
          <Tab
            icon={<Speed />}
            label="Temps Réel"
            id="dashboard-tab-0"
            aria-controls="dashboard-tabpanel-0"
          />
          <Tab
            icon={<ShowChart />}
            label="Performance"
            id="dashboard-tab-1"
            aria-controls="dashboard-tabpanel-1"
          />
          <Tab
            icon={<Analytics />}
            label="Champions"
            id="dashboard-tab-2"
            aria-controls="dashboard-tabpanel-2"
          />
          <Tab
            icon={<Timeline />}
            label="Progression"
            id="dashboard-tab-3"
            aria-controls="dashboard-tabpanel-3"
          />
        </Tabs>
      </Box>

      {/* Contenu des onglets */}
      <TabPanel value={currentTab} index={0}>
        <RealTimePerformanceCharts
          puuid={puuid}
          refreshInterval={autoRefresh ? 30000 : 0}
          showLiveData={autoRefresh}
        />
      </TabPanel>

      <TabPanel value={currentTab} index={1}>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <InteractiveCharts
              title="Vue d'ensemble des performances"
              series={performanceOverviewData}
              chartType="line"
              height={400}
              showControls={true}
              onDataPointClick={handleChartClick}
            />
          </Grid>
          
          <Grid item xs={12} md={6}>
            <InteractiveCharts
              title="Répartition par rôle"
              series={roleDistributionData}
              chartType="doughnut"
              height={350}
              showControls={false}
            />
          </Grid>
          
          <Grid item xs={12} md={6}>
            <InteractiveCharts
              title="Performance comparative"
              series={performanceOverviewData}
              chartType="radar"
              height={350}
              showControls={false}
            />
          </Grid>
        </Grid>
      </TabPanel>

      <TabPanel value={currentTab} index={2}>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <InteractiveCharts
              title="Performance par champion"
              series={championPerformanceData}
              chartType="bar"
              height={400}
              showControls={true}
              onDataPointClick={handleChartClick}
            />
          </Grid>
          
          <Grid item xs={12} md={6}>
            <InteractiveCharts
              title="Analyse scatter performance/games"
              series={[
                {
                  name: 'Champions',
                  data: championPerformanceData[0].data.map(point => ({
                    x: point.metadata?.games || 0,
                    y: point.y,
                    label: point.x as string,
                    metadata: point.metadata,
                  })),
                  color: '#1976d2',
                }
              ]}
              chartType="scatter"
              height={350}
              showControls={false}
              onDataPointClick={handleChartClick}
            />
          </Grid>
          
          <Grid item xs={12} md={6}>
            <InteractiveCharts
              title="Heatmap Champion/Performance"
              series={[]}
              chartType="heatmap"
              height={350}
              showControls={false}
            />
          </Grid>
        </Grid>
      </TabPanel>

      <TabPanel value={currentTab} index={3}>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <InteractiveCharts
              title="Progression du rang (LP)"
              series={rankProgressionData}
              chartType="line"
              height={400}
              showControls={true}
              onDataPointClick={handleChartClick}
            />
          </Grid>
          
          <Grid item xs={12} md={6}>
            <InteractiveCharts
              title="Évolution du score de performance"
              series={performanceOverviewData.slice(0, 1)}
              chartType="polar"
              height={350}
              showControls={false}
            />
          </Grid>
          
          <Grid item xs={12} md={6}>
            <InteractiveCharts
              title="Timeline détaillée"
              series={rankProgressionData}
              chartType="d3-timeline"
              height={350}
              showControls={false}
              onDataPointClick={handleChartClick}
            />
          </Grid>
        </Grid>
      </TabPanel>

      {/* Informations additionnelles */}
      <Box mt={4}>
        <Card>
          <CardHeader
            title="Informations sur les données"
            titleTypographyProps={{ variant: 'h6' }}
          />
          <CardContent>
            <Grid container spacing={2}>
              <Grid item xs={12} md={4}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Source des données
                </Typography>
                <Typography variant="body1">
                  API Riot Games + Analytics IA
                </Typography>
              </Grid>
              <Grid item xs={12} md={4}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Fréquence de mise à jour
                </Typography>
                <Typography variant="body1">
                  {autoRefresh ? 'Temps réel (30s)' : 'Manuel'}
                </Typography>
              </Grid>
              <Grid item xs={12} md={4}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Période analysée
                </Typography>
                <Typography variant="body1">
                  30 derniers jours
                </Typography>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      </Box>
    </Box>
  );
};