import React from 'react';
import { Line } from 'react-chartjs-2';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Avatar,
  Grid,
  Chip,
  useTheme,
  LinearProgress,
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
} from '@mui/icons-material';
import { getDefaultOptions, seriesColors } from './ChartConfig';
import { leagueColors } from '../../theme/leagueTheme';

interface PlayerStats {
  riot_id: string;
  riot_tag: string;
  rank?: string;
  lp?: number;
  winrate: number;
  kda: number;
  recent_games: number;
  trend: 'up' | 'down' | 'stable';
  performance_history: number[]; // Last 10 games winrate
}

interface PlayerPerformanceWidgetProps {
  player: PlayerStats;
  compact?: boolean;
}

const PlayerPerformanceWidget: React.FC<PlayerPerformanceWidgetProps> = ({ 
  player, 
  compact = false 
}) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';

  const getTrendIcon = () => {
    switch (player.trend) {
      case 'up':
        return <TrendingUpIcon sx={{ fontSize: 16, color: leagueColors.win }} />;
      case 'down':
        return <TrendingDownIcon sx={{ fontSize: 16, color: leagueColors.loss }} />;
      default:
        return null;
    }
  };

  const getTrendColor = () => {
    switch (player.trend) {
      case 'up':
        return leagueColors.win;
      case 'down':
        return leagueColors.loss;
      default:
        return 'text.secondary';
    }
  };

  const performanceChartData = {
    labels: Array.from({ length: player.performance_history.length }, (_, i) => `G${i + 1}`),
    datasets: [
      {
        label: 'Performance',
        data: player.performance_history,
        borderColor: player.winrate >= 0.5 ? leagueColors.win : leagueColors.blue[500],
        backgroundColor: 'transparent',
        borderWidth: 2,
        pointRadius: 2,
        pointHoverRadius: 4,
        tension: 0.4,
      },
    ],
  };

  const performanceChartOptions = {
    ...getDefaultOptions(isDarkMode),
    plugins: {
      ...getDefaultOptions(isDarkMode).plugins,
      legend: {
        display: false,
      },
      tooltip: {
        ...getDefaultOptions(isDarkMode).plugins?.tooltip,
        callbacks: {
          label: (context: any) => {
            return `Partie ${context.dataIndex + 1}: ${context.parsed.y > 0 ? 'Victoire' : 'Défaite'}`;
          },
        },
      },
    },
    scales: {
      x: {
        display: false,
      },
      y: {
        display: false,
        min: 0,
        max: 1,
      },
    },
    elements: {
      point: {
        radius: 3,
        hoverRadius: 5,
      },
    },
    maintainAspectRatio: false,
  };

  if (compact) {
    return (
      <Card
        sx={{
          background: isDarkMode
            ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
            : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
          border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
        }}
      >
        <CardContent sx={{ p: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Avatar
              sx={{
                width: 40,
                height: 40,
                background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
                color: '#fff',
                fontWeight: 600,
                fontSize: '1rem',
              }}
            >
              {player.riot_id.charAt(0).toUpperCase()}
            </Avatar>
            <Box sx={{ flexGrow: 1, minWidth: 0 }}>
              <Typography variant="body1" sx={{ fontWeight: 500 }} noWrap>
                {player.riot_id}#{player.riot_tag}
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Chip
                  label={`${(player.winrate * 100).toFixed(1)}%`}
                  size="small"
                  sx={{
                    height: 20,
                    fontSize: '0.7rem',
                    background: `${player.winrate >= 0.5 ? leagueColors.win : leagueColors.loss}20`,
                    color: player.winrate >= 0.5 ? leagueColors.win : leagueColors.loss,
                    border: 'none',
                  }}
                />
                {getTrendIcon()}
              </Box>
            </Box>
            <Box sx={{ width: 60, height: 30 }}>
              <Line data={performanceChartData} options={performanceChartOptions} />
            </Box>
          </Box>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card
      sx={{
        background: isDarkMode
          ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
          : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
        border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
      }}
    >
      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
          <Avatar
            sx={{
              width: 56,
              height: 56,
              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
              color: '#fff',
              fontWeight: 600,
              fontSize: '1.2rem',
            }}
          >
            {player.riot_id.charAt(0).toUpperCase()}
          </Avatar>
          <Box sx={{ flexGrow: 1 }}>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>
              {player.riot_id}#{player.riot_tag}
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
              {player.rank && (
                <Chip
                  label={player.rank}
                  size="small"
                  variant="outlined"
                  sx={{ fontSize: '0.7rem' }}
                />
              )}
              {player.lp && (
                <Chip
                  label={`${player.lp} LP`}
                  size="small"
                  variant="outlined"
                  sx={{ fontSize: '0.7rem' }}
                />
              )}
            </Box>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            {getTrendIcon()}
            <Typography variant="body2" sx={{ color: getTrendColor() }}>
              {player.trend === 'up' ? 'En progression' : 
               player.trend === 'down' ? 'En régression' : 'Stable'}
            </Typography>
          </Box>
        </Box>

        {/* Stats Grid */}
        <Grid container spacing={2} sx={{ mb: 3 }}>
          <Grid item xs={6} sm={3}>
            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="h4" sx={{ fontWeight: 700, color: player.winrate >= 0.5 ? leagueColors.win : leagueColors.loss }}>
                {(player.winrate * 100).toFixed(1)}%
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Taux de victoire
              </Typography>
            </Box>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="h4" sx={{ fontWeight: 700 }}>
                {player.kda.toFixed(2)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                KDA Moyen
              </Typography>
            </Box>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="h4" sx={{ fontWeight: 700 }}>
                {player.recent_games}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Parties récentes
              </Typography>
            </Box>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="h4" sx={{ fontWeight: 700, color: leagueColors.blue[500] }}>
                {Math.floor(Math.random() * 50) + 70}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Score de forme
              </Typography>
            </Box>
          </Grid>
        </Grid>

        {/* Performance Progress */}
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
            <Typography variant="body2" color="text.secondary">
              Performance générale
            </Typography>
            <Typography variant="body2" sx={{ fontWeight: 500 }}>
              {(player.winrate * 100).toFixed(1)}%
            </Typography>
          </Box>
          <LinearProgress
            variant="determinate"
            value={player.winrate * 100}
            sx={{
              height: 8,
              borderRadius: 4,
              backgroundColor: `${leagueColors.loss}30`,
              '& .MuiLinearProgress-bar': {
                backgroundColor: player.winrate >= 0.5 ? leagueColors.win : leagueColors.loss,
                borderRadius: 4,
              },
            }}
          />
        </Box>

        {/* Performance History Chart */}
        <Box>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
            Historique des performances (10 dernières parties)
          </Typography>
          <Box sx={{ height: 80 }}>
            <Line data={performanceChartData} options={performanceChartOptions} />
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default PlayerPerformanceWidget;