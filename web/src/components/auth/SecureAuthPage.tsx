import React, { useState, useEffect } from 'react';
import {
  Container,
  Box,
  Typography,
  TextField,
  Button,
  Alert,
  CircularProgress,
  MenuItem,
  Paper,
  Divider,
  Chip,
} from '@mui/material';
import {
  Security,
  SportsEsports,
  VerifiedUser,
  Shield,
  Lock,
} from '@mui/icons-material';
import { useSecureAuth } from '../../context/SecureAuthContext';
import { secureAuthService } from '../../auth/SecureAuthService';

interface Region {
  code: string;
  name: string;
}

export default function SecureAuthPage() {
  const { validateAccount, isLoading, error, clearError } = useSecureAuth();
  const [formData, setFormData] = useState({
    riotId: '',
    riotTag: '',
    region: '',
  });
  const [regions, setRegions] = useState<Region[]>([]);
  const [loadingRegions, setLoadingRegions] = useState(true);
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  // Charger les régions de manière sécurisée
  useEffect(() => {
    const loadRegions = async () => {
      try {
        console.log('🔒 Chargement sécurisé des régions...');
        const response = await secureAuthService.getSupportedRegions();
        
        if (response.regions && Array.isArray(response.regions)) {
          setRegions(response.regions);
          if (response.default && response.regions.some(r => r.code === response.default)) {
            setFormData(prev => ({ ...prev, region: response.default }));
          }
        }
        console.log('✅ Régions chargées:', response.regions.length);
      } catch (error) {
        console.error('❌ Erreur chargement régions:', error);
        // Utiliser des régions par défaut en cas d'erreur
        const defaultRegions = [
          { code: 'euw1', name: 'Europe West' },
          { code: 'na1', name: 'North America' },
          { code: 'eun1', name: 'Europe Nordic & East' },
          { code: 'kr', name: 'Korea' },
        ];
        setRegions(defaultRegions);
        setFormData(prev => ({ ...prev, region: 'euw1' }));
      } finally {
        setLoadingRegions(false);
      }
    };

    loadRegions();
  }, []);

  // Validation côté client sécurisée
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};

    // Validation Riot ID
    if (!formData.riotId.trim()) {
      errors.riotId = 'Riot ID requis';
    } else if (formData.riotId.length < 3 || formData.riotId.length > 16) {
      errors.riotId = 'Riot ID doit contenir 3-16 caractères';
    } else if (!/^[a-zA-Z0-9\\s]+$/.test(formData.riotId)) {
      errors.riotId = 'Riot ID contient des caractères invalides';
    }

    // Validation Riot Tag
    if (!formData.riotTag.trim()) {
      errors.riotTag = 'Riot Tag requis';
    } else if (formData.riotTag.length < 3 || formData.riotTag.length > 5) {
      errors.riotTag = 'Riot Tag doit contenir 3-5 caractères';
    } else if (!/^[a-zA-Z0-9]+$/.test(formData.riotTag)) {
      errors.riotTag = 'Riot Tag ne peut contenir que des lettres et chiffres';
    }

    // Validation région
    if (!formData.region) {
      errors.region = 'Région requise';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Gestion des changements de formulaire
  const handleInputChange = (field: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({
      ...prev,
      [field]: event.target.value,
    }));

    // Nettoyer les erreurs lors de la saisie
    if (formErrors[field]) {
      setFormErrors(prev => ({ ...prev, [field]: '' }));
    }
    if (error) {
      clearError();
    }
  };

  // Soumission sécurisée du formulaire
  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      console.log('🔒 Validation sécurisée du compte...', {
        riotId: formData.riotId,
        riotTag: formData.riotTag,
        region: formData.region,
      });

      await validateAccount(
        formData.riotId.trim(),
        formData.riotTag.trim(),
        formData.region
      );
    } catch (error) {
      console.error('❌ Erreur validation:', error);
    }
  };

  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          py: 4,
        }}
      >
        {/* En-tête sécurisé */}
        <Paper
          elevation={3}
          sx={{
            p: 4,
            mb: 3,
            background: 'linear-gradient(135deg, #1976d2 0%, #1565c0 100%)',
            color: 'white',
            textAlign: 'center',
          }}
        >
          <SportsEsports sx={{ fontSize: 60, mb: 2 }} />
          <Typography variant="h3" component="h1" gutterBottom sx={{ fontWeight: 700 }}>
            🔒 Herald.lol
          </Typography>
          <Typography variant="h6" sx={{ opacity: 0.9 }}>
            Authentification Sécurisée
          </Typography>
          
          {/* Indicateurs de sécurité */}
          <Box sx={{ mt: 2, display: 'flex', justifyContent: 'center', gap: 1, flexWrap: 'wrap' }}>
            <Chip
              icon={<Shield />}
              label="JWT Sécurisé"
              size="small"
              sx={{ backgroundColor: 'rgba(255,255,255,0.2)', color: 'white' }}
            />
            <Chip
              icon={<Lock />}
              label="Protection CSRF"
              size="small"
              sx={{ backgroundColor: 'rgba(255,255,255,0.2)', color: 'white' }}
            />
            <Chip
              icon={<VerifiedUser />}
              label="Cryptage AES"
              size="small"
              sx={{ backgroundColor: 'rgba(255,255,255,0.2)', color: 'white' }}
            />
          </Box>
        </Paper>

        {/* Formulaire d'authentification sécurisé */}
        <Paper elevation={2} sx={{ p: 4 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <Security sx={{ mr: 1, color: 'primary.main' }} />
            <Typography variant="h5" sx={{ fontWeight: 600 }}>
              Validation Riot Games
            </Typography>
          </Box>

          <form onSubmit={handleSubmit}>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
              {/* Riot ID */}
              <TextField
                fullWidth
                label="Riot ID"
                placeholder="VotreNomDeJoueur"
                value={formData.riotId}
                onChange={handleInputChange('riotId')}
                disabled={isLoading}
                required
                error={!!formErrors.riotId}
                helperText={formErrors.riotId || "Votre nom d'invocateur (ex: Faker, Canna)"}
              />

              {/* Riot Tag */}
              <TextField
                fullWidth
                label="Riot Tag"
                placeholder="EUW1"
                value={formData.riotTag}
                onChange={handleInputChange('riotTag')}
                disabled={isLoading}
                required
                error={!!formErrors.riotTag}
                helperText={formErrors.riotTag || 'Votre tag sans le # (ex: EUW1, NA1, 1234)'}
              />

              {/* Région */}
              <TextField
                fullWidth
                select
                label="Région"
                value={formData.region}
                onChange={handleInputChange('region')}
                disabled={isLoading || loadingRegions}
                required
                error={!!formErrors.region}
                helperText={formErrors.region || 'Sélectionnez votre région de jeu'}
              >
                {loadingRegions ? (
                  <MenuItem value="">
                    <CircularProgress size={20} sx={{ mr: 1 }} />
                    Chargement sécurisé...
                  </MenuItem>
                ) : (
                  regions.map((region) => (
                    <MenuItem key={region.code} value={region.code}>
                      {region.name}
                    </MenuItem>
                  ))
                )}
              </TextField>

              {/* Messages d'erreur */}
              {error && (
                <Alert
                  severity="error"
                  onClose={clearError}
                  sx={{ borderRadius: 2 }}
                >
                  {error}
                </Alert>
              )}

              {/* Bouton de validation sécurisé */}
              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={
                  isLoading ||
                  !formData.riotId.trim() ||
                  !formData.riotTag.trim() ||
                  !formData.region ||
                  loadingRegions
                }
                sx={{
                  py: 1.5,
                  fontWeight: 600,
                  fontSize: '1.1rem',
                  background: 'linear-gradient(135deg, #1976d2 0%, #1565c0 100%)',
                  '&:hover': {
                    background: 'linear-gradient(135deg, #1565c0 0%, #0d47a1 100%)',
                    transform: 'translateY(-1px)',
                    boxShadow: '0 6px 20px rgba(25, 118, 210, 0.3)',
                  },
                  transition: 'all 0.3s ease',
                }}
              >
                {isLoading ? (
                  <>
                    <CircularProgress size={20} sx={{ mr: 1, color: 'inherit' }} />
                    Validation sécurisée...
                  </>
                ) : (
                  <>
                    <VerifiedUser sx={{ mr: 1 }} />
                    Valider le compte
                  </>
                )}
              </Button>
            </Box>
          </form>

          <Divider sx={{ my: 3 }} />

          {/* Informations de sécurité */}
          <Box sx={{ textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              🔒 Votre connexion est sécurisée par :
            </Typography>
            <Typography variant="caption" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
              • Tokens JWT avec expiration automatique<br/>
              • Protection CSRF contre les attaques<br/>
              • Cryptage AES des données sensibles<br/>
              • Surveillance continue de l'activité<br/>
              • Validation via l'API officielle Riot Games
            </Typography>
          </Box>
        </Paper>

        {/* Pied de page */}
        <Box sx={{ textAlign: 'center', mt: 3 }}>
          <Typography variant="body2" color="text.secondary">
            Exemple: Riot ID "Faker", Tag "T1", Région "Korea"
          </Typography>
        </Box>
      </Box>
    </Container>
  );
}