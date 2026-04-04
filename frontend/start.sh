#!/bin/sh
# Resolves BACKEND_URL to an IPv4 address before starting nginx, and re-checks
# every 30 s so a backend redeploy (new IP) is picked up without restarting the
# frontend container.

set -e

# Split "host:port" → separate variables
BACKEND_HOST="${BACKEND_URL%%:*}"
BACKEND_PORT="${BACKEND_URL#*:}"
# If there was no colon, port var equals the whole URL — clear it
[ "$BACKEND_PORT" = "$BACKEND_URL" ] && BACKEND_PORT=""

resolve_ipv4() {
    # getent hosts on Alpine/musl prefers A records; grep -v ':' drops any IPv6 lines
    getent hosts "$BACKEND_HOST" 2>/dev/null | grep -v ':' | awk 'NR==1{print $1}'
}

CURRENT_IP=$(resolve_ipv4)
if [ -z "$CURRENT_IP" ]; then
    echo "[start] WARNING: could not resolve $BACKEND_HOST to IPv4, falling back to hostname"
    CURRENT_IP="$BACKEND_HOST"
else
    echo "[start] Resolved $BACKEND_HOST -> $CURRENT_IP"
fi

if [ -n "$BACKEND_PORT" ]; then
    RESOLVED_BACKEND="${CURRENT_IP}:${BACKEND_PORT}"
else
    RESOLVED_BACKEND="$CURRENT_IP"
fi
export RESOLVED_BACKEND

envsubst '${RESOLVED_BACKEND}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf
echo "[start] nginx config written with backend: $RESOLVED_BACKEND"

# Background loop: re-resolve every 30 s; reload nginx if the IP changes
(
    LOOP_IP="$CURRENT_IP"
    while true; do
        sleep 30
        NEW_IP=$(resolve_ipv4)
        if [ -n "$NEW_IP" ] && [ "$NEW_IP" != "$LOOP_IP" ]; then
            echo "[start] Backend IP changed $LOOP_IP -> $NEW_IP, reloading nginx"
            LOOP_IP="$NEW_IP"
            if [ -n "$BACKEND_PORT" ]; then
                RESOLVED_BACKEND="${LOOP_IP}:${BACKEND_PORT}"
            else
                RESOLVED_BACKEND="$LOOP_IP"
            fi
            export RESOLVED_BACKEND
            envsubst '${RESOLVED_BACKEND}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf
            nginx -s reload
        fi
    done
) &

exec nginx -g 'daemon off;'
