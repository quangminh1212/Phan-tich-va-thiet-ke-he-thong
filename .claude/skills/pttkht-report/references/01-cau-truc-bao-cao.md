# Khung cấu trúc báo cáo PTTKHT - Hệ thống Quản lý Kho Hàng (cố định theo mục lục đã duyệt)

Cấu trúc 7 chương dưới đây là **BẮT BUỘC** - lấy từ mục lục bản báo cáo đã duyệt (`BaoCao/main.toc`). KHÔNG thêm/bớt/đổi tên chương hay mục cấp 2; mục cấp 3 (1.1.1) giữ đúng như liệt kê. Đánh số chương/mục phân cấp (1, 1.1, 1.1.1).

## Phần đầu (front matter) - KHÔNG đánh số chương

1. **Trang bìa** - bố cục căn giữa theo đúng thứ tự dọc sau:
   1. **ĐẠI HỌC BÁCH KHOA HÀ NỘI** (in đậm)
   2. **TRƯỜNG CÔNG NGHỆ THÔNG TIN VÀ TRUYỀN THÔNG** (in đậm)
   3. **Logo huy hiệu Bách Khoa** (căn giữa)
   4. **BÀI TẬP LỚN** (in đậm, cỡ vừa)
   5. **PHÂN TÍCH VÀ THIẾT KẾ HỆ THỐNG QUẢN LÝ KHO HÀNG** (in đậm, cỡ lớn nhất)
   6. **Bảng thông tin** 2 cột (không tô màu header): các dòng `GVHD`, `Lớp`, `Nhóm <N>`; ô bên phải dòng Nhóm liệt kê thành viên dạng `MSSV: Họ tên`, mỗi người một dòng (dùng `<br>` xuống dòng trong ô bảng).

   Logo có sẵn: `assets/hust-logo.png`. **Copy logo vào cùng thư mục báo cáo** và nhúng `![Bách Khoa](hust-logo.png)` (đường dẫn tương đối so với file Markdown).
2. **Lời nói đầu / Lời cảm ơn**.
3. **Mục lục**: liệt kê chi tiết đến cấp 1.1.1 (để PLACEHOLDER, sinh tự động khi xuất bản).
4. **Danh mục hình vẽ** (PLACEHOLDER, sinh khi xuất bản) và **Danh mục thuật ngữ & viết tắt**.

## Chương 1: KHẢO SÁT HIỆN TRẠNG
* 1.1 Giới thiệu doanh nghiệp
* 1.2 Phạm vi và đối tượng của hệ thống
    * 1.2.1 Phạm vi hệ thống
    * 1.2.2 Đối tượng sử dụng
* 1.3 Thực trạng quy trình quản lý hiện tại và các hạn chế
    * 1.3.1 Quy trình quản lý hiện tại
    * 1.3.2 Các hạn chế của quy trình hiện tại
* 1.4 Mục tiêu và yêu cầu của hệ thống
    * 1.4.1 Mục tiêu chung
    * 1.4.2 Yêu cầu chức năng (đánh mã FR-01, FR-02...)
    * 1.4.3 Yêu cầu phi chức năng (đánh mã NFR-01...)
* 1.5 Lựa chọn công nghệ và công cụ
    * 1.5.1 Lựa chọn công nghệ triển khai
    * 1.5.2 Hướng phân tích - thiết kế: hướng đối tượng (OO)
* 1.6 Kết luận chương

## Chương 2: MÔ TẢ NGHIỆP VỤ
* 2.1 Quy trình phân tích - thiết kế (hướng đối tượng)
* 2.2 Tổng quan nghiệp vụ
* 2.3 Nghiệp vụ 1: Quản lý tài khoản
* 2.4 Nghiệp vụ 2: Quản lý nhà cung cấp
* 2.5 Nghiệp vụ 3: Quản lý hàng hóa
* 2.6 Nghiệp vụ 4: Quản lý kho
* 2.7 Nghiệp vụ 5: Nhập kho
* 2.8 Nghiệp vụ 6: Xuất kho
* 2.9 Nghiệp vụ 7: Kiểm kê tồn kho
* 2.10 Nghiệp vụ 8: Báo cáo nhập - xuất - tồn
* 2.11 Biểu đồ hoạt động (theo nghiệp vụ chính)
    * 2.11.1 Nghiệp vụ nhập kho
    * 2.11.2 Nghiệp vụ xuất kho
    * 2.11.3 Nghiệp vụ kiểm kê
* 2.12 Kết luận chương

> Lưu ý: mục lục cũ có 2 mục "Kết luận chương" (2.11 và 2.13) - bản viết lại chuẩn hóa còn MỘT mục kết luận đặt cuối chương như trên.

## Chương 3: PHÂN TÍCH CHỨC NĂNG
* 3.1 Mục đích giai đoạn
* 3.2 Xác định tác nhân
* 3.3 Nguyên tắc tách use case
* 3.4 Phân tích chức năng bằng biểu đồ use case tổng quát
* 3.5 Phân tích chức năng bằng biểu đồ use case phân rã
    * 3.5.1 Nhóm truy cập và tài khoản
    * 3.5.2 Nhóm đối tác
    * 3.5.3 Nhóm hàng hóa và kho
    * 3.5.4 Nhóm nhập, xuất, kiểm kê, báo cáo
* 3.6 Bảng tổng hợp use case
* 3.7 Đặc tả use case - đặc tả VĂN BẢN ĐẦY ĐỦ cho TẤT CẢ 34 UC, mỗi UC một subsection 3.7.x theo mẫu bảng 7 dòng (`references/07-mau-bang-chuan-hust.md`):
    * 3.7.1 UC01 - Đăng nhập
    * 3.7.2 UC02 - Đăng xuất
    * 3.7.3-3.7.7 UC03-UC07 - Thêm / Sửa / Xóa / Tìm kiếm / Xem danh sách tài khoản
    * 3.7.8-3.7.12 UC08-UC12 - Thêm / Sửa / Xóa / Tìm kiếm / Xem danh sách nhà cung cấp
    * 3.7.13-3.7.17 UC13-UC17 - Thêm / Sửa / Xóa / Tìm kiếm / Xem danh sách hàng hóa
    * 3.7.18-3.7.22 UC18-UC22 - Thêm / Sửa / Xóa / Tìm kiếm / Xem danh sách kho
    * 3.7.23-3.7.27 UC23-UC27 - Lập / Sửa / Xóa / Tìm kiếm / Xem danh sách phiếu nhập
    * 3.7.28-3.7.32 UC28-UC32 - Lập / Sửa / Xóa / Tìm kiếm / Xem danh sách phiếu xuất
    * 3.7.33 UC33 - Lập phiếu kiểm kê
    * 3.7.34 UC34 - Xem báo cáo nhập xuất tồn

## Chương 4: PHÂN TÍCH CẤU TRÚC
* 4.1 Giới thiệu giai đoạn
* 4.2 Xác định các lớp chính
* 4.3 Phân tích cấu trúc bằng biểu đồ lớp tổng thể
    * 4.3.1 Nhóm danh mục (Master Data)
    * 4.3.2 Nhóm nghiệp vụ nhập - xuất
    * 4.3.3 Nhóm kiểm kê và lịch sử
* 4.4 Phân tích cấu trúc theo nhóm entity (biểu đồ lớp cắt lát)
    * 4.4.1 Nhóm tài khoản / xác thực (UC01-UC07)
    * 4.4.2 Nhóm nhà cung cấp (UC08-UC12)
    * 4.4.3 Nhóm hàng hóa (UC13-UC17)
    * 4.4.4 Nhóm kho (UC18-UC22)
    * 4.4.5 Nhóm phiếu nhập (UC23-UC27)
    * 4.4.6 Nhóm phiếu xuất (UC28-UC32)
    * 4.4.7 Nhóm kiểm kê (UC33)
    * 4.4.8 Nhóm báo cáo NXT (UC34)
* 4.5 Bảng mô tả thuộc tính các lớp
* 4.6 Giải thích các phương thức của lớp
    * 4.6.1 Lớp NguoiDung
    * 4.6.2 Lớp HangHoa
    * 4.6.3 Lớp TonKho
    * 4.6.4 Lớp PhieuNhapKho
    * 4.6.5 Lớp PhieuXuatKho
    * 4.6.6 Lớp PhieuKiemKe
    * 4.6.7 Lớp NhaCungCap / Kho (và các danh mục tương tự)
    * 4.6.8 Lớp LichSuThaoTac
    * 4.6.9 Liên hệ phương thức với use case (tóm tắt)
* 4.7 Tóm tắt quan hệ quan trọng
* 4.8 Kết luận chương

## Chương 5: PHÂN TÍCH HÀNH VI
* 5.1 Phân tích hành vi bằng biểu đồ trình tự - MỖI use case một biểu đồ tuần tự mermaid, subsection 5.1.1 đến 5.1.34 tương ứng UC01-UC34 (tên mục trùng tên UC ở mục 3.7)
* 5.2 Phân tích hành vi đối tượng phiếu (biểu đồ trạng thái)
    * 5.2.1 Phiếu nhập kho
    * 5.2.2 Phiếu xuất kho
* 5.3 Kết luận chương

## Chương 6: THIẾT KẾ CHI TIẾT
* 6.1 Component Diagram
* 6.2 Quan hệ chính giữa các lớp (thiết kế)
* 6.3 Kết luận chương

## Chương 7: THIẾT KẾ CƠ SỞ DỮ LIỆU VÀ GIAO DIỆN
* 7.1 Danh sách bảng
* 7.2 Đặc tả một số bảng quan trọng
    * 7.2.1 Bảng HangHoa
    * 7.2.2 Bảng TonKho
    * 7.2.3 Bảng PhieuNhapKho / PhieuXuatKho
    * 7.2.4 Bảng PhieuKiemKe
* 7.3 Thiết kế giao diện
    * 7.3.1 Sơ đồ điều hướng màn hình
* 7.4 Kết luận chương

## Phần cuối
* **Tài liệu tham khảo**: theo chuẩn IEEE/APA.

## Ánh xạ pha phân tích - chương (khung nhóm 3)

| Giai đoạn | Biểu đồ | Chương |
|-----------|---------|--------|
| Khảo sát hiện trạng | - | Ch.1 |
| Mô tả nghiệp vụ | Activity | Ch.2 |
| Phân tích chức năng | Use Case | Ch.3 |
| Phân tích cấu trúc | Class | Ch.4 |
| Phân tích hành vi | Sequence + State | Ch.5 |
| Thiết kế chi tiết | Component | Ch.6 |
| Thiết kế CSDL & giao diện | Bảng CSDL, điều hướng màn hình | Ch.7 |

## Quy ước trình bày (định dạng IEEE/học thuật)

- Font: Times New Roman, cỡ 13, giãn dòng 1.5 (áp dụng khi xuất .docx/.pdf).
- Biểu đồ: mermaid trong Markdown; rõ ràng, đúng ký pháp UML; không dùng ảnh chụp màn hình chất lượng thấp.
- Đánh số heading phân cấp; mỗi chương bắt đầu trang mới.
- **Hình**: đánh số theo chương (Hình 3.2), caption đặt **dưới** hình, căn giữa.
- **Bảng**: đánh số theo chương (Bảng 2.1), caption đặt **trên** bảng, căn giữa.
- Mọi hình/bảng phải được tham chiếu trong văn bản ("như Hình 3.2...").
- **Đánh số trang**: chân trang; KHÔNG đánh số trang bìa và mục lục.
- Thuật ngữ tiếng Anh giữ nguyên + giải thích lần đầu; viết tắt định nghĩa lần đầu.
- Nhất quán tên gọi xuyên suốt: NguoiDung, NhaCungCap, HangHoa, Kho, TonKho, PhieuNhapKho, PhieuXuatKho, PhieuKiemKe, LichSuThaoTac.
