// Location suggestion dataset for the autocomplete in RequestForm.
// Reflects the routes and cities seeded by migration 000011 so the form stays
// consistent with what's actually in the database.

export const AIRPORTS = [
  { code: 'SGN', name: 'Sài Gòn (TP. Hồ Chí Minh)', city: 'Ho Chi Minh' },
  { code: 'HAN', name: 'Hà Nội (Nội Bài)', city: 'Hanoi' },
  { code: 'DAD', name: 'Đà Nẵng', city: 'Da Nang' },
  { code: 'PQC', name: 'Phú Quốc', city: 'Phu Quoc' },
  { code: 'CXR', name: 'Cam Ranh (Nha Trang)', city: 'Nha Trang' },
  { code: 'VII', name: 'Vinh', city: 'Vinh' },
];

export const CITIES = [
  { code: 'Ho Chi Minh', name: 'TP. Hồ Chí Minh', city: 'Ho Chi Minh' },
  { code: 'Hanoi', name: 'Hà Nội', city: 'Hanoi' },
  { code: 'Da Nang', name: 'Đà Nẵng', city: 'Da Nang' },
  { code: 'Phu Quoc', name: 'Phú Quốc', city: 'Phu Quoc' },
  { code: 'Nha Trang', name: 'Nha Trang', city: 'Nha Trang' },
  { code: 'Hue', name: 'Huế', city: 'Hue' },
];

// Routes that have flights in the seeded data — returned as a sanity hint
// when the autocomplete has no exact match.
export const SEED_ROUTES = [
  { origin: 'SGN', destination: 'HAN', date: '2026-07-20' },
  { origin: 'HAN', destination: 'SGN', date: '2026-07-20' },
  { origin: 'SGN', destination: 'DAD', date: '2026-07-21' },
  { origin: 'DAD', destination: 'SGN', date: '2026-07-21' },
  { origin: 'HAN', destination: 'DAD', date: '2026-07-22' },
  { origin: 'SGN', destination: 'PQC', date: '2026-07-23' },
];

// Resolve a typed query to one of the dataset entries by code, name, or city.
// Returns null if nothing matches. Used by the autocomplete to highlight which
// entries apply to the dataset.
export function findAirport(query) {
  if (!query) return null;
  const q = query.trim().toLowerCase();
  return (
    AIRPORTS.find(
      (a) =>
        a.code.toLowerCase() === q ||
        a.name.toLowerCase().includes(q) ||
        a.city.toLowerCase().includes(q),
    ) || null
  );
}

export function findCity(query) {
  if (!query) return null;
  const q = query.trim().toLowerCase();
  return (
    CITIES.find(
      (c) =>
        c.code.toLowerCase() === q ||
        c.name.toLowerCase().includes(q) ||
        c.city.toLowerCase().includes(q),
    ) || null
  );
}
