import React, { useState } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import type { Player } from '../types';

export const PlayerLayout: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [currentPlayer, setCurrentPlayer] = useState<Player | null>(() => {
        const saved = localStorage.getItem('currentPlayer');
        return saved ? JSON.parse(saved) : null;
    });

    const isActive = (path: string) => location.pathname === path;

    const handleLogout = () => {
        localStorage.removeItem('currentPlayer');
        setCurrentPlayer(null);
        navigate('/player/login');
    };

    return (
        <div className="app">
            <nav className="navbar">
                <div className="nav-container">
                    <h1
                        className="logo"
                        onClick={() => navigate(currentPlayer ? '/player/dashboard' : '/player/login')}
                        style={{ cursor: 'pointer' }}
                    >
                        🎮 Maushold
                    </h1>
                    <div className="nav-buttons">
                        {currentPlayer ? (
                            <>
                                <button
                                    onClick={() => navigate('/player/dashboard')}
                                    className={`nav-btn ${isActive('/player/dashboard') ? 'active' : ''}`}
                                >
                                    Dashboard
                                </button>
                                <button
                                    onClick={() => navigate('/player/profile')}
                                    className={`nav-btn ${isActive('/player/profile') ? 'active' : ''}`}
                                >
                                    My Profile
                                </button>
                                <button
                                    onClick={() => navigate('/player/battle')}
                                    className={`nav-btn ${isActive('/player/battle') ? 'active' : ''}`}
                                >
                                    Battle
                                </button>
                                <button
                                    onClick={handleLogout}
                                    className="nav-btn"
                                    style={{ background: '#ef4444' }}
                                >
                                    Logout ({currentPlayer.username})
                                </button>
                            </>
                        ) : (
                            <button
                                onClick={() => navigate('/player/login')}
                                className="nav-btn"
                            >
                                Login
                            </button>
                        )}
                    </div>
                </div>
            </nav>

            <div className="container">
                <Outlet context={{ currentPlayer, setCurrentPlayer }} />
            </div>
        </div>
    );
};