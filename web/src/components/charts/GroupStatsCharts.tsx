import React from 'react';
import {
  Bar,
  Doughnut,
  Line,
} from 'react-chartjs-2';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  useTheme,
} from '@mui/material';
import {
  getDefaultOptions,
  getPieOptions,
  seriesColors,
  formatChartValue,
} from './ChartConfig';
import { leagueColors } from '../../theme/leagueTheme';
import { GroupStats } from '../../services/groupApi';

interface GroupStatsChartsProps {
  stats: GroupStats;
}

const GroupStatsCharts: React.FC<GroupStatsChartsProps> = ({ stats }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';

  // Champions Performance Chart
  const getChampionsChart = () => {
    if (!stats.top_champions || stats.top_champions.length === 0) return null;

    const championsData = {
      labels: stats.top_champions.slice(0, 8).map(champion => champion.champion_name),
      datasets: [
        {
          label: 'Taux de Victoire (%)',
          data: stats.top_champions.slice(0, 8).map(champion => champion.win_rate * 100),
          backgroundColor: stats.top_champions.slice(0, 8).map((_, index) => 
            `${seriesColors[index % seriesColors.length]}80`
          ),
          borderColor: stats.top_champions.slice(0, 8).map((_, index) => 
            seriesColors[index % seriesColors.length]
          ),
          borderWidth: 2,
        },
        {
          label: 'KDA Moyen',
          data: stats.top_champions.slice(0, 8).map(champion => champion.avg_kda),
          backgroundColor: `${leagueColors.gold[500]}60`,
          borderColor: leagueColors.gold[500],
          borderWidth: 2,
          yAxisID: 'y1',
        },
      ],
    };

    const championsOptions = {
      ...getDefaultOptions(isDarkMode),
      scales: {
        ...getDefaultOptions(isDarkMode).scales,
        y: {
          ...getDefaultOptions(isDarkMode).scales?.y,
          type: 'linear' as const,
          display: true,
          position: 'left' as const,
          title: {
            display: true,
            text: 'Taux de Victoire (%)',
            color: isDarkMode ? '#ffffff' : '#333333',
          },
          max: 100,
        },
        y1: {
          ...getDefaultOptions(isDarkMode).scales?.y,
          type: 'linear' as const,
          display: true,
          position: 'right' as const,
          title: {
            display: true,
            text: 'KDA Moyen',
            color: isDarkMode ? '#ffffff' : '#333333',
          },
          grid: {
            drawOnChartArea: false,
          },
          max: Math.max(...stats.top_champions.map(c => c.avg_kda)) + 1,
        },
      },
      plugins: {
        ...getDefaultOptions(isDarkMode).plugins,
        tooltip: {
          ...getDefaultOptions(isDarkMode).plugins?.tooltip,
          callbacks: {
            label: (context: any) => {
              const label = context.dataset.label;
              const value = context.parsed.y;
              if (label.includes('Victoire')) {
                return `${label}: ${value.toFixed(1)}%`;
              } else {
                return `${label}: ${value.toFixed(2)}`;
              }
            },
          },
        },
      },
    };

    return { data: championsData, options: championsOptions };
  };

  // Roles Distribution Chart
  const getRolesChart = () => {
    if (!stats.popular_roles || stats.popular_roles.length === 0) return null;

    const rolesData = {
      labels: stats.popular_roles.map(role => {
        const roleNames: { [key: string]: string } = {
          'TOP': 'Top Lane',
          'JUNGLE': 'Jungle',
          'MIDDLE': 'Mid Lane',
          'BOTTOM': 'ADC',
          'UTILITY': 'Support',
        };
        return roleNames[role.role] || role.role;
      }),
      datasets: [
        {
          label: 'Parties Jouées',
          data: stats.popular_roles.map(role => role.play_count),
          backgroundColor: stats.popular_roles.map((_, index) => 
            `${seriesColors[index % seriesColors.length]}cc`
          ),
          borderColor: stats.popular_roles.map((_, index) => 
            seriesColors[index % seriesColors.length]
          ),
          borderWidth: 2,
        },
      ],
    };

    const rolesOptions = {
      ...getPieOptions(isDarkMode),
      plugins: {
        ...getPieOptions(isDarkMode).plugins,
        tooltip: {
          ...getPieOptions(isDarkMode).plugins?.tooltip,
          callbacks: {
            label: (context: any) => {
              const role = stats.popular_roles[context.dataIndex];
              const total = stats.popular_roles.reduce((sum, r) => sum + r.play_count, 0);
              const percentage = ((role.play_count / total) * 100).toFixed(1);
              const winrate = (role.win_rate * 100).toFixed(1);
              return [
                `${context.label}: ${role.play_count} parties (${percentage}%)`,
                `Taux de victoire: ${winrate}%`
              ];
            },
          },
        },
      },
    };

    return { data: rolesData, options: rolesOptions };
  };

  // Winrate Comparison Chart
  const getWinrateChart = () => {
    if (!stats.winrate_comparison || Object.keys(stats.winrate_comparison).length === 0) return null;

    const winrateEntries = Object.entries(stats.winrate_comparison);
    
    const winrateData = {
      labels: winrateEntries.map(([name]) => name),
      datasets: [
        {
          label: 'Taux de Victoire (%)',
          data: winrateEntries.map(([, winrate]) => winrate * 100),
          backgroundColor: winrateEntries.map(([, winrate]) => 
            winrate >= 0.5 ? `${leagueColors.win}80` : `${leagueColors.loss}80`
          ),
          borderColor: winrateEntries.map(([, winrate]) => 
            winrate >= 0.5 ? leagueColors.win : leagueColors.loss
          ),
          borderWidth: 2,
        },
      ],
    };

    const winrateOptions = {
      ...getDefaultOptions(isDarkMode),
      scales: {
        ...getDefaultOptions(isDarkMode).scales,
        y: {
          ...getDefaultOptions(isDarkMode).scales?.y,
          min: 0,
          max: 100,
          title: {
            display: true,
            text: 'Taux de Victoire (%)',
            color: isDarkMode ? '#ffffff' : '#333333',
          },
        },
      },
      plugins: {
        ...getDefaultOptions(isDarkMode).plugins,
        tooltip: {
          ...getDefaultOptions(isDarkMode).plugins?.tooltip,
          callbacks: {
            label: (context: any) => {
              return `Taux de victoire: ${context.parsed.y.toFixed(1)}%`;
            },
          },
        },
      },
    };

    return { data: winrateData, options: winrateOptions };
  };

  // Activity Overview Chart (mock data for demonstration)
  const getActivityChart = () => {
    // Generate mock activity data for the last 7 days
    const days = ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'];
    const mockActivityData = days.map(() => Math.floor(Math.random() * 20) + 5);

    const activityData = {
      labels: days,
      datasets: [
        {
          label: 'Parties par jour',
          data: mockActivityData,
          backgroundColor: `${leagueColors.blue[500]}60`,
          borderColor: leagueColors.blue[500],
          borderWidth: 3,
          fill: true,
          tension: 0.4,
        },
      ],
    };

    const activityOptions = {
      ...getDefaultOptions(isDarkMode),
      plugins: {
        ...getDefaultOptions(isDarkMode).plugins,
        tooltip: {
          ...getDefaultOptions(isDarkMode).plugins?.tooltip,
          callbacks: {
            label: (context: any) => {
              return `Parties jouées: ${context.parsed.y}`;
            },
          },
        },
      },
      scales: {
        ...getDefaultOptions(isDarkMode).scales,
        y: {
          ...getDefaultOptions(isDarkMode).scales?.y,
          beginAtZero: true,
          title: {
            display: true,
            text: 'Nombre de parties',
            color: isDarkMode ? '#ffffff' : '#333333',
          },
        },
      },
    };

    return { data: activityData, options: activityOptions };
  };

  const championsChart = getChampionsChart();
  const rolesChart = getRolesChart();
  const winrateChart = getWinrateChart();
  const activityChart = getActivityChart();

  return (
    <Grid container spacing={3} sx={{ mt: 2 }}>
      {/* Champions Performance */}
      {championsChart && (
        <Grid item xs={12}>
          <Card
            sx={{
              height: 400,
              border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
            }}
          >
            <CardContent sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Performance des Champions Populaires
              </Typography>
              <Box sx={{ flexGrow: 1, position: 'relative', minHeight: 0 }}>
                <Bar data={championsChart.data} options={championsChart.options} />
              </Box>
            </CardContent>
          </Card>
        </Grid>
      )}

      {/* Roles Distribution and Winrate Comparison */}
      <Grid item xs={12} md={6}>
        {rolesChart && (
          <Card
            sx={{
              height: 350,
              border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
            }}
          >
            <CardContent sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Distribution des Rôles
              </Typography>
              <Box sx={{ flexGrow: 1, position: 'relative', minHeight: 0 }}>
                <Doughnut data={rolesChart.data} options={rolesChart.options} />
              </Box>
            </CardContent>
          </Card>
        )}
      </Grid>

      <Grid item xs={12} md={6}>
        {winrateChart && (
          <Card
            sx={{
              height: 350,
              border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
            }}
          >
            <CardContent sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Comparaison des Taux de Victoire
              </Typography>
              <Box sx={{ flexGrow: 1, position: 'relative', minHeight: 0 }}>
                <Bar data={winrateChart.data} options={winrateChart.options} />
              </Box>
            </CardContent>
          </Card>
        )}
      </Grid>

      {/* Activity Chart */}
      {activityChart && (
        <Grid item xs={12}>
          <Card
            sx={{
              height: 300,
              border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
            }}
          >
            <CardContent sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Activité du Groupe (7 derniers jours)
              </Typography>
              <Box sx={{ flexGrow: 1, position: 'relative', minHeight: 0 }}>
                <Line data={activityChart.data} options={activityChart.options} />
              </Box>
            </CardContent>
          </Card>
        </Grid>
      )}
    </Grid>
  );
};

export default GroupStatsCharts;