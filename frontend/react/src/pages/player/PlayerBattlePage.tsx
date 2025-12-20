import React, { useState, useEffect } from 'react';
import { useNavigate, useOutletContext, Navigate } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { PlayerContextType, Player, PlayerMonster } from '../../types';

export const PlayerBattlePage: React.FC = () => {
    const navigate = useNavigate();
    const { currentPlayer } = useOutletContext<PlayerContextType>();
    const [myMonsters, setMyMonsters] = useState<PlayerMonster[]>([]);
    const [allPlayers, setAllPlayers] = useState<Player[]>([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedMyMonster, setSelectedMyMonster] = useState<number | null>(null);
    const [selectedOpponent, setSelectedOpponent] = useState<Player | null>(null);
    const [opponentMonsters, setOpponentMonsters] = useState<PlayerMonster[]>([]);
    const [selectedOpponentMonster, setSelectedOpponentMonster] = useState<number | null>(null);
    const [loading, setLoading] = useState(true);
    const [battling, setBattling] = useState(false);

    useEffect(() => {
        if (currentPlayer) {
            loadData();
        }
    }, [currentPlayer]);

    const loadData = async () => {
        if (!currentPlayer) return;

        setLoading(true);
        try {
            const [monsters, players] = await Promise.all([
                apiService.getPlayerMonsters(currentPlayer.id),
                apiService.getPlayers()
            ]);
            setMyMonsters(monsters);
            // Filter out current player from the list
            setAllPlayers(players.filter(p => p.id !== currentPlayer.id));
        } catch (error) {
            console.error('Error loading data:', error);
        } finally {
            setLoading(false);
        }
    };

    const selectOpponent = async (opponent: Player) => {
        setSelectedOpponent(opponent);
        setSelectedOpponentMonster(null);
        try {
            const monsters = await apiService.getPlayerMonsters(opponent.id);
            setOpponentMonsters(monsters);
        } catch (error) {
            console.error('Error loading opponent monsters:', error);
        }
    };

    const startBattle = async () => {
        if (!currentPlayer || !selectedMyMonster || !selectedOpponent || !selectedOpponentMonster) {
            return;
        }

        setBattling(true);
        try {
            await apiService.createBattle(
                currentPlayer.id,
                selectedOpponent.id,
                selectedMyMonster,
                selectedOpponentMonster
            );
            // Refresh player data and navigate to dashboard
            await loadData();
            alert('Battle completed! Check your dashboard for results.');
            navigate('/player/dashboard');
        } catch (error) {
            console.error('Error starting battle:', error);
            alert('Failed to start battle. Please try again.');
        } finally {
            setBattling(false);
        }
    };

    if (!currentPlayer) {
        return <Navigate to="/player/login" replace />;
    }

    if (myMonsters.length === 0) {
        return (
            <div className="view">
                <div className="card" style={{ background: '#fef3c7', borderLeft: '4px solid #f59e0b' }}>
                    <h3 style={{ color: '#92400e', marginBottom: '8px' }}>⚠️ No Monsters!</h3>
                    <p style={{ color: '#78350f', marginBottom: '12px' }}>
                        You need at least one monster to battle.
                    </p>
                    <button onClick={() => navigate('/player/profile')} className="btn-primary">
                        Add Monsters to Your Team
                    </button>
                </div>
            </div>
        );
    }

    const filteredPlayers = allPlayers.filter(p =>
        p.username.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const canBattle = selectedMyMonster && selectedOpponent && selectedOpponentMonster;

    return (
        <div className="view">
            <h2 className="battle-title">⚔️ Battle Arena</h2>

            <div className="battle-grid">
                {/* Your Monster Selection */}
                <div className="battle-side blue">
                    <h3 className="side-title">Your Monster</h3>
                    <div className="selection-list">
                        {myMonsters.map(m => (
                            <div
                                key={m.id}
                                onClick={() => setSelectedMyMonster(m.id)}
                                className={`selection-item ${selectedMyMonster === m.id ? 'selected' : ''}`}
                            >
                                <p className="selection-name">{m.nickname}</p>
                                <p className="selection-stats">
                                    HP: {m.hp} | ATK: {m.attack} | DEF: {m.defense} | SPD: {m.speed}
                                </p>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Opponent Selection */}
                <div className="battle-side red">
                    <h3 className="side-title">Select Opponent</h3>

                    {/* Search Bar */}
                    <input
                        type="text"
                        placeholder="Search players..."
                        value={searchTerm}
                        onChange={e => setSearchTerm(e.target.value)}
                        className="input"
                        style={{ marginBottom: '12px' }}
                    />

                    <div className="selection-list" style={{ maxHeight: '400px', overflowY: 'auto' }}>
                        {loading ? (
                            <p style={{ textAlign: 'center', color: '#666' }}>Loading players...</p>
                        ) : filteredPlayers.length === 0 ? (
                            <p style={{ textAlign: 'center', color: '#666' }}>
                                {searchTerm ? 'No players found' : 'No other players available'}
                            </p>
                        ) : (
                            filteredPlayers.map(p => (
                                <div key={p.id}>
                                    <div
                                        onClick={() => selectOpponent(p)}
                                        className={`selection-item ${selectedOpponent?.id === p.id ? 'selected' : ''}`}
                                    >
                                        <p className="selection-name">{p.username}</p>
                                        <p className="selection-stats">⭐ {p.points} points</p>
                                    </div>

                                    {/* Show opponent's monsters if selected */}
                                    {selectedOpponent?.id === p.id && opponentMonsters.length > 0 && (
                                        <div className="sub-selection">
                                            <p style={{ fontSize: '0.75rem', marginBottom: '8px', color: '#999' }}>
                                                Select opponent's monster:
                                            </p>
                                            {opponentMonsters.map(om => (
                                                <div
                                                    key={om.id}
                                                    onClick={() => setSelectedOpponentMonster(om.id)}
                                                    className={`sub-item ${selectedOpponentMonster === om.id ? 'selected' : ''}`}
                                                >
                                                    {om.nickname} (HP: {om.hp} | ATK: {om.attack})
                                                </div>
                                            ))}
                                        </div>
                                    )}

                                    {selectedOpponent?.id === p.id && opponentMonsters.length === 0 && (
                                        <div className="sub-selection">
                                            <p style={{ fontSize: '0.75rem', color: '#ef4444' }}>
                                                This player has no monsters
                                            </p>
                                        </div>
                                    )}
                                </div>
                            ))
                        )}
                    </div>
                </div>
            </div>

            {/* Battle Button */}
            <div className="battle-action">
                <button
                    onClick={startBattle}
                    disabled={!canBattle || battling}
                    className={`btn-battle-start ${!canBattle || battling ? 'disabled' : ''}`}
                >
                    {battling ? '⚔️ BATTLING...' : '⚔️ START BATTLE!'}
                </button>
                {!canBattle && (
                    <p style={{ marginTop: '12px', color: '#666', fontSize: '0.875rem' }}>
                        Select your monster, an opponent, and their monster to battle
                    </p>
                )}
            </div>
        </div>
    );
};