import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  Switch,
  Grid,
  Alert,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  useTheme,
} from '@mui/material';
import {
  Settings as SettingsIcon,
  Save as SaveIcon,
  Delete as DeleteIcon,
  Warning as WarningIcon,
  Public as PublicIcon,
  Lock as LockIcon,
  GroupAdd as GroupAddIcon,
  ContentCopy as ContentCopyIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import { Group } from '../../services/groupApi';

interface GroupSettingsProps {
  group: Group;
  onUpdate: () => void;
}

const GroupSettings: React.FC<GroupSettingsProps> = ({ group, onUpdate }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [formData, setFormData] = useState({
    name: group.name,
    description: group.description,
    privacy: group.privacy,
  });
  
  const [settings, setSettings] = useState({
    allowMemberInvites: true,
    autoAcceptJoinRequests: false,
    showMemberStats: true,
    allowComparisons: true,
    publicProfile: group.privacy === 'public',
    enableNotifications: true,
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [inviteCodeCopied, setInviteCodeCopied] = useState(false);

  const handleSaveBasicInfo = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSuccess('Informations mises à jour avec succès');
      onUpdate();
    } catch (err) {
      setError('Erreur lors de la mise à jour');
    } finally {
      setLoading(false);
    }
  };

  const handleSaveSettings = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSuccess('Paramètres sauvegardés avec succès');
    } catch (err) {
      setError('Erreur lors de la sauvegarde');
    } finally {
      setLoading(false);
    }
  };

  const handleRegenerateInviteCode = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSuccess('Nouveau code d\'invitation généré');
      onUpdate();
    } catch (err) {
      setError('Erreur lors de la génération du code');
    } finally {
      setLoading(false);
    }
  };

  const handleCopyInviteCode = async () => {
    if (group.invite_code) {
      try {
        await navigator.clipboard.writeText(group.invite_code);
        setInviteCodeCopied(true);
        setTimeout(() => setInviteCodeCopied(false), 2000);
      } catch (err) {
        console.error('Failed to copy invite code:', err);
      }
    }
  };

  const handleDeleteGroup = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setDeleteDialogOpen(false);
      // Group deleted, close dialog or redirect
    } catch (err) {
      setError('Erreur lors de la suppression du groupe');
    } finally {
      setLoading(false);
    }
  };

  const privacyOptions = [
    {
      value: 'private',
      label: 'Privé',
      description: 'Seuls les membres invités peuvent rejoindre',
      icon: <LockIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.gold[500],
    },
    {
      value: 'invite_only',
      label: 'Sur invitation',
      description: 'Visible mais nécessite une invitation pour rejoindre',
      icon: <GroupAddIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.dark[400],
    },
    {
      value: 'public',
      label: 'Public',
      description: 'Tout le monde peut rechercher et rejoindre ce groupe',
      icon: <PublicIcon sx={{ fontSize: 20 }} />,
      color: leagueColors.blue[500],
    },
  ];

  return (
    <Box sx={{ p: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography 
          variant="h5" 
          sx={{ 
            fontWeight: 700,
            background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
            backgroundClip: 'text',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            mb: 1,
            display: 'flex',
            alignItems: 'center',
            gap: 1,
          }}
        >
          <SettingsIcon sx={{ color: leagueColors.blue[500] }} />
          Paramètres du Groupe
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Gérez les informations et les paramètres de votre groupe
        </Typography>
      </Box>

      {/* Alerts */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}
      {success && (
        <Alert severity="success" sx={{ mb: 3 }} onClose={() => setSuccess(null)}>
          {success}
        </Alert>
      )}

      <Grid container spacing={3}>
        {/* Basic Information */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
                Informations de base
              </Typography>

              <Grid container spacing={3}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="Nom du groupe"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    disabled={loading}
                    inputProps={{ maxLength: 50 }}
                    helperText={`${formData.name.length}/50 caractères`}
                  />
                </Grid>

                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    label="Description"
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    multiline
                    rows={3}
                    disabled={loading}
                    inputProps={{ maxLength: 500 }}
                    helperText={`${formData.description.length}/500 caractères`}
                  />
                </Grid>

                <Grid item xs={12}>
                  <Button
                    variant="contained"
                    startIcon={<SaveIcon />}
                    onClick={handleSaveBasicInfo}
                    disabled={loading || !formData.name.trim()}
                    sx={{
                      background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
                    }}
                  >
                    Sauvegarder les informations
                  </Button>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Privacy Settings */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
                Confidentialité
              </Typography>

              <FormControl component="fieldset" sx={{ width: '100%' }}>
                <RadioGroup
                  value={formData.privacy}
                  onChange={(e) => setFormData({ ...formData, privacy: e.target.value as any })}
                  disabled={loading}
                >
                  {privacyOptions.map((option) => (
                    <Box key={option.value} sx={{ mb: 2 }}>
                      <FormControlLabel
                        value={option.value}
                        control={<Radio sx={{ color: option.color }} />}
                        label={
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <Box sx={{ color: option.color }}>
                              {option.icon}
                            </Box>
                            <Box>
                              <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                {option.label}
                              </Typography>
                              <Typography variant="body2" color="text.secondary">
                                {option.description}
                              </Typography>
                            </Box>
                          </Box>
                        }
                        sx={{
                          alignItems: 'flex-start',
                          border: `1px solid ${formData.privacy === option.value ? option.color : 'transparent'}`,
                          borderRadius: 2,
                          p: 1,
                          m: 0,
                          width: '100%',
                          '& .MuiFormControlLabel-label': {
                            mt: 0.5,
                            width: '100%',
                          },
                        }}
                      />
                    </Box>
                  ))}
                </RadioGroup>
              </FormControl>
            </CardContent>
          </Card>
        </Grid>

        {/* Invite Code Management */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
                Code d'invitation
              </Typography>

              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                <TextField
                  label="Code actuel"
                  value={group.invite_code || 'Aucun code généré'}
                  disabled
                  sx={{ flexGrow: 1 }}
                />
                <Button
                  variant="outlined"
                  startIcon={inviteCodeCopied ? <SaveIcon /> : <ContentCopyIcon />}
                  onClick={handleCopyInviteCode}
                  disabled={!group.invite_code}
                  sx={{
                    color: inviteCodeCopied ? leagueColors.win : leagueColors.blue[500],
                    borderColor: inviteCodeCopied ? leagueColors.win : leagueColors.blue[500],
                  }}
                >
                  {inviteCodeCopied ? 'Copié!' : 'Copier'}
                </Button>
                <Button
                  variant="outlined"
                  startIcon={<RefreshIcon />}
                  onClick={handleRegenerateInviteCode}
                  disabled={loading}
                  sx={{
                    color: leagueColors.gold[500],
                    borderColor: leagueColors.gold[500],
                  }}
                >
                  Régénérer
                </Button>
              </Box>

              <Alert severity="info">
                Le code d'invitation permet aux autres joueurs de rejoindre directement votre groupe.
              </Alert>
            </CardContent>
          </Card>
        </Grid>

        {/* Group Settings */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
                Paramètres avancés
              </Typography>

              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body1" sx={{ fontWeight: 500 }}>
                        Invitations par les membres
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Les membres peuvent inviter d'autres joueurs
                      </Typography>
                    </Box>
                    <Switch
                      checked={settings.allowMemberInvites}
                      onChange={(e) => setSettings({ ...settings, allowMemberInvites: e.target.checked })}
                      disabled={loading}
                    />
                  </Box>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body1" sx={{ fontWeight: 500 }}>
                        Acceptation automatique
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Accepter automatiquement les demandes de rejoindre
                      </Typography>
                    </Box>
                    <Switch
                      checked={settings.autoAcceptJoinRequests}
                      onChange={(e) => setSettings({ ...settings, autoAcceptJoinRequests: e.target.checked })}
                      disabled={loading}
                    />
                  </Box>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body1" sx={{ fontWeight: 500 }}>
                        Statistiques des membres
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Afficher les stats individuelles aux membres
                      </Typography>
                    </Box>
                    <Switch
                      checked={settings.showMemberStats}
                      onChange={(e) => setSettings({ ...settings, showMemberStats: e.target.checked })}
                      disabled={loading}
                    />
                  </Box>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body1" sx={{ fontWeight: 500 }}>
                        Comparaisons
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Permettre la création de comparaisons
                      </Typography>
                    </Box>
                    <Switch
                      checked={settings.allowComparisons}
                      onChange={(e) => setSettings({ ...settings, allowComparisons: e.target.checked })}
                      disabled={loading}
                    />
                  </Box>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body1" sx={{ fontWeight: 500 }}>
                        Profil public
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Rendre le groupe visible dans les recherches
                      </Typography>
                    </Box>
                    <Switch
                      checked={settings.publicProfile}
                      onChange={(e) => setSettings({ ...settings, publicProfile: e.target.checked })}
                      disabled={loading}
                    />
                  </Box>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="body1" sx={{ fontWeight: 500 }}>
                        Notifications
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Recevoir des notifications d'activité du groupe
                      </Typography>
                    </Box>
                    <Switch
                      checked={settings.enableNotifications}
                      onChange={(e) => setSettings({ ...settings, enableNotifications: e.target.checked })}
                      disabled={loading}
                    />
                  </Box>
                </Grid>

                <Grid item xs={12}>
                  <Divider sx={{ my: 2 }} />
                  <Button
                    variant="contained"
                    startIcon={<SaveIcon />}
                    onClick={handleSaveSettings}
                    disabled={loading}
                    sx={{
                      background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
                    }}
                  >
                    Sauvegarder les paramètres
                  </Button>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Danger Zone */}
        <Grid item xs={12}>
          <Card sx={{ border: `1px solid ${leagueColors.loss}` }}>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3, color: leagueColors.loss }}>
                Zone dangereuse
              </Typography>

              <Alert severity="warning" sx={{ mb: 3 }}>
                <Typography variant="body2">
                  Ces actions sont irréversibles. Assurez-vous de bien comprendre les conséquences.
                </Typography>
              </Alert>

              <Button
                variant="outlined"
                startIcon={<DeleteIcon />}
                onClick={() => setDeleteDialogOpen(true)}
                disabled={loading}
                sx={{
                  color: leagueColors.loss,
                  borderColor: leagueColors.loss,
                  '&:hover': {
                    borderColor: leagueColors.loss,
                    backgroundColor: `${leagueColors.loss}10`,
                  },
                }}
              >
                Supprimer le groupe
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <WarningIcon sx={{ color: leagueColors.loss }} />
          Confirmer la suppression
        </DialogTitle>
        <DialogContent>
          <Typography variant="body1" sx={{ mb: 2 }}>
            Êtes-vous sûr de vouloir supprimer le groupe <strong>{group.name}</strong> ?
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Cette action est irréversible. Toutes les données, statistiques, et comparaisons 
            associées au groupe seront définitivement supprimées.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => setDeleteDialogOpen(false)}
            disabled={loading}
          >
            Annuler
          </Button>
          <Button
            onClick={handleDeleteGroup}
            disabled={loading}
            variant="contained"
            sx={{
              backgroundColor: leagueColors.loss,
              '&:hover': {
                backgroundColor: '#d32f2f',
              },
            }}
          >
            Supprimer définitivement
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default GroupSettings;