import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { hotels, ApiError } from '../lib/api.js';
import { IconArrowRight, IconPin } from '../components/Icons.jsx';
import RequestForm from '../components/RequestForm.jsx';
import { EmptyState, Skeleton } from '../components/Primitives.jsx';
import { formatDate, formatVnd, nightsLabel } from '../lib/format.js';

export default function HotelsResults() {
  const [params] = useSearchParams();
  const navigate = useNavigate();

  const city = params.get('city') || '';
  const checkin = params.get('checkin') || '';
  const checkout = params.get('checkout') || '';
  const rooms = Number(params.get('rooms') || '1');

  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!city || !checkin || !checkout) {
      setData(null);
      setError(null);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(null);
    hotels
      .search({ city, checkin, checkout, page: 1, pageSize: 12 })
      .then((res) => !cancelled && setData(res))
      .catch((err) => {
        if (cancelled) return;
        if (err instanceof ApiError) setError(err.message);
        else setError('Could not load stays');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  }, [city, checkin, checkout]);

  const list = data?.data || [];
  const nights = nightsLabel(checkin, checkout);

  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-12">
        <div className="grid lg:grid-cols-2 gap-10 items-start">
          <div>
            <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
              Step 1 of 3 · Pick a stay
            </div>
            <h1 className="display-xl">Stays in {city || 'your destination'}</h1>
            <p className="mt-3 body-md text-body">
              {checkin && checkout
                ? `${formatDate(checkin)} → ${formatDate(checkout)} · ${nights} ${
                    nights === 1 ? 'night' : 'nights'
                  } · ${rooms} ${rooms === 1 ? 'room' : 'rooms'}`
                : 'Pick a city and dates to see prices.'}
            </p>
            <p className="mt-2 caption text-mute">
              Gợi ý: thử các thành phố đã có khách sạn —
              Ho Chi Minh, Hanoi, Da Nang.
            </p>
          </div>
          <div className="flex justify-start lg:justify-end">
            <RequestForm />
          </div>
        </div>
      </div>

      <div className="container-x pb-8 md:pb-12">
        <div className="h-divider mb-6" />

        {!city || !checkin || !checkout ? (
          <EmptyState
            title="Start with a city and dates"
            description="Use the form above to find rooms that match your dates."
            action={
              <Link to="/" className="btn-primary">
                Back to home
              </Link>
            }
          />
        ) : loading ? (
          <HotelListSkeleton />
        ) : error ? (
          <EmptyState title="We couldn't load stays" description={error} />
        ) : list.length === 0 ? (
          <EmptyState
            title="No rooms match those dates"
            description={`There are no rooms in ${city} from ${formatDate(
              checkin,
            )} to ${formatDate(checkout)}. Try different dates or another city.`}
          />
        ) : (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-5">
            {list.map((hotel) => (
              <button
                key={hotel.id}
                type="button"
                onClick={() =>
                  navigate(
                    `/hotels/${hotel.id}?checkin=${checkin}&checkout=${checkout}&rooms=${rooms}`,
                  )
                }
                className="text-left card-elevated hover:shadow-l2 transition-shadow overflow-hidden flex flex-col"
              >
                <HotelIllustration city={hotel.city} />
                <div className="p-5 flex-1 flex flex-col">
                  <div className="flex items-start gap-2">
                    <IconPin className="h-4 w-4 text-mute mt-1 shrink-0" />
                    <div>
                      <div className="display-md">{hotel.name}</div>
                      <div className="caption text-mute">{hotel.city}</div>
                    </div>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    {(hotel.room_types || []).slice(0, 3).map((rt) => (
                      <span key={rt.id} className="category-pill">
                        {rt.name}
                      </span>
                    ))}
                  </div>
                  <div className="mt-auto pt-4 flex items-center justify-between">
                    <div>
                      <div className="caption text-mute">From</div>
                      <div className="display-md tabular-nums">
                        {hotel.room_types && hotel.room_types[0]
                          ? formatVnd(hotel.room_types[0].base_price)
                          : '—'}
                      </div>
                      <div className="caption text-mute">per room / night</div>
                    </div>
                    <span className="inline-flex items-center gap-2 body-md-strong text-ink">
                      See rooms <IconArrowRight className="h-4 w-4" />
                    </span>
                  </div>
                </div>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function HotelListSkeleton() {
  return (
    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-5">
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="card-elevated overflow-hidden">
          <Skeleton className="h-40 w-full rounded-none" />
          <div className="p-5">
            <Skeleton className="h-6 w-3/4 mb-2" />
            <Skeleton className="h-4 w-1/2" />
          </div>
        </div>
      ))}
    </div>
  );
}

function HotelIllustration({ city }) {
  // City-based illustration variant so hotels read visually distinct
  const seed = (city || 'City').length;
  const variant = seed % 4;
  return (
    <svg viewBox="0 0 320 160" xmlns="http://www.w3.org/2000/svg" className="w-full h-40">
      <rect width="320" height="160" fill="#efefef" />
      {variant === 0 ? (
        <g>
          <rect x="40" y="80" width="40" height="80" fill="#000" />
          <rect x="88" y="50" width="48" height="110" fill="#000" />
          <rect x="144" y="70" width="36" height="90" fill="#000" />
          <rect x="188" y="40" width="56" height="120" fill="#000" />
          <rect x="252" y="70" width="40" height="90" fill="#000" />
        </g>
      ) : variant === 1 ? (
        <g>
          <circle cx="240" cy="50" r="24" fill="#000" />
          <rect x="60" y="100" width="80" height="60" fill="#000" />
          <rect x="160" y="80" width="40" height="80" fill="#000" />
          <rect x="220" y="120" width="60" height="40" fill="#000" />
        </g>
      ) : variant === 2 ? (
        <g>
          <path d="M0 130 Q160 70 320 130 L320 160 L0 160 Z" fill="#000" />
          <rect x="40" y="100" width="240" height="30" fill="#000" />
          <circle cx="80" cy="80" r="14" fill="#000" />
          <circle cx="240" cy="90" r="20" fill="#000" />
        </g>
      ) : (
        <g>
          <rect x="40" y="60" width="240" height="100" fill="#000" />
          <rect x="80" y="80" width="40" height="20" fill="#fff" />
          <rect x="140" y="80" width="40" height="20" fill="#fff" />
          <rect x="200" y="80" width="40" height="20" fill="#fff" />
          <rect x="60" y="120" width="20" height="40" fill="#fff" />
          <rect x="240" y="120" width="20" height="40" fill="#fff" />
        </g>
      )}
    </svg>
  );
}
