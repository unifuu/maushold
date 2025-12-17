import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useOutletContext } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { AdminContextType, Player, PlayerMonster } from '../../types';

export const AdminProfilePage: React.FC = () => {
    const { playerId } = useParams<{ playerId: string }>();
    const navigate = useNavigate();
    const { monsters, refreshData } = useOutletContext<AdminContextType>();
    const [player, setPlayer] = useState<Player | null>(null);
    const [playerMonsters, setPlayerMonsters] = useState<PlayerMonster[]>([]);
    const [showAdd, setShowAdd] = useState(false);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadPlayerData();
    }, [playerId]);

    const loadPlayerData = async () => {
        if (!playerId) return;

        setLoading(true);
        try {
            const [playerData, monstersData] = await Promise.all([
                apiService.getPlayer(parseInt(playerId)),
                apiService.getPlayerMonsters(parseInt(playerId))
            ]);
            setPlayer(playerData);
            setPlayerMonsters(monstersData);
        } catch (error) {
            console.error('Error loading player:', error);
        } finally {
            setLoading(false);
        }
    };

    const addMonster = async (monsterId: number) => {
        if (!player) return;

        const monsterData = monsters.find(m => m.id === monsterId);
        if (!monsterData) return;

        try {
            await apiService.addMonsterToPlayer(player.id, {
                monster_id: monsterId,
                nickname: monsterData.name,
                level: 1,
                hp: monsterData.base_hp,
                attack: monsterData.base_attack,
                defense: monsterData.base_defense,
                speed: monsterData.base_speed
            });
            await loadPlayerData();
            setShowAdd(false);
        } catch (error) {
            console.error('Error adding monster:', error);
        }
    };

    const getMonsterDetails = (pm: PlayerMonster) => {
        const monsterData = monsters.find(m => m.id === pm.monster_id);
        return {
            ...pm,
            monsterName: monsterData?.name || 'Unknown',
            monsterType: monsterData?.type1 || 'Unknown'
        };
    };

    const handleDelete = async () => {
        if (!player) return;

        const confirmed = window.confirm(
            `Are you sure you want to delete player "${player.username}"? This action cannot be undone and will delete all their monsters.`
        );

        if (!confirmed) return;

        try {
            await apiService.deletePlayer(player.id);
            await refreshData();
            navigate('/admin');
        } catch (error) {
            console.error('Error deleting player:', error);
            alert('Failed to delete player. Please try again.');
        }
    };

    if (loading) return <div className="view"><p>Loading...</p></div>;
    if (!player) return <div className="view"><p>Player not found</p></div>;

    return (
        <div className="view">
            <button onClick={() => navigate('/admin')} className="btn-secondary" style={{ marginBottom: '16px' }}>
                ← Back to Admin Home
            </button>

            <div className="profile-header">
                <div>
                    <h2 className="profile-name">{player.username}</h2>
                    <p className="profile-points">⭐ {player.points} Points</p>
                    <p style={{ fontSize: '0.875rem', color: '#666' }}>Player ID: {player.id}</p>
                </div>
                <button
                    onClick={handleDelete}
                    className="btn-secondary"
                    style={{
                        background: '#ef4444',
                        alignSelf: 'flex-start',
                        marginTop: '8px'
                    }}
                >
                    🗑️ Delete Player
                </button>
            </div>

            <div className="card">
                <div className="card-header">
                    <h3 className="card-title">Monster Team</h3>
                    <div className="button-group">
                        <button onClick={() => setShowAdd(!showAdd)} className="btn-primary">
                            Add Monster
                        </button>
                        {playerMonsters.length >= 1 && (
                            <button onClick={() => navigate(`/admin/battle/${player.id}`)} className="btn-battle">
                                ⚔️ Start Battle
                            </button>
                        )}
                    </div>
                </div>

                {showAdd && (
                    <div className="add-monster">
                        <h4 className="section-title">Available Monsters</h4>
                        <div className="monster-grid">
                            {monsters.map(m => (
                                <div key={m.id} onClick={() => addMonster(m.id)} className="monster-card">
                                    <p className="monster-name">{m.name}</p>
                                    <p className="monster-type">{m.type1}</p>
                                    <div style={{ fontSize: '0.75rem', marginTop: '4px', color: '#666' }}>
                                        <div>HP: {m.base_hp} | ATK: {m.base_attack}</div>
                                        <div>DEF: {m.base_defense} | SPD: {m.base_speed}</div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {playerMonsters.length === 0 ? (
                    <p className="empty-message">No monsters yet. Add some to the team!</p>
                ) : (
                    <div className="team-grid">
                        {playerMonsters.map(pm => {
                            const details = getMonsterDetails(pm);
                            return (
                                <div key={pm.id} className="team-card">
                                    <h4 className="team-name">{details.nickname || details.monsterName}</h4>
                                    <p style={{ fontSize: '0.85rem', color: '#666', marginBottom: '8px' }}>
                                        {details.monsterName} ({details.monsterType})
                                    </p>
                                    <div className="stats">
                                        <p>❤️ HP: {pm.hp || 0}</p>
                                        <p>⚔️ Attack: {pm.attack || 0}</p>
                                        <p>🛡️ Defense: {pm.defense || 0}</p>
                                        <p>⚡ Speed: {pm.speed || 0}</p>
                                        <p>📊 Level: {pm.level || 1}</p>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
};
