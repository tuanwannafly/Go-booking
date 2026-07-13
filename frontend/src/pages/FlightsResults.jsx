import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';
import { flights, ApiError } from '../lib/api.js';
import { IconArrowRight, IconSwap, IconPlane } from '../components/Icons.jsx';
import RequestForm from '../components/RequestForm.jsx';
import { EmptyState, Skeleton } from '../components/Primitives.jsx';
import { durationLabel, formatDate, formatTime, formatVnd } from '../lib/format.js';

const sortOptions = [
  { id: 'price', label: 'Lowest price' },
  { id: 'duration', label: 'Shortest' },
  { id: 'departure', label: 'Earliest' },
];

export default function FlightsResults() {
  const [params] = useSearchParams();
  const navigate = useNavigate();

  const origin = params.get('origin') || '';
  const destination = params.get('destination') || '';
  const date = params.get('date') || '';
  const passengers = Number(params.get('passengers') || '1');

  const [sort, setSort] = useState('price');
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!origin || !destination || !date) {
      setData(null);
      setError(null);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(null);
    flights
      .search({ origin, destination, date, page, pageSize })
      .then((res) => !cancelled && setData(res))
      .catch((err) => {
        if (cancelled) return;
        if (err instanceof ApiError) setError(err.message);
        else setError('Could not load flights');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  }, [origin, destination, date, page, pageSize]);

  const flights2 = data?.data || [];

  const sortedFlights = useMemo(() => {
    const arr = [...flights2];
    if (sort === 'price') arr.sort((a, b) => a.base_price - b.base_price);
    if (sort === 'duration') {
      arr.sort((a, b) => {
        const da = new Date(a.arrival_time) - new Date(a.departure_time);
        const db = new Date(b.arrival_time) - new Date(b.departure_time);
        return da - db;
      });
    }
    if (sort === 'departure')
      arr.sort((a, b) => new Date(a.departure_time) - new Date(b.departure_time));
    return arr;
  }, [flights2, sort]);

  const totalPages = data?.total_pages || 1;

  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-12">
        <div className="grid lg:grid-cols-2 gap-10 items-start">
          <div>
            <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
              Step 1 of 3 · Pick a flight
            </div>
            <h1 className="display-xl">
              {origin || 'Anywhere'}
              <span className="mx-2 inline-flex items-center text-body">
                <IconSwap className="h-6 w-6" />
              </span>
              {destination || 'Anywhere'}
            </h1>
            <p className="mt-3 body-md text-body">
              {date
                ? `${formatDate(date, { weekday: 'short' })} · ${passengers} ${
                    passengers === 1 ? 'passenger' : 'passengers'
                  }`
                : 'Choose an origin, a destination and a date to see fares.'}
            </p>
            <p className="mt-2 caption text-mute">
              Gợi ý: thử các chuyến đã có trong hệ thống —
              SGN↔HAN (2026-07-20), SGN↔DAD (2026-07-21), HAN↔DAD (2026-07-22), SGN→PQC (2026-07-23).
            </p>
          </div>
          <div className="flex justify-start lg:justify-end">
            <RequestForm />
          </div>
        </div>
      </div>

      <div className="container-x pb-8 md:pb-12">
        <div className="h-divider mb-6" />

        <div className="flex flex-wrap items-center gap-3 mb-6">
          <div className="caption text-mute">Sort by</div>
          <div className="flex flex-wrap gap-2">
            {sortOptions.map((opt) => (
              <button
                key={opt.id}
                type="button"
                onClick={() => setSort(opt.id)}
                className={sort === opt.id ? 'category-pill-active' : 'category-pill'}
              >
                {opt.label}
              </button>
            ))}
          </div>
        </div>

        {!origin || !destination || !date ? (
          <EmptyState
            title="Start with where you're going"
            description="Use the form above to pick an origin, a destination and a date."
            action={
              <Link to="/" className="btn-primary">
                Back to home
              </Link>
            }
          />
        ) : loading ? (
          <FlightListSkeleton />
        ) : error ? (
          <EmptyState title="We couldn't load flights" description={error} />
        ) : sortedFlights.length === 0 ? (
          <EmptyState
            title="No flights matched"
            description={`There are no flights between ${origin} and ${destination} on ${formatDate(
              date,
            )}. Try the next date or change airports.`}
          />
        ) : (
          <div className="grid gap-4">
            {sortedFlights.map((flight) => (
              <button
                key={flight.id}
                type="button"
                onClick={() =>
                  navigate(
                    `/flights/${flight.id}?date=${date}&passengers=${passengers}`,
                  )
                }
                className="text-left card-elevated hover:shadow-l2 transition-shadow p-5 md:p-6 flex flex-col md:flex-row md:items-center gap-5"
              >
                <FlightIdentity flight={flight} />
                <FlightTimes flight={flight} />
                <FlightPrice flight={flight} passengers={passengers} />
                <span className="inline-flex items-center gap-2 ml-auto md:ml-0 body-md-strong text-ink">
                  Select <IconArrowRight className="h-4 w-4" />
                </span>
              </button>
            ))}
          </div>
        )}

        {data && totalPages > 1 ? (
          <Pagination page={page} totalPages={totalPages} onChange={setPage} />
        ) : null}
      </div>
    </div>
  );
}

function FlightListSkeleton() {
  return (
    <div className="grid gap-4">
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i} className="card-elevated p-6 flex flex-col md:flex-row gap-5">
          <Skeleton className="h-12 w-12 rounded-md" />
          <div className="flex-1 grid gap-2">
            <Skeleton className="h-5 w-32" />
            <Skeleton className="h-4 w-48" />
          </div>
          <Skeleton className="h-10 w-32 rounded-pill" />
        </div>
      ))}
    </div>
  );
}

function FlightIdentity({ flight }) {
  return (
    <div className="flex items-center gap-3 w-full md:w-48">
      <div className="rounded-full bg-canvasSoft w-12 h-12 flex items-center justify-center">
        <IconPlane className="h-6 w-6" />
      </div>
      <div>
        <div className="body-sm-strong">{flight.code}</div>
        <div className="caption text-mute">
          {flight.origin} → {flight.destination}
        </div>
      </div>
    </div>
  );
}

function FlightTimes({ flight }) {
  return (
    <div className="flex items-center gap-4 md:gap-6 flex-1 min-w-[260px]">
      <div className="flex flex-col">
        <span className="display-md tabular-nums">{formatTime(flight.departure_time)}</span>
        <span className="caption text-mute">{flight.origin}</span>
      </div>
      <div className="flex-1 flex flex-col items-center gap-1">
        <span className="caption text-body">
          {durationLabel(flight.departure_time, flight.arrival_time)}
        </span>
        <div className="w-full h-px bg-ink/20 relative">
          <span className="absolute left-0 top-1/2 -translate-y-1/2 w-1.5 h-1.5 rounded-full bg-ink" />
          <span className="absolute right-0 top-1/2 -translate-y-1/2 w-1.5 h-1.5 rounded-full bg-ink" />
        </div>
        <span className="caption text-mute">Direct</span>
      </div>
      <div className="flex flex-col items-end">
        <span className="display-md tabular-nums">{formatTime(flight.arrival_time)}</span>
        <span className="caption text-mute">{flight.destination}</span>
      </div>
    </div>
  );
}

function FlightPrice({ flight, passengers }) {
  const total = (flight.base_price || 0) * Math.max(1, passengers);
  return (
    <div className="flex flex-col md:items-end">
      <span className="display-md tabular-nums">{formatVnd(flight.base_price)}</span>
      <span className="caption text-mute">
        {formatDate(flight.departure_time, {
          weekday: 'short',
          day: '2-digit',
          month: 'short',
        })}
      </span>
      {passengers > 1 ? (
        <span className="caption text-mute">
          Total · {formatVnd(total)} for {passengers}
        </span>
      ) : null}
    </div>
  );
}

function Pagination({ page, totalPages, onChange }) {
  const pages = useMemo(() => {
    const arr = [];
    for (let i = 1; i <= totalPages; i += 1) arr.push(i);
    return arr.slice(0, 5);
  }, [totalPages]);
  return (
    <div className="mt-8 flex items-center gap-2 justify-center">
      <button
        type="button"
        onClick={() => onChange(Math.max(1, page - 1))}
        className="category-pill disabled:opacity-40"
        disabled={page === 1}
      >
        Previous
      </button>
      {pages.map((p) => (
        <button
          key={p}
          type="button"
          onClick={() => onChange(p)}
          className={p === page ? 'category-pill-active' : 'category-pill'}
        >
          {p}
        </button>
      ))}
      <button
        type="button"
        onClick={() => onChange(Math.min(totalPages, page + 1))}
        className="category-pill disabled:opacity-40"
        disabled={page === totalPages}
      >
        Next
      </button>
    </div>
  );
}
