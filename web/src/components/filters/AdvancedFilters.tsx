import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Typography,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Checkbox,
  FormControlLabel,
  FormGroup,
  Chip,
  IconButton,
  Tooltip,
  Button,
  Grid,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Slider,
  ToggleButton,
  ToggleButtonGroup,
  Autocomplete,
  DatePicker,
  Switch,
  Divider,
  Badge,
  Collapse,
} from '@mui/material';
import {
  ExpandMore,
  FilterList,
  Clear,
  Search,
  Save,
  Restore,
  Tune,
  Close,
  Add,
  Remove,
} from '@mui/icons-material';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { DatePicker as MUIDatePicker } from '@mui/x-date-pickers/DatePicker';
import dayjs, { Dayjs } from 'dayjs';

export interface FilterCriteria {
  // Filtres de recherche de base
  searchQuery: string;
  
  // Filtres temporels
  dateRange: {
    start: Dayjs | null;
    end: Dayjs | null;
  };
  season: string[];
  patch: string[];
  
  // Filtres de jeu
  champions: string[];
  roles: string[];
  gameMode: string[];
  mapId: string[];
  queueType: string[];
  
  // Filtres de performance
  winRate: {
    min: number;
    max: number;
  };
  kda: {
    min: number;
    max: number;
  };
  gameDuration: {
    min: number;
    max: number;
  };
  rank: string[];
  
  // Filtres avancés
  killParticipation: {
    min: number;
    max: number;
  };
  csPerMinute: {
    min: number;
    max: number;
  };
  visionScore: {
    min: number;
    max: number;
  };
  goldDifference: {
    min: number;
    max: number;
  };
  
  // Filtres spéciaux
  firstBlood: boolean | null;
  pentaKills: boolean | null;
  multikills: string[];
  itemBuilds: string[];
  summoners: string[];
  
  // Options d'affichage
  sortBy: string;
  sortOrder: 'asc' | 'desc';
  groupBy: string;
  showInactive: boolean;
}

export interface AdvancedFiltersProps {
  onFiltersChange: (filters: FilterCriteria) => void;
  onSavePreset?: (name: string, filters: FilterCriteria) => void;
  onLoadPreset?: (filters: FilterCriteria) => void;
  savedPresets?: Array<{ name: string; filters: FilterCriteria }>;
  initialFilters?: Partial<FilterCriteria>;
  showAdvanced?: boolean;
}

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

// Données de référence
const CHAMPIONS = [
  'Jinx', 'Yasuo', 'Thresh', 'Leona', 'Zed', 'Caitlyn', 'Ahri', 'Lee Sin',
  'Vayne', 'Blitzcrank', 'Katarina', 'Darius', 'Lux', 'Ezreal', 'Morgana',
];

const ROLES = ['Top', 'Jungle', 'Mid', 'ADC', 'Support'];

const GAME_MODES = [
  'Ranked Solo', 'Ranked Flex', 'Normal', 'ARAM', 'URF', 'TFT', 'Arena'
];

const RANKS = [
  'Iron', 'Bronze', 'Silver', 'Gold', 'Platinum', 'Diamond', 'Master', 'Grandmaster', 'Challenger'
];

const SEASONS = ['2024', '2023', '2022', '2021'];

const PATCHES = ['14.15', '14.14', '14.13', '14.12', '14.11', '14.10'];

export const AdvancedFilters: React.FC<AdvancedFiltersProps> = ({
  onFiltersChange,
  onSavePreset,
  onLoadPreset,
  savedPresets = [],
  initialFilters = {},
  showAdvanced = true,
}) => {
  const [filters, setFilters] = useState<FilterCriteria>({
    ...defaultFilters,
    ...initialFilters,
  });
  
  const [expanded, setExpanded] = useState<string | false>('basic');
  const [showPresets, setShowPresets] = useState(false);
  const [presetName, setPresetName] = useState('');
  const [activeFiltersCount, setActiveFiltersCount] = useState(0);

  // Calcul du nombre de filtres actifs
  const calculateActiveFilters = useCallback((currentFilters: FilterCriteria) => {
    let count = 0;
    
    if (currentFilters.searchQuery) count++;
    if (currentFilters.dateRange.start || currentFilters.dateRange.end) count++;
    if (currentFilters.champions.length > 0) count++;
    if (currentFilters.roles.length > 0) count++;
    if (currentFilters.gameMode.length > 0) count++;
    if (currentFilters.rank.length > 0) count++;
    if (currentFilters.winRate.min > 0 || currentFilters.winRate.max < 100) count++;
    if (currentFilters.kda.min > 0 || currentFilters.kda.max < 10) count++;
    if (currentFilters.gameDuration.min > 0 || currentFilters.gameDuration.max < 60) count++;
    if (currentFilters.firstBlood !== null) count++;
    if (currentFilters.pentaKills !== null) count++;
    if (currentFilters.multikills.length > 0) count++;
    
    return count;
  }, []);

  // Mise à jour des filtres
  const updateFilters = useCallback((newFilters: Partial<FilterCriteria>) => {
    const updatedFilters = { ...filters, ...newFilters };
    setFilters(updatedFilters);
    setActiveFiltersCount(calculateActiveFilters(updatedFilters));
    onFiltersChange(updatedFilters);
  }, [filters, calculateActiveFilters, onFiltersChange]);

  // Réinitialisation des filtres
  const resetFilters = () => {
    setFilters(defaultFilters);
    setActiveFiltersCount(0);
    onFiltersChange(defaultFilters);
  };

  // Sauvegarde de preset
  const savePreset = () => {
    if (presetName.trim() && onSavePreset) {
      onSavePreset(presetName.trim(), filters);
      setPresetName('');
      setShowPresets(false);
    }
  };

  // Gestion des accordéons
  const handleAccordionChange = (panel: string) => (
    event: React.SyntheticEvent,
    isExpanded: boolean
  ) => {
    setExpanded(isExpanded ? panel : false);
  };

  useEffect(() => {
    setActiveFiltersCount(calculateActiveFilters(filters));
  }, [filters, calculateActiveFilters]);

  return (
    <LocalizationProvider dateAdapter={AdapterDayjs}>
      <Card>
        <CardHeader
          title={
            <Box display="flex" alignItems="center" gap={1}>
              <Badge badgeContent={activeFiltersCount} color="primary">
                <FilterList />
              </Badge>
              <Typography variant="h6">Filtres avancés</Typography>
            </Box>
          }
          action={
            <Box display="flex" gap={1}>
              <Tooltip title="Presets">
                <IconButton onClick={() => setShowPresets(!showPresets)}>
                  <Save />
                </IconButton>
              </Tooltip>
              <Tooltip title="Réinitialiser">
                <IconButton onClick={resetFilters}>
                  <Clear />
                </IconButton>
              </Tooltip>
            </Box>
          }
        />
        
        <CardContent>
          {/* Section des presets */}
          <Collapse in={showPresets}>
            <Box sx={{ mb: 2, p: 2, bgcolor: 'background.paper', borderRadius: 1 }}>
              <Typography variant="subtitle2" gutterBottom>
                Presets sauvegardés
              </Typography>
              <Grid container spacing={1} sx={{ mb: 2 }}>
                {savedPresets.map((preset, index) => (
                  <Grid item key={index}>
                    <Chip
                      label={preset.name}
                      onClick={() => {
                        setFilters(preset.filters);
                        onLoadPreset?.(preset.filters);
                      }}
                      variant="outlined"
                    />
                  </Grid>
                ))}
              </Grid>
              <Box display="flex" gap={1}>
                <TextField
                  size="small"
                  placeholder="Nom du preset"
                  value={presetName}
                  onChange={(e) => setPresetName(e.target.value)}
                />
                <Button onClick={savePreset} disabled={!presetName.trim()}>
                  Sauvegarder
                </Button>
              </Box>
            </Box>
          </Collapse>

          {/* Recherche rapide */}
          <Box sx={{ mb: 2 }}>
            <TextField
              fullWidth
              placeholder="Rechercher par champion, joueur, ou mots-clés..."
              value={filters.searchQuery}
              onChange={(e) => updateFilters({ searchQuery: e.target.value })}
              InputProps={{
                startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />,
                endAdornment: filters.searchQuery && (
                  <IconButton
                    size="small"
                    onClick={() => updateFilters({ searchQuery: '' })}
                  >
                    <Close />
                  </IconButton>
                ),
              }}
            />
          </Box>

          {/* Filtres de base */}
          <Accordion
            expanded={expanded === 'basic'}
            onChange={handleAccordionChange('basic')}
          >
            <AccordionSummary expandIcon={<ExpandMore />}>
              <Typography>Filtres de base</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                  <MUIDatePicker
                    label="Date de début"
                    value={filters.dateRange.start}
                    onChange={(date) =>
                      updateFilters({
                        dateRange: { ...filters.dateRange, start: date },
                      })
                    }
                    slotProps={{ textField: { fullWidth: true, size: 'small' } }}
                  />
                </Grid>
                <Grid item xs={12} md={6}>
                  <MUIDatePicker
                    label="Date de fin"
                    value={filters.dateRange.end}
                    onChange={(date) =>
                      updateFilters({
                        dateRange: { ...filters.dateRange, end: date },
                      })
                    }
                    slotProps={{ textField: { fullWidth: true, size: 'small' } }}
                  />
                </Grid>
                
                <Grid item xs={12} md={6}>
                  <Autocomplete
                    multiple
                    options={CHAMPIONS}
                    value={filters.champions}
                    onChange={(_, newValue) => updateFilters({ champions: newValue })}
                    renderInput={(params) => (
                      <TextField {...params} label="Champions" size="small" />
                    )}
                    renderTags={(value, getTagProps) =>
                      value.map((option, index) => (
                        <Chip
                          label={option}
                          size="small"
                          {...getTagProps({ index })}
                          key={option}
                        />
                      ))
                    }
                  />
                </Grid>
                
                <Grid item xs={12} md={6}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Rôles</InputLabel>
                    <Select
                      multiple
                      value={filters.roles}
                      onChange={(e) =>
                        updateFilters({
                          roles: typeof e.target.value === 'string'
                            ? e.target.value.split(',')
                            : e.target.value,
                        })
                      }
                      renderValue={(selected) => (
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {(selected as string[]).map((value) => (
                            <Chip key={value} label={value} size="small" />
                          ))}
                        </Box>
                      )}
                    >
                      {ROLES.map((role) => (
                        <MenuItem key={role} value={role}>
                          <Checkbox checked={filters.roles.includes(role)} />
                          {role}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
                
                <Grid item xs={12} md={6}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Mode de jeu</InputLabel>
                    <Select
                      multiple
                      value={filters.gameMode}
                      onChange={(e) =>
                        updateFilters({
                          gameMode: typeof e.target.value === 'string'
                            ? e.target.value.split(',')
                            : e.target.value,
                        })
                      }
                    >
                      {GAME_MODES.map((mode) => (
                        <MenuItem key={mode} value={mode}>
                          {mode}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
                
                <Grid item xs={12} md={6}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Rang</InputLabel>
                    <Select
                      multiple
                      value={filters.rank}
                      onChange={(e) =>
                        updateFilters({
                          rank: typeof e.target.value === 'string'
                            ? e.target.value.split(',')
                            : e.target.value,
                        })
                      }
                    >
                      {RANKS.map((rank) => (
                        <MenuItem key={rank} value={rank}>
                          {rank}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>
            </AccordionDetails>
          </Accordion>

          {/* Filtres de performance */}
          <Accordion
            expanded={expanded === 'performance'}
            onChange={handleAccordionChange('performance')}
          >
            <AccordionSummary expandIcon={<ExpandMore />}>
              <Typography>Filtres de performance</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <Typography gutterBottom>Winrate (%)</Typography>
                  <Slider
                    value={[filters.winRate.min, filters.winRate.max]}
                    onChange={(_, newValue) =>
                      updateFilters({
                        winRate: {
                          min: (newValue as number[])[0],
                          max: (newValue as number[])[1],
                        },
                      })
                    }
                    valueLabelDisplay="auto"
                    min={0}
                    max={100}
                    marks={[
                      { value: 0, label: '0%' },
                      { value: 50, label: '50%' },
                      { value: 100, label: '100%' },
                    ]}
                  />
                </Grid>
                
                <Grid item xs={12}>
                  <Typography gutterBottom>KDA</Typography>
                  <Slider
                    value={[filters.kda.min, filters.kda.max]}
                    onChange={(_, newValue) =>
                      updateFilters({
                        kda: {
                          min: (newValue as number[])[0],
                          max: (newValue as number[])[1],
                        },
                      })
                    }
                    valueLabelDisplay="auto"
                    min={0}
                    max={10}
                    step={0.1}
                    marks={[
                      { value: 0, label: '0' },
                      { value: 2, label: '2' },
                      { value: 5, label: '5' },
                      { value: 10, label: '10+' },
                    ]}
                  />
                </Grid>
                
                <Grid item xs={12}>
                  <Typography gutterBottom>Durée de jeu (minutes)</Typography>
                  <Slider
                    value={[filters.gameDuration.min, filters.gameDuration.max]}
                    onChange={(_, newValue) =>
                      updateFilters({
                        gameDuration: {
                          min: (newValue as number[])[0],
                          max: (newValue as number[])[1],
                        },
                      })
                    }
                    valueLabelDisplay="auto"
                    min={0}
                    max={60}
                    marks={[
                      { value: 0, label: '0min' },
                      { value: 20, label: '20min' },
                      { value: 40, label: '40min' },
                      { value: 60, label: '60min+' },
                    ]}
                  />
                </Grid>
              </Grid>
            </AccordionDetails>
          </Accordion>

          {/* Filtres avancés */}
          {showAdvanced && (
            <Accordion
              expanded={expanded === 'advanced'}
              onChange={handleAccordionChange('advanced')}
            >
              <AccordionSummary expandIcon={<ExpandMore />}>
                <Typography>Filtres avancés</Typography>
              </AccordionSummary>
              <AccordionDetails>
                <Grid container spacing={2}>
                  <Grid item xs={12} md={6}>
                    <FormGroup>
                      <FormControlLabel
                        control={
                          <Checkbox
                            checked={filters.firstBlood === true}
                            indeterminate={filters.firstBlood === null}
                            onChange={(e) =>
                              updateFilters({
                                firstBlood: e.target.checked ? true : null,
                              })
                            }
                          />
                        }
                        label="First Blood uniquement"
                      />
                      <FormControlLabel
                        control={
                          <Checkbox
                            checked={filters.pentaKills === true}
                            indeterminate={filters.pentaKills === null}
                            onChange={(e) =>
                              updateFilters({
                                pentaKills: e.target.checked ? true : null,
                              })
                            }
                          />
                        }
                        label="Penta kills uniquement"
                      />
                      <FormControlLabel
                        control={
                          <Switch
                            checked={filters.showInactive}
                            onChange={(e) =>
                              updateFilters({ showInactive: e.target.checked })
                            }
                          />
                        }
                        label="Inclure matches inactifs"
                      />
                    </FormGroup>
                  </Grid>
                  
                  <Grid item xs={12} md={6}>
                    <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                      <InputLabel>Trier par</InputLabel>
                      <Select
                        value={filters.sortBy}
                        onChange={(e) =>
                          updateFilters({ sortBy: e.target.value })
                        }
                      >
                        <MenuItem value="date">Date</MenuItem>
                        <MenuItem value="winrate">Winrate</MenuItem>
                        <MenuItem value="kda">KDA</MenuItem>
                        <MenuItem value="duration">Durée</MenuItem>
                        <MenuItem value="performance">Performance</MenuItem>
                      </Select>
                    </FormControl>
                    
                    <ToggleButtonGroup
                      value={filters.sortOrder}
                      exclusive
                      onChange={(_, newOrder) =>
                        newOrder && updateFilters({ sortOrder: newOrder })
                      }
                      size="small"
                      fullWidth
                    >
                      <ToggleButton value="asc">Croissant</ToggleButton>
                      <ToggleButton value="desc">Décroissant</ToggleButton>
                    </ToggleButtonGroup>
                  </Grid>
                </Grid>
              </AccordionDetails>
            </Accordion>
          )}

          {/* Résumé des filtres actifs */}
          {activeFiltersCount > 0 && (
            <Box sx={{ mt: 2, p: 2, bgcolor: 'primary.light', borderRadius: 1 }}>
              <Typography variant="body2" color="primary.dark">
                {activeFiltersCount} filtre(s) actif(s)
              </Typography>
            </Box>
          )}
        </CardContent>
      </Card>
    </LocalizationProvider>
  );
};