#!/bin/bash

KONG_ADMIN_URL="http://localhost:18001"
MAX_ATTEMPTS=30
ATTEMPT=1

echo "Waiting for Kong Admin API to be ready..."

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
    if curl -s "$KONG_ADMIN_URL" > /dev/null 2>&1; then
        echo "Kong is ready!"
        exit 0
    fi
    
    echo "Attempt $ATTEMPT/$MAX_ATTEMPTS: Kong not ready yet, waiting..."
    sleep 2
    ATTEMPT=$((ATTEMPT + 1))
done

echo "Kong failed to start within expected time"
exit 1