# Mẫu bảng chuẩn (rút từ báo cáo HUST thực tế đạt yêu cầu)

Các quy ước dưới đây lấy từ một báo cáo PTTKHT thật của HUST (đề tài Quản lý phòng GYM) — gọn, nhất quán, dễ chấm. **Ưu tiên dùng các mẫu này.**

## 1. Mã use case phân cấp theo gói chức năng

Đặt mã UC theo nhóm chức năng: `UC-<nhóm>.<thứ tự>`. Ví dụ:
- Nhóm 1 Quản trị hệ thống: UC-1.1 Thêm nhân sự, UC-1.2 Chuẩn hóa địa chỉ...
- Nhóm 2 Vận hành & CSKH: UC-2.1 Đăng ký hội viên...
- Nhóm 3 Dịch vụ & Đào tạo: UC-3.1 Đặt lịch tập PT...
- Nhóm 4 Tài chính & Thanh toán: UC-4.1 Kiểm tra hóa đơn chờ...

Nhóm này khớp với các package trong biểu đồ use case tổng quát → truy vết rõ ràng.

## 2. Bảng đặc tả use case (7 dòng, dọc "Thành phần | Nội dung")

| Thành phần | Nội dung |
|---|---|
| Mã Use Case | UC-1.1 |
| Tên Use Case | Thêm mới nhân sự |
| Tác nhân chính | Quản trị viên (ADMIN) |
| Mô tả tóm tắt | Quản trị viên khởi tạo tài khoản và phân quyền cho nhân sự |
| Điều kiện | Quản trị viên đã đăng nhập và xác thực thành công |
| Thao tác chính | 1. ... 2. ... 3. ... (đánh số bước; số bước khớp message trong sequence) |

> Với UC phức tạp, bổ sung dòng **Luồng ngoại lệ** và **Hậu điều kiện**.

## 3. Mô tả chi tiết từng lớp = 2 bảng/lớp

Mỗi lớp một mục riêng (vd "Lớp Staff", "Lớp Member"...), gồm:

**Bảng thuộc tính** — `Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa`

| Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa |
|---|---|---|
| Id | int | Mã định danh nhân sự (Khóa chính) |
| Username | string | Tên đăng nhập hệ thống |
| Role | string | Vai trò phân quyền (ADMIN, STAFF, PT) |

**Bảng phương thức** — `Tên phương thức | Ý nghĩa` (ghi kèm mã UC để truy vết)

| Tên phương thức | Ý nghĩa |
|---|---|
| Login() | Kiểm tra thông tin đăng nhập và xác thực phân quyền |
| CreateStaff() | Thêm mới tài khoản nhân sự và phân quyền (UC-1.1) |
| DisableAccount() | Vô hiệu hóa tài khoản khi nhân sự nghỉ việc (UC-1.3) |

> Ghi `(UC-x.y)` trong cột Ý nghĩa của phương thức để liên kết phương thức ↔ use case. Phương thức ở đây phải khớp biểu đồ lớp thiết kế và message trong sequence (chạy `scripts/check-sync.py`).

## 4. Từ điển dữ liệu (đặc tả bảng CSDL) — 4 cột

`Tên trường | Kiểu dữ liệu | Kích thước | Ràng buộc toàn vẹn`

| Tên trường | Kiểu dữ liệu | Kích thước | Ràng buộc toàn vẹn |
|---|---|---|---|
| id | int |  | Khóa chính |
| phone | varchar(255) | 255 ký tự |  |
| ward_id | varchar(20) | 20 ký tự | Khóa ngoại → ward |
| status | int |  |  |
| create_at | timestamp |  |  |

## 5. Quy ước đặt tên (nhất quán toàn báo cáo)

- **Lớp**: tiếng Anh, PascalCase — `Staff`, `Member`, `GymPackage`, `GymDevice`, `PTSchedule`, `Invoice`, `Payment`.
- **Thuộc tính lớp**: PascalCase/camelCase — `Id`, `Username`, `WardId`, `CreateAt`.
- **Cột CSDL**: snake_case — `id`, `ward_id`, `create_at`. Ánh xạ rõ ràng thuộc tính↔cột (Id↔id, WardId↔ward_id).
- Giữ một tên duy nhất cho mỗi khái niệm xuyên suốt (xem `04-review-checklist.md`).

## 6. Tổ chức chương "Phân tích hành vi"

Gom các biểu đồ động (sequence, activity, state machine) vào một mục **"Phân tích hành vi"** (Ch3.3) — thay vì rải rác. Vẫn theo quy tắc HUST: chỉ ~4 use case tiêu biểu cho mỗi loại biểu đồ. (Biểu đồ lớp thiết kế chi tiết thuộc Ch5.1; CSDL/đặc tả bảng thuộc Ch5.2.)
