import { useState } from 'react';
import { Link as RouterLink, useNavigate, useLocation } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Typography,
  Link,
  Alert,
  Container,
  InputAdornment,
  IconButton,
  Divider,
  Stack,
  Chip,
  Paper,
} from '@mui/material';
import { 
  Visibility, 
  VisibilityOff, 
  Email, 
  Lock,
  SportsEsports,
  TrendingUp,
  Security,
  Google,
} from '@mui/icons-material';
import { useAuth } from '@/hooks/useAuth';
import { LoginRequest } from '@/types';

const LoginPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  
  const [formData, setFormData] = useState<LoginRequest>({
    email: '',
    password: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const [step, setStep] = useState<'login' | 'mfa'>('login');
  const [mfaCode, setMfaCode] = useState('');

  // Get the intended destination after login
  const from = location.state?.from?.pathname || '/dashboard';

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    // Clear errors when user types
    if (errors) setErrors('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setErrors('');

    try {
      // Simulate MFA requirement for demo
      if (formData.email.includes('mfa')) {
        setStep('mfa');
        setIsLoading(false);
        return;
      }
      
      await login(formData);
      navigate(from, { replace: true });
    } catch (error: any) {
      setErrors(error?.message || 'Login failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleMfaSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!mfaCode || mfaCode.length !== 6) {
      setErrors('Please enter a valid 6-digit code');
      return;
    }

    setIsLoading(true);
    try {
      // Simulate MFA verification
      await new Promise(resolve => setTimeout(resolve, 1000));
      navigate(from, { replace: true });
    } catch (error: any) {
      setErrors('Invalid verification code. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleLogin = () => {
    // Simulate Google OAuth
    navigate(from, { replace: true });
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #0a0e13 0%, #1e2328 100%)',
        p: 2,
      }}
    >
      {/* Background decorations */}
      <Box
        sx={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          overflow: 'hidden',
          zIndex: 0,
        }}
      >
        <Box
          sx={{
            position: 'absolute',
            top: '10%',
            left: '10%',
            width: 200,
            height: 200,
            borderRadius: '50%',
            background: 'radial-gradient(circle, rgba(25,118,210,0.1) 0%, transparent 70%)',
            animation: 'pulse 4s ease-in-out infinite',
          }}
        />
        <Box
          sx={{
            position: 'absolute',
            bottom: '20%',
            right: '15%',
            width: 150,
            height: 150,
            borderRadius: '50%',
            background: 'radial-gradient(circle, rgba(200,155,60,0.1) 0%, transparent 70%)',
            animation: 'pulse 3s ease-in-out infinite reverse',
          }}
        />
      </Box>

      <Card
        sx={{
          maxWidth: 440,
          width: '100%',
          zIndex: 1,
          backdropFilter: 'blur(10px)',
          border: '1px solid rgba(255,255,255,0.1)',
        }}
      >
        <CardContent sx={{ p: 4 }}>
          {/* Logo and Header */}
          <Box sx={{ textAlign: 'center', mb: 4 }}>
            <Box
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: 1,
                mb: 2,
              }}
            >
              <SportsEsports sx={{ fontSize: 40, color: 'primary.main' }} />
              <Typography
                variant="h3"
                sx={{
                  fontWeight: 'bold',
                  background: 'linear-gradient(45deg, #1976d2, #c89b3c)',
                  backgroundClip: 'text',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                }}
              >
                Herald.lol
              </Typography>
            </Box>
            <Typography variant="h5" sx={{ fontWeight: 600, mb: 1 }}>
              {step === 'login' ? 'Welcome Back!' : 'Verify Your Identity'}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {step === 'login'
                ? 'Sign in to access your gaming analytics'
                : 'Enter the verification code from your authenticator app'}
            </Typography>
          </Box>

          {/* Features highlight */}
          {step === 'login' && (
            <Paper sx={{ p: 2, mb: 3, bgcolor: 'rgba(25,118,210,0.05)' }}>
              <Stack direction="row" spacing={2} justifyContent="center" flexWrap="wrap">
                <Chip
                  icon={<TrendingUp />}
                  label="Real-time Analytics"
                  size="small"
                  variant="outlined"
                />
                <Chip
                  icon={<Security />}
                  label="Secure Gaming Data"
                  size="small"
                  variant="outlined"
                />
                <Chip
                  icon={<SportsEsports />}
                  label="LoL Integration"
                  size="small"
                  variant="outlined"
                />
              </Stack>
            </Paper>
          )}

          {errors && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {errors}
            </Alert>
          )}

          {/* Login Form */}
          {step === 'login' && (
            <Box component="form" onSubmit={handleSubmit}>
              <TextField
                fullWidth
                name="email"
                label="Email Address"
                type="email"
                value={formData.email}
                onChange={handleChange}
                required
                autoComplete="email"
                autoFocus
                margin="normal"
                placeholder="summoner@herald.lol"
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Email color="action" />
                    </InputAdornment>
                  ),
                }}
              />

              <TextField
                fullWidth
                name="password"
                label="Password"
                type={showPassword ? 'text' : 'password'}
                value={formData.password}
                onChange={handleChange}
                required
                autoComplete="current-password"
                margin="normal"
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Lock color="action" />
                    </InputAdornment>
                  ),
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        aria-label="toggle password visibility"
                        onClick={togglePasswordVisibility}
                        edge="end"
                      >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />

              <Button
                type="submit"
                fullWidth
                variant="contained"
                disabled={isLoading}
                sx={{ mt: 3, mb: 2, py: 1.5 }}
              >
                {isLoading ? 'Signing In...' : 'Sign In'}
              </Button>
            </Box>
          )}

          {/* MFA Form */}
          {step === 'mfa' && (
            <Box component="form" onSubmit={handleMfaSubmit}>
              <TextField
                fullWidth
                label="Verification Code"
                value={mfaCode}
                onChange={(e) => {
                  const value = e.target.value.replace(/\D/g, '').slice(0, 6);
                  setMfaCode(value);
                  if (errors) setErrors('');
                }}
                margin="normal"
                required
                placeholder="123456"
                inputProps={{ maxLength: 6, style: { textAlign: 'center', fontSize: '1.5em', letterSpacing: '0.5em' } }}
              />

              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={isLoading || mfaCode.length !== 6}
                sx={{ mt: 3, mb: 2, py: 1.5 }}
              >
                {isLoading ? 'Verifying...' : 'Verify Code'}
              </Button>

              <Button
                fullWidth
                variant="text"
                onClick={() => {
                  setStep('login');
                  setMfaCode('');
                  setErrors('');
                }}
              >
                Back to Login
              </Button>
            </Box>
          )}

          {/* Google Login & Divider */}
          {step === 'login' && (
            <>
              <Divider sx={{ my: 3 }}>or</Divider>

              <Button
                fullWidth
                variant="outlined"
                size="large"
                startIcon={<Google />}
                onClick={handleGoogleLogin}
                sx={{ mb: 3, py: 1.5 }}
              >
                Continue with Google
              </Button>
            </>
          )}

          {/* Footer Links */}
          {step === 'login' && (
            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="body2" color="text.secondary">
                Don't have an account?{' '}
                <Link component={RouterLink} to="/register" color="primary">
                  Sign up for free
                </Link>
              </Typography>
              <Link
                component={RouterLink}
                to="/forgot-password"
                variant="body2"
                color="text.secondary"
                sx={{ mt: 1, display: 'block' }}
              >
                Forgot your password?
              </Link>
            </Box>
          )}
        </CardContent>
      </Card>
    </Box>
  );
};

export default LoginPage;