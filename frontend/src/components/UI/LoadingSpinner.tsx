import { Box, CircularProgress, Typography } from '@mui/material';
import { SxProps, Theme } from '@mui/material/styles';

interface LoadingSpinnerProps {
  size?: number;
  message?: string;
  color?: 'primary' | 'secondary' | 'inherit';
  sx?: SxProps<Theme>;
}

const LoadingSpinner = ({ 
  size = 40, 
  message = 'Loading...', 
  color = 'primary',
  sx = {}
}: LoadingSpinnerProps) => {
  return (
    <Box
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      gap={2}
      sx={sx}
    >
      <CircularProgress 
        size={size} 
        color={color}
        thickness={4}
      />
      {message && (
        <Typography variant="body2" color="text.secondary">
          {message}
        </Typography>
      )}
    </Box>
  );
};

export default LoadingSpinner;