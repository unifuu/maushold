import React, { useState, useEffect } from 'react';
import './App.css';

const API = {
  PLAYER: 'http://localhost:8001',
  POKEMON: 'http://localhost:8002',
  BATTLE: 'http://localhost:8003',
  RANKING: 'http://localhost:8004'
};

interface Player {
  id: number;
  username: string;
  points: number;
  created_at: string;
  updated_at: string;
}

interface Pokemon {
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

interface PlayerPokemon {
  id: number;
  player_id: number;
  pokemon_id: number;
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
  pokemon1_id: number;
  pokemon2_id: number;
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
  const [pokemon, setPokemon] = useState<Pokemon[]>([]);
  const [playerPokemon, setPlayerPokemon] = useState<PlayerPokemon[]>([]);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [recentBattles, setRecentBattles] = useState<Battle[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    try {
      const [playersRes, pokemonRes, leaderboardRes] = await Promise.all([
        fetch(`${API.PLAYER}/players`),
        fetch(`${API.POKEMON}/pokemon`),
        fetch(`${API.RANKING}/rankings`)
      ]);
      
      setPlayers(await playersRes.json());
      setPokemon(await pokemonRes.json());
      setLeaderboard(await leaderboardRes.json());
    } catch (error) {
      console.error('Error loading data:', error);
    }
  };

  const createPlayer = async (username: string) => {
    try {
      console.log('Creating player:', username);
      const response = await fetch(`${API.PLAYER}/players`, {
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
      setPlayers([...players, newPlayer]);
      setCurrentPlayer(newPlayer);
      setView('profile');
    } catch (error) {
      console.error('Error creating player:', error);
      alert(`Error: ${error}`);
    }
  };

  const selectPlayer = async (playerId: number) => {
    try {
      const [playerRes, pokemonRes] = await Promise.all([
        fetch(`${API.PLAYER}/players/${playerId}`),
        fetch(`${API.PLAYER}/players/${playerId}/pokemon`)
      ]);
      
      const player: Player = await playerRes.json();
      const pokemon: PlayerPokemon[] = await pokemonRes.json();
      
      setCurrentPlayer(player);
      setPlayerPokemon(pokemon);
      setView('profile');
    } catch (error) {
      console.error('Error selecting player:', error);
    }
  };

  const addPokemonToPlayer = async (pokemonId: number) => {
    if (!currentPlayer) return;
    
    const pokemonData = pokemon.find(p => p.id === pokemonId);
    if (!pokemonData) return;

    try {
      const response = await fetch(`${API.PLAYER}/players/${currentPlayer.id}/pokemon`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          pokemon_id: pokemonId,
          nickname: pokemonData.name,
          level: 1,
          hp: pokemonData.base_hp,
          attack: pokemonData.base_attack,
          defense: pokemonData.base_defense,
          speed: pokemonData.base_speed
        })
      });
      
      const newPokemon: PlayerPokemon = await response.json();
      setPlayerPokemon([...playerPokemon, newPokemon]);
    } catch (error) {
      console.error('Error adding Pokemon:', error);
    }
  };

  const startBattle = async (opponentId: number, myPokemonId: number, opponentPokemonId: number) => {
    if (!currentPlayer) return;
    
    setLoading(true);
    try {
      const response = await fetch(`${API.BATTLE}/battles`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          player1_id: currentPlayer.id,
          player2_id: opponentId,
          pokemon1_id: myPokemonId,
          pokemon2_id: opponentPokemonId
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
        {view === 'home' && <HomeView players={players} createPlayer={createPlayer} selectPlayer={selectPlayer} />}
        {view === 'profile' && currentPlayer && (
          <ProfileView 
            currentPlayer={currentPlayer} 
            playerPokemon={playerPokemon} 
            pokemon={pokemon} 
            addPokemonToPlayer={addPokemonToPlayer} 
            setView={setView} 
          />
        )}
        {view === 'battle' && currentPlayer && (
          <BattleView 
            currentPlayer={currentPlayer} 
            playerPokemon={playerPokemon} 
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
}> = ({ players, createPlayer, selectPlayer }) => {
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
        <p className="subtitle">Battle with Pokémon and climb the rankings!</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Select Your Player</h3>
          <button onClick={() => setShowCreate(!showCreate)} className="btn-primary">
            Create New Player
          </button>
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
          {players.map(player => (
            <div 
              key={player.id} 
              onClick={() => selectPlayer(player.id)} 
              className="player-card"
            >
              <h4 className="player-name">{player.username}</h4>
              <p className="player-points">⭐ {player.points} points</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

// Profile View Component
const ProfileView: React.FC<{
  currentPlayer: Player;
  playerPokemon: PlayerPokemon[];
  pokemon: Pokemon[];
  addPokemonToPlayer: (pokemonId: number) => void;
  setView: (view: View) => void;
}> = ({ currentPlayer, playerPokemon, pokemon, addPokemonToPlayer, setView }) => {
  const [showAdd, setShowAdd] = useState(false);

  return (
    <div className="view">
      <div className="profile-header">
        <h2 className="profile-name">{currentPlayer.username}</h2>
        <p className="profile-points">⭐ {currentPlayer.points} Points</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Your Pokémon Team</h3>
          <div className="button-group">
            <button onClick={() => setShowAdd(!showAdd)} className="btn-primary">
              Add Pokémon
            </button>
            {playerPokemon.length >= 1 && (
              <button onClick={() => setView('battle')} className="btn-battle">
                ⚔️ Battle!
              </button>
            )}
          </div>
        </div>

        {showAdd && (
          <div className="add-pokemon">
            <h4 className="section-title">Available Pokémon</h4>
            <div className="pokemon-grid">
              {pokemon.map(p => (
                <div 
                  key={p.id} 
                  onClick={() => { addPokemonToPlayer(p.id); setShowAdd(false); }} 
                  className="pokemon-card"
                >
                  <p className="pokemon-name">{p.name}</p>
                  <p className="pokemon-type">{p.type1}</p>
                </div>
              ))}
            </div>
          </div>
        )}

        {playerPokemon.length === 0 ? (
          <p className="empty-message">No Pokémon yet. Add some to your team!</p>
        ) : (
          <div className="team-grid">
            {playerPokemon.map(p => (
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
  playerPokemon: PlayerPokemon[];
  players: Player[];
  startBattle: (opponentId: number, myPokemonId: number, opponentPokemonId: number) => void;
  loading: boolean;
}> = ({ currentPlayer, playerPokemon, players, startBattle, loading }) => {
  const [selectedMyPokemon, setSelectedMyPokemon] = useState<number | null>(null);
  const [selectedOpponent, setSelectedOpponent] = useState<number | null>(null);
  const [opponentPokemon, setOpponentPokemon] = useState<PlayerPokemon[]>([]);
  const [selectedOpponentPokemon, setSelectedOpponentPokemon] = useState<number | null>(null);

  const selectOpponent = async (opponentId: number) => {
    setSelectedOpponent(opponentId);
    try {
      const response = await fetch(`${API.PLAYER}/players/${opponentId}/pokemon`);
      const pokemon: PlayerPokemon[] = await response.json();
      setOpponentPokemon(pokemon);
    } catch (error) {
      console.error('Error fetching opponent Pokemon:', error);
    }
  };

  const canBattle = selectedMyPokemon && selectedOpponent && selectedOpponentPokemon;

  return (
    <div className="view">
      <h2 className="battle-title">⚔️ Battle Arena</h2>

      <div className="battle-grid">
        <div className="battle-side blue">
          <h3 className="side-title">Your Pokémon</h3>
          <div className="selection-list">
            {playerPokemon.map(p => (
              <div 
                key={p.id} 
                onClick={() => setSelectedMyPokemon(p.id)} 
                className={`selection-item ${selectedMyPokemon === p.id ? 'selected' : ''}`}
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
                
                {selectedOpponent === p.id && opponentPokemon.length > 0 && (
                  <div className="sub-selection">
                    {opponentPokemon.map(op => (
                      <div 
                        key={op.id} 
                        onClick={() => setSelectedOpponentPokemon(op.id)} 
                        className={`sub-item ${selectedOpponentPokemon === op.id ? 'selected' : ''}`}
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
          onClick={() => selectedMyPokemon && selectedOpponent && selectedOpponentPokemon && 
                        startBattle(selectedOpponent, selectedMyPokemon, selectedOpponentPokemon)} 
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