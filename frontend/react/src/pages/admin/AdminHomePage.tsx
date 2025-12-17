import React, { useState } from 'react';
import { useNavigate, useOutletContext } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { AdminContextType } from '../../types';

export const AdminHomePage: React.FC = () => {
    const navigate = useNavigate();
    const { players, loading, refreshData } = useOutletContext<AdminContextType>();
    const [showCreate, setShowCreate] = useState(false);
    const [username, setUsername] = useState('');
    const [creating, setCreating] = useState(false);

    const handleCreate = async () => {
        if (!username.trim()) return;

        setCreating(true);
        try {
            const newPlayer = await apiService.createPlayer(username);
            await refreshData();
            setShowCreate(false);
            setUsername('');
            navigate(`/admin/profile/${newPlayer.id}`);
        } catch (error) {
            console.error('Error creating player:', error);
            alert(`Error: ${error}`);
        } finally {
            setCreating(false);
        }
    };

    return (
        <div className="view">
            <div className="header">
                <h2 className="title">Admin Panel - Player Management</h2>
                <p className="subtitle">Manage all players and monitor game activity</p>
            </div>

            <div className="card">
                <div className="card-header">
                    <h3 className="card-title">All Players</h3>
                    <div className="button-group">
                        <button onClick={() => setShowCreate(!showCreate)} className="btn-primary">
                            ➕ Create New Player
                        </button>
                        <button onClick={refreshData} className="btn-secondary">
                            🔄 Refresh
                        </button>
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
                            disabled={creating}
                        />
                        <button onClick={handleCreate} className="btn-submit" disabled={creating}>
                            {creating ? 'Creating...' : 'Create'}
                        </button>
                    </div>
                )}

                <div className="player-grid">
                    {loading ? (
                        <p>Loading players...</p>
                    ) : players.length === 0 ? (
                        <p>No players found. Create the first player!</p>
                    ) : (
                        players.map(player => (
                            <div
                                key={player.id}
                                onClick={() => navigate(`/admin/profile/${player.id}`)}
                                className="player-card"
                            >
                                <h4 className="player-name">{player.username}</h4>
                                <p className="player-points">⭐ {player.points} points</p>
                                <p style={{ fontSize: '0.75rem', color: '#666', marginTop: '4px' }}>
                                    ID: {player.id}
                                </p>
                            </div>
                        ))
                    )}
                </div>
            </div>
        </div>
    );
};