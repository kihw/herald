/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_APP_ENV: string;
  readonly VITE_APP_VERSION: string;
  readonly VITE_RIOT_API_KEY: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}