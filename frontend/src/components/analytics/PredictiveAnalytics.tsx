import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Tabs,
  Tab,
  Grid,
  Chip,
  LinearProgress,
  Avatar,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Alert,
  Tooltip,
  IconButton,
  Button,
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  useTheme,
} from '@mui/material';
import {
  ExpandMore,
  TrendingUp,
  TrendingDown,
  Timeline,
  Psychology,
  EmojiEvents,
  Group,
  Lightbulb,
  Star,
  Warning,
  Info,
  CheckCircle,
  RadioButtonUnchecked,
  Schedule,
  Speed,
  Assessment,
  School,
  Refresh,
} from '@mui/icons-material';
import { Line, Radar, Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip as ChartTooltip,
  Legend,
  RadialLinearScale,
  BarElement,
} from 'chart.js';
import { predictiveService } from '../../services/predictiveService';
import type {
  PredictiveAnalysis,
  PerformancePredictionData,
  RankProgressionPrediction,
  SkillDevelopmentForecast,
  ChampionRecommendationData,
  CareerTrajectoryForecast,
  PlayerPotentialAssessment,
  ActionableInsight,
} from '../../types/predictive';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  ChartTooltip,
  Legend,
  RadialLinearScale,
  BarElement
);

interface PredictiveAnalyticsProps {
  summonerId: string;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`predictive-tabpanel-${index}`}
      aria-labelledby={`predictive-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

function a11yProps(index: number) {
  return {
    id: `predictive-tab-${index}`,
    'aria-controls': `predictive-tabpanel-${index}`,
  };
}

const PredictiveAnalytics: React.FC<PredictiveAnalyticsProps> = ({ summonerId }) => {
  const theme = useTheme();
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [analysis, setAnalysis] = useState<PredictiveAnalysis | null>(null);

  useEffect(() => {
    loadPredictiveAnalysis();
  }, [summonerId]);

  const loadPredictiveAnalysis = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await predictiveService.getComprehensivePredictiveAnalysis(summonerId);
      setAnalysis(data);
    } catch (err) {
      setError('Failed to load predictive analysis');
      console.error('Predictive analysis error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 0.8) return theme.palette.success.main;
    if (confidence >= 0.6) return theme.palette.warning.main;
    return theme.palette.error.main;
  };

  const getConfidenceLabel = (confidence: number) => {
    if (confidence >= 0.8) return 'High';
    if (confidence >= 0.6) return 'Medium';
    return 'Low';
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'rising':
      case 'improving':
      case 'ascending':
        return <TrendingUp color="success" />;
      case 'declining':
      case 'falling':
      case 'descending':
        return <TrendingDown color="error" />;
      default:
        return <Timeline color="action" />;
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'critical':
      case 'high':
        return theme.palette.error.main;
      case 'medium':
        return theme.palette.warning.main;
      case 'low':
        return theme.palette.success.main;
      default:
        return theme.palette.grey[500];
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="h6" sx={{ ml: 2 }}>
          Generating Predictive Analysis...
        </Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Alert 
        severity="error" 
        action={
          <Button color="inherit" size="small" onClick={loadPredictiveAnalysis}>
            Retry
          </Button>
        }
      >
        {error}
      </Alert>
    );
  }

  if (!analysis) {
    return (
      <Alert severity="info">
        No predictive analysis data available for this summoner.
      </Alert>
    );
  }

  const renderPerformancePrediction = () => {
    const performance = analysis.performancePrediction;
    if (!performance) return <Typography>Performance prediction not available</Typography>;

    const progressData = {
      labels: ['Short Term', 'Medium Term', 'Long Term'],
      datasets: [
        {
          label: 'Win Rate Forecast',
          data: [
            performance.shortTermForecasts.winRateEstimate,
            performance.mediumTermForecasts.skillRatingChange + 50,
            performance.nextGameWinProbability,
          ],
          borderColor: theme.palette.primary.main,
          backgroundColor: `${theme.palette.primary.main}20`,
          tension: 0.4,
          fill: true,
        },
      ],
    };

    return (
      <Grid container spacing={3}>
        {/* Next Game Prediction */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Next Game Prediction
              </Typography>
              <Box display="flex" alignItems="center" mb={2}>
                <CircularProgress
                  variant="determinate"
                  value={performance.nextGameWinProbability}
                  size={80}
                  thickness={6}
                />
                <Box ml={2}>
                  <Typography variant="h4" color="primary">
                    {Math.round(performance.nextGameWinProbability)}%
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Win Probability
                  </Typography>
                </Box>
              </Box>
              <Typography variant="subtitle2" gutterBottom>
                Expected KDA: {performance.expectedKDA.kdaRatio.toFixed(1)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                K: {performance.expectedKDA.kills.toFixed(1)} / 
                D: {performance.expectedKDA.deaths.toFixed(1)} / 
                A: {performance.expectedKDA.assists.toFixed(1)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Performance Factors */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Performance Factors
              </Typography>
              <List dense>
                <ListItem>
                  <ListItemText
                    primary="Champion Pool Diversity"
                    secondary={`${performance.performanceFactors.championPool.diversity.toFixed(1)}%`}
                  />
                  <LinearProgress
                    variant="determinate"
                    value={performance.performanceFactors.championPool.diversity}
                    sx={{ width: 100, ml: 1 }}
                  />
                </ListItem>
                <ListItem>
                  <ListItemText
                    primary="Playstyle Consistency"
                    secondary={`${performance.performanceFactors.playstyleConsistency.toFixed(1)}%`}
                  />
                  <LinearProgress
                    variant="determinate"
                    value={performance.performanceFactors.playstyleConsistency}
                    sx={{ width: 100, ml: 1 }}
                  />
                </ListItem>
                <ListItem>
                  <ListItemText
                    primary="Adaptability"
                    secondary={`${performance.performanceFactors.adaptabilityScore.toFixed(1)}%`}
                  />
                  <LinearProgress
                    variant="determinate"
                    value={performance.performanceFactors.adaptabilityScore}
                    sx={{ width: 100, ml: 1 }}
                  />
                </ListItem>
                <ListItem>
                  <ListItemText
                    primary="Learning Rate"
                    secondary={`${performance.performanceFactors.learningRate.toFixed(1)}%`}
                  />
                  <LinearProgress
                    variant="determinate"
                    value={performance.performanceFactors.learningRate}
                    sx={{ width: 100, ml: 1 }}
                  />
                </ListItem>
              </List>
            </CardContent>
          </Card>
        </Grid>

        {/* Forecast Chart */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Performance Forecast
              </Typography>
              <Box height={300}>
                <Line
                  data={progressData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                      legend: {
                        position: 'top' as const,
                      },
                      title: {
                        display: true,
                        text: 'Performance Prediction Timeline',
                      },
                    },
                    scales: {
                      y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                          callback: function(value) {
                            return value + '%';
                          },
                        },
                      },
                    },
                  }}
                />
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  };

  const renderRankProgression = () => {
    const rankProgression = analysis.rankProgression;
    if (!rankProgression) return <Typography>Rank progression not available</Typography>;

    return (
      <Grid container spacing={3}>
        {/* Current Status */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Current Rank
              </Typography>
              <Box display="flex" alignItems="center" mb={2}>
                <Avatar sx={{ bgcolor: theme.palette.primary.main, mr: 2 }}>
                  <EmojiEvents />
                </Avatar>
                <Box>
                  <Typography variant="h4">
                    {rankProgression.currentRank}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {rankProgression.currentLP} LP
                  </Typography>
                </Box>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Trend: {getTrendIcon(rankProgression.progressionFactors.currentTrend)}
                <Chip
                  label={rankProgression.progressionFactors.currentTrend}
                  size="small"
                  sx={{ ml: 1 }}
                />
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Progression Scenarios */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Progression Scenarios
              </Typography>
              <TableContainer>
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Scenario</TableCell>
                      <TableCell align="right">Time to Target</TableCell>
                      <TableCell align="right">Games Required</TableCell>
                      <TableCell align="right">Win Rate Needed</TableCell>
                      <TableCell align="right">Probability</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {rankProgression.progressionScenarios.map((scenario, index) => (
                      <TableRow key={index}>
                        <TableCell component="th" scope="row">
                          <Chip
                            label={scenario.scenario}
                            size="small"
                            color={
                              scenario.scenario === 'optimistic' ? 'success' :
                              scenario.scenario === 'realistic' ? 'primary' : 'default'
                            }
                          />
                        </TableCell>
                        <TableCell align="right">{scenario.timeToTarget} days</TableCell>
                        <TableCell align="right">{scenario.gamesRequired}</TableCell>
                        <TableCell align="right">{(scenario.winRateRequired * 100).toFixed(1)}%</TableCell>
                        <TableCell align="right">
                          <Chip
                            label={`${(scenario.probability * 100).toFixed(1)}%`}
                            size="small"
                            style={{ backgroundColor: getConfidenceColor(scenario.probability) }}
                          />
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Barriers and Accelerators */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom color="error">
                <Warning sx={{ mr: 1, verticalAlign: 'middle' }} />
                Barriers to Progress
              </Typography>
              {rankProgression.barriers.map((barrier, index) => (
                <Accordion key={index}>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Box display="flex" alignItems="center" width="100%">
                      <Typography sx={{ flexGrow: 1 }}>{barrier.factor}</Typography>
                      <Chip
                        label={barrier.impact}
                        size="small"
                        color={barrier.impact === 'high' ? 'error' : barrier.impact === 'medium' ? 'warning' : 'default'}
                        sx={{ mr: 1 }}
                      />
                    </Box>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Typography variant="body2" gutterBottom>
                      {barrier.description}
                    </Typography>
                    <Typography variant="subtitle2" gutterBottom>
                      Mitigation Strategies:
                    </Typography>
                    <List dense>
                      {barrier.mitigation.map((strategy, i) => (
                        <ListItem key={i}>
                          <ListItemIcon>
                            <CheckCircle color="success" fontSize="small" />
                          </ListItemIcon>
                          <ListItemText primary={strategy} />
                        </ListItem>
                      ))}
                    </List>
                  </AccordionDetails>
                </Accordion>
              ))}
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom color="success.main">
                <Speed sx={{ mr: 1, verticalAlign: 'middle' }} />
                Accelerators
              </Typography>
              {rankProgression.accelerators.map((accelerator, index) => (
                <Accordion key={index}>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Box display="flex" alignItems="center" width="100%">
                      <Typography sx={{ flexGrow: 1 }}>{accelerator.factor}</Typography>
                      <Chip
                        label={accelerator.impact}
                        size="small"
                        color={accelerator.impact === 'high' ? 'success' : accelerator.impact === 'medium' ? 'warning' : 'default'}
                        sx={{ mr: 1 }}
                      />
                    </Box>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Typography variant="body2" gutterBottom>
                      {accelerator.description}
                    </Typography>
                    <Typography variant="subtitle2" gutterBottom>
                      Activation Steps:
                    </Typography>
                    <List dense>
                      {accelerator.activation.map((step, i) => (
                        <ListItem key={i}>
                          <ListItemIcon>
                            <Star color="warning" fontSize="small" />
                          </ListItemIcon>
                          <ListItemText primary={step} />
                        </ListItem>
                      ))}
                    </List>
                  </AccordionDetails>
                </Accordion>
              ))}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  };

  const renderSkillDevelopment = () => {
    const skillDevelopment = analysis.skillDevelopment;
    if (!skillDevelopment) return <Typography>Skill development not available</Typography>;

    const skillRadarData = {
      labels: skillDevelopment.skillForecasts.map(skill => skill.category),
      datasets: [
        {
          label: 'Current Level',
          data: skillDevelopment.skillForecasts.map(skill => skill.currentLevel),
          borderColor: theme.palette.primary.main,
          backgroundColor: `${theme.palette.primary.main}30`,
          pointBackgroundColor: theme.palette.primary.main,
          pointBorderColor: theme.palette.primary.main,
          pointHoverBackgroundColor: theme.palette.primary.dark,
          pointHoverBorderColor: theme.palette.primary.dark,
        },
        {
          label: 'Predicted Level',
          data: skillDevelopment.skillForecasts.map(skill => skill.predictedLevel),
          borderColor: theme.palette.success.main,
          backgroundColor: `${theme.palette.success.main}30`,
          pointBackgroundColor: theme.palette.success.main,
          pointBorderColor: theme.palette.success.main,
          pointHoverBackgroundColor: theme.palette.success.dark,
          pointHoverBorderColor: theme.palette.success.dark,
        },
      ],
    };

    return (
      <Grid container spacing={3}>
        {/* Skill Radar Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Skill Development Radar
              </Typography>
              <Box height={400}>
                <Radar
                  data={skillRadarData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                      legend: {
                        position: 'top' as const,
                      },
                    },
                    scales: {
                      r: {
                        beginAtZero: true,
                        max: 100,
                      },
                    },
                  }}
                />
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Overall Trajectory */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Overall Development Trajectory
              </Typography>
              <Box mb={3}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Current Skill Level
                </Typography>
                <LinearProgress
                  variant="determinate"
                  value={skillDevelopment.overallTrajectory.currentSkillLevel}
                  sx={{ height: 8, borderRadius: 4, mb: 1 }}
                />
                <Typography variant="caption">
                  {skillDevelopment.overallTrajectory.currentSkillLevel.toFixed(1)}/100
                </Typography>
              </Box>
              <Box mb={3}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Predicted Skill Level ({skillDevelopment.forecastPeriod} days)
                </Typography>
                <LinearProgress
                  variant="determinate"
                  value={skillDevelopment.overallTrajectory.predictedSkillLevel}
                  color="success"
                  sx={{ height: 8, borderRadius: 4, mb: 1 }}
                />
                <Typography variant="caption">
                  {skillDevelopment.overallTrajectory.predictedSkillLevel.toFixed(1)}/100
                </Typography>
              </Box>
              <Box>
                <Typography variant="subtitle2" gutterBottom>
                  Development Rate: 
                  <Chip
                    label={skillDevelopment.overallTrajectory.developmentRate}
                    size="small"
                    color={
                      skillDevelopment.overallTrajectory.developmentRate === 'accelerating' ? 'success' :
                      skillDevelopment.overallTrajectory.developmentRate === 'steady' ? 'primary' : 'warning'
                    }
                    sx={{ ml: 1 }}
                  />
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Skill Balance: {skillDevelopment.overallTrajectory.skillBalance.toFixed(1)}%
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Individual Skill Forecasts */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Individual Skill Forecasts
              </Typography>
              {skillDevelopment.skillForecasts.map((skill, index) => (
                <Accordion key={index}>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Box display="flex" alignItems="center" width="100%">
                      <Typography sx={{ flexGrow: 1 }}>{skill.category}</Typography>
                      <Chip
                        label={`${skill.currentLevel.toFixed(1)} → ${skill.predictedLevel.toFixed(1)}`}
                        size="small"
                        color="primary"
                        sx={{ mr: 2 }}
                      />
                      <Chip
                        label={skill.learningCurve}
                        size="small"
                        variant="outlined"
                      />
                    </Box>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Grid container spacing={2}>
                      <Grid item xs={12} md={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          Development Plan:
                        </Typography>
                        <List dense>
                          {skill.developmentPlan.focusAreas.map((area, i) => (
                            <ListItem key={i}>
                              <ListItemIcon>
                                <RadioButtonUnchecked fontSize="small" />
                              </ListItemIcon>
                              <ListItemText primary={area} />
                            </ListItem>
                          ))}
                        </List>
                      </Grid>
                      <Grid item xs={12} md={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          Practice Recommendations:
                        </Typography>
                        <List dense>
                          {skill.developmentPlan.practiceRecommendations.map((rec, i) => (
                            <ListItem key={i}>
                              <ListItemIcon>
                                <School fontSize="small" />
                              </ListItemIcon>
                              <ListItemText primary={rec} />
                            </ListItem>
                          ))}
                        </List>
                      </Grid>
                    </Grid>
                  </AccordionDetails>
                </Accordion>
              ))}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  };

  const renderChampionRecommendations = () => {
    const recommendations = analysis.championRecommendations;
    if (!recommendations || recommendations.length === 0) {
      return <Typography>Champion recommendations not available</Typography>;
    }

    return (
      <Grid container spacing={3}>
        {recommendations.slice(0, 6).map((rec, index) => (
          <Grid item xs={12} md={6} lg={4} key={index}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" mb={2}>
                  <Avatar sx={{ mr: 2, bgcolor: theme.palette.secondary.main }}>
                    {rec.champion.charAt(0)}
                  </Avatar>
                  <Box>
                    <Typography variant="h6">{rec.champion}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      {rec.role} • {rec.recommendationType.replace('_', ' ')}
                    </Typography>
                  </Box>
                </Box>
                
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Overall Score
                  </Typography>
                  <LinearProgress
                    variant="determinate"
                    value={rec.overallScore}
                    sx={{ height: 6, borderRadius: 3 }}
                  />
                  <Typography variant="caption" color="text.secondary">
                    {rec.overallScore.toFixed(1)}/100
                  </Typography>
                </Box>

                <Grid container spacing={1} mb={2}>
                  <Grid item xs={6}>
                    <Typography variant="caption" color="text.secondary">
                      Predicted Win Rate
                    </Typography>
                    <Typography variant="body2" fontWeight="bold">
                      {(rec.predictedWinRate * 100).toFixed(1)}%
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="caption" color="text.secondary">
                      Learning Difficulty
                    </Typography>
                    <Chip
                      label={rec.learningDifficulty}
                      size="small"
                      color={
                        rec.learningDifficulty === 'easy' ? 'success' :
                        rec.learningDifficulty === 'medium' ? 'warning' : 'error'
                      }
                    />
                  </Grid>
                </Grid>

                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Meta Rating: 
                    <Chip
                      label={rec.currentMetaRating}
                      size="small"
                      sx={{ ml: 1 }}
                    />
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Time to Mastery: {rec.timeToMastery} games
                  </Typography>
                </Box>

                <Accordion>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Typography variant="subtitle2">
                      Why Recommended?
                    </Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <List dense>
                      {rec.reasons.slice(0, 3).map((reason, i) => (
                        <ListItem key={i}>
                          <ListItemIcon>
                            <Lightbulb 
                              fontSize="small" 
                              color={reason.importance === 'high' ? 'warning' : 'action'} 
                            />
                          </ListItemIcon>
                          <ListItemText 
                            primary={reason.factor}
                            secondary={reason.description}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </AccordionDetails>
                </Accordion>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    );
  };

  const renderActionableInsights = () => {
    const insights = analysis.actionableInsights;
    if (!insights || insights.length === 0) {
      return <Typography>No actionable insights available</Typography>;
    }

    return (
      <Grid container spacing={3}>
        {insights.map((insight, index) => (
          <Grid item xs={12} key={index}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" mb={2}>
                  <Box 
                    sx={{
                      bgcolor: getPriorityColor(insight.priority),
                      borderRadius: 1,
                      p: 1,
                      mr: 2
                    }}
                  >
                    {insight.type === 'improvement' && <TrendingUp sx={{ color: 'white' }} />}
                    {insight.type === 'opportunity' && <Lightbulb sx={{ color: 'white' }} />}
                    {insight.type === 'warning' && <Warning sx={{ color: 'white' }} />}
                    {insight.type === 'recommendation' && <Info sx={{ color: 'white' }} />}
                  </Box>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="h6">{insight.title}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      {insight.category} • {insight.timeframe.replace('_', ' ')}
                    </Typography>
                  </Box>
                  <Box display="flex" gap={1}>
                    <Chip
                      label={insight.priority}
                      size="small"
                      style={{ backgroundColor: getPriorityColor(insight.priority), color: 'white' }}
                    />
                    <Chip
                      label={insight.impact}
                      size="small"
                      color="primary"
                      variant="outlined"
                    />
                    <Chip
                      label={insight.difficulty}
                      size="small"
                      color="secondary"
                      variant="outlined"
                    />
                  </Box>
                </Box>

                <Typography variant="body2" paragraph>
                  {insight.description}
                </Typography>

                <Accordion>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Typography variant="subtitle2">Action Steps</Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Grid container spacing={2}>
                      <Grid item xs={12} md={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          Steps to Take:
                        </Typography>
                        <List dense>
                          {insight.actionSteps.map((step, i) => (
                            <ListItem key={i}>
                              <ListItemIcon>
                                <CheckCircle color="success" fontSize="small" />
                              </ListItemIcon>
                              <ListItemText primary={step} />
                            </ListItem>
                          ))}
                        </List>
                      </Grid>
                      <Grid item xs={12} md={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          Success Metrics:
                        </Typography>
                        <List dense>
                          {insight.successMetrics.map((metric, i) => (
                            <ListItem key={i}>
                              <ListItemIcon>
                                <Assessment color="primary" fontSize="small" />
                              </ListItemIcon>
                              <ListItemText primary={metric} />
                            </ListItem>
                          ))}
                        </List>
                        <Box mt={2}>
                          <Typography variant="body2" color="text.secondary">
                            Expected Outcome: {insight.expectedOutcome}
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            Confidence: {(insight.confidence * 100).toFixed(1)}%
                          </Typography>
                        </Box>
                      </Grid>
                    </Grid>
                  </AccordionDetails>
                </Accordion>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    );
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">
          <Psychology sx={{ mr: 1, verticalAlign: 'middle' }} />
          Predictive Analytics
        </Typography>
        <Box display="flex" alignItems="center" gap={2}>
          <Chip
            label={`Confidence: ${getConfidenceLabel(analysis.modelConfidence?.overallConfidence || 0.85)}`}
            color={analysis.modelConfidence?.overallConfidence >= 0.8 ? 'success' : 
                   analysis.modelConfidence?.overallConfidence >= 0.6 ? 'warning' : 'error'}
          />
          <Tooltip title="Refresh Analysis">
            <IconButton onClick={loadPredictiveAnalysis}>
              <Refresh />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={tabValue} onChange={handleTabChange} aria-label="predictive analytics tabs">
          <Tab label="Performance Prediction" {...a11yProps(0)} />
          <Tab label="Rank Progression" {...a11yProps(1)} />
          <Tab label="Skill Development" {...a11yProps(2)} />
          <Tab label="Champion Recommendations" {...a11yProps(3)} />
          <Tab label="Actionable Insights" {...a11yProps(4)} />
        </Tabs>
      </Box>

      <TabPanel value={tabValue} index={0}>
        {renderPerformancePrediction()}
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        {renderRankProgression()}
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        {renderSkillDevelopment()}
      </TabPanel>

      <TabPanel value={tabValue} index={3}>
        {renderChampionRecommendations()}
      </TabPanel>

      <TabPanel value={tabValue} index={4}>
        {renderActionableInsights()}
      </TabPanel>
    </Box>
  );
};

export default PredictiveAnalytics;