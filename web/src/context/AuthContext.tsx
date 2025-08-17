import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react';
import { apiService, User, AuthResponse } from '../services/api';

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
  isLoading: true, // Démarrer avec loading=true pour vérifier l'état au démarrage
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
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Provider du contexte d'authentification
interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [state, dispatch] = useReducer(authReducer, initialState);

  // Vérifier l'état d'authentification au démarrage
  const checkAuthStatus = async () => {
    try {
      dispatch({ type: 'AUTH_START' });
      const response = await apiService.checkSession();
      if (response.authenticated && response.user) {
        dispatch({ type: 'AUTH_SUCCESS', payload: response.user });
      } else {
        dispatch({ type: 'AUTH_ERROR', payload: '' });
      }
    } catch (error) {
      dispatch({ type: 'AUTH_ERROR', payload: '' });
    }
  };

  // Validation de compte Riot
  const validateAccount = async (riotId: string, riotTag: string, region: string) => {
    try {
      dispatch({ type: 'AUTH_START' });
      const response = await apiService.validateAccount({ riot_id: riotId, riot_tag: riotTag, region });
      if (response.valid && response.user) {
        dispatch({ type: 'AUTH_SUCCESS', payload: response.user });
      } else {
        dispatch({ type: 'AUTH_ERROR', payload: response.error_message || 'Validation échouée' });
      }
    } catch (error) {
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

  // Vérifier l'authentification au démarrage
  useEffect(() => {
    checkAuthStatus();
  }, []);

  const value: AuthContextType = {
    state,
    validateAccount,
    logout,
    clearError,
    checkAuthStatus,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

// Hook pour utiliser le contexte d'authentification
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

export default AuthContext;
