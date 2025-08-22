import { createTheme, ThemeOptions } from '@mui/material/styles';

// Couleurs League of Legends officielles
const leagueColors = {
  // Couleurs principales
  blue: {
    25: '#f0f8ff',
    50: '#e3f2fd',
    100: '#bbdefb',
    200: '#90caf9',
    300: '#64b5f6',
    400: '#42a5f5',
    500: '#1976d2', // Bleu principal LoL
    600: '#1565c0',
    700: '#0d47a1',
    800: '#0a3d91',
    900: '#063281',
  },
  gold: {
    25: '#fffffe',
    50: '#fffef7',
    100: '#fffbeb',
    200: '#fef3c7',
    300: '#fde68a',
    400: '#fcd34d',
    500: '#f59e0b', // Or League of Legends
    600: '#d97706',
    700: '#b45309',
    800: '#92400e',
    900: '#78350f',
  },
  // Couleurs de rang
  iron: '#6d4c41',
  bronze: '#8d6e63',
  silver: '#90a4ae',
  goldRank: '#ffc107',
  platinum: '#00bcd4',
  diamond: '#3f51b5',
  master: '#9c27b0',
  grandmaster: '#e91e63',
  challenger: '#f44336',
  // Couleurs d'état
  win: '#4caf50',
  loss: '#f44336',
  // Couleurs sombres
  dark: {
    50: '#263238',
    100: '#37474f',
    200: '#455a64',
    300: '#546e7a',
    400: '#607d8b',
    500: '#78909c',
    600: '#90a4ae',
    700: '#b0bec5',
    800: '#cfd8dc',
    900: '#eceff1',
  }
};

// Configuration du thème
const createLeagueTheme = (mode: 'light' | 'dark' = 'light'): ThemeOptions => {
  const isDark = mode === 'dark';

  return {
    palette: {
      mode,
      primary: {
        main: leagueColors.blue[500],
        light: leagueColors.blue[400],
        dark: leagueColors.blue[700],
        contrastText: '#ffffff',
      },
      secondary: {
        main: leagueColors.gold[500],
        light: leagueColors.gold[400],
        dark: leagueColors.gold[700],
        contrastText: '#000000',
      },
      success: {
        main: leagueColors.win,
        light: '#81c784',
        dark: '#388e3c',
      },
      error: {
        main: leagueColors.loss,
        light: '#e57373',
        dark: '#d32f2f',
      },
      warning: {
        main: leagueColors.gold[600],
        light: leagueColors.gold[400],
        dark: leagueColors.gold[800],
      },
      info: {
        main: leagueColors.blue[400],
        light: leagueColors.blue[300],
        dark: leagueColors.blue[600],
      },
      background: {
        default: isDark ? '#0a1428' : '#f5f7fa',
        paper: isDark ? '#1e2328' : '#ffffff',
      },
      text: {
        primary: isDark ? '#f0e6d2' : '#1e2328',
        secondary: isDark ? '#c9aa71' : '#5a6b7d',
      },
      divider: isDark ? '#463714' : '#e0e6ed',
    },
    typography: {
      fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
      h1: {
        fontSize: '2.5rem',
        fontWeight: 700,
        lineHeight: 1.2,
        color: isDark ? '#f0e6d2' : '#1e2328',
      },
      h2: {
        fontSize: '2rem',
        fontWeight: 600,
        lineHeight: 1.3,
        color: isDark ? '#f0e6d2' : '#1e2328',
      },
      h3: {
        fontSize: '1.5rem',
        fontWeight: 600,
        lineHeight: 1.4,
        color: isDark ? '#f0e6d2' : '#1e2328',
      },
      h4: {
        fontSize: '1.25rem',
        fontWeight: 600,
        lineHeight: 1.4,
        color: isDark ? '#f0e6d2' : '#1e2328',
      },
      h5: {
        fontSize: '1.125rem',
        fontWeight: 500,
        lineHeight: 1.4,
        color: isDark ? '#f0e6d2' : '#1e2328',
      },
      h6: {
        fontSize: '1rem',
        fontWeight: 500,
        lineHeight: 1.4,
        color: isDark ? '#f0e6d2' : '#1e2328',
      },
      body1: {
        fontSize: '1rem',
        lineHeight: 1.6,
        color: isDark ? '#f0e6d2' : '#1e2328',
      },
      body2: {
        fontSize: '0.875rem',
        lineHeight: 1.5,
        color: isDark ? '#c9aa71' : '#5a6b7d',
      },
      button: {
        fontSize: '0.875rem',
        fontWeight: 500,
        textTransform: 'none' as const,
      },
    },
    shape: {
      borderRadius: 8,
    },
    shadows: isDark ? [
      'none',
      '0px 1px 3px rgba(0, 0, 0, 0.4)',
      '0px 2px 6px rgba(0, 0, 0, 0.4)',
      '0px 3px 8px rgba(0, 0, 0, 0.4)',
      '0px 4px 10px rgba(0, 0, 0, 0.4)',
      '0px 5px 12px rgba(0, 0, 0, 0.4)',
      '0px 6px 14px rgba(0, 0, 0, 0.4)',
      '0px 7px 16px rgba(0, 0, 0, 0.4)',
      '0px 8px 18px rgba(0, 0, 0, 0.4)',
      '0px 9px 20px rgba(0, 0, 0, 0.4)',
      '0px 10px 22px rgba(0, 0, 0, 0.4)',
      '0px 11px 24px rgba(0, 0, 0, 0.4)',
      '0px 12px 26px rgba(0, 0, 0, 0.4)',
      '0px 13px 28px rgba(0, 0, 0, 0.4)',
      '0px 14px 30px rgba(0, 0, 0, 0.4)',
      '0px 15px 32px rgba(0, 0, 0, 0.4)',
      '0px 16px 34px rgba(0, 0, 0, 0.4)',
      '0px 17px 36px rgba(0, 0, 0, 0.4)',
      '0px 18px 38px rgba(0, 0, 0, 0.4)',
      '0px 19px 40px rgba(0, 0, 0, 0.4)',
      '0px 20px 42px rgba(0, 0, 0, 0.4)',
      '0px 21px 44px rgba(0, 0, 0, 0.4)',
      '0px 22px 46px rgba(0, 0, 0, 0.4)',
      '0px 23px 48px rgba(0, 0, 0, 0.4)',
      '0px 24px 50px rgba(0, 0, 0, 0.4)',
    ] : undefined,
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          body: {
            scrollbarWidth: 'thin',
            scrollbarColor: isDark ? '#463714 #1e2328' : '#e0e6ed #f5f7fa',
            '&::-webkit-scrollbar, & *::-webkit-scrollbar': {
              width: '8px',
              height: '8px',
            },
            '&::-webkit-scrollbar-thumb, & *::-webkit-scrollbar-thumb': {
              borderRadius: '4px',
              backgroundColor: isDark ? '#463714' : '#c4c4c4',
              minHeight: '24px',
            },
            '&::-webkit-scrollbar-track, & *::-webkit-scrollbar-track': {
              backgroundColor: isDark ? '#1e2328' : '#f5f7fa',
            },
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: '8px',
            padding: '8px 24px',
            fontWeight: 500,
            textTransform: 'none',
            boxShadow: 'none',
            '&:hover': {
              boxShadow: '0px 2px 8px rgba(0, 0, 0, 0.15)',
            },
          },
          containedPrimary: {
            background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
            '&:hover': {
              background: `linear-gradient(135deg, ${leagueColors.blue[600]} 0%, ${leagueColors.blue[700]} 100%)`,
            },
          },
          containedSecondary: {
            background: `linear-gradient(135deg, ${leagueColors.gold[500]} 0%, ${leagueColors.gold[600]} 100%)`,
            color: '#000000',
            '&:hover': {
              background: `linear-gradient(135deg, ${leagueColors.gold[600]} 0%, ${leagueColors.gold[700]} 100%)`,
            },
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            borderRadius: '12px',
            border: isDark ? `1px solid #463714` : `1px solid #e0e6ed`,
            background: isDark 
              ? 'linear-gradient(135deg, #1e2328 0%, #0a1428 100%)'
              : 'linear-gradient(135deg, #ffffff 0%, #f8fafc 100%)',
            boxShadow: isDark 
              ? '0px 4px 12px rgba(0, 0, 0, 0.3)'
              : '0px 2px 8px rgba(0, 0, 0, 0.1)',
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            borderRadius: '12px',
            backgroundImage: 'none',
          },
        },
      },
      MuiTab: {
        styleOverrides: {
          root: {
            textTransform: 'none',
            minWidth: 0,
            fontWeight: 500,
            color: isDark ? '#c9aa71' : '#5a6b7d',
            '&.Mui-selected': {
              color: leagueColors.blue[500],
            },
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            borderRadius: '6px',
            fontWeight: 500,
          },
          colorPrimary: {
            backgroundColor: leagueColors.blue[100],
            color: leagueColors.blue[800],
          },
          colorSecondary: {
            backgroundColor: leagueColors.gold[100],
            color: leagueColors.gold[800],
          },
        },
      },
    },
  };
};

// Couleurs personnalisées pour les rangs
export const rankColors = {
  IRON: leagueColors.iron,
  BRONZE: leagueColors.bronze,
  SILVER: leagueColors.silver,
  GOLD: leagueColors.goldRank,
  PLATINUM: leagueColors.platinum,
  DIAMOND: leagueColors.diamond,
  MASTER: leagueColors.master,
  GRANDMASTER: leagueColors.grandmaster,
  CHALLENGER: leagueColors.challenger,
};

// Création des thèmes
export const lightTheme = createTheme(createLeagueTheme('light'));
export const darkTheme = createTheme(createLeagueTheme('dark'));

export { leagueColors };
export default lightTheme;