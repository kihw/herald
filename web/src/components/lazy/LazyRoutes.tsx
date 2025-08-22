import { lazy } from 'react';
import { withLazy } from '../common/LazyComponent';

// Lazy load group components for better performance
export const LazyGroupManagement = withLazy(
  () => import('../groups/GroupManagement'),
  {
    minLoadTime: 200,
  }
);

export const LazyCreateGroupDialog = withLazy(
  () => import('../groups/CreateGroupDialog'),
  {
    minLoadTime: 100,
  }
);

export const LazyJoinGroupDialog = withLazy(
  () => import('../groups/JoinGroupDialog'),
  {
    minLoadTime: 100,
  }
);

export const LazyGroupDetailsDialog = withLazy(
  () => import('../groups/GroupDetailsDialog'),
  {
    minLoadTime: 150,
  }
);

export const LazyComparisonManager = withLazy(
  () => import('../groups/ComparisonManager'),
  {
    minLoadTime: 200,
  }
);

export const LazyComparisonCharts = withLazy(
  () => import('../charts/ComparisonCharts'),
  {
    minLoadTime: 300,
  }
);

export const LazyGroupStatsCharts = withLazy(
  () => import('../charts/GroupStatsCharts'),
  {
    minLoadTime: 300,
  }
);

export const LazyPlayerPerformanceWidget = withLazy(
  () => import('../charts/PlayerPerformanceWidget'),
  {
    minLoadTime: 150,
  }
);

// Dashboard components
export const LazyMainDashboard = withLazy(
  () => import('../dashboard/MainDashboard'),
  {
    minLoadTime: 200,
  }
);

export const LazyAnalyticsDashboard = withLazy(
  () => import('../dashboard/AnalyticsDashboard'),
  {
    minLoadTime: 250,
  }
);

export const LazyEnhancedDashboard = withLazy(
  () => import('../dashboard/EnhancedDashboard'),
  {
    minLoadTime: 250,
  }
);

// Auth components
export const LazyAuthPage = withLazy(
  () => import('../auth/AuthPage'),
  {
    minLoadTime: 100,
  }
);

export const LazyGoogleAuth = withLazy(
  () => import('../auth/GoogleAuth'),
  {
    minLoadTime: 100,
  }
);

export const LazyRiotValidationForm = withLazy(
  () => import('../auth/RiotValidationForm'),
  {
    minLoadTime: 150,
  }
);

// Analytics components
export const LazyChampionsAnalytics = withLazy(
  () => import('../ChampionsAnalytics'),
  {
    minLoadTime: 300,
  }
);

export const LazyMMRAnalytics = withLazy(
  () => import('../MMRAnalytics'),
  {
    minLoadTime: 300,
  }
);

export const LazyExporterMUI = withLazy(
  () => import('../ExporterMUI'),
  {
    minLoadTime: 150,
  }
);

// Export all for convenience
export const LazyComponents = {
  // Groups
  GroupManagement: LazyGroupManagement,
  CreateGroupDialog: LazyCreateGroupDialog,
  JoinGroupDialog: LazyJoinGroupDialog,
  GroupDetailsDialog: LazyGroupDetailsDialog,
  ComparisonManager: LazyComparisonManager,
  
  // Charts
  ComparisonCharts: LazyComparisonCharts,
  GroupStatsCharts: LazyGroupStatsCharts,
  PlayerPerformanceWidget: LazyPlayerPerformanceWidget,
  
  // Dashboard
  MainDashboard: LazyMainDashboard,
  AnalyticsDashboard: LazyAnalyticsDashboard,
  EnhancedDashboard: LazyEnhancedDashboard,
  
  // Auth
  AuthPage: LazyAuthPage,
  GoogleAuth: LazyGoogleAuth,
  RiotValidationForm: LazyRiotValidationForm,
  
  // Analytics
  ChampionsAnalytics: LazyChampionsAnalytics,
  MMRAnalytics: LazyMMRAnalytics,
  ExporterMUI: LazyExporterMUI,
};

export default LazyComponents;