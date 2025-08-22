import React, { useState } from 'react';
import { useExport } from '../../hooks/useExport';
import { useExportFeedback } from '../ExportFeedback';
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Box,
  Container,
  useTheme,
  useMediaQuery,
  Divider,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Tooltip,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Brightness4,
  Brightness7,
  SportsEsports,
  EmojiEvents,
  Analytics,
  Settings,
  GetApp,
  Dashboard,
  Groups,
  TrendingUp,
} from '@mui/icons-material';
import { useTheme as useAppTheme } from '../../App';

const DRAWER_WIDTH = 280;

// Types pour la navigation
export type ViewType = 'overview' | 'roles' | 'champions' | 'details' | 'analytics-dashboard' | 'analytics-champions' | 'analytics-mmr' | 'comparator' | 'heatmap' | 'groups';

export interface AppLayoutProps {
  children: React.ReactNode;
  currentView: ViewType;
  onViewChange: (view: ViewType) => void;
  selectedRole?: string;
  selectedChampion?: string;
  onRoleSelect?: (role: string | undefined) => void;
  onChampionSelect?: (champion: string | undefined) => void;
  data?: any[]; // Données pour l'export global
}

// Configuration des serveurs/régions
const REGIONS = [
  { value: 'euw1', label: 'Europe West (EUW)' },
  { value: 'eun1', label: 'Europe Nordic & East (EUNE)' },
  { value: 'na1', label: 'North America (NA)' },
  { value: 'kr', label: 'Korea (KR)' },
  { value: 'br1', label: 'Brazil (BR)' },
  { value: 'jp1', label: 'Japan (JP)' },
  { value: 'oc1', label: 'Oceania (OCE)' },
  { value: 'tr1', label: 'Turkey (TR)' },
  { value: 'ru', label: 'Russia (RU)' },
];

// Saisons disponibles
const SEASONS = [
  { value: '2024', label: 'Saison 2024' },
  { value: '2023', label: 'Saison 2023' },
  { value: '2022', label: 'Saison 2022' },
  { value: '2021', label: 'Saison 2021' },
];

export const AppLayout: React.FC<AppLayoutProps> = ({
  children,
  currentView,
  onViewChange,
  selectedRole,
  selectedChampion,
  onRoleSelect,
  onChampionSelect,
  data = [],
}) => {
  const theme = useTheme();
  const { isDarkMode, toggleTheme } = useAppTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('lg'));
  
  const [mobileOpen, setMobileOpen] = useState(false);
  const [selectedRegion, setSelectedRegion] = useState('euw1');
  const [selectedSeason, setSelectedSeason] = useState('2024');
  
  const { showSuccess, showError } = useExportFeedback();
  const { exportCombined } = useExport({
    onSuccess: showSuccess,
    onError: showError,
  });

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  // Configuration des éléments de navigation
  const navigationItems = [
    {
      key: 'overview',
      label: 'Vue d\'ensemble',
      icon: <Analytics />,
      view: 'overview' as ViewType,
    },
    {
      key: 'roles',
      label: 'Rôles & Lanes',
      icon: <SportsEsports />,
      view: 'roles' as ViewType,
    },
    {
      key: 'groups',
      label: 'Groupes d\'Amis',
      icon: <Groups />,
      view: 'groups' as ViewType,
    },
  ];

  // Navigation pour les analytics avancées
  const analyticsItems = [
    {
      key: 'analytics-dashboard',
      label: 'Dashboard Analytics',
      icon: <Dashboard />,
      view: 'analytics-dashboard' as ViewType,
    },
    {
      key: 'analytics-champions',
      label: 'Analyse Champions',
      icon: <EmojiEvents />,
      view: 'analytics-champions' as ViewType,
    },
    {
      key: 'analytics-mmr',
      label: 'Analyse MMR',
      icon: <TrendingUp />,
      view: 'analytics-mmr' as ViewType,
    },
    {
      key: 'comparator',
      label: 'Comparateur',
      icon: <Settings />,
      view: 'comparator' as ViewType,
    },
    {
      key: 'heatmap',
      label: 'Heatmap',
      icon: <Analytics />,
      view: 'heatmap' as ViewType,
    },
  ];

  // Breadcrumb dynamique
  const getBreadcrumb = () => {
    const parts = ['Dashboard'];
    
    // Analytics views
    if (currentView === 'analytics-dashboard') {
      parts.push('Analytics', 'Dashboard');
    } else if (currentView === 'analytics-champions') {
      parts.push('Analytics', 'Champions');
    } else if (currentView === 'analytics-mmr') {
      parts.push('Analytics', 'MMR');
    }
    // Legacy views
    else if (currentView === 'roles' || selectedRole) {
      parts.push('Rôles');
    }
    
    if (selectedRole) {
      parts.push(selectedRole);
    }
    
    if (currentView === 'champions' || selectedChampion) {
      parts.push('Champions');
    }
    
    if (selectedChampion) {
      parts.push(selectedChampion);
    }
    
    return parts.join(' › ');
  };

  // Export global de la vue actuelle
  const handleGlobalExport = () => {
    if (!data.length) {
      showError('Aucune donnée disponible pour l\'export');
      return;
    }

    let elementId = '';
    let exportType: 'roles' | 'champions' = 'roles';

    switch (currentView) {
      case 'roles':
        elementId = 'roles-view';
        exportType = 'roles';
        break;
      case 'champions':
        elementId = 'champions-view';
        exportType = 'champions';
        break;
      case 'details':
        elementId = 'champion-details-view';
        exportType = 'champions';
        break;
      default:
        elementId = 'overview-view';
        exportType = 'roles';
    }

    exportCombined(elementId, data, exportType, { selectedRole });
  };

  const drawer = (
    <Box sx={{ mt: 1 }}>
      {/* En-tête avec logo */}
      <Box sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
        <EmojiEvents sx={{ color: 'primary.main', fontSize: 28 }} />
        <Typography variant="h6" sx={{ fontWeight: 700, color: 'primary.main' }}>
          LoL Analytics
        </Typography>
      </Box>
      
      <Divider />
      
      {/* Sélecteurs globaux */}
      <Box sx={{ p: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
        <FormControl size="small" fullWidth>
          <InputLabel>Région</InputLabel>
          <Select
            value={selectedRegion}
            label="Région"
            onChange={(e) => setSelectedRegion(e.target.value)}
          >
            {REGIONS.map((region) => (
              <MenuItem key={region.value} value={region.value}>
                {region.label}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        
        <FormControl size="small" fullWidth>
          <InputLabel>Saison</InputLabel>
          <Select
            value={selectedSeason}
            label="Saison"
            onChange={(e) => setSelectedSeason(e.target.value)}
          >
            {SEASONS.map((season) => (
              <MenuItem key={season.value} value={season.value}>
                {season.label}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>
      
      <Divider />
      
      {/* Navigation principale */}
      <List>
        {navigationItems.map((item) => (
          <ListItem key={item.key} disablePadding>
            <ListItemButton
              selected={currentView === item.view}
              onClick={() => {
                onViewChange(item.view);
                if (isMobile) setMobileOpen(false);
              }}
              sx={{
                '&.Mui-selected': {
                  backgroundColor: 'primary.main',
                  color: 'primary.contrastText',
                  '&:hover': {
                    backgroundColor: 'primary.dark',
                  },
                },
              }}
            >
              <ListItemIcon
                sx={{
                  color: currentView === item.view ? 'inherit' : 'text.secondary',
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.label} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      
      <Divider sx={{ my: 1 }} />
      
      {/* Analytics avancées */}
      <Box sx={{ px: 2, py: 1 }}>
        <Typography variant="overline" color="text.secondary" sx={{ fontWeight: 600 }}>
          Analytics Avancées
        </Typography>
      </Box>
      <List>
        {analyticsItems.map((item) => (
          <ListItem key={item.key} disablePadding>
            <ListItemButton
              selected={currentView === item.view}
              onClick={() => {
                onViewChange(item.view);
                if (isMobile) setMobileOpen(false);
              }}
              sx={{
                '&.Mui-selected': {
                  backgroundColor: 'secondary.main',
                  color: 'secondary.contrastText',
                  '&:hover': {
                    backgroundColor: 'secondary.dark',
                  },
                },
              }}
            >
              <ListItemIcon
                sx={{
                  color: currentView === item.view ? 'inherit' : 'text.secondary',
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.label} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      
      <Divider sx={{ mt: 2 }} />
      
      {/* Breadcrumb dans le sidebar pour mobile */}
      {(selectedRole || selectedChampion) && (
        <Box sx={{ p: 2 }}>
          <Typography variant="caption" color="text.secondary" sx={{ mb: 1, display: 'block' }}>
            Navigation actuelle:
          </Typography>
          <Typography variant="body2" sx={{ fontWeight: 500 }}>
            {getBreadcrumb()}
          </Typography>
          
          {selectedChampion && (
            <Box sx={{ mt: 1 }}>
              <Typography 
                variant="body2" 
                sx={{ 
                  color: 'primary.main', 
                  cursor: 'pointer',
                  textDecoration: 'underline',
                }}
                onClick={() => {
                  onChampionSelect?.(undefined);
                  onViewChange('champions');
                }}
              >
                ← Retour aux champions
              </Typography>
            </Box>
          )}
          
          {selectedRole && !selectedChampion && (
            <Box sx={{ mt: 1 }}>
              <Typography 
                variant="body2" 
                sx={{ 
                  color: 'primary.main', 
                  cursor: 'pointer',
                  textDecoration: 'underline',
                }}
                onClick={() => {
                  onRoleSelect?.(undefined);
                  onViewChange('roles');
                }}
              >
                ← Retour aux rôles
              </Typography>
            </Box>
          )}
        </Box>
      )}
    </Box>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      {/* AppBar */}
      <AppBar
        position="fixed"
        sx={{
          width: { lg: `calc(100% - ${DRAWER_WIDTH}px)` },
          ml: { lg: `${DRAWER_WIDTH}px` },
          backgroundColor: 'background.paper',
          color: 'text.primary',
          boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
        }}
      >
        <Toolbar>
          <IconButton
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { lg: 'none' } }}
            aria-label="Ouvrir le menu de navigation"
          >
            <MenuIcon />
          </IconButton>
          
          <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 600 }}>
            {getBreadcrumb()}
          </Typography>
          
          <Tooltip title="Export global">
            <IconButton 
              sx={{ mr: 1 }} 
              onClick={handleGlobalExport}
              aria-label="Exporter la vue actuelle en PNG et Excel"
            >
              <GetApp />
            </IconButton>
          </Tooltip>
          
          <Tooltip title="Paramètres">
            <IconButton 
              sx={{ mr: 1 }}
              aria-label="Ouvrir les paramètres"
            >
              <Settings />
            </IconButton>
          </Tooltip>
          
          <Tooltip title={`Passer au thème ${isDarkMode ? 'clair' : 'sombre'}`}>
            <IconButton 
              onClick={toggleTheme}
              aria-label={`Activer le thème ${isDarkMode ? 'clair' : 'sombre'}`}
            >
              {isDarkMode ? <Brightness7 /> : <Brightness4 />}
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBar>

      {/* Drawer */}
      <Box
        component="nav"
        sx={{ width: { lg: DRAWER_WIDTH }, flexShrink: { lg: 0 } }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{ keepMounted: true }}
          sx={{
            display: { xs: 'block', lg: 'none' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: DRAWER_WIDTH,
              backgroundColor: 'background.paper',
            },
          }}
        >
          {drawer}
        </Drawer>
        
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', lg: 'block' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: DRAWER_WIDTH,
              backgroundColor: 'background.paper',
              borderRight: '1px solid',
              borderColor: 'divider',
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>

      {/* Contenu principal */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          width: { lg: `calc(100% - ${DRAWER_WIDTH}px)` },
          minHeight: '100vh',
          backgroundColor: 'background.default',
        }}
      >
        <Toolbar />
        <Container maxWidth="xl" sx={{ py: 3 }}>
          {children}
        </Container>
      </Box>
    </Box>
  );
};