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

  // Helper function to get monster details
  const getMonsterDetails = (playerMon: PlayerMonster) => {
    const monsterData = monster.find(m => m.id === playerMon.monster_id);
    return {
      ...playerMon,
      monsterName: monsterData?.name || 'Unknown',
      monsterType: monsterData?.type1 || 'Unknown'
    };
  };

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
              {monster.map(m => (
                <div
                  key={m.id}
                  onClick={() => {
                    addMonsterToPlayer(m.id);
                    setShowAdd(false);
                  }}
                  className="monster-card"
                >
                  <p className="monster-name">{m.name}</p>
                  <p className="monster-type">{m.type1}</p>
                  <div style={{ fontSize: '0.75rem', marginTop: '4px', color: '#666' }}>
                    <div>HP: {m.base_hp} | ATK: {m.base_attack}</div>
                    <div>DEF: {m.base_defense} | SPD: {m.base_speed}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {playerMonster.length === 0 ? (
          <p className="empty-message">No Monster yet. Add some to your team!</p>
        ) : (
          <div className="team-grid">
            {playerMonster.map(pm => {
              const details = getMonsterDetails(pm);
              return (
                <div key={pm.id} className="team-card">
                  <h4 className="team-name">{details.nickname || details.monsterName}</h4>
                  <p style={{ fontSize: '0.85rem', color: '#666', marginBottom: '8px' }}>
                    {details.monsterName} ({details.monsterType})
                  </p>
                  <div className="stats">
                    <p>❤️ HP: {pm.hp || 0}</p>
                    <p>⚔️ Attack: {pm.attack || 0}</p>
                    <p>🛡️ Defense: {pm.defense || 0}</p>
                    <p>⚡ Speed: {pm.speed || 0}</p>
                    <p>📊 Level: {pm.level || 1}</p>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};