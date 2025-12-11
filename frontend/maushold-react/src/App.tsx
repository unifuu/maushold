import React, { useState, useEffect } from 'react';

const API = {
  PLAYER: 'http://localhost:8001',
  POKEMON: 'http://localhost:8002',
  BATTLE: 'http://localhost:8003',
  RANKING: 'http://localhost:8004'
};

interface Player {
  id: number;
  username: string;
  email: string;
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
  image_url?: string;
}

interface PlayerPokemon {
  id: number;
  player_id: number;
  pokemon_id: number;
  nickname: string;
  level: number;
  experience: number;
  hp: number;
  attack: number;
  defense: number;
  speed: number;
  created_at: string;
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
  completed_at?: string;
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

export default function MausholdApp() {
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

  const createPlayer = async (username: string, email: string) => {
    try {
      const response = await fetch(`${API.PLAYER}/players`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email })
      });
      const newPlayer: Player = await response.json();
      setPlayers([...players, newPlayer]);
      setCurrentPlayer(newPlayer);
      return newPlayer;
    } catch (error) {
      console.error('Error creating player:', error);
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
      
      return battle;
    } catch (error) {
      console.error('Error starting battle:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-pink-900 text-white">
      <nav className="bg-black bg-opacity-50 backdrop-blur-lg p-4 shadow-xl">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-3xl font-bold bg-gradient-to-r from-yellow-400 to-pink-500 bg-clip-text text-transparent">
            🎮 Maushold
          </h1>
          <div className="space-x-4">
            <button onClick={() => setView('home')} className="px-4 py-2 hover:bg-white hover:bg-opacity-10 rounded">
              Home
            </button>
            <button onClick={() => setView('leaderboard')} className="px-4 py-2 hover:bg-white hover:bg-opacity-10 rounded">
              Leaderboard
            </button>
            {currentPlayer && (
              <button onClick={() => setView('profile')} className="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded">
                {currentPlayer.username}
              </button>
            )}
          </div>
        </div>
      </nav>

      <div className="container mx-auto p-8">
        {view === 'home' && <HomeView players={players} createPlayer={createPlayer} selectPlayer={selectPlayer} />}
        {view === 'profile' && currentPlayer && <ProfileView currentPlayer={currentPlayer} playerPokemon={playerPokemon} pokemon={pokemon} addPokemonToPlayer={addPokemonToPlayer} setView={setView} />}
        {view === 'battle' && currentPlayer && <BattleView currentPlayer={currentPlayer} playerPokemon={playerPokemon} players={players} startBattle={startBattle} loading={loading} />}
        {view === 'leaderboard' && <LeaderboardView leaderboard={leaderboard} />}
        {view === 'battle-result' && recentBattles.length > 0 && currentPlayer && <BattleResultView battle={recentBattles[0]} currentPlayer={currentPlayer} setView={setView} />}
      </div>
    </div>
  );
}

interface HomeViewProps {
  players: Player[];
  createPlayer: (username: string, email: string) => void;
  selectPlayer: (playerId: number) => void;
}

function HomeView({ players, createPlayer, selectPlayer }: HomeViewProps) {
  const [showCreate, setShowCreate] = useState(false);
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');

  const handleCreate = async () => {
    if (username && email) {
      await createPlayer(username, email);
      setShowCreate(false);
      setUsername('');
      setEmail('');
    }
  };

  return (
    <div className="space-y-8">
      <div className="text-center space-y-4">
        <h2 className="text-5xl font-bold">Welcome to Maushold</h2>
        <p className="text-xl text-gray-300">Battle with Pokémon and climb the rankings!</p>
      </div>

      <div className="bg-white bg-opacity-10 backdrop-blur-lg rounded-xl p-8">
        <div className="flex justify-between items-center mb-6">
          <h3 className="text-2xl font-bold">Select Your Player</h3>
          <button onClick={() => setShowCreate(!showCreate)} className="px-6 py-2 bg-green-600 hover:bg-green-700 rounded-lg">
            Create New Player
          </button>
        </div>

        {showCreate && (
          <div className="mb-6 p-4 bg-black bg-opacity-30 rounded-lg space-y-4">
            <input 
              type="text" 
              placeholder="Username" 
              value={username} 
              onChange={e => setUsername(e.target.value)} 
              className="w-full px-4 py-2 bg-white bg-opacity-20 rounded text-white placeholder-gray-400"
            />
            <input 
              type="email" 
              placeholder="Email" 
              value={email} 
              onChange={e => setEmail(e.target.value)} 
              className="w-full px-4 py-2 bg-white bg-opacity-20 rounded text-white placeholder-gray-400"
            />
            <button onClick={handleCreate} className="px-6 py-2 bg-blue-600 hover:bg-blue-700 rounded">
              Create
            </button>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {players.map(player => (
            <div 
              key={player.id} 
              onClick={() => selectPlayer(player.id)} 
              className="p-6 bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg cursor-pointer hover:scale-105 transition transform"
            >
              <h4 className="text-xl font-bold">{player.username}</h4>
              <p className="text-yellow-300">⭐ {player.points} points</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

interface ProfileViewProps {
  currentPlayer: Player;
  playerPokemon: PlayerPokemon[];
  pokemon: Pokemon[];
  addPokemonToPlayer: (pokemonId: number) => void;
  setView: (view: View) => void;
}

function ProfileView({ currentPlayer, playerPokemon, pokemon, addPokemonToPlayer, setView }: ProfileViewProps) {
  const [showAdd, setShowAdd] = useState(false);

  return (
    <div className="space-y-8">
      <div className="bg-gradient-to-r from-purple-600 to-pink-600 rounded-xl p-8 shadow-2xl">
        <h2 className="text-4xl font-bold mb-2">{currentPlayer.username}</h2>
        <p className="text-2xl text-yellow-300">⭐ {currentPlayer.points} Points</p>
        <p className="text-gray-200">{currentPlayer.email}</p>
      </div>

      <div className="bg-white bg-opacity-10 backdrop-blur-lg rounded-xl p-8">
        <div className="flex justify-between items-center mb-6">
          <h3 className="text-2xl font-bold">Your Pokémon Team</h3>
          <div className="space-x-4">
            <button onClick={() => setShowAdd(!showAdd)} className="px-6 py-2 bg-green-600 hover:bg-green-700 rounded-lg">
              Add Pokémon
            </button>
            {playerPokemon.length >= 1 && (
              <button onClick={() => setView('battle')} className="px-6 py-2 bg-red-600 hover:bg-red-700 rounded-lg">
                ⚔️ Battle!
              </button>
            )}
          </div>
        </div>

        {showAdd && (
          <div className="mb-6 p-4 bg-black bg-opacity-30 rounded-lg">
            <h4 className="text-xl mb-4">Available Pokémon</h4>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
              {pokemon.map(p => (
                <div 
                  key={p.id} 
                  onClick={() => { addPokemonToPlayer(p.id); setShowAdd(false); }} 
                  className="p-4 bg-gradient-to-br from-blue-500 to-purple-500 rounded-lg cursor-pointer hover:scale-105 transition"
                >
                  <p className="font-bold text-center">{p.name}</p>
                  <p className="text-sm text-center text-gray-200">{p.type1}</p>
                </div>
              ))}
            </div>
          </div>
        )}

        {playerPokemon.length === 0 ? (
          <p className="text-gray-400 text-center py-8">No Pokémon yet. Add some to your team!</p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {playerPokemon.map(p => (
              <div key={p.id} className="bg-gradient-to-br from-indigo-600 to-purple-600 rounded-lg p-6 shadow-lg">
                <h4 className="text-xl font-bold mb-2">{p.nickname}</h4>
                <div className="space-y-1 text-sm">
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
}

interface BattleViewProps {
  currentPlayer: Player;
  playerPokemon: PlayerPokemon[];
  players: Player[];
  startBattle: (opponentId: number, myPokemonId: number, opponentPokemonId: number) => void;
  loading: boolean;
}

function BattleView({ currentPlayer, playerPokemon, players, startBattle, loading }: BattleViewProps) {
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
    <div className="space-y-8">
      <h2 className="text-4xl font-bold text-center">⚔️ Battle Arena</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="bg-blue-900 bg-opacity-50 rounded-xl p-6">
          <h3 className="text-2xl font-bold mb-4">Your Pokémon</h3>
          <div className="space-y-4">
            {playerPokemon.map(p => (
              <div 
                key={p.id} 
                onClick={() => setSelectedMyPokemon(p.id)} 
                className={`p-4 rounded-lg cursor-pointer transition ${selectedMyPokemon === p.id ? 'bg-green-600' : 'bg-white bg-opacity-10 hover:bg-opacity-20'}`}
              >
                <p className="font-bold">{p.nickname}</p>
                <p className="text-sm">HP: {p.hp} | ATK: {p.attack} | DEF: {p.defense}</p>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-red-900 bg-opacity-50 rounded-xl p-6">
          <h3 className="text-2xl font-bold mb-4">Select Opponent</h3>
          <div className="space-y-4">
            {players.filter(p => p.id !== currentPlayer.id).map(p => (
              <div key={p.id}>
                <div 
                  onClick={() => selectOpponent(p.id)} 
                  className={`p-4 rounded-lg cursor-pointer transition ${selectedOpponent === p.id ? 'bg-red-600' : 'bg-white bg-opacity-10 hover:bg-opacity-20'}`}
                >
                  <p className="font-bold">{p.username}</p>
                  <p className="text-sm">⭐ {p.points} points</p>
                </div>
                
                {selectedOpponent === p.id && opponentPokemon.length > 0 && (
                  <div className="ml-4 mt-2 space-y-2">
                    {opponentPokemon.map(op => (
                      <div 
                        key={op.id} 
                        onClick={() => setSelectedOpponentPokemon(op.id)} 
                        className={`p-3 rounded cursor-pointer text-sm ${selectedOpponentPokemon === op.id ? 'bg-yellow-600' : 'bg-white bg-opacity-10 hover:bg-opacity-20'}`}
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

      <div className="text-center">
        <button 
          onClick={() => selectedMyPokemon && selectedOpponent && selectedOpponentPokemon && startBattle(selectedOpponent, selectedMyPokemon, selectedOpponentPokemon)} 
          disabled={!canBattle || loading} 
          className={`px-12 py-4 text-2xl font-bold rounded-xl transition ${canBattle && !loading ? 'bg-gradient-to-r from-red-600 to-yellow-600 hover:from-red-700 hover:to-yellow-700' : 'bg-gray-600 cursor-not-allowed'}`}
        >
          {loading ? 'Battling...' : '⚔️ START BATTLE!'}
        </button>
      </div>
    </div>
  );
}

interface BattleResultViewProps {
  battle: Battle;
  currentPlayer: Player;
  setView: (view: View) => void;
}

function BattleResultView({ battle, currentPlayer, setView }: BattleResultViewProps) {
  const isWinner = battle.winner_id === currentPlayer.id;

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div className={`text-center p-8 rounded-xl ${isWinner ? 'bg-gradient-to-r from-green-600 to-blue-600' : 'bg-gradient-to-r from-red-600 to-gray-600'}`}>
        <h2 className="text-5xl font-bold mb-4">{isWinner ? '🎉 Victory!' : '💔 Defeat'}</h2>
        <p className="text-2xl">{isWinner ? `+${battle.points_won}` : `-${battle.points_lost}`} Points</p>
      </div>

      <div className="bg-white bg-opacity-10 backdrop-blur-lg rounded-xl p-6">
        <h3 className="text-2xl font-bold mb-4">Battle Log</h3>
        <pre className="bg-black bg-opacity-50 p-4 rounded overflow-auto whitespace-pre-wrap text-sm font-mono">
          {battle.battle_log}
        </pre>
      </div>

      <div className="flex justify-center space-x-4">
        <button onClick={() => setView('battle')} className="px-8 py-3 bg-red-600 hover:bg-red-700 rounded-lg">
          Battle Again
        </button>
        <button onClick={() => setView('profile')} className="px-8 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg">
          Back to Profile
        </button>
      </div>
    </div>
  );
}

interface LeaderboardViewProps {
  leaderboard: LeaderboardEntry[];
}

function LeaderboardView({ leaderboard }: LeaderboardViewProps) {
  return (
    <div className="space-y-8">
      <h2 className="text-4xl font-bold text-center">🏆 Global Leaderboard</h2>
      
      <div className="bg-white bg-opacity-10 backdrop-blur-lg rounded-xl overflow-hidden">
        <table className="w-full">
          <thead className="bg-black bg-opacity-50">
            <tr>
              <th className="p-4 text-left">Rank</th>
              <th className="p-4 text-left">Player</th>
              <th className="p-4 text-right">Points</th>
              <th className="p-4 text-right">Battles</th>
              <th className="p-4 text-right">Win Rate</th>
            </tr>
          </thead>
          <tbody>
            {leaderboard.map((entry, i) => (
              <tr 
                key={entry.player_id} 
                className={`border-t border-white border-opacity-10 hover:bg-white hover:bg-opacity-5 ${i < 3 ? 'bg-yellow-600 bg-opacity-20' : ''}`}
              >
                <td className="p-4">
                  {i === 0 && '🥇'}
                  {i === 1 && '🥈'}
                  {i === 2 && '🥉'}
                  {i >= 3 && `#${i + 1}`}
                </td>
                <td className="p-4 font-bold">{entry.username}</td>
                <td className="p-4 text-right text-yellow-300">{entry.total_points}</td>
                <td className="p-4 text-right">{entry.wins + entry.losses}</td>
                <td className="p-4 text-right">{entry.win_rate ? entry.win_rate.toFixed(1) : '0.0'}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}