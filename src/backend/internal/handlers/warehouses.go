package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

var errDuplicate = errors.New("duplicate key")

// UC18-UC22: quản lý kho (chỉ QL kho).

func (h *H) ListWarehouses(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.Kho{})
	if kw := c.Query("q"); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("ma_kho ILIKE ? OR ten_kho ILIKE ? OR dia_chi ILIKE ?", like, like, like)
	}
	var total int64
	q.Count(&total)
	var items []models.Kho
	if err := q.Order("ma_kho").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	listResponse(c, items, total, page)
}

// CreateWarehouse - UC18: kiểm tra trùng mã; khởi tạo tồn = 0 cho mọi hàng đang dùng.
func (h *H) CreateWarehouse(c *gin.Context) {
	var k models.Kho
	if err := c.ShouldBindJSON(&k); err != nil || k.TenKho == "" {
		badRequest(c, "cần nhập tối thiểu tên kho")
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if k.MaKho == "" {
			code, err := nextCode(tx, "kho", "ma_kho", "K")
			if err != nil {
				return err
			}
			k.MaKho = code
		} else {
			var cnt int64
			tx.Model(&models.Kho{}).Where("ma_kho = ?", k.MaKho).Count(&cnt)
			if cnt > 0 {
				return errDuplicate
			}
		}
		k.TrangThai = true
		if err := tx.Create(&k).Error; err != nil {
			return err
		}
		var products []models.HangHoa
		tx.Where("trang_thai = true").Find(&products)
		for _, p := range products {
			if err := services.KhoiTaoTon(tx, k.MaKho, p.MaHangHoa); err != nil {
				return err
			}
		}
		services.GhiLog(tx, c.GetString("user_id"), "them_kho", "Kho", k.MaKho)
		return nil
	})
	if err == errDuplicate {
		badRequest(c, "mã kho đã tồn tại")
		return
	}
	if err != nil {
		serverErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, k)
}

// UpdateWarehouse - UC19 sửa kho.
func (h *H) UpdateWarehouse(c *gin.Context) {
	var k models.Kho
	if err := h.DB.First(&k, "ma_kho = ?", c.Param("id")).Error; err != nil {
		notFound(c)
		return
	}
	var req models.Kho
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}
	k.TenKho, k.DiaChi, k.NguoiPhuTrach, k.TrangThai = req.TenKho, req.DiaChi, req.NguoiPhuTrach, req.TrangThai
	if err := h.DB.Save(&k).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "sua_kho", "Kho", k.MaKho)
	c.JSON(http.StatusOK, k)
}

// DeleteWarehouse - UC20: kho còn tồn > 0 hoặc còn phiếu tham chiếu thì chỉ NGỪNG hoạt động.
func (h *H) DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	var k models.Kho
	if err := h.DB.First(&k, "ma_kho = ?", id).Error; err != nil {
		notFound(c)
		return
	}
	var tonDuong, refs, n2 int64
	h.DB.Model(&models.TonKho{}).Where("ma_kho = ? AND so_luong_ton > 0", id).Count(&tonDuong)
	h.DB.Model(&models.PhieuNhapKho{}).Where("ma_kho = ?", id).Count(&refs)
	h.DB.Model(&models.PhieuXuatKho{}).Where("ma_kho = ?", id).Count(&n2)
	refs += n2
	if tonDuong > 0 || refs > 0 {
		h.DB.Model(&k).Update("trang_thai", false)
		services.GhiLog(h.DB, c.GetString("user_id"), "ngung_kho", "Kho", id)
		c.JSON(http.StatusOK, gin.H{"message": "kho còn tồn / phiếu, chuyển sang NGỪNG hoạt động", "locked": true})
		return
	}
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("ma_kho = ?", id).Delete(&models.TonKho{}).Error; err != nil {
			return err
		}
		return tx.Delete(&k).Error
	}); err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "xoa_kho", "Kho", id)
	c.JSON(http.StatusOK, gin.H{"message": "đã xóa kho", "locked": false})
}

// ListStocks - tra cứu tồn kho (phục vụ UC16 kèm tồn, form phiếu xuất, trang tồn kho).
func (h *H) ListStocks(c *gin.Context) {
	q := h.DB.Table("ton_kho t").
		Select("t.ma_kho, k.ten_kho, t.ma_hang_hoa, hh.ten_hang_hoa, t.so_luong_ton, t.ngay_cap_nhat").
		Joins("JOIN kho k ON k.ma_kho = t.ma_kho").
		Joins("JOIN hang_hoa hh ON hh.ma_hang_hoa = t.ma_hang_hoa")
	if v := c.Query("ma_kho"); v != "" {
		q = q.Where("t.ma_kho = ?", v)
	}
	if v := c.Query("ma_hang_hoa"); v != "" {
		q = q.Where("t.ma_hang_hoa = ?", v)
	}
	var rows []map[string]any
	if err := q.Order("t.ma_kho, t.ma_hang_hoa").Find(&rows).Error; err != nil {
		serverErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rows})
}

// ListGroups / ListUnits - danh mục phụ cho form hàng hóa.
func (h *H) ListGroups(c *gin.Context) {
	var items []models.NhomHangHoa
	h.DB.Order("ma_nhom_hang").Find(&items)
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *H) ListUnits(c *gin.Context) {
	var items []models.DonViTinh
	h.DB.Order("ma_don_vi_tinh").Find(&items)
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// ListLogs - tra cứu lịch sử thao tác (QTV).
func (h *H) ListLogs(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.LichSuThaoTac{})
	if v := c.Query("ma_nguoi_dung"); v != "" {
		q = q.Where("ma_nguoi_dung = ?", v)
	}
	var total int64
	q.Count(&total)
	var items []models.LichSuThaoTac
	if err := q.Order("ma_log DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	listResponse(c, items, total, page)
}
