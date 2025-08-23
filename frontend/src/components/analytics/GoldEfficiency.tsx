import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Paper,
  Chip,
  LinearProgress,
  Tooltip,
  IconButton,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Alert,
  CircularProgress,
  ToggleButton,
  ToggleButtonGroup,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider
} from '@mui/material';
import {
  AccountBalance as GoldIcon,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  ShoppingCart as ShoppingIcon,
  Agriculture as FarmingIcon,
  Timeline as TrendIcon,
  CompareArrows as CompareIcon,
  Lightbulb as OptimizationIcon,
  Download as DownloadIcon,
  Refresh as RefreshIcon,
  MonetizationOn as CoinIcon,
  AccessTime as TimingIcon,
  Speed as EfficiencyIcon,
  CheckCircle as CheckIcon,
  Warning as WarningIcon,
  Info as InfoIcon
} from '@mui/icons-material';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip, Legend, ResponsiveContainer, BarChart, Bar, PieChart, Pie, Cell, AreaChart, Area } from 'recharts';

interface GoldAnalysis {
  player_id: string;
  champion?: string;
  position?: string;
  time_range: string;
  average_gold_earned: number;
  average_gold_per_minute: number;
  gold_efficiency_score: number;
  economy_rating: string;
  gold_sources: GoldSourcesData;
  item_efficiency: ItemEfficiencyData;
  spending_patterns: SpendingPatternsData;
  early_game_gold: GoldPhaseData;
  mid_game_gold: GoldPhaseData;
  late_game_gold: GoldPhaseData;
  role_benchmark: GoldBenchmark;
  rank_benchmark: GoldBenchmark;
  global_benchmark: GoldBenchmark;
  gold_advantage_win_rate: number;
  gold_disadvantage_win_rate: number;
  gold_impact_score: number;
  trend_direction: string;
  trend_slope: number;
  trend_confidence: number;
  trend_data: GoldTrendPoint[];
  income_optimization: IncomeOptimizationData;
  spending_optimization: SpendingOptimizationData;
  strength_areas: string[];
  improvement_areas: string[];
  recommendations: GoldRecommendation[];
  recent_matches: MatchGoldData[];
}

interface GoldSourcesData {
  farming_gold: number;
  farming_percent: number;
  kills_gold: number;
  kills_percent: number;
  assists_gold: number;
  assists_percent: number;
  objective_gold: number;
  objective_percent: number;
  passive_gold: number;
  passive_percent: number;
  items_gold: number;
  items_percent: number;
  cs_gold_per_minute: number;
  kill_gold_efficiency: number;
  objective_gold_share: number;
}

interface ItemEfficiencyData {
  average_items_completed: number;
  item_completion_speed: number;
  gold_spent_on_items: number;
  item_value_efficiency: number;
  damage_items_percent: number;
  defensive_items_percent: number;
  utility_items_percent: number;
  first_item_timing: number;
  core_items_timing: number;
  six_items_timing: number;
  optimal_item_order: boolean;
  counter_build_efficiency: number;
  component_utilization: number;
}

interface SpendingPatternsData {
  control_wards_percent: number;
  consumables_percent: number;
  back_timing: BackTimingData;
  gold_efficiency_by_phase: PhaseGoldEfficiency[];
  average_shopping_time: number;
  optimal_back_timing: number;
  emergency_backs: number;
}

interface BackTimingData {
  average_back_timing: number;
  optimal_backs: number;
  suboptimal_backs: number;
  forced_backs: number;
  gold_per_back: number;
}

interface PhaseGoldEfficiency {
  phase: string;
  gold_per_minute: number;
  spending_efficiency: number;
  income_efficiency: number;
  economy_score: number;
}

interface GoldPhaseData {
  phase: string;
  average_gold_per_minute: number;
  gold_advantage: number;
  farming_efficiency: number;
  kill_participation: number;
  objective_participation: number;
  spending_score: number;
  efficiency_rating: string;
}

interface GoldBenchmark {
  category: string;
  average_gold_per_minute: number;
  top_10_percent: number;
  top_25_percent: number;
  median: number;
  player_percentile: number;
  efficiency_average: number;
}

interface IncomeOptimizationData {
  cs_improvement_potential: number;
  kp_improvement_potential: number;
  objective_improvement_potential: number;
  early_farming_suggestions: string[];
  mid_game_position_suggestions: string[];
  late_game_focus_suggestions: string[];
  expected_gpm_increase: number;
  expected_win_rate_increase: number;
}

interface SpendingOptimizationData {
  item_order_optimization: string[];
  back_timing_optimization: string[];
  gold_allocation_suggestions: string[];
  component_buying_tips: string[];
  power_spike_timing: string[];
  early_game_priorities: string[];
  mid_game_priorities: string[];
  late_game_priorities: string[];
}

interface GoldTrendPoint {
  date: string;
  gold_per_minute: number;
  gold_efficiency: number;
  farming_efficiency: number;
  spending_efficiency: number;
  moving_average: number;
}

interface GoldRecommendation {
  priority: 'high' | 'medium' | 'low';
  category: string;
  title: string;
  description: string;
  impact: string;
  game_phase: string[];
  expected_gpm_increase: number;
  implementation_difficulty: string;
}

interface MatchGoldData {
  match_id: string;
  champion: string;
  position: string;
  total_gold_earned: number;
  gold_per_minute: number;
  gold_efficiency_score: number;
  farming_gold: number;
  kill_gold: number;
  objective_gold: number;
  items_completed: number;
  control_wards_spent: number;
  game_duration: number;
  result: string;
  date: string;
  gold_advantage_at_15: number;
  team_gold_share: number;
}

interface GoldEfficiencyProps {
  playerData: GoldAnalysis | null;
  timeRange: '7d' | '30d' | '90d';
  onTimeRangeChange: (timeRange: string) => void;
  loading?: boolean;
  error?: string;
}

const GOLD_SOURCE_COLORS = {
  farming: '#4caf50',
  kills: '#f44336',
  assists: '#ff9800',
  objectives: '#9c27b0',
  passive: '#607d8b',
  items: '#795548'
};

const PHASE_COLORS = {
  early: '#4cc9f0',
  mid: '#f9c74f',
  late: '#f8961e'
};

export const GoldEfficiency: React.FC<GoldEfficiencyProps> = ({
  playerData,
  timeRange,
  onTimeRangeChange,
  loading = false,
  error
}) => {
  const [selectedTab, setSelectedTab] = useState(0);
  const [benchmarkType, setBenchmarkType] = useState<'role' | 'rank' | 'global'>('role');

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };

  const getGoldSourceData = () => {
    if (!playerData?.gold_sources) return [];
    
    return [
      { name: 'Farming', value: playerData.gold_sources.farming_percent, color: GOLD_SOURCE_COLORS.farming, amount: playerData.gold_sources.farming_gold },
      { name: 'Kills', value: playerData.gold_sources.kills_percent, color: GOLD_SOURCE_COLORS.kills, amount: playerData.gold_sources.kills_gold },
      { name: 'Assists', value: playerData.gold_sources.assists_percent, color: GOLD_SOURCE_COLORS.assists, amount: playerData.gold_sources.assists_gold },
      { name: 'Objectives', value: playerData.gold_sources.objective_percent, color: GOLD_SOURCE_COLORS.objectives, amount: playerData.gold_sources.objective_gold },
      { name: 'Passive', value: playerData.gold_sources.passive_percent, color: GOLD_SOURCE_COLORS.passive, amount: playerData.gold_sources.passive_gold }
    ];
  };

  const getPhaseComparisonData = () => {
    if (!playerData) return [];
    
    return [
      {
        phase: 'Early Game',
        gpm: playerData.early_game_gold.average_gold_per_minute,
        efficiency: playerData.early_game_gold.spending_score,
        advantage: playerData.early_game_gold.gold_advantage
      },
      {
        phase: 'Mid Game',
        gpm: playerData.mid_game_gold.average_gold_per_minute,
        efficiency: playerData.mid_game_gold.spending_score,
        advantage: playerData.mid_game_gold.gold_advantage
      },
      {
        phase: 'Late Game',
        gpm: playerData.late_game_gold.average_gold_per_minute,
        efficiency: playerData.late_game_gold.spending_score,
        advantage: playerData.late_game_gold.gold_advantage
      }
    ];
  };

  const getEfficiencyColor = (rating: string) => {
    switch (rating) {
      case 'excellent': return 'success';
      case 'good': return 'primary';
      case 'average': return 'warning';
      case 'poor': return 'error';
      default: return 'default';
    }
  };

  const getPriorityIcon = (priority: string) => {
    switch (priority) {
      case 'high': return <WarningIcon color="error" />;
      case 'medium': return <InfoIcon color="warning" />;
      case 'low': return <InfoIcon color="info" />;
      default: return <InfoIcon />;
    }
  };

  const formatGold = (gold: number) => {
    if (gold >= 1000) {
      return `${(gold / 1000).toFixed(1)}K`;
    }
    return Math.round(gold).toString();
  };

  const getBenchmarkValue = (metric: string) => {
    const benchmark = benchmarkType === 'role' ? playerData?.role_benchmark : 
                     benchmarkType === 'rank' ? playerData?.rank_benchmark : 
                     playerData?.global_benchmark;
    
    switch (metric) {
      case 'gpm':
        return benchmark?.average_gold_per_minute || 400;
      default:
        return 400;
    }
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
            <CircularProgress />
            <Typography variant="body2" sx={{ ml: 2 }}>
              Analyzing gold efficiency...
            </Typography>
          </Box>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  if (!playerData) {
    return (
      <Alert severity="info" sx={{ mb: 2 }}>
        No gold efficiency data available for the selected time range.
      </Alert>
    );
  }

  return (
    <Box>
      {/* Header Controls */}
      <Card sx={{ mb: 2 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} sm={6} md={4}>
              <ToggleButtonGroup
                value={timeRange}
                exclusive
                onChange={(_, value) => value && onTimeRangeChange(value)}
                size="small"
              >
                <ToggleButton value="7d">7 Days</ToggleButton>
                <ToggleButton value="30d">30 Days</ToggleButton>
                <ToggleButton value="90d">90 Days</ToggleButton>
              </ToggleButtonGroup>
            </Grid>

            <Grid item xs={12} sm={6} md={4}>
              <FormControl fullWidth size="small">
                <InputLabel>Benchmark</InputLabel>
                <Select
                  value={benchmarkType}
                  label="Benchmark"
                  onChange={(e) => setBenchmarkType(e.target.value as any)}
                >
                  <MenuItem value="role">Role Average</MenuItem>
                  <MenuItem value="rank">Rank Average</MenuItem>
                  <MenuItem value="global">Global Average</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} sm={12} md={4}>
              <Box display="flex" gap={1} justifyContent="flex-end">
                <Tooltip title="Refresh Data">
                  <IconButton size="small">
                    <RefreshIcon />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Export Report">
                  <IconButton size="small">
                    <DownloadIcon />
                  </IconButton>
                </Tooltip>
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Key Metrics Overview */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <CoinIcon sx={{ fontSize: 40, color: 'primary.main', mb: 1 }} />
            <Typography variant="h4" color="primary">
              {Math.round(playerData.average_gold_per_minute)}
            </Typography>
            <Typography variant="caption" color="textSecondary">
              Gold per Minute
            </Typography>
            <Box display="flex" alignItems="center" justifyContent="center" mt={1}>
              {playerData.average_gold_per_minute > getBenchmarkValue('gpm') ? (
                <TrendingUpIcon color="success" fontSize="small" />
              ) : (
                <TrendingDownIcon color="error" fontSize="small" />
              )}
              <Typography variant="caption" sx={{ ml: 0.5 }}>
                vs {benchmarkType} avg
              </Typography>
            </Box>
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <EfficiencyIcon sx={{ fontSize: 40, color: 'warning.main', mb: 1 }} />
            <Typography variant="h4" color="warning.main">
              {playerData.gold_efficiency_score.toFixed(0)}
            </Typography>
            <Typography variant="caption" color="textSecondary">
              Efficiency Score
            </Typography>
            <Chip 
              label={playerData.economy_rating} 
              size="small" 
              color={getEfficiencyColor(playerData.economy_rating) as any}
              sx={{ mt: 1 }}
            />
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <GoldIcon sx={{ fontSize: 40, color: 'secondary.main', mb: 1 }} />
            <Typography variant="h4" color="secondary.main">
              {playerData.gold_impact_score.toFixed(0)}
            </Typography>
            <Typography variant="caption" color="textSecondary">
              Gold Impact Score
            </Typography>
            <LinearProgress 
              variant="determinate" 
              value={playerData.gold_impact_score} 
              sx={{ mt: 1, height: 8, borderRadius: 4 }}
            />
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <FarmingIcon sx={{ fontSize: 40, color: 'success.main', mb: 1 }} />
            <Typography variant="h4" color="success.main">
              {playerData.gold_sources.farming_percent.toFixed(0)}%
            </Typography>
            <Typography variant="caption" color="textSecondary">
              From Farming
            </Typography>
            <Typography variant="body2" sx={{ mt: 1, fontSize: '0.75rem' }}>
              {formatGold(playerData.gold_sources.farming_gold)} avg
            </Typography>
          </Paper>
        </Grid>
      </Grid>

      {/* Tabbed Analytics */}
      <Card>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={selectedTab} onChange={handleTabChange}>
            <Tab label="Overview" />
            <Tab label="Gold Sources" />
            <Tab label="Spending" />
            <Tab label="Game Phases" />
            <Tab label="Trends" />
            <Tab label="Optimization" />
          </Tabs>
        </Box>

        <CardContent>
          {selectedTab === 0 && (
            <Grid container spacing={3}>
              {/* Gold Performance vs Benchmark */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Performance vs {benchmarkType.charAt(0).toUpperCase() + benchmarkType.slice(1)} Average
                </Typography>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="body2" color="textSecondary">
                    Gold per Minute
                  </Typography>
                  <LinearProgress
                    variant="determinate"
                    value={Math.min((playerData.average_gold_per_minute / getBenchmarkValue('gpm')) * 100, 150)}
                    sx={{ height: 8, borderRadius: 4, mb: 1 }}
                  />
                  <Typography variant="caption">
                    {Math.round(playerData.average_gold_per_minute)} / {Math.round(getBenchmarkValue('gpm'))} GPM
                  </Typography>
                </Box>
                
                <Box sx={{ mb: 2 }}>
                  <Typography variant="body2" color="textSecondary">
                    Efficiency Score
                  </Typography>
                  <LinearProgress
                    variant="determinate"
                    value={playerData.gold_efficiency_score}
                    color={getEfficiencyColor(playerData.economy_rating) as any}
                    sx={{ height: 8, borderRadius: 4, mb: 1 }}
                  />
                  <Typography variant="caption">
                    {playerData.gold_efficiency_score.toFixed(0)}/100 ({playerData.economy_rating})
                  </Typography>
                </Box>
              </Grid>

              {/* Win Rate Impact */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Gold Impact on Win Rate
                </Typography>
                <Grid container spacing={2}>
                  <Grid item xs={6}>
                    <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'success.light' }}>
                      <Typography variant="h5" color="success.contrastText">
                        {playerData.gold_advantage_win_rate.toFixed(0)}%
                      </Typography>
                      <Typography variant="caption" color="success.contrastText">
                        Gold Advantage Games
                      </Typography>
                    </Paper>
                  </Grid>
                  <Grid item xs={6}>
                    <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'error.light' }}>
                      <Typography variant="h5" color="error.contrastText">
                        {playerData.gold_disadvantage_win_rate.toFixed(0)}%
                      </Typography>
                      <Typography variant="caption" color="error.contrastText">
                        Gold Disadvantage Games
                      </Typography>
                    </Paper>
                  </Grid>
                </Grid>
                
                <Box sx={{ mt: 2 }}>
                  <Typography variant="body2" color="textSecondary">
                    Gold Impact: +{(playerData.gold_advantage_win_rate - playerData.gold_disadvantage_win_rate).toFixed(0)}% win rate with gold advantage
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          )}

          {selectedTab === 1 && (
            <Grid container spacing={3}>
              {/* Gold Sources Pie Chart */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Gold Income Sources
                </Typography>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={getGoldSourceData()}
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      fill="#8884d8"
                      dataKey="value"
                      label={({ name, value }) => `${name}: ${value.toFixed(1)}%`}
                    >
                      {getGoldSourceData().map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <RechartsTooltip 
                      formatter={(value: any, name: any, props: any) => [
                        `${value.toFixed(1)}% (${formatGold(props.payload.amount)}g)`,
                        name
                      ]}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </Grid>

              {/* Gold Sources Breakdown */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Income Analysis
                </Typography>
                <List>
                  <ListItem>
                    <ListItemIcon>
                      <FarmingIcon sx={{ color: GOLD_SOURCE_COLORS.farming }} />
                    </ListItemIcon>
                    <ListItemText
                      primary={`Farming: ${playerData.gold_sources.farming_percent.toFixed(1)}%`}
                      secondary={`${formatGold(playerData.gold_sources.farming_gold)} gold • ${playerData.gold_sources.cs_gold_per_minute.toFixed(0)} GPM from CS`}
                    />
                  </ListItem>
                  <ListItem>
                    <ListItemIcon>
                      <CoinIcon sx={{ color: GOLD_SOURCE_COLORS.kills }} />
                    </ListItemIcon>
                    <ListItemText
                      primary={`Combat: ${(playerData.gold_sources.kills_percent + playerData.gold_sources.assists_percent).toFixed(1)}%`}
                      secondary={`${formatGold(playerData.gold_sources.kills_gold + playerData.gold_sources.assists_gold)} gold • ${playerData.gold_sources.kill_gold_efficiency.toFixed(0)}% efficiency`}
                    />
                  </ListItem>
                  <ListItem>
                    <ListItemIcon>
                      <GoldIcon sx={{ color: GOLD_SOURCE_COLORS.objectives }} />
                    </ListItemIcon>
                    <ListItemText
                      primary={`Objectives: ${playerData.gold_sources.objective_percent.toFixed(1)}%`}
                      secondary={`${formatGold(playerData.gold_sources.objective_gold)} gold • ${playerData.gold_sources.objective_gold_share.toFixed(0)}% team share`}
                    />
                  </ListItem>
                </List>
              </Grid>
            </Grid>
          )}

          {selectedTab === 2 && (
            <Grid container spacing={3}>
              {/* Item Efficiency */}
              <Grid item xs={12} md={6}>
                <Typography variant="h6" gutterBottom>
                  Item Efficiency
                </Typography>
                <Paper sx={{ p: 2 }}>
                  <Grid container spacing={2}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="textSecondary">
                        Items Completed
                      </Typography>
                      <Typography variant="h6">
                        {playerData.item_efficiency.average_items_completed}
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="textSecondary">
                        First Item
                      </Typography>
                      <Typography variant="h6">
                        {playerData.item_efficiency.first_item_timing.toFixed(1)}m
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="textSecondary">
                        Core Items
                      </Typography>
                      <Typography variant="h6">
                        {playerData.item_efficiency.core_items_timing.toFixed(1)}m
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="textSecondary">
                        Value Efficiency
                      </Typography>
                      <Typography variant="h6">
                        {playerData.item_efficiency.item_value_efficiency.toFixed(0)}%
                      </Typography>
                    </Grid>
                  </Grid>
                </Paper>
              </Grid>

              {/* Spending Patterns */}
              <Grid item xs={12} md={6}>
                <Typography variant="h6" gutterBottom>
                  Spending Patterns
                </Typography>
                <Paper sx={{ p: 2 }}>
                  <Grid container spacing={2}>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="textSecondary">
                        Control Wards
                      </Typography>
                      <Typography variant="h6">
                        {playerData.spending_patterns.control_wards_percent.toFixed(1)}%
                      </Typography>
                    </Grid>
                    <Grid item xs={6}>
                      <Typography variant="body2" color="textSecondary">
                        Back Timing
                      </Typography>
                      <Typography variant="h6">
                        {playerData.spending_patterns.optimal_back_timing.toFixed(0)}%
                      </Typography>
                    </Grid>
                    <Grid item xs={12}>
                      <Divider sx={{ my: 1 }} />
                      <Typography variant="body2" color="textSecondary">
                        Optimal Backs: {playerData.spending_patterns.back_timing.optimal_backs}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        Emergency Backs: {playerData.spending_patterns.emergency_backs}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        Gold per Back: {formatGold(playerData.spending_patterns.back_timing.gold_per_back)}
                      </Typography>
                    </Grid>
                  </Grid>
                </Paper>
              </Grid>
            </Grid>
          )}

          {selectedTab === 3 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Gold Efficiency by Game Phase
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={getPhaseComparisonData()}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="phase" />
                  <YAxis />
                  <RechartsTooltip />
                  <Legend />
                  <Bar dataKey="gpm" fill={PHASE_COLORS.early} name="Gold per Minute" />
                  <Bar dataKey="efficiency" fill={PHASE_COLORS.mid} name="Efficiency Score" />
                </BarChart>
              </ResponsiveContainer>

              <TableContainer component={Paper} sx={{ mt: 2 }}>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Phase</TableCell>
                      <TableCell align="right">GPM</TableCell>
                      <TableCell align="right">Gold Advantage</TableCell>
                      <TableCell align="right">Farming %</TableCell>
                      <TableCell align="center">Rating</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        <Chip label="Early Game" size="small" sx={{ bgcolor: PHASE_COLORS.early, color: 'white' }} />
                      </TableCell>
                      <TableCell align="right">{Math.round(playerData.early_game_gold.average_gold_per_minute)}</TableCell>
                      <TableCell align="right">
                        <Typography color={playerData.early_game_gold.gold_advantage > 0 ? 'success.main' : 'error.main'}>
                          {playerData.early_game_gold.gold_advantage > 0 ? '+' : ''}{Math.round(playerData.early_game_gold.gold_advantage)}g
                        </Typography>
                      </TableCell>
                      <TableCell align="right">{playerData.early_game_gold.farming_efficiency.toFixed(0)}%</TableCell>
                      <TableCell align="center">
                        <Chip 
                          label={playerData.early_game_gold.efficiency_rating}
                          size="small"
                          color={getEfficiencyColor(playerData.early_game_gold.efficiency_rating) as any}
                        />
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        <Chip label="Mid Game" size="small" sx={{ bgcolor: PHASE_COLORS.mid, color: 'white' }} />
                      </TableCell>
                      <TableCell align="right">{Math.round(playerData.mid_game_gold.average_gold_per_minute)}</TableCell>
                      <TableCell align="right">
                        <Typography color={playerData.mid_game_gold.gold_advantage > 0 ? 'success.main' : 'error.main'}>
                          {playerData.mid_game_gold.gold_advantage > 0 ? '+' : ''}{Math.round(playerData.mid_game_gold.gold_advantage)}g
                        </Typography>
                      </TableCell>
                      <TableCell align="right">{playerData.mid_game_gold.farming_efficiency.toFixed(0)}%</TableCell>
                      <TableCell align="center">
                        <Chip 
                          label={playerData.mid_game_gold.efficiency_rating}
                          size="small"
                          color={getEfficiencyColor(playerData.mid_game_gold.efficiency_rating) as any}
                        />
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        <Chip label="Late Game" size="small" sx={{ bgcolor: PHASE_COLORS.late, color: 'white' }} />
                      </TableCell>
                      <TableCell align="right">{Math.round(playerData.late_game_gold.average_gold_per_minute)}</TableCell>
                      <TableCell align="right">
                        <Typography color={playerData.late_game_gold.gold_advantage > 0 ? 'success.main' : 'error.main'}>
                          {playerData.late_game_gold.gold_advantage > 0 ? '+' : ''}{Math.round(playerData.late_game_gold.gold_advantage)}g
                        </Typography>
                      </TableCell>
                      <TableCell align="right">{playerData.late_game_gold.farming_efficiency.toFixed(0)}%</TableCell>
                      <TableCell align="center">
                        <Chip 
                          label={playerData.late_game_gold.efficiency_rating}
                          size="small"
                          color={getEfficiencyColor(playerData.late_game_gold.efficiency_rating) as any}
                        />
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          )}

          {selectedTab === 4 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Gold Efficiency Trends
              </Typography>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={playerData.trend_data}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis 
                    dataKey="date" 
                    tickFormatter={(value) => new Date(value).toLocaleDateString()}
                  />
                  <YAxis />
                  <RechartsTooltip 
                    labelFormatter={(value) => new Date(value).toLocaleDateString()}
                  />
                  <Legend />
                  <Line 
                    type="monotone" 
                    dataKey="gold_per_minute" 
                    stroke="#8884d8" 
                    name="GPM" 
                  />
                  <Line 
                    type="monotone" 
                    dataKey="gold_efficiency" 
                    stroke="#82ca9d" 
                    name="Efficiency"
                  />
                  <Line 
                    type="monotone" 
                    dataKey="moving_average" 
                    stroke="#ffc658" 
                    strokeDasharray="5 5"
                    name="Moving Average"
                  />
                </LineChart>
              </ResponsiveContainer>

              <Box sx={{ mt: 2 }}>
                <Grid container spacing={2}>
                  <Grid item xs={12} sm={4}>
                    <Paper sx={{ p: 2, textAlign: 'center' }}>
                      <Typography variant="h6" color={
                        playerData.trend_direction === 'improving' ? 'success.main' :
                        playerData.trend_direction === 'declining' ? 'error.main' : 'warning.main'
                      }>
                        {playerData.trend_direction.charAt(0).toUpperCase() + playerData.trend_direction.slice(1)}
                      </Typography>
                      <Typography variant="caption">Trend Direction</Typography>
                    </Paper>
                  </Grid>
                  <Grid item xs={12} sm={4}>
                    <Paper sx={{ p: 2, textAlign: 'center' }}>
                      <Typography variant="h6">
                        {playerData.trend_slope > 0 ? '+' : ''}{playerData.trend_slope.toFixed(1)}
                      </Typography>
                      <Typography variant="caption">GPM Change Rate</Typography>
                    </Paper>
                  </Grid>
                  <Grid item xs={12} sm={4}>
                    <Paper sx={{ p: 2, textAlign: 'center' }}>
                      <Typography variant="h6">
                        {(playerData.trend_confidence * 100).toFixed(0)}%
                      </Typography>
                      <Typography variant="caption">Confidence</Typography>
                    </Paper>
                  </Grid>
                </Grid>
              </Box>
            </Box>
          )}

          {selectedTab === 5 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Gold Optimization Recommendations
              </Typography>
              
              {/* Income Optimization */}
              <Paper sx={{ p: 2, mb: 2 }}>
                <Typography variant="subtitle1" gutterBottom>
                  <OptimizationIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
                  Income Improvement Potential
                </Typography>
                <Grid container spacing={2}>
                  <Grid item xs={12} sm={4}>
                    <Typography variant="body2" color="textSecondary">CS Improvement</Typography>
                    <Typography variant="h6">+{playerData.income_optimization.cs_improvement_potential.toFixed(0)} GPM</Typography>
                  </Grid>
                  <Grid item xs={12} sm={4}>
                    <Typography variant="body2" color="textSecondary">Kill Participation</Typography>
                    <Typography variant="h6">+{playerData.income_optimization.kp_improvement_potential.toFixed(0)} GPM</Typography>
                  </Grid>
                  <Grid item xs={12} sm={4}>
                    <Typography variant="body2" color="textSecondary">Objectives</Typography>
                    <Typography variant="h6">+{playerData.income_optimization.objective_improvement_potential.toFixed(0)} GPM</Typography>
                  </Grid>
                </Grid>
                <Typography variant="body2" sx={{ mt: 1, fontWeight: 'bold' }}>
                  Total Expected: +{playerData.income_optimization.expected_gpm_increase.toFixed(0)} GPM 
                  ({playerData.income_optimization.expected_win_rate_increase.toFixed(1)}% WR increase)
                </Typography>
              </Paper>

              {/* Recommendations List */}
              <Grid container spacing={2}>
                {playerData.recommendations.map((rec, index) => (
                  <Grid item xs={12} md={6} key={index}>
                    <Paper sx={{ p: 2, height: '100%' }}>
                      <Box display="flex" alignItems="flex-start" mb={1}>
                        {getPriorityIcon(rec.priority)}
                        <Box sx={{ ml: 1, flex: 1 }}>
                          <Box display="flex" alignItems="center" mb={1}>
                            <Chip 
                              label={rec.priority.toUpperCase()} 
                              size="small"
                              color={rec.priority === 'high' ? 'error' : rec.priority === 'medium' ? 'warning' : 'default'}
                              sx={{ mr: 1 }}
                            />
                            <Chip 
                              label={rec.category} 
                              size="small"
                              variant="outlined"
                            />
                          </Box>
                          <Typography variant="subtitle1" gutterBottom>
                            {rec.title}
                          </Typography>
                          <Typography variant="body2" color="textSecondary" paragraph>
                            {rec.description}
                          </Typography>
                          <Typography variant="caption" color="primary">
                            Impact: {rec.impact}
                          </Typography>
                          {rec.expected_gpm_increase > 0 && (
                            <Typography variant="caption" display="block" color="success.main">
                              Expected: +{rec.expected_gpm_increase.toFixed(0)} GPM
                            </Typography>
                          )}
                          <Box mt={1}>
                            {rec.game_phase.map((phase) => (
                              <Chip 
                                key={phase}
                                label={`${phase} game`}
                                size="small"
                                sx={{ 
                                  mr: 0.5, 
                                  bgcolor: PHASE_COLORS[phase as keyof typeof PHASE_COLORS], 
                                  color: 'white' 
                                }}
                              />
                            ))}
                          </Box>
                        </Box>
                      </Box>
                    </Paper>
                  </Grid>
                ))}
              </Grid>
            </Box>
          )}
        </CardContent>
      </Card>
    </Box>
  );
};

export default GoldEfficiency;