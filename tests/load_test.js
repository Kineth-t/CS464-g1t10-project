import http from 'k6/http';
import { check, sleep } from 'k6';

// BASE_URL can be overridden via the -e flag:
//   k6 run -e BASE_URL=https://ringr.up.railway.app/api tests/load_test.js
//
// Defaults to localhost for running against docker-compose locally.
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// ─── Scenario ────────────────────────────────────────────────────────────────
// Uncomment ONE of the option blocks below depending on what you want to test.

// Option 1: Quick smoke test — 1 user, 10 requests
export const options = {
  vus: 1,
  iterations: 10,
};

// Option 2: Staged load test — ramps up to prove rate limiting (429s) kicks in
// export const options = {
//   stages: [
//     { duration: '30s', target: 50 },   // Ramp up to normal traffic
//     { duration: '1m',  target: 500 },  // High traffic — hits global safety throttle
//     { duration: '30s', target: 1000 }, // Stress — expect many 429 Too Many Requests
//     { duration: '30s', target: 0 },    // Ramp down
//   ],
// };

// ─── Helpers ─────────────────────────────────────────────────────────────────

function jsonHeaders() {
  return { headers: { 'Content-Type': 'application/json' } };
}

// ─── Test Scenarios ───────────────────────────────────────────────────────────

// Scenario A: Public phone catalog (no auth required)
function testPhoneListing() {
  const res = http.get(`${BASE_URL}/phones`, jsonHeaders());
  check(res, {
    'GET /phones → 200': (r) => r.status === 200,
    'GET /phones → returns array': (r) => {
      try { return Array.isArray(JSON.parse(r.body)); } catch { return false; }
    },
  });
}

// Scenario B: Login with wrong password → expect 401
// Also proves the per-user login throttle (5 attempts/min) fires a 429 under load
function testLoginThrottle() {
  const payload = JSON.stringify({ username: 'admin', password: 'wrongpassword' });
  const res = http.post(`${BASE_URL}/auth/login`, payload, jsonHeaders());
  check(res, {
    'POST /auth/login wrong password → 401 or 429': (r) =>
      r.status === 401 || r.status === 429,
  });
}

// ─── Main ─────────────────────────────────────────────────────────────────────

export default function () {
  testPhoneListing();
  sleep(0.5);

  testLoginThrottle();
  sleep(1);
}
