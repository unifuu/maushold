import React, { useState } from 'react';
import { useNavigate, useOutletContext } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { PlayerContextType } from '../../types';

export const PlayerLoginPage: React.FC = () => {
    const navigate = useNavigate();
    const { setCurrentPlayer } = useOutletContext<PlayerContextType>();
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [isLogin, setIsLogin] = useState(true);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            if (isLogin) {
                // For now, just find player by username (you'll need proper authentication later)
                const players = await apiService.getPlayers();
                const player = players.find(p => p.username === username);

                if (player) {
                    setCurrentPlayer(player);
                    localStorage.setItem('currentPlayer', JSON.stringify(player));
                    navigate('/player/dashboard');
                } else {
                    setError('Player not found');
                }
            } else {
                // Register new player
                const newPlayer = await apiService.createPlayer(username);
                setCurrentPlayer(newPlayer);
                localStorage.setItem('currentPlayer', JSON.stringify(newPlayer));
                navigate('/player/dashboard');
            }
        } catch (err) {
            setError('An error occurred. Please try again.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="view">
            <div style={{ maxWidth: '400px', margin: '0 auto', paddingTop: '60px' }}>
                <div className="card">
                    <h2 className="card-title" style={{ textAlign: 'center', marginBottom: '24px' }}>
                        {isLogin ? '🎮 Player Login' : '🎮 Create Account'}
                    </h2>

                    <form onSubmit={handleSubmit}>
                        <div style={{ marginBottom: '16px' }}>
                            <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                                Username
                            </label>
                            <input
                                type="text"
                                value={username}
                                onChange={e => setUsername(e.target.value)}
                                className="input"
                                required
                                disabled={loading}
                            />
                        </div>

                        {isLogin && (
                            <div style={{ marginBottom: '16px' }}>
                                <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                                    Password
                                </label>
                                <input
                                    type="password"
                                    value={password}
                                    onChange={e => setPassword(e.target.value)}
                                    className="input"
                                    disabled={loading}
                                />
                            </div>
                        )}

                        {error && (
                            <div style={{
                                padding: '12px',
                                background: '#fee2e2',
                                color: '#dc2626',
                                borderRadius: '6px',
                                marginBottom: '16px',
                                fontSize: '0.875rem'
                            }}>
                                {error}
                            </div>
                        )}

                        <button
                            type="submit"
                            className="btn-primary"
                            style={{ width: '100%', marginBottom: '12px' }}
                            disabled={loading}
                        >
                            {loading ? 'Processing...' : (isLogin ? 'Login' : 'Create Account')}
                        </button>

                        <button
                            type="button"
                            onClick={() => {
                                setIsLogin(!isLogin);
                                setError('');
                            }}
                            className="btn-secondary"
                            style={{ width: '100%' }}
                            disabled={loading}
                        >
                            {isLogin ? 'Need an account? Register' : 'Have an account? Login'}
                        </button>
                    </form>

                    <div style={{
                        marginTop: '24px',
                        paddingTop: '24px',
                        borderTop: '1px solid #e5e7eb',
                        textAlign: 'center',
                        fontSize: '0.875rem',
                        color: '#666'
                    }}>
                        <p>Demo: Just enter any username to login</p>
                        <p style={{ marginTop: '8px' }}>
                            <a
                                href="/admin"
                                style={{ color: '#6366f1', textDecoration: 'none' }}
                            >
                                Go to Admin Panel →
                            </a>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};