import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Avatar,
  LinearProgress,
  Chip,
  IconButton,
  Tooltip,
  Badge,
  Alert,
  Snackbar,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Paper,
  Switch,
  FormControlLabel
} from '@mui/material';
import {
  PlayArrow as PlayIcon,
  Pause as PauseIcon,
  Visibility as WatchIcon,
  VisibilityOff as UnwatchIcon,
  Notifications as NotificationIcon,
  TrendingUp as TrendingUpIcon,
  Timeline as TimelineIcon
} from '@mui/icons-material';
import { useWebSocket, MatchUpdateData, ParticipantData, MessageType } from '../../services/websocket';
import { useAuthStore } from '../../store/authStore';
import { formatGameTime, formatGold, getChampionImage, getItemImage } from '../../utils/gameUtils';

interface LiveMatchTrackerProps {
  gameId?: string;
  autoWatch?: boolean;
  showNotifications?: boolean;
}

interface LiveMatchState {
  isWatching: boolean;
  matchData: MatchUpdateData | null;
  lastUpdate: Date | null;
  updateCount: number;
  isLive: boolean;
  notifications: Array<{
    id: string;
    type: 'kill' | 'death' | 'assist' | 'level' | 'item' | 'objective';
    message: string;
    timestamp: Date;
    priority: 'high' | 'medium' | 'low';
  }>;
}

export const LiveMatchTracker: React.FC<LiveMatchTrackerProps> = ({
  gameId,
  autoWatch = true,
  showNotifications = true
}) => {
  const { user } = useAuthStore();
  const { webSocket, isConnected } = useWebSocket(user?.id, user?.token);
  
  const [matchState, setMatchState] = useState<LiveMatchState>({
    isWatching: false,
    matchData: null,
    lastUpdate: null,
    updateCount: 0,
    isLive: false,
    notifications: []
  });
  
  const [showAlert, setShowAlert] = useState(false);
  const [alertMessage, setAlertMessage] = useState('');
  const [enableSoundNotifications, setEnableSoundNotifications] = useState(true);

  // Start watching match
  const startWatching = useCallback((matchId: string) => {
    if (!webSocket || !isConnected) return;
    
    webSocket.watchMatch(matchId);
    setMatchState(prev => ({ ...prev, isWatching: true, isLive: true }));
    
    if (showNotifications) {
      setAlertMessage('Started watching live match');
      setShowAlert(true);
    }
  }, [webSocket, isConnected, showNotifications]);

  // Stop watching match
  const stopWatching = useCallback((matchId: string) => {
    if (!webSocket || !isConnected) return;
    
    webSocket.unwatchMatch(matchId);
    setMatchState(prev => ({ ...prev, isWatching: false, isLive: false }));
    
    if (showNotifications) {
      setAlertMessage('Stopped watching match');
      setShowAlert(true);
    }
  }, [webSocket, isConnected, showNotifications]);

  // Handle match updates
  const handleMatchUpdate = useCallback((data: MatchUpdateData) => {
    setMatchState(prev => {
      const newNotifications = generateUpdateNotifications(prev.matchData, data, prev.notifications);
      
      return {
        ...prev,
        matchData: data,
        lastUpdate: new Date(),
        updateCount: prev.updateCount + 1,
        notifications: newNotifications
      };
    });

    // Play notification sound if enabled
    if (enableSoundNotifications && showNotifications) {
      playNotificationSound('match_update');
    }
  }, [enableSoundNotifications, showNotifications]);

  // Generate notifications from match updates
  const generateUpdateNotifications = (
    oldData: MatchUpdateData | null,
    newData: MatchUpdateData,
    currentNotifications: LiveMatchState['notifications']
  ): LiveMatchState['notifications'] => {
    const notifications = [...currentNotifications];
    
    if (!oldData) return notifications;

    // Check for kills, deaths, assists
    newData.participants.forEach((newParticipant, index) => {
      const oldParticipant = oldData.participants[index];
      if (!oldParticipant) return;

      // Kill notification
      if (newParticipant.kills > oldParticipant.kills) {
        notifications.push({
          id: `kill-${Date.now()}-${index}`,
          type: 'kill',
          message: `${newParticipant.summoner_name} got a kill! (${newParticipant.kills}/${newParticipant.deaths}/${newParticipant.assists})`,
          timestamp: new Date(),
          priority: 'high'
        });
      }

      // Death notification
      if (newParticipant.deaths > oldParticipant.deaths) {
        notifications.push({
          id: `death-${Date.now()}-${index}`,
          type: 'death',
          message: `${newParticipant.summoner_name} died (${newParticipant.kills}/${newParticipant.deaths}/${newParticipant.assists})`,
          timestamp: new Date(),
          priority: 'medium'
        });
      }

      // Level up notification
      if (newParticipant.level > oldParticipant.level) {
        notifications.push({
          id: `level-${Date.now()}-${index}`,
          type: 'level',
          message: `${newParticipant.summoner_name} reached level ${newParticipant.level}`,
          timestamp: new Date(),
          priority: 'low'
        });
      }
    });

    // Keep only last 20 notifications
    return notifications.slice(-20);
  };

  // Play notification sounds
  const playNotificationSound = (type: string) => {
    if (!('Audio' in window)) return;
    
    try {
      const audio = new Audio(`/sounds/${type}.mp3`);
      audio.volume = 0.3;
      audio.play().catch(() => {
        // Ignore audio play errors (user interaction required)
      });
    } catch (error) {
      console.warn('Could not play notification sound:', error);
    }
  };

  // Setup WebSocket event listeners
  useEffect(() => {
    if (!webSocket) return;

    const handleMatchUpdate = (data: MatchUpdateData) => {
      // Only process updates for our watched match
      if (!gameId || data.game_id === gameId) {
        handleMatchUpdate(data);
      }
    };

    webSocket.on(MessageType.MATCH_UPDATE, handleMatchUpdate);

    return () => {
      webSocket.off(MessageType.MATCH_UPDATE, handleMatchUpdate);
    };
  }, [webSocket, gameId, handleMatchUpdate]);

  // Auto-watch match if gameId provided and autoWatch enabled
  useEffect(() => {
    if (gameId && autoWatch && isConnected && !matchState.isWatching) {
      startWatching(gameId);
    }
  }, [gameId, autoWatch, isConnected, matchState.isWatching, startWatching]);

  if (!isConnected) {
    return (
      <Alert severity="warning">
        WebSocket connection required for live match tracking
      </Alert>
    );
  }

  if (!matchState.matchData && matchState.isWatching) {
    return (
      <Card>
        <CardContent>
          <Box display="flex" alignItems="center" justifyContent="center" py={4}>
            <LinearProgress sx={{ width: '100%', mr: 2 }} />
            <Typography variant="body2" color="textSecondary">
              Waiting for match data...
            </Typography>
          </Box>
        </CardContent>
      </Card>
    );
  }

  if (!matchState.matchData) {
    return (
      <Card>
        <CardContent>
          <Box textAlign="center" py={4}>
            <Typography variant="h6" gutterBottom>
              Live Match Tracker
            </Typography>
            <Typography variant="body2" color="textSecondary" gutterBottom>
              Start watching a live match to see real-time updates
            </Typography>
            {gameId && (
              <IconButton
                color="primary"
                onClick={() => startWatching(gameId)}
                size="large"
              >
                <WatchIcon />
              </IconButton>
            )}
          </Box>
        </CardContent>
      </Card>
    );
  }

  const { matchData } = matchState;

  return (
    <Box>
      {/* Match Header */}
      <Card sx={{ mb: 2 }}>
        <CardContent>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Box>
              <Typography variant="h6" gutterBottom>
                Live Match - {formatGameTime(matchData.game_time)}
              </Typography>
              <Box display="flex" alignItems="center" gap={1}>
                <Badge
                  color={matchState.isLive ? "success" : "default"}
                  variant="dot"
                >
                  <Chip
                    label={matchData.status}
                    color={matchData.status === 'in_progress' ? 'success' : 'default'}
                    size="small"
                  />
                </Badge>
                <Typography variant="body2" color="textSecondary">
                  Updates: {matchState.updateCount}
                </Typography>
                {matchState.lastUpdate && (
                  <Typography variant="body2" color="textSecondary">
                    Last: {matchState.lastUpdate.toLocaleTimeString()}
                  </Typography>
                )}
              </Box>
            </Box>
            
            <Box display="flex" alignItems="center" gap={1}>
              <FormControlLabel
                control={
                  <Switch
                    checked={enableSoundNotifications}
                    onChange={(e) => setEnableSoundNotifications(e.target.checked)}
                    size="small"
                  />
                }
                label="Sound"
              />
              <Tooltip title={matchState.isWatching ? "Stop watching" : "Start watching"}>
                <IconButton
                  color={matchState.isWatching ? "success" : "default"}
                  onClick={() => 
                    matchState.isWatching 
                      ? stopWatching(matchData.game_id)
                      : startWatching(matchData.game_id)
                  }
                >
                  {matchState.isWatching ? <WatchIcon /> : <UnwatchIcon />}
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
        </CardContent>
      </Card>

      <Grid container spacing={2}>
        {/* Team Stats */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Team Performance
              </Typography>
              
              {/* Blue Team */}
              <Box mb={2}>
                <Typography variant="subtitle2" color="primary" gutterBottom>
                  Blue Team
                </Typography>
                <Grid container spacing={1}>
                  {matchData.participants
                    .filter(p => p.summoner_name.includes('Blue')) // This is a simplified filter
                    .map((participant, index) => (
                      <Grid item key={index}>
                        <ParticipantCard participant={participant} />
                      </Grid>
                    ))}
                </Grid>
              </Box>

              {/* Red Team */}
              <Box>
                <Typography variant="subtitle2" color="error" gutterBottom>
                  Red Team
                </Typography>
                <Grid container spacing={1}>
                  {matchData.participants
                    .filter(p => !p.summoner_name.includes('Blue')) // Simplified filter
                    .map((participant, index) => (
                      <Grid item key={index}>
                        <ParticipantCard participant={participant} />
                      </Grid>
                    ))}
                </Grid>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Live Updates & Notifications */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Live Updates
              </Typography>
              
              <Paper sx={{ maxHeight: 400, overflow: 'auto' }}>
                <List dense>
                  {matchState.notifications.length === 0 ? (
                    <ListItem>
                      <ListItemText 
                        primary="No updates yet"
                        secondary="Live events will appear here"
                      />
                    </ListItem>
                  ) : (
                    matchState.notifications
                      .slice()
                      .reverse()
                      .map((notification) => (
                        <ListItem key={notification.id}>
                          <ListItemAvatar>
                            <Avatar 
                              sx={{ 
                                bgcolor: getNotificationColor(notification.type),
                                width: 32,
                                height: 32
                              }}
                            >
                              {getNotificationIcon(notification.type)}
                            </Avatar>
                          </ListItemAvatar>
                          <ListItemText
                            primary={notification.message}
                            secondary={notification.timestamp.toLocaleTimeString()}
                          />
                        </ListItem>
                      ))
                  )}
                </List>
              </Paper>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Notification Snackbar */}
      <Snackbar
        open={showAlert}
        autoHideDuration={3000}
        onClose={() => setShowAlert(false)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={() => setShowAlert(false)} severity="info">
          {alertMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
};

// Participant Card Component
const ParticipantCard: React.FC<{ participant: ParticipantData }> = ({ participant }) => (
  <Card variant="outlined" sx={{ minWidth: 200 }}>
    <CardContent sx={{ p: 1 }}>
      <Box display="flex" alignItems="center" gap={1} mb={1}>
        <Avatar 
          src={getChampionImage(participant.champion_name)}
          sx={{ width: 32, height: 32 }}
        />
        <Box>
          <Typography variant="body2" fontWeight="bold">
            {participant.summoner_name}
          </Typography>
          <Typography variant="caption" color="textSecondary">
            {participant.champion_name} (Lv.{participant.level})
          </Typography>
        </Box>
      </Box>
      
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
        <Typography variant="body2">
          {participant.kills}/{participant.deaths}/{participant.assists}
        </Typography>
        <Chip 
          label={`${participant.kda.toFixed(2)} KDA`}
          size="small"
          color={participant.kda >= 2 ? 'success' : participant.kda >= 1 ? 'warning' : 'error'}
        />
      </Box>
      
      <Box display="flex" justifyContent="space-between" mb={1}>
        <Typography variant="caption">
          CS: {participant.cs}
        </Typography>
        <Typography variant="caption">
          Gold: {formatGold(participant.gold)}
        </Typography>
      </Box>
      
      {/* Items */}
      <Box display="flex" gap={0.5}>
        {participant.items.slice(0, 6).map((itemId, index) => (
          <Avatar
            key={index}
            src={getItemImage(itemId)}
            sx={{ width: 20, height: 20 }}
          />
        ))}
      </Box>
    </CardContent>
  </Card>
);

// Helper functions
const getNotificationColor = (type: string): string => {
  switch (type) {
    case 'kill': return '#4caf50';
    case 'death': return '#f44336';
    case 'assist': return '#2196f3';
    case 'level': return '#ff9800';
    case 'item': return '#9c27b0';
    case 'objective': return '#795548';
    default: return '#757575';
  }
};

const getNotificationIcon = (type: string): React.ReactNode => {
  switch (type) {
    case 'kill': return '⚔️';
    case 'death': return '💀';
    case 'assist': return '🤝';
    case 'level': return '📈';
    case 'item': return '🎒';
    case 'objective': return '🏆';
    default: return '📝';
  }
};

export default LiveMatchTracker;