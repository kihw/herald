import React, { useState, useEffect, useRef, useMemo } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Typography,
  IconButton,
  Tooltip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
  Chip,
  Avatar,
  useTheme,
  alpha,
  Button,
  Grid,
  Paper,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  Visibility,
  VisibilityOff,
  Download,
  Refresh,
  Settings,
  ZoomIn,
  ZoomOut,
} from '@mui/icons-material';
import * as d3 from 'd3';

export interface HeatmapDataPoint {
  x: string; // Champion/Role
  y: string; // Role/Champion
  value: number; // Performance metric (winrate, KDA, etc.)
  count?: number; // Number of games
  metadata?: {
    totalGames: number;
    wins: number;
    losses: number;
    avgKDA: number;
    avgKills: number;
    avgDeaths: number;
    avgAssists: number;
    avgCS: number;
    lastPlayed?: Date;
  };
}

export interface PerformanceHeatmapProps {
  data: HeatmapDataPoint[];
  title?: string;
  xAxisLabel?: string;
  yAxisLabel?: string;
  metric?: 'winrate' | 'kda' | 'games' | 'performance';
  colorScheme?: 'blues' | 'reds' | 'greens' | 'oranges' | 'purples' | 'rainbow';
  showLabels?: boolean;
  showTooltip?: boolean;
  showLegend?: boolean;
  interactive?: boolean;
  height?: number;
  width?: number;
  onCellClick?: (dataPoint: HeatmapDataPoint) => void;
  onCellHover?: (dataPoint: HeatmapDataPoint | null) => void;
}

export const PerformanceHeatmap: React.FC<PerformanceHeatmapProps> = ({
  data,
  title = 'Performance Heatmap',
  xAxisLabel = 'Champions',
  yAxisLabel = 'Roles',
  metric = 'winrate',
  colorScheme = 'blues',
  showLabels = true,
  showTooltip = true,
  showLegend = true,
  interactive = true,
  height = 400,
  width = 600,
  onCellClick,
  onCellHover,
}) => {
  const theme = useTheme();
  const svgRef = useRef<SVGSVGElement>(null);
  const [selectedMetric, setSelectedMetric] = useState(metric);
  const [selectedColorScheme, setSelectedColorScheme] = useState(colorScheme);
  const [showValues, setShowValues] = useState(showLabels);
  const [hoveredCell, setHoveredCell] = useState<HeatmapDataPoint | null>(null);
  const [zoomLevel, setZoomLevel] = useState(1);

  // Configuration des couleurs selon le thème
  const getColorScale = (scheme: string) => {
    const isDark = theme.palette.mode === 'dark';
    
    switch (scheme) {
      case 'blues':
        return isDark 
          ? d3.scaleSequential(d3.interpolateBlues).domain([0, 100])
          : d3.scaleSequential(d3.interpolateBlues).domain([0, 100]);
      case 'reds':
        return d3.scaleSequential(d3.interpolateReds).domain([0, 100]);
      case 'greens':
        return d3.scaleSequential(d3.interpolateGreens).domain([0, 100]);
      case 'oranges':
        return d3.scaleSequential(d3.interpolateOranges).domain([0, 100]);
      case 'purples':
        return d3.scaleSequential(d3.interpolatePurples).domain([0, 100]);
      case 'rainbow':
        return d3.scaleSequential(d3.interpolateViridis).domain([0, 100]);
      default:
        return d3.scaleSequential(d3.interpolateBlues).domain([0, 100]);
    }
  };

  // Préparation des données pour la heatmap
  const processedData = useMemo(() => {
    if (!data.length) return { matrix: [], xLabels: [], yLabels: [], extent: [0, 100] };

    // Extraire les labels uniques
    const xLabels = Array.from(new Set(data.map(d => d.x))).sort();
    const yLabels = Array.from(new Set(data.map(d => d.y))).sort();

    // Créer une matrice
    const matrix: Array<Array<HeatmapDataPoint | null>> = yLabels.map(() => 
      xLabels.map(() => null)
    );

    // Remplir la matrice
    data.forEach(point => {
      const xIndex = xLabels.indexOf(point.x);
      const yIndex = yLabels.indexOf(point.y);
      if (xIndex !== -1 && yIndex !== -1) {
        matrix[yIndex][xIndex] = point;
      }
    });

    // Calculer l'étendue des valeurs
    const values = data.map(d => d.value).filter(v => !isNaN(v));
    const extent = d3.extent(values) as [number, number];

    return { matrix, xLabels, yLabels, extent: extent || [0, 100] };
  }, [data]);

  // Rendu de la heatmap avec D3
  useEffect(() => {
    if (!svgRef.current || !processedData.matrix.length) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    const margin = { top: 80, right: 60, bottom: 80, left: 100 };
    const chartWidth = width - margin.left - margin.right;
    const chartHeight = height - margin.top - margin.bottom;

    const container = svg
      .attr('width', width)
      .attr('height', height)
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Échelles
    const xScale = d3.scaleBand()
      .domain(processedData.xLabels)
      .range([0, chartWidth])
      .padding(0.05);

    const yScale = d3.scaleBand()
      .domain(processedData.yLabels)
      .range([0, chartHeight])
      .padding(0.05);

    const colorScale = getColorScale(selectedColorScheme);

    // Dégradé pour les cellules vides
    const defs = svg.append('defs');
    const pattern = defs.append('pattern')
      .attr('id', 'diagonalHatch')
      .attr('patternUnits', 'userSpaceOnUse')
      .attr('width', 8)
      .attr('height', 8);

    pattern.append('path')
      .attr('d', 'M0,8 L8,0')
      .attr('stroke', theme.palette.divider)
      .attr('stroke-width', 1);

    // Cellules de la heatmap
    const cells = container.selectAll('.cell')
      .data(processedData.matrix.flat().filter(Boolean) as HeatmapDataPoint[])
      .enter()
      .append('g')
      .attr('class', 'cell');

    // Rectangles des cellules
    cells.append('rect')
      .attr('x', d => xScale(d.x) || 0)
      .attr('y', d => yScale(d.y) || 0)
      .attr('width', xScale.bandwidth())
      .attr('height', yScale.bandwidth())
      .attr('fill', d => isNaN(d.value) ? 'url(#diagonalHatch)' : colorScale(d.value))
      .attr('stroke', theme.palette.background.paper)
      .attr('stroke-width', 1)
      .attr('rx', 2)
      .style('cursor', interactive ? 'pointer' : 'default')
      .on('click', (event, d) => {
        if (interactive && onCellClick) {
          onCellClick(d);
        }
      })
      .on('mouseover', (event, d) => {
        if (interactive) {
          setHoveredCell(d);
          onCellHover?.(d);
          
          // Effet de survol
          d3.select(event.target)
            .transition()
            .duration(200)
            .attr('stroke-width', 3)
            .attr('stroke', theme.palette.primary.main);
        }
      })
      .on('mouseout', (event, d) => {
        if (interactive) {
          setHoveredCell(null);
          onCellHover?.(null);
          
          d3.select(event.target)
            .transition()
            .duration(200)
            .attr('stroke-width', 1)
            .attr('stroke', theme.palette.background.paper);
        }
      });

    // Labels de valeurs
    if (showValues) {
      cells.append('text')
        .attr('x', d => (xScale(d.x) || 0) + xScale.bandwidth() / 2)
        .attr('y', d => (yScale(d.y) || 0) + yScale.bandwidth() / 2)
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'middle')
        .style('font-size', `${Math.min(xScale.bandwidth(), yScale.bandwidth()) / 6}px`)
        .style('font-weight', 'bold')
        .style('fill', d => {
          const color = d3.color(colorScale(d.value));
          if (!color) return theme.palette.text.primary;
          const brightness = color.r * 0.299 + color.g * 0.587 + color.b * 0.114;
          return brightness > 128 ? '#000' : '#fff';
        })
        .style('pointer-events', 'none')
        .text(d => {
          if (isNaN(d.value)) return 'N/A';
          return selectedMetric === 'winrate' ? `${d.value.toFixed(0)}%` :
                 selectedMetric === 'kda' ? d.value.toFixed(1) :
                 selectedMetric === 'games' ? d.value.toString() :
                 d.value.toFixed(1);
        });
    }

    // Axes
    const xAxis = d3.axisBottom(xScale)
      .tickSize(0)
      .tickPadding(10);

    const yAxis = d3.axisLeft(yScale)
      .tickSize(0)
      .tickPadding(10);

    container.append('g')
      .attr('transform', `translate(0, ${chartHeight})`)
      .call(xAxis)
      .selectAll('text')
      .style('fill', theme.palette.text.primary)
      .style('font-size', '12px')
      .attr('transform', 'rotate(-45)')
      .style('text-anchor', 'end');

    container.append('g')
      .call(yAxis)
      .selectAll('text')
      .style('fill', theme.palette.text.primary)
      .style('font-size', '12px');

    // Labels des axes
    container.append('text')
      .attr('x', chartWidth / 2)
      .attr('y', chartHeight + 60)
      .attr('text-anchor', 'middle')
      .style('font-size', '14px')
      .style('font-weight', 'bold')
      .style('fill', theme.palette.text.primary)
      .text(xAxisLabel);

    container.append('text')
      .attr('x', -chartHeight / 2)
      .attr('y', -60)
      .attr('text-anchor', 'middle')
      .attr('transform', 'rotate(-90)')
      .style('font-size', '14px')
      .style('font-weight', 'bold')
      .style('fill', theme.palette.text.primary)
      .text(yAxisLabel);

    // Légende de couleur
    if (showLegend) {
      const legendWidth = 200;
      const legendHeight = 10;
      const legendX = chartWidth - legendWidth;
      const legendY = -40;

      const legendScale = d3.scaleLinear()
        .domain(processedData.extent)
        .range([0, legendWidth]);

      const legendAxis = d3.axisBottom(legendScale)
        .tickSize(4)
        .tickValues([processedData.extent[0], processedData.extent[1]])
        .tickFormat(d => selectedMetric === 'winrate' ? `${d}%` : d.toString());

      // Gradient pour la légende
      const legendGradient = defs.append('linearGradient')
        .attr('id', 'legend-gradient')
        .attr('x1', '0%').attr('y1', '0%')
        .attr('x2', '100%').attr('y2', '0%');

      const steps = 20;
      for (let i = 0; i <= steps; i++) {
        const value = processedData.extent[0] + (processedData.extent[1] - processedData.extent[0]) * (i / steps);
        legendGradient.append('stop')
          .attr('offset', `${(i / steps) * 100}%`)
          .attr('stop-color', colorScale(value));
      }

      const legend = container.append('g')
        .attr('transform', `translate(${legendX}, ${legendY})`);

      legend.append('rect')
        .attr('width', legendWidth)
        .attr('height', legendHeight)
        .attr('fill', 'url(#legend-gradient)')
        .attr('stroke', theme.palette.divider);

      legend.append('g')
        .attr('transform', `translate(0, ${legendHeight})`)
        .call(legendAxis)
        .selectAll('text')
        .style('fill', theme.palette.text.primary)
        .style('font-size', '10px');
    }

  }, [processedData, selectedColorScheme, showValues, showLegend, theme, width, height, interactive, selectedMetric, onCellClick, onCellHover]);

  // Métriques disponibles
  const metrics = [
    { value: 'winrate', label: 'Winrate (%)' },
    { value: 'kda', label: 'KDA Moyen' },
    { value: 'games', label: 'Nombre de parties' },
    { value: 'performance', label: 'Score de performance' },
  ];

  // Schémas de couleurs
  const colorSchemes = [
    { value: 'blues', label: 'Bleus' },
    { value: 'reds', label: 'Rouges' },
    { value: 'greens', label: 'Verts' },
    { value: 'oranges', label: 'Oranges' },
    { value: 'purples', label: 'Violets' },
    { value: 'rainbow', label: 'Arc-en-ciel' },
  ];

  return (
    <Card>
      <CardHeader
        title={
          <Box display="flex" alignItems="center" gap={1}>
            <TrendingUp color="primary" />
            <Typography variant="h6">{title}</Typography>
            <Chip 
              label={`${data.length} points`} 
              size="small" 
              color="primary" 
              variant="outlined" 
            />
          </Box>
        }
        action={
          <Box display="flex" gap={1}>
            <Tooltip title="Zoom avant">
              <IconButton onClick={() => setZoomLevel(prev => Math.min(prev * 1.2, 3))}>
                <ZoomIn />
              </IconButton>
            </Tooltip>
            <Tooltip title="Zoom arrière">
              <IconButton onClick={() => setZoomLevel(prev => Math.max(prev / 1.2, 0.5))}>
                <ZoomOut />
              </IconButton>
            </Tooltip>
            <Tooltip title="Exporter">
              <IconButton>
                <Download />
              </IconButton>
            </Tooltip>
          </Box>
        }
      />
      
      <CardContent>
        {/* Contrôles */}
        <Grid container spacing={2} sx={{ mb: 3 }} alignItems="center">
          <Grid item xs={12} sm={6} md={3}>
            <FormControl fullWidth size="small">
              <InputLabel>Métrique</InputLabel>
              <Select
                value={selectedMetric}
                label="Métrique"
                onChange={(e) => setSelectedMetric(e.target.value as any)}
              >
                {metrics.map(metric => (
                  <MenuItem key={metric.value} value={metric.value}>
                    {metric.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} sm={6} md={3}>
            <FormControl fullWidth size="small">
              <InputLabel>Couleurs</InputLabel>
              <Select
                value={selectedColorScheme}
                label="Couleurs"
                onChange={(e) => setSelectedColorScheme(e.target.value as any)}
              >
                {colorSchemes.map(scheme => (
                  <MenuItem key={scheme.value} value={scheme.value}>
                    {scheme.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} sm={6} md={3}>
            <FormControlLabel
              control={
                <Switch
                  checked={showValues}
                  onChange={(e) => setShowValues(e.target.checked)}
                />
              }
              label="Afficher valeurs"
            />
          </Grid>
          
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="body2" color="text.secondary">
              Zoom: {(zoomLevel * 100).toFixed(0)}%
            </Typography>
          </Grid>
        </Grid>

        {/* Heatmap */}
        <Box 
          sx={{ 
            overflow: 'auto',
            display: 'flex',
            justifyContent: 'center',
            transform: `scale(${zoomLevel})`,
            transformOrigin: 'center top',
            transition: 'transform 0.3s ease',
          }}
        >
          <svg ref={svgRef} />
        </Box>

        {/* Tooltip d'information */}
        {hoveredCell && showTooltip && (
          <Paper 
            elevation={3} 
            sx={{ 
              p: 2, 
              mt: 2, 
              bgcolor: alpha(theme.palette.background.paper, 0.95),
              border: 1,
              borderColor: 'divider',
            }}
          >
            <Typography variant="h6" gutterBottom>
              {hoveredCell.x} - {hoveredCell.y}
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="body2" color="text.secondary">
                  Valeur principale
                </Typography>
                <Typography variant="h5" color="primary">
                  {selectedMetric === 'winrate' ? `${hoveredCell.value.toFixed(1)}%` :
                   selectedMetric === 'kda' ? hoveredCell.value.toFixed(2) :
                   hoveredCell.value.toFixed(1)}
                </Typography>
              </Grid>
              {hoveredCell.metadata && (
                <>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Parties jouées
                    </Typography>
                    <Typography variant="h6">
                      {hoveredCell.metadata.totalGames}
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Victoires/Défaites
                    </Typography>
                    <Typography variant="body1">
                      {hoveredCell.metadata.wins}W / {hoveredCell.metadata.losses}L
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      KDA Moyen
                    </Typography>
                    <Typography variant="body1">
                      {hoveredCell.metadata.avgKDA.toFixed(2)}
                    </Typography>
                  </Grid>
                </>
              )}
            </Grid>
          </Paper>
        )}

        {/* Message si pas de données */}
        {data.length === 0 && (
          <Box textAlign="center" py={6}>
            <TrendingUp sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" gutterBottom>
              Aucune donnée disponible
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Importez des données pour voir la heatmap de performance
            </Typography>
          </Box>
        )}
      </CardContent>
    </Card>
  );
};