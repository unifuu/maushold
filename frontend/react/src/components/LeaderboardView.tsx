import React from 'react';
import type { LeaderboardEntry } from '../types';

interface LeaderboardViewProps {
  leaderboard: LeaderboardEntry[];
}

export const LeaderboardView: React.FC<LeaderboardViewProps> = ({ leaderboard }) => {
  return (
    <div className="view">
      <h2 className="page-title">🏆 Global Leaderboard</h2>

      <div className="card">
        {leaderboard.length === 0 ? (
          <p className="empty-message">No rankings available yet. Play some battles!</p>
        ) : (
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
        )}
      </div>
    </div>
  );
};