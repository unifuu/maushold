import React, { useState, useEffect } from 'react';
import { useNavigate, useOutletContext, Navigate } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { PlayerContextType, PlayerMonster, Monster } from '../../types';

export const PlayerProfilePage: React.FC = () => {
    const navigate = useNavigate();
    const { currentPlayer } = useOutletContext<PlayerContextType>();
    const [myMonsters, setMyMonsters] = useState<PlayerMonster[]>([]);
    const [availableMonsters, setAvailableMonsters] = useState<Monster[]>([]);
    const [showAdd, setShowAdd] = useState(false);
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
            const [monsters, allMonsters] = await Promise.all([
                apiService.getPlayerMonsters(currentPlayer.id),
                apiService.getMonsters()
            ]);
            setMyMonsters(monsters);
            setAvailableMonsters(allMonsters);
        } catch (error) {
            console.error('Error loading data:', error);
        } finally {
            setLoading(false);
        }
    };

    const addMonster = async (monsterId: number) => {
        if (!currentPlayer) return;

        const monsterData = availableMonsters.find(m => m.id === monsterId);
        if (!monsterData) return;

        try {
            await apiService.addMonsterToPlayer(currentPlayer.id, {
                monster_id: monsterId,
                nickname: monsterData.name,
                level: 1,
                hp: monsterData.base_hp,
                attack: monsterData.base_attack,
                defense: monsterData.base_defense,
                speed: monsterData.base_speed
            });
            await loadData();
            setShowAdd(false);
        } catch (error) {
            console.error('Error adding monster:', error);
        }
    };

    const getMonsterDetails = (pm: PlayerMonster) => {
        const monsterData = availableMonsters.find(m => m.id === pm.monster_id);
        return {
            ...pm,
            monsterName: monsterData?.name || 'Unknown',
            monsterType: monsterData?.type1 || 'Unknown'
        };
    };

    if (!currentPlayer) {
        return <Navigate to="/player/login" replace />;
    }

    return (
        <div className="view">
            <div className="profile-header">
                <h2 className="profile-name">{currentPlayer.username}</h2>
                <p className="profile-points">⭐ {currentPlayer.points} Points</p>
            </div>

            <div className="card">
                <div className="card-header">
                    <h3 className="card-title">My Monster Team</h3>
                    <div className="button-group">
                        <button onClick={() => setShowAdd(!showAdd)} className="btn-primary">
                            ➕ Add Monster
                        </button>
                        {myMonsters.length >= 1 && (
                            <button onClick={() => navigate('/player/battle')} className="btn-battle">
                                ⚔️ Battle!
                            </button>
                        )}
                    </div>
                </div>

                {showAdd && (
                    <div className="add-monster">
                        <h4 className="section-title">Available Monsters</h4>
                        <div className="monster-grid">
                            {availableMonsters.map(m => (
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

                {myMonsters.length === 0 ? (
                    <p className="empty-message">No monsters yet. Add some to your team!</p>
                ) : (
                    <div className="team-grid">
                        {myMonsters.map(pm => {
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