// Services pour communiquer avec l'API backend - Version Sécurisée avec JWT
import { getApiUrl } from '../utils/api-config';
import { jwtVerify, SignJWT } from 'jose';
import CryptoJS from 'crypto-js';
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
  default: string;
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

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  csrfToken: string;
}

export interface SecureValidationResponse {
  valid: boolean;
  user?: User;
  tokens?: AuthTokens;
  error_message?: string;
}

class ApiService {
  private baseUrl = getApiUrl();
  private csrfToken: string | null = null;
  private readonly JWT_SECRET = new TextEncoder().encode('herald-jwt-secret-production');
  private readonly CSRF_SECRET = 'herald-csrf-secret-production';

  constructor() {
    this.initializeSecurity();
  }

  // Initialisation sécurisée
  private initializeSecurity(): void {
    this.csrfToken = CryptoJS.lib.WordArray.random(32).toString();
    this.setupSecurityMonitoring();
  }

  // Surveillance de sécurité continue
  private setupSecurityMonitoring(): void {
    // Surveillance de l'activité suspecte
    setInterval(() => {
      this.checkSuspiciousActivity();
    }, 60000); // Vérifier toutes les minutes
  }

  // Détection d'activité suspecte
  private checkSuspiciousActivity(): boolean {
    const lastActivity = sessionStorage.getItem('last_activity');
    const currentTime = Date.now();
    
    if (lastActivity) {
      const timeDiff = currentTime - parseInt(lastActivity);
      // Si plus de 30 minutes d'inactivité, considérer comme suspect
      if (timeDiff > 30 * 60 * 1000) {
        console.warn('🚨 Activité suspecte détectée - Déconnexion sécurisée');
        this.secureLogout();
        return true;
      }
    }
    
    sessionStorage.setItem('last_activity', currentTime.toString());
    return false;
  }

  // Déconnexion sécurisée forcée
  private secureLogout(): void {
    this.clearSecureTokens();
    window.location.reload();
  }

  // Cryptage sécurisé des données sensibles
  private encryptSensitiveData(data: string): string {
    return CryptoJS.AES.encrypt(data, this.CSRF_SECRET).toString();
  }

  // Décryptage des données sensibles
  private decryptSensitiveData(encryptedData: string): string {
    const bytes = CryptoJS.AES.decrypt(encryptedData, this.CSRF_SECRET);
    return bytes.toString(CryptoJS.enc.Utf8);
  }

  // Stockage sécurisé des tokens
  private storeSecureTokens(tokens: AuthTokens): void {
    sessionStorage.setItem('auth_access_token', tokens.accessToken);
    sessionStorage.setItem('auth_refresh_token', this.encryptSensitiveData(tokens.refreshToken));
    sessionStorage.setItem('auth_csrf_token', tokens.csrfToken);
    sessionStorage.setItem('last_activity', Date.now().toString());
  }

  // Récupération sécurisée des tokens
  private getStoredToken(): string | null {
    return sessionStorage.getItem('auth_access_token');
  }

  private getStoredRefreshToken(): string | null {
    const encryptedToken = sessionStorage.getItem('auth_refresh_token');
    if (!encryptedToken) return null;
    
    try {
      return this.decryptSensitiveData(encryptedToken);
    } catch {
      return null;
    }
  }

  // Nettoyage sécurisé des tokens
  private clearSecureTokens(): void {
    sessionStorage.removeItem('auth_access_token');
    sessionStorage.removeItem('auth_refresh_token');
    sessionStorage.removeItem('auth_csrf_token');
    sessionStorage.removeItem('last_activity');
  }

  // Headers sécurisés pour les requêtes
  private getSecureHeaders(): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Requested-With': 'XMLHttpRequest',
    };

    if (this.csrfToken) {
      headers['X-CSRF-Token'] = this.csrfToken;
    }

    const accessToken = this.getStoredToken();
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    return headers;
  }

  // Validation défensive des données reçues
  private validateResponse(data: any): any {
    if (!data || typeof data !== 'object') {
      throw new Error('Réponse API invalide ou corrompue');
    }
    
    // Validation spécifique pour les réponses utilisateur
    if (data.user && typeof data.user === 'object') {
      if (!data.user.riot_id || !data.user.riot_tag || !data.user.region) {
        console.warn('Données utilisateur incomplètes détectées:', data.user);
        // Nettoyer l'objet utilisateur pour éviter les erreurs downstream
        return { ...data, user: null, authenticated: false }; // Invalider la session
      }
    }
    
    // Validation des tableaux pour éviter les erreurs "Cannot read properties of undefined (reading '0')"
    if (data.regions && Array.isArray(data.regions)) {
      data.regions = data.regions.filter(region => region && typeof region === 'object' && region.code);
    }
    
    // Validation des données de match
    if (data.matches && Array.isArray(data.matches)) {
      data.matches = data.matches.filter(match => match && typeof match === 'object');
    }
    
    return data;
  }

  // Requête sécurisée avec retry automatique et gestion JWT
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const secureHeaders = {
      ...this.getSecureHeaders(),
      ...options.headers,
    };

    const requestOptions: RequestInit = {
      ...options,
      headers: secureHeaders,
      credentials: 'include', // Important pour les cookies de session
    };

    try {
      const response = await fetch(url, requestOptions);

      // Gestion des erreurs d'authentification avec retry automatique
      if (response.status === 401) {
        const refreshToken = this.getStoredRefreshToken();
        if (refreshToken) {
          try {
            await this.refreshAuthToken();
            // Retry de la requête avec le nouveau token
            const retryHeaders = {
              ...this.getSecureHeaders(),
              ...options.headers,
            };
            const retryResponse = await fetch(url, {
              ...requestOptions,
              headers: retryHeaders,
            });
            
            if (!retryResponse.ok) {
              throw new Error(`HTTP ${retryResponse.status}`);
            }
            
            return retryResponse.json();
          } catch {
            this.clearSecureTokens();
            throw new Error('Session expirée, veuillez vous reconnecter');
          }
        } else {
          this.clearSecureTokens();
          throw new Error('Non autorisé');
        }
      }

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Network error' }));
        throw new Error(errorData.error || errorData.message || `HTTP ${response.status}`);
      }

      // Mettre à jour l'activité utilisateur
      sessionStorage.setItem('last_activity', Date.now().toString());

      const responseData = await response.json();
      return this.validateResponse(responseData);
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Erreur de connexion');
    }
  }

  // Refresh sécurisé du token JWT
  private async refreshAuthToken(): Promise<void> {
    const refreshToken = this.getStoredRefreshToken();
    if (!refreshToken) {
      throw new Error('Aucun refresh token disponible');
    }

    const response = await fetch(`${this.baseUrl}/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': this.csrfToken || '',
      },
      credentials: 'include',
      body: JSON.stringify({
        refresh_token: refreshToken,
        csrf_token: this.csrfToken,
      }),
    });

    if (!response.ok) {
      throw new Error('Impossible de renouveler le token');
    }

    const data = await response.json();
    if (data.tokens) {
      this.storeSecureTokens(data.tokens);
    }
  }

  // Validation sécurisée de compte Riot avec tokens JWT
  async validateAccount(data: ValidationRequest): Promise<SecureValidationResponse> {
    const secureData = {
      ...data,
      csrf_token: this.csrfToken,
      timestamp: Date.now(),
    };

    const response = await this.request<SecureValidationResponse>('/auth/validate', {
      method: 'POST',
      body: JSON.stringify(secureData),
    });

    // Stocker les tokens de sécurité si la validation réussit
    if (response.valid && response.tokens) {
      this.storeSecureTokens(response.tokens);
    }

    return response;
  }

  // Vérifier la session
  async checkSession(): Promise<SessionResponse> {
    return this.request<SessionResponse>('/auth/session');
  }

  // Récupérer les régions supportées
  async getSupportedRegions(): Promise<RegionsResponse> {
    return this.request<RegionsResponse>('/auth/regions');
  }

  // Déconnexion sécurisée avec nettoyage complet
  async logout(): Promise<void> {
    try {
      await this.request('/auth/logout', { 
        method: 'POST',
        body: JSON.stringify({
          csrf_token: this.csrfToken,
        }),
      });
    } catch (error) {
      console.error('Erreur lors de la déconnexion:', error);
    } finally {
      this.clearSecureTokens();
      this.csrfToken = CryptoJS.lib.WordArray.random(32).toString();
    }
  }

  // Profil utilisateur
  async getProfile(): Promise<User> {
    return this.request<User>('/profile');
  }

  // Dashboard
  async getDashboardStats(): Promise<DashboardStats> {
    return this.request<DashboardStats>('/dashboard');
  }

  // Get matches
  async getMatches(limit: number = 20, offset: number = 0): Promise<any> {
    return this.request<any>(`/matches?limit=${limit}&offset=${offset}`);
  }

  // Get settings
  async getSettings(): Promise<any> {
    return this.request<any>('/settings');
  }

  // Update settings
  async updateSettings(settings: any): Promise<any> {
    return this.request<any>('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }

  // Synchronisation
  async syncMatches(count: number = 20): Promise<any> {
    return this.request<any>('/sync/matches', {
      method: 'POST',
      body: JSON.stringify({ count }),
    });
  }

  // Get sync status
  async getSyncStatus(jobId: string): Promise<any> {
    return this.request<any>(`/sync/status/${jobId}`);
  }

  // Santé de l'API
  async healthCheck(): Promise<{ service: string; status: string }> {
    return this.request<{ service: string; status: string }>('/health');
  }
}

export const apiService = new ApiService();
export default apiService;
