package db

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"warehouse-backend/internal/models"
)

func Connect(dsn string) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := gdb.AutoMigrate(
		&models.NguoiDung{}, &models.NhomHangHoa{}, &models.DonViTinh{}, &models.HangHoa{},
		&models.NhaCungCap{}, &models.Kho{}, &models.TonKho{},
		&models.PhieuNhapKho{}, &models.ChiTietPhieuNhap{},
		&models.PhieuXuatKho{}, &models.ChiTietPhieuXuat{},
		&models.PhieuKiemKe{}, &models.ChiTietKiemKe{}, &models.LichSuThaoTac{},
	); err != nil {
		return nil, err
	}
	return gdb, nil
}

func hash(pw string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h)
}

// Seed tạo dữ liệu mẫu tối thiểu nếu CSDL trống: 3 tài khoản 3 vai trò,
// nhóm hàng / ĐVT / hàng hóa, NCC, 2 kho và tồn kho khởi tạo = 0.
func Seed(gdb *gorm.DB) {
	var n int64
	gdb.Model(&models.NguoiDung{}).Count(&n)
	if n > 0 {
		return
	}
	log.Println("seeding sample data...")
	users := []models.NguoiDung{
		{MaNguoiDung: "ND01", TenDangNhap: "admin", MatKhau: hash("admin123"), HoTen: "Quản trị viên", VaiTro: models.RoleQTV, TrangThai: true},
		{MaNguoiDung: "ND02", TenDangNhap: "quanly", MatKhau: hash("quanly123"), HoTen: "Quản lý kho", VaiTro: models.RoleQLKho, TrangThai: true},
		{MaNguoiDung: "ND03", TenDangNhap: "nhanvien", MatKhau: hash("nhanvien123"), HoTen: "Nhân viên kho", VaiTro: models.RoleNVKho, TrangThai: true},
	}
	gdb.Create(&users)

	groups := []models.NhomHangHoa{
		{MaNhomHang: "NH01", TenNhomHang: "Nguyên vật liệu"},
		{MaNhomHang: "NH02", TenNhomHang: "Bán thành phẩm"},
		{MaNhomHang: "NH03", TenNhomHang: "Thành phẩm"},
	}
	gdb.Create(&groups)

	units := []models.DonViTinh{
		{MaDonViTinh: "DVT01", TenDonViTinh: "Cái"},
		{MaDonViTinh: "DVT02", TenDonViTinh: "Kg"},
		{MaDonViTinh: "DVT03", TenDonViTinh: "Thùng"},
	}
	gdb.Create(&units)

	products := []models.HangHoa{
		{MaHangHoa: "HH01", TenHangHoa: "Thép cuộn D10", MaNhomHang: "NH01", MaDonViTinh: "DVT02", GiaNhap: 15000, GiaXuat: 18000, TrangThai: true},
		{MaHangHoa: "HH02", TenHangHoa: "Bao bì carton A3", MaNhomHang: "NH01", MaDonViTinh: "DVT01", GiaNhap: 3000, GiaXuat: 4500, TrangThai: true},
		{MaHangHoa: "HH03", TenHangHoa: "Sản phẩm hoàn thiện X1", MaNhomHang: "NH03", MaDonViTinh: "DVT03", GiaNhap: 220000, GiaXuat: 275000, TrangThai: true},
	}
	gdb.Create(&products)

	suppliers := []models.NhaCungCap{
		{MaNCC: "NCC01", TenNCC: "Công ty Thép Miền Bắc", DiaChi: "Hà Nội", SDT: "0241111111", NguoiLienHe: "Nguyễn Văn A", TrangThai: true},
		{MaNCC: "NCC02", TenNCC: "Công ty Bao bì Sông Hồng", DiaChi: "Hưng Yên", SDT: "0242222222", NguoiLienHe: "Trần Thị B", TrangThai: true},
	}
	gdb.Create(&suppliers)

	warehouses := []models.Kho{
		{MaKho: "K01", TenKho: "Kho trung tâm Hà Nội", DiaChi: "KCN Thăng Long", NguoiPhuTrach: "Quản lý kho", TrangThai: true},
		{MaKho: "K02", TenKho: "Kho chi nhánh Bắc Ninh", DiaChi: "KCN Quế Võ", NguoiPhuTrach: "Quản lý kho", TrangThai: true},
	}
	gdb.Create(&warehouses)

	now := time.Now()
	for _, k := range warehouses {
		for _, p := range products {
			gdb.Create(&models.TonKho{MaKho: k.MaKho, MaHangHoa: p.MaHangHoa, SoLuongTon: 0, NgayCapNhat: now})
		}
	}
	log.Println("seed done: admin/admin123, quanly/quanly123, nhanvien/nhanvien123")
}
