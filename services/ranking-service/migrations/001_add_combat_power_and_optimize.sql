-- Migration: Add combat_power and optimize indexes for high-performance ranking
-- Date: 2025-12-18

-- Add combat_power column
ALTER TABLE player_rankings 
ADD COLUMN IF NOT EXISTS combat_power BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS last_battle_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();

-- Update existing records to use total_points as initial combat_power
UPDATE player_rankings 
SET combat_power = total_points * 100 
WHERE combat_power = 0;

-- Drop old index on total_points
DROP INDEX IF EXISTS idx_player_rankings_total_points;

-- Create optimized index for combat_power (descending order)
CREATE INDEX IF NOT EXISTS idx_combat_power_desc ON player_rankings(combat_power DESC, id ASC);

-- Rename existing player_id index for clarity
DROP INDEX IF EXISTS idx_player_rankings_player_id;
CREATE INDEX IF NOT EXISTS idx_player_id ON player_rankings(player_id);

-- Create index for updated_at (for sync operations)
CREATE INDEX IF NOT EXISTS idx_updated_at ON player_rankings(updated_at);

-- Create partial index for top 10K players (optimization)
-- This will be updated periodically based on the threshold
CREATE INDEX IF NOT EXISTS idx_top_players ON player_rankings(combat_power DESC, id ASC)
WHERE combat_power >= 100000; -- Initial threshold, will be updated

-- Create materialized view for top 10K players
DROP MATERIALIZED VIEW IF EXISTS top_10k_players;
CREATE MATERIALIZED VIEW top_10k_players AS
SELECT 
    player_id,
    username,
    combat_power,
    total_points,
    wins,
    losses,
    win_rate,
    updated_at,
    ROW_NUMBER() OVER (ORDER BY combat_power DESC, id ASC) as rank
FROM player_rankings
ORDER BY combat_power DESC, id ASC
LIMIT 10000;

-- Create unique index on materialized view for concurrent refresh
CREATE UNIQUE INDEX idx_mv_top_10k_rank ON top_10k_players(rank);
CREATE INDEX idx_mv_top_10k_player_id ON top_10k_players(player_id);
CREATE INDEX idx_mv_top_10k_combat_power ON top_10k_players(combat_power DESC);

-- Function to refresh materialized view (can be called by cron or app)
CREATE OR REPLACE FUNCTION refresh_top_10k_leaderboard()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY top_10k_players;
END;
$$ LANGUAGE plpgsql;

-- Optional: Create a trigger to track last update time
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_player_rankings_updated_at ON player_rankings;
CREATE TRIGGER update_player_rankings_updated_at
    BEFORE UPDATE ON player_rankings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON COLUMN player_rankings.combat_power IS 'Primary ranking metric, indexed for fast sorting';
COMMENT ON INDEX idx_combat_power_desc IS 'Main index for leaderboard queries, sorted descending';
COMMENT ON INDEX idx_top_players IS 'Partial index for top 10K players only, reduces index size by 99%+';
COMMENT ON MATERIALIZED VIEW top_10k_players IS 'Pre-computed top 10K leaderboard, refresh periodically';
