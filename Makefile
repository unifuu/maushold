# Maushold Management Makefile

# Load environment variables from .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: help clear sync refresh

help:
	@echo "Available commands:"
	@echo "  make clear    - Truncate all database tables and flush Redis"
	@echo "  make seed     - Populate monster database with seed data"
	@echo "  make sync     - Trigger initial sync for ranking service"
	@echo "  make refresh  - Refresh ranking materialized view"

clear:
	@echo "🧹 Clearing all database data..."
	@docker exec maushold-player-db-1 psql -U $(DB_USER) -d $(PLAYER_DB_NAME) -c "TRUNCATE TABLE players, player_monsters CASCADE;"
	@docker exec maushold-monster-db-1 psql -U $(DB_USER) -d $(MONSTER_DB_NAME) -c "TRUNCATE TABLE monsters CASCADE;"
	@docker exec maushold-battle-db-1 psql -U $(DB_USER) -d $(BATTLE_DB_NAME) -c "TRUNCATE TABLE battles CASCADE;"
	@docker exec maushold-ranking-db-1 psql -U $(DB_USER) -d $(RANKING_DB_NAME) -c "TRUNCATE TABLE player_rankings, leaderboard_entries, players CASCADE;"
	@echo "🧹 Flushing Redis cache..."
	@docker exec maushold-redis-1 redis-cli -a $(REDIS_PASSWORD) FLUSHALL
	@echo "✨ Refreshing materialized views..."
	@docker exec maushold-ranking-db-1 psql -U $(DB_USER) -d $(RANKING_DB_NAME) -c "REFRESH MATERIALIZED VIEW top_10k_players;"
	@echo "✅ All data cleared successfully!"

seed:
	@echo "🌱 Seeding monsters database..."
	@cat scripts/seed_monsters.sql | docker exec -i maushold-monster-db-1 psql -U $(DB_USER) -d $(MONSTER_DB_NAME)
	@echo "🧹 Invalidating monster cache..."
	@docker exec maushold-redis-1 redis-cli -a $(REDIS_PASSWORD) DEL monster:all
	@echo "✅ Monsters seeded successfully!"

sync:
	@echo "🔄 Triggering ranking service sync..."
	@curl -X POST http://localhost:8004/api/rankings/sync
	@echo "\n✅ Sync triggered."

refresh:
	@echo "🔄 Refreshing ranking materialized view..."
	@curl -X POST http://localhost:8004/api/rankings/refresh
	@echo "\n✅ Refresh triggered."