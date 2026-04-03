import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 1, // Just one person
  iterations: 10, // Just ten clicks

  // stages: [
  //   { duration: '30s', target: 50 },  // Normal traffic (User Throttle range)
  //   { duration: '1m', target: 500 },  // High traffic (Global Safety range)
  //   { duration: '30s', target: 1000 }, // Stress test (Should see many 429s)
  //   { duration: '30s', target: 0 },    // Ramp down
  // ],
};

export default function () {
  const url = 'http://Ringr-Mobile-Backend:8080/auth/login'; // Update to your API URL
  const payload = JSON.stringify({ username: 'admin', password: 'wrongpassword' });
  const params = { headers: { 'Content-Type': 'application/json' } };
  const res = http.post(url, payload, params);

  // const url = 'http://Ringr-Mobile-Backend:8080/phones';
  // const res = http.post(url, payload, params);

  check(res, {
    'is status 200': (r) => r.status === 200,
    'is status 401': (r) => r.status === 401,
    'is status 429': (r) => r.status === 429, // This proves throttling works!
  });

  sleep(1); // Simulate a human waiting 1s between clicks
}