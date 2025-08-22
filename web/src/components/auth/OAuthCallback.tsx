import React, { useEffect, useState } from 'react';
import { Box, CircularProgress, Typography, Alert } from '@mui/material';

interface OAuthCallbackProps {
  onAuthSuccess?: (user: any) => void;
  onAuthError?: (error: string) => void;
}

export const OAuthCallback: React.FC<OAuthCallbackProps> = ({ 
  onAuthSuccess, 
  onAuthError 
}) => {
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');
  const [userInfo, setUserInfo] = useState<any>(null);

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    
    // Vérifier les paramètres OAuth dans l'URL
    const oauthSuccess = urlParams.get('oauth_success');
    const oauthError = urlParams.get('oauth_error');
    
    if (oauthSuccess) {
      // Succès OAuth
      const user = urlParams.get('user');
      const email = urlParams.get('email');
      const picture = urlParams.get('picture');
      
      const userData = {
        name: user,
        email: email,
        picture: picture,
        authenticated: true,
      };
      
      setUserInfo(userData);
      setStatus('success');
      setMessage(`Bienvenue, ${user}!`);
      
      // Stocker dans localStorage pour persistence
      localStorage.setItem('herald_user', JSON.stringify(userData));
      
      // Callback de succès
      onAuthSuccess?.(userData);
      
      // Nettoyer l'URL
      window.history.replaceState({}, document.title, window.location.pathname);
      
    } else if (oauthError) {
      // Erreur OAuth
      setStatus('error');
      
      const errorMessages: { [key: string]: string } = {
        'access_denied': 'Accès refusé par l\'utilisateur',
        'invalid_request': 'Requête invalide',
        'invalid_scope': 'Portée invalide',
        'server_error': 'Erreur du serveur',
        'temporarily_unavailable': 'Service temporairement indisponible',
        'missing_code': 'Code d\'autorisation manquant',
        'invalid_state': 'Token de sécurité invalide',
        'token_exchange_failed': 'Échec de l\'échange de token',
        'userinfo_failed': 'Échec de récupération des informations utilisateur',
      };
      
      const errorMsg = errorMessages[oauthError] || `Erreur OAuth: ${oauthError}`;
      setMessage(errorMsg);
      
      // Callback d'erreur
      onAuthError?.(errorMsg);
      
      // Nettoyer l'URL
      window.history.replaceState({}, document.title, window.location.pathname);
      
    } else {
      // Pas de paramètres OAuth, vérifier localStorage
      const storedUser = localStorage.getItem('herald_user');
      if (storedUser) {
        try {
          const userData = JSON.parse(storedUser);
          setUserInfo(userData);
          setStatus('success');
          setMessage(`Déjà connecté en tant que ${userData.name}`);
          onAuthSuccess?.(userData);
        } catch (e) {
          localStorage.removeItem('herald_user');
          setStatus('loading');
        }
      } else {
        setStatus('loading');
      }
    }
  }, [onAuthSuccess, onAuthError]);

  if (status === 'loading') {
    return (
      <Box sx={{ 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'center',
        minHeight: '200px',
        flexDirection: 'column',
        gap: 2 
      }}>
        <CircularProgress />
        <Typography variant="body1">
          Traitement de l'authentification...
        </Typography>
      </Box>
    );
  }

  if (status === 'error') {
    return (
      <Box sx={{ p: 2 }}>
        <Alert severity="error">
          <Typography variant="h6">Erreur d'authentification</Typography>
          {message}
        </Alert>
      </Box>
    );
  }

  if (status === 'success' && userInfo) {
    return (
      <Box sx={{ p: 2 }}>
        <Alert severity="success">
          <Typography variant="h6">Authentification réussie</Typography>
          {message}
          {userInfo.email && (
            <Typography variant="body2" sx={{ mt: 1 }}>
              Email: {userInfo.email}
            </Typography>
          )}
        </Alert>
      </Box>
    );
  }

  return null;
};

export default OAuthCallback;