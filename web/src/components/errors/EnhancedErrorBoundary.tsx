import React, { Component, ReactNode, ErrorInfo } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Alert,
  AlertTitle,
  Collapse,
  IconButton,
  LinearProgress,
  Chip,
  Divider,
  Fade,
  useTheme,
} from '@mui/material';
import {
  ErrorOutline,
  Refresh,
  ExpandMore,
  ExpandLess,
  BugReport,
  RestartAlt,
  Warning,
  Info,
} from '@mui/icons-material';

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
  retryCount: number;
  isRetrying: boolean;
  showDetails: boolean;
  errorId: string;
}

interface EnhancedErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  maxRetries?: number;
  retryDelay?: number;
  autoRetry?: boolean;
  onError?: (error: Error, errorInfo: ErrorInfo, errorId: string) => void;
  onRetry?: (retryCount: number) => void;
  onMaxRetriesReached?: (error: Error) => void;
  enableErrorReporting?: boolean;
  showErrorDetails?: boolean;
  errorLevel?: 'low' | 'medium' | 'high' | 'critical';
}

class EnhancedErrorBoundaryClass extends Component<
  EnhancedErrorBoundaryProps,
  ErrorBoundaryState
> {
  private retryTimeout: NodeJS.Timeout | null = null;

  constructor(props: EnhancedErrorBoundaryProps) {
    super(props);
    
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      retryCount: 0,
      isRetrying: false,
      showDetails: false,
      errorId: '',
    };
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    const errorId = `ERR_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    return {
      hasError: true,
      error,
      errorId,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    this.setState({
      error,
      errorInfo,
    });

    // Log error
    console.error('Enhanced Error Boundary caught an error:', error, errorInfo);

    // Call onError callback
    if (this.props.onError) {
      this.props.onError(error, errorInfo, this.state.errorId);
    }

    // Auto retry if enabled
    if (this.props.autoRetry && this.state.retryCount < (this.props.maxRetries || 3)) {
      this.handleAutoRetry();
    }
  }

  componentWillUnmount() {
    if (this.retryTimeout) {
      clearTimeout(this.retryTimeout);
    }
  }

  handleAutoRetry = () => {
    const delay = this.props.retryDelay || 1000;
    
    this.setState({ isRetrying: true });
    
    this.retryTimeout = setTimeout(() => {
      this.handleRetry();
    }, delay);
  };

  handleRetry = () => {
    const newRetryCount = this.state.retryCount + 1;
    const maxRetries = this.props.maxRetries || 3;

    if (newRetryCount >= maxRetries) {
      if (this.props.onMaxRetriesReached) {
        this.props.onMaxRetriesReached(this.state.error!);
      }
      this.setState({ isRetrying: false });
      return;
    }

    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      retryCount: newRetryCount,
      isRetrying: false,
    });

    if (this.props.onRetry) {
      this.props.onRetry(newRetryCount);
    }
  };

  handleManualRetry = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      isRetrying: false,
    });

    if (this.props.onRetry) {
      this.props.onRetry(this.state.retryCount);
    }
  };

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      retryCount: 0,
      isRetrying: false,
      showDetails: false,
      errorId: '',
    });
  };

  toggleDetails = () => {
    this.setState(prev => ({ showDetails: !prev.showDetails }));
  };

  getErrorSeverity = () => {
    const { error } = this.state;
    const { errorLevel } = this.props;
    
    if (errorLevel) return errorLevel;
    
    // Déterminer automatiquement la sévérité
    if (!error) return 'low';
    
    if (error.name === 'ChunkLoadError' || error.message.includes('Loading chunk')) {
      return 'medium';
    }
    
    if (error.name === 'TypeError' || error.name === 'ReferenceError') {
      return 'high';
    }
    
    if (error.message.includes('Network') || error.message.includes('fetch')) {
      return 'medium';
    }
    
    return 'high';
  };

  getSeverityColor = () => {
    const severity = this.getErrorSeverity();
    
    switch (severity) {
      case 'low': return 'info';
      case 'medium': return 'warning';
      case 'high': return 'error';
      case 'critical': return 'error';
      default: return 'error';
    }
  };

  getSeverityIcon = () => {
    const severity = this.getErrorSeverity();
    
    switch (severity) {
      case 'low': return <Info />;
      case 'medium': return <Warning />;
      case 'high': return <ErrorOutline />;
      case 'critical': return <BugReport />;
      default: return <ErrorOutline />;
    }
  };

  getErrorMessage = () => {
    const { error } = this.state;
    if (!error) return 'Une erreur inattendue s\'est produite';
    
    // Messages d'erreur personnalisés
    if (error.name === 'ChunkLoadError') {
      return 'Erreur de chargement des ressources. Veuillez rafraîchir la page.';
    }
    
    if (error.message.includes('Network')) {
      return 'Erreur de connexion réseau. Vérifiez votre connexion internet.';
    }
    
    if (error.message.includes('fetch')) {
      return 'Erreur lors de la récupération des données. Veuillez réessayer.';
    }
    
    return error.message || 'Une erreur inattendue s\'est produite';
  };

  render() {
    if (this.state.hasError) {
      const { maxRetries = 3, showErrorDetails = true } = this.props;
      const severity = this.getErrorSeverity();
      const severityColor = this.getSeverityColor();
      
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <Fade in={true}>
          <Box sx={{ p: 3 }}>
            <Card elevation={3}>
              <CardContent>
                <Box display="flex" alignItems="flex-start" gap={2} mb={3}>
                  {this.getSeverityIcon()}
                  <Box flex={1}>
                    <Typography variant="h6" color={`${severityColor}.main`} gutterBottom>
                      Erreur de l'application
                    </Typography>
                    <Typography variant="body1" color="text.secondary" paragraph>
                      {this.getErrorMessage()}
                    </Typography>
                    
                    {/* Informations sur les tentatives */}
                    <Box display="flex" alignItems="center" gap={1} mb={2}>
                      <Chip
                        label={`Tentative ${this.state.retryCount}/${maxRetries}`}
                        size="small"
                        color={this.state.retryCount >= maxRetries ? 'error' : 'default'}
                      />
                      <Chip
                        label={`ID: ${this.state.errorId}`}
                        size="small"
                        variant="outlined"
                      />
                      <Chip
                        label={severity.toUpperCase()}
                        size="small"
                        color={severityColor}
                      />
                    </Box>
                  </Box>
                </Box>

                {/* Barre de progression pour retry automatique */}
                {this.state.isRetrying && (
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      Tentative de récupération automatique...
                    </Typography>
                    <LinearProgress />
                  </Box>
                )}

                {/* Actions */}
                <Box display="flex" gap={2} mb={showErrorDetails ? 2 : 0}>
                  <Button
                    variant="contained"
                    startIcon={<Refresh />}
                    onClick={this.handleManualRetry}
                    disabled={this.state.isRetrying}
                  >
                    Réessayer
                  </Button>
                  
                  <Button
                    variant="outlined"
                    startIcon={<RestartAlt />}
                    onClick={this.handleReset}
                  >
                    Réinitialiser
                  </Button>
                  
                  <Button
                    variant="text"
                    onClick={() => window.location.reload()}
                  >
                    Rafraîchir la page
                  </Button>

                  {showErrorDetails && (
                    <Button
                      variant="text"
                      endIcon={this.state.showDetails ? <ExpandLess /> : <ExpandMore />}
                      onClick={this.toggleDetails}
                      sx={{ ml: 'auto' }}
                    >
                      Détails techniques
                    </Button>
                  )}
                </Box>

                {/* Détails de l'erreur */}
                {showErrorDetails && (
                  <Collapse in={this.state.showDetails}>
                    <Divider sx={{ my: 2 }} />
                    <Alert severity="warning" sx={{ mb: 2 }}>
                      <AlertTitle>Informations de débogage</AlertTitle>
                      Ces informations peuvent aider les développeurs à identifier le problème.
                    </Alert>
                    
                    {this.state.error && (
                      <Box sx={{ mb: 2 }}>
                        <Typography variant="subtitle2" gutterBottom>
                          Erreur:
                        </Typography>
                        <Box
                          sx={{
                            p: 2,
                            bgcolor: 'grey.100',
                            borderRadius: 1,
                            fontFamily: 'monospace',
                            fontSize: '0.875rem',
                            overflow: 'auto',
                            maxHeight: 200,
                          }}
                        >
                          <strong>{this.state.error.name}:</strong> {this.state.error.message}
                          {this.state.error.stack && (
                            <pre style={{ whiteSpace: 'pre-wrap', margin: 0 }}>
                              {this.state.error.stack}
                            </pre>
                          )}
                        </Box>
                      </Box>
                    )}

                    {this.state.errorInfo && (
                      <Box>
                        <Typography variant="subtitle2" gutterBottom>
                          Stack de composants:
                        </Typography>
                        <Box
                          sx={{
                            p: 2,
                            bgcolor: 'grey.100',
                            borderRadius: 1,
                            fontFamily: 'monospace',
                            fontSize: '0.875rem',
                            overflow: 'auto',
                            maxHeight: 200,
                          }}
                        >
                          <pre style={{ whiteSpace: 'pre-wrap', margin: 0 }}>
                            {this.state.errorInfo.componentStack}
                          </pre>
                        </Box>
                      </Box>
                    )}
                  </Collapse>
                )}
              </CardContent>
            </Card>
          </Box>
        </Fade>
      );
    }

    return this.props.children;
  }
}

// Wrapper fonctionnel pour une utilisation plus facile
export const EnhancedErrorBoundary: React.FC<EnhancedErrorBoundaryProps> = (props) => {
  return <EnhancedErrorBoundaryClass {...props} />;
};

// Hook pour la gestion d'erreurs dans les composants fonctionnels
export const useErrorHandler = () => {
  const [error, setError] = React.useState<Error | null>(null);

  const handleError = React.useCallback((error: Error) => {
    console.error('Error caught by useErrorHandler:', error);
    setError(error);
  }, []);

  const clearError = React.useCallback(() => {
    setError(null);
  }, []);

  // Relancer l'erreur pour qu'elle soit capturée par l'Error Boundary
  React.useEffect(() => {
    if (error) {
      throw error;
    }
  }, [error]);

  return { handleError, clearError, error };
};

// Composant d'erreur pour les erreurs spécifiques
export const ErrorDisplay: React.FC<{
  error: Error;
  onRetry?: () => void;
  onDismiss?: () => void;
  severity?: 'error' | 'warning' | 'info';
}> = ({ error, onRetry, onDismiss, severity = 'error' }) => {
  return (
    <Alert 
      severity={severity}
      action={
        <Box>
          {onRetry && (
            <Button color="inherit" size="small" onClick={onRetry}>
              Réessayer
            </Button>
          )}
          {onDismiss && (
            <Button color="inherit" size="small" onClick={onDismiss}>
              Fermer
            </Button>
          )}
        </Box>
      }
    >
      <AlertTitle>Erreur</AlertTitle>
      {error.message}
    </Alert>
  );
};

export default EnhancedErrorBoundary;