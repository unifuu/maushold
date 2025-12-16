#!/bin/bash

# Kong Admin API URL
KONG_ADMIN_URL="http://localhost:18001"

echo "Setting up Kong API Gateway routes (simple approach)..."

# Wait for Kong to be ready
bash scripts/wait-for-kong.sh

echo "Kong is ready! Setting up services and routes..."

# Create Player Service
echo "Creating Player Service..."
PLAYER_SERVICE_ID=$(curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=player-service" \
  --data "url=http://player-service:8001" | jq -r '.id')

# Create Player Service Route
echo "Creating Player Service route..."
PLAYER_ROUTE_ID=$(curl -s -X POST $KONG_ADMIN_URL/services/player-service/routes \
  --data "paths[]=/api/players" \
  --data "methods[]=GET" \
  --data "methods[]=POST" \
  --data "methods[]=PUT" \
  --data "methods[]=DELETE" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=true" | jq -r '.id')

# Add request transformer to rewrite path
echo "Adding request transformer for players..."
curl -s -X POST $KONG_ADMIN_URL/routes/$PLAYER_ROUTE_ID/plugins \
  --data "name=request-transformer" \
  --data "config.replace.uri=/players" > /dev/null

# Create Monster Service
echo "Creating Monster Service..."
MONSTER_SERVICE_ID=$(curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=monster-service" \
  --data "url=http://monster-service:8002" | jq -r '.id')

# Create Monster Service Route
echo "Creating Monster Service route..."
MONSTER_ROUTE_ID=$(curl -s -X POST $KONG_ADMIN_URL/services/monster-service/routes \
  --data "paths[]=/api/monster" \
  --data "methods[]=GET" \
  --data "methods[]=POST" \
  --data "methods[]=PUT" \
  --data "methods[]=DELETE" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=true" | jq -r '.id')

# Add request transformer to rewrite path
echo "Adding request transformer for monsters..."
curl -s -X POST $KONG_ADMIN_URL/routes/$MONSTER_ROUTE_ID/plugins \
  --data "name=request-transformer" \
  --data "config.replace.uri=/monster" > /dev/null

# Enable CORS plugin globally
echo "Enabling CORS plugin..."
curl -s -X POST $KONG_ADMIN_URL/plugins/ \
  --data "name=cors" \
  --data "config.origins=*" \
  --data "config.methods=GET" \
  --data "config.methods=POST" \
  --data "config.methods=PUT" \
  --data "config.methods=DELETE" \
  --data "config.methods=OPTIONS" \
  --data "config.headers=Accept,Accept-Version,Content-Length,Content-MD5,Content-Type,Date,X-Auth-Token,Authorization" \
  --data "config.exposed_headers=X-Auth-Token" \
  --data "config.credentials=true" \
  --data "config.max_age=3600" > /dev/null

echo "Kong setup complete!"
echo "You can now access your APIs through:"
echo "  - Players: http://localhost:8000/api/players"
echo "  - Monsters: http://localhost:8000/api/monster"

echo ""
echo "Testing the setup..."
echo "Players endpoint:"
curl -s http://localhost:8000/api/players | head -c 200
echo ""
echo "Monsters endpoint:"
curl -s http://localhost:8000/api/monster | head -c 200
echo ""