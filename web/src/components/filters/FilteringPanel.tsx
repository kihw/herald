import React, { useState, useCallback, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  IconButton,
  Tooltip,
  Chip,
  Grid,
  Fade,
  Collapse,
  Alert,
  LinearProgress,
  Divider,
} from '@mui/material';
import {
  Tune,
  Search,
  Clear,
  FilterList,
  Save,
  Restore,
  ViewList,
  ViewModule,
  Sort,
} from '@mui/icons-material';
import dayjs from 'dayjs';
import { AdvancedFilters, FilterCriteria } from './AdvancedFilters';
import { SmartSearch, SearchResult } from '../search/SmartSearch';

export interface FilteringPanelProps {
  onFiltersChange: (filters: FilterCriteria) => void;
  onSearch: (query: string) => void;
  onViewModeChange?: (mode: 'list' | 'grid' | 'cards') => void;
  data?: any[];
  loading?: boolean;
  totalResults?: number;
  showAdvancedFilters?: boolean;
  showSearch?: boolean;
  showViewControls?: boolean;
  savedPresets?: Array<{ name: string; filters: FilterCriteria }>;
  onSavePreset?: (name: string, filters: FilterCriteria) => void;
  onLoadPreset?: (filters: FilterCriteria) => void;
}

export const FilteringPanel: React.FC<FilteringPanelProps> = ({
  onFiltersChange,
  onSearch,
  onViewModeChange,
  data = [],
  loading = false,
  totalResults = 0,
  showAdvancedFilters = true,
  showSearch = true,
  showViewControls = true,
  savedPresets = [],
  onSavePreset,
  onLoadPreset,
}) => {
  const [showFilters, setShowFilters] = useState(false);
  const [activeFilters, setActiveFilters] = useState<FilterCriteria>();
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState<'list' | 'grid' | 'cards'>('list');
  const [appliedFiltersCount, setAppliedFiltersCount] = useState(0);
  const [recentSearches, setRecentSearches] = useState<string[]>([]);
  const [popularSearches] = useState<string[]>([
    'Jinx ADC',
    'Yasuo Mid',
    'Thresh Support',
    'filter:winrate>70',
    'rank:diamond',
  ]);

  // Calcul du nombre de filtres appliqués
  const calculateAppliedFilters = useCallback((filters: FilterCriteria) => {
    let count = 0;
    
    if (filters.searchQuery) count++;
    if (filters.dateRange.start || filters.dateRange.end) count++;
    if (filters.champions.length > 0) count++;
    if (filters.roles.length > 0) count++;
    if (filters.gameMode.length > 0) count++;
    if (filters.rank.length > 0) count++;
    if (filters.winRate.min > 0 || filters.winRate.max < 100) count++;
    if (filters.kda.min > 0 || filters.kda.max < 10) count++;
    if (filters.gameDuration.min > 0 || filters.gameDuration.max < 60) count++;
    if (filters.firstBlood !== null) count++;
    if (filters.pentaKills !== null) count++;
    if (filters.multikills.length > 0) count++;
    
    return count;
  }, []);

  // Gestion des changements de filtres
  const handleFiltersChange = useCallback((filters: FilterCriteria) => {
    setActiveFilters(filters);
    setAppliedFiltersCount(calculateAppliedFilters(filters));
    onFiltersChange(filters);
  }, [onFiltersChange, calculateAppliedFilters]);

  // Gestion de la recherche
  const handleSearch = useCallback((query: string) => {
    setSearchQuery(query);
    
    // Ajouter à l'historique
    if (query.trim() && !recentSearches.includes(query)) {
      setRecentSearches(prev => [query, ...prev.slice(0, 9)]);
    }
    
    onSearch(query);
  }, [onSearch, recentSearches]);

  // Gestion de la sélection de résultats de recherche
  const handleSearchResultSelect = useCallback((result: SearchResult) => {
    switch (result.type) {
      case 'champion':
        // Ajouter le champion aux filtres
        if (activeFilters) {
          const newFilters = {
            ...activeFilters,
            champions: [...activeFilters.champions, result.title],
          };
          handleFiltersChange(newFilters);
        }
        break;
      
      case 'player':
        // Lancer une recherche de joueur
        handleSearch(result.title);
        break;
      
      case 'filter':
        // Appliquer le filtre
        if (result.metadata && activeFilters) {
          const newFilters = { ...activeFilters };
          
          switch (result.metadata.filterType) {
            case 'winrate':
              newFilters.winRate = { min: result.metadata.value, max: 100 };
              break;
            case 'rank':
              newFilters.rank = [result.metadata.value];
              break;
            case 'time':
              // Gérer les filtres temporels
              const now = new Date();
              switch (result.metadata.period) {
                case 'today':
                  newFilters.dateRange = {
                    start: dayjs(now).startOf('day'),
                    end: dayjs(now).endOf('day'),
                  };
                  break;
                case 'week':
                  newFilters.dateRange = {
                    start: dayjs(now).subtract(7, 'days'),
                    end: dayjs(now),
                  };
                  break;
                case 'month':
                  newFilters.dateRange = {
                    start: dayjs(now).subtract(30, 'days'),
                    end: dayjs(now),
                  };
                  break;
              }
              break;
          }
          
          handleFiltersChange(newFilters);
        }
        break;
      
      case 'command':
        // Traiter les commandes
        if (result.title.startsWith('filter:') || result.title.startsWith('sort:')) {
          handleSearch(result.title);
        }
        break;
    }
  }, [activeFilters, handleFiltersChange, handleSearch]);

  // Application de filtres rapides
  const handleQuickFilter = (type: string, value: any) => {
    if (!activeFilters) return;
    
    const newFilters = { ...activeFilters };
    
    switch (type) {
      case 'clear':
        // Réinitialiser tous les filtres
        const defaultFilters: FilterCriteria = {
          searchQuery: '',
          dateRange: { start: null, end: null },
          season: [],
          patch: [],
          champions: [],
          roles: [],
          gameMode: [],
          mapId: [],
          queueType: [],
          winRate: { min: 0, max: 100 },
          kda: { min: 0, max: 10 },
          gameDuration: { min: 0, max: 60 },
          rank: [],
          killParticipation: { min: 0, max: 100 },
          csPerMinute: { min: 0, max: 15 },
          visionScore: { min: 0, max: 200 },
          goldDifference: { min: -10000, max: 10000 },
          firstBlood: null,
          pentaKills: null,
          multikills: [],
          itemBuilds: [],
          summoners: [],
          sortBy: 'date',
          sortOrder: 'desc',
          groupBy: 'none',
          showInactive: false,
        };
        handleFiltersChange(defaultFilters);
        setSearchQuery('');
        break;
      
      case 'role':
        newFilters.roles = newFilters.roles.includes(value)
          ? newFilters.roles.filter(r => r !== value)
          : [...newFilters.roles, value];
        break;
      
      case 'rank':
        newFilters.rank = newFilters.rank.includes(value)
          ? newFilters.rank.filter(r => r !== value)
          : [...newFilters.rank, value];
        break;
    }
    
    handleFiltersChange(newFilters);
  };

  // Gestion des modes d'affichage
  const handleViewModeChange = (mode: 'list' | 'grid' | 'cards') => {
    setViewMode(mode);
    onViewModeChange?.(mode);
  };

  // Filtres rapides prédéfinis
  const quickFilters = [
    { type: 'role', value: 'ADC', label: 'ADC' },
    { type: 'role', value: 'Support', label: 'Support' },
    { type: 'role', value: 'Mid', label: 'Mid' },
    { type: 'rank', value: 'Gold', label: 'Gold+' },
    { type: 'rank', value: 'Platinum', label: 'Platinum+' },
  ];

  return (
    <Paper elevation={2} sx={{ p: 2, mb: 2 }}>
      {/* Barre de recherche et contrôles principaux */}
      <Grid container spacing={2} alignItems="center">
        {showSearch && (
          <Grid item xs={12} md={6}>
            <SmartSearch
              placeholder="Rechercher champions, joueurs, filtres..."
              onSearch={handleSearch}
              onResultSelect={handleSearchResultSelect}
              recentSearches={recentSearches}
              popularSearches={popularSearches}
              enableHistory={true}
              enableSuggestions={true}
              enableFilters={true}
            />
          </Grid>
        )}
        
        <Grid item xs={12} md={6}>
          <Box display="flex" justifyContent="flex-end" alignItems="center" gap={1}>
            {/* Statistiques */}
            <Box display="flex" alignItems="center" gap={1} mr={2}>
              <Typography variant="body2" color="text.secondary">
                {totalResults} résultat{totalResults !== 1 ? 's' : ''}
              </Typography>
              {appliedFiltersCount > 0 && (
                <Chip
                  size="small"
                  label={`${appliedFiltersCount} filtre${appliedFiltersCount !== 1 ? 's' : ''}`}
                  color="primary"
                  variant="outlined"
                />
              )}
            </Box>
            
            {/* Contrôles d'affichage */}
            {showViewControls && (
              <Box display="flex" gap={0.5}>
                <Tooltip title="Vue liste">
                  <IconButton
                    size="small"
                    color={viewMode === 'list' ? 'primary' : 'default'}
                    onClick={() => handleViewModeChange('list')}
                  >
                    <ViewList />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Vue grille">
                  <IconButton
                    size="small"
                    color={viewMode === 'grid' ? 'primary' : 'default'}
                    onClick={() => handleViewModeChange('grid')}
                  >
                    <ViewModule />
                  </IconButton>
                </Tooltip>
              </Box>
            )}
            
            {/* Boutons d'action */}
            <Tooltip title="Filtres avancés">
              <IconButton
                onClick={() => setShowFilters(!showFilters)}
                color={showFilters ? 'primary' : 'default'}
              >
                <Tune />
              </IconButton>
            </Tooltip>
            
            <Tooltip title="Effacer tous les filtres">
              <IconButton
                onClick={() => handleQuickFilter('clear', null)}
                disabled={appliedFiltersCount === 0}
              >
                <Clear />
              </IconButton>
            </Tooltip>
          </Box>
        </Grid>
      </Grid>

      {/* Indicateur de chargement */}
      {loading && (
        <Box sx={{ mt: 1 }}>
          <LinearProgress />
        </Box>
      )}

      {/* Filtres rapides */}
      <Box sx={{ mt: 2 }}>
        <Box display="flex" flexWrap="wrap" gap={1} alignItems="center">
          <Typography variant="body2" color="text.secondary" sx={{ mr: 1 }}>
            Filtres rapides:
          </Typography>
          {quickFilters.map((filter) => {
            const isActive = activeFilters && (
              (filter.type === 'role' && activeFilters.roles.includes(filter.value)) ||
              (filter.type === 'rank' && activeFilters.rank.includes(filter.value))
            );
            
            return (
              <Chip
                key={`${filter.type}-${filter.value}`}
                label={filter.label}
                size="small"
                clickable
                color={isActive ? 'primary' : 'default'}
                variant={isActive ? 'filled' : 'outlined'}
                onClick={() => handleQuickFilter(filter.type, filter.value)}
              />
            );
          })}
        </Box>
      </Box>

      {/* Panneau de filtres avancés */}
      <Collapse in={showFilters} timeout={300}>
        <Box sx={{ mt: 2 }}>
          <Divider sx={{ mb: 2 }} />
          <Fade in={showFilters}>
            <div>
              {showAdvancedFilters && (
                <AdvancedFilters
                  onFiltersChange={handleFiltersChange}
                  onSavePreset={onSavePreset}
                  onLoadPreset={onLoadPreset}
                  savedPresets={savedPresets}
                  showAdvanced={true}
                />
              )}
            </div>
          </Fade>
        </Box>
      </Collapse>

      {/* Messages d'information */}
      {totalResults === 0 && !loading && (
        <Alert severity="info" sx={{ mt: 2 }}>
          Aucun résultat trouvé. Essayez d'ajuster vos filtres ou votre recherche.
        </Alert>
      )}
    </Paper>
  );
};