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
  TableRow
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  LocalFireDepartment as DamageIcon,
  Group as TeamIcon,
  EmojiEvents as CarryIcon,
  Timeline as TrendIcon,
  CompareArrows as CompareIcon,
  Speed as EfficiencyIcon,
  Download as DownloadIcon,
  Refresh as RefreshIcon
} from '@mui/icons-material';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip, Legend, ResponsiveContainer, BarChart, Bar, PieChart, Pie, Cell, RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar } from 'recharts';

interface DamageAnalysis {
  player_id: string;
  champion?: string;
  position?: string;
  time_range: string;
  damage_share: number;
  damage_per_minute: number;
  total_damage: number;
  physical_damage_percent: number;
  magic_damage_percent: number;
  true_damage_percent: number;
  carry_potential: number;
  efficiency_rating: string;
  damage_consistency: number;
  team_contribution: TeamContributionData;
  damage_distribution: DamageDistribution;
  game_phase_analysis: GamePhaseAnalysis;
  high_damage_win_rate: number;
  low_damage_win_rate: number;
  role_benchmark: DamageBenchmark;
  rank_benchmark: DamageBenchmark;
  global_benchmark: DamageBenchmark;
  trend_data: DamageTrendPoint[];
  recommendations: DamageRecommendation[];
}

interface TeamContributionData {
  damage_contribution_score: number;
  kill_participation: number;
  solo_kill_rate: number;
  team_fight_damage_share: number;
  objective_damage_share: number;
  consistent_damage_games: number;
  clutch_performance_score: number;
}

interface DamageDistribution {
  champion_damage_percent: number;
  structure_damage_percent: number;
  monster_damage_percent: number;
  damage_by_type: {
    physical: number;
    magical: number;
    true_damage: number;
  };
  target_priority: {
    carries: number;
    tanks: number;
    supports: number;
  };
}

interface GamePhaseAnalysis {
  early_game: {
    damage_per_minute: number;
    damage_share: number;
    carry_potential: number;
    efficiency: string;
  };
  mid_game: {
    damage_per_minute: number;
    damage_share: number;
    carry_potential: number;
    efficiency: string;
  };
  late_game: {
    damage_per_minute: number;
    damage_share: number;
    carry_potential: number;
    efficiency: string;
  };
}

interface DamageBenchmark {
  category: string;
  average_damage_share: number;
  top_10_percent: number;
  top_25_percent: number;
  median: number;
  player_percentile: number;
}

interface DamageTrendPoint {
  date: string;
  damage_per_minute: number;
  damage_share: number;
  carry_potential: number;
  moving_average: number;
  efficiency: number;
}

interface DamageRecommendation {
  priority: 'high' | 'medium' | 'low';
  category: string;
  title: string;
  description: string;
  impact: string;
  game_phase: string[];
}

interface DamageAnalyticsProps {
  playerData: DamageAnalysis | null;
  timeRange: '7d' | '30d' | '90d';
  onTimeRangeChange: (timeRange: string) => void;
  loading?: boolean;
  error?: string;
}

const DAMAGE_COLORS = {
  physical: '#ff6b35',
  magical: '#4361ee',
  true: '#7209b7',
  total: '#06ffa5'
};

const PHASE_COLORS = {
  early: '#4cc9f0',
  mid: '#f9c74f',
  late: '#f8961e'
};

export const DamageAnalytics: React.FC<DamageAnalyticsProps> = ({
  playerData,
  timeRange,
  onTimeRangeChange,
  loading = false,
  error
}) => {
  const [selectedTab, setSelectedTab] = useState(0);
  const [selectedMetric, setSelectedMetric] = useState<'damage_share' | 'dpm' | 'carry_potential'>('damage_share');
  const [benchmarkType, setBenchmarkType] = useState<'role' | 'rank' | 'global'>('role');

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };

  const getDamageTypeData = () => {
    if (!playerData?.damage_distribution) return [];
    
    return [
      { name: 'Physical', value: playerData.damage_distribution.damage_by_type.physical, color: DAMAGE_COLORS.physical },
      { name: 'Magical', value: playerData.damage_distribution.damage_by_type.magical, color: DAMAGE_COLORS.magical },
      { name: 'True', value: playerData.damage_distribution.damage_by_type.true_damage, color: DAMAGE_COLORS.true }
    ];
  };

  const getRadarData = () => {
    if (!playerData) return [];
    
    return [
      {
        subject: 'Damage Share',
        player: playerData.damage_share,
        benchmark: getBenchmarkValue('damage_share'),
        fullMark: 100
      },
      {
        subject: 'DPM',
        player: (playerData.damage_per_minute / 1000) * 100,
        benchmark: 50,
        fullMark: 100
      },
      {
        subject: 'Consistency',
        player: playerData.damage_consistency,
        benchmark: 70,
        fullMark: 100
      },
      {
        subject: 'Carry Potential',
        player: playerData.carry_potential,
        benchmark: 50,
        fullMark: 100
      },
      {
        subject: 'Team Contribution',
        player: playerData.team_contribution.damage_contribution_score,
        benchmark: 50,
        fullMark: 100
      }
    ];
  };

  const getBenchmarkValue = (metric: string) => {
    const benchmark = benchmarkType === 'role' ? playerData?.role_benchmark : 
                     benchmarkType === 'rank' ? playerData?.rank_benchmark : 
                     playerData?.global_benchmark;
    
    switch (metric) {
      case 'damage_share':
        return benchmark?.average_damage_share || 25;
      default:
        return 50;
    }
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

  const formatDamage = (damage: number) => {
    if (damage >= 1000000) {
      return `${(damage / 1000000).toFixed(1)}M`;
    } else if (damage >= 1000) {
      return `${(damage / 1000).toFixed(0)}K`;
    }
    return damage.toString();
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
            <CircularProgress />
            <Typography variant="body2" sx={{ ml: 2 }}>
              Analyzing damage metrics...
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
        No damage data available for the selected time range.
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
            <DamageIcon sx={{ fontSize: 40, color: 'primary.main', mb: 1 }} />
            <Typography variant="h4" color="primary">
              {playerData.damage_share.toFixed(1)}%
            </Typography>
            <Typography variant="caption" color="textSecondary">
              Damage Share
            </Typography>
            <Box display="flex" alignItems="center" justifyContent="center" mt={1}>
              {playerData.damage_share > getBenchmarkValue('damage_share') ? (
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
            <TrendIcon sx={{ fontSize: 40, color: 'warning.main', mb: 1 }} />
            <Typography variant="h4" color="warning.main">
              {formatDamage(playerData.damage_per_minute)}
            </Typography>
            <Typography variant="caption" color="textSecondary">
              Damage per Minute
            </Typography>
            <Chip 
              label={playerData.efficiency_rating} 
              size="small" 
              color={getEfficiencyColor(playerData.efficiency_rating) as any}
              sx={{ mt: 1 }}
            />
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <CarryIcon sx={{ fontSize: 40, color: 'secondary.main', mb: 1 }} />
            <Typography variant="h4" color="secondary.main">
              {playerData.carry_potential.toFixed(0)}
            </Typography>
            <Typography variant="caption" color="textSecondary">
              Carry Potential
            </Typography>
            <LinearProgress 
              variant="determinate" 
              value={playerData.carry_potential} 
              sx={{ mt: 1, height: 8, borderRadius: 4 }}
            />
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <TeamIcon sx={{ fontSize: 40, color: 'info.main', mb: 1 }} />
            <Typography variant="h4" color="info.main">
              {playerData.team_contribution.damage_contribution_score.toFixed(0)}
            </Typography>
            <Typography variant="caption" color="textSecondary">
              Team Contribution
            </Typography>
            <Typography variant="body2" sx={{ mt: 1, fontSize: '0.75rem' }}>
              {playerData.team_contribution.kill_participation.toFixed(0)}% KP
            </Typography>
          </Paper>
        </Grid>
      </Grid>

      {/* Tabbed Analytics */}
      <Card>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={selectedTab} onChange={handleTabChange}>
            <Tab label="Overview" />
            <Tab label="Distribution" />
            <Tab label="Game Phases" />
            <Tab label="Trends" />
            <Tab label="Recommendations" />
          </Tabs>
        </Box>

        <CardContent>
          {selectedTab === 0 && (
            <Grid container spacing={3}>
              {/* Damage Performance Radar */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Performance Radar
                </Typography>
                <ResponsiveContainer width="100%" height={300}>
                  <RadarChart data={getRadarData()}>
                    <PolarGrid />
                    <PolarAngleAxis dataKey="subject" />
                    <PolarRadiusAxis angle={0} domain={[0, 100]} />
                    <Radar
                      name="Player"
                      dataKey="player"
                      stroke="#8884d8"
                      fill="#8884d8"
                      fillOpacity={0.3}
                    />
                    <Radar
                      name="Benchmark"
                      dataKey="benchmark"
                      stroke="#82ca9d"
                      fill="transparent"
                      strokeDasharray="5 5"
                    />
                    <Legend />
                  </RadarChart>
                </ResponsiveContainer>
              </Grid>

              {/* Win Rate Correlation */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Win Rate by Damage Performance
                </Typography>
                <Grid container spacing={2}>
                  <Grid item xs={6}>
                    <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'success.light' }}>
                      <Typography variant="h5" color="success.contrastText">
                        {playerData.high_damage_win_rate.toFixed(0)}%
                      </Typography>
                      <Typography variant="caption" color="success.contrastText">
                        High Damage Games
                      </Typography>
                    </Paper>
                  </Grid>
                  <Grid item xs={6}>
                    <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'error.light' }}>
                      <Typography variant="h5" color="error.contrastText">
                        {playerData.low_damage_win_rate.toFixed(0)}%
                      </Typography>
                      <Typography variant="caption" color="error.contrastText">
                        Low Damage Games
                      </Typography>
                    </Paper>
                  </Grid>
                </Grid>
                
                <Box sx={{ mt: 2 }}>
                  <Typography variant="body2" color="textSecondary">
                    Damage Impact: +{(playerData.high_damage_win_rate - playerData.low_damage_win_rate).toFixed(0)}% win rate increase
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          )}

          {selectedTab === 1 && (
            <Grid container spacing={3}>
              {/* Damage Type Distribution */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Damage Type Distribution
                </Typography>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={getDamageTypeData()}
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      fill="#8884d8"
                      dataKey="value"
                      label={({ name, value }) => `${name}: ${value.toFixed(1)}%`}
                    >
                      {getDamageTypeData().map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <RechartsTooltip />
                  </PieChart>
                </ResponsiveContainer>
              </Grid>

              {/* Target Distribution */}
              <Grid item xs={12} lg={6}>
                <Typography variant="h6" gutterBottom>
                  Damage Target Distribution
                </Typography>
                <Grid container spacing={1} sx={{ mt: 1 }}>
                  <Grid item xs={4}>
                    <Paper sx={{ p: 1, textAlign: 'center' }}>
                      <Typography variant="h6" color="primary">
                        {playerData.damage_distribution.champion_damage_percent.toFixed(0)}%
                      </Typography>
                      <Typography variant="caption">Champions</Typography>
                    </Paper>
                  </Grid>
                  <Grid item xs={4}>
                    <Paper sx={{ p: 1, textAlign: 'center' }}>
                      <Typography variant="h6" color="warning.main">
                        {playerData.damage_distribution.structure_damage_percent.toFixed(0)}%
                      </Typography>
                      <Typography variant="caption">Structures</Typography>
                    </Paper>
                  </Grid>
                  <Grid item xs={4}>
                    <Paper sx={{ p: 1, textAlign: 'center' }}>
                      <Typography variant="h6" color="info.main">
                        {playerData.damage_distribution.monster_damage_percent.toFixed(0)}%
                      </Typography>
                      <Typography variant="caption">Monsters</Typography>
                    </Paper>
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          )}

          {selectedTab === 2 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Game Phase Analysis
              </Typography>
              <TableContainer component={Paper}>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Phase</TableCell>
                      <TableCell align="right">Damage/Min</TableCell>
                      <TableCell align="right">Damage Share</TableCell>
                      <TableCell align="right">Carry Potential</TableCell>
                      <TableCell align="center">Efficiency</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        <Chip label="Early Game" size="small" sx={{ bgcolor: PHASE_COLORS.early, color: 'white' }} />
                      </TableCell>
                      <TableCell align="right">{formatDamage(playerData.game_phase_analysis.early_game.damage_per_minute)}</TableCell>
                      <TableCell align="right">{playerData.game_phase_analysis.early_game.damage_share.toFixed(1)}%</TableCell>
                      <TableCell align="right">{playerData.game_phase_analysis.early_game.carry_potential.toFixed(0)}</TableCell>
                      <TableCell align="center">
                        <Chip 
                          label={playerData.game_phase_analysis.early_game.efficiency}
                          size="small"
                          color={getEfficiencyColor(playerData.game_phase_analysis.early_game.efficiency) as any}
                        />
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        <Chip label="Mid Game" size="small" sx={{ bgcolor: PHASE_COLORS.mid, color: 'white' }} />
                      </TableCell>
                      <TableCell align="right">{formatDamage(playerData.game_phase_analysis.mid_game.damage_per_minute)}</TableCell>
                      <TableCell align="right">{playerData.game_phase_analysis.mid_game.damage_share.toFixed(1)}%</TableCell>
                      <TableCell align="right">{playerData.game_phase_analysis.mid_game.carry_potential.toFixed(0)}</TableCell>
                      <TableCell align="center">
                        <Chip 
                          label={playerData.game_phase_analysis.mid_game.efficiency}
                          size="small"
                          color={getEfficiencyColor(playerData.game_phase_analysis.mid_game.efficiency) as any}
                        />
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        <Chip label="Late Game" size="small" sx={{ bgcolor: PHASE_COLORS.late, color: 'white' }} />
                      </TableCell>
                      <TableCell align="right">{formatDamage(playerData.game_phase_analysis.late_game.damage_per_minute)}</TableCell>
                      <TableCell align="right">{playerData.game_phase_analysis.late_game.damage_share.toFixed(1)}%</TableCell>
                      <TableCell align="right">{playerData.game_phase_analysis.late_game.carry_potential.toFixed(0)}</TableCell>
                      <TableCell align="center">
                        <Chip 
                          label={playerData.game_phase_analysis.late_game.efficiency}
                          size="small"
                          color={getEfficiencyColor(playerData.game_phase_analysis.late_game.efficiency) as any}
                        />
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          )}

          {selectedTab === 3 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Performance Trends
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
                    dataKey="damage_share" 
                    stroke="#8884d8" 
                    name="Damage Share %" 
                  />
                  <Line 
                    type="monotone" 
                    dataKey="carry_potential" 
                    stroke="#82ca9d" 
                    name="Carry Potential"
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
            </Box>
          )}

          {selectedTab === 4 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Improvement Recommendations
              </Typography>
              <Grid container spacing={2}>
                {playerData.recommendations.map((rec, index) => (
                  <Grid item xs={12} md={6} key={index}>
                    <Paper sx={{ p: 2 }}>
                      <Box display="flex" alignItems="flex-start" mb={1}>
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

export default DamageAnalytics;