import React from 'react';
import { Container, Box, Typography } from '@mui/material';

export function SimpleAuthPage() {
  return (
    <Container component="main" maxWidth="sm">
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
        }}
      >
        <Typography variant="h4" align="center">
          Simple Auth Page - No Errors
        </Typography>
      </Box>
    </Container>
  );
}