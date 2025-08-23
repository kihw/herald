import { createTheme, Theme } from '@mui/material/styles';

// League of Legends inspired color palette
export const colors = {
  // Primary - Blue theme inspired by LoL
  primary: {
    50: '#e3f2fd',
    100: '#bbdefb',
    200: '#90caf9',
    300: '#64b5f6',
    400: '#42a5f5',
    500: '#1976d2', // Main blue
    600: '#1565c0',
    700: '#0d47a1',
    800: '#0a3d91',
    900: '#063281',
  },
  
  // Gold theme inspired by LoL gold/yellow
  secondary: {
    50: '#fefcf3',
    100: '#fef7e0',
    200: '#fdecc8',
    300: '#fbd679',
    400: '#f7c52d',
    500: '#c89b3c', // LoL gold
    600: '#b8860b',
    700: '#9a6f0a',
    800: '#805c08',
    900: '#6b4a06',
  },
  
  // Dark theme colors
  dark: {
    background: {
      default: '#0a0e13', // Very dark blue-black
      paper: '#1e2328',   // Dark blue-grey
      elevated: '#3c3c41', // Elevated surface
    },
    text: {
      primary: '#f0e6d2',   // Cream white
      secondary: '#cdbe91', // Light gold
      disabled: '#5bc0de',  // Muted blue
    },
    divider: '#463714',
  },
  
  // Light theme colors  
  light: {
    background: {
      default: '#f5f5f5',
      paper: '#ffffff',
      elevated: '#fafafa',
    },
    text: {
      primary: '#1e2328',
      secondary: '#3c3c41',
      disabled: '#a09b8c',
    },
    divider: '#e0e0e0',
  },
  
  // Status colors
  success: {
    main: '#4caf50',
    light: '#81c784',
    dark: '#388e3c',
  },
  error: {
    main: '#f44336',
    light: '#e57373',
    dark: '#d32f2f',
  },
  warning: {
    main: '#ff9800',
    light: '#ffb74d',
    dark: '#f57c00',
  },
  info: {
    main: '#2196f3',
    light: '#64b5f6',
    dark: '#1976d2',
  },
};

// Custom theme for League of Legends styling
const createLoLTheme = (mode: 'light' | 'dark'): Theme => {
  const isDark = mode === 'dark';
  
  return createTheme({
    palette: {
      mode,
      primary: {
        ...colors.primary,
        main: colors.primary[500],
      },
      secondary: {
        ...colors.secondary,
        main: colors.secondary[500],
      },
      background: {
        default: isDark ? colors.dark.background.default : colors.light.background.default,
        paper: isDark ? colors.dark.background.paper : colors.light.background.paper,
      },
      text: {
        primary: isDark ? colors.dark.text.primary : colors.light.text.primary,
        secondary: isDark ? colors.dark.text.secondary : colors.light.text.secondary,
        disabled: isDark ? colors.dark.text.disabled : colors.light.text.disabled,
      },
      divider: isDark ? colors.dark.divider : colors.light.divider,
      success: colors.success,
      error: colors.error,
      warning: colors.warning,
      info: colors.info,
    },
    typography: {
      fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
      h1: {
        fontFamily: '"Poppins", "Inter", sans-serif',
        fontWeight: 700,
        fontSize: '3rem',
        lineHeight: 1.2,
        color: isDark ? colors.dark.text.primary : colors.light.text.primary,
      },
      h2: {
        fontFamily: '"Poppins", "Inter", sans-serif',
        fontWeight: 600,
        fontSize: '2.25rem',
        lineHeight: 1.3,
        color: isDark ? colors.dark.text.primary : colors.light.text.primary,
      },
      h3: {
        fontFamily: '"Poppins", "Inter", sans-serif',
        fontWeight: 600,
        fontSize: '1.875rem',
        lineHeight: 1.4,
        color: isDark ? colors.dark.text.primary : colors.light.text.primary,
      },
      h4: {
        fontFamily: '"Poppins", "Inter", sans-serif',
        fontWeight: 600,
        fontSize: '1.5rem',
        lineHeight: 1.4,
        color: isDark ? colors.dark.text.primary : colors.light.text.primary,
      },
      h5: {
        fontFamily: '"Poppins", "Inter", sans-serif',
        fontWeight: 600,
        fontSize: '1.25rem',
        lineHeight: 1.5,
        color: isDark ? colors.dark.text.primary : colors.light.text.primary,
      },
      h6: {
        fontFamily: '"Poppins", "Inter", sans-serif',
        fontWeight: 600,
        fontSize: '1.125rem',
        lineHeight: 1.5,
        color: isDark ? colors.dark.text.primary : colors.light.text.primary,
      },
      body1: {
        fontSize: '1rem',
        lineHeight: 1.6,
        color: isDark ? colors.dark.text.primary : colors.light.text.primary,
      },
      body2: {
        fontSize: '0.875rem',
        lineHeight: 1.6,
        color: isDark ? colors.dark.text.secondary : colors.light.text.secondary,
      },
      button: {
        fontWeight: 600,
        textTransform: 'none',
        fontSize: '0.875rem',
      },
      caption: {
        fontSize: '0.75rem',
        lineHeight: 1.4,
        color: isDark ? colors.dark.text.disabled : colors.light.text.disabled,
      },
    },
    shape: {
      borderRadius: 8,
    },
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          '*': {
            scrollbarWidth: 'thin',
            scrollbarColor: `${isDark ? '#3c3c41' : '#c0c0c0'} transparent`,
          },
          '*::-webkit-scrollbar': {
            width: '8px',
            height: '8px',
          },
          '*::-webkit-scrollbar-track': {
            background: 'transparent',
          },
          '*::-webkit-scrollbar-thumb': {
            backgroundColor: isDark ? '#3c3c41' : '#c0c0c0',
            borderRadius: '4px',
            '&:hover': {
              backgroundColor: isDark ? '#5a5a5f' : '#a0a0a0',
            },
          },
          body: {
            backgroundColor: isDark ? colors.dark.background.default : colors.light.background.default,
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 6,
            textTransform: 'none',
            fontWeight: 600,
            padding: '10px 24px',
            transition: 'all 0.2s ease-in-out',
            '&:hover': {
              transform: 'translateY(-1px)',
              boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
            },
          },
          contained: {
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
            '&:hover': {
              boxShadow: '0 4px 16px rgba(0,0,0,0.2)',
            },
          },
          outlined: {
            borderWidth: '2px',
            '&:hover': {
              borderWidth: '2px',
            },
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            borderRadius: 12,
            border: `1px solid ${isDark ? '#3c3c41' : '#e0e0e0'}`,
            backgroundColor: isDark ? colors.dark.background.paper : colors.light.background.paper,
            boxShadow: isDark 
              ? '0 4px 16px rgba(0,0,0,0.3)'
              : '0 2px 12px rgba(0,0,0,0.08)',
            transition: 'all 0.2s ease-in-out',
            '&:hover': {
              transform: 'translateY(-2px)',
              boxShadow: isDark
                ? '0 8px 24px rgba(0,0,0,0.4)'
                : '0 4px 20px rgba(0,0,0,0.12)',
            },
          },
        },
      },
      MuiTextField: {
        styleOverrides: {
          root: {
            '& .MuiOutlinedInput-root': {
              borderRadius: 8,
              backgroundColor: isDark ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.02)',
              '&:hover .MuiOutlinedInput-notchedOutline': {
                borderColor: colors.primary[400],
              },
              '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                borderColor: colors.primary[500],
                borderWidth: '2px',
              },
            },
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            borderRadius: 6,
            fontWeight: 500,
          },
          colorPrimary: {
            backgroundColor: colors.primary[500],
            color: 'white',
          },
          colorSecondary: {
            backgroundColor: colors.secondary[500],
            color: 'white',
          },
        },
      },
      MuiAppBar: {
        styleOverrides: {
          root: {
            backgroundColor: isDark ? colors.dark.background.paper : colors.light.background.paper,
            borderBottom: `1px solid ${isDark ? '#3c3c41' : '#e0e0e0'}`,
            boxShadow: isDark
              ? '0 2px 8px rgba(0,0,0,0.3)'
              : '0 1px 4px rgba(0,0,0,0.1)',
          },
        },
      },
      MuiDrawer: {
        styleOverrides: {
          paper: {
            backgroundColor: isDark ? colors.dark.background.paper : colors.light.background.paper,
            borderColor: isDark ? '#3c3c41' : '#e0e0e0',
          },
        },
      },
      MuiListItem: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            margin: '2px 8px',
            '&:hover': {
              backgroundColor: isDark ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.04)',
            },
            '&.Mui-selected': {
              backgroundColor: `${colors.primary[500]}20`,
              '&:hover': {
                backgroundColor: `${colors.primary[500]}30`,
              },
            },
          },
        },
      },
      MuiTableCell: {
        styleOverrides: {
          head: {
            fontWeight: 600,
            color: isDark ? colors.dark.text.secondary : colors.light.text.secondary,
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            backgroundImage: 'none',
          },
        },
      },
    },
    breakpoints: {
      values: {
        xs: 0,
        sm: 600,
        md: 960,
        lg: 1280,
        xl: 1920,
      },
    },
    transitions: {
      easing: {
        easeInOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
        easeOut: 'cubic-bezier(0.0, 0, 0.2, 1)',
        easeIn: 'cubic-bezier(0.4, 0, 1, 1)',
        sharp: 'cubic-bezier(0.4, 0, 0.6, 1)',
      },
    },
  });
};

export const lightTheme = createLoLTheme('light');
export const darkTheme = createLoLTheme('dark');

export default { lightTheme, darkTheme };

// Theme hook will be used to switch between themes
export type ThemeMode = 'light' | 'dark' | 'auto';