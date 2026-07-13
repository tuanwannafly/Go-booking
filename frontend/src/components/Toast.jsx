import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { IconCheck, IconClose, IconSpinner } from './Icons.jsx';

const ToastContext = createContext(null);

let pushId = 0;

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([]);

  const dismiss = useCallback((id) => {
    setToasts((list) => list.filter((t) => t.id !== id));
  }, []);

  const push = useCallback(
    (toast) => {
      pushId += 1;
      const id = pushId;
      const next = { id, kind: 'info', timeout: 3500, ...toast };
      setToasts((list) => [...list, next]);
      if (next.timeout > 0) {
        setTimeout(() => dismiss(id), next.timeout);
      }
      return id;
    },
    [dismiss],
  );

  return (
    <ToastContext.Provider
      value={{
        push,
        success: (msg, opts = {}) => push({ message: msg, kind: 'success', ...opts }),
        error: (msg, opts = {}) => push({ message: msg, kind: 'error', ...opts }),
        info: (msg, opts = {}) => push({ message: msg, kind: 'info', ...opts }),
        dismiss,
      }}
    >
      {children}
      <div className="pointer-events-none fixed top-4 right-4 z-[60] flex w-[min(360px,calc(100vw-2rem))] flex-col gap-2">
        {toasts.map((t) => (
          <ToastCard key={t.id} toast={t} onDismiss={() => dismiss(t.id)} />
        ))}
      </div>
    </ToastContext.Provider>
  );
}

function ToastCard({ toast, onDismiss }) {
  const kindStyle =
    toast.kind === 'success'
      ? 'bg-ink text-onPrimary'
      : toast.kind === 'error'
        ? 'bg-canvas text-ink border border-ink/10'
        : 'bg-canvas text-ink border border-ink/10';

  return (
    <div
      className={`pointer-events-auto shadow-l2 rounded-xl px-4 py-3 flex items-start gap-3 ${kindStyle}`}
      role="status"
    >
      <div className="pt-0.5">
        {toast.kind === 'success' ? (
          <IconCheck className="h-5 w-5" />
        ) : toast.kind === 'loading' ? (
          <IconSpinner className="h-5 w-5" />
        ) : (
          <IconClose className="h-5 w-5" />
        )}
      </div>
      <div className="flex-1 body-sm">{toast.message}</div>
      <button
        type="button"
        onClick={onDismiss}
        className="rounded-full p-1 hover:bg-canvasSoft/30 ring-focus"
        aria-label="Dismiss"
      >
        <IconClose className="h-4 w-4" />
      </button>
    </div>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) {
    throw new Error('useToast must be used inside ToastProvider');
  }
  return ctx;
}

// Avoid unused-import linter warning when only kind info is referenced.
void IconSpinner;
