import React, { useState, useEffect, useRef, useMemo } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Typography,
  IconButton,
  Tooltip,
  Slider,
  Button,
  Chip,
  Avatar,
  Fade,
  Grow,
  Slide,
  Zoom,
  useTheme,
  alpha,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
} from '@mui/material';
import {
  PlayArrow,
  Pause,
  SkipNext,
  SkipPrevious,
  Speed,
  Timeline as TimelineIcon,
  ZoomIn,
  ZoomOut,
  Fullscreen,
  Settings,
} from '@mui/icons-material';
import * as d3 from 'd3';

export interface TimelineEvent {
  id: string;
  timestamp: Date;
  type: 'match' | 'rank_change' | 'achievement' | 'milestone';
  title: string;
  description?: string;
  value?: number;
  metadata?: Record<string, any>;
  color?: string;
  icon?: React.ReactNode;
  importance?: 'low' | 'medium' | 'high' | 'critical';
}

export interface InteractiveTimelineProps {
  events: TimelineEvent[];
  title?: string;
  autoPlay?: boolean;
  playbackSpeed?: number;
  showControls?: boolean;
  showFilters?: boolean;
  height?: number;
  onEventClick?: (event: TimelineEvent) => void;
  onTimeRangeChange?: (start: Date, end: Date) => void;
}

type AnimationStyle = 'fade' | 'slide' | 'zoom' | 'bounce';
type ViewMode = 'detailed' | 'compact' | 'minimal';

export const InteractiveTimeline: React.FC<InteractiveTimelineProps> = ({
  events,
  title = 'Progression Timeline',
  autoPlay = false,
  playbackSpeed = 1000,
  showControls = true,
  showFilters = true,
  height = 400,
  onEventClick,
  onTimeRangeChange,
}) => {
  const theme = useTheme();
  const svgRef = useRef<SVGSVGElement>(null);
  const [isPlaying, setIsPlaying] = useState(autoPlay);
  const [currentTime, setCurrentTime] = useState(0);
  const [playSpeed, setPlaySpeed] = useState(playbackSpeed);
  const [selectedEventTypes, setSelectedEventTypes] = useState<string[]>(['match', 'rank_change', 'achievement', 'milestone']);
  const [animationStyle, setAnimationStyle] = useState<AnimationStyle>('fade');
  const [viewMode, setViewMode] = useState<ViewMode>('detailed');
  const [zoomLevel, setZoomLevel] = useState(1);
  const [showAnimations, setShowAnimations] = useState(true);
  const [highlightedEvents, setHighlightedEvents] = useState<Set<string>>(new Set());

  // Filtrage des événements
  const filteredEvents = useMemo(() => {
    return events
      .filter(event => selectedEventTypes.includes(event.type))
      .sort((a, b) => a.timestamp.getTime() - b.timestamp.getTime());
  }, [events, selectedEventTypes]);

  // Configuration D3
  const margin = { top: 50, right: 50, bottom: 50, left: 50 };
  const width = 800;
  const chartHeight = height - margin.top - margin.bottom;

  // Rendu de la timeline D3
  useEffect(() => {
    if (!svgRef.current || filteredEvents.length === 0) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    const container = svg
      .attr('width', width)
      .attr('height', height)
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Échelles
    const timeExtent = d3.extent(filteredEvents, d => d.timestamp) as [Date, Date];
    const xScale = d3.scaleTime()
      .domain(timeExtent)
      .range([0, width - margin.left - margin.right]);

    const yScale = d3.scaleBand()
      .domain(filteredEvents.map(d => d.type))
      .range([0, chartHeight])
      .padding(0.2);

    // Ligne principale de la timeline
    container
      .append('line')
      .attr('x1', 0)
      .attr('x2', width - margin.left - margin.right)
      .attr('y1', chartHeight / 2)
      .attr('y2', chartHeight / 2)
      .attr('stroke', theme.palette.divider)
      .attr('stroke-width', 2);

    // Gradient pour les événements
    const gradient = svg.append('defs')
      .append('linearGradient')
      .attr('id', 'event-gradient')
      .attr('gradientUnits', 'userSpaceOnUse')
      .attr('x1', '0%').attr('y1', '0%')
      .attr('x2', '100%').attr('y2', '0%');

    gradient.append('stop')
      .attr('offset', '0%')
      .attr('stop-color', theme.palette.primary.main)
      .attr('stop-opacity', 0.8);

    gradient.append('stop')
      .attr('offset', '100%')
      .attr('stop-color', theme.palette.secondary.main)
      .attr('stop-opacity', 0.3);

    // Événements
    const eventGroups = container
      .selectAll('.event-group')
      .data(filteredEvents)
      .enter()
      .append('g')
      .attr('class', 'event-group')
      .attr('transform', d => `translate(${xScale(d.timestamp)}, ${chartHeight / 2})`);

    // Cercles des événements
    eventGroups
      .append('circle')
      .attr('r', d => {
        switch (d.importance) {
          case 'critical': return 12;
          case 'high': return 10;
          case 'medium': return 8;
          default: return 6;
        }
      })
      .attr('fill', d => d.color || theme.palette.primary.main)
      .attr('stroke', theme.palette.background.paper)
      .attr('stroke-width', 2)
      .attr('opacity', 0)
      .style('cursor', 'pointer')
      .on('click', (event, d) => {
        setHighlightedEvents(new Set([d.id]));
        onEventClick?.(d);
      })
      .on('mouseover', function(event, d) {
        d3.select(this)
          .transition()
          .duration(200)
          .attr('r', d => {
            const baseR = d.importance === 'critical' ? 12 : d.importance === 'high' ? 10 : d.importance === 'medium' ? 8 : 6;
            return baseR * 1.5;
          });
        
        // Tooltip
        const tooltip = container
          .append('g')
          .attr('class', 'tooltip')
          .attr('transform', `translate(${xScale(d.timestamp)}, ${chartHeight / 2 - 40})`);

        const rect = tooltip
          .append('rect')
          .attr('x', -60)
          .attr('y', -25)
          .attr('width', 120)
          .attr('height', 40)
          .attr('rx', 5)
          .attr('fill', alpha(theme.palette.background.paper, 0.95))
          .attr('stroke', theme.palette.divider);

        tooltip
          .append('text')
          .attr('text-anchor', 'middle')
          .attr('y', -10)
          .style('font-size', '12px')
          .style('fill', theme.palette.text.primary)
          .text(d.title);

        tooltip
          .append('text')
          .attr('text-anchor', 'middle')
          .attr('y', 5)
          .style('font-size', '10px')
          .style('fill', theme.palette.text.secondary)
          .text(d.timestamp.toLocaleDateString());
      })
      .on('mouseout', function(event, d) {
        d3.select(this)
          .transition()
          .duration(200)
          .attr('r', d => {
            switch (d.importance) {
              case 'critical': return 12;
              case 'high': return 10;
              case 'medium': return 8;
              default: return 6;
            }
          });
        
        container.selectAll('.tooltip').remove();
      });

    // Lignes de connexion pour les événements importants
    eventGroups
      .filter(d => d.importance === 'critical' || d.importance === 'high')
      .append('line')
      .attr('x1', 0)
      .attr('x2', 0)
      .attr('y1', -20)
      .attr('y2', 20)
      .attr('stroke', d => d.color || theme.palette.primary.main)
      .attr('stroke-width', 2)
      .attr('opacity', 0);

    // Labels pour les événements importants
    eventGroups
      .filter(d => d.importance === 'critical' || d.importance === 'high')
      .append('text')
      .attr('text-anchor', 'middle')
      .attr('y', -25)
      .style('font-size', '10px')
      .style('fill', theme.palette.text.primary)
      .style('font-weight', 'bold')
      .text(d => d.title.length > 15 ? d.title.substring(0, 15) + '...' : d.title)
      .attr('opacity', 0);

    // Axe temporel
    const xAxis = d3.axisBottom(xScale)
      .tickFormat(d3.timeFormat('%d/%m'))
      .ticks(6);

    container
      .append('g')
      .attr('transform', `translate(0, ${chartHeight + 10})`)
      .call(xAxis)
      .selectAll('text')
      .style('fill', theme.palette.text.secondary);

    // Animation d'apparition
    if (showAnimations) {
      eventGroups.selectAll('circle')
        .transition()
        .delay((d, i) => i * 100)
        .duration(600)
        .ease(d3.easeBackOut.overshoot(1.7))
        .attr('opacity', 1);

      eventGroups.selectAll('line')
        .transition()
        .delay((d, i) => i * 100 + 300)
        .duration(400)
        .attr('opacity', 0.6);

      eventGroups.selectAll('text')
        .transition()
        .delay((d, i) => i * 100 + 500)
        .duration(400)
        .attr('opacity', 1);
    } else {
      eventGroups.selectAll('circle').attr('opacity', 1);
      eventGroups.selectAll('line').attr('opacity', 0.6);
      eventGroups.selectAll('text').attr('opacity', 1);
    }

  }, [filteredEvents, theme, height, width, chartHeight, showAnimations, onEventClick]);

  // Animation de lecture automatique
  useEffect(() => {
    if (!isPlaying || filteredEvents.length === 0) return;

    const interval = setInterval(() => {
      setCurrentTime(prev => {
        const nextTime = prev + 1;
        if (nextTime >= filteredEvents.length) {
          setIsPlaying(false);
          return 0;
        }
        
        // Mettre en surbrillance l'événement actuel
        setHighlightedEvents(new Set([filteredEvents[nextTime].id]));
        
        return nextTime;
      });
    }, playSpeed);

    return () => clearInterval(interval);
  }, [isPlaying, playSpeed, filteredEvents]);

  // Contrôles de lecture
  const handlePlay = () => setIsPlaying(true);
  const handlePause = () => setIsPlaying(false);
  const handleNext = () => {
    setCurrentTime(prev => Math.min(prev + 1, filteredEvents.length - 1));
  };
  const handlePrevious = () => {
    setCurrentTime(prev => Math.max(prev - 1, 0));
  };

  // Types d'événements uniques
  const eventTypes = Array.from(new Set(events.map(e => e.type)));

  return (
    <Card>
      <CardHeader
        title={
          <Box display="flex" alignItems="center" gap={1}>
            <TimelineIcon color="primary" />
            <Typography variant="h6">{title}</Typography>
            <Chip label={`${filteredEvents.length} événements`} size="small" />
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
            <Tooltip title="Plein écran">
              <IconButton>
                <Fullscreen />
              </IconButton>
            </Tooltip>
          </Box>
        }
      />
      
      <CardContent>
        {/* Contrôles et filtres */}
        {(showControls || showFilters) && (
          <Box sx={{ mb: 3 }}>
            {showControls && (
              <Box display="flex" alignItems="center" gap={2} sx={{ mb: 2 }}>
                <Box display="flex" gap={1}>
                  <IconButton onClick={handlePrevious} disabled={currentTime === 0}>
                    <SkipPrevious />
                  </IconButton>
                  <IconButton onClick={isPlaying ? handlePause : handlePlay}>
                    {isPlaying ? <Pause /> : <PlayArrow />}
                  </IconButton>
                  <IconButton onClick={handleNext} disabled={currentTime >= filteredEvents.length - 1}>
                    <SkipNext />
                  </IconButton>
                </Box>
                
                <Box flex={1} sx={{ mx: 2 }}>
                  <Slider
                    value={currentTime}
                    onChange={(_, value) => setCurrentTime(value as number)}
                    min={0}
                    max={Math.max(0, filteredEvents.length - 1)}
                    step={1}
                    marks
                    valueLabelDisplay="auto"
                    valueLabelFormat={(value) => 
                      filteredEvents[value]?.timestamp.toLocaleDateString() || ''
                    }
                  />
                </Box>
                
                <Box display="flex" alignItems="center" gap={1}>
                  <Speed />
                  <Slider
                    value={playSpeed}
                    onChange={(_, value) => setPlaySpeed(value as number)}
                    min={200}
                    max={2000}
                    step={200}
                    sx={{ width: 100 }}
                    valueLabelDisplay="auto"
                    valueLabelFormat={(value) => `${value}ms`}
                  />
                </Box>
              </Box>
            )}
            
            {showFilters && (
              <Box display="flex" flexWrap="wrap" gap={2} alignItems="center">
                <FormControl size="small" sx={{ minWidth: 120 }}>
                  <InputLabel>Vue</InputLabel>
                  <Select
                    value={viewMode}
                    label="Vue"
                    onChange={(e) => setViewMode(e.target.value as ViewMode)}
                  >
                    <MenuItem value="detailed">Détaillée</MenuItem>
                    <MenuItem value="compact">Compacte</MenuItem>
                    <MenuItem value="minimal">Minimale</MenuItem>
                  </Select>
                </FormControl>
                
                <FormControl size="small" sx={{ minWidth: 120 }}>
                  <InputLabel>Animation</InputLabel>
                  <Select
                    value={animationStyle}
                    label="Animation"
                    onChange={(e) => setAnimationStyle(e.target.value as AnimationStyle)}
                  >
                    <MenuItem value="fade">Fondu</MenuItem>
                    <MenuItem value="slide">Glissement</MenuItem>
                    <MenuItem value="zoom">Zoom</MenuItem>
                    <MenuItem value="bounce">Rebond</MenuItem>
                  </Select>
                </FormControl>
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={showAnimations}
                      onChange={(e) => setShowAnimations(e.target.checked)}
                    />
                  }
                  label="Animations"
                />
                
                <Box display="flex" flexWrap="wrap" gap={1}>
                  {eventTypes.map(type => (
                    <Chip
                      key={type}
                      label={type}
                      size="small"
                      clickable
                      color={selectedEventTypes.includes(type) ? 'primary' : 'default'}
                      onClick={() => {
                        setSelectedEventTypes(prev =>
                          prev.includes(type)
                            ? prev.filter(t => t !== type)
                            : [...prev, type]
                        );
                      }}
                    />
                  ))}
                </Box>
              </Box>
            )}
          </Box>
        )}

        {/* Timeline SVG */}
        <Box 
          sx={{ 
            overflow: 'auto',
            transform: `scale(${zoomLevel})`,
            transformOrigin: 'top left',
            transition: 'transform 0.3s ease',
          }}
        >
          <svg ref={svgRef} />
        </Box>

        {/* Événement actuel */}
        {filteredEvents[currentTime] && (
          <Fade in={true} key={currentTime}>
            <Box sx={{ mt: 2, p: 2, bgcolor: 'background.paper', borderRadius: 1, border: 1, borderColor: 'divider' }}>
              <Typography variant="h6" gutterBottom>
                {filteredEvents[currentTime].title}
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                {filteredEvents[currentTime].timestamp.toLocaleString()}
              </Typography>
              {filteredEvents[currentTime].description && (
                <Typography variant="body2">
                  {filteredEvents[currentTime].description}
                </Typography>
              )}
            </Box>
          </Fade>
        )}
      </CardContent>
    </Card>
  );
};