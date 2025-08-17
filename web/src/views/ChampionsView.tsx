import React, { useMemo } from 'react';
import { useExport } from '../hooks/useExport';
import { useExportFeedback } from '../components/ExportFeedback';
import { 
  Typography, 
  Box,
  Alert,
  Grid,
  Card,
  CardContent,
  CardHeader,
  Skeleton,
  Chip,
} from '@mui/material';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ScatterChart,
  Scatter,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
} from 'recharts';
import { Row } from '../types';
import { ChampionsTable } from '../components/tables/ChampionsTable';

interface ChampionsViewProps {
  data: Row[];
  loading: boolean;
  error: string | null;
  selectedRole?: string;
  onChampionSelect: (champion: string) => void;
}

const COLORS = ['#C89B3C', '#0596AA', '#F0E6D2', '#A0752B', '#5BC0DE'];

export const ChampionsView: React.FC<ChampionsViewProps> = ({
  data,
  loading,
  error,
  selectedRole,
  onChampionSelect,
}) => {
  const { showSuccess, showError, showProgress, hideProgress } = useExportFeedback();
  
  const { exportToPNG, exportChampionsToExcel, exportCombined, isExporting, exportProgress } = useExport({
    onSuccess: showSuccess,
    onError: showError,
  });

  // Mettre à jour la progression si en cours d'export
  React.useEffect(() => {
    if (isExporting && exportProgress > 0) {
      showProgress('Export en cours...', exportProgress);
    } else {
      hideProgress();
    }
  }, [isExporting, exportProgress, showProgress, hideProgress]);
  const { chartData, topChampions, performanceData } = useMemo(() => {
    if (!data.length) return { chartData: [], topChampions: [], performanceData: [] };

    const championStats = data.reduce((acc, row) => {
      const champion = row.champion || 'Unknown';
      if (!acc[champion]) {
        acc[champion] = {
          champion,
          games: 0,
          wins: 0,
          totalKda: 0,
          totalKp: 0,
          totalCs: 0,
          totalGpm: 0,
          totalDpm: 0,
          totalVision: 0,
        };
      }

      const stat = acc[champion];
      stat.games++;
      if (row.win) stat.wins++;
      if (typeof row.kda === 'number') stat.totalKda += row.kda;
      if (typeof row.kp === 'number') stat.totalKp += row.kp;
      if (typeof row.cs_per_min === 'number') stat.totalCs += row.cs_per_min;
      if (typeof row.gpm === 'number') stat.totalGpm += row.gpm;
      if (typeof row.dpm === 'number') stat.totalDpm += row.dpm;
      if (typeof row.vision_score === 'number') stat.totalVision += row.vision_score;

      return acc;
    }, {} as Record<string, any>);

    const processedData = Object.values(championStats).map((stat: any) => ({
      champion: stat.champion,
      games: stat.games,
      winrate: (stat.wins / stat.games) * 100,
      avgKda: stat.totalKda / stat.games,
      avgKp: (stat.totalKp / stat.games) * 100,
      avgCs: stat.totalCs / stat.games,
      avgGpm: stat.totalGpm / stat.games,
      avgDpm: stat.totalDpm / stat.games,
      avgVision: stat.totalVision / stat.games,
    })).sort((a, b) => b.games - a.games);

    const topChampions = processedData.slice(0, 10);

    // Données pour le radar chart (top 5 champions)
    const performanceData = processedData.slice(0, 5).map(champ => ({
      champion: champ.champion,
      winrate: champ.winrate,
      kda: Math.min(champ.avgKda * 20, 100), // Normaliser KDA sur 100
      kp: champ.avgKp,
      cs: Math.min((champ.avgCs / 10) * 100, 100), // Normaliser CS/min sur 100
      vision: Math.min((champ.avgVision / 50) * 100, 100), // Normaliser vision sur 100
    }));

    return { chartData: processedData, topChampions, performanceData };
  }, [data]);

  if (loading) {
    return (
      <Box>
        <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
          Champions {selectedRole && `- ${selectedRole}`}
        </Typography>
        
        <Grid container spacing={3}>
          {[...Array(4)].map((_, i) => (
            <Grid item xs={12} md={6} key={i}>
              <Card>
                <CardContent>
                  <Skeleton variant="text" width="60%" height={32} />
                  <Skeleton variant="rectangular" height={250} sx={{ mt: 2 }} />
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 3 }}>
        {error}
      </Alert>
    );
  }

  if (!data.length) {
    return (
      <Box>
        <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
          Champions {selectedRole && `- ${selectedRole}`}
        </Typography>
        <Alert severity="info">
          Aucune donnée disponible pour l'analyse des champions
          {selectedRole && ` dans le rôle ${selectedRole}`}
        </Alert>
      </Box>
    );
  }

  return (
    <Box id="champions-view">
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
        <Typography variant="h4" sx={{ fontWeight: 700 }}>
          Champions
        </Typography>
        {selectedRole && (
          <Chip 
            label={selectedRole} 
            color="primary" 
            variant="outlined"
            sx={{ fontWeight: 600 }}
          />
        )}
      </Box>
      
      {/* Graphiques de synthèse */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        {/* Top 10 champions par nombre de games */}
        <Grid item xs={12} lg={6}>
          <Card>
            <CardHeader title="Top 10 champions les plus joués" />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={topChampions} layout="horizontal">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" />
                  <YAxis dataKey="champion" type="category" width={80} />
                  <Tooltip 
                    formatter={(value, name) => [
                      name === 'games' ? `${value} matchs` : value,
                      name === 'games' ? 'Matchs joués' : name
                    ]} 
                  />
                  <Bar dataKey="games" fill="#C89B3C" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Scatter plot: Winrate vs Games */}
        <Grid item xs={12} lg={6}>
          <Card>
            <CardHeader title="Relation Taux de victoire / Expérience" />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <ScatterChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis 
                    type="number" 
                    dataKey="games" 
                    name="Matchs joués"
                    domain={['dataMin', 'dataMax']}
                  />
                  <YAxis 
                    type="number" 
                    dataKey="winrate" 
                    name="Taux de victoire"
                    domain={[0, 100]}
                  />
                  <Tooltip 
                    cursor={{ strokeDasharray: '3 3' }}
                    formatter={(value, name) => [
                      name === 'winrate' ? `${typeof value === 'number' ? value.toFixed(1) : value}%` : value,
                      name === 'winrate' ? 'Taux de victoire' : 'Matchs joués'
                    ]}
                    labelFormatter={(label) => `Champion: ${label}`}
                  />
                  <Scatter dataKey="winrate" fill="#0596AA" />
                </ScatterChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Radar chart performance top 5 */}
        <Grid item xs={12} lg={6}>
          <Card>
            <CardHeader title="Performance multi-critères (Top 5)" />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <RadarChart data={performanceData}>
                  <PolarGrid />
                  <PolarAngleAxis dataKey="champion" />
                  <PolarRadiusAxis domain={[0, 100]} tick={false} />
                  <Radar
                    name="Winrate"
                    dataKey="winrate"
                    stroke="#C89B3C"
                    fill="#C89B3C"
                    fillOpacity={0.3}
                  />
                  <Radar
                    name="KDA"
                    dataKey="kda"
                    stroke="#0596AA"
                    fill="#0596AA"
                    fillOpacity={0.3}
                  />
                  <Tooltip 
                    formatter={(value) => [`${typeof value === 'number' ? value.toFixed(1) : value}`, '']}
                  />
                </RadarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Distribution winrate */}
        <Grid item xs={12} lg={6}>
          <Card>
            <CardHeader title="Distribution des taux de victoire" />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={topChampions}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="champion" angle={-45} textAnchor="end" height={80} />
                  <YAxis domain={[0, 100]} />
                  <Tooltip 
                    formatter={(value) => [`${typeof value === 'number' ? value.toFixed(1) : value}%`, 'Taux de victoire']}
                  />
                  <Bar 
                    dataKey="winrate" 
                    fill="#A0752B"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Tableau détaillé des champions */}
      <ChampionsTable
        data={data}
        selectedRole={selectedRole}
        onChampionSelect={onChampionSelect}
        onExportPNG={() => {
          const filename = selectedRole 
            ? `lol-analytics-champions-${selectedRole.toLowerCase()}`
            : 'lol-analytics-champions';
          exportToPNG('champions-view', { filename });
        }}
        onExportExcel={() => exportChampionsToExcel(data, selectedRole)}
      />
    </Box>
  );
};