import { createTheme, ThemeOptions } from '@mui/material/styles';

// Palette League of Legends
const lolColors = {
  primary: {
    main: '#C89B3C', // Or champion (gold)
    light: '#F0E6D2',
    dark: '#A0752B',
    contrastText: '#0F2027',
  },
  secondary: {
    main: '#0596AA', // Bleu clair LoL
    light: '#5BC0DE',
    dark: '#0F2027',
    contrastText: '#F0E6D2',
  },
  success: {
    main: '#0596AA', // Victoire
    light: '#5BC0DE',
    dark: '#0F2027',
  },
  error: {
    main: '#C8AA6E', // Défaite/Erreur
    light: '#F0E6D2',
    dark: '#A0752B',
  },
  warning: {
    main: '#C89B3C',
    light: '#F0E6D2',
    dark: '#A0752B',
  },
  info: {
    main: '#0596AA',
    light: '#5BC0DE',
    dark: '#0F2027',
  },
  background: {
    default: '#0F2027',
    paper: '#1E2328',
  },
  text: {
    primary: '#F0E6D2',
    secondary: '#C9AA71',
  },
  divider: '#3C3C41',
};

const lightColors = {
  primary: {
    main: '#C89B3C',
    light: '#F0E6D2',
    dark: '#A0752B',
    contrastText: '#0F2027',
  },
  secondary: {
    main: '#0596AA',
    light: '#5BC0DE',
    dark: '#0F2027',
    contrastText: '#F0E6D2',
  },
  success: {
    main: '#0596AA',
    light: '#5BC0DE',
    dark: '#0F2027',
  },
  error: {
    main: '#D32F2F',
    light: '#FFCDD2',
    dark: '#B71C1C',
  },
  warning: {
    main: '#FF9800',
    light: '#FFE0B2',
    dark: '#E65100',
  },
  info: {
    main: '#0596AA',
    light: '#5BC0DE',
    dark: '#0F2027',
  },
  background: {
    default: '#F0E6D2',
    paper: '#FFFFFF',
  },
  text: {
    primary: '#0F2027',
    secondary: '#5A5A5A',
  },
  divider: '#E0E0E0',
};

const baseTheme: ThemeOptions = {
  typography: {
    fontFamily: [
      'Spiegel',
      'Roboto',
      '"Segoe UI"',
      '"Helvetica Neue"',
      'Arial',
      'sans-serif',
    ].join(','),
    h1: {
      fontSize: '2.5rem',
      fontWeight: 700,
      letterSpacing: '-0.01562em',
    },
    h2: {
      fontSize: '2rem',
      fontWeight: 600,
      letterSpacing: '-0.00833em',
    },
    h3: {
      fontSize: '1.75rem',
      fontWeight: 600,
      letterSpacing: '0em',
    },
    h4: {
      fontSize: '1.5rem',
      fontWeight: 600,
      letterSpacing: '0.00735em',
    },
    h5: {
      fontSize: '1.25rem',
      fontWeight: 600,
      letterSpacing: '0em',
    },
    h6: {
      fontSize: '1.125rem',
      fontWeight: 600,
      letterSpacing: '0.0075em',
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.5,
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.43,
    },
    button: {
      fontSize: '0.875rem',
      fontWeight: 600,
      textTransform: 'none' as const,
    },
  },
  shape: {
    borderRadius: 8,
  },
  spacing: 8,
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          scrollbarWidth: 'thin',
          '&::-webkit-scrollbar': {
            width: '8px',
            height: '8px',
          },
          '&::-webkit-scrollbar-track': {
            background: 'rgba(0,0,0,0.1)',
          },
          '&::-webkit-scrollbar-thumb': {
            background: 'rgba(0,0,0,0.3)',
            borderRadius: '4px',
            '&:hover': {
              background: 'rgba(0,0,0,0.5)',
            },
          },
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 6,
          padding: '8px 16px',
          fontWeight: 600,
        },
        contained: {
          boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
          '&:hover': {
            boxShadow: '0 4px 8px rgba(0,0,0,0.3)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 6,
          fontWeight: 500,
        },
      },
    },
  },
};

export const createAppTheme = (mode: 'light' | 'dark') => {
  const colors = mode === 'dark' ? lolColors : lightColors;
  
  return createTheme({
    ...baseTheme,
    palette: {
      mode,
      ...colors,
    },
  });
};

export const darkTheme = createAppTheme('dark');
export const lightTheme = createAppTheme('light');