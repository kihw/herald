import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  FormGroup,
  FormControlLabel,
  Switch,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Button,
  Divider,
  Alert,
  CircularProgress,
} from '@mui/material';
import {
  Save,
  RestoreFromTrash,
} from '@mui/icons-material';

interface UserSettings {
  include_timeline: boolean;
  include_all_data: boolean;
  light_mode: boolean;
  auto_sync_enabled: boolean;
  sync_frequency_hours: number;
}

interface SettingsPanelProps {
  settings: UserSettings;
  onSave: (settings: UserSettings) => void;
  loading?: boolean;
  saving?: boolean;
}

const SettingsPanel: React.FC<SettingsPanelProps> = ({
  settings,
  onSave,
  loading = false,
  saving = false,
}) => {
  const [localSettings, setLocalSettings] = useState<UserSettings>(settings);
  const [hasChanges, setHasChanges] = useState(false);

  const handleSettingChange = (key: keyof UserSettings, value: any) => {
    const newSettings = { ...localSettings, [key]: value };
    setLocalSettings(newSettings);
    setHasChanges(JSON.stringify(newSettings) !== JSON.stringify(settings));
  };

  const handleSave = () => {
    onSave(localSettings);
    setHasChanges(false);
  };

  const handleReset = () => {
    setLocalSettings(settings);
    setHasChanges(false);
  };

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h6" component="div">
            User Settings
          </Typography>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              variant="outlined"
              startIcon={<RestoreFromTrash />}
              onClick={handleReset}
              disabled={!hasChanges || saving}
              size="small"
            >
              Reset
            </Button>
            <Button
              variant="contained"
              startIcon={saving ? <CircularProgress size={16} /> : <Save />}
              onClick={handleSave}
              disabled={!hasChanges || saving}
              size="small"
            >
              {saving ? 'Saving...' : 'Save Changes'}
            </Button>
          </Box>
        </Box>

        {hasChanges && (
          <Alert severity="warning" sx={{ mb: 3 }}>
            You have unsaved changes. Don't forget to save!
          </Alert>
        )}

        {/* Data Collection Settings */}
        <Typography variant="h6" sx={{ mb: 2 }}>
          Data Collection
        </Typography>
        <FormGroup sx={{ mb: 3 }}>
          <FormControlLabel
            control={
              <Switch
                checked={localSettings.include_timeline}
                onChange={(e) => handleSettingChange('include_timeline', e.target.checked)}
              />
            }
            label="Include match timeline data"
          />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 4, mb: 2 }}>
            Collect detailed timeline events for each match (recommended for advanced analytics)
          </Typography>

          <FormControlLabel
            control={
              <Switch
                checked={localSettings.include_all_data}
                onChange={(e) => handleSettingChange('include_all_data', e.target.checked)}
              />
            }
            label="Include all match data"
          />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 4 }}>
            Collect complete match data including all players (larger file sizes)
          </Typography>
        </FormGroup>

        <Divider sx={{ my: 3 }} />

        {/* Sync Settings */}
        <Typography variant="h6" sx={{ mb: 2 }}>
          Synchronization
        </Typography>
        <FormGroup sx={{ mb: 3 }}>
          <FormControlLabel
            control={
              <Switch
                checked={localSettings.auto_sync_enabled}
                onChange={(e) => handleSettingChange('auto_sync_enabled', e.target.checked)}
              />
            }
            label="Enable automatic sync"
          />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 4, mb: 2 }}>
            Automatically synchronize new matches at regular intervals
          </Typography>

          <Box sx={{ ml: 4, maxWidth: 300 }}>
            <FormControl fullWidth disabled={!localSettings.auto_sync_enabled}>
              <InputLabel>Sync Frequency</InputLabel>
              <Select
                value={localSettings.sync_frequency_hours}
                label="Sync Frequency"
                onChange={(e) => handleSettingChange('sync_frequency_hours', e.target.value)}
              >
                <MenuItem value={1}>Every hour</MenuItem>
                <MenuItem value={6}>Every 6 hours</MenuItem>
                <MenuItem value={12}>Every 12 hours</MenuItem>
                <MenuItem value={24}>Daily</MenuItem>
                <MenuItem value={72}>Every 3 days</MenuItem>
                <MenuItem value={168}>Weekly</MenuItem>
              </Select>
            </FormControl>
          </Box>
        </FormGroup>

        <Divider sx={{ my: 3 }} />

        {/* Appearance Settings */}
        <Typography variant="h6" sx={{ mb: 2 }}>
          Appearance
        </Typography>
        <FormGroup>
          <FormControlLabel
            control={
              <Switch
                checked={localSettings.light_mode}
                onChange={(e) => handleSettingChange('light_mode', e.target.checked)}
              />
            }
            label="Light mode"
          />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 4 }}>
            Use light theme instead of dark theme
          </Typography>
        </FormGroup>

        <Divider sx={{ my: 3 }} />

        {/* Additional Settings */}
        <Typography variant="h6" sx={{ mb: 2 }}>
          Advanced
        </Typography>
        <Alert severity="info">
          Additional settings will be available in future updates. This includes
          export formats, data retention policies, and integration options.
        </Alert>
      </CardContent>
    </Card>
  );
};

export default SettingsPanel;
