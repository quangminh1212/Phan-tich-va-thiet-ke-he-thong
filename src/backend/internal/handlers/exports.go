package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

// UC28-UC32: phiếu xuất kho. Không gắn khách hàng - chỉ lý do xuất / ghi chú.
// Điểm khác phiếu nhập: kiểm tra tồn đủ (kiemTraDu) khi lập / sửa và khi duyệt thì GIẢM tồn.

type exportReq struct {
	MaKho    string       `json:"ma_kho" binding:"required"`
	NgayXuat string       `json:"ngay_xuat"`
	LyDoXuat string       `json:"ly_do_xuat" binding:"required"`
	GhiChu   string       `json:"ghi_chu"`
	ChiTiet  []ticketLine `json:"chi_tiet" binding:"required"`
}

// checkStockLines kiểm tra tồn đủ cho toàn bộ dòng hàng (UC28 bước kiemTraDu).
func checkStockLines(tx *gorm.DB, maKho string, lines []ticketLine) error {
	for _, l := range lines {
		if l.SoLuong <= 0 {
			return errors.New("số lượng dòng hàng phải > 0")
		}
		ok, ton, err := services.KiemTraDu(tx, maKho, l.MaHangHoa, l.SoLuong)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("tồn kho không đủ: %s tại %s còn %d, cần xuất %d", l.MaHangHoa, maKho, ton, l.SoLuong)
		}
	}
	return nil
}

// ListExports - UC32 xem danh sách + UC31 tìm kiếm.
func (h *H) ListExports(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.PhieuXuatKho{})
	if kw := c.Query("q"); kw != "" {
		q = q.Where("ma_phieu_xuat ILIKE ?", "%"+kw+"%")
	}
	if v := c.Query("trang_thai"); v != "" {
		q = q.Where("trang_thai = ?", v)
	}
	if v := c.Query("ma_kho"); v != "" {
		q = q.Where("ma_kho = ?", v)
	}
	var total int64
	q.Count(&total)
	var items []models.PhieuXuatKho
	if err := q.Order("ma_phieu_xuat DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	listResponse(c, items, total, page)
}

func (h *H) GetExport(c *gin.Context) {
	var p models.PhieuXuatKho
	if err := h.DB.Preload("ChiTiet").First(&p, "ma_phieu_xuat = ?", c.Param("id")).Error; err != nil {
		notFound(c)
		return
	}
	c.JSON(http.StatusOK, p)
}

// CreateExport - UC28 lập phiếu xuất: kiểm tra tồn đủ từng dòng rồi mới lưu, trạng thái "Chờ duyệt".
func (h *H) CreateExport(c *gin.Context) {
	var req exportReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.ChiTiet) == 0 {
		badRequest(c, "phiếu cần kho, lý do xuất và ít nhất một dòng hàng")
		return
	}
	var p models.PhieuXuatKho
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := checkStockLines(tx, req.MaKho, req.ChiTiet); err != nil {
			return err
		}
		code, err := nextCode(tx, "phieu_xuat_kho", "ma_phieu_xuat", "PX")
		if err != nil {
			return err
		}
		p = models.PhieuXuatKho{
			MaPhieuXuat: code, MaKho: req.MaKho, MaNguoiDung: c.GetString("user_id"),
			NgayXuat: parseDate(req.NgayXuat), LyDoXuat: req.LyDoXuat, GhiChu: req.GhiChu,
			TrangThai: models.StatusChoDuyet,
		}
		for _, l := range req.ChiTiet {
			ct := models.ChiTietPhieuXuat{
				MaPhieuXuat: code, MaHangHoa: l.MaHangHoa,
				SoLuong: l.SoLuong, DonGia: l.DonGia, ThanhTien: float64(l.SoLuong) * l.DonGia,
			}
			p.ChiTiet = append(p.ChiTiet, ct)
			p.TongTien += ct.ThanhTien
		}
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "lap_phieu_xuat", "PhieuXuatKho", code)
		return nil
	})
	if err != nil {
		badRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, p)
}

// UpdateExport - UC29 sửa phiếu xuất: chỉ khi "Chờ duyệt"; đổi số lượng thì kiểm tra tồn lại.
func (h *H) UpdateExport(c *gin.Context) {
	id := c.Param("id")
	var req exportReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.ChiTiet) == 0 {
		badRequest(c, "phiếu cần kho, lý do xuất và ít nhất một dòng hàng")
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuXuatKho
		if err := tx.First(&p, "ma_phieu_xuat = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		if err := checkStockLines(tx, req.MaKho, req.ChiTiet); err != nil {
			return err
		}
		if err := tx.Where("ma_phieu_xuat = ?", id).Delete(&models.ChiTietPhieuXuat{}).Error; err != nil {
			return err
		}
		p.MaKho, p.NgayXuat, p.LyDoXuat, p.GhiChu, p.TongTien =
			req.MaKho, parseDate(req.NgayXuat), req.LyDoXuat, req.GhiChu, 0
		for _, l := range req.ChiTiet {
			ct := models.ChiTietPhieuXuat{
				MaPhieuXuat: id, MaHangHoa: l.MaHangHoa,
				SoLuong: l.SoLuong, DonGia: l.DonGia, ThanhTien: float64(l.SoLuong) * l.DonGia,
			}
			if err := tx.Create(&ct).Error; err != nil {
				return err
			}
			p.TongTien += ct.ThanhTien
		}
		if err := tx.Save(&p).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "sua_phieu_xuat", "PhieuXuatKho", id)
		return nil
	})
	h.ticketErr(c, err, id, "phiếu chỉ sửa được khi Chờ duyệt")
}

// DeleteExport - UC30 xóa phiếu xuất: chỉ khi chưa duyệt (không cần hoàn tồn vì chưa trừ).
func (h *H) DeleteExport(c *gin.Context) {
	id := c.Param("id")
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuXuatKho
		if err := tx.First(&p, "ma_phieu_xuat = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai == models.StatusDaDuyet {
			return errBadState
		}
		if err := tx.Where("ma_phieu_xuat = ?", id).Delete(&models.ChiTietPhieuXuat{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&p).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "xoa_phieu_xuat", "PhieuXuatKho", id)
		return nil
	})
	h.ticketErr(c, err, id, "phiếu Đã duyệt không được xóa")
}

// ApproveExport - QL kho duyệt: kiểm tra tồn lần cuối rồi GIẢM tồn từng dòng.
func (h *H) ApproveExport(c *gin.Context) {
	id := c.Param("id")
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuXuatKho
		if err := tx.Preload("ChiTiet").First(&p, "ma_phieu_xuat = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		for _, ct := range p.ChiTiet {
			if err := services.GiamTon(tx, p.MaKho, ct.MaHangHoa, ct.SoLuong); err != nil {
				return err
			}
		}
		if err := tx.Model(&p).Update("trang_thai", models.StatusDaDuyet).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "duyet_phieu_xuat", "PhieuXuatKho", id)
		return nil
	})
	h.ticketErr(c, err, id, "chỉ duyệt được phiếu Chờ duyệt")
}

// RejectExport - từ chối phiếu xuất: tồn không đổi.
func (h *H) RejectExport(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		LyDo string `json:"ly_do"`
	}
	_ = c.ShouldBindJSON(&body)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuXuatKho
		if err := tx.First(&p, "ma_phieu_xuat = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		if err := tx.Model(&p).Update("trang_thai", models.StatusTuChoi).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "tu_choi_phieu_xuat", "PhieuXuatKho", id+" - "+body.LyDo)
		return nil
	})
	h.ticketErr(c, err, id, "chỉ từ chối được phiếu Chờ duyệt")
}
