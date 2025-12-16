import React, { useState } from 'react';
import { apiService } from '../services/api';
import { API_CONFIG } from '../config/api.config';
import type { Player, PlayerMonster } from '../types';

interface BattleViewProps {
  currentPlayer: Player;
  playerMonster: PlayerMonster[];
  players: Player[];
  startBattle: (opponentId: number, myMonsterId: number, opponentMonsterId: number) => void;
  loading: boolean;
}

export const BattleView: React.FC<BattleViewProps> = ({
  currentPlayer,
  playerMonster,
  players,
  startBattle,
  loading
}) => {
  const [selectedMyMonster, setSelectedMyMonster] = useState<number | null>(null);
  const [selectedOpponent, setSelectedOpponent] = useState<number | null>(null);
  const [opponentMonster, setOpponentMonster] = useState<PlayerMonster[]>([]);
  const [selectedOpponentMonster, setSelectedOpponentMonster] = useState<number | null>(null);

  const selectOpponent = async (opponentId: number) => {
    setSelectedOpponent(opponentId);
    try {
      const monsters = await apiService.getPlayerMonsters(opponentId);
      setOpponentMonster(monsters);
    } catch (error) {
      console.error('Error fetching opponent Monster:', error);
    }
  };

  const canBattle = selectedMyMonster && selectedOpponent && selectedOpponentMonster;

  return (
    <div className="view">
      <h2 className="battle-title">⚔️ Battle Arena</h2>

      <div className="battle-grid">
        <div className="battle-side blue">
          <h3 className="side-title">Your Monster</h3>
          <div className="selection-list">
            {playerMonster.map(p => (
              <div
                key={p.id}
                onClick={() => setSelectedMyMonster(p.id)}
                className={`selection-item ${selectedMyMonster === p.id ? 'selected' : ''}`}
              >
                <p className="selection-name">{p.nickname}</p>
                <p className="selection-stats">HP: {p.hp} | ATK: {p.attack} | DEF: {p.defense}</p>
              </div>
            ))}
          </div>
        </div>

        <div className="battle-side red">
          <h3 className="side-title">Select Opponent</h3>
          <div className="selection-list">
            {players.filter(p => p.id !== currentPlayer.id).map(p => (
              <div key={p.id}>
                <div
                  onClick={() => selectOpponent(p.id)}
                  className={`selection-item ${selectedOpponent === p.id ? 'selected' : ''}`}
                >
                  <p className="selection-name">{p.username}</p>
                  <p className="selection-stats">⭐ {p.points} points</p>
                </div>

                {selectedOpponent === p.id && opponentMonster.length > 0 && (
                  <div className="sub-selection">
                    {opponentMonster.map(op => (
                      <div
                        key={op.id}
                        onClick={() => setSelectedOpponentMonster(op.id)}
                        className={`sub-item ${selectedOpponentMonster === op.id ? 'selected' : ''}`}
                      >
                        {op.nickname} (HP: {op.hp})
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="battle-action">
        <button
          onClick={() =>
            selectedMyMonster &&
            selectedOpponent &&
            selectedOpponentMonster &&
            startBattle(selectedOpponent, selectedMyMonster, selectedOpponentMonster)
          }
          disabled={!canBattle || loading}
          className={`btn-battle-start ${!canBattle || loading ? 'disabled' : ''}`}
        >
          {loading ? 'Battling...' : '⚔️ START BATTLE!'}
        </button>
      </div>
    </div>
  );
};
