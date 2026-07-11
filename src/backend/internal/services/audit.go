package services

import (
	"time"

	"gorm.io/gorm"

	"warehouse-backend/internal/models"
)

// GhiLog ghi nhận thao tác vào lịch sử (LichSuThaoTac.ghiLog - FR-10).
func GhiLog(tx *gorm.DB, maNguoiDung, hanhDong, doiTuong, moTa string) {
	tx.Create(&models.LichSuThaoTac{
		MaNguoiDung: maNguoiDung,
		HanhDong:    hanhDong,
		ThoiGian:    time.Now(),
		DoiTuong:    doiTuong,
		MoTa:        moTa,
	})
}
