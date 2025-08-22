import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  IconButton,
  Box,
  Typography,
  useTheme,
  Tabs,
  Tab,
  Button,
  Chip,
  Card,
  CardContent,
  Avatar,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  ListItemSecondaryAction,
  Divider,
  Alert,
  CircularProgress,
  Tooltip,
} from '@mui/material';
import {
  Close as CloseIcon,
  People as PeopleIcon,
  TrendingUp as TrendingUpIcon,
  CompareArrows as CompareArrowsIcon,
  Settings as SettingsIcon,
  ContentCopy as ContentCopyIcon,
  Share as ShareIcon,
  AdminPanelSettings as AdminIcon,
  Person as PersonIcon,
  Crown as CrownIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import { Group, GroupMember, groupApi } from '../../services/groupApi';
import ComparisonManager from './ComparisonManager';
import GroupStatsOverview from './GroupStatsOverview';
import GroupSettings from './GroupSettings';

interface GroupDetailsDialogProps {
  open: boolean;
  onClose: () => void;
  group: Group;
  onUpdate: () => void;
}

const GroupDetailsDialog: React.FC<GroupDetailsDialogProps> = ({ open, onClose, group, onUpdate }) => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  
  const [activeTab, setActiveTab] = useState(0);
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copySuccess, setCopySuccess] = useState(false);

  useEffect(() => {
    if (open) {
      loadGroupMembers();
    }
  }, [open, group.id]);

  const loadGroupMembers = async () => {
    try {
      setLoading(true);
      const groupMembers = await groupApi.getGroupMembers(group.id);
      setMembers(groupMembers);
      setError(null);
    } catch (err) {
      setError('Erreur lors du chargement des membres');
      console.error('Error loading members:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCopyInviteCode = async () => {
    if (group.invite_code) {
      try {
        await navigator.clipboard.writeText(group.invite_code);
        setCopySuccess(true);
        setTimeout(() => setCopySuccess(false), 2000);
      } catch (err) {
        console.error('Failed to copy invite code:', err);
      }
    }
  };

  const getRoleIcon = (role: string) => {
    switch (role) {
      case 'owner':
        return <CrownIcon sx={{ fontSize: 18, color: leagueColors.gold[500] }} />;
      case 'admin':
        return <AdminIcon sx={{ fontSize: 18, color: leagueColors.blue[500] }} />;
      default:
        return <PersonIcon sx={{ fontSize: 18, color: 'text.secondary' }} />;
    }
  };

  const getRoleLabel = (role: string) => {
    switch (role) {
      case 'owner':
        return 'Propriétaire';
      case 'admin':
        return 'Administrateur';
      default:
        return 'Membre';
    }
  };

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'owner':
        return leagueColors.gold[500];
      case 'admin':
        return leagueColors.blue[500];
      default:
        return 'text.secondary';
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="md"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: 3,
          border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          background: isDarkMode
            ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
            : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
          minHeight: 600,
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
        <Box>
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
            {group.name}
          </Typography>
          <Typography variant="body2" sx={{ opacity: 0.9 }}>
            {group.member_count} membre{group.member_count > 1 ? 's' : ''}
          </Typography>
        </Box>
        <IconButton
          onClick={onClose}
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
            icon={<PeopleIcon />} 
            label="Membres" 
            iconPosition="start"
          />
          <Tab 
            icon={<TrendingUpIcon />} 
            label="Statistiques" 
            iconPosition="start"
          />
          <Tab 
            icon={<CompareArrowsIcon />} 
            label="Comparaisons" 
            iconPosition="start"
          />
          <Tab 
            icon={<SettingsIcon />} 
            label="Paramètres" 
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

        {/* Tab 0: Members */}
        {activeTab === 0 && (
          <Box>
            {/* Group Info */}
            <Card sx={{ mb: 3, border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}` }}>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>
                    Informations du groupe
                  </Typography>
                  {group.invite_code && (
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      <Button
                        size="small"
                        startIcon={<ContentCopyIcon />}
                        onClick={handleCopyInviteCode}
                        sx={{
                          color: copySuccess ? leagueColors.win : leagueColors.blue[500],
                          borderColor: copySuccess ? leagueColors.win : leagueColors.blue[500],
                        }}
                        variant="outlined"
                      >
                        {copySuccess ? 'Copié!' : `Code: ${group.invite_code}`}
                      </Button>
                      <Tooltip title="Partager le groupe">
                        <IconButton size="small">
                          <ShareIcon />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  )}
                </Box>

                {group.description && (
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    {group.description}
                  </Typography>
                )}

                <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                  <Chip
                    label={`${group.privacy === 'public' ? 'Public' : group.privacy === 'private' ? 'Privé' : 'Sur invitation'}`}
                    size="small"
                    color="primary"
                  />
                  <Chip
                    label={`Créé le ${new Date(group.created_at).toLocaleDateString()}`}
                    size="small"
                    variant="outlined"
                  />
                </Box>
              </CardContent>
            </Card>

            {/* Members List */}
            <Card sx={{ border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}` }}>
              <CardContent>
                <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                  Membres ({members.length})
                </Typography>

                {loading ? (
                  <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
                    <CircularProgress />
                  </Box>
                ) : (
                  <List>
                    {members.map((member, index) => (
                      <React.Fragment key={member.id}>
                        <ListItem sx={{ px: 0 }}>
                          <ListItemAvatar>
                            <Avatar sx={{ 
                              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
                              color: '#fff',
                              fontWeight: 600,
                            }}>
                              {member.user.riot_id.charAt(0).toUpperCase()}
                            </Avatar>
                          </ListItemAvatar>
                          
                          <ListItemText
                            primary={
                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                  {member.user.riot_id}#{member.user.riot_tag}
                                </Typography>
                                {getRoleIcon(member.role)}
                              </Box>
                            }
                            secondary={
                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mt: 0.5 }}>
                                <Chip
                                  label={getRoleLabel(member.role)}
                                  size="small"
                                  sx={{
                                    height: 20,
                                    fontSize: '0.7rem',
                                    color: getRoleColor(member.role),
                                    borderColor: getRoleColor(member.role),
                                  }}
                                  variant="outlined"
                                />
                                <Typography variant="caption" color="text.secondary">
                                  {member.user.region?.toUpperCase()} • Membre depuis {new Date(member.joined_at).toLocaleDateString()}
                                </Typography>
                              </Box>
                            }
                          />

                          <ListItemSecondaryAction>
                            {member.role === 'owner' && (
                              <Chip
                                icon={<CrownIcon />}
                                label="Propriétaire"
                                size="small"
                                sx={{
                                  background: `linear-gradient(135deg, ${leagueColors.gold[500]} 0%, ${leagueColors.gold[600]} 100%)`,
                                  color: '#000',
                                  fontWeight: 600,
                                }}
                              />
                            )}
                          </ListItemSecondaryAction>
                        </ListItem>
                        
                        {index < members.length - 1 && <Divider variant="inset" component="li" />}
                      </React.Fragment>
                    ))}
                  </List>
                )}
              </CardContent>
            </Card>
          </Box>
        )}

        {/* Tab 1: Statistics */}
        {activeTab === 1 && (
          <GroupStatsOverview groupId={group.id} />
        )}

        {/* Tab 2: Comparisons */}
        {activeTab === 2 && (
          <ComparisonManager groupId={group.id} />
        )}

        {/* Tab 3: Settings */}
        {activeTab === 3 && (
          <GroupSettings group={group} onUpdate={onUpdate} />
        )}
      </DialogContent>
    </Dialog>
  );
};

export default GroupDetailsDialog;