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

// UC33: kiểm kê tồn kho. NV kho tạo phiếu (hệ thống nạp SL sổ sách từ TonKho),
// nhập SL thực tế; QL kho duyệt điều chỉnh thì tồn được GÁN = số thực tế.

func (h *H) ListStocktakes(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.PhieuKiemKe{})
	if v := c.Query("trang_thai"); v != "" {
		q = q.Where("trang_thai = ?", v)
	}
	if v := c.Query("ma_kho"); v != "" {
		q = q.Where("ma_kho = ?", v)
	}
	var total int64
	q.Count(&total)
	var items []models.PhieuKiemKe
	if err := q.Order("ma_phieu_kiem_ke DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	listResponse(c, items, total, page)
}

func (h *H) GetStocktake(c *gin.Context) {
	var p models.PhieuKiemKe
	if err := h.DB.Preload("ChiTiet").First(&p, "ma_phieu_kiem_ke = ?", c.Param("id")).Error; err != nil {
		notFound(c)
		return
	}
	c.JSON(http.StatusOK, p)
}

// CreateStocktake - taoPhieu: tạo phiếu theo kho, nạp SL sổ sách từ TonKho (layTon).
func (h *H) CreateStocktake(c *gin.Context) {
	var req struct {
		MaKho string `json:"ma_kho" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "cần chọn kho kiểm kê")
		return
	}
	var p models.PhieuKiemKe
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var tons []models.TonKho
		if err := tx.Where("ma_kho = ?", req.MaKho).Order("ma_hang_hoa").Find(&tons).Error; err != nil {
			return err
		}
		if len(tons) == 0 {
			return errors.New("kho chưa có dữ liệu tồn để kiểm kê")
		}
		code, err := nextCode(tx, "phieu_kiem_ke", "ma_phieu_kiem_ke", "KK")
		if err != nil {
			return err
		}
		p = models.PhieuKiemKe{
			MaPhieuKiemKe: code, MaKho: req.MaKho, MaNguoiDung: c.GetString("user_id"),
			NgayKiem: time.Now(), TrangThai: models.StatusChoDuyet,
		}
		for _, t := range tons {
			p.ChiTiet = append(p.ChiTiet, models.ChiTietKiemKe{
				MaPhieuKiemKe: code, MaHangHoa: t.MaHangHoa,
				SLSoSach: t.SoLuongTon, SLThucTe: t.SoLuongTon, ChenhLech: 0,
			})
		}
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "lap_phieu_kiem_ke", "PhieuKiemKe", code)
		return nil
	})
	if err != nil {
		badRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusCreated, p)
}

// UpdateStocktake - nhapSoThucTe + tinhChenhLech: ghi số thực tế cho các dòng (khi Chờ duyệt).
func (h *H) UpdateStocktake(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ChiTiet []struct {
			MaHangHoa string `json:"ma_hang_hoa" binding:"required"`
			SLThucTe  int    `json:"sl_thuc_te"`
		} `json:"chi_tiet" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuKiemKe
		if err := tx.First(&p, "ma_phieu_kiem_ke = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		for _, l := range req.ChiTiet {
			if l.SLThucTe < 0 {
				return errors.New("số lượng thực tế phải >= 0")
			}
			res := tx.Model(&models.ChiTietKiemKe{}).
				Where("ma_phieu_kiem_ke = ? AND ma_hang_hoa = ?", id, l.MaHangHoa).
				Updates(map[string]any{
					"sl_thuc_te": l.SLThucTe,
					"chenh_lech": gorm.Expr("? - sl_so_sach", l.SLThucTe),
				})
			if res.Error != nil {
				return res.Error
			}
		}
		services.GhiLog(tx, c.GetString("user_id"), "nhap_so_thuc_te", "PhieuKiemKe", id)
		return nil
	})
	h.ticketErr(c, err, id, "phiếu kiểm kê chỉ sửa được khi Chờ duyệt")
}

// ApproveStocktake - duyetDieuChinh: QL kho duyệt, gán TonKho = số thực tế (ganTheoKiemKe).
func (h *H) ApproveStocktake(c *gin.Context) {
	id := c.Param("id")
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuKiemKe
		if err := tx.Preload("ChiTiet").First(&p, "ma_phieu_kiem_ke = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		for _, ct := range p.ChiTiet {
			if ct.ChenhLech != 0 {
				if err := services.GanTheoKiemKe(tx, p.MaKho, ct.MaHangHoa, ct.SLThucTe); err != nil {
					return err
				}
			}
		}
		if err := tx.Model(&p).Update("trang_thai", models.StatusDaDuyet).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "duyet_kiem_ke", "PhieuKiemKe", id)
		return nil
	})
	h.ticketErr(c, err, id, "chỉ duyệt được phiếu Chờ duyệt")
}

// RejectStocktake - không duyệt: yêu cầu kiểm lại, tồn không đổi.
func (h *H) RejectStocktake(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		LyDo string `json:"ly_do"`
	}
	_ = c.ShouldBindJSON(&body)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var p models.PhieuKiemKe
		if err := tx.First(&p, "ma_phieu_kiem_ke = ?", id).Error; err != nil {
			return gorm.ErrRecordNotFound
		}
		if p.TrangThai != models.StatusChoDuyet {
			return errBadState
		}
		if err := tx.Model(&p).Update("trang_thai", models.StatusTuChoi).Error; err != nil {
			return err
		}
		services.GhiLog(tx, c.GetString("user_id"), "tu_choi_kiem_ke", "PhieuKiemKe", id+" - "+body.LyDo)
		return nil
	})
	h.ticketErr(c, err, id, "chỉ từ chối được phiếu Chờ duyệt")
}
