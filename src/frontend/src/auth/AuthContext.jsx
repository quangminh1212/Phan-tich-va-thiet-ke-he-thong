import { createContext, useContext, useState } from 'react';
import { api } from '../api/client';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    try {
      return JSON.parse(localStorage.getItem('user')) || null;
    } catch {
      return null;
    }
  });

  // UC01 Đăng nhập
  async function login(tenDangNhap, matKhau) {
    const data = await api('/auth/login', {
      method: 'POST',
      body: { ten_dang_nhap: tenDangNhap, mat_khau: matKhau },
    });
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    setUser(data.user);
    return data.user;
  }

  // UC02 Đăng xuất
  async function logout() {
    try {
      await api('/auth/logout', { method: 'POST' });
    } catch {
      /* token có thể đã hết hạn - vẫn xóa phía client */
    }
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
  }

  return <AuthContext.Provider value={{ user, login, logout }}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}
