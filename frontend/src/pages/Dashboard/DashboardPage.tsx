import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Avatar,
  Chip,
  LinearProgress,
  IconButton,
  Stack,
  Divider,
  Paper,
  Button,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  SportsSoccer,
  Timeline,
  EmojiEvents,
  Visibility,
  AttachMoney,
  PersonAdd,
  PlayArrow,
  Refresh,
} from '@mui/icons-material';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';

// Mock data - in production this would come from the backend
const mockPlayerStats = {
  summoner: {
    name: 'Herald User',
    tagLine: 'EUW',
    level: 247,
    profileIconId: 4641,
    tier: 'GOLD',
    rank: 'III',
    leaguePoints: 65,
    wins: 47,
    losses: 39,
  },
  currentMatch: null,
  recentPerformance: {
    kda: 2.4,
    kdaTrend: 'up',
    winRate: 54.7,
    winRateTrend: 'up',
    csPerMin: 6.8,
    csPerMinTrend: 'down',
    visionScore: 45.2,
    visionTrend: 'up',
    goldPerMin: 412,
    goldTrend: 'up',
  },
  matchHistory: [
    { id: 1, champion: 'Jinx', result: 'Victory', kda: '12/4/8', duration: '28:45', ago: '2 hours' },
    { id: 2, champion: 'Caitlyn', result: 'Defeat', kda: '8/6/12', duration: '35:22', ago: '4 hours' },
    { id: 3, champion: 'Vayne', result: 'Victory', kda: '15/2/6', duration: '32:18', ago: '6 hours' },
    { id: 4, champion: 'Ashe', result: 'Victory', kda: '7/3/14', duration: '29:34', ago: '1 day' },
    { id: 5, champion: 'Kai\'Sa', result: 'Defeat', kda: '9/8/5', duration: '41:12', ago: '1 day' },
  ],
  rankProgress: [
    { date: '2024-01', lp: 1245 },
    { date: '2024-02', lp: 1367 },
    { date: '2024-03', lp: 1423 },
    { date: '2024-04', lp: 1465 },
    { date: '2024-05', lp: 1398 },
    { date: '2024-06', lp: 1465 },
  ],
  championStats: [
    { name: 'Jinx', games: 15, winRate: 73, color: '#ff6b9d' },
    { name: 'Caitlyn', games: 12, winRate: 58, color: '#4ecdc4' },
    { name: 'Vayne', games: 8, winRate: 62, color: '#45b7d1' },
    { name: 'Ashe', games: 7, winRate: 43, color: '#96ceb4' },
    { name: 'Others', games: 14, winRate: 50, color: '#feca57' },
  ],
};

interface StatCardProps {
  title: string;
  value: string | number;
  trend?: 'up' | 'down' | 'neutral';
  icon: React.ReactNode;
  color: string;
}

const StatCard = ({ title, value, trend, icon, color }: StatCardProps) => {
  const getTrendIcon = () => {
    if (trend === 'up') return <TrendingUp sx={{ color: 'success.main', fontSize: 20 }} />;
    if (trend === 'down') return <TrendingDown sx={{ color: 'error.main', fontSize: 20 }} />;
    return null;
  };

  return (
    <Card sx={{ height: '100%', transition: 'all 0.3s ease' }}>
      <CardContent>
        <Stack direction="row" justifyContent="space-between" alignItems="flex-start">
          <Box>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              {title}
            </Typography>
            <Typography variant="h4" component="div" sx={{ fontWeight: 'bold', color }}>
              {value}
            </Typography>
          </Box>
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1 }}>
            {icon}
            {getTrendIcon()}
          </Box>
        </Stack>
      </CardContent>
    </Card>
  );
};

const DashboardPage = () => {
  const [isLive, setIsLive] = useState(false);
  const [lastUpdate, setLastUpdate] = useState(new Date());

  // Simulate real-time updates
  useEffect(() => {
    const interval = setInterval(() => {
      setLastUpdate(new Date());
    }, 30000); // Update every 30 seconds

    return () => clearInterval(interval);
  }, []);

  const handleRefresh = () => {
    setLastUpdate(new Date());
  };

  return (
    <Box>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Typography variant="h4" gutterBottom sx={{ fontWeight: 'bold', color: 'primary.main' }}>
            🎮 Gaming Dashboard
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Welcome back, {mockPlayerStats.summoner.name}! Ready to climb?
          </Typography>
        </Box>
        <Stack direction="row" spacing={2} alignItems="center">
          <Typography variant="caption" color="text.secondary">
            Last updated: {lastUpdate.toLocaleTimeString()}
          </Typography>
          <IconButton onClick={handleRefresh} size="small">
            <Refresh />
          </IconButton>
          {isLive && (
            <Chip 
              label="🔴 LIVE" 
              color="error" 
              variant="filled" 
              sx={{ animation: 'pulse 2s infinite' }}
            />
          )}
        </Stack>
      </Box>

      <Grid container spacing={3}>
        {/* Summoner Info Card */}
        <Grid item xs={12} lg={4}>
          <Card sx={{ height: 'fit-content' }}>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center" mb={2}>
                <Avatar
                  src={`/api/cdn/img/profileicon/${mockPlayerStats.summoner.profileIconId}.png`}
                  sx={{ width: 60, height: 60, border: '3px solid', borderColor: 'secondary.main' }}
                />
                <Box flexGrow={1}>
                  <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                    {mockPlayerStats.summoner.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    #{mockPlayerStats.summoner.tagLine} • Level {mockPlayerStats.summoner.level}
                  </Typography>
                  <Stack direction="row" spacing={1} mt={1}>
                    <Chip 
                      label={`${mockPlayerStats.summoner.tier} ${mockPlayerStats.summoner.rank}`}
                      color="secondary" 
                      size="small" 
                    />
                    <Chip 
                      label={`${mockPlayerStats.summoner.leaguePoints} LP`} 
                      variant="outlined" 
                      size="small"
                    />
                  </Stack>
                </Box>
              </Stack>
              
              <Divider sx={{ mb: 2 }} />
              
              <Box>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Ranked Progress
                </Typography>
                <LinearProgress 
                  variant="determinate" 
                  value={(mockPlayerStats.summoner.leaguePoints / 100) * 100} 
                  sx={{ height: 8, borderRadius: 4, mb: 1 }}
                />
                <Stack direction="row" justifyContent="space-between">
                  <Typography variant="caption">
                    {mockPlayerStats.summoner.wins}W / {mockPlayerStats.summoner.losses}L
                  </Typography>
                  <Typography variant="caption">
                    {Math.round((mockPlayerStats.summoner.wins / (mockPlayerStats.summoner.wins + mockPlayerStats.summoner.losses)) * 100)}% WR
                  </Typography>
                </Stack>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Performance Stats */}
        <Grid item xs={12} lg={8}>
          <Grid container spacing={2}>
            <Grid item xs={6} md={4}>
              <StatCard
                title="KDA Ratio"
                value={mockPlayerStats.recentPerformance.kda}
                trend={mockPlayerStats.recentPerformance.kdaTrend as 'up' | 'down'}
                icon={<EmojiEvents sx={{ color: 'warning.main', fontSize: 28 }} />}
                color="warning.main"
              />
            </Grid>
            <Grid item xs={6} md={4}>
              <StatCard
                title="Win Rate"
                value={`${mockPlayerStats.recentPerformance.winRate}%`}
                trend={mockPlayerStats.recentPerformance.winRateTrend as 'up' | 'down'}
                icon={<TrendingUp sx={{ color: 'success.main', fontSize: 28 }} />}
                color="success.main"
              />
            </Grid>
            <Grid item xs={6} md={4}>
              <StatCard
                title="CS/min"
                value={mockPlayerStats.recentPerformance.csPerMin}
                trend={mockPlayerStats.recentPerformance.csPerMinTrend as 'up' | 'down'}
                icon={<SportsSoccer sx={{ color: 'info.main', fontSize: 28 }} />}
                color="info.main"
              />
            </Grid>
            <Grid item xs={6} md={4}>
              <StatCard
                title="Vision Score"
                value={mockPlayerStats.recentPerformance.visionScore}
                trend={mockPlayerStats.recentPerformance.visionTrend as 'up' | 'down'}
                icon={<Visibility sx={{ color: 'secondary.main', fontSize: 28 }} />}
                color="secondary.main"
              />
            </Grid>
            <Grid item xs={6} md={4}>
              <StatCard
                title="Gold/min"
                value={mockPlayerStats.recentPerformance.goldPerMin}
                trend={mockPlayerStats.recentPerformance.goldTrend as 'up' | 'down'}
                icon={<AttachMoney sx={{ color: 'warning.main', fontSize: 28 }} />}
                color="warning.main"
              />
            </Grid>
            <Grid item xs={6} md={4}>
              <Card sx={{ height: '100%' }}>
                <CardContent sx={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center' }}>
                  {mockPlayerStats.currentMatch ? (
                    <Button
                      variant="contained"
                      color="success"
                      startIcon={<PlayArrow />}
                      fullWidth
                      sx={{ mb: 1 }}
                    >
                      Watch Live
                    </Button>
                  ) : (
                    <Button
                      variant="outlined"
                      startIcon={<PersonAdd />}
                      fullWidth
                      sx={{ mb: 1 }}
                    >
                      Queue Up
                    </Button>
                  )}
                  <Typography variant="caption" color="text.secondary" textAlign="center">
                    {mockPlayerStats.currentMatch ? 'Currently in game' : 'Ready to play'}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Grid>

        {/* Rank Progress Chart */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Timeline color="primary" />
                Rank Progress (Last 6 Months)
              </Typography>
              <Box sx={{ width: '100%', height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={mockPlayerStats.rankProgress}>
                    <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
                    <XAxis 
                      dataKey="date" 
                      axisLine={false} 
                      tickLine={false}
                      tick={{ fontSize: 12 }}
                    />
                    <YAxis 
                      axisLine={false} 
                      tickLine={false}
                      tick={{ fontSize: 12 }}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: '#1e2328', 
                        border: 'none', 
                        borderRadius: 8,
                        color: '#f0e6d2'
                      }}
                    />
                    <Line 
                      type="monotone" 
                      dataKey="lp" 
                      stroke="#1976d2" 
                      strokeWidth={3}
                      dot={{ fill: '#1976d2', strokeWidth: 2, r: 4 }}
                      activeDot={{ r: 6, stroke: '#1976d2', strokeWidth: 2 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Champion Performance */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Champion Performance
              </Typography>
              <Box sx={{ width: '100%', height: 250 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={mockPlayerStats.championStats}
                      cx="50%"
                      cy="50%"
                      innerRadius={40}
                      outerRadius={80}
                      paddingAngle={5}
                      dataKey="games"
                    >
                      {mockPlayerStats.championStats.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: '#1e2328', 
                        border: 'none', 
                        borderRadius: 8,
                        color: '#f0e6d2'
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </Box>
              <Stack spacing={1} mt={2}>
                {mockPlayerStats.championStats.map((champ) => (
                  <Stack key={champ.name} direction="row" justifyContent="space-between" alignItems="center">
                    <Stack direction="row" alignItems="center" spacing={1}>
                      <Box sx={{ width: 12, height: 12, bgcolor: champ.color, borderRadius: '50%' }} />
                      <Typography variant="body2">{champ.name}</Typography>
                    </Stack>
                    <Typography variant="body2" color="text.secondary">
                      {champ.games}G • {champ.winRate}%
                    </Typography>
                  </Stack>
                ))}
              </Stack>
            </CardContent>
          </Card>
        </Grid>

        {/* Recent Matches */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Recent Matches
              </Typography>
              <Stack spacing={1}>
                {mockPlayerStats.matchHistory.map((match) => (
                  <Paper key={match.id} sx={{ p: 2, bgcolor: 'background.paper' }}>
                    <Stack direction="row" justifyContent="space-between" alignItems="center">
                      <Stack direction="row" alignItems="center" spacing={2}>
                        <Avatar
                          src={`/api/cdn/champion/${match.champion}.png`}
                          sx={{ width: 40, height: 40 }}
                        />
                        <Box>
                          <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                            {match.champion}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {match.duration} • {match.ago}
                          </Typography>
                        </Box>
                        <Chip 
                          label={match.result}
                          color={match.result === 'Victory' ? 'success' : 'error'}
                          size="small"
                          variant="outlined"
                        />
                      </Stack>
                      <Stack direction="row" alignItems="center" spacing={2}>
                        <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                          {match.kda}
                        </Typography>
                        <Button size="small" variant="text">
                          View Details
                        </Button>
                      </Stack>
                    </Stack>
                  </Paper>
                ))}
              </Stack>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default DashboardPage;