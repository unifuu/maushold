#!/bin/bash

# Kong Admin API URL
KONG_ADMIN_URL="http://localhost:18001"

echo "Cleaning up Kong routes and services..."

# Delete all routes
echo "Deleting all routes..."
curl -s $KONG_ADMIN_URL/routes | jq -r '.data[].id' | while read route_id; do
  echo "Deleting route: $route_id"
  curl -X DELETE $KONG_ADMIN_URL/routes/$route_id
done

# Delete all services
echo "Deleting all services..."
curl -s $KONG_ADMIN_URL/services | jq -r '.data[].id' | while read service_id; do
  echo "Deleting service: $service_id"
  curl -X DELETE $KONG_ADMIN_URL/services/$service_id
done

echo "Cleanup complete!"