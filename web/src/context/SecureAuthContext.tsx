import React, { createContext, useContext, useEffect, ReactNode } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { secureAuthService, SecureUser, AuthTokens } from '../auth/SecureAuthService';

// Interface pour l'état d'authentification sécurisé
interface SecureAuthContextType {
  // État
  user: SecureUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  validateAccount: (riotId: string, riotTag: string, region: string) => Promise<void>;
  logout: () => Promise<void>;
  clearError: () => void;
  
  // Sécurité
  checkSuspiciousActivity: () => boolean;
}

// Création du contexte sécurisé
const SecureAuthContext = createContext<SecureAuthContextType | undefined>(undefined);

// Provider du contexte d'authentification sécurisé
interface SecureAuthProviderProps {
  children: ReactNode;
}

export function SecureAuthProvider({ children }: SecureAuthProviderProps) {
  const queryClient = useQueryClient();

  // Query pour vérifier la session
  const {
    data: sessionData,
    isLoading,
    error: sessionError,
  } = useQuery({
    queryKey: ['auth', 'session'],
    queryFn: async () => {
      console.log('🔒 Vérification sécurisée de la session...');
      
      // Vérifier l'activité suspecte
      if (secureAuthService.detectSuspiciousActivity()) {
        console.warn('⚠️ Activité suspecte détectée, déconnexion de sécurité');
        await secureAuthService.logout();
        throw new Error('Session expirée pour des raisons de sécurité');
      }
      
      const result = await secureAuthService.checkSession();
      console.log('✅ Vérification de session terminée:', result);
      return result;
    },
    retry: 1,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  });

  // Mutation pour la validation de compte
  const validateAccountMutation = useMutation({
    mutationFn: async ({
      riotId,
      riotTag,
      region,
    }: {
      riotId: string;
      riotTag: string;
      region: string;
    }) => {
      console.log('🔒 Validation sécurisée du compte Riot...', { riotId, riotTag, region });
      return secureAuthService.validateRiotAccount(riotId, riotTag, region);
    },
    onSuccess: (data) => {
      console.log('✅ Validation réussie:', data.user);
      // Mettre à jour le cache de session
      queryClient.setQueryData(['auth', 'session'], {
        user: data.user,
        isAuthenticated: true,
      });
    },
    onError: (error) => {
      console.error('❌ Erreur de validation:', error);
    },
  });

  // Mutation pour la déconnexion
  const logoutMutation = useMutation({
    mutationFn: async () => {
      console.log('🔒 Déconnexion sécurisée...');
      await secureAuthService.logout();
    },
    onSuccess: () => {
      console.log('✅ Déconnexion réussie');
      // Nettoyer tout le cache
      queryClient.clear();
      queryClient.setQueryData(['auth', 'session'], {
        user: null,
        isAuthenticated: false,
      });
    },
    onError: (error) => {
      console.error('❌ Erreur de déconnexion:', error);
      // Même en cas d'erreur, nettoyer le cache local
      queryClient.clear();
    },
  });

  // Actions du contexte
  const validateAccount = async (riotId: string, riotTag: string, region: string) => {
    await validateAccountMutation.mutateAsync({ riotId, riotTag, region });
  };

  const logout = async () => {
    await logoutMutation.mutateAsync();
  };

  const clearError = () => {
    queryClient.setQueryData(['auth', 'session'], (old: any) => ({
      ...old,
      error: null,
    }));
  };

  const checkSuspiciousActivity = () => {
    return secureAuthService.detectSuspiciousActivity();
  };

  // Surveillance continue de la sécurité
  useEffect(() => {
    const securityInterval = setInterval(() => {
      if (sessionData?.isAuthenticated && checkSuspiciousActivity()) {
        console.warn('⚠️ Activité suspecte détectée, déconnexion automatique');
        logout();
      }
    }, 60000); // Vérifier toutes les minutes

    return () => clearInterval(securityInterval);
  }, [sessionData?.isAuthenticated]);

  // Valeurs du contexte
  const contextValue: SecureAuthContextType = {
    user: sessionData?.user || null,
    isAuthenticated: sessionData?.isAuthenticated || false,
    isLoading: isLoading || validateAccountMutation.isPending || logoutMutation.isPending,
    error: sessionError?.message || validateAccountMutation.error?.message || null,
    validateAccount,
    logout,
    clearError,
    checkSuspiciousActivity,
  };

  return (
    <SecureAuthContext.Provider value={contextValue}>
      {children}
    </SecureAuthContext.Provider>
  );
}

// Hook pour utiliser le contexte d'authentification sécurisé
export function useSecureAuth() {
  const context = useContext(SecureAuthContext);
  if (context === undefined) {
    throw new Error('useSecureAuth must be used within a SecureAuthProvider');
  }
  return context;
}

export default SecureAuthContext;