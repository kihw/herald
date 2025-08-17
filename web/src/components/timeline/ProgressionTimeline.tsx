import React, { useMemo } from 'react';
import {
  Box,
  Typography,
  Chip,
  Avatar,
} from '@mui/material';
import {
  EmojiEvents,
  TrendingUp,
  TrendingDown,
  Star,
  Timeline as TimelineIcon,
  SportsEsports,
} from '@mui/icons-material';
import { InteractiveTimeline, TimelineEvent } from './InteractiveTimeline';
import { Row } from '../../types';

export interface ProgressionTimelineProps {
  matches: Row[];
  puuid?: string;
  title?: string;
  showControls?: boolean;
  onEventClick?: (event: TimelineEvent) => void;
}

export const ProgressionTimeline: React.FC<ProgressionTimelineProps> = ({
  matches,
  puuid,
  title = 'Progression des Performances',
  showControls = true,
  onEventClick,
}) => {
  // Conversion des matches en événements de timeline
  const timelineEvents = useMemo((): TimelineEvent[] => {
    if (!matches.length) return [];

    const events: TimelineEvent[] = [];
    
    // Tri des matches par date
    const sortedMatches = [...matches].sort((a, b) => 
      new Date(a.game_creation).getTime() - new Date(b.game_creation).getTime()
    );

    // Calcul des statistiques de progression
    let cumulativeWins = 0;
    let totalGames = 0;
    let previousWinRate = 0;
    let bestKDA = 0;
    let currentStreak = 0;
    let bestStreak = 0;
    let streakType: 'win' | 'loss' | null = null;

    sortedMatches.forEach((match, index) => {
      totalGames++;
      if (match.win) {
        cumulativeWins++;
        if (streakType === 'win') {
          currentStreak++;
        } else {
          currentStreak = 1;
          streakType = 'win';
        }
      } else {
        if (streakType === 'loss') {
          currentStreak++;
        } else {
          currentStreak = 1;
          streakType = 'loss';
        }
      }

      const currentWinRate = (cumulativeWins / totalGames) * 100;
      const kda = match.deaths > 0 ? (match.kills + match.assists) / match.deaths : match.kills + match.assists;
      
      if (kda > bestKDA) {
        bestKDA = kda;
      }

      if (streakType === 'win' && currentStreak > bestStreak) {
        bestStreak = currentStreak;
      }

      // Événement de match standard
      events.push({
        id: `match-${match.game_id}`,
        timestamp: new Date(match.game_creation),
        type: 'match',
        title: `${match.win ? 'Victoire' : 'Défaite'} - ${match.champion}`,
        description: `${match.kills}/${match.deaths}/${match.assists} • ${match.lane} • ${Math.round(match.game_duration / 60)}min`,
        value: currentWinRate,
        color: match.win ? '#4caf50' : '#f44336',
        importance: 'low',
        metadata: {
          match,
          winRate: currentWinRate,
          kda,
          totalGames,
          streak: currentStreak,
          streakType,
        },
        icon: <SportsEsports />,
      });

      // Événements spéciaux
      
      // Changement significatif de winrate
      if (Math.abs(currentWinRate - previousWinRate) >= 5 && totalGames > 5) {
        events.push({
          id: `winrate-change-${match.game_id}`,
          timestamp: new Date(match.game_creation),
          type: 'milestone',
          title: `Winrate ${currentWinRate > previousWinRate ? 'amélioration' : 'baisse'}`,
          description: `${previousWinRate.toFixed(1)}% → ${currentWinRate.toFixed(1)}%`,
          value: currentWinRate,
          color: currentWinRate > previousWinRate ? '#2196f3' : '#ff9800',
          importance: 'medium',
          icon: currentWinRate > previousWinRate ? <TrendingUp /> : <TrendingDown />,
        });
      }

      // Streak important (5+ victoires ou défaites)
      if (currentStreak >= 5 && (
        (streakType === 'win' && currentStreak > 5) ||
        (streakType === 'loss' && currentStreak === 5)
      )) {
        events.push({
          id: `streak-${match.game_id}`,
          timestamp: new Date(match.game_creation),
          type: 'achievement',
          title: `${currentStreak} ${streakType === 'win' ? 'victoires' : 'défaites'} consécutives`,
          description: `Série ${streakType === 'win' ? 'victorieuse' : 'de défaites'} remarquable`,
          value: currentStreak,
          color: streakType === 'win' ? '#4caf50' : '#f44336',
          importance: currentStreak >= 10 ? 'critical' : 'high',
          icon: <Star />,
        });
      }

      // KDA exceptionnel (4.0+)
      if (kda >= 4.0) {
        events.push({
          id: `high-kda-${match.game_id}`,
          timestamp: new Date(match.game_creation),
          type: 'achievement',
          title: `KDA Exceptionnel: ${kda.toFixed(1)}`,
          description: `Performance remarquable avec ${match.champion}`,
          value: kda,
          color: '#9c27b0',
          importance: kda >= 6.0 ? 'critical' : 'high',
          icon: <EmojiEvents />,
        });
      }

      // Premier match avec un nouveau champion
      const previousChampions = sortedMatches.slice(0, index).map(m => m.champion);
      if (!previousChampions.includes(match.champion)) {
        events.push({
          id: `new-champion-${match.game_id}`,
          timestamp: new Date(match.game_creation),
          type: 'milestone',
          title: `Nouveau champion: ${match.champion}`,
          description: `Premier match avec ${match.champion}`,
          value: 1,
          color: '#ff5722',
          importance: 'low',
          icon: <TimelineIcon />,
        });
      }

      // Jalons de nombre de parties (10, 25, 50, 100, etc.)
      if ([10, 25, 50, 100, 200, 500].includes(totalGames)) {
        events.push({
          id: `games-milestone-${totalGames}`,
          timestamp: new Date(match.game_creation),
          type: 'milestone',
          title: `${totalGames} parties jouées`,
          description: `Jalon atteint • Winrate: ${currentWinRate.toFixed(1)}%`,
          value: totalGames,
          color: '#607d8b',
          importance: totalGames >= 100 ? 'high' : 'medium',
          icon: <EmojiEvents />,
        });
      }

      previousWinRate = currentWinRate;
    });

    return events;
  }, [matches]);

  const handleEventClick = (event: TimelineEvent) => {
    console.log('Timeline event clicked:', event);
    onEventClick?.(event);
  };

  if (matches.length === 0) {
    return (
      <Box textAlign="center" py={6}>
        <TimelineIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
        <Typography variant="h6" color="text.secondary" gutterBottom>
          Aucune donnée disponible
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Importez des données de match pour voir votre progression
        </Typography>
      </Box>
    );
  }

  // Statistiques de résumé
  const totalMatches = matches.length;
  const wins = matches.filter(m => m.win).length;
  const winRate = (wins / totalMatches) * 100;
  const avgKDA = matches.reduce((sum, m) => {
    const kda = m.deaths > 0 ? (m.kills + m.assists) / m.deaths : m.kills + m.assists;
    return sum + kda;
  }, 0) / totalMatches;

  const uniqueChampions = new Set(matches.map(m => m.champion)).size;
  const timeSpan = new Date(Math.max(...matches.map(m => new Date(m.game_creation).getTime()))) 
    .getTime() - new Date(Math.min(...matches.map(m => new Date(m.game_creation).getTime()))).getTime();
  const daysSpan = Math.ceil(timeSpan / (1000 * 60 * 60 * 24));

  return (
    <Box>
      {/* Statistiques de résumé */}
      <Box display="flex" flexWrap="wrap" gap={1} mb={3}>
        <Chip 
          label={`${totalMatches} parties`} 
          color="primary" 
          variant="outlined" 
        />
        <Chip 
          label={`${winRate.toFixed(1)}% WR`} 
          color={winRate >= 50 ? "success" : "error"} 
          variant="outlined" 
        />
        <Chip 
          label={`${avgKDA.toFixed(2)} KDA`} 
          color="secondary" 
          variant="outlined" 
        />
        <Chip 
          label={`${uniqueChampions} champions`} 
          color="info" 
          variant="outlined" 
        />
        <Chip 
          label={`${daysSpan} jours`} 
          color="default" 
          variant="outlined" 
        />
        <Chip 
          label={`${timelineEvents.length} événements`} 
          color="warning" 
          variant="outlined" 
        />
      </Box>

      {/* Timeline interactive */}
      <InteractiveTimeline
        events={timelineEvents}
        title={title}
        autoPlay={false}
        playbackSpeed={800}
        showControls={showControls}
        showFilters={true}
        height={500}
        onEventClick={handleEventClick}
      />

      {/* Légende des types d'événements */}
      <Box mt={3}>
        <Typography variant="subtitle2" gutterBottom>
          Légende des événements:
        </Typography>
        <Box display="flex" flexWrap="wrap" gap={2}>
          <Box display="flex" alignItems="center" gap={1}>
            <Box 
              sx={{ 
                width: 12, 
                height: 12, 
                borderRadius: '50%', 
                bgcolor: '#4caf50' 
              }} 
            />
            <Typography variant="caption">Victoires</Typography>
          </Box>
          <Box display="flex" alignItems="center" gap={1}>
            <Box 
              sx={{ 
                width: 12, 
                height: 12, 
                borderRadius: '50%', 
                bgcolor: '#f44336' 
              }} 
            />
            <Typography variant="caption">Défaites</Typography>
          </Box>
          <Box display="flex" alignItems="center" gap={1}>
            <Box 
              sx={{ 
                width: 12, 
                height: 12, 
                borderRadius: '50%', 
                bgcolor: '#2196f3' 
              }} 
            />
            <Typography variant="caption">Jalons</Typography>
          </Box>
          <Box display="flex" alignItems="center" gap={1}>
            <Box 
              sx={{ 
                width: 12, 
                height: 12, 
                borderRadius: '50%', 
                bgcolor: '#9c27b0' 
              }} 
            />
            <Typography variant="caption">Achievements</Typography>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};