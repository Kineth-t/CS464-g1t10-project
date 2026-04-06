#!/bin/sh
# Write nginx config and start nginx. A background loop reloads nginx every 30s
# so that stale backend IPs (from Railway private network redeployments) are
# picked up automatically — nginx re-resolves proxy_pass hostnames on reload.

set -e

envsubst '${BACKEND_URL}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf
echo "[start] nginx config written (backend: ${BACKEND_URL})"

# Background: reload nginx every 30s to re-resolve backend.railway.internal
( while true; do sleep 30; nginx -s reload 2>/dev/null; done ) &

exec nginx -g 'daemon off;'
