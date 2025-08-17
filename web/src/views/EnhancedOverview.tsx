import React, { useState, useMemo, useCallback } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  Skeleton,
  Alert,
  Chip,
  Avatar,
  LinearProgress,
  Divider,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  SportsEsports,
  Person,
  Timeline,
  EmojiEvents,
  Star,
  Visibility,
} from '@mui/icons-material';
import { Row } from '../types';
import { FilteringPanel } from '../components/filters/FilteringPanel';
import { FilterCriteria } from '../components/filters/AdvancedFilters';
import { AnimatedCard, StaggeredList, CounterAnimation } from '../components/animations/AnimatedComponents';
import { EnhancedErrorBoundary } from '../components/errors/EnhancedErrorBoundary';
import { StatCardSkeleton, ChartSkeleton } from '../components/loading/LoadingStates';

interface EnhancedOverviewProps {
  allData: Row[];
  data: Row[];
  loading: boolean;
  error: string | null;
  onRoleSelect: (role: string) => void;
}

interface StatCard {
  title: string;
  value: string | number;
  subtitle?: string;
  trend?: 'up' | 'down' | 'stable';
  color?: string;
  icon?: React.ReactNode;
  onClick?: () => void;
}

export const EnhancedOverview: React.FC<EnhancedOverviewProps> = ({
  allData,
  data: initialData,
  loading,
  error,
  onRoleSelect,
}) => {
  const [filteredData, setFilteredData] = useState<Row[]>(initialData);
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState<'list' | 'grid' | 'cards'>('grid');
  const [savedPresets, setSavedPresets] = useState<Array<{ name: string; filters: FilterCriteria }>>([]);

  // Application des filtres
  const applyFilters = useCallback((data: Row[], filters: FilterCriteria, query: string): Row[] => {
    let result = [...data];

    // Recherche textuelle
    if (query.trim()) {
      const lowerQuery = query.toLowerCase();
      result = result.filter(row => 
        row.champion?.toLowerCase().includes(lowerQuery) ||
        row.lane?.toLowerCase().includes(lowerQuery) ||
        row.summoner_name?.toLowerCase().includes(lowerQuery)
      );
    }

    // Filtres de champions
    if (filters.champions.length > 0) {
      result = result.filter(row => filters.champions.includes(row.champion));
    }

    // Filtres de rôles
    if (filters.roles.length > 0) {
      result = result.filter(row => filters.roles.includes(row.lane));
    }

    // Filtres de performance
    if (filters.winRate.min > 0 || filters.winRate.max < 100) {
      result = result.filter(row => {
        const winRate = row.win ? 100 : 0; // Simplification - dans un vrai cas, calculer le winrate
        return winRate >= filters.winRate.min && winRate <= filters.winRate.max;
      });
    }

    // Filtres KDA
    if (filters.kda.min > 0 || filters.kda.max < 10) {
      result = result.filter(row => {
        const kda = row.deaths > 0 ? (row.kills + row.assists) / row.deaths : row.kills + row.assists;
        return kda >= filters.kda.min && kda <= filters.kda.max;
      });
    }

    // Filtres de durée
    if (filters.gameDuration.min > 0 || filters.gameDuration.max < 60) {
      result = result.filter(row => {
        const durationMinutes = row.game_duration / 60;
        return durationMinutes >= filters.gameDuration.min && durationMinutes <= filters.gameDuration.max;
      });
    }

    // Tri
    result.sort((a, b) => {
      let aValue: any, bValue: any;
      
      switch (filters.sortBy) {
        case 'date':
          aValue = new Date(a.game_creation);
          bValue = new Date(b.game_creation);
          break;
        case 'kda':
          aValue = a.deaths > 0 ? (a.kills + a.assists) / a.deaths : a.kills + a.assists;
          bValue = b.deaths > 0 ? (b.kills + b.assists) / b.deaths : b.kills + b.assists;
          break;
        case 'duration':
          aValue = a.game_duration;
          bValue = b.game_duration;
          break;
        default:
          return 0;
      }
      
      const comparison = aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
      return filters.sortOrder === 'desc' ? -comparison : comparison;
    });

    return result;
  }, []);

  // Gestion des changements de filtres
  const handleFiltersChange = useCallback((filters: FilterCriteria) => {
    const filtered = applyFilters(allData, filters, searchQuery);
    setFilteredData(filtered);
  }, [allData, searchQuery, applyFilters]);

  // Gestion de la recherche
  const handleSearch = useCallback((query: string) => {
    setSearchQuery(query);
    // Les filtres seront appliqués via handleFiltersChange
  }, []);

  // Sauvegarde de presets
  const handleSavePreset = useCallback((name: string, filters: FilterCriteria) => {
    setSavedPresets(prev => [...prev, { name, filters }]);
  }, []);

  // Chargement de presets
  const handleLoadPreset = useCallback((filters: FilterCriteria) => {
    handleFiltersChange(filters);
  }, [handleFiltersChange]);

  // Calcul des statistiques
  const stats = useMemo(() => {
    const totalGames = filteredData.length;
    const wins = filteredData.filter(row => row.win).length;
    const winRate = totalGames > 0 ? (wins / totalGames) * 100 : 0;
    
    const totalKills = filteredData.reduce((sum, row) => sum + row.kills, 0);
    const totalDeaths = filteredData.reduce((sum, row) => sum + row.deaths, 0);
    const totalAssists = filteredData.reduce((sum, row) => sum + row.assists, 0);
    const avgKDA = totalDeaths > 0 ? (totalKills + totalAssists) / totalDeaths : totalKills + totalAssists;
    
    const avgGameDuration = totalGames > 0 
      ? filteredData.reduce((sum, row) => sum + row.game_duration, 0) / totalGames / 60 
      : 0;

    // Champions les plus joués
    const championCounts = filteredData.reduce((acc, row) => {
      acc[row.champion] = (acc[row.champion] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
    
    const topChampions = Object.entries(championCounts)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5)
      .map(([champion, games]) => ({ champion, games }));

    // Rôles les plus joués
    const roleCounts = filteredData.reduce((acc, row) => {
      acc[row.lane] = (acc[row.lane] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
    
    const topRoles = Object.entries(roleCounts)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5)
      .map(([role, games]) => ({ role, games }));

    return {
      totalGames,
      winRate,
      avgKDA,
      avgGameDuration,
      topChampions,
      topRoles,
    };
  }, [filteredData]);

  // Cartes de statistiques
  const statCards: StatCard[] = [
    {
      title: 'Total Games',
      value: stats.totalGames,
      icon: <SportsEsports />,
      color: 'primary.main',
    },
    {
      title: 'Winrate',
      value: `${stats.winRate.toFixed(1)}%`,
      subtitle: `${filteredData.filter(r => r.win).length}W / ${filteredData.filter(r => !r.win).length}L`,
      trend: stats.winRate >= 50 ? 'up' : 'down',
      icon: stats.winRate >= 50 ? <TrendingUp /> : <TrendingDown />,
      color: stats.winRate >= 50 ? 'success.main' : 'error.main',
    },
    {
      title: 'KDA Moyen',
      value: stats.avgKDA.toFixed(2),
      trend: stats.avgKDA >= 2 ? 'up' : stats.avgKDA >= 1 ? 'stable' : 'down',
      icon: <EmojiEvents />,
      color: stats.avgKDA >= 2 ? 'success.main' : stats.avgKDA >= 1 ? 'warning.main' : 'error.main',
    },
    {
      title: 'Durée Moyenne',
      value: `${stats.avgGameDuration.toFixed(1)}min`,
      icon: <Timeline />,
      color: 'info.main',
    },
  ];

  if (loading) {
    return (
      <Box>
        <EnhancedErrorBoundary
          maxRetries={1}
          autoRetry={false}
          errorLevel="low"
        >
          <FilteringPanel
            onFiltersChange={handleFiltersChange}
            onSearch={handleSearch}
            onViewModeChange={setViewMode}
            loading={loading}
            totalResults={0}
            savedPresets={savedPresets}
            onSavePreset={handleSavePreset}
            onLoadPreset={handleLoadPreset}
          />
        </EnhancedErrorBoundary>
        
        <Grid container spacing={3}>
          {[...Array(4)].map((_, i) => (
            <Grid item xs={12} sm={6} md={3} key={i}>
              <StatCardSkeleton />
            </Grid>
          ))}
        </Grid>
        
        <Grid container spacing={3} sx={{ mt: 2 }}>
          <Grid item xs={12} md={6}>
            <ChartSkeleton height={300} />
          </Grid>
          <Grid item xs={12} md={6}>
            <ChartSkeleton height={300} />
          </Grid>
        </Grid>
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 3 }}>
        {error}
      </Alert>
    );
  }

  return (
    <Box>
      {/* Panneau de filtrage avec protection d'erreur */}
      <EnhancedErrorBoundary
        maxRetries={2}
        autoRetry={false}
        errorLevel="medium"
        fallback={
          <Alert severity="warning" sx={{ mb: 3 }}>
            <AlertTitle>Erreur de filtrage</AlertTitle>
            Les filtres sont temporairement indisponibles. Les données sont affichées sans filtrage.
          </Alert>
        }
      >
        <FilteringPanel
          onFiltersChange={handleFiltersChange}
          onSearch={handleSearch}
          onViewModeChange={setViewMode}
          data={filteredData}
          loading={loading}
          totalResults={filteredData.length}
          savedPresets={savedPresets}
          onSavePreset={handleSavePreset}
          onLoadPreset={handleLoadPreset}
        />
      </EnhancedErrorBoundary>

      {/* Cartes de statistiques avec protection d'erreur */}
      <EnhancedErrorBoundary
        maxRetries={1}
        autoRetry={true}
        errorLevel="low"
      >
        <Grid container spacing={3} sx={{ mb: 3 }}>
          <StaggeredList staggerDelay={150} animationType="zoom">
            {statCards.map((stat, index) => (
              <Grid item xs={12} sm={6} md={3} key={index}>
                <AnimatedCard 
                  animationDelay={index * 100}
                  hoverEffect={true}
                  sx={{ 
                    cursor: stat.onClick ? 'pointer' : 'default',
                  }}
                  onClick={stat.onClick}
                >
                  <CardContent>
                    <Box display="flex" alignItems="center" justifyContent="space-between">
                      <Box>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          {stat.title}
                        </Typography>
                        <Typography variant="h5" fontWeight="bold" color={stat.color}>
                          {typeof stat.value === 'number' ? (
                            <CounterAnimation 
                              value={stat.value} 
                              duration={1000 + index * 200}
                            />
                          ) : (
                            stat.value
                          )}
                        </Typography>
                        {stat.subtitle && (
                          <Typography variant="caption" color="text.secondary">
                            {stat.subtitle}
                          </Typography>
                        )}
                      </Box>
                      <Box color={stat.color}>
                        {stat.icon}
                      </Box>
                    </Box>
                  </CardContent>
                </AnimatedCard>
              </Grid>
            ))}
          </StaggeredList>
        </Grid>
      </EnhancedErrorBoundary>

      {/* Top Champions et Rôles avec protection d'erreur */}
      <EnhancedErrorBoundary
        maxRetries={2}
        autoRetry={true}
        errorLevel="medium"
        fallback={
          <Alert severity="info" sx={{ mt: 3 }}>
            <AlertTitle>Données indisponibles</AlertTitle>
            Les classements des champions et rôles ne peuvent pas être affichés pour le moment.
          </Alert>
        }
      >
        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                  <Star color="warning" />
                  Champions les plus joués
                </Typography>
                <List dense>
                  {stats.topChampions.map(({ champion, games }, index) => (
                    <ListItem
                      key={champion}
                      button
                      onClick={() => console.log('Champion selected:', champion)}
                    >
                      <ListItemAvatar>
                        <Avatar sx={{ bgcolor: 'primary.main' }}>
                          {champion[0]}
                        </Avatar>
                      </ListItemAvatar>
                      <ListItemText
                        primary={champion}
                        secondary={`${games} partie${games !== 1 ? 's' : ''}`}
                      />
                      <Chip 
                        label={`#${index + 1}`} 
                        size="small" 
                        color="primary" 
                        variant="outlined" 
                      />
                    </ListItem>
                  ))}
                </List>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                  <SportsEsports color="secondary" />
                  Rôles les plus joués
                </Typography>
                <List dense>
                  {stats.topRoles.map(({ role, games }, index) => (
                    <ListItem
                      key={role}
                      button
                      onClick={() => onRoleSelect(role)}
                    >
                      <ListItemAvatar>
                        <Avatar sx={{ bgcolor: 'secondary.main' }}>
                          {role[0]}
                        </Avatar>
                      </ListItemAvatar>
                      <ListItemText
                        primary={role}
                        secondary={`${games} partie${games !== 1 ? 's' : ''}`}
                      />
                      <Box display="flex" alignItems="center" gap={1}>
                        <LinearProgress
                          variant="determinate"
                          value={(games / stats.totalGames) * 100}
                          sx={{ width: 60, mr: 1 }}
                        />
                        <Typography variant="caption">
                          {((games / stats.totalGames) * 100).toFixed(0)}%
                        </Typography>
                      </Box>
                    </ListItem>
                  ))}
                </List>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </EnhancedErrorBoundary>

      {/* Message si aucune donnée */}
      {filteredData.length === 0 && !loading && (
        <Alert severity="info" sx={{ mt: 3 }}>
          Aucune donnée disponible avec les filtres actuels. 
          Essayez d'ajuster vos critères de recherche.
        </Alert>
      )}
    </Box>
  );
};