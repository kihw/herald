import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Grid,
  Chip,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  useTheme,
  Alert,
  CircularProgress,
  Fab,
  Tooltip,
} from '@mui/material';
import {
  Add as AddIcon,
  MoreVert as MoreVertIcon,
  CompareArrows as CompareArrowsIcon,
  TrendingUp as TrendingUpIcon,
  EmojiEvents as EmojiEventsIcon,
  Person as PersonIcon,
  SportsEsports as SportsEsportsIcon,
  Analytics as AnalyticsIcon,
  Refresh as RefreshIcon,
  Share as ShareIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import CreateComparisonDialog from './CreateComparisonDialog';
import ComparisonDetailsDialog from './ComparisonDetailsDialog';
import { groupApi, GroupComparison } from '../../services/groupApi';

interface ComparisonManagerProps {
  groupId: number;
}

const ComparisonManager: React.FC<ComparisonManagerProps> = ({ groupId }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [comparisons, setComparisons] = useState<GroupComparison[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Dialog states
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [selectedComparison, setSelectedComparison] = useState<GroupComparison | null>(null);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  
  // Menu state
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [menuComparisonId, setMenuComparisonId] = useState<number | null>(null);

  useEffect(() => {
    if (groupId) {
      loadComparisons();
    }
  }, [groupId]);

  const loadComparisons = async () => {
    try {
      setLoading(true);
      const groupComparisons = await groupApi.getGroupComparisons(groupId);
      setComparisons(groupComparisons);
      setError(null);
    } catch (err) {
      setError('Erreur lors du chargement des comparaisons');
      console.error('Error loading comparisons:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, comparisonId: number) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
    setMenuComparisonId(comparisonId);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setMenuComparisonId(null);
  };

  const handleComparisonDetails = (comparison: GroupComparison) => {
    setSelectedComparison(comparison);
    setDetailsDialogOpen(true);
    handleMenuClose();
  };

  const handleCreateComparison = async (comparisonData: any) => {
    try {
      await groupApi.createComparison(groupId, comparisonData);
      await loadComparisons();
      setCreateDialogOpen(false);
    } catch (err) {
      console.error('Error creating comparison:', err);
    }
  };

  const handleRegenerateComparison = async (comparisonId: number) => {
    try {
      await groupApi.regenerateComparison(groupId, comparisonId);
      await loadComparisons();
      handleMenuClose();
    } catch (err) {
      setError('Erreur lors de la régénération de la comparaison');
      console.error('Error regenerating comparison:', err);
    }
  };

  const getComparisonIcon = (compareType: string) => {
    switch (compareType) {
      case 'champions':
        return <EmojiEventsIcon sx={{ fontSize: 20, color: leagueColors.gold[500] }} />;
      case 'roles':
        return <SportsEsportsIcon sx={{ fontSize: 20, color: leagueColors.blue[500] }} />;
      case 'performance':
        return <TrendingUpIcon sx={{ fontSize: 20, color: leagueColors.win }} />;
      case 'trends':
        return <AnalyticsIcon sx={{ fontSize: 20, color: leagueColors.dark[400] }} />;
      default:
        return <CompareArrowsIcon sx={{ fontSize: 20, color: 'text.secondary' }} />;
    }
  };

  const getComparisonTypeLabel = (compareType: string) => {
    switch (compareType) {
      case 'champions':
        return 'Champions';
      case 'roles':
        return 'Rôles';
      case 'performance':
        return 'Performance';
      case 'trends':
        return 'Tendances';
      default:
        return 'Comparaison';
    }
  };

  const getComparisonTypeColor = (compareType: string) => {
    switch (compareType) {
      case 'champions':
        return leagueColors.gold[500];
      case 'roles':
        return leagueColors.blue[500];
      case 'performance':
        return leagueColors.win;
      case 'trends':
        return leagueColors.dark[400];
      default:
        return 'text.secondary';
    }
  };

  if (loading) {
    return (
      <Box 
        sx={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          minHeight: 300 
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <Box>
          <Typography 
            variant="h5" 
            sx={{ 
              fontWeight: 700,
              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              mb: 1,
            }}
          >
            Comparaisons de Groupe
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Analysez et comparez les performances des membres du groupe
          </Typography>
        </Box>
        
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setCreateDialogOpen(true)}
          sx={{
            background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
            '&:hover': {
              background: `linear-gradient(135deg, ${leagueColors.blue[600]} 0%, ${leagueColors.blue[700]} 100%)`,
            },
          }}
        >
          Créer une Comparaison
        </Button>
      </Box>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Comparisons Grid */}
      {comparisons.length === 0 ? (
        <Card 
          sx={{ 
            textAlign: 'center', 
            py: 6,
            background: isDarkMode
              ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
              : `linear-gradient(135deg, ${leagueColors.blue[25]} 0%, #ffffff 100%)`,
            border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          }}
        >
          <CardContent>
            <CompareArrowsIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              Aucune comparaison créée
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              Créez votre première comparaison pour analyser les performances entre membres
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => setCreateDialogOpen(true)}
              sx={{
                background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
              }}
            >
              Créer une Comparaison
            </Button>
          </CardContent>
        </Card>
      ) : (
        <Grid container spacing={3}>
          {comparisons.map((comparison) => (
            <Grid item xs={12} sm={6} md={4} key={comparison.id}>
              <Card
                sx={{
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  background: isDarkMode
                    ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
                    : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
                  border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                  transition: 'all 0.3s ease',
                  cursor: 'pointer',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: `0 8px 25px ${isDarkMode ? 'rgba(0,0,0,0.3)' : 'rgba(25, 118, 210, 0.15)'}`,
                    borderColor: leagueColors.blue[300],
                  },
                }}
                onClick={() => handleComparisonDetails(comparison)}
              >
                <CardContent sx={{ flexGrow: 1 }}>
                  {/* Header */}
                  <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', mb: 2 }}>
                    <Box sx={{ flexGrow: 1 }}>
                      <Typography variant="h6" sx={{ fontWeight: 600, mb: 0.5 }}>
                        {comparison.name}
                      </Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                        {getComparisonIcon(comparison.compare_type)}
                        <Chip
                          label={getComparisonTypeLabel(comparison.compare_type)}
                          size="small"
                          sx={{
                            height: 20,
                            fontSize: '0.7rem',
                            background: `${getComparisonTypeColor(comparison.compare_type)}20`,
                            color: getComparisonTypeColor(comparison.compare_type),
                            border: `1px solid ${getComparisonTypeColor(comparison.compare_type)}40`,
                          }}
                        />
                      </Box>
                    </Box>
                    <IconButton
                      size="small"
                      onClick={(e) => handleMenuOpen(e, comparison.id)}
                    >
                      <MoreVertIcon />
                    </IconButton>
                  </Box>

                  {/* Description */}
                  {comparison.description && (
                    <Typography 
                      variant="body2" 
                      color="text.secondary" 
                      sx={{ mb: 2, lineHeight: 1.4 }}
                    >
                      {comparison.description}
                    </Typography>
                  )}

                  {/* Creator */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                    <PersonIcon sx={{ fontSize: 16, color: 'text.secondary' }} />
                    <Typography variant="body2" color="text.secondary">
                      par {comparison.creator.riot_id}#{comparison.creator.riot_tag}
                    </Typography>
                  </Box>

                  {/* Results Status */}
                  {comparison.results ? (
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Chip
                        icon={<AnalyticsIcon />}
                        label="Résultats disponibles"
                        size="small"
                        sx={{
                          background: `linear-gradient(135deg, ${leagueColors.win}20 0%, ${leagueColors.win}10 100%)`,
                          color: leagueColors.win,
                          border: `1px solid ${leagueColors.win}40`,
                        }}
                      />
                    </Box>
                  ) : (
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Chip
                        icon={<CircularProgress size={12} />}
                        label="En cours d'analyse"
                        size="small"
                        sx={{
                          background: `${leagueColors.dark[400]}20`,
                          color: leagueColors.dark[400],
                          border: `1px solid ${leagueColors.dark[400]}40`,
                        }}
                      />
                    </Box>
                  )}
                </CardContent>

                {/* Footer */}
                <Box 
                  sx={{ 
                    px: 2, 
                    py: 1.5, 
                    borderTop: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                    background: isDarkMode 
                      ? `${leagueColors.dark[200]}20` 
                      : `${leagueColors.blue[50]}50`,
                  }}
                >
                  <Typography variant="caption" color="text.secondary">
                    Créé le {new Date(comparison.created_at).toLocaleDateString()}
                  </Typography>
                </Box>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {/* Comparison Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        PaperProps={{
          sx: {
            borderRadius: 2,
            border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          },
        }}
      >
        <MenuItem onClick={() => {
          const comparison = comparisons.find(c => c.id === menuComparisonId);
          if (comparison) handleComparisonDetails(comparison);
        }}>
          <AnalyticsIcon sx={{ mr: 1, fontSize: 18 }} />
          Voir les résultats
        </MenuItem>
        <MenuItem onClick={() => {
          if (menuComparisonId) handleRegenerateComparison(menuComparisonId);
        }}>
          <RefreshIcon sx={{ mr: 1, fontSize: 18 }} />
          Régénérer
        </MenuItem>
        <MenuItem onClick={handleMenuClose}>
          <ShareIcon sx={{ mr: 1, fontSize: 18 }} />
          Partager
        </MenuItem>
        <MenuItem onClick={handleMenuClose} sx={{ color: 'error.main' }}>
          <DeleteIcon sx={{ mr: 1, fontSize: 18 }} />
          Supprimer
        </MenuItem>
      </Menu>

      {/* Dialogs */}
      <CreateComparisonDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onCreateComparison={handleCreateComparison}
        groupId={groupId}
      />

      {selectedComparison && (
        <ComparisonDetailsDialog
          open={detailsDialogOpen}
          onClose={() => {
            setDetailsDialogOpen(false);
            setSelectedComparison(null);
          }}
          comparison={selectedComparison}
          groupId={groupId}
          onUpdate={loadComparisons}
        />
      )}
    </Box>
  );
};

export default ComparisonManager;