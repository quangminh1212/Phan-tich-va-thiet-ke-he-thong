package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"warehouse-backend/internal/models"
)

// Các thao tác trên TonKho - tương ứng phương thức lớp TonKho ở Chương 4:
// khoiTaoTon / tangTon / giamTon / kiemTraDu / ganTheoKiemKe.
// Tồn kho CHỈ được thay đổi qua các hàm này, trong transaction duyệt phiếu.

func KhoiTaoTon(tx *gorm.DB, maKho, maHang string) error {
	var cnt int64
	tx.Model(&models.TonKho{}).Where("ma_kho = ? AND ma_hang_hoa = ?", maKho, maHang).Count(&cnt)
	if cnt > 0 {
		return nil
	}
	return tx.Create(&models.TonKho{MaKho: maKho, MaHangHoa: maHang, SoLuongTon: 0, NgayCapNhat: time.Now()}).Error
}

func TangTon(tx *gorm.DB, maKho, maHang string, soLuong int) error {
	if err := KhoiTaoTon(tx, maKho, maHang); err != nil {
		return err
	}
	return tx.Model(&models.TonKho{}).
		Where("ma_kho = ? AND ma_hang_hoa = ?", maKho, maHang).
		Updates(map[string]any{
			"so_luong_ton":  gorm.Expr("so_luong_ton + ?", soLuong),
			"ngay_cap_nhat": time.Now(),
		}).Error
}

func GiamTon(tx *gorm.DB, maKho, maHang string, soLuong int) error {
	ok, ton, err := KiemTraDu(tx, maKho, maHang, soLuong)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("tồn kho không đủ: %s tại %s còn %d, cần xuất %d", maHang, maKho, ton, soLuong)
	}
	return tx.Model(&models.TonKho{}).
		Where("ma_kho = ? AND ma_hang_hoa = ?", maKho, maHang).
		Updates(map[string]any{
			"so_luong_ton":  gorm.Expr("so_luong_ton - ?", soLuong),
			"ngay_cap_nhat": time.Now(),
		}).Error
}

// KiemTraDu trả về (đủ hay không, số tồn hiện tại).
func KiemTraDu(tx *gorm.DB, maKho, maHang string, soLuong int) (bool, int, error) {
	var t models.TonKho
	err := tx.Where("ma_kho = ? AND ma_hang_hoa = ?", maKho, maHang).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return t.SoLuongTon >= soLuong, t.SoLuongTon, nil
}

func GanTheoKiemKe(tx *gorm.DB, maKho, maHang string, slThucTe int) error {
	if err := KhoiTaoTon(tx, maKho, maHang); err != nil {
		return err
	}
	return tx.Model(&models.TonKho{}).
		Where("ma_kho = ? AND ma_hang_hoa = ?", maKho, maHang).
		Updates(map[string]any{"so_luong_ton": slThucTe, "ngay_cap_nhat": time.Now()}).Error
}
