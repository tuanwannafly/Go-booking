import { Link } from 'react-router-dom';
import {
  IconIllustrationCity,
  IconIllustrationDriver,
  IconIllustrationRider,
} from './Icons.jsx';

const categories = [
  { label: 'Reserve', icon: '📅' },
  { label: 'Rentals', icon: '🚗' },
  { label: 'Teens', icon: '🎓' },
  { label: 'Group rides', icon: '👥' },
];

export function CategoryRow() {
  return (
    <div className="container-x py-6 md:py-8">
      <div className="flex flex-wrap gap-2">
        {categories.map((c, i) => (
          <Link
            key={c.label}
            to={i === 0 ? '/hotels' : '#'}
            className="category-pill"
          >
            <span aria-hidden>{c.icon}</span>
            {c.label}
          </Link>
        ))}
      </div>
    </div>
  );
}

export default CategoryRow;

export function PromoBanner({ title, subtitle, ctaText, to, illustration }) {
  return (
    <div className="container-x py-8 md:py-12">
      <Link
        to={to}
        className="block rounded-xl bg-canvasSoft p-6 md:p-10 hover:bg-surfacePressed/70 transition-colors"
      >
        <div className="flex flex-col md:flex-row items-start md:items-center gap-6 md:gap-10">
          <div className="flex-1">
            <div className="display-xl">{title}</div>
            {subtitle ? (
              <div className="mt-3 body-md text-body max-w-prose">{subtitle}</div>
            ) : null}
            <div className="mt-5">
              <span className="btn-primary">{ctaText}</span>
            </div>
          </div>
          <div className="w-full md:w-72 lg:w-80 shrink-0">{illustration}</div>
        </div>
      </Link>
    </div>
  );
}

export function PromoLight({ title, body, ctaText, ctaTo, illustration, reverse }) {
  return (
    <div className="container-x py-8 md:py-12">
      <div
        className={`rounded-xl bg-canvas p-6 md:p-10 grid md:grid-cols-2 gap-6 md:gap-12 items-center ${
          reverse ? 'md:[&>:first-child]:order-2' : ''
        }`}
      >
        <div>
          <div className="display-lg">{title}</div>
          {body ? <p className="mt-3 body-md text-body max-w-prose">{body}</p> : null}
          {ctaText ? (
            <Link to={ctaTo || '#'} className="mt-5 inline-flex btn-primary">
              {ctaText}
            </Link>
          ) : null}
        </div>
        <div className="rounded-xl overflow-hidden">{illustration}</div>
      </div>
    </div>
  );
}

export function PromoDark({ title, body, ctaText, ctaTo }) {
  return (
    <div className="bg-ink text-onPrimary">
      <div className="container-x py-12 md:py-20 grid md:grid-cols-2 gap-8 md:gap-12 items-center">
        <div>
          <div className="display-xl">{title}</div>
          {body ? <p className="mt-3 body-lg text-mute max-w-prose">{body}</p> : null}
          {ctaText ? (
            <Link to={ctaTo || '#'} className="mt-6 inline-flex btn-secondary">
              {ctaText}
            </Link>
          ) : null}
        </div>
        <div className="rounded-xl overflow-hidden bg-blackElevated">
          <svg viewBox="0 0 320 240" xmlns="http://www.w3.org/2000/svg">
            <rect width="320" height="240" rx="16" fill="#282828" />
            <g fill="#fff">
              <circle cx="80" cy="120" r="22" />
              <rect x="60" y="150" width="40" height="48" rx="6" />
              <circle cx="240" cy="120" r="22" />
              <rect x="220" y="150" width="40" height="48" rx="6" />
            </g>
            <rect x="148" y="40" width="24" height="160" rx="12" fill="#fff" />
          </svg>
        </div>
      </div>
    </div>
  );
}

export function HeroBand({ children }) {
  return (
    <div className="bg-canvas">
      <div className="container-x py-8 md:py-16 grid lg:grid-cols-2 gap-10 items-center">
        {children}
      </div>
    </div>
  );
}

export { IconIllustrationCity, IconIllustrationDriver };
