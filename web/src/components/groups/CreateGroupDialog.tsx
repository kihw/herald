import React, { useState } from 'react';
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
} from '@mui/material';
import {
  Close as CloseIcon,
  Public as PublicIcon,
  Lock as LockIcon,
  GroupAdd as GroupAddIcon,
  Info as InfoIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';

interface CreateGroupDialogProps {
  open: boolean;
  onClose: () => void;
  onCreate: (groupData: {
    name: string;
    description: string;
    privacy: 'public' | 'private' | 'invite_only';
  }) => Promise<void>;
}

const CreateGroupDialog: React.FC<CreateGroupDialogProps> = ({ open, onClose, onCreate }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    privacy: 'private' as 'public' | 'private' | 'invite_only',
  });
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name.trim()) {
      setError('Le nom du groupe est requis');
      return;
    }

    if (formData.name.length < 3) {
      setError('Le nom du groupe doit contenir au moins 3 caractères');
      return;
    }

    if (formData.name.length > 50) {
      setError('Le nom du groupe ne peut pas dépasser 50 caractères');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      await onCreate({
        name: formData.name.trim(),
        description: formData.description.trim(),
        privacy: formData.privacy,
      });
      
      // Reset form
      setFormData({
        name: '',
        description: '',
        privacy: 'private',
      });
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur lors de la création du groupe');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      setFormData({
        name: '',
        description: '',
        privacy: 'private',
      });
      setError(null);
      onClose();
    }
  };

  const privacyOptions = [
    {
      value: 'private',
      label: 'Privé',
      description: 'Seuls les membres invités peuvent rejoindre',
      icon: <LockIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.gold[500],
    },
    {
      value: 'invite_only',
      label: 'Sur invitation',
      description: 'Visible mais nécessite une invitation pour rejoindre',
      icon: <GroupAddIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.dark[400],
    },
    {
      value: 'public',
      label: 'Public',
      description: 'Tout le monde peut rechercher et rejoindre ce groupe',
      icon: <PublicIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.blue[500],
    },
  ];

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: 3,
          border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          background: isDarkMode
            ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
            : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
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
          Créer un Nouveau Groupe
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

          {/* Group Name */}
          <TextField
            fullWidth
            label="Nom du groupe"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="Ex: Les Pros de la Rift"
            required
            disabled={loading}
            sx={{ mb: 3 }}
            helperText={`${formData.name.length}/50 caractères`}
            inputProps={{ maxLength: 50 }}
          />

          {/* Description */}
          <TextField
            fullWidth
            label="Description (optionnel)"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Décrivez votre groupe et ses objectifs..."
            multiline
            rows={3}
            disabled={loading}
            sx={{ mb: 3 }}
            helperText={`${formData.description.length}/500 caractères`}
            inputProps={{ maxLength: 500 }}
          />

          {/* Privacy Settings */}
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
              Confidentialité du groupe
            </FormLabel>
            
            <RadioGroup
              value={formData.privacy}
              onChange={(e) => setFormData({ ...formData, privacy: e.target.value as any })}
              disabled={loading}
            >
              {privacyOptions.map((option) => (
                <Box key={option.value} sx={{ mb: 2 }}>
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
                      '& .MuiFormControlLabel-label': {
                        mt: 0.5,
                      },
                    }}
                  />
                </Box>
              ))}
            </RadioGroup>
          </FormControl>

          {/* Info Box */}
          <Box
            sx={{
              mt: 3,
              p: 2,
              borderRadius: 2,
              background: `${leagueColors.blue[500]}10`,
              border: `1px solid ${leagueColors.blue[500]}30`,
              display: 'flex',
              gap: 1,
            }}
          >
            <InfoIcon sx={{ color: leagueColors.blue[500], fontSize: 20, mt: 0.5 }} />
            <Box>
              <Typography variant="body2" sx={{ fontWeight: 500, mb: 0.5 }}>
                À propos des groupes
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Les groupes vous permettent de comparer vos performances League of Legends avec vos amis, 
                de créer des comparaisons personnalisées et de suivre vos progrès ensemble.
              </Typography>
            </Box>
          </Box>
        </DialogContent>

        <DialogActions sx={{ px: 3, pb: 3 }}>
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
            disabled={loading || !formData.name.trim()}
            startIcon={loading ? <CircularProgress size={16} /> : null}
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
            {loading ? 'Création...' : 'Créer le Groupe'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default CreateGroupDialog;