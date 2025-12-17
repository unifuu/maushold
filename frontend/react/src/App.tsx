import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AdminLayout } from './layouts/AdminLayout';
import { PlayerLayout } from './layouts/PlayerLayout';
import { AdminHomePage } from './pages/admin/AdminHomePage';
import { AdminProfilePage } from './pages/admin/AdminProfilePage';
import { AdminBattlePage } from './pages/admin/AdminBattlePage';
import { AdminBattleResultPage } from './pages/admin/AdminBattleResultPage';
import { AdminLeaderboardPage } from './pages/admin/AdminLeaderboardPage';
import { PlayerLoginPage } from './pages/player/PlayerLoginPage';
import { PlayerDashboardPage } from './pages/player/PlayerDashboardPage';
import { PlayerProfilePage } from './pages/player/PlayerProfilePage';
import { PlayerBattlePage } from './pages/player/PlayerBattlePage';
import './App.css';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        {/* Admin Routes */}
        <Route path="/admin" element={<AdminLayout />}>
          <Route index element={<AdminHomePage />} />
          <Route path="profile/:playerId" element={<AdminProfilePage />} />
          <Route path="battle/:playerId" element={<AdminBattlePage />} />
          <Route path="battle-result/:battleId" element={<AdminBattleResultPage />} />
          <Route path="leaderboard" element={<AdminLeaderboardPage />} />
        </Route>

        {/* Player Routes */}
        <Route path="/player" element={<PlayerLayout />}>
          <Route index element={<Navigate to="/player/login" replace />} />
          <Route path="login" element={<PlayerLoginPage />} />
          <Route path="dashboard" element={<PlayerDashboardPage />} />
          <Route path="profile" element={<PlayerProfilePage />} />
          <Route path="battle" element={<PlayerBattlePage />} />
        </Route>

        {/* Default redirect to player login */}
        <Route path="/" element={<Navigate to="/player/login" replace />} />
        <Route path="*" element={<Navigate to="/player/login" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;