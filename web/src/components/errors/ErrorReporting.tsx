import React, { useState, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  AlertTitle,
  Chip,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  IconButton,
  Tooltip,
  Snackbar,
  LinearProgress,
} from '@mui/material';
import {
  BugReport,
  Send,
  Copy,
  ExpandMore,
  Close,
  CheckCircle,
  Error as ErrorIcon,
  Info,
  Warning,
} from '@mui/icons-material';

interface ErrorReportingProps {
  error: Error;
  errorInfo?: any;
  errorId: string;
  onClose?: () => void;
  onSubmit?: (report: ErrorReport) => Promise<boolean>;
}

interface ErrorReport {
  errorId: string;
  errorMessage: string;
  errorStack?: string;
  userDescription: string;
  userEmail?: string;
  reproducibilitySteps: string;
  browserInfo: string;
  timestamp: number;
  severity: 'low' | 'medium' | 'high' | 'critical';
  category: string;
}

interface ErrorCategory {
  name: string;
  label: string;
  description: string;
  icon: React.ReactNode;
}

const ERROR_CATEGORIES: ErrorCategory[] = [
  {
    name: 'ui',
    label: 'Interface utilisateur',
    description: 'Problème d\'affichage, mise en page, composants',
    icon: <ErrorIcon />,
  },
  {
    name: 'data',
    label: 'Données',
    description: 'Erreur de chargement, parsing, filtrage des données',
    icon: <Warning />,
  },
  {
    name: 'performance',
    label: 'Performance',
    description: 'Lenteur, freeze, consommation mémoire',
    icon: <Info />,
  },
  {
    name: 'api',
    label: 'API/Réseau',
    description: 'Erreur de connexion, timeout, réponse invalide',
    icon: <BugReport />,
  },
  {
    name: 'other',
    label: 'Autre',
    description: 'Problème non catégorisé',
    icon: <ErrorIcon />,
  },
];

export const ErrorReporting: React.FC<ErrorReportingProps> = ({
  error,
  errorInfo,
  errorId,
  onClose,
  onSubmit,
}) => {
  const [open, setOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [showDetails, setShowDetails] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  
  // Form state
  const [userDescription, setUserDescription] = useState('');
  const [userEmail, setUserEmail] = useState('');
  const [reproducibilitySteps, setReproducibilitySteps] = useState('');
  const [severity, setSeverity] = useState<'low' | 'medium' | 'high' | 'critical'>('medium');
  const [category, setCategory] = useState('other');

  // Browser info
  const browserInfo = React.useMemo(() => {
    if (typeof window === 'undefined') return 'Server-side';
    
    const nav = window.navigator;
    return JSON.stringify({
      userAgent: nav.userAgent,
      platform: nav.platform,
      language: nav.language,
      cookieEnabled: nav.cookieEnabled,
      onLine: nav.onLine,
      viewport: {
        width: window.innerWidth,
        height: window.innerHeight,
      },
      screen: {
        width: window.screen.width,
        height: window.screen.height,
        colorDepth: window.screen.colorDepth,
      },
      timestamp: new Date().toISOString(),
    }, null, 2);
  }, []);

  const handleSubmit = async () => {
    if (!userDescription.trim()) {
      return;
    }

    setSubmitting(true);

    const report: ErrorReport = {
      errorId,
      errorMessage: error.message,
      errorStack: error.stack,
      userDescription: userDescription.trim(),
      userEmail: userEmail.trim() || undefined,
      reproducibilitySteps: reproducibilitySteps.trim(),
      browserInfo,
      timestamp: Date.now(),
      severity,
      category,
    };

    try {
      const success = onSubmit ? await onSubmit(report) : true;
      
      if (success) {
        setSubmitted(true);
        setTimeout(() => {
          setOpen(false);
          onClose?.();
        }, 2000);
      } else {
        throw new Error('Échec de l\'envoi du rapport');
      }
    } catch (err) {
      console.error('Error submitting report:', err);
      setSnackbarOpen(true);
    } finally {
      setSubmitting(false);
    }
  };

  const handleCopyError = useCallback(() => {
    const errorText = `
Erreur ID: ${errorId}
Message: ${error.message}
Stack: ${error.stack || 'Non disponible'}
Timestamp: ${new Date().toISOString()}
Navigateur: ${browserInfo}
    `.trim();

    navigator.clipboard?.writeText(errorText).then(() => {
      setSnackbarOpen(true);
    });
  }, [error, errorId, browserInfo]);

  const getSeverityColor = () => {
    switch (severity) {
      case 'low': return 'info';
      case 'medium': return 'warning';
      case 'high': return 'error';
      case 'critical': return 'error';
      default: return 'warning';
    }
  };

  if (submitted) {
    return (
      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="sm" fullWidth>
        <DialogContent sx={{ textAlign: 'center', py: 4 }}>
          <CheckCircle color="success" sx={{ fontSize: 64, mb: 2 }} />
          <Typography variant="h6" gutterBottom>
            Rapport envoyé avec succès
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Merci pour votre retour ! Notre équipe va examiner le problème.
          </Typography>
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <>
      <Card sx={{ mb: 2 }}>
        <CardContent>
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center" gap={1}>
              <BugReport color="error" />
              <Typography variant="h6">
                Erreur détectée
              </Typography>
              <Chip 
                label={errorId} 
                size="small" 
                variant="outlined" 
                color="error" 
              />
            </Box>
            <Box>
              <Tooltip title="Copier les détails de l'erreur">
                <IconButton onClick={handleCopyError} size="small">
                  <Copy />
                </IconButton>
              </Tooltip>
              <Button
                variant="outlined"
                startIcon={<BugReport />}
                onClick={() => setOpen(true)}
                size="small"
              >
                Signaler
              </Button>
            </Box>
          </Box>

          <Alert severity="error">
            <AlertTitle>Message d'erreur</AlertTitle>
            {error.message}
          </Alert>

          <Accordion sx={{ mt: 2 }}>
            <AccordionSummary expandIcon={<ExpandMore />}>
              <Typography variant="body2">
                Détails techniques
              </Typography>
            </AccordionSummary>
            <AccordionDetails>
              <Box
                sx={{
                  bgcolor: 'grey.100',
                  p: 2,
                  borderRadius: 1,
                  fontFamily: 'monospace',
                  fontSize: '0.875rem',
                  overflow: 'auto',
                  maxHeight: 200,
                }}
              >
                <strong>Stack Trace:</strong>
                <pre style={{ whiteSpace: 'pre-wrap', margin: 0 }}>
                  {error.stack || 'Non disponible'}
                </pre>
              </Box>
            </AccordionDetails>
          </Accordion>
        </CardContent>
      </Card>

      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          <Box display="flex" alignItems="center" justifyContent="space-between">
            <Box display="flex" alignItems="center" gap={1}>
              <BugReport />
              Signaler une erreur
            </Box>
            <IconButton onClick={() => setOpen(false)} size="small">
              <Close />
            </IconButton>
          </Box>
        </DialogTitle>

        <DialogContent>
          {submitting && <LinearProgress sx={{ mb: 2 }} />}
          
          <Box sx={{ mb: 3 }}>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Aidez-nous à résoudre ce problème en décrivant ce qui s'est passé.
            </Typography>
          </Box>

          <TextField
            fullWidth
            multiline
            rows={4}
            label="Description du problème *"
            placeholder="Décrivez ce que vous faisiez quand l'erreur s'est produite..."
            value={userDescription}
            onChange={(e) => setUserDescription(e.target.value)}
            sx={{ mb: 3 }}
            required
          />

          <TextField
            fullWidth
            multiline
            rows={3}
            label="Étapes pour reproduire"
            placeholder="1. J'ai cliqué sur...\n2. Puis j'ai...\n3. L'erreur s'est produite..."
            value={reproducibilitySteps}
            onChange={(e) => setReproducibilitySteps(e.target.value)}
            sx={{ mb: 3 }}
          />

          <TextField
            fullWidth
            label="Email (optionnel)"
            placeholder="votre@email.com"
            value={userEmail}
            onChange={(e) => setUserEmail(e.target.value)}
            helperText="Pour un suivi personnalisé du problème"
            sx={{ mb: 3 }}
          />

          <Box sx={{ mb: 3 }}>
            <Typography variant="subtitle2" gutterBottom>
              Catégorie du problème
            </Typography>
            <Box display="flex" flexWrap="wrap" gap={1}>
              {ERROR_CATEGORIES.map((cat) => (
                <Chip
                  key={cat.name}
                  label={cat.label}
                  icon={cat.icon}
                  onClick={() => setCategory(cat.name)}
                  color={category === cat.name ? 'primary' : 'default'}
                  variant={category === cat.name ? 'filled' : 'outlined'}
                />
              ))}
            </Box>
          </Box>

          <Box sx={{ mb: 3 }}>
            <Typography variant="subtitle2" gutterBottom>
              Sévérité
            </Typography>
            <Box display="flex" gap={1}>
              {[
                { value: 'low', label: 'Faible', color: 'info' },
                { value: 'medium', label: 'Moyenne', color: 'warning' },
                { value: 'high', label: 'Élevée', color: 'error' },
                { value: 'critical', label: 'Critique', color: 'error' },
              ].map((sev) => (
                <Chip
                  key={sev.value}
                  label={sev.label}
                  onClick={() => setSeverity(sev.value as any)}
                  color={severity === sev.value ? sev.color as any : 'default'}
                  variant={severity === sev.value ? 'filled' : 'outlined'}
                />
              ))}
            </Box>
          </Box>

          <Accordion>
            <AccordionSummary expandIcon={<ExpandMore />}>
              <Typography variant="body2">
                Informations techniques qui seront incluses
              </Typography>
            </AccordionSummary>
            <AccordionDetails>
              <Box
                sx={{
                  bgcolor: 'grey.50',
                  p: 2,
                  borderRadius: 1,
                  fontFamily: 'monospace',
                  fontSize: '0.75rem',
                  overflow: 'auto',
                  maxHeight: 150,
                }}
              >
                <pre style={{ whiteSpace: 'pre-wrap', margin: 0 }}>
                  {`ID: ${errorId}\nMessage: ${error.message}\n\nNavigateur:\n${browserInfo}`}
                </pre>
              </Box>
            </AccordionDetails>
          </Accordion>
        </DialogContent>

        <DialogActions>
          <Button onClick={() => setOpen(false)}>
            Annuler
          </Button>
          <Button
            variant="contained"
            startIcon={<Send />}
            onClick={handleSubmit}
            disabled={!userDescription.trim() || submitting}
          >
            {submitting ? 'Envoi...' : 'Envoyer le rapport'}
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbarOpen}
        autoHideDuration={3000}
        onClose={() => setSnackbarOpen(false)}
        message="Copié dans le presse-papiers"
      />
    </>
  );
};

// Hook pour l'utilisation simplifiée
export const useErrorReporting = () => {
  const [errors, setErrors] = useState<Array<{
    error: Error;
    errorInfo?: any;
    errorId: string;
    timestamp: number;
  }>>([]);

  const reportError = useCallback((error: Error, errorInfo?: any) => {
    const errorId = `ERR_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    setErrors(prev => [...prev, {
      error,
      errorInfo,
      errorId,
      timestamp: Date.now(),
    }]);

    return errorId;
  }, []);

  const clearError = useCallback((errorId: string) => {
    setErrors(prev => prev.filter(e => e.errorId !== errorId));
  }, []);

  const clearAllErrors = useCallback(() => {
    setErrors([]);
  }, []);

  return {
    errors,
    reportError,
    clearError,
    clearAllErrors,
  };
};

export default ErrorReporting;