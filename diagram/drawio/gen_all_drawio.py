# -*- coding: utf-8 -*-
"""Sinh file .drawio cho báo cáo v2 - Hệ thống Quản lý Kho Hàng.

Spec-driven: CLASSES + USE_CASES định nghĩa một chỗ, dùng chung cho:
  - biểu đồ lớp (tổng thể + cắt lát)
  - biểu đồ trình tự (34 UC)
  - check đồng bộ Sequence-Class-CSDL (check_sync.py import file này)

AN TOÀN VỚI FILE SỬA TAY: mặc định CHỈ tạo file còn thiếu, KHÔNG ghi đè
file .drawio đã tồn tại (để không mất chỉnh sửa thủ công trong draw.io).

Chạy:
  python3 gen_all_drawio.py                     # chỉ tạo file còn thiếu
  python3 gen_all_drawio.py --force             # ghi đè TẤT CẢ (có backup)
  python3 gen_all_drawio.py --force --only 4.1  # chỉ ghi đè file tên bắt đầu "4.1"

Trước khi ghi đè, bản hiện tại được sao lưu vào _backup/<timestamp>/.
"""
from __future__ import annotations

import html
import shutil
import sys
import time
import uuid
from pathlib import Path

OUT = Path(__file__).resolve().parent

FORCE = "--force" in sys.argv
ONLY: list[str] = []
if "--only" in sys.argv:
    _i = sys.argv.index("--only")
    ONLY = [a for a in sys.argv[_i + 1:] if not a.startswith("--")]
_BACKUP_DIR = OUT / "_backup" / time.strftime("%Y%m%d-%H%M%S")
_SKIPPED: list[str] = []

# ───────────────────────── Styles ─────────────────────────
S_ACTOR = "shape=umlActor;verticalLabelPosition=bottom;verticalAlign=top;html=1;outlineConnect=0;fillColor=#FFFFFF;strokeColor=#000000;"
S_UC = "ellipse;whiteSpace=wrap;html=1;fillColor=#FFFFFF;strokeColor=#000000;fontSize=11;align=center;"
S_SYS = "rounded=1;whiteSpace=wrap;html=1;dashed=1;dashPattern=8 8;fillColor=none;strokeColor=#666666;verticalAlign=top;fontStyle=1;fontSize=13;align=center;spacingTop=8;"
S_RECT = "rounded=0;whiteSpace=wrap;html=1;fillColor=#FFFFFF;strokeColor=#000000;fontSize=11;align=center;"
S_ROUNDED = "rounded=1;whiteSpace=wrap;html=1;fillColor=#F5F5F5;strokeColor=#000000;fontSize=11;align=center;fontStyle=1;"
S_PROCESS = "rounded=1;whiteSpace=wrap;html=1;fillColor=#FFFFFF;strokeColor=#000000;fontSize=11;align=center;"
S_DECISION = "rhombus;whiteSpace=wrap;html=1;fillColor=#F0F0F0;strokeColor=#000000;fontSize=10;align=center;"
S_CLASS = "swimlane;fontStyle=1;align=center;verticalAlign=top;childLayout=stackLayout;horizontal=1;startSize=26;horizontalStack=0;resizeParent=1;resizeParentMax=0;resizeLast=0;collapsible=0;marginBottom=0;whiteSpace=wrap;html=1;fillColor=#FFFFFF;strokeColor=#000000;"
S_CLASS_ATTR = "text;strokeColor=none;fillColor=none;align=left;verticalAlign=top;spacingLeft=4;spacingRight=4;overflow=hidden;rotatable=0;points=[[0,0.5],[1,0.5]];portConstraint=eastwest;fontSize=10;"
S_CLASS_SEP = "line;strokeWidth=1;fillColor=none;align=left;verticalAlign=middle;spacingTop=-1;spacingLeft=3;spacingRight=3;rotatable=0;labelPosition=right;points=[];portConstraint=eastwest;strokeColor=#000000;"
S_EDGE = "endArrow=block;html=1;rounded=0;endFill=1;strokeColor=#000000;"
S_EDGE_OPEN = "endArrow=open;html=1;rounded=0;endFill=0;strokeColor=#000000;fontSize=10;"
S_EDGE_NONE = "endArrow=none;html=1;rounded=0;strokeColor=#000000;"
S_EDGE_INCLUDE = "endArrow=open;html=1;rounded=0;dashed=1;endFill=0;strokeColor=#000000;fontSize=10;"
S_EDGE_DEP = "endArrow=open;html=1;rounded=0;dashed=1;endFill=0;strokeColor=#000000;fontSize=10;"
S_MSG_CALL = "html=1;verticalAlign=bottom;endArrow=block;endFill=1;rounded=0;strokeColor=#000000;fontSize=10;"
S_MSG_RET = "html=1;verticalAlign=bottom;endArrow=open;endFill=0;dashed=1;rounded=0;strokeColor=#000000;fontSize=10;fontStyle=2;"
S_LIFELINE = "endArrow=none;dashed=1;html=1;rounded=0;strokeColor=#555555;"
S_POINT = "point;html=1;fillColor=none;strokeColor=none;"
S_STATE = "rounded=1;whiteSpace=wrap;html=1;fillColor=#FFFFFF;strokeColor=#000000;fontSize=12;align=center;arcSize=40;"
S_STATE_START = "ellipse;fillColor=#000000;strokeColor=#000000;"
S_STATE_END = "ellipse;shape=doubleEllipse;fillColor=#000000;strokeColor=#000000;"
S_CYL = "shape=cylinder3;whiteSpace=wrap;html=1;boundedLbl=1;backgroundOutline=1;size=12;fillColor=#F0F0F0;strokeColor=#000000;fontSize=10;align=center;"
S_NOTE = "text;html=1;strokeColor=none;fillColor=none;align=center;fontSize=12;fontStyle=1;"
S_COMPONENT = "shape=component;align=center;spacingLeft=10;whiteSpace=wrap;html=1;fillColor=#FFFFFF;strokeColor=#000000;fontSize=11;"


class D:
    def __init__(self, name: str, w: int = 1200, h: int = 900):
        self.name = name
        self.w = w
        self.h = h
        self.cells: list[str] = []
        self._n = 2

    def _id(self) -> str:
        i = str(self._n)
        self._n += 1
        return i

    def node(self, value: str, x: float, y: float, w: float, h: float, style: str, parent: str = "1") -> str:
        cid = self._id()
        v = html.escape(value).replace("\n", "&#xa;")
        self.cells.append(
            f'<mxCell id="{cid}" value="{v}" style="{style}" vertex="1" parent="{parent}">'
            f'<mxGeometry x="{x}" y="{y}" width="{w}" height="{h}" as="geometry"/></mxCell>'
        )
        return cid

    def edge(self, src: str, tgt: str, label: str = "", style: str = S_EDGE) -> str:
        cid = self._id()
        v = html.escape(label).replace("\n", "&#xa;") if label else ""
        self.cells.append(
            f'<mxCell id="{cid}" value="{v}" style="{style}" edge="1" parent="1" source="{src}" target="{tgt}">'
            f'<mxGeometry relative="1" as="geometry"/></mxCell>'
        )
        return cid

    def class_box(self, title: str, attrs: list[str], methods: list[str], x: float, y: float, w: float = 200) -> str:
        h_row = 18
        n_rows = max(len(attrs), 0) + max(len(methods), 0) + (1 if methods else 0)
        h = 26 + h_row * max(n_rows, 1) + 4
        cid = self._id()
        t = html.escape(title)
        self.cells.append(
            f'<mxCell id="{cid}" value="{t}" style="{S_CLASS}" vertex="1" parent="1">'
            f'<mxGeometry x="{x}" y="{y}" width="{w}" height="{h}" as="geometry"/></mxCell>'
        )
        rows = [("attr", a) for a in attrs]
        if methods:
            rows.append(("sep", ""))
            rows += [("attr", m) for m in methods]
        yoff = 26
        for kind, a in rows:
            rid = self._id()
            style = S_CLASS_SEP if kind == "sep" else S_CLASS_ATTR
            hh = 8 if kind == "sep" else h_row
            self.cells.append(
                f'<mxCell id="{rid}" value="{html.escape(a)}" style="{style}" vertex="1" parent="{cid}">'
                f'<mxGeometry x="0" y="{yoff}" width="{w}" height="{hh}" as="geometry"/></mxCell>'
            )
            yoff += hh
        return cid

    def save(self, filename: str):
        target = OUT / filename
        if target.exists():
            if not FORCE or (ONLY and not any(filename.startswith(p) for p in ONLY)):
                _SKIPPED.append(filename)
                return
            # ghi đè có chủ đích (--force): sao lưu bản hiện tại trước
            _BACKUP_DIR.mkdir(parents=True, exist_ok=True)
            shutil.copy2(target, _BACKUP_DIR / filename)
        body = "\n".join(self.cells)
        xml = f'''<?xml version="1.0" encoding="UTF-8"?>
<mxfile host="app.diagrams.net" agent="gen_all_drawio" version="24.0.0">
  <diagram id="{uuid.uuid4().hex[:8]}" name="{html.escape(self.name)}">
    <mxGraphModel dx="1200" dy="800" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="{self.w}" pageHeight="{self.h}" math="0" shadow="0">
      <root>
        <mxCell id="0"/>
        <mxCell id="1" parent="0"/>
{body}
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
'''
        (OUT / filename).write_text(xml, encoding="utf-8")
        print("wrote", filename)


# ════════════════════════ SPEC (nguồn sự thật) ════════════════════════

# Lớp: kind = entity | boundary | control ; entity ánh xạ bảng CSDL cùng tên
CLASSES: dict[str, dict] = {
    # ---- entity ----
    "NguoiDung": {
        "kind": "entity",
        "attrs": ["MaNguoiDung", "TenDangNhap", "MatKhau", "HoTen", "Email", "VaiTro", "TrangThai"],
        "methods": ["dangNhap(tenDangNhap, matKhau)", "dangXuat()", "themTaiKhoan(thongTin)",
                    "suaTaiKhoan(ma, thongTin)", "xoaTaiKhoan(ma)", "khoaMoTaiKhoan(ma, trangThai)",
                    "timKiemTaiKhoan(tuKhoa, boLoc)", "layDanhSachTaiKhoan()", "kiemTraQuyen(hanhDong)",
                    "kiemTraTrung(tenDangNhap)", "kiemTraRangBuoc(ma)"],
    },
    "NhaCungCap": {
        "kind": "entity",
        "attrs": ["MaNCC", "TenNCC", "DiaChi", "SDT", "Email", "NguoiLienHe", "TrangThai"],
        "methods": ["themMoi(thongTin)", "capNhat(ma, thongTin)", "xoa(ma)", "ngungHoatDong(ma)",
                    "timKiem(tuKhoa)", "layDanhSach()", "kiemTraTrung(ma)", "kiemTraRangBuoc(ma)"],
    },
    "HangHoa": {
        "kind": "entity",
        "attrs": ["MaHangHoa", "TenHangHoa", "MaNhomHang", "MaDonViTinh", "GiaNhap", "GiaXuat", "MoTa", "TrangThai"],
        "methods": ["themHangHoa(thongTin)", "capNhatThongTin(ma, thongTin)", "xoa(ma)", "ngungSuDung(ma)",
                    "timKiem(tuKhoa, boLoc)", "layDanhSach()", "layThongTin(ma)",
                    "kiemTraTrung(ma)", "kiemTraRangBuoc(ma)"],
    },
    "NhomHangHoa": {"kind": "entity", "attrs": ["MaNhomHang", "TenNhomHang", "MoTa"], "methods": ["layDanhSach()"]},
    "DonViTinh": {"kind": "entity", "attrs": ["MaDonViTinh", "TenDonViTinh"], "methods": ["layDanhSach()"]},
    "Kho": {
        "kind": "entity",
        "attrs": ["MaKho", "TenKho", "DiaChi", "NguoiPhuTrach", "TrangThai"],
        "methods": ["themMoi(thongTin)", "capNhat(ma, thongTin)", "xoa(ma)", "ngungHoatDong(ma)",
                    "timKiem(tuKhoa)", "layDanhSach()", "kiemTraTrung(ma)", "kiemTraRangBuoc(ma)"],
    },
    "TonKho": {
        "kind": "entity",
        "attrs": ["MaKho", "MaHangHoa", "SoLuongTon", "NgayCapNhat"],
        "methods": ["khoiTaoTon(maKho, maHang)", "tangTon(maKho, maHang, soLuong)",
                    "giamTon(maKho, maHang, soLuong)", "kiemTraDu(maKho, maHang, soLuong)",
                    "ganTheoKiemKe(maKho, maHang, slThucTe)", "layTon(maKho, maHang)"],
    },
    "PhieuNhapKho": {
        "kind": "entity",
        "attrs": ["MaPhieuNhap", "MaNCC", "MaKho", "MaNguoiDung", "NgayNhap", "TongTien", "TrangThai"],
        "methods": ["lapPhieu(thongTinChung)", "themChiTiet(dongHang)", "tinhTongTien()", "guiDuyet()",
                    "duyetPhieu()", "tuChoiPhieu(lyDo)", "capNhatPhieu(ma, thongTin)", "xoaPhieu(ma)",
                    "timKiem(dieuKien)", "layDanhSach()", "kiemTraTrangThai(ma)"],
    },
    "ChiTietPhieuNhap": {"kind": "entity", "attrs": ["MaPhieuNhap", "MaHangHoa", "SoLuong", "DonGia", "ThanhTien"],
                         "methods": ["tinhThanhTien()"]},
    "PhieuXuatKho": {
        "kind": "entity",
        "attrs": ["MaPhieuXuat", "MaKho", "MaNguoiDung", "NgayXuat", "LyDoXuat", "GhiChu", "TongTien", "TrangThai"],
        "methods": ["lapPhieu(thongTinChung)", "themChiTiet(dongHang)", "tinhTongTien()", "guiDuyet()",
                    "duyetPhieu()", "tuChoiPhieu(lyDo)", "capNhatPhieu(ma, thongTin)", "xoaPhieu(ma)",
                    "timKiem(dieuKien)", "layDanhSach()", "kiemTraTrangThai(ma)"],
    },
    "ChiTietPhieuXuat": {"kind": "entity", "attrs": ["MaPhieuXuat", "MaHangHoa", "SoLuong", "DonGia", "ThanhTien"],
                         "methods": ["tinhThanhTien()"]},
    "PhieuKiemKe": {
        "kind": "entity",
        "attrs": ["MaPhieuKiemKe", "MaKho", "MaNguoiDung", "NgayKiem", "TrangThai"],
        "methods": ["taoPhieu(maKho, ngay)", "nhapSoThucTe(maHang, soLuong)", "tinhChenhLech()",
                    "guiDuyet()", "duyetDieuChinh()", "tuChoiPhieu(lyDo)"],
    },
    "ChiTietKiemKe": {"kind": "entity", "attrs": ["MaPhieuKiemKe", "MaHangHoa", "SLThucTe", "SLSoSach", "ChenhLech"],
                      "methods": ["tinhChenhLech()"]},
    "LichSuThaoTac": {
        "kind": "entity",
        "attrs": ["MaLog", "MaNguoiDung", "HanhDong", "ThoiGian", "DoiTuong", "MoTa"],
        "methods": ["ghiLog(nguoiDung, hanhDong, doiTuong, moTa)"],
    },
    # ---- control đặc biệt: lớp báo cáo (không có bảng CSDL riêng) ----
    "BaoCaoNXT": {
        "kind": "control",
        "attrs": ["kyBaoCao", "maKho", "maHang"],
        "methods": ["tongHopNXT(dieuKien)", "xuatExcel(duLieu)"],
    },
}

# Boundary / Control theo nhóm chức năng
_B_COMMON = ["nhapThongTin(duLieu)", "chonBanGhi(ma)", "nhapDieuKien(dieuKien)",
             "chonChucNang(tenChucNang)", "hienThiKetQua(danhSach)"]
GROUPS: dict[str, dict] = {
    "xacthuc": {"boundary": "ManHinhDangNhap", "control": "DieuKhienXacThuc", "entity": "NguoiDung"},
    "taikhoan": {"boundary": "ManHinhTaiKhoan", "control": "DieuKhienTaiKhoan", "entity": "NguoiDung"},
    "ncc": {"boundary": "ManHinhNhaCungCap", "control": "DieuKhienNhaCungCap", "entity": "NhaCungCap"},
    "hanghoa": {"boundary": "ManHinhHangHoa", "control": "DieuKhienHangHoa", "entity": "HangHoa"},
    "kho": {"boundary": "ManHinhKho", "control": "DieuKhienKho", "entity": "Kho"},
    "phieunhap": {"boundary": "ManHinhPhieuNhap", "control": "DieuKhienPhieuNhap", "entity": "PhieuNhapKho"},
    "phieuxuat": {"boundary": "ManHinhPhieuXuat", "control": "DieuKhienPhieuXuat", "entity": "PhieuXuatKho"},
    "kiemke": {"boundary": "ManHinhKiemKe", "control": "DieuKhienKiemKe", "entity": "PhieuKiemKe"},
    "baocao": {"boundary": "ManHinhBaoCao", "control": "DieuKhienBaoCao", "entity": "BaoCaoNXT"},
}

_CONTROL_METHODS = {
    "DieuKhienXacThuc": ["yeuCauDangNhap(tenDangNhap, matKhau)", "yeuCauDangXuat()"],
    "DieuKhienTaiKhoan": ["yeuCauThem(thongTin)", "yeuCauSua(ma, thongTin)", "yeuCauXoa(ma)",
                          "yeuCauTimKiem(dieuKien)", "yeuCauDanhSach()"],
    "DieuKhienNhaCungCap": ["yeuCauThem(thongTin)", "yeuCauSua(ma, thongTin)", "yeuCauXoa(ma)",
                            "yeuCauTimKiem(dieuKien)", "yeuCauDanhSach()"],
    "DieuKhienHangHoa": ["yeuCauThem(thongTin)", "yeuCauSua(ma, thongTin)", "yeuCauXoa(ma)",
                         "yeuCauTimKiem(dieuKien)", "yeuCauDanhSach()"],
    "DieuKhienKho": ["yeuCauThem(thongTin)", "yeuCauSua(ma, thongTin)", "yeuCauXoa(ma)",
                     "yeuCauTimKiem(dieuKien)", "yeuCauDanhSach()"],
    "DieuKhienPhieuNhap": ["yeuCauLapPhieu(thongTinPhieu)", "yeuCauSua(ma, thongTin)", "yeuCauXoa(ma)",
                           "yeuCauTimKiem(dieuKien)", "yeuCauDanhSach()", "yeuCauDuyet(maPhieu)"],
    "DieuKhienPhieuXuat": ["yeuCauLapPhieu(thongTinPhieu)", "yeuCauSua(ma, thongTin)", "yeuCauXoa(ma)",
                           "yeuCauTimKiem(dieuKien)", "yeuCauDanhSach()", "yeuCauDuyet(maPhieu)"],
    "DieuKhienKiemKe": ["yeuCauTaoPhieu(maKho)", "yeuCauNhapThucTe(maHang, soLuong)", "yeuCauDuyet(maPhieu)"],
    "DieuKhienBaoCao": ["yeuCauBaoCao(dieuKien)"],
}

for g in GROUPS.values():
    b, c = g["boundary"], g["control"]
    if b not in CLASSES:
        extra = ["xacNhanDuyet(maPhieu)"] if "Phieu" in b or "KiemKe" in b else []
        extra += ["chonKho(maKho)"] if "KiemKe" in b else []
        CLASSES[b] = {"kind": "boundary", "attrs": [], "methods": _B_COMMON + extra}
    if c not in CLASSES:
        CLASSES[c] = {"kind": "control", "attrs": [], "methods": _CONTROL_METHODS[c]}

# Bảng CSDL (Chương 7) - entity nào cũng phải ánh xạ được
DB_TABLES = ["NguoiDung", "HangHoa", "NhomHangHoa", "DonViTinh", "NhaCungCap", "Kho",
             "PhieuNhapKho", "ChiTietPhieuNhap", "PhieuXuatKho", "ChiTietPhieuXuat",
             "PhieuKiemKe", "ChiTietKiemKe", "TonKho", "LichSuThaoTac"]

ACTORS = ["Quản trị viên", "Quản lý kho", "Nhân viên kho"]


# ---- helper xây message: (src, dst, label, kind) kind = call | ret ----
def _crud_uc(ucid, ten, actor, group, flavor, entity_ops):
    """entity_ops: dict flavor -> (list message tuples sau bước boundary->control, dacta bổ sung)"""
    B, C, E = GROUPS[group]["boundary"], GROUPS[group]["control"], GROUPS[group]["entity"]
    A = actor
    msgs, dacta = [], []
    if flavor == "them":
        msgs = [(A, B, "nhapThongTin(duLieu)", "call"),
                (B, C, "yeuCauThem(thongTin)", "call")] + entity_ops + [
               (C, "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
               (C, B, "ketQua", "ret"), (B, A, "thông báo thành công", "ret")]
        dacta = [f"{A} chọn \"Thêm mới\" trên màn hình, nhập thông tin bắt buộc và nhấn \"Lưu\".",
                 "Màn hình gửi yêu cầu thêm mới tới lớp điều khiển."]
    elif flavor == "sua":
        msgs = [(A, B, "chonBanGhi(ma)", "call"), (A, B, "nhapThongTin(duLieu)", "call"),
                (B, C, "yeuCauSua(ma, thongTin)", "call")] + entity_ops + [
               (C, "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
               (C, B, "ketQua", "ret"), (B, A, "thông báo thành công", "ret")]
        dacta = [f"{A} chọn bản ghi cần sửa trên danh sách và bấm \"Sửa\"; hệ thống nạp form với dữ liệu hiện tại.",
                 f"{A} chỉnh các trường cho phép rồi nhấn \"Lưu\".",
                 "Màn hình gửi yêu cầu cập nhật tới lớp điều khiển."]
    elif flavor == "xoa":
        msgs = [(A, B, "chonBanGhi(ma)", "call"),
                (B, C, "yeuCauXoa(ma)", "call")] + entity_ops + [
               (C, "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
               (C, B, "ketQua", "ret"), (B, A, "thông báo kết quả", "ret")]
        dacta = [f"{A} chọn bản ghi, bấm \"Xóa\" và xác nhận trên hộp thoại.",
                 "Màn hình gửi yêu cầu xóa tới lớp điều khiển."]
    elif flavor == "tim":
        msgs = [(A, B, "nhapDieuKien(dieuKien)", "call"),
                (B, C, "yeuCauTimKiem(dieuKien)", "call")] + entity_ops + [
               (C, B, "danhSachKetQua", "ret"), (B, A, "hiển thị kết quả", "ret")]
        dacta = [f"{A} nhập từ khóa / bộ lọc và nhấn \"Tìm kiếm\".",
                 "Màn hình gửi điều kiện tìm tới lớp điều khiển."]
    elif flavor == "xem":
        msgs = [(A, B, "chonChucNang(tenChucNang)", "call"),
                (B, C, "yeuCauDanhSach()", "call")] + entity_ops + [
               (C, B, "danhSach", "ret"), (B, A, "hiển thị danh sách", "ret")]
        dacta = [f"{A} chọn menu chức năng tương ứng.",
                 "Màn hình gửi yêu cầu tải danh sách tới lớp điều khiển."]
    return {"ten": ten, "actor": actor, "group": group, "flavor": flavor,
            "messages": msgs, "dacta_prefix": dacta}


def _catalog_group(group, ucs, actor, e, kw):
    """Sinh 5 UC CRUD cho một nhóm danh mục. kw: từ khóa tên UC (vd 'nhà cung cấp')."""
    C = GROUPS[group]["control"]
    them, sua, xoa, tim, xem = ucs
    ops = {
        "them": [(C, e, kw["kt_trung"], "call"), (C, e, kw["them"], "call"), (e, C, "ketQua", "ret")],
        "sua": [(C, e, kw["sua"], "call"), (e, C, "ketQua", "ret")],
        "xoa": [(C, e, kw["kt_rb"], "call"), (C, e, kw["xoa"], "call"), (C, e, kw["ngung"], "call"),
                (e, C, "ketQua", "ret")],
        "tim": [(C, e, kw["tim"], "call"), (e, C, "danhSachKetQua", "ret")],
        "xem": [(C, e, kw["xem"], "call"), (e, C, "danhSach", "ret")],
    }
    out = {}
    for ucid, flavor, ten in [(them, "them", kw["ten_them"]), (sua, "sua", kw["ten_sua"]),
                              (xoa, "xoa", kw["ten_xoa"]), (tim, "tim", kw["ten_tim"]),
                              (xem, "xem", kw["ten_xem"])]:
        out[ucid] = _crud_uc(ucid, ten, actor, group, flavor, ops[flavor])
    return out


USE_CASES: dict[str, dict] = {}

# UC01 / UC02 - xác thực
USE_CASES["UC01"] = {
    "ten": "Đăng nhập", "actor": "Người dùng", "group": "xacthuc", "flavor": "dacbiet",
    "messages": [
        ("Người dùng", "ManHinhDangNhap", "nhapThongTin(duLieu)", "call"),
        ("ManHinhDangNhap", "DieuKhienXacThuc", "yeuCauDangNhap(tenDangNhap, matKhau)", "call"),
        ("DieuKhienXacThuc", "NguoiDung", "dangNhap(tenDangNhap, matKhau)", "call"),
        ("NguoiDung", "DieuKhienXacThuc", "ketQuaXacThuc", "ret"),
        ("DieuKhienXacThuc", "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
        ("DieuKhienXacThuc", "ManHinhDangNhap", "ketQua / phiên làm việc", "ret"),
        ("ManHinhDangNhap", "Người dùng", "chuyển trang chủ theo vai trò", "ret"),
    ],
    "dacta_prefix": [],
}
USE_CASES["UC02"] = {
    "ten": "Đăng xuất", "actor": "Người dùng", "group": "xacthuc", "flavor": "dacbiet",
    "messages": [
        ("Người dùng", "ManHinhDangNhap", "chonChucNang(tenChucNang)", "call"),
        ("ManHinhDangNhap", "DieuKhienXacThuc", "yeuCauDangXuat()", "call"),
        ("DieuKhienXacThuc", "NguoiDung", "dangXuat()", "call"),
        ("DieuKhienXacThuc", "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
        ("DieuKhienXacThuc", "ManHinhDangNhap", "ketQua", "ret"),
        ("ManHinhDangNhap", "Người dùng", "chuyển về trang đăng nhập", "ret"),
    ],
    "dacta_prefix": [],
}

# UC03-07 tài khoản (QTV)
USE_CASES.update(_catalog_group(
    "taikhoan", ["UC03", "UC04", "UC05", "UC06", "UC07"], "Quản trị viên", "NguoiDung",
    {"ten_them": "Thêm tài khoản", "ten_sua": "Sửa tài khoản", "ten_xoa": "Xóa tài khoản",
     "ten_tim": "Tìm kiếm tài khoản", "ten_xem": "Xem danh sách tài khoản",
     "kt_trung": "kiemTraTrung(tenDangNhap)", "them": "themTaiKhoan(thongTin)",
     "sua": "suaTaiKhoan(ma, thongTin)", "kt_rb": "kiemTraRangBuoc(ma)",
     "xoa": "xoaTaiKhoan(ma)", "ngung": "khoaMoTaiKhoan(ma, trangThai)",
     "tim": "timKiemTaiKhoan(tuKhoa, boLoc)", "xem": "layDanhSachTaiKhoan()"}))

# UC08-12 NCC (NV, QL)
USE_CASES.update(_catalog_group(
    "ncc", ["UC08", "UC09", "UC10", "UC11", "UC12"], "Nhân viên kho", "NhaCungCap",
    {"ten_them": "Thêm nhà cung cấp", "ten_sua": "Sửa nhà cung cấp", "ten_xoa": "Xóa nhà cung cấp",
     "ten_tim": "Tìm kiếm nhà cung cấp", "ten_xem": "Xem danh sách nhà cung cấp",
     "kt_trung": "kiemTraTrung(ma)", "them": "themMoi(thongTin)", "sua": "capNhat(ma, thongTin)",
     "kt_rb": "kiemTraRangBuoc(ma)", "xoa": "xoa(ma)", "ngung": "ngungHoatDong(ma)",
     "tim": "timKiem(tuKhoa)", "xem": "layDanhSach()"}))

# UC13-17 hàng hóa (NV, QL) - thêm hàng có khởi tạo tồn
USE_CASES.update(_catalog_group(
    "hanghoa", ["UC13", "UC14", "UC15", "UC16", "UC17"], "Nhân viên kho", "HangHoa",
    {"ten_them": "Thêm hàng hóa", "ten_sua": "Sửa hàng hóa", "ten_xoa": "Xóa hàng hóa",
     "ten_tim": "Tìm kiếm hàng hóa", "ten_xem": "Xem danh sách hàng hóa",
     "kt_trung": "kiemTraTrung(ma)", "them": "themHangHoa(thongTin)",
     "sua": "capNhatThongTin(ma, thongTin)", "kt_rb": "kiemTraRangBuoc(ma)",
     "xoa": "xoa(ma)", "ngung": "ngungSuDung(ma)",
     "tim": "timKiem(tuKhoa, boLoc)", "xem": "layDanhSach()"}))
# UC13: chèn khởi tạo tồn sau themHangHoa
_m = USE_CASES["UC13"]["messages"]
_i = next(i for i, m in enumerate(_m) if m[2].startswith("themHangHoa"))
_m.insert(_i + 1, ("DieuKhienHangHoa", "TonKho", "khoiTaoTon(maKho, maHang)", "call"))
# UC16: tìm kiếm kèm đọc tồn
_m = USE_CASES["UC16"]["messages"]
_i = next(i for i, m in enumerate(_m) if m[2].startswith("timKiem"))
_m.insert(_i + 1, ("DieuKhienHangHoa", "TonKho", "layTon(maKho, maHang)", "call"))

# UC18-22 kho (QL)
USE_CASES.update(_catalog_group(
    "kho", ["UC18", "UC19", "UC20", "UC21", "UC22"], "Quản lý kho", "Kho",
    {"ten_them": "Thêm kho", "ten_sua": "Sửa kho", "ten_xoa": "Xóa kho",
     "ten_tim": "Tìm kiếm kho", "ten_xem": "Xem danh sách kho",
     "kt_trung": "kiemTraTrung(ma)", "them": "themMoi(thongTin)", "sua": "capNhat(ma, thongTin)",
     "kt_rb": "kiemTraRangBuoc(ma)", "xoa": "xoa(ma)", "ngung": "ngungHoatDong(ma)",
     "tim": "timKiem(tuKhoa)", "xem": "layDanhSach()"}))

# UC23-27 phiếu nhập
def _phieu_group(prefix, ucs, group, e, ten_phieu, ton_msg):
    B, C = GROUPS[group]["boundary"], GROUPS[group]["control"]
    lap, sua, xoa, tim, xem = ucs
    out = {}
    lap_msgs = [
        ("Nhân viên kho", B, "nhapThongTin(duLieu)", "call"),
        (B, C, "yeuCauLapPhieu(thongTinPhieu)", "call"),
        (C, e, "lapPhieu(thongTinChung)", "call"),
    ]
    if group == "phieuxuat":
        lap_msgs.append((C, "TonKho", "kiemTraDu(maKho, maHang, soLuong)", "call"))
        lap_msgs.append(("TonKho", C, "ketQuaKiemTra", "ret"))
    lap_msgs += [
        (C, e, "themChiTiet(dongHang)", "call"),
        (C, e, "tinhTongTien()", "call"),
        (C, e, "guiDuyet()", "call"),
        (e, C, "ketQua (Chờ duyệt)", "ret"),
        (C, B, "ketQua", "ret"),
        ("Quản lý kho", B, "xacNhanDuyet(maPhieu)", "call"),
        (B, C, "yeuCauDuyet(maPhieu)", "call"),
        (C, e, "duyetPhieu()", "call"),
        (e, "TonKho", ton_msg, "call"),
        (C, "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
        (C, B, "ketQua (Đã duyệt)", "ret"),
        (B, "Nhân viên kho", "thông báo kết quả", "ret"),
    ]
    out[lap] = {"ten": f"Lập {ten_phieu}", "actor": "Nhân viên kho", "group": group,
                "flavor": "lap", "messages": lap_msgs, "dacta_prefix": []}
    out[sua] = {"ten": f"Sửa {ten_phieu}", "actor": "Nhân viên kho", "group": group, "flavor": "sua",
                "messages": [
                    ("Nhân viên kho", B, "chonBanGhi(ma)", "call"),
                    ("Nhân viên kho", B, "nhapThongTin(duLieu)", "call"),
                    (B, C, "yeuCauSua(ma, thongTin)", "call"),
                    (C, e, "kiemTraTrangThai(ma)", "call"),
                    (C, e, "capNhatPhieu(ma, thongTin)", "call"),
                    (e, C, "ketQua", "ret"),
                    (C, "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
                    (C, B, "ketQua", "ret"),
                    (B, "Nhân viên kho", "thông báo kết quả", "ret")],
                "dacta_prefix": []}
    out[xoa] = {"ten": f"Xóa {ten_phieu}", "actor": "Nhân viên kho", "group": group, "flavor": "xoa",
                "messages": [
                    ("Nhân viên kho", B, "chonBanGhi(ma)", "call"),
                    (B, C, "yeuCauXoa(ma)", "call"),
                    (C, e, "kiemTraTrangThai(ma)", "call"),
                    (C, e, "xoaPhieu(ma)", "call"),
                    (e, C, "ketQua", "ret"),
                    (C, "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
                    (C, B, "ketQua", "ret"),
                    (B, "Nhân viên kho", "thông báo kết quả", "ret")],
                "dacta_prefix": []}
    out[tim] = {"ten": f"Tìm kiếm {ten_phieu}", "actor": "Nhân viên kho", "group": group, "flavor": "tim",
                "messages": [
                    ("Nhân viên kho", B, "nhapDieuKien(dieuKien)", "call"),
                    (B, C, "yeuCauTimKiem(dieuKien)", "call"),
                    (C, e, "timKiem(dieuKien)", "call"),
                    (e, C, "danhSachKetQua", "ret"),
                    (C, B, "danhSachKetQua", "ret"),
                    (B, "Nhân viên kho", "hiển thị kết quả", "ret")],
                "dacta_prefix": []}
    out[xem] = {"ten": f"Xem danh sách {ten_phieu}", "actor": "Nhân viên kho", "group": group, "flavor": "xem",
                "messages": [
                    ("Nhân viên kho", B, "chonChucNang(tenChucNang)", "call"),
                    (B, C, "yeuCauDanhSach()", "call"),
                    (C, e, "layDanhSach()", "call"),
                    (e, C, "danhSach", "ret"),
                    (C, B, "danhSach", "ret"),
                    (B, "Nhân viên kho", "hiển thị danh sách", "ret")],
                "dacta_prefix": []}
    return out


USE_CASES.update(_phieu_group("PN", ["UC23", "UC24", "UC25", "UC26", "UC27"], "phieunhap",
                              "PhieuNhapKho", "phiếu nhập kho", "tangTon(maKho, maHang, soLuong)"))
USE_CASES.update(_phieu_group("PX", ["UC28", "UC29", "UC30", "UC31", "UC32"], "phieuxuat",
                              "PhieuXuatKho", "phiếu xuất kho", "giamTon(maKho, maHang, soLuong)"))
# UC29 sửa phiếu xuất: kiểm tra tồn lại khi đổi số lượng
_m = USE_CASES["UC29"]["messages"]
_i = next(i for i, m in enumerate(_m) if m[2].startswith("kiemTraTrangThai"))
_m.insert(_i + 1, ("DieuKhienPhieuXuat", "TonKho", "kiemTraDu(maKho, maHang, soLuong)", "call"))

# Tên phiếu trong UC26/27, UC31/32 theo bảng tổng hợp: "phiếu nhập"/"phiếu xuất"
USE_CASES["UC26"]["ten"] = "Tìm kiếm phiếu nhập"
USE_CASES["UC27"]["ten"] = "Xem danh sách phiếu nhập"
USE_CASES["UC31"]["ten"] = "Tìm kiếm phiếu xuất"
USE_CASES["UC32"]["ten"] = "Xem danh sách phiếu xuất"

# UC33 kiểm kê
USE_CASES["UC33"] = {
    "ten": "Lập phiếu kiểm kê", "actor": "Nhân viên kho", "group": "kiemke", "flavor": "lap",
    "messages": [
        ("Nhân viên kho", "ManHinhKiemKe", "chonKho(maKho)", "call"),
        ("ManHinhKiemKe", "DieuKhienKiemKe", "yeuCauTaoPhieu(maKho)", "call"),
        ("DieuKhienKiemKe", "PhieuKiemKe", "taoPhieu(maKho, ngay)", "call"),
        ("PhieuKiemKe", "TonKho", "layTon(maKho, maHang)", "call"),
        ("TonKho", "PhieuKiemKe", "soLuongSoSach", "ret"),
        ("Nhân viên kho", "ManHinhKiemKe", "nhapThongTin(duLieu)", "call"),
        ("ManHinhKiemKe", "DieuKhienKiemKe", "yeuCauNhapThucTe(maHang, soLuong)", "call"),
        ("DieuKhienKiemKe", "PhieuKiemKe", "nhapSoThucTe(maHang, soLuong)", "call"),
        ("DieuKhienKiemKe", "PhieuKiemKe", "tinhChenhLech()", "call"),
        ("DieuKhienKiemKe", "PhieuKiemKe", "guiDuyet()", "call"),
        ("Quản lý kho", "ManHinhKiemKe", "xacNhanDuyet(maPhieu)", "call"),
        ("ManHinhKiemKe", "DieuKhienKiemKe", "yeuCauDuyet(maPhieu)", "call"),
        ("DieuKhienKiemKe", "PhieuKiemKe", "duyetDieuChinh()", "call"),
        ("PhieuKiemKe", "TonKho", "ganTheoKiemKe(maKho, maHang, slThucTe)", "call"),
        ("DieuKhienKiemKe", "LichSuThaoTac", "ghiLog(nguoiDung, hanhDong, doiTuong, moTa)", "call"),
        ("ManHinhKiemKe", "Nhân viên kho", "thông báo kết quả", "ret"),
    ],
    "dacta_prefix": [],
}

# UC34 báo cáo NXT
USE_CASES["UC34"] = {
    "ten": "Xem báo cáo nhập xuất tồn", "actor": "Quản lý kho", "group": "baocao", "flavor": "xem",
    "messages": [
        ("Quản lý kho", "ManHinhBaoCao", "nhapDieuKien(dieuKien)", "call"),
        ("ManHinhBaoCao", "DieuKhienBaoCao", "yeuCauBaoCao(dieuKien)", "call"),
        ("DieuKhienBaoCao", "BaoCaoNXT", "tongHopNXT(dieuKien)", "call"),
        ("BaoCaoNXT", "PhieuNhapKho", "layDanhSach()", "call"),
        ("BaoCaoNXT", "PhieuXuatKho", "layDanhSach()", "call"),
        ("BaoCaoNXT", "TonKho", "layTon(maKho, maHang)", "call"),
        ("BaoCaoNXT", "DieuKhienBaoCao", "duLieuBaoCao", "ret"),
        ("DieuKhienBaoCao", "ManHinhBaoCao", "duLieuBaoCao", "ret"),
        ("ManHinhBaoCao", "Quản lý kho", "hiển thị bảng báo cáo", "ret"),
    ],
    "dacta_prefix": [],
}

# Actor thực hiện theo bảng tổng hợp UC (dùng cho md + use case diagram)
UC_ACTORS = {}
for uc in ["UC01", "UC02"]:
    UC_ACTORS[uc] = "Tất cả người dùng"
for uc in ["UC03", "UC04", "UC05", "UC06", "UC07"]:
    UC_ACTORS[uc] = "Quản trị viên"
for uc in ["UC08", "UC09", "UC10", "UC11", "UC12", "UC13", "UC14", "UC15", "UC16", "UC17"]:
    UC_ACTORS[uc] = "Nhân viên kho, Quản lý kho"
for uc in ["UC18", "UC19", "UC20", "UC21", "UC22"]:
    UC_ACTORS[uc] = "Quản lý kho"
for uc in ["UC23", "UC24", "UC25", "UC26", "UC27", "UC28", "UC29", "UC30", "UC31", "UC32"]:
    UC_ACTORS[uc] = "Nhân viên kho, Quản lý kho"
UC_ACTORS["UC33"] = "Nhân viên kho (lập), Quản lý kho (duyệt)"
UC_ACTORS["UC34"] = "Quản lý kho"


# ════════════════════════ Renderers ════════════════════════

def _wrap_label(label: str, maxlen: int) -> str:
    """Nhãn mũi tên quá dài: tách làm 2 dòng tại khoảng trắng gần giữa nhất
    (không cắt quá sát hai đầu để tránh dòng cụt kiểu '1:')."""
    if len(label) <= maxlen:
        return label
    mid = len(label) // 2
    spaces = [i for i, ch in enumerate(label) if ch == " " and 8 <= i <= len(label) - 6]
    if not spaces:
        return label
    cut = min(spaces, key=lambda i: abs(i - mid))
    return label[:cut] + "\n" + label[cut + 1:]


def seq_diagram(title, filename, messages):
    """Vẽ sequence UML: lifeline theo thứ tự xuất hiện, actor = hình người.
    Biểu đồ nhiều participant (>5) dùng layout gọn: lifeline sát nhau hơn,
    hộp tên nhỏ hơn, nhãn mũi tên dài tự xuống dòng."""
    order: list[str] = []
    for s, t, _, _ in messages:
        for p in (s, t):
            if p not in order:
                order.append(p)
    is_actor = {p: (p in ACTORS or p in ("Người dùng", "Tất cả người dùng")) for p in order}
    compact = len(order) > 5
    span = 150 if compact else 200
    head_w = 120 if compact else 160
    wrap_at = 24 if compact else 34
    x0 = 40 if compact else 60
    top = 40
    msg_y0 = top + 90
    step = 46
    h = msg_y0 + step * len(messages) + 60
    w = x0 + span * len(order) + 40
    d = D(title, w, h)
    centers = {}
    for i, p in enumerate(order):
        cx = x0 + i * span + head_w // 2
        centers[p] = cx
        if is_actor[p]:
            d.node(p, cx - 20, top, 40, 60, S_ACTOR)
            head_bottom = top + 78
        else:
            d.node(p, cx - head_w // 2, top, head_w, 40, S_RECT)
            head_bottom = top + 40
        a1 = d.node("", cx, head_bottom, 1, 1, S_POINT)
        a2 = d.node("", cx, h - 30, 1, 1, S_POINT)
        d.edge(a1, a2, "", S_LIFELINE)
    y = msg_y0
    n = 0
    for s, t, label, kind in messages:
        if kind == "call":
            n += 1
            lab = f"{n}: {label}"
            style = S_MSG_CALL
        else:
            lab = label
            style = S_MSG_RET
        lab = _wrap_label(lab, wrap_at)
        p1 = d.node("", centers[s], y, 1, 1, S_POINT)
        p2 = d.node("", centers[t], y, 1, 1, S_POINT)
        d.edge(p1, p2, lab, style)
        y += step
    d.save(filename)


def class_slice(title, filename, boxes, edges, notes=None, width=1150):
    """boxes: list (ClassName, x, y, w, max_methods hoặc None = tất cả)"""
    d = D(title, width, 760)
    ids = {}
    for name, x, y, w, maxm in boxes:
        spec = CLASSES[name]
        methods = ["+ " + m for m in spec["methods"]]
        if maxm is not None:
            methods = methods[:maxm]
        attrs = ["+ " + a for a in spec["attrs"]]
        stereo = {"boundary": "«boundary»\n", "control": "«control»\n", "entity": ""}[spec["kind"]]
        ids[name] = d.class_box(stereo + name, attrs, methods, x, y, w)
    for a, b, label, style in edges:
        d.edge(ids[a], ids[b], label, style)
    for text, x, y, w in (notes or []):
        d.node(text, x, y, w, 30, S_NOTE)
    d.save(filename)


# ───────────────────────── Chương 2: Activity ─────────────────────────

def gen_activities():
    def flow(title, filename, labels, branches, w=820, h=1150):
        d = D(title, w, h)
        y, x = 30, 260
        nodes = []
        for text, style, bw, bh, lab in labels:
            n = d.node(text, x + (260 - bw) / 2, y, bw, bh, style)
            if nodes:
                d.edge(nodes[-1], n, lab)
            nodes.append(n)
            y += bh + 28
        for dec_idx, err_text, back_idx, ey in branches:
            err = d.node(err_text, 570, ey, 200, 44, S_PROCESS)
            d.edge(nodes[dec_idx], err, "Không")
            d.edge(err, nodes[back_idx])
        d.save(filename)

    flow("Activity - Nghiệp vụ nhập kho", "2.1_activity_nhapkho.drawio", [
        ("Bắt đầu", S_ROUNDED, 140, 40, ""),
        ("Nhân viên kho lập phiếu nhập", S_PROCESS, 260, 44, ""),
        ("Nhập NCC, kho, ngày nhập", S_PROCESS, 260, 44, ""),
        ("Thêm chi tiết hàng hóa nhập", S_PROCESS, 260, 44, ""),
        ("Dữ liệu hợp lệ?", S_DECISION, 150, 70, ""),
        ("Gửi phiếu nhập để duyệt", S_PROCESS, 260, 44, "Có"),
        ("Quản lý kho duyệt?", S_DECISION, 150, 70, ""),
        ("Cập nhật tồn kho (tăng)", S_PROCESS, 260, 44, "Có"),
        ("In / lưu phiếu nhập", S_PROCESS, 260, 44, ""),
        ("Kết thúc", S_ROUNDED, 140, 40, ""),
    ], [(4, "Báo lỗi, yêu cầu nhập lại", 2, 300), (6, "Từ chối, trả lại nhân viên", 1, 520)])

    flow("Activity - Nghiệp vụ xuất kho", "2.2_activity_xuatkho.drawio", [
        ("Bắt đầu", S_ROUNDED, 140, 40, ""),
        ("Nhân viên kho lập phiếu xuất", S_PROCESS, 260, 44, ""),
        ("Nhập kho xuất, ngày, lý do / ghi chú", S_PROCESS, 260, 44, ""),
        ("Thêm chi tiết hàng xuất", S_PROCESS, 260, 44, ""),
        ("Tồn kho đủ?", S_DECISION, 150, 70, ""),
        ("Gửi phiếu xuất để duyệt", S_PROCESS, 260, 44, "Có"),
        ("Quản lý kho duyệt?", S_DECISION, 150, 70, ""),
        ("Cập nhật tồn kho (giảm)", S_PROCESS, 260, 44, "Có"),
        ("In / lưu phiếu xuất", S_PROCESS, 260, 44, ""),
        ("Kết thúc", S_ROUNDED, 140, 40, ""),
    ], [(4, "Báo tồn không đủ, sửa lại dòng hàng", 3, 300), (6, "Từ chối, trả lại nhân viên", 1, 520)])

    flow("Activity - Nghiệp vụ kiểm kê", "2.3_activity_kiemke.drawio", [
        ("Bắt đầu", S_ROUNDED, 140, 40, ""),
        ("Tạo phiếu kiểm kê theo kho", S_PROCESS, 260, 44, ""),
        ("Hệ thống nạp SL theo sổ (TonKho)", S_PROCESS, 260, 44, ""),
        ("Nhân viên nhập SL thực tế", S_PROCESS, 260, 44, ""),
        ("Tính chênh lệch", S_PROCESS, 260, 44, ""),
        ("Gửi duyệt", S_PROCESS, 260, 44, ""),
        ("Quản lý duyệt điều chỉnh?", S_DECISION, 170, 70, ""),
        ("Cập nhật TonKho = số thực tế", S_PROCESS, 260, 44, "Có"),
        ("Kết thúc", S_ROUNDED, 140, 40, ""),
    ], [(6, "Yêu cầu kiểm đếm lại", 3, 470)])


# ───────────────────────── Chương 3: Use case ─────────────────────────

def gen_usecases():
    # 3.1 tổng quát
    d = D("Use case tổng quát", 1150, 860)
    d.node("Hệ thống quản lý kho hàng", 220, 40, 700, 660, S_SYS)
    left = [(300, 100, "Đăng nhập / Đăng xuất"), (300, 190, "Quản lý tài khoản"),
            (300, 280, "Quản lý nhà cung cấp"), (300, 370, "Quản lý hàng hóa"),
            (300, 460, "Quản lý kho")]
    right = [(640, 100, "Nhập kho"), (640, 190, "Xuất kho"), (640, 280, "Kiểm kê tồn kho"),
             (640, 370, "Báo cáo nhập xuất tồn")]
    lids = [d.node(t, x, y, 200, 50, S_UC) for x, y, t in left]
    rids = [d.node(t, x, y, 200, 50, S_UC) for x, y, t in right]
    admin = d.node("Quản trị viên", 60, 230, 70, 90, S_ACTOR)
    ql = d.node("Quản lý kho", 1010, 230, 70, 90, S_ACTOR)
    nv = d.node("Nhân viên kho", 530, 740, 70, 90, S_ACTOR)
    d.edge(admin, lids[0], style=S_EDGE_NONE)
    d.edge(admin, lids[1], style=S_EDGE_NONE)
    for i in [0, 2, 3, 4]:
        d.edge(ql, lids[i], style=S_EDGE_NONE)
    for r in rids:
        d.edge(ql, r, style=S_EDGE_NONE)
    for i in [0, 2, 3]:
        d.edge(nv, lids[i], style=S_EDGE_NONE)
    for r in rids[:3]:
        d.edge(nv, r, style=S_EDGE_NONE)
    d.save("3.1_usecase_tongquat.drawio")

    # 3.2 tài khoản + xác thực
    d = D("UC phân rã: truy cập và tài khoản", 1050, 780)
    d.node("Nhóm truy cập và tài khoản", 170, 30, 700, 640, S_SYS)
    items = [(220, 100, "UC03 Thêm tài khoản"), (220, 180, "UC04 Sửa tài khoản"),
             (220, 260, "UC05 Xóa tài khoản"), (220, 340, "UC06 Tìm kiếm tài khoản"),
             (220, 420, "UC07 Xem danh sách tài khoản")]
    rights = [(590, 160, "UC01 Đăng nhập"), (590, 280, "UC02 Đăng xuất")]
    ids = [d.node(t, x, y, 230, 48, S_UC) for x, y, t in items]
    rids = [d.node(t, x, y, 210, 48, S_UC) for x, y, t in rights]
    admin = d.node("Quản trị viên", 40, 260, 70, 90, S_ACTOR)
    user = d.node("Người dùng", 930, 210, 70, 90, S_ACTOR)
    for i in ids:
        d.edge(admin, i, style=S_EDGE_NONE)
        d.edge(i, rids[0], "«include»", S_EDGE_INCLUDE)
    d.edge(user, rids[0], style=S_EDGE_NONE)
    d.edge(user, rids[1], style=S_EDGE_NONE)
    d.save("3.2_usecase_taikhoan.drawio")

    # 3.3 NCC
    d = D("UC phân rã: nhà cung cấp", 1000, 720)
    d.node("Nhóm đối tác (nhà cung cấp)", 190, 30, 600, 600, S_SYS)
    items = ["UC08 Thêm nhà cung cấp", "UC09 Sửa nhà cung cấp", "UC10 Xóa nhà cung cấp",
             "UC11 Tìm kiếm nhà cung cấp", "UC12 Xem danh sách nhà cung cấp"]
    ids = []
    y = 80
    for t in items:
        ids.append(d.node(t, 280, y, 250, 48, S_UC))
        y += 76
    inc = d.node("UC01 Đăng nhập", 600, 480, 170, 44, S_UC)
    nv = d.node("Nhân viên kho", 50, 250, 70, 90, S_ACTOR)
    ql = d.node("Quản lý kho", 880, 250, 70, 90, S_ACTOR)
    for i in ids:
        d.edge(nv, i, style=S_EDGE_NONE)
        d.edge(ql, i, style=S_EDGE_NONE)
    d.edge(ids[0], inc, "«include»", S_EDGE_INCLUDE)
    d.save("3.3_usecase_ncc.drawio")

    # 3.4 hàng hóa + kho
    d = D("UC phân rã: hàng hóa và kho", 1150, 780)
    d.node("Nhóm hàng hóa và kho", 140, 20, 860, 650, S_SYS)
    d.node("Hàng hóa", 250, 55, 100, 24, S_NOTE)
    d.node("Kho", 680, 55, 60, 24, S_NOTE)
    h = ["UC13 Thêm hàng hóa", "UC14 Sửa hàng hóa", "UC15 Xóa hàng hóa",
         "UC16 Tìm kiếm hàng hóa", "UC17 Xem danh sách hàng hóa"]
    k = ["UC18 Thêm kho", "UC19 Sửa kho", "UC20 Xóa kho", "UC21 Tìm kiếm kho", "UC22 Xem danh sách kho"]
    hids, kids = [], []
    y = 95
    for t in h:
        hids.append(d.node(t, 200, y, 220, 46, S_UC))
        y += 76
    y = 95
    for t in k:
        kids.append(d.node(t, 620, y, 220, 46, S_UC))
        y += 76
    nv = d.node("Nhân viên kho", 40, 300, 70, 90, S_ACTOR)
    ql = d.node("Quản lý kho", 1030, 300, 70, 90, S_ACTOR)
    for i in hids:
        d.edge(nv, i, style=S_EDGE_NONE)
        d.edge(ql, i, style=S_EDGE_NONE)
    for i in kids:
        d.edge(ql, i, style=S_EDGE_NONE)
    d.save("3.4_usecase_hang_kho.drawio")

    # 3.5 phiếu + kiểm kê + báo cáo
    d = D("UC phân rã: nhập, xuất, kiểm kê, báo cáo", 1250, 880)
    d.node("Nhóm nhập, xuất, kiểm kê, báo cáo", 120, 15, 1000, 760, S_SYS)
    d.node("Phiếu nhập", 250, 50, 100, 22, S_NOTE)
    d.node("Phiếu xuất", 760, 50, 100, 22, S_NOTE)
    pn = ["UC23 Lập phiếu nhập kho", "UC24 Sửa phiếu nhập kho", "UC25 Xóa phiếu nhập kho",
          "UC26 Tìm kiếm phiếu nhập", "UC27 Xem danh sách phiếu nhập"]
    px = ["UC28 Lập phiếu xuất kho", "UC29 Sửa phiếu xuất kho", "UC30 Xóa phiếu xuất kho",
          "UC31 Tìm kiếm phiếu xuất", "UC32 Xem danh sách phiếu xuất"]
    pids, xids = [], []
    y = 90
    for t in pn:
        pids.append(d.node(t, 200, y, 230, 44, S_UC))
        y += 70
    y = 90
    for t in px:
        xids.append(d.node(t, 710, y, 230, 44, S_UC))
        y += 70
    kk = d.node("UC33 Lập phiếu kiểm kê", 200, 480, 230, 44, S_UC)
    bc = d.node("UC34 Xem báo cáo nhập xuất tồn", 710, 480, 260, 44, S_UC)
    nv = d.node("Nhân viên kho", 40, 260, 70, 90, S_ACTOR)
    ql = d.node("Quản lý kho", 1140, 260, 70, 90, S_ACTOR)
    for i in pids + xids + [kk]:
        d.edge(nv, i, style=S_EDGE_NONE)
        d.edge(ql, i, style=S_EDGE_NONE)
    d.edge(ql, bc, style=S_EDGE_NONE)
    d.save("3.5_usecase_phieu_baocao.drawio")


# ───────────────────────── Chương 4: Class ─────────────────────────

def gen_classes_overview():
    # 4.1 master data (entity thuần)
    d = D("Class - Nhóm danh mục (Master Data)", 1200, 820)
    def eb(name, x, y, w=200, maxm=3):
        spec = CLASSES[name]
        return d.class_box(name, ["+ " + a for a in spec["attrs"]],
                           ["+ " + m for m in spec["methods"][:maxm]], x, y, w)
    c_nhom = eb("NhomHangHoa", 40, 60, 190, 1)
    c_hh = eb("HangHoa", 320, 40, 230, 3)
    c_dvt = eb("DonViTinh", 640, 60, 190, 1)
    c_kho = eb("Kho", 40, 380, 200, 3)
    c_nd = eb("NguoiDung", 320, 360, 230, 2)
    c_ncc = eb("NhaCungCap", 640, 380, 210, 3)
    c_ton = eb("TonKho", 920, 200, 220, 4)
    d.edge(c_nhom, c_hh, "1..*", S_EDGE_OPEN)
    d.edge(c_dvt, c_hh, "1..*", S_EDGE_OPEN)
    d.edge(c_hh, c_ton, "1..*", S_EDGE_OPEN)
    d.edge(c_kho, c_ton, "1..*", S_EDGE_OPEN)
    d.save("4.1_class_master.drawio")

    # 4.2 nhập xuất
    d = D("Class - Nhóm nghiệp vụ nhập xuất", 1200, 820)
    pn = eb2(d, "PhieuNhapKho", 60, 40, 240, 6)
    ctn = eb2(d, "ChiTietPhieuNhap", 430, 60, 210, 1)
    px = eb2(d, "PhieuXuatKho", 60, 420, 240, 6)
    ctx = eb2(d, "ChiTietPhieuXuat", 430, 440, 210, 1)
    ncc = eb2(d, "NhaCungCap", 750, 40, 200, 0)
    kho = eb2(d, "Kho", 750, 280, 190, 0)
    hh = eb2(d, "HangHoa", 750, 480, 210, 0)
    ton = eb2(d, "TonKho", 990, 280, 200, 4)
    nd = eb2(d, "NguoiDung", 430, 260, 200, 0)
    d.edge(pn, ctn, "1..*", S_EDGE_OPEN)
    d.edge(px, ctx, "1..*", S_EDGE_OPEN)
    d.edge(ncc, pn, "1..*", S_EDGE_OPEN)
    d.edge(kho, pn, "1..*", S_EDGE_OPEN)
    d.edge(kho, px, "1..*", S_EDGE_OPEN)
    d.edge(nd, pn, "1..*", S_EDGE_OPEN)
    d.edge(nd, px, "1..*", S_EDGE_OPEN)
    d.edge(ctn, hh, "n..1", S_EDGE_OPEN)
    d.edge(ctx, hh, "n..1", S_EDGE_OPEN)
    d.edge(kho, ton, "1..*", S_EDGE_OPEN)
    d.edge(hh, ton, "1..*", S_EDGE_OPEN)
    d.save("4.2_class_nhapxuat.drawio")

    # 4.3 kiểm kê + lịch sử
    d = D("Class - Nhóm kiểm kê và lịch sử", 1050, 640)
    pk = eb2(d, "PhieuKiemKe", 60, 60, 240, 6)
    ct = eb2(d, "ChiTietKiemKe", 420, 80, 220, 1)
    ton = eb2(d, "TonKho", 720, 80, 210, 4)
    log = eb2(d, "LichSuThaoTac", 420, 380, 240, 1)
    nd = eb2(d, "NguoiDung", 60, 400, 220, 0)
    d.edge(pk, ct, "1..*", S_EDGE_OPEN)
    d.edge(ct, ton, "đối chiếu", S_EDGE_DEP)
    d.edge(nd, pk, "1..*", S_EDGE_OPEN)
    d.edge(nd, log, "1..*", S_EDGE_OPEN)
    d.save("4.3_class_kiemke_lichsu.drawio")


def eb2(d: D, name, x, y, w=200, maxm=3):
    spec = CLASSES[name]
    methods = ["+ " + m for m in spec["methods"][:maxm]] if maxm else []
    return d.class_box(name, ["+ " + a for a in spec["attrs"]], methods, x, y, w)


def gen_class_slices():
    dep = S_EDGE_DEP
    asc = S_EDGE_OPEN
    class_slice("Class cắt lát - Tài khoản / xác thực (UC01-UC07)", "4.4_class_taikhoan.drawio",
                [("ManHinhTaiKhoan", 40, 60, 230, 4), ("DieuKhienTaiKhoan", 340, 60, 240, None),
                 ("NguoiDung", 660, 40, 260, None), ("LichSuThaoTac", 660, 480, 240, 1)],
                [("ManHinhTaiKhoan", "DieuKhienTaiKhoan", "", dep),
                 ("DieuKhienTaiKhoan", "NguoiDung", "", dep),
                 ("DieuKhienTaiKhoan", "LichSuThaoTac", "", dep),
                 ("NguoiDung", "LichSuThaoTac", "1..*", asc)])
    class_slice("Class cắt lát - Nhà cung cấp (UC08-UC12)", "4.5_class_ncc.drawio",
                [("ManHinhNhaCungCap", 40, 60, 230, 4), ("DieuKhienNhaCungCap", 340, 60, 240, None),
                 ("NhaCungCap", 660, 40, 250, None), ("PhieuNhapKho", 660, 420, 240, 0),
                 ("LichSuThaoTac", 340, 420, 240, 1)],
                [("ManHinhNhaCungCap", "DieuKhienNhaCungCap", "", dep),
                 ("DieuKhienNhaCungCap", "NhaCungCap", "", dep),
                 ("DieuKhienNhaCungCap", "LichSuThaoTac", "", dep),
                 ("NhaCungCap", "PhieuNhapKho", "1..*", asc)])
    class_slice("Class cắt lát - Hàng hóa (UC13-UC17)", "4.6_class_hanghoa.drawio",
                [("ManHinhHangHoa", 40, 60, 230, 4), ("DieuKhienHangHoa", 340, 60, 240, None),
                 ("HangHoa", 660, 40, 260, None), ("TonKho", 960, 60, 220, 4),
                 ("NhomHangHoa", 660, 480, 200, 1), ("DonViTinh", 900, 480, 190, 1),
                 ("LichSuThaoTac", 340, 480, 240, 1)],
                [("ManHinhHangHoa", "DieuKhienHangHoa", "", dep),
                 ("DieuKhienHangHoa", "HangHoa", "", dep),
                 ("DieuKhienHangHoa", "TonKho", "", dep),
                 ("DieuKhienHangHoa", "LichSuThaoTac", "", dep),
                 ("NhomHangHoa", "HangHoa", "1..*", asc),
                 ("DonViTinh", "HangHoa", "1..*", asc),
                 ("HangHoa", "TonKho", "1..*", asc)], width=1250)
    class_slice("Class cắt lát - Kho (UC18-UC22)", "4.7_class_kho.drawio",
                [("ManHinhKho", 40, 60, 230, 4), ("DieuKhienKho", 340, 60, 240, None),
                 ("Kho", 660, 40, 240, None), ("TonKho", 960, 60, 220, 4),
                 ("LichSuThaoTac", 340, 440, 240, 1)],
                [("ManHinhKho", "DieuKhienKho", "", dep),
                 ("DieuKhienKho", "Kho", "", dep),
                 ("DieuKhienKho", "LichSuThaoTac", "", dep),
                 ("Kho", "TonKho", "1..*", asc)], width=1250)
    class_slice("Class cắt lát - Phiếu nhập (UC23-UC27)", "4.8_class_phieunhap.drawio",
                [("ManHinhPhieuNhap", 40, 60, 240, 5), ("DieuKhienPhieuNhap", 350, 60, 250, None),
                 ("PhieuNhapKho", 680, 40, 260, None), ("ChiTietPhieuNhap", 990, 60, 210, 1),
                 ("NhaCungCap", 40, 460, 210, 0), ("Kho", 300, 460, 190, 0),
                 ("HangHoa", 540, 500, 210, 0), ("TonKho", 800, 500, 210, 4),
                 ("LichSuThaoTac", 1050, 500, 200, 1)],
                [("ManHinhPhieuNhap", "DieuKhienPhieuNhap", "", dep),
                 ("DieuKhienPhieuNhap", "PhieuNhapKho", "", dep),
                 ("DieuKhienPhieuNhap", "LichSuThaoTac", "", dep),
                 ("PhieuNhapKho", "ChiTietPhieuNhap", "1..*", asc),
                 ("NhaCungCap", "PhieuNhapKho", "1..*", asc),
                 ("Kho", "PhieuNhapKho", "1..*", asc),
                 ("ChiTietPhieuNhap", "HangHoa", "n..1", asc),
                 ("PhieuNhapKho", "TonKho", "tangTon khi duyệt", dep)], width=1300)
    class_slice("Class cắt lát - Phiếu xuất (UC28-UC32)", "4.9_class_phieuxuat.drawio",
                [("ManHinhPhieuXuat", 40, 60, 240, 5), ("DieuKhienPhieuXuat", 350, 60, 250, None),
                 ("PhieuXuatKho", 680, 40, 260, None), ("ChiTietPhieuXuat", 990, 60, 210, 1),
                 ("Kho", 300, 480, 190, 0), ("HangHoa", 540, 520, 210, 0),
                 ("TonKho", 800, 520, 210, 4), ("LichSuThaoTac", 1050, 520, 200, 1)],
                [("ManHinhPhieuXuat", "DieuKhienPhieuXuat", "", dep),
                 ("DieuKhienPhieuXuat", "PhieuXuatKho", "", dep),
                 ("DieuKhienPhieuXuat", "LichSuThaoTac", "", dep),
                 ("PhieuXuatKho", "ChiTietPhieuXuat", "1..*", asc),
                 ("Kho", "PhieuXuatKho", "1..*", asc),
                 ("ChiTietPhieuXuat", "HangHoa", "n..1", asc),
                 ("PhieuXuatKho", "TonKho", "kiemTraDu / giamTon", dep)],
                width=1300)
    class_slice("Class cắt lát - Kiểm kê (UC33)", "4.10_class_kiemke.drawio",
                [("ManHinhKiemKe", 40, 60, 240, 5), ("DieuKhienKiemKe", 350, 60, 260, None),
                 ("PhieuKiemKe", 690, 40, 250, None), ("ChiTietKiemKe", 990, 60, 220, 1),
                 ("Kho", 350, 440, 190, 0), ("TonKho", 690, 440, 220, 4),
                 ("LichSuThaoTac", 990, 440, 200, 1)],
                [("ManHinhKiemKe", "DieuKhienKiemKe", "", dep),
                 ("DieuKhienKiemKe", "PhieuKiemKe", "", dep),
                 ("DieuKhienKiemKe", "LichSuThaoTac", "", dep),
                 ("PhieuKiemKe", "ChiTietKiemKe", "1..*", asc),
                 ("Kho", "PhieuKiemKe", "1..*", asc),
                 ("PhieuKiemKe", "TonKho", "layTon / ganTheoKiemKe", dep)], width=1300)
    class_slice("Class cắt lát - Báo cáo NXT (UC34)", "4.11_class_baocao.drawio",
                [("ManHinhBaoCao", 40, 60, 240, 4), ("DieuKhienBaoCao", 350, 60, 240, None),
                 ("BaoCaoNXT", 660, 40, 240, None), ("PhieuNhapKho", 960, 40, 220, 0),
                 ("PhieuXuatKho", 960, 300, 220, 0), ("TonKho", 660, 400, 220, 4)],
                [("ManHinhBaoCao", "DieuKhienBaoCao", "", dep),
                 ("DieuKhienBaoCao", "BaoCaoNXT", "", dep),
                 ("BaoCaoNXT", "PhieuNhapKho", "đọc", dep),
                 ("BaoCaoNXT", "PhieuXuatKho", "đọc", dep),
                 ("BaoCaoNXT", "TonKho", "đọc", dep)], width=1250)


# ───────────────────────── Chương 5: Sequence + State ─────────────────────────

def gen_sequences():
    for i in range(1, 35):
        ucid = f"UC{i:02d}"
        uc = USE_CASES[ucid]
        seq_diagram(f"Sequence {ucid} - {uc['ten']}", f"5.{i:02d}_sequence_uc{i:02d}.drawio",
                    uc["messages"])


def gen_states():
    def state(title, filename, done_label, ton_label):
        d = D(title, 1150, 560)
        start = d.node("", 40, 135, 24, 24, S_STATE_START)
        s1 = d.node("Mới tạo", 120, 120, 140, 56, S_STATE)
        s2 = d.node("Chờ duyệt", 340, 120, 140, 56, S_STATE)
        s3 = d.node("Đã duyệt", 560, 120, 140, 56, S_STATE)
        s4 = d.node(done_label, 780, 120, 160, 56, S_STATE)
        s5 = d.node("Từ chối", 340, 330, 140, 56, S_STATE)
        s6 = d.node("Đã hủy", 120, 330, 140, 56, S_STATE)
        end = d.node("", 1000, 136, 26, 26, S_STATE_END)
        d.edge(start, s1)
        d.edge(s1, s2, "guiDuyet()")
        d.edge(s2, s3, "duyetPhieu()")
        d.edge(s3, s4, ton_label)
        d.edge(s2, s5, "tuChoiPhieu(lyDo)")
        d.edge(s1, s6, "xoaPhieu(ma)")
        d.edge(s5, s6, "xoaPhieu(ma)")
        d.edge(s4, end)
        d.save(filename)

    state("State - Phiếu nhập kho", "5.35_state_phieunhap.drawio",
          "Đã nhập kho", "tangTon(...)")
    state("State - Phiếu xuất kho", "5.36_state_phieuxuat.drawio",
          "Đã xuất kho", "giamTon(...)")


# ───────────────────────── Chương 6 + 7 ─────────────────────────

def gen_component():
    d = D("Component - Kiến trúc phân lớp", 980, 800)
    layers = [
        (300, 40, "Giao diện người dùng\n(UI / Web Browser)"),
        (300, 150, "Controller Layer\n(API / Routing)"),
        (300, 260, "Service Layer\n(Business Logic)"),
        (300, 370, "Repository Layer\n(Data Access)"),
    ]
    ids = [d.node(t, x, y, 250, 64, S_COMPONENT) for x, y, t in layers]
    db = d.node("Cơ sở dữ liệu\n(PostgreSQL)", 355, 500, 150, 90, S_CYL)
    labs = ["HTTP/HTTPS", "Gọi dịch vụ", "Truy vấn"]
    for i in range(len(ids) - 1):
        d.edge(ids[i], ids[i + 1], labs[i])
    d.edge(ids[-1], db, "SQL")
    auth = d.node("Authentication\nService", 660, 150, 170, 60, S_COMPONENT)
    inv = d.node("Inventory Service\n(tính tồn kho)", 660, 280, 170, 60, S_COMPONENT)
    ntf = d.node("Notification Service\n(cảnh báo tồn thấp)", 660, 390, 180, 60, S_COMPONENT)
    d.edge(ids[1], auth, "xác thực")
    d.edge(ids[2], inv, "tồn kho")
    d.edge(ids[2], ntf, "cảnh báo")
    d.save("6.1_component.drawio")


def gen_ui_nav():
    d = D("Sơ đồ điều hướng màn hình", 1150, 650)
    login = d.node("Đăng nhập\n(UC01)", 480, 30, 130, 52, S_ROUNDED)
    home = d.node("Trang chủ\n(Dashboard)", 480, 130, 130, 52, S_RECT)
    d.edge(login, home)
    menus = [
        (40, 260, "Tài khoản\n(UC03-UC07)"),
        (190, 260, "Nhà cung cấp\n(UC08-UC12)"),
        (340, 260, "Hàng hóa\n(UC13-UC17)"),
        (490, 260, "Kho\n(UC18-UC22)"),
        (640, 260, "Báo cáo NXT\n(UC34)"),
        (190, 420, "Phiếu nhập\n(UC23-UC27)"),
        (400, 420, "Phiếu xuất\n(UC28-UC32)"),
        (610, 420, "Kiểm kê\n(UC33)"),
    ]
    for x, y, t in menus:
        m = d.node(t, x, y, 135, 56, S_RECT)
        d.edge(home, m)
    d.node("Menu hiển thị theo vai trò: QTV chỉ thấy Tài khoản; NV kho không thấy Kho và Báo cáo",
           250, 560, 620, 30, S_NOTE)
    d.save("8.1_ui_navigation.drawio")


# ───────────────────────── Manifest ─────────────────────────

FILES: dict[str, str] = {
    "2.1_activity_nhapkho": "Biểu đồ hoạt động nghiệp vụ nhập kho",
    "2.2_activity_xuatkho": "Biểu đồ hoạt động nghiệp vụ xuất kho",
    "2.3_activity_kiemke": "Biểu đồ hoạt động nghiệp vụ kiểm kê",
    "3.1_usecase_tongquat": "Biểu đồ use case tổng quát",
    "3.2_usecase_taikhoan": "Use case phân rã nhóm truy cập và tài khoản",
    "3.3_usecase_ncc": "Use case phân rã nhóm đối tác",
    "3.4_usecase_hang_kho": "Use case phân rã nhóm hàng hóa và kho",
    "3.5_usecase_phieu_baocao": "Use case phân rã nhóm nhập, xuất, kiểm kê, báo cáo",
    "4.1_class_master": "Biểu đồ lớp nhóm danh mục (Master Data)",
    "4.2_class_nhapxuat": "Biểu đồ lớp nhóm nghiệp vụ nhập - xuất",
    "4.3_class_kiemke_lichsu": "Biểu đồ lớp nhóm kiểm kê và lịch sử",
    "4.4_class_taikhoan": "Biểu đồ lớp cắt lát nhóm tài khoản / xác thực",
    "4.5_class_ncc": "Biểu đồ lớp cắt lát nhóm nhà cung cấp",
    "4.6_class_hanghoa": "Biểu đồ lớp cắt lát nhóm hàng hóa",
    "4.7_class_kho": "Biểu đồ lớp cắt lát nhóm kho",
    "4.8_class_phieunhap": "Biểu đồ lớp cắt lát nhóm phiếu nhập",
    "4.9_class_phieuxuat": "Biểu đồ lớp cắt lát nhóm phiếu xuất",
    "4.10_class_kiemke": "Biểu đồ lớp cắt lát nhóm kiểm kê",
    "4.11_class_baocao": "Biểu đồ lớp cắt lát nhóm báo cáo NXT",
    "6.1_component": "Component diagram kiến trúc phân lớp",
    "8.1_ui_navigation": "Sơ đồ điều hướng màn hình",
}
for i in range(1, 35):
    FILES[f"5.{i:02d}_sequence_uc{i:02d}"] = f"Biểu đồ trình tự UC{i:02d} - {USE_CASES[f'UC{i:02d}']['ten']}"
FILES["5.35_state_phieunhap"] = "Biểu đồ trạng thái phiếu nhập kho"
FILES["5.36_state_phieuxuat"] = "Biểu đồ trạng thái phiếu xuất kho"


def main():
    gen_activities()
    gen_usecases()
    gen_classes_overview()
    gen_class_slices()
    gen_sequences()
    gen_states()
    gen_component()
    gen_ui_nav()
    files = sorted(OUT.glob("*.drawio"))
    missing = [f for f in FILES if not (OUT / f"{f}.drawio").exists()]
    print(f"\nTotal: {len(files)} drawio files (manifest {len(FILES)})")
    if _SKIPPED:
        print(f"Giữ nguyên {len(_SKIPPED)} file đã tồn tại (không ghi đè). "
              f"Muốn sinh lại: --force [--only <tên>] - bản cũ sẽ được sao lưu vào _backup/.")
    if FORCE and _BACKUP_DIR.exists():
        print(f"Đã sao lưu bản trước khi ghi đè vào: {_BACKUP_DIR}")
    if missing:
        print("MISSING:", missing)


if __name__ == "__main__":
    main()
