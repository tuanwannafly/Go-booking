import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { bookings, ApiError } from '../lib/api.js';
import {
  IconArrowRight,
  IconCalendar,
  IconCheck,
  IconClock,
  IconReceipt,
} from '../components/Icons.jsx';
import { EmptyState, Modal, Skeleton, StatusPill } from '../components/Primitives.jsx';
import { useAuth } from '../lib/auth.jsx';
import { useRequireAuth } from '../lib/useAuthCheck.js';
import { useToast } from '../components/Toast.jsx';
import { formatDate, formatTime, formatVnd } from '../lib/format.js';

const TIMELINE = ['pending', 'confirmed', 'cancelled'];

export default function BookingDetail() {
  const { id } = useParams();
  const { isAuthed } = useAuth();
  const { isReady } = useRequireAuth({
    redirectTo: `/bookings/${id}`,
    showToast: false,
  });
  const navigate = useNavigate();
  const toast = useToast();

  const [booking, setBooking] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [working, setWorking] = useState(false);

  const [confirmOpen, setConfirmOpen] = useState(false);
  const [cancelOpen, setCancelOpen] = useState(false);
  const [rescheduleOpen, setRescheduleOpen] = useState(false);

  const refresh = () => {
    if (!isAuthed) return;
    let cancelled = false;
    setLoading(true);
    bookings
      .get(id)
      .then((b) => !cancelled && setBooking(b))
      .catch((err) => {
        if (cancelled) return;
        if (err instanceof ApiError) setError(err.message);
        else setError('Không thể tải thông tin đặt chỗ');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  };

  useEffect(() => {
    if (!isReady) return;
    if (!isAuthed) {
      navigate(`/login?next=${encodeURIComponent(window.location.pathname)}`);
      return;
    }
    const cleanup = refresh();
    return cleanup;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, isAuthed, isReady]);

  if (loading) return <DetailSkeleton />;
  if (error || !booking) {
    return (
      <div className="container-x py-12">
        <EmptyState
          title="We couldn't load this booking"
          description={error || 'Try again from your bookings list.'}
          action={
            <Link to="/bookings" className="btn-primary">
              Back to bookings
            </Link>
          }
        />
      </div>
    );
  }

  const onConfirm = async (paymentMethod) => {
    setWorking(true);
    try {
      const updated = await bookings.confirm(booking.id, paymentMethod);
      setBooking(updated);
      toast.success('Booking confirmed. Mock payment captured.');
      setConfirmOpen(false);
    } catch (err) {
      toast.error(err.message || 'Could not confirm booking');
    } finally {
      setWorking(false);
    }
  };

  const onCancel = async (reason) => {
    setWorking(true);
    try {
      const updated = await bookings.cancel(booking.id, reason);
      setBooking(updated);
      toast.success('Booking cancelled. Inventory returned to the pool.');
      setCancelOpen(false);
    } catch (err) {
      toast.error(err.message || 'Could not cancel booking');
    } finally {
      setWorking(false);
    }
  };

  const onReschedule = async (when) => {
    setWorking(true);
    try {
      const updated = await bookings.schedule(booking.id, when);
      setBooking(updated);
      toast.success(`Đã đổi lịch sang ${formatDate(when, { weekday: 'short' })}`);
      setRescheduleOpen(false);
    } catch (err) {
      toast.error(err.message || 'Không thể đổi lịch');
    } finally {
      setWorking(false);
    }
  };

  const stepIndex = TIMELINE.indexOf(booking.status);
  const total = booking.total_amount;

  const actionLabel = booking.status === 'pending'
    ? 'Xác nhận trước khi hết thời gian giữ chỗ.'
    : booking.status === 'confirmed'
      ? 'Thay đổi hoặc hủy theo chính sách.'
      : 'Không còn thao tác nào khả dụng.';

  const statusLabel = booking.status === 'confirmed'
    ? 'Đã thanh toán · mock provider'
    : booking.status === 'cancelled'
      ? 'Hoàn tiền đang xử lý'
      : 'Đang chờ xác nhận';

  const confirmTitle = 'Xác nhận và thanh toán';
  const notYetLabel = 'Chưa';
  const payLabel = `Thanh toán ${formatVnd(total)}`;

  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-12">
        <Link to="/bookings" className="body-sm text-body hover:underline">
          ← Back to bookings
        </Link>

        <header className="mt-4 flex flex-col md:flex-row md:items-end md:justify-between gap-4">
          <div>
            <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
              Booking #{booking.id?.slice(-8).toUpperCase()}
            </div>
            <h1 className="display-xl">Your trip, in one column</h1>
            <div className="mt-3 body-md text-body inline-flex items-center gap-2 flex-wrap">
              <StatusPill status={booking.status} />
              <span aria-hidden>·</span>
              <IconCalendar className="h-4 w-4" />
              {booking.scheduled_at
                ? formatDate(booking.scheduled_at, { weekday: 'short' })
                : 'No date set yet'}
              <span aria-hidden>·</span>
              <span>
                {booking.items?.length || 0} item{booking.items?.length === 1 ? '' : 's'}
              </span>
            </div>
          </div>
          <div className="flex flex-col items-start md:items-end">
            <span className="display-md tabular-nums">{formatVnd(booking.total_amount)}</span>
            <span className="caption text-mute">{statusLabel}</span>
          </div>
        </header>
      </div>

      <div className="container-x pb-12 md:pb-20 grid lg:grid-cols-[1fr_400px] gap-8">
        <div className="flex flex-col gap-8">
          <StatusTimeline status={booking.status} />

          <section className="card p-6 md:p-8">
            <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
              Items
            </div>
            <div className="grid gap-3">
              {(booking.items || []).map((item, i) => (
                <div
                  key={item.id || i}
                  className="flex items-center justify-between p-3 rounded-md bg-canvasSoft"
                >
                  <div>
                    <div className="body-md-strong">Item #{i + 1}</div>
                    <div className="caption text-mute">
                      Unit: {item.inventory_unit_id?.slice(-8) || '—'}
                    </div>
                  </div>
                  <div className="tabular-nums body-md-strong">
                    {formatVnd(item.price)}
                  </div>
                </div>
              ))}
              {(!booking.items || booking.items.length === 0) ? (
                <div className="caption text-mute">No items in this booking.</div>
              ) : null}
            </div>
          </section>

          <section className="card p-6 md:p-8">
            <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
              Timeline
            </div>
            <ol className="relative pl-6 body-md">
              <li className="relative pb-4">
                <span className="absolute left-0 top-2 w-3 h-3 rounded-full bg-ink" />
                <div className="body-md-strong">Created</div>
                <div className="caption text-mute">
                  {formatDate(booking.created_at, { weekday: 'short' })} ·{' '}
                  {formatTime(booking.created_at)}
                </div>
              </li>
              <li className="relative pb-4">
                <span
                  className={`absolute left-0 top-2 w-3 h-3 rounded-full ${
                    booking.status === 'confirmed' ? 'bg-ink' : 'bg-surfacePressed'
                  }`}
                />
                <div className="body-md-strong">Confirmed</div>
                <div className="caption text-mute">
                  {booking.status === 'confirmed'
                    ? `${formatDate(booking.updated_at)} · ${formatTime(booking.updated_at)}`
                    : 'Awaiting payment'}
                </div>
              </li>
              <li className="relative pb-2">
                <span
                  className={`absolute left-0 top-2 w-3 h-3 rounded-full ${
                    booking.status === 'cancelled' ? 'bg-ink' : 'bg-surfacePressed'
                  }`}
                />
                <div className="body-md-strong">Cancelled</div>
                <div className="caption text-mute">
                  {booking.status === 'cancelled'
                    ? `${formatDate(booking.updated_at)} · ${formatTime(booking.updated_at)}`
                    : 'Not cancelled'}
                </div>
              </li>
            </ol>
          </section>
        </div>

        <aside className="lg:sticky lg:top-24 h-fit">
          <div className="card-elevated p-6 md:p-7 flex flex-col gap-4">
            <div className="flex items-center gap-3">
              <div className="rounded-full bg-canvasSoft w-10 h-10 flex items-center justify-center">
                <IconReceipt className="h-5 w-5" />
              </div>
              <div>
                <div className="body-sm-strong">Thao tác</div>
                <div className="caption text-mute">{actionLabel}</div>
              </div>
            </div>

            {booking.status === 'pending' ? (
              <>
                <button
                  type="button"
                  onClick={() => setConfirmOpen(true)}
                  className="btn-primary w-full"
                >
                  Xác nhận và thanh toán
                  <IconArrowRight className="h-4 w-4" />
                </button>
                <button
                  type="button"
                  onClick={() => setRescheduleOpen(true)}
                  className="btn-secondary w-full"
                >
                  Đổi lịch
                </button>
                <button
                  type="button"
                  onClick={() => setCancelOpen(true)}
                  className="btn-subtle w-full"
                >
                  Hủy đặt chỗ
                </button>
              </>
            ) : booking.status === 'confirmed' ? (
              <>
                <button
                  type="button"
                  onClick={() => setRescheduleOpen(true)}
                  className="btn-primary w-full"
                >
                  Đổi lịch
                </button>
                <button
                  type="button"
                  onClick={() => setCancelOpen(true)}
                  className="btn-secondary w-full"
                >
                  Hủy đặt chỗ
                </button>
              </>
            ) : (
              <Link to="/flights" className="btn-primary w-full">
                Tìm chuyến khác
              </Link>
            )}

            <div className="h-divider" />

            <div className="grid grid-cols-2 gap-3 body-sm text-body">
              <div>
                <div className="caption text-mute">Booking ID</div>
                <div className="body-md-strong tabular-nums break-all">
                  {booking.id}
                </div>
              </div>
              {booking.idempotency_key ? (
                <div>
                  <div className="caption text-mute">Idempotency key</div>
                  <div className="body-md-strong tabular-nums break-all">
                    {booking.idempotency_key}
                  </div>
                </div>
              ) : null}
            </div>
          </div>
        </aside>
      </div>

      <ConfirmDialog
        open={confirmOpen}
        onClose={() => setConfirmOpen(false)}
        working={working}
        total={booking.total_amount}
        onConfirm={() => onConfirm('mock_card')}
      />
      <CancelDialog
        open={cancelOpen}
        onClose={() => setCancelOpen(false)}
        working={working}
        onConfirm={(reason) => onCancel(reason)}
      />
      <RescheduleDialog
        open={rescheduleOpen}
        onClose={() => setRescheduleOpen(false)}
        working={working}
        initial={booking.scheduled_at}
        onConfirm={(when) => onReschedule(when)}
      />
    </div>
  );
}

function StatusTimeline({ status }) {
  const pending = status === 'pending';
  return (
    <div className="card p-6 md:p-7">
      <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
        Status
      </div>
      <div className="grid grid-cols-3 gap-2">
        {TIMELINE.map((s, i) => {
          const reached = i <= TIMELINE.indexOf(status);
          const current = s === status;
          return (
            <div
              key={s}
              className={`rounded-md px-3 py-2 body-sm-strong capitalize ${
                current
                  ? 'bg-ink text-onPrimary'
                  : reached
                    ? 'bg-canvasSoft text-ink'
                    : 'bg-canvas text-mute border border-surfacePressed'
              }`}
            >
              {s}
            </div>
          );
        })}
      </div>
      {pending ? (
        <div className="mt-4 body-sm text-body inline-flex items-center gap-2">
          <IconClock className="h-4 w-4" />
          Hold expires soon. Confirm to keep this booking.
        </div>
      ) : null}
    </div>
  );
}

function ConfirmDialog({ open, onClose, working, total, onConfirm }) {
  return (
    <Modal
      open={open}
      onClose={onClose}
      title="Xác nhận và thanh toán"
      footer={
        <>
          <button className="btn-secondary" onClick={onClose} disabled={working}>
            {notYetLabel}
          </button>
          <button className="btn-primary" onClick={onConfirm} disabled={working}>
            {working ? 'Đang xử lý…' : payLabel}
            {!working ? <IconCheck className="h-4 w-4" /> : null}
          </button>
        </>
      }
    >
      <p className="body-md text-body">
        Chúng tôi sẽ giả lập thanh toán {formatVnd(total)} cho thẻ của bạn. Ghế sẽ được giữ cho bạn trong suốt ngày.
      </p>
      <p className="mt-3 body-sm text-body">
        Hủy trước 24 giờ trước giờ khởi hành để được hoàn tiền 100%; hủy từ 6-24 giờ hoàn 50%; dưới 6 giờ không hoàn tiền.
      </p>
    </Modal>
  );
}

function CancelDialog({ open, onClose, working, onConfirm }) {
  const [reason, setReason] = useState('');
  return (
    <Modal
      open={open}
      onClose={onClose}
      title="Cancel this booking?"
      footer={
        <>
          <button className="btn-secondary" onClick={onClose} disabled={working}>
            Keep booking
          </button>
          <button
            className="btn-primary"
            onClick={() => onConfirm(reason || 'User requested cancellation')}
            disabled={working}
          >
            {working ? 'Cancelling…' : 'Confirm cancellation'}
          </button>
        </>
      }
    >
      <p className="body-md text-body">
        Your inventory will return to the pool immediately. The refund amount
        depends on how close the booking is to its scheduled time.
      </p>
      <label className="block mt-5">
        <span className="body-sm-strong">Reason (optional)</span>
        <textarea
          rows={3}
          className="input mt-2"
          value={reason}
          onChange={(e) => setReason(e.target.value)}
          placeholder="Anything we should know?"
        />
      </label>
    </Modal>
  );
}

function RescheduleDialog({ open, onClose, working, initial, onConfirm }) {
  const [iso, setIso] = useState(() => {
    if (!initial) return '';
    const d = new Date(initial);
    const yyyy = d.getFullYear();
    const mm = `${d.getMonth() + 1}`.padStart(2, '0');
    const dd = `${d.getDate()}`.padStart(2, '0');
    return `${yyyy}-${mm}-${dd}T${`${d.getHours()}`.padStart(2, '0')}:${`${d.getMinutes()}`.padStart(2, '0')}`;
  });

  useEffect(() => {
    if (!initial) return;
    const d = new Date(initial);
    const yyyy = d.getFullYear();
    const mm = `${d.getMonth() + 1}`.padStart(2, '0');
    const dd = `${d.getDate()}`.padStart(2, '0');
    setIso(
      `${yyyy}-${mm}-${dd}T${`${d.getHours()}`.padStart(2, '0')}:${`${d.getMinutes()}`.padStart(
        2,
        '0',
      )}`,
    );
  }, [initial, open]);

  return (
    <Modal
      open={open}
      onClose={onClose}
      title="Reschedule"
      footer={
        <>
          <button className="btn-secondary" onClick={onClose} disabled={working}>
            Cancel
          </button>
          <button
            className="btn-primary"
            onClick={() => iso && onConfirm(new Date(iso).toISOString())}
            disabled={working || !iso}
          >
            {working ? 'Saving…' : 'Save new time'}
          </button>
        </>
      }
    >
      <p className="body-md text-body">
        Pick the new departure or check-in time.
      </p>
      <input
        type="datetime-local"
        className="input mt-4"
        value={iso}
        onChange={(e) => setIso(e.target.value)}
      />
      <p className="mt-3 body-sm text-body">
        Pick a future time. We'll keep the rest of the booking the same.
      </p>
    </Modal>
  );
}

function DetailSkeleton() {
  return (
    <div className="container-x py-8 md:py-12">
      <Skeleton className="h-6 w-32 mb-4" />
      <Skeleton className="h-10 w-72 mb-6" />
      <div className="grid lg:grid-cols-[1fr_400px] gap-8">
        <Skeleton className="h-72 rounded-xl" />
        <Skeleton className="h-72 rounded-xl" />
      </div>
    </div>
  );
}
