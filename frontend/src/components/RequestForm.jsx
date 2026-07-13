import { useNavigate, useSearchParams } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';
import {
  IconCalendar,
  IconClock,
  IconDot,
  IconPlane,
  IconHotel,
  IconSwap,
  IconArrowRight,
} from './Icons.jsx';
import LocationInput from './LocationInput.jsx';
import { AIRPORTS, CITIES } from '../lib/locations.js';

const today = () => {
  const d = new Date();
  d.setHours(0, 0, 0, 0);
  return d;
};

const ymd = (d) => {
  if (!d) return '';
  const dt = d instanceof Date ? d : new Date(d);
  const yyyy = dt.getFullYear();
  const mm = `${dt.getMonth() + 1}`.padStart(2, '0');
  const dd = `${dt.getDate()}`.padStart(2, '0');
  return `${yyyy}-${mm}-${dd}`;
};

const tomorrow = () => {
  const t = today();
  t.setDate(t.getDate() + 1);
  return t;
};

const dayAfter = (days) => {
  const t = today();
  t.setDate(t.getDate() + days);
  return t;
};

const weekdayLabel = (d) => {
  const dt = d instanceof Date ? d : new Date(d);
  return dt.toLocaleDateString(undefined, { weekday: 'short' });
};

const isoWeek = (d) => {
  const dt = d instanceof Date ? d : new Date(d);
  const day = dt.getDay() || 7;
  dt.setDate(dt.getDate() + 4 - day);
  const yearStart = new Date(dt.getFullYear(), 0, 1);
  return Math.ceil(((dt - yearStart) / 86400000 + 1) / 7);
};

const prettyDate = (d) => {
  const dt = d instanceof Date ? d : new Date(d);
  return dt.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
  });
};

export default function RequestForm() {
  const [mode, setMode] = useState('flights'); // flights | hotels
  const [params] = useSearchParams();
  const navigate = useNavigate();

  // flight state
  const [origin, setOrigin] = useState('');
  const [originLabel, setOriginLabel] = useState('');
  const [destination, setDestination] = useState('');
  const [destinationLabel, setDestinationLabel] = useState('');
  const [departureDate, setDepartureDate] = useState(tomorrow());
  const [passengers, setPassengers] = useState(1);

  // hotel state
  const [city, setCity] = useState('');
  const [cityLabel, setCityLabel] = useState('');
  const [checkIn, setCheckIn] = useState(tomorrow());
  const [checkOut, setCheckOut] = useState(dayAfter(3));
  const [rooms, setRooms] = useState(1);

  useEffect(() => {
    const tab = params.get('tab');
    if (tab === 'hotels') setMode('hotels');
    if (tab === 'flights') setMode('flights');
  }, [params]);

  const flightNights = useMemo(
    () => Math.max(1, Math.round((+departureDate - +today()) / 86400000) || 1),
    [departureDate],
  );

  const hotelNights = useMemo(
    () => Math.max(1, Math.round((+checkOut - +checkIn) / 86400000)),
    [checkIn, checkOut],
  );

  const submit = (e) => {
    e.preventDefault();
    if (mode === 'flights') {
      const qs = new URLSearchParams({
        origin: origin.trim(),
        destination: destination.trim(),
        date: ymd(departureDate),
        passengers: String(passengers),
      });
      navigate(`/flights?${qs.toString()}`);
    } else {
      const qs = new URLSearchParams({
        city: city.trim(),
        checkin: ymd(checkIn),
        checkout: ymd(checkOut),
        rooms: String(rooms),
      });
      navigate(`/hotels?${qs.toString()}`);
    }
  };

  const swap = () => {
    if (mode === 'flights') {
      setOrigin(destination);
      setOriginLabel(destinationLabel);
      setDestination(origin);
      setDestinationLabel(originLabel);
    }
  };

  return (
    <div className="request-form-card p-4 md:p-5 w-full max-w-[490px]">
      {/* Pill-tab toggle */}
      <div className="flex items-center gap-2 mb-4 md:mb-5">
        <div className="bg-canvasSoft rounded-pillTab p-1 inline-flex">
          <button
            type="button"
            onClick={() => setMode('flights')}
            className={`btn-tab-translucent ${mode === 'flights' ? 'shadow-l3' : ''}`}
            aria-pressed={mode === 'flights'}
          >
            <IconPlane className="h-4 w-4" />
            Flights
          </button>
          <button
            type="button"
            onClick={() => setMode('hotels')}
            className={`btn-tab-translucent ${mode === 'hotels' ? 'shadow-l3' : ''}`}
            aria-pressed={mode === 'hotels'}
          >
            <IconHotel className="h-4 w-4" />
            Stays
          </button>
        </div>
      </div>

      <form onSubmit={submit} className="flex flex-col gap-3">
        {mode === 'flights' ? (
          <>
            <LocationInput
              label="From"
              value={origin}
              onChange={(v) => {
                setOrigin(v.toUpperCase());
              }}
              onSelect={(item) => {
                setOrigin(item.code);
                setOriginLabel(item.name);
              }}
              placeholder="Sân bay đi (SGN, HAN…)"
              suggestions={AIRPORTS}
              required
            />
            <div className="flex justify-center -my-1.5 relative">
              <button
                type="button"
                onClick={swap}
                className="rounded-full bg-canvas shadow-l3 p-2 hover:bg-canvasSoft ring-focus"
                aria-label="Swap origin and destination"
              >
                <IconSwap className="h-4 w-4" />
              </button>
            </div>
            <LocationInput
              label="To"
              value={destination}
              onChange={(v) => {
                setDestination(v.toUpperCase());
              }}
              onSelect={(item) => {
                setDestination(item.code);
                setDestinationLabel(item.name);
              }}
              placeholder="Sân bay đến (SGN, HAN…)"
              suggestions={AIRPORTS.filter((a) => a.code !== origin)}
              required
            />

            <DateChooser label="Departure" value={departureDate} onChange={setDepartureDate} />

            <Counter
              label="Passengers"
              value={passengers}
              min={1}
              max={9}
              onChange={setPassengers}
            />
          </>
        ) : (
          <>
            <LocationInput
              label="City"
              value={city}
              onChange={(v) => setCity(v)}
              onSelect={(item) => {
                setCity(item.code);
                setCityLabel(item.name);
              }}
              placeholder="Thành phố (Ho Chi Minh, Hanoi…)"
              suggestions={CITIES}
              required
            />

            <div className="grid grid-cols-2 gap-3">
              <div className="request-form-row !rounded-md">
                <IconCalendar className="h-5 w-5 text-ink/70 shrink-0" />
                <input
                  required
                  type="date"
                  value={ymd(checkIn)}
                  min={ymd(today())}
                  onChange={(e) => {
                    const d = new Date(e.target.value);
                    if (!Number.isNaN(d.getTime())) {
                      setCheckIn(d);
                      if (+d >= +checkOut) {
                        const out = new Date(d);
                        out.setDate(out.getDate() + 1);
                        setCheckOut(out);
                      }
                    }
                  }}
                  className="bg-transparent flex-1 outline-none body-md"
                />
              </div>
              <div className="request-form-row !rounded-md">
                <IconCalendar className="h-5 w-5 text-ink/70 shrink-0" />
                <input
                  required
                  type="date"
                  value={ymd(checkOut)}
                  min={ymd(new Date(+checkIn + 86400000))}
                  onChange={(e) => {
                    const d = new Date(e.target.value);
                    if (!Number.isNaN(d.getTime())) setCheckOut(d);
                  }}
                  className="bg-transparent flex-1 outline-none body-md"
                />
              </div>
            </div>

            <Counter label="Rooms" value={rooms} min={1} max={6} onChange={setRooms} />
          </>
        )}

        <div className="flex items-center justify-between pt-2">
          <div className="body-sm text-body">
            {mode === 'flights'
              ? `${flightNights} ${flightNights === 1 ? 'day' : 'days'} · ${passengers} ${passengers === 1 ? 'passenger' : 'passengers'}`
              : `${hotelNights} ${hotelNights === 1 ? 'night' : 'nights'} · ${rooms} ${rooms === 1 ? 'room' : 'rooms'}`}
          </div>
          <button type="submit" className="btn-primary px-7">
            See prices
            <IconArrowRight className="h-4 w-4" />
          </button>
        </div>
      </form>
    </div>
  );
}

function Counter({ label, value, min, max, onChange }) {
  return (
    <div className="request-form-row">
      <IconClock className="h-5 w-5 text-ink/70 shrink-0" />
      <span className="flex-1 body-md">{label}</span>
      <button
        type="button"
        onClick={() => onChange(Math.max(min, value - 1))}
        className="rounded-full bg-canvas w-8 h-8 flex items-center justify-center hover:bg-surfacePressed ring-focus"
        aria-label={`Decrease ${label}`}
      >
        –
      </button>
      <span className="w-6 text-center body-md-strong tabular-nums">{value}</span>
      <button
        type="button"
        onClick={() => onChange(Math.min(max, value + 1))}
        className="rounded-full bg-canvas w-8 h-8 flex items-center justify-center hover:bg-surfacePressed ring-focus"
        aria-label={`Increase ${label}`}
      >
        +
      </button>
    </div>
  );
}

function DateChooser({ label, value, onChange }) {
  const days = useMemo(() => {
    const t = today();
    return Array.from({ length: 8 }, (_, i) => {
      const d = new Date(t);
      d.setDate(d.getDate() + i);
      return d;
    });
  }, []);

  const changeOffset = (offset) => {
    const d = new Date(value);
    d.setDate(d.getDate() + offset);
    onChange(d);
  };

  return (
    <div className="request-form-row !rounded-md flex-col items-stretch gap-3">
      <div className="flex items-center gap-3">
        <IconCalendar className="h-5 w-5 text-ink/70 shrink-0" />
        <span className="flex-1 body-md">{label}</span>
        <input
          required
          type="date"
          value={ymd(value)}
          onChange={(e) => {
            const d = new Date(e.target.value);
            if (!Number.isNaN(d.getTime())) onChange(d);
          }}
          min={ymd(today())}
          className="bg-canvas text-ink rounded-md px-2 py-1 body-sm ring-focus"
        />
      </div>
      <div className="flex gap-2 overflow-x-auto thin-scroll">
        <button
          type="button"
          onClick={() => changeOffset(-1)}
          className="rounded-md bg-canvas w-8 h-9 flex items-center justify-center hover:bg-canvasSoft ring-focus text-body"
          aria-label="Previous day"
        >
          ‹
        </button>
        {days.map((d) => {
          const isSelected = ymd(d) === ymd(value);
          return (
            <button
              type="button"
              key={ymd(d)}
              onClick={() => onChange(d)}
              className={`shrink-0 rounded-md px-3 h-9 flex flex-col items-center justify-center body-sm-strong ring-focus ${
                isSelected ? 'bg-ink text-onPrimary' : 'bg-canvas hover:bg-canvasSoft text-ink'
              }`}
            >
              <span className="caption">{weekdayLabel(d)}</span>
              <span>{prettyDate(d)}</span>
            </button>
          );
        })}
        <button
          type="button"
          onClick={() => changeOffset(1)}
          className="rounded-md bg-canvas w-8 h-9 flex items-center justify-center hover:bg-canvasSoft ring-focus text-body"
          aria-label="Next day"
        >
          ›
        </button>
      </div>
    </div>
  );
}

// Used icons referenced via JSX above. Re-export to keep tree-shaking honest.
void IconDot;
