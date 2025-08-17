export interface Row {
  [k: string]: any;
  date?: string;
  win?: boolean;
  kills?: number; deaths?: number; assists?: number; kda?: number; kp?: number;
  champion?: string; queue?: string; duration_s?: number; cs?: number; gold?: number;
  dmg_to_champs?: number; vision_score?: number; side?: string; lane?: string; role?: string;
  role_std?: string;
  cs_per_min?: number; gpm?: number; dpm?: number;
  gold_share?: number; dmg_share?: number; vision_share?: number;
  cs10?: number; gold10?: number; xp10?: number;
  matchId?: string;
  kills_per_min?: number; wards_placed?: number; wards_killed?: number;
}