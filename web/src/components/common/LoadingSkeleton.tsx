import React from 'react';
import {
  Card,
  CardContent,
  Grid,
  Skeleton,
  Box,
} from '@mui/material';

interface LoadingSkeletonProps {
  type: 'stats' | 'matches' | 'settings';
  count?: number;
}

const LoadingSkeleton: React.FC<LoadingSkeletonProps> = ({ type, count = 1 }) => {
  if (type === 'stats') {
    return (
      <Grid container spacing={2}>
        {[1, 2, 3, 4].map((i) => (
          <Grid item xs={12} sm={6} md={3} key={i}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  <Skeleton variant="circular" width={40} height={40} sx={{ mr: 2 }} />
                  <Skeleton variant="text" sx={{ fontSize: '1.2rem', width: '60%' }} />
                </Box>
                <Skeleton variant="text" sx={{ fontSize: '2rem', mb: 1, width: '40%' }} />
                <Skeleton variant="text" sx={{ fontSize: '0.9rem', width: '80%' }} />
              </CardContent>
            </Card>
          </Grid>
        ))}
        
        {/* Additional cards for favorite champion and recent performance */}
        <Grid item xs={12} sm={6}>
          <Card>
            <CardContent>
              <Skeleton variant="text" sx={{ fontSize: '1.2rem', mb: 2 }} />
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Skeleton variant="circular" width={56} height={56} sx={{ mr: 2 }} />
                <Box sx={{ flexGrow: 1 }}>
                  <Skeleton variant="text" sx={{ fontSize: '1.5rem', mb: 1 }} />
                  <Skeleton variant="text" sx={{ fontSize: '0.9rem' }} />
                </Box>
              </Box>
              <Skeleton variant="rectangular" height={4} />
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <Card>
            <CardContent>
              <Skeleton variant="text" sx={{ fontSize: '1.2rem', mb: 2 }} />
              <Box sx={{ mb: 3 }}>
                <Skeleton variant="text" sx={{ mb: 1 }} />
                <Skeleton variant="text" sx={{ mb: 1, width: '60%' }} />
                <Skeleton variant="rectangular" height={4} />
              </Box>
              <Box>
                <Skeleton variant="text" sx={{ mb: 1 }} />
                <Skeleton variant="text" sx={{ mb: 1, width: '60%' }} />
                <Skeleton variant="rectangular" height={4} />
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }

  if (type === 'matches') {
    return (
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
            <Skeleton variant="text" sx={{ fontSize: '1.2rem', width: '30%' }} />
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Skeleton variant="rectangular" width={80} height={32} />
              <Skeleton variant="rectangular" width={80} height={32} />
              <Skeleton variant="rectangular" width={120} height={32} />
            </Box>
          </Box>
          
          {Array.from({ length: count }).map((_, index) => (
            <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 2, p: 2, border: '1px solid #e0e0e0', borderRadius: 1 }}>
              <Skeleton variant="circular" width={32} height={32} sx={{ mr: 2 }} />
              <Box sx={{ flexGrow: 1, display: 'flex', gap: 2 }}>
                <Skeleton variant="text" sx={{ width: '15%' }} />
                <Skeleton variant="rectangular" width={80} height={20} />
                <Skeleton variant="rectangular" width={60} height={20} />
                <Skeleton variant="text" sx={{ width: '10%' }} />
                <Skeleton variant="text" sx={{ width: '10%' }} />
                <Skeleton variant="text" sx={{ width: '10%' }} />
                <Skeleton variant="rectangular" width={60} height={20} />
              </Box>
              <Skeleton variant="circular" width={24} height={24} />
            </Box>
          ))}
        </CardContent>
      </Card>
    );
  }

  if (type === 'settings') {
    return (
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
            <Skeleton variant="text" sx={{ fontSize: '1.2rem', width: '30%' }} />
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Skeleton variant="rectangular" width={80} height={32} />
              <Skeleton variant="rectangular" width={120} height={32} />
            </Box>
          </Box>

          {/* Settings sections */}
          {[1, 2, 3].map((section) => (
            <Box key={section} sx={{ mb: 4 }}>
              <Skeleton variant="text" sx={{ fontSize: '1.1rem', mb: 2, width: '25%' }} />
              {[1, 2].map((item) => (
                <Box key={item} sx={{ mb: 2 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <Skeleton variant="rectangular" width={38} height={20} sx={{ mr: 2 }} />
                    <Skeleton variant="text" sx={{ width: '40%' }} />
                  </Box>
                  <Skeleton variant="text" sx={{ fontSize: '0.8rem', width: '60%', ml: 5 }} />
                </Box>
              ))}
            </Box>
          ))}
        </CardContent>
      </Card>
    );
  }

  return null;
};

export default LoadingSkeleton;
