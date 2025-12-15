import React, { useState, useEffect } from 'react';
import './App.css';

const API = {
  GATEWAY: 'http://localhost:8000',
  PLAYER: 'http://localhost:8000/api/players',
  MONSTER: 'http://localhost:8000/api/monster',
  BATTLE: 'http://localhost:8000/api/battles',
  RANKING: 'http://localhost:8000/api/rankings'
};

interface Player {
  id: number;
  username: string;
  points: number;
  created_at: string;
  updated_at: string;
}

interface Monster {
  id: number;
  name: string;
  type1: string;
  type2?: string;
  base_hp: number;
  base_attack: number;
  base_defense: number;
  base_speed: number;
  description?: string;
}

interface PlayerMonster {
  id: number;
  player_id: number;
  monster_id: number;
  nickname: string;
  level: number;
  hp: number;
  attack: number;
  defense: number;
  speed: number;
}

interface Battle {
  id: number;
  player1_id: number;
  player2_id: number;
  monster1_id: number;
  monster2_id: number;
  winner_id: number;
  status: string;
  battle_log: string;
  points_won: number;
  points_lost: number;
  created_at: string;
}

interface LeaderboardEntry {
  player_id: number;
  username: string;
  total_points: number;
  wins: number;
  losses: number;
  win_rate: number;
  rank: number;
}

type View = 'home' | 'profile' | 'battle' | 'battle-result' | 'leaderboard';

const App: React.FC = () => {
  const [view, setView] = useState<View>('home');
  const [currentPlayer, setCurrentPlayer] = useState<Player | null>(null);
  const [players, setPlayers] = useState<Player[]>([]);
  const [monster, setMonster] = useState<Monster[]>([]);
  const [playerMonster, setPlayerMonster] = useState<PlayerMonster[]>([]);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [recentBattles, setRecentBattles] = useState<Battle[]>([]);
  const [loading, setLoading] = useState(false);
  const [dataLoading, setDataLoading] = useState(true);

  useEffect(() => {
    console.log('App mounted, loading initial data...');
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    setDataLoading(true);
    try {
      console.log('Starting to load initial data...');

      // Load players and monster (these services are running)
      const [playersRes, monsterRes] = await Promise.all([
        fetch(`${API.PLAYER}`),
        fetch(`${API.MONSTER}`)
      ]);

      console.log('Players response status:', playersRes.status);
      console.log('Monster response status:', monsterRes.status);

      if (playersRes.ok) {
        const playersData = await playersRes.json();
        console.log('Loaded players:', playersData);
        setPlayers(playersData || []);
      } else {
        console.error('Failed to load players:', playersRes.status, await playersRes.text());
        setPlayers([]);
      }

      if (monsterRes.ok) {
        const monsterData = await monsterRes.json();
        console.log('Loaded monster:', monsterData);
        setMonster(monsterData || []);
      } else {
        console.error('Failed to load monster:', monsterRes.status, await monsterRes.text());
        setMonster([]);
      }

      // Try to load leaderboard (ranking service might not be running)
      try {
        const leaderboardRes = await fetch(`${API.RANKING}/rankings`);
        if (leaderboardRes.ok) {
          const leaderboardData = await leaderboardRes.json();
          setLeaderboard(leaderboardData || []);
        }
      } catch (error) {
        console.log('Ranking service not available:', error);
        setLeaderboard([]);
      }
    } catch (error) {
      console.error('Error loading data:', error);
      setPlayers([]);
      setMonster([]);
      setLeaderboard([]);
    } finally {
      setDataLoading(false);
    }
  };

  const createPlayer = async (username: string) => {
    try {
      console.log('Creating player:', username);
      const response = await fetch(`${API.PLAYER}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username })
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Error:', errorText);
        alert(`Failed to create player: ${errorText}`);
        return;
      }

      const newPlayer: Player = await response.json();
      console.log('Player created:', newPlayer);

      // Refresh the player list to make sure we have the latest data
      await loadInitialData();

      setCurrentPlayer(newPlayer);
      setView('profile');
    } catch (error) {
      console.error('Error creating player:', error);
      alert(`Error: ${error}`);
    }
  };

  const selectPlayer = async (playerId: number) => {
    try {
      const [playerRes, monsterRes] = await Promise.all([
        fetch(`${API.PLAYER}/${playerId}`),
        fetch(`${API.PLAYER}/${playerId}/monster`)
      ]);

      const player: Player = await playerRes.json();
      const monster: PlayerMonster[] = await monsterRes.json();

      setCurrentPlayer(player);
      setPlayerMonster(monster);
      setView('profile');
    } catch (error) {
      console.error('Error selecting player:', error);
    }
  };

  const addMonsterToPlayer = async (monsterId: number) => {
    if (!currentPlayer) return;

    const monsterData = monster.find(p => p.id === monsterId);
    if (!monsterData) return;

    try {
      const response = await fetch(`${API.PLAYER}/${currentPlayer.id}/monster`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          monster_id: monsterId,
          nickname: monsterData.name,
          level: 1,
          hp: monsterData.base_hp,
          attack: monsterData.base_attack,
          defense: monsterData.base_defense,
          speed: monsterData.base_speed
        })
      });

      const newMonster: PlayerMonster = await response.json();
      setPlayerMonster([...playerMonster, newMonster]);
    } catch (error) {
      console.error('Error adding Monster:', error);
    }
  };

  const startBattle = async (opponentId: number, myMonsterId: number, opponentMonsterId: number) => {
    if (!currentPlayer) return;

    setLoading(true);
    try {
      const response = await fetch(`${API.BATTLE}/battles`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          player1_id: currentPlayer.id,
          player2_id: opponentId,
          monster1_id: myMonsterId,
          monster2_id: opponentMonsterId
        })
      });

      const battle: Battle = await response.json();

      await loadInitialData();
      await selectPlayer(currentPlayer.id);

      setRecentBattles([battle, ...recentBattles]);
      setView('battle-result');
    } catch (error) {
      console.error('Error starting battle:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app">
      <nav className="navbar">
        <div className="nav-container">
          <h1 className="logo">🎮 Maushold</h1>
          <div className="nav-buttons">
            <button onClick={() => setView('home')} className="nav-btn">Home</button>
            <button onClick={() => setView('leaderboard')} className="nav-btn">Leaderboard</button>
            {currentPlayer && (
              <button onClick={() => setView('profile')} className="nav-btn-profile">
                {currentPlayer.username}
              </button>
            )}
          </div>
        </div>
      </nav>

      <div className="container">
        {view === 'home' && <HomeView players={players} createPlayer={createPlayer} selectPlayer={selectPlayer} refreshData={loadInitialData} dataLoading={dataLoading} />}
        {view === 'profile' && currentPlayer && (
          <ProfileView
            currentPlayer={currentPlayer}
            playerMonster={playerMonster}
            monster={monster}
            addMonsterToPlayer={addMonsterToPlayer}
            setView={setView}
          />
        )}
        {view === 'battle' && currentPlayer && (
          <BattleView
            currentPlayer={currentPlayer}
            playerMonster={playerMonster}
            players={players}
            startBattle={startBattle}
            loading={loading}
          />
        )}
        {view === 'leaderboard' && <LeaderboardView leaderboard={leaderboard} />}
        {view === 'battle-result' && recentBattles.length > 0 && currentPlayer && (
          <BattleResultView battle={recentBattles[0]} currentPlayer={currentPlayer} setView={setView} />
        )}
      </div>
    </div>
  );
};

// Home View Component
const HomeView: React.FC<{
  players: Player[];
  createPlayer: (username: string) => void;
  selectPlayer: (playerId: number) => void;
  refreshData?: () => void;
  dataLoading?: boolean;
}> = ({ players, createPlayer, selectPlayer, refreshData, dataLoading }) => {
  const [showCreate, setShowCreate] = useState(false);
  const [username, setUsername] = useState('');

  const handleCreate = () => {
    if (username.trim()) {
      createPlayer(username);
      setShowCreate(false);
      setUsername('');
    }
  };

  return (
    <div className="view">
      <div className="header">
        <h2 className="title">Welcome to Maushold</h2>
        <p className="subtitle">Battle with Monster and climb the rankings!</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Select Your Player</h3>
          <div className="button-group">
            <button onClick={() => setShowCreate(!showCreate)} className="btn-primary">
              Create New Player
            </button>
            {refreshData && (
              <button onClick={refreshData} className="btn-secondary">
                🔄 Refresh
              </button>
            )}
          </div>
        </div>

        {showCreate && (
          <div className="create-form">
            <input
              type="text"
              placeholder="Enter username"
              value={username}
              onChange={e => setUsername(e.target.value)}
              onKeyPress={e => e.key === 'Enter' && handleCreate()}
              className="input"
            />
            <button onClick={handleCreate} className="btn-submit">Create</button>
          </div>
        )}

        <div className="player-grid">
          {dataLoading ? (
            <p>Loading players...</p>
          ) : players.length === 0 ? (
            <p>No players found. Create the first player!</p>
          ) : (
            players.map(player => (
              <div
                key={player.id}
                onClick={() => selectPlayer(player.id)}
                className="player-card"
              >
                <h4 className="player-name">{player.username}</h4>
                <p className="player-points">⭐ {player.points} points</p>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
};

// Profile View Component
const ProfileView: React.FC<{
  currentPlayer: Player;
  playerMonster: PlayerMonster[];
  monster: Monster[];
  addMonsterToPlayer: (monsterId: number) => void;
  setView: (view: View) => void;
}> = ({ currentPlayer, playerMonster, monster, addMonsterToPlayer, setView }) => {
  const [showAdd, setShowAdd] = useState(false);

  return (
    <div className="view">
      <div className="profile-header">
        <h2 className="profile-name">{currentPlayer.username}</h2>
        <p className="profile-points">⭐ {currentPlayer.points} Points</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Your Monster Team</h3>
          <div className="button-group">
            <button onClick={() => setShowAdd(!showAdd)} className="btn-primary">
              Add Monster
            </button>
            {playerMonster.length >= 1 && (
              <button onClick={() => setView('battle')} className="btn-battle">
                ⚔️ Battle!
              </button>
            )}
          </div>
        </div>

        {showAdd && (
          <div className="add-monster">
            <h4 className="section-title">Available Monster</h4>
            <div className="monster-grid">
              {monster.map(p => (
                <div
                  key={p.id}
                  onClick={() => { addMonsterToPlayer(p.id); setShowAdd(false); }}
                  className="monster-card"
                >
                  <p className="monster-name">{p.name}</p>
                  <p className="monster-type">{p.type1}</p>
                </div>
              ))}
            </div>
          </div>
        )}

        {playerMonster.length === 0 ? (
          <p className="empty-message">No Monster yet. Add some to your team!</p>
        ) : (
          <div className="team-grid">
            {playerMonster.map(p => (
              <div key={p.id} className="team-card">
                <h4 className="team-name">{p.nickname}</h4>
                <div className="stats">
                  <p>❤️ HP: {p.hp}</p>
                  <p>⚔️ Attack: {p.attack}</p>
                  <p>🛡️ Defense: {p.defense}</p>
                  <p>⚡ Speed: {p.speed}</p>
                  <p>📊 Level: {p.level}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

// Battle View Component
const BattleView: React.FC<{
  currentPlayer: Player;
  playerMonster: PlayerMonster[];
  players: Player[];
  startBattle: (opponentId: number, myMonsterId: number, opponentMonsterId: number) => void;
  loading: boolean;
}> = ({ currentPlayer, playerMonster, players, startBattle, loading }) => {
  const [selectedMyMonster, setSelectedMyMonster] = useState<number | null>(null);
  const [selectedOpponent, setSelectedOpponent] = useState<number | null>(null);
  const [opponentMonster, setOpponentMonster] = useState<PlayerMonster[]>([]);
  const [selectedOpponentMonster, setSelectedOpponentMonster] = useState<number | null>(null);

  const selectOpponent = async (opponentId: number) => {
    setSelectedOpponent(opponentId);
    try {
      const response = await fetch(`${API.PLAYER}/${opponentId}/monster`);
      const monster: PlayerMonster[] = await response.json();
      setOpponentMonster(monster);
    } catch (error) {
      console.error('Error fetching opponent Monster:', error);
    }
  };

  const canBattle = selectedMyMonster && selectedOpponent && selectedOpponentMonster;

  return (
    <div className="view">
      <h2 className="battle-title">⚔️ Battle Arena</h2>

      <div className="battle-grid">
        <div className="battle-side blue">
          <h3 className="side-title">Your Monster</h3>
          <div className="selection-list">
            {playerMonster.map(p => (
              <div
                key={p.id}
                onClick={() => setSelectedMyMonster(p.id)}
                className={`selection-item ${selectedMyMonster === p.id ? 'selected' : ''}`}
              >
                <p className="selection-name">{p.nickname}</p>
                <p className="selection-stats">HP: {p.hp} | ATK: {p.attack} | DEF: {p.defense}</p>
              </div>
            ))}
          </div>
        </div>

        <div className="battle-side red">
          <h3 className="side-title">Select Opponent</h3>
          <div className="selection-list">
            {players.filter(p => p.id !== currentPlayer.id).map(p => (
              <div key={p.id}>
                <div
                  onClick={() => selectOpponent(p.id)}
                  className={`selection-item ${selectedOpponent === p.id ? 'selected' : ''}`}
                >
                  <p className="selection-name">{p.username}</p>
                  <p className="selection-stats">⭐ {p.points} points</p>
                </div>

                {selectedOpponent === p.id && opponentMonster.length > 0 && (
                  <div className="sub-selection">
                    {opponentMonster.map(op => (
                      <div
                        key={op.id}
                        onClick={() => setSelectedOpponentMonster(op.id)}
                        className={`sub-item ${selectedOpponentMonster === op.id ? 'selected' : ''}`}
                      >
                        {op.nickname} (HP: {op.hp})
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
          onClick={() => selectedMyMonster && selectedOpponent && selectedOpponentMonster &&
            startBattle(selectedOpponent, selectedMyMonster, selectedOpponentMonster)}
          disabled={!canBattle || loading}
          className={`btn-battle-start ${!canBattle || loading ? 'disabled' : ''}`}
        >
          {loading ? 'Battling...' : '⚔️ START BATTLE!'}
        </button>
      </div>
    </div>
  );
};

// Battle Result View
const BattleResultView: React.FC<{
  battle: Battle;
  currentPlayer: Player;
  setView: (view: View) => void;
}> = ({ battle, currentPlayer, setView }) => {
  const isWinner = battle.winner_id === currentPlayer.id;

  return (
    <div className="view">
      <div className={`result-header ${isWinner ? 'win' : 'lose'}`}>
        <h2 className="result-title">{isWinner ? '🎉 Victory!' : '💔 Defeat'}</h2>
        <p className="result-points">{isWinner ? `+${battle.points_won}` : `-${battle.points_lost}`} Points</p>
      </div>

      <div className="card">
        <h3 className="card-title">Battle Log</h3>
        <pre className="battle-log">{battle.battle_log}</pre>
      </div>

      <div className="result-actions">
        <button onClick={() => setView('battle')} className="btn-secondary">Battle Again</button>
        <button onClick={() => setView('profile')} className="btn-primary">Back to Profile</button>
      </div>
    </div>
  );
};

// Leaderboard View
const LeaderboardView: React.FC<{ leaderboard: LeaderboardEntry[] }> = ({ leaderboard }) => {
  return (
    <div className="view">
      <h2 className="page-title">🏆 Global Leaderboard</h2>

      <div className="card">
        <table className="leaderboard-table">
          <thead>
            <tr>
              <th>Rank</th>
              <th>Player</th>
              <th>Points</th>
              <th>Battles</th>
              <th>Win Rate</th>
            </tr>
          </thead>
          <tbody>
            {leaderboard.map((entry, i) => (
              <tr key={entry.player_id} className={i < 3 ? 'top-three' : ''}>
                <td>
                  {i === 0 && '🥇'}
                  {i === 1 && '🥈'}
                  {i === 2 && '🥉'}
                  {i >= 3 && `#${i + 1}`}
                </td>
                <td className="player-name-cell">{entry.username}</td>
                <td className="points-cell">{entry.total_points}</td>
                <td>{entry.wins + entry.losses}</td>
                <td>{entry.win_rate ? entry.win_rate.toFixed(1) : '0.0'}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default App;