import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  RadialLinearScale,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';
import { leagueColors } from '../../theme/leagueTheme';

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  RadialLinearScale,
  Title,
  Tooltip,
  Legend,
  Filler
);

// League of Legends themed color palette
export const chartColors = {
  primary: leagueColors.blue[500],
  secondary: leagueColors.gold[500],
  success: leagueColors.win,
  warning: leagueColors.gold[400],
  error: leagueColors.loss,
  info: leagueColors.blue[300],
  dark: leagueColors.dark[400],
  gradient: {
    blue: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
    gold: `linear-gradient(135deg, ${leagueColors.gold[500]} 0%, ${leagueColors.gold[600]} 100%)`,
    win: `linear-gradient(135deg, ${leagueColors.win} 0%, #4caf50 100%)`,
    loss: `linear-gradient(135deg, ${leagueColors.loss} 0%, #f44336 100%)`,
  },
};

// Color palette for multiple data series
export const seriesColors = [
  leagueColors.blue[500],
  leagueColors.gold[500],
  leagueColors.win,
  leagueColors.loss,
  leagueColors.blue[300],
  leagueColors.gold[300],
  leagueColors.dark[400],
  '#9c27b0', // Purple
  '#ff9800', // Orange
  '#795548', // Brown
];

// Default chart options
export const getDefaultOptions = (isDarkMode: boolean) => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        color: isDarkMode ? '#ffffff' : '#333333',
        usePointStyle: true,
        padding: 20,
        font: {
          family: '"Roboto", "Helvetica", "Arial", sans-serif',
          size: 12,
          weight: '500',
        },
      },
    },
    title: {
      display: false, // We handle titles externally
    },
    tooltip: {
      backgroundColor: isDarkMode ? 'rgba(0, 0, 0, 0.8)' : 'rgba(255, 255, 255, 0.95)',
      titleColor: isDarkMode ? '#ffffff' : '#333333',
      bodyColor: isDarkMode ? '#ffffff' : '#333333',
      borderColor: leagueColors.blue[500],
      borderWidth: 1,
      cornerRadius: 8,
      padding: 12,
      titleFont: {
        size: 14,
        weight: 'bold',
      },
      bodyFont: {
        size: 13,
      },
    },
  },
  scales: {
    x: {
      grid: {
        color: isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)',
        lineWidth: 1,
      },
      ticks: {
        color: isDarkMode ? '#ffffff' : '#666666',
        font: {
          size: 11,
        },
      },
    },
    y: {
      grid: {
        color: isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)',
        lineWidth: 1,
      },
      ticks: {
        color: isDarkMode ? '#ffffff' : '#666666',
        font: {
          size: 11,
        },
      },
    },
  },
  elements: {
    bar: {
      borderRadius: 4,
      borderSkipped: false,
    },
    line: {
      tension: 0.4,
      borderWidth: 3,
      fill: false,
    },
    point: {
      radius: 5,
      hoverRadius: 8,
      borderWidth: 2,
    },
  },
  interaction: {
    intersect: false,
    mode: 'index' as const,
  },
  animation: {
    duration: 1000,
    easing: 'easeInOutQuart' as const,
  },
});

// Radar chart specific options
export const getRadarOptions = (isDarkMode: boolean) => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        color: isDarkMode ? '#ffffff' : '#333333',
        usePointStyle: true,
        padding: 20,
      },
    },
    tooltip: {
      backgroundColor: isDarkMode ? 'rgba(0, 0, 0, 0.8)' : 'rgba(255, 255, 255, 0.95)',
      titleColor: isDarkMode ? '#ffffff' : '#333333',
      bodyColor: isDarkMode ? '#ffffff' : '#333333',
      borderColor: leagueColors.blue[500],
      borderWidth: 1,
      cornerRadius: 8,
      padding: 12,
    },
  },
  scales: {
    r: {
      angleLines: {
        color: isDarkMode ? 'rgba(255, 255, 255, 0.2)' : 'rgba(0, 0, 0, 0.2)',
      },
      grid: {
        color: isDarkMode ? 'rgba(255, 255, 255, 0.2)' : 'rgba(0, 0, 0, 0.2)',
      },
      pointLabels: {
        color: isDarkMode ? '#ffffff' : '#333333',
        font: {
          size: 12,
          weight: '500',
        },
      },
      ticks: {
        color: isDarkMode ? '#ffffff' : '#666666',
        backdropColor: 'transparent',
        font: {
          size: 10,
        },
      },
      min: 0,
    },
  },
  elements: {
    line: {
      borderWidth: 3,
    },
    point: {
      borderWidth: 2,
      radius: 4,
      hoverRadius: 6,
    },
  },
  animation: {
    duration: 1200,
    easing: 'easeInOutQuart' as const,
  },
});

// Pie/Doughnut chart specific options  
export const getPieOptions = (isDarkMode: boolean) => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'right' as const,
      labels: {
        color: isDarkMode ? '#ffffff' : '#333333',
        usePointStyle: true,
        padding: 15,
        generateLabels: (chart: any) => {
          const data = chart.data;
          if (data.labels && data.datasets.length) {
            return data.labels.map((label: string, i: number) => ({
              text: label,
              fillStyle: data.datasets[0].backgroundColor[i],
              strokeStyle: data.datasets[0].borderColor[i],
              lineWidth: 2,
              hidden: isNaN(data.datasets[0].data[i]) || chart.getDatasetMeta(0).data[i].hidden,
              index: i,
              pointStyle: 'circle',
            }));
          }
          return [];
        },
      },
    },
    tooltip: {
      backgroundColor: isDarkMode ? 'rgba(0, 0, 0, 0.8)' : 'rgba(255, 255, 255, 0.95)',
      titleColor: isDarkMode ? '#ffffff' : '#333333',
      bodyColor: isDarkMode ? '#ffffff' : '#333333',
      borderColor: leagueColors.blue[500],
      borderWidth: 1,
      cornerRadius: 8,
      padding: 12,
      callbacks: {
        label: (context: any) => {
          const label = context.label || '';
          const value = context.parsed;
          const total = context.dataset.data.reduce((a: number, b: number) => a + b, 0);
          const percentage = ((value / total) * 100).toFixed(1);
          return `${label}: ${value} (${percentage}%)`;
        },
      },
    },
  },
  elements: {
    arc: {
      borderWidth: 2,
      hoverBorderWidth: 3,
    },
  },
  animation: {
    animateRotate: true,
    animateScale: true,
    duration: 1000,
    easing: 'easeInOutQuart' as const,
  },
});

// Generate gradient background for canvas
export const createGradient = (ctx: CanvasRenderingContext2D, startColor: string, endColor: string, vertical = true) => {
  const gradient = vertical 
    ? ctx.createLinearGradient(0, 0, 0, ctx.canvas.height)
    : ctx.createLinearGradient(0, 0, ctx.canvas.width, 0);
  
  gradient.addColorStop(0, startColor);
  gradient.addColorStop(1, endColor);
  
  return gradient;
};

// Helper to format chart data values
export const formatChartValue = (value: number, type: 'percentage' | 'decimal' | 'integer' | 'kda' = 'decimal') => {
  switch (type) {
    case 'percentage':
      return `${(value * 100).toFixed(1)}%`;
    case 'integer':
      return Math.round(value).toString();
    case 'kda':
      return value.toFixed(2);
    case 'decimal':
    default:
      return value.toFixed(1);
  }
};