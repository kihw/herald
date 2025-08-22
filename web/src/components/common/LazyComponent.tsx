import React, { Suspense, lazy, ComponentType } from 'react';
import {
  Box,
  Skeleton,
  Grid,
  Card,
  CardContent,
  CircularProgress,
} from '@mui/material';
import { usePerformance } from '../../hooks/usePerformance';

interface LazyComponentProps {
  fallback?: React.ReactNode;
  error?: React.ReactNode;
  minLoadTime?: number;
}

// Higher-order component for lazy loading with performance optimization
export const withLazy = <P extends object>(
  importFunc: () => Promise<{ default: ComponentType<P> }>,
  options: LazyComponentProps = {}
) => {
  const LazyComponent = lazy(async () => {
    const [component, delay] = await Promise.all([
      importFunc(),
      // Minimum load time to prevent flash
      options.minLoadTime ? new Promise(resolve => setTimeout(resolve, options.minLoadTime)) : Promise.resolve(),
    ]);
    return component;
  });

  return (props: P) => {
    const { shouldLazyLoad } = usePerformance();
    
    const defaultFallback = (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 200 }}>
        <CircularProgress />
      </Box>
    );

    if (!shouldLazyLoad) {
      return (
        <Suspense fallback={options.fallback || defaultFallback}>
          <LazyComponent {...props} />
        </Suspense>
      );
    }

    return (
      <Suspense fallback={options.fallback || defaultFallback}>
        <LazyComponent {...props} />
      </Suspense>
    );
  };
};

// Skeleton components for different layouts
export const CardSkeleton: React.FC<{ count?: number; height?: number }> = ({ 
  count = 1, 
  height = 200 
}) => (
  <Grid container spacing={3}>
    {Array.from({ length: count }, (_, index) => (
      <Grid item xs={12} sm={6} md={4} key={index}>
        <Card>
          <CardContent>
            <Skeleton variant="rectangular" height={height} />
            <Skeleton variant="text" sx={{ mt: 2 }} />
            <Skeleton variant="text" width="60%" />
          </CardContent>
        </Card>
      </Grid>
    ))}
  </Grid>
);

export const ListSkeleton: React.FC<{ count?: number }> = ({ count = 5 }) => (
  <Box>
    {Array.from({ length: count }, (_, index) => (
      <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <Skeleton variant="circular" width={40} height={40} sx={{ mr: 2 }} />
        <Box sx={{ flexGrow: 1 }}>
          <Skeleton variant="text" />
          <Skeleton variant="text" width="60%" />
        </Box>
      </Box>
    ))}
  </Box>
);

export const ChartSkeleton: React.FC<{ height?: number }> = ({ height = 300 }) => (
  <Card>
    <CardContent>
      <Skeleton variant="text" width="40%" sx={{ mb: 2 }} />
      <Skeleton variant="rectangular" height={height} />
    </CardContent>
  </Card>
);

export const TableSkeleton: React.FC<{ rows?: number; columns?: number }> = ({ 
  rows = 5, 
  columns = 4 
}) => (
  <Box>
    {Array.from({ length: rows }, (_, rowIndex) => (
      <Box key={rowIndex} sx={{ display: 'flex', gap: 2, mb: 1 }}>
        {Array.from({ length: columns }, (_, colIndex) => (
          <Skeleton key={colIndex} variant="text" sx={{ flex: 1 }} />
        ))}
      </Box>
    ))}
  </Box>
);

// Progressive loading component
interface ProgressiveImageProps {
  src: string;
  alt: string;
  placeholder?: string;
  className?: string;
  style?: React.CSSProperties;
}

export const ProgressiveImage: React.FC<ProgressiveImageProps> = ({
  src,
  alt,
  placeholder = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIwIiBoZWlnaHQ9IjI0MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PC9zdmc+',
  className,
  style,
}) => {
  const { shouldLazyLoad } = usePerformance();
  const [isLoaded, setIsLoaded] = React.useState(false);
  const [isError, setIsError] = React.useState(false);

  const handleLoad = () => {
    setIsLoaded(true);
  };

  const handleError = () => {
    setIsError(true);
  };

  if (!shouldLazyLoad) {
    return (
      <img
        src={src}
        alt={alt}
        className={className}
        style={style}
        onLoad={handleLoad}
        onError={handleError}
      />
    );
  }

  return (
    <Box sx={{ position: 'relative', overflow: 'hidden' }}>
      {!isLoaded && !isError && (
        <img
          src={placeholder}
          alt=""
          style={{
            width: '100%',
            height: 'auto',
            filter: 'blur(5px)',
            transform: 'scale(1.1)',
            transition: 'all 0.3s ease',
            ...style,
          }}
        />
      )}
      <img
        src={src}
        alt={alt}
        className={className}
        style={{
          ...style,
          opacity: isLoaded ? 1 : 0,
          transition: 'opacity 0.3s ease',
          position: isLoaded ? 'static' : 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: 'auto',
        }}
        onLoad={handleLoad}
        onError={handleError}
      />
      {isError && (
        <Skeleton
          variant="rectangular"
          sx={{
            width: '100%',
            height: 200,
            position: 'absolute',
            top: 0,
            left: 0,
          }}
        />
      )}
    </Box>
  );
};

export default withLazy;