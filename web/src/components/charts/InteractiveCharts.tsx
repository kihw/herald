import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Typography,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Grid,
  IconButton,
  Tooltip,
  Chip,
  Switch,
  FormControlLabel,
  useTheme,
  alpha,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  ShowChart,
  BarChart,
  PieChart,
  ScatterPlot,
  Timeline,
  Fullscreen,
  Download,
  Refresh,
} from '@mui/icons-material';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip as ChartTooltip,
  Legend,
  Filler,
} from 'chart.js';
import {
  Line,
  Bar,
  Doughnut,
  Scatter,
  Radar,
  PolarArea,
} from 'react-chartjs-2';
import * as d3 from 'd3';

// Enregistrement des composants Chart.js
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  ChartTooltip,
  Legend,
  Filler
);

export type ChartType = 'line' | 'bar' | 'doughnut' | 'scatter' | 'radar' | 'polar' | 'heatmap' | 'd3-timeline';

export interface ChartDataPoint {
  x: any;
  y: number;
  label?: string;
  category?: string;
  metadata?: Record<string, any>;
}

export interface ChartSeries {
  name: string;
  data: ChartDataPoint[];
  color?: string;
  type?: ChartType;
}

export interface InteractiveChartsProps {
  title: string;
  series: ChartSeries[];
  chartType?: ChartType;
  height?: number;
  showControls?: boolean;
  realTimeUpdate?: boolean;
  onDataPointClick?: (point: ChartDataPoint, series: ChartSeries) => void;
  onExport?: (format: 'png' | 'svg' | 'pdf') => void;
}

export const InteractiveCharts: React.FC<InteractiveChartsProps> = ({
  title,
  series,
  chartType = 'line',
  height = 400,
  showControls = true,
  realTimeUpdate = false,
  onDataPointClick,
  onExport,
}) => {
  const theme = useTheme();
  const [selectedChartType, setSelectedChartType] = useState<ChartType>(chartType);
  const [showAnimation, setShowAnimation] = useState(true);
  const [showGrid, setShowGrid] = useState(true);
  const [selectedSeries, setSelectedSeries] = useState<string[]>(series.map(s => s.name));
  const d3ContainerRef = useRef<HTMLDivElement>(null);

  // Configuration des couleurs du thème
  const chartColors = [
    theme.palette.primary.main,
    theme.palette.secondary.main,
    theme.palette.success.main,
    theme.palette.warning.main,
    theme.palette.error.main,
    theme.palette.info.main,
  ];

  // Filtrage des séries sélectionnées
  const filteredSeries = series.filter(s => selectedSeries.includes(s.name));

  // Configuration des options communes
  const commonOptions = {
    responsive: true,
    maintainAspectRatio: false,
    animation: {
      duration: showAnimation ? 750 : 0,
    },
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          color: theme.palette.text.primary,
          usePointStyle: true,
        },
      },
      tooltip: {
        backgroundColor: alpha(theme.palette.background.paper, 0.95),
        titleColor: theme.palette.text.primary,
        bodyColor: theme.palette.text.secondary,
        borderColor: theme.palette.divider,
        borderWidth: 1,
      },
    },
    scales: {
      x: {
        grid: {
          display: showGrid,
          color: alpha(theme.palette.divider, 0.1),
        },
        ticks: {
          color: theme.palette.text.secondary,
        },
      },
      y: {
        grid: {
          display: showGrid,
          color: alpha(theme.palette.divider, 0.1),
        },
        ticks: {
          color: theme.palette.text.secondary,
        },
      },
    },
    onClick: (event: any, elements: any[]) => {
      if (elements.length > 0 && onDataPointClick) {
        const elementIndex = elements[0].index;
        const datasetIndex = elements[0].datasetIndex;
        const series = filteredSeries[datasetIndex];
        const point = series.data[elementIndex];
        onDataPointClick(point, series);
      }
    },
  };

  // Conversion des données pour Chart.js
  const getChartData = () => {
    const labels = filteredSeries[0]?.data.map(point => point.x) || [];
    
    const datasets = filteredSeries.map((serie, index) => ({
      label: serie.name,
      data: serie.data.map(point => point.y),
      backgroundColor: selectedChartType === 'line' 
        ? alpha(chartColors[index % chartColors.length], 0.1)
        : chartColors[index % chartColors.length],
      borderColor: chartColors[index % chartColors.length],
      borderWidth: 2,
      fill: selectedChartType === 'line',
      tension: 0.3,
      pointBackgroundColor: chartColors[index % chartColors.length],
      pointBorderColor: theme.palette.background.paper,
      pointBorderWidth: 2,
      pointRadius: 4,
      pointHoverRadius: 6,
    }));

    return { labels, datasets };
  };

  // Rendu D3.js pour les graphiques avancés
  const renderD3Chart = () => {
    if (!d3ContainerRef.current || selectedChartType !== 'd3-timeline') return;

    const container = d3.select(d3ContainerRef.current);
    container.selectAll('*').remove();

    const margin = { top: 20, right: 30, bottom: 40, left: 50 };
    const width = container.node()?.getBoundingClientRect().width || 800;
    const chartHeight = height - margin.top - margin.bottom;

    const svg = container
      .append('svg')
      .attr('width', width)
      .attr('height', height);

    const g = svg
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Timeline pour les données temporelles
    if (filteredSeries.length > 0) {
      const data = filteredSeries[0].data;
      const xScale = d3.scaleTime()
        .domain(d3.extent(data, d => new Date(d.x)) as [Date, Date])
        .range([0, width - margin.left - margin.right]);

      const yScale = d3.scaleLinear()
        .domain(d3.extent(data, d => d.y) as [number, number])
        .range([chartHeight, 0]);

      // Ligne de base
      const line = d3.line<ChartDataPoint>()
        .x(d => xScale(new Date(d.x)))
        .y(d => yScale(d.y))
        .curve(d3.curveMonotoneX);

      // Gradient
      const gradient = svg.append('defs')
        .append('linearGradient')
        .attr('id', 'timeline-gradient')
        .attr('gradientUnits', 'userSpaceOnUse')
        .attr('x1', 0).attr('y1', chartHeight)
        .attr('x2', 0).attr('y2', 0);

      gradient.append('stop')
        .attr('offset', '0%')
        .attr('stop-color', theme.palette.primary.main)
        .attr('stop-opacity', 0);

      gradient.append('stop')
        .attr('offset', '100%')
        .attr('stop-color', theme.palette.primary.main)
        .attr('stop-opacity', 0.3);

      // Area fill
      const area = d3.area<ChartDataPoint>()
        .x(d => xScale(new Date(d.x)))
        .y0(chartHeight)
        .y1(d => yScale(d.y))
        .curve(d3.curveMonotoneX);

      g.append('path')
        .datum(data)
        .attr('fill', 'url(#timeline-gradient)')
        .attr('d', area);

      // Line
      g.append('path')
        .datum(data)
        .attr('fill', 'none')
        .attr('stroke', theme.palette.primary.main)
        .attr('stroke-width', 2)
        .attr('d', line);

      // Points
      g.selectAll('.dot')
        .data(data)
        .enter().append('circle')
        .attr('class', 'dot')
        .attr('cx', d => xScale(new Date(d.x)))
        .attr('cy', d => yScale(d.y))
        .attr('r', 4)
        .attr('fill', theme.palette.primary.main)
        .on('click', (event, d) => {
          if (onDataPointClick) {
            onDataPointClick(d, filteredSeries[0]);
          }
        })
        .on('mouseover', function(event, d) {
          d3.select(this).attr('r', 6);
        })
        .on('mouseout', function() {
          d3.select(this).attr('r', 4);
        });

      // Axes
      g.append('g')
        .attr('transform', `translate(0,${chartHeight})`)
        .call(d3.axisBottom(xScale))
        .selectAll('text')
        .style('fill', theme.palette.text.secondary);

      g.append('g')
        .call(d3.axisLeft(yScale))
        .selectAll('text')
        .style('fill', theme.palette.text.secondary);
    }
  };

  // Heatmap avec D3
  const renderHeatmap = () => {
    if (!d3ContainerRef.current || selectedChartType !== 'heatmap') return;

    const container = d3.select(d3ContainerRef.current);
    container.selectAll('*').remove();

    // Simulation de données de heatmap (champions x rôles)
    const heatmapData = [
      { champion: 'Jinx', role: 'ADC', value: 85 },
      { champion: 'Jinx', role: 'Mid', value: 65 },
      { champion: 'Yasuo', role: 'Mid', value: 90 },
      { champion: 'Yasuo', role: 'Top', value: 75 },
      { champion: 'Thresh', role: 'Support', value: 88 },
      { champion: 'Leona', role: 'Support', value: 82 },
    ];

    const margin = { top: 80, right: 25, bottom: 30, left: 70 };
    const width = container.node()?.getBoundingClientRect().width || 600;
    const chartHeight = height - margin.top - margin.bottom;

    const svg = container
      .append('svg')
      .attr('width', width)
      .attr('height', height);

    const champions = Array.from(new Set(heatmapData.map(d => d.champion)));
    const roles = Array.from(new Set(heatmapData.map(d => d.role)));

    const xScale = d3.scaleBand()
      .range([0, width - margin.left - margin.right])
      .domain(champions)
      .padding(0.05);

    const yScale = d3.scaleBand()
      .range([chartHeight, 0])
      .domain(roles)
      .padding(0.05);

    const colorScale = d3.scaleSequential()
      .interpolator(d3.interpolateBlues)
      .domain([0, 100]);

    const g = svg.append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Rectangles de la heatmap
    g.selectAll()
      .data(heatmapData)
      .enter()
      .append('rect')
      .attr('x', d => xScale(d.champion) || 0)
      .attr('y', d => yScale(d.role) || 0)
      .attr('width', xScale.bandwidth())
      .attr('height', yScale.bandwidth())
      .style('fill', d => colorScale(d.value))
      .style('stroke', theme.palette.background.paper)
      .style('stroke-width', 2)
      .on('mouseover', function(event, d) {
        d3.select(this).style('stroke-width', 4);
      })
      .on('mouseout', function() {
        d3.select(this).style('stroke-width', 2);
      });

    // Labels
    g.append('g')
      .selectAll('text')
      .data(heatmapData)
      .enter()
      .append('text')
      .attr('x', d => (xScale(d.champion) || 0) + xScale.bandwidth() / 2)
      .attr('y', d => (yScale(d.role) || 0) + yScale.bandwidth() / 2)
      .attr('text-anchor', 'middle')
      .attr('dominant-baseline', 'middle')
      .style('fill', d => d.value > 50 ? 'white' : 'black')
      .style('font-size', '12px')
      .text(d => d.value);

    // Axes
    g.append('g')
      .style('font-size', '14px')
      .attr('transform', `translate(0, ${chartHeight})`)
      .call(d3.axisBottom(xScale).tickSize(0))
      .select('.domain').remove();

    g.append('g')
      .style('font-size', '14px')
      .call(d3.axisLeft(yScale).tickSize(0))
      .select('.domain').remove();
  };

  // Effet pour les graphiques D3
  useEffect(() => {
    if (selectedChartType === 'd3-timeline') {
      renderD3Chart();
    } else if (selectedChartType === 'heatmap') {
      renderHeatmap();
    }
  }, [selectedChartType, filteredSeries, theme, height, showGrid]);

  // Rendu du graphique selon le type
  const renderChart = () => {
    const data = getChartData();

    switch (selectedChartType) {
      case 'line':
        return <Line data={data} options={commonOptions} />;
      case 'bar':
        return <Bar data={data} options={commonOptions} />;
      case 'doughnut':
        return <Doughnut data={data} options={{ ...commonOptions, scales: undefined }} />;
      case 'scatter':
        return <Scatter data={data} options={commonOptions} />;
      case 'radar':
        return <Radar data={data} options={{ ...commonOptions, scales: undefined }} />;
      case 'polar':
        return <PolarArea data={data} options={{ ...commonOptions, scales: undefined }} />;
      case 'heatmap':
      case 'd3-timeline':
        return <div ref={d3ContainerRef} style={{ width: '100%', height: height }} />;
      default:
        return <Line data={data} options={commonOptions} />;
    }
  };

  return (
    <Card>
      <CardHeader
        title={
          <Box display="flex" alignItems="center" gap={1}>
            <ShowChart color="primary" />
            <Typography variant="h6">{title}</Typography>
            {realTimeUpdate && (
              <Chip 
                label="Live" 
                color="success" 
                size="small" 
                icon={<TrendingUp />} 
              />
            )}
          </Box>
        }
        action={
          showControls && (
            <Box display="flex" gap={1}>
              <Tooltip title="Rafraîchir">
                <IconButton size="small">
                  <Refresh />
                </IconButton>
              </Tooltip>
              <Tooltip title="Plein écran">
                <IconButton size="small">
                  <Fullscreen />
                </IconButton>
              </Tooltip>
              <Tooltip title="Exporter">
                <IconButton size="small" onClick={() => onExport?.('png')}>
                  <Download />
                </IconButton>
              </Tooltip>
            </Box>
          )
        }
      />
      <CardContent>
        {showControls && (
          <Grid container spacing={2} sx={{ mb: 2 }}>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Type de graphique</InputLabel>
                <Select
                  value={selectedChartType}
                  label="Type de graphique"
                  onChange={(e) => setSelectedChartType(e.target.value as ChartType)}
                >
                  <MenuItem value="line">
                    <Box display="flex" alignItems="center" gap={1}>
                      <ShowChart /> Ligne
                    </Box>
                  </MenuItem>
                  <MenuItem value="bar">
                    <Box display="flex" alignItems="center" gap={1}>
                      <BarChart /> Barres
                    </Box>
                  </MenuItem>
                  <MenuItem value="doughnut">
                    <Box display="flex" alignItems="center" gap={1}>
                      <PieChart /> Donut
                    </Box>
                  </MenuItem>
                  <MenuItem value="scatter">
                    <Box display="flex" alignItems="center" gap={1}>
                      <ScatterPlot /> Nuage
                    </Box>
                  </MenuItem>
                  <MenuItem value="radar">
                    <Box display="flex" alignItems="center" gap={1}>
                      <Timeline /> Radar
                    </Box>
                  </MenuItem>
                  <MenuItem value="heatmap">
                    <Box display="flex" alignItems="center" gap={1}>
                      <Timeline /> Heatmap
                    </Box>
                  </MenuItem>
                  <MenuItem value="d3-timeline">
                    <Box display="flex" alignItems="center" gap={1}>
                      <Timeline /> Timeline D3
                    </Box>
                  </MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControlLabel
                control={
                  <Switch
                    checked={showAnimation}
                    onChange={(e) => setShowAnimation(e.target.checked)}
                  />
                }
                label="Animations"
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControlLabel
                control={
                  <Switch
                    checked={showGrid}
                    onChange={(e) => setShowGrid(e.target.checked)}
                  />
                }
                label="Grille"
              />
            </Grid>
            <Grid item xs={12} md={3}>
              <Box display="flex" flexWrap="wrap" gap={0.5}>
                {series.map((serie) => (
                  <Chip
                    key={serie.name}
                    label={serie.name}
                    clickable
                    color={selectedSeries.includes(serie.name) ? 'primary' : 'default'}
                    onClick={() => {
                      setSelectedSeries(prev => 
                        prev.includes(serie.name)
                          ? prev.filter(s => s !== serie.name)
                          : [...prev, serie.name]
                      );
                    }}
                    size="small"
                  />
                ))}
              </Box>
            </Grid>
          </Grid>
        )}
        
        <Box sx={{ height, position: 'relative' }}>
          {renderChart()}
        </Box>
      </CardContent>
    </Card>
  );
};