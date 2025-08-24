import type { Preview } from '@storybook/react-vite'
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material'
import React from 'react'

// Herald.lol Gaming Theme
const heraldTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#C89B3C', // League of Legends gold
      dark: '#A0792E',
      light: '#D4AF58',
    },
    secondary: {
      main: '#3C8CE7', // Gaming blue
      dark: '#2E6AB3',
      light: '#5BA3EA',
    },
    background: {
      default: '#0F1419', // Dark gaming background
      paper: '#1A1A1A',
    },
    text: {
      primary: '#CDBE91', // LoL golden text
      secondary: '#A09B8C',
    },
    success: {
      main: '#00C851', // Gaming green
    },
    error: {
      main: '#FF4444', // Gaming red
    },
    warning: {
      main: '#FF8800', // Gaming orange
    },
  },
  typography: {
    fontFamily: '"Roboto", "Arial", sans-serif',
    h1: {
      fontWeight: 700,
      fontSize: '2.5rem',
      color: '#C89B3C',
    },
    h2: {
      fontWeight: 600,
      fontSize: '2rem',
      color: '#C89B3C',
    },
    h3: {
      fontWeight: 600,
      fontSize: '1.5rem',
      color: '#CDBE91',
    },
    body1: {
      color: '#CDBE91',
    },
    button: {
      fontWeight: 600,
      textTransform: 'none',
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiCard: {
      styleOverrides: {
        root: {
          backgroundColor: '#1A1A1A',
          border: '1px solid #463714',
          '&:hover': {
            borderColor: '#C89B3C',
          },
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          padding: '8px 16px',
          '&.MuiButton-containedPrimary': {
            background: 'linear-gradient(45deg, #C89B3C 30%, #D4AF58 90%)',
            '&:hover': {
              background: 'linear-gradient(45deg, #A0792E 30%, #C89B3C 90%)',
            },
          },
        },
      },
    },
  },
})

const lightTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#C89B3C',
      dark: '#A0792E',
      light: '#D4AF58',
    },
    secondary: {
      main: '#3C8CE7',
      dark: '#2E6AB3',
      light: '#5BA3EA',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Arial", sans-serif',
    button: {
      textTransform: 'none',
    },
  },
  shape: {
    borderRadius: 8,
  },
})

const preview: Preview = {
  parameters: {
    actions: { argTypesRegex: '^on[A-Z].*' },
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    backgrounds: {
      default: 'herald-dark',
      values: [
        {
          name: 'herald-dark',
          value: '#0F1419',
        },
        {
          name: 'herald-light',
          value: '#FFFFFF',
        },
        {
          name: 'gaming-dark',
          value: '#1A1A1A',
        },
      ],
    },
    viewport: {
      viewports: {
        mobile: {
          name: 'Mobile',
          styles: {
            width: '375px',
            height: '667px',
          },
        },
        tablet: {
          name: 'Tablet',
          styles: {
            width: '768px',
            height: '1024px',
          },
        },
        desktop: {
          name: 'Desktop',
          styles: {
            width: '1440px',
            height: '900px',
          },
        },
        gaming4k: {
          name: 'Gaming 4K',
          styles: {
            width: '3840px',
            height: '2160px',
          },
        },
      },
    },
    docs: {
      theme: heraldTheme,
    },
  },
  decorators: [
    (Story, context) => {
      const theme = context.globals.theme === 'light' ? lightTheme : heraldTheme
      return React.createElement(
        ThemeProvider,
        { theme },
        React.createElement(CssBaseline),
        React.createElement(Story)
      )
    },
  ],
  globalTypes: {
    theme: {
      description: 'Herald.lol Theme',
      defaultValue: 'dark',
      toolbar: {
        title: 'Theme',
        icon: 'paintbrush',
        items: [
          { value: 'light', title: 'Light', icon: 'sun' },
          { value: 'dark', title: 'Dark (Gaming)', icon: 'moon' },
        ],
        dynamicTitle: true,
      },
    },
  },
}

export default preview