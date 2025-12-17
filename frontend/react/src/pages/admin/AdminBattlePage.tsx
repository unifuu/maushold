import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useOutletContext } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { AdminContextType, Player, PlayerMonster } from '../../types';

export const AdminBattlePage: React.FC = () => {
    const { playerId } = useParams<{ playerId: string }>();
    const navigate = useNavigate();
    const { players } = useOutletContext<AdminContextType>();
    const [player, setPlayer] = useState<Player | null>(null);
    const [playerMonsters, setPlayerMonsters] = useState<PlayerMonster[]>([]);
    const [selectedMyMonster, setSelectedMyMonster] = useState<number | null>(null);
    const [selectedOpponent, setSelectedOpponent] = useState<number | null>(null);
    const [opponentMonsters, setOpponentMonsters] = useState<PlayerMonster[]>([]);
    const [selectedOpponentMonster, setSelectedOpponentMonster] = useState<number | null>(null);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        loadData();
    }, [playerId]);

    const loadData = async () => {
        if (!playerId) return;
        try {
            const [playerData, monstersData] = await Promise.all([
                apiService.getPlayer(parseInt(playerId)),
                apiService.getPlayerMonsters(parseInt(playerId))
            ]);
            setPlayer(playerData);
            setPlayerMonsters(monstersData);
        } catch (error) {
            console.error('Error loading data:', error);
        }
    };

    const selectOpponent = async (opponentId: number) => {
        setSelectedOpponent(opponentId);
        try {
            const monsters = await apiService.getPlayerMonsters(opponentId);
            setOpponentMonsters(monsters);
        } catch (error) {
            console.error('Error fetching opponent monsters:', error);
        }
    };

    const startBattle = async () => {
        if (!player || !selectedMyMonster || !selectedOpponent || !selectedOpponentMonster) return;

        setLoading(true);
        try {
            const battle = await apiService.createBattle(
                player.id,
                selectedOpponent,
                selectedMyMonster,
                selectedOpponentMonster
            );
            navigate(`/admin/battle-result/${battle.id}`);
        } catch (error) {
            console.error('Error starting battle:', error);
            alert('Failed to start battle');
        } finally {
            setLoading(false);
        }
    };

    if (!player) return <div className="view"><p>Loading...</p></div>;

    const canBattle = selectedMyMonster && selectedOpponent && selectedOpponentMonster;

    return (
        <div className="view">
            <button onClick={() => navigate(`/admin/profile/${playerId}`)} className="btn-secondary" style={{ marginBottom: '16px' }}>
                ← Back to Profile
            </button>

            <h2 className="battle-title">⚔️ Battle Arena - {player.username}</h2>

            <div className="battle-grid">
                <div className="battle-side blue">
                    <h3 className="side-title">Your Monster</h3>
                    <div className="selection-list">
                        {playerMonsters.map(m => (
                            <div
                                key={m.id}
                                onClick={() => setSelectedMyMonster(m.id)}
                                className={`selection-item ${selectedMyMonster === m.id ? 'selected' : ''}`}
                            >
                                <p className="selection-name">{m.nickname}</p>
                                <p className="selection-stats">HP: {m.hp} | ATK: {m.attack} | DEF: {m.defense}</p>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="battle-side red">
                    <h3 className="side-title">Select Opponent</h3>
                    <div className="selection-list">
                        {players.filter(p => p.id !== player.id).map(p => (
                            <div key={p.id}>
                                <div
                                    onClick={() => selectOpponent(p.id)}
                                    className={`selection-item ${selectedOpponent === p.id ? 'selected' : ''}`}
                                >
                                    <p className="selection-name">{p.username}</p>
                                    <p className="selection-stats">⭐ {p.points} points</p>
                                </div>

                                {selectedOpponent === p.id && opponentMonsters.length > 0 && (
                                    <div className="sub-selection">
                                        {opponentMonsters.map(om => (
                                            <div
                                                key={om.id}
                                                onClick={() => setSelectedOpponentMonster(om.id)}
                                                className={`sub-item ${selectedOpponentMonster === om.id ? 'selected' : ''}`}
                                            >
                                                {om.nickname} (HP: {om.hp})
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            <div className="battle-action">
                <button
                    onClick={startBattle}
                    disabled={!canBattle || loading}
                    className={`btn-battle-start ${!canBattle || loading ? 'disabled' : ''}`}
                >
                    {loading ? 'Battling...' : '⚔️ START BATTLE!'}
                </button>
            </div>
        </div>
    );
};