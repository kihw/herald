import React, { useState } from 'react';
import {
  Box,
  Paper,
  TextField,
  Button,
  Typography,
  Alert,
  CircularProgress,
  Divider,
  InputAdornment,
} from '@mui/material';
import { useAuth } from '../../context/AuthContext';

interface RegisterFormProps {
  onSwitchToLogin: () => void;
}

export function RegisterForm({ onSwitchToLogin }: RegisterFormProps) {
  const { state, register, clearError } = useAuth();
  const [formData, setFormData] = useState({
    username: '',
    tagline: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // Effacer les erreurs
    if (state.error) clearError();
    if (formErrors[name]) {
      setFormErrors(prev => ({ ...prev, [name]: '' }));
    }
  };

  const validateForm = () => {
    const errors: Record<string, string> = {};

    if (!formData.username || formData.username.length < 3) {
      errors.username = 'Le nom d\'utilisateur doit faire au moins 3 caractères';
    }

    if (!formData.tagline || formData.tagline.length < 2) {
      errors.tagline = 'Le tagline doit faire au moins 2 caractères';
    }

    if (!formData.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = 'Email invalide';
    }

    if (!formData.password || formData.password.length < 6) {
      errors.password = 'Le mot de passe doit faire au moins 6 caractères';
    }

    if (formData.password !== formData.confirmPassword) {
      errors.confirmPassword = 'Les mots de passe ne correspondent pas';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    await register(formData.username, formData.tagline, formData.email, formData.password);
  };

  return (
    <Paper elevation={3} sx={{ p: 4, maxWidth: 400, mx: 'auto', mt: 8 }}>
      <Typography variant="h4" component="h1" gutterBottom align="center">
        Inscription
      </Typography>
      
      {state.error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {state.error}
        </Alert>
      )}

      <Box component="form" onSubmit={handleSubmit} noValidate>
        <TextField
          margin="normal"
          required
          fullWidth
          id="username"
          label="Nom d'utilisateur Riot"
          name="username"
          autoComplete="username"
          autoFocus
          value={formData.username}
          onChange={handleChange}
          disabled={state.isLoading}
          error={!!formErrors.username}
          helperText={formErrors.username}
          placeholder="Ex: Faker"
        />
        
        <TextField
          margin="normal"
          required
          fullWidth
          id="tagline"
          label="Tagline"
          name="tagline"
          value={formData.tagline}
          onChange={handleChange}
          disabled={state.isLoading}
          error={!!formErrors.tagline}
          helperText={formErrors.tagline}
          placeholder="Ex: EUW"
          InputProps={{
            startAdornment: <InputAdornment position="start">#</InputAdornment>,
          }}
        />

        <TextField
          margin="normal"
          required
          fullWidth
          id="email"
          label="Email"
          name="email"
          autoComplete="email"
          value={formData.email}
          onChange={handleChange}
          disabled={state.isLoading}
          error={!!formErrors.email}
          helperText={formErrors.email}
        />

        <TextField
          margin="normal"
          required
          fullWidth
          name="password"
          label="Mot de passe"
          type="password"
          id="password"
          autoComplete="new-password"
          value={formData.password}
          onChange={handleChange}
          disabled={state.isLoading}
          error={!!formErrors.password}
          helperText={formErrors.password}
        />

        <TextField
          margin="normal"
          required
          fullWidth
          name="confirmPassword"
          label="Confirmer le mot de passe"
          type="password"
          id="confirmPassword"
          value={formData.confirmPassword}
          onChange={handleChange}
          disabled={state.isLoading}
          error={!!formErrors.confirmPassword}
          helperText={formErrors.confirmPassword}
        />
        
        <Button
          type="submit"
          fullWidth
          variant="contained"
          sx={{ mt: 3, mb: 2 }}
          disabled={state.isLoading}
        >
          {state.isLoading ? <CircularProgress size={24} /> : 'S\'inscrire'}
        </Button>

        <Divider sx={{ my: 2 }} />

        <Box textAlign="center">
          <Typography variant="body2">
            Déjà un compte ?{' '}
            <Button
              variant="text"
              size="small"
              onClick={onSwitchToLogin}
              disabled={state.isLoading}
            >
              Se connecter
            </Button>
          </Typography>
        </Box>
      </Box>
    </Paper>
  );
}
