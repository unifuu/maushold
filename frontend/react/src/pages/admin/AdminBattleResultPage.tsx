import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { Battle } from '../../types';

export const AdminBattleResultPage: React.FC = () => {
    const { battleId } = useParams<{ battleId: string }>();
    const navigate = useNavigate();
    const [battle, setBattle] = useState<Battle | null>(null);

    useEffect(() => {
        if (battleId) {
            apiService.getBattle(parseInt(battleId))
                .then(setBattle)
                .catch(err => {
                    console.error('Error fetching battle:', err);
                    alert('Failed to load battle results');
                });
        }
    }, [battleId]);

    if (!battle) return <div className="view"><p>Loading battle results...</p></div>;

    return (
        <div className="view">
            <div className="result-header win">
                <h2 className="result-title">🎉 Battle Complete!</h2>
                <p className="result-points">Winner: Player {battle.winner_id}</p>
            </div>

            <div className="card">
                <h3 className="card-title">Battle Log</h3>
                <pre className="battle-log">{battle.battle_log}</pre>
            </div>

            <div className="result-actions">
                <button onClick={() => navigate('/admin')} className="btn-primary">
                    Back to Admin Home
                </button>
            </div>
        </div>
    );
};