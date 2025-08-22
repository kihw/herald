import React from 'react';
import {
  Snackbar,
  Alert,
  Typography,
  Slide,
  SlideProps,
} from '@mui/material';

interface ToastNotificationProps {
  open: boolean;
  message: string;
  severity: 'success' | 'error' | 'warning' | 'info';
  title?: string;
  duration?: number;
  onClose: () => void;
}

function SlideTransition(props: SlideProps) {
  return <Slide {...props} direction="down" />;
}

const ToastNotification: React.FC<ToastNotificationProps> = ({
  open,
  message,
  severity,
  title,
  duration = 6000,
  onClose,
}) => {
  return (
    <Snackbar
      open={open}
      autoHideDuration={duration}
      onClose={onClose}
      TransitionComponent={SlideTransition}
      anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      sx={{ marginTop: 8 }}
    >
      <Alert
        onClose={onClose}
        severity={severity}
        variant="filled"
        sx={{
          width: '100%',
          minWidth: 300,
          boxShadow: 3,
        }}
      >
        {title && <Typography variant="h6" component="div" sx={{ fontWeight: 'bold' }}>{title}</Typography>}
        {message}
      </Alert>
    </Snackbar>
  );
};

export default ToastNotification;
