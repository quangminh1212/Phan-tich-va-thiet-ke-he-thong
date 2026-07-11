---
name: pttkht-report
description: Soạn "Báo cáo Bài tập lớn Phân tích và Thiết kế Hệ thống Quản lý Kho Hàng" hướng đối tượng (UML) - bản tùy biến riêng cho project này, cấu trúc 7 chương cố định theo mục lục đã duyệt. Use when user wants to write/rewrite/review the PTTKHT warehouse-management report in this repo, đặc tả use case, biểu đồ lớp/tuần tự/trạng thái.
---

# Báo cáo PTTKHT - Hệ thống Quản lý Kho Hàng (bản tùy biến project)

Bản tùy biến của skill pttkht-report cho đề tài **Quản lý Kho Hàng** trong repo này. Khác bản gốc: cấu trúc **7 chương CỐ ĐỊNH** theo mục lục đã duyệt (xem [references/01-cau-truc-bao-cao.md](references/01-cau-truc-bao-cao.md)), **KHÔNG kiểm tra độ dài trang**. Mục tiêu: **chính xác, đầy đủ, đồng bộ giữa các biểu đồ, đạt mọi điều kiện review.**

Hướng tiếp cận: **đối tượng (UML)**, khung phân tích "nhóm 3": phân tích chức năng (Use Case) --> phân tích cấu trúc (Class) --> phân tích hành vi (Sequence + State).

## Quy trình làm việc (bắt buộc theo thứ tự)

1. **Dựng khung tài liệu** - Tạo đúng cấu trúc 7 chương + front matter theo [references/01-cau-truc-bao-cao.md](references/01-cau-truc-bao-cao.md). KHÔNG thêm/bớt/đổi tên chương hay mục cấp 2. Báo cáo viết bằng Markdown + mermaid, đặt tại `v2/`.
2. **Điền nội dung theo pha** - đi tuần tự 7 chương: Khảo sát hiện trạng --> Mô tả nghiệp vụ (8 nghiệp vụ + activity) --> Phân tích chức năng (use case tổng quát/phân rã, đặc tả đủ 34 UC) --> Phân tích cấu trúc (lớp tổng thể + lớp cắt lát theo nhóm entity, thuộc tính, phương thức) --> Phân tích hành vi (34 sequence + 2 state) --> Thiết kế chi tiết (component, quan hệ lớp) --> Thiết kế CSDL và giao diện.
3. **Đồng bộ chéo các biểu đồ** - sau khi điền, **BẮT BUỘC chạy** `python3 scripts/check-sync.py <file>.md` (thoát mã phải = 0). Script kiểm tự động Sequence-Class (mọi message khớp phương thức khai báo) và Class-CSDL. KHÔNG được tự khẳng định "đã đồng bộ" mà chưa chạy script. Quy tắc chi tiết: [references/03-dong-bo-bieu-do.md](references/03-dong-bo-bieu-do.md).
4. **Tự kiểm duyệt (Review)** - đối chiếu TỪNG mục trong [references/04-review-checklist.md](references/04-review-checklist.md). Chỉ báo "hoàn thành" khi mọi mục đạt. KHÔNG kiểm tra độ dài trang.
5. **Xuất bản & nộp** - chuyển Markdown sang `.docx`/`.pdf` bằng công cụ tùy chọn. Đảm bảo các yêu cầu định dạng (trang bìa căn giữa đủ chữ, lề ĐHBK, đánh số trang bỏ trang bìa, mục lục đến cấp 1.1.1) - chi tiết ở [references/01-cau-truc-bao-cao.md](references/01-cau-truc-bao-cao.md).
   - **Verify trực quan**: mở bản `.docx`/`.pdf` kiểm tra bằng mắt - trang bìa hiện đủ chữ, có mục lục, biểu đồ đọc được, không tràn trang; chạy `scripts/check-sync.py` cho phần đồng bộ.
   - Theo quy định HUST cần **cả PDF lẫn DOCX**; nộp repo theo [references/06-quy-dinh-nop-hust.md](references/06-quy-dinh-nop-hust.md).

## Nguyên tắc cốt lõi (để qua review)

- **Đầy đủ**: đặc tả văn bản cho **đủ 34 use case** (mục 3.7.1-3.7.34); biểu đồ tuần tự cho **đủ 34 use case** (mục 5.1.1-5.1.34); activity cho 3 nghiệp vụ chính (nhập kho, xuất kho, kiểm kê); state cho phiếu nhập và phiếu xuất. Mỗi tác nhân phải xuất hiện trong biểu đồ use case tổng quát.
- **Truy vết (traceability)**: yêu cầu chức năng --> use case --> kịch bản --> biểu đồ tuần tự --> lớp/phương thức --> bảng CSDL. Mọi phần tử phải truy ngược được về một yêu cầu.
- **Nhất quán tên gọi**: 1 khái niệm = 1 tên duy nhất xuyên suốt (tác nhân, use case, lớp, thuộc tính). Lớp trong Sequence = lớp trong Class diagram. Bộ lớp chuẩn của đề tài: NguoiDung, NhaCungCap, HangHoa, Kho, TonKho, PhieuNhapKho, PhieuXuatKho, PhieuKiemKe, LichSuThaoTac.
- **Đúng ký pháp UML**: chiều mũi tên, loại quan hệ (association/generalization/dependency, include/extend), bội số (multiplicity).
- **Định dạng học thuật**: mọi hình/bảng được đánh số theo chương + chú thích (caption) + được tham chiếu trong văn bản ("như Hình 3.2..."). Có mục lục, danh mục hình vẽ, danh mục thuật ngữ & viết tắt.

## Mẫu tham khảo

- **Mẫu bảng chuẩn** (đặc tả UC 7 dòng, mô tả lớp 2 bảng, từ điển dữ liệu, quy ước đặt tên) rút từ báo cáo HUST thực tế: [references/07-mau-bang-chuan-hust.md](references/07-mau-bang-chuan-hust.md).
- Template use case và mermaid: [references/03-dong-bo-bieu-do.md](references/03-dong-bo-bieu-do.md).
- Nội dung nghiệp vụ/UC/lớp/CSDL của đề tài: tham khảo bản LaTeX hiện có tại `BaoCao/chapters/*.tex` (nguồn sự thật về tên UC, tên lớp, thuộc tính, bảng CSDL).

**Logo trang bìa**: `assets/hust-logo.png` là **huy hiệu Đại học Bách Khoa** - logo chính thức dùng cho trang bìa. Copy cạnh file báo cáo và nhúng `![Bách Khoa](hust-logo.png)` (căn giữa, trên cùng trang bìa).
