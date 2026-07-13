import { Link, NavLink, useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import {
  IconClose,
  IconLogoMark,
  IconMenu,
  IconPlane,
  IconHotel,
  IconReceipt,
  IconUser,
} from './Icons.jsx';
import { useAuth } from '../lib/auth.jsx';

const links = [
  { to: '/flights', label: 'Flights', icon: IconPlane },
  { to: '/hotels', label: 'Stays', icon: IconHotel },
  { to: '/bookings', label: 'Bookings', icon: IconReceipt },
];

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const { user, loading, signOut } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    document.body.style.overflow = open ? 'hidden' : '';
    return () => {
      document.body.style.overflow = '';
    };
  }, [open]);

  const onSignOut = () => {
    signOut();
    setOpen(false);
    navigate('/');
  };

  return (
    <header className="sticky top-0 z-40 bg-canvas/95 backdrop-blur border-b border-surfacePressed/60">
      <div className="container-x flex h-16 items-center gap-8">
        <Link to="/" className="flex items-center gap-2 text-ink" aria-label="GoBooking home">
          <IconLogoMark className="h-7 w-7" />
          <span className="font-display text-[18px] font-bold tracking-tighter">
            GoBooking
          </span>
        </Link>

        <nav className="hidden md:flex items-center gap-7" aria-label="Primary">
          {links.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              className={({ isActive }) =>
                `nav-link inline-flex items-center gap-1.5 ${
                  isActive ? 'text-ink' : 'text-ink/80'
                }`
              }
            >
              <span>{link.label}</span>
            </NavLink>
          ))}
        </nav>

        <div className="ml-auto flex items-center gap-2">
          <Link to="/about" className="hidden md:inline-flex nav-link">
            About
          </Link>

          {user ? (
            <div className="hidden md:flex items-center gap-2">
              <Link to="/bookings" className="nav-link inline-flex items-center gap-2">
                <IconUser className="h-5 w-5" />
                <span>{user.email?.split('@')[0] || 'Account'}</span>
              </Link>
              <button onClick={onSignOut} className="btn-subtle">
                Đăng xuất
              </button>
            </div>
          ) : loading ? (
            <div className="hidden md:flex items-center gap-2">
              <div className="h-9 w-24 rounded-full bg-canvasSoft animate-pulse" />
            </div>
          ) : (
            <div className="hidden md:flex items-center gap-2">
              <Link to="/login" className="btn-secondary">
                Đăng nhập
              </Link>
              <Link to="/login?signup=1" className="btn-primary">
                Đăng ký
              </Link>
            </div>
          )}

          <button
            type="button"
            className="md:hidden icon-button-circular"
            aria-label={open ? 'Close menu' : 'Open menu'}
            onClick={() => setOpen((v) => !v)}
          >
            {open ? <IconClose className="h-5 w-5" /> : <IconMenu className="h-5 w-5" />}
          </button>
        </div>
      </div>

      {open ? (
        <div className="md:hidden absolute inset-x-0 top-16 bg-canvas border-b border-surfacePressed">
          <div className="container-x py-6 flex flex-col gap-5">
            <nav className="flex flex-col gap-4" aria-label="Mobile">
              {links.map((link) => (
                <Link
                  key={link.to}
                  to={link.to}
                  onClick={() => setOpen(false)}
                  className="display-sm flex items-center gap-3"
                >
                  <link.icon className="h-5 w-5" />
                  {link.label}
                </Link>
              ))}
              <Link
                to="/about"
                onClick={() => setOpen(false)}
                className="display-sm flex items-center gap-3"
              >
                About
              </Link>
            </nav>

            <div className="h-px bg-surfacePressed" />

            {user ? (
              <div className="flex flex-col gap-3">
                <div className="body-sm text-body">Đã đăng nhập: {user.email}</div>
                <button onClick={onSignOut} className="btn-primary w-full">
                  Đăng xuất
                </button>
              </div>
            ) : loading ? (
              <div className="flex flex-col gap-3">
                <div className="h-10 w-full rounded-full bg-canvasSoft animate-pulse" />
              </div>
            ) : (
              <div className="flex flex-col gap-3">
                <Link
                  to="/login"
                  onClick={() => setOpen(false)}
                  className="btn-secondary w-full"
                >
                  Đăng nhập
                </Link>
                <Link
                  to="/login?signup=1"
                  onClick={() => setOpen(false)}
                  className="btn-primary w-full"
                >
                  Đăng ký
                </Link>
              </div>
            )}
          </div>
        </div>
      ) : null}
    </header>
  );
}
