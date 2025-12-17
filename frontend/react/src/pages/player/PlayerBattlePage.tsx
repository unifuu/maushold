import React from 'react';
import { useNavigate, useOutletContext, Navigate } from 'react-router-dom';
import type { PlayerContextType } from '../../types';

export const PlayerBattlePage: React.FC = () => {
    const navigate = useNavigate();
    const { currentPlayer } = useOutletContext<PlayerContextType>();

    if (!currentPlayer) {
        return <Navigate to="/player/login" replace />;
    }

    return (
        <div className="view">
            <button onClick={() => navigate('/player/dashboard')} className="btn-secondary" style={{ marginBottom: '16px' }}>
                ← Back to Dashboard
            </button>

            <h2 className="page-title">⚔️ Battle Arena</h2>

            <div className="card">
                <p style={{ textAlign: 'center', padding: '40px', color: '#666' }}>
                    Battle functionality coming soon!
                </p>
                <p style={{ textAlign: 'center', color: '#666' }}>
                    For now, use the Admin Panel to initiate battles.
                </p>
                <div style={{ textAlign: 'center', marginTop: '20px' }}>
                    <button onClick={() => navigate('/admin')} className="btn-primary">
                        Go to Admin Panel
                    </button>
                </div>
            </div>
        </div>
    );
};