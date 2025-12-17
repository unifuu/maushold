import { API_CONFIG } from '../config/api.config';
import type { Player, Monster, PlayerMonster, Battle, LeaderboardEntry } from '../types';

const { BASE_URL, ENDPOINTS } = API_CONFIG;

class ApiService {
  // Players
  async getPlayers(): Promise<Player[]> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.PLAYERS}`);
    if (!response.ok) throw new Error('Failed to fetch players');
    return response.json();
  }

  async getPlayer(id: number): Promise<Player> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.PLAYERS}/${id}`);
    if (!response.ok) throw new Error('Failed to fetch player');
    return response.json();
  }

  async createPlayer(username: string): Promise<Player> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.PLAYERS}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username })
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText);
    }
    return response.json();
  }

  async deletePlayer(id: number): Promise<void> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.PLAYERS}/${id}`, {
      method: 'DELETE'
    });
    if (!response.ok) {
      throw new Error('Failed to delete player');
    }
  }

  // Monsters
  async getMonsters(): Promise<Monster[]> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.MONSTERS}`);
    if (!response.ok) throw new Error('Failed to fetch monsters');
    return response.json();
  }

  async getPlayerMonsters(playerId: number): Promise<PlayerMonster[]> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.PLAYERS}/${playerId}/monster`);
    if (!response.ok) throw new Error('Failed to fetch player monsters');
    return response.json();
  }

  async addMonsterToPlayer(
    playerId: number,
    monsterData: {
      monster_id: number;
      nickname: string;
      level: number;
      hp: number;
      attack: number;
      defense: number;
      speed: number;
    }
  ): Promise<PlayerMonster> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.PLAYERS}/${playerId}/monster`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(monsterData)
    });
    if (!response.ok) throw new Error('Failed to add monster');
    return response.json();
  }

  // Battles
  async createBattle(
    player1Id: number,
    player2Id: number,
    monster1Id: number,
    monster2Id: number
  ): Promise<Battle> {
    const response = await fetch(`${BASE_URL}${ENDPOINTS.BATTLES}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        player1_id: player1Id,
        player2_id: player2Id,
        monster1_id: monster1Id,
        monster2_id: monster2Id
      })
    });
    if (!response.ok) throw new Error('Failed to create battle');
    return response.json();
  }

  // Rankings
  async getLeaderboard(): Promise<LeaderboardEntry[]> {
    try {
      const response = await fetch(`${BASE_URL}${ENDPOINTS.RANKINGS}`);
      if (!response.ok) return [];
      return response.json();
    } catch {
      return [];
    }
  }
}

export const apiService = new ApiService();