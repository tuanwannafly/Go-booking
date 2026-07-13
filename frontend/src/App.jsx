import { Route, Routes, useLocation } from 'react-router-dom';
import { useEffect } from 'react';
import Navbar from './components/Navbar.jsx';
import Footer from './components/Footer.jsx';
import Home from './pages/Home.jsx';
import FlightsResults from './pages/FlightsResults.jsx';
import FlightDetail from './pages/FlightDetail.jsx';
import HotelsResults from './pages/HotelsResults.jsx';
import HotelDetail from './pages/HotelDetail.jsx';
import Bookings from './pages/Bookings.jsx';
import BookingDetail from './pages/BookingDetail.jsx';
import Login from './pages/Login.jsx';
import About from './pages/About.jsx';
import NotFound from './pages/NotFound.jsx';
import { AuthProvider } from './lib/auth.jsx';

function ScrollToTop() {
  const { pathname } = useLocation();
  useEffect(() => {
    window.scrollTo({ top: 0, behavior: 'instant' });
  }, [pathname]);
  return null;
}

export default function App() {
  return (
    <AuthProvider>
      <ScrollToTop />
      <div className="min-h-screen flex flex-col bg-canvas text-ink">
        <Navbar />
        <main className="flex-1 page-enter">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/flights" element={<FlightsResults />} />
            <Route path="/flights/:id" element={<FlightDetail />} />
            <Route path="/hotels" element={<HotelsResults />} />
            <Route path="/hotels/:id" element={<HotelDetail />} />
            <Route path="/bookings" element={<Bookings />} />
            <Route path="/bookings/:id" element={<BookingDetail />} />
            <Route path="/login" element={<Login />} />
            <Route path="/about" element={<About />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </main>
        <Footer />
      </div>
    </AuthProvider>
  );
}
