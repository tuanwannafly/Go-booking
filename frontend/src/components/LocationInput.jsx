// Autocomplete input used by the request form for origin / destination / city.
// Renders a floating list of suggestions below the field so the user can pick
// from the seeded dataset instead of typing blind guesses.

import { useEffect, useMemo, useRef, useState } from 'react';
import { IconPin } from './Icons.jsx';

export default function LocationInput({
  value,
  onChange,
  onSelect,
  placeholder,
  suggestions,
  // Used for `aria-label` and inside the suggestion rows. The label describes
  // what kind of location this field expects: e.g. "From" or "To".
  label,
  required,
}) {
  const [open, setOpen] = useState(false);
  const [active, setActive] = useState(0);
  const wrapperRef = useRef(null);
  const inputRef = useRef(null);

  const filtered = useMemo(() => {
    const q = String(value || '').trim().toLowerCase();
    const list = suggestions || [];
    if (!q) return list;
    return list.filter(
      (s) =>
        s.code.toLowerCase().includes(q) ||
        s.name.toLowerCase().includes(q) ||
        s.city.toLowerCase().includes(q),
    );
  }, [value, suggestions]);

  useEffect(() => {
    const onClickAway = (e) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', onClickAway);
    return () => document.removeEventListener('mousedown', onClickAway);
  }, []);

  const pick = (item) => {
    onSelect?.(item);
    setOpen(false);
    inputRef.current?.focus();
  };

  const onKeyDown = (e) => {
    if (!open || filtered.length === 0) {
      if (e.key === 'ArrowDown') setOpen(true);
      return;
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActive((i) => Math.min(filtered.length - 1, i + 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActive((i) => Math.max(0, i - 1));
    } else if (e.key === 'Enter') {
      e.preventDefault();
      const item = filtered[active];
      if (item) pick(item);
    } else if (e.key === 'Escape') {
      setOpen(false);
    }
  };

  return (
    <div ref={wrapperRef} className="relative">
      <div className="request-form-row">
        <IconPin className="h-5 w-5 text-ink/70 shrink-0" />
        <input
          ref={inputRef}
          required={required}
          value={value}
          onChange={(e) => {
            onChange(e.target.value);
            setOpen(true);
            setActive(0);
          }}
          onFocus={() => setOpen(true)}
          onKeyDown={onKeyDown}
          placeholder={placeholder}
          aria-label={label}
          aria-autocomplete="list"
          aria-expanded={open}
          aria-controls={`${label}-listbox`}
          className="bg-transparent flex-1 outline-none body-md placeholder:text-mute"
        />
      </div>
      {open && filtered.length > 0 ? (
        <ul
          id={`${label}-listbox`}
          role="listbox"
          className="absolute z-30 mt-2 w-full bg-canvas border border-surfacePressed rounded-md shadow-l1 max-h-64 overflow-auto"
        >
          {filtered.map((item, i) => {
            const selected = i === active;
            return (
              <li
                key={`${item.code}-${i}`}
                role="option"
                aria-selected={selected}
                onMouseEnter={() => setActive(i)}
                onMouseDown={(e) => {
                  // mousedown so the input keeps focus
                  e.preventDefault();
                  pick(item);
                }}
                className={`flex items-center gap-3 px-3 py-2 cursor-pointer body-sm ${
                  selected ? 'bg-canvasSoft' : 'bg-canvas hover:bg-canvasSoft'
                }`}
              >
                <span className="rounded-pill bg-canvasSoft px-2 py-0.5 body-sm-strong tabular-nums min-w-[3rem] text-center">
                  {item.code}
                </span>
                <span className="flex-1 truncate">{item.name}</span>
                <span className="caption text-mute">{item.city}</span>
              </li>
            );
          })}
        </ul>
      ) : null}
      {open && filtered.length === 0 && value ? (
        <div
          role="status"
          className="absolute z-30 mt-2 w-full bg-canvas border border-surfacePressed rounded-md px-3 py-2 body-sm text-mute"
        >
          Không có địa điểm phù hợp. Hãy chọn từ gợi ý.
        </div>
      ) : null}
    </div>
  );
}
