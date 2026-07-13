// Cross-cutting helpers: price, date, duration formatters.
// Prices for the GoBooking API are stored as VND — int with currency implied.

const vndFormatter = new Intl.NumberFormat('vi-VN', {
  style: 'currency',
  currency: 'VND',
  maximumFractionDigits: 0,
});

const vndShort = new Intl.NumberFormat('en', {
  notation: 'compact',
  maximumFractionDigits: 1,
});

export function formatVnd(value, { short = false } = {}) {
  const n = Number(value || 0);
  if (Number.isNaN(n)) return '—';
  if (short) {
    const compact = vndShort.format(n);
    return `₫${compact}`;
  }
  return vndFormatter.format(n);
}

export function formatUsd(value) {
  const n = Number(value || 0);
  if (Number.isNaN(n)) return '—';
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(n / 24000); // mock conversion for the hero price view only
}

export function formatDate(iso, options = {}) {
  if (!iso) return '—';
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return '—';
  return d.toLocaleDateString(undefined, {
    weekday: 'short',
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    ...options,
  });
}

export function formatTime(iso) {
  if (!iso) return '—';
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return '—';
  return d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });
}

export function durationLabel(startIso, endIso) {
  if (!startIso || !endIso) return '—';
  const ms = new Date(endIso).getTime() - new Date(startIso).getTime();
  if (Number.isNaN(ms) || ms <= 0) return '—';
  const mins = Math.round(ms / 60000);
  const h = Math.floor(mins / 60);
  const m = mins % 60;
  if (h && m) return `${h}h ${m}m`;
  if (h) return `${h}h`;
  return `${m}m`;
}

export function nightsLabel(checkInIso, checkOutIso) {
  if (!checkInIso || !checkOutIso) return 0;
  const ms = new Date(checkOutIso).getTime() - new Date(checkInIso).getTime();
  return Math.max(0, Math.round(ms / 86400000));
}

export function isoFromDate(d) {
  if (!d) return '';
  const dt = d instanceof Date ? d : new Date(d);
  const yyyy = dt.getFullYear();
  const mm = `${dt.getMonth() + 1}`.padStart(2, '0');
  const dd = `${dt.getDate()}`.padStart(2, '0');
  return `${yyyy}-${mm}-${dd}`;
}
