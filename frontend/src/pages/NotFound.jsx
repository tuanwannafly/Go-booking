import { Link } from 'react-router-dom';
import { IconArrowRight } from '../components/Icons.jsx';

export default function NotFound() {
  return (
    <div className="container-x py-12 md:py-20 text-center">
      <div className="caption uppercase tracking-[0.15em] text-mute mb-3">
        Error 404
      </div>
      <h1 className="display-xxl">We searched everywhere. Nothing here.</h1>
      <p className="mt-4 body-md text-body max-w-prose mx-auto">
        The URL you tried doesn't match any of our routes. Head back home and
        try a fresh search.
      </p>
      <div className="mt-8 flex justify-center gap-3">
        <Link to="/" className="btn-primary">
          Back to home <IconArrowRight className="h-4 w-4" />
        </Link>
        <Link to="/flights" className="btn-secondary">
          Search flights
        </Link>
      </div>
    </div>
  );
}
