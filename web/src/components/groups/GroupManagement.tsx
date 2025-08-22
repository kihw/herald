import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Grid,
  Chip,
  Avatar,
  AvatarGroup,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  useTheme,
  Fab,
  Tooltip,
  Alert,
  CircularProgress,
} from '@mui/material';
import useResponsive from '../../hooks/useResponsive';
import { usePerformance } from '../../hooks/usePerformance';
import { ResponsiveContainer, ResponsiveCardContainer } from '../common/ResponsiveContainer';
import { CardSkeleton } from '../common/LazyComponent';
import {
  Add as AddIcon,
  People as PeopleIcon,
  Settings as SettingsIcon,
  MoreVert as MoreVertIcon,
  Public as PublicIcon,
  Lock as LockIcon,
  Search as SearchIcon,
  GroupAdd as GroupAddIcon,
  TrendingUp as TrendingUpIcon,
} from '@mui/icons-material';
import { leagueColors } from '../../theme/leagueTheme';
import CreateGroupDialog from './CreateGroupDialog';
import JoinGroupDialog from './JoinGroupDialog';
import GroupDetailsDialog from './GroupDetailsDialog';
import { groupApi } from '../../services/groupApi';

interface Group {
  id: number;
  name: string;
  description: string;
  privacy: 'public' | 'private' | 'invite_only';
  member_count: number;
  owner: {
    riot_id: string;
    riot_tag: string;
  };
  members?: GroupMember[];
  invite_code?: string;
  created_at: string;
}

interface GroupMember {
  id: number;
  user: {
    riot_id: string;
    riot_tag: string;
    region: string;
  };
  role: 'owner' | 'admin' | 'member';
  joined_at: string;
}

const GroupManagement: React.FC = () => {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  const { isMobile, getGridSpacing, getCardPadding } = useResponsive();
  const { shouldReduceAnimations, measureRenderTime } = usePerformance();
  
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Dialog states
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [joinDialogOpen, setJoinDialogOpen] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  
  // Menu state
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [menuGroupId, setMenuGroupId] = useState<number | null>(null);

  useEffect(() => {
    loadUserGroups();
  }, []);

  const loadUserGroups = async () => {
    try {
      setLoading(true);
      const userGroups = await groupApi.getUserGroups();
      setGroups(userGroups);
      setError(null);
    } catch (err) {
      setError('Erreur lors du chargement des groupes');
      console.error('Error loading groups:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, groupId: number) => {
    setAnchorEl(event.currentTarget);
    setMenuGroupId(groupId);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setMenuGroupId(null);
  };

  const handleGroupDetails = (group: Group) => {
    setSelectedGroup(group);
    setDetailsDialogOpen(true);
    handleMenuClose();
  };

  const handleCreateGroup = async (groupData: any) => {
    try {
      await groupApi.createGroup(groupData);
      await loadUserGroups();
      setCreateDialogOpen(false);
    } catch (err) {
      console.error('Error creating group:', err);
    }
  };

  const handleJoinGroup = async (inviteCode: string) => {
    try {
      await groupApi.joinGroup(inviteCode);
      await loadUserGroups();
      setJoinDialogOpen(false);
    } catch (err) {
      console.error('Error joining group:', err);
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

  const getPrivacyColor = (privacy: string) => {
    switch (privacy) {
      case 'public':
        return leagueColors.blue[500];
      case 'private':
        return leagueColors.gold[500];
      default:
        return leagueColors.dark[400];
    }
  };

  if (loading) {
    return (
      <ResponsiveContainer>
        <CardSkeleton count={6} height={200} />
      </ResponsiveContainer>
    );
  }

  return (
    <ResponsiveContainer>
      {/* Header */}
      <Box sx={{ 
        mb: getGridSpacing(), 
        display: 'flex', 
        alignItems: isMobile ? 'flex-start' : 'center', 
        justifyContent: 'space-between',
        flexDirection: isMobile ? 'column' : 'row',
        gap: isMobile ? 2 : 0,
      }}>
        <Box sx={{ flexGrow: isMobile ? 1 : 0 }}>
          <Typography 
            variant={isMobile ? "h5" : "h4"}
            sx={{ 
              fontWeight: 700,
              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.gold[500]} 100%)`,
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              mb: 1,
            }}
          >
            Mes Groupes d'Amis
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Créez des groupes pour comparer vos performances avec vos amis League of Legends
          </Typography>
        </Box>
        
        <Box sx={{ 
          display: 'flex', 
          gap: 1,
          flexDirection: isMobile ? 'row' : 'row',
          width: isMobile ? '100%' : 'auto',
        }}>
          <Button
            variant="outlined"
            startIcon={!isMobile ? <SearchIcon /> : undefined}
            onClick={() => setJoinDialogOpen(true)}
            fullWidth={isMobile}
            size={isMobile ? 'large' : 'medium'}
            sx={{
              borderColor: leagueColors.blue[500],
              color: leagueColors.blue[500],
              '&:hover': {
                borderColor: leagueColors.blue[600],
                backgroundColor: `${leagueColors.blue[500]}10`,
              },
            }}
          >
            {isMobile ? <SearchIcon sx={{ mr: 1 }} /> : null}
            Rejoindre
          </Button>
          <Button
            variant="contained"
            startIcon={!isMobile ? <AddIcon /> : undefined}
            onClick={() => setCreateDialogOpen(true)}
            fullWidth={isMobile}
            size={isMobile ? 'large' : 'medium'}
            sx={{
              background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
              '&:hover': {
                background: `linear-gradient(135deg, ${leagueColors.blue[600]} 0%, ${leagueColors.blue[700]} 100%)`,
              },
            }}
          >
            {isMobile ? <AddIcon sx={{ mr: 1 }} /> : null}
            Créer un Groupe
          </Button>
        </Box>
      </Box>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Groups Grid */}
      {groups.length === 0 ? (
        <Card 
          sx={{ 
            textAlign: 'center', 
            py: isMobile ? 4 : 6,
            background: isDarkMode
              ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
              : `linear-gradient(135deg, ${leagueColors.blue[25]} 0%, #ffffff 100%)`,
            border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          }}
        >
          <CardContent sx={{ p: getCardPadding() }}>
            <PeopleIcon sx={{ fontSize: isMobile ? 48 : 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant={isMobile ? "subtitle1" : "h6"} gutterBottom>
              Aucun groupe pour l'instant
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              Créez votre premier groupe ou rejoignez un groupe existant pour commencer à comparer vos performances
            </Typography>
            <Box sx={{ 
              display: 'flex', 
              gap: 2, 
              justifyContent: 'center',
              flexDirection: isMobile ? 'column' : 'row',
            }}>
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={() => setCreateDialogOpen(true)}
                fullWidth={isMobile}
                size={isMobile ? 'large' : 'medium'}
                sx={{
                  background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
                }}
              >
                Créer un Groupe
              </Button>
              <Button
                variant="outlined"
                startIcon={<SearchIcon />}
                onClick={() => setJoinDialogOpen(true)}
                fullWidth={isMobile}
                size={isMobile ? 'large' : 'medium'}
                sx={{
                  borderColor: leagueColors.blue[500],
                  color: leagueColors.blue[500],
                }}
              >
                Rejoindre un Groupe
              </Button>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <ResponsiveCardContainer
          columns={{ xs: 1, sm: 2, md: 3, lg: 4 }}
          spacing={getGridSpacing()}
        >
          {groups.map((group) => (
            <Card
              key={group.id}
              sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                background: isDarkMode
                  ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
                  : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
                border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                transition: shouldReduceAnimations ? 'none' : 'all 0.3s ease',
                cursor: 'pointer',
                '&:hover': shouldReduceAnimations ? {} : {
                  transform: 'translateY(-4px)',
                  boxShadow: `0 8px 25px ${isDarkMode ? 'rgba(0,0,0,0.3)' : 'rgba(25, 118, 210, 0.15)'}`,
                  borderColor: leagueColors.blue[300],
                },
              }}
              onClick={() => handleGroupDetails(group)}
            >
                <CardContent sx={{ flexGrow: 1 }}>
                  {/* Header */}
                  <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', mb: 2 }}>
                    <Box sx={{ flexGrow: 1 }}>
                      <Typography variant="h6" sx={{ fontWeight: 600, mb: 0.5 }}>
                        {group.name}
                      </Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                        {getPrivacyIcon(group.privacy)}
                        <Typography variant="caption" sx={{ textTransform: 'capitalize' }}>
                          {group.privacy === 'public' ? 'Public' : 
                           group.privacy === 'private' ? 'Privé' : 'Sur invitation'}
                        </Typography>
                      </Box>
                    </Box>
                    <IconButton
                      size="small"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleMenuOpen(e, group.id);
                      }}
                    >
                      <MoreVertIcon />
                    </IconButton>
                  </Box>

                  {/* Description */}
                  {group.description && (
                    <Typography 
                      variant="body2" 
                      color="text.secondary" 
                      sx={{ mb: 2, lineHeight: 1.4 }}
                    >
                      {group.description}
                    </Typography>
                  )}

                  {/* Owner */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                    <Avatar sx={{ width: 24, height: 24, fontSize: '0.8rem' }}>
                      {group.owner.riot_id.charAt(0).toUpperCase()}
                    </Avatar>
                    <Typography variant="body2" color="text.secondary">
                      {group.owner.riot_id}#{group.owner.riot_tag}
                    </Typography>
                  </Box>

                  {/* Members */}
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <PeopleIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
                      <Typography variant="body2" color="text.secondary">
                        {group.member_count} membre{group.member_count > 1 ? 's' : ''}
                      </Typography>
                    </Box>
                    
                    {group.members && group.members.length > 0 && (
                      <AvatarGroup max={3} sx={{ '& .MuiAvatar-root': { width: 28, height: 28, fontSize: '0.7rem' } }}>
                        {group.members.slice(0, 3).map((member) => (
                          <Avatar key={member.id}>
                            {member.user.riot_id.charAt(0).toUpperCase()}
                          </Avatar>
                        ))}
                      </AvatarGroup>
                    )}
                  </Box>
                </CardContent>

                {/* Footer with stats */}
                <Box 
                  sx={{ 
                    px: 2, 
                    py: 1.5, 
                    borderTop: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                    background: isDarkMode 
                      ? `${leagueColors.dark[200]}20` 
                      : `${leagueColors.blue[50]}50`,
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Typography variant="caption" color="text.secondary">
                      Créé le {new Date(group.created_at).toLocaleDateString()}
                    </Typography>
                    <Chip
                      icon={<TrendingUpIcon />}
                      label="Voir Stats"
                      size="small"
                      sx={{
                        height: 24,
                        fontSize: '0.7rem',
                        background: `linear-gradient(135deg, ${getPrivacyColor(group.privacy)}20 0%, ${getPrivacyColor(group.privacy)}10 100%)`,
                        color: getPrivacyColor(group.privacy),
                        border: `1px solid ${getPrivacyColor(group.privacy)}40`,
                      }}
                    />
                  </Box>
                </Box>
              </Card>
            ))}
        </ResponsiveCardContainer>
      )}

      {/* Group Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        PaperProps={{
          sx: {
            borderRadius: 2,
            border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
          },
        }}
      >
        <MenuItem onClick={() => {
          const group = groups.find(g => g.id === menuGroupId);
          if (group) handleGroupDetails(group);
        }}>
          <PeopleIcon sx={{ mr: 1, fontSize: 18 }} />
          Voir le groupe
        </MenuItem>
        <MenuItem onClick={handleMenuClose}>
          <SettingsIcon sx={{ mr: 1, fontSize: 18 }} />
          Paramètres
        </MenuItem>
        <MenuItem onClick={handleMenuClose}>
          <TrendingUpIcon sx={{ mr: 1, fontSize: 18 }} />
          Statistiques
        </MenuItem>
      </Menu>

      {/* Dialogs */}
      <CreateGroupDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onCreate={handleCreateGroup}
      />

      <JoinGroupDialog
        open={joinDialogOpen}
        onClose={() => setJoinDialogOpen(false)}
        onJoin={handleJoinGroup}
      />

      {selectedGroup && (
        <GroupDetailsDialog
          open={detailsDialogOpen}
          onClose={() => {
            setDetailsDialogOpen(false);
            setSelectedGroup(null);
          }}
          group={selectedGroup}
          onUpdate={loadUserGroups}
        />
      )}
    </ResponsiveContainer>
  );
};

export default GroupManagement;