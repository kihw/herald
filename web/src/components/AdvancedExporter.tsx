import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Button,
  Box,
  Typography,
  LinearProgress,
  Alert,
  Checkbox,
  FormControlLabel,
  Grid,
  Chip,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Switch,
  FormGroup,
  RadioGroup,
  Radio,
  FormLabel,
  Slider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  List,
  ListItem,
  ListItemText,
  IconButton,
} from '@mui/material';
import {
  ExpandMore,
  Settings,
  Download,
  PlayArrow,
  Stop,
  History,
  FilterList,
} from '@mui/icons-material';
// Date pickers simplifiés avec input HTML

// Types pour les options d'export avancé
export interface ExportFilter {
  date_from?: string;
  date_to?: string;
  queues?: number[];
  champions?: number[];
  min_duration?: number;
  max_duration?: number;
  win_only?: boolean;
  ranked_only?: boolean;
  recent_first: boolean;
  include_remake: boolean;
}

export interface ExportOptions {
  format: 'csv' | 'json' | 'parquet' | 'xlsx';
  filter: ExportFilter;
  columns?: string[];
  filename: string;
  compression: boolean;
  metadata: boolean;
}

interface AdvancedExportRequest {
  username: string;
  tagline: string;
  options: ExportOptions;
}

interface JobStatus {
  id: string;
  status: string;
  progress: number;
  log_lines: string[];
  start_time: string;
  end_time?: string;
  error?: string;
  zip_path?: string;
}

interface AdvancedExporterProps {
  onLoadingChange?: (loading: boolean) => void;
  onErrorChange?: (error: string | null) => void;
}

const QUEUE_OPTIONS = [
  { id: 420, name: 'Ranked Solo/Duo' },
  { id: 440, name: 'Ranked Flex' },
  { id: 400, name: 'Normal Draft' },
  { id: 430, name: 'Normal Blind' },
  { id: 450, name: 'ARAM' },
];

const AVAILABLE_COLUMNS = [
  'match_id', 'game_creation', 'game_duration', 'queue_id', 'game_mode',
  'champion_name', 'role', 'lane', 'win', 'kills', 'deaths', 'assists',
  'kda', 'cs', 'gold', 'damage', 'vision', 'rank', 'lp', 'mmr'
];

export const AdvancedExporter: React.FC<AdvancedExporterProps> = ({
  onLoadingChange,
  onErrorChange,
}) => {
  const [username, setUsername] = useState('');
  const [tagline, setTagline] = useState('');
  const [currentJob, setCurrentJob] = useState<JobStatus | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [supportedFormats, setSupportedFormats] = useState<string[]>([]);
  const [exportHistory, setExportHistory] = useState<any[]>([]);
  const [showHistory, setShowHistory] = useState(false);
  const [templates, setTemplates] = useState<any[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<string>('');
  
  // Options d'export
  const [exportOptions, setExportOptions] = useState<ExportOptions>({
    format: 'csv',
    filter: {
      recent_first: true,
      include_remake: false,
    },
    filename: '',
    compression: false,
    metadata: true,
  });

  // Filtres de date
  const [dateFrom, setDateFrom] = useState<Date | null>(null);
  const [dateTo, setDateTo] = useState<Date | null>(null);

  // Colonnes sélectionnées
  const [selectedColumns, setSelectedColumns] = useState<string[]>([]);

  useEffect(() => {
    loadSupportedFormats();
    loadExportHistory();
  }, []);

  useEffect(() => {
    if (onLoadingChange) {
      onLoadingChange(isLoading);
    }
  }, [isLoading, onLoadingChange]);

  useEffect(() => {
    if (onErrorChange) {
      onErrorChange(error);
    }
  }, [error, onErrorChange]);

  const loadSupportedFormats = async () => {
    try {
      const response = await fetch('/api/export/formats');
      const data = await response.json();
      setSupportedFormats(data.formats || []);
    } catch (err) {
      console.error('Erreur lors du chargement des formats:', err);
    }
  };

  const loadExportHistory = async () => {
    try {
      const response = await fetch('/api/export/history');
      const data = await response.json();
      setExportHistory(data.history || []);
    } catch (err) {
      console.error('Erreur lors du chargement de l\'historique:', err);
    }
  };

  const validateOptions = async (): Promise<boolean> => {
    try {
      const response = await fetch('/api/export/validate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(exportOptions),
      });
      const data = await response.json();
      return data.valid;
    } catch (err) {
      console.error('Erreur lors de la validation:', err);
      return false;
    }
  };

  const startAdvancedExport = async () => {
    if (!username.trim() || !tagline.trim()) {
      setError('Nom d\'utilisateur et tagline requis');
      return;
    }

    // Préparer les options avec les filtres de date
    const finalOptions = { ...exportOptions };
    if (dateFrom) {
      finalOptions.filter.date_from = dateFrom.toISOString();
    }
    if (dateTo) {
      finalOptions.filter.date_to = dateTo.toISOString();
    }
    if (selectedColumns.length > 0) {
      finalOptions.columns = selectedColumns;
    }

    // Valider les options
    const isValid = await validateOptions();
    if (!isValid) {
      setError('Options d\'export invalides');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const request: AdvancedExportRequest = {
        username: username.trim(),
        tagline: tagline.trim(),
        options: finalOptions,
      };

      const response = await fetch('/api/export/advanced', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Erreur lors du démarrage de l\'export');
      }

      const data = await response.json();
      pollJobStatus(data.job_id);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur inconnue');
      setIsLoading(false);
    }
  };

  const pollJobStatus = async (jobId: string) => {
    const poll = async () => {
      try {
        const response = await fetch(`/api/jobs/${jobId}`);
        const job: JobStatus = await response.json();
        setCurrentJob(job);

        if (job.status === 'completed' || job.status === 'failed' || job.status === 'cancelled') {
          setIsLoading(false);
          if (job.status === 'completed') {
            loadExportHistory(); // Rafraîchir l'historique
          }
          return;
        }

        setTimeout(poll, 2000);
      } catch (err) {
        console.error('Erreur lors de la vérification du statut:', err);
        setIsLoading(false);
      }
    };

    poll();
  };

  const handleQueueChange = (queueId: number, checked: boolean) => {
    setExportOptions(prev => ({
      ...prev,
      filter: {
        ...prev.filter,
        queues: checked
          ? [...(prev.filter.queues || []), queueId]
          : (prev.filter.queues || []).filter(id => id !== queueId)
      }
    }));
  };

  const handleColumnChange = (column: string, checked: boolean) => {
    setSelectedColumns(prev =>
      checked
        ? [...prev, column]
        : prev.filter(col => col !== column)
    );
  };

  const downloadResult = async (jobId: string) => {
    try {
      const response = await fetch(`/api/jobs/${jobId}/download`);
      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `export_${jobId}.zip`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
      }
    } catch (err) {
      console.error('Erreur lors du téléchargement:', err);
    }
  };

  return (
    <div>
      <Card>
        <CardHeader 
          title="Export Avancé" 
          action={
            <Box>
              <IconButton onClick={() => setShowHistory(true)}>
                <History />
              </IconButton>
            </Box>
          }
        />
        <CardContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <Grid container spacing={3}>
            {/* Informations utilisateur */}
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Nom d'utilisateur"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                disabled={isLoading}
                margin="normal"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Tagline"
                value={tagline}
                onChange={(e) => setTagline(e.target.value)}
                disabled={isLoading}
                margin="normal"
              />
            </Grid>

            {/* Options de format */}
            <Grid item xs={12} md={6}>
              <FormControl fullWidth margin="normal">
                <InputLabel>Format d'export</InputLabel>
                <Select
                  value={exportOptions.format}
                  onChange={(e) => setExportOptions(prev => ({
                    ...prev,
                    format: e.target.value as any
                  }))}
                  disabled={isLoading}
                >
                  {supportedFormats.map(format => (
                    <MenuItem key={format} value={format}>
                      {format.toUpperCase()}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Nom du fichier (optionnel)"
                value={exportOptions.filename}
                onChange={(e) => setExportOptions(prev => ({
                  ...prev,
                  filename: e.target.value
                }))}
                disabled={isLoading}
                margin="normal"
              />
            </Grid>

            {/* Filtres avancés */}
            <Grid item xs={12}>
              <Accordion>
                <AccordionSummary expandIcon={<ExpandMore />}>
                  <Typography>Filtres avancés</Typography>
                </AccordionSummary>
                <AccordionDetails>
                  <Grid container spacing={2}>
                    {/* Filtres de date */}
                    <Grid item xs={12} md={6}>
                      <TextField
                        fullWidth
                        label="Date de début"
                        type="date"
                        value={dateFrom ? dateFrom.toISOString().split('T')[0] : ''}
                        onChange={(e) => setDateFrom(e.target.value ? new Date(e.target.value) : null)}
                        disabled={isLoading}
                        InputLabelProps={{ shrink: true }}
                      />
                    </Grid>
                    <Grid item xs={12} md={6}>
                      <TextField
                        fullWidth
                        label="Date de fin"
                        type="date"
                        value={dateTo ? dateTo.toISOString().split('T')[0] : ''}
                        onChange={(e) => setDateTo(e.target.value ? new Date(e.target.value) : null)}
                        disabled={isLoading}
                        InputLabelProps={{ shrink: true }}
                      />
                    </Grid>

                    {/* Filtres de queue */}
                    <Grid item xs={12}>
                      <Typography variant="subtitle2" gutterBottom>
                        Types de parties
                      </Typography>
                      <FormGroup row>
                        {QUEUE_OPTIONS.map(queue => (
                          <FormControlLabel
                            key={queue.id}
                            control={
                              <Checkbox
                                checked={(exportOptions.filter.queues || []).includes(queue.id)}
                                onChange={(e) => handleQueueChange(queue.id, e.target.checked)}
                                disabled={isLoading}
                              />
                            }
                            label={queue.name}
                          />
                        ))}
                      </FormGroup>
                    </Grid>

                    {/* Options diverses */}
                    <Grid item xs={12}>
                      <FormGroup>
                        <FormControlLabel
                          control={
                            <Switch
                              checked={exportOptions.filter.win_only || false}
                              onChange={(e) => setExportOptions(prev => ({
                                ...prev,
                                filter: { ...prev.filter, win_only: e.target.checked }
                              }))}
                              disabled={isLoading}
                            />
                          }
                          label="Victoires uniquement"
                        />
                        <FormControlLabel
                          control={
                            <Switch
                              checked={exportOptions.filter.ranked_only || false}
                              onChange={(e) => setExportOptions(prev => ({
                                ...prev,
                                filter: { ...prev.filter, ranked_only: e.target.checked }
                              }))}
                              disabled={isLoading}
                            />
                          }
                          label="Parties classées uniquement"
                        />
                        <FormControlLabel
                          control={
                            <Switch
                              checked={exportOptions.filter.include_remake}
                              onChange={(e) => setExportOptions(prev => ({
                                ...prev,
                                filter: { ...prev.filter, include_remake: e.target.checked }
                              }))}
                              disabled={isLoading}
                            />
                          }
                          label="Inclure les remakes"
                        />
                      </FormGroup>
                    </Grid>
                  </Grid>
                </AccordionDetails>
              </Accordion>
            </Grid>

            {/* Sélection des colonnes */}
            <Grid item xs={12}>
              <Accordion>
                <AccordionSummary expandIcon={<ExpandMore />}>
                  <Typography>Colonnes à exporter ({selectedColumns.length} sélectionnées)</Typography>
                </AccordionSummary>
                <AccordionDetails>
                  <FormGroup>
                    <Grid container>
                      {AVAILABLE_COLUMNS.map(column => (
                        <Grid item xs={6} md={4} key={column}>
                          <FormControlLabel
                            control={
                              <Checkbox
                                checked={selectedColumns.includes(column)}
                                onChange={(e) => handleColumnChange(column, e.target.checked)}
                                disabled={isLoading}
                              />
                            }
                            label={column.replace('_', ' ')}
                          />
                        </Grid>
                      ))}
                    </Grid>
                  </FormGroup>
                </AccordionDetails>
              </Accordion>
            </Grid>

            {/* Options d'export */}
            <Grid item xs={12}>
              <FormGroup row>
                <FormControlLabel
                  control={
                    <Switch
                      checked={exportOptions.compression}
                      onChange={(e) => setExportOptions(prev => ({
                        ...prev,
                        compression: e.target.checked
                      }))}
                      disabled={isLoading}
                    />
                  }
                  label="Compression"
                />
                <FormControlLabel
                  control={
                    <Switch
                      checked={exportOptions.metadata}
                      onChange={(e) => setExportOptions(prev => ({
                        ...prev,
                        metadata: e.target.checked
                      }))}
                      disabled={isLoading}
                    />
                  }
                  label="Inclure les métadonnées"
                />
              </FormGroup>
            </Grid>
          </Grid>

          {/* Bouton de lancement */}
          <Box sx={{ mt: 3 }}>
            <Button
              variant="contained"
              size="large"
              startIcon={isLoading ? <Stop /> : <PlayArrow />}
              onClick={startAdvancedExport}
              disabled={isLoading || !username.trim() || !tagline.trim()}
              fullWidth
            >
              {isLoading ? 'Arrêter l\'export' : 'Démarrer l\'export avancé'}
            </Button>
          </Box>

          {/* Statut du job */}
          {currentJob && (
            <Box sx={{ mt: 3 }}>
              <Typography variant="h6" gutterBottom>
                Statut de l'export
              </Typography>
              <LinearProgress 
                variant="determinate" 
                value={currentJob.progress} 
                sx={{ mb: 1 }} 
              />
              <Typography variant="body2" color="text.secondary">
                {currentJob.status} - {currentJob.progress}%
              </Typography>
              
              {currentJob.status === 'completed' && currentJob.zip_path && (
                <Button
                  variant="outlined"
                  startIcon={<Download />}
                  onClick={() => downloadResult(currentJob.id)}
                  sx={{ mt: 1 }}
                >
                  Télécharger le résultat
                </Button>
              )}

              {currentJob.error && (
                <Alert severity="error" sx={{ mt: 1 }}>
                  {currentJob.error}
                </Alert>
              )}
            </Box>
          )}
        </CardContent>
      </Card>

      {/* Dialog historique */}
      <Dialog open={showHistory} onClose={() => setShowHistory(false)} maxWidth="md" fullWidth>
        <DialogTitle>Historique des exports</DialogTitle>
        <DialogContent>
          <List>
            {exportHistory.map((item, index) => (
              <ListItem key={index}>
                <ListItemText
                  primary={`Export ${item.id}`}
                  secondary={`${item.status} - ${new Date(item.start_time).toLocaleString()}`}
                />
                {item.download_available && (
                  <IconButton onClick={() => downloadResult(item.id)}>
                    <Download />
                  </IconButton>
                )}
              </ListItem>
            ))}
          </List>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowHistory(false)}>Fermer</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};