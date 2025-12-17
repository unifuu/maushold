export interface Player {
  id: number;
  username: string;
  points: number;
  created_at: string;
  updated_at: string;
}

export interface Monster {
  id: number;
  name: string;
  type1: string;
  type2?: string;
  base_hp: number;
  base_attack: number;
  base_defense: number;
  base_speed: number;
  description?: string;
}

export interface PlayerMonster {
  id: number;
  player_id: number;
  monster_id: number;
  nickname: string;
  level: number;
  hp: number;
  attack: number;
  defense: number;
  speed: number;
}

export interface Battle {
  id: number;
  player1_id: number;
  player2_id: number;
  monster1_id: number;
  monster2_id: number;
  winner_id: number;
  status: string;
  battle_log: string;
  points_won: number;
  points_lost: number;
  created_at: string;
}

export interface LeaderboardEntry {
  player_id: number;
  username: string;
  total_points: number;
  wins: number;
  losses: number;
  win_rate: number;
  rank: number;
}

export interface AdminContextType {
  players: Player[];
  monsters: Monster[];
  leaderboard: LeaderboardEntry[];
  loading: boolean;
  refreshData: () => Promise<void>;
}

export interface PlayerContextType {
  currentPlayer: Player | null;
  setCurrentPlayer: (player: Player | null) => void;
}

export type View = 'home' | 'profile' | 'battle' | 'battle-result' | 'leaderboard';