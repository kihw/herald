import { getApiUrl } from '../utils/api-config';

// Types pour les groupes
export interface Group {
  id: number;
  name: string;
  description: string;
  privacy: 'public' | 'private' | 'invite_only';
  member_count: number;
  owner: GroupUser;
  members?: GroupMember[];
  invite_code?: string;
  created_at: string;
  updated_at: string;
  stats?: GroupStats;
}

export interface GroupUser {
  id: number;
  riot_id: string;
  riot_tag: string;
  region: string;
  rank?: string;
  lp?: number;
  mmr?: number;
}

export interface GroupMember {
  id: number;
  user: GroupUser;
  role: 'owner' | 'admin' | 'member';
  status: 'active' | 'pending' | 'banned';
  joined_at: string;
  nickname?: string;
}

export interface GroupStats {
  total_members: number;
  active_members: number;
  average_rank?: string;
  average_mmr?: number;
  top_champions: ChampionStat[];
  popular_roles: RoleStat[];
  winrate_comparison: { [key: string]: number };
  last_updated: string;
}

export interface ChampionStat {
  champion_id: number;
  champion_name: string;
  play_count: number;
  win_rate: number;
  avg_kda: number;
}

export interface RoleStat {
  role: string;
  play_count: number;
  win_rate: number;
}

export interface GroupComparison {
  id: number;
  name: string;
  description?: string;
  compare_type: 'champions' | 'roles' | 'performance' | 'trends';
  creator: GroupUser;
  created_at: string;
  results?: ComparisonResults;
}

export interface ComparisonResults {
  summary: {
    top_performer: string;
    best_metric: string;
    average_win_rate: number;
    total_games_compared: number;
    time_span: string;
  };
  member_stats: { [key: string]: any };
  charts: ChartData[];
  rankings: MemberRanking[];
  insights: string[];
  generated_at: string;
}

export interface ChartData {
  type: 'bar' | 'line' | 'radar' | 'pie';
  title: string;
  labels: string[];
  datasets: ChartDataset[];
  options?: { [key: string]: any };
}

export interface ChartDataset {
  label: string;
  data: number[];
  background_color?: string[];
  border_color?: string[];
}

export interface MemberRanking {
  user_id: number;
  username: string;
  rank: number;
  score: number;
  metric: string;
  change: 'up' | 'down' | 'same';
}

export interface CreateGroupRequest {
  name: string;
  description?: string;
  privacy: 'public' | 'private' | 'invite_only';
}

export interface CreateComparisonRequest {
  name: string;
  description?: string;
  compare_type: 'champions' | 'roles' | 'performance' | 'trends';
  parameters: {
    member_ids: number[];
    time_range: string;
    champions?: number[];
    roles?: string[];
    game_modes?: number[];
    metrics: string[];
    min_games: number;
  };
}

// Service API pour les groupes
class GroupApiService {
  private baseUrl = getApiUrl('/api/groups');

  // Créer un nouveau groupe
  async createGroup(groupData: CreateGroupRequest): Promise<Group> {
    const response = await fetch(this.baseUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(groupData),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Erreur lors de la création du groupe');
    }

    const result = await response.json();
    return result.group;
  }

  // Récupérer les groupes de l'utilisateur
  async getUserGroups(): Promise<Group[]> {
    const response = await fetch(`${this.baseUrl}/my`);

    if (!response.ok) {
      throw new Error('Erreur lors de la récupération des groupes');
    }

    const result = await response.json();
    return result.groups || [];
  }

  // Rechercher des groupes publics
  async searchGroups(query: string, limit = 20): Promise<Group[]> {
    const params = new URLSearchParams({ q: query, limit: limit.toString() });
    const response = await fetch(`${this.baseUrl}/search?${params}`);

    if (!response.ok) {
      throw new Error('Erreur lors de la recherche de groupes');
    }

    const result = await response.json();
    return result.groups || [];
  }

  // Rejoindre un groupe via code d'invitation
  async joinGroup(inviteCode: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/join`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ invite_code: inviteCode }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Erreur lors de la tentative de rejoindre le groupe');
    }
  }

  // Récupérer les détails d'un groupe
  async getGroup(groupId: number): Promise<Group> {
    const response = await fetch(`${this.baseUrl}/${groupId}`);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Groupe non trouvé');
    }

    const result = await response.json();
    return result.group;
  }

  // Récupérer les membres d'un groupe
  async getGroupMembers(groupId: number): Promise<GroupMember[]> {
    const response = await fetch(`${this.baseUrl}/${groupId}/members`);

    if (!response.ok) {
      throw new Error('Erreur lors de la récupération des membres');
    }

    const result = await response.json();
    return result.members || [];
  }

  // Inviter un utilisateur dans un groupe
  async inviteToGroup(groupId: number, email: string, message?: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/${groupId}/invite`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, message }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Erreur lors de l\'envoi de l\'invitation');
    }
  }

  // Retirer un membre du groupe
  async removeMember(groupId: number, userId: number): Promise<void> {
    const response = await fetch(`${this.baseUrl}/${groupId}/members`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ user_id: userId }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Erreur lors de la suppression du membre');
    }
  }

  // Récupérer les statistiques d'un groupe
  async getGroupStats(groupId: number): Promise<GroupStats> {
    const response = await fetch(`${this.baseUrl}/${groupId}/stats`);

    if (!response.ok) {
      throw new Error('Erreur lors de la récupération des statistiques');
    }

    const result = await response.json();
    return result.stats;
  }

  // Créer une comparaison
  async createComparison(groupId: number, comparisonData: CreateComparisonRequest): Promise<GroupComparison> {
    const response = await fetch(`${this.baseUrl}/${groupId}/comparisons`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(comparisonData),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Erreur lors de la création de la comparaison');
    }

    const result = await response.json();
    return result.comparison;
  }

  // Récupérer les comparaisons d'un groupe
  async getGroupComparisons(groupId: number, limit = 20): Promise<GroupComparison[]> {
    const params = new URLSearchParams({ limit: limit.toString() });
    const response = await fetch(`${this.baseUrl}/${groupId}/comparisons?${params}`);

    if (!response.ok) {
      throw new Error('Erreur lors de la récupération des comparaisons');
    }

    const result = await response.json();
    return result.comparisons || [];
  }

  // Récupérer une comparaison spécifique
  async getComparison(groupId: number, comparisonId: number): Promise<GroupComparison> {
    const response = await fetch(`${this.baseUrl}/${groupId}/comparisons/${comparisonId}`);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Comparaison non trouvée');
    }

    const result = await response.json();
    return result.comparison;
  }

  // Régénérer les résultats d'une comparaison
  async regenerateComparison(groupId: number, comparisonId: number): Promise<GroupComparison> {
    const response = await fetch(`${this.baseUrl}/${groupId}/comparisons/${comparisonId}/regenerate`, {
      method: 'POST',
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Erreur lors de la régénération');
    }

    const result = await response.json();
    return result.comparison;
  }

  // Mettre à jour les paramètres d'un groupe
  async updateGroupSettings(groupId: number, settings: any): Promise<void> {
    const response = await fetch(`${this.baseUrl}/${groupId}/settings`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ settings }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Erreur lors de la mise à jour des paramètres');
    }
  }
}

// Instance singleton du service
export const groupApi = new GroupApiService();