import { useState, useEffect, createContext, useContext } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiService } from '@/services/api';
import { User, LoginRequest, RegisterRequest } from '@/types';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    // Fallback implementation when used outside provider
    const queryClient = useQueryClient();
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    // Get user profile query
    const { data: user, isLoading, refetch } = useQuery({
      queryKey: ['auth', 'profile'],
      queryFn: apiService.getProfile,
      enabled: apiService.isAuthenticated(),
      retry: (failureCount, error: any) => {
        if (error?.status === 401) return false;
        return failureCount < 2;
      },
      staleTime: 5 * 60 * 1000, // 5 minutes
    });

    // Login mutation
    const loginMutation = useMutation({
      mutationFn: apiService.login,
      onSuccess: (data) => {
        setIsAuthenticated(true);
        queryClient.setQueryData(['auth', 'profile'], data.user);
        queryClient.invalidateQueries({ queryKey: ['auth'] });
      },
      onError: () => {
        setIsAuthenticated(false);
      },
    });

    // Register mutation
    const registerMutation = useMutation({
      mutationFn: apiService.register,
      onSuccess: (data) => {
        setIsAuthenticated(true);
        queryClient.setQueryData(['auth', 'profile'], data.user);
        queryClient.invalidateQueries({ queryKey: ['auth'] });
      },
      onError: () => {
        setIsAuthenticated(false);
      },
    });

    // Logout mutation
    const logoutMutation = useMutation({
      mutationFn: apiService.logout,
      onSettled: () => {
        setIsAuthenticated(false);
        queryClient.clear();
      },
    });

    // Update authentication state based on user data
    useEffect(() => {
      setIsAuthenticated(!!user);
    }, [user]);

    // Check initial authentication state
    useEffect(() => {
      setIsAuthenticated(apiService.isAuthenticated());
    }, []);

    const login = async (credentials: LoginRequest) => {
      await loginMutation.mutateAsync(credentials);
    };

    const register = async (userData: RegisterRequest) => {
      await registerMutation.mutateAsync(userData);
    };

    const logout = async () => {
      await logoutMutation.mutateAsync();
    };

    const refreshUser = () => {
      refetch();
    };

    return {
      user: user || null,
      isAuthenticated,
      isLoading: isLoading || loginMutation.isPending || registerMutation.isPending,
      login,
      register,
      logout,
      refreshUser,
    };
  }
  
  return context;
};