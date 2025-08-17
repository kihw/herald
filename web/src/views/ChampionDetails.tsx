import React, { useMemo } from 'react';
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
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Avatar,
  LinearProgress,
} from '@mui/material';
import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import { Row } from '../types';
import { ProgressionTimeline } from '../components/timeline/ProgressionTimeline';
import dayjs from 'dayjs';

interface ChampionDetailsProps {
  data: Row[];
  loading: boolean;
  error: string | null;
  selectedChampion?: string;
  selectedRole?: string;
}

const COLORS = ['#C89B3C', '#0596AA', '#F0E6D2', '#A0752B', '#5BC0DE'];

// Fonction pour obtenir l'URL de l'icône champion
const getChampionIconUrl = (championName: string, patch = '14.23.1') => {
  const champId = championName.replace(/[^A-Za-z0-9]/g, '');
  return `https://ddragon.leagueoflegends.com/cdn/${patch}/img/champion/${champId}.png`;
};

export const ChampionDetails: React.FC<ChampionDetailsProps> = ({
  data,
  loading,
  error,
  selectedChampion,
  selectedRole,
}) => {
  const championData = useMemo(() => {
    if (!selectedChampion || !data.length) return null;

    const championMatches = data.filter(row => row.champion === selectedChampion);
    if (!championMatches.length) return null;

    // Statistiques globales
    const totalGames = championMatches.length;
    const wins = championMatches.filter(match => match.win).length;
    const winrate = (wins / totalGames) * 100;

    // Moyennes
    const avgKda = championMatches.reduce((sum, match) => sum + (match.kda || 0), 0) / totalGames;
    const avgKp = championMatches.reduce((sum, match) => sum + (match.kp || 0), 0) / totalGames;
    const avgCs = championMatches.reduce((sum, match) => sum + (match.cs_per_min || 0), 0) / totalGames;
    const avgGpm = championMatches.reduce((sum, match) => sum + (match.gpm || 0), 0) / totalGames;
    const avgDpm = championMatches.reduce((sum, match) => sum + (match.dpm || 0), 0) / totalGames;
    const avgVision = championMatches.reduce((sum, match) => sum + (match.vision_score || 0), 0) / totalGames;

    // Évolution dans le temps (derniers 20 matchs)
    const timelineData = championMatches
      .slice(-20)
      .map((match, index) => ({
        game: index + 1,
        winrate: match.win ? 100 : 0,
        kda: match.kda || 0,
        kp: (match.kp || 0) * 100,
        cs: match.cs_per_min || 0,
        gpm: match.gpm || 0,
        date: match.date ? dayjs(match.date).format('DD/MM') : `Game ${index + 1}`,
      }));

    // Distribution par résultat
    const resultData = [
      { name: 'Victoires', value: wins, color: '#0596AA' },
      { name: 'Défaites', value: totalGames - wins, color: '#C89B3C' },
    ];

    // Matchs récents (derniers 10)
    const recentMatches = championMatches
      .slice(-10)
      .reverse()
      .map((match, index) => ({
        id: index,
        date: match.date ? dayjs(match.date).format('DD/MM/YYYY HH:mm') : 'N/A',
        result: match.win ? 'Victoire' : 'Défaite',
        kda: `${match.kills || 0}/${match.deaths || 0}/${match.assists || 0}`,
        kdaRatio: match.kda || 0,
        cs: match.cs || 0,
        csPerMin: match.cs_per_min || 0,
        gpm: match.gpm || 0,
        dpm: match.dpm || 0,
        visionScore: match.vision_score || 0,
        duration: match.duration_s ? `${Math.floor(match.duration_s / 60)}:${String(match.duration_s % 60).padStart(2, '0')}` : 'N/A',
      }));

    return {
      totalGames,
      wins,
      losses: totalGames - wins,
      winrate,
      avgKda,
      avgKp,
      avgCs,
      avgGpm,
      avgDpm,
      avgVision,
      timelineData,
      resultData,
      recentMatches,
    };
  }, [data, selectedChampion]);

  if (loading) {
    return (
      <Box>
        <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
          {selectedChampion} {selectedRole && `(${selectedRole})`}
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

  if (!selectedChampion || !championData) {
    return (
      <Box>
        <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
          Détails du champion
        </Typography>
        <Alert severity="info">
          Aucun champion sélectionné ou aucune donnée disponible
        </Alert>
      </Box>
    );
  }

  return (
    <Box>
      {/* En-tête avec icône du champion */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 3, mb: 4 }}>
        <Avatar
          src={getChampionIconUrl(selectedChampion)}
          alt={selectedChampion}
          sx={{ 
            width: 80, 
            height: 80, 
            border: '3px solid',
            borderColor: 'primary.main',
          }}
        >
          {selectedChampion.slice(0, 2).toUpperCase()}
        </Avatar>
        <Box>
          <Typography variant="h3" sx={{ fontWeight: 700, mb: 1 }}>
            {selectedChampion}
          </Typography>
          <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
            {selectedRole && (
              <Chip 
                label={selectedRole} 
                color="primary" 
                variant="outlined"
              />
            )}
            <Chip 
              label={`${championData.totalGames} matchs`} 
              variant="outlined"
            />
            <Chip 
              label={`${championData.winrate.toFixed(1)}% WR`} 
              color={championData.winrate >= 50 ? 'success' : 'error'}
            />
          </Box>
        </Box>
      </Box>

      {/* Statistiques principales */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={6} sm={3}>
          <Card>
            <CardContent sx={{ textAlign: 'center' }}>
              <Typography variant="h4" color="primary" sx={{ fontWeight: 700 }}>
                {championData.totalGames}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Matchs joués
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6} sm={3}>
          <Card>
            <CardContent sx={{ textAlign: 'center' }}>
              <Typography 
                variant="h4" 
                sx={{ 
                  fontWeight: 700,
                  color: championData.winrate >= 50 ? 'success.main' : 'error.main'
                }}
              >
                {championData.winrate.toFixed(1)}%
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Taux de victoire
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6} sm={3}>
          <Card>
            <CardContent sx={{ textAlign: 'center' }}>
              <Typography variant="h4" sx={{ fontWeight: 700 }}>
                {championData.avgKda.toFixed(2)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                KDA moyen
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6} sm={3}>
          <Card>
            <CardContent sx={{ textAlign: 'center' }}>
              <Typography variant="h4" sx={{ fontWeight: 700 }}>
                {(championData.avgKp * 100).toFixed(1)}%
              </Typography>
              <Typography variant="body2" color="text.secondary">
                KP moyen
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Graphiques d'analyse */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {/* Évolution performance */}
        <Grid item xs={12} lg={8}>
          <Card>
            <CardHeader title="Évolution des performances (20 derniers matchs)" />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={championData.timelineData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Line 
                    type="monotone" 
                    dataKey="kda" 
                    stroke="#C89B3C" 
                    strokeWidth={2}
                    name="KDA"
                  />
                  <Line 
                    type="monotone" 
                    dataKey="cs" 
                    stroke="#0596AA" 
                    strokeWidth={2}
                    name="CS/min"
                  />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Répartition victoires/défaites */}
        <Grid item xs={12} lg={4}>
          <Card>
            <CardHeader title="Répartition des résultats" />
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={championData.resultData}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    label={({ name, value }) => `${name}: ${value}`}
                  >
                    {championData.resultData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Timeline de progression */}
      <Box sx={{ mb: 4 }}>
        <ProgressionTimeline
          matches={data.filter(match => match.champion === selectedChampion)}
          title={`Progression avec ${selectedChampion}`}
          showControls={true}
        />
      </Box>

      {/* Tableau des matchs récents */}
      <Card>
        <CardHeader title="Matchs récents" />
        <CardContent sx={{ p: 0 }}>
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Date</TableCell>
                  <TableCell>Résultat</TableCell>
                  <TableCell>KDA</TableCell>
                  <TableCell>Ratio</TableCell>
                  <TableCell>CS</TableCell>
                  <TableCell>CS/min</TableCell>
                  <TableCell>GPM</TableCell>
                  <TableCell>DPM</TableCell>
                  <TableCell>Vision</TableCell>
                  <TableCell>Durée</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {championData.recentMatches.map((match) => (
                  <TableRow key={match.id}>
                    <TableCell>{match.date}</TableCell>
                    <TableCell>
                      <Chip
                        label={match.result}
                        color={match.result === 'Victoire' ? 'success' : 'error'}
                        size="small"
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell>{match.kda}</TableCell>
                    <TableCell>{match.kdaRatio.toFixed(2)}</TableCell>
                    <TableCell>{match.cs}</TableCell>
                    <TableCell>{match.csPerMin.toFixed(1)}</TableCell>
                    <TableCell>{match.gpm.toFixed(0)}</TableCell>
                    <TableCell>{match.dpm.toFixed(0)}</TableCell>
                    <TableCell>{match.visionScore}</TableCell>
                    <TableCell>{match.duration}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    </Box>
  );
};