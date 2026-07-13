// Auth context với mock backend integration.
// Token được lưu trong localStorage để duy trì session giữa các lần refresh.
// Backend verify token qua Authorization header (AuthMiddleware).

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';

const AuthContext = createContext(null);

const STORAGE_KEY = 'gobooking.user';
const TOKEN_KEY = 'gobooking.token';

// Mock user database - trong production sẽ replace bằng real backend call
const MOCK_USERS = new Map();

function loadInitial() {
  if (typeof window === 'undefined') return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

function persist(user) {
  if (typeof window === 'undefined') return;
  if (user) localStorage.setItem(STORAGE_KEY, JSON.stringify(user));
  else localStorage.removeItem(STORAGE_KEY);
}

function persistToken(token) {
  if (typeof window === 'undefined') return;
  if (token) localStorage.setItem(TOKEN_KEY, token);
  else localStorage.removeItem(TOKEN_KEY);
}

function generateToken() {
  return `token_${Date.now()}_${Math.random().toString(36).slice(2, 12)}`;
}

function validateEmail(email) {
  const trimmed = String(email || '').trim().toLowerCase();
  if (!trimmed) return { valid: false, error: 'Email là bắt buộc' };
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(trimmed)) return { valid: false, error: 'Email không hợp lệ' };
  return { valid: true, email: trimmed };
}

function validatePassword(password) {
  if (!password) return { valid: false, error: 'Mật khẩu là bắt buộc' };
  if (password.length < 6) return { valid: false, error: 'Mật khẩu phải có ít nhất 6 ký tự' };
  return { valid: true };
}

function validateName(name) {
  const trimmed = String(name || '').trim();
  if (!trimmed) return { valid: false, error: 'Tên là bắt buộc' };
  if (trimmed.length < 2) return { valid: false, error: 'Tên phải có ít nhất 2 ký tự' };
  if (trimmed.length > 50) return { valid: false, error: 'Tên không được quá 50 ký tự' };
  return { valid: true, name: trimmed };
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(loadInitial);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    persist(user);
  }, [user]);

  const signIn = useCallback(async ({ email, password }) => {
    // Validate inputs
    const emailResult = validateEmail(email);
    if (!emailResult.valid) {
      throw new Error(emailResult.error);
    }

    const passwordResult = validatePassword(password);
    if (!passwordResult.valid) {
      throw new Error(passwordResult.error);
    }

    setLoading(true);
    try {
      // Mock backend verification - trong production sẽ gọi API thực
      // Check nếu user đã tồn tại
      const storedUser = MOCK_USERS.get(emailResult.email);

      if (!storedUser) {
        throw new Error('Email chưa được đăng ký. Vui lòng đăng ký tài khoản mới.');
      }

      // Verify password
      if (storedUser.password !== password) {
        throw new Error('Mật khẩu không đúng. Vui lòng thử lại.');
      }

      // Generate token và login
      const token = generateToken();
      persistToken(token);

      const next = {
        id: storedUser.id,
        email: storedUser.email,
        name: storedUser.name,
      };
      setUser(next);
      return next;
    } finally {
      setLoading(false);
    }
  }, []);

  const signUp = useCallback(async ({ email, password, name }) => {
    // Validate inputs
    const emailResult = validateEmail(email);
    if (!emailResult.valid) {
      throw new Error(emailResult.error);
    }

    const passwordResult = validatePassword(password);
    if (!passwordResult.valid) {
      throw new Error(passwordResult.error);
    }

    const nameResult = validateName(name);
    if (!nameResult.valid) {
      throw new Error(nameResult.error);
    }

    setLoading(true);
    try {
      // Check nếu email đã được đăng ký
      if (MOCK_USERS.has(emailResult.email)) {
        throw new Error('Email này đã được đăng ký. Vui lòng đăng nhập.');
      }

      // Tạo user mới
      const newUser = {
        id: `user_${Date.now()}_${Math.random().toString(36).slice(2, 9)}`,
        email: emailResult.email,
        name: nameResult.name,
        password: password, // Lưu ý: trong production KHÔNG BAO GIỜ lưu plain password
      };
      MOCK_USERS.set(emailResult.email, newUser);

      // Generate token và auto login
      const token = generateToken();
      persistToken(token);

      const next = {
        id: newUser.id,
        email: newUser.email,
        name: newUser.name,
      };
      setUser(next);
      return next;
    } finally {
      setLoading(false);
    }
  }, []);

  const signOut = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    setUser(null);
  }, []);

  const getToken = useCallback(() => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(TOKEN_KEY);
  }, []);

  const value = useMemo(
    () => ({
      user,
      isAuthed: !!user,
      loading,
      signIn,
      signUp,
      signOut,
      getToken,
    }),
    [user, loading, signIn, signUp, signOut, getToken],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider');
  return ctx;
}
