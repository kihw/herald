import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  Avatar,
  IconButton,
  Pagination,
  Button,
  CircularProgress,
  Alert,
} from '@mui/material';
import {
  Sync,
  Visibility,
  Download,
  FilterList,
} from '@mui/icons-material';

interface MatchSummary {
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
}

interface MatchesTableProps {
  matches: MatchSummary[];
  loading?: boolean;
  onSync?: () => void;
  onViewMatch?: (matchId: string) => void;
  onExport?: () => void;
  syncLoading?: boolean;
}

const MatchesTable: React.FC<MatchesTableProps> = ({
  matches,
  loading = false,
  onSync,
  onViewMatch,
  onExport,
  syncLoading = false,
}) => {
  const [page, setPage] = useState(1);
  const matchesPerPage = 10;
  const totalPages = Math.ceil(matches.length / matchesPerPage);

  const currentMatches = matches.slice(
    (page - 1) * matchesPerPage,
    page * matchesPerPage
  );

  const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
  };

  const getResultColor = (win: boolean) => {
    return win ? 'success' : 'error';
  };

  const getChampionImage = (championName: string | null) => {
    if (!championName) return '';
    return `https://ddragon.leagueoflegends.com/cdn/14.21.1/img/champion/${championName}.png`;
  };

  const formatDuration = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  const formatDate = (timestamp: number) => {
    return new Date(timestamp * 1000).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: '2-digit'
    });
  };

  const formatKDA = (kills: number, deaths: number, assists: number) => {
    return `${kills}/${deaths}/${assists}`;
  };

  const getGameMode = (queueId: number | null, gameMode: string | null) => {
    if (queueId === 420) return 'Ranked Solo';
    if (queueId === 440) return 'Ranked Flex';
    if (queueId === 400) return 'Normal Draft';
    if (queueId === 430) return 'Normal Blind';
    if (queueId === 450) return 'ARAM';
    return gameMode || 'Unknown';
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h6" component="div">
            Match History ({matches.length} matches)
          </Typography>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              variant="outlined"
              startIcon={<FilterList />}
              size="small"
            >
              Filter
            </Button>
            <Button
              variant="outlined"
              startIcon={<Download />}
              onClick={onExport}
              size="small"
            >
              Export
            </Button>
            <Button
              variant="contained"
              startIcon={syncLoading ? <CircularProgress size={16} /> : <Sync />}
              onClick={onSync}
              disabled={syncLoading}
              size="small"
            >
              {syncLoading ? 'Syncing...' : 'Sync Matches'}
            </Button>
          </Box>
        </Box>

        {/* Matches Table */}
        {matches.length === 0 ? (
          <Alert severity="info" sx={{ mt: 2 }}>
            No matches found. Click "Sync Matches" to load your recent games.
          </Alert>
        ) : (
          <>
            <TableContainer component={Paper} variant="outlined">
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Champion</TableCell>
                    <TableCell>Game Mode</TableCell>
                    <TableCell>Result</TableCell>
                    <TableCell>KDA</TableCell>
                    <TableCell>Duration</TableCell>
                    <TableCell>Date</TableCell>
                    <TableCell>Rank</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {currentMatches.map((matchData) => (
                    <TableRow key={matchData.match.match_id} hover>
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center' }}>
                          <Avatar
                            src={getChampionImage(matchData.participant.champion_name)}
                            sx={{ width: 32, height: 32, mr: 2 }}
                          >
                            {matchData.participant.champion_name?.[0] || '?'}
                          </Avatar>
                          <Typography variant="body2">
                            {matchData.participant.champion_name || 'Unknown'}
                          </Typography>
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={getGameMode(matchData.match.queue_id, matchData.match.game_mode)}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={matchData.participant.win ? 'WIN' : 'LOSS'}
                          color={getResultColor(matchData.participant.win)}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2" fontFamily="monospace">
                          {formatKDA(
                            matchData.participant.kills,
                            matchData.participant.deaths,
                            matchData.participant.assists
                          )}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {formatDuration(matchData.match.game_duration)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {formatDate(matchData.match.game_creation)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label="Unranked"
                          size="small"
                          color="default"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>
                        <IconButton
                          size="small"
                          onClick={() => onViewMatch?.(matchData.match.match_id)}
                          title="View Details"
                        >
                          <Visibility />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>

            {/* Pagination */}
            {totalPages > 1 && (
              <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
                <Pagination
                  count={totalPages}
                  page={page}
                  onChange={handlePageChange}
                  color="primary"
                />
              </Box>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
};

export default MatchesTable;
