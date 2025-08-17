import React from 'react';
import { 
  Grid, 
  Card, 
  CardContent, 
  Typography, 
  Box,
  Skeleton,
  Alert,
  Chip,
} from '@mui/material';
import { Row } from '../types';

interface OverviewProps {
  allData: Row[];
  data: Row[];
  loading: boolean;
  error: string | null;
  onRoleSelect: (role: string) => void;
}

export const Overview: React.FC<OverviewProps> = ({
  allData,
  loading,
  error,
  onRoleSelect,
}) => {
  if (loading) {
    return (
      <Grid container spacing={3}>
        {[...Array(6)].map((_, i) => (
          <Grid item xs={12} sm={6} md={4} key={i}>
            <Card>
              <CardContent>
                <Skeleton variant="text" width="60%" height={32} />
                <Skeleton variant="text" width="40%" height={24} sx={{ mt: 1 }} />
                <Skeleton variant="rectangular" height={60} sx={{ mt: 2 }} />
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 3 }}>
        {error}
      </Alert>
    );
  }

  if (allData.length === 0) {
    return (
      <Card>
        <CardContent sx={{ textAlign: 'center', py: 6 }}>
          <Typography variant="h6" color="text.secondary">
            Aucune donnée disponible
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Lancez un export pour commencer l'analyse de vos performances
          </Typography>
        </CardContent>
      </Card>
    );
  }

  // Calculs de base pour la vue d'ensemble
  const totalMatches = allData.length;
  const wins = allData.filter(row => row.win).length;
  const winrate = totalMatches > 0 ? (wins / totalMatches) * 100 : 0;

  // Répartition par rôle
  const roleStats = allData.reduce((acc, row) => {
    const role = row.lane || 'Unknown';
    if (!acc[role]) {
      acc[role] = { matches: 0, wins: 0 };
    }
    acc[role].matches++;
    if (row.win) acc[role].wins++;
    return acc;
  }, {} as Record<string, { matches: number; wins: number }>);

  // Champions les plus joués
  const championStats = allData.reduce((acc, row) => {
    const champion = row.champion || 'Unknown';
    if (!acc[champion]) {
      acc[champion] = { matches: 0, wins: 0 };
    }
    acc[champion].matches++;
    if (row.win) acc[champion].wins++;
    return acc;
  }, {} as Record<string, { matches: number; wins: number }>);

  const topChampions = Object.entries(championStats)
    .sort(([,a], [,b]) => b.matches - a.matches)
    .slice(0, 5);

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
        Vue d'ensemble
      </Typography>

      <Grid container spacing={3}>
        {/* Statistiques globales */}
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                Matches totaux
              </Typography>
              <Typography variant="h3" sx={{ fontWeight: 700 }}>
                {totalMatches}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                Taux de victoire
              </Typography>
              <Typography 
                variant="h3" 
                sx={{ 
                  fontWeight: 700,
                  color: winrate >= 50 ? 'success.main' : 'error.main'
                }}
              >
                {winrate.toFixed(1)}%
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                Victoires
              </Typography>
              <Typography variant="h3" sx={{ fontWeight: 700, color: 'success.main' }}>
                {wins}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                Défaites
              </Typography>
              <Typography variant="h3" sx={{ fontWeight: 700, color: 'error.main' }}>
                {totalMatches - wins}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Répartition par rôle */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 2 }}>
                Performance par rôle
              </Typography>
              {Object.entries(roleStats).map(([role, stats]) => {
                const roleWinrate = (stats.wins / stats.matches) * 100;
                return (
                  <Box 
                    key={role}
                    sx={{ 
                      display: 'flex', 
                      justifyContent: 'space-between', 
                      alignItems: 'center',
                      mb: 1,
                      p: 1,
                      borderRadius: 1,
                      cursor: 'pointer',
                      '&:hover': {
                        backgroundColor: 'action.hover',
                      }
                    }}
                    onClick={() => onRoleSelect(role)}
                  >
                    <Typography variant="body1" sx={{ fontWeight: 500 }}>
                      {role}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Chip 
                        label={`${stats.matches} games`}
                        size="small"
                        variant="outlined"
                      />
                      <Chip 
                        label={`${roleWinrate.toFixed(1)}%`}
                        size="small"
                        color={roleWinrate >= 50 ? 'success' : 'error'}
                      />
                    </Box>
                  </Box>
                );
              })}
            </CardContent>
          </Card>
        </Grid>

        {/* Top champions */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 2 }}>
                Champions les plus joués
              </Typography>
              {topChampions.map(([champion, stats]) => {
                const championWinrate = (stats.wins / stats.matches) * 100;
                return (
                  <Box 
                    key={champion}
                    sx={{ 
                      display: 'flex', 
                      justifyContent: 'space-between', 
                      alignItems: 'center',
                      mb: 1,
                      p: 1,
                      borderRadius: 1,
                    }}
                  >
                    <Typography variant="body1" sx={{ fontWeight: 500 }}>
                      {champion}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Chip 
                        label={`${stats.matches} games`}
                        size="small"
                        variant="outlined"
                      />
                      <Chip 
                        label={`${championWinrate.toFixed(1)}%`}
                        size="small"
                        color={championWinrate >= 50 ? 'success' : 'error'}
                      />
                    </Box>
                  </Box>
                );
              })}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};