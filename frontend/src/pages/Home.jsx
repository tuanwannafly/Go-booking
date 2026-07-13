import { Link } from 'react-router-dom';
import {
  IconArrowRight,
  IconPlane,
  IconHotel,
  IconIllustrationCity,
  IconIllustrationDriver,
} from '../components/Icons.jsx';
import RequestForm from '../components/RequestForm.jsx';
import { CategoryRow, HeroBand, PromoDark, PromoLight } from '../components/Marketing.jsx';

const savings = [
  {
    eyebrow: 'Why GoBooking',
    title: 'Travel and stays, simplified',
    body: 'One place to find flights and rooms across Vietnam and beyond. No surprise fees at checkout, no second brand colour in your inbox.',
  },
  {
    eyebrow: 'How it works',
    title: 'See the price, hold the seat',
    body: 'Tap See prices to compare fares, hold the seat you like for ten minutes, and confirm only when you are happy with the total.',
  },
];

export default function Home() {
  return (
    <div>
      <HeroBand>
        <div>
          <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
            Travel · Stays · Anywhere
          </div>
          <h1 className="display-xxl max-w-[12ch]">
            Go anywhere, with one black-and-white promise.
          </h1>
          <p className="mt-5 body-lg text-body max-w-prose">
            Search flights and hotels side-by-side. Hold a seat while you decide.
            Confirm, cancel, or reschedule without talking to a call centre.
          </p>
          <div className="mt-7 flex flex-wrap gap-3">
            <Link to="/flights" className="btn-primary">
              Search flights <IconArrowRight className="h-4 w-4" />
            </Link>
            <Link to="/hotels" className="btn-secondary">
              Find stays <IconArrowRight className="h-4 w-4" />
            </Link>
          </div>

          <div className="mt-10 grid grid-cols-3 gap-4 max-w-md">
            <HeroStat label="Routes" value="12+" />
            <HeroStat label="Cities" value="5" />
            <HeroStat label="Avg hold time" value="10m" />
          </div>
        </div>

        <div className="flex justify-center lg:justify-end">
          <RequestForm />
        </div>
      </HeroBand>

      <CategoryRow />

      <PromoLight
        title="Your flight, on your schedule"
        body="Reschedule any pending booking up to the day of departure. Cancel inside the policy window and we will refund the rest, no phone tag required."
        ctaText="Explore flights"
        ctaTo="/flights"
        illustration={<IconIllustrationCity className="w-full h-auto" />}
      />

      <PromoDark
        title="From the city to the beach, in one tab"
        body="Whether you are flying for work or switching off for the weekend, compare every fare and every room with the same stripped-down interface."
        ctaText="See all stays"
        ctaTo="/hotels"
      />

      <PromoLight
        reverse
        title="Hold first, pay once you are sure"
        body="Tap a flight or a room to reserve the specific seat or unit for ten minutes. No double-bookings, no overbooked flights, no last-minute apologies."
        ctaText="Try a search"
        ctaTo="/flights"
        illustration={<IconIllustrationDriver className="w-full h-auto" />}
      />

      {/* Why section */}
      <section className="container-x py-12 md:py-16">
        <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
          The promise
        </div>
        <h2 className="display-xl max-w-[18ch]">
          One interface. No second brand colour in your inbox.
        </h2>
        <p className="mt-4 body-md text-body max-w-[60ch]">
          We borrow the discipline of great transportation apps: a single black
          CTA, a single pill shape, and a list of things you can do today.
        </p>

        <div className="mt-10 grid gap-5 md:grid-cols-3">
          {savings.map((s, i) => (
            <div key={s.title} className="card-elevated p-6 md:p-8 h-full flex flex-col">
              <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
                {String(i + 1).padStart(2, '0')} · {s.eyebrow}
              </div>
              <h3 className="display-md mb-3">{s.title}</h3>
              <p className="body-md text-body">{s.body}</p>
            </div>
          ))}
        </div>
      </section>

      <PromoDark
        title="Safety, simplified"
        body="Every payment is a mock provider with idempotent responses. Every hold has a hard expiry. Every booking keeps the same status grammar across the entire app."
        ctaText="How it works"
        ctaTo="/about"
      />

      {/* App download pills */}
      <section className="container-x py-12 md:py-16">
        <div className="grid md:grid-cols-2 gap-5">
          <div className="card-soft p-8 md:p-10 flex items-center gap-6">
            <div className="rounded-xl bg-canvas w-16 h-16 flex items-center justify-center">
              <IconPlane className="h-7 w-7" />
            </div>
            <div className="flex-1">
              <div className="caption uppercase tracking-[0.15em] text-mute">For travellers</div>
              <div className="display-md">Download the Traveller app</div>
              <p className="body-md text-body mt-2">Hold, confirm and reschedule on the go.</p>
            </div>
            <button type="button" className="app-pill">Download iOS</button>
          </div>
          <div className="card-soft p-8 md:p-10 flex items-center gap-6">
            <div className="rounded-xl bg-canvas w-16 h-16 flex items-center justify-center">
              <IconHotel className="h-7 w-7" />
            </div>
            <div className="flex-1">
              <div className="caption uppercase tracking-[0.15em] text-mute">For partners</div>
              <div className="display-md">List your inventory</div>
              <p className="body-md text-body mt-2">Plug your inventory API into GoBooking.</p>
            </div>
            <button type="button" className="app-pill">Become a partner</button>
          </div>
        </div>
      </section>

      {/* FAQ */}
      <section className="container-x py-12 md:py-16">
        <h2 className="display-xl mb-8">Frequently asked</h2>
        <Faq
          q="How long does a seat hold last?"
          a="Ten minutes from the moment you tap Hold. We ring a countdown so you know when it expires."
        />
        <Faq
          q="What happens if my hold expires?"
          a="The seat returns to the available pool automatically — no manual cleanup needed."
        />
        <Faq
          q="Can I reschedule my booking later?"
          a="Yes, as long as the booking is still pending and the new departure time is in the future."
        />
        <Faq
          q="How is my payment handled?"
          a="The mock payment provider charges the booking's total amount on confirm and refunds according to the cancellation policy."
        />
      </section>
    </div>
  );
}

function HeroStat({ label, value }) {
  return (
    <div className="flex flex-col">
      <div className="display-md">{value}</div>
      <div className="caption text-body">{label}</div>
    </div>
  );
}

function Faq({ q, a }) {
  return (
    <details className="group border-b border-surfacePressed">
      <summary className="body-md-strong flex items-center justify-between cursor-pointer py-5 list-none">
        <span>{q}</span>
        <span className="ml-4 inline-flex h-6 w-6 items-center justify-center rounded-full bg-canvasSoft group-open:rotate-45 transition-transform">
          +
        </span>
      </summary>
      <div className="pb-5 body-md text-body max-w-prose">{a}</div>
    </details>
  );
}
