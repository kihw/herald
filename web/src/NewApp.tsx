import React, { useState, useMemo } from 'react';
import { ThemeContextProvider } from './context/ThemeContext';
import { ExportFeedbackProvider } from './components/ExportFeedback';
import { AnimationProvider } from './components/animations/AnimationProvider';
import { EnhancedErrorBoundary } from './components/errors/EnhancedErrorBoundary';
import { GlobalErrorProvider } from './components/errors/GlobalErrorHandler';
import { AppLayout, ViewType } from './components/layout/AppLayout';
import { Overview } from './views/Overview';
import { EnhancedOverview } from './views/EnhancedOverview';
import { RolesView } from './views/RolesView';
import { ChampionsView } from './views/ChampionsView';
import { ChampionDetails } from './views/ChampionDetails';
import { ComparatorView } from './views/ComparatorView';
import { HeatmapView } from './views/HeatmapView';
import { ExporterMUI } from './components/ExporterMUI';
import { Dashboard } from './components/Dashboard';
import { EnhancedDashboard } from './components/dashboard/EnhancedDashboard';
import { ChampionsAnalytics } from './components/ChampionsAnalytics';
import { MMRAnalytics } from './components/MMRAnalytics';
import { Row } from './types';
import { Box, Fade, Alert, AlertTitle } from '@mui/material';

// Types pour la gestion des données
export interface AppState {
  rows: Row[];
  loading: boolean;
  error: string | null;
  currentPuuid: string | null; // For analytics
}

export default function NewApp() {
  // État principal des données
  const [appState, setAppState] = useState<AppState>({
    rows: [],
    loading: false,
    error: null,
    currentPuuid: null,
  });

  // État de navigation
  const [currentView, setCurrentView] = useState<ViewType>('overview');
  const [selectedRole, setSelectedRole] = useState<string | undefined>();
  const [selectedChampion, setSelectedChampion] = useState<string | undefined>();

  // État de gestion d'erreurs globales
  const [retryCount, setRetryCount] = useState(0);
  const [lastErrorTime, setLastErrorTime] = useState<number | null>(null);

  // Gestion des changements de vue avec reset des sélections
  const handleViewChange = (view: ViewType) => {
    setCurrentView(view);
    
    // Reset des sélections selon la vue
    if (view === 'overview') {
      setSelectedRole(undefined);
      setSelectedChampion(undefined);
    } else if (view === 'roles') {
      setSelectedChampion(undefined);
    }
  };

  const handleRoleSelect = (role: string | undefined) => {
    setSelectedRole(role);
    if (role) {
      setCurrentView('champions');
      setSelectedChampion(undefined);
    }
  };

  const handleChampionSelect = (champion: string | undefined) => {
    setSelectedChampion(champion);
    if (champion) {
      setCurrentView('details');
    }
  };

  // Données filtrées selon les sélections
  const filteredData = useMemo(() => {
    let data = appState.rows;
    
    if (selectedRole) {
      data = data.filter(row => row.lane === selectedRole);
    }
    
    if (selectedChampion) {
      data = data.filter(row => row.champion === selectedChampion);
    }
    
    return data;
  }, [appState.rows, selectedRole, selectedChampion]);

  // Gestion du chargement des données depuis l'exporteur
  const handleLoadRows = (rows: Row[], puuid?: string) => {
    setAppState(prev => ({
      ...prev,
      rows,
      loading: false,
      error: null,
      currentPuuid: puuid || prev.currentPuuid,
    }));
  };

  const handleLoadingChange = (loading: boolean) => {
    setAppState(prev => ({ ...prev, loading }));
  };

  const handleErrorChange = (error: string | null) => {
    setAppState(prev => ({ ...prev, error }));
    if (error) {
      setLastErrorTime(Date.now());
    }
  };

  // Gestion des erreurs globales avec retry
  const handleGlobalError = (error: Error, errorInfo: any, errorId: string) => {
    console.error('Application Error:', error, errorInfo, { errorId, retryCount });
    
    // Analytics d'erreur (optionnel)
    if (typeof window !== 'undefined' && (window as any).gtag) {
      (window as any).gtag('event', 'exception', {
        description: `${error.name}: ${error.message}`,
        fatal: false,
        error_id: errorId,
        retry_count: retryCount,
      });
    }
  };

  const handleRetry = (currentRetryCount: number) => {
    setRetryCount(currentRetryCount);
    console.log(`Application retry attempt ${currentRetryCount}`);
    
    // Reset de l'erreur après un délai
    setTimeout(() => {
      setAppState(prev => ({ ...prev, error: null }));
    }, 1000);
  };

  const handleMaxRetriesReached = (error: Error) => {
    console.error('Max retries reached for application error:', error);
    
    // Fallback vers une page d'erreur ou reload complet
    const shouldReload = window.confirm(
      'L\'application a rencontré une erreur persistante. Voulez-vous recharger la page ?'
    );
    
    if (shouldReload) {
      window.location.reload();
    }
  };

  // Rendu conditionnel des vues
  const renderCurrentView = () => {
    const commonProps = {
      data: filteredData,
      loading: appState.loading,
      error: appState.error,
    };

    switch (currentView) {
      case 'overview':
        return (
          <Fade in timeout={300}>
            <Box>
              <EnhancedOverview 
                {...commonProps}
                allData={appState.rows}
                onRoleSelect={handleRoleSelect}
              />
            </Box>
          </Fade>
        );
        
      case 'roles':
        return (
          <Fade in timeout={300}>
            <Box>
              <RolesView 
                {...commonProps}
                onRoleSelect={handleRoleSelect}
              />
            </Box>
          </Fade>
        );
        
      case 'champions':
        return (
          <Fade in timeout={300}>
            <Box>
              <ChampionsView 
                {...commonProps}
                selectedRole={selectedRole}
                onChampionSelect={handleChampionSelect}
              />
            </Box>
          </Fade>
        );
        
      case 'details':
        return (
          <Fade in timeout={300}>
            <Box>
              <ChampionDetails 
                {...commonProps}
                selectedChampion={selectedChampion}
                selectedRole={selectedRole}
              />
            </Box>
          </Fade>
        );

      case 'analytics-dashboard':
        return (
          <Fade in timeout={300}>
            <Box>
              {appState.currentPuuid ? (
                <EnhancedDashboard puuid={appState.currentPuuid} />
              ) : (
                <Box sx={{ p: 3, textAlign: 'center' }}>
                  <h3>Effectuez d'abord un export pour accéder aux analytics</h3>
                  <p>Le dashboard analytics nécessite des données d'un utilisateur spécifique.</p>
                </Box>
              )}
            </Box>
          </Fade>
        );

      case 'analytics-champions':
        return (
          <Fade in timeout={300}>
            <Box>
              {appState.currentPuuid ? (
                <ChampionsAnalytics puuid={appState.currentPuuid} />
              ) : (
                <Box sx={{ p: 3, textAlign: 'center' }}>
                  <h3>Effectuez d'abord un export pour accéder aux analytics</h3>
                  <p>L'analyse des champions nécessite des données d'un utilisateur spécifique.</p>
                </Box>
              )}
            </Box>
          </Fade>
        );

      case 'analytics-mmr':
        return (
          <Fade in timeout={300}>
            <Box>
              {appState.currentPuuid ? (
                <MMRAnalytics puuid={appState.currentPuuid} />
              ) : (
                <Box sx={{ p: 3, textAlign: 'center' }}>
                  <h3>Effectuez d'abord un export pour accéder aux analytics</h3>
                  <p>L'analyse MMR nécessite des données d'un utilisateur spécifique.</p>
                </Box>
              )}
            </Box>
          </Fade>
        );

      case 'comparator':
        return (
          <Fade in timeout={300}>
            <Box>
              <ComparatorView />
            </Box>
          </Fade>
        );

      case 'heatmap':
        return (
          <Fade in timeout={300}>
            <Box>
              <HeatmapView 
                data={appState.rows}
                loading={appState.loading}
                error={appState.error}
              />
            </Box>
          </Fade>
        );
        
      default:
        return (
          <Fade in timeout={300}>
            <Box>
              <Overview 
                {...commonProps}
                allData={appState.rows}
                onRoleSelect={handleRoleSelect}
              />
            </Box>
          </Fade>
        );
    }
  };

  return (
    <ThemeContextProvider>
      <GlobalErrorProvider
        enableNotifications={true}
        enableFloatingButton={true}
        maxErrors={20}
        onError={(error) => {
          console.log('Global error captured:', error);
          // Optionnel: envoi à un service de monitoring
        }}
      >
        <AnimationProvider initialConfig={{ speed: 'normal', enabled: true }}>
          <ExportFeedbackProvider>
            <EnhancedErrorBoundary
              maxRetries={3}
              retryDelay={2000}
              autoRetry={true}
              onError={handleGlobalError}
              onRetry={handleRetry}
              onMaxRetriesReached={handleMaxRetriesReached}
              enableErrorReporting={true}
              showErrorDetails={true}
              errorLevel="high"
            >
              <AppLayout
                currentView={currentView}
                onViewChange={handleViewChange}
                selectedRole={selectedRole}
                selectedChampion={selectedChampion}
                onRoleSelect={handleRoleSelect}
                onChampionSelect={handleChampionSelect}
                data={filteredData}
              >
                {/* Interface d'export - toujours visible */}
                <Box sx={{ mb: 3 }}>
                  <EnhancedErrorBoundary
                    maxRetries={2}
                    retryDelay={1000}
                    autoRetry={false}
                    errorLevel="medium"
                    fallback={
                      <Alert severity="warning" sx={{ mb: 2 }}>
                        <AlertTitle>Erreur d'export</AlertTitle>
                        L'interface d'export est temporairement indisponible. 
                        Veuillez rafraîchir la page.
                      </Alert>
                    }
                  >
                    <ExporterMUI 
                      onLoadRows={handleLoadRows}
                      onLoadingChange={handleLoadingChange}
                      onErrorChange={handleErrorChange}
                    />
                  </EnhancedErrorBoundary>
                </Box>

                {/* Vue principale avec protection d'erreur */}
                <EnhancedErrorBoundary
                  maxRetries={2}
                  retryDelay={1500}
                  autoRetry={true}
                  errorLevel="medium"
                  onError={(error, errorInfo, errorId) => {
                    console.error(`View error (${currentView}):`, error, errorInfo, errorId);
                  }}
                >
                  {renderCurrentView()}
                </EnhancedErrorBoundary>
              </AppLayout>
            </EnhancedErrorBoundary>
          </ExportFeedbackProvider>
        </AnimationProvider>
      </GlobalErrorProvider>
    </ThemeContextProvider>
  );
}