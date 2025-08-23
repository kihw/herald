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
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Collapse,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  TrendingFlat,
  Star,
  Block,
  Visibility,
  ExpandMore,
  ExpandLess,
  Psychology,
  Timeline,
  Assessment,
  Lightbulb,
  AutoGraph,
} from '@mui/icons-material';
import { Bar, Line, Doughnut, Radar } from 'react-chartjs-2';
import { metaService } from '../../services/metaService';
import type {
  MetaAnalysis,
  ChampionTierList,
  MetaTrends,
  ChampionMetaStats,
  MetaPredictions
} from '../../types/meta';

interface MetaAnalyticsProps {
  patch?: string;
  region?: string;
  rank?: string;
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

const MetaAnalytics: React.FC<MetaAnalyticsProps> = ({
  patch = '14.1',
  region = 'all',
  rank = 'all'
}) => {
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [metaAnalysis, setMetaAnalysis] = useState<MetaAnalysis | null>(null);
  const [selectedRole, setSelectedRole] = useState<string>('ALL');
  const [selectedTier, setSelectedTier] = useState<string>('all');
  const [expandedTiers, setExpandedTiers] = useState<Record<string, boolean>>({});

  useEffect(() => {
    fetchMetaData();
  }, [patch, region, rank, selectedRole]);

  const fetchMetaData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch comprehensive meta analysis
      const analysisData = await metaService.getMetaAnalysis(
        patch,
        region,
        rank,
        '7d'
      );
      setMetaAnalysis(analysisData);

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load meta analytics');
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const getTierColor = (tier: string) => {
    switch (tier.toLowerCase()) {
      case 's+': return '#ff4444';
      case 's': return '#ff6666';
      case 'a+': return '#ff9800';
      case 'a': return '#ffab00';
      case 'b+': return '#ffc107';
      case 'b': return '#ffeb3b';
      case 'c+': return '#8bc34a';
      case 'c': return '#4caf50';
      case 'd': return '#9e9e9e';
      default: return '#757575';
    }
  };

  const getTrendIcon = (direction: string) => {
    switch (direction) {
      case 'rising': return <TrendingUp color="success" />;
      case 'falling': return <TrendingDown color="error" />;
      case 'declining': return <TrendingDown color="error" />;
      default: return <TrendingFlat color="info" />;
    }
  };

  const toggleTierExpansion = (tier: string) => {
    setExpandedTiers(prev => ({
      ...prev,
      [tier]: !prev[tier]
    }));
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Loading Meta Analytics...
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

  if (!metaAnalysis) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6">
            No meta data available for patch {patch}.
          </Typography>
        </CardContent>
      </Card>
    );
  }

  // Prepare chart data
  const tierDistributionData = {
    labels: ['S+', 'S', 'A+', 'A', 'B+', 'B', 'C+', 'C', 'D'],
    datasets: [
      {
        label: 'Champions',
        data: [
          metaAnalysis.tierList.sPlusTier?.length || 0,
          metaAnalysis.tierList.sTier?.length || 0,
          metaAnalysis.tierList.aPlusTier?.length || 0,
          metaAnalysis.tierList.aTier?.length || 0,
          metaAnalysis.tierList.bPlusTier?.length || 0,
          metaAnalysis.tierList.bTier?.length || 0,
          metaAnalysis.tierList.cPlusTier?.length || 0,
          metaAnalysis.tierList.cTier?.length || 0,
          metaAnalysis.tierList.dTier?.length || 0,
        ],
        backgroundColor: [
          '#ff4444', '#ff6666', '#ff9800', '#ffab00',
          '#ffc107', '#ffeb3b', '#8bc34a', '#4caf50', '#9e9e9e'
        ],
        borderColor: [
          '#d32f2f', '#f44336', '#f57c00', '#ff8f00',
          '#f9a825', '#f57f17', '#689f38', '#388e3c', '#616161'
        ],
        borderWidth: 2
      }
    ]
  };

  const championTypeData = {
    labels: ['Tanks', 'Fighters', 'Assassins', 'Mages', 'Marksmen', 'Support'],
    datasets: [
      {
        label: 'Pick Rate',
        data: [
          metaAnalysis.metaTrends.championTypes.tanks.pickRate,
          metaAnalysis.metaTrends.championTypes.fighters.pickRate,
          metaAnalysis.metaTrends.championTypes.assassins.pickRate,
          metaAnalysis.metaTrends.championTypes.mages.pickRate,
          metaAnalysis.metaTrends.championTypes.marksmen.pickRate,
          metaAnalysis.metaTrends.championTypes.support.pickRate
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

  const renderTierSection = (tierName: string, champions: any[], tierColor: string) => {
    if (!champions || champions.length === 0) return null;

    return (
      <Card key={tierName} variant="outlined" sx={{ mb: 1 }}>
        <Box
          sx={{
            bgcolor: tierColor,
            color: 'white',
            p: 2,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            cursor: 'pointer'
          }}
          onClick={() => toggleTierExpansion(tierName)}
        >
          <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
            {tierName} Tier ({champions.length} champions)
          </Typography>
          <IconButton sx={{ color: 'white' }}>
            {expandedTiers[tierName] ? <ExpandLess /> : <ExpandMore />}
          </IconButton>
        </Box>
        <Collapse in={expandedTiers[tierName]}>
          <List>
            {champions.map((champion, index) => (
              <ListItem key={index} divider={index < champions.length - 1}>
                <ListItemIcon>
                  <Avatar sx={{ width: 32, height: 32, bgcolor: tierColor }}>
                    {champion.champion.charAt(0)}
                  </Avatar>
                </ListItemIcon>
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
                        {champion.champion}
                      </Typography>
                      {getTrendIcon(champion.trendDirection)}
                      {champion.tierMovement !== 0 && (
                        <Chip
                          label={champion.tierMovement > 0 ? `+${champion.tierMovement}` : champion.tierMovement}
                          color={champion.tierMovement > 0 ? 'success' : 'error'}
                          size="small"
                        />
                      )}
                    </Box>
                  }
                  secondary={
                    <Box sx={{ display: 'flex', gap: 2, mt: 1 }}>
                      <Chip
                        label={`${champion.winRate.toFixed(1)}% WR`}
                        color="success"
                        size="small"
                        variant="outlined"
                      />
                      <Chip
                        label={`${champion.pickRate.toFixed(1)}% PR`}
                        color="primary"
                        size="small"
                        variant="outlined"
                      />
                      <Chip
                        label={`${champion.banRate.toFixed(1)}% BR`}
                        color="error"
                        size="small"
                        variant="outlined"
                      />
                      <Chip
                        label={`${champion.carryPotential.toFixed(0)} CP`}
                        color="secondary"
                        size="small"
                        variant="outlined"
                      />
                    </Box>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Collapse>
      </Card>
    );
  };

  return (
    <Box>
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Assessment color="primary" />
              Meta Analytics - Patch {patch}
            </Typography>
            
            <Box sx={{ display: 'flex', gap: 2 }}>
              <FormControl size="small" sx={{ minWidth: 120 }}>
                <InputLabel>Role</InputLabel>
                <Select
                  value={selectedRole}
                  onChange={(e) => setSelectedRole(e.target.value)}
                >
                  <MenuItem value="ALL">All Roles</MenuItem>
                  <MenuItem value="TOP">Top</MenuItem>
                  <MenuItem value="JUNGLE">Jungle</MenuItem>
                  <MenuItem value="MID">Mid</MenuItem>
                  <MenuItem value="ADC">ADC</MenuItem>
                  <MenuItem value="SUPPORT">Support</MenuItem>
                </Select>
              </FormControl>
            </Box>
          </Box>

          <Grid container spacing={3} sx={{ mb: 3 }}>
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'primary.light', color: 'white' }}>
                <Typography variant="h4">
                  {metaAnalysis.tierList.confidence.toFixed(1)}%
                </Typography>
                <Typography variant="body2">Tier List Confidence</Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'success.light', color: 'white' }}>
                <Typography variant="h4">
                  {metaAnalysis.dataQuality.toFixed(1)}%
                </Typography>
                <Typography variant="body2">Data Quality</Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'warning.light', color: 'white' }}>
                <Typography variant="h4">
                  {metaAnalysis.emergingPicks.length}
                </Typography>
                <Typography variant="body2">Emerging Picks</Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'error.light', color: 'white' }}>
                <Typography variant="h4">
                  {metaAnalysis.deciningPicks.length}
                </Typography>
                <Typography variant="body2">Declining Picks</Typography>
              </Paper>
            </Grid>
          </Grid>

          <Tabs value={tabValue} onChange={handleTabChange} variant="scrollable" scrollButtons="auto">
            <Tab label="Tier List" icon={<Star />} />
            <Tab label="Meta Trends" icon={<AutoGraph />} />
            <Tab label="Ban/Pick Analysis" icon={<Block />} />
            <Tab label="Champion Stats" icon={<Assessment />} />
            <Tab label="Predictions" icon={<Psychology />} />
            <Tab label="Recommendations" icon={<Lightbulb />} />
          </Tabs>

          <TabPanel value={tabValue} index={0}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={8}>
                <Box>
                  {renderTierSection('S+', metaAnalysis.tierList.sPlusTier, '#ff4444')}
                  {renderTierSection('S', metaAnalysis.tierList.sTier, '#ff6666')}
                  {renderTierSection('A+', metaAnalysis.tierList.aPlusTier, '#ff9800')}
                  {renderTierSection('A', metaAnalysis.tierList.aTier, '#ffab00')}
                  {renderTierSection('B+', metaAnalysis.tierList.bPlusTier, '#ffc107')}
                  {renderTierSection('B', metaAnalysis.tierList.bTier, '#ffeb3b')}
                  {renderTierSection('C+', metaAnalysis.tierList.cPlusTier, '#8bc34a')}
                  {renderTierSection('C', metaAnalysis.tierList.cTier, '#4caf50')}
                  {renderTierSection('D', metaAnalysis.tierList.dTier, '#9e9e9e')}
                </Box>
              </Grid>
              
              <Grid item xs={12} md={4}>
                <Card variant="outlined" sx={{ mb: 2 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Tier Distribution
                    </Typography>
                    <Box sx={{ height: 300 }}>
                      <Bar 
                        data={tierDistributionData} 
                        options={{
                          responsive: true,
                          maintainAspectRatio: false,
                          scales: {
                            y: {
                              beginAtZero: true,
                              title: { display: true, text: 'Champions' }
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

                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Tier List Info
                    </Typography>
                    <List dense>
                      <ListItem>
                        <ListItemText 
                          primary="Sample Size" 
                          secondary={metaAnalysis.tierList.sampleSize.toLocaleString()} 
                        />
                      </ListItem>
                      <ListItem>
                        <ListItemText 
                          primary="Last Updated" 
                          secondary={new Date(metaAnalysis.tierList.lastUpdated).toLocaleDateString()} 
                        />
                      </ListItem>
                      <ListItem>
                        <ListItemText 
                          primary="Confidence" 
                          secondary={`${metaAnalysis.tierList.confidence.toFixed(1)}%`} 
                        />
                      </ListItem>
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
                      Champion Type Trends
                    </Typography>
                    <Box sx={{ height: 300 }}>
                      <Radar 
                        data={championTypeData} 
                        options={{
                          responsive: true,
                          maintainAspectRatio: false,
                          scales: {
                            r: {
                              beginAtZero: true,
                              max: 30,
                              ticks: { stepSize: 5 }
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
                      Dominant Strategies
                    </Typography>
                    <List>
                      {metaAnalysis.metaTrends.dominantStrategies.map((strategy, index) => (
                        <ListItem key={index} sx={{ flexDirection: 'column', alignItems: 'flex-start', border: '1px solid', borderColor: 'divider', borderRadius: 1, mb: 1 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}>
                            {getTrendIcon(strategy.trend)}
                            {strategy.strategy}
                          </Typography>
                          <Typography variant="body2" color="textSecondary" gutterBottom>
                            {strategy.description}
                          </Typography>
                          <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
                            <Chip 
                              label={`${strategy.popularity.toFixed(1)}% Popular`} 
                              color="primary" 
                              size="small" 
                            />
                            <Chip 
                              label={`${strategy.winRate.toFixed(1)}% WR`} 
                              color="success" 
                              size="small" 
                            />
                            <Chip 
                              label={strategy.trend} 
                              color={strategy.trend === 'rising' ? 'success' : strategy.trend === 'falling' ? 'error' : 'default'} 
                              size="small" 
                            />
                          </Box>
                          <Typography variant="caption" sx={{ mt: 1 }}>
                            Key Champions: {strategy.champions.join(', ')}
                          </Typography>
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>

            <Grid container spacing={3} sx={{ mt: 2 }}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ color: 'success.main' }}>
                      Emerging Picks
                    </Typography>
                    <List>
                      {metaAnalysis.emergingPicks.map((pick, index) => (
                        <ListItem key={index} sx={{ border: '1px solid', borderColor: 'success.light', borderRadius: 1, mb: 1 }}>
                          <ListItemIcon>
                            <TrendingUp color="success" />
                          </ListItemIcon>
                          <ListItemText 
                            primary={
                              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <Typography variant="subtitle1">{pick.champion}</Typography>
                                <Box sx={{ display: 'flex', gap: 1 }}>
                                  <Chip label={pick.role} color="primary" size="small" />
                                  <Chip label={`${pick.currentTier} → ${pick.projectedTier}`} color="success" size="small" />
                                </Box>
                              </Box>
                            }
                            secondary={
                              <Box>
                                <Typography variant="body2" gutterBottom>
                                  Reasons: {pick.reasonForRise.join(', ')}
                                </Typography>
                                <Typography variant="caption">
                                  Confidence: {(pick.confidence * 100).toFixed(0)}%
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
                      Declining Picks
                    </Typography>
                    <List>
                      {metaAnalysis.deciningPicks.map((pick, index) => (
                        <ListItem key={index} sx={{ border: '1px solid', borderColor: 'error.light', borderRadius: 1, mb: 1 }}>
                          <ListItemIcon>
                            <TrendingDown color="error" />
                          </ListItemIcon>
                          <ListItemText 
                            primary={
                              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <Typography variant="subtitle1">{pick.champion}</Typography>
                                <Box sx={{ display: 'flex', gap: 1 }}>
                                  <Chip label={pick.role} color="primary" size="small" />
                                  <Chip label={`${pick.previousTier} → ${pick.projectedTier}`} color="error" size="small" />
                                </Box>
                              </Box>
                            }
                            secondary={
                              <Box>
                                <Typography variant="body2" gutterBottom>
                                  Reasons: {pick.reasonForDecline.join(', ')}
                                </Typography>
                                <Typography variant="caption">
                                  Confidence: {(pick.confidence * 100).toFixed(0)}%
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
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={2}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Top Banned Champions
                    </Typography>
                    <TableContainer>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Champion</TableCell>
                            <TableCell align="right">Ban Rate</TableCell>
                            <TableCell align="right">Threat Level</TableCell>
                            <TableCell align="right">Priority</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {metaAnalysis.banAnalysis.topBannedChampions.map((champion, index) => (
                            <TableRow key={index}>
                              <TableCell>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                  <Avatar sx={{ width: 24, height: 24 }}>{champion.champion.charAt(0)}</Avatar>
                                  {champion.champion}
                                </Box>
                              </TableCell>
                              <TableCell align="right">
                                <Chip 
                                  label={`${champion.banRate.toFixed(1)}%`} 
                                  color="error" 
                                  size="small" 
                                />
                              </TableCell>
                              <TableCell align="right">
                                <LinearProgress 
                                  variant="determinate" 
                                  value={champion.threatLevel} 
                                  sx={{ width: 60, height: 6 }}
                                  color="error"
                                />
                              </TableCell>
                              <TableCell align="right">
                                {champion.banPriority.toFixed(1)}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Role Targeting in Bans
                    </Typography>
                    <Box>
                      {Object.entries(metaAnalysis.banAnalysis.roleTargeting).map(([role, percentage]) => (
                        <Box key={role} sx={{ mb: 2 }}>
                          <Typography variant="body2" gutterBottom>
                            {role}: {percentage.toFixed(1)}%
                          </Typography>
                          <LinearProgress 
                            variant="determinate" 
                            value={percentage} 
                            sx={{ height: 8, borderRadius: 4 }}
                          />
                        </Box>
                      ))}
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={3}>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Champion Performance Overview
                    </Typography>
                    <TableContainer>
                      <Table>
                        <TableHead>
                          <TableRow>
                            <TableCell>Champion</TableCell>
                            <TableCell>Role</TableCell>
                            <TableCell>Tier</TableCell>
                            <TableCell align="right">Win Rate</TableCell>
                            <TableCell align="right">Pick Rate</TableCell>
                            <TableCell align="right">Ban Rate</TableCell>
                            <TableCell align="right">Presence</TableCell>
                            <TableCell>Trend</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {metaAnalysis.championStats.slice(0, 10).map((champion, index) => (
                            <TableRow key={index} hover>
                              <TableCell>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                  <Avatar sx={{ width: 32, height: 32, bgcolor: getTierColor(champion.tier) }}>
                                    {champion.champion.charAt(0)}
                                  </Avatar>
                                  {champion.champion}
                                </Box>
                              </TableCell>
                              <TableCell>{champion.role}</TableCell>
                              <TableCell>
                                <Chip 
                                  label={champion.tier} 
                                  sx={{ bgcolor: getTierColor(champion.tier), color: 'white' }}
                                  size="small"
                                />
                              </TableCell>
                              <TableCell align="right">
                                <Chip 
                                  label={`${champion.winRate.toFixed(1)}%`} 
                                  color={champion.winRate >= 52 ? 'success' : champion.winRate >= 48 ? 'warning' : 'error'}
                                  size="small"
                                />
                              </TableCell>
                              <TableCell align="right">{champion.pickRate.toFixed(1)}%</TableCell>
                              <TableCell align="right">{champion.banRate.toFixed(1)}%</TableCell>
                              <TableCell align="right">{champion.presenceRate.toFixed(1)}%</TableCell>
                              <TableCell>{getTrendIcon(champion.trendDirection)}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
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
                    <Typography variant="h6" gutterBottom>
                      Next Patch Predictions
                    </Typography>
                    <List>
                      {metaAnalysis.predictions.nextPatchPredictions.map((prediction, index) => (
                        <ListItem key={index} sx={{ flexDirection: 'column', alignItems: 'flex-start', border: '1px solid', borderColor: 'divider', borderRadius: 1, mb: 1 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
                            {prediction.champion}: {prediction.currentTier} → {prediction.predictedTier}
                          </Typography>
                          <Typography variant="body2" color="textSecondary" gutterBottom>
                            Factors: {prediction.reasoningFactors.join(', ')}
                          </Typography>
                          <Chip 
                            label={`${(prediction.confidence * 100).toFixed(0)}% Confidence`} 
                            color="info" 
                            size="small" 
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
                    <Typography variant="h6" gutterBottom>
                      Prediction Accuracy
                    </Typography>
                    <Box sx={{ textAlign: 'center', mb: 3 }}>
                      <Typography variant="h3" color="primary">
                        {(metaAnalysis.predictions.predictionAccuracy * 100).toFixed(1)}%
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        Historical Accuracy
                      </Typography>
                    </Box>
                    
                    <Typography variant="subtitle1" gutterBottom>
                      Emerging Champions
                    </Typography>
                    <List>
                      {metaAnalysis.predictions.emergingChampions.map((prediction, index) => (
                        <ListItem key={index}>
                          <ListItemIcon>
                            <TrendingUp color="success" />
                          </ListItemIcon>
                          <ListItemText 
                            primary={`${prediction.champion} (${prediction.role})`}
                            secondary={
                              <Box>
                                <Typography variant="caption">
                                  Timeline: {prediction.timeline}
                                </Typography>
                                <br />
                                <Typography variant="caption">
                                  Confidence: {(prediction.confidence * 100).toFixed(0)}%
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
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={5}>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Meta Recommendations
                    </Typography>
                    <List>
                      {metaAnalysis.recommendations.map((rec, index) => (
                        <ListItem key={index} sx={{ flexDirection: 'column', alignItems: 'flex-start', border: '1px solid', borderColor: 'divider', borderRadius: 1, mb: 2, p: 2 }}>
                          <Typography variant="h6" sx={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}>
                            <Lightbulb color="primary" />
                            {rec.title}
                          </Typography>
                          <Typography variant="body1" gutterBottom sx={{ mt: 1 }}>
                            {rec.description}
                          </Typography>
                          <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                            <Chip 
                              label={rec.priority} 
                              color={rec.priority === 'high' ? 'error' : rec.priority === 'medium' ? 'warning' : 'default'} 
                              size="small" 
                            />
                            <Chip 
                              label={rec.type} 
                              color="primary" 
                              size="small" 
                              variant="outlined"
                            />
                            <Chip 
                              label={rec.targetRank} 
                              color="secondary" 
                              size="small" 
                              variant="outlined"
                            />
                          </Box>
                          {rec.champions && rec.champions.length > 0 && (
                            <Typography variant="body2" color="textSecondary">
                              <strong>Recommended Champions:</strong> {rec.champions.join(', ')}
                            </Typography>
                          )}
                          {rec.strategies && rec.strategies.length > 0 && (
                            <Typography variant="body2" color="textSecondary">
                              <strong>Key Strategies:</strong> {rec.strategies.join(', ')}
                            </Typography>
                          )}
                          {rec.expected && (
                            <Typography variant="body2" color="success.main" sx={{ mt: 1, fontWeight: 'bold' }}>
                              Expected Result: {rec.expected}
                            </Typography>
                          )}
                        </ListItem>
                      ))}
                    </List>
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

export default MetaAnalytics;