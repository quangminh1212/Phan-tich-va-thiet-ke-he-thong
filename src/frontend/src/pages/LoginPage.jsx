import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

// UC01 Đăng nhập / UC02 Đăng xuất (nút ở layout)
export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [tenDangNhap, setTenDangNhap] = useState('');
  const [matKhau, setMatKhau] = useState('');
  const [err, setErr] = useState('');

  async function submit(e) {
    e.preventDefault();
    setErr('');
    try {
      const user = await login(tenDangNhap, matKhau);
      navigate(user.vai_tro === 'QTV' ? '/users' : '/products');
    } catch (ex) {
      setErr(ex.message);
    }
  }

  return (
    <div className="login-wrap">
      <form className="login-card" onSubmit={submit}>
        <h1>Hệ thống Quản lý Kho Hàng</h1>
        <div style={{ display: 'grid', gap: 10 }}>
          <input placeholder="Tên đăng nhập" value={tenDangNhap} onChange={(e) => setTenDangNhap(e.target.value)} />
          <input type="password" placeholder="Mật khẩu" value={matKhau} onChange={(e) => setMatKhau(e.target.value)} />
          {err && <div className="error">{err}</div>}
          <button type="submit">Đăng nhập</button>
          <div className="muted">Tài khoản mẫu: admin / quanly / nhanvien (mật khẩu: tên + 123)</div>
        </div>
      </form>
    </div>
  );
}
