import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Alert,
  CircularProgress,
  Chip,
} from '@mui/material';
import {
  Sync as SyncIcon,
  Person as PersonIcon,
  TrendingUp as TrendingUpIcon,
  SportsMma as SportsIcon,
  Schedule as ScheduleIcon,
} from '@mui/icons-material';
import { useAuth } from '../../context/AuthContext';
import { apiService, DashboardStats } from '../../services/api';

export function NewDashboard() {
  const { user, isAuthenticated, logout, error: authError } = useAuth();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadStats = async () => {
    try {
      setLoading(true);
      setError(null);
      const dashboardStats = await apiService.getDashboardStats();
      setStats(dashboardStats);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Erreur lors du chargement des statistiques';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  const handleSync = async () => {
    try {
      setSyncing(true);
      setError(null);
      await apiService.syncMatches();
      // Recharger les stats après la synchronisation
      await loadStats();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Erreur lors de la synchronisation';
      setError(message);
    } finally {
      setSyncing(false);
    }
  };

  useEffect(() => {
    loadStats();
  }, []);

  const formatDate = (dateString: string | null) => {
    if (!dateString) return 'Jamais';
    return new Date(dateString).toLocaleString('fr-FR');
  };

  const formatWinRate = (winRate: number) => {
    return `${Math.round(winRate * 100)}%`;
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      {/* Header avec info utilisateur */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box display="flex" alignItems="center" gap={2}>
            <PersonIcon color="primary" />
            <Box>
              <Typography variant="h5">
                {user?.riot_id}#{user?.riot_tag}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {user?.region?.toUpperCase()}
              </Typography>
            </Box>
          </Box>
          <Box display="flex" gap={2}>
            <Button
              variant="contained"
              startIcon={syncing ? <CircularProgress size={20} /> : <SyncIcon />}
              onClick={handleSync}
              disabled={syncing}
            >
              {syncing ? 'Synchronisation...' : 'Synchroniser'}
            </Button>
            <Button variant="outlined" onClick={logout}>
              Déconnexion
            </Button>
          </Box>
        </Box>
      </Paper>

      {/* Messages d'erreur */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Statistiques principales */}
      {stats && (
        <Grid container spacing={3}>
          {/* Nombre total de matches */}
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" gap={2}>
                  <SportsIcon color="primary" />
                  <Box>
                    <Typography variant="h4" component="div">
                      {stats.total_matches}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Matches totaux
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* Taux de victoire */}
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" gap={2}>
                  <TrendingUpIcon color="success" />
                  <Box>
                    <Typography variant="h4" component="div">
                      {formatWinRate(stats.win_rate)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Taux de victoire
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* KDA moyen */}
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" gap={2}>
                  <SportsIcon color="warning" />
                  <Box>
                    <Typography variant="h4" component="div">
                      {stats.average_kda.toFixed(2)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      KDA moyen
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* Champion favori */}
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" gap={2}>
                  <PersonIcon color="secondary" />
                  <Box>
                    <Typography variant="h6" component="div">
                      {stats.favorite_champion || 'Aucun'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Champion favori
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* Informations de synchronisation */}
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Synchronisation
                </Typography>
                <Box display="flex" gap={2} flexWrap="wrap">
                  <Chip
                    icon={<ScheduleIcon />}
                    label={`Dernière sync: ${formatDate(stats.last_sync_at)}`}
                    variant="outlined"
                  />
                  <Chip
                    icon={<ScheduleIcon />}
                    label={`Prochaine sync: ${formatDate(stats.next_sync_at)}`}
                    variant="outlined"
                    color="primary"
                  />
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}

      {/* Message si pas de données */}
      {stats && stats.total_matches === 0 && (
        <Paper sx={{ p: 4, textAlign: 'center', mt: 3 }}>
          <Typography variant="h6" gutterBottom>
            Aucune donnée disponible
          </Typography>
          <Typography variant="body1" color="text.secondary" gutterBottom>
            Lancez une synchronisation pour récupérer vos matches League of Legends.
          </Typography>
          <Button
            variant="contained"
            startIcon={syncing ? <CircularProgress size={20} /> : <SyncIcon />}
            onClick={handleSync}
            disabled={syncing}
            sx={{ mt: 2 }}
          >
            {syncing ? 'Synchronisation...' : 'Première synchronisation'}
          </Button>
        </Paper>
      )}
    </Box>
  );
}
