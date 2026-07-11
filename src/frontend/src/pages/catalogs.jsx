import { useEffect, useState } from 'react';
import CatalogPage from '../components/CatalogPage';
import { api } from '../api/client';

const activeBadge = (v, onLabel = 'Hoạt động', offLabel = 'Ngừng / Khóa') => (
  <span className={`badge ${v ? 'on' : 'off'}`}>{v ? onLabel : offLabel}</span>
);

// UC03-UC07: quản lý tài khoản (QTV)
export function UsersPage() {
  return (
    <CatalogPage
      title="Quản lý tài khoản"
      endpoint="/users"
      idKey="ma_nguoi_dung"
      columns={[
        { key: 'ma_nguoi_dung', label: 'Mã' },
        { key: 'ten_dang_nhap', label: 'Tên đăng nhập' },
        { key: 'ho_ten', label: 'Họ tên' },
        { key: 'email', label: 'Email' },
        { key: 'vai_tro', label: 'Vai trò' },
        { key: 'trang_thai', label: 'Trạng thái', render: (r) => activeBadge(r.trang_thai) },
      ]}
      fields={[
        { key: 'ten_dang_nhap', label: 'Tên đăng nhập', required: true },
        { key: 'mat_khau', label: 'Mật khẩu', type: 'password' },
        { key: 'ho_ten', label: 'Họ tên', required: true },
        { key: 'email', label: 'Email' },
        {
          key: 'vai_tro', label: 'Vai trò', type: 'select', required: true,
          options: [
            { value: 'QTV', label: 'Quản trị viên' },
            { value: 'QLKho', label: 'Quản lý kho' },
            { value: 'NVKho', label: 'Nhân viên kho' },
          ],
        },
        { key: 'trang_thai', label: 'Hoạt động', type: 'checkbox' },
      ]}
    />
  );
}

// UC08-UC12: quản lý nhà cung cấp
export function SuppliersPage() {
  return (
    <CatalogPage
      title="Quản lý nhà cung cấp"
      endpoint="/suppliers"
      idKey="ma_ncc"
      columns={[
        { key: 'ma_ncc', label: 'Mã' },
        { key: 'ten_ncc', label: 'Tên NCC' },
        { key: 'dia_chi', label: 'Địa chỉ' },
        { key: 'sdt', label: 'SĐT' },
        { key: 'nguoi_lien_he', label: 'Người liên hệ' },
        { key: 'trang_thai', label: 'Trạng thái', render: (r) => activeBadge(r.trang_thai) },
      ]}
      fields={[
        { key: 'ten_ncc', label: 'Tên NCC', required: true },
        { key: 'dia_chi', label: 'Địa chỉ' },
        { key: 'sdt', label: 'SĐT' },
        { key: 'email', label: 'Email' },
        { key: 'nguoi_lien_he', label: 'Người liên hệ' },
        { key: 'trang_thai', label: 'Hoạt động', type: 'checkbox' },
      ]}
    />
  );
}

// UC13-UC17: quản lý hàng hóa (thêm hàng tự khởi tạo tồn = 0)
export function ProductsPage() {
  const [groups, setGroups] = useState([]);
  const [units, setUnits] = useState([]);
  useEffect(() => {
    api('/groups').then((d) => setGroups(d.items || [])).catch(() => {});
    api('/units').then((d) => setUnits(d.items || [])).catch(() => {});
  }, []);
  return (
    <CatalogPage
      title="Quản lý hàng hóa"
      endpoint="/products"
      idKey="ma_hang_hoa"
      columns={[
        { key: 'ma_hang_hoa', label: 'Mã' },
        { key: 'ten_hang_hoa', label: 'Tên hàng' },
        { key: 'ma_nhom_hang', label: 'Nhóm' },
        { key: 'ma_don_vi_tinh', label: 'ĐVT' },
        { key: 'gia_nhap', label: 'Giá nhập', render: (r) => Number(r.gia_nhap).toLocaleString('vi-VN') },
        { key: 'gia_xuat', label: 'Giá xuất', render: (r) => Number(r.gia_xuat).toLocaleString('vi-VN') },
        { key: 'tong_ton', label: 'Tổng tồn' },
        { key: 'trang_thai', label: 'Trạng thái', render: (r) => activeBadge(r.trang_thai, 'Đang dùng', 'Ngừng') },
      ]}
      fields={[
        { key: 'ten_hang_hoa', label: 'Tên hàng', required: true },
        {
          key: 'ma_nhom_hang', label: 'Nhóm hàng', type: 'select',
          options: groups.map((g) => ({ value: g.ma_nhom_hang, label: g.ten_nhom_hang })),
        },
        {
          key: 'ma_don_vi_tinh', label: 'Đơn vị tính', type: 'select',
          options: units.map((u) => ({ value: u.ma_don_vi_tinh, label: u.ten_don_vi_tinh })),
        },
        { key: 'gia_nhap', label: 'Giá nhập', type: 'number' },
        { key: 'gia_xuat', label: 'Giá xuất', type: 'number' },
        { key: 'mo_ta', label: 'Mô tả' },
        { key: 'trang_thai', label: 'Đang dùng', type: 'checkbox' },
      ]}
    />
  );
}

// UC18-UC22: quản lý kho (chỉ QL kho)
export function WarehousesPage() {
  return (
    <CatalogPage
      title="Quản lý kho"
      endpoint="/warehouses"
      idKey="ma_kho"
      columns={[
        { key: 'ma_kho', label: 'Mã' },
        { key: 'ten_kho', label: 'Tên kho' },
        { key: 'dia_chi', label: 'Địa chỉ' },
        { key: 'nguoi_phu_trach', label: 'Người phụ trách' },
        { key: 'trang_thai', label: 'Trạng thái', render: (r) => activeBadge(r.trang_thai) },
      ]}
      fields={[
        { key: 'ten_kho', label: 'Tên kho', required: true },
        { key: 'dia_chi', label: 'Địa chỉ' },
        { key: 'nguoi_phu_trach', label: 'Người phụ trách' },
        { key: 'trang_thai', label: 'Hoạt động', type: 'checkbox' },
      ]}
    />
  );
}
