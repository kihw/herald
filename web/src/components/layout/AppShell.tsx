import React, { useState, useEffect } from 'react';
import { Box } from '@mui/material';
import AppLayout from './AppLayout';
import { ViewType } from './AppLayout';
import { LazyComponents } from '../lazy/LazyRoutes';
import { preloadRoute } from '../../utils/performance';
import useResponsive from '../../hooks/useResponsive';
import { usePerformance } from '../../hooks/usePerformance';

const AppShell: React.FC = () => {
  const [currentView, setCurrentView] = useState<ViewType>('overview');
  const [selectedRole, setSelectedRole] = useState<string | undefined>();
  const [selectedChampion, setSelectedChampion] = useState<string | undefined>();
  
  const { isMobile } = useResponsive();
  const { shouldLazyLoad } = usePerformance();

  // Preload routes based on user interaction patterns
  useEffect(() => {
    if (!shouldLazyLoad) {
      // Preload likely next routes
      const timer = setTimeout(() => {
        preloadRoute('groups');
        if (!isMobile) {
          preloadRoute('charts');
          preloadRoute('dashboard');
        }
      }, 2000);

      return () => clearTimeout(timer);
    }
  }, [shouldLazyLoad, isMobile]);

  const handleViewChange = (view: ViewType) => {
    setCurrentView(view);
    
    // Reset selections when changing main views
    if (view === 'overview' || view.startsWith('analytics-') || view === 'groups') {
      setSelectedRole(undefined);
      setSelectedChampion(undefined);
    }

    // Preload related routes
    if (view === 'groups' && !shouldLazyLoad) {
      preloadRoute('charts');
    }
  };

  const handleRoleSelect = (role: string | undefined) => {
    setSelectedRole(role);
    if (role) {
      setCurrentView('champions');
    }
  };

  const handleChampionSelect = (champion: string | undefined) => {
    setSelectedChampion(champion);
    if (champion) {
      setCurrentView('details');
    }
  };

  const renderContent = () => {
    switch (currentView) {
      case 'groups':
        return <LazyComponents.GroupManagement />;
      case 'analytics-dashboard':
        return <LazyComponents.AnalyticsDashboard />;
      case 'analytics-champions':
        return <LazyComponents.ChampionsAnalytics />;
      case 'analytics-mmr':
        return <LazyComponents.MMRAnalytics />;
      default:
        return <LazyComponents.MainDashboard />;
    }
  };

  return (
    <AppLayout
      currentView={currentView}
      onViewChange={handleViewChange}
      selectedRole={selectedRole}
      selectedChampion={selectedChampion}
      onRoleSelect={handleRoleSelect}
      onChampionSelect={handleChampionSelect}
    >
      {renderContent()}
    </AppLayout>
  );
};

export default AppShell;