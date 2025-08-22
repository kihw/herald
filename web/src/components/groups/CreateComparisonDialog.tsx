import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  Box,
  Typography,
  IconButton,
  useTheme,
  Alert,
  CircularProgress,
  Chip,
  Autocomplete,
  Slider,
  Grid,
  Card,
  CardContent,
  Divider,
} from '@mui/material';
import {
  Close as CloseIcon,
  EmojiEvents as EmojiEventsIcon,
  SportsEsports as SportsEsportsIcon,
  TrendingUp as TrendingUpIcon,
  Analytics as AnalyticsIcon,
  Person as PersonIcon,
  Group as GroupIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import { groupApi, GroupMember } from '../../services/groupApi';

interface CreateComparisonDialogProps {
  open: boolean;
  onClose: () => void;
  onCreateComparison: (comparisonData: any) => Promise<void>;
  groupId: number;
}

const CreateComparisonDialog: React.FC<CreateComparisonDialogProps> = ({ 
  open, 
  onClose, 
  onCreateComparison,
  groupId 
}) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    compare_type: 'champions' as 'champions' | 'roles' | 'performance' | 'trends',
    parameters: {
      member_ids: [] as number[],
      time_range: '30d',
      champions: [] as number[],
      roles: [] as string[],
      game_modes: [] as number[],
      metrics: ['winrate', 'kda', 'cs'] as string[],
      min_games: 5,
    },
  });
  
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (open) {
      loadGroupMembers();
    }
  }, [open, groupId]);

  const loadGroupMembers = async () => {
    try {
      const groupMembers = await groupApi.getGroupMembers(groupId);
      setMembers(groupMembers);
    } catch (err) {
      setError('Erreur lors du chargement des membres');
      console.error('Error loading members:', err);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name.trim()) {
      setError('Le nom de la comparaison est requis');
      return;
    }

    if (formData.parameters.member_ids.length < 2) {
      setError('Sélectionnez au moins 2 membres pour la comparaison');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      await onCreateComparison({
        name: formData.name.trim(),
        description: formData.description.trim(),
        compare_type: formData.compare_type,
        parameters: formData.parameters,
      });
      
      // Reset form
      setFormData({
        name: '',
        description: '',
        compare_type: 'champions',
        parameters: {
          member_ids: [],
          time_range: '30d',
          champions: [],
          roles: [],
          game_modes: [],
          metrics: ['winrate', 'kda', 'cs'],
          min_games: 5,
        },
      });
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur lors de la création de la comparaison');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      setFormData({
        name: '',
        description: '',
        compare_type: 'champions',
        parameters: {
          member_ids: [],
          time_range: '30d',
          champions: [],
          roles: [],
          game_modes: [],
          metrics: ['winrate', 'kda', 'cs'],
          min_games: 5,
        },
      });
      setError(null);
      onClose();
    }
  };

  const comparisonTypes = [
    {
      value: 'champions',
      label: 'Champions',
      description: 'Comparer les performances sur les champions favoris',
      icon: <EmojiEventsIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.gold[500],
    },
    {
      value: 'roles',
      label: 'Rôles',
      description: 'Analyser les performances par rôle (ADC, Support, etc.)',
      icon: <SportsEsportsIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.blue[500],
    },
    {
      value: 'performance',
      label: 'Performance',
      description: 'Comparer les métriques globales (KDA, CS, vision)',
      icon: <TrendingUpIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.win,
    },
    {
      value: 'trends',
      label: 'Tendances',
      description: 'Évolution des performances dans le temps',
      icon: <AnalyticsIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.dark[400],
    },
  ];

  const timeRanges = [
    { value: '7d', label: '7 derniers jours' },
    { value: '30d', label: '30 derniers jours' },
    { value: '90d', label: '3 derniers mois' },
    { value: 'season', label: 'Saison actuelle' },
  ];

  const availableMetrics = [
    { value: 'winrate', label: 'Taux de victoire' },
    { value: 'kda', label: 'KDA moyen' },
    { value: 'cs', label: 'CS par minute' },
    { value: 'vision', label: 'Score de vision' },
    { value: 'damage', label: 'Dégâts par minute' },
    { value: 'gold', label: 'Or par minute' },
  ];

  const availableRoles = [
    { value: 'TOP', label: 'Top Lane' },
    { value: 'JUNGLE', label: 'Jungle' },
    { value: 'MIDDLE', label: 'Mid Lane' },
    { value: 'BOTTOM', label: 'ADC' },
    { value: 'UTILITY', label: 'Support' },
  ];

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="md"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: 3,
          border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          background: isDarkMode
            ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
            : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
          minHeight: 600,
        },
      }}
    >
      <DialogTitle
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          pb: 1,
          background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
          color: '#fff',
          m: -3,
          mb: 3,
          px: 3,
          py: 2,
        }}
      >
        <Typography variant="h6" sx={{ fontWeight: 600 }}>
          Créer une Nouvelle Comparaison
        </Typography>
        <IconButton
          onClick={handleClose}
          disabled={loading}
          sx={{ color: '#fff' }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <form onSubmit={handleSubmit}>
        <DialogContent sx={{ pt: 0 }}>
          {/* Error Alert */}
          {error && (
            <Alert 
              severity="error" 
              sx={{ mb: 3 }}
              onClose={() => setError(null)}
            >
              {error}
            </Alert>
          )}

          <Grid container spacing={3}>
            {/* Basic Information */}
            <Grid item xs={12}>
              <Card sx={{ border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}` }}>
                <CardContent>
                  <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                    Informations de base
                  </Typography>

                  <TextField
                    fullWidth
                    label="Nom de la comparaison"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="Ex: Comparaison Champions S14"
                    required
                    disabled={loading}
                    sx={{ mb: 3 }}
                    helperText={`${formData.name.length}/100 caractères`}
                    inputProps={{ maxLength: 100 }}
                  />

                  <TextField
                    fullWidth
                    label="Description (optionnel)"
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    placeholder="Décrivez l'objectif de cette comparaison..."
                    multiline
                    rows={2}
                    disabled={loading}
                    sx={{ mb: 2 }}
                    helperText={`${formData.description.length}/300 caractères`}
                    inputProps={{ maxLength: 300 }}
                  />
                </CardContent>
              </Card>
            </Grid>

            {/* Comparison Type */}
            <Grid item xs={12}>
              <Card sx={{ border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}` }}>
                <CardContent>
                  <FormControl component="fieldset" sx={{ width: '100%' }}>
                    <FormLabel 
                      component="legend" 
                      sx={{ 
                        mb: 2, 
                        fontWeight: 600,
                        color: 'text.primary',
                        '&.Mui-focused': { color: 'text.primary' },
                      }}
                    >
                      Type de comparaison
                    </FormLabel>
                    
                    <RadioGroup
                      value={formData.compare_type}
                      onChange={(e) => setFormData({ 
                        ...formData, 
                        compare_type: e.target.value as any
                      })}
                      disabled={loading}
                    >
                      <Grid container spacing={2}>
                        {comparisonTypes.map((option) => (
                          <Grid item xs={12} sm={6} key={option.value}>
                            <FormControlLabel
                              value={option.value}
                              control={<Radio sx={{ color: option.color }} />}
                              label={
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                  <Box sx={{ color: option.color }}>
                                    {option.icon}
                                  </Box>
                                  <Box>
                                    <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                      {option.label}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                      {option.description}
                                    </Typography>
                                  </Box>
                                </Box>
                              }
                              sx={{
                                alignItems: 'flex-start',
                                border: `1px solid ${formData.compare_type === option.value ? option.color : 'transparent'}`,
                                borderRadius: 2,
                                p: 1,
                                m: 0,
                                width: '100%',
                                '& .MuiFormControlLabel-label': {
                                  mt: 0.5,
                                  width: '100%',
                                },
                              }}
                            />
                          </Grid>
                        ))}
                      </Grid>
                    </RadioGroup>
                  </FormControl>
                </CardContent>
              </Card>
            </Grid>

            {/* Members Selection */}
            <Grid item xs={12}>
              <Card sx={{ border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}` }}>
                <CardContent>
                  <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                    Sélection des membres
                  </Typography>

                  <Autocomplete
                    multiple
                    options={members}
                    getOptionLabel={(member) => `${member.user.riot_id}#${member.user.riot_tag}`}
                    value={members.filter(member => formData.parameters.member_ids.includes(member.user.id))}
                    onChange={(event, newValue) => {
                      setFormData({
                        ...formData,
                        parameters: {
                          ...formData.parameters,
                          member_ids: newValue.map(member => member.user.id),
                        },
                      });
                    }}
                    renderTags={(value, getTagProps) =>
                      value.map((option, index) => (
                        <Chip
                          variant="outlined"
                          label={`${option.user.riot_id}#${option.user.riot_tag}`}
                          {...getTagProps({ index })}
                          icon={<PersonIcon />}
                        />
                      ))
                    }
                    renderOption={(props, option) => (
                      <li {...props}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <PersonIcon />
                          <Box>
                            <Typography variant="body2">
                              {option.user.riot_id}#{option.user.riot_tag}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                              {option.user.region?.toUpperCase()} • {option.role}
                            </Typography>
                          </Box>
                        </Box>
                      </li>
                    )}
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        label="Membres à comparer"
                        placeholder="Sélectionnez au moins 2 membres"
                        helperText={`${formData.parameters.member_ids.length} membre(s) sélectionné(s)`}
                      />
                    )}
                    disabled={loading}
                  />
                </CardContent>
              </Card>
            </Grid>

            {/* Parameters */}
            <Grid item xs={12}>
              <Card sx={{ border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}` }}>
                <CardContent>
                  <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                    Paramètres de comparaison
                  </Typography>

                  <Grid container spacing={3}>
                    {/* Time Range */}
                    <Grid item xs={12} sm={6}>
                      <FormControl fullWidth>
                        <TextField
                          select
                          label="Période"
                          value={formData.parameters.time_range}
                          onChange={(e) => setFormData({
                            ...formData,
                            parameters: {
                              ...formData.parameters,
                              time_range: e.target.value,
                            },
                          })}
                          SelectProps={{ native: true }}
                          disabled={loading}
                        >
                          {timeRanges.map((option) => (
                            <option key={option.value} value={option.value}>
                              {option.label}
                            </option>
                          ))}
                        </TextField>
                      </FormControl>
                    </Grid>

                    {/* Minimum Games */}
                    <Grid item xs={12} sm={6}>
                      <Typography gutterBottom>
                        Minimum de parties: {formData.parameters.min_games}
                      </Typography>
                      <Slider
                        value={formData.parameters.min_games}
                        onChange={(e, newValue) => setFormData({
                          ...formData,
                          parameters: {
                            ...formData.parameters,
                            min_games: newValue as number,
                          },
                        })}
                        min={1}
                        max={50}
                        step={1}
                        marks={[
                          { value: 1, label: '1' },
                          { value: 25, label: '25' },
                          { value: 50, label: '50' },
                        ]}
                        disabled={loading}
                      />
                    </Grid>

                    {/* Metrics */}
                    <Grid item xs={12}>
                      <Autocomplete
                        multiple
                        options={availableMetrics}
                        getOptionLabel={(option) => option.label}
                        value={availableMetrics.filter(metric => formData.parameters.metrics.includes(metric.value))}
                        onChange={(event, newValue) => {
                          setFormData({
                            ...formData,
                            parameters: {
                              ...formData.parameters,
                              metrics: newValue.map(metric => metric.value),
                            },
                          });
                        }}
                        renderTags={(value, getTagProps) =>
                          value.map((option, index) => (
                            <Chip
                              variant="outlined"
                              label={option.label}
                              {...getTagProps({ index })}
                            />
                          ))
                        }
                        renderInput={(params) => (
                          <TextField
                            {...params}
                            label="Métriques à comparer"
                            placeholder="Sélectionnez les statistiques"
                          />
                        )}
                        disabled={loading}
                      />
                    </Grid>

                    {/* Role-specific options */}
                    {formData.compare_type === 'roles' && (
                      <Grid item xs={12}>
                        <Autocomplete
                          multiple
                          options={availableRoles}
                          getOptionLabel={(option) => option.label}
                          value={availableRoles.filter(role => formData.parameters.roles.includes(role.value))}
                          onChange={(event, newValue) => {
                            setFormData({
                              ...formData,
                              parameters: {
                                ...formData.parameters,
                                roles: newValue.map(role => role.value),
                              },
                            });
                          }}
                          renderTags={(value, getTagProps) =>
                            value.map((option, index) => (
                              <Chip
                                variant="outlined"
                                label={option.label}
                                {...getTagProps({ index })}
                                icon={<SportsEsportsIcon />}
                              />
                            ))
                          }
                          renderInput={(params) => (
                            <TextField
                              {...params}
                              label="Rôles à analyser"
                              placeholder="Tous les rôles si vide"
                            />
                          )}
                          disabled={loading}
                        />
                      </Grid>
                    )}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </DialogContent>

        <Divider />

        <DialogActions sx={{ px: 3, py: 2 }}>
          <Button
            onClick={handleClose}
            disabled={loading}
            sx={{
              color: 'text.secondary',
              '&:hover': {
                backgroundColor: `${leagueColors.dark[400]}10`,
              },
            }}
          >
            Annuler
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading || !formData.name.trim() || formData.parameters.member_ids.length < 2}
            startIcon={loading ? <CircularProgress size={16} /> : <AnalyticsIcon />}
            sx={{
              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
              '&:hover': {
                background: `linear-gradient(135deg, ${leagueColors.blue[600]} 0%, ${leagueColors.blue[700]} 100%)`,
              },
              '&:disabled': {
                background: 'rgba(0,0,0,0.12)',
              },
            }}
          >
            {loading ? 'Création...' : 'Créer la Comparaison'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default CreateComparisonDialog;