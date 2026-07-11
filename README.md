# Báo cáo PTTKHT - Hệ thống Quản lý Kho Hàng (v2)

Báo cáo Bài tập lớn Phân tích và Thiết kế Hệ thống, viết bằng Markdown (dialect
[md-to-docx](https://www.npmjs.com/package/@luytbq43/md-to-docx)), biểu đồ UML vẽ bằng draw.io
và nhúng dưới dạng PNG.

## Cấu trúc thư mục

```
BaoCao_QuanLyKhoHang.md      Nguồn báo cáo (Markdown + directive md-to-docx)
BaoCao_QuanLyKhoHang.docx    Bản DOCX build từ file .md
hust-logo.png                Logo trang bìa
diagram/
  drawio/                    57 file .drawio (nguồn biểu đồ, mở bằng draw.io)
    gen_all_drawio.py        Generator sinh toàn bộ .drawio từ SPEC (CLASSES / USE_CASES)
    check_sync.py            Kiểm tra đồng bộ Sequence-Class-CSDL-Markdown
  png/                       57 file PNG xuất từ drawio/ (được nhúng vào báo cáo)
script/
  generate_docx_from_md.bat  Build DOCX trên Windows
  generate_docx_from_md.sh   Build DOCX trên macOS / Linux
v1/                          Bản báo cáo LaTeX cũ (tham chiếu)
README.md
```

## 1. Build DOCX từ Markdown

Yêu cầu: Node.js (có `npx`).

**Windows** - chạy:

```bat
script\generate_docx_from_md.bat
```

**macOS / Linux** - chạy:

```bash
script/generate_docx_from_md.sh
```

(tương đương: `npx -y -p @luytbq43/md-to-docx md-to-docx BaoCao_QuanLyKhoHang.md -o BaoCao_QuanLyKhoHang.docx` từ thư mục gốc)

Converter in ra cảnh báo (`var` / `link` / `image` / `style`...) nếu nguồn có lỗi - sửa nguồn
tới khi **0 warning**. Xuất PDF: mở file DOCX bằng Word rồi Save As PDF (HUST yêu cầu nộp cả
DOCX lẫn PDF).

## 2. Sửa / sinh lại biểu đồ

Mọi biểu đồ sinh tự động từ **SPEC** trong `diagram/drawio/gen_all_drawio.py`
(`CLASSES` = lớp / thuộc tính / phương thức, `USE_CASES` = message của 34 sequence,
`DB_TABLES` = bảng CSDL). Sửa spec rồi sinh lại:

```bash
cd diagram/drawio
python3 gen_all_drawio.py                     # AN TOÀN: chỉ tạo file còn thiếu, KHÔNG ghi đè
python3 gen_all_drawio.py --force             # ghi đè tất cả (bản cũ tự sao lưu vào _backup/)
python3 gen_all_drawio.py --force --only 4.1  # chỉ sinh lại file tên bắt đầu bằng "4.1"
```

File `.drawio` chỉnh tay bằng https://app.diagrams.net hoặc draw.io desktop sẽ **không** bị
generator động vào (mặc định bỏ qua file đã tồn tại). Ngoài ra draw.io desktop tự lưu backup
ẩn dạng `.$tên.drawio.bkp` cạnh mỗi file - có thể khôi phục từ đó nếu lỡ tay.

## 3. Xuất PNG từ drawio

Cần draw.io CLI (`brew install --cask drawio` trên macOS; trên Windows cài draw.io Desktop và
dùng `draw.io.exe`):

```bash
cd diagram
for f in drawio/*.drawio; do
  b=$(basename "${f%.drawio}")
  drawio -x -f png -s 2 -o "png/$b.png" "$f"
done
```

Windows (PowerShell):

```powershell
cd diagram
Get-ChildItem drawio\*.drawio | ForEach-Object {
  & "C:\Program Files\draw.io\draw.io.exe" -x -f png -s 2 -o ("png\" + $_.BaseName + ".png") $_.FullName
}
```

## 4. Kiểm tra đồng bộ (bắt buộc trước khi nộp)

```bash
cd diagram/drawio
python3 check_sync.py            # exit 0 = đạt
```

Script kiểm: mọi message trong 34 sequence khớp phương thức khai báo của lớp đích; lifeline
không đặt tên chung chung; 14 lớp entity ánh xạ đủ 14 bảng CSDL; mọi PNG được nhúng trong
báo cáo và ngược lại; đủ 34 mục đặc tả (3.7.x) + 34 mục sequence (5.1.x) với tên UC trùng khớp.

## 5. Quy trình sửa nội dung

1. Sửa `BaoCao_QuanLyKhoHang.md` (nếu đổi tên UC / lớp / phương thức thì sửa cả SPEC trong
   `gen_all_drawio.py` rồi sinh lại drawio + PNG).
2. Chạy `check_sync.py` tới khi exit 0.
3. Build DOCX (mục 1), mở kiểm tra bằng mắt, xuất PDF.
