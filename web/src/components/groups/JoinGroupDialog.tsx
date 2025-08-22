import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Typography,
  IconButton,
  useTheme,
  Alert,
  CircularProgress,
  Tabs,
  Tab,
  Card,
  CardContent,
  Chip,
  InputAdornment,
  Divider,
} from '@mui/material';
import {
  Close as CloseIcon,
  Search as SearchIcon,
  Public as PublicIcon,
  Lock as LockIcon,
  GroupAdd as GroupAddIcon,
  People as PeopleIcon,
  VpnKey as VpnKeyIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import { groupApi, Group } from '../../services/groupApi';

interface JoinGroupDialogProps {
  open: boolean;
  onClose: () => void;
  onJoin: (inviteCode: string) => Promise<void>;
}

const JoinGroupDialog: React.FC<JoinGroupDialogProps> = ({ open, onClose, onJoin }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [activeTab, setActiveTab] = useState(0);
  const [inviteCode, setInviteCode] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [publicGroups, setPublicGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(false);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (open && activeTab === 1 && searchQuery.length >= 2) {
      searchPublicGroups();
    }
  }, [searchQuery, activeTab, open]);

  const searchPublicGroups = async () => {
    try {
      setSearching(true);
      setError(null);
      const groups = await groupApi.searchGroups(searchQuery);
      setPublicGroups(groups);
    } catch (err) {
      setError('Erreur lors de la recherche de groupes');
      console.error('Error searching groups:', err);
    } finally {
      setSearching(false);
    }
  };

  const handleJoinByCode = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!inviteCode.trim()) {
      setError('Veuillez entrer un code d\'invitation');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      await onJoin(inviteCode.trim().toUpperCase());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Code d\'invitation invalide');
    } finally {
      setLoading(false);
    }
  };

  const handleJoinPublicGroup = async (group: Group) => {
    if (group.invite_code) {
      try {
        setLoading(true);
        setError(null);
        await onJoin(group.invite_code);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Erreur lors de la tentative de rejoindre le groupe');
      } finally {
        setLoading(false);
      }
    }
  };

  const handleClose = () => {
    if (!loading) {
      setInviteCode('');
      setSearchQuery('');
      setPublicGroups([]);
      setError(null);
      setActiveTab(0);
      onClose();
    }
  };

  const getPrivacyIcon = (privacy: string) => {
    switch (privacy) {
      case 'public':
        return <PublicIcon sx={{ fontSize: 16, color: leagueColors.blue[500] }} />;
      case 'private':
        return <LockIcon sx={{ fontSize: 16, color: leagueColors.gold[500] }} />;
      default:
        return <GroupAddIcon sx={{ fontSize: 16, color: leagueColors.dark[400] }} />;
    }
  };

  const formatInviteCode = (value: string) => {
    // Format: XXXX-XXXX ou XXXXXXXX
    const cleaned = value.replace(/[^A-Za-z0-9]/g, '').toUpperCase();
    if (cleaned.length <= 8) {
      return cleaned.length > 4 ? `${cleaned.slice(0, 4)}-${cleaned.slice(4)}` : cleaned;
    }
    return cleaned.slice(0, 8);
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="md"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: 3,
          border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          background: isDarkMode
            ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
            : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
          minHeight: 500,
        },
      }}
    >
      <DialogTitle
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          pb: 0,
          background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
          color: '#fff',
          m: -3,
          mb: 0,
          px: 3,
          py: 2,
        }}
      >
        <Typography variant="h6" sx={{ fontWeight: 600 }}>
          Rejoindre un Groupe
        </Typography>
        <IconButton
          onClick={handleClose}
          disabled={loading}
          sx={{ color: '#fff' }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs
          value={activeTab}
          onChange={(e, newValue) => setActiveTab(newValue)}
          sx={{
            px: 3,
            '& .MuiTab-root': {
              fontWeight: 500,
            },
          }}
        >
          <Tab 
            icon={<VpnKeyIcon />} 
            label="Code d'invitation" 
            iconPosition="start"
          />
          <Tab 
            icon={<SearchIcon />} 
            label="Rechercher des groupes publics" 
            iconPosition="start"
          />
        </Tabs>
      </Box>

      <DialogContent sx={{ pt: 3 }}>
        {/* Error Alert */}
        {error && (
          <Alert 
            severity="error" 
            sx={{ mb: 3 }}
            onClose={() => setError(null)}
          >
            {error}
          </Alert>
        )}

        {/* Tab 0: Join by invite code */}
        {activeTab === 0 && (
          <Box>
            <Box sx={{ textAlign: 'center', mb: 4 }}>
              <VpnKeyIcon sx={{ fontSize: 48, color: leagueColors.gold[500], mb: 2 }} />
              <Typography variant="h6" gutterBottom>
                Entrez votre code d'invitation
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Votre ami doit vous avoir fourni un code d'invitation de 8 caractères
              </Typography>
            </Box>

            <form onSubmit={handleJoinByCode}>
              <TextField
                fullWidth
                label="Code d'invitation"
                value={inviteCode}
                onChange={(e) => setInviteCode(formatInviteCode(e.target.value))}
                placeholder="DEMO2025"
                disabled={loading}
                sx={{ mb: 3 }}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <VpnKeyIcon color="primary" />
                    </InputAdornment>
                  ),
                }}
                helperText="Les codes d'invitation sont fournis par les propriétaires de groupes"
              />

              <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                <Button
                  type="submit"
                  variant="contained"
                  disabled={loading || !inviteCode.trim()}
                  startIcon={loading ? <CircularProgress size={16} /> : <GroupAddIcon />}
                  sx={{
                    background: `linear-gradient(135deg, ${leagueColors.gold[500]} 0%, ${leagueColors.gold[600]} 100%)`,
                    color: '#000',
                    fontWeight: 600,
                    '&:hover': {
                      background: `linear-gradient(135deg, ${leagueColors.gold[600]} 0%, ${leagueColors.gold[700]} 100%)`,
                    },
                    '&:disabled': {
                      background: 'rgba(0,0,0,0.12)',
                      color: 'rgba(0,0,0,0.26)',
                    },
                  }}
                >
                  {loading ? 'Rejoindre...' : 'Rejoindre le Groupe'}
                </Button>
              </Box>
            </form>
          </Box>
        )}

        {/* Tab 1: Search public groups */}
        {activeTab === 1 && (
          <Box>
            <Box sx={{ mb: 3 }}>
              <TextField
                fullWidth
                label="Rechercher des groupes"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Nom du groupe, description..."
                disabled={loading}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon color="primary" />
                    </InputAdornment>
                  ),
                  endAdornment: searching && (
                    <InputAdornment position="end">
                      <CircularProgress size={20} />
                    </InputAdornment>
                  ),
                }}
                helperText="Tapez au moins 2 caractères pour commencer la recherche"
              />
            </Box>

            {searchQuery.length >= 2 && !searching && publicGroups.length === 0 && (
              <Box sx={{ textAlign: 'center', py: 4 }}>
                <SearchIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
                <Typography variant="h6" color="text.secondary">
                  Aucun groupe trouvé
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Essayez avec d'autres mots-clés
                </Typography>
              </Box>
            )}

            {searchQuery.length < 2 && (
              <Box sx={{ textAlign: 'center', py: 4 }}>
                <SearchIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
                <Typography variant="h6" color="text.secondary">
                  Recherchez des groupes publics
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Découvrez des groupes ouverts à tous les joueurs
                </Typography>
              </Box>
            )}

            {/* Public Groups Results */}
            {publicGroups.length > 0 && (
              <Box sx={{ maxHeight: 400, overflowY: 'auto' }}>
                {publicGroups.map((group) => (
                  <Card
                    key={group.id}
                    sx={{
                      mb: 2,
                      border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                      '&:hover': {
                        borderColor: leagueColors.blue[300],
                        boxShadow: `0 4px 12px ${isDarkMode ? 'rgba(0,0,0,0.2)' : 'rgba(25, 118, 210, 0.1)'}`,
                      },
                    }}
                  >
                    <CardContent>
                      <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', mb: 2 }}>
                        <Box sx={{ flexGrow: 1 }}>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                            <Typography variant="h6" sx={{ fontWeight: 600 }}>
                              {group.name}
                            </Typography>
                            {getPrivacyIcon(group.privacy)}
                          </Box>
                          
                          {group.description && (
                            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                              {group.description}
                            </Typography>
                          )}

                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
                            <Chip
                              icon={<PeopleIcon />}
                              label={`${group.member_count} membre${group.member_count > 1 ? 's' : ''}`}
                              size="small"
                              variant="outlined"
                            />
                            <Typography variant="caption" color="text.secondary">
                              par {group.owner.riot_id}#{group.owner.riot_tag}
                            </Typography>
                          </Box>
                        </Box>

                        <Button
                          variant="contained"
                          size="small"
                          onClick={() => handleJoinPublicGroup(group)}
                          disabled={loading}
                          startIcon={loading ? <CircularProgress size={16} /> : <GroupAddIcon />}
                          sx={{
                            background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
                            '&:hover': {
                              background: `linear-gradient(135deg, ${leagueColors.blue[600]} 0%, ${leagueColors.blue[700]} 100%)`,
                            },
                          }}
                        >
                          Rejoindre
                        </Button>
                      </Box>
                    </CardContent>
                  </Card>
                ))}
              </Box>
            )}
          </Box>
        )}
      </DialogContent>

      <Divider />
      
      <DialogActions sx={{ px: 3, py: 2 }}>
        <Button
          onClick={handleClose}
          disabled={loading}
          sx={{
            color: 'text.secondary',
            '&:hover': {
              backgroundColor: `${leagueColors.dark[400]}10`,
            },
          }}
        >
          Fermer
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default JoinGroupDialog;