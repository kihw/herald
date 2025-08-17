import React, { useState } from 'react';
import {
  Box,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  IconButton,
} from '@mui/material';
import {
  Add,
  Person,
  Close,
} from '@mui/icons-material';
import { PerformanceComparator, PlayerStats } from '../components/comparator/PerformanceComparator';

// Données de joueurs simulées
const MOCK_PLAYERS: PlayerStats[] = [
  {
    id: '1',
    name: 'ProPlayer123',
    rank: 'Diamond II',
    totalGames: 156,
    winRate: 73.2,
    avgKDA: 2.4,
    avgKills: 8.2,
    avgDeaths: 4.1,
    avgAssists: 5.8,
    avgCS: 7.2,
    avgVisionScore: 32,
    avgGameDuration: 28.5,
    avgGoldPerMinute: 425,
    avgDamagePerMinute: 580,
    topChampions: [
      { champion: 'Jinx', games: 42, winRate: 78.6, avgKDA: 2.8 },
      { champion: 'Caitlyn', games: 38, winRate: 71.1, avgKDA: 2.3 },
      { champion: 'Vayne', games: 28, winRate: 67.9, avgKDA: 2.1 },
    ],
    performanceByRole: {
      'ADC': { games: 108, winRate: 75.0, avgKDA: 2.5 },
      'Mid': { games: 32, winRate: 68.8, avgKDA: 2.2 },
      'Support': { games: 16, winRate: 68.8, avgKDA: 2.0 },
    },
    monthlyProgress: [
      { month: 'Jan', winRate: 65, avgKDA: 2.0, games: 28 },
      { month: 'Feb', winRate: 68, avgKDA: 2.1, games: 32 },
      { month: 'Mar', winRate: 71, avgKDA: 2.3, games: 35 },
      { month: 'Apr', winRate: 73, avgKDA: 2.4, games: 38 },
      { month: 'May', winRate: 75, avgKDA: 2.5, games: 23 },
    ],
  },
  {
    id: '2',
    name: 'CasualGamer',
    rank: 'Gold I',
    totalGames: 89,
    winRate: 58.4,
    avgKDA: 1.8,
    avgKills: 6.5,
    avgDeaths: 5.2,
    avgAssists: 7.8,
    avgCS: 5.8,
    avgVisionScore: 28,
    avgGameDuration: 31.2,
    avgGoldPerMinute: 380,
    avgDamagePerMinute: 520,
    topChampions: [
      { champion: 'Thresh', games: 35, winRate: 62.9, avgKDA: 2.1 },
      { champion: 'Leona', games: 28, winRate: 57.1, avgKDA: 1.8 },
      { champion: 'Morgana', games: 16, winRate: 50.0, avgKDA: 1.5 },
    ],
    performanceByRole: {
      'Support': { games: 79, winRate: 59.5, avgKDA: 1.9 },
      'Jungle': { games: 10, winRate: 50.0, avgKDA: 1.4 },
    },
    monthlyProgress: [
      { month: 'Jan', winRate: 52, avgKDA: 1.5, games: 18 },
      { month: 'Feb', winRate: 55, avgKDA: 1.6, games: 20 },
      { month: 'Mar', winRate: 57, avgKDA: 1.7, games: 22 },
      { month: 'Apr', winRate: 58, avgKDA: 1.8, games: 18 },
      { month: 'May', winRate: 61, avgKDA: 1.9, games: 11 },
    ],
  },
  {
    id: '3',
    name: 'MidLaneKing',
    rank: 'Platinum III',
    totalGames: 127,
    winRate: 64.6,
    avgKDA: 2.1,
    avgKills: 9.1,
    avgDeaths: 5.8,
    avgAssists: 3.1,
    avgCS: 8.1,
    avgVisionScore: 18,
    avgGameDuration: 26.8,
    avgGoldPerMinute: 445,
    avgDamagePerMinute: 640,
    topChampions: [
      { champion: 'Yasuo', games: 45, winRate: 68.9, avgKDA: 2.3 },
      { champion: 'Zed', games: 38, winRate: 63.2, avgKDA: 2.0 },
      { champion: 'Ahri', games: 28, winRate: 60.7, avgKDA: 1.9 },
    ],
    performanceByRole: {
      'Mid': { games: 111, winRate: 65.8, avgKDA: 2.2 },
      'Top': { games: 16, winRate: 56.3, avgKDA: 1.7 },
    },
    monthlyProgress: [
      { month: 'Jan', winRate: 58, avgKDA: 1.8, games: 25 },
      { month: 'Feb', winRate: 61, avgKDA: 1.9, games: 28 },
      { month: 'Mar', winRate: 63, avgKDA: 2.0, games: 30 },
      { month: 'Apr', winRate: 65, avgKDA: 2.1, games: 27 },
      { month: 'May', winRate: 67, avgKDA: 2.2, games: 17 },
    ],
  },
  {
    id: '4',
    name: 'TopLaneBeast',
    rank: 'Platinum I',
    totalGames: 98,
    winRate: 69.4,
    avgKDA: 2.3,
    avgKills: 7.8,
    avgDeaths: 4.5,
    avgAssists: 2.6,
    avgCS: 7.9,
    avgVisionScore: 15,
    avgGameDuration: 29.1,
    avgGoldPerMinute: 420,
    avgDamagePerMinute: 590,
    topChampions: [
      { champion: 'Darius', games: 32, winRate: 75.0, avgKDA: 2.6 },
      { champion: 'Garen', games: 28, winRate: 67.9, avgKDA: 2.2 },
      { champion: 'Fiora', games: 22, winRate: 63.6, avgKDA: 2.0 },
    ],
    performanceByRole: {
      'Top': { games: 82, winRate: 71.0, avgKDA: 2.4 },
      'Jungle': { games: 16, winRate: 62.5, avgKDA: 1.9 },
    },
    monthlyProgress: [
      { month: 'Jan', winRate: 62, avgKDA: 2.0, games: 20 },
      { month: 'Feb', winRate: 65, avgKDA: 2.1, games: 22 },
      { month: 'Mar', winRate: 68, avgKDA: 2.2, games: 24 },
      { month: 'Apr', winRate: 70, avgKDA: 2.3, games: 18 },
      { month: 'May', winRate: 72, avgKDA: 2.4, games: 14 },
    ],
  },
];

export const ComparatorView: React.FC = () => {
  const [selectedPlayers, setSelectedPlayers] = useState<PlayerStats[]>([
    MOCK_PLAYERS[0],
    MOCK_PLAYERS[1],
  ]);
  const [showPlayerDialog, setShowPlayerDialog] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  // Filtrage des joueurs disponibles
  const availablePlayers = MOCK_PLAYERS.filter(
    player => !selectedPlayers.some(selected => selected.id === player.id)
  ).filter(
    player => player.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // Ajout d'un joueur
  const handlePlayerAdd = (player?: PlayerStats) => {
    if (player) {
      setSelectedPlayers(prev => [...prev, player]);
      setShowPlayerDialog(false);
      setSearchQuery('');
    } else {
      setShowPlayerDialog(true);
    }
  };

  // Suppression d'un joueur
  const handlePlayerRemove = (playerId: string) => {
    setSelectedPlayers(prev => prev.filter(p => p.id !== playerId));
  };

  return (
    <Box>
      {/* En-tête */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Comparateur de Performance
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Comparez les statistiques de performance entre différents joueurs
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={() => handlePlayerAdd()}
          disabled={selectedPlayers.length >= 4}
        >
          Ajouter un joueur
        </Button>
      </Box>

      {/* Information sur l'utilisation */}
      {selectedPlayers.length === 0 && (
        <Alert severity="info" sx={{ mb: 3 }}>
          Sélectionnez au moins 2 joueurs pour commencer la comparaison de leurs performances.
        </Alert>
      )}

      {/* Composant de comparaison */}
      <PerformanceComparator
        players={selectedPlayers}
        onPlayerAdd={() => handlePlayerAdd()}
        onPlayerRemove={handlePlayerRemove}
        maxPlayers={4}
        showCharts={true}
        showDetailedStats={true}
      />

      {/* Joueurs suggérés */}
      {selectedPlayers.length > 0 && selectedPlayers.length < 4 && (
        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Joueurs suggérés pour la comparaison
            </Typography>
            <Grid container spacing={2}>
              {availablePlayers.slice(0, 3).map(player => (
                <Grid item xs={12} sm={6} md={4} key={player.id}>
                  <Card variant="outlined" sx={{ cursor: 'pointer' }}>
                    <CardContent onClick={() => handlePlayerAdd(player)}>
                      <Box display="flex" alignItems="center" gap={2}>
                        <Avatar sx={{ bgcolor: 'primary.main' }}>
                          {player.name[0]}
                        </Avatar>
                        <Box flex={1}>
                          <Typography variant="subtitle1">{player.name}</Typography>
                          <Typography variant="body2" color="text.secondary">
                            {player.rank} • {player.winRate.toFixed(1)}% WR
                          </Typography>
                        </Box>
                        <Button size="small" variant="outlined">
                          Ajouter
                        </Button>
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </CardContent>
        </Card>
      )}

      {/* Dialog de sélection de joueur */}
      <Dialog 
        open={showPlayerDialog} 
        onClose={() => setShowPlayerDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          Ajouter un joueur à la comparaison
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mb: 2 }}>
            <TextField
              fullWidth
              placeholder="Rechercher un joueur..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              size="small"
            />
          </Box>
          <List>
            {availablePlayers.map(player => (
              <ListItem
                key={player.id}
                button
                onClick={() => handlePlayerAdd(player)}
              >
                <ListItemAvatar>
                  <Avatar sx={{ bgcolor: 'primary.main' }}>
                    {player.name[0]}
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={player.name}
                  secondary={`${player.rank} • ${player.totalGames} parties • ${player.winRate.toFixed(1)}% WR`}
                />
              </ListItem>
            ))}
            {availablePlayers.length === 0 && (
              <ListItem>
                <ListItemText
                  primary="Aucun joueur trouvé"
                  secondary="Essayez un autre terme de recherche"
                />
              </ListItem>
            )}
          </List>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowPlayerDialog(false)}>
            Annuler
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};