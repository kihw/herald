import React, { useState, createContext, useContext, ReactNode } from 'react';
import {
  Snackbar,
  Alert,
  AlertColor,
  Slide,
  SlideProps,
  LinearProgress,
  Box,
  Typography,
} from '@mui/material';

interface ExportFeedbackContextType {
  showSuccess: (message: string) => void;
  showError: (message: string) => void;
  showProgress: (message: string, progress: number) => void;
  hideProgress: () => void;
}

const ExportFeedbackContext = createContext<ExportFeedbackContextType | undefined>(undefined);

export const useExportFeedback = () => {
  const context = useContext(ExportFeedbackContext);
  if (!context) {
    throw new Error('useExportFeedback must be used within ExportFeedbackProvider');
  }
  return context;
};

interface NotificationState {
  open: boolean;
  message: string;
  severity: AlertColor;
}

interface ProgressState {
  open: boolean;
  message: string;
  progress: number;
}

function SlideTransition(props: SlideProps) {
  return <Slide {...props} direction="up" />;
}

interface ExportFeedbackProviderProps {
  children: ReactNode;
}

export const ExportFeedbackProvider: React.FC<ExportFeedbackProviderProps> = ({ children }) => {
  const [notification, setNotification] = useState<NotificationState>({
    open: false,
    message: '',
    severity: 'success',
  });

  const [progress, setProgress] = useState<ProgressState>({
    open: false,
    message: '',
    progress: 0,
  });

  const showSuccess = (message: string) => {
    setNotification({
      open: true,
      message,
      severity: 'success',
    });
  };

  const showError = (message: string) => {
    setNotification({
      open: true,
      message,
      severity: 'error',
    });
  };

  const showProgress = (message: string, progressValue: number) => {
    setProgress({
      open: true,
      message,
      progress: progressValue,
    });
  };

  const hideProgress = () => {
    setProgress(prev => ({ ...prev, open: false }));
  };

  const handleNotificationClose = () => {
    setNotification(prev => ({ ...prev, open: false }));
  };

  const handleProgressClose = () => {
    setProgress(prev => ({ ...prev, open: false }));
  };

  return (
    <ExportFeedbackContext.Provider 
      value={{ showSuccess, showError, showProgress, hideProgress }}
    >
      {children}
      
      {/* Notification de succès/erreur */}
      <Snackbar
        open={notification.open}
        autoHideDuration={4000}
        onClose={handleNotificationClose}
        TransitionComponent={SlideTransition}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert
          onClose={handleNotificationClose}
          severity={notification.severity}
          variant="filled"
          sx={{ width: '100%' }}
        >
          {notification.message}
        </Alert>
      </Snackbar>

      {/* Indicateur de progression */}
      <Snackbar
        open={progress.open}
        TransitionComponent={SlideTransition}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert
          severity="info"
          variant="filled"
          sx={{ 
            width: '100%',
            minWidth: 300,
          }}
        >
          <Box>
            <Typography variant="body2" sx={{ mb: 1 }}>
              {progress.message}
            </Typography>
            <LinearProgress 
              variant="determinate" 
              value={progress.progress}
              sx={{
                backgroundColor: 'rgba(255, 255, 255, 0.2)',
                '& .MuiLinearProgress-bar': {
                  backgroundColor: 'rgba(255, 255, 255, 0.8)',
                },
              }}
            />
            <Typography variant="caption" sx={{ mt: 0.5, display: 'block' }}>
              {progress.progress.toFixed(0)}%
            </Typography>
          </Box>
        </Alert>
      </Snackbar>
    </ExportFeedbackContext.Provider>
  );
};