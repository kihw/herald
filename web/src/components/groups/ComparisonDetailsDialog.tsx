import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  IconButton,
  Box,
  Typography,
  useTheme,
  Tabs,
  Tab,
  Card,
  CardContent,
  Grid,
  Chip,
  Avatar,
  List,
  ListItem,
  ListItemText,
  Divider,
  Alert,
  CircularProgress,
  Button,
  Paper,
} from '@mui/material';
import {
  Close as CloseIcon,
  TrendingUp as TrendingUpIcon,
  Analytics as AnalyticsIcon,
  CompareArrows as CompareArrowsIcon,
  EmojiEvents as EmojiEventsIcon,
  Person as PersonIcon,
  Share as ShareIcon,
  Refresh as RefreshIcon,
  GetApp as GetAppIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import { GroupComparison, groupApi } from '../../services/groupApi';
import ComparisonCharts from '../charts/ComparisonCharts';

interface ComparisonDetailsDialogProps {
  open: boolean;
  onClose: () => void;
  comparison: GroupComparison;
  groupId: number;
  onUpdate: () => void;
}

const ComparisonDetailsDialog: React.FC<ComparisonDetailsDialogProps> = ({ 
  open, 
  onClose, 
  comparison, 
  groupId, 
  onUpdate 
}) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [activeTab, setActiveTab] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [comparisonData, setComparisonData] = useState<GroupComparison>(comparison);

  useEffect(() => {
    setComparisonData(comparison);
  }, [comparison]);

  const handleRegenerateResults = async () => {
    try {
      setLoading(true);
      setError(null);
      const updatedComparison = await groupApi.regenerateComparison(groupId, comparison.id);
      setComparisonData(updatedComparison);
      onUpdate();
    } catch (err) {
      setError('Erreur lors de la régénération des résultats');
      console.error('Error regenerating comparison:', err);
    } finally {
      setLoading(false);
    }
  };

  const getComparisonIcon = (compareType: string) => {
    switch (compareType) {
      case 'champions':
        return <EmojiEventsIcon sx={{ fontSize: 24, color: leagueColors.gold[500] }} />;
      case 'roles':
        return <CompareArrowsIcon sx={{ fontSize: 24, color: leagueColors.blue[500] }} />;
      case 'performance':
        return <TrendingUpIcon sx={{ fontSize: 24, color: leagueColors.win }} />;
      case 'trends':
        return <AnalyticsIcon sx={{ fontSize: 24, color: leagueColors.dark[400] }} />;
      default:
        return <CompareArrowsIcon sx={{ fontSize: 24, color: 'text.secondary' }} />;
    }
  };

  const getComparisonTypeLabel = (compareType: string) => {
    switch (compareType) {
      case 'champions':
        return 'Comparaison Champions';
      case 'roles':
        return 'Comparaison Rôles';
      case 'performance':
        return 'Comparaison Performance';
      case 'trends':
        return 'Analyse des Tendances';
      default:
        return 'Comparaison';
    }
  };

  const renderSummaryTab = () => {
    if (!comparisonData.results) {
      return (
        <Box sx={{ textAlign: 'center', py: 6 }}>
          <CircularProgress size={48} sx={{ mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            Analyse en cours
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Les résultats seront disponibles sous peu
          </Typography>
          <Button
            variant="outlined"
            startIcon={<RefreshIcon />}
            onClick={handleRegenerateResults}
            disabled={loading}
          >
            Vérifier les résultats
          </Button>
        </Box>
      );
    }

    const { results } = comparisonData;

    return (
      <Grid container spacing={3}>
        {/* Summary Cards */}
        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.win}20 0%, ${leagueColors.win}10 100%)`,
              border: `1px solid ${leagueColors.win}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <EmojiEventsIcon sx={{ fontSize: 40, color: leagueColors.win, mb: 1 }} />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                Meilleur Joueur
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {results.summary.top_performer}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.blue[500]}20 0%, ${leagueColors.blue[500]}10 100%)`,
              border: `1px solid ${leagueColors.blue[500]}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <TrendingUpIcon sx={{ fontSize: 40, color: leagueColors.blue[500], mb: 1 }} />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                Meilleure Métrique
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {results.summary.best_metric}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.gold[500]}20 0%, ${leagueColors.gold[500]}10 100%)`,
              border: `1px solid ${leagueColors.gold[500]}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <AnalyticsIcon sx={{ fontSize: 40, color: leagueColors.gold[500], mb: 1 }} />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                Winrate Moyen
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {(results.summary.average_win_rate * 100).toFixed(1)}%
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card 
            sx={{ 
              background: `linear-gradient(135deg, ${leagueColors.dark[400]}20 0%, ${leagueColors.dark[400]}10 100%)`,
              border: `1px solid ${leagueColors.dark[400]}30`,
            }}
          >
            <CardContent sx={{ textAlign: 'center' }}>
              <CompareArrowsIcon sx={{ fontSize: 40, color: leagueColors.dark[400], mb: 1 }} />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                Parties Analysées
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {results.summary.total_games_compared}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Insights */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Analyses et Observations
              </Typography>
              <List>
                {results.insights.map((insight, index) => (
                  <ListItem key={index} sx={{ px: 0 }}>
                    <ListItemText
                      primary={
                        <Typography variant="body1" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Box
                            sx={{
                              width: 6,
                              height: 6,
                              borderRadius: '50%',
                              backgroundColor: leagueColors.blue[500],
                              flexShrink: 0,
                            }}
                          />
                          {insight}
                        </Typography>
                      }
                    />
                  </ListItem>
                ))}
              </List>
            </CardContent>
          </Card>
        </Grid>

        {/* Rankings */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Classement des Joueurs
              </Typography>
              <List>
                {results.rankings.map((ranking, index) => (
                  <React.Fragment key={ranking.user_id}>
                    <ListItem sx={{ px: 0 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
                        <Chip
                          label={`#${ranking.rank}`}
                          size="small"
                          sx={{
                            minWidth: 40,
                            fontWeight: 600,
                            background: ranking.rank === 1 
                              ? `linear-gradient(135deg, ${leagueColors.gold[500]} 0%, ${leagueColors.gold[600]} 100%)`
                              : ranking.rank === 2
                              ? `linear-gradient(135deg, #C0C0C0 0%, #A8A8A8 100%)`
                              : ranking.rank === 3
                              ? `linear-gradient(135deg, #CD7F32 0%, #B8860B 100%)`
                              : 'rgba(0,0,0,0.1)',
                            color: ranking.rank <= 3 ? '#000' : 'text.primary',
                          }}
                        />
                        <Avatar sx={{ width: 32, height: 32 }}>
                          {ranking.username.charAt(0).toUpperCase()}
                        </Avatar>
                        <Box sx={{ flexGrow: 1 }}>
                          <Typography variant="body1" sx={{ fontWeight: 500 }}>
                            {ranking.username}
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            {ranking.metric}: {ranking.score.toFixed(2)}
                          </Typography>
                        </Box>
                        <Box sx={{ textAlign: 'right' }}>
                          {ranking.change === 'up' && (
                            <TrendingUpIcon sx={{ color: leagueColors.win, fontSize: 20 }} />
                          )}
                          {ranking.change === 'down' && (
                            <TrendingUpIcon sx={{ 
                              color: leagueColors.loss, 
                              fontSize: 20, 
                              transform: 'rotate(180deg)' 
                            }} />
                          )}
                        </Box>
                      </Box>
                    </ListItem>
                    {index < results.rankings.length - 1 && <Divider variant="inset" component="li" />}
                  </React.Fragment>
                ))}
              </List>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  };

  const renderChartsTab = () => {
    if (!comparisonData.results?.charts || comparisonData.results.charts.length === 0) {
      // Generate mock charts for demonstration when no real data is available
      const mockCharts = [
        {
          type: 'bar' as const,
          title: 'Comparaison des Taux de Victoire',
          labels: ['Player1', 'Player2', 'Player3'],
          datasets: [
            {
              label: 'Taux de Victoire (%)',
              data: [65, 72, 58],
              background_color: undefined,
              border_color: undefined,
            },
          ],
        },
        {
          type: 'radar' as const,
          title: 'Performance Multi-Critères',
          labels: ['KDA', 'CS/min', 'Vision', 'Dégâts', 'Or/min'],
          datasets: [
            {
              label: 'Player1',
              data: [7, 8, 6, 9, 7],
              background_color: undefined,
              border_color: undefined,
            },
            {
              label: 'Player2',
              data: [8, 7, 8, 7, 8],
              background_color: undefined,
              border_color: undefined,
            },
          ],
        },
        {
          type: 'line' as const,
          title: 'Évolution des Performances',
          labels: ['Semaine 1', 'Semaine 2', 'Semaine 3', 'Semaine 4'],
          datasets: [
            {
              label: 'Player1',
              data: [60, 65, 63, 68],
              background_color: undefined,
              border_color: undefined,
            },
            {
              label: 'Player2',
              data: [55, 62, 70, 72],
              background_color: undefined,
              border_color: undefined,
            },
          ],
        },
        {
          type: 'pie' as const,
          title: 'Répartition des Rôles Joués',
          labels: ['ADC', 'Support', 'Mid', 'Top', 'Jungle'],
          datasets: [
            {
              label: 'Parties',
              data: [35, 25, 20, 15, 5],
              background_color: undefined,
              border_color: undefined,
            },
          ],
        },
      ];

      return <ComparisonCharts charts={mockCharts} />;
    }

    return <ComparisonCharts charts={comparisonData.results.charts} />;
  };

  const renderDetailsTab = () => {
    return (
      <Grid container spacing={3}>
        {/* Comparison Info */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Détails de la Comparaison
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" color="text.secondary">
                    Type de comparaison
                  </Typography>
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {getComparisonTypeLabel(comparisonData.compare_type)}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" color="text.secondary">
                    Créé par
                  </Typography>
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {comparisonData.creator.riot_id}#{comparisonData.creator.riot_tag}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" color="text.secondary">
                    Date de création
                  </Typography>
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {new Date(comparisonData.created_at).toLocaleDateString()}
                  </Typography>
                </Grid>
                {comparisonData.results && (
                  <Grid item xs={12} sm={6}>
                    <Typography variant="body2" color="text.secondary">
                      Dernière analyse
                    </Typography>
                    <Typography variant="body1" sx={{ fontWeight: 500 }}>
                      {new Date(comparisonData.results.generated_at).toLocaleDateString()}
                    </Typography>
                  </Grid>
                )}
              </Grid>

              {comparisonData.description && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="body2" color="text.secondary">
                    Description
                  </Typography>
                  <Typography variant="body1">
                    {comparisonData.description}
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Raw Data */}
        {comparisonData.results && (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                  Données Brutes
                </Typography>
                <Paper
                  sx={{
                    p: 2,
                    backgroundColor: isDarkMode ? leagueColors.dark[50] : leagueColors.blue[25],
                    border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                    maxHeight: 400,
                    overflow: 'auto',
                  }}
                >
                  <Typography
                    component="pre"
                    variant="body2"
                    sx={{
                      fontFamily: 'monospace',
                      fontSize: '0.75rem',
                      whiteSpace: 'pre-wrap',
                      wordBreak: 'break-word',
                    }}
                  >
                    {JSON.stringify(comparisonData.results.member_stats, null, 2)}
                  </Typography>
                </Paper>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    );
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="lg"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: 3,
          border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          background: isDarkMode
            ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
            : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
          minHeight: 700,
        },
      }}
    >
      <DialogTitle
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          pb: 0,
          background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
          color: '#fff',
          m: -3,
          mb: 0,
          px: 3,
          py: 2,
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          {getComparisonIcon(comparisonData.compare_type)}
          <Box>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>
              {comparisonData.name}
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.9 }}>
              {getComparisonTypeLabel(comparisonData.compare_type)}
            </Typography>
          </Box>
        </Box>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <IconButton 
            sx={{ color: '#fff' }}
            onClick={handleRegenerateResults}
            disabled={loading}
          >
            <RefreshIcon />
          </IconButton>
          <IconButton sx={{ color: '#fff' }}>
            <ShareIcon />
          </IconButton>
          <IconButton sx={{ color: '#fff' }}>
            <GetAppIcon />
          </IconButton>
          <IconButton onClick={onClose} sx={{ color: '#fff' }}>
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs
          value={activeTab}
          onChange={(e, newValue) => setActiveTab(newValue)}
          sx={{
            px: 3,
            '& .MuiTab-root': {
              fontWeight: 500,
            },
          }}
        >
          <Tab 
            icon={<TrendingUpIcon />} 
            label="Résumé" 
            iconPosition="start"
          />
          <Tab 
            icon={<AnalyticsIcon />} 
            label="Graphiques" 
            iconPosition="start"
          />
          <Tab 
            icon={<CompareArrowsIcon />} 
            label="Détails" 
            iconPosition="start"
          />
        </Tabs>
      </Box>

      <DialogContent sx={{ pt: 3 }}>
        {/* Error Alert */}
        {error && (
          <Alert 
            severity="error" 
            sx={{ mb: 3 }}
            onClose={() => setError(null)}
          >
            {error}
          </Alert>
        )}

        {/* Tab Content */}
        {activeTab === 0 && renderSummaryTab()}
        {activeTab === 1 && renderChartsTab()}
        {activeTab === 2 && renderDetailsTab()}
      </DialogContent>
    </Dialog>
  );
};

export default ComparisonDetailsDialog;