// Services pour communiquer avec l'API backend
export interface User {
  id: number;
  riot_id: string;
  riot_tag: string;
  riot_puuid: string;
  summoner_id?: string;
  account_id?: string;
  profile_icon_id: number;
  summoner_level: number;
  region: string;
  is_validated: boolean;
  created_at: string;
  updated_at: string;
  last_sync?: string;
}

export interface ValidationRequest {
  riot_id: string;
  riot_tag: string;
  region: string;
}

export interface ValidationResponse {
  valid: boolean;
  user?: User;
  error_message?: string;
}

export interface SessionResponse {
  authenticated: boolean;
  user?: User;
}

export interface Region {
  code: string;
  name: string;
}

export interface RegionsResponse {
  regions: Region[];
}

export interface DashboardStats {
  total_matches: number;
  win_rate: number;
  average_kda: number;
  favorite_champion: string;
  last_sync_at: string | null;
  next_sync_at: string | null;
}

export interface ApiError {
  error: string;
}

class ApiService {
  private baseUrl = (import.meta.env?.VITE_API_BASE || 'http://localhost:8004') + '/api';

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    // Ajouter les headers personnalisés
    if (options.headers) {
      Object.assign(headers, options.headers);
    }

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: 'include', // Important pour les cookies de session
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Network error' }));
      throw new Error(errorData.error || errorData.message || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // Validation de compte Riot
  async validateAccount(data: ValidationRequest): Promise<ValidationResponse> {
    return this.request<ValidationResponse>('/auth/validate', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Vérifier la session
  async checkSession(): Promise<SessionResponse> {
    return this.request<SessionResponse>('/auth/session');
  }

  // Récupérer les régions supportées
  async getSupportedRegions(): Promise<RegionsResponse> {
    return this.request<RegionsResponse>('/auth/regions');
  }

  // Déconnexion
  async logout(): Promise<void> {
    await this.request('/auth/logout', { method: 'POST' });
  }

  // Profil utilisateur
  async getProfile(): Promise<User> {
    return this.request<User>('/profile');
  }

  // Dashboard
  async getDashboardStats(): Promise<DashboardStats> {
    return this.request<DashboardStats>('/dashboard/stats');
  }

  // Get matches
  async getMatches(page: number = 1, limit: number = 10): Promise<any> {
    return this.request<any>(`/dashboard/matches?page=${page}&limit=${limit}`);
  }

  // Get settings
  async getSettings(): Promise<any> {
    return this.request<any>('/dashboard/settings');
  }

  // Update settings
  async updateSettings(settings: any): Promise<any> {
    return this.request<any>('/dashboard/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }

  // Synchronisation
  async syncMatches(): Promise<{ message: string }> {
    return this.request<{ message: string }>('/dashboard/sync', {
      method: 'POST',
    });
  }

  // Santé de l'API
  async healthCheck(): Promise<{ service: string; status: string }> {
    return this.request<{ service: string; status: string }>('/health');
  }
}

export const apiService = new ApiService();
export default apiService;
