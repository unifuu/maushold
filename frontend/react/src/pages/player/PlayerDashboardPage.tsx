import React, { useState, useEffect } from 'react';
import { useNavigate, useOutletContext, Navigate } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { PlayerContextType, PlayerMonster, Battle } from '../../types';

export const PlayerDashboardPage: React.FC = () => {
    const navigate = useNavigate();
    const { currentPlayer } = useOutletContext<PlayerContextType>();
    const [myMonsters, setMyMonsters] = useState<PlayerMonster[]>([]);
    const [recentBattles, setRecentBattles] = useState<Battle[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (currentPlayer) {
            loadData();
        }
    }, [currentPlayer]);

    const loadData = async () => {
        if (!currentPlayer) return;

        setLoading(true);
        try {
            const monsters = await apiService.getPlayerMonsters(currentPlayer.id);
            setMyMonsters(monsters);
            // Load recent battles if you have that endpoint
            // const battles = await apiService.getPlayerBattles(currentPlayer.id);
            // setRecentBattles(battles);
        } catch (error) {
            console.error('Error loading data:', error);
        } finally {
            setLoading(false);
        }
    };

    if (!currentPlayer) {
        return <Navigate to="/player/login" replace />;
    }

    return (
        <div className="view">
            <div className="header">
                <h2 className="title">Welcome back, {currentPlayer.username}! 👋</h2>
                <p className="subtitle">Your personal dashboard</p>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '20px', marginBottom: '24px' }}>
                <div className="card" style={{ textAlign: 'center', padding: '24px' }}>
                    <h3 style={{ fontSize: '2rem', margin: '0' }}>⭐</h3>
                    <p style={{ fontSize: '2rem', fontWeight: 'bold', margin: '8px 0' }}>{currentPlayer.points}</p>
                    <p style={{ color: '#666', fontSize: '0.875rem' }}>Total Points</p>
                </div>

                <div className="card" style={{ textAlign: 'center', padding: '24px' }}>
                    <h3 style={{ fontSize: '2rem', margin: '0' }}>🎮</h3>
                    <p style={{ fontSize: '2rem', fontWeight: 'bold', margin: '8px 0' }}>{myMonsters.length}</p>
                    <p style={{ color: '#666', fontSize: '0.875rem' }}>Monsters</p>
                </div>

                <div className="card" style={{ textAlign: 'center', padding: '24px' }}>
                    <h3 style={{ fontSize: '2rem', margin: '0' }}>⚔️</h3>
                    <p style={{ fontSize: '2rem', fontWeight: 'bold', margin: '8px 0' }}>{recentBattles.length}</p>
                    <p style={{ color: '#666', fontSize: '0.875rem' }}>Battles</p>
                </div>
            </div>

            <div className="card">
                <h3 className="card-title">Quick Actions</h3>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '12px' }}>
                    <button onClick={() => navigate('/player/profile')} className="btn-primary">
                        📝 Manage My Monsters
                    </button>
                    <button onClick={() => navigate('/player/battle')} className="btn-battle">
                        ⚔️ Start Battle
                    </button>
                    <button onClick={() => navigate('/admin/leaderboard')} className="btn-secondary">
                        🏆 View Leaderboard
                    </button>
                </div>
            </div>

            {myMonsters.length === 0 && (
                <div className="card" style={{ background: '#fef3c7', borderLeft: '4px solid #f59e0b' }}>
                    <h3 style={{ color: '#92400e', marginBottom: '8px' }}>⚠️ No Monsters Yet!</h3>
                    <p style={{ color: '#78350f', marginBottom: '12px' }}>
                        You need to add monsters to your team before you can battle.
                    </p>
                    <button onClick={() => navigate('/player/profile')} className="btn-primary">
                        Add Your First Monster
                    </button>
                </div>
            )}
        </div>
    );
};
