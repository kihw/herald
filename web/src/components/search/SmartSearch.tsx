import React, { useState, useEffect, useCallback, useMemo } from 'react';
import {
  Box,
  TextField,
  Paper,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Typography,
  Chip,
  IconButton,
  InputAdornment,
  Popper,
  ClickAwayListener,
  Fade,
  Avatar,
  Divider,
  ListSubheader,
} from '@mui/material';
import {
  Search,
  Close,
  History,
  Person,
  SportsEsports,
  Timeline,
  EmojiEvents,
  Star,
  TrendingUp,
  Schedule,
} from '@mui/icons-material';

export interface SearchResult {
  id: string;
  type: 'champion' | 'player' | 'match' | 'filter' | 'command';
  title: string;
  subtitle?: string;
  description?: string;
  metadata?: Record<string, any>;
  icon?: React.ReactNode;
  tags?: string[];
  score?: number;
}

export interface SearchSuggestion {
  query: string;
  type: 'recent' | 'popular' | 'autocomplete';
  timestamp?: Date;
  count?: number;
}

export interface SmartSearchProps {
  placeholder?: string;
  onSearch: (query: string) => void;
  onResultSelect: (result: SearchResult) => void;
  onFilterApply?: (filter: any) => void;
  searchResults?: SearchResult[];
  loading?: boolean;
  recentSearches?: string[];
  popularSearches?: string[];
  enableHistory?: boolean;
  enableSuggestions?: boolean;
  enableFilters?: boolean;
  maxResults?: number;
}

// Données simulées pour les suggestions
const MOCK_CHAMPIONS = [
  { name: 'Jinx', role: 'ADC', popularity: 85 },
  { name: 'Yasuo', role: 'Mid', popularity: 92 },
  { name: 'Thresh', role: 'Support', popularity: 78 },
  { name: 'Lee Sin', role: 'Jungle', popularity: 81 },
  { name: 'Caitlyn', role: 'ADC', popularity: 76 },
  { name: 'Ahri', role: 'Mid', popularity: 73 },
  { name: 'Leona', role: 'Support', popularity: 69 },
  { name: 'Darius', role: 'Top', popularity: 72 },
];

const MOCK_COMMANDS = [
  { command: 'filter:winrate>70', description: 'Filtrer par winrate supérieur à 70%' },
  { command: 'filter:role:adc', description: 'Filtrer par rôle ADC' },
  { command: 'filter:rank:gold', description: 'Filtrer par rang Gold' },
  { command: 'sort:kda:desc', description: 'Trier par KDA décroissant' },
  { command: 'export:csv', description: 'Exporter en CSV' },
  { command: 'analyze:trends', description: 'Analyser les tendances' },
];

export const SmartSearch: React.FC<SmartSearchProps> = ({
  placeholder = 'Rechercher champions, joueurs, filtres...',
  onSearch,
  onResultSelect,
  onFilterApply,
  searchResults = [],
  loading = false,
  recentSearches = [],
  popularSearches = [],
  enableHistory = true,
  enableSuggestions = true,
  enableFilters = true,
  maxResults = 10,
}) => {
  const [query, setQuery] = useState('');
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const [suggestions, setSuggestions] = useState<SearchResult[]>([]);
  const [recentQueries, setRecentQueries] = useState<string[]>(recentSearches);
  const [showSuggestions, setShowSuggestions] = useState(false);

  // Génération intelligente de suggestions
  const generateSuggestions = useCallback((searchQuery: string): SearchResult[] => {
    if (!searchQuery.trim()) {
      // Suggestions par défaut (historique + populaires)
      const results: SearchResult[] = [];
      
      if (enableHistory && recentQueries.length > 0) {
        results.push(
          ...recentQueries.slice(0, 3).map((q, index) => ({
            id: `recent-${index}`,
            type: 'command' as const,
            title: q,
            subtitle: 'Recherche récente',
            icon: <History color="action" />,
            score: 100 - index,
          }))
        );
      }
      
      if (popularSearches.length > 0) {
        results.push(
          ...popularSearches.slice(0, 3).map((q, index) => ({
            id: `popular-${index}`,
            type: 'command' as const,
            title: q,
            subtitle: 'Recherche populaire',
            icon: <TrendingUp color="action" />,
            score: 90 - index,
          }))
        );
      }
      
      return results;
    }

    const results: SearchResult[] = [];
    const lowerQuery = searchQuery.toLowerCase();
    
    // Recherche de champions
    MOCK_CHAMPIONS.forEach(champion => {
      const championMatch = champion.name.toLowerCase().includes(lowerQuery);
      const roleMatch = champion.role.toLowerCase().includes(lowerQuery);
      
      if (championMatch || roleMatch) {
        results.push({
          id: `champion-${champion.name}`,
          type: 'champion',
          title: champion.name,
          subtitle: `${champion.role} • Popularité: ${champion.popularity}%`,
          description: `Champion populaire en ${champion.role}`,
          icon: <Avatar sx={{ width: 24, height: 24 }}>{champion.name[0]}</Avatar>,
          tags: [champion.role, 'champion'],
          score: championMatch ? champion.popularity + 20 : champion.popularity,
          metadata: { role: champion.role, popularity: champion.popularity },
        });
      }
    });
    
    // Détection de commandes
    if (lowerQuery.includes('filter:') || lowerQuery.includes('sort:') || lowerQuery.includes('export:')) {
      MOCK_COMMANDS.forEach((cmd, index) => {
        if (cmd.command.toLowerCase().includes(lowerQuery) || lowerQuery.includes(cmd.command.split(':')[0])) {
          results.push({
            id: `command-${index}`,
            type: 'command',
            title: cmd.command,
            subtitle: 'Commande',
            description: cmd.description,
            icon: <SportsEsports color="primary" />,
            score: 95,
          });
        }
      });
    }
    
    // Recherche de filtres intelligents
    if (enableFilters) {
      const filterSuggestions = generateFilterSuggestions(lowerQuery);
      results.push(...filterSuggestions);
    }
    
    // Recherche de joueurs simulée
    if (lowerQuery.length >= 3 && !lowerQuery.includes(':')) {
      results.push({
        id: `player-${lowerQuery}`,
        type: 'player',
        title: lowerQuery,
        subtitle: 'Rechercher ce joueur',
        description: 'Voir les statistiques de ce joueur',
        icon: <Person color="secondary" />,
        score: 80,
      });
    }
    
    // Tri par pertinence
    return results
      .sort((a, b) => (b.score || 0) - (a.score || 0))
      .slice(0, maxResults);
  }, [recentQueries, popularSearches, enableHistory, enableFilters, maxResults]);

  // Génération de suggestions de filtres
  const generateFilterSuggestions = (query: string): SearchResult[] => {
    const suggestions: SearchResult[] = [];
    
    // Suggestions de winrate
    if (query.includes('win') || query.includes('victoire')) {
      suggestions.push({
        id: 'filter-winrate-high',
        type: 'filter',
        title: 'Winrate > 70%',
        subtitle: 'Filtre de performance',
        description: 'Afficher uniquement les matches avec plus de 70% de winrate',
        icon: <EmojiEvents color="success" />,
        score: 85,
        metadata: { filterType: 'winrate', operator: '>', value: 70 },
      });
    }
    
    // Suggestions de rang
    if (query.includes('gold') || query.includes('plat') || query.includes('diamond')) {
      const rank = query.includes('gold') ? 'Gold' : query.includes('plat') ? 'Platinum' : 'Diamond';
      suggestions.push({
        id: `filter-rank-${rank.toLowerCase()}`,
        type: 'filter',
        title: `Rang: ${rank}`,
        subtitle: 'Filtre de rang',
        description: `Afficher uniquement les matches en ${rank}`,
        icon: <Star color="warning" />,
        score: 80,
        metadata: { filterType: 'rank', value: rank },
      });
    }
    
    // Suggestions temporelles
    if (query.includes('aujourd') || query.includes('semaine') || query.includes('mois')) {
      const period = query.includes('aujourd') ? 'today' : query.includes('semaine') ? 'week' : 'month';
      const label = period === 'today' ? "Aujourd'hui" : period === 'week' ? 'Cette semaine' : 'Ce mois';
      
      suggestions.push({
        id: `filter-time-${period}`,
        type: 'filter',
        title: label,
        subtitle: 'Filtre temporel',
        description: `Afficher uniquement les matches de ${label.toLowerCase()}`,
        icon: <Schedule color="info" />,
        score: 75,
        metadata: { filterType: 'time', period },
      });
    }
    
    return suggestions;
  };

  // Mise à jour des suggestions en temps réel
  useEffect(() => {
    const newSuggestions = generateSuggestions(query);
    setSuggestions(newSuggestions);
  }, [query, generateSuggestions]);

  // Gestion de la recherche
  const handleSearch = useCallback((searchQuery: string) => {
    if (searchQuery.trim()) {
      // Ajouter à l'historique
      setRecentQueries(prev => {
        const updated = [searchQuery, ...prev.filter(q => q !== searchQuery)];
        return updated.slice(0, 10); // Garder seulement les 10 dernières
      });
      
      onSearch(searchQuery);
      setShowSuggestions(false);
    }
  }, [onSearch]);

  // Sélection d'un résultat
  const handleResultSelect = (result: SearchResult) => {
    setQuery(result.title);
    
    if (result.type === 'filter' && onFilterApply) {
      onFilterApply(result.metadata);
    } else if (result.type === 'command') {
      if (result.title.startsWith('filter:') || result.title.startsWith('sort:')) {
        // Traitement des commandes
        handleSearch(result.title);
      } else {
        setQuery(result.title);
        handleSearch(result.title);
      }
    } else {
      onResultSelect(result);
    }
    
    setShowSuggestions(false);
  };

  // Gestion du focus
  const handleFocus = (event: React.FocusEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
    setShowSuggestions(true);
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      event.preventDefault();
      handleSearch(query);
    } else if (event.key === 'Escape') {
      setShowSuggestions(false);
    }
  };

  const renderSuggestionIcon = (result: SearchResult) => {
    if (result.icon) return result.icon;
    
    switch (result.type) {
      case 'champion':
        return <SportsEsports color="primary" />;
      case 'player':
        return <Person color="secondary" />;
      case 'filter':
        return <Timeline color="info" />;
      case 'command':
        return <History color="action" />;
      default:
        return <Search color="action" />;
    }
  };

  const groupedSuggestions = useMemo(() => {
    const groups: Record<string, SearchResult[]> = {};
    
    suggestions.forEach(suggestion => {
      const group = suggestion.type === 'command' && query === '' 
        ? (suggestion.subtitle === 'Recherche récente' ? 'Récentes' : 'Populaires')
        : suggestion.type === 'champion' ? 'Champions'
        : suggestion.type === 'player' ? 'Joueurs'
        : suggestion.type === 'filter' ? 'Filtres'
        : 'Commandes';
      
      if (!groups[group]) groups[group] = [];
      groups[group].push(suggestion);
    });
    
    return groups;
  }, [suggestions, query]);

  return (
    <Box sx={{ position: 'relative', width: '100%' }}>
      <TextField
        fullWidth
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onFocus={handleFocus}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <Search color="action" />
            </InputAdornment>
          ),
          endAdornment: query && (
            <InputAdornment position="end">
              <IconButton
                size="small"
                onClick={() => {
                  setQuery('');
                  setShowSuggestions(false);
                }}
              >
                <Close />
              </IconButton>
            </InputAdornment>
          ),
        }}
      />
      
      <Popper
        open={showSuggestions && (suggestions.length > 0 || query === '')}
        anchorEl={anchorEl}
        placement="bottom-start"
        style={{ width: anchorEl?.clientWidth, zIndex: 1300 }}
        transition
      >
        {({ TransitionProps }) => (
          <Fade {...TransitionProps} timeout={200}>
            <Paper elevation={8} sx={{ maxHeight: 400, overflow: 'auto', mt: 0.5 }}>
              <ClickAwayListener onClickAway={() => setShowSuggestions(false)}>
                <List dense>
                  {Object.entries(groupedSuggestions).map(([groupName, groupResults]) => (
                    <React.Fragment key={groupName}>
                      {Object.keys(groupedSuggestions).length > 1 && (
                        <ListSubheader component="div" sx={{ bgcolor: 'background.paper' }}>
                          {groupName}
                        </ListSubheader>
                      )}
                      {groupResults.map((result) => (
                        <ListItem
                          key={result.id}
                          button
                          onClick={() => handleResultSelect(result)}
                          sx={{
                            '&:hover': {
                              backgroundColor: 'action.hover',
                            },
                          }}
                        >
                          <ListItemIcon sx={{ minWidth: 40 }}>
                            {renderSuggestionIcon(result)}
                          </ListItemIcon>
                          <ListItemText
                            primary={
                              <Box display="flex" alignItems="center" gap={1}>
                                <Typography variant="body2">
                                  {result.title}
                                </Typography>
                                {result.tags && result.tags.map(tag => (
                                  <Chip
                                    key={tag}
                                    label={tag}
                                    size="small"
                                    variant="outlined"
                                    sx={{ height: 20, fontSize: '0.75rem' }}
                                  />
                                ))}
                              </Box>
                            }
                            secondary={
                              result.subtitle || result.description ? (
                                <Typography variant="caption" color="text.secondary">
                                  {result.subtitle}
                                  {result.description && ` • ${result.description}`}
                                </Typography>
                              ) : undefined
                            }
                          />
                        </ListItem>
                      ))}
                      {Object.keys(groupedSuggestions).length > 1 && 
                       groupName !== Object.keys(groupedSuggestions)[Object.keys(groupedSuggestions).length - 1] && (
                        <Divider />
                      )}
                    </React.Fragment>
                  ))}
                  
                  {suggestions.length === 0 && query !== '' && (
                    <ListItem>
                      <ListItemText
                        primary="Aucun résultat trouvé"
                        secondary="Essayez une autre recherche ou utilisez les filtres avancés"
                      />
                    </ListItem>
                  )}
                </List>
              </ClickAwayListener>
            </Paper>
          </Fade>
        )}
      </Popper>
    </Box>
  );
};