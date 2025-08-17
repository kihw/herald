import React, { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { useApi } from '../hooks/useApi';

interface DashboardStats {
  total_matches: number;
  win_rate: number;
  average_kda: number;
  favorite_champion: string;
  last_sync_at?: string;
  next_sync_at?: string;
}

export const DashboardPage: React.FC = () => {
  const { user, logout } = useAuth();
  const { apiCall } = useApi();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [syncLoading, setSyncLoading] = useState(false);
  const [syncError, setSyncError] = useState('');
  const [cooldownTime, setCooldownTime] = useState(0);

  useEffect(() => {
    loadStats();
  }, []);

  useEffect(() => {
    // Countdown timer for sync cooldown
    if (cooldownTime > 0) {
      const timer = setTimeout(() => setCooldownTime(cooldownTime - 1), 1000);
      return () => clearTimeout(timer);
    }
  }, [cooldownTime]);

  const loadStats = async () => {
    try {
      const response = await apiCall('/api/dashboard/stats', 'GET');
      setStats(response);
    } catch (error) {
      console.error('Failed to load stats:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSync = async () => {
    setSyncLoading(true);
    setSyncError('');
    
    try {
      const response = await apiCall('/api/sync', 'POST', {});
      
      if (response.status === 'started') {
        // Sync started successfully
        await loadStats(); // Refresh stats
        setCooldownTime(120); // 2 minutes cooldown
      } else if (response.status === 'cooldown') {
        setSyncError(`Veuillez attendre ${Math.ceil(response.remaining_seconds / 60)} minutes avant la prochaine synchronisation`);
      }
    } catch (error) {
      setSyncError(error instanceof Error ? error.message : 'Échec de la synchronisation');
    } finally {
      setSyncLoading(false);
    }
  };

  const formatCooldownTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  if (loading) {
    return (
      <div className="dashboard-loading">
        <div className="spinner"></div>
        <p>Chargement de vos statistiques...</p>
      </div>
    );
  }

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div className="user-info">
          <h1>Bienvenue, {user?.username}#{user?.tagline}</h1>
          <p>Dernière synchronisation: {stats?.last_sync_at ? new Date(stats.last_sync_at).toLocaleString('fr-FR') : 'Jamais'}</p>
        </div>
        <div className="header-actions">
          <button 
            onClick={handleSync} 
            disabled={syncLoading || cooldownTime > 0}
            className="sync-button"
          >
            {syncLoading ? 'Synchronisation...' : 
             cooldownTime > 0 ? `Actualiser (${formatCooldownTime(cooldownTime)})` : 
             'Actualiser les données'}
          </button>
          <button onClick={logout} className="logout-button">
            Déconnexion
          </button>
        </div>
      </header>

      {syncError && (
        <div className="error-banner">
          {syncError}
        </div>
      )}

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-number">{stats?.total_matches || 0}</div>
          <div className="stat-label">Matchs totaux</div>
        </div>

        <div className="stat-card">
          <div className="stat-number">{stats?.win_rate?.toFixed(1) || 0}%</div>
          <div className="stat-label">Taux de victoire</div>
        </div>

        <div className="stat-card">
          <div className="stat-number">{stats?.average_kda?.toFixed(2) || 0}</div>
          <div className="stat-label">KDA moyen</div>
        </div>

        <div className="stat-card">
          <div className="stat-number">{stats?.favorite_champion || 'Aucun'}</div>
          <div className="stat-label">Champion favori</div>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="recent-matches">
          <h2>Matchs récents</h2>
          <div className="matches-placeholder">
            <p>Vos matchs récents apparaîtront ici après synchronisation</p>
          </div>
        </div>

        <div className="quick-actions">
          <h2>Actions rapides</h2>
          <div className="actions-grid">
            <button className="action-card">
              <h3>Paramètres</h3>
              <p>Configurer vos préférences d'export</p>
            </button>
            
            <button className="action-card">
              <h3>Profil</h3>
              <p>Modifier vos informations</p>
            </button>
            
            <button className="action-card">
              <h3>Historique</h3>
              <p>Voir tous vos matchs</p>
            </button>
          </div>
        </div>
      </div>

      {stats?.next_sync_at && (
        <div className="next-sync-info">
          <p>Prochaine synchronisation automatique: {new Date(stats.next_sync_at).toLocaleString('fr-FR')}</p>
        </div>
      )}
    </div>
  );
};
