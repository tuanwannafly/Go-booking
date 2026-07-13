import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../lib/auth.jsx';
import { useToast } from '../components/Toast.jsx';
import { IconLogoMark, IconArrowRight } from '../components/Icons.jsx';
import { FieldLabel } from '../components/Primitives.jsx';

export default function Login() {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { signIn, signUp, loading } = useAuth();
  const toast = useToast();

  const initialMode = params.get('signup') ? 'signup' : 'login';
  const next = params.get('next') || '/bookings';

  const [mode, setMode] = useState(initialMode);
  const [email, setEmail] = useState('');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const submit = async (e) => {
    e.preventDefault();
    try {
      if (mode === 'signup') {
        await signUp({ email, name, password });
        toast.success('Đăng ký thành công! Chào mừng đến với GoBooking.');
      } else {
        await signIn({ email, password });
        toast.success('Đăng nhập thành công!');
      }
      navigate(next);
    } catch (err) {
      toast.error(err.message || 'Không thể đăng nhập');
    }
  };

  const switchMode = (newMode) => {
    setMode(newMode);
    setPassword(''); // Clear password when switching modes
  };

  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-12 grid lg:grid-cols-2 gap-10 items-start">
        <div className="hidden lg:flex flex-col justify-between min-h-[520px]">
          <div>
            <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
              GoBooking
            </div>
            <h1 className="display-xl max-w-[16ch]">
              {mode === 'signup'
                ? 'Tạo tài khoản để đặt chỗ dễ dàng hơn.'
                : 'Đăng nhập để quản lý đặt chỗ của bạn.'}
            </h1>
            <p className="mt-4 body-md text-body max-w-prose">
              {mode === 'signup'
                ? 'Tạo tài khoản miễn phí để giữ ghế, xác nhận đặt chỗ và thay đổi lịch trình.'
                : 'Truy cập vào tài khoản để xem tất cả đặt chỗ, giữ ghế trong 10 phút và xác nhận thanh toán.'}
            </p>
          </div>
          <div className="card-soft p-6 max-w-sm mt-10">
            <div className="caption uppercase tracking-[0.15em] text-mute mb-2">
              Lợi ích khi đăng nhập
            </div>
            <ul className="grid gap-2 body-md">
              <li className="flex items-start gap-2">
                <span className="text-ink mt-0.5">•</span>
                Giữ ghế và phòng trong 10 phút
              </li>
              <li className="flex items-start gap-2">
                <span className="text-ink mt-0.5">•</span>
                Xem và quản lý tất cả đặt chỗ
              </li>
              <li className="flex items-start gap-2">
                <span className="text-ink mt-0.5">•</span>
                Thay đổi hoặc hủy đặt chỗ dễ dàng
              </li>
            </ul>
          </div>
        </div>

        <div className="card-elevated p-6 md:p-8 max-w-lg w-full justify-self-end w-full">
          <div className="flex items-center gap-2 mb-6">
            <IconLogoMark className="h-7 w-7" />
            <span className="display-sm">
              {mode === 'signup' ? 'Tạo tài khoản mới' : 'Đăng nhập GoBooking'}
            </span>
          </div>

          {/* Mode toggle tabs */}
          <div className="flex items-center gap-2 mb-6">
            <button
              type="button"
              onClick={() => switchMode('login')}
              className={`btn-tab-translucent ${mode === 'login' ? 'shadow-l3' : ''}`}
              aria-pressed={mode === 'login'}
            >
              Đăng nhập
            </button>
            <button
              type="button"
              onClick={() => switchMode('signup')}
              className={`btn-tab-translucent ${mode === 'signup' ? 'shadow-l3' : ''}`}
              aria-pressed={mode === 'signup'}
            >
              Đăng ký
            </button>
          </div>

          <form onSubmit={submit} className="grid gap-4">
            {/* Email field */}
            <div>
              <FieldLabel>Email</FieldLabel>
              <input
                type="email"
                required
                autoFocus
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input"
                placeholder="email@example.com"
              />
            </div>

            {/* Name field - only for signup */}
            {mode === 'signup' && (
              <div>
                <FieldLabel>Họ và tên</FieldLabel>
                <input
                  type="text"
                  required
                  autoComplete="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="input"
                  placeholder="Nhập họ và tên của bạn"
                />
              </div>
            )}

            {/* Password field */}
            <div>
              <FieldLabel>Mật khẩu</FieldLabel>
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  required
                  autoComplete={mode === 'signup' ? 'new-password' : 'current-password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="input pr-12"
                  placeholder={mode === 'signup' ? 'Tạo mật khẩu (ít nhất 6 ký tự)' : 'Nhập mật khẩu'}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-mute hover:text-ink transition-colors"
                  aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
                >
                  {showPassword ? (
                    <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                    </svg>
                  ) : (
                    <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  )}
                </button>
              </div>
            </div>

            {/* Submit button */}
            <button
              type="submit"
              className="btn-primary mt-2 w-full"
              disabled={loading}
            >
              {loading ? 'Đang xử lý…' : mode === 'signup' ? 'Tạo tài khoản' : 'Đăng nhập'}
              {!loading && <IconArrowRight className="h-4 w-4" />}
            </button>
          </form>

          <div className="mt-6 h-divider" />
          <div className="mt-6 body-sm text-body">
            {mode === 'signup' ? (
              <p>
                Bằng việc đăng ký, bạn đồng ý với{' '}
                <Link to="/about" className="link-blue">
                  Điều khoản sử dụng
                </Link>{' '}
                và{' '}
                <Link to="/about" className="link-blue">
                  Chính sách bảo mật
                </Link>{' '}
                của GoBooking.
              </p>
            ) : (
              <p>
                Chưa có tài khoản?{' '}
                <button
                  type="button"
                  onClick={() => switchMode('signup')}
                  className="link-blue"
                >
                  Đăng ký ngay
                </button>
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
