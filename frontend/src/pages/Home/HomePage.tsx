import { 
  Box, 
  Typography, 
  Button, 
  Card, 
  CardContent, 
  Grid,
  Container,
  useTheme
} from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

const HomePage = () => {
  const theme = useTheme();
  const { isAuthenticated } = useAuth();

  return (
    <Container maxWidth="xl">
      <Box sx={{ py: 8 }}>
        {/* Hero Section */}
        <Box textAlign="center" mb={8}>
          <Typography 
            variant="h1" 
            component="h1" 
            gutterBottom
            sx={{
              background: `linear-gradient(135deg, ${theme.palette.primary.main} 0%, ${theme.palette.secondary.main} 100%)`,
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              color: 'transparent',
              mb: 3,
            }}
          >
            Herald.lol
          </Typography>
          
          <Typography variant="h4" color="text.secondary" gutterBottom>
            Professional Gaming Analytics
          </Typography>
          
          <Typography variant="h6" color="text.primary" sx={{ mb: 4, maxWidth: 600, mx: 'auto' }}>
            Transform your League of Legends gameplay with AI-powered analytics, 
            personalized coaching, and professional-grade performance insights.
          </Typography>

          <Box display="flex" gap={2} justifyContent="center" flexWrap="wrap">
            {!isAuthenticated ? (
              <>
                <Button
                  component={RouterLink}
                  to="/register"
                  variant="contained"
                  size="large"
                  sx={{ minWidth: 120 }}
                >
                  Get Started
                </Button>
                <Button
                  component={RouterLink}
                  to="/login"
                  variant="outlined"
                  size="large"
                  sx={{ minWidth: 120 }}
                >
                  Sign In
                </Button>
              </>
            ) : (
              <Button
                component={RouterLink}
                to="/dashboard"
                variant="contained"
                size="large"
                sx={{ minWidth: 120 }}
              >
                Go to Dashboard
              </Button>
            )}
          </Box>
        </Box>

        {/* Features Section */}
        <Grid container spacing={4} sx={{ mb: 8 }}>
          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
              <CardContent>
                <Typography variant="h5" component="h3" gutterBottom color="primary">
                  🎯 Performance Analytics
                </Typography>
                <Typography variant="body1">
                  Deep dive into your gameplay with advanced metrics: KDA analysis, 
                  CS/min optimization, vision score tracking, and damage efficiency.
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
              <CardContent>
                <Typography variant="h5" component="h3" gutterBottom color="primary">
                  🤖 AI Coaching
                </Typography>
                <Typography variant="body1">
                  Get personalized improvement recommendations powered by machine learning 
                  and professional gameplay analysis.
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} md={4}>
            <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
              <CardContent>
                <Typography variant="h5" component="h3" gutterBottom color="primary">
                  ⚡ Real-time Insights
                </Typography>
                <Typography variant="body1">
                  Track your progress in real-time with automated match analysis 
                  and instant performance feedback after every game.
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Stats Section */}
        <Box 
          textAlign="center" 
          sx={{ 
            py: 6, 
            bgcolor: 'background.paper', 
            borderRadius: 2,
            border: `1px solid ${theme.palette.divider}`,
          }}
        >
          <Typography variant="h4" gutterBottom>
            Join Thousands of Players
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
            Herald.lol helps players improve their gameplay and climb the ranked ladder
          </Typography>
          
          <Grid container spacing={4} justifyContent="center">
            <Grid item xs={6} sm={3}>
              <Typography variant="h3" color="primary" fontWeight="bold">
                50k+
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Active Users
              </Typography>
            </Grid>
            <Grid item xs={6} sm={3}>
              <Typography variant="h3" color="primary" fontWeight="bold">
                1M+
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Matches Analyzed
              </Typography>
            </Grid>
            <Grid item xs={6} sm={3}>
              <Typography variant="h3" color="primary" fontWeight="bold">
                85%
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Rank Improvement
              </Typography>
            </Grid>
            <Grid item xs={6} sm={3}>
              <Typography variant="h3" color="primary" fontWeight="bold">
                &lt;5s
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Analysis Speed
              </Typography>
            </Grid>
          </Grid>
        </Box>
      </Box>
    </Container>
  );
};

export default HomePage;