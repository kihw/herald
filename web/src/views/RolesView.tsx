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
} from '@mui/material';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
} from 'recharts';
import { Row } from '../types';
import { RolesTable } from '../components/tables/RolesTable';

interface RolesViewProps {
  data: Row[];
  loading: boolean;
  error: string | null;
  onRoleSelect: (role: string) => void;
}

const COLORS = ['#C89B3C', '#0596AA', '#F0E6D2', '#A0752B', '#5BC0DE'];

export const RolesView: React.FC<RolesViewProps> = ({
  data,
  loading,
  error,
  onRoleSelect,
}) => {
  const { showSuccess, showError, showProgress, hideProgress } = useExportFeedback();
  
  const { exportToPNG, exportRolesToExcel, exportCombined, isExporting, exportProgress } = useExport({
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
  const chartData = useMemo(() => {
    if (!data.length) return [];

    const roleStats = data.reduce((acc, row) => {
      const role = row.lane || 'Unknown';
      if (!acc[role]) {
        acc[role] = {
          role,
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

      const stat = acc[role];
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

    return Object.values(roleStats).map((stat: any) => ({
      role: stat.role,
      games: stat.games,
      winrate: (stat.wins / stat.games) * 100,
      avgKda: stat.totalKda / stat.games,
      avgKp: (stat.totalKp / stat.games) * 100,
      avgCs: stat.totalCs / stat.games,
      avgGpm: stat.totalGpm / stat.games,
      avgDpm: stat.totalDpm / stat.games,
      avgVision: stat.totalVision / stat.games,
    })).sort((a, b) => b.games - a.games);
  }, [data]);

  if (loading) {
    return (
      <Box>
        <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
          Analyse par Rôles
        </Typography>
        
        <Grid container spacing={3}>
          {[...Array(6)].map((_, i) => (
            <Grid item xs={12} md={6} key={i}>
              <Card>
                <CardContent>
                  <Skeleton variant="text" width="60%" height={32} />
                  <Skeleton variant="rectangular" height={200} sx={{ mt: 2 }} />
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
          Analyse par Rôles
        </Typography>
        <Alert severity="info">
          Aucune donnée disponible pour l'analyse des rôles
        </Alert>
      </Box>
    );
  }

  return (
    <Box id="roles-view">
      <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
        Analyse par Rôles
      </Typography>
      
      {/* Graphiques de synthèse */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        {/* Répartition des games par rôle */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader title="Répartition des matchs par rôle" />
            <CardContent>
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie
                    data={chartData}
                    dataKey="games"
                    nameKey="role"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    label={({ role, games }) => `${role}: ${games}`}
                  >
                    {chartData.map((_, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Taux de victoire par rôle */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader title="Taux de victoire par rôle" />
            <CardContent>
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="role" />
                  <YAxis domain={[0, 100]} />
                  <Tooltip formatter={(value) => [`${typeof value === 'number' ? value.toFixed(1) : value}%`, 'Taux de victoire']} />
                  <Bar dataKey="winrate" fill="#C89B3C" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Performance KDA par rôle */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader title="KDA moyen par rôle" />
            <CardContent>
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="role" />
                  <YAxis />
                  <Tooltip formatter={(value) => [typeof value === 'number' ? value.toFixed(2) : value, 'KDA moyen']} />
                  <Bar dataKey="avgKda" fill="#0596AA" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* CS par minute par rôle */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader title="CS/min moyen par rôle" />
            <CardContent>
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="role" />
                  <YAxis />
                  <Tooltip formatter={(value) => [typeof value === 'number' ? value.toFixed(1) : value, 'CS/min']} />
                  <Bar dataKey="avgCs" fill="#A0752B" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Tableau détaillé des rôles */}
      <RolesTable
        data={data}
        onRoleSelect={onRoleSelect}
        onExportPNG={() => exportToPNG('roles-view', { filename: 'lol-analytics-roles' })}
        onExportExcel={() => exportRolesToExcel(data)}
      />
    </Box>
  );
};