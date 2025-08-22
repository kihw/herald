import React, { createContext, useContext, ReactNode } from 'react';
import { QueryClient, QueryClientProvider, useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiService, User, AuthTokens } from '../services/api';

// Types pour l'état d'authentification sécurisé
export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  tokens: AuthTokens | null;
  lastActivity: number | null;
  suspiciousActivity: boolean;
}

// Surveillance sécurisée
interface SecurityMonitoring {
  failedAttempts: number;
  lastFailedAttempt: number | null;
  suspiciousPatterns: string[];
}

// Interface du contexte sécurisé
interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  validateAccount: (riotId: string, riotTag: string, region: string) => void;
  logout: () => void;
  clearError: () => void;
  securityStatus: {
    failedAttempts: number;
    suspiciousActivity: boolean;
    lastActivity: number | null;
  };
}

// Client React Query sécurisé
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error: any) => {
        // Ne pas retry sur les erreurs d'auth
        if (error?.message?.includes('Non autorisé') || error?.message?.includes('Session expirée')) {
          return false;
        }
        return failureCount < 2;
      },
      staleTime: 5 * 60 * 1000, // 5 minutes
      cacheTime: 10 * 60 * 1000, // 10 minutes
    },
    mutations: {
      retry: 1,
    },
  },
});

// Surveillance de sécurité globale
let securityMonitoring: SecurityMonitoring = {
  failedAttempts: 0,
  lastFailedAttempt: null,
  suspiciousPatterns: [],
};

// Fonction de détection d'activité suspecte
function detectSuspiciousActivity(error: string): boolean {
  const suspiciousPatterns = [
    'brute force',
    'too many requests',
    'rate limit',
    'suspicious activity',
    'multiple failed',
  ];
  
  return suspiciousPatterns.some(pattern => 
    error.toLowerCase().includes(pattern)
  );
}

// Création du contexte sécurisé
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Gestionnaire d'état sécurisé pour éviter les données corrompues
interface SecureAuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  lastUpdate: number;
}

// État par défaut stable
const defaultAuthState: SecureAuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: true,
  error: null,
  lastUpdate: Date.now()
};

// Hook sécurisé pour la session utilisateur avec gestion d'état robuste
function useSecureSession() {
  const [authState, setAuthState] = React.useState<SecureAuthState>(defaultAuthState);

  const query = useQuery({
    queryKey: ['auth-session'],
    queryFn: async () => {
      try {
        const response = await apiService.checkSession();
        
        // Validation stricte des données reçues
        if (response && typeof response === 'object') {
          if (response.authenticated && response.user && typeof response.user === 'object') {
            // Validation des champs utilisateur requis
            if (response.user.riot_id && response.user.riot_tag && response.user.region) {
              securityMonitoring.failedAttempts = 0;
              return {
                user: response.user,
                isAuthenticated: true,
                isLoading: false,
                error: null,
                lastUpdate: Date.now()
              };
            }
          }
        }
        
        // État non-authentifié stable
        return {
          user: null,
          isAuthenticated: false,
          isLoading: false,
          error: null,
          lastUpdate: Date.now()
        };
      } catch (error) {
        console.error('Erreur de vérification de session:', error);
        return {
          user: null,
          isAuthenticated: false,
          isLoading: false,
          error: error instanceof Error ? error.message : 'Erreur de session',
          lastUpdate: Date.now()
        };
      }
    },
    staleTime: 2 * 60 * 1000,
    cacheTime: 5 * 60 * 1000,
    // Callbacks pour synchroniser l'état local
    onSuccess: (data) => {
      if (data && typeof data === 'object') {
        setAuthState(data);
      }
    },
    onError: (error) => {
      setAuthState({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Erreur de connexion',
        lastUpdate: Date.now()
      });
    }
  });

  // Retourner l'état local stable au lieu des données React Query directes
  return {
    data: authState.user,
    isLoading: query.isLoading || authState.isLoading,
    error: query.error || authState.error,
    authState
  };
}

// Provider du contexte d'authentification sécurisé
interface AuthProviderProps {
  children: ReactNode;
}

function AuthProviderContent({ children }: AuthProviderProps) {
  const queryClient = useQueryClient();
  const { data: user, isLoading, error, authState } = useSecureSession();
  const [authError, setAuthError] = React.useState<string | null>(null);
  const [stableState, setStableState] = React.useState<SecureAuthState>(defaultAuthState);

  // Synchroniser l'état stable avec authState quand il change
  React.useEffect(() => {
    if (authState && typeof authState === 'object') {
      setStableState(prev => ({
        ...authState,
        // Préserver l'état précédent si nouvelles données invalides
        user: authState.user || prev.user,
        isAuthenticated: authState.user ? true : false,
        lastUpdate: Date.now()
      }));
    }
  }, [authState]);

  // Mutation sécurisée pour la validation de compte
  const validateAccountMutation = useMutation({
    mutationFn: async ({ riotId, riotTag, region }: { riotId: string; riotTag: string; region: string }) => {
      const response = await apiService.validateAccount({ riot_id: riotId, riot_tag: riotTag, region });
      return response;
    },
    onSuccess: (response) => {
      if (response && response.valid && response.user && typeof response.user === 'object') {
        // Validation stricte des données utilisateur reçues
        if (response.user.riot_id && response.user.riot_tag && response.user.region) {
          // Réinitialiser la surveillance de sécurité en cas de succès
          securityMonitoring.failedAttempts = 0;
          securityMonitoring.lastFailedAttempt = null;
          setAuthError(null);
          
          // Mettre à jour l'état stable
          setStableState({
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
            lastUpdate: Date.now()
          });
          
          // Mettre à jour le cache React Query avec un état complet
          queryClient.setQueryData(['auth-session'], {
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
            lastUpdate: Date.now()
          });
        } else {
          // Données utilisateur incomplètes
          const errorMessage = 'Données utilisateur incomplètes reçues';
          setAuthError(errorMessage);
          setStableState(prev => ({ ...prev, error: errorMessage, isLoading: false }));
        }
      } else {
        const errorMessage = response.error_message || 'Validation échouée';
        setAuthError(errorMessage);
        
        // Surveiller les tentatives d'échec
        securityMonitoring.failedAttempts++;
        securityMonitoring.lastFailedAttempt = Date.now();
        
        if (detectSuspiciousActivity(errorMessage)) {
          securityMonitoring.suspiciousPatterns.push(errorMessage);
          console.warn('🚨 Activité suspecte détectée lors de la validation');
        }
      }
    },
    onError: (error: Error) => {
      const errorMessage = error.message || 'Erreur de validation';
      setAuthError(errorMessage);
      
      // Surveiller les échecs
      securityMonitoring.failedAttempts++;
      securityMonitoring.lastFailedAttempt = Date.now();
      
      if (detectSuspiciousActivity(errorMessage)) {
        securityMonitoring.suspiciousPatterns.push(errorMessage);
        console.warn('🚨 Tentative d\'authentification suspecte détectée');
      }
    },
  });

  // Mutation sécurisée pour la déconnexion
  const logoutMutation = useMutation({
    mutationFn: async () => {
      await apiService.logout();
    },
    onSettled: () => {
      // Nettoyer le cache et l'état en toutes circonstances
      queryClient.clear();
      setAuthError(null);
      securityMonitoring.failedAttempts = 0;
      securityMonitoring.lastFailedAttempt = null;
      securityMonitoring.suspiciousPatterns = [];
    },
  });

  // Fonctions d'interface
  const validateAccount = (riotId: string, riotTag: string, region: string) => {
    setAuthError(null); // Effacer l'erreur précédente
    validateAccountMutation.mutate({ riotId, riotTag, region });
  };

  const logout = () => {
    logoutMutation.mutate();
  };

  const clearError = () => {
    setAuthError(null);
  };

  // Calcul de l'état de sécurité
  const securityStatus = {
    failedAttempts: securityMonitoring.failedAttempts,
    suspiciousActivity: securityMonitoring.suspiciousPatterns.length > 0 || securityMonitoring.failedAttempts >= 3,
    lastActivity: securityMonitoring.lastFailedAttempt,
  };

  // Utiliser l'état stable pour éviter les données corrompues
  const value: AuthContextType = {
    user: stableState.user,
    isAuthenticated: stableState.isAuthenticated,
    isLoading: stableState.isLoading || validateAccountMutation.isPending || logoutMutation.isPending,
    error: authError || stableState.error,
    validateAccount,
    logout,
    clearError,
    securityStatus,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

export function AuthProvider({ children }: AuthProviderProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProviderContent>
        {children}
      </AuthProviderContent>
    </QueryClientProvider>
  );
}

// Hook sécurisé pour utiliser le contexte d'authentification
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  
  // Surveillance automatique de sécurité
  React.useEffect(() => {
    if (context.securityStatus.suspiciousActivity) {
      console.warn('🚨 Activité suspecte détectée - Surveillance renforcée activée');
    }
  }, [context.securityStatus.suspiciousActivity]);
  
  return context;
}

export default AuthContext;
