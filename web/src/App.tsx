import React from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { CssBaseline, Box, CircularProgress } from '@mui/material';
import { AuthProvider, useAuth } from './context/AuthContext';
import { AuthPage } from './components/auth/AuthPage';
import MainDashboard from './components/dashboard/MainDashboard';

// Thème Material-UI
const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

// Composant principal de l'application (après authentification)
function AppContent() {
  const { state } = useAuth();

  // Écran de chargement pendant la vérification de l'authentification
  if (state.isLoading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="100vh"
      >
        <CircularProgress />
      </Box>
    );
  }

  // Affichage conditionnel selon l'état d'authentification
  return state.isAuthenticated ? <MainDashboard /> : <AuthPage />;
}

// Composant racine de l'application
export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </ThemeProvider>
  );
}
