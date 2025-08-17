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

interface Match {
  id: string;
  gameMode: string;
  championName: string;
  result: 'WIN' | 'LOSS';
  kda: string;
  duration: string;
  date: string;
  rank?: string;
}

interface MatchesTableProps {
  matches: Match[];
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

  const getResultColor = (result: string) => {
    return result === 'WIN' ? 'success' : 'error';
  };

  const getChampionImage = (championName: string) => {
    return `https://ddragon.leagueoflegends.com/cdn/14.17.1/img/champion/${championName}.png`;
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
                  {currentMatches.map((match) => (
                    <TableRow key={match.id} hover>
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center' }}>
                          <Avatar
                            src={getChampionImage(match.championName)}
                            sx={{ width: 32, height: 32, mr: 2 }}
                          >
                            {match.championName[0]}
                          </Avatar>
                          <Typography variant="body2">
                            {match.championName}
                          </Typography>
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={match.gameMode}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={match.result}
                          color={getResultColor(match.result)}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2" fontFamily="monospace">
                          {match.kda}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {match.duration}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {match.date}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        {match.rank && (
                          <Chip
                            label={match.rank}
                            size="small"
                            color="primary"
                            variant="outlined"
                          />
                        )}
                      </TableCell>
                      <TableCell>
                        <IconButton
                          size="small"
                          onClick={() => onViewMatch?.(match.id)}
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
