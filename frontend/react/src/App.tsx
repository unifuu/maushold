import React, { useState, useEffect } from 'react';
import { Navbar } from './components/Navbar';
import { HomeView } from './components/HomeView';
import { ProfileView } from './components/ProfileView';
import { BattleView } from './components/BattleView';
import { BattleResultView } from './components/BattleResultView';
import { LeaderboardView } from './components/LeaderboardView';
import { apiService } from './services/api';
import type { Player, Monster, PlayerMonster, Battle, LeaderboardEntry, View } from './types';
import './App.css';

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
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    setDataLoading(true);
    try {
      const [playersData, monsterData] = await Promise.all([
        apiService.getPlayers(),
        apiService.getMonsters()
      ]);

      setPlayers(playersData);
      setMonster(monsterData);

      const leaderboardData = await apiService.getLeaderboard();
      setLeaderboard(leaderboardData);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setDataLoading(false);
    }
  };

  const createPlayer = async (username: string) => {
    try {
      const newPlayer = await apiService.createPlayer(username);
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
      const [player, monsters] = await Promise.all([
        apiService.getPlayer(playerId),
        apiService.getPlayerMonsters(playerId)
      ]);
      setCurrentPlayer(player);
      setPlayerMonster(monsters);
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
      const newMonster = await apiService.addMonsterToPlayer(currentPlayer.id, {
        monster_id: monsterId,
        nickname: monsterData.name,
        level: 1,
        hp: monsterData.base_hp,
        attack: monsterData.base_attack,
        defense: monsterData.base_defense,
        speed: monsterData.base_speed
      });
      setPlayerMonster([...playerMonster, newMonster]);
    } catch (error) {
      console.error('Error adding Monster:', error);
    }
  };

  const startBattle = async (
    opponentId: number,
    myMonsterId: number,
    opponentMonsterId: number
  ) => {
    if (!currentPlayer) return;

    setLoading(true);
    try {
      const battle = await apiService.createBattle(
        currentPlayer.id,
        opponentId,
        myMonsterId,
        opponentMonsterId
      );

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
      <Navbar currentPlayer={currentPlayer} setView={setView} />
      <div className="container">
        {view === 'home' && (
          <HomeView
            players={players}
            createPlayer={createPlayer}
            selectPlayer={selectPlayer}
            refreshData={loadInitialData}
            dataLoading={dataLoading}
          />
        )}
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
          <BattleResultView
            battle={recentBattles[0]}
            currentPlayer={currentPlayer}
            setView={setView}
          />
        )}
      </div>
    </div>
  );
};

export default App;