import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { Alert, AlertTitle, Snackbar, Box, Fab, Badge } from '@mui/material';
import { BugReport, Close } from '@mui/icons-material';
import { ErrorReporting, useErrorReporting } from './ErrorReporting';

interface GlobalError {
  id: string;
  error: Error;
  errorInfo?: any;
  context?: string;
  timestamp: number;
  severity: 'low' | 'medium' | 'high' | 'critical';
  resolved: boolean;
}

interface GlobalErrorContextType {
  errors: GlobalError[];
  reportError: (error: Error, context?: string, severity?: 'low' | 'medium' | 'high' | 'critical') => string;
  resolveError: (errorId: string) => void;
  clearAllErrors: () => void;
  unreadCount: number;
}

const GlobalErrorContext = createContext<GlobalErrorContextType | undefined>(undefined);

interface GlobalErrorProviderProps {
  children: ReactNode;
  enableNotifications?: boolean;
  enableFloatingButton?: boolean;
  maxErrors?: number;
  onError?: (error: GlobalError) => void;
}

export const GlobalErrorProvider: React.FC<GlobalErrorProviderProps> = ({
  children,
  enableNotifications = true,
  enableFloatingButton = true,
  maxErrors = 10,
  onError,
}) => {
  const [errors, setErrors] = useState<GlobalError[]>([]);
  const [showReporting, setShowReporting] = useState(false);
  const [selectedError, setSelectedError] = useState<GlobalError | null>(null);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [latestError, setLatestError] = useState<GlobalError | null>(null);

  const reportError = useCallback((
    error: Error, 
    context?: string, 
    severity: 'low' | 'medium' | 'high' | 'critical' = 'medium'
  ) => {
    const errorId = `GLOBAL_ERR_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    const globalError: GlobalError = {
      id: errorId,
      error,
      context,
      timestamp: Date.now(),
      severity,
      resolved: false,
    };

    setErrors(prev => {
      const newErrors = [globalError, ...prev].slice(0, maxErrors);
      return newErrors;
    });

    // Notification pour les erreurs importantes
    if (enableNotifications && (severity === 'high' || severity === 'critical')) {
      setLatestError(globalError);
      setSnackbarOpen(true);
    }

    // Callback externe
    onError?.(globalError);

    // Log pour debugging
    console.error('Global Error Reported:', {
      id: errorId,
      message: error.message,
      context,
      severity,
      stack: error.stack,
    });

    return errorId;
  }, [maxErrors, enableNotifications, onError]);

  const resolveError = useCallback((errorId: string) => {
    setErrors(prev => 
      prev.map(err => 
        err.id === errorId 
          ? { ...err, resolved: true }
          : err
      )
    );
  }, []);

  const clearAllErrors = useCallback(() => {
    setErrors([]);
    setSnackbarOpen(false);
    setShowReporting(false);
    setSelectedError(null);
  }, []);

  const unreadCount = errors.filter(err => !err.resolved).length;

  const handleOpenReporting = (error?: GlobalError) => {
    if (error) {
      setSelectedError(error);
    } else if (errors.length > 0) {
      setSelectedError(errors[0]);
    }
    setShowReporting(true);
  };

  const handleCloseReporting = () => {
    setShowReporting(false);
    setSelectedError(null);
  };

  const handleReportSubmit = async (report: any): Promise<boolean> => {
    try {
      // Ici, vous pouvez envoyer le rapport à votre service de logging
      console.log('Error report submitted:', report);
      
      // Simuler l'envoi
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Marquer l'erreur comme résolue
      if (selectedError) {
        resolveError(selectedError.id);
      }
      
      return true;
    } catch (err) {
      console.error('Failed to submit error report:', err);
      return false;
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'low': return 'info';
      case 'medium': return 'warning';
      case 'high': return 'error';
      case 'critical': return 'error';
      default: return 'warning';
    }
  };

  const contextValue: GlobalErrorContextType = {
    errors,
    reportError,
    resolveError,
    clearAllErrors,
    unreadCount,
  };

  return (
    <GlobalErrorContext.Provider value={contextValue}>
      {children}

      {/* Notification Snackbar pour erreurs importantes */}
      {enableNotifications && latestError && (
        <Snackbar
          open={snackbarOpen}
          autoHideDuration={6000}
          onClose={() => setSnackbarOpen(false)}
          anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
        >
          <Alert 
            severity={getSeverityColor(latestError.severity) as any}
            onClose={() => setSnackbarOpen(false)}
            action={
              <Box>
                <button
                  onClick={() => {
                    handleOpenReporting(latestError);
                    setSnackbarOpen(false);
                  }}
                  style={{
                    background: 'none',
                    border: 'none',
                    color: 'inherit',
                    cursor: 'pointer',
                    textDecoration: 'underline',
                    padding: 0,
                    margin: '0 8px',
                  }}
                >
                  Signaler
                </button>
              </Box>
            }
          >
            <AlertTitle>
              Erreur {latestError.severity === 'critical' ? 'critique' : 'importante'} détectée
            </AlertTitle>
            {latestError.context && `${latestError.context}: `}
            {latestError.error.message}
          </Alert>
        </Snackbar>
      )}

      {/* Bouton flottant pour les erreurs */}
      {enableFloatingButton && unreadCount > 0 && (
        <Fab
          color="error"
          sx={{
            position: 'fixed',
            bottom: 16,
            right: 16,
            zIndex: 1000,
          }}
          onClick={() => handleOpenReporting()}
        >
          <Badge badgeContent={unreadCount} color="warning">
            <BugReport />
          </Badge>
        </Fab>
      )}

      {/* Interface de reporting */}
      {showReporting && selectedError && (
        <Box
          sx={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            bgcolor: 'rgba(0, 0, 0, 0.5)',
            zIndex: 2000,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            p: 2,
          }}
        >
          <Box
            sx={{
              bgcolor: 'background.paper',
              borderRadius: 2,
              maxWidth: 800,
              width: '100%',
              maxHeight: '90vh',
              overflow: 'auto',
              p: 3,
            }}
          >
            <ErrorReporting
              error={selectedError.error}
              errorInfo={selectedError}
              errorId={selectedError.id}
              onClose={handleCloseReporting}
              onSubmit={handleReportSubmit}
            />
          </Box>
        </Box>
      )}
    </GlobalErrorContext.Provider>
  );
};

// Hook pour utiliser le contexte d'erreurs globales
export const useGlobalError = () => {
  const context = useContext(GlobalErrorContext);
  if (!context) {
    throw new Error('useGlobalError must be used within a GlobalErrorProvider');
  }
  return context;
};

// Hook pour signaler facilement des erreurs
export const useErrorReporter = () => {
  const { reportError } = useGlobalError();

  const reportAPIError = useCallback((error: Error, endpoint?: string) => {
    return reportError(error, `API: ${endpoint}`, 'medium');
  }, [reportError]);

  const reportUIError = useCallback((error: Error, component?: string) => {
    return reportError(error, `UI: ${component}`, 'low');
  }, [reportError]);

  const reportDataError = useCallback((error: Error, operation?: string) => {
    return reportError(error, `Data: ${operation}`, 'high');
  }, [reportError]);

  const reportCriticalError = useCallback((error: Error, context?: string) => {
    return reportError(error, context, 'critical');
  }, [reportError]);

  return {
    reportError,
    reportAPIError,
    reportUIError,
    reportDataError,
    reportCriticalError,
  };
};

// HOC pour capturer automatiquement les erreurs de composants
export const withErrorReporting = <P extends object>(
  Component: React.ComponentType<P>,
  context?: string
) => {
  return React.forwardRef<any, P>((props, ref) => {
    const { reportError } = useGlobalError();

    React.useEffect(() => {
      const handleError = (event: ErrorEvent) => {
        reportError(
          new Error(event.message),
          context || Component.displayName || Component.name,
          'medium'
        );
      };

      const handleUnhandledRejection = (event: PromiseRejectionEvent) => {
        reportError(
          new Error(event.reason?.message || 'Unhandled Promise Rejection'),
          context || Component.displayName || Component.name,
          'high'
        );
      };

      window.addEventListener('error', handleError);
      window.addEventListener('unhandledrejection', handleUnhandledRejection);

      return () => {
        window.removeEventListener('error', handleError);
        window.removeEventListener('unhandledrejection', handleUnhandledRejection);
      };
    }, [reportError]);

    return <Component {...props} ref={ref} />;
  });
};

export default GlobalErrorProvider;