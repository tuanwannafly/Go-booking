// Lighter shared UI components — icon button, modal, empty state.

import { useEffect, useRef } from 'react';
import { IconClose } from './Icons.jsx';

export function IconButton({ className = '', children, ...props }) {
  return (
    <button
      type="button"
      className={`inline-flex h-10 w-10 items-center justify-center rounded-full bg-canvasSoft text-ink hover:bg-surfacePressed ring-focus ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}

export function Modal({ open, onClose, title, children, footer }) {
  const dialogRef = useRef(null);

  useEffect(() => {
    if (!open) return undefined;
    const onKey = (e) => {
      if (e.key === 'Escape') onClose?.();
    };
    document.addEventListener('keydown', onKey);
    document.body.style.overflow = 'hidden';
    return () => {
      document.removeEventListener('keydown', onKey);
      document.body.style.overflow = '';
    };
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-end md:items-center justify-center p-0 md:p-6 bg-ink/40">
      <div
        ref={dialogRef}
        role="dialog"
        aria-modal="true"
        className="w-full md:max-w-lg card-elevated rounded-t-xl md:rounded-xl bg-canvas text-ink"
      >
        <header className="flex items-center justify-between p-6 border-b border-surfacePressed">
          <h3 className="display-sm">{title}</h3>
          <button
            type="button"
            onClick={onClose}
            className="rounded-full p-2 hover:bg-canvasSoft ring-focus"
            aria-label="Close"
          >
            <IconClose className="h-5 w-5" />
          </button>
        </header>
        <div className="p-6 body-md">{children}</div>
        {footer ? (
          <footer className="px-6 py-4 border-t border-surfacePressed flex flex-col md:flex-row gap-3 md:justify-end">
            {footer}
          </footer>
        ) : null}
      </div>
    </div>
  );
}

export function EmptyState({ illustration, title, description, action }) {
  return (
    <div className="card-soft p-8 md:p-10 text-center flex flex-col items-center gap-4">
      {illustration ? (
        <div className="w-32 h-32 md:w-40 md:h-40 rounded-xl bg-canvas flex items-center justify-center">
          {illustration}
        </div>
      ) : null}
      <div className="flex flex-col gap-2 max-w-md">
        <div className="display-md">{title}</div>
        {description ? (
          <div className="body-md text-body">{description}</div>
        ) : null}
      </div>
      {action ? <div className="mt-2">{action}</div> : null}
    </div>
  );
}

export function Skeleton({ className = '', ...props }) {
  return (
    <div
      className={`animate-pulse bg-canvasSoft rounded-md ${className}`}
      {...props}
    />
  );
}

export function FieldLabel({ children, hint }) {
  return (
    <div className="flex items-center justify-between mb-1.5">
      <label className="body-sm-strong text-ink">{children}</label>
      {hint ? <span className="caption text-mute">{hint}</span> : null}
    </div>
  );
}

export function StatusPill({ status }) {
  const map = {
    pending: 'bg-canvasSoft text-ink',
    confirmed: 'bg-ink text-onPrimary',
    cancelled: 'bg-canvas text-ink border border-ink/20',
    expired: 'bg-canvas text-body border border-ink/10',
    completed: 'bg-ink text-onPrimary',
    failed: 'bg-canvas text-ink border border-ink/20',
    refunded: 'bg-canvas text-ink border border-ink/20',
  };
  const cls = map[status] || 'bg-canvas text-ink border border-ink/20';
  return (
    <span className={`inline-flex items-center gap-1 rounded-pill px-3 py-1 body-sm-strong ${cls}`}>
      <span className="capitalize">{status}</span>
    </span>
  );
}
