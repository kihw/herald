import { useTheme, useMediaQuery, Breakpoint } from '@mui/material';
import { useMemo } from 'react';

// Hook for responsive design utilities
export const useResponsive = () => {
  const theme = useTheme();
  
  // Breakpoint queries
  const isXs = useMediaQuery(theme.breakpoints.only('xs'));
  const isSm = useMediaQuery(theme.breakpoints.only('sm'));
  const isMd = useMediaQuery(theme.breakpoints.only('md'));
  const isLg = useMediaQuery(theme.breakpoints.only('lg'));
  const isXl = useMediaQuery(theme.breakpoints.only('xl'));

  // Direction queries
  const isSmUp = useMediaQuery(theme.breakpoints.up('sm'));
  const isMdUp = useMediaQuery(theme.breakpoints.up('md'));
  const isLgUp = useMediaQuery(theme.breakpoints.up('lg'));
  const isXlUp = useMediaQuery(theme.breakpoints.up('xl'));

  const isSmDown = useMediaQuery(theme.breakpoints.down('sm'));
  const isMdDown = useMediaQuery(theme.breakpoints.down('md'));
  const isLgDown = useMediaQuery(theme.breakpoints.down('lg'));

  // Mobile detection
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const isTablet = useMediaQuery(theme.breakpoints.between('sm', 'md'));
  const isDesktop = useMediaQuery(theme.breakpoints.up('lg'));

  // Current breakpoint
  const currentBreakpoint = useMemo((): Breakpoint => {
    if (isXs) return 'xs';
    if (isSm) return 'sm';
    if (isMd) return 'md';
    if (isLg) return 'lg';
    return 'xl';
  }, [isXs, isSm, isMd, isLg]);

  // Responsive values helper
  const getResponsiveValue = <T>(values: {
    xs?: T;
    sm?: T;
    md?: T;
    lg?: T;
    xl?: T;
    mobile?: T;
    tablet?: T;
    desktop?: T;
  }): T | undefined => {
    // Priority: specific breakpoint > device type > fallback
    if (values[currentBreakpoint] !== undefined) {
      return values[currentBreakpoint];
    }
    
    if (isMobile && values.mobile !== undefined) {
      return values.mobile;
    }
    
    if (isTablet && values.tablet !== undefined) {
      return values.tablet;
    }
    
    if (isDesktop && values.desktop !== undefined) {
      return values.desktop;
    }

    // Fallback to largest available value
    return values.xl || values.lg || values.md || values.sm || values.xs;
  };

  // Grid spacing helper
  const getGridSpacing = () => getResponsiveValue({
    xs: 1,
    sm: 2,
    md: 3,
    lg: 3,
  }) || 2;

  // Card padding helper
  const getCardPadding = () => getResponsiveValue({
    xs: 2,
    sm: 3,
    md: 3,
    lg: 4,
  }) || 3;

  // Typography scale helper
  const getTypographyScale = () => getResponsiveValue({
    xs: 0.8,
    sm: 0.9,
    md: 1,
    lg: 1,
  }) || 1;

  // Dialog sizing helper
  const getDialogMaxWidth = () => getResponsiveValue({
    xs: 'xs' as const,
    sm: 'sm' as const,
    md: 'md' as const,
    lg: 'lg' as const,
    xl: 'xl' as const,
  }) || 'md';

  // Chart height helper
  const getChartHeight = () => getResponsiveValue({
    xs: 250,
    sm: 300,
    md: 350,
    lg: 400,
  }) || 350;

  return {
    // Breakpoint flags
    isXs,
    isSm,
    isMd,
    isLg,
    isXl,
    isSmUp,
    isMdUp,
    isLgUp,
    isXlUp,
    isSmDown,
    isMdDown,
    isLgDown,
    
    // Device type flags
    isMobile,
    isTablet,
    isDesktop,
    
    // Current breakpoint
    currentBreakpoint,
    
    // Helper functions
    getResponsiveValue,
    getGridSpacing,
    getCardPadding,
    getTypographyScale,
    getDialogMaxWidth,
    getChartHeight,
  };
};

export default useResponsive;