import React from 'react';
import {
  Box,
  Container,
  Grid,
  useTheme,
  SxProps,
  Theme,
} from '@mui/material';
import useResponsive from '../../hooks/useResponsive';

interface ResponsiveContainerProps {
  children: React.ReactNode;
  maxWidth?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | false;
  disableGutters?: boolean;
  sx?: SxProps<Theme>;
  spacing?: number;
}

export const ResponsiveContainer: React.FC<ResponsiveContainerProps> = ({
  children,
  maxWidth,
  disableGutters,
  sx,
  spacing,
}) => {
  const { getGridSpacing, isMobile } = useResponsive();
  const theme = useTheme();

  const responsiveMaxWidth = maxWidth || (isMobile ? 'sm' : 'xl');
  const responsiveSpacing = spacing ?? getGridSpacing();

  return (
    <Container
      maxWidth={responsiveMaxWidth}
      disableGutters={disableGutters}
      sx={{
        px: isMobile ? 1 : 3,
        py: responsiveSpacing,
        ...sx,
      }}
    >
      {children}
    </Container>
  );
};

interface ResponsiveGridProps {
  children: React.ReactNode;
  spacing?: number;
  sx?: SxProps<Theme>;
}

export const ResponsiveGrid: React.FC<ResponsiveGridProps> = ({
  children,
  spacing,
  sx,
}) => {
  const { getGridSpacing } = useResponsive();
  
  return (
    <Grid
      container
      spacing={spacing ?? getGridSpacing()}
      sx={sx}
    >
      {children}
    </Grid>
  );
};

interface ResponsiveCardContainerProps {
  children: React.ReactNode;
  columns?: {
    xs?: number;
    sm?: number;
    md?: number;
    lg?: number;
    xl?: number;
  };
  minHeight?: number;
  spacing?: number;
}

export const ResponsiveCardContainer: React.FC<ResponsiveCardContainerProps> = ({
  children,
  columns = { xs: 1, sm: 2, md: 3, lg: 4 },
  minHeight,
  spacing,
}) => {
  const { getGridSpacing } = useResponsive();

  return (
    <ResponsiveGrid spacing={spacing ?? getGridSpacing()}>
      {React.Children.map(children, (child, index) => (
        <Grid
          item
          xs={12 / (columns.xs || 1)}
          sm={12 / (columns.sm || columns.xs || 1)}
          md={12 / (columns.md || columns.sm || columns.xs || 1)}
          lg={12 / (columns.lg || columns.md || columns.sm || columns.xs || 1)}
          xl={12 / (columns.xl || columns.lg || columns.md || columns.sm || columns.xs || 1)}
          key={index}
          sx={{
            minHeight: minHeight ? `${minHeight}px` : 'auto',
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          {child}
        </Grid>
      ))}
    </ResponsiveGrid>
  );
};

interface ResponsiveDialogContainerProps {
  children: React.ReactNode;
  fullScreenOnMobile?: boolean;
}

export const ResponsiveDialogContainer: React.FC<ResponsiveDialogContainerProps> = ({
  children,
  fullScreenOnMobile = true,
}) => {
  const { isMobile, getDialogMaxWidth } = useResponsive();
  
  return (
    <Box
      sx={{
        width: '100%',
        maxWidth: isMobile && fullScreenOnMobile ? '100vw' : undefined,
        height: isMobile && fullScreenOnMobile ? '100vh' : 'auto',
        maxHeight: isMobile && fullScreenOnMobile ? '100vh' : '90vh',
        overflow: 'auto',
      }}
    >
      {children}
    </Box>
  );
};

interface AdaptiveLayoutProps {
  sidebar?: React.ReactNode;
  main: React.ReactNode;
  sidebarWidth?: number;
  collapsible?: boolean;
  defaultCollapsed?: boolean;
}

export const AdaptiveLayout: React.FC<AdaptiveLayoutProps> = ({
  sidebar,
  main,
  sidebarWidth = 280,
  collapsible = true,
  defaultCollapsed = false,
}) => {
  const { isMobile, isTablet } = useResponsive();
  const [collapsed, setCollapsed] = React.useState(defaultCollapsed);

  // Auto-collapse on mobile/tablet
  React.useEffect(() => {
    if (isMobile || isTablet) {
      setCollapsed(true);
    }
  }, [isMobile, isTablet]);

  const sidebarComponent = sidebar && (
    <Box
      sx={{
        width: collapsed ? 0 : sidebarWidth,
        minWidth: collapsed ? 0 : sidebarWidth,
        height: '100vh',
        overflow: 'hidden',
        transition: 'width 0.3s ease',
        borderRight: collapsed ? 'none' : '1px solid',
        borderColor: 'divider',
        display: collapsed && (isMobile || isTablet) ? 'none' : 'block',
      }}
    >
      {sidebar}
    </Box>
  );

  const mainComponent = (
    <Box
      sx={{
        flexGrow: 1,
        minWidth: 0,
        height: '100vh',
        overflow: 'auto',
        ml: collapsed || !sidebar ? 0 : `${sidebarWidth}px`,
        transition: 'margin-left 0.3s ease',
      }}
    >
      {main}
    </Box>
  );

  return (
    <Box sx={{ display: 'flex', width: '100%', height: '100vh' }}>
      {sidebarComponent}
      {mainComponent}
    </Box>
  );
};

interface FlexibleStackProps {
  children: React.ReactNode;
  direction?: 'row' | 'column' | 'responsive';
  spacing?: number;
  align?: 'flex-start' | 'center' | 'flex-end' | 'stretch';
  justify?: 'flex-start' | 'center' | 'flex-end' | 'space-between' | 'space-around';
  wrap?: boolean;
  sx?: SxProps<Theme>;
}

export const FlexibleStack: React.FC<FlexibleStackProps> = ({
  children,
  direction = 'responsive',
  spacing = 2,
  align = 'stretch',
  justify = 'flex-start',
  wrap = false,
  sx,
}) => {
  const { isMobile } = useResponsive();
  
  const flexDirection = direction === 'responsive' ? (isMobile ? 'column' : 'row') : direction;

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection,
        alignItems: align,
        justifyContent: justify,
        flexWrap: wrap ? 'wrap' : 'nowrap',
        gap: spacing,
        ...sx,
      }}
    >
      {children}
    </Box>
  );
};

export default ResponsiveContainer;