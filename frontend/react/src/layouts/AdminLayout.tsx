import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { apiService } from '../services/api';
import type { Player, Monster, LeaderboardEntry } from '../types';

export const AdminLayout: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [players, setPlayers] = useState<Player[]>([]);
    const [monsters, setMonsters] = useState<Monster[]>([]);
    const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        try {
            const [playersData, monstersData, leaderboardData] = await Promise.all([
                apiService.getPlayers(),
                apiService.getMonsters(),
                apiService.getLeaderboard()
            ]);
            setPlayers(playersData);
            setMonsters(monstersData);
            setLeaderboard(leaderboardData);
        } catch (error) {
            console.error('Error loading data:', error);
        } finally {
            setLoading(false);
        }
    };

    const isActive = (path: string) => location.pathname === path;

    return (
        <div className="app">
            <nav className="navbar">
                <div className="nav-container">
                    <h1
                        className="logo"
                        onClick={() => navigate('/admin')}
                        style={{ cursor: 'pointer' }}
                    >
                        🎮 Maushold
                    </h1>
                    <div className="nav-buttons">
                        <button
                            onClick={() => navigate('/admin')}
                            className={`nav-btn ${isActive('/admin') ? 'active' : ''}`}
                        >
                            Admin Home
                        </button>
                        <button
                            onClick={() => navigate('/admin/leaderboard')}
                            className={`nav-btn ${isActive('/admin/leaderboard') ? 'active' : ''}`}
                        >
                            Leaderboard
                        </button>
                        <button
                            onClick={() => navigate('/admin/history')}
                            className={`nav-btn ${isActive('/admin/history') ? 'active' : ''}`}
                        >
                            Battle History
                        </button>
                        <button
                            onClick={() => navigate('/player/login')}
                            className="nav-btn"
                            style={{ marginLeft: '20px', background: '#10b981' }}
                        >
                            👤 Player Portal
                        </button>
                    </div>
                </div>
            </nav>

            <div className="container">
                <Outlet context={{ players, monsters, leaderboard, loading, refreshData: loadData }} />
            </div>
        </div>
    );
};