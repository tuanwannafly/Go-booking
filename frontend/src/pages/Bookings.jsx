import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { bookings, ApiError } from '../lib/api.js';
import {
  IconArrowRight,
  IconBag,
  IconCalendar,
  IconPlane,
  IconHotel,
  IconReceipt,
} from '../components/Icons.jsx';
import { EmptyState, Skeleton, StatusPill } from '../components/Primitives.jsx';
import { useAuth } from '../lib/auth.jsx';
import { useRequireAuth } from '../lib/useAuthCheck.js';
import { formatDate, formatVnd } from '../lib/format.js';

const tabs = [
  { id: 'all', label: 'All bookings' },
  { id: 'pending', label: 'Pending' },
  { id: 'confirmed', label: 'Confirmed' },
  { id: 'cancelled', label: 'Cancelled' },
];

export default function Bookings() {
  const { isAuthed } = useAuth();
  const { isReady } = useRequireAuth({ toastMessage: 'Vui lòng đăng nhập để xem đặt chỗ' });
  const navigate = useNavigate();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [tab, setTab] = useState('all');

  useEffect(() => {
    if (!isAuthed) return;
    let cancelled = false;
    setLoading(true);
    setError(null);
    bookings
      .list({ page: 1, pageSize: 20 })
      .then((res) => !cancelled && setData(res))
      .catch((err) => {
        if (cancelled) return;
        if (err instanceof ApiError) setError(err.message);
        else setError('Could not load bookings');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  }, [isAuthed]);

  // Hiển thị loading state trong khi chờ auth check
  if (!isReady) {
    return (
      <div className="container-x py-12 md:py-20">
        <div className="grid gap-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="card-elevated p-6 flex gap-5">
              <Skeleton className="h-12 w-12 rounded-md" />
              <div className="flex-1 grid gap-2">
                <Skeleton className="h-5 w-48" />
                <Skeleton className="h-4 w-72" />
              </div>
              <Skeleton className="h-10 w-24 rounded-pill" />
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (!isAuthed) {
    return (
      <div className="container-x py-12 md:py-20">
        <div className="grid md:grid-cols-2 gap-10 items-center">
          <div>
            <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
              Chuyến đi của bạn
            </div>
            <h1 className="display-xl max-w-[16ch]">
              Đăng nhập để xem tất cả đặt chỗ của bạn.
            </h1>
            <p className="mt-4 body-md text-body max-w-prose">
              Giữ ghế, xác nhận đặt phòng và thay đổi lịch trình mà không cần mở ticket hỗ trợ.
            </p>
            <div className="mt-6 flex flex-wrap gap-3">
              <Link to="/login" className="btn-primary">
                Đăng nhập
              </Link>
              <Link to="/flights" className="btn-secondary">
                Tìm chuyến bay
              </Link>
            </div>
          </div>
          <div className="card-soft p-8 md:p-10 hidden md:block">
            <ul className="grid gap-4">
              <li className="flex items-start gap-3">
                <span className="rounded-pill bg-ink text-onPrimary w-8 h-8 inline-flex items-center justify-center body-sm-strong">
                  1
                </span>
                <div>
                  <div className="body-md-strong">Tìm và giữ chỗ</div>
                  <div className="body-sm text-body">
                    Tìm chuyến bay, chọn ghế, giữ trong 10 phút.
                  </div>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <span className="rounded-pill bg-ink text-onPrimary w-8 h-8 inline-flex items-center justify-center body-sm-strong">
                  2
                </span>
                <div>
                  <div className="body-md-strong">Xác nhận chỉ với một chạm</div>
                  <div className="body-sm text-body">
                    Thanh toán mock, retries idempotent.
                  </div>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <span className="rounded-pill bg-ink text-onPrimary w-8 h-8 inline-flex items-center justify-center body-sm-strong">
                  3
                </span>
                <div>
                  <div className="body-md-strong">Thay đổi, hủy, hoàn tiền</div>
                  <div className="body-sm text-body">
                    Điều chỉnh lịch khởi hành và hoàn tiền theo chính sách.
                  </div>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>
    );
  }

  const list = (data?.data || []).filter((b) => tab === 'all' || b.status === tab);

  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-12">
        <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
          Your account
        </div>
        <h1 className="display-xl">Bookings</h1>
        <p className="mt-3 body-md text-body max-w-prose">
          Every booking, every status, one column. Hold, confirm, reschedule.
        </p>

        <div className="mt-8 flex flex-wrap gap-2">
          {tabs.map((t) => (
            <button
              key={t.id}
              type="button"
              onClick={() => setTab(t.id)}
              className={tab === t.id ? 'category-pill-active' : 'category-pill'}
            >
              {t.label}
            </button>
          ))}
        </div>
      </div>

      <div className="container-x pb-12 md:pb-20">
        {loading ? (
          <div className="grid gap-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="card-elevated p-6 flex gap-5">
                <Skeleton className="h-12 w-12 rounded-md" />
                <div className="flex-1 grid gap-2">
                  <Skeleton className="h-5 w-48" />
                  <Skeleton className="h-4 w-72" />
                </div>
                <Skeleton className="h-10 w-24 rounded-pill" />
              </div>
            ))}
          </div>
        ) : error ? (
          <EmptyState title="Could not load your bookings" description={error} />
        ) : list.length === 0 ? (
          <EmptyState
            title="No bookings here yet"
            description="Find a flight or a stay and start your first booking."
            action={
              <button onClick={() => navigate('/')} className="btn-primary">
                Start searching
              </button>
            }
          />
        ) : (
          <div className="grid gap-4">
            {list.map((b) => (
              <BookingRow key={b.id} booking={b} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function BookingRow({ booking }) {
  return (
    <Link
      to={`/bookings/${booking.id}`}
      className="card-elevated p-5 md:p-6 flex flex-col md:flex-row md:items-center gap-5 hover:shadow-l2 transition-shadow"
    >
      <div className="rounded-full bg-canvasSoft w-12 h-12 flex items-center justify-center">
        {booking.scheduled_at ? (
          <IconPlane className="h-6 w-6" />
        ) : (
          <IconHotel className="h-6 w-6" />
        )}
      </div>
      <div className="flex-1 min-w-[200px]">
        <div className="flex flex-wrap items-center gap-3 mb-1">
          <div className="display-sm">Booking #{booking.id?.slice(-6).toUpperCase()}</div>
          <StatusPill status={booking.status} />
        </div>
        <div className="body-sm text-body inline-flex items-center gap-2 flex-wrap">
          <IconCalendar className="h-4 w-4" />
          {booking.scheduled_at
            ? formatDate(booking.scheduled_at, { weekday: 'short' })
            : 'No date set'}
          <span aria-hidden>·</span>
          <IconBag className="h-4 w-4" />
          {booking.items?.length || 0} item(s)
        </div>
      </div>
      <div className="text-left md:text-right">
        <div className="display-md tabular-nums">{formatVnd(booking.total_amount)}</div>
        <div className="caption text-mute">
          Created {formatDate(booking.created_at)}
        </div>
      </div>
      <span className="inline-flex items-center gap-2 body-md-strong text-ink md:ml-4">
        Open <IconArrowRight className="h-4 w-4" />
      </span>
    </Link>
  );
}
