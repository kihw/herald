import React from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Skeleton,
  Grid,
  Typography,
  LinearProgress,
  CircularProgress,
  Fade,
  useTheme,
  alpha,
  keyframes,
  styled,
} from '@mui/material';

// Animations pour les loading states
const shimmer = keyframes`
  0% {
    background-position: -200px 0;
  }
  100% {
    background-position: calc(200px + 100%) 0;
  }
`;

const pulse = keyframes`
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
`;

const wave = keyframes`
  0%, 60%, 100% {
    transform: initial;
  }
  30% {
    transform: translateY(-15px);
  }
`;

// Styled components pour les animations
const ShimmerBox = styled(Box)(({ theme }) => ({
  background: `linear-gradient(90deg, ${alpha(theme.palette.grey[300], 0.2)} 25%, ${alpha(theme.palette.grey[300], 0.4)} 50%, ${alpha(theme.palette.grey[300], 0.2)} 75%)`,
  backgroundSize: '200px 100%',
  animation: `${shimmer} 2s infinite linear`,
  borderRadius: theme.shape.borderRadius,
}));

const PulsingBox = styled(Box)({
  animation: `${pulse} 1.5s ease-in-out infinite`,
});

const WaveBox = styled(Box)({
  display: 'inline-block',
  animation: `${wave} 1.3s ease-in-out infinite`,
});

// Loading skeleton pour les cartes de statistiques
export const StatCardSkeleton: React.FC = () => {
  return (
    <Card>
      <CardContent>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Box flex={1}>
            <Skeleton variant="text" width="60%" height={20} sx={{ mb: 1 }} />
            <Skeleton variant="text" width="40%" height={32} />
            <Skeleton variant="text" width="80%" height={16} />
          </Box>
          <Skeleton variant="circular" width={40} height={40} />
        </Box>
      </CardContent>
    </Card>
  );
};

// Loading skeleton pour les graphiques
export const ChartSkeleton: React.FC<{ height?: number }> = ({ height = 300 }) => {
  const theme = useTheme();
  
  return (
    <Card>
      <CardHeader
        title={<Skeleton variant="text" width="40%" height={24} />}
        action={<Skeleton variant="circular" width={24} height={24} />}
      />
      <CardContent>
        <Box position="relative" height={height}>
          {/* Axes simulés */}
          <Box
            position="absolute"
            bottom={40}
            left={40}
            right={20}
            height={2}
            bgcolor={alpha(theme.palette.grey[300], 0.3)}
          />
          <Box
            position="absolute"
            bottom={40}
            left={40}
            width={2}
            top={20}
            bgcolor={alpha(theme.palette.grey[300], 0.3)}
          />
          
          {/* Barres de données simulées */}
          {[...Array(6)].map((_, i) => (
            <PulsingBox
              key={i}
              position="absolute"
              bottom={42}
              left={60 + i * 80}
              width={40}
              height={Math.random() * 150 + 50}
              bgcolor={alpha(theme.palette.primary.main, 0.2)}
              sx={{ animationDelay: `${i * 0.2}s` }}
            />
          ))}
          
          {/* Labels d'axes */}
          <Box position="absolute" bottom={10} left={40} right={20} display="flex" justifyContent="space-between">
            {[...Array(6)].map((_, i) => (
              <Skeleton key={i} variant="text" width={30} height={16} />
            ))}
          </Box>
          <Box position="absolute" top={20} left={5} display="flex" flexDirection="column" gap={2}>
            {[...Array(4)].map((_, i) => (
              <Skeleton key={i} variant="text" width={25} height={16} />
            ))}
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

// Loading skeleton pour les tableaux
export const TableSkeleton: React.FC<{ rows?: number; columns?: number }> = ({ 
  rows = 5, 
  columns = 4 
}) => {
  return (
    <Card>
      <CardHeader
        title={<Skeleton variant="text" width="30%" height={24} />}
      />
      <CardContent sx={{ p: 0 }}>
        <Box p={2}>
          {/* En-têtes */}
          <Grid container spacing={2} sx={{ mb: 2 }}>
            {[...Array(columns)].map((_, i) => (
              <Grid item xs={12 / columns} key={i}>
                <Skeleton variant="text" width="80%" height={20} />
              </Grid>
            ))}
          </Grid>
          
          {/* Lignes */}
          {[...Array(rows)].map((_, rowIndex) => (
            <Grid container spacing={2} key={rowIndex} sx={{ mb: 1 }}>
              {[...Array(columns)].map((_, colIndex) => (
                <Grid item xs={12 / columns} key={colIndex}>
                  <Skeleton 
                    variant="text" 
                    width={`${60 + Math.random() * 30}%`} 
                    height={16}
                    sx={{ animationDelay: `${(rowIndex * columns + colIndex) * 0.1}s` }}
                  />
                </Grid>
              ))}
            </Grid>
          ))}
        </Box>
      </CardContent>
    </Card>
  );
};

// Loading skeleton pour les listes
export const ListSkeleton: React.FC<{ items?: number }> = ({ items = 5 }) => {
  return (
    <Card>
      <CardHeader
        title={<Skeleton variant="text" width="40%" height={24} />}
      />
      <CardContent>
        {[...Array(items)].map((_, i) => (
          <Box key={i} display="flex" alignItems="center" gap={2} sx={{ mb: 2 }}>
            <Skeleton variant="circular" width={40} height={40} />
            <Box flex={1}>
              <Skeleton variant="text" width="60%" height={20} />
              <Skeleton variant="text" width="40%" height={16} />
            </Box>
            <Skeleton variant="rectangular" width={60} height={24} sx={{ borderRadius: 1 }} />
          </Box>
        ))}
      </CardContent>
    </Card>
  );
};

// Loading skeleton pour la heatmap
export const HeatmapSkeleton: React.FC = () => {
  const theme = useTheme();
  
  return (
    <Card>
      <CardHeader
        title={<Skeleton variant="text" width="40%" height={24} />}
        action={
          <Box display="flex" gap={1}>
            <Skeleton variant="circular" width={32} height={32} />
            <Skeleton variant="circular" width={32} height={32} />
          </Box>
        }
      />
      <CardContent>
        {/* Contrôles */}
        <Box display="flex" gap={2} mb={3}>
          <Skeleton variant="rectangular" width={120} height={32} sx={{ borderRadius: 1 }} />
          <Skeleton variant="rectangular" width={120} height={32} sx={{ borderRadius: 1 }} />
          <Skeleton variant="rectangular" width={80} height={32} sx={{ borderRadius: 1 }} />
        </Box>
        
        {/* Grille heatmap */}
        <Box>
          {[...Array(5)].map((_, rowIndex) => (
            <Box key={rowIndex} display="flex" gap={1} mb={1}>
              {[...Array(8)].map((_, colIndex) => (
                <ShimmerBox
                  key={colIndex}
                  width={60}
                  height={40}
                  sx={{ 
                    animationDelay: `${(rowIndex * 8 + colIndex) * 0.1}s`,
                    bgcolor: alpha(theme.palette.primary.main, Math.random() * 0.5 + 0.1),
                  }}
                />
              ))}
            </Box>
          ))}
        </Box>
      </CardContent>
    </Card>
  );
};

// Loading skeleton pour la timeline
export const TimelineSkeleton: React.FC = () => {
  return (
    <Card>
      <CardHeader
        title={<Skeleton variant="text" width="50%" height={24} />}
        action={
          <Box display="flex" gap={1}>
            {[...Array(3)].map((_, i) => (
              <Skeleton key={i} variant="circular" width={32} height={32} />
            ))}
          </Box>
        }
      />
      <CardContent>
        {/* Contrôles de lecture */}
        <Box display="flex" alignItems="center" gap={2} mb={3}>
          <Box display="flex" gap={1}>
            {[...Array(3)].map((_, i) => (
              <Skeleton key={i} variant="circular" width={32} height={32} />
            ))}
          </Box>
          <Skeleton variant="rectangular" width="100%" height={4} sx={{ borderRadius: 2 }} />
          <Skeleton variant="rectangular" width={80} height={32} sx={{ borderRadius: 1 }} />
        </Box>
        
        {/* Timeline */}
        <Box position="relative" height={300}>
          {/* Ligne principale */}
          <Box
            position="absolute"
            top="50%"
            left={0}
            right={0}
            height={2}
            bgcolor="grey.300"
          />
          
          {/* Points d'événements */}
          {[...Array(8)].map((_, i) => (
            <PulsingBox
              key={i}
              position="absolute"
              top="45%"
              left={`${10 + i * 10}%`}
              width={16}
              height={16}
              bgcolor="primary.main"
              borderRadius="50%"
              sx={{ animationDelay: `${i * 0.3}s` }}
            />
          ))}
        </Box>
        
        {/* Événement actuel */}
        <Box mt={2} p={2} bgcolor="grey.50" borderRadius={1}>
          <Skeleton variant="text" width="40%" height={24} />
          <Skeleton variant="text" width="20%" height={16} sx={{ my: 1 }} />
          <Skeleton variant="text" width="80%" height={16} />
        </Box>
      </CardContent>
    </Card>
  );
};

// Composant de loading avec progression
export const ProgressiveLoading: React.FC<{
  steps: string[];
  currentStep: number;
  message?: string;
}> = ({ steps, currentStep, message }) => {
  const progress = (currentStep / steps.length) * 100;
  
  return (
    <Box textAlign="center" p={4}>
      <CircularProgress size={60} sx={{ mb: 2 }} />
      
      <Typography variant="h6" gutterBottom>
        {message || 'Chargement en cours...'}
      </Typography>
      
      <Box sx={{ mb: 2 }}>
        <LinearProgress 
          variant="determinate" 
          value={progress} 
          sx={{ 
            height: 8, 
            borderRadius: 4,
            '& .MuiLinearProgress-bar': {
              borderRadius: 4,
            }
          }} 
        />
      </Box>
      
      <Typography variant="body2" color="text.secondary">
        Étape {currentStep + 1} sur {steps.length}: {steps[currentStep] || 'Finalisation...'}
      </Typography>
      
      {/* Animation de points */}
      <Box sx={{ mt: 2 }}>
        {[...Array(3)].map((_, i) => (
          <WaveBox 
            key={i} 
            sx={{ 
              animationDelay: `${i * 0.2}s`,
              mx: 0.5,
              fontSize: '1.5rem',
            }}
          >
            •
          </WaveBox>
        ))}
      </Box>
    </Box>
  );
};

// Loading skeleton pour le dashboard complet
export const DashboardSkeleton: React.FC = () => {
  return (
    <Box>
      {/* En-tête */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Skeleton variant="text" width={300} height={32} />
          <Skeleton variant="text" width={200} height={20} />
        </Box>
        <Box display="flex" gap={1}>
          <Skeleton variant="rectangular" width={80} height={32} sx={{ borderRadius: 1 }} />
          <Skeleton variant="circular" width={32} height={32} />
        </Box>
      </Box>
      
      {/* Cartes de statistiques */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        {[...Array(4)].map((_, i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <StatCardSkeleton />
          </Grid>
        ))}
      </Grid>
      
      {/* Graphiques */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <ChartSkeleton height={400} />
        </Grid>
        <Grid item xs={12} md={4}>
          <ListSkeleton items={6} />
        </Grid>
      </Grid>
    </Box>
  );
};