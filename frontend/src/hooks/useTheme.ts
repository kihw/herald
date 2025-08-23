import { useState, useEffect, createContext, useContext } from 'react';
import { Theme } from '@mui/material/styles';
import { lightTheme, darkTheme, ThemeMode } from '@/theme';

interface ThemeContextType {
  themeMode: ThemeMode;
  theme: Theme;
  toggleTheme: () => void;
  setThemeMode: (mode: ThemeMode) => void;
}

export const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    // Fallback for when hook is used outside provider
    const [themeMode, setThemeMode] = useState<ThemeMode>(() => {
      const saved = localStorage.getItem('herald_theme_mode');
      return (saved as ThemeMode) || 'dark';
    });

    const [theme, setTheme] = useState<Theme>(() => {
      if (themeMode === 'auto') {
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? darkTheme : lightTheme;
      }
      return themeMode === 'dark' ? darkTheme : lightTheme;
    });

    useEffect(() => {
      localStorage.setItem('herald_theme_mode', themeMode);
      
      if (themeMode === 'auto') {
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
        const handleChange = () => {
          setTheme(mediaQuery.matches ? darkTheme : lightTheme);
        };
        
        handleChange();
        mediaQuery.addEventListener('change', handleChange);
        return () => mediaQuery.removeEventListener('change', handleChange);
      } else {
        setTheme(themeMode === 'dark' ? darkTheme : lightTheme);
      }
    }, [themeMode]);

    const toggleTheme = () => {
      setThemeMode(current => current === 'dark' ? 'light' : 'dark');
    };

    return {
      themeMode,
      theme,
      toggleTheme,
      setThemeMode,
    };
  }
  
  return context;
};