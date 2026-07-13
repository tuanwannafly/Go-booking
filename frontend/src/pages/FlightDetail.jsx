import { useEffect, useState, useMemo } from 'react';
import { Link, useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { flights, holds, bookings, ApiError } from '../lib/api.js';
import {
  IconArrowRight,
  IconCheck,
  IconClock,
  IconPlane,
} from '../components/Icons.jsx';
import { EmptyState, Skeleton, StatusPill } from '../components/Primitives.jsx';
import { useAuth } from '../lib/auth.jsx';
import { useToast } from '../components/Toast.jsx';
import { durationLabel, formatDate, formatTime, formatVnd } from '../lib/format.js';

const seatClassMap = {
  J: { label: 'Business', per: 'Upgrade to lie-flat for a small premium' },
  Y: { label: 'Economy', per: 'Free baggage, lounge snacks, and earn miles' },
};

export default function FlightDetail() {
  const { id } = useParams();
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { isAuthed } = useAuth();
  const toast = useToast();

  const date = params.get('date') || '';
  const passengers = Number(params.get('passengers') || '1');

  const [flight, setFlight] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedSeat, setSelectedSeat] = useState(null);
  const [holding, setHolding] = useState(false);
  const [hold, setHold] = useState(null);
  const [seatPlaceholder] = useState(() =>
    generateSeatLayout(120, 6),
  );

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    flights
      .get(id)
      .then((f) => !cancelled && setFlight(f))
      .catch((err) => {
        if (cancelled) return;
        if (err instanceof ApiError) setError(err.message);
        else setError('Could not load this flight');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  }, [id]);

  const groupedSeats = useMemo(() => groupSeats(seatPlaceholder), [seatPlaceholder]);

  const chooseSeat = (seat) => {
    setSelectedSeat(seat);
    setHold(null);
  };

  const onHold = async () => {
    if (!isAuthed) {
      toast.info('Vui lòng đăng nhập để giữ ghế');
      navigate(`/login?next=${encodeURIComponent(window.location.pathname + window.location.search)}`);
      return;
    }
    if (!selectedSeat) {
      toast.info('Vui lòng chọn ghế trước');
      return;
    }
    setHolding(true);
    try {
      const res = await holds.seat(id, selectedSeat.code, { duration: 10 });
      setHold(res);
      toast.success('Đã giữ ghế trong 10 phút');
    } catch (err) {
      toast.error(err.message || 'Không thể giữ ghế này');
    } finally {
      setHolding(false);
    }
  };

  const onConfirm = async () => {
    if (!hold) {
      toast.info('Vui lòng giữ ghế trước');
      return;
    }
    try {
      const created = await bookings.create([
        { inventory_unit_id: hold.inventory_unit?.id || hold.inventory_unit_id || '00000000-0000-0000-0000-000000000000' },
      ]);
      toast.success('Đã tạo đặt chỗ. Đang chuyển đến trang thanh toán…');
      navigate(`/bookings/${created.id}`);
    } catch (err) {
      toast.error(err.message || 'Không thể tạo đặt chỗ');
    }
  };

  if (loading) {
    return <DetailSkeleton />;
  }
  if (error || !flight) {
    return (
      <div className="container-x py-12">
        <EmptyState
          title="We couldn't load this flight"
          description={error || 'Try again or pick a different flight.'}
          action={
            <Link to="/flights" className="btn-primary">
              Back to results
            </Link>
          }
        />
      </div>
    );
  }

  const perSeat = flight.base_price;
  const total = perSeat * Math.max(1, passengers);

  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-12">
        <Link to="/flights" className="body-sm text-body hover:underline">
          ← Back to results
        </Link>

        <header className="mt-4 flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
          <div>
            <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
              {flight.code} · {flight.origin} → {flight.destination}
            </div>
            <h1 className="display-xl max-w-[18ch]">
              {formatDate(flight.departure_time, {
                weekday: 'short',
                day: '2-digit',
                month: 'long',
              })}
            </h1>
            <div className="mt-2 body-md text-body">
              {durationLabel(flight.departure_time, flight.arrival_time)} ·{' '}
              {formatTime(flight.departure_time)} → {formatTime(flight.arrival_time)}
            </div>
          </div>
          <div className="flex flex-col items-end">
            <span className="display-md tabular-nums">{formatVnd(perSeat)}</span>
            <span className="caption text-mute">per traveller · taxes included</span>
          </div>
        </header>
      </div>

      <div className="container-x pb-12 md:pb-20 grid lg:grid-cols-[1fr_400px] gap-8">
        <div className="flex flex-col gap-8">
          {/* Seat map */}
          <section className="card-elevated p-6 md:p-8">
            <div className="flex flex-col md:flex-row md:items-end md:justify-between gap-3 mb-6">
              <div>
                <div className="caption uppercase tracking-[0.15em] text-mute mb-1">
                  Step 2 of 3 · Pick a seat
                </div>
                <h2 className="display-lg">Choose where you sit</h2>
              </div>
              <div className="flex flex-wrap gap-3 body-sm text-body">
                <Legend swatch="bg-canvas border border-surfacePressed">Available</Legend>
                <Legend swatch="bg-ink text-onPrimary">Selected</Legend>
                <Legend swatch="bg-canvasSoft">Held</Legend>
              </div>
            </div>

            <div className="grid md:grid-cols-2 gap-6">
              {Object.entries(groupedSeats).map(([cls, seats]) => {
                const meta = seatClassMap[cls] || { label: cls, per: '' };
                return (
                  <div key={cls} className="card-soft p-5">
                    <div className="flex items-center justify-between mb-3">
                      <div>
                        <div className="body-md-strong">{meta.label}</div>
                        <div className="caption text-mute">{meta.per}</div>
                      </div>
                      <div className="caption text-mute">{seats.length} seats</div>
                    </div>
                    <div className="grid grid-cols-6 gap-1.5">
                      {seats.map((seat) => {
                        const isSelected = selectedSeat?.code === seat.code;
                        return (
                          <button
                            key={seat.code}
                            type="button"
                            onClick={() => chooseSeat(seat)}
                            className={`h-9 rounded-md text-[11px] font-medium tabular-nums ring-focus ${
                              isSelected
                                ? 'bg-ink text-onPrimary'
                                : 'bg-canvas hover:bg-canvasSoft border border-surfacePressed'
                            }`}
                            aria-pressed={isSelected}
                            aria-label={`Seat ${seat.code}`}
                          >
                            {seat.code.replace(/^[A-Z]/, '')}
                          </button>
                        );
                      })}
                    </div>
                  </div>
                );
              })}
            </div>

            {selectedSeat ? (
              <div className="mt-6 flex flex-wrap items-center gap-4 p-4 rounded-xl bg-canvasSoft">
                <div className="flex items-center gap-3">
                  <IconCheck className="h-5 w-5" />
                  <div>
                    <div className="body-sm-strong">Seat {selectedSeat.code}</div>
                    <div className="caption text-mute">
                      {seatClassMap[selectedSeat.code[0]]?.label || 'Seat'} · {selectedSeat.position}
                    </div>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={onHold}
                  disabled={holding}
                  className="btn-primary ml-auto"
                >
                  {holding ? 'Holding…' : 'Hold this seat'}
                  {!holding ? <IconArrowRight className="h-4 w-4" /> : null}
                </button>
              </div>
            ) : null}
          </section>

          {/* Itinerary */}
          <section className="card p-6 md:p-8 grid md:grid-cols-2 gap-6">
            <div>
              <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
                Itinerary
              </div>
              <div className="space-y-4">
                <ItineraryRow
                  code={flight.origin}
                  time={flight.departure_time}
                  label="Departure"
                />
                <ItineraryRow
                  code={flight.destination}
                  time={flight.arrival_time}
                  label="Arrival"
                />
              </div>
            </div>
            <div className="card-soft p-5">
              <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
                What you get
              </div>
              <ul className="grid gap-2 body-md">
                <li className="flex items-center gap-2">
                  <IconCheck className="h-4 w-4" /> 7kg cabin baggage
                </li>
                <li className="flex items-center gap-2">
                  <IconCheck className="h-4 w-4" /> Free reschedule (24h notice)
                </li>
                <li className="flex items-center gap-2">
                  <IconCheck className="h-4 w-4" /> Cancellation within policy
                </li>
                <li className="flex items-center gap-2">
                  <IconCheck className="h-4 w-4" /> Mock payment gateway
                </li>
              </ul>
            </div>
          </section>
        </div>

        {/* Right rail summary */}
        <aside className="lg:sticky lg:top-24 h-fit">
          <div className="card-elevated p-6 md:p-7 flex flex-col gap-4">
            <div className="flex items-center gap-3">
              <div className="rounded-full bg-canvasSoft w-10 h-10 flex items-center justify-center">
                <IconPlane className="h-5 w-5" />
              </div>
              <div>
                <div className="body-sm-strong">{flight.code}</div>
                <div className="caption text-mute">
                  {flight.origin} → {flight.destination}
                </div>
              </div>
            </div>

            <div className="h-divider" />

            <div className="grid grid-cols-2 gap-3 body-md">
              <SummaryRow label="Date" value={formatDate(flight.departure_time)} />
              <SummaryRow
                label="Time"
                value={`${formatTime(flight.departure_time)}`}
              />
              <SummaryRow
                label="Duration"
                value={durationLabel(flight.departure_time, flight.arrival_time)}
              />
              <SummaryRow
                label="Seats left"
                value={flight.available_seats ?? '—'}
              />
              {passengers > 1 ? (
                <SummaryRow label="Passengers" value={passengers} />
              ) : null}
              {selectedSeat ? <SummaryRow label="Seat" value={selectedSeat.code} /> : null}
            </div>

            <div className="h-divider" />

            <div className="flex items-baseline justify-between">
              <span className="body-md-strong">Total</span>
              <span className="display-lg tabular-nums">{formatVnd(total)}</span>
            </div>
            <div className="caption text-mute">
              {passengers > 1
                ? `${passengers} passengers × ${formatVnd(perSeat)}`
                : 'Inclusive of taxes'}
            </div>

            {hold ? (
              <div className="card-soft p-4 flex items-center gap-3">
                <IconClock className="h-5 w-5" />
                <div className="flex-1">
                  <div className="body-sm-strong">Seat on hold</div>
                  <div className="caption text-body">
                    Expires {hold.held_until ? formatTime(hold.held_until) : 'soon'}
                  </div>
                </div>
                <StatusPill status="pending" />
              </div>
            ) : null}

            <button
              type="button"
              onClick={onConfirm}
              disabled={!hold}
              className="btn-primary w-full"
            >
              {hold ? 'Continue to confirmation' : 'Hold a seat first'}
            </button>
            <button
              type="button"
              onClick={() => navigate(-1)}
              className="btn-secondary w-full"
            >
              Back
            </button>

            <p className="caption text-mute">
              Hold is held by your account for 10 minutes. Confirm before the timer
              runs out and the seat returns to the pool.
            </p>
          </div>
        </aside>
      </div>
    </div>
  );
}

function Legend({ swatch, children }) {
  return (
    <span className="inline-flex items-center gap-2">
      <span className={`inline-block w-3 h-3 rounded ${swatch}`} />
      {children}
    </span>
  );
}

function SummaryRow({ label, value }) {
  return (
    <div className="flex flex-col">
      <span className="caption text-mute">{label}</span>
      <span className="body-md-strong">{value}</span>
    </div>
  );
}

function ItineraryRow({ code, time, label }) {
  return (
    <div className="flex items-center gap-3">
      <div className="display-md tabular-nums w-16">{formatTime(time)}</div>
      <div className="flex-1 h-px bg-ink/15 relative">
        <span className="absolute -top-1 left-0 w-2 h-2 bg-ink rounded-full" />
        <span className="absolute -top-1 right-0 w-2 h-2 bg-ink rounded-full" />
      </div>
      <div className="text-right">
        <div className="body-sm-strong">{code}</div>
        <div className="caption text-mute">{label}</div>
      </div>
    </div>
  );
}

function generateSeatLayout(rows, perRow) {
  const codeOf = (row, letter) => `${row}${letter}`;
  const out = [];
  for (let r = 1; r <= rows; r += 1) {
    const cls = r <= 12 ? 'J' : 'Y';
    for (let l = 0; l < perRow; l += 1) {
      const letter = String.fromCharCode(65 + l);
      out.push({
        code: codeOf(cls === 'J' ? `J${r}` : r, letter),
        cls,
        position: l === 0 || l === perRow - 1 ? 'Window' : l === 1 || l === perRow - 2 ? 'Middle' : 'Aisle',
      });
    }
  }
  return out;
}

function groupSeats(seats) {
  const out = {};
  for (const s of seats) {
    if (!out[s.cls]) out[s.cls] = [];
    out[s.cls].push(s);
  }
  return out;
}

function DetailSkeleton() {
  return (
    <div className="container-x py-8 md:py-12">
      <Skeleton className="h-6 w-32 mb-4" />
      <Skeleton className="h-10 w-72 mb-6" />
      <div className="grid lg:grid-cols-[1fr_400px] gap-8">
        <Skeleton className="h-96 rounded-xl" />
        <Skeleton className="h-72 rounded-xl" />
      </div>
    </div>
  );
}
