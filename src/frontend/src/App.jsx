import { BrowserRouter, Navigate, NavLink, Route, Routes, useNavigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './auth/AuthContext';
import LoginPage from './pages/LoginPage';
import { ProductsPage, SuppliersPage, UsersPage, WarehousesPage } from './pages/catalogs';
import TicketPage from './components/TicketPage';
import StocktakePage from './pages/StocktakePage';
import ReportPage from './pages/ReportPage';

// Menu theo vai trò - khớp sơ đồ điều hướng Hình 7.1 của báo cáo:
// QTV chỉ thấy Tài khoản; NV kho không thấy Kho và Báo cáo.
const MENU = [
  { to: '/users', label: 'Tài khoản', roles: ['QTV'] },
  { to: '/suppliers', label: 'Nhà cung cấp', roles: ['NVKho', 'QLKho'] },
  { to: '/products', label: 'Hàng hóa', roles: ['NVKho', 'QLKho'] },
  { to: '/warehouses', label: 'Kho', roles: ['QLKho'] },
  { to: '/imports', label: 'Phiếu nhập', roles: ['NVKho', 'QLKho'] },
  { to: '/exports', label: 'Phiếu xuất', roles: ['NVKho', 'QLKho'] },
  { to: '/stocktakes', label: 'Kiểm kê', roles: ['NVKho', 'QLKho'] },
  { to: '/reports', label: 'Báo cáo NXT', roles: ['QLKho'] },
];

const ROLE_LABEL = { QTV: 'Quản trị viên', QLKho: 'Quản lý kho', NVKho: 'Nhân viên kho' };

function Layout({ children, allow }) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  if (!user) return <Navigate to="/login" replace />;
  if (allow && !allow.includes(user.vai_tro)) return <Navigate to="/" replace />;

  async function doLogout() {
    await logout();
    navigate('/login');
  }

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="brand">Quản lý Kho Hàng</div>
        <nav>
          {MENU.filter((m) => m.roles.includes(user.vai_tro)).map((m) => (
            <NavLink key={m.to} to={m.to} className={({ isActive }) => (isActive ? 'active' : '')}>
              {m.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main className="content">
        <div className="topbar">
          <span className="who">
            {user.ho_ten} - {ROLE_LABEL[user.vai_tro] || user.vai_tro}
          </span>
          <button className="secondary" onClick={doLogout}>Đăng xuất</button>
        </div>
        {children}
      </main>
    </div>
  );
}

function Home() {
  const { user } = useAuth();
  if (!user) return <Navigate to="/login" replace />;
  return <Navigate to={user.vai_tro === 'QTV' ? '/users' : '/products'} replace />;
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/" element={<Home />} />
          <Route path="/users" element={<Layout allow={['QTV']}><UsersPage /></Layout>} />
          <Route path="/suppliers" element={<Layout allow={['NVKho', 'QLKho']}><SuppliersPage /></Layout>} />
          <Route path="/products" element={<Layout allow={['NVKho', 'QLKho']}><ProductsPage /></Layout>} />
          <Route path="/warehouses" element={<Layout allow={['QLKho']}><WarehousesPage /></Layout>} />
          <Route path="/imports" element={<Layout allow={['NVKho', 'QLKho']}><TicketPage kind="import" /></Layout>} />
          <Route path="/exports" element={<Layout allow={['NVKho', 'QLKho']}><TicketPage kind="export" /></Layout>} />
          <Route path="/stocktakes" element={<Layout allow={['NVKho', 'QLKho']}><StocktakePage /></Layout>} />
          <Route path="/reports" element={<Layout allow={['QLKho']}><ReportPage /></Layout>} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
