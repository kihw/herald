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
  Button,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Alert,
  Stepper,
  Step,
  StepLabel,
  StepContent,
  Tooltip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Badge,
  Avatar,
  useTheme,
} from '@mui/material';
import {
  ExpandMore,
  TrendingUp,
  TrendingDown,
  CheckCircle,
  RadioButtonUnchecked,
  Star,
  Warning,
  Info,
  Lightbulb,
  PlayArrow,
  Pause,
  Schedule,
  Assessment,
  EmojiEvents,
  School,
  Psychology,
  Refresh,
  Timeline,
  Speed,
  Assignment,
  PersonalVideo,
  MenuBook,
  Build,
} from '@mui/icons-material';
import { Line, Radar, Bar } from 'react-chartjs-2';
import { improvementService } from '../../services/improvementService';
import type {
  ImprovementRecommendation,
  PlayerAnalysisResult,
  OverallProgress,
  ImprovementInsight,
  QuickWin,
  CoachingPlan,
  ProgressUpdateRequest,
} from '../../types/improvement';

interface ImprovementRecommendationsProps {
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
      id={`improvement-tabpanel-${index}`}
      aria-labelledby={`improvement-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

function a11yProps(index: number) {
  return {
    id: `improvement-tab-${index}`,
    'aria-controls': `improvement-tabpanel-${index}`,
  };
}

const ImprovementRecommendations: React.FC<ImprovementRecommendationsProps> = ({ summonerId }) => {
  const theme = useTheme();
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [recommendations, setRecommendations] = useState<ImprovementRecommendation[]>([]);
  const [analysis, setAnalysis] = useState<PlayerAnalysisResult | null>(null);
  const [progress, setProgress] = useState<OverallProgress | null>(null);
  const [insights, setInsights] = useState<ImprovementInsight[]>([]);
  const [quickWins, setQuickWins] = useState<QuickWin[]>([]);
  const [activeRecommendation, setActiveRecommendation] = useState<ImprovementRecommendation | null>(null);
  const [progressDialog, setProgressDialog] = useState(false);

  useEffect(() => {
    loadImprovementData();
  }, [summonerId]);

  const loadImprovementData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [
        recommendationsData,
        analysisData,
        progressData,
        insightsData,
        quickWinsData,
      ] = await Promise.all([
        improvementService.getPersonalizedRecommendations(summonerId),
        improvementService.getPlayerAnalysis(summonerId),
        improvementService.getOverallProgress(summonerId),
        improvementService.getImprovementInsights(summonerId),
        improvementService.getQuickWins(summonerId),
      ]);

      setRecommendations(recommendationsData.recommendations);
      setAnalysis(analysisData.analysis);
      setProgress(progressData);
      setInsights(insightsData.insights);
      setQuickWins(quickWinsData.quickWins);
    } catch (err) {
      setError('Failed to load improvement recommendations');
      console.error('Improvement recommendations error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'critical':
        return theme.palette.error.main;
      case 'high':
        return theme.palette.warning.main;
      case 'medium':
        return theme.palette.info.main;
      case 'low':
        return theme.palette.success.main;
      default:
        return theme.palette.grey[500];
    }
  };

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'very_easy':
        return theme.palette.success.main;
      case 'easy':
        return theme.palette.success.light;
      case 'medium':
        return theme.palette.warning.main;
      case 'hard':
        return theme.palette.error.light;
      case 'expert':
        return theme.palette.error.main;
      default:
        return theme.palette.grey[500];
    }
  };

  const getInsightIcon = (type: string) => {
    switch (type) {
      case 'opportunity':
        return <Lightbulb color="warning" />;
      case 'trend':
        return <Timeline color="info" />;
      case 'warning':
        return <Warning color="error" />;
      case 'achievement':
        return <EmojiEvents color="success" />;
      default:
        return <Info color="action" />;
    }
  };

  const handleStartRecommendation = async (recommendation: ImprovementRecommendation) => {
    setActiveRecommendation(recommendation);
    setProgressDialog(true);
  };

  const handleProgressUpdate = async (progressData: ProgressUpdateRequest) => {
    if (!activeRecommendation) return;
    
    try {
      await improvementService.updateRecommendationProgress(
        activeRecommendation.id,
        progressData
      );
      setProgressDialog(false);
      setActiveRecommendation(null);
      // Reload data to get updated progress
      await loadImprovementData();
    } catch (err) {
      console.error('Failed to update progress:', err);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="h6" sx={{ ml: 2 }}>
          Analyzing Your Performance...
        </Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Alert 
        severity="error" 
        action={
          <Button color="inherit" size="small" onClick={loadImprovementData}>
            Retry
          </Button>
        }
      >
        {error}
      </Alert>
    );
  }

  const renderRecommendations = () => {
    if (recommendations.length === 0) {
      return (
        <Alert severity="info">
          Great job! No critical improvement areas found. Keep up the excellent work!
        </Alert>
      );
    }

    return (
      <Grid container spacing={3}>
        {recommendations.map((rec, index) => (
          <Grid item xs={12} key={rec.id}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
                  <Box display="flex" alignItems="center">
                    <Avatar
                      sx={{ 
                        bgcolor: getPriorityColor(rec.priority),
                        mr: 2,
                        width: 56,
                        height: 56,
                      }}
                    >
                      <Assignment />
                    </Avatar>
                    <Box>
                      <Typography variant="h6">{rec.title}</Typography>
                      <Box display="flex" gap={1} mt={1}>
                        <Chip
                          label={rec.priority}
                          size="small"
                          style={{ 
                            backgroundColor: getPriorityColor(rec.priority),
                            color: 'white' 
                          }}
                        />
                        <Chip
                          label={rec.difficultyLevel}
                          size="small"
                          style={{ 
                            backgroundColor: getDifficultyColor(rec.difficultyLevel),
                            color: 'white' 
                          }}
                        />
                        <Chip
                          label={`${rec.timeToSeeResults} days`}
                          size="small"
                          icon={<Schedule />}
                          variant="outlined"
                        />
                      </Box>
                    </Box>
                  </Box>
                  <Box textAlign="center">
                    <Typography variant="h4" color="primary">
                      {Math.round(rec.impactScore)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Impact Score
                    </Typography>
                  </Box>
                </Box>

                <Typography variant="body2" color="text.secondary" paragraph>
                  {rec.description}
                </Typography>

                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Expected ROI: {rec.estimatedROI.toFixed(1)}%
                  </Typography>
                  <LinearProgress
                    variant="determinate"
                    value={rec.progressTracking.overallProgress}
                    sx={{ height: 8, borderRadius: 4 }}
                  />
                  <Typography variant="caption" color="text.secondary">
                    Progress: {rec.progressTracking.overallProgress.toFixed(1)}%
                  </Typography>
                </Box>

                <Accordion>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Typography variant="subtitle1">Action Plan</Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Box mb={2}>
                      <Typography variant="subtitle2" gutterBottom>
                        Primary Objective:
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {rec.actionPlan.primaryObjective}
                      </Typography>
                    </Box>

                    <Stepper orientation="vertical">
                      {rec.actionPlan.actionSteps.map((step, stepIndex) => (
                        <Step 
                          key={stepIndex} 
                          active={stepIndex === 0}
                          completed={rec.progressTracking.completedSteps.includes(step.stepNumber)}
                        >
                          <StepLabel>
                            {step.title}
                          </StepLabel>
                          <StepContent>
                            <Typography variant="body2" gutterBottom>
                              {step.description}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                              Duration: {step.duration} • Frequency: {step.frequency}
                            </Typography>
                            {step.tools.length > 0 && (
                              <Box mt={1}>
                                <Typography variant="caption" color="text.secondary">
                                  Tools: {step.tools.join(', ')}
                                </Typography>
                              </Box>
                            )}
                          </StepContent>
                        </Step>
                      ))}
                    </Stepper>
                  </AccordionDetails>
                </Accordion>

                <Box display="flex" justifyContent="space-between" alignItems="center" mt={2}>
                  <Box>
                    <Typography variant="caption" color="text.secondary">
                      Confidence: {rec.recommendationContext.confidenceScore.toFixed(1)}%
                    </Typography>
                  </Box>
                  <Button
                    variant="contained"
                    startIcon={<PlayArrow />}
                    onClick={() => handleStartRecommendation(rec)}
                    disabled={rec.status === 'completed'}
                  >
                    {rec.status === 'completed' ? 'Completed' : 'Start Working On This'}
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    );
  };

  const renderQuickWins = () => {
    return (
      <Grid container spacing={3}>
        {quickWins.map((quickWin, index) => (
          <Grid item xs={12} md={6} key={index}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" mb={2}>
                  <Speed color="warning" sx={{ mr: 2 }} />
                  <Typography variant="h6">{quickWin.title}</Typography>
                </Box>
                
                <Typography variant="body2" color="text.secondary" paragraph>
                  {quickWin.description}
                </Typography>

                <Grid container spacing={2} mb={2}>
                  <Grid item xs={6}>
                    <Typography variant="caption" color="text.secondary">
                      Expected Impact
                    </Typography>
                    <Typography variant="body2" fontWeight="bold">
                      {quickWin.expectedImpact}
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="caption" color="text.secondary">
                      Time to Implement
                    </Typography>
                    <Typography variant="body2" fontWeight="bold">
                      {quickWin.timeToImplement}
                    </Typography>
                  </Grid>
                </Grid>

                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    ROI Score
                  </Typography>
                  <LinearProgress
                    variant="determinate"
                    value={quickWin.roiScore}
                    color="warning"
                    sx={{ height: 6, borderRadius: 3 }}
                  />
                  <Typography variant="caption">
                    {quickWin.roiScore.toFixed(1)}/100
                  </Typography>
                </Box>

                <Accordion>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Typography variant="subtitle2">Instructions</Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <List dense>
                      {quickWin.instructions.map((instruction, i) => (
                        <ListItem key={i}>
                          <ListItemIcon>
                            <CheckCircle color="success" fontSize="small" />
                          </ListItemIcon>
                          <ListItemText primary={instruction} />
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

  const renderInsights = () => {
    return (
      <Grid container spacing={3}>
        {insights.map((insight, index) => (
          <Grid item xs={12} key={index}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" mb={2}>
                  {getInsightIcon(insight.type)}
                  <Box ml={2} flexGrow={1}>
                    <Typography variant="h6">{insight.title}</Typography>
                    <Box display="flex" gap={1} mt={1}>
                      <Chip
                        label={insight.type}
                        size="small"
                        variant="outlined"
                      />
                      <Chip
                        label={insight.priority}
                        size="small"
                        style={{ backgroundColor: getPriorityColor(insight.priority), color: 'white' }}
                      />
                      <Chip
                        label={insight.impact}
                        size="small"
                        color="secondary"
                        variant="outlined"
                      />
                    </Box>
                  </Box>
                </Box>

                <Typography variant="body2" color="text.secondary" paragraph>
                  {insight.description}
                </Typography>

                {insight.actionSteps && insight.actionSteps.length > 0 && (
                  <Accordion>
                    <AccordionSummary expandIcon={<ExpandMore />}>
                      <Typography variant="subtitle2">Action Steps</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                      <List dense>
                        {insight.actionSteps.map((step, i) => (
                          <ListItem key={i}>
                            <ListItemIcon>
                              <RadioButtonUnchecked fontSize="small" />
                            </ListItemIcon>
                            <ListItemText primary={step} />
                          </ListItem>
                        ))}
                      </List>
                    </AccordionDetails>
                  </Accordion>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    );
  };

  const renderAnalysisDashboard = () => {
    if (!analysis) return <Typography>Analysis data not available</Typography>;

    const skillData = {
      labels: Object.keys(analysis.skillBreakdown),
      datasets: [
        {
          label: 'Current Level',
          data: Object.values(analysis.skillBreakdown),
          borderColor: theme.palette.primary.main,
          backgroundColor: `${theme.palette.primary.main}30`,
          pointBackgroundColor: theme.palette.primary.main,
          pointBorderColor: theme.palette.primary.main,
        },
        {
          label: 'Potential Level',
          data: Object.keys(analysis.skillBreakdown).map(skill => 
            (analysis.skillBreakdown[skill] + (analysis.improvementPotential[skill] || 0))
          ),
          borderColor: theme.palette.success.main,
          backgroundColor: `${theme.palette.success.main}30`,
          pointBackgroundColor: theme.palette.success.main,
          pointBorderColor: theme.palette.success.main,
        },
      ],
    };

    return (
      <Grid container spacing={3}>
        {/* Overall Rating */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Overall Rating
              </Typography>
              <Box display="flex" alignItems="center">
                <CircularProgress
                  variant="determinate"
                  value={analysis.overallRating}
                  size={100}
                  thickness={6}
                />
                <Box ml={3}>
                  <Typography variant="h3" color="primary">
                    {Math.round(analysis.overallRating)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    out of 100
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Competitive Benchmark */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Competitive Benchmark
              </Typography>
              <Box mb={2}>
                <Typography variant="body2" color="text.secondary">
                  Rank: {analysis.competitiveBenchmark.rankTier} • 
                  Percentile: {analysis.competitiveBenchmark.regionalPercentile.toFixed(1)}%
                </Typography>
              </Box>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" color="success.main" gutterBottom>
                    Stronger Than Peers
                  </Typography>
                  <List dense>
                    {analysis.competitiveBenchmark.strongerThanPeers.map((strength, i) => (
                      <ListItem key={i}>
                        <ListItemIcon>
                          <TrendingUp color="success" fontSize="small" />
                        </ListItemIcon>
                        <ListItemText primary={strength.replace('_', ' ')} />
                      </ListItem>
                    ))}
                  </List>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" color="warning.main" gutterBottom>
                    Areas for Improvement
                  </Typography>
                  <List dense>
                    {analysis.competitiveBenchmark.weakerThanPeers.map((weakness, i) => (
                      <ListItem key={i}>
                        <ListItemIcon>
                          <TrendingDown color="warning" fontSize="small" />
                        </ListItemIcon>
                        <ListItemText primary={weakness.replace('_', ' ')} />
                      </ListItem>
                    ))}
                  </List>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Skill Breakdown Radar */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Skill Analysis
              </Typography>
              <Box height={400}>
                <Radar
                  data={skillData}
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

        {/* Critical Weaknesses */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Critical Weaknesses
              </Typography>
              {analysis.criticalWeaknesses.map((weakness, index) => (
                <Accordion key={index}>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Box display="flex" alignItems="center" width="100%">
                      <Typography sx={{ flexGrow: 1 }}>
                        {weakness.area.replace('_', ' ')}
                      </Typography>
                      <Chip
                        label={`${weakness.impactOnWinRate.toFixed(1)}% WR impact`}
                        size="small"
                        color="error"
                        sx={{ mr: 1 }}
                      />
                      <Chip
                        label={weakness.severity}
                        size="small"
                        style={{ backgroundColor: getPriorityColor(weakness.severity), color: 'white' }}
                      />
                    </Box>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          Root Causes:
                        </Typography>
                        <List dense>
                          {weakness.rootCauses.map((cause, i) => (
                            <ListItem key={i}>
                              <ListItemIcon>
                                <Warning color="error" fontSize="small" />
                              </ListItemIcon>
                              <ListItemText primary={cause} />
                            </ListItem>
                          ))}
                        </List>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="subtitle2" gutterBottom>
                          Quick Wins:
                        </Typography>
                        <List dense>
                          {weakness.quickWins.map((win, i) => (
                            <ListItem key={i}>
                              <ListItemIcon>
                                <Star color="warning" fontSize="small" />
                              </ListItemIcon>
                              <ListItemText primary={win} />
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

  const renderProgressTracking = () => {
    if (!progress) return <Typography>Progress data not available</Typography>;

    return (
      <Grid container spacing={3}>
        {/* Overall Progress Summary */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Progress Summary
              </Typography>
              <Box mb={2}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Improvement Score
                </Typography>
                <LinearProgress
                  variant="determinate"
                  value={progress.overallProgress.improvementScore}
                  sx={{ height: 8, borderRadius: 4, mb: 1 }}
                />
                <Typography variant="caption">
                  {progress.overallProgress.improvementScore.toFixed(1)}/100
                </Typography>
              </Box>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="h4" color="primary">
                    {progress.overallProgress.activeRecommendations}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Active
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="h4" color="success.main">
                    {progress.overallProgress.completedRecommendations}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Completed
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Recent Achievements */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Recent Achievements
              </Typography>
              {progress.recentAchievements.map((achievement, index) => (
                <Box key={index} display="flex" alignItems="center" mb={2}>
                  <EmojiEvents color="warning" sx={{ mr: 2 }} />
                  <Box>
                    <Typography variant="body2" fontWeight="bold">
                      {achievement.achievement}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {achievement.description} • {achievement.impact}
                    </Typography>
                  </Box>
                </Box>
              ))}
            </CardContent>
          </Card>
        </Grid>

        {/* Skill Progress */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Skill Progress Tracking
              </Typography>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Skill</TableCell>
                      <TableCell align="right">Baseline</TableCell>
                      <TableCell align="right">Current</TableCell>
                      <TableCell align="right">Target</TableCell>
                      <TableCell align="right">Progress</TableCell>
                      <TableCell align="right">Trend</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {Object.entries(progress.skillProgress).map(([skill, data]) => (
                      <TableRow key={skill}>
                        <TableCell component="th" scope="row">
                          {skill.replace('_', ' ')}
                        </TableCell>
                        <TableCell align="right">{data.baseline.toFixed(1)}</TableCell>
                        <TableCell align="right">{data.current.toFixed(1)}</TableCell>
                        <TableCell align="right">{data.target.toFixed(1)}</TableCell>
                        <TableCell align="right">
                          <Box display="flex" alignItems="center">
                            <LinearProgress
                              variant="determinate"
                              value={Math.max(0, Math.min(100, data.progress))}
                              sx={{ width: 60, mr: 1 }}
                            />
                            {data.progress.toFixed(1)}%
                          </Box>
                        </TableCell>
                        <TableCell align="right">
                          <Chip
                            size="small"
                            label={data.trend}
                            color={
                              data.trend === 'improving' ? 'success' :
                              data.trend === 'stable' ? 'default' : 'error'
                            }
                            icon={
                              data.trend === 'improving' ? <TrendingUp /> :
                              data.trend === 'declining' ? <TrendingDown /> : <Timeline />
                            }
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
      </Grid>
    );
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">
          <Psychology sx={{ mr: 1, verticalAlign: 'middle' }} />
          Improvement Recommendations
        </Typography>
        <Box display="flex" alignItems="center" gap={2}>
          <Badge badgeContent={recommendations.length} color="primary">
            <Assignment />
          </Badge>
          <Tooltip title="Refresh Recommendations">
            <IconButton onClick={loadImprovementData}>
              <Refresh />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={tabValue} onChange={handleTabChange} aria-label="improvement tabs">
          <Tab label="Recommendations" {...a11yProps(0)} />
          <Tab label="Quick Wins" {...a11yProps(1)} />
          <Tab label="Analysis" {...a11yProps(2)} />
          <Tab label="Insights" {...a11yProps(3)} />
          <Tab label="Progress" {...a11yProps(4)} />
        </Tabs>
      </Box>

      <TabPanel value={tabValue} index={0}>
        {renderRecommendations()}
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        {renderQuickWins()}
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        {renderAnalysisDashboard()}
      </TabPanel>

      <TabPanel value={tabValue} index={3}>
        {renderInsights()}
      </TabPanel>

      <TabPanel value={tabValue} index={4}>
        {renderProgressTracking()}
      </TabPanel>

      {/* Progress Update Dialog */}
      <Dialog
        open={progressDialog}
        onClose={() => setProgressDialog(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          Update Progress: {activeRecommendation?.title}
        </DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" paragraph>
            Track your progress on this recommendation to see your improvement over time.
          </Typography>
          {/* Progress update form would go here */}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setProgressDialog(false)}>Cancel</Button>
          <Button 
            variant="contained" 
            onClick={() => handleProgressUpdate({
              overallProgress: 25,
              completedSteps: [1],
              currentMilestone: 1,
              performanceMetrics: {},
            })}
          >
            Update Progress
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default ImprovementRecommendations;