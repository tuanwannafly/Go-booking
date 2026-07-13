import { Link } from 'react-router-dom';
import { IconLogoMark } from './Icons.jsx';

const columns = [
  {
    title: 'Company',
    links: [
      { label: 'About', to: '/about' },
      { label: 'Newsroom', to: '/about' },
      { label: 'Investors', to: '/about' },
      { label: 'Sustainability', to: '/about' },
      { label: 'Careers', to: '/about' },
    ],
  },
  {
    title: 'Products',
    links: [
      { label: 'Flights', to: '/flights' },
      { label: 'Stays', to: '/hotels' },
      { label: 'Bookings', to: '/bookings' },
      { label: 'Gift cards', to: '/about' },
    ],
  },
  {
    title: 'Support',
    links: [
      { label: 'Help center', to: '/about' },
      { label: 'Safety', to: '/about' },
      { label: 'Cancellation policy', to: '/about' },
      { label: 'Contact support', to: '/about' },
    ],
  },
  {
    title: 'Legal',
    links: [
      { label: 'Terms', to: '/about' },
      { label: 'Privacy', to: '/about' },
      { label: 'Cookies', to: '/about' },
      { label: 'Accessibility', to: '/about' },
    ],
  },
];

export default function Footer() {
  return (
    <footer className="bg-ink text-onPrimary">
      <div className="container-x py-12 md:py-16">
        <div className="grid grid-cols-2 gap-10 md:grid-cols-5">
          <div className="col-span-2 md:col-span-1 flex flex-col gap-4">
            <Link to="/" className="inline-flex items-center gap-2 text-onPrimary">
              <IconLogoMark className="h-7 w-7" />
              <span className="font-display text-[18px] font-bold tracking-tighter">
                GoBooking
              </span>
            </Link>
            <p className="body-sm text-mute max-w-[200px]">
              Travel and stays, simplified. One black-and-white promise.
            </p>
          </div>

          {columns.map((col) => (
            <div key={col.title} className="flex flex-col gap-3">
              <div className="body-md-strong text-onPrimary">{col.title}</div>
              <ul className="flex flex-col gap-2">
                {col.links.map((l) => (
                  <li key={l.label}>
                    <Link to={l.to} className="link-mute">
                      {l.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="mt-12 pt-8 border-t border-onPrimary/15 flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
          <div className="body-sm text-mute">
            © {new Date().getFullYear()} GoBooking Technologies Inc.
          </div>
          <div className="flex flex-wrap gap-3 caption text-mute">
            <span>Made with care in Saigon · Hanoi · Da Nang</span>
          </div>
          <div className="flex gap-3">
            <button type="button" className="app-pill text-sm">
              Download iOS
            </button>
            <button type="button" className="app-pill text-sm">
              Download Android
            </button>
          </div>
        </div>
      </div>
    </footer>
  );
}
