// Configuration dynamique de l'API
export function getApiBaseUrl(): string {
  // Utilise la configuration runtime depuis config.js
  const appConfig = (window as any).APP_CONFIG;
  if (appConfig?.API_BASE) {
    return appConfig.API_BASE;
  }

  // Fallback sur les variables d'environnement Vite
  const viteApiBase = import.meta.env?.VITE_API_BASE || (window as any).VITE_API_BASE;
  if (viteApiBase) {
    return viteApiBase;
  }

  // Fallback intelligent basé sur l'URL actuelle
  return `${window.location.protocol}//${window.location.host}`;
}

export function getApiUrl(endpoint: string = ''): string {
  const baseUrl = getApiBaseUrl();
  const apiPath = (window as any).APP_CONFIG?.API_PATH || '/api';
  
  // Si l'endpoint commence déjà par /api, ne pas le dupliquer
  if (endpoint.startsWith('/api')) {
    return `${baseUrl}${endpoint}`;
  }
  
  // Nettoyer l'endpoint pour éviter les doubles slashes
  const cleanEndpoint = endpoint.replace(/^\/+/, ''); // enlever les slashes au début
  
  // Construire l'URL finale
  if (cleanEndpoint) {
    return `${baseUrl}${apiPath}/${cleanEndpoint}`;
  } else {
    return `${baseUrl}${apiPath}`;
  }
}