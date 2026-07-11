// Package models: 14 bảng CSDL theo Chương 7 của báo cáo PTTKHT.
package models

import "time"

// Vai trò người dùng (Bảng 3.1)
const (
	RoleQTV   = "QTV"   // Quản trị viên
	RoleQLKho = "QLKho" // Quản lý kho
	RoleNVKho = "NVKho" // Nhân viên kho
)

// Trạng thái phiếu (biểu đồ trạng thái Hình 5.35 / 5.36)
const (
	StatusChoDuyet = "cho_duyet"
	StatusDaDuyet  = "da_duyet"
	StatusTuChoi   = "tu_choi"
)

type NguoiDung struct {
	MaNguoiDung string `gorm:"primaryKey;size:20" json:"ma_nguoi_dung"`
	TenDangNhap string `gorm:"size:50;uniqueIndex;not null" json:"ten_dang_nhap"`
	MatKhau     string `gorm:"size:255;not null" json:"-"`
	HoTen       string `gorm:"size:100" json:"ho_ten"`
	Email       string `gorm:"size:100" json:"email"`
	VaiTro      string `gorm:"size:20;not null" json:"vai_tro"`
	TrangThai   bool   `gorm:"default:true" json:"trang_thai"`
}

func (NguoiDung) TableName() string { return "nguoi_dung" }

type NhomHangHoa struct {
	MaNhomHang  string `gorm:"primaryKey;size:20" json:"ma_nhom_hang"`
	TenNhomHang string `gorm:"size:100;not null" json:"ten_nhom_hang"`
	MoTa        string `gorm:"size:255" json:"mo_ta"`
}

func (NhomHangHoa) TableName() string { return "nhom_hang_hoa" }

type DonViTinh struct {
	MaDonViTinh  string `gorm:"primaryKey;size:20" json:"ma_don_vi_tinh"`
	TenDonViTinh string `gorm:"size:50;not null" json:"ten_don_vi_tinh"`
}

func (DonViTinh) TableName() string { return "don_vi_tinh" }

type HangHoa struct {
	MaHangHoa   string  `gorm:"primaryKey;size:20" json:"ma_hang_hoa"`
	TenHangHoa  string  `gorm:"size:255;not null" json:"ten_hang_hoa"`
	MaNhomHang  string  `gorm:"size:20" json:"ma_nhom_hang"`
	MaDonViTinh string  `gorm:"size:20" json:"ma_don_vi_tinh"`
	GiaNhap     float64 `gorm:"type:decimal(18,2);default:0" json:"gia_nhap"`
	GiaXuat     float64 `gorm:"type:decimal(18,2);default:0" json:"gia_xuat"`
	MoTa        string  `gorm:"size:500" json:"mo_ta"`
	TrangThai   bool    `gorm:"default:true" json:"trang_thai"`
}

func (HangHoa) TableName() string { return "hang_hoa" }

type NhaCungCap struct {
	MaNCC       string `gorm:"primaryKey;size:20;column:ma_ncc" json:"ma_ncc"`
	TenNCC      string `gorm:"size:255;not null;column:ten_ncc" json:"ten_ncc"`
	DiaChi      string `gorm:"size:255" json:"dia_chi"`
	SDT         string `gorm:"size:20;column:sdt" json:"sdt"`
	Email       string `gorm:"size:100" json:"email"`
	NguoiLienHe string `gorm:"size:100" json:"nguoi_lien_he"`
	TrangThai   bool   `gorm:"default:true" json:"trang_thai"`
}

func (NhaCungCap) TableName() string { return "nha_cung_cap" }

type Kho struct {
	MaKho         string `gorm:"primaryKey;size:20" json:"ma_kho"`
	TenKho        string `gorm:"size:255;not null" json:"ten_kho"`
	DiaChi        string `gorm:"size:255" json:"dia_chi"`
	NguoiPhuTrach string `gorm:"size:100" json:"nguoi_phu_trach"`
	TrangThai     bool   `gorm:"default:true" json:"trang_thai"`
}

func (Kho) TableName() string { return "kho" }

type TonKho struct {
	MaKho       string    `gorm:"primaryKey;size:20" json:"ma_kho"`
	MaHangHoa   string    `gorm:"primaryKey;size:20" json:"ma_hang_hoa"`
	SoLuongTon  int       `gorm:"default:0" json:"so_luong_ton"`
	NgayCapNhat time.Time `json:"ngay_cap_nhat"`
}

func (TonKho) TableName() string { return "ton_kho" }

type PhieuNhapKho struct {
	MaPhieuNhap string             `gorm:"primaryKey;size:20" json:"ma_phieu_nhap"`
	MaNCC       string             `gorm:"size:20;column:ma_ncc" json:"ma_ncc"`
	MaKho       string             `gorm:"size:20;not null" json:"ma_kho"`
	MaNguoiDung string             `gorm:"size:20;not null" json:"ma_nguoi_dung"`
	NgayNhap    time.Time          `json:"ngay_nhap"`
	TongTien    float64            `gorm:"type:decimal(18,2);default:0" json:"tong_tien"`
	TrangThai   string             `gorm:"size:20;default:cho_duyet" json:"trang_thai"`
	ChiTiet     []ChiTietPhieuNhap `gorm:"foreignKey:MaPhieuNhap;constraint:OnDelete:CASCADE" json:"chi_tiet"`
}

func (PhieuNhapKho) TableName() string { return "phieu_nhap_kho" }

type ChiTietPhieuNhap struct {
	MaPhieuNhap string  `gorm:"primaryKey;size:20" json:"ma_phieu_nhap"`
	MaHangHoa   string  `gorm:"primaryKey;size:20" json:"ma_hang_hoa"`
	SoLuong     int     `gorm:"not null" json:"so_luong"`
	DonGia      float64 `gorm:"type:decimal(18,2);default:0" json:"don_gia"`
	ThanhTien   float64 `gorm:"type:decimal(18,2);default:0" json:"thanh_tien"`
}

func (ChiTietPhieuNhap) TableName() string { return "chi_tiet_phieu_nhap" }

type PhieuXuatKho struct {
	MaPhieuXuat string             `gorm:"primaryKey;size:20" json:"ma_phieu_xuat"`
	MaKho       string             `gorm:"size:20;not null" json:"ma_kho"`
	MaNguoiDung string             `gorm:"size:20;not null" json:"ma_nguoi_dung"`
	NgayXuat    time.Time          `json:"ngay_xuat"`
	LyDoXuat    string             `gorm:"size:255" json:"ly_do_xuat"`
	GhiChu      string             `gorm:"size:500" json:"ghi_chu"`
	TongTien    float64            `gorm:"type:decimal(18,2);default:0" json:"tong_tien"`
	TrangThai   string             `gorm:"size:20;default:cho_duyet" json:"trang_thai"`
	ChiTiet     []ChiTietPhieuXuat `gorm:"foreignKey:MaPhieuXuat;constraint:OnDelete:CASCADE" json:"chi_tiet"`
}

func (PhieuXuatKho) TableName() string { return "phieu_xuat_kho" }

type ChiTietPhieuXuat struct {
	MaPhieuXuat string  `gorm:"primaryKey;size:20" json:"ma_phieu_xuat"`
	MaHangHoa   string  `gorm:"primaryKey;size:20" json:"ma_hang_hoa"`
	SoLuong     int     `gorm:"not null" json:"so_luong"`
	DonGia      float64 `gorm:"type:decimal(18,2);default:0" json:"don_gia"`
	ThanhTien   float64 `gorm:"type:decimal(18,2);default:0" json:"thanh_tien"`
}

func (ChiTietPhieuXuat) TableName() string { return "chi_tiet_phieu_xuat" }

type PhieuKiemKe struct {
	MaPhieuKiemKe string           `gorm:"primaryKey;size:20" json:"ma_phieu_kiem_ke"`
	MaKho         string           `gorm:"size:20;not null" json:"ma_kho"`
	MaNguoiDung   string           `gorm:"size:20;not null" json:"ma_nguoi_dung"`
	NgayKiem      time.Time        `json:"ngay_kiem"`
	TrangThai     string           `gorm:"size:20;default:cho_duyet" json:"trang_thai"`
	ChiTiet       []ChiTietKiemKe  `gorm:"foreignKey:MaPhieuKiemKe;constraint:OnDelete:CASCADE" json:"chi_tiet"`
}

func (PhieuKiemKe) TableName() string { return "phieu_kiem_ke" }

type ChiTietKiemKe struct {
	MaPhieuKiemKe string `gorm:"primaryKey;size:20" json:"ma_phieu_kiem_ke"`
	MaHangHoa     string `gorm:"primaryKey;size:20" json:"ma_hang_hoa"`
	SLThucTe      int    `gorm:"column:sl_thuc_te;default:0" json:"sl_thuc_te"`
	SLSoSach      int    `gorm:"column:sl_so_sach;default:0" json:"sl_so_sach"`
	ChenhLech     int    `gorm:"default:0" json:"chenh_lech"`
}

func (ChiTietKiemKe) TableName() string { return "chi_tiet_kiem_ke" }

type LichSuThaoTac struct {
	MaLog       uint      `gorm:"primaryKey;autoIncrement" json:"ma_log"`
	MaNguoiDung string    `gorm:"size:20" json:"ma_nguoi_dung"`
	HanhDong    string    `gorm:"size:100" json:"hanh_dong"`
	ThoiGian    time.Time `json:"thoi_gian"`
	DoiTuong    string    `gorm:"size:100" json:"doi_tuong"`
	MoTa        string    `gorm:"size:500" json:"mo_ta"`
}

func (LichSuThaoTac) TableName() string { return "lich_su_thao_tac" }
