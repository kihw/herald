import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { 
  User, 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse, 
  ChangePasswordRequest,
  ApiResponse,
  ApiError 
} from '@/types';

class ApiService {
  private api: AxiosInstance;
  private baseURL: string;

  constructor() {
    this.baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
    
    this.api = axios.create({
      baseURL: this.baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor to add auth token
    this.api.interceptors.request.use(
      (config) => {
        const token = this.getToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor to handle errors and token refresh
    this.api.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;

        // Handle 401 errors (token expired)
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          const refreshToken = this.getRefreshToken();
          if (refreshToken) {
            try {
              const response = await this.refreshToken(refreshToken);
              this.setTokens(response.token, response.refresh_token);
              
              // Retry original request
              originalRequest.headers.Authorization = `Bearer ${response.token}`;
              return this.api(originalRequest);
            } catch (refreshError) {
              // Refresh failed, clear tokens and redirect to login
              this.clearTokens();
              window.location.href = '/login';
              return Promise.reject(refreshError);
            }
          } else {
            // No refresh token, redirect to login
            this.clearTokens();
            window.location.href = '/login';
          }
        }

        // Handle rate limiting
        if (error.response?.status === 429) {
          const retryAfter = error.response.headers['retry-after'] || 1;
          await new Promise(resolve => setTimeout(resolve, retryAfter * 1000));
          return this.api(originalRequest);
        }

        return Promise.reject(this.formatError(error));
      }
    );
  }

  private formatError(error: any): ApiError {
    if (error.response?.data) {
      return {
        error: error.response.data.error || 'Unknown error',
        message: error.response.data.message || error.message,
        status_code: error.response.status,
      };
    }
    
    return {
      error: 'Network Error',
      message: error.message || 'An unexpected error occurred',
    };
  }

  // Token management
  private getToken(): string | null {
    return localStorage.getItem('herald_token');
  }

  private getRefreshToken(): string | null {
    return localStorage.getItem('herald_refresh_token');
  }

  private setTokens(token: string, refreshToken: string): void {
    localStorage.setItem('herald_token', token);
    localStorage.setItem('herald_refresh_token', refreshToken);
  }

  private clearTokens(): void {
    localStorage.removeItem('herald_token');
    localStorage.removeItem('herald_refresh_token');
  }

  // Generic API methods
  private async get<T>(endpoint: string, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.api.get(endpoint, config);
    return response.data;
  }

  private async post<T>(endpoint: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.api.post(endpoint, data, config);
    return response.data;
  }

  private async put<T>(endpoint: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.api.put(endpoint, data, config);
    return response.data;
  }

  private async delete<T>(endpoint: string, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.api.delete(endpoint, config);
    return response.data;
  }

  private async patch<T>(endpoint: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.api.patch(endpoint, data, config);
    return response.data;
  }

  // Authentication endpoints
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await this.post<AuthResponse>('/auth/login', credentials);
    this.setTokens(response.token, response.refresh_token);
    return response;
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response = await this.post<AuthResponse>('/auth/register', userData);
    this.setTokens(response.token, response.refresh_token);
    return response;
  }

  async logout(): Promise<void> {
    try {
      await this.post<void>('/auth/logout');
    } finally {
      this.clearTokens();
    }
  }

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    return this.post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken });
  }

  async getProfile(): Promise<User> {
    return this.get<User>('/auth/profile');
  }

  async changePassword(data: ChangePasswordRequest): Promise<ApiResponse<void>> {
    return this.post<ApiResponse<void>>('/auth/change-password', data);
  }

  async resetPassword(email: string): Promise<ApiResponse<void>> {
    return this.post<ApiResponse<void>>('/auth/reset-password', { email });
  }

  // User management endpoints
  async updateProfile(data: Partial<User>): Promise<User> {
    return this.put<User>('/users/profile', data);
  }

  async updatePreferences(preferences: any): Promise<void> {
    return this.put<void>('/users/preferences', preferences);
  }

  async deleteAccount(): Promise<void> {
    return this.delete<void>('/users/account');
  }

  // Riot API endpoints
  async linkRiotAccount(data: {
    region: string;
    game_name: string;
    tag_line: string;
  }): Promise<any> {
    return this.post<any>('/riot/link', data);
  }

  async getLinkedAccounts(): Promise<any[]> {
    return this.get<any[]>('/riot/accounts');
  }

  async syncMatches(accountId: string, count = 20): Promise<ApiResponse<void>> {
    return this.post<ApiResponse<void>>(`/riot/accounts/${accountId}/sync`, { count });
  }

  async getSummonerInfo(data: {
    region: string;
    game_name: string;
    tag_line: string;
  }): Promise<any> {
    return this.post<any>('/riot/summoner', data);
  }

  async getRankedInfo(data: {
    region: string;
    summoner_id: string;
  }): Promise<any[]> {
    return this.post<any[]>('/riot/ranked', data);
  }

  async getMatchHistory(params: {
    puuid: string;
    region: string;
    count?: number;
  }): Promise<any> {
    return this.get<any>('/riot/matches', { params });
  }

  async getMatchDetails(matchId: string, region: string): Promise<any> {
    return this.get<any>(`/riot/matches/${matchId}`, { params: { region } });
  }

  async getRateLimitStatus(): Promise<any> {
    return this.get<any>('/riot/rate-limit');
  }

  // Analytics endpoints
  async getPlayerStats(accountId: string, timeframe = '30d'): Promise<any> {
    return this.get<any>(`/analytics/players/${accountId}/stats`, {
      params: { timeframe }
    });
  }

  async getChampionStats(accountId: string, timeframe = '30d'): Promise<any> {
    return this.get<any>(`/analytics/players/${accountId}/champions`, {
      params: { timeframe }
    });
  }

  async getPerformanceTrends(accountId: string, timeframe = '30d'): Promise<any> {
    return this.get<any>(`/analytics/players/${accountId}/trends`, {
      params: { timeframe }
    });
  }

  async getRankProgression(accountId: string): Promise<any> {
    return this.get<any>(`/analytics/players/${accountId}/rank-progression`);
  }

  async getMatchAnalysis(matchId: string): Promise<any> {
    return this.get<any>(`/analytics/matches/${matchId}/analysis`);
  }

  async getCoachingInsights(accountId: string): Promise<any> {
    return this.get<any>(`/analytics/players/${accountId}/coaching`);
  }

  // Match endpoints
  async getMatches(params: {
    account_id?: string;
    queue_id?: number;
    champion_id?: number;
    limit?: number;
    offset?: number;
  } = {}): Promise<any> {
    return this.get<any>('/matches', { params });
  }

  async getMatch(matchId: string): Promise<any> {
    return this.get<any>(`/matches/${matchId}`);
  }

  // Health check
  async healthCheck(): Promise<any> {
    return this.get<any>('/health');
  }

  // Utility methods
  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  getAuthHeader(): string | null {
    const token = this.getToken();
    return token ? `Bearer ${token}` : null;
  }
}

// Export singleton instance
export const apiService = new ApiService();
export default apiService;