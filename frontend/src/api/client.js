const BASE_URL = '/api';

function getToken() {
  return localStorage.getItem('token');
}

async function request(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  const token = getToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;

  try {
    const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });
    const text = await res.text();
    let data = null;
    try { data = text ? JSON.parse(text) : null; } catch { data = { error: text }; }

    if (!res.ok) {
      const msg = data?.error || data?.message || text || `Request failed (${res.status})`;
      throw new Error(msg);
    }
    return data;
  } catch (e) {
    console.error(`API Error: ${path}`, e);
    throw e;
  }
}

// Auth
export const authAPI = {
  login: (username, password) =>
    request('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) }),
  register: (payload) =>
    request('/auth/register', { method: 'POST', body: JSON.stringify(payload) }),
};

// Phones
export const phonesAPI = {
  list: () => request('/phones').then((data) => Array.isArray(data) ? data : []),
  get: (id) => request(`/phones/${id}`),
  create: (payload) => request('/phones', { method: 'POST', body: JSON.stringify(payload) }),
  update: (id, payload) => request(`/phones/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
  delete: (id) => request(`/phones/${id}`, { method: 'DELETE' }),
};

// Cart
export const cartAPI = {
  get: () => request('/cart'),
  add: (phone_id, quantity) =>
    request('/cart', { method: 'POST', body: JSON.stringify({ phone_id, quantity }) }),
  remove: (itemId) => request(`/cart/${itemId}`, { method: 'DELETE' }),
  checkout: () => request('/cart/checkout', { method: 'POST' }),
};

// Checkout payment
export const paymentAPI = {
  pay: ({ payment_method_id }) =>
    request('/pay', {
      method: 'POST',
      body: JSON.stringify({
        payment_method_id
      }),
    }),

  getOrders: () =>
    request('/orders'),

  getOrder: (id) =>
    request(`/orders/${id}`),
}