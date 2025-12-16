import React from 'react';
import type { Battle, Player, View } from '../types';

interface BattleResultViewProps {
  battle: Battle;
  currentPlayer: Player;
  setView: (view: View) => void;
}

export const BattleResultView: React.FC<BattleResultViewProps> = ({
  battle,
  currentPlayer,
  setView
}) => {
  const isWinner = battle.winner_id === currentPlayer.id;

  return (
    <div className="view">
      <div className={`result-header ${isWinner ? 'win' : 'lose'}`}>
        <h2 className="result-title">{isWinner ? '🎉 Victory!' : '💔 Defeat'}</h2>
        <p className="result-points">
          {isWinner ? `+${battle.points_won}` : `-${battle.points_lost}`} Points
        </p>
      </div>

      <div className="card">
        <h3 className="card-title">Battle Log</h3>
        <pre className="battle-log">{battle.battle_log}</pre>
      </div>

      <div className="result-actions">
        <button onClick={() => setView('battle')} className="btn-secondary">
          Battle Again
        </button>
        <button onClick={() => setView('profile')} className="btn-primary">
          Back to Profile
        </button>
      </div>
    </div>
  );
};