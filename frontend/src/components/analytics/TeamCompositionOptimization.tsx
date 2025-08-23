// Team Composition Optimization Component for Herald.lol
import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Tabs,
  Tab,
  Button,
  Chip,
  LinearProgress,
  Alert,
  IconButton,
  Tooltip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Switch,
  FormControlLabel,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Divider,
  Paper
} from '@mui/material';
import {
  ExpandMore as ExpandMoreIcon,
  Analytics as AnalyticsIcon,
  Group as GroupIcon,
  Star as StarIcon,
  TrendingUp as TrendingUpIcon,
  Security as SecurityIcon,
  Speed as SpeedIcon,
  Build as BuildIcon,
  Timeline as TimelineIcon,
  Psychology as PsychologyIcon,
  AutoAwesome as AutoAwesomeIcon,
  Compare as CompareIcon,
  PlayArrow as PlayArrowIcon,
  Block as BlockIcon
} from '@mui/icons-material';
import {
  Radar,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Legend,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  Area,
  AreaChart
} from 'recharts';
import { useAuth } from '../../contexts/AuthContext';
import teamCompositionService from '../../services/teamCompositionService';
import {
  TeamCompositionOptimization,
  CompositionRecommendation,
  CompositionAnalysis,
  OptimizationStrategy,
  MetaComposition,
  PlayerComfortData,
  BanStrategy,
  SynergyAnalysis,
  ScalingProfile
} from '../../types/teamComposition';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div role="tabpanel" hidden={value !== index} {...other}>
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

const TeamCompositionOptimization: React.FC = () => {
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Optimization Data
  const [optimization, setOptimization] = useState<TeamCompositionOptimization | null>(null);
  const [recommendations, setRecommendations] = useState<CompositionRecommendation[]>([]);
  const [compositionAnalysis, setCompositionAnalysis] = useState<CompositionAnalysis | null>(null);
  const [metaCompositions, setMetaCompositions] = useState<MetaComposition[]>([]);
  const [playerComfort, setPlayerComfort] = useState<PlayerComfortData | null>(null);
  const [banStrategy, setBanStrategy] = useState<BanStrategy | null>(null);

  // Configuration
  const [selectedStrategy, setSelectedStrategy] = useState<OptimizationStrategy>('balanced');
  const [gameMode, setGameMode] = useState('ranked');
  const [playerRoles, setPlayerRoles] = useState({
    top: true,
    jungle: true,
    mid: true,
    adc: true,
    support: true
  });
  const [constraints, setConstraints] = useState({
    maxNewChampions: 2,
    requireTank: false,
    requireADC: true,
    preferLateGame: false,
    preferEarlyGame: false
  });
  const [bannedChampions, setBannedChampions] = useState<string[]>([]);
  const [requiredChampions, setRequiredChampions] = useState<string[]>([]);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  const loadOptimization = useCallback(async () => {
    if (!user?.summonerId) return;

    setLoading(true);
    setError(null);

    try {
      // Load team composition optimization
      const playerData = Object.entries(playerRoles)
        .filter(([role, enabled]) => enabled)
        .map(([role]) => ({
          summonerId: user.summonerId,
          role: role,
          championPool: [], // Will be populated by service
          comfortLevel: 8,
          recentGames: 50
        }));

      const optimizationRequest = {
        playerData,
        strategy: selectedStrategy,
        constraints,
        preferences: {
          prioritizeMeta: 7,
          prioritizeSynergy: 8,
          prioritizeComfort: 6,
          prioritizeBalance: 8,
          prioritizeFlexibility: 5,
          avoidHighBanRate: true,
          preferProvenCombos: true,
          allowExperimental: false
        },
        gameMode,
        bannedChampions,
        requiredChampions
      };

      const optimizationResult = await teamCompositionService.optimizeTeamComposition(optimizationRequest);
      setOptimization(optimizationResult);

      // Load additional data in parallel
      const [
        suggestionsResult,
        metaResult,
        comfortResult
      ] = await Promise.all([
        teamCompositionService.getCompositionSuggestions(user.summonerId, undefined, gameMode, selectedStrategy, 15),
        teamCompositionService.getMetaCompositions(gameMode, 'all', 'global', undefined, 20),
        teamCompositionService.getPlayerComfortPicks(user.summonerId, undefined, gameMode, 50, 25)
      ]);

      setRecommendations(suggestionsResult);
      setMetaCompositions(metaResult);
      setPlayerComfort(comfortResult);

    } catch (err) {
      console.error('Error loading team composition optimization:', err);
      setError(err instanceof Error ? err.message : 'Failed to load optimization data');
    } finally {
      setLoading(false);
    }
  }, [user?.summonerId, selectedStrategy, gameMode, playerRoles, constraints, bannedChampions, requiredChampions]);

  const analyzeComposition = async (composition: Array<{ champion: string; role: string }>) => {
    if (!user?.summonerId) return;

    setLoading(true);
    try {
      const analysis = await teamCompositionService.analyzeComposition(composition, undefined, gameMode);
      setCompositionAnalysis(analysis);
    } catch (err) {
      console.error('Error analyzing composition:', err);
      setError(err instanceof Error ? err.message : 'Failed to analyze composition');
    } finally {
      setLoading(false);
    }
  };

  const generateBanStrategy = async () => {
    if (!user?.summonerId) return;

    setLoading(true);
    try {
      const playerData = Object.entries(playerRoles)
        .filter(([role, enabled]) => enabled)
        .map(([role]) => ({
          summonerId: user.summonerId,
          role: role
        }));

      const strategy = await teamCompositionService.getBanStrategy(
        playerData,
        [], // Enemy data (unknown in this context)
        'first_ban',
        bannedChampions,
        gameMode,
        'balanced'
      );
      setBanStrategy(strategy);
    } catch (err) {
      console.error('Error generating ban strategy:', err);
      setError(err instanceof Error ? err.message : 'Failed to generate ban strategy');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadOptimization();
  }, [loadOptimization]);

  const createSynergyChartData = (synergy: SynergyAnalysis) => {
    return [
      { subject: 'Combo Potential', value: synergy.overall || 0, fullMark: 100 },
      { subject: 'Engage Synergy', value: (synergy.chainSynergies?.filter(s => s.type === 'engage').length || 0) * 20, fullMark: 100 },
      { subject: 'Protection', value: (synergy.chainSynergies?.filter(s => s.type === 'protection').length || 0) * 25, fullMark: 100 },
      { subject: 'Amplification', value: (synergy.chainSynergies?.filter(s => s.type === 'amplification').length || 0) * 20, fullMark: 100 },
      { subject: 'Follow-up', value: (synergy.chainSynergies?.filter(s => s.type === 'follow_up').length || 0) * 15, fullMark: 100 }
    ];
  };

  const createScalingChartData = (profile: ScalingProfile) => {
    return profile.scalingCurve.map(point => ({
      minute: point.minute,
      power: point.powerLevel,
      factors: point.keyFactors.join(', ')
    }));
  };

  if (loading && !optimization) {
    return (
      <Box sx={{ width: '100%', mt: 2 }}>
        <LinearProgress />
        <Typography sx={{ mt: 2, textAlign: 'center' }}>
          Optimizing team composition...
        </Typography>
      </Box>
    );
  }

  if (error && !optimization) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        {error}
      </Alert>
    );
  }

  return (
    <Box sx={{ width: '100%' }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <GroupIcon color="primary" />
        Team Composition Optimization
      </Typography>

      {/* Configuration Panel */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>Configuration</Typography>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Strategy</InputLabel>
                <Select
                  value={selectedStrategy}
                  onChange={(e) => setSelectedStrategy(e.target.value as OptimizationStrategy)}
                  label="Strategy"
                >
                  <MenuItem value="balanced">Balanced</MenuItem>
                  <MenuItem value="meta_optimal">Meta Optimal</MenuItem>
                  <MenuItem value="synergy_focused">Synergy Focused</MenuItem>
                  <MenuItem value="comfort_picks">Comfort Picks</MenuItem>
                  <MenuItem value="counter_focused">Counter Focused</MenuItem>
                  <MenuItem value="scaling_focused">Late Game</MenuItem>
                  <MenuItem value="early_focused">Early Game</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Game Mode</InputLabel>
                <Select
                  value={gameMode}
                  onChange={(e) => setGameMode(e.target.value)}
                  label="Game Mode"
                >
                  <MenuItem value="ranked">Ranked</MenuItem>
                  <MenuItem value="normal">Normal</MenuItem>
                  <MenuItem value="tournament">Tournament</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <TextField
                fullWidth
                size="small"
                type="number"
                label="Max New Champions"
                value={constraints.maxNewChampions}
                onChange={(e) => setConstraints(prev => ({
                  ...prev,
                  maxNewChampions: parseInt(e.target.value) || 0
                }))}
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Button
                variant="contained"
                startIcon={<AutoAwesomeIcon />}
                onClick={loadOptimization}
                disabled={loading}
                fullWidth
              >
                Optimize
              </Button>
            </Grid>
            <Grid item xs={12}>
              <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={constraints.requireTank}
                      onChange={(e) => setConstraints(prev => ({ ...prev, requireTank: e.target.checked }))}
                    />
                  }
                  label="Require Tank"
                />
                <FormControlLabel
                  control={
                    <Switch
                      checked={constraints.requireADC}
                      onChange={(e) => setConstraints(prev => ({ ...prev, requireADC: e.target.checked }))}
                    />
                  }
                  label="Require ADC"
                />
                <FormControlLabel
                  control={
                    <Switch
                      checked={constraints.preferLateGame}
                      onChange={(e) => setConstraints(prev => ({ ...prev, preferLateGame: e.target.checked }))}
                    />
                  }
                  label="Prefer Late Game"
                />
                <FormControlLabel
                  control={
                    <Switch
                      checked={constraints.preferEarlyGame}
                      onChange={(e) => setConstraints(prev => ({ ...prev, preferEarlyGame: e.target.checked }))}
                    />
                  }
                  label="Prefer Early Game"
                />
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onChange={handleTabChange} sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tab icon={<StarIcon />} label="Recommendations" />
        <Tab icon={<AnalyticsIcon />} label="Analysis" />
        <Tab icon={<TrendingUpIcon />} label="Meta Compositions" />
        <Tab icon={<PsychologyIcon />} label="Player Comfort" />
        <Tab icon={<BlockIcon />} label="Ban Strategy" />
      </Tabs>

      {/* Recommendations Tab */}
      <TabPanel value={activeTab} index={0}>
        <Grid container spacing={3}>
          {recommendations.map((recommendation, index) => (
            <Grid item xs={12} key={recommendation.id}>
              <Card sx={{ position: 'relative' }}>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'between', alignItems: 'flex-start', mb: 2 }}>
                    <Typography variant="h6">{recommendation.name}</Typography>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      <Chip
                        label={`${Math.round(recommendation.overallScore)}%`}
                        color="primary"
                        size="small"
                      />
                      <Chip
                        label={`Synergy: ${Math.round(recommendation.synergyScore)}%`}
                        color="secondary"
                        size="small"
                      />
                      <Chip
                        label={`Meta: ${Math.round(recommendation.metaScore)}%`}
                        color="success"
                        size="small"
                      />
                    </Box>
                  </Box>

                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    {recommendation.description}
                  </Typography>

                  {/* Champion Composition */}
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>Team Composition</Typography>
                    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      {recommendation.composition.map((pick, pickIndex) => (
                        <Chip
                          key={pickIndex}
                          label={`${pick.champion} (${pick.role})`}
                          variant={pick.priority === 'must_pick' ? 'filled' : 'outlined'}
                          color={
                            pick.priority === 'must_pick' ? 'error' :
                            pick.priority === 'preferred' ? 'primary' :
                            pick.priority === 'situational' ? 'secondary' : 'default'
                          }
                          size="small"
                        />
                      ))}
                    </Box>
                  </Box>

                  {/* Strengths and Weaknesses */}
                  <Grid container spacing={2}>
                    <Grid item xs={12} md={6}>
                      <Typography variant="subtitle2" color="success.main" gutterBottom>
                        Strengths
                      </Typography>
                      <List dense>
                        {recommendation.strengths.slice(0, 3).map((strength, strengthIndex) => (
                          <ListItem key={strengthIndex} sx={{ py: 0 }}>
                            <ListItemIcon sx={{ minWidth: 20 }}>
                              <Box
                                sx={{
                                  width: 6,
                                  height: 6,
                                  borderRadius: '50%',
                                  bgcolor: 'success.main'
                                }}
                              />
                            </ListItemIcon>
                            <ListItemText
                              primary={strength}
                              primaryTypographyProps={{ variant: 'body2' }}
                            />
                          </ListItem>
                        ))}
                      </List>
                    </Grid>
                    <Grid item xs={12} md={6}>
                      <Typography variant="subtitle2" color="warning.main" gutterBottom>
                        Weaknesses
                      </Typography>
                      <List dense>
                        {recommendation.weaknesses.slice(0, 3).map((weakness, weaknessIndex) => (
                          <ListItem key={weaknessIndex} sx={{ py: 0 }}>
                            <ListItemIcon sx={{ minWidth: 20 }}>
                              <Box
                                sx={{
                                  width: 6,
                                  height: 6,
                                  borderRadius: '50%',
                                  bgcolor: 'warning.main'
                                }}
                              />
                            </ListItemIcon>
                            <ListItemText
                              primary={weakness}
                              primaryTypographyProps={{ variant: 'body2' }}
                            />
                          </ListItem>
                        ))}
                      </List>
                    </Grid>
                  </Grid>

                  {/* Win Conditions */}
                  <Accordion sx={{ mt: 2 }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                      <Typography variant="subtitle2">Win Conditions & Strategy</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                      <Grid container spacing={2}>
                        <Grid item xs={12} md={6}>
                          <Typography variant="body2" fontWeight="medium" gutterBottom>
                            Win Conditions
                          </Typography>
                          {recommendation.winConditions.slice(0, 3).map((condition, conditionIndex) => (
                            <Box key={conditionIndex} sx={{ mb: 1 }}>
                              <Typography variant="body2">
                                {condition.condition} ({Math.round(condition.probability)}% success rate)
                              </Typography>
                              <Typography variant="caption" color="text.secondary">
                                {condition.timeline}
                              </Typography>
                            </Box>
                          ))}
                        </Grid>
                        <Grid item xs={12} md={6}>
                          <Typography variant="body2" fontWeight="medium" gutterBottom>
                            Key Strategy
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            {recommendation.playStyle.description}
                          </Typography>
                          <Box sx={{ mt: 1 }}>
                            {recommendation.playStyle.keyTactics.slice(0, 3).map((tactic, tacticIndex) => (
                              <Chip
                                key={tacticIndex}
                                label={tactic}
                                size="small"
                                variant="outlined"
                                sx={{ mr: 0.5, mb: 0.5 }}
                              />
                            ))}
                          </Box>
                        </Grid>
                      </Grid>
                    </AccordionDetails>
                  </Accordion>

                  {/* Action Button */}
                  <Box sx={{ display: 'flex', gap: 1, mt: 2 }}>
                    <Button
                      variant="contained"
                      size="small"
                      startIcon={<AnalyticsIcon />}
                      onClick={() => analyzeComposition(recommendation.composition)}
                    >
                      Analyze
                    </Button>
                    <Button
                      variant="outlined"
                      size="small"
                      startIcon={<CompareIcon />}
                      onClick={() => {
                        // Add to comparison - could be implemented later
                      }}
                    >
                      Compare
                    </Button>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </TabPanel>

      {/* Analysis Tab */}
      <TabPanel value={activeTab} index={1}>
        {compositionAnalysis ? (
          <Grid container spacing={3}>
            {/* Overall Rating */}
            <Grid item xs={12} md={4}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>Overall Rating</Typography>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Typography variant="h3" color="primary">
                      {Math.round(compositionAnalysis.overallRating)}
                    </Typography>
                    <Box>
                      <Typography variant="body1">
                        Tier {compositionAnalysis.tierRating}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Meta Fit: {Math.round(compositionAnalysis.metaFit)}%
                      </Typography>
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            </Grid>

            {/* Role Balance */}
            <Grid item xs={12} md={8}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>Role Balance</Typography>
                  <Grid container spacing={2}>
                    {[
                      { label: 'Tankiness', value: compositionAnalysis.roleBalance.tankiness },
                      { label: 'Damage', value: compositionAnalysis.roleBalance.damage },
                      { label: 'Utility', value: compositionAnalysis.roleBalance.utility },
                      { label: 'Engage', value: compositionAnalysis.roleBalance.engage },
                      { label: 'Peel', value: compositionAnalysis.roleBalance.peel },
                      { label: 'Wave Clear', value: compositionAnalysis.roleBalance.waveClear }
                    ].map((stat, index) => (
                      <Grid item xs={6} sm={4} key={index}>
                        <Box sx={{ mb: 1 }}>
                          <Typography variant="body2">{stat.label}</Typography>
                          <LinearProgress
                            variant="determinate"
                            value={stat.value}
                            sx={{ height: 8, borderRadius: 1 }}
                          />
                          <Typography variant="caption" color="text.secondary">
                            {stat.value}%
                          </Typography>
                        </Box>
                      </Grid>
                    ))}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>

            {/* Synergy Radar Chart */}
            {compositionAnalysis.synergy && (
              <Grid item xs={12} md={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>Team Synergy</Typography>
                    <ResponsiveContainer width="100%" height={300}>
                      <RadarChart data={createSynergyChartData(compositionAnalysis.synergy)}>
                        <PolarGrid />
                        <PolarAngleAxis dataKey="subject" />
                        <PolarRadiusAxis angle={90} domain={[0, 100]} />
                        <Radar
                          name="Synergy"
                          dataKey="value"
                          stroke="#8884d8"
                          fill="#8884d8"
                          fillOpacity={0.3}
                        />
                      </RadarChart>
                    </ResponsiveContainer>
                    <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                      Overall Synergy: {Math.round(compositionAnalysis.synergy.overall)}%
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            )}

            {/* Performance Prediction */}
            <Grid item xs={12} md={6}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>Performance Prediction</Typography>
                  <Grid container spacing={1}>
                    {[
                      { label: 'Early Game', value: compositionAnalysis.expectedPerformance.earlyGame, icon: <SpeedIcon /> },
                      { label: 'Mid Game', value: compositionAnalysis.expectedPerformance.midGame, icon: <BuildIcon /> },
                      { label: 'Late Game', value: compositionAnalysis.expectedPerformance.lateGame, icon: <TimelineIcon /> },
                      { label: 'Team Fighting', value: compositionAnalysis.expectedPerformance.teamFighting, icon: <GroupIcon /> },
                      { label: 'Objective Control', value: compositionAnalysis.expectedPerformance.objectiveControl, icon: <SecurityIcon /> }
                    ].map((phase, index) => (
                      <Grid item xs={12} key={index}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
                          {phase.icon}
                          <Typography variant="body2" sx={{ minWidth: 100 }}>
                            {phase.label}
                          </Typography>
                          <LinearProgress
                            variant="determinate"
                            value={phase.value}
                            sx={{ flexGrow: 1, height: 8, borderRadius: 1 }}
                          />
                          <Typography variant="caption" sx={{ minWidth: 40 }}>
                            {Math.round(phase.value)}%
                          </Typography>
                        </Box>
                      </Grid>
                    ))}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        ) : (
          <Alert severity="info">
            Select a team composition from the recommendations to see detailed analysis.
          </Alert>
        )}
      </TabPanel>

      {/* Meta Compositions Tab */}
      <TabPanel value={activeTab} index={2}>
        <Grid container spacing={3}>
          {metaCompositions.map((comp, index) => (
            <Grid item xs={12} md={6} key={comp.id}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'between', alignItems: 'flex-start', mb: 2 }}>
                    <Typography variant="h6">{comp.name}</Typography>
                    <Chip
                      label={`${Math.round(comp.winRate)}% WR`}
                      color={comp.winRate >= 55 ? 'success' : comp.winRate >= 50 ? 'primary' : 'warning'}
                      size="small"
                    />
                  </Box>

                  <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 2 }}>
                    {comp.composition.map((pick, pickIndex) => (
                      <Chip
                        key={pickIndex}
                        label={`${pick.champion} (${pick.role})`}
                        size="small"
                        variant="outlined"
                      />
                    ))}
                  </Box>

                  <Grid container spacing={2}>
                    <Grid item xs={6}>
                      <Typography variant="caption" color="text.secondary">Pick Rate</Typography>
                      <LinearProgress
                        variant="determinate"
                        value={comp.pickRate}
                        sx={{ height: 4, borderRadius: 1 }}
                      />
                      <Typography variant="caption">{Math.round(comp.pickRate)}%</Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="caption" color="text.secondary">Ban Rate</Typography>
                      <LinearProgress
                        variant="determinate"
                        value={comp.banRate}
                        color="warning"
                        sx={{ height: 4, borderRadius: 1 }}
                      />
                      <Typography variant="caption">{Math.round(comp.banRate)}%</Typography>
                    </Grid>
                  </Grid>

                  <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
                    <Chip
                      label={comp.trend === 'rising' ? '↗ Rising' : comp.trend === 'declining' ? '↘ Declining' : '→ Stable'}
                      color={comp.trend === 'rising' ? 'success' : comp.trend === 'declining' ? 'error' : 'default'}
                      size="small"
                    />
                    {comp.proPlayPresence > 20 && (
                      <Chip label="Pro Play" color="info" size="small" />
                    )}
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </TabPanel>

      {/* Player Comfort Tab */}
      <TabPanel value={activeTab} index={3}>
        {playerComfort ? (
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>Champion Comfort Analysis</Typography>
            </Grid>

            {/* Master Tier Champions */}
            <Grid item xs={12} md={6}>
              <Card>
                <CardContent>
                  <Typography variant="h6" color="success.main" gutterBottom>
                    Master Tier (S-Tier Comfort)
                  </Typography>
                  {playerComfort.masterTier.map((champion, index) => (
                    <Box key={index} sx={{ mb: 2, p: 2, bgcolor: 'success.light', borderRadius: 1 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'between', alignItems: 'center' }}>
                        <Typography variant="subtitle1">{champion.champion}</Typography>
                        <Chip
                          label={`${Math.round(champion.winRate)}% WR`}
                          size="small"
                          color="success"
                        />
                      </Box>
                      <Typography variant="body2" color="text.secondary">
                        {champion.masteryPoints.toLocaleString()} mastery • {champion.recentGames} recent games
                      </Typography>
                      <LinearProgress
                        variant="determinate"
                        value={champion.comfortScore}
                        sx={{ mt: 1, height: 6, borderRadius: 1 }}
                      />
                    </Box>
                  ))}
                </CardContent>
              </Card>
            </Grid>

            {/* Comfort Tier Champions */}
            <Grid item xs={12} md={6}>
              <Card>
                <CardContent>
                  <Typography variant="h6" color="primary.main" gutterBottom>
                    Comfort Tier (A-Tier Comfort)
                  </Typography>
                  {playerComfort.comfortTier.map((champion, index) => (
                    <Box key={index} sx={{ mb: 2, p: 2, bgcolor: 'primary.light', borderRadius: 1 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'between', alignItems: 'center' }}>
                        <Typography variant="subtitle1">{champion.champion}</Typography>
                        <Chip
                          label={`${Math.round(champion.winRate)}% WR`}
                          size="small"
                          color="primary"
                        />
                      </Box>
                      <Typography variant="body2" color="text.secondary">
                        {champion.masteryPoints.toLocaleString()} mastery • {champion.recentGames} recent games
                      </Typography>
                      <LinearProgress
                        variant="determinate"
                        value={champion.comfortScore}
                        sx={{ mt: 1, height: 6, borderRadius: 1 }}
                      />
                    </Box>
                  ))}
                </CardContent>
              </Card>
            </Grid>

            {/* Champion Pool Analysis */}
            <Grid item xs={12}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>Pool Analysis</Typography>
                  <Grid container spacing={2}>
                    <Grid item xs={12} sm={4}>
                      <Box sx={{ textAlign: 'center' }}>
                        <Typography variant="h4" color="primary">
                          {Math.round(playerComfort.poolDepth)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Pool Depth
                        </Typography>
                      </Box>
                    </Grid>
                    <Grid item xs={12} sm={4}>
                      <Box sx={{ textAlign: 'center' }}>
                        <Typography variant="h4" color="success.main">
                          {Math.round(playerComfort.metaAlignment)}%
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Meta Alignment
                        </Typography>
                      </Box>
                    </Grid>
                    <Grid item xs={12} sm={4}>
                      <Box sx={{ textAlign: 'center' }}>
                        <Typography variant="h4" color="info.main">
                          {Math.round(playerComfort.flexibility)}%
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Flexibility
                        </Typography>
                      </Box>
                    </Grid>
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        ) : (
          <Alert severity="info">Loading player comfort data...</Alert>
        )}
      </TabPanel>

      {/* Ban Strategy Tab */}
      <TabPanel value={activeTab} index={4}>
        <Box sx={{ mb: 2 }}>
          <Button
            variant="contained"
            startIcon={<BlockIcon />}
            onClick={generateBanStrategy}
            disabled={loading}
          >
            Generate Ban Strategy
          </Button>
        </Box>

        {banStrategy ? (
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    Recommended Bans ({banStrategy.strategy.replace('_', ' ').toUpperCase()})
                  </Typography>
                  <Grid container spacing={2}>
                    {banStrategy.recommendations.map((ban, index) => (
                      <Grid item xs={12} sm={6} md={4} key={index}>
                        <Paper sx={{ p: 2, bgcolor: 'error.light' }}>
                          <Box sx={{ display: 'flex', justifyContent: 'between', alignItems: 'center', mb: 1 }}>
                            <Typography variant="subtitle1">{ban.champion}</Typography>
                            <Chip
                              label={`Priority ${ban.priority}`}
                              size="small"
                              color={ban.priority >= 8 ? 'error' : ban.priority >= 6 ? 'warning' : 'default'}
                            />
                          </Box>
                          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                            Target: {ban.target.replace('_', ' ')}
                          </Typography>
                          <Typography variant="body2" sx={{ mb: 1 }}>
                            Effectiveness: {Math.round(ban.effectiveness)}%
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {ban.reasoning[0]}
                          </Typography>
                        </Paper>
                      </Grid>
                    ))}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>

            {/* Strategy Analysis */}
            <Grid item xs={12}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>Strategy Analysis</Typography>
                  <Grid container spacing={2}>
                    <Grid item xs={12} md={6}>
                      <Typography variant="subtitle2" gutterBottom>Expected Impact</Typography>
                      <List dense>
                        {banStrategy.analysis.expectedImpact.map((impact, index) => (
                          <ListItem key={index}>
                            <ListItemIcon>
                              <StarIcon color="primary" fontSize="small" />
                            </ListItemIcon>
                            <ListItemText primary={impact} />
                          </ListItem>
                        ))}
                      </List>
                    </Grid>
                    <Grid item xs={12} md={6}>
                      <Typography variant="subtitle2" gutterBottom>Risk Assessment</Typography>
                      <List dense>
                        {banStrategy.analysis.riskAssessment.map((risk, index) => (
                          <ListItem key={index}>
                            <ListItemIcon>
                              <SecurityIcon color="warning" fontSize="small" />
                            </ListItemIcon>
                            <ListItemText primary={risk} />
                          </ListItem>
                        ))}
                      </List>
                    </Grid>
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        ) : (
          <Alert severity="info">
            Click "Generate Ban Strategy" to get personalized ban recommendations based on your team composition and enemy threats.
          </Alert>
        )}
      </TabPanel>
    </Box>
  );
};

export default TeamCompositionOptimization;