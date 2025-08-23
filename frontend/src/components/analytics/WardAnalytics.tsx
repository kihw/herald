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
  Switch,
  FormControlLabel,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material';
import {
  Visibility,
  VisibilityOff,
  LocationOn,
  TrendingUp,
  TrendingDown,
  RemoveRedEye,
  Shield,
  Target,
  Map,
  Timeline,
  Settings,
} from '@mui/icons-material';
import { Line, Radar, Pie, Scatter } from 'react-chartjs-2';
import { wardService } from '../../services/wardService';
import type { WardAnalysis, WardTrendPoint } from '../../types/ward';

interface WardAnalyticsProps {
  playerId: string;
  timeRange?: string;
  champion?: string;
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

const WardAnalytics: React.FC<WardAnalyticsProps> = ({
  playerId,
  timeRange = '30d',
  champion = '',
  position = ''
}) => {
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [wardAnalysis, setWardAnalysis] = useState<WardAnalysis | null>(null);
  const [trendData, setTrendData] = useState<WardTrendPoint[]>([]);
  const [selectedZone, setSelectedZone] = useState<string>('');
  const [selectedMetric, setSelectedMetric] = useState<string>('');
  const [showHeatmap, setShowHeatmap] = useState(true);

  useEffect(() => {
    fetchWardData();
  }, [playerId, timeRange, champion, position]);

  const fetchWardData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch comprehensive ward analysis
      const analysisData = await wardService.getWardAnalysis(
        playerId,
        timeRange,
        champion,
        position
      );
      setWardAnalysis(analysisData);

      // Fetch trend data for all key metrics
      const [
        placedTrends,
        killedTrends,
        controlTrends,
        efficiencyTrends
      ] = await Promise.all([
        wardService.getWardTrends(playerId, 'wards_placed', 'daily', 30),
        wardService.getWardTrends(playerId, 'wards_killed', 'daily', 30),
        wardService.getWardTrends(playerId, 'map_control', 'daily', 30),
        wardService.getWardTrends(playerId, 'efficiency', 'daily', 30)
      ]);
      
      setTrendData(placedTrends); // Default to wards placed trends
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load ward analytics');
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const getMapControlColor = (score: number) => {
    if (score >= 80) return '#4caf50';
    if (score >= 60) return '#ff9800';
    return '#f44336';
  };

  const getEfficiencyRating = (score: number) => {
    if (score >= 85) return { label: 'Excellent', color: '#4caf50' };
    if (score >= 70) return { label: 'Good', color: '#8bc34a' };
    if (score >= 55) return { label: 'Average', color: '#ff9800' };
    if (score >= 40) return { label: 'Below Average', color: '#ff5722' };
    return { label: 'Poor', color: '#f44336' };
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Loading Ward Analytics...
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

  if (!wardAnalysis) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6">
            No ward data available for the selected criteria.
          </Typography>
        </CardContent>
      </Card>
    );
  }

  const efficiencyRating = getEfficiencyRating(wardAnalysis.wardEfficiency);

  // Prepare chart data
  const mapControlRadarData = {
    labels: ['River Control', 'Jungle Control', 'Dragon Area', 'Baron Area', 'Strategic Coverage', 'Safety Provided'],
    datasets: [
      {
        label: 'Your Performance',
        data: [
          wardAnalysis.riverControl.controlScore,
          wardAnalysis.jungleControl.controlScore,
          wardAnalysis.zoneControl?.['Dragon Pit']?.controlScore || 0,
          wardAnalysis.zoneControl?.['Baron Pit']?.controlScore || 0,
          wardAnalysis.strategicCoverage.overallScore,
          wardAnalysis.safetyProvided.overallSafetyScore
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

  const wardTypePieData = {
    labels: ['Yellow Wards', 'Control Wards', 'Blue Trinket'],
    datasets: [
      {
        data: [
          wardAnalysis.yellowWardsAnalysis.totalPlaced,
          wardAnalysis.controlWardsAnalysis.totalPlaced,
          wardAnalysis.blueWardAnalysis.totalUsed
        ],
        backgroundColor: ['#ffeb3b', '#9c27b0', '#2196f3'],
        hoverBackgroundColor: ['#fdd835', '#8e24aa', '#1976d2']
      }
    ]
  };

  const trendChartData = {
    labels: trendData.map(point => new Date(point.date).toLocaleDateString()),
    datasets: [
      {
        label: 'Map Control Score',
        data: trendData.map(point => point.value),
        borderColor: '#1976d2',
        backgroundColor: 'rgba(25, 118, 210, 0.1)',
        tension: 0.4,
        fill: true
      }
    ]
  };

  return (
    <Box>
      <Card>
        <CardContent>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <RemoveRedEye color="primary" />
            Ward Placement & Map Control Analytics
          </Typography>
          
          <Grid container spacing={3} sx={{ mb: 3 }}>
            <Grid item xs={12} md={3}>
              <Card variant="outlined">
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h4" color="primary">
                    {wardAnalysis.mapControlScore.toFixed(1)}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Map Control Score
                  </Typography>
                  <LinearProgress 
                    variant="determinate" 
                    value={wardAnalysis.mapControlScore} 
                    sx={{ mt: 1, height: 6, borderRadius: 3 }}
                    color={wardAnalysis.mapControlScore >= 70 ? 'success' : wardAnalysis.mapControlScore >= 50 ? 'warning' : 'error'}
                  />
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Card variant="outlined">
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h4" sx={{ color: efficiencyRating.color }}>
                    {wardAnalysis.wardEfficiency.toFixed(1)}%
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Ward Efficiency
                  </Typography>
                  <Chip 
                    label={efficiencyRating.label} 
                    size="small" 
                    sx={{ 
                      mt: 1, 
                      bgcolor: efficiencyRating.color, 
                      color: 'white' 
                    }} 
                  />
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Card variant="outlined">
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h4" color="primary">
                    {wardAnalysis.counterWardingScore.toFixed(1)}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Counter-Warding Score
                  </Typography>
                  <Typography variant="caption" sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', mt: 0.5, gap: 0.5 }}>
                    <Shield fontSize="small" />
                    Vision Denial
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} md={3}>
              <Card variant="outlined">
                <CardContent sx={{ textAlign: 'center' }}>
                  <Typography variant="h4" color="primary">
                    {wardAnalysis.territoryControlled.toFixed(1)}%
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Territory Controlled
                  </Typography>
                  <Typography variant="caption" sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', mt: 0.5, gap: 0.5 }}>
                    <Map fontSize="small" />
                    Map Coverage
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>

          <Tabs value={tabValue} onChange={handleTabChange} variant="scrollable" scrollButtons="auto">
            <Tab label="Overview" icon={<Visibility />} />
            <Tab label="Map Control" icon={<Map />} />
            <Tab label="Ward Types" icon={<Target />} />
            <Tab label="Placement Patterns" icon={<LocationOn />} />
            <Tab label="Trends" icon={<Timeline />} />
            <Tab label="Optimization" icon={<Settings />} />
          </Tabs>

          <TabPanel value={tabValue} index={0}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={8}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Map Control Performance
                    </Typography>
                    <Box sx={{ height: 400 }}>
                      <Radar 
                        data={mapControlRadarData} 
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
                    <Typography variant="h6" gutterBottom>
                      Key Strengths
                    </Typography>
                    <List dense>
                      {wardAnalysis.recommendations.map((rec, index) => (
                        rec.type === 'strength' && (
                          <ListItem key={index}>
                            <ListItemIcon>
                              <TrendingUp color="success" />
                            </ListItemIcon>
                            <ListItemText 
                              primary={rec.title}
                              secondary={rec.description}
                            />
                          </ListItem>
                        )
                      ))}
                    </List>
                  </CardContent>
                </Card>
                
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Improvement Areas
                    </Typography>
                    <List dense>
                      {wardAnalysis.recommendations.map((rec, index) => (
                        rec.type === 'improvement' && (
                          <ListItem key={index}>
                            <ListItemIcon>
                              <TrendingDown color="warning" />
                            </ListItemIcon>
                            <ListItemText 
                              primary={rec.title}
                              secondary={rec.description}
                            />
                          </ListItem>
                        )
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
                      Zone Control Analysis
                    </Typography>
                    <FormControl fullWidth sx={{ mb: 2 }}>
                      <InputLabel>Focus Zone</InputLabel>
                      <Select
                        value={selectedZone}
                        onChange={(e) => setSelectedZone(e.target.value)}
                      >
                        <MenuItem value="">All Zones</MenuItem>
                        <MenuItem value="river">River Control</MenuItem>
                        <MenuItem value="jungle">Jungle Control</MenuItem>
                        <MenuItem value="dragon">Dragon Pit</MenuItem>
                        <MenuItem value="baron">Baron Pit</MenuItem>
                      </Select>
                    </FormControl>
                    
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" gutterBottom>
                        River Control: {wardAnalysis.riverControl.controlScore.toFixed(1)}%
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={wardAnalysis.riverControl.controlScore} 
                        sx={{ mb: 1, height: 8, borderRadius: 4 }}
                        color="primary"
                      />
                    </Box>
                    
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" gutterBottom>
                        Jungle Control: {wardAnalysis.jungleControl.controlScore.toFixed(1)}%
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={wardAnalysis.jungleControl.controlScore} 
                        sx={{ mb: 1, height: 8, borderRadius: 4 }}
                        color="success"
                      />
                    </Box>
                    
                    {Object.entries(wardAnalysis.zoneControl).map(([zone, data]) => (
                      <Box key={zone} sx={{ mb: 2 }}>
                        <Typography variant="body2" gutterBottom>
                          {zone}: {data.controlScore.toFixed(1)}%
                        </Typography>
                        <LinearProgress 
                          variant="determinate" 
                          value={data.controlScore} 
                          sx={{ mb: 1, height: 8, borderRadius: 4 }}
                          color="warning"
                        />
                      </Box>
                    ))}
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Strategic Coverage
                    </Typography>
                    
                    <Box sx={{ mb: 3 }}>
                      <Typography variant="body2" gutterBottom>
                        Overall Coverage Score
                      </Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <LinearProgress 
                          variant="determinate" 
                          value={wardAnalysis.strategicCoverage.overallScore} 
                          sx={{ flexGrow: 1, height: 10, borderRadius: 5 }}
                        />
                        <Typography variant="h6" color="primary">
                          {wardAnalysis.strategicCoverage.overallScore.toFixed(1)}%
                        </Typography>
                      </Box>
                    </Box>
                    
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Box sx={{ textAlign: 'center', p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
                          <Typography variant="h4" color="primary">
                            {wardAnalysis.strategicCoverage.dragonPitCoverage.toFixed(1)}%
                          </Typography>
                          <Typography variant="body2">
                            Dragon Coverage
                          </Typography>
                        </Box>
                      </Grid>
                      <Grid item xs={6}>
                        <Box sx={{ textAlign: 'center', p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
                          <Typography variant="h4" color="primary">
                            {wardAnalysis.strategicCoverage.baronPitCoverage.toFixed(1)}%
                          </Typography>
                          <Typography variant="body2">
                            Baron Coverage
                          </Typography>
                        </Box>
                      </Grid>
                    </Grid>
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
                      Ward Type Distribution
                    </Typography>
                    <Box sx={{ height: 300 }}>
                      <Pie 
                        data={wardTypePieData} 
                        options={{
                          responsive: true,
                          maintainAspectRatio: false,
                          plugins: {
                            legend: { position: 'bottom' }
                          }
                        }} 
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Card variant="outlined" sx={{ mb: 2 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ color: '#ffeb3b' }}>
                      Yellow Wards Analysis
                    </Typography>
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Typography variant="h4" color="primary">
                          {wardAnalysis.yellowWardsAnalysis.totalPlaced}
                        </Typography>
                        <Typography variant="body2">Total Placed</Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="h4" color="primary">
                          {wardAnalysis.yellowWardsAnalysis.averageLifespan.toFixed(1)}s
                        </Typography>
                        <Typography variant="body2">Avg Lifespan</Typography>
                      </Grid>
                    </Grid>
                  </CardContent>
                </Card>
                
                <Card variant="outlined" sx={{ mb: 2 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ color: '#9c27b0' }}>
                      Control Wards Analysis
                    </Typography>
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Typography variant="h4" color="primary">
                          {wardAnalysis.controlWardsAnalysis.totalPlaced}
                        </Typography>
                        <Typography variant="body2">Total Placed</Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="h4" color="primary">
                          {wardAnalysis.controlWardsAnalysis.averageValue.toFixed(0)}g
                        </Typography>
                        <Typography variant="body2">Avg Value</Typography>
                      </Grid>
                    </Grid>
                  </CardContent>
                </Card>
                
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom sx={{ color: '#2196f3' }}>
                      Blue Trinket Analysis
                    </Typography>
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <Typography variant="h4" color="primary">
                          {wardAnalysis.blueWardAnalysis.totalUsed}
                        </Typography>
                        <Typography variant="body2">Total Used</Typography>
                      </Grid>
                      <Grid item xs={6}>
                        <Typography variant="h4" color="primary">
                          {wardAnalysis.blueWardAnalysis.accuracyRate.toFixed(1)}%
                        </Typography>
                        <Typography variant="body2">Accuracy Rate</Typography>
                      </Grid>
                    </Grid>
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
                      Ward Placement Patterns
                    </Typography>
                    
                    <FormControlLabel
                      control={
                        <Switch
                          checked={showHeatmap}
                          onChange={(e) => setShowHeatmap(e.target.checked)}
                        />
                      }
                      label="Show Heatmap"
                    />
                    
                    <Grid container spacing={3} sx={{ mt: 1 }}>
                      <Grid item xs={12} md={4}>
                        <Typography variant="subtitle1" gutterBottom>
                          Optimal Placements
                        </Typography>
                        <List dense>
                          {wardAnalysis.optimalPlacements.map((placement, index) => (
                            <ListItem key={index}>
                              <ListItemIcon>
                                <LocationOn color="success" />
                              </ListItemIcon>
                              <ListItemText 
                                primary={placement.location}
                                secondary={`${placement.frequency.toFixed(1)}% frequency`}
                              />
                            </ListItem>
                          ))}
                        </List>
                      </Grid>
                      
                      <Grid item xs={12} md={4}>
                        <Typography variant="subtitle1" gutterBottom>
                          Placement Timing
                        </Typography>
                        <Box sx={{ mb: 2 }}>
                          <Typography variant="body2" gutterBottom>
                            Early Game: {wardAnalysis.placementTiming.earlyGameFrequency.toFixed(1)}/min
                          </Typography>
                          <LinearProgress 
                            variant="determinate" 
                            value={wardAnalysis.placementTiming.earlyGameFrequency * 10} 
                            sx={{ mb: 1 }}
                          />
                        </Box>
                        <Box sx={{ mb: 2 }}>
                          <Typography variant="body2" gutterBottom>
                            Mid Game: {wardAnalysis.placementTiming.midGameFrequency.toFixed(1)}/min
                          </Typography>
                          <LinearProgress 
                            variant="determinate" 
                            value={wardAnalysis.placementTiming.midGameFrequency * 10} 
                            sx={{ mb: 1 }}
                          />
                        </Box>
                        <Box>
                          <Typography variant="body2" gutterBottom>
                            Late Game: {wardAnalysis.placementTiming.lateGameFrequency.toFixed(1)}/min
                          </Typography>
                          <LinearProgress 
                            variant="determinate" 
                            value={wardAnalysis.placementTiming.lateGameFrequency * 10} 
                          />
                        </Box>
                      </Grid>
                      
                      <Grid item xs={12} md={4}>
                        <Typography variant="subtitle1" gutterBottom>
                          Placement Quality
                        </Typography>
                        <Box sx={{ textAlign: 'center', p: 3, bgcolor: 'background.default', borderRadius: 2 }}>
                          <Typography variant="h3" color="primary">
                            {wardAnalysis.placementOptimization.currentQualityScore.toFixed(1)}%
                          </Typography>
                          <Typography variant="body2" gutterBottom>
                            Placement Quality Score
                          </Typography>
                          <Chip 
                            label={wardAnalysis.placementOptimization.currentQualityScore >= 70 ? 'Good' : 'Needs Improvement'} 
                            color={wardAnalysis.placementOptimization.currentQualityScore >= 70 ? 'success' : 'warning'}
                          />
                        </Box>
                      </Grid>
                    </Grid>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={4}>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Ward Performance Trends
                    </Typography>
                    
                    <FormControl sx={{ mb: 3, minWidth: 200 }}>
                      <InputLabel>Trend Metric</InputLabel>
                      <Select
                        value={selectedMetric}
                        onChange={(e) => setSelectedMetric(e.target.value)}
                      >
                        <MenuItem value="wards_placed">Wards Placed</MenuItem>
                        <MenuItem value="wards_killed">Wards Killed</MenuItem>
                        <MenuItem value="map_control">Map Control</MenuItem>
                        <MenuItem value="efficiency">Ward Efficiency</MenuItem>
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
                              beginAtZero: true,
                              max: 100,
                              title: { display: true, text: 'Score' }
                            },
                            x: {
                              title: { display: true, text: 'Date' }
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
            </Grid>
          </TabPanel>

          <TabPanel value={tabValue} index={5}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Placement Optimization
                    </Typography>
                    <Box sx={{ mb: 3 }}>
                      <Typography variant="body2" gutterBottom>
                        Expected Control Gain: +{wardAnalysis.placementOptimization.expectedControlGain.toFixed(1)}%
                      </Typography>
                      <Typography variant="body2" gutterBottom>
                        Expected Safety Gain: +{wardAnalysis.placementOptimization.expectedSafetyGain.toFixed(1)}%
                      </Typography>
                    </Box>
                    
                    <List>
                      {wardAnalysis.placementOptimization.suggestions.map((suggestion, index) => (
                        <ListItem key={index}>
                          <ListItemIcon>
                            <Target color="primary" />
                          </ListItemIcon>
                          <ListItemText 
                            primary={suggestion.area}
                            secondary={suggestion.reason}
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
                      Clearing Optimization
                    </Typography>
                    <Box sx={{ mb: 3 }}>
                      <Typography variant="body2" gutterBottom>
                        Expected Denial Gain: +{wardAnalysis.clearingOptimization.expectedDenialGain.toFixed(1)}%
                      </Typography>
                      <Typography variant="body2" gutterBottom>
                        Counter-Warding Efficiency: {wardAnalysis.clearingOptimization.currentCounterWardingEfficiency.toFixed(1)}%
                      </Typography>
                    </Box>
                    
                    <List>
                      {wardAnalysis.clearingOptimization.suggestions.map((suggestion, index) => (
                        <ListItem key={index}>
                          <ListItemIcon>
                            <VisibilityOff color="secondary" />
                          </ListItemIcon>
                          <ListItemText 
                            primary={suggestion.area}
                            secondary={suggestion.reason}
                          />
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

export default WardAnalytics;