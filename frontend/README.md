# GoBooking — Frontend

A React + TailwindCSS single-page app for the GoBooking backend.

The design language follows `DESIGN-uber.md`: black-and-white duet, sentence-case
displays, signature `999px` pill CTAs, `16px` cards, ride-request-style form card,
mid-page black promo bands, deep-black footer.

## Stack
- React 18 (JavaScript, no TypeScript to keep the entry low-friction)
- React Router 6
- Vite 5
- TailwindCSS 3 (custom design tokens mapped to the Uber spec)
- Native `fetch` for the API client — no extra runtime dependencies

## Layout

```
frontend/
├── src/
│   ├── App.jsx                     # routes
│   ├── main.jsx                    # bootstraps React + providers
│   ├── index.css                   # Tailwind layers + design-token components
│   ├── lib/
│   │   ├── api.js                  # fetch wrapper, search/hold/booking endpoints
│   │   ├── auth.jsx                # mock auth context (mirrors backend stub)
│   │   └── format.js               # VND price, date, duration, nights helpers
│   ├── components/
│   │   ├── Icons.jsx               # SVG icon kit used across the app
│   │   ├── Navbar.jsx              # sticky nav, hamburger, sign-in pills
│   │   ├── Footer.jsx              # black band with company columns
│   │   ├── Primitives.jsx          # IconButton, Modal, EmptyState, StatusPill
│   │   ├── Toast.jsx               # toast provider + hook
│   │   ├── RequestForm.jsx         # ride-request-style booking form card
│   │   └── Marketing.jsx           # PromoLight, PromoDark, CategoryRow
│   └── pages/
│       ├── Home.jsx                # hero + categories + promise + FAQ
│       ├── FlightsResults.jsx      # search list with sort + pagination
│       ├── FlightDetail.jsx        # itinerary + seat map + summary card
│       ├── HotelsResults.jsx       # hotel grid with city illustrations
│       ├── HotelDetail.jsx         # room types + hold + summary card
│       ├── Bookings.jsx            # bookings list (auth required)
│       ├── BookingDetail.jsx       # timeline + confirm/cancel/reschedule dialogs
│       ├── Login.jsx               # mock sign-in / sign-up
│       ├── About.jsx               # project tenets + API surface preview
│       └── NotFound.jsx            # 404
├── index.html
├── tailwind.config.js              # Uber tokens → Tailwind theme
├── postcss.config.js
├── vite.config.js                  # /api → :8080 dev proxy
└── package.json
```

## Backend wiring

The frontend talks to the Go service on `/api/v1/*`. In dev, Vite proxies `/api`
to `http://localhost:8080`. In production build, set `VITE_API_BASE` to your
deployed API prefix.

```bash
# terminal 1 — backend (GoBooking repository root)
docker compose up -d postgres redis
go run ./cmd/api

# terminal 2 — frontend
cd frontend
npm install
npm run dev
```

The pages map to the API surface like so:

| Page              | Endpoint(s)                                                                 |
|-------------------|-----------------------------------------------------------------------------|
| Home              | (no API calls)                                                              |
| Flights search    | `GET /api/v1/flights/search?origin=&destination=&date=`                     |
| Flight detail     | `GET /api/v1/flights/:id` → `POST /api/v1/flights/:id/seats/:seatId/hold`   |
| Stays search      | `GET /api/v1/hotels/search?city=&checkin=&checkout=`                        |
| Hotel detail      | `GET /api/v1/hotels/:id` → `POST /api/v1/hotels/rooms/:roomId/hold`         |
| Booking flow      | `POST /api/v1/bookings` (auto `Idempotency-Key`)                            |
| Bookings list     | `GET /api/v1/bookings`                                                      |
| Booking detail    | `POST /:id/confirm · /:id/cancel · /:id/schedule`                          |

Mock authentication is intentional: any email + name mocks a Bearer token so the
`AuthMiddleware` (`internal/http/middleware`) accepts requests without a full JWT
flow (Sprint 5 deliverable). The token persists in `localStorage` so refreshing
the page keeps you signed in.

## Design tokens

`tailwind.config.js` exposes the spec's colors, radii and shadow levels:

| Token                              | Tailwind class           |
|------------------------------------|--------------------------|
| Ink black                          | `bg-ink`, `text-ink`     |
| Canvas / on-primary                | `bg-canvas`, `text-onPrimary` |
| Body, mute, hairline-mid           | `text-body`, `text-mute` |
| Pill radius (999px)                | `rounded-pill`           |
| Card radius (16px)                 | `rounded-xl`             |
| Pill-tab radius (36px)             | `rounded-pillTab`        |
| Level 2 card drop                  | `shadow-l2`              |
| Display weights 700 Inter          | `display-xxl|xl|lg|md|sm` |
| Body, button, link weights 400/500 | `body-lg|md|md-strong|sm|sm-strong|caption` |

The companion CSS in `src/index.css` adds `.btn-primary`, `.btn-secondary`,
`.btn-subtle`, `.card-elevated`, `.card-soft`, `.request-form-card` and friends
so page code stays declarative and aligned with the spec.

## Run scripts

```bash
npm install      # install
npm run dev      # vite dev server on :5173
npm run build    # production build → dist/
npm run preview  # preview the production build
```

## Demo flow

1. Land on the home page → click **See prices** with origin/destination/city set.
2. The Flights search page lists fares sorted by price. Pick one.
3. On the flight detail, choose a seat and tap **Hold this seat**.
4. Tap **Continue to confirmation** to push the held unit into a booking.
5. Sign in (any email works) when prompted; the bookings tab shows your booking.
6. From the booking detail, confirm / cancel / reschedule — the toasts reflect
   the server response, including an idempotency-replay 200 on retries.
