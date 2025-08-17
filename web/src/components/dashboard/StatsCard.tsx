import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Grid,
  Box,
  Avatar,
  Chip,
  LinearProgress,
} from '@mui/material';
import {
  TrendingUp,
  Person,
  EmojiEvents,
  Timeline,
} from '@mui/icons-material';

interface StatsData {
  totalMatches: number;
  winRate: number;
  mainRole: string;
  favoriteChampion: {
    name: string;
    winRate: number;
    matches: number;
  };
  recentPerformance: {
    last7Days: {
      matches: number;
      wins: number;
      winRate: number;
    };
    last30Days: {
      matches: number;
      wins: number;
      winRate: number;
    };
  };
  rankInfo: {
    tier: string;
    division: string;
    lp: number;
  };
}

interface StatsCardProps {
  stats: StatsData;
  loading?: boolean;
}

const StatsCard: React.FC<StatsCardProps> = ({ stats, loading = false }) => {
  if (loading) {
    return (
      <Grid container spacing={2}>
        {[1, 2, 3, 4].map((i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <Card>
              <CardContent>
                <Box sx={{ height: 100 }}>
                  <LinearProgress />
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    );
  }

  return (
    <Grid container spacing={2}>
      {/* Total Matches */}
      <Grid item xs={12} sm={6} md={3}>
        <Card sx={{ height: '100%' }}>
          <CardContent>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <Avatar sx={{ bgcolor: 'primary.main', mr: 2 }}>
                <Timeline />
              </Avatar>
              <Typography variant="h6" component="div">
                Total Matches
              </Typography>
            </Box>
            <Typography variant="h4" component="div" sx={{ mb: 1 }}>
              {stats.totalMatches}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              All time matches played
            </Typography>
          </CardContent>
        </Card>
      </Grid>

      {/* Win Rate */}
      <Grid item xs={12} sm={6} md={3}>
        <Card sx={{ height: '100%' }}>
          <CardContent>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <Avatar sx={{ bgcolor: 'success.main', mr: 2 }}>
                <TrendingUp />
              </Avatar>
              <Typography variant="h6" component="div">
                Win Rate
              </Typography>
            </Box>
            <Typography variant="h4" component="div" sx={{ mb: 1 }}>
              {stats.winRate}%
            </Typography>
            <LinearProgress
              variant="determinate"
              value={stats.winRate}
              sx={{ mt: 1 }}
              color={stats.winRate >= 60 ? 'success' : stats.winRate >= 50 ? 'warning' : 'error'}
            />
          </CardContent>
        </Card>
      </Grid>

      {/* Main Role */}
      <Grid item xs={12} sm={6} md={3}>
        <Card sx={{ height: '100%' }}>
          <CardContent>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <Avatar sx={{ bgcolor: 'info.main', mr: 2 }}>
                <Person />
              </Avatar>
              <Typography variant="h6" component="div">
                Main Role
              </Typography>
            </Box>
            <Chip
              label={stats.mainRole}
              color="primary"
              sx={{ mb: 2, fontSize: '1.2rem', p: 2 }}
            />
            <Typography variant="body2" color="text.secondary">
              Most played position
            </Typography>
          </CardContent>
        </Card>
      </Grid>

      {/* Rank */}
      <Grid item xs={12} sm={6} md={3}>
        <Card sx={{ height: '100%' }}>
          <CardContent>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <Avatar sx={{ bgcolor: 'warning.main', mr: 2 }}>
                <EmojiEvents />
              </Avatar>
              <Typography variant="h6" component="div">
                Current Rank
              </Typography>
            </Box>
            <Typography variant="h4" component="div" sx={{ mb: 1 }}>
              {stats.rankInfo.tier} {stats.rankInfo.division}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {stats.rankInfo.lp} LP
            </Typography>
          </CardContent>
        </Card>
      </Grid>

      {/* Favorite Champion */}
      <Grid item xs={12} sm={6}>
        <Card sx={{ height: '100%' }}>
          <CardContent>
            <Typography variant="h6" component="div" sx={{ mb: 2 }}>
              Favorite Champion
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <Avatar
                sx={{ width: 56, height: 56, mr: 2 }}
                src={`https://ddragon.leagueoflegends.com/cdn/14.17.1/img/champion/${stats.favoriteChampion.name}.png`}
              >
                {stats.favoriteChampion.name[0]}
              </Avatar>
              <Box>
                <Typography variant="h5" component="div">
                  {stats.favoriteChampion.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {stats.favoriteChampion.matches} matches played
                </Typography>
              </Box>
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <Typography variant="body1" sx={{ mr: 2 }}>
                Win Rate: {stats.favoriteChampion.winRate}%
              </Typography>
              <LinearProgress
                variant="determinate"
                value={stats.favoriteChampion.winRate}
                sx={{ flexGrow: 1 }}
                color={stats.favoriteChampion.winRate >= 60 ? 'success' : 'warning'}
              />
            </Box>
          </CardContent>
        </Card>
      </Grid>

      {/* Recent Performance */}
      <Grid item xs={12} sm={6}>
        <Card sx={{ height: '100%' }}>
          <CardContent>
            <Typography variant="h6" component="div" sx={{ mb: 3 }}>
              Recent Performance
            </Typography>
            
            <Box sx={{ mb: 3 }}>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>
                Last 7 Days
              </Typography>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                <Typography variant="body2">
                  {stats.recentPerformance.last7Days.wins}W / {stats.recentPerformance.last7Days.matches - stats.recentPerformance.last7Days.wins}L
                </Typography>
                <Typography variant="body2" color="primary">
                  {stats.recentPerformance.last7Days.winRate}%
                </Typography>
              </Box>
              <LinearProgress
                variant="determinate"
                value={stats.recentPerformance.last7Days.winRate}
                color={stats.recentPerformance.last7Days.winRate >= 60 ? 'success' : 'warning'}
              />
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>
                Last 30 Days
              </Typography>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                <Typography variant="body2">
                  {stats.recentPerformance.last30Days.wins}W / {stats.recentPerformance.last30Days.matches - stats.recentPerformance.last30Days.wins}L
                </Typography>
                <Typography variant="body2" color="primary">
                  {stats.recentPerformance.last30Days.winRate}%
                </Typography>
              </Box>
              <LinearProgress
                variant="determinate"
                value={stats.recentPerformance.last30Days.winRate}
                color={stats.recentPerformance.last30Days.winRate >= 60 ? 'success' : 'warning'}
              />
            </Box>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
};

export default StatsCard;
