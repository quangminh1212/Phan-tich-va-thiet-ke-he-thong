# Hệ thống Quản lý Kho Hàng - mã nguồn

Hiện thực hóa thiết kế trong báo cáo PTTKHT (34 use case, 14 bảng CSDL):
**Go (Gin + GORM) + PostgreSQL + React (Vite)**, xác thực JWT, phân quyền 3 vai trò.

## Cấu trúc

```
src/
  docker-compose.yml   PostgreSQL 16
  backend/             API Go (Gin + GORM), cổng 8080
    cmd/server/        main
    internal/
      models/          14 model = 14 bảng Chương 7 báo cáo
      handlers/        auth, users, suppliers, products, warehouses,
                       imports, exports, stocktakes, reports (UC01-UC34)
      services/        stock (tangTon/giamTon/kiemTraDu/ganTheoKiemKe), audit (ghiLog)
      middleware/      JWT + RequireRoles
  frontend/            React (Vite), cổng 5173, proxy /api -> 8080
```

## Chạy hệ thống

```bash
# 1. CSDL
cd src && docker compose up -d

# 2. Backend (tự migrate + seed dữ liệu mẫu lần đầu)
cd backend && go run ./cmd/server

# 3. Frontend
cd ../frontend && npm install && npm run dev
# mở http://localhost:5173
```

Tài khoản mẫu (seed): `admin/admin123` (QTV), `quanly/quanly123` (QL kho), `nhanvien/nhanvien123` (NV kho).

## Phân quyền (theo Bảng 3.2 của báo cáo)

| Chức năng | QTV | QL kho | NV kho |
|---|---|---|---|
| Đăng nhập / đăng xuất (UC01-02) | x | x | x |
| Quản lý tài khoản (UC03-07) | x | | |
| NCC / Hàng hóa (UC08-17) | | x | x |
| Kho (UC18-22) | | x | |
| Lập / sửa / xóa / tìm / xem phiếu nhập - xuất (UC23-32) | | x | x |
| Duyệt / từ chối phiếu | | x | |
| Kiểm kê (UC33) - duyệt điều chỉnh chỉ QL | | x | x |
| Báo cáo NXT (UC34) | | x | |

## Quy tắc nghiệp vụ chính

- Tồn kho **chỉ** thay đổi trong transaction duyệt phiếu: duyệt nhập tăng tồn, duyệt xuất giảm tồn, duyệt kiểm kê gán tồn = số thực tế.
- Phiếu xuất kiểm tra tồn đủ khi lập / sửa và kiểm tra lại khi duyệt.
- Thêm hàng hóa khởi tạo tồn = 0 tại mọi kho; thêm kho khởi tạo tồn = 0 cho mọi hàng.
- Xóa bản ghi đã có dữ liệu tham chiếu sẽ chuyển thành khóa / ngừng hoạt động.
- Mọi thao tác ghi được lưu vào lich_su_thao_tac (QTV xem qua GET /api/logs).
