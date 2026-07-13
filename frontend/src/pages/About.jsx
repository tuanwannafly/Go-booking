import { Link } from 'react-router-dom';
import { IconArrowRight } from '../components/Icons.jsx';
import { PromoDark, PromoLight } from '../components/Marketing.jsx';

const stats = [
  { label: 'Routes live', value: '12+' },
  { label: 'Hotels seeded', value: '5' },
  { label: 'Hold window', value: '10m' },
  { label: 'Concurrency tested', value: '100 goroutines' },
];

const tenets = [
  {
    n: '01',
    title: 'Two-layer locking',
    body: 'Flights use pessimistic locking (SELECT FOR UPDATE + Redis distributed lock). Rooms use optimistic locking (version column, retry on conflict). Each resource type gets the strategy that matches its contention pattern.',
  },
  {
    n: '02',
    title: 'Idempotent by default',
    body: 'Every booking creation reads the Idempotency-Key header, persists the response, and replays it on retry. The same request can fire twice without producing two bookings.',
  },
  {
    n: '03',
    title: 'Bounded time',
    body: 'Holds expire in ten minutes. Bookings that are never confirmed expire on a worker tick. Inventory always returns to the available pool.',
  },
  {
    n: '04',
    title: 'Refund policy as code',
    body: 'Cancel more than 24h before for a full refund, 6-24h for 50%, under 6h none. The policy is a single struct that the service layer reads from.',
  },
];

export default function About() {
  return (
    <div className="bg-canvas">
      <section className="container-x py-8 md:py-16">
        <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
          About GoBooking
        </div>
        <h1 className="display-xxl max-w-[18ch]">
          A travel super-app built like a transportation app.
        </h1>
        <p className="mt-5 body-lg text-body max-w-[60ch]">
          GoBooking is a portfolio project that models the real-world trade-offs
          of booking systems: how to keep two travellers from buying the same
          seat, how to charge a card exactly once even when the network drops the
          response, and how to reschedule without opening a ticket.
        </p>

        <div className="mt-10 grid grid-cols-2 md:grid-cols-4 gap-4">
          {stats.map((s) => (
            <div key={s.label} className="card-soft p-5">
              <div className="display-md tabular-nums">{s.value}</div>
              <div className="caption text-mute">{s.label}</div>
            </div>
          ))}
        </div>
      </section>

      <PromoLight
        title="What you can do today"
        body="Search flights by origin, destination and date. Hold a specific seat for ten minutes. Confirm a booking, reschedule a departure, or cancel inside the policy window."
        ctaText="Start searching"
        ctaTo="/flights"
        illustration={<IconIllustrationRider className="w-full h-auto" />}
      />

      <PromoDark
        title="What you should not do yet"
        body="Real payment gateway, OAuth, and an admin dashboard. Those are scheduled for Sprint 5, after the concurrency and idempotency stories have shipped."
        ctaText="See the API"
        ctaTo="/about"
      />

      <PromoLight
        reverse
        title="How the frontend serves it"
        body="A React SPA that calls the same JSON the rest of your stack uses. Pill-shaped CTAs, sentence-case displays, no second brand colour. Black-and-white, by design."
        ctaText="Read the design"
        ctaTo="/about"
        illustration={<IconIllustrationDriver className="w-full h-auto" />}
      />

      <section className="container-x py-12 md:py-16">
        <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
          The tenets
        </div>
        <h2 className="display-xl max-w-[18ch]">
          Four ideas that shape every endpoint.
        </h2>

        <div className="mt-10 grid md:grid-cols-2 gap-5">
          {tenets.map((t) => (
            <div key={t.n} className="card-elevated p-6 md:p-8">
              <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
                {t.n}
              </div>
              <h3 className="display-md">{t.title}</h3>
              <p className="mt-3 body-md text-body">{t.body}</p>
            </div>
          ))}
        </div>
      </section>

      <PromoDark
        title="Want to run it locally?"
        body="Two terminals: docker compose up for the database and Redis, then go run ./cmd/api for the service, then npm run dev for the frontend. You will have a working end-to-end flow in under five minutes."
        ctaText="Back to home"
        ctaTo="/"
      />

      <section className="container-x py-12 md:py-16">
        <h2 className="display-xl mb-8">A peek at the API surface</h2>
        <div className="grid md:grid-cols-2 gap-5">
          <div className="card p-6 md:p-8">
            <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
              Public · Search
            </div>
            <ul className="grid gap-3 body-md">
              <li>
                <code className="body-sm-strong">GET /api/v1/flights/search</code> · origin, destination, date
              </li>
              <li>
                <code className="body-sm-strong">GET /api/v1/hotels/search</code> · city, checkin, checkout
              </li>
              <li>
                <code className="body-sm-strong">GET /api/v1/flights/:id</code>
              </li>
              <li>
                <code className="body-sm-strong">GET /api/v1/hotels/:id</code>
              </li>
            </ul>
          </div>
          <div className="card p-6 md:p-8">
            <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
              Protected · Bookings
            </div>
            <ul className="grid gap-3 body-md">
              <li>
                <code className="body-sm-strong">POST /api/v1/bookings</code> · Idempotency-Key
              </li>
              <li>
                <code className="body-sm-strong">POST /api/v1/bookings/:id/confirm</code>
              </li>
              <li>
                <code className="body-sm-strong">POST /api/v1/bookings/:id/cancel</code>
              </li>
              <li>
                <code className="body-sm-strong">POST /api/v1/bookings/:id/schedule</code>
              </li>
            </ul>
          </div>
        </div>

        <div className="mt-10 flex flex-col md:flex-row gap-3">
          <Link to="/" className="btn-primary">
            Back to home <IconArrowRight className="h-4 w-4" />
          </Link>
          <Link to="/flights" className="btn-secondary">
            Search flights
          </Link>
        </div>
      </section>
    </div>
  );
}
