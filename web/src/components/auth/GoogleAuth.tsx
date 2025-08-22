import React, { useState } from 'react';
import { Button, Box, Typography, Alert, useTheme } from '@mui/material';
import { Google as GoogleIcon } from '@mui/icons-material';
import { getApiUrl } from '../../utils/api-config';
import { leagueColors } from '../../theme/leagueTheme';

interface GoogleAuthProps {
  onSuccess?: (user: any) => void;
  onError?: (error: string) => void;
}

export const GoogleAuth: React.FC<GoogleAuthProps> = ({ onSuccess, onError }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGoogleLogin = async () => {
    setLoading(true);
    setError(null);

    try {
      // Initier le flow OAuth Google
      const response = await fetch(getApiUrl('/auth/google/init'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Erreur lors de l\'initialisation OAuth');
      }

      const data = await response.json();
      
      if (data.auth_url) {
        // Rediriger vers Google OAuth
        window.location.href = data.auth_url;
      } else {
        throw new Error('URL d\'authentification non reçue');
      }
    } catch (err: any) {
      const errorMessage = err.message || 'Erreur de connexion Google';
      setError(errorMessage);
      onError?.(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box 
      sx={{ 
        textAlign: 'center', 
        p: 3,
        borderRadius: 3,
        background: isDarkMode
          ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
          : `linear-gradient(135deg, ${leagueColors.blue[50]} 0%, #ffffff 100%)`,
        border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
      }}
    >
      <Typography 
        variant="h5" 
        gutterBottom 
        sx={{ 
          fontWeight: 600,
          background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
          backgroundClip: 'text',
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent',
          mb: 3,
        }}
      >
        Connexion Herald.lol
      </Typography>
      
      {error && (
        <Alert 
          severity="error" 
          sx={{ 
            mb: 3,
            borderRadius: 2,
            border: `1px solid ${leagueColors.loss}`,
          }}
        >
          {error}
        </Alert>
      )}
      
      <Button
        variant="contained"
        size="large"
        startIcon={<GoogleIcon />}
        onClick={handleGoogleLogin}
        disabled={loading}
        sx={{
          background: `linear-gradient(135deg, #4285f4 0%, #3367d6 100%)`,
          '&:hover': {
            background: `linear-gradient(135deg, #3367d6 0%, #2955c7 100%)`,
            transform: 'translateY(-1px)',
            boxShadow: '0 6px 20px rgba(66, 133, 244, 0.3)',
          },
          '&:disabled': {
            background: 'rgba(66, 133, 244, 0.5)',
          },
          textTransform: 'none',
          fontSize: '16px',
          fontWeight: 600,
          py: 1.5,
          px: 4,
          borderRadius: 2,
          transition: 'all 0.3s ease',
          minWidth: 280,
        }}
      >
        {loading ? 'Connexion en cours...' : 'Se connecter avec Google'}
      </Button>
      
      <Typography 
        variant="body2" 
        color="text.secondary" 
        sx={{ 
          mt: 3,
          fontSize: '0.9rem',
          lineHeight: 1.5,
        }}
      >
        Connectez-vous avec votre compte Google pour accéder à vos statistiques League of Legends et débloquer toutes les fonctionnalités d'analyse.
      </Typography>
    </Box>
  );
};

export default GoogleAuth;