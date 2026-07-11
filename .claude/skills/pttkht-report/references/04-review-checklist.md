# Điều kiện kiểm duyệt (Review Checklist) - bản Quản lý Kho Hàng, 7 chương

Đối chiếu TỪNG mục. Chỉ kết luận "đạt" khi tất cả ✅. Chia 3 nhóm: Cấu trúc - Nội dung/Đồng bộ - Định dạng. KHÔNG kiểm tra độ dài trang.

## A. Cấu trúc tài liệu

- [ ] Có đủ front matter: trang bìa, lời nói đầu, mục lục (đến cấp 1.1.1), danh mục hình vẽ, danh mục thuật ngữ & viết tắt.
- [ ] Có đủ **7 chương đúng tên và đúng mục** theo `01-cau-truc-bao-cao.md`; không thiếu/gộp/thêm/đổi tên chương hay mục cấp 2.
- [ ] Mục 3.7 có đủ 34 subsection đặc tả UC01-UC34; mục 5.1 có đủ 34 subsection sequence UC01-UC34; tên UC ở 3.7 và 5.1 trùng khớp từng chữ.
- [ ] Chương 2 chỉ có MỘT mục "Kết luận chương" (2.12), đặt cuối chương.
- [ ] Có tài liệu tham khảo.
- [ ] Mỗi chương mở đầu bằng đoạn dẫn, kết bằng mục "Kết luận chương".
- [ ] Đánh số chương/mục phân cấp, không nhảy số.

## B. Nội dung & Đồng bộ (quan trọng nhất khi chấm)

- [ ] Mục tiêu & phạm vi (Ch.1) rõ ràng, đo được; các kết luận chương đối chiếu lại được.
- [ ] Yêu cầu chức năng có mã (FR-xx) và đầy đủ; có cả yêu cầu phi chức năng (NFR-xx); mỗi FR ánh xạ được tới ít nhất một UC.
- [ ] Ch.2: đủ 8 nghiệp vụ; 3 biểu đồ hoạt động (nhập kho, xuất kho, kiểm kê) khớp mô tả nghiệp vụ tương ứng.
- [ ] Biểu đồ use case tổng quát (3.4) đầy đủ tác nhân & nhóm chức năng; biểu đồ phân rã (3.5) phủ đủ 34 UC; quan hệ include/extend đúng.
- [ ] **Đủ 34 UC** có đặc tả kịch bản theo mẫu bảng 7 dòng (tác nhân, tiền/hậu điều kiện, luồng chính, luồng ngoại lệ).
- [ ] **Đủ 34 UC** có biểu đồ tuần tự (5.1.1-5.1.34); số bước luồng chính trong đặc tả khớp message trong sequence.
- [ ] Ch.4: biểu đồ lớp tổng thể (4.3, chia 3 nhóm) + 8 biểu đồ lớp cắt lát theo nhóm entity (4.4) nhất quán với nhau; bảng thuộc tính (4.5) và giải thích phương thức (4.6) khớp biểu đồ.
- [ ] **Sequence - Class**: chạy `python3 scripts/check-sync.py <file>.md` với kết quả thoát mã 0 (BẮT BUỘC, không khai khống). Mọi message khớp phương thức của lớp đích.
- [ ] **Class - CSDL**: ánh xạ lớp sang bảng đầy đủ (Ch.7); khóa chính/ngoại, ràng buộc, kiểu dữ liệu rõ.
- [ ] Có biểu đồ trạng thái cho PhieuNhapKho và PhieuXuatKho (5.2); trạng thái khớp với luồng trong đặc tả UC và sequence liên quan.
- [ ] Ch.6: component diagram và quan hệ lớp thiết kế nhất quán với lớp ở Ch.4.
- [ ] Thiết kế giao diện (7.3) ánh xạ tới use case; sơ đồ điều hướng màn hình phủ các nhóm chức năng chính.
- [ ] **Tên nhất quán** xuyên suốt (actor, UC, lớp, thuộc tính, bảng): NguoiDung, NhaCungCap, HangHoa, Kho, TonKho, PhieuNhapKho, PhieuXuatKho, PhieuKiemKe, LichSuThaoTac.
- [ ] Truy vết được Yêu cầu --> UC --> Sequence --> Lớp --> Bảng; không phần tử "mồ côi".
- [ ] Không mâu thuẫn logic (vd phương thức gọi trong sequence không tồn tại; bội số xung đột; tác nhân thừa).

## C. Định dạng (IEEE/SoICT)

- [ ] Có đủ bộ biểu đồ: Use Case (tổng quát + 4 phân rã), Activity (3), Class (tổng thể + 8 cắt lát), Sequence (34), State (2), Component, thiết kế CSDL, sơ đồ điều hướng màn hình.
- [ ] Có đặc tả bảng CSDL cho các bảng quan trọng (7.2) theo mẫu từ điển dữ liệu ở `07-mau-bang-chuan-hust.md`.
- [ ] Tuân thủ quy định nộp (xem `06-quy-dinh-nop-hust.md`): xuất **cả PDF và DOCX**; nếu nộp repo thì đúng cấu trúc thư mục.
- [ ] Mọi hình có số + caption (dưới hình) + được tham chiếu trong văn bản.
- [ ] Mọi bảng có số + caption (trên bảng) + được tham chiếu.
- [ ] Hình/bảng đánh số theo chương (Hình 3.2, Bảng 2.1).
- [ ] Biểu đồ UML đúng ký pháp (mũi tên, loại quan hệ, lifeline, multiplicity).
- [ ] Font/giãn dòng/lề/căn lề nhất quán; mỗi chương sang trang mới.
- [ ] Thuật ngữ tiếng Anh & viết tắt được giải thích ở lần dùng đầu.
- [ ] Trích dẫn tài liệu tham khảo nhất quán một chuẩn (IEEE/APA) và xuất hiện trong thân bài.
- [ ] Chính tả, ngữ pháp tiếng Việt; văn phong khách quan, không ngôi thứ nhất số ít.

## Cách dùng khi review
1. Chạy qua A --> B --> C, đánh dấu từng mục.
2. Với mục FAIL: nêu rõ vị trí (chương/hình/bảng) + lý do + cách sửa.
3. Sửa, rồi chạy lại nhóm B (đồng bộ) vì sửa một biểu đồ thường kéo theo biểu đồ khác.
4. Chỉ báo "hoàn thành" khi 100% mục B đạt và toàn bộ A, C đạt.
