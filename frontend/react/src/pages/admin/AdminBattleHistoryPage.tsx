import React, { useState, useEffect } from 'react';
import { useNavigate, useOutletContext } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { Battle, AdminContextType } from '../../types';

export const AdminBattleHistoryPage: React.FC = () => {
    const navigate = useNavigate();
    const { players, monsters, loading: contextLoading } = useOutletContext<AdminContextType>();
    const [battles, setBattles] = useState<Battle[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        apiService.getBattles()
            .then(setBattles)
            .catch(err => console.error('Error fetching battles:', err))
            .finally(() => setLoading(false));
    }, []);

    const getPlayerName = (id: number) => {
        return players.find(p => p.id === id)?.username || `ID: ${id}`;
    };

    const getMonsterName = (id: number) => {
        return monsters.find(m => m.id === id)?.name || `ID: ${id}`;
    };

    const formatDate = (dateStr: string) => {
        return new Date(dateStr).toLocaleString();
    };

    if (loading || contextLoading) return <div className="view"><p>Loading battle history...</p></div>;

    return (
        <div className="view">
            <div className="header">
                <h2 className="title">Battle History</h2>
                <p className="subtitle">Review all past combat results</p>
            </div>

            <div className="card">
                {battles.length === 0 ? (
                    <p className="empty-message">No battles recorded yet.</p>
                ) : (
                    <table className="leaderboard-table">
                        <thead>
                            <tr>
                                <th>Players</th>
                                <th>Monsters</th>
                                <th>Winner</th>
                                <th>Date</th>
                            </tr>
                        </thead>
                        <tbody>
                            {battles.map((battle) => (
                                <tr
                                    key={battle.id}
                                    onClick={() => navigate(`/admin/battle-result/${battle.id}`)}
                                    style={{ cursor: 'pointer' }}
                                >
                                    <td>
                                        <div style={{ display: 'flex', flexDirection: 'column' }}>
                                            <span>{getPlayerName(battle.player1_id)}</span>
                                            <span style={{ fontSize: '0.8rem', color: '#666' }}>vs</span>
                                            <span>{getPlayerName(battle.player2_id)}</span>
                                        </div>
                                    </td>
                                    <td>
                                        <div style={{ display: 'flex', flexDirection: 'column' }}>
                                            <span>{getMonsterName(battle.monster1_id)}</span>
                                            <span style={{ fontSize: '0.8rem', color: '#666' }}>&nbsp;</span>
                                            <span>{getMonsterName(battle.monster2_id)}</span>
                                        </div>
                                    </td>
                                    <td className="points-cell">
                                        <span className="tag winner">
                                            {getPlayerName(battle.winner_id)}
                                        </span>
                                    </td>
                                    <td style={{ fontSize: '0.85rem' }}>{formatDate(battle.created_at)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
};
