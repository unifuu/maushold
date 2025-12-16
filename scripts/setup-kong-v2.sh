#!/bin/bash

# Kong Admin API URL
KONG_ADMIN_URL="http://localhost:18001"

echo "Setting up Kong API Gateway routes (v2)..."

# Wait for Kong to be ready
bash scripts/wait-for-kong.sh

echo "Kong is ready! Setting up services and routes..."

# Create Player Service
echo "Creating Player Service..."
curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=player-service" \
  --data "url=http://player-service:8001"

# Create specific routes for player endpoints
echo "Creating Player Service routes..."

# Route for /api/players -> /players
curl -s -X POST $KONG_ADMIN_URL/services/player-service/routes \
  --data "name=players-list" \
  --data "paths[]=/api/players" \
  --data "methods[]=GET" \
  --data "methods[]=POST" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false" \
  --data "regex_priority=100"

# Route for /api/players/{id} -> /players/{id}
curl -s -X POST $KONG_ADMIN_URL/services/player-service/routes \
  --data "name=players-detail" \
  --data 'paths[]=/api/players/(?<id>\d+)' \
  --data "methods[]=GET" \
  --data "methods[]=PUT" \
  --data "methods[]=DELETE" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false" \
  --data "regex_priority=200"

# Route for /api/players/{id}/monster -> /players/{id}/monster
curl -s -X POST $KONG_ADMIN_URL/services/player-service/routes \
  --data "name=players-monster" \
  --data 'paths[]=/api/players/(?<id>\d+)/monster' \
  --data "methods[]=GET" \
  --data "methods[]=POST" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false" \
  --data "regex_priority=300"

# Create Monster Service
echo "Creating Monster Service..."
curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=monster-service" \
  --data "url=http://monster-service:8002"

# Create Monster Service Routes
echo "Creating Monster Service routes..."
curl -s -X POST $KONG_ADMIN_URL/services/monster-service/routes \
  --data "name=monster-list" \
  --data "paths[]=/api/monster" \
  --data "methods[]=GET" \
  --data "methods[]=POST" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false"

curl -s -X POST $KONG_ADMIN_URL/services/monster-service/routes \
  --data "name=monster-detail" \
  --data 'paths[]=/api/monster/(?<id>\d+)' \
  --data "methods[]=GET" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false"

curl -s -X POST $KONG_ADMIN_URL/services/monster-service/routes \
  --data "name=monster-random" \
  --data "paths[]=/api/monster/random" \
  --data "methods[]=GET" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false"

# Create Battle Service
echo "Creating Battle Service..."
curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=battle-service" \
  --data "url=http://battle-service:8003"

# Create Battle Service Routes
echo "Creating Battle Service routes..."
curl -s -X POST $KONG_ADMIN_URL/services/battle-service/routes \
  --data "name=battles" \
  --data "paths[]=/api/battles" \
  --data "methods[]=GET" \
  --data "methods[]=POST" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false"

# Create Ranking Service
echo "Creating Ranking Service..."
curl -s -X POST $KONG_ADMIN_URL/services/ \
  --data "name=ranking-service" \
  --data "url=http://ranking-service:8004"

# Create Ranking Service Routes
echo "Creating Ranking Service routes..."
curl -s -X POST $KONG_ADMIN_URL/services/ranking-service/routes \
  --data "name=rankings" \
  --data "paths[]=/api/rankings" \
  --data "methods[]=GET" \
  --data "methods[]=POST" \
  --data "methods[]=OPTIONS" \
  --data "strip_path=false"

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
echo "  - Battles: http://localhost:8000/api/battles"
echo "  - Rankings: http://localhost:8000/api/rankings"

echo ""
echo "Testing the setup..."
echo "Players endpoint:"
curl -s http://localhost:8000/api/players | head -c 200
echo ""
echo "Monsters endpoint:"
curl -s http://localhost:8000/api/monster | head -c 200
echo ""