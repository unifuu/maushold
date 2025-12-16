import React, { useState } from 'react';
import type { Player } from '../types';

interface HomeViewProps {
  players: Player[];
  createPlayer: (username: string) => void;
  selectPlayer: (playerId: number) => void;
  refreshData?: () => void;
  dataLoading?: boolean;
}

export const HomeView: React.FC<HomeViewProps> = ({
  players,
  createPlayer,
  selectPlayer,
  refreshData,
  dataLoading
}) => {
  const [showCreate, setShowCreate] = useState(false);
  const [username, setUsername] = useState('');

  const handleCreate = () => {
    if (username.trim()) {
      createPlayer(username);
      setShowCreate(false);
      setUsername('');
    }
  };

  return (
    <div className="view">
      <div className="header">
        <h2 className="title">Welcome to Maushold</h2>
        <p className="subtitle">Battle with Monster and climb the rankings!</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Select Your Player</h3>
          <div className="button-group">
            <button onClick={() => setShowCreate(!showCreate)} className="btn-primary">
              Create New Player
            </button>
            {refreshData && (
              <button onClick={refreshData} className="btn-secondary">
                🔄 Refresh
              </button>
            )}
          </div>
        </div>

        {showCreate && (
          <div className="create-form">
            <input
              type="text"
              placeholder="Enter username"
              value={username}
              onChange={e => setUsername(e.target.value)}
              onKeyPress={e => e.key === 'Enter' && handleCreate()}
              className="input"
            />
            <button onClick={handleCreate} className="btn-submit">
              Create
            </button>
          </div>
        )}

        <div className="player-grid">
          {dataLoading ? (
            <p>Loading players...</p>
          ) : players.length === 0 ? (
            <p>No players found. Create the first player!</p>
          ) : (
            players.map(player => (
              <div
                key={player.id}
                onClick={() => selectPlayer(player.id)}
                className="player-card"
              >
                <h4 className="player-name">{player.username}</h4>
                <p className="player-points">⭐ {player.points} points</p>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
};