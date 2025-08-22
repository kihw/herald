import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Chip,
  Avatar,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  LinearProgress,
  useTheme,
  Alert,
  CircularProgress,
  Divider,
  Paper,
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  EmojiEvents as EmojiEventsIcon,
  Person as PersonIcon,
  Group as GroupIcon,
  Analytics as AnalyticsIcon,
  SportsEsports as SportsEsportsIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import { groupApi, GroupStats } from '../../services/groupApi';
import GroupStatsCharts from '../charts/GroupStatsCharts';

interface GroupStatsOverviewProps {
  groupId: number;
}

const GroupStatsOverview: React.FC<GroupStatsOverviewProps> = ({ groupId }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [stats, setStats] = useState<GroupStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (groupId) {
      loadGroupStats();
    }
  }, [groupId]);

  const loadGroupStats = async () => {
    try {
      setLoading(true);
      const groupStats = await groupApi.getGroupStats(groupId);
      setStats(groupStats);
      setError(null);
    } catch (err) {
      setError('Erreur lors du chargement des statistiques');
      console.error('Error loading group stats:', err);
    } finally {
      setLoading(false);
    }
  };

  const getRankColor = (rank: string) => {
    const rankLower = rank.toLowerCase();
    if (rankLower.includes('challenger')) return leagueColors.gold[500];
    if (rankLower.includes('grandmaster')) return leagueColors.gold[400];
    if (rankLower.includes('master')) return leagueColors.gold[300];
    if (rankLower.includes('diamond')) return leagueColors.blue[400];
    if (rankLower.includes('platinum')) return leagueColors.blue[300];
    if (rankLower.includes('gold')) return leagueColors.gold[200];
    if (rankLower.includes('silver')) return '#C0C0C0';
    if (rankLower.includes('bronze')) return '#CD7F32';
    if (rankLower.includes('iron')) return '#8B4513';
    return 'text.secondary';
  };

  if (loading) {
    return (
      <Box 
        sx={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          minHeight: 300 
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ m: 3 }}>
        {error}
      </Alert>
    );
  }

  if (!stats) {
    return (
      <Box sx={{ textAlign: 'center', py: 6 }}>
        <AnalyticsIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
        <Typography variant="h6" color="text.secondary" gutterBottom>
          Statistiques non disponibles
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Les statistiques du groupe seront calculées automatiquement
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      {/* Overview Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.blue[500]}20 0%, ${leagueColors.blue[500]}10 100%)`,
              border: `1px solid ${leagueColors.blue[500]}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <GroupIcon sx={{ fontSize: 40, color: leagueColors.blue[500], mb: 1 }} />
              <Typography variant="h4" sx={{ fontWeight: 700 }}>
                {stats.total_members}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Membres Total
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.win}20 0%, ${leagueColors.win}10 100%)`,
              border: `1px solid ${leagueColors.win}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <PersonIcon sx={{ fontSize: 40, color: leagueColors.win, mb: 1 }} />
              <Typography variant="h4" sx={{ fontWeight: 700 }}>
                {stats.active_members}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Membres Actifs
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.gold[500]}20 0%, ${leagueColors.gold[500]}10 100%)`,
              border: `1px solid ${leagueColors.gold[500]}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <EmojiEventsIcon sx={{ fontSize: 40, color: leagueColors.gold[500], mb: 1 }} />
              <Typography variant="h4" sx={{ fontWeight: 700, fontSize: '1.5rem' }}>
                {stats.average_rank || 'N/A'}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Rang Moyen
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.dark[400]}20 0%, ${leagueColors.dark[400]}10 100%)`,
              border: `1px solid ${leagueColors.dark[400]}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <TrendingUpIcon sx={{ fontSize: 40, color: leagueColors.dark[400], mb: 1 }} />
              <Typography variant="h4" sx={{ fontWeight: 700 }}>
                {stats.average_mmr ? Math.round(stats.average_mmr) : 'N/A'}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                MMR Moyen
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Grid container spacing={3}>
        {/* Top Champions */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <EmojiEventsIcon color="primary" />
                Champions Populaires
              </Typography>
              
              {stats.top_champions && stats.top_champions.length > 0 ? (
                <List>
                  {stats.top_champions.slice(0, 5).map((champion, index) => (
                    <React.Fragment key={champion.champion_id}>
                      <ListItem sx={{ px: 0 }}>
                        <ListItemAvatar>
                          <Avatar 
                            sx={{ 
                              background: `linear-gradient(135deg, ${leagueColors.gold[500]} 0%, ${leagueColors.gold[600]} 100%)`,
                              color: '#000',
                              fontWeight: 600,
                              fontSize: '0.8rem',
                            }}
                          >
                            #{index + 1}
                          </Avatar>
                        </ListItemAvatar>
                        <ListItemText
                          primary={champion.champion_name}
                          secondary={
                            <Box sx={{ mt: 0.5 }}>
                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
                                <Chip
                                  label={`${champion.play_count} parties`}
                                  size="small"
                                  variant="outlined"
                                />
                                <Chip
                                  label={`${(champion.win_rate * 100).toFixed(1)}% WR`}
                                  size="small"
                                  sx={{
                                    background: `${leagueColors.win}20`,
                                    color: leagueColors.win,
                                    border: 'none',
                                  }}
                                />
                                <Chip
                                  label={`${champion.avg_kda.toFixed(1)} KDA`}
                                  size="small"
                                  variant="outlined"
                                />
                              </Box>
                              <LinearProgress
                                variant="determinate"
                                value={champion.win_rate * 100}
                                sx={{
                                  height: 4,
                                  borderRadius: 2,
                                  backgroundColor: `${leagueColors.loss}30`,
                                  '& .MuiLinearProgress-bar': {
                                    backgroundColor: leagueColors.win,
                                    borderRadius: 2,
                                  },
                                }}
                              />
                            </Box>
                          }
                        />
                      </ListItem>
                      {index < Math.min(stats.top_champions.length, 5) - 1 && <Divider variant="inset" component="li" />}
                    </React.Fragment>
                  ))}
                </List>
              ) : (
                <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 3 }}>
                  Aucune donnée de champion disponible
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Popular Roles */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <SportsEsportsIcon color="primary" />
                Rôles Favoris
              </Typography>
              
              {stats.popular_roles && stats.popular_roles.length > 0 ? (
                <List>
                  {stats.popular_roles.map((role, index) => (
                    <React.Fragment key={role.role}>
                      <ListItem sx={{ px: 0 }}>
                        <ListItemAvatar>
                          <Avatar 
                            sx={{ 
                              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
                              color: '#fff',
                              fontWeight: 600,
                              fontSize: '0.8rem',
                            }}
                          >
                            {role.role.charAt(0)}
                          </Avatar>
                        </ListItemAvatar>
                        <ListItemText
                          primary={
                            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                              <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                {role.role}
                              </Typography>
                              <Typography variant="body2" color="text.secondary">
                                {role.play_count} parties
                              </Typography>
                            </Box>
                          }
                          secondary={
                            <Box sx={{ mt: 0.5 }}>
                              <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                                <Typography variant="body2" color="text.secondary">
                                  Taux de victoire
                                </Typography>
                                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                  {(role.win_rate * 100).toFixed(1)}%
                                </Typography>
                              </Box>
                              <LinearProgress
                                variant="determinate"
                                value={role.win_rate * 100}
                                sx={{
                                  height: 4,
                                  borderRadius: 2,
                                  backgroundColor: `${leagueColors.loss}30`,
                                  '& .MuiLinearProgress-bar': {
                                    backgroundColor: role.win_rate >= 0.5 ? leagueColors.win : leagueColors.loss,
                                    borderRadius: 2,
                                  },
                                }}
                              />
                            </Box>
                          }
                        />
                      </ListItem>
                      {index < stats.popular_roles.length - 1 && <Divider variant="inset" component="li" />}
                    </React.Fragment>
                  ))}
                </List>
              ) : (
                <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 3 }}>
                  Aucune donnée de rôle disponible
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Winrate Comparison */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <TrendingUpIcon color="primary" />
                Comparaison des Taux de Victoire
              </Typography>
              
              {stats.winrate_comparison && Object.keys(stats.winrate_comparison).length > 0 ? (
                <Grid container spacing={2}>
                  {Object.entries(stats.winrate_comparison).map(([memberName, winrate]) => (
                    <Grid item xs={12} sm={6} md={4} key={memberName}>
                      <Paper
                        sx={{
                          p: 2,
                          textAlign: 'center',
                          background: isDarkMode
                            ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
                            : `linear-gradient(135deg, ${leagueColors.blue[25]} 0%, #ffffff 100%)`,
                          border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                        }}
                      >
                        <Avatar
                          sx={{
                            width: 48,
                            height: 48,
                            margin: '0 auto 8px auto',
                            background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
                            color: '#fff',
                            fontWeight: 600,
                          }}
                        >
                          {memberName.charAt(0).toUpperCase()}
                        </Avatar>
                        <Typography variant="body1" sx={{ fontWeight: 500, mb: 1 }}>
                          {memberName}
                        </Typography>
                        <Typography 
                          variant="h6" 
                          sx={{ 
                            fontWeight: 700,
                            color: winrate >= 0.5 ? leagueColors.win : leagueColors.loss,
                          }}
                        >
                          {(winrate * 100).toFixed(1)}%
                        </Typography>
                        <LinearProgress
                          variant="determinate"
                          value={winrate * 100}
                          sx={{
                            mt: 1,
                            height: 6,
                            borderRadius: 3,
                            backgroundColor: `${leagueColors.loss}30`,
                            '& .MuiLinearProgress-bar': {
                              backgroundColor: winrate >= 0.5 ? leagueColors.win : leagueColors.loss,
                              borderRadius: 3,
                            },
                          }}
                        />
                      </Paper>
                    </Grid>
                  ))}
                </Grid>
              ) : (
                <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 3 }}>
                  Aucune donnée de comparaison disponible
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Charts Section */}
        <Grid item xs={12}>
          <GroupStatsCharts stats={stats} />
        </Grid>

        {/* Last Updated */}
        <Grid item xs={12}>
          <Box sx={{ textAlign: 'center', py: 2 }}>
            <Typography variant="caption" color="text.secondary">
              Dernière mise à jour: {new Date(stats.last_updated).toLocaleString()}
            </Typography>
          </Box>
        </Grid>
      </Grid>
    </Box>
  );
};

export default GroupStatsOverview;