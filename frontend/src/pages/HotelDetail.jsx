import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { hotels, holds, bookings, ApiError } from '../lib/api.js';
import {
  IconArrowRight,
  IconCheck,
  IconClock,
  IconIllustrationCity,
  IconPin,
  IconUser,
} from '../components/Icons.jsx';
import { EmptyState, Skeleton, StatusPill } from '../components/Primitives.jsx';
import { useAuth } from '../lib/auth.jsx';
import { useToast } from '../components/Toast.jsx';
import { formatDate, formatVnd, nightsLabel } from '../lib/format.js';

export default function HotelDetail() {
  const { id } = useParams();
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { isAuthed } = useAuth();
  const toast = useToast();

  const checkin = params.get('checkin') || '';
  const checkout = params.get('checkout') || '';
  const rooms = Number(params.get('rooms') || '1');

  const [hotel, setHotel] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedType, setSelectedType] = useState(null);
  const [heldRooms, setHeldRooms] = useState([]);
  const [holding, setHolding] = useState(false);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    hotels
      .get(id, { checkin, checkout })
      .then((h) => !cancelled && setHotel(h))
      .catch((err) => {
        if (cancelled) return;
        if (err instanceof ApiError) setError(err.message);
        else setError('Could not load this hotel');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  }, [id, checkin, checkout]);

  if (loading) return <DetailSkeleton />;
  if (error || !hotel) {
    return (
      <div className="container-x py-12">
      <EmptyState
        title="Không tải được thông tin khách sạn"
        description={error || 'Hãy thử lựa chọn khác.'}
          action={
            <Link to="/hotels" className="btn-primary">
              Back to results
            </Link>
          }
        />
      </div>
    );
  }

  const nights = nightsLabel(checkin, checkout);
  const roomTypes = hotel.room_types || [];

  const onHold = async () => {
    if (!selectedType) {
      toast.info('Vui lòng chọn loại phòng trước');
      return;
    }
    setHolding(true);
    try {
      const heldUnits = [];
      for (let i = 0; i < Math.max(1, rooms); i += 1) {
        // eslint-disable-next-line no-await-in-loop
        const res = await holds.room(selectedType.id || '00000000-0000-0000-0000-000000000000');
        heldUnits.push(res);
      }
      setHeldRooms(heldUnits);
      toast.success('Đã giữ phòng trong 10 phút');
    } catch (err) {
      toast.error(err.message || 'Không thể giữ phòng này');
    } finally {
      setHolding(false);
    }
  };

  const onConfirm = async () => {
    if (!isAuthed) {
      toast.info('Vui lòng đăng nhập để xác nhận');
      navigate(`/login?next=${encodeURIComponent(window.location.pathname + window.location.search)}`);
      return;
    }
    if (heldRooms.length === 0) {
      toast.info('Vui lòng giữ phòng trước');
      return;
    }
    try {
      const items = heldRooms.map((r) => ({
        inventory_unit_id: r.inventory_unit?.id || r.inventory_unit_id || '00000000-0000-0000-0000-000000000000',
      }));
      const created = await bookings.create(items);
      toast.success('Đã tạo đặt chỗ. Đang chuyển đến trang thanh toán…');
      navigate(`/bookings/${created.id}`);
    } catch (err) {
      toast.error(err.message || 'Không thể tạo đặt chỗ');
    }
  };

  const total = selectedType && nights ? selectedType.base_price * nights * rooms : 0;

  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-12">
        <Link to="/hotels" className="body-sm text-body hover:underline">
          ← Back to results
        </Link>

        <header className="mt-4 grid lg:grid-cols-[1fr_400px] gap-8">
          <div>
            <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
              {hotel.city}
            </div>
            <h1 className="display-xl">{hotel.name}</h1>
            <div className="mt-3 body-md text-body inline-flex items-center gap-2">
              <IconPin className="h-4 w-4" /> {hotel.address}
            </div>

            <div className="mt-6 rounded-xl overflow-hidden">
              <HotelHero city={hotel.city} />
            </div>
          </div>
          <aside>
            <div className="card-elevated p-6 flex flex-col gap-4 sticky top-24">
              <div className="body-md-strong">Thông tin phòng</div>
              <div className="grid grid-cols-2 gap-3 body-md">
                <SummaryRow label="Nhận phòng" value={formatDate(checkin)} />
                <SummaryRow label="Trả phòng" value={formatDate(checkout)} />
                <SummaryRow label="Số đêm" value={nights} />
                <SummaryRow label="Số phòng" value={rooms} />
              </div>
              <div className="h-divider" />
              <div className="flex items-baseline justify-between">
                <span className="body-md-strong">Tổng cộng</span>
                <span className="display-lg tabular-nums">{formatVnd(total)}</span>
              </div>
              <div className="caption text-mute">
                {selectedType
                  ? `${rooms} × ${formatVnd(selectedType.base_price)} × ${nights} đêm`
                  : 'Chọn phòng để xem tổng'}
              </div>
              <button
                type="button"
                onClick={onConfirm}
                disabled={heldRooms.length === 0}
                className="btn-primary w-full"
              >
                {heldRooms.length > 0 ? 'Tiếp tục xác nhận' : 'Giữ phòng trước'}
              </button>
              <button
                type="button"
                onClick={() => navigate(-1)}
                className="btn-secondary w-full"
              >
                Quay lại
              </button>
            </div>
          </aside>
        </header>
      </div>

      <div className="container-x pb-12 md:pb-20">
        <section>
          <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
            Bước 2/3 · Chọn phòng
          </div>
          <h2 className="display-lg mb-6">Chọn loại phòng</h2>
          <div className="grid md:grid-cols-2 gap-5">
            {roomTypes.map((rt) => {
              const isSelected = selectedType?.id === rt.id;
              const summary = `${nights} đêm × ${rooms} phòng`;
              const subtotal = rt.base_price * nights * rooms;
              return (
                <button
                  key={rt.id}
                  type="button"
                  onClick={() => {
                    setSelectedType(rt);
                    setHeldRooms([]);
                  }}
                  className={`text-left card-elevated p-6 transition-shadow hover:shadow-l2 ${
                    isSelected ? 'ring-2 ring-ink' : ''
                  }`}
                  aria-pressed={isSelected}
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="display-md">{rt.name}</div>
                      <div className="caption text-mute mt-1">
                        Tối đa {rt.capacity || 2} người · Hủy miễn phí (24h)
                      </div>
                    </div>
                    <div className="flex flex-col items-end">
                      <span className="display-md tabular-nums">
                        {formatVnd(rt.base_price)}
                      </span>
                      <span className="caption text-mute">mỗi đêm</span>
                    </div>
                  </div>
                  <div className="mt-4 grid grid-cols-2 gap-3 body-md">
                    <Bullet>Wi-Fi miễn phí</Bullet>
                    <Bullet>Bữa sáng bao gồm</Bullet>
                    <Bullet>1 giường king</Bullet>
                    <Bullet>Tầm nhìn thành phố</Bullet>
                  </div>
                  <div className="mt-5 flex items-center justify-between">
                    <div className="caption text-mute">{summary}</div>
                    <div className="body-md-strong tabular-nums">
                      {formatVnd(subtotal)}
                    </div>
                  </div>
                </button>
              );
            })}
          </div>

          {selectedType ? (
            <div className="mt-8 card-soft p-5 flex flex-col md:flex-row md:items-center gap-3">
              <div className="flex-1">
                <div className="body-md-strong">
                  {selectedType.name} × {rooms}
                </div>
                <div className="caption text-mute">
                  Giữ 10 phút, xác nhận để khóa giá.
                </div>
              </div>
              <button
                type="button"
                onClick={onHold}
                disabled={holding}
                className="btn-primary"
              >
                {holding ? 'Đang giữ…' : 'Giữ phòng này'}
                {!holding ? <IconArrowRight className="h-4 w-4" /> : null}
              </button>
            </div>
          ) : null}

          {heldRooms.length > 0 ? (
            <div className="mt-4 card-elevated p-5 flex items-center gap-3">
              <IconClock className="h-5 w-5" />
              <div className="flex-1">
                <div className="body-md-strong">{heldRooms.length} phòng đang được giữ</div>
                <div className="caption text-body">
                  Giữ đến {formatDate(heldRooms[0].held_until, { weekday: 'short' })} · {heldRooms[0].held_until?.slice(11, 16) || 'sớm'}
                </div>
              </div>
              <StatusPill status="pending" />
            </div>
          ) : null}
        </section>

        <section className="mt-12 md:mt-16">
          <h2 className="display-lg mb-6">Về chỗ ở này</h2>
          <div className="grid md:grid-cols-3 gap-5">
            <AboutCard
              icon={<IconCheck className="h-5 w-5" />}
              title="Hủy miễn phí"
              body="Hủy trước 24 giờ để nhận hoàn tiền đầy đủ."
            />
            <AboutCard
              icon={<IconUser className="h-5 w-5" />}
              title="Hỗ trợ 24/7"
              body="Liên hệ hỗ trợ qua chat hoặc điện thoại trong 3 phút."
            />
            <AboutCard
              icon={<IconPin className="h-5 w-5" />}
              title="Giá tốt nhất"
              body="Chúng tôi hoàn tiền nếu bạn tìm thấy giá thấp hơn trong 24 giờ."
            />
          </div>
        </section>
      </div>
    </div>
  );
}

function Bullet({ children }) {
  return (
    <div className="flex items-center gap-2 body-sm">
      <IconCheck className="h-4 w-4 text-ink" />
      {children}
    </div>
  );
}

function AboutCard({ icon, title, body }) {
  return (
    <div className="card p-5">
      <div className="rounded-full bg-canvasSoft w-10 h-10 flex items-center justify-center mb-3">
        {icon}
      </div>
      <div className="display-sm">{title}</div>
      <div className="mt-2 body-md text-body">{body}</div>
    </div>
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

function HotelHero({ city }) {
  return (
    <div className="aspect-[16/9] bg-canvasSoft rounded-xl overflow-hidden">
      <IconIllustrationCity className="w-full h-full" />
    </div>
  );
}

function DetailSkeleton() {
  return (
    <div className="container-x py-8 md:py-12">
      <Skeleton className="h-6 w-32 mb-4" />
      <Skeleton className="h-10 w-72 mb-6" />
      <Skeleton className="h-72 w-full rounded-xl mb-6" />
    </div>
  );
}
