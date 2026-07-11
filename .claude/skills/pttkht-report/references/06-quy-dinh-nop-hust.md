# Quy định nộp bài & tổ chức repo (HUST / SoICT)

Nguồn: quy định môn **IT3120 Phân tích & Thiết kế Hệ thống** (SoICT) và mẫu **SINNOProjectTemplate** của SOICT Innovation Club. Áp dụng khi đề bài yêu cầu nộp theo repo GitHub.

## Định dạng nộp
- Báo cáo cuối cùng phải có **cả 2 định dạng: `BaoCao.pdf` và `BaoCao.docx`**.
- Soạn nội dung bằng Markdown rồi chuyển sang `.docx` bằng công cụ tùy chọn; xuất PDF bằng cách mở .docx và "Save as PDF", hoặc `soffice --headless --convert-to pdf bao-cao.docx` nếu có LibreOffice.

## Cấu trúc repo GitHub (bắt buộc với BTL IT3120)
```
<repo>/
├── Use Case Diagram/
├── Activity Diagram/
├── Sequence Diagram/
├── State Machine Diagram/
├── Class Diagram/
├── Database Design/
├── BaoCao.pdf
├── BaoCao.docx
└── README.md            # tổng quan đề tài, thành viên, hướng dẫn
```
Mỗi thư mục biểu đồ chứa file nguồn (Draw.io/.drawio, .puml) và ảnh xuất (.png).

## Quy định quá trình (ảnh hưởng điểm)
- **Commit định kỳ ~2 tuần/lần** theo tiến độ; nhóm chỉ commit dồn cuối kỳ bị **trừ 30% điểm**.
- **Mọi thành viên phải có commit** (thể hiện đóng góp) và **phải tham gia bảo vệ**; giảng viên hỏi về phần đã làm — không hiểu phần mình làm có thể bị trừ điểm cá nhân.
- README.md nêu rõ tổng quan dự án.

## Quy tắc nội dung (SoICT)
- "Không cần vẽ tất cả biểu đồ cho mọi use case — mỗi loại biểu đồ tập trung **~4 use case tiêu biểu/hấp dẫn nhất**." Ưu tiên **chất lượng hơn số lượng**.
- Mỗi chương mở đầu bằng đoạn tóm tắt, hạn chế liệt kê gạch đầu dòng tràn lan.
- Biểu đồ vẽ bằng công cụ chuyên dụng (Draw.io/Lucidchart), mockup giao diện bằng **Balsamiq**; **không** chèn ảnh chụp màn hình biểu đồ.

## Công cụ hỗ trợ (SINNO khuyến nghị)
- Sinh tài liệu thiết kế chi tiết tự động từ mã nguồn: **Doxygen + Graphviz** (sơ đồ lớp, cây gọi hàm).
- Quản lý công việc: MS Planner (có sẵn trong Office365 của SV Bách Khoa), Trello.
- Quy trình nhánh: GitFlow.

> Lưu ý: với bài tập lớn thuần phân tích–thiết kế (không code), phần "sinh tự động từ mã nguồn" có thể không áp dụng; khi đó vẽ thủ công bằng Draw.io là đủ.
