import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Tabs,
  Tab,
  Grid,
  LinearProgress,
  Chip,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Avatar,
  IconButton,
  Tooltip,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Divider,
  Paper,
  Rating,
} from '@mui/material';
import {
  EmojiEvents,
  TrendingUp,
  TrendingDown,
  Build,
  Psychology,
  SportsEsports,
  Timeline,
  School,
  CompareArrows,
  Star,
  Shield,
  Bolt,
  Speed,
} from '@mui/icons-material';
import { Line, Radar, Bar, Doughnut } from 'react-chartjs-2';
import { championService } from '../../services/championService';
import type { 
  ChampionAnalysis, 
  ChampionTrendPoint, 
  ChampionMasteryRanking,
  ChampionComparisonData 
} from '../../types/champion';

interface ChampionAnalyticsProps {
  playerId: string;
  champion: string;
  timeRange?: string;
  position?: string;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel({ children, value, index }: TabPanelProps) {
  return (
    <div hidden={value !== index} style={{ paddingTop: '24px' }}>
      {value === index && children}
    </div>
  );
}

const ChampionAnalytics: React.FC<ChampionAnalyticsProps> = ({
  playerId,
  champion,
  timeRange = '30d',
  position = ''
}) => {
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [championAnalysis, setChampionAnalysis] = useState<ChampionAnalysis | null>(null);
  const [trendData, setTrendData] = useState<ChampionTrendPoint[]>([]);
  const [masteryRankings, setMasteryRankings] = useState<ChampionMasteryRanking[]>([]);
  const [selectedMetric, setSelectedMetric] = useState<string>('rating');

  useEffect(() => {
    fetchChampionData();
  }, [playerId, champion, timeRange, position]);

  const fetchChampionData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch comprehensive champion analysis
      const analysisData = await championService.getChampionAnalysis(
        playerId,
        timeRange,
        champion,
        position
      );
      setChampionAnalysis(analysisData);

      // Fetch trend data
      const trendsData = await championService.getChampionTrends(
        playerId,
        champion,
        'rating',
        'daily',
        30
      );
      setTrendData(trendsData);

      // Fetch mastery rankings
      const masteryData = await championService.getChampionMastery(
        playerId,
        timeRange,
        10
      );
      setMasteryRankings(masteryData);

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load champion analytics');
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const getRatingColor = (rating: number) => {
    if (rating >= 90) return '#4caf50';
    if (rating >= 80) return '#8bc34a';
    if (rating >= 70) return '#ff9800';
    if (rating >= 60) return '#ff5722';
    return '#f44336';
  };

  const getMasteryLevel = (points: number) => {
    if (points >= 500000) return { level: 'Master', color: '#9c27b0' };
    if (points >= 200000) return { level: 'Diamond', color: '#2196f3' };
    if (points >= 100000) return { level: 'Platinum', color: '#00bcd4' };
    if (points >= 50000) return { level: 'Gold', color: '#ff9800' };
    if (points >= 25000) return { level: 'Silver', color: '#9e9e9e' };
    return { level: 'Bronze', color: '#795548' };
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Loading Champion Analytics...
          </Typography>
          <LinearProgress />
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" color="error">
            Error: {error}
          </Typography>
        </CardContent>
      </Card>
    );
  }

  if (!championAnalysis) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6">
            No champion data available for {champion}.
          </Typography>
        </CardContent>
      </Card>
    );
  }

  const masteryLevel = getMasteryLevel(championAnalysis.masteryPoints);

  // Prepare chart data
  const performanceRadarData = {
    labels: ['Mechanics', 'Game Knowledge', 'Consistency', 'Adaptability', 'Team Fighting', 'Positioning'],
    datasets: [
      {
        label: 'Your Performance',
        data: [
          championAnalysis.mechanicsScore,
          championAnalysis.gameKnowledgeScore,
          championAnalysis.consistencyScore,
          championAnalysis.adaptabilityScore,
          championAnalysis.teamFightRating,
          championAnalysis.teamFightStats.positioningScore
        ],
        borderColor: '#1976d2',
        backgroundColor: 'rgba(25, 118, 210, 0.2)',
        pointBackgroundColor: '#1976d2',
        pointBorderColor: '#fff',
        pointHoverBackgroundColor: '#fff',
        pointHoverBorderColor: '#1976d2'
      }
    ]
  };

  const trendChartData = {
    labels: trendData.map(point => new Date(point.date).toLocaleDateString()),
    datasets: [
      {
        label: 'Overall Rating',
        data: trendData.map(point => point.overallRating),
        borderColor: '#1976d2',
        backgroundColor: 'rgba(25, 118, 210, 0.1)',
        tension: 0.4,
        fill: true
      },
      {
        label: 'Win Rate',
        data: trendData.map(point => point.winRate),
        borderColor: '#4caf50',
        backgroundColor: 'rgba(76, 175, 80, 0.1)',
        tension: 0.4,
        fill: false
      }
    ]
  };

  const powerSpikeData = {
    labels: championAnalysis.powerSpikes.map(spike => `Level ${spike.level}`),
    datasets: [
      {
        label: 'Power Rating',
        data: championAnalysis.powerSpikes.map(spike => spike.powerRating),
        backgroundColor: ['#ff6b6b', '#4ecdc4', '#45b7d1', '#96ceb4', '#feca57'],
        borderColor: ['#ff5252', '#26a69a', '#1976d2', '#66bb6a', '#ff9800'],
        borderWidth: 2
      }
    ]
  };

  const gamePhaseData = {
    labels: ['Early Game', 'Mid Game', 'Late Game'],
    datasets: [
      {
        label: 'Phase Rating',
        data: [
          championAnalysis.lanePhasePerformance.phaseRating,
          championAnalysis.midGamePerformance.phaseRating,
          championAnalysis.lateGamePerformance.phaseRating
        ],
        backgroundColor: ['#ff9800', '#2196f3', '#4caf50'],
        borderColor: ['#f57c00', '#1976d2', '#388e3c'],
        borderWidth: 2
      }
    ]
  };

  return (
    <Box>
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
            <Avatar
              sx={{
                width: 80,
                height: 80,
                bgcolor: masteryLevel.color,
                fontSize: '2rem'
              }}
            >
              {champion.charAt(0)}
            </Avatar>
            <Box sx={{ flexGrow: 1 }}>
              <Typography variant="h4" gutterBottom>
                {champion} Performance Analysis
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
                <Chip
                  label={`Level ${championAnalysis.masteryLevel}`}
                  color="primary"
                  variant="outlined"
                />
                <Chip
                  label={masteryLevel.level}
                  sx={{ bgcolor: masteryLevel.color, color: 'white' }}
                />
                <Chip
                  label={`${championAnalysis.masteryPoints.toLocaleString()} MP`}
                  color="secondary"
                  variant="outlined"
                />
              </Box>
              <Typography variant="body2" color="textSecondary">
                {championAnalysis.totalGames} games • {championAnalysis.playRate}% play rate • Recent: {championAnalysis.recentForm}
              </Typography>
            </Box>
            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="h3" sx={{ color: getRatingColor(championAnalysis.overallRating) }}>
                {championAnalysis.overallRating.toFixed(1)}
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Overall Rating
              </Typography>
              <Rating
                value={championAnalysis.overallRating / 20}
                readOnly
                precision={0.1}
                size="small"
              />
            </Box>
          </Box>

          <Grid container spacing={3} sx={{ mb: 3 }}>
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'success.light', color: 'white' }}>
                <Typography variant="h4">
                  {championAnalysis.winRate.toFixed(1)}%
                </Typography>
                <Typography variant="body2">Win Rate</Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'warning.light', color: 'white' }}>
                <Typography variant="h4">
                  {championAnalysis.championStats.averageKDA.toFixed(2)}
                </Typography>
                <Typography variant="body2">Average KDA</Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'info.light', color: 'white' }}>
                <Typography variant="h4">
                  {championAnalysis.carryPotential.toFixed(1)}
                </Typography>
                <Typography variant="body2">Carry Potential</Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'secondary.light', color: 'white' }}>
                <Typography variant="h4">
                  {championAnalysis.clutchFactor.toFixed(1)}
                </Typography>
                <Typography variant="body2">Clutch Factor</Typography>
              </Paper>
            </Grid>
          </Grid>

          <Tabs value={tabValue} onChange={handleTabChange} variant="scrollable" scrollButtons="auto">
            <Tab label="Overview" icon={<SportsEsports />} />
            <Tab label="Performance" icon={<TrendingUp />} />
            <Tab label="Power Spikes" icon={<Bolt />} />
            <Tab label="Builds & Runes" icon={<Build />} />
            <Tab label="Matchups" icon={<CompareArrows />} />
            <Tab label="Coaching" icon={<School />} />
            <Tab label="Trends" icon={<Timeline />} />
          </Tabs>

          <TabPanel value={tabValue} index={0}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={8}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Performance Radar
                    </Typography>
                    <Box sx={{ height: 400 }}>
                      <Radar 
                        data={performanceRadarData} 
                        options={{
                          responsive: true,
                          maintainAspectRatio: false,
                          scales: {
                            r: {
                              beginAtZero: true,
                              max: 100,
                              ticks: { stepSize: 20 }
                            }
                          },
                          plugins: {
                            legend: { display: false }
                          }
                        }} 
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={4}>
                <Card variant="outlined" sx={{ mb: 2 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Star color="primary" />
                      Core Strengths
                    </Typography>
                    <List dense>
                      {championAnalysis.coreStrengths.map((strength, index) => (
                        <ListItem key={index}>
                          <ListItemIcon>
                            <TrendingUp color="success" />
                          </ListItemIcon>
                          <ListItemText 
                            primary={strength.title}
                            secondary={strength.description}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
                
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Psychology color="warning" />
                      Improvement Areas
                    </Typography>
                    <List dense>
                      {championAnalysis.improvementAreas.map((area, index) => (
                        <ListItem key={index}>
                          <ListItemIcon>
                            <TrendingDown color="warning" />
                          </ListItemIcon>
                          <ListItemText 
                            primary={area.title}
                            secondary={area.description}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={1}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Game Phase Performance
                    </Typography>
                    <Box sx={{ height: 300 }}>
                      <Bar 
                        data={gamePhaseData} 
                        options={{
                          responsive: true,
                          maintainAspectRatio: false,
                          scales: {
                            y: {
                              beginAtZero: true,
                              max: 100,
                              title: { display: true, text: 'Rating' }
                            }
                          },
                          plugins: {
                            legend: { display: false }
                          }
                        }} 
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Champion Statistics
                    </Typography>
                    
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Box sx={{ textAlign: 'center', p: 1 }}>
                          <Typography variant="h5" color="primary">
                            {championAnalysis.championStats.csPerMinute.toFixed(1)}
                          </Typography>
                          <Typography variant="body2">CS/min</Typography>
                        </Box>
                      </Grid>
                      <Grid item xs={6}>
                        <Box sx={{ textAlign: 'center', p: 1 }}>
                          <Typography variant="h5" color="primary">
                            {championAnalysis.championStats.damagePerMinute.toFixed(0)}
                          </Typography>
                          <Typography variant="body2">DPM</Typography>
                        </Box>
                      </Grid>
                      <Grid item xs={6}>
                        <Box sx={{ textAlign: 'center', p: 1 }}>
                          <Typography variant="h5" color="primary">
                            {championAnalysis.championStats.goldPerMinute.toFixed(0)}
                          </Typography>
                          <Typography variant="body2">GPM</Typography>
                        </Box>
                      </Grid>
                      <Grid item xs={6}>
                        <Box sx={{ textAlign: 'center', p: 1 }}>
                          <Typography variant="h5" color="primary">
                            {championAnalysis.championStats.killParticipation.toFixed(1)}%
                          </Typography>
                          <Typography variant="body2">KP</Typography>
                        </Box>
                      </Grid>
                    </Grid>

                    <Divider sx={{ my: 2 }} />

                    <Typography variant="subtitle1" gutterBottom>
                      Team Fight Analysis
                    </Typography>
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" gutterBottom>
                        Participation: {championAnalysis.teamFightStats.participationRate.toFixed(1)}%
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={championAnalysis.teamFightStats.participationRate} 
                        sx={{ mb: 1 }}
                      />
                    </Box>
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" gutterBottom>
                        Survival Rate: {championAnalysis.teamFightStats.survivalRate.toFixed(1)}%
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={championAnalysis.teamFightStats.survivalRate} 
                        sx={{ mb: 1 }}
                        color="success"
                      />
                    </Box>
                    <Box>
                      <Typography variant="body2" gutterBottom>
                        Damage Contribution: {championAnalysis.teamFightStats.damageContribution.toFixed(1)}%
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={championAnalysis.teamFightStats.damageContribution} 
                        color="warning"
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={2}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={8}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Power Spike Analysis
                    </Typography>
                    <Box sx={{ height: 300 }}>
                      <Bar 
                        data={powerSpikeData} 
                        options={{
                          responsive: true,
                          maintainAspectRatio: false,
                          scales: {
                            y: {
                              beginAtZero: true,
                              max: 100,
                              title: { display: true, text: 'Power Rating' }
                            }
                          },
                          plugins: {
                            legend: { display: false }
                          }
                        }} 
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={4}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Power Spike Details
                    </Typography>
                    <List>
                      {championAnalysis.powerSpikes.map((spike, index) => (
                        <ListItem key={index} sx={{ flexDirection: 'column', alignItems: 'flex-start', p: 2, border: '1px solid', borderColor: 'divider', borderRadius: 1, mb: 1 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}>
                            <Bolt color="warning" />
                            Level {spike.level}
                          </Typography>
                          <Typography variant="body2" color="textSecondary" gutterBottom>
                            {spike.itemThreshold}
                          </Typography>
                          <Box sx={{ display: 'flex', justifyContent: 'space-between', width: '100%', mt: 1 }}>
                            <Typography variant="caption">
                              Power: {spike.powerRating.toFixed(1)}
                            </Typography>
                            <Typography variant="caption">
                              Win Rate: +{spike.winRateIncrease.toFixed(1)}%
                            </Typography>
                          </Box>
                          <Box sx={{ width: '100%', mt: 1 }}>
                            <Typography variant="caption" color="textSecondary">
                              Optimal Timing: {Math.floor(spike.optimalTiming / 60)}:{(spike.optimalTiming % 60).toString().padStart(2, '0')}
                            </Typography>
                          </Box>
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={3}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Item Build Analysis
                    </Typography>
                    
                    <Box sx={{ mb: 3 }}>
                      <Typography variant="subtitle1" gutterBottom>
                        Most Successful Build
                      </Typography>
                      {championAnalysis.itemBuilds.mostSuccessfulBuild.map((build, index) => (
                        <Box key={index} sx={{ p: 2, bgcolor: 'background.default', borderRadius: 1, mb: 2 }}>
                          <Typography variant="body2" gutterBottom>
                            Items: {build.items.join(' → ')}
                          </Typography>
                          <Box sx={{ display: 'flex', gap: 2 }}>
                            <Chip label={`${build.winRate.toFixed(1)}% WR`} color="success" size="small" />
                            <Chip label={`${build.playRate.toFixed(1)}% PR`} color="info" size="small" />
                          </Box>
                        </Box>
                      ))}
                    </Box>

                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" gutterBottom>
                        Build Adaptability: {championAnalysis.itemBuilds.adaptabilityScore.toFixed(1)}%
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={championAnalysis.itemBuilds.adaptabilityScore} 
                        sx={{ mb: 1 }}
                      />
                    </Box>

                    <Typography variant="subtitle2" gutterBottom>
                      Core Item Timing
                    </Typography>
                    {Object.entries(championAnalysis.itemBuilds.coreItemTiming).map(([item, timing]) => (
                      <Box key={item} sx={{ display: 'flex', justifyContent: 'space-between', py: 0.5 }}>
                        <Typography variant="body2">{item}</Typography>
                        <Typography variant="body2" color="primary">{timing.toFixed(1)} min</Typography>
                      </Box>
                    ))}
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Card variant="outlined" sx={{ mb: 2 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Rune Optimization
                    </Typography>
                    
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="subtitle1" gutterBottom>
                        Most Successful Setup
                      </Typography>
                      <Box sx={{ p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
                        <Typography variant="body2" gutterBottom>
                          <strong>Primary:</strong> {championAnalysis.runeOptimization.mostSuccessfulSetup.primaryTree}
                        </Typography>
                        <Typography variant="body2" gutterBottom>
                          <strong>Secondary:</strong> {championAnalysis.runeOptimization.mostSuccessfulSetup.secondaryTree}
                        </Typography>
                        <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
                          <Chip 
                            label={`${championAnalysis.runeOptimization.mostSuccessfulSetup.winRate.toFixed(1)}% WR`} 
                            color="success" 
                            size="small" 
                          />
                          <Chip 
                            label={`${championAnalysis.runeOptimization.mostSuccessfulSetup.playRate.toFixed(1)}% PR`} 
                            color="info" 
                            size="small" 
                          />
                        </Box>
                      </Box>
                    </Box>

                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" gutterBottom>
                        Keystone Optimality: {championAnalysis.runeOptimization.keystoneOptimality.toFixed(1)}%
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={championAnalysis.runeOptimization.keystoneOptimality} 
                        color="primary"
                      />
                    </Box>
                  </CardContent>
                </Card>
                
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Skill Order Analysis
                    </Typography>
                    
                    <Typography variant="subtitle2" gutterBottom>
                      Most Common Order
                    </Typography>
                    <Typography variant="body1" sx={{ 
                      p: 1, 
                      bgcolor: 'primary.light', 
                      color: 'white', 
                      borderRadius: 1, 
                      textAlign: 'center', 
                      mb: 2,
                      fontFamily: 'monospace',
                      fontSize: '1.2rem'
                    }}>
                      {championAnalysis.skillOrder.mostCommonOrder}
                    </Typography>
                    
                    <Typography variant="subtitle2" gutterBottom>
                      Optimal Order
                    </Typography>
                    <Typography variant="body1" sx={{ 
                      p: 1, 
                      bgcolor: 'success.light', 
                      color: 'white', 
                      borderRadius: 1, 
                      textAlign: 'center',
                      fontFamily: 'monospace',
                      fontSize: '1.2rem'
                    }}>
                      {championAnalysis.skillOrder.optimalOrder}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={4}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ color: 'success.main' }}>
                      Strength Matchups
                    </Typography>
                    <List>
                      {championAnalysis.strengthMatchups.map((matchup, index) => (
                        <ListItem key={index} sx={{ border: '1px solid', borderColor: 'success.light', borderRadius: 1, mb: 1 }}>
                          <ListItemText 
                            primary={
                              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <Typography variant="subtitle1">{matchup.opponentChampion}</Typography>
                                <Chip 
                                  label={`${matchup.winRate.toFixed(1)}% WR`} 
                                  color="success" 
                                  size="small" 
                                />
                              </Box>
                            }
                            secondary={
                              <Box>
                                <Typography variant="body2" gutterBottom>
                                  {matchup.gamesPlayed} games played
                                </Typography>
                                <Typography variant="body2">
                                  Lane Phase: {matchup.lanePhasePerformance.toFixed(1)} | 
                                  CS Advantage: +{matchup.averageCSAdvantage.toFixed(1)}
                                </Typography>
                              </Box>
                            }
                          />
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ color: 'error.main' }}>
                      Challenging Matchups
                    </Typography>
                    <List>
                      {championAnalysis.weaknessMatchups.map((matchup, index) => (
                        <ListItem key={index} sx={{ border: '1px solid', borderColor: 'error.light', borderRadius: 1, mb: 1 }}>
                          <ListItemText 
                            primary={
                              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <Typography variant="subtitle1">{matchup.opponentChampion}</Typography>
                                <Chip 
                                  label={`${matchup.winRate.toFixed(1)}% WR`} 
                                  color="error" 
                                  size="small" 
                                />
                              </Box>
                            }
                            secondary={
                              <Box>
                                <Typography variant="body2" gutterBottom>
                                  {matchup.gamesPlayed} games played
                                </Typography>
                                <Typography variant="body2">
                                  Lane Phase: {matchup.lanePhasePerformance.toFixed(1)} | 
                                  CS Advantage: {matchup.averageCSAdvantage.toFixed(1)}
                                </Typography>
                                {matchup.commonMistakes && matchup.commonMistakes.length > 0 && (
                                  <Typography variant="caption" color="error">
                                    Common mistakes: {matchup.commonMistakes.join(', ')}
                                  </Typography>
                                )}
                              </Box>
                            }
                          />
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={5}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Play Style Recommendations
                    </Typography>
                    <List>
                      {championAnalysis.playStyleRecommendations.map((rec, index) => (
                        <ListItem key={index} sx={{ flexDirection: 'column', alignItems: 'flex-start', border: '1px solid', borderColor: 'divider', borderRadius: 1, mb: 2 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
                            {rec.title}
                          </Typography>
                          <Typography variant="body2" color="textSecondary" gutterBottom>
                            {rec.description}
                          </Typography>
                          <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
                            <Chip label={rec.priority} color={rec.priority === 'high' ? 'error' : rec.priority === 'medium' ? 'warning' : 'default'} size="small" />
                            <Chip label={rec.difficulty} variant="outlined" size="small" />
                            <Chip label={`+${rec.expectedImprovement.toFixed(1)} points`} color="success" size="small" />
                          </Box>
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Training Recommendations
                    </Typography>
                    <List>
                      {championAnalysis.trainingRecommendations.map((training, index) => (
                        <ListItem key={index} sx={{ flexDirection: 'column', alignItems: 'flex-start', border: '1px solid', borderColor: 'divider', borderRadius: 1, mb: 2 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
                            {training.title}
                          </Typography>
                          <Typography variant="body2" color="textSecondary" gutterBottom>
                            {training.description}
                          </Typography>
                          <Box sx={{ display: 'flex', gap: 1, mb: 1 }}>
                            <Chip label={training.duration} variant="outlined" size="small" />
                            <Chip label={training.frequency} variant="outlined" size="small" />
                          </Box>
                          <Typography variant="caption" color="primary">
                            Expected improvement timeline: {training.expectedTimeline}
                          </Typography>
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>

            <Grid item xs={12}>
              <Card variant="outlined" sx={{ mt: 3 }}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    Learning Curve Analysis
                  </Typography>
                  <Grid container spacing={3}>
                    <Grid item xs={12} md={4}>
                      <Box sx={{ textAlign: 'center' }}>
                        <Typography variant="h4" color="primary">
                          {championAnalysis.learningCurve.currentStage}
                        </Typography>
                        <Typography variant="body2">Current Stage</Typography>
                      </Box>
                    </Grid>
                    <Grid item xs={12} md={4}>
                      <Box sx={{ textAlign: 'center' }}>
                        <Typography variant="h4" color="primary">
                          {championAnalysis.learningCurve.progressScore.toFixed(1)}%
                        </Typography>
                        <Typography variant="body2">Progress Score</Typography>
                      </Box>
                    </Grid>
                    <Grid item xs={12} md={4}>
                      <Box sx={{ textAlign: 'center' }}>
                        <Typography variant="h4" color="primary">
                          {championAnalysis.learningCurve.estimatedTimeToMastery}
                        </Typography>
                        <Typography variant="body2">Games to Mastery</Typography>
                      </Box>
                    </Grid>
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={6}>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Performance Trends
                    </Typography>
                    
                    <FormControl sx={{ mb: 3, minWidth: 200 }}>
                      <InputLabel>Trend Metric</InputLabel>
                      <Select
                        value={selectedMetric}
                        onChange={(e) => setSelectedMetric(e.target.value)}
                      >
                        <MenuItem value="rating">Overall Rating</MenuItem>
                        <MenuItem value="winrate">Win Rate</MenuItem>
                        <MenuItem value="kda">KDA</MenuItem>
                        <MenuItem value="dpm">Damage Per Minute</MenuItem>
                        <MenuItem value="cs">CS Per Minute</MenuItem>
                      </Select>
                    </FormControl>
                    
                    <Box sx={{ height: 400 }}>
                      <Line 
                        data={trendChartData} 
                        options={{
                          responsive: true,
                          maintainAspectRatio: false,
                          scales: {
                            y: {
                              beginAtZero: false,
                              title: { display: true, text: 'Value' }
                            },
                            x: {
                              title: { display: true, text: 'Date' }
                            }
                          },
                          plugins: {
                            legend: { position: 'top' }
                          }
                        }} 
                      />
                    </Box>
                    
                    <Box sx={{ mt: 2, display: 'flex', justifyContent: 'center', gap: 2 }}>
                      <Chip 
                        label={`Trend: ${championAnalysis.trendDirection}`}
                        color={championAnalysis.trendDirection === 'improving' ? 'success' : championAnalysis.trendDirection === 'declining' ? 'error' : 'default'}
                        icon={championAnalysis.trendDirection === 'improving' ? <TrendingUp /> : championAnalysis.trendDirection === 'declining' ? <TrendingDown /> : <Speed />}
                      />
                      <Chip 
                        label={`Confidence: ${(championAnalysis.trendConfidence * 100).toFixed(1)}%`}
                        variant="outlined"
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>
        </CardContent>
      </Card>
    </Box>
  );
};

export default ChampionAnalytics;