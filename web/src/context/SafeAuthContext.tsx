import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react';
import { apiService, User } from '../services/api';

// Types pour l'état d'authentification
export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

// Actions pour le reducer
type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: User }
  | { type: 'AUTH_ERROR'; payload: string }
  | { type: 'AUTH_LOGOUT' }
  | { type: 'CLEAR_ERROR' };

// État initial
const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: true,
  error: null,
};

// Reducer pour gérer l'état d'authentification
function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case 'AUTH_START':
      return {
        ...state,
        isLoading: true,
        error: null,
      };
    case 'AUTH_SUCCESS':
      return {
        ...state,
        user: action.payload,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      };
    case 'AUTH_ERROR':
      return {
        ...state,
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: action.payload,
      };
    case 'AUTH_LOGOUT':
      return {
        ...state,
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: null,
      };
    case 'CLEAR_ERROR':
      return {
        ...state,
        error: null,
      };
    default:
      return state;
  }
}

// Interface du contexte
interface AuthContextType {
  state: AuthState;
  validateAccount: (riotId: string, riotTag: string, region: string) => Promise<void>;
  logout: () => Promise<void>;
  clearError: () => void;
  checkAuthStatus: () => Promise<void>;
}

// Création du contexte
const SafeAuthContext = createContext<AuthContextType | undefined>(undefined);

// Provider du contexte d'authentification
interface AuthProviderProps {
  children: ReactNode;
}

export function SafeAuthProvider({ children }: AuthProviderProps) {
  const [state, dispatch] = useReducer(authReducer, initialState);

  // Vérifier l'état d'authentification au démarrage avec plus de sécurité
  const checkAuthStatus = async () => {
    try {
      console.log('🔍 Checking auth status...');
      dispatch({ type: 'AUTH_START' });
      
      const response = await apiService.checkSession();
      console.log('✅ Session response:', response);
      
      // Vérification très défensive
      if (response && typeof response === 'object') {
        if (response.authenticated === true && response.user && typeof response.user === 'object') {
          console.log('✅ User authenticated:', response.user);
          dispatch({ type: 'AUTH_SUCCESS', payload: response.user });
        } else {
          console.log('ℹ️ User not authenticated');
          dispatch({ type: 'AUTH_ERROR', payload: '' });
        }
      } else {
        console.error('❌ Invalid response format:', response);
        dispatch({ type: 'AUTH_ERROR', payload: 'Invalid response format' });
      }
    } catch (error) {
      console.error('❌ Auth check error:', error);
      dispatch({ type: 'AUTH_ERROR', payload: error instanceof Error ? error.message : 'Unknown error' });
    }
  };

  // Validation de compte Riot avec plus de sécurité
  const validateAccount = async (riotId: string, riotTag: string, region: string) => {
    try {
      console.log('🔍 Validating account:', { riotId, riotTag, region });
      dispatch({ type: 'AUTH_START' });
      
      const response = await apiService.validateAccount({ riot_id: riotId, riot_tag: riotTag, region });
      console.log('✅ Validation response:', response);
      
      if (response && typeof response === 'object') {
        if (response.valid === true && response.user && typeof response.user === 'object') {
          console.log('✅ Account validated:', response.user);
          dispatch({ type: 'AUTH_SUCCESS', payload: response.user });
        } else {
          const errorMsg = response.error_message || 'Validation échouée';
          console.log('❌ Validation failed:', errorMsg);
          dispatch({ type: 'AUTH_ERROR', payload: errorMsg });
        }
      } else {
        console.error('❌ Invalid validation response:', response);
        dispatch({ type: 'AUTH_ERROR', payload: 'Invalid response format' });
      }
    } catch (error) {
      console.error('❌ Validation error:', error);
      const message = error instanceof Error ? error.message : 'Erreur de validation';
      dispatch({ type: 'AUTH_ERROR', payload: message });
    }
  };

  // Déconnexion
  const logout = async () => {
    try {
      await apiService.logout();
    } catch (error) {
      console.error('Erreur lors de la déconnexion:', error);
    } finally {
      dispatch({ type: 'AUTH_LOGOUT' });
    }
  };

  // Effacer l'erreur
  const clearError = () => {
    dispatch({ type: 'CLEAR_ERROR' });
  };

  // Vérifier l'authentification au démarrage (on initial mount only)
  useEffect(() => {
    console.log('🚀 AuthContext mounting, checking auth status...');
    checkAuthStatus();
  }, []); // Empty dependency array

  const value: AuthContextType = {
    state,
    validateAccount,
    logout,
    clearError,
    checkAuthStatus,
  };

  return (
    <SafeAuthContext.Provider value={value}>
      {children}
    </SafeAuthContext.Provider>
  );
}

// Hook pour utiliser le contexte d'authentification
export function useSafeAuth() {
  const context = useContext(SafeAuthContext);
  if (context === undefined) {
    throw new Error('useSafeAuth must be used within a SafeAuthProvider');
  }
  return context;
}

export default SafeAuthContext;