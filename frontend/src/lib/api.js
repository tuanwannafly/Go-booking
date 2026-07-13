// Lightweight fetch wrapper that knows about the GoBooking API surface.
// All calls go through here so we have one place to inject auth, idempotency
// keys, and to translate non-2xx responses into structured errors.

const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1';

function getToken() {
  return localStorage.getItem('gobooking.token');
}

async function request(path, { method = 'GET', body, headers = {}, query, signal, raw } = {}) {
  const url = new URL(
    path.startsWith('http') ? path : `${API_BASE}${path}`,
    window.location.origin,
  );
  if (query) {
    Object.entries(query).forEach(([k, v]) => {
      if (v == null || v === '') return;
      url.searchParams.append(k, String(v));
    });
  }

  const token = getToken();
  const finalHeaders = {
    Accept: 'application/json',
    ...headers,
  };
  if (body && !(body instanceof FormData)) {
    finalHeaders['Content-Type'] = 'application/json';
  }
  if (token) finalHeaders.Authorization = `Bearer ${token}`;

  let payload;
  if (body == null) {
    payload = undefined;
  } else if (body instanceof FormData) {
    payload = body;
  } else {
    payload = JSON.stringify(body);
  }

  let res;
  try {
    res = await fetch(raw ? url.toString() : url.toString().replace(window.location.origin, ''), {
      method,
      headers: finalHeaders,
      body: payload,
      signal,
    });
  } catch (err) {
    if (err.name === 'AbortError') throw err;
    throw new ApiError('network', 0, 'Cannot reach the booking service');
  }

  const isJson = res.headers.get('content-type')?.includes('application/json');
  const data = isJson ? await res.json() : await res.text();

  if (!res.ok) {
    const msg = (isJson && data && (data.error || data.message)) || res.statusText;
    throw new ApiError(msg, res.status, data);
  }
  return data;
}

export class ApiError extends Error {
  constructor(message, status, details) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.details = details;
  }
}

// Inject the API base path so consumers can call fetch directly through window when needed
// (e.g., Stripe-style page navigation). Keeps URL construction in one place.
export const apiBaseUrl = API_BASE;

// Public search endpoints ------------------------------------------------------
export const flights = {
  search: (params) =>
    request('/flights/search', {
      query: {
        origin: params.origin?.toUpperCase(),
        destination: params.destination?.toUpperCase(),
        date: params.date,
        page: params.page || 1,
        page_size: params.pageSize || 10,
      },
    }),
  get: (id) => request(`/flights/${id}`),
};

export const hotels = {
  search: (params) =>
    request('/hotels/search', {
      query: {
        city: params.city,
        checkin: params.checkin,
        checkout: params.checkout,
        page: params.page || 1,
        page_size: params.pageSize || 10,
      },
    }),
  get: (id, params = {}) =>
    request(`/hotels/${id}`, {
      query: { checkin: params.checkin, checkout: params.checkout },
    }),
};

// Authenticated operations ----------------------------------------------------
function ensureAuth() {
  if (!getToken()) {
    throw new ApiError('Please sign in to continue', 401);
  }
}

function idempotencyKey() {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) return crypto.randomUUID();
  return `idem-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

export const holds = {
  seat: (flightId, seatId, { duration = 10 } = {}) => {
    ensureAuth();
    return request(`/flights/${flightId}/seats/${seatId}/hold`, {
      method: 'POST',
      body: { hold_duration: duration },
    });
  },
  room: (roomId, { duration = 10 } = {}) => {
    ensureAuth();
    return request(`/hotels/rooms/${roomId}/hold`, {
      method: 'POST',
      body: { hold_duration: duration },
    });
  },
  release: (unitId) => {
    ensureAuth();
    return request(`/holds/${unitId}`, { method: 'DELETE' });
  },
};

export const bookings = {
  create: (items, { scheduledAt } = {}) => {
    ensureAuth();
    return request('/bookings', {
      method: 'POST',
      body: { items, scheduled_at: scheduledAt },
      headers: { 'Idempotency-Key': idempotencyKey() },
    });
  },
  list: ({ page = 1, pageSize = 10 } = {}) => {
    ensureAuth();
    return request('/bookings', { query: { page, page_size: pageSize } });
  },
  get: (id) => {
    ensureAuth();
    return request(`/bookings/${id}`);
  },
  confirm: (id, paymentMethod) => {
    ensureAuth();
    return request(`/bookings/${id}/confirm`, {
      method: 'POST',
      body: { payment_method: paymentMethod },
    });
  },
  cancel: (id, reason) => {
    ensureAuth();
    return request(`/bookings/${id}/cancel`, {
      method: 'POST',
      body: { reason },
    });
  },
  schedule: (id, scheduledAt) => {
    ensureAuth();
    return request(`/bookings/${id}/schedule`, {
      method: 'POST',
      body: { scheduled_at: scheduledAt },
    });
  },
};

// Health checks used by the splash screen. Live on /api/v1's parent (server
// health) — we keep the convention simple by exposing them from the same base.
export const health = {
  ready: () => request('/readyz', { raw: true }),
};

export default { flights, hotels, holds, bookings, health };
