import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  LinearProgress,
  Chip,
  Alert,
  Tooltip,
  IconButton,
  Switch,
  FormControlLabel,
  Paper,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Avatar,
  Badge,
  Divider
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  Speed as SpeedIcon,
  Visibility as VisibilityIcon,
  AttachMoney as GoldIcon,
  Psychology as CoachingIcon,
  Notifications as NotificationIcon,
  Settings as SettingsIcon,
  Timeline as TimelineIcon
} from '@mui/icons-material';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip, ResponsiveContainer, AreaChart, Area } from 'recharts';
import { useWebSocket, PerformanceUpdateData, MessageType, ClientPreferences } from '../../services/websocket';
import { useAuthStore } from '../../store/authStore';

interface PerformanceMonitorProps {
  showNotifications?: boolean;
  autoUpdate?: boolean;
  updateInterval?: number;
}

interface PerformanceMetric {
  name: string;
  current: number;
  average: number;
  trend: 'up' | 'down' | 'stable';
  improvement: number;
  color: string;
  icon: React.ReactNode;
  unit: string;
}

interface PerformanceHistory {
  timestamp: Date;
  kda: number;
  cs_per_minute: number;
  vision_score: number;
  damage_share: number;
  gold_efficiency: number;
}

export const PerformanceMonitor: React.FC<PerformanceMonitorProps> = ({
  showNotifications = true,
  autoUpdate = true,
  updateInterval = 30000
}) => {
  const { user } = useAuthStore();
  const { webSocket, isConnected } = useWebSocket(user?.id, user?.token);
  
  const [currentPerformance, setCurrentPerformance] = useState<PerformanceUpdateData | null>(null);
  const [performanceHistory, setPerformanceHistory] = useState<PerformanceHistory[]>([]);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);
  const [suggestions, setSuggestions] = useState<string[]>([]);
  
  // Settings
  const [preferences, setPreferences] = useState<ClientPreferences>({
    match_updates: true,
    rank_updates: true,
    friend_activity: false,
    coaching_suggestions: true,
    system_notifications: true
  });
  const [showAdvancedMetrics, setShowAdvancedMetrics] = useState(false);

  // Handle performance updates
  const handlePerformanceUpdate = useCallback((data: PerformanceUpdateData) => {
    setCurrentPerformance(data);
    setLastUpdate(new Date());
    
    // Add to history
    setPerformanceHistory(prev => {
      const newEntry: PerformanceHistory = {
        timestamp: new Date(),
        kda: data.current_kda,
        cs_per_minute: data.cs_per_minute,
        vision_score: data.vision_score,
        damage_share: data.damage_share,
        gold_efficiency: data.gold_efficiency
      };
      
      // Keep last 50 data points
      return [...prev, newEntry].slice(-50);
    });
    
    // Update suggestions
    if (data.improvement_suggestion) {
      setSuggestions(prev => {
        const newSuggestions = [data.improvement_suggestion, ...prev.slice(0, 4)];
        return Array.from(new Set(newSuggestions)); // Remove duplicates
      });
    }
    
    // Show notification if enabled
    if (showNotifications && preferences.coaching_suggestions) {
      showPerformanceNotification(data);
    }
  }, [showNotifications, preferences.coaching_suggestions]);

  // Show performance notification
  const showPerformanceNotification = (data: PerformanceUpdateData) => {
    if (!('Notification' in window) || Notification.permission !== 'granted') {
      return;
    }

    const isImprovement = data.current_kda > data.average_kda;
    const title = isImprovement ? 'Performance Improving! 📈' : 'Performance Update 📊';
    const body = `KDA: ${data.current_kda.toFixed(2)} | ${data.improvement_suggestion}`;

    new Notification(title, {
      body,
      icon: '/herald-logo.png',
      tag: 'performance-update',
      renotify: false
    });
  };

  // Update preferences
  const updatePreferences = (newPrefs: Partial<ClientPreferences>) => {
    const updated = { ...preferences, ...newPrefs };
    setPreferences(updated);
    webSocket?.updatePreferences(updated);
  };

  // Setup WebSocket listeners
  useEffect(() => {
    if (!webSocket) return;

    webSocket.on(MessageType.PERFORMANCE_UPDATE, handlePerformanceUpdate);
    
    // Subscribe to updates
    webSocket.subscribe();

    return () => {
      webSocket.off(MessageType.PERFORMANCE_UPDATE, handlePerformanceUpdate);
    };
  }, [webSocket, handlePerformanceUpdate]);

  // Auto-update interval
  useEffect(() => {
    if (!autoUpdate || !webSocket?.isConnected) return;

    const interval = setInterval(() => {
      // Trigger performance update request
      // This would typically be done through an API call
      console.log('Requesting performance update...');
    }, updateInterval);

    return () => clearInterval(interval);
  }, [autoUpdate, webSocket?.isConnected, updateInterval]);

  if (!isConnected) {
    return (
      <Alert severity="warning">
        WebSocket connection required for real-time performance monitoring
      </Alert>
    );
  }

  if (!currentPerformance) {
    return (
      <Card>
        <CardContent>
          <Box display="flex" alignItems="center" justifyContent="center" py={4}>
            <LinearProgress sx={{ width: '100%', mr: 2 }} />
            <Typography variant="body2" color="textSecondary">
              Loading performance data...
            </Typography>
          </Box>
        </CardContent>
      </Card>
    );
  }

  // Calculate metrics
  const metrics: PerformanceMetric[] = [
    {
      name: 'KDA Ratio',
      current: currentPerformance.current_kda,
      average: currentPerformance.average_kda,
      trend: getTrend(currentPerformance.current_kda, currentPerformance.average_kda),
      improvement: ((currentPerformance.current_kda - currentPerformance.average_kda) / currentPerformance.average_kda) * 100,
      color: '#4caf50',
      icon: <TrendingUpIcon />,
      unit: ''
    },
    {
      name: 'CS/min',
      current: currentPerformance.cs_per_minute,
      average: 7.5, // Static average for now
      trend: getTrend(currentPerformance.cs_per_minute, 7.5),
      improvement: ((currentPerformance.cs_per_minute - 7.5) / 7.5) * 100,
      color: '#2196f3',
      icon: <SpeedIcon />,
      unit: '/min'
    },
    {
      name: 'Vision Score',
      current: currentPerformance.vision_score,
      average: 35, // Static average
      trend: getTrend(currentPerformance.vision_score, 35),
      improvement: ((currentPerformance.vision_score - 35) / 35) * 100,
      color: '#ff9800',
      icon: <VisibilityIcon />,
      unit: ''
    },
    {
      name: 'Damage Share',
      current: currentPerformance.damage_share,
      average: 25, // Static average
      trend: getTrend(currentPerformance.damage_share, 25),
      improvement: ((currentPerformance.damage_share - 25) / 25) * 100,
      color: '#f44336',
      icon: <TimelineIcon />,
      unit: '%'
    },
    {
      name: 'Gold Efficiency',
      current: currentPerformance.gold_efficiency,
      average: 85, // Static average
      trend: getTrend(currentPerformance.gold_efficiency, 85),
      improvement: ((currentPerformance.gold_efficiency - 85) / 85) * 100,
      color: '#ffc107',
      icon: <GoldIcon />,
      unit: '%'
    }
  ];

  return (
    <Box>
      {/* Header with Settings */}
      <Card sx={{ mb: 2 }}>
        <CardContent>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Box>
              <Typography variant="h6" gutterBottom>
                Real-time Performance Monitor
              </Typography>
              {lastUpdate && (
                <Typography variant="body2" color="textSecondary">
                  Last updated: {lastUpdate.toLocaleTimeString()}
                </Typography>
              )}
            </Box>
            
            <Box display="flex" alignItems="center" gap={1}>
              <FormControlLabel
                control={
                  <Switch
                    checked={showAdvancedMetrics}
                    onChange={(e) => setShowAdvancedMetrics(e.target.checked)}
                    size="small"
                  />
                }
                label="Advanced"
              />
              <FormControlLabel
                control={
                  <Switch
                    checked={preferences.coaching_suggestions}
                    onChange={(e) => updatePreferences({ coaching_suggestions: e.target.checked })}
                    size="small"
                  />
                }
                label="Notifications"
              />
              <Badge color="success" variant="dot" invisible={!isConnected}>
                <IconButton size="small">
                  <SettingsIcon />
                </IconButton>
              </Badge>
            </Box>
          </Box>
        </CardContent>
      </Card>

      <Grid container spacing={2}>
        {/* Performance Metrics */}
        <Grid item xs={12} lg={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Current Session Metrics
              </Typography>
              
              <Grid container spacing={2}>
                {metrics.map((metric) => (
                  <Grid item xs={12} sm={6} md={4} key={metric.name}>
                    <MetricCard metric={metric} />
                  </Grid>
                ))}
              </Grid>

              {/* Performance Chart */}
              {showAdvancedMetrics && performanceHistory.length > 1 && (
                <Box mt={3}>
                  <Typography variant="subtitle1" gutterBottom>
                    Performance Trends
                  </Typography>
                  <ResponsiveContainer width="100%" height={300}>
                    <LineChart data={performanceHistory}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis 
                        dataKey="timestamp"
                        tickFormatter={(time) => new Date(time).toLocaleTimeString()}
                      />
                      <YAxis />
                      <RechartsTooltip 
                        labelFormatter={(time) => new Date(time).toLocaleString()}
                        formatter={(value: number, name: string) => [value.toFixed(2), name]}
                      />
                      <Line
                        type="monotone"
                        dataKey="kda"
                        stroke="#4caf50"
                        strokeWidth={2}
                        name="KDA"
                      />
                      <Line
                        type="monotone"
                        dataKey="cs_per_minute"
                        stroke="#2196f3"
                        strokeWidth={2}
                        name="CS/min"
                      />
                      <Line
                        type="monotone"
                        dataKey="vision_score"
                        stroke="#ff9800"
                        strokeWidth={2}
                        name="Vision Score"
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Coaching Suggestions */}
        <Grid item xs={12} lg={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                <CoachingIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
                Improvement Suggestions
              </Typography>
              
              <Paper sx={{ maxHeight: 300, overflow: 'auto' }}>
                <List dense>
                  {suggestions.length === 0 ? (
                    <ListItem>
                      <ListItemText
                        primary="No suggestions yet"
                        secondary="Play some games to get personalized tips"
                      />
                    </ListItem>
                  ) : (
                    suggestions.map((suggestion, index) => (
                      <ListItem key={index}>
                        <ListItemIcon>
                          <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.main' }}>
                            {index + 1}
                          </Avatar>
                        </ListItemIcon>
                        <ListItemText
                          primary={suggestion}
                          secondary={`Priority: ${index === 0 ? 'High' : index === 1 ? 'Medium' : 'Low'}`}
                        />
                      </ListItem>
                    ))
                  )}
                </List>
              </Paper>

              <Divider sx={{ my: 2 }} />

              <Typography variant="subtitle2" gutterBottom>
                Notification Settings
              </Typography>
              
              <FormControlLabel
                control={
                  <Switch
                    checked={preferences.match_updates}
                    onChange={(e) => updatePreferences({ match_updates: e.target.checked })}
                    size="small"
                  />
                }
                label="Match Updates"
              />
              
              <FormControlLabel
                control={
                  <Switch
                    checked={preferences.rank_updates}
                    onChange={(e) => updatePreferences({ rank_updates: e.target.checked })}
                    size="small"
                  />
                }
                label="Rank Changes"
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

// Metric Card Component
const MetricCard: React.FC<{ metric: PerformanceMetric }> = ({ metric }) => (
  <Card variant="outlined">
    <CardContent sx={{ p: 2 }}>
      <Box display="flex" justifyContent="between" alignItems="center" mb={1}>
        <Box display="flex" alignItems="center" gap={1}>
          <Avatar sx={{ bgcolor: metric.color, width: 32, height: 32 }}>
            {metric.icon}
          </Avatar>
          <Typography variant="body2" color="textSecondary">
            {metric.name}
          </Typography>
        </Box>
        
        <Tooltip title={`${metric.improvement > 0 ? '+' : ''}${metric.improvement.toFixed(1)}% vs average`}>
          {metric.trend === 'up' ? (
            <TrendingUpIcon color="success" />
          ) : metric.trend === 'down' ? (
            <TrendingDownIcon color="error" />
          ) : (
            <TrendingUpIcon color="disabled" />
          )}
        </Tooltip>
      </Box>
      
      <Typography variant="h6" component="div">
        {metric.current.toFixed(2)}{metric.unit}
      </Typography>
      
      <Box display="flex" alignItems="center" gap={1} mt={1}>
        <LinearProgress
          variant="determinate"
          value={Math.min(100, (metric.current / (metric.average * 1.5)) * 100)}
          sx={{
            flexGrow: 1,
            height: 6,
            borderRadius: 3,
            bgcolor: 'grey.200',
            '& .MuiLinearProgress-bar': {
              bgcolor: metric.color
            }
          }}
        />
        <Typography variant="caption" color="textSecondary">
          Avg: {metric.average.toFixed(1)}
        </Typography>
      </Box>
      
      <Chip
        label={`${metric.improvement > 0 ? '+' : ''}${metric.improvement.toFixed(1)}%`}
        size="small"
        color={metric.improvement > 0 ? 'success' : metric.improvement < -5 ? 'error' : 'default'}
        sx={{ mt: 1 }}
      />
    </CardContent>
  </Card>
);

// Helper function
const getTrend = (current: number, average: number): 'up' | 'down' | 'stable' => {
  const diff = ((current - average) / average) * 100;
  if (diff > 5) return 'up';
  if (diff < -5) return 'down';
  return 'stable';
};

export default PerformanceMonitor;