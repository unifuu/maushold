echo "🚀 Setting up Kong API Gateway..."
echo ""

# Wait for Kong Admin API to be ready
echo "⏳ Waiting for Kong Admin API..."
until curl -s http://localhost:18001/ > /dev/null 2>&1; do
  echo "   Waiting for Kong to be ready..."
  sleep 2
done

echo "✅ Kong is ready!"
echo ""

# Function to create service and route
create_service_and_route() {
  SERVICE_NAME=$1
  SERVICE_URL=$2
  ROUTE_PATH=$3
  
  echo "📝 Creating $SERVICE_NAME..."
  
  # Create service
  curl -s -X POST http://localhost:18001/services \
    --data "name=$SERVICE_NAME" \
    --data "url=$SERVICE_URL" > /dev/null
  
  # Create route
  curl -s -X POST http://localhost:18001/services/$SERVICE_NAME/routes \
    --data "paths[]=$ROUTE_PATH" \
    --data "strip_path=true" > /dev/null
  
  echo "   ✅ $SERVICE_NAME created at $ROUTE_PATH"
}

# Create Player Service
create_service_and_route "player-service" "http://player-service:8001" "/api/players"

# Create Monster Service  
create_service_and_route "monster-service" "http://monster-service:8002" "/api/monster"

# Create Battle Service
create_service_and_route "battle-service" "http://battle-service:8003" "/api/battles"

# Create Ranking Service
create_service_and_route "ranking-service" "http://ranking-service:8004" "/api/rankings"

echo ""
echo "🔌 Adding Plugins..."

# Add CORS Plugin
curl -s -X POST http://localhost:18001/plugins \
  --data "name=cors" \
  --data "config.origins=http://localhost:3000" \
  --data "config.methods=GET,POST,PUT,DELETE,OPTIONS,PATCH" \
  --data "config.headers=Accept,Accept-Version,Content-Length,Content-MD5,Content-Type,Date,Authorization,X-Requested-With" \
  --data "config.exposed_headers=X-Auth-Token" \
  --data "config.credentials=true" \
  --data "config.max_age=3600" \
  --data "config.preflight_continue=false" > /dev/null

echo "   ✅ CORS plugin enabled"

# Add Rate Limiting (optional)
curl -s -X POST http://localhost:18001/plugins \
  --data "name=rate-limiting" \
  --data "config.minute=1000" \
  --data "config.policy=local" > /dev/null

echo "   ✅ Rate limiting enabled (1000 req/min)"

echo ""
echo "✅ Kong setup complete!"
echo ""
echo "📊 Endpoints available:"
echo "   http://localhost:8000/api/players"
echo "   http://localhost:8000/api/monster"
echo "   http://localhost:8000/api/battles"
echo "   http://localhost:8000/api/rankings"
echo ""
echo "🧪 Test with:"
echo "   curl http://localhost:8000/api/players/health"
echo "   curl http://localhost:8000/api/monster/health"
echo ""