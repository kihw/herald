import React, { useState } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline, Box, CircularProgress } from '@mui/material';
import { AuthProvider, useAuth } from './context/AuthContext';
import { AuthPage } from './components/auth/AuthPage';
import MainDashboard from './components/dashboard/MainDashboard';
import { lightTheme, darkTheme } from './theme/leagueTheme';

// Context pour le thème
const ThemeContext = React.createContext({
  isDarkMode: false,
  toggleTheme: () => {},
});

export const useTheme = () => React.useContext(ThemeContext);

// Composant principal de l'application (après authentification)
function AppContent() {
  const { user, isAuthenticated, isLoading } = useAuth();

  // Écran de chargement pendant la vérification de l'authentification
  if (isLoading) {
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
  return isAuthenticated ? <MainDashboard /> : <AuthPage />;
}

// Composant racine de l'application
export default function App() {
  const [isDarkMode, setIsDarkMode] = useState(
    localStorage.getItem('herald-theme') === 'dark'
  );

  const toggleTheme = () => {
    const newMode = !isDarkMode;
    setIsDarkMode(newMode);
    localStorage.setItem('herald-theme', newMode ? 'dark' : 'light');
  };

  const currentTheme = isDarkMode ? darkTheme : lightTheme;

  return (
    <ThemeContext.Provider value={{ isDarkMode, toggleTheme }}>
      <ThemeProvider theme={currentTheme}>
        <CssBaseline />
        <AuthProvider>
          <AppContent />
        </AuthProvider>
      </ThemeProvider>
    </ThemeContext.Provider>
  );
}
