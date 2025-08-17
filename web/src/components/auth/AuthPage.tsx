import React from 'react';
import { Container, Box } from '@mui/material';
import { RiotValidationForm } from './RiotValidationForm';

export function AuthPage() {
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
        <RiotValidationForm />
      </Box>
    </Container>
  );
}
