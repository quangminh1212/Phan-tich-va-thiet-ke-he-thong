package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

// UC23-UC27: phiếu nhập kho. Lập / sửa / xóa: NV kho + QL kho; duyệt / từ chối: QL kho.

type ticketLine struct {
	MaHangHoa string  `json:"ma_hang_hoa" binding:"required"`
	SoLuong   int     `json:"so_luong" binding:"required"`
	DonGia    float64 `json:"don_gia"`
}

type importReq struct {
	MaNCC   string       `json:"ma_ncc" binding:"required"`
	MaKho   string       `json:"ma_kho" binding:"required"`
	NgayNhap string      `json:"ngay_nhap"`
	ChiTiet []ticketLine `json:"chi_tiet" binding:"required"`
}

func parseDate(s string) time.Time {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t
	}
	return time.Now()
}

var errBadState = errors.New("trạng thái phiếu không cho phép thao tác")

// ListImports - UC27 xem danh sách + UC26 tìm kiếm (?q, ?trang_thai, ?ma_kho, ?ma_ncc).
func (h *H) ListImports(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.PhieuNhapKho{})
	if kw := c.Query("q"); kw != "" {
		q = q.Where("ma_phieu_nhap ILIKE ?", "%"+kw+"%")
	}
	if v := c.Query("trang_thai"); v != "" {
		q = q.Where("trang_thai = ?", v)
	}
	if v := c.Query("ma_kho"); v != "" {
		q = q.Where("ma_kho = ?", v)
	}
	if v := c.Query("ma_ncc"); v != "" {
		q = q.Where("ma_ncc = ?", v)
	}
	var total int64
	q.Count(&total)
	var items []models.PhieuNhapKho
	if err := q.Order("ma_phieu_nhap DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	listResponse(c, items, total, page)
}

func (h *H) GetImport(c *gin.Context) {
	var p models.PhieuNhapKho
	if err := h.DB.Preload("ChiTiet").First(&p, "ma_phieu_nhap = ?", c.Param("id")).Error; err != nil {
		notFound(c)
		return
	}
	c.JSON(http.StatusOK, p)
}

// CreateImport - UC23 lập phiếu nhập: tạo phiếu + chi tiết, tính tổng tiền,
// trạng thái "Chờ duyệt" (lapPhieu + themChiTiet + tinhTongTien + guiDuyet).
func (h *H) CreateImport(c *gin.Context) {
	var req importReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.ChiTiet) == 0 {
		badRequest(c, "phiếu cần NCC, kho và ít nhất một dòng hàng")
		return
	}
	var p models.PhieuNhapKho
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		code, err := nextCode(tx, "phieu_nhap_kho", "ma_phieu_nhap", "PN")
		if err != nil {
			return err
		}
		p = models.PhieuNhapKho{
			MaPhieuNhap: code, MaNCC: req.MaNCC, MaKho: req.MaKho,
			MaNguoiDung: c.GetString("user_id"), NgayNhap: parseDate(req.NgayNhap),
			TrangThai: models.StatusChoDuyet,
		}
		for _, l := range req.ChiTiet {
			if l.SoLuong <= 0 {
				return errors.New("số lượng dòng hàng phải > 0")
			}
			ct := models.ChiTietPhieuNhap{
				MaPhieuNhap: code, MaHangHoa: l.MaHangHoa,
				SoLuong: l.SoLuong, DonGia: l.DonGia, ThanhTien: float64(l.SoLuong) * l.DonGia,
			}
			p.ChiTiet = append(p.ChiTiet, ct)
			p.TongTien += ct.ThanhTien
		}
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "lap_phieu_nhap", "PhieuNhapKho", code)
		return nil
	})
	if err != nil {
		badRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, p)
}

// UpdateImport - UC24 sửa phiếu nhập: chỉ khi "Chờ duyệt"; thay toàn bộ header + chi tiết.
func (h *H) UpdateImport(c *gin.Context) {
	id := c.Param("id")
	var req importReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.ChiTiet) == 0 {
		badRequest(c, "phiếu cần NCC, kho và ít nhất một dòng hàng")
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuNhapKho
		if err := tx.First(&p, "ma_phieu_nhap = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		if err := tx.Where("ma_phieu_nhap = ?", id).Delete(&models.ChiTietPhieuNhap{}).Error; err != nil {
			return err
		}
		p.MaNCC, p.MaKho, p.NgayNhap, p.TongTien = req.MaNCC, req.MaKho, parseDate(req.NgayNhap), 0
		for _, l := range req.ChiTiet {
			if l.SoLuong <= 0 {
				return errors.New("số lượng dòng hàng phải > 0")
			}
			ct := models.ChiTietPhieuNhap{
				MaPhieuNhap: id, MaHangHoa: l.MaHangHoa,
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
		services.GhiLog(tx, c.GetString("user_id"), "sua_phieu_nhap", "PhieuNhapKho", id)
		return nil
	})
	h.ticketErr(c, err, id, "phiếu chỉ sửa được khi Chờ duyệt")
}

// DeleteImport - UC25 xóa phiếu nhập: chỉ khi "Chờ duyệt" hoặc "Từ chối" (tồn chưa đổi).
func (h *H) DeleteImport(c *gin.Context) {
	id := c.Param("id")
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuNhapKho
		if err := tx.First(&p, "ma_phieu_nhap = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai == models.StatusDaDuyet {
			return errBadState
		}
		if err := tx.Where("ma_phieu_nhap = ?", id).Delete(&models.ChiTietPhieuNhap{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&p).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "xoa_phieu_nhap", "PhieuNhapKho", id)
		return nil
	})
	h.ticketErr(c, err, id, "phiếu Đã duyệt không được xóa")
}

// ApproveImport - QL kho duyệt (duyetPhieu): chuyển "Đã duyệt" và TĂNG tồn từng dòng.
func (h *H) ApproveImport(c *gin.Context) {
	id := c.Param("id")
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuNhapKho
		if err := tx.Preload("ChiTiet").First(&p, "ma_phieu_nhap = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		for _, ct := range p.ChiTiet {
			if err := services.TangTon(tx, p.MaKho, ct.MaHangHoa, ct.SoLuong); err != nil {
				return err
			}
		}
		if err := tx.Model(&p).Update("trang_thai", models.StatusDaDuyet).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "duyet_phieu_nhap", "PhieuNhapKho", id)
		return nil
	})
	h.ticketErr(c, err, id, "chỉ duyệt được phiếu Chờ duyệt")
}

// RejectImport - QL kho từ chối (tuChoiPhieu): tồn không đổi, lý do ghi vào log.
func (h *H) RejectImport(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		LyDo string `json:"ly_do"`
	}
	_ = c.ShouldBindJSON(&body)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuNhapKho
		if err := tx.First(&p, "ma_phieu_nhap = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		if err := tx.Model(&p).Update("trang_thai", models.StatusTuChoi).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "tu_choi_phieu_nhap", "PhieuNhapKho", id+" - "+body.LyDo)
		return nil
	})
	h.ticketErr(c, err, id, "chỉ từ chối được phiếu Chờ duyệt")
}

// ticketErr quy đổi lỗi transaction phiếu thành HTTP response thống nhất.
func (h *H) ticketErr(c *gin.Context, err error, id, badStateMsg string) {
	switch {
	case err == nil:
		c.JSON(http.StatusOK, gin.H{"message": "thành công", "ma_phieu": id})
	case errors.Is(err, gorm.ErrRecordNotFound):
		notFound(c)
	case errors.Is(err, errBadState):
		badRequest(c, badStateMsg)
	default:
		badRequest(c, err.Error())
	}
}
