package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

// SeedDemo nạp dữ liệu demo khối lượng lớn: 5 kho, 100 loại hàng hóa,
// 50 phiếu nhập, 50 phiếu xuất. Idempotent: bỏ qua nếu đã đủ 100 hàng hóa.
//
// Dữ liệu sinh ra tôn trọng bất biến nghiệp vụ của báo cáo:
// tồn kho chỉ thay đổi khi phiếu ĐÃ DUYỆT (TangTon / GiamTon),
// phiếu xuất không bao giờ vượt tồn tại thời điểm lập.
func SeedDemo(gdb *gorm.DB) {
	var nProducts int64
	gdb.Model(&models.HangHoa{}).Count(&nProducts)
	if nProducts >= 100 {
		return
	}
	log.Println("seeding demo data (5 kho, 100 hàng hóa, 50 phiếu nhập, 50 phiếu xuất)...")
	rng := rand.New(rand.NewSource(42)) // cố định để tái lập được

	// ---- 5 kho ----
	khoNames := []string{"Kho trung tâm Hà Nội", "Kho chi nhánh Bắc Ninh", "Kho Hải Phòng", "Kho Đà Nẵng", "Kho TP.HCM"}
	khoAddrs := []string{"KCN Thăng Long", "KCN Quế Võ", "KCN Đình Vũ", "KCN Hòa Khánh", "KCN Tân Bình"}
	for i := 0; i < 5; i++ {
		code := fmt.Sprintf("K%02d", i+1)
		var cnt int64
		gdb.Model(&models.Kho{}).Where("ma_kho = ?", code).Count(&cnt)
		if cnt == 0 {
			gdb.Create(&models.Kho{MaKho: code, TenKho: khoNames[i], DiaChi: khoAddrs[i],
				NguoiPhuTrach: "Quản lý kho", TrangThai: true})
		}
	}

	// ---- 100 loại hàng hóa ----
	groups := []string{"NH01", "NH02", "NH03"}
	units := []string{"DVT01", "DVT02", "DVT03"}
	tenMau := []string{"Thép cuộn", "Thép tấm", "Nhôm định hình", "Bao bì carton", "Màng PE",
		"Ốc vít M6", "Bulong M10", "Sơn tĩnh điện", "Dầu thủy lực", "Vòng bi",
		"Dây điện 2x1.5", "Ống nhựa PVC", "Gioăng cao su", "Keo dán công nghiệp", "Pallet gỗ",
		"Thùng nhựa", "Băng dính", "Găng tay bảo hộ", "Linh kiện X", "Thành phẩm Y"}
	for i := 1; i <= 200 && nProducts < 100; i++ {
		code := fmt.Sprintf("HH%03d", i)
		var cnt int64
		gdb.Model(&models.HangHoa{}).Where("ma_hang_hoa = ?", code).Count(&cnt)
		if cnt > 0 {
			continue
		}
		nProducts++
		giaNhap := float64((rng.Intn(490) + 10) * 1000)
		gdb.Create(&models.HangHoa{
			MaHangHoa:   code,
			TenHangHoa:  fmt.Sprintf("%s loại %d", tenMau[(i-1)%len(tenMau)], (i-1)/len(tenMau)+1),
			MaNhomHang:  groups[i%len(groups)],
			MaDonViTinh: units[i%len(units)],
			GiaNhap:     giaNhap,
			GiaXuat:     giaNhap * 1.25,
			TrangThai:   true,
		})
	}

	// ---- khởi tạo tồn = 0 cho mọi cặp (kho, hàng) ----
	var khoList []models.Kho
	gdb.Find(&khoList)
	var products []models.HangHoa
	gdb.Order("ma_hang_hoa").Find(&products)
	gdb.Transaction(func(tx *gorm.DB) error {
		for _, k := range khoList {
			for _, p := range products {
				if err := services.KhoiTaoTon(tx, k.MaKho, p.MaHangHoa); err != nil {
					return err
				}
			}
		}
		return nil
	})

	var suppliers []models.NhaCungCap
	gdb.Where("trang_thai = true").Find(&suppliers)
	if len(suppliers) == 0 {
		gdb.Create(&models.NhaCungCap{MaNCC: "NCC01", TenNCC: "NCC mặc định", TrangThai: true})
		gdb.Where("trang_thai = true").Find(&suppliers)
	}

	pick := func(n int) int { return rng.Intn(n) }
	day := func(daysAgo int) time.Time { return time.Now().AddDate(0, 0, -daysAgo) }
	// trạng thái: 40 đã duyệt / 5 chờ duyệt / 5 từ chối cho mỗi loại phiếu
	statusFor := func(i int) string {
		switch {
		case i < 40:
			return models.StatusDaDuyet
		case i < 45:
			return models.StatusChoDuyet
		default:
			return models.StatusTuChoi
		}
	}

	// ---- 50 phiếu nhập ----
	var nPN int64
	gdb.Model(&models.PhieuNhapKho{}).Count(&nPN)
	for i := 0; int64(i) < 50-nPN; i++ {
		st := statusFor(i)
		err := gdb.Transaction(func(tx *gorm.DB) error {
			code := fmt.Sprintf("PN%05d", int(nPN)+i+1)
			kho := khoList[pick(len(khoList))]
			p := models.PhieuNhapKho{
				MaPhieuNhap: code,
				MaNCC:       suppliers[pick(len(suppliers))].MaNCC,
				MaKho:       kho.MaKho,
				MaNguoiDung: "ND03",
				NgayNhap:    day(pick(30)),
				TrangThai:   st,
			}
			nLines := pick(3) + 1
			seen := map[string]bool{}
			for j := 0; j < nLines; j++ {
				hh := products[pick(len(products))]
				if seen[hh.MaHangHoa] {
					continue
				}
				seen[hh.MaHangHoa] = true
				sl := pick(90) + 10
				ct := models.ChiTietPhieuNhap{
					MaPhieuNhap: code, MaHangHoa: hh.MaHangHoa,
					SoLuong: sl, DonGia: hh.GiaNhap, ThanhTien: float64(sl) * hh.GiaNhap,
				}
				p.ChiTiet = append(p.ChiTiet, ct)
				p.TongTien += ct.ThanhTien
			}
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
			if st == models.StatusDaDuyet {
				for _, ct := range p.ChiTiet {
					if err := services.TangTon(tx, p.MaKho, ct.MaHangHoa, ct.SoLuong); err != nil {
						return err
					}
				}
				services.GhiLog(tx, "ND02", "duyet_phieu_nhap", "PhieuNhapKho", code)
			}
			services.GhiLog(tx, "ND03", "lap_phieu_nhap", "PhieuNhapKho", code)
			return nil
		})
		if err != nil {
			log.Printf("seed demo phiếu nhập lỗi: %v", err)
		}
	}

	// ---- 50 phiếu xuất (chỉ xuất trong phạm vi tồn hiện có) ----
	lyDo := []string{"Xuất bán lẻ", "Xuất chuyển nội bộ", "Xuất cho sản xuất", "Xuất trả nhà cung cấp", "Xuất hàng mẫu"}
	var nPX int64
	gdb.Model(&models.PhieuXuatKho{}).Count(&nPX)
	for i := 0; int64(i) < 50-nPX; i++ {
		st := statusFor(i)
		err := gdb.Transaction(func(tx *gorm.DB) error {
			code := fmt.Sprintf("PX%05d", int(nPX)+i+1)
			kho := khoList[pick(len(khoList))]
			// chọn các dòng tồn dương của kho này
			var tons []models.TonKho
			if err := tx.Where("ma_kho = ? AND so_luong_ton > 0", kho.MaKho).
				Order("ma_hang_hoa").Find(&tons).Error; err != nil {
				return err
			}
			if len(tons) == 0 {
				return nil // kho chưa có tồn - bỏ qua, vòng sau kho khác
			}
			p := models.PhieuXuatKho{
				MaPhieuXuat: code, MaKho: kho.MaKho, MaNguoiDung: "ND03",
				NgayXuat: day(pick(30)), LyDoXuat: lyDo[pick(len(lyDo))],
				TrangThai: st,
			}
			nLines := pick(3) + 1
			seen := map[string]bool{}
			for j := 0; j < nLines && j < len(tons); j++ {
				t := tons[pick(len(tons))]
				if seen[t.MaHangHoa] || t.SoLuongTon < 1 {
					continue
				}
				seen[t.MaHangHoa] = true
				sl := pick(t.SoLuongTon/2+1) + 1 // luôn <= tồn
				var hh models.HangHoa
				tx.First(&hh, "ma_hang_hoa = ?", t.MaHangHoa)
				ct := models.ChiTietPhieuXuat{
					MaPhieuXuat: code, MaHangHoa: t.MaHangHoa,
					SoLuong: sl, DonGia: hh.GiaXuat, ThanhTien: float64(sl) * hh.GiaXuat,
				}
				p.ChiTiet = append(p.ChiTiet, ct)
				p.TongTien += ct.ThanhTien
			}
			if len(p.ChiTiet) == 0 {
				return nil
			}
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
			if st == models.StatusDaDuyet {
				for _, ct := range p.ChiTiet {
					if err := services.GiamTon(tx, p.MaKho, ct.MaHangHoa, ct.SoLuong); err != nil {
						return err
					}
				}
				services.GhiLog(tx, "ND02", "duyet_phieu_xuat", "PhieuXuatKho", code)
			}
			services.GhiLog(tx, "ND03", "lap_phieu_xuat", "PhieuXuatKho", code)
			return nil
		})
		if err != nil {
			log.Printf("seed demo phiếu xuất lỗi: %v", err)
		}
	}

	var cP, cK, cPN, cPX int64
	gdb.Model(&models.HangHoa{}).Count(&cP)
	gdb.Model(&models.Kho{}).Count(&cK)
	gdb.Model(&models.PhieuNhapKho{}).Count(&cPN)
	gdb.Model(&models.PhieuXuatKho{}).Count(&cPX)
	log.Printf("seed demo xong: %d hàng hóa, %d kho, %d phiếu nhập, %d phiếu xuất", cP, cK, cPN, cPX)
}
