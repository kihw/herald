import React, { useState, useEffect } from 'react';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  useTheme,
  Tab,
  Tabs,
  Fade,
  Chip,
  Avatar,
  AvatarGroup,
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  EmojiEvents as EmojiEventsIcon,
  Group as GroupIcon,
  Analytics as AnalyticsIcon,
  Timer as TimerIcon,
  Visibility as VisibilityIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import { Group, GroupStats } from '../../services/groupApi';
import PlayerPerformanceWidget from '../charts/PlayerPerformanceWidget';
import GroupStatsCharts from '../charts/GroupStatsCharts';

interface EnhancedGroupOverviewProps {
  group: Group;
  stats?: GroupStats;
}

const EnhancedGroupOverview: React.FC<EnhancedGroupOverviewProps> = ({ group, stats }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [activeTab, setActiveTab] = useState(0);
  const [liveData, setLiveData] = useState({
    activeNow: Math.floor(Math.random() * 3) + 1,
    gamesInProgress: Math.floor(Math.random() * 2),
    recentActivity: Math.floor(Math.random() * 50) + 10,
    averageGameTime: 28 + Math.floor(Math.random() * 10),
  });

  // Mock player data for demonstration
  const mockPlayers = Array.from({ length: Math.min(group.member_count, 6) }, (_, i) => ({
    riot_id: `Player${i + 1}`,
    riot_tag: `TAG${i + 1}`,
    rank: ['Bronze II', 'Silver III', 'Gold I', 'Platinum IV', 'Diamond III'][i % 5],
    lp: 45 + (i * 23) % 100,
    winrate: 0.45 + (Math.random() * 0.3),
    kda: 1.5 + (Math.random() * 1.5),
    recent_games: 8 + Math.floor(Math.random() * 12),
    trend: ['up', 'down', 'stable'][i % 3] as 'up' | 'down' | 'stable',
    performance_history: Array.from({ length: 10 }, () => Math.random() > 0.5 ? 1 : 0),
  }));

  useEffect(() => {
    // Simulate live data updates
    const interval = setInterval(() => {
      setLiveData(prev => ({
        ...prev,
        activeNow: Math.max(0, prev.activeNow + (Math.random() > 0.7 ? (Math.random() > 0.5 ? 1 : -1) : 0)),
        gamesInProgress: Math.max(0, prev.gamesInProgress + (Math.random() > 0.8 ? (Math.random() > 0.5 ? 1 : -1) : 0)),
        recentActivity: prev.recentActivity + Math.floor(Math.random() * 3) - 1,
      }));
    }, 5000);

    return () => clearInterval(interval);
  }, []);

  const getPrivacyIcon = (privacy: string) => {
    switch (privacy) {
      case 'public':
        return <VisibilityIcon sx={{ fontSize: 16, color: leagueColors.blue[500] }} />;
      case 'private':
        return <GroupIcon sx={{ fontSize: 16, color: leagueColors.gold[500] }} />;
      default:
        return <EmojiEventsIcon sx={{ fontSize: 16, color: leagueColors.dark[400] }} />;
    }
  };

  const renderOverviewTab = () => (
    <Fade in={activeTab === 0}>
      <Grid container spacing={3}>
        {/* Live Metrics */}
        <Grid item xs={12}>
          <Card
            sx={{
              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
              color: '#fff',
              border: `1px solid ${leagueColors.blue[400]}`,
            }}
          >
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <TimerIcon />
                Activité en Temps Réel
              </Typography>
              <Grid container spacing={3}>
                <Grid item xs={6} sm={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h3" sx={{ fontWeight: 700, mb: 0.5 }}>
                      {liveData.activeNow}
                    </Typography>
                    <Typography variant="body2" sx={{ opacity: 0.9 }}>
                      En ligne maintenant
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h3" sx={{ fontWeight: 700, mb: 0.5 }}>
                      {liveData.gamesInProgress}
                    </Typography>
                    <Typography variant="body2" sx={{ opacity: 0.9 }}>
                      Parties en cours
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h3" sx={{ fontWeight: 700, mb: 0.5 }}>
                      {liveData.recentActivity}
                    </Typography>
                    <Typography variant="body2" sx={{ opacity: 0.9 }}>
                      Parties (7j)
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h3" sx={{ fontWeight: 700, mb: 0.5 }}>
                      {liveData.averageGameTime}m
                    </Typography>
                    <Typography variant="body2" sx={{ opacity: 0.9 }}>
                      Durée moyenne
                    </Typography>
                  </Box>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Group Info */}
        <Grid item xs={12} md={4}>
          <Card
            sx={{
              height: '100%',
              background: isDarkMode
                ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
                : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
              border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
            }}
          >
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Informations du Groupe
              </Typography>
              
              <Box sx={{ mb: 3 }}>
                <Typography variant="h4" sx={{ fontWeight: 700, mb: 1 }}>
                  {group.name}
                </Typography>
                {group.description && (
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    {group.description}
                  </Typography>
                )}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap' }}>
                  <Chip
                    icon={getPrivacyIcon(group.privacy)}
                    label={
                      group.privacy === 'public' ? 'Public' : 
                      group.privacy === 'private' ? 'Privé' : 'Sur invitation'
                    }
                    size="small"
                    variant="outlined"
                  />
                  <Chip
                    label={`${group.member_count} membre${group.member_count > 1 ? 's' : ''}`}
                    size="small"
                    variant="outlined"
                  />
                </Box>
              </Box>

              <Box>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  Propriétaire
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Avatar sx={{ width: 32, height: 32 }}>
                    {group.owner.riot_id.charAt(0).toUpperCase()}
                  </Avatar>
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {group.owner.riot_id}#{group.owner.riot_tag}
                  </Typography>
                </Box>
              </Box>

              <Box sx={{ mt: 3 }}>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  Créé le
                </Typography>
                <Typography variant="body1">
                  {new Date(group.created_at).toLocaleDateString()}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Quick Stats */}
        <Grid item xs={12} md={8}>
          <Grid container spacing={2} sx={{ height: '100%' }}>
            <Grid item xs={6} sm={3}>
              <Card
                sx={{
                  height: '100%',
                  textAlign: 'center',
                  background: `linear-gradient(135deg, ${leagueColors.win}20 0%, ${leagueColors.win}10 100%)`,
                  border: `1px solid ${leagueColors.win}30`,
                }}
              >
                <CardContent>
                  <EmojiEventsIcon sx={{ fontSize: 32, color: leagueColors.win, mb: 1 }} />
                  <Typography variant="h4" sx={{ fontWeight: 700 }}>
                    {stats?.average_rank || 'Silver'}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Rang Moyen
                  </Typography>
                </CardContent>
              </Card>
            </Grid>

            <Grid item xs={6} sm={3}>
              <Card
                sx={{
                  height: '100%',
                  textAlign: 'center',
                  background: `linear-gradient(135deg, ${leagueColors.gold[500]}20 0%, ${leagueColors.gold[500]}10 100%)`,
                  border: `1px solid ${leagueColors.gold[500]}30`,
                }}
              >
                <CardContent>
                  <TrendingUpIcon sx={{ fontSize: 32, color: leagueColors.gold[500], mb: 1 }} />
                  <Typography variant="h4" sx={{ fontWeight: 700 }}>
                    {stats?.average_mmr ? Math.round(stats.average_mmr) : 1450}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    MMR Moyen
                  </Typography>
                </CardContent>
              </Card>
            </Grid>

            <Grid item xs={6} sm={3}>
              <Card
                sx={{
                  height: '100%',
                  textAlign: 'center',
                  background: `linear-gradient(135deg, ${leagueColors.blue[500]}20 0%, ${leagueColors.blue[500]}10 100%)`,
                  border: `1px solid ${leagueColors.blue[500]}30`,
                }}
              >
                <CardContent>
                  <GroupIcon sx={{ fontSize: 32, color: leagueColors.blue[500], mb: 1 }} />
                  <Typography variant="h4" sx={{ fontWeight: 700 }}>
                    {stats?.active_members || Math.floor(group.member_count * 0.8)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Actifs
                  </Typography>
                </CardContent>
              </Card>
            </Grid>

            <Grid item xs={6} sm={3}>
              <Card
                sx={{
                  height: '100%',
                  textAlign: 'center',
                  background: `linear-gradient(135deg, ${leagueColors.dark[400]}20 0%, ${leagueColors.dark[400]}10 100%)`,
                  border: `1px solid ${leagueColors.dark[400]}30`,
                }}
              >
                <CardContent>
                  <AnalyticsIcon sx={{ fontSize: 32, color: leagueColors.dark[400], mb: 1 }} />
                  <Typography variant="h4" sx={{ fontWeight: 700 }}>
                    {Math.floor(Math.random() * 200) + 100}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Parties Total
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Grid>

        {/* Recent Members */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <GroupIcon />
                Membres Récents
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <AvatarGroup max={8}>
                  {mockPlayers.map((player, index) => (
                    <Avatar 
                      key={index}
                      sx={{ 
                        width: 40, 
                        height: 40,
                        background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
                        color: '#fff',
                        fontWeight: 600,
                      }}
                    >
                      {player.riot_id.charAt(0)}
                    </Avatar>
                  ))}
                </AvatarGroup>
                <Typography variant="body2" color="text.secondary">
                  et {Math.max(0, group.member_count - 8)} autres membres
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Fade>
  );

  const renderPlayersTab = () => (
    <Fade in={activeTab === 1}>
      <Grid container spacing={3}>
        {mockPlayers.map((player, index) => (
          <Grid item xs={12} md={6} lg={4} key={index}>
            <PlayerPerformanceWidget player={player} />
          </Grid>
        ))}
      </Grid>
    </Fade>
  );

  const renderChartsTab = () => (
    <Fade in={activeTab === 2}>
      <Box>
        {stats ? (
          <GroupStatsCharts stats={stats} />
        ) : (
          <Box sx={{ textAlign: 'center', py: 6 }}>
            <AnalyticsIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" gutterBottom>
              Chargement des statistiques...
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Les graphiques seront disponibles dans quelques instants
            </Typography>
          </Box>
        )}
      </Box>
    </Fade>
  );

  return (
    <Box sx={{ p: 3 }}>
      {/* Tabs Navigation */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs
          value={activeTab}
          onChange={(e, newValue) => setActiveTab(newValue)}
          sx={{
            '& .MuiTab-root': {
              fontWeight: 500,
            },
          }}
        >
          <Tab 
            icon={<GroupIcon />} 
            label="Vue d'ensemble" 
            iconPosition="start"
          />
          <Tab 
            icon={<EmojiEventsIcon />} 
            label="Performances" 
            iconPosition="start"
          />
          <Tab 
            icon={<AnalyticsIcon />} 
            label="Graphiques" 
            iconPosition="start"
          />
        </Tabs>
      </Box>

      {/* Tab Content */}
      <Box>
        {activeTab === 0 && renderOverviewTab()}
        {activeTab === 1 && renderPlayersTab()}
        {activeTab === 2 && renderChartsTab()}
      </Box>
    </Box>
  );
};

export default EnhancedGroupOverview;