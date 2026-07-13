// SVG icons used across the app. Single-color, currentColor driven.
// Stroke icons keep to 1.5 width to match the clean engineering-grade feel.

const baseProps = {
  width: 20,
  height: 20,
  viewBox: '0 0 24 24',
  fill: 'none',
  stroke: 'currentColor',
  strokeWidth: 1.6,
  strokeLinecap: 'round',
  strokeLinejoin: 'round',
};

export function IconPlane(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M21 16v-2l-8-5V3.5a1.5 1.5 0 1 0-3 0V9L2 14v2l8-2.5V19l-2 1.5V22l3.5-1L15 22v-1.5L13 19v-5.5L21 16Z" />
    </svg>
  );
}

export function IconHotel(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M3 20V8a2 2 0 0 1 2-2h6v4h10v10" />
      <path d="M3 20h18" />
      <path d="M9 12h.01M9 16h.01M14 12h.01M14 16h.01M19 12h.01M19 16h.01" />
    </svg>
  );
}

export function IconPin(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M12 21s7-6 7-11a7 7 0 1 0-14 0c0 5 7 11 7 11Z" />
      <circle cx="12" cy="10" r="2.5" />
    </svg>
  );
}

export function IconDot(props) {
  return (
    <svg {...baseProps} {...props}>
      <circle cx="12" cy="12" r="4" fill="currentColor" stroke="none" />
    </svg>
  );
}

export function IconCalendar(props) {
  return (
    <svg {...baseProps} {...props}>
      <rect x="3" y="5" width="18" height="16" rx="2" />
      <path d="M3 10h18M8 3v4M16 3v4" />
    </svg>
  );
}

export function IconClock(props) {
  return (
    <svg {...baseProps} {...props}>
      <circle cx="12" cy="12" r="9" />
      <path d="M12 7v5l3 2" />
    </svg>
  );
}

export function IconUser(props) {
  return (
    <svg {...baseProps} {...props}>
      <circle cx="12" cy="8" r="4" />
      <path d="M4 21a8 8 0 1 1 16 0" />
    </svg>
  );
}

export function IconArrowRight(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M5 12h14M13 6l6 6-6 6" />
    </svg>
  );
}

export function IconArrowLeft(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M19 12H5M11 18l-6-6 6-6" />
    </svg>
  );
}

export function IconSwap(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M7 7h12l-3-3M17 17H5l3 3" />
    </svg>
  );
}

export function IconCheck(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M5 12l4 4L19 7" />
    </svg>
  );
}

export function IconClose(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M6 6l12 12M18 6L6 18" />
    </svg>
  );
}

export function IconSearch(props) {
  return (
    <svg {...baseProps} {...props}>
      <circle cx="11" cy="11" r="7" />
      <path d="m20 20-3.5-3.5" />
    </svg>
  );
}

export function IconPlus(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M12 5v14M5 12h14" />
    </svg>
  );
}

export function IconMinus(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M5 12h14" />
    </svg>
  );
}

export function IconMenu(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M4 6h16M4 12h16M4 18h16" />
    </svg>
  );
}

export function IconChevronDown(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="m6 9 6 6 6-6" />
    </svg>
  );
}

export function IconChevronRight(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="m9 6 6 6-6 6" />
    </svg>
  );
}

export function IconBookmark(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M6 4h12v17l-6-4-6 4V4Z" />
    </svg>
  );
}

export function IconBag(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M4 7h16l-1 13H5L4 7Z" />
      <path d="M8 7V5a4 4 0 1 1 8 0v2" />
    </svg>
  );
}

export function IconReceipt(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M6 3v18l3-2 3 2 3-2 3 2V3H6Z" />
      <path d="M9 8h6M9 12h6M9 16h4" />
    </svg>
  );
}

export function IconShield(props) {
  return (
    <svg {...baseProps} {...props}>
      <path d="M12 3 5 6v6a9 9 0 0 0 7 9 9 9 0 0 0 7-9V6l-7-3Z" />
      <path d="m9 12 2 2 4-4" />
    </svg>
  );
}

export function IconGlobe(props) {
  return (
    <svg {...baseProps} {...props}>
      <circle cx="12" cy="12" r="9" />
      <path d="M3 12h18M12 3a14 14 0 0 1 0 18M12 3a14 14 0 0 0 0 18" />
    </svg>
  );
}

export function IconLogoMark(props) {
  // Cube / luggage emblem used as the brand mark in nav and footer
  return (
    <svg
      width={26}
      height={26}
      viewBox="0 0 32 32"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect x="2" y="6" width="28" height="20" rx="6" fill="currentColor" />
      <path
        d="M10 6V4h12v2"
        stroke="#fff"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <circle cx="16" cy="16" r="3" fill="#fff" />
      <path
        d="M11 22h10"
        stroke="#fff"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  );
}

export function IconIllustrationCity(props) {
  return (
    <svg viewBox="0 0 320 240" xmlns="http://www.w3.org/2000/svg" {...props}>
      <rect width="320" height="240" rx="16" fill="#efefef" />
      <g fill="#000">
        <rect x="32" y="120" width="40" height="100" />
        <rect x="80" y="80" width="48" height="140" />
        <rect x="136" y="100" width="36" height="120" />
        <rect x="180" y="60" width="56" height="160" />
        <rect x="244" y="120" width="44" height="100" />
      </g>
      <g fill="#fff">
        <rect x="38" y="130" width="6" height="6" />
        <rect x="50" y="130" width="6" height="6" />
        <rect x="38" y="146" width="6" height="6" />
        <rect x="50" y="146" width="6" height="6" />
        <rect x="38" y="162" width="6" height="6" />
        <rect x="50" y="162" width="6" height="6" />
        <rect x="86" y="90" width="6" height="6" />
        <rect x="100" y="90" width="6" height="6" />
        <rect x="114" y="90" width="6" height="6" />
        <rect x="86" y="106" width="6" height="6" />
        <rect x="100" y="106" width="6" height="6" />
        <rect x="86" y="122" width="6" height="6" />
        <rect x="100" y="122" width="6" height="6" />
        <rect x="114" y="122" width="6" height="6" />
        <rect x="142" y="110" width="6" height="6" />
        <rect x="156" y="110" width="6" height="6" />
        <rect x="142" y="126" width="6" height="6" />
        <rect x="156" y="126" width="6" height="6" />
        <rect x="186" y="72" width="6" height="6" />
        <rect x="200" y="72" width="6" height="6" />
        <rect x="214" y="72" width="6" height="6" />
        <rect x="186" y="88" width="6" height="6" />
        <rect x="200" y="88" width="6" height="6" />
        <rect x="214" y="88" width="6" height="6" />
        <rect x="250" y="132" width="6" height="6" />
        <rect x="264" y="132" width="6" height="6" />
        <rect x="278" y="132" width="6" height="6" />
      </g>
      <circle cx="280" cy="40" r="14" fill="#000" />
      <circle cx="40" cy="30" r="6" fill="#000" />
    </svg>
  );
}

export function IconIllustrationDriver(props) {
  return (
    <svg viewBox="0 0 320 240" xmlns="http://www.w3.org/2000/svg" {...props}>
      <rect width="320" height="240" rx="16" fill="#efefef" />
      <g fill="#000">
        <rect x="40" y="80" width="240" height="120" rx="14" />
      </g>
      <g fill="#fff">
        <rect x="76" y="108" width="48" height="36" rx="4" />
        <rect x="196" y="108" width="48" height="36" rx="4" />
        <rect x="132" y="116" width="56" height="20" rx="10" />
      </g>
      <circle cx="120" cy="80" r="24" fill="#000" />
      <rect x="220" y="60" width="40" height="6" rx="3" fill="#000" />
      <rect x="220" y="72" width="40" height="6" rx="3" fill="#000" />
      <rect x="220" y="84" width="40" height="6" rx="3" fill="#000" />
      <rect x="20" y="20" width="14" height="14" rx="2" fill="#000" />
      <rect x="286" y="180" width="14" height="14" rx="2" fill="#000" />
    </svg>
  );
}

export function IconIllustrationRider(props) {
  return (
    <svg viewBox="0 0 320 240" xmlns="http://www.w3.org/2000/svg" {...props}>
      <rect width="320" height="240" rx="16" fill="#efefef" />
      <circle cx="80" cy="80" r="20" fill="#000" />
      <path
        d="M40 180c0-22 18-40 40-40s40 18 40 40v20H40v-20Z"
        fill="#000"
      />
      <circle cx="240" cy="80" r="20" fill="#000" />
      <path
        d="M200 180c0-22 18-40 40-40s40 18 40 40v20H200v-20Z"
        fill="#000"
      />
      <rect x="20" y="200" width="280" height="6" fill="#000" />
      <rect x="140" y="40" width="40" height="40" rx="20" fill="#000" />
    </svg>
  );
}

export function IconIllustrationShield(props) {
  return (
    <svg viewBox="0 0 320 240" xmlns="http://www.w3.org/2000/svg" {...props}>
      <rect width="320" height="240" rx="16" fill="#efefef" />
      <path
        d="M160 30 90 60v60c0 36 30 70 70 90 40-20 70-54 70-90V60l-70-30Z"
        fill="#000"
      />
      <path
        d="m130 122 22 22 38-44"
        fill="none"
        stroke="#fff"
        strokeWidth="10"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

export function IconIllustrationBag(props) {
  return (
    <svg viewBox="0 0 320 240" xmlns="http://www.w3.org/2000/svg" {...props}>
      <rect width="320" height="240" rx="16" fill="#efefef" />
      <rect x="70" y="80" width="180" height="120" rx="14" fill="#000" />
      <path
        d="M110 80V58a30 30 0 1 1 60 0v22"
        stroke="#000"
        strokeWidth="10"
        strokeLinecap="round"
        fill="none"
      />
      <rect x="130" y="140" width="60" height="6" rx="3" fill="#fff" />
    </svg>
  );
}

export function IconIllustrationGlobe(props) {
  return (
    <svg viewBox="0 0 320 240" xmlns="http://www.w3.org/2000/svg" {...props}>
      <rect width="320" height="240" rx="16" fill="#efefef" />
      <circle cx="160" cy="120" r="80" fill="#000" />
      <ellipse
        cx="160"
        cy="120"
        rx="36"
        ry="80"
        fill="none"
        stroke="#fff"
        strokeWidth="6"
      />
      <ellipse
        cx="160"
        cy="120"
        rx="80"
        ry="36"
        fill="none"
        stroke="#fff"
        strokeWidth="6"
      />
      <path
        d="M80 120h160"
        stroke="#fff"
        strokeWidth="6"
      />
    </svg>
  );
}

export function IconEmptySeat(props) {
  return (
    <svg viewBox="0 0 24 24" {...props}>
      <rect
        x="5"
        y="3"
        width="14"
        height="18"
        rx="2"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.4"
      />
      <path
        d="M9 7h6M9 11h6"
        stroke="currentColor"
        strokeWidth="1.4"
        strokeLinecap="round"
      />
    </svg>
  );
}

export function IconSpinner(props) {
  return (
    <svg
      {...baseProps}
      {...props}
      className={`animate-spin ${props.className || ''}`}
    >
      <path d="M21 12a9 9 0 1 1-6.2-8.55" />
    </svg>
  );
}
