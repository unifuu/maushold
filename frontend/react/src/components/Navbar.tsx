import React from 'react';
import type { Player, View } from '../types';

interface NavbarProps {
  currentPlayer: Player | null;
  setView: (view: View) => void;
}

export const Navbar: React.FC<NavbarProps> = ({ currentPlayer, setView }) => {
  return (
    <nav className="navbar">
      <div className="nav-container">
        <h1 className="logo">🎮 Maushold</h1>
        <div className="nav-buttons">
          <button onClick={() => setView('home')} className="nav-btn">
            Home
          </button>
          <button onClick={() => setView('leaderboard')} className="nav-btn">
            Leaderboard
          </button>
          {currentPlayer && (
            <button onClick={() => setView('profile')} className="nav-btn-profile">
              {currentPlayer.username}
            </button>
          )}
        </div>
      </div>
    </nav>
  );
};