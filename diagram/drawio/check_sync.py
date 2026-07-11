#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Kiểm tra ĐỒNG BỘ chéo cho báo cáo v2 (biểu đồ drawio + Markdown).

Bản thích nghi của scripts/check-sync.py (skill pttkht-report): báo cáo v2 dùng
biểu đồ drawio/PNG thay vì mermaid, nên kiểm tra chạy trực tiếp trên SPEC
(CLASSES / USE_CASES / DB_TABLES / FILES trong gen_all_drawio.py) - cùng nguồn
sinh ra biểu đồ - và đối chiếu file Markdown báo cáo.

Kiểm tra:
  1) Sequence - Class: mọi message `method(...)` trong 34 sequence phải khớp một
     phương thức khai báo của lớp đích trong CLASSES.
  2) Lifeline không được là tên chung chung (UI/Controller/Hệ thống/CSDL...);
     phải là lớp cụ thể trong CLASSES hoặc actor.
  3) Class - CSDL: mọi lớp entity ánh xạ được tới một bảng trong DB_TABLES
     (Chương 7) và ngược lại.
  4) Markdown - PNG: mọi file diagram/png/*.png được tham chiếu trong báo cáo,
     và mọi ảnh tham chiếu trong báo cáo tồn tại trên đĩa.
  5) Markdown - đủ 34 mục đặc tả (3.7.x) và 34 mục sequence (5.1.x), tên UC
     ở hai nơi trùng khớp với spec.

Dùng:  python3 check_sync.py [đường_dẫn_báo_cáo.md]
Thoát mã 0 nếu đồng bộ, 1 nếu có lỗi.
"""
import re
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))
from gen_all_drawio import CLASSES, USE_CASES, DB_TABLES, FILES, ACTORS  # noqa: E402

MD_DEFAULT = HERE.parent.parent / "BaoCao_QuanLyKhoHang.md"
PNG_DIR = HERE.parent / "png"

ACTOR_NAMES = set(ACTORS) | {"Người dùng", "Tất cả người dùng"}
GENERIC = {"giaodien", "dieukhien", "hethong", "system", "ui", "controller",
           "database", "csdl", "nguoidung?", "actor", "boundary", "control", "entity"}


def norm(s: str) -> str:
    return re.sub(r"[^a-z0-9]", "", s.lower())


def main() -> int:
    errors: list[str] = []
    warnings: list[str] = []

    # 1) + 2) Sequence - Class
    declared = {name: {m.split("(")[0] for m in spec["methods"]} for name, spec in CLASSES.items()}
    for ucid in sorted(USE_CASES):
        uc = USE_CASES[ucid]
        for src, dst, label, kind in uc["messages"]:
            for p in (src, dst):
                if p in ACTOR_NAMES:
                    continue
                if p not in CLASSES:
                    errors.append(f"{ucid}: lifeline '{p}' không phải lớp khai báo trong CLASSES")
                elif norm(p) in GENERIC:
                    errors.append(f"{ucid}: lifeline '{p}' đặt tên chung chung")
            if kind != "call":
                continue
            m = re.match(r"(\w+)\s*\(", label)
            if not m:
                errors.append(f"{ucid}: message call '{label}' không có dạng method(...)")
                continue
            method = m.group(1)
            if dst in ACTOR_NAMES:
                errors.append(f"{ucid}: message call '{label}' trỏ tới actor '{dst}'")
            elif dst in declared and method not in declared[dst]:
                errors.append(f"{ucid}: '{method}()' không khai báo trong lớp {dst}")

    # 3) Class - CSDL
    entity_classes = [n for n, s in CLASSES.items() if s["kind"] == "entity"]
    tables_norm = {norm(t): t for t in DB_TABLES}
    for c in entity_classes:
        if norm(c) not in tables_norm:
            errors.append(f"Lớp entity '{c}' không ánh xạ tới bảng CSDL nào (DB_TABLES)")
    class_norm = {norm(c) for c in entity_classes}
    for t in DB_TABLES:
        if norm(t) not in class_norm:
            errors.append(f"Bảng '{t}' không có lớp entity tương ứng")

    # 4) + 5) Markdown
    md_path = Path(sys.argv[1]) if len(sys.argv) > 1 else MD_DEFAULT
    if not md_path.exists():
        warnings.append(f"Chưa có file báo cáo {md_path} - bỏ qua kiểm tra Markdown")
    else:
        md = md_path.read_text(encoding="utf-8")
        refs = set(re.findall(r"!\[[^\]]*\]\(diagram/png/([^)\s]+)\.png(?:\s+=[^)]*)?\)", md))
        for f in FILES:
            if f not in refs:
                errors.append(f"PNG 'diagram/png/{f}.png' không được nhúng trong báo cáo")
        for r in sorted(refs):
            if not (PNG_DIR / f"{r}.png").exists():
                errors.append(f"Báo cáo nhúng 'diagram/png/{r}.png' nhưng file không tồn tại")
        # đủ 34 đặc tả + 34 sequence, tên trùng spec
        for i in range(1, 35):
            ucid = f"UC{i:02d}"
            ten = USE_CASES[ucid]["ten"]
            spec_head = rf"^###\s+3\.7\.{i}\.?\s+{re.escape(ucid)}\s+-\s+{re.escape(ten)}\s*$"
            seq_head = rf"^###\s+5\.1\.{i}\.?\s+{re.escape(ucid)}\s+-\s+{re.escape(ten)}\s*$"
            if not re.search(spec_head, md, re.M):
                errors.append(f"Thiếu/sai tên mục đặc tả 3.7.{i} ({ucid} - {ten})")
            if not re.search(seq_head, md, re.M):
                errors.append(f"Thiếu/sai tên mục sequence 5.1.{i} ({ucid} - {ten})")

    for w in warnings:
        print(f"[WARN] {w}")
    if errors:
        print(f"\nFAIL - {len(errors)} lỗi đồng bộ:")
        for e in errors:
            print("  -", e)
        return 1
    n_msgs = sum(len(u["messages"]) for u in USE_CASES.values())
    print(f"OK - đồng bộ: {len(USE_CASES)} UC, {n_msgs} message, "
          f"{len(entity_classes)} lớp entity ↔ {len(DB_TABLES)} bảng, {len(FILES)} biểu đồ.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
