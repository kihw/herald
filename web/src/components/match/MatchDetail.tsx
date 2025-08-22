import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Avatar,
  Chip,
  Divider,
  IconButton,
} from '@mui/material';
import { Close } from '@mui/icons-material';

interface MatchDetailProps {
  open: boolean;
  onClose: () => void;
  matchData?: {
    match: {
      id: number;
      match_id: string;
      platform: string;
      game_creation: number;
      game_duration: number;
      game_mode: string | null;
      game_type: string | null;
      queue_id: number | null;
      created_at: string;
    };
    participant: {
      champion_id: number;
      champion_name: string | null;
      kills: number;
      deaths: number;
      assists: number;
      total_damage_dealt_to_champions: number;
      gold_earned: number;
      total_minions_killed: number;
      vision_score: number;
      win: boolean;
    };
  } | null;
}

const MatchDetail: React.FC<MatchDetailProps> = ({ open, onClose, matchData }) => {
  if (!matchData) return null;

  const { match, participant } = matchData;

  const formatDuration = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  const formatDate = (timestamp: number) => {
    return new Date(timestamp * 1000).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getChampionImage = (championName: string | null) => {
    if (!championName) return '';
    return `https://ddragon.leagueoflegends.com/cdn/14.21.1/img/champion/${championName}.png`;
  };

  const getGameMode = (queueId: number | null, gameMode: string | null) => {
    if (queueId === 420) return 'Ranked Solo/Duo';
    if (queueId === 440) return 'Ranked Flex';
    if (queueId === 400) return 'Normal Draft';
    if (queueId === 430) return 'Normal Blind';
    if (queueId === 450) return 'ARAM';
    return gameMode || 'Custom Game';
  };

  const calculateKDA = () => {
    const kda = participant.deaths === 0 
      ? participant.kills + participant.assists 
      : (participant.kills + participant.assists) / participant.deaths;
    return kda.toFixed(2);
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">Match Details</Typography>
          <IconButton onClick={onClose}>
            <Close />
          </IconButton>
        </Box>
      </DialogTitle>
      <DialogContent>
        <Box sx={{ py: 2 }}>
          {/* Match Overview */}
          <Card sx={{ mb: 3 }}>
            <CardContent>
              <Grid container spacing={3}>
                <Grid item xs={12} md={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <Avatar
                      src={getChampionImage(participant.champion_name)}
                      sx={{ width: 64, height: 64, mr: 2 }}
                    >
                      {participant.champion_name?.[0] || '?'}
                    </Avatar>
                    <Box>
                      <Typography variant="h5">
                        {participant.champion_name || 'Unknown Champion'}
                      </Typography>
                      <Chip
                        label={participant.win ? 'VICTORY' : 'DEFEAT'}
                        color={participant.win ? 'success' : 'error'}
                        sx={{ mt: 1 }}
                      />
                    </Box>
                  </Box>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Typography variant="body1" color="text.secondary" gutterBottom>
                    {getGameMode(match.queue_id, match.game_mode)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {formatDate(match.game_creation)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Duration: {formatDuration(match.game_duration)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ fontFamily: 'monospace' }}>
                    Match ID: {match.match_id}
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>

          {/* Player Performance */}
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Your Performance
              </Typography>
              <Divider sx={{ mb: 2 }} />
              
              <Grid container spacing={3}>
                <Grid item xs={6} md={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h4" color="primary">
                      {participant.kills}/{participant.deaths}/{participant.assists}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      KDA
                    </Typography>
                    <Typography variant="body1" color="text.secondary">
                      {calculateKDA()} ratio
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={6} md={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h4" color="primary">
                      {participant.total_damage_dealt_to_champions.toLocaleString()}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Damage to Champions
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={6} md={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h4" color="primary">
                      {participant.gold_earned.toLocaleString()}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Gold Earned
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={6} md={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h4" color="primary">
                      {participant.total_minions_killed}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      CS
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={6} md={3}>
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h4" color="primary">
                      {participant.vision_score}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Vision Score
                    </Typography>
                  </Box>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

export default MatchDetail;