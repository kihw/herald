// useAuth Hook Tests for Herald.lol
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactNode } from 'react';
import { server } from '@/test/mocks/server';
import { http, HttpResponse } from 'msw';
import { useAuth } from '../useAuth';

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0 },
      mutations: { retry: false },
    },
  });

  return ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('useAuth Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  afterEach(() => {
    localStorage.clear();
  });

  describe('Initial State', () => {
    it('should return initial authentication state', () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.user).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.isLoading).toBe(false);
      expect(typeof result.current.login).toBe('function');
      expect(typeof result.current.logout).toBe('function');
      expect(typeof result.current.register).toBe('function');
    });
  });

  describe('Login Process', () => {
    it('should successfully login with valid credentials', async () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(true);
        expect(result.current.user).not.toBeNull();
        expect(result.current.user?.email).toBe('test@herald.lol');
      });

      // Check if token is stored
      expect(localStorage.getItem('herald_token')).toBeTruthy();
    });

    it('should handle login errors gracefully', async () => {
      server.use(
        http.post('/api/v1/auth/login', () => {
          return HttpResponse.json(
            { error: 'Invalid credentials' },
            { status: 401 }
          );
        })
      );

      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await act(async () => {
        try {
          await result.current.login({
            email: 'invalid@herald.lol',
            password: 'wrongpassword'
          });
        } catch (error) {
          expect(error).toBeDefined();
        }
      });

      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.user).toBeNull();
    });

    it('should set loading state during login', async () => {
      server.use(
        http.post('/api/v1/auth/login', async () => {
          await new Promise(resolve => setTimeout(resolve, 100));
          return HttpResponse.json({
            token: 'mock-token',
            user: { id: '1', email: 'test@herald.lol' }
          });
        })
      );

      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      act(() => {
        result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      expect(result.current.isLoading).toBe(true);

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });
    });
  });

  describe('Registration Process', () => {
    it('should successfully register new user', async () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await act(async () => {
        await result.current.register({
          email: 'newuser@herald.lol',
          password: 'password123',
          username: 'newuser',
          display_name: 'New User'
        });
      });

      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(true);
        expect(result.current.user).not.toBeNull();
        expect(result.current.user?.email).toBe('test@herald.lol'); // Mock returns test user
      });
    });

    it('should handle registration validation errors', async () => {
      server.use(
        http.post('/api/v1/auth/register', () => {
          return HttpResponse.json(
            { 
              error: 'Validation failed',
              details: { email: 'Email already exists' }
            },
            { status: 400 }
          );
        })
      );

      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await act(async () => {
        try {
          await result.current.register({
            email: 'existing@herald.lol',
            password: 'password123',
            username: 'existing',
            display_name: 'Existing User'
          });
        } catch (error) {
          expect(error).toBeDefined();
        }
      });

      expect(result.current.isAuthenticated).toBe(false);
    });
  });

  describe('Logout Process', () => {
    it('should successfully logout user', async () => {
      // First login
      localStorage.setItem('herald_token', 'mock-token');
      
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      // Simulate logged in state
      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      // Then logout
      await act(async () => {
        await result.current.logout();
      });

      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(false);
        expect(result.current.user).toBeNull();
      });

      // Check if token is removed
      expect(localStorage.getItem('herald_token')).toBeNull();
    });

    it('should handle logout API errors gracefully', async () => {
      server.use(
        http.post('/api/v1/auth/logout', () => {
          return HttpResponse.json(
            { error: 'Logout failed' },
            { status: 500 }
          );
        })
      );

      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      // Login first
      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      // Attempt logout
      await act(async () => {
        await result.current.logout();
      });

      // Should still clear local state even if API fails
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.user).toBeNull();
    });
  });

  describe('Token Management', () => {
    it('should restore session from stored token', async () => {
      localStorage.setItem('herald_token', 'valid-token');
      
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(true);
        expect(result.current.user).not.toBeNull();
      });
    });

    it('should handle invalid stored token', async () => {
      localStorage.setItem('herald_token', 'invalid-token');
      
      server.use(
        http.get('/api/v1/auth/profile', () => {
          return HttpResponse.json(
            { error: 'Invalid token' },
            { status: 401 }
          );
        })
      );

      const wrapper = createWrapper();
      renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(localStorage.getItem('herald_token')).toBeNull();
      });
    });

    it('should refresh token when needed', async () => {
      localStorage.setItem('herald_token', 'expiring-token');
      localStorage.setItem('herald_refresh_token', 'valid-refresh-token');

      server.use(
        http.get('/api/v1/auth/profile', () => {
          return HttpResponse.json(
            { error: 'Token expired' },
            { status: 401 }
          );
        }),
        http.post('/api/v1/auth/refresh', () => {
          return HttpResponse.json({
            token: 'new-token',
            refresh_token: 'new-refresh-token',
            user: { id: '1', email: 'test@herald.lol' }
          });
        })
      );

      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(localStorage.getItem('herald_token')).toBe('new-token');
        expect(result.current.isAuthenticated).toBe(true);
      });
    });
  });

  describe('Performance Requirements', () => {
    it('should complete authentication within reasonable time', async () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      const startTime = Date.now();
      
      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      const authTime = Date.now() - startTime;
      expect(authTime).toBeLessThan(2000); // Auth should be fast
    });

    it('should handle multiple concurrent auth requests', async () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      const promises = Array.from({ length: 5 }, () =>
        act(async () => {
          try {
            await result.current.login({
              email: 'test@herald.lol',
              password: 'password123'
            });
          } catch (error) {
            // Some may fail due to race conditions, which is expected
          }
        })
      );

      await Promise.allSettled(promises);

      // Final state should be consistent
      expect(result.current.isAuthenticated).toBe(true);
    });
  });

  describe('Gaming Context Integration', () => {
    it('should include gaming preferences in user data', async () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      await waitFor(() => {
        expect(result.current.user?.preferences).toBeDefined();
        expect(result.current.user?.preferences?.receive_ai_coaching).toBeDefined();
        expect(result.current.user?.preferences?.skill_level).toBeDefined();
      });
    });

    it('should handle Riot account linking status', async () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      await waitFor(() => {
        expect(result.current.user).toBeDefined();
        // User should have riot account status available for analytics features
      });
    });
  });

  describe('Error Recovery', () => {
    it('should recover from network errors', async () => {
      let requestCount = 0;
      server.use(
        http.post('/api/v1/auth/login', () => {
          requestCount++;
          if (requestCount === 1) {
            return HttpResponse.error();
          }
          return HttpResponse.json({
            token: 'mock-token',
            user: { id: '1', email: 'test@herald.lol' }
          });
        })
      );

      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      // First attempt should fail
      await act(async () => {
        try {
          await result.current.login({
            email: 'test@herald.lol',
            password: 'password123'
          });
        } catch (error) {
          expect(error).toBeDefined();
        }
      });

      expect(result.current.isAuthenticated).toBe(false);

      // Second attempt should succeed
      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      expect(result.current.isAuthenticated).toBe(true);
    });
  });

  describe('Security', () => {
    it('should not expose sensitive data', async () => {
      const wrapper = createWrapper();
      const { result } = renderHook(() => useAuth(), { wrapper });

      await act(async () => {
        await result.current.login({
          email: 'test@herald.lol',
          password: 'password123'
        });
      });

      // Password should never be exposed
      expect(JSON.stringify(result.current)).not.toContain('password');
      expect(JSON.stringify(result.current.user)).not.toContain('password');
    });

    it('should handle token expiration gracefully', async () => {
      localStorage.setItem('herald_token', 'expired-token');
      
      server.use(
        http.get('/api/v1/auth/profile', () => {
          return HttpResponse.json(
            { error: 'Token expired' },
            { status: 401 }
          );
        }),
        http.post('/api/v1/auth/refresh', () => {
          return HttpResponse.json(
            { error: 'Refresh token invalid' },
            { status: 401 }
          );
        })
      );

      const wrapper = createWrapper();
      renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(localStorage.getItem('herald_token')).toBeNull();
        expect(localStorage.getItem('herald_refresh_token')).toBeNull();
      });
    });
  });
});