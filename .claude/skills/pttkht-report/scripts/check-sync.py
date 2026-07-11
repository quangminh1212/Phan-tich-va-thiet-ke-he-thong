#!/usr/bin/env python3
"""Kiểm tra ĐỒNG BỘ chéo giữa các biểu đồ trong báo cáo PTTKHT (Markdown + mermaid).
Bắt các lỗi review phổ biến nhất mà mắt thường/agent hay bỏ sót hoặc khai khống.

Kiểm tra:
  1) Sequence ↔ Class: mọi message dạng `A->>B: method(...)` trong sequenceDiagram
     phải khớp một phương thức `+method(...)`/`-method(...)` khai báo trong classDiagram
     thiết kế (classDiagram dài nhất = bản thiết kế chi tiết).
  2) Lớp thực thể ↔ Bảng CSDL: mọi class có thuộc tính nên ánh xạ tới một bảng erDiagram
     (so khớp gần đúng theo tên chuẩn hóa snake_case).

Bỏ qua: message không có dấu ngoặc `()` (nhãn/return thuần), và các nhãn return qua `-->>`.

Dùng: python3 check-sync.py <file.md>
Thoát mã 0 nếu đồng bộ, 1 nếu có lỗi (dùng làm gate trong quy trình).
"""
import re, sys

def norm(s):
    return re.sub(r'[^a-z0-9]', '', s.lower())

def main(path):
    md = open(path, encoding='utf-8').read()
    class_blocks = re.findall(r'```mermaid\s*\n(classDiagram.*?)```', md, re.S)
    seq_blocks = re.findall(r'```mermaid\s*\n(sequenceDiagram.*?)```', md, re.S)
    er_blocks = re.findall(r'```mermaid\s*\n(erDiagram.*?)```', md, re.S)

    errors = []

    # --- 1) Sequence ↔ Class ---
    if not class_blocks:
        errors.append("Không tìm thấy classDiagram nào.")
        declared = set()
    else:
        design = max(class_blocks, key=len)  # bản thiết kế chi tiết = dài nhất
        declared = set(re.findall(r'[+\-#]\s*(\w+)\s*\(', design))

    # Tên lifeline chung chung KHÔNG được chấp nhận — phải là tên lớp thiết kế cụ thể
    GENERIC = ['giao diện', 'giao dien', 'điều khiển', 'dieu khien', 'hệ thống', 'he thong',
               'system', 'ui', 'controller', 'database', 'csdl', 'người dùng', 'nguoi dung',
               'actor', 'boundary', 'control', 'entity']
    class_names = set(re.findall(r'class\s+(\w+)', max(class_blocks, key=len))) if class_blocks else set()

    generic_lifelines = set()
    for s in seq_blocks:
        for alias, label in re.findall(r'(?:participant|actor)\s+(\w+)\s+as\s+([^\n]+)', s):
            lab = label.strip()
            if norm(lab) in {norm(g) for g in GENERIC} and lab not in class_names:
                generic_lifelines.add(lab)
    if generic_lifelines:
        errors.append("Sequence: lifeline đặt tên CHUNG CHUNG (phải là tên lớp thiết kế cụ thể, "
                      "vd DieuKhienDatVe/BenhNhanService, không phải vai trò): "
                      + ", ".join(sorted(generic_lifelines)))

    called = []  # (method, snippet)
    for s in seq_blocks:
        for line in s.splitlines():
            m = re.search(r'(?<!-)->>\s*(\w+)\s*:\s*(\w+)\s*\(', line)
            if m:
                called.append((m.group(2), line.strip()))
    # Thông báo/return viết bằng mũi tên đặc ->> là sai ký pháp (phải dùng -->> và không gọi method)
    NOTIFY = re.compile(r'^(thongBao|hienThi|hopLe|chiTiet|ketQua|traVe)', re.I)
    bad_returns = sorted({meth for meth, _ in called if NOTIFY.match(meth)})
    if bad_returns:
        errors.append("Sequence: thông báo/return dùng mũi tên đặc `->>` kèm `()` — sai ký pháp; "
                      "dùng return `-->>` (nét đứt), không gọi như phương thức:\n    - "
                      + "\n    - ".join(bad_returns))
    missing = sorted({meth for meth, _ in called if meth not in declared and not NOTIFY.match(meth)})
    if missing:
        errors.append("Sequence↔Class: %d phương thức gọi trong sequence nhưng CHƯA khai báo "
                      "trong class thiết kế:\n    - " % len(missing) + "\n    - ".join(missing))

    # --- 2) Entity class ↔ Table ---
    if er_blocks:
        tables = set()
        for e in er_blocks:
            tables |= set(re.findall(r'^\s*([A-Za-z_]\w*)\s*\{', e, re.M))
        tnorm = {norm(t) for t in tables}
        # lớp thực thể = class có ≥1 thuộc tính (dòng bắt đầu bằng -/+ và KHÔNG có ngoặc)
        entity_classes = set()
        for c in class_blocks:
            for cm in re.finditer(r'class\s+(\w+)\s*\{(.*?)\}', c, re.S):
                name, body = cm.group(1), cm.group(2)
                has_attr = any(re.match(r'\s*[-+#]\s*\w+\s*:', l) or
                               (re.match(r'\s*[-+#]\s*\w+\s*$', l)) for l in body.splitlines())
                has_method = '(' in body
                if has_attr and not has_method:
                    entity_classes.add(name)
        unmapped = sorted(n for n in entity_classes if norm(n) not in tnorm
                          and not any(norm(n) in t or t in norm(n) for t in tnorm))
        if unmapped:
            errors.append("Class↔CSDL: lớp thực thể chưa thấy bảng tương ứng trong erDiagram: "
                          + ", ".join(unmapped))

    # --- Báo cáo ---
    print(f"classDiagram: {len(class_blocks)} | sequenceDiagram: {len(seq_blocks)} | "
          f"erDiagram: {len(er_blocks)}")
    print(f"Phương thức khai báo (class thiết kế): {len(declared)} | "
          f"message gọi (sequence): {len({m for m,_ in called})}")
    if errors:
        print("\n❌ PHÁT HIỆN LỖI ĐỒNG BỘ:\n")
        for e in errors:
            print("• " + e + "\n")
        print("=> Sửa: thêm phương thức còn thiếu vào class thiết kế, hoặc đổi message cho khớp.")
        return 1
    print("\n✅ Đồng bộ Sequence↔Class và Class↔CSDL: KHÔNG phát hiện lỗi.")
    return 0

if __name__ == '__main__':
    path = sys.argv[1] if len(sys.argv) > 1 else 'bao-cao.md'
    sys.exit(main(path))
