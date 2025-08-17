import React, { Component, ErrorInfo, ReactNode } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Alert,
  Collapse,
  IconButton,
  Stack,
  Chip
} from '@mui/material';
import {
  ErrorOutline,
  Refresh,
  BugReport,
  ExpandMore,
  ExpandLess,
  ContentCopy
} from '@mui/icons-material';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
  expanded: boolean;
  retryCount: number;
}

/**
 * Enhanced Error Boundary with detailed error reporting and recovery options
 */
class ErrorBoundary extends Component<Props, State> {
  private maxRetries = 3;

  constructor(props: Props) {
    super(props);
    
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      expanded: false,
      retryCount: 0
    };
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    return {
      hasError: true,
      error
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Log error to console in development
    if (process.env.NODE_ENV === 'development') {
      console.error('🚨 Error Boundary caught an error:', error);
      console.error('Error Info:', errorInfo);
    }

    this.setState({
      error,
      errorInfo
    });

    // Call custom error handler if provided
    this.props.onError?.(error, errorInfo);

    // Log to external service in production
    this.logErrorToService(error, errorInfo);
  }

  private logErrorToService = (error: Error, errorInfo: ErrorInfo) => {
    // In a real app, you would send this to your error tracking service
    // like Sentry, LogRocket, etc.
    if (process.env.NODE_ENV === 'production') {
      // Example: Sentry.captureException(error, { extra: errorInfo });
      console.warn('Error would be logged to external service:', {
        message: error.message,
        stack: error.stack,
        componentStack: errorInfo.componentStack
      });
    }
  };

  private handleRetry = () => {
    if (this.state.retryCount < this.maxRetries) {
      this.setState(prevState => ({
        hasError: false,
        error: null,
        errorInfo: null,
        retryCount: prevState.retryCount + 1
      }));
    }
  };

  private handleReload = () => {
    window.location.reload();
  };

  private toggleExpanded = () => {
    this.setState(prevState => ({
      expanded: !prevState.expanded
    }));
  };

  private copyErrorToClipboard = async () => {
    const { error, errorInfo } = this.state;
    if (!error) return;

    const errorText = `
Error: ${error.message}

Stack Trace:
${error.stack}

Component Stack:
${errorInfo?.componentStack}

User Agent: ${navigator.userAgent}
Timestamp: ${new Date().toISOString()}
URL: ${window.location.href}
    `.trim();

    try {
      await navigator.clipboard.writeText(errorText);
      // Could show a toast notification here
    } catch (err) {
      console.error('Failed to copy error to clipboard:', err);
    }
  };

  private getErrorCategory = (error: Error): string => {
    const message = error.message.toLowerCase();
    
    if (message.includes('network') || message.includes('fetch')) {
      return 'Network Error';
    }
    
    if (message.includes('chunk') || message.includes('loading')) {
      return 'Loading Error';
    }
    
    if (message.includes('permission') || message.includes('unauthorized')) {
      return 'Permission Error';
    }
    
    if (message.includes('timeout')) {
      return 'Timeout Error';
    }
    
    return 'Application Error';
  };

  private getSuggestedActions = (error: Error): string[] => {
    const message = error.message.toLowerCase();
    const suggestions: string[] = [];
    
    if (message.includes('network') || message.includes('fetch')) {
      suggestions.push('Check your internet connection');
      suggestions.push('Try refreshing the page');
    }
    
    if (message.includes('chunk') || message.includes('loading')) {
      suggestions.push('Clear your browser cache');
      suggestions.push('Try a hard refresh (Ctrl+F5)');
    }
    
    if (message.includes('timeout')) {
      suggestions.push('Wait a moment and try again');
      suggestions.push('Check if the server is responding');
    }
    
    if (suggestions.length === 0) {
      suggestions.push('Try refreshing the page');
      suggestions.push('Clear browser cache and cookies');
      suggestions.push('Contact support if the issue persists');
    }
    
    return suggestions;
  };

  render() {
    if (this.state.hasError) {
      // Custom fallback if provided
      if (this.props.fallback) {
        return this.props.fallback;
      }

      const { error, errorInfo, expanded, retryCount } = this.state;
      const canRetry = retryCount < this.maxRetries;
      const errorCategory = error ? this.getErrorCategory(error) : 'Unknown Error';
      const suggestions = error ? this.getSuggestedActions(error) : [];

      return (
        <Box
          display="flex"
          justifyContent="center"
          alignItems="center"
          minHeight="400px"
          p={3}
        >
          <Card sx={{ maxWidth: 600, width: '100%' }}>
            <CardContent>
              {/* Error Header */}
              <Stack direction="row" alignItems="center" spacing={2} mb={2}>
                <ErrorOutline color="error" fontSize="large" />
                <Box>
                  <Typography variant="h5" color="error" gutterBottom>
                    Oops! Something went wrong
                  </Typography>
                  <Chip 
                    label={errorCategory} 
                    color="error" 
                    variant="outlined" 
                    size="small" 
                  />
                </Box>
              </Stack>

              {/* Error Message */}
              <Alert severity="error" sx={{ mb: 2 }}>
                <Typography variant="body1">
                  {error?.message || 'An unexpected error occurred'}
                </Typography>
              </Alert>

              {/* Suggested Actions */}
              {suggestions.length > 0 && (
                <Box mb={2}>
                  <Typography variant="subtitle2" gutterBottom>
                    Try these solutions:
                  </Typography>
                  <Box component="ul" sx={{ pl: 2, mb: 0 }}>
                    {suggestions.map((suggestion, index) => (
                      <Typography 
                        key={index}
                        component="li" 
                        variant="body2" 
                        color="text.secondary"
                      >
                        {suggestion}
                      </Typography>
                    ))}
                  </Box>
                </Box>
              )}

              {/* Action Buttons */}
              <Stack direction="row" spacing={2} mb={2}>
                {canRetry && (
                  <Button
                    variant="contained"
                    startIcon={<Refresh />}
                    onClick={this.handleRetry}
                    color="primary"
                  >
                    Try Again {retryCount > 0 && `(${retryCount}/${this.maxRetries})`}
                  </Button>
                )}
                
                <Button
                  variant="outlined"
                  startIcon={<Refresh />}
                  onClick={this.handleReload}
                >
                  Reload Page
                </Button>

                <Button
                  variant="text"
                  startIcon={<ContentCopy />}
                  onClick={this.copyErrorToClipboard}
                  size="small"
                >
                  Copy Error
                </Button>
              </Stack>

              {/* Technical Details (Collapsible) */}
              {process.env.NODE_ENV === 'development' && error && (
                <>
                  <Button
                    variant="text"
                    startIcon={<BugReport />}
                    endIcon={expanded ? <ExpandLess /> : <ExpandMore />}
                    onClick={this.toggleExpanded}
                    size="small"
                    sx={{ mb: 1 }}
                  >
                    Technical Details
                  </Button>
                  
                  <Collapse in={expanded}>
                    <Alert severity="info" sx={{ mb: 2 }}>
                      <Typography variant="caption" component="div">
                        <strong>Error Stack:</strong>
                      </Typography>
                      <Box
                        component="pre"
                        sx={{
                          fontSize: '0.75rem',
                          whiteSpace: 'pre-wrap',
                          overflow: 'auto',
                          maxHeight: 200,
                          mt: 1,
                          p: 1,
                          bgcolor: 'background.paper',
                          border: '1px solid',
                          borderColor: 'divider',
                          borderRadius: 1
                        }}
                      >
                        {error.stack}
                      </Box>
                      
                      {errorInfo?.componentStack && (
                        <>
                          <Typography variant="caption" component="div" sx={{ mt: 2 }}>
                            <strong>Component Stack:</strong>
                          </Typography>
                          <Box
                            component="pre"
                            sx={{
                              fontSize: '0.75rem',
                              whiteSpace: 'pre-wrap',
                              overflow: 'auto',
                              maxHeight: 200,
                              mt: 1,
                              p: 1,
                              bgcolor: 'background.paper',
                              border: '1px solid',
                              borderColor: 'divider',
                              borderRadius: 1
                            }}
                          >
                            {errorInfo.componentStack}
                          </Box>
                        </>
                      )}
                    </Alert>
                  </Collapse>
                </>
              )}

              {/* Contact Support */}
              <Typography variant="caption" color="text.secondary">
                If the problem persists, please contact support with the error details above.
              </Typography>
            </CardContent>
          </Card>
        </Box>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
