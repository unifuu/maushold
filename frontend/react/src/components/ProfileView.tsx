import React, { useState } from 'react';
import type { Player, PlayerMonster, Monster, View } from '../types';

interface ProfileViewProps {
  currentPlayer: Player;
  playerMonster: PlayerMonster[];
  monster: Monster[];
  addMonsterToPlayer: (monsterId: number) => void;
  setView: (view: View) => void;
}

export const ProfileView: React.FC<ProfileViewProps> = ({
  currentPlayer,
  playerMonster,
  monster,
  addMonsterToPlayer,
  setView
}) => {
  const [showAdd, setShowAdd] = useState(false);

  return (
    <div className="view">
      <div className="profile-header">
        <h2 className="profile-name">{currentPlayer.username}</h2>
        <p className="profile-points">⭐ {currentPlayer.points} Points</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Your Monster Team</h3>
          <div className="button-group">
            <button onClick={() => setShowAdd(!showAdd)} className="btn-primary">
              Add Monster
            </button>
            {playerMonster.length >= 1 && (
              <button onClick={() => setView('battle')} className="btn-battle">
                ⚔️ Battle!
              </button>
            )}
          </div>
        </div>

        {showAdd && (
          <div className="add-monster">
            <h4 className="section-title">Available Monster</h4>
            <div className="monster-grid">
              {monster.map(p => (
                <div
                  key={p.id}
                  onClick={() => {
                    addMonsterToPlayer(p.id);
                    setShowAdd(false);
                  }}
                  className="monster-card"
                >
                  <p className="monster-name">{p.name}</p>
                  <p className="monster-type">{p.type1}</p>
                </div>
              ))}
            </div>
          </div>
        )}

        {playerMonster.length === 0 ? (
          <p className="empty-message">No Monster yet. Add some to your team!</p>
        ) : (
          <div className="team-grid">
            {playerMonster.map(p => (
              <div key={p.id} className="team-card">
                <h4 className="team-name">{p.nickname}</h4>
                <div className="stats">
                  <p>❤️ HP: {p.hp}</p>
                  <p>⚔️ Attack: {p.attack}</p>
                  <p>🛡️ Defense: {p.defense}</p>
                  <p>⚡ Speed: {p.speed}</p>
                  <p>📊 Level: {p.level}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};