// Auth Types
export interface User {
  id: string;
  email: string;
  username: string;
  display_name?: string;
  avatar?: string;
  bio?: string;
  timezone: string;
  language: string;
  is_active: boolean;
  is_premium: boolean;
  last_login: string;
  login_count: number;
  riot_accounts?: RiotAccount[];
  preferences: UserPreferences;
  subscription?: Subscription;
  total_matches: number;
  last_sync_at: string;
  favorite_champion?: string;
  main_role?: string;
  current_rank?: string;
  created_at: string;
  updated_at: string;
}

export interface RiotAccount {
  id: string;
  user_id: string;
  puuid: string;
  summoner_name: string;
  tag_line: string;
  summoner_id: string;
  account_id: string;
  region: string;
  platform: string;
  is_verified: boolean;
  is_primary: boolean;
  last_sync_at: string;
  solo_queue_rank?: string;
  flex_queue_rank?: string;
  tft_rank?: string;
  arena_rank?: string;
  summoner_level: number;
  profile_icon: number;
  total_mastery_score: number;
  created_at: string;
  updated_at: string;
}

export interface UserPreferences {
  id: string;
  user_id: string;
  theme: 'dark' | 'light' | 'auto';
  compact_mode: boolean;
  show_detailed_stats: boolean;
  default_timeframe: '1d' | '7d' | '30d' | 'season';
  email_notifications: boolean;
  push_notifications: boolean;
  match_notifications: boolean;
  rank_change_notifications: boolean;
  auto_sync_matches: boolean;
  sync_interval: number;
  include_normal_games: boolean;
  include_aram_games: boolean;
  favorite_game_modes?: string;
  public_profile: boolean;
  show_in_leaderboards: boolean;
  allow_data_export: boolean;
  receive_ai_coaching: boolean;
  coaching_focus?: string;
  skill_level: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  preferred_coaching_style: 'gentle' | 'direct' | 'balanced';
  created_at: string;
  updated_at: string;
}

export interface Subscription {
  id: string;
  user_id: string;
  plan: 'free' | 'premium' | 'elite' | 'enterprise';
  status: 'active' | 'canceled' | 'expired' | 'trial';
  started_at: string;
  expires_at: string;
  trial_ends_at?: string;
  amount: number;
  currency: string;
  interval: 'monthly' | 'yearly';
  payment_method?: string;
  max_riot_accounts: number;
  unlimited_analytics: boolean;
  ai_coaching_access: boolean;
  advanced_metrics: boolean;
  data_export_access: boolean;
  priority_support: boolean;
  created_at: string;
  updated_at: string;
}

// Auth API Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  display_name?: string;
}

export interface AuthResponse {
  token: string;
  refresh_token: string;
  user: User;
  expires_in: number;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

// Match Types
export interface Match {
  id: string;
  match_id: string;
  game_id: number;
  platform_id: string;
  game_mode: string;
  game_type: string;
  queue_id: number;
  map_id: number;
  game_start_timestamp: number;
  game_end_timestamp: number;
  game_duration: number;
  game_version: string;
  participants: MatchParticipant[];
  is_processed: boolean;
  processed_at?: string;
  is_analyzed: boolean;
  analyzed_at?: string;
  winning_team: number;
  created_at: string;
  updated_at: string;
}

export interface MatchParticipant {
  id: string;
  match_id: string;
  puuid: string;
  summoner_name: string;
  summoner_id: string;
  participant_id: number;
  team_id: number;
  team_position: string;
  champion_id: number;
  champion_name: string;
  spell1_id: number;
  spell2_id: number;
  primary_rune_style: number;
  sub_rune_style: number;
  kills: number;
  deaths: number;
  assists: number;
  champion_level: number;
  won: boolean;
  total_damage_dealt: number;
  total_damage_dealt_to_champions: number;
  total_damage_taken: number;
  total_heal: number;
  total_heals_on_teammates: number;
  damage_dealt_to_objectives: number;
  damage_dealt_to_turrets: number;
  gold_earned: number;
  gold_spent: number;
  total_cs: number;
  cs_per_minute: number;
  vision_score: number;
  wards_placed: number;
  wards_killed: number;
  control_wards_placed: number;
  vision_wards_bought_in_game: number;
  item0: number;
  item1: number;
  item2: number;
  item3: number;
  item4: number;
  item5: number;
  item6: number;
  kda: number;
  kill_participation: number;
  damage_share: number;
  gold_share: number;
  performance_score: number;
  turret_kills: number;
  inhibitor_kills: number;
  dragon_kills: number;
  baron_kills: number;
  first_blood_kill: boolean;
  first_blood_assist: boolean;
  largest_killing_spree: number;
  largest_multi_kill: number;
  early_game_performance: number;
  mid_game_performance: number;
  late_game_performance: number;
  jungle_cs: number;
  lane_cs: number;
  support_item_quest: boolean;
  roaming_score: number;
  teamfight_score: number;
  economic_efficiency: number;
  vision_efficiency: number;
  objective_contribution: number;
  created_at: string;
  updated_at: string;
}

// Analytics Types
export interface PlayerStats {
  total_matches: number;
  wins: number;
  losses: number;
  win_rate: number;
  avg_kda: number;
  avg_kills: number;
  avg_deaths: number;
  avg_assists: number;
  avg_cs: number;
  avg_vision_score: number;
  avg_gold_earned: number;
  avg_damage_dealt: number;
  avg_game_duration: number;
  favorite_champions: ChampionStats[];
  recent_performance: PerformanceTrend[];
  rank_progression: RankProgression[];
}

export interface ChampionStats {
  champion_id: number;
  champion_name: string;
  games_played: number;
  wins: number;
  losses: number;
  win_rate: number;
  avg_kda: number;
  avg_cs: number;
  avg_damage: number;
  performance_score: number;
}

export interface PerformanceTrend {
  date: string;
  matches: number;
  wins: number;
  avg_kda: number;
  avg_performance: number;
}

export interface RankProgression {
  date: string;
  tier: string;
  division: string;
  lp: number;
  wins: number;
  losses: number;
}

// Riot API Types
export interface SummonerInfo {
  id: string;
  account_id: string;
  puuid: string;
  name: string;
  profile_icon_id: number;
  revision_date: number;
  summoner_level: number;
}

export interface LeagueEntry {
  league_id: string;
  summoner_id: string;
  summoner_name: string;
  queue_type: string;
  tier: string;
  rank: string;
  league_points: number;
  wins: number;
  losses: number;
  hot_streak: boolean;
  veteran: boolean;
  fresh_blood: boolean;
  inactive: boolean;
}

export interface MatchHistory {
  match_ids: string[];
}

// UI Types
export interface Theme {
  mode: 'light' | 'dark';
  primary: string;
  secondary: string;
  background: string;
  surface: string;
  text: string;
  error: string;
  warning: string;
  info: string;
  success: string;
}

export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message: string;
  timestamp: string;
  read: boolean;
  action?: {
    label: string;
    url: string;
  };
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface ApiError {
  error: string;
  message: string;
  status_code?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

// Form Types
export interface FormErrors {
  [key: string]: string[];
}

export interface LoadingState {
  [key: string]: boolean;
}

// Navigation Types
export interface NavItem {
  label: string;
  path: string;
  icon?: string;
  children?: NavItem[];
  requiresAuth?: boolean;
  requiresPremium?: boolean;
}

// Gaming Constants
export const REGIONS = [
  { code: 'na1', name: 'North America' },
  { code: 'euw1', name: 'Europe West' },
  { code: 'eun1', name: 'Europe Nordic & East' },
  { code: 'kr', name: 'Korea' },
  { code: 'jp1', name: 'Japan' },
  { code: 'br1', name: 'Brazil' },
  { code: 'la1', name: 'Latin America North' },
  { code: 'la2', name: 'Latin America South' },
  { code: 'oc1', name: 'Oceania' },
  { code: 'tr1', name: 'Turkey' },
  { code: 'ru', name: 'Russia' },
] as const;

export const QUEUE_TYPES = {
  420: 'Ranked Solo/Duo',
  440: 'Ranked Flex',
  450: 'ARAM',
  400: 'Draft Pick',
  430: 'Blind Pick',
  700: 'Clash',
  1700: 'Arena',
} as const;

export const RANKED_TIERS = [
  'IRON',
  'BRONZE', 
  'SILVER',
  'GOLD',
  'PLATINUM',
  'EMERALD',
  'DIAMOND',
  'MASTER',
  'GRANDMASTER',
  'CHALLENGER',
] as const;

export const RANKED_DIVISIONS = ['IV', 'III', 'II', 'I'] as const;