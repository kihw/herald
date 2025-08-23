import React, { useState, useEffect, useRef, useMemo } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  ToggleButton,
  ToggleButtonGroup,
  Slider,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Grid,
  Paper,
  Tooltip,
  IconButton,
  Alert,
  CircularProgress,
  Switch,
  FormControlLabel
} from '@mui/material';
import {
  Visibility as VisionIcon,
  MyLocation as MyLocationIcon,
  FilterAlt as FilterIcon,
  Download as DownloadIcon,
  Fullscreen as FullscreenIcon,
  Settings as SettingsIcon
} from '@mui/icons-material';
import { Canvas, useFrame, useThree } from '@react-three/fiber';
import { OrthographicCamera } from '@react-three/drei';
import * as THREE from 'three';

interface HeatmapPoint {
  x: number;
  y: number;
  frequency: number;
  weight: number;
  zone: string;
}

interface HeatmapData {
  map_side: string;
  data_points: HeatmapPoint[];
  intensity: { [zone: string]: number };
  coverage: number;
}

interface VisionHeatmapProps {
  playerData: HeatmapData | null;
  wardType: 'YELLOW' | 'CONTROL' | 'BLUE_TRINKET' | 'ALL';
  timeRange: '7d' | '30d' | '90d';
  onWardTypeChange: (wardType: string) => void;
  onTimeRangeChange: (timeRange: string) => void;
  loading?: boolean;
  error?: string;
}

interface MapZone {
  name: string;
  coordinates: number[][];
  color: string;
  strategic: boolean;
}

// League of Legends map zones
const MAP_ZONES: MapZone[] = [
  {
    name: 'Dragon Pit',
    coordinates: [[9800, 4200], [10200, 4200], [10200, 4600], [9800, 4600]],
    color: '#ff6b35',
    strategic: true
  },
  {
    name: 'Baron Pit',
    coordinates: [[4800, 10200], [5200, 10200], [5200, 10600], [4800, 10600]],
    color: '#7209b7',
    strategic: true
  },
  {
    name: 'Blue Side Blue Buff',
    coordinates: [[3800, 8000], [4200, 8000], [4200, 8400], [3800, 8400]],
    color: '#4361ee',
    strategic: true
  },
  {
    name: 'Red Side Red Buff',
    coordinates: [[10800, 6600], [11200, 6600], [11200, 7000], [10800, 7000]],
    color: '#f72585',
    strategic: true
  },
  {
    name: 'River',
    coordinates: [[6000, 6000], [9000, 6000], [9000, 9000], [6000, 9000]],
    color: '#4cc9f0',
    strategic: true
  }
];

// Heatmap intensity colors
const HEAT_COLORS = [
  { value: 0, color: [0, 0, 255, 0] },      // Transparent
  { value: 0.2, color: [0, 0, 255, 100] },  // Blue
  { value: 0.4, color: [0, 255, 0, 150] },  // Green
  { value: 0.6, color: [255, 255, 0, 200] }, // Yellow
  { value: 0.8, color: [255, 165, 0, 230] }, // Orange
  { value: 1.0, color: [255, 0, 0, 255] }    // Red
];

export const VisionHeatmap: React.FC<VisionHeatmapProps> = ({
  playerData,
  wardType,
  timeRange,
  onWardTypeChange,
  onTimeRangeChange,
  loading = false,
  error
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [heatmapIntensity, setHeatmapIntensity] = useState(0.7);
  const [showZones, setShowZones] = useState(true);
  const [showStats, setShowStats] = useState(true);
  const [selectedZone, setSelectedZone] = useState<string>('');
  const [mapView, setMapView] = useState<'blue' | 'red' | 'both'>('both');

  // Process heatmap data for visualization
  const processedData = useMemo(() => {
    if (!playerData || !playerData.data_points) return null;

    const points = playerData.data_points;
    if (points.length === 0) return null;

    // Normalize coordinates to map size (14870x14870 Summoner's Rift)
    const mapWidth = 14870;
    const mapHeight = 14870;
    
    // Find max frequency for normalization
    const maxFrequency = Math.max(...points.map(p => p.frequency));
    
    return {
      points: points.map(point => ({
        ...point,
        normalizedX: (point.x / mapWidth) * 512, // Scale to canvas size
        normalizedY: (point.y / mapHeight) * 512,
        normalizedIntensity: point.frequency / maxFrequency
      })),
      maxFrequency,
      totalPoints: points.length
    };
  }, [playerData]);

  // Draw heatmap on canvas
  useEffect(() => {
    if (!canvasRef.current || !processedData) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Draw background (Summoner's Rift map would go here)
    ctx.fillStyle = '#0a1a0a';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Draw map zones if enabled
    if (showZones) {
      MAP_ZONES.forEach(zone => {
        if (zone.strategic) {
          ctx.strokeStyle = zone.color;
          ctx.lineWidth = 2;
          ctx.globalAlpha = 0.3;
          
          ctx.beginPath();
          zone.coordinates.forEach((coord, index) => {
            const x = (coord[0] / 14870) * canvas.width;
            const y = canvas.height - (coord[1] / 14870) * canvas.height; // Flip Y
            
            if (index === 0) {
              ctx.moveTo(x, y);
            } else {
              ctx.lineTo(x, y);
            }
          });
          ctx.closePath();
          ctx.stroke();
          
          // Zone label
          const centerX = zone.coordinates.reduce((sum, coord) => sum + coord[0], 0) / zone.coordinates.length;
          const centerY = zone.coordinates.reduce((sum, coord) => sum + coord[1], 0) / zone.coordinates.length;
          
          ctx.fillStyle = zone.color;
          ctx.font = '10px Arial';
          ctx.globalAlpha = 0.8;
          ctx.fillText(
            zone.name,
            (centerX / 14870) * canvas.width,
            canvas.height - (centerY / 14870) * canvas.height
          );
        }
      });
      ctx.globalAlpha = 1.0;
    }

    // Draw heatmap points
    processedData.points.forEach(point => {
      const radius = 15 + (point.weight * 10);
      const intensity = point.normalizedIntensity * heatmapIntensity;
      
      // Create radial gradient for heat effect
      const gradient = ctx.createRadialGradient(
        point.normalizedX, point.normalizedY, 0,
        point.normalizedX, point.normalizedY, radius
      );
      
      // Get color based on intensity
      const color = getHeatColor(intensity);
      gradient.addColorStop(0, `rgba(${color[0]}, ${color[1]}, ${color[2]}, ${color[3] / 255})`);
      gradient.addColorStop(1, 'rgba(0, 0, 0, 0)');
      
      ctx.fillStyle = gradient;
      ctx.beginPath();
      ctx.arc(point.normalizedX, point.normalizedY, radius, 0, 2 * Math.PI);
      ctx.fill();
      
      // Draw ward icon for high-frequency points
      if (point.frequency > 5) {
        ctx.fillStyle = '#ffffff';
        ctx.font = '12px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('👁️', point.normalizedX, point.normalizedY + 4);
      }
    });

  }, [processedData, heatmapIntensity, showZones]);

  // Get heat color based on intensity
  const getHeatColor = (intensity: number): number[] => {
    for (let i = 0; i < HEAT_COLORS.length - 1; i++) {
      const current = HEAT_COLORS[i];
      const next = HEAT_COLORS[i + 1];
      
      if (intensity >= current.value && intensity <= next.value) {
        // Interpolate between colors
        const ratio = (intensity - current.value) / (next.value - current.value);
        return [
          Math.round(current.color[0] + (next.color[0] - current.color[0]) * ratio),
          Math.round(current.color[1] + (next.color[1] - current.color[1]) * ratio),
          Math.round(current.color[2] + (next.color[2] - current.color[2]) * ratio),
          Math.round(current.color[3] + (next.color[3] - current.color[3]) * ratio)
        ];
      }
    }
    return HEAT_COLORS[HEAT_COLORS.length - 1].color;
  };

  // Download heatmap as image
  const downloadHeatmap = () => {
    if (!canvasRef.current) return;
    
    const link = document.createElement('a');
    link.download = `herald-vision-heatmap-${wardType}-${timeRange}.png`;
    link.href = canvasRef.current.toDataURL();
    link.click();
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
            <CircularProgress />
            <Typography variant="body2" sx={{ ml: 2 }}>
              Generating vision heatmap...
            </Typography>
          </Box>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  return (
    <Box>
      {/* Controls */}
      <Card sx={{ mb: 2 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            {/* Ward Type Selection */}
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Ward Type</InputLabel>
                <Select
                  value={wardType}
                  label="Ward Type"
                  onChange={(e) => onWardTypeChange(e.target.value)}
                >
                  <MenuItem value="ALL">All Wards</MenuItem>
                  <MenuItem value="YELLOW">Yellow Wards</MenuItem>
                  <MenuItem value="CONTROL">Control Wards</MenuItem>
                  <MenuItem value="BLUE_TRINKET">Blue Trinket</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            {/* Time Range */}
            <Grid item xs={12} sm={6} md={3}>
              <ToggleButtonGroup
                value={timeRange}
                exclusive
                onChange={(_, value) => value && onTimeRangeChange(value)}
                size="small"
              >
                <ToggleButton value="7d">7 Days</ToggleButton>
                <ToggleButton value="30d">30 Days</ToggleButton>
                <ToggleButton value="90d">90 Days</ToggleButton>
              </ToggleButtonGroup>
            </Grid>

            {/* Map View */}
            <Grid item xs={12} sm={6} md={3}>
              <ToggleButtonGroup
                value={mapView}
                exclusive
                onChange={(_, value) => value && setMapView(value)}
                size="small"
              >
                <ToggleButton value="blue">Blue Side</ToggleButton>
                <ToggleButton value="red">Red Side</ToggleButton>
                <ToggleButton value="both">Both</ToggleButton>
              </ToggleButtonGroup>
            </Grid>

            {/* Actions */}
            <Grid item xs={12} sm={6} md={3}>
              <Box display="flex" gap={1}>
                <Tooltip title="Download Heatmap">
                  <IconButton onClick={downloadHeatmap} size="small">
                    <DownloadIcon />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Fullscreen">
                  <IconButton size="small">
                    <FullscreenIcon />
                  </IconButton>
                </Tooltip>
              </Box>
            </Grid>
          </Grid>

          {/* Settings */}
          <Box sx={{ mt: 2 }}>
            <Grid container spacing={2} alignItems="center">
              <Grid item xs={12} sm={4}>
                <Typography variant="body2" gutterBottom>
                  Heatmap Intensity
                </Typography>
                <Slider
                  value={heatmapIntensity}
                  onChange={(_, value) => setHeatmapIntensity(value as number)}
                  min={0.1}
                  max={1.0}
                  step={0.1}
                  marks
                  valueLabelDisplay="auto"
                  size="small"
                />
              </Grid>

              <Grid item xs={12} sm={4}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={showZones}
                      onChange={(e) => setShowZones(e.target.checked)}
                      size="small"
                    />
                  }
                  label="Show Map Zones"
                />
              </Grid>

              <Grid item xs={12} sm={4}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={showStats}
                      onChange={(e) => setShowStats(e.target.checked)}
                      size="small"
                    />
                  }
                  label="Show Statistics"
                />
              </Grid>
            </Grid>
          </Box>
        </CardContent>
      </Card>

      <Grid container spacing={2}>
        {/* Heatmap Visualization */}
        <Grid item xs={12} lg={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                <VisionIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
                Vision Control Heatmap
              </Typography>
              
              <Box
                sx={{
                  position: 'relative',
                  width: '100%',
                  height: 512,
                  border: '2px solid',
                  borderColor: 'divider',
                  borderRadius: 1,
                  overflow: 'hidden'
                }}
              >
                <canvas
                  ref={canvasRef}
                  width={512}
                  height={512}
                  style={{
                    width: '100%',
                    height: '100%',
                    cursor: 'crosshair'
                  }}
                />
                
                {/* Legend */}
                <Box
                  sx={{
                    position: 'absolute',
                    top: 8,
                    right: 8,
                    bgcolor: 'rgba(0,0,0,0.7)',
                    color: 'white',
                    p: 1,
                    borderRadius: 1,
                    fontSize: '12px'
                  }}
                >
                  <Typography variant="caption" display="block">
                    Heat Intensity
                  </Typography>
                  <Box display="flex" alignItems="center" gap={0.5} mt={0.5}>
                    <Box width={10} height={10} bgcolor="rgba(0,0,255,0.6)" />
                    <Typography variant="caption">Low</Typography>
                    <Box width={10} height={10} bgcolor="rgba(255,255,0,0.8)" />
                    <Typography variant="caption">Med</Typography>
                    <Box width={10} height={10} bgcolor="rgba(255,0,0,1)" />
                    <Typography variant="caption">High</Typography>
                  </Box>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Statistics Panel */}
        {showStats && (
          <Grid item xs={12} lg={4}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Vision Statistics
                </Typography>
                
                {playerData && (
                  <Box>
                    <Paper sx={{ p: 2, mb: 2 }}>
                      <Typography variant="body2" color="textSecondary" gutterBottom>
                        Map Coverage
                      </Typography>
                      <Typography variant="h4" color="primary">
                        {playerData.coverage.toFixed(1)}%
                      </Typography>
                      <Typography variant="caption" color="textSecondary">
                        of strategic areas
                      </Typography>
                    </Paper>

                    <Typography variant="subtitle2" gutterBottom>
                      Zone Activity
                    </Typography>
                    
                    <Box sx={{ mb: 2 }}>
                      {Object.entries(playerData.intensity || {})
                        .sort((a, b) => b[1] - a[1])
                        .slice(0, 5)
                        .map(([zone, intensity]) => (
                          <Box key={zone} display="flex" justifyContent="space-between" alignItems="center" mb={0.5}>
                            <Typography variant="body2">{zone}</Typography>
                            <Chip
                              label={intensity}
                              size="small"
                              color={intensity > 10 ? 'success' : intensity > 5 ? 'warning' : 'default'}
                            />
                          </Box>
                        ))}
                    </Box>

                    {processedData && (
                      <Paper sx={{ p: 2 }}>
                        <Typography variant="body2" color="textSecondary" gutterBottom>
                          Total Ward Placements
                        </Typography>
                        <Typography variant="h5" color="primary">
                          {processedData.totalPoints}
                        </Typography>
                        <Typography variant="caption" color="textSecondary">
                          in selected time range
                        </Typography>
                      </Paper>
                    )}
                  </Box>
                )}
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>

      {/* Heat Color Legend */}
      <Card sx={{ mt: 2 }}>
        <CardContent>
          <Typography variant="subtitle2" gutterBottom>
            Understanding Your Vision Heatmap
          </Typography>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" paragraph>
                • <strong>Red/Hot areas:</strong> Frequent ward placements - your preferred vision spots
              </Typography>
              <Typography variant="body2" paragraph>
                • <strong>Yellow areas:</strong> Moderate activity - occasional warding
              </Typography>
              <Typography variant="body2" paragraph>
                • <strong>Blue/Cold areas:</strong> Low activity - consider improving coverage
              </Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" paragraph>
                • <strong>Strategic zones:</strong> High-value areas for vision control
              </Typography>
              <Typography variant="body2" paragraph>
                • <strong>Coverage percentage:</strong> How well you cover strategic areas
              </Typography>
              <Typography variant="body2" paragraph>
                • <strong>Zone activity:</strong> Your most frequently warded locations
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </Card>
    </Box>
  );
};

export default VisionHeatmap;