export const API_CONFIG = {
  BASE_URL: process.env.REACT_APP_API_BASE_URL || 'http://localhost:8000',
  ENDPOINTS: {
    PLAYERS: '/api/players',
    MONSTERS: '/api/monster',
    BATTLES: '/api/battles',
    RANKINGS: '/api/rankings'
  }
};