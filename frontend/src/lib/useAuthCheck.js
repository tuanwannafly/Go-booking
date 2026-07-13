// Custom hook cho việc bảo vệ các trang yêu cầu đăng nhập.
// Sử dụng hook này thay vì tự implement logic redirect trong từng component.

import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from './auth.jsx';
import { useToast } from '../components/Toast.jsx';

/**
 * Hook để bảo vệ các trang yêu cầu đăng nhập.
 * Tự động redirect về trang login với `next` param nếu chưa đăng nhập.
 *
 * @param {Object} options
 * @param {string} options.redirectTo - URL để redirect sau khi login thành công (mặc định: trang hiện tại)
 * @param {boolean} options.showToast - Có hiện toast thông báo không (mặc định: true)
 * @param {string} options.toastMessage - Message cho toast (mặc định: 'Vui lòng đăng nhập để tiếp tục')
 * @returns {Object} - { isReady: boolean } để biết khi nào auth state đã sẵn sàng
 */
export function useRequireAuth({
  redirectTo,
  showToast = true,
  toastMessage = 'Vui lòng đăng nhập để tiếp tục',
} = {}) {
  const { isAuthed, loading } = useAuth();
  const navigate = useNavigate();
  const toast = useToast();

  useEffect(() => {
    // Đợi cho auth state load xong
    if (loading) return;

    // Nếu chưa đăng nhập, redirect về login
    if (!isAuthed) {
      if (showToast) {
        toast.info(toastMessage);
      }

      const currentPath = window.location.pathname + window.location.search;
      const nextPath = redirectTo || currentPath;
      const loginPath = nextPath !== '/' ? `/login?next=${encodeURIComponent(nextPath)}` : '/login';

      navigate(loginPath);
    }
  }, [isAuthed, loading, navigate, showToast, toastMessage, redirectTo, toast]);

  // Trả về true khi auth state đã sẵn sàng (đã load xong)
  return {
    isReady: !loading,
    isAuthed,
  };
}

/**
 * Hook để kiểm tra auth state mà không redirect.
 * Dùng cho các trang có partial functionality cần auth.
 *
 * @returns {Object} - { isAuthed, isReady }
 */
export function useAuthCheck() {
  const { isAuthed, loading } = useAuth();
  return {
    isAuthed,
    isReady: !loading,
  };
}
