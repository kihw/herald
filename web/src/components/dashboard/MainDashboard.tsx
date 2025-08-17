import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  Tabs,
  Tab,
  Alert,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  History,
  Settings,
  Sync,
  Speed,
  Analytics,
  CloudDownload,
} from '@mui/icons-material';

import StatsCard from './StatsCard';
import MatchesTable from './MatchesTable';
import SettingsPanel from './SettingsPanel';
import AnalyticsDashboard from './AnalyticsDashboard';
import ToastNotification from '../common/ToastNotification';
import LoadingSkeleton from '../common/LoadingSkeleton';
import HelpPanel from '../common/HelpPanel';
import PerformanceMonitor from '../common/PerformanceMonitor';
import ErrorBoundary from '../common/ErrorBoundary';
import { AdvancedExporter } from '../AdvancedExporter';
import { useAuth } from '../../context/AuthContext';
import { apiService } from '../../services/api';
import { useOptimizedData } from '../../utils/DataProcessor';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`dashboard-tabpanel-${index}`}
      aria-labelledby={`dashboard-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
}

const MainDashboard: React.FC = () => {
  const { state } = useAuth();
  const { user } = state;
  const [tabValue, setTabValue] = useState(0);
  
  // Stats state
  const [stats, setStats] = useState<any>(null);
  const [statsLoading, setStatsLoading] = useState(true);
  
  // Matches state
  const [matches, setMatches] = useState([]);
  const [matchesLoading, setMatchesLoading] = useState(false);
  const [syncLoading, setSyncLoading] = useState(false);
  
  // Settings state
  const [settings, setSettings] = useState({
    include_timeline: true,
    include_all_data: true,
    light_mode: true,
    auto_sync_enabled: true,
    sync_frequency_hours: 24,
  });
  const [settingsLoading, setSettingsLoading] = useState(false);
  const [settingsSaving, setSettingsSaving] = useState(false);
  
  // Notification state
  const [notification, setNotification] = useState({
    open: false,
    message: '',
    title: '',
    severity: 'success' as 'success' | 'error' | 'warning' | 'info',
  });

  // Performance monitoring state
  const [showPerformanceStats, setShowPerformanceStats] = useState(false);
  const [performanceMetrics, setPerformanceMetrics] = useState<any[]>([]);

  // Use optimized data processing
  const { 
    matches: optimizedMatches, 
    championStats, 
    insights, 
    performanceSummary,
    cacheStats 
  } = useOptimizedData(matches);

  // Load initial data
  useEffect(() => {
    loadStats();
    loadSettings();
  }, []);

  const showNotification = (message: string, severity: 'success' | 'error' | 'warning' | 'info' = 'success', title?: string) => {
    setNotification({ open: true, message, title: title || '', severity });
  };

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };

  // Performance metrics handler
  const handlePerformanceMetrics = (metrics: any) => {
    setPerformanceMetrics(prev => [...prev.slice(-9), metrics]); // Keep last 10 metrics
    
    // Show warning for slow operations
    if (metrics.renderTime > 100) {
      showNotification(
        `Slow render detected: ${metrics.renderTime.toFixed(2)}ms`,
        'warning',
        'Performance Warning'
      );
    }
  };

  // Toggle performance stats display
  const togglePerformanceStats = () => {
    setShowPerformanceStats(prev => !prev);
    if (!showPerformanceStats) {
      const memUsage = (performance as any).memory?.usedJSHeapSize 
        ? ((performance as any).memory.usedJSHeapSize / 1024 / 1024).toFixed(1) + 'MB' 
        : 'N/A';
      const stats = cacheStats;
      showNotification(
        `Cache: ${stats.size} entries | Memory: ${memUsage}`,
        'info',
        'Performance Stats'
      );
    }
  };

  const loadStats = async () => {
    try {
      setStatsLoading(true);
      const response = await apiService.getDashboardStats();
      
      // Transform API response to match StatsCard interface
      const transformedStats = {
        totalMatches: response.total_matches || 0,
        winRate: response.win_rate || 0,
        mainRole: "ADC", // Mock data for now
        favoriteChampion: {
          name: response.favorite_champion || "Jinx",
          winRate: 72.3,
          matches: 23,
        },
        recentPerformance: {
          last7Days: {
            matches: 8,
            wins: 6,
            winRate: 75.0,
          },
          last30Days: {
            matches: 34,
            wins: 22,
            winRate: 64.7,
          },
        },
        rankInfo: {
          tier: "Gold",
          division: "II",
          lp: 67,
        },
      };
      
      setStats(transformedStats);
    } catch (error) {
      console.error('Failed to load stats:', error);
      showNotification('Failed to load statistics', 'error');
    } finally {
      setStatsLoading(false);
    }
  };

  const loadMatches = async () => {
    try {
      setMatchesLoading(true);
      const response = await apiService.getMatches();
      setMatches(response.matches || []);
    } catch (error) {
      console.error('Failed to load matches:', error);
      showNotification('Failed to load match history', 'error');
    } finally {
      setMatchesLoading(false);
    }
  };

  const loadSettings = async () => {
    try {
      setSettingsLoading(true);
      const response = await apiService.getSettings();
      setSettings(response);
    } catch (error) {
      console.error('Failed to load settings:', error);
      showNotification('Failed to load settings', 'error');
    } finally {
      setSettingsLoading(false);
    }
  };

  const handleSyncMatches = async () => {
    try {
      setSyncLoading(true);
      console.log('🔄 Starting match synchronization...');
      const response = await apiService.syncMatches();
      console.log('✅ Sync response:', response);
      
      if (response.success) {
        const { newMatches, totalMatches } = response.data;
        console.log(`📊 Sync completed: ${newMatches} new matches, ${totalMatches} total`);
        showNotification(
          `Synchronized ${newMatches} new matches! Total: ${totalMatches}`, 
          'success', 
          'Sync Complete'
        );
      } else {
        console.warn('⚠️ Sync returned success=false:', response);
        showNotification('Sync completed but no data returned', 'warning');
      }
      
      // Reload matches and stats after sync
      console.log('🔄 Reloading matches and stats...');
      await Promise.all([loadMatches(), loadStats()]);
      console.log('✅ Data reloaded successfully');
    } catch (error) {
      console.error('❌ Failed to sync matches:', error);
      showNotification('Failed to sync matches', 'error');
    } finally {
      setSyncLoading(false);
    }
  };

  const handleSaveSettings = async (newSettings: typeof settings) => {
    try {
      setSettingsSaving(true);
      await apiService.updateSettings(newSettings);
      setSettings(newSettings);
      showNotification('Your preferences have been updated successfully!', 'success', 'Settings Saved');
    } catch (error) {
      console.error('Failed to save settings:', error);
      showNotification('Failed to save settings', 'error');
    } finally {
      setSettingsSaving(false);
    }
  };

  const handleViewMatch = (matchId: string) => {
    // TODO: Implement match detail view
    showNotification(`Viewing match ${matchId} (coming soon)`, 'info');
  };

  const handleExportMatches = () => {
    // TODO: Implement match export
    showNotification('Export feature coming soon!', 'info');
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
    
    // Load data when switching to matches tab
    if (newValue === 1 && matches.length === 0) {
      loadMatches();
    }
  };

  if (!user) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error">
          Please log in to access the dashboard.
        </Alert>
      </Container>
    );
  }

  return (
    <ErrorBoundary 
      onError={(error, errorInfo) => {
        console.error('Dashboard Error:', error, errorInfo);
        showNotification('An unexpected error occurred', 'error', 'Dashboard Error');
      }}
    >
      {/* Temporarily disabled PerformanceMonitor to fix memory leak */}
      {/* <PerformanceMonitor 
        componentName="MainDashboard"
        onMetricsUpdate={handlePerformanceMetrics}
      /> */}
      <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* Header */}
      <Box sx={{ mb: 4, position: 'relative' }}>
        <Typography variant="h4" component="h1" sx={{ mb: 1 }}>
          Welcome back, {user.riot_id}
        </Typography>
        <Typography variant="body1" color="text.secondary">
          #{user.riot_tag} • {user.region.toUpperCase()}
        </Typography>
        
        {/* Performance Monitor Button (Development only) */}
        {process.env.NODE_ENV === 'development' && (
          <Tooltip title="Performance Stats">
            <IconButton
              onClick={togglePerformanceStats}
              sx={{ 
                position: 'absolute', 
                top: 0, 
                right: 0,
                color: showPerformanceStats ? 'primary.main' : 'text.secondary'
              }}
            >
              <Speed />
            </IconButton>
          </Tooltip>
        )}
      </Box>

      {/* Help Panel */}
      <HelpPanel />

      {/* Navigation Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab
            icon={<DashboardIcon />}
            label="Overview"
            iconPosition="start"
          />
          <Tab
            icon={<History />}
            label="Match History"
            iconPosition="start"
          />
          <Tab
            icon={<Analytics />}
            label="Analytics"
            iconPosition="start"
          />
          <Tab
            icon={<CloudDownload />}
            label="Export Avancé"
            iconPosition="start"
          />
          <Tab
            icon={<Settings />}
            label="Settings"
            iconPosition="start"
          />
        </Tabs>
      </Box>

      {/* Tab Panels */}
      <TabPanel value={tabValue} index={0}>
        {statsLoading ? (
          <LoadingSkeleton type="stats" />
        ) : stats ? (
          <StatsCard stats={stats} loading={statsLoading} />
        ) : (
          <LoadingSkeleton type="stats" />
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        {matchesLoading ? (
          <LoadingSkeleton type="matches" count={5} />
        ) : (
          <MatchesTable
            matches={matches}
            loading={matchesLoading}
            onSync={handleSyncMatches}
            onViewMatch={handleViewMatch}
            onExport={handleExportMatches}
            syncLoading={syncLoading}
          />
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        <AnalyticsDashboard />
      </TabPanel>

      <TabPanel value={tabValue} index={3}>
        <AdvancedExporter 
          onLoadingChange={(loading) => {
            // Optionnel: gérer l'état de chargement global
          }}
          onErrorChange={(error) => {
            if (error) {
              setNotification({ type: 'error', message: error });
            }
          }}
        />
      </TabPanel>

      <TabPanel value={tabValue} index={4}>
        {settingsLoading ? (
          <LoadingSkeleton type="settings" />
        ) : (
          <SettingsPanel
            settings={settings}
            onSave={handleSaveSettings}
            loading={settingsLoading}
            saving={settingsSaving}
          />
        )}
      </TabPanel>

      {/* Notifications */}
      <ToastNotification
        open={notification.open}
        message={notification.message}
        title={notification.title}
        severity={notification.severity}
        onClose={handleCloseNotification}
      />
    </Container>
    </ErrorBoundary>
  );
};

export default MainDashboard;
