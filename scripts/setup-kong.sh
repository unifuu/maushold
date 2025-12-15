echo "🚀 Setting up Kong API Gateway..."

# Wait for Kong to be ready
until curl -s http://localhost:8001/ > /dev/null 2>&1; do
  echo "Waiting for Kong Admin API..."
  sleep 2
done

echo "✅ Kong is ready!"

# Create Player Service
echo "📝 Creating Player Service..."
curl -i -X POST http://localhost:8001/services \
  --data name=player-service \
  --data url=http://player-service:8001

curl -i -X POST http://localhost:8001/services/player-service/routes \
  --data "paths[]=/api/players" \
  --data "strip_path=false"

# Create Monster Service
echo "📝 Creating Monster Service..."
curl -i -X POST http://localhost:8001/services \
  --data name=monster-service \
  --data url=http://monster-service:8002

curl -i -X POST http://localhost:8001/services/monster-service/routes \
  --data "paths[]=/api/monster" \
  --data "strip_path=false"

# Create Battle Service
echo "📝 Creating Battle Service..."
curl -i -X POST http://localhost:8001/services \
  --data name=battle-service \
  --data url=http://battle-service:8003

curl -i -X POST http://localhost:8001/services/battle-service/routes \
  --data "paths[]=/api/battles" \
  --data "strip_path=false"

# Create Ranking Service
echo "📝 Creating Ranking Service..."
curl -i -X POST http://localhost:8001/services \
  --data name=ranking-service \
  --data url=http://ranking-service:8004

curl -i -X POST http://localhost:8001/services/ranking-service/routes \
  --data "paths[]=/api/rankings" \
  --data "strip_path=false"

# Add Rate Limiting Plugin (Optional)
echo "🔌 Adding Rate Limiting..."
curl -i -X POST http://localhost:8001/plugins \
  --data "name=rate-limiting" \
  --data "config.minute=100" \
  --data "config.policy=local"

# Add CORS Plugin
echo "🔌 Adding CORS..."
curl -i -X POST http://localhost:8001/plugins \
  --data "name=cors" \
  --data "config.origins=*" \
  --data "config.methods=GET,POST,PUT,DELETE,OPTIONS" \
  --data "config.headers=Content-Type,Authorization"

# Add Request/Response Logging (Optional)
echo "🔌 Adding Logging..."
curl -i -X POST http://localhost:8001/plugins \
  --data "name=file-log" \
  --data "config.path=/tmp/kong.log"

echo "✅ Kong setup complete!"
echo ""
echo "📊 Access points:"
echo "  - API Gateway:  http://localhost:8000"
echo "  - Kong Admin:   http://localhost:8001"
echo "  - Konga UI:     http://localhost:1337"
echo ""
echo "🔍 Test endpoints:"
echo "  curl http://localhost:8000/api/players"
echo "  curl http://localhost:8000/api/monster"
echo "  curl http://localhost:8000/api/battles"
echo "  curl http://localhost:8000/api/rankings"

# ==================================
# frontend/src/App.tsx (UPDATED for API Gateway)
# ==================================
// Change API endpoints to use Kong Gateway
const API = {
  GATEWAY: 'http://localhost:8000', // Kong Gateway
  PLAYER: 'http://localhost:8000/api/players',
  MONSTER: 'http://localhost:8000/api/monster',
  BATTLE: 'http://localhost:8000/api/battles',
  RANKING: 'http://localhost:8000/api/rankings'
};

// Or use direct access (bypass gateway)
const API_DIRECT = {
  PLAYER: 'http://localhost:8001',
  MONSTER: 'http://localhost:8002',
  BATTLE: 'http://localhost:8003',
  RANKING: 'http://localhost:8004'
};