import React from 'react';
import {
  Bar,
  Line,
  Radar,
  Doughnut,
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
  getRadarOptions,
  getPieOptions,
  seriesColors,
  chartColors,
  formatChartValue,
} from './ChartConfig';
import { leagueColors } from '../../theme/leagueTheme';
import { ChartData, ChartDataset } from '../../services/groupApi';

interface ComparisonChartsProps {
  charts: ChartData[];
}

const ComparisonCharts: React.FC<ComparisonChartsProps> = ({ charts }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';

  const processChartData = (chart: ChartData) => {
    const processedData = {
      labels: chart.labels,
      datasets: chart.datasets.map((dataset, index) => ({
        ...dataset,
        backgroundColor: dataset.background_color || [seriesColors[index % seriesColors.length]],
        borderColor: dataset.border_color || [seriesColors[index % seriesColors.length]],
        borderWidth: 2,
        ...(chart.type === 'bar' && {
          backgroundColor: `${seriesColors[index % seriesColors.length]}80`,
          borderColor: seriesColors[index % seriesColors.length],
          hoverBackgroundColor: `${seriesColors[index % seriesColors.length]}cc`,
        }),
        ...(chart.type === 'line' && {
          backgroundColor: 'transparent',
          borderColor: seriesColors[index % seriesColors.length],
          pointBackgroundColor: seriesColors[index % seriesColors.length],
          pointBorderColor: '#fff',
          pointBorderWidth: 2,
          fill: false,
          tension: 0.4,
        }),
        ...(chart.type === 'radar' && {
          backgroundColor: `${seriesColors[index % seriesColors.length]}20`,
          borderColor: seriesColors[index % seriesColors.length],
          pointBackgroundColor: seriesColors[index % seriesColors.length],
          pointBorderColor: '#fff',
          pointHoverBackgroundColor: '#fff',
          pointHoverBorderColor: seriesColors[index % seriesColors.length],
        }),
        ...(chart.type === 'pie' && {
          backgroundColor: chart.labels.map((_, i) => seriesColors[i % seriesColors.length]),
          borderColor: chart.labels.map(() => isDarkMode ? '#1e1e1e' : '#ffffff'),
          hoverBackgroundColor: chart.labels.map((_, i) => `${seriesColors[i % seriesColors.length]}cc`),
        }),
      })),
    };

    return processedData;
  };

  const getChartOptions = (chart: ChartData) => {
    const baseOptions = chart.type === 'radar' 
      ? getRadarOptions(isDarkMode)
      : chart.type === 'pie'
      ? getPieOptions(isDarkMode)
      : getDefaultOptions(isDarkMode);

    // Merge with custom options from chart data
    return {
      ...baseOptions,
      ...chart.options,
      plugins: {
        ...baseOptions.plugins,
        ...(chart.options?.plugins || {}),
        tooltip: {
          ...baseOptions.plugins?.tooltip,
          callbacks: {
            ...baseOptions.plugins?.tooltip?.callbacks,
            label: (context: any) => {
              const label = context.dataset.label || '';
              const value = context.parsed.y !== undefined ? context.parsed.y : context.parsed;
              
              // Custom formatting based on chart type and data
              let formattedValue = value;
              if (chart.title.toLowerCase().includes('winrate') || chart.title.toLowerCase().includes('victoire')) {
                formattedValue = formatChartValue(value / 100, 'percentage');
              } else if (chart.title.toLowerCase().includes('kda')) {
                formattedValue = formatChartValue(value, 'kda');
              } else if (chart.title.toLowerCase().includes('cs') || chart.title.toLowerCase().includes('gold')) {
                formattedValue = formatChartValue(value, 'decimal');
              } else {
                formattedValue = formatChartValue(value, 'integer');
              }
              
              return `${label}: ${formattedValue}`;
            },
          },
        },
      },
    };
  };

  const renderChart = (chart: ChartData) => {
    const chartData = processChartData(chart);
    const options = getChartOptions(chart);

    switch (chart.type) {
      case 'bar':
        return <Bar data={chartData} options={options} />;
      case 'line':
        return <Line data={chartData} options={options} />;
      case 'radar':
        return <Radar data={chartData} options={options} />;
      case 'pie':
        return <Doughnut data={chartData} options={options} />;
      default:
        return <Bar data={chartData} options={options} />;
    }
  };

  if (!charts || charts.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 6 }}>
        <Typography variant="h6" color="text.secondary" gutterBottom>
          Aucun graphique disponible
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Les visualisations seront générées automatiquement avec les données de comparaison
        </Typography>
      </Box>
    );
  }

  return (
    <Grid container spacing={3}>
      {charts.map((chart, index) => (
        <Grid 
          item 
          xs={12} 
          md={chart.type === 'pie' ? 6 : chart.type === 'radar' ? 6 : 12} 
          key={index}
        >
          <Card
            sx={{
              height: chart.type === 'line' ? 400 : 350,
              border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
              background: isDarkMode
                ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
                : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
            }}
          >
            <CardContent sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              <Typography 
                variant="h6" 
                sx={{ 
                  fontWeight: 600, 
                  mb: 2,
                  background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
                  backgroundClip: 'text',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                }}
              >
                {chart.title}
              </Typography>
              <Box sx={{ flexGrow: 1, position: 'relative', minHeight: 0 }}>
                {renderChart(chart)}
              </Box>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
};

export default ComparisonCharts;