package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

// UC08-UC12: quản lý nhà cung cấp (NV kho, QL kho).

func (h *H) ListSuppliers(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.NhaCungCap{})
	if kw := c.Query("q"); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("ma_ncc ILIKE ? OR ten_ncc ILIKE ? OR sdt ILIKE ?", like, like, like)
	}
	var total int64
	q.Count(&total)
	var items []models.NhaCungCap
	if err := q.Order("ma_ncc").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	listResponse(c, items, total, page)
}

// CreateSupplier - UC08: kiểm tra trùng mã rồi lưu.
func (h *H) CreateSupplier(c *gin.Context) {
	var s models.NhaCungCap
	if err := c.ShouldBindJSON(&s); err != nil || s.TenNCC == "" {
		badRequest(c, "cần nhập tối thiểu tên nhà cung cấp")
		return
	}
	if s.MaNCC == "" {
		code, err := nextCode(h.DB, "nha_cung_cap", "ma_ncc", "NCC")
		if err != nil {
			serverErr(c, err)
			return
		}
		s.MaNCC = code
	} else {
		var cnt int64
		h.DB.Model(&models.NhaCungCap{}).Where("ma_ncc = ?", s.MaNCC).Count(&cnt)
		if cnt > 0 {
			badRequest(c, "mã nhà cung cấp đã tồn tại")
			return
		}
	}
	s.TrangThai = true
	if err := h.DB.Create(&s).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "them_ncc", "NhaCungCap", s.MaNCC)
	c.JSON(http.StatusCreated, s)
}

// UpdateSupplier - UC09: mã không đổi, cập nhật các trường còn lại.
func (h *H) UpdateSupplier(c *gin.Context) {
	var s models.NhaCungCap
	if err := h.DB.First(&s, "ma_ncc = ?", c.Param("id")).Error; err != nil {
		notFound(c)
		return
	}
	var req models.NhaCungCap
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}
	s.TenNCC, s.DiaChi, s.SDT, s.Email, s.NguoiLienHe, s.TrangThai =
		req.TenNCC, req.DiaChi, req.SDT, req.Email, req.NguoiLienHe, req.TrangThai
	if err := h.DB.Save(&s).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "sua_ncc", "NhaCungCap", s.MaNCC)
	c.JSON(http.StatusOK, s)
}

// DeleteSupplier - UC10: đã có phiếu nhập tham chiếu thì chỉ NGỪNG hoạt động.
func (h *H) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	var s models.NhaCungCap
	if err := h.DB.First(&s, "ma_ncc = ?", id).Error; err != nil {
		notFound(c)
		return
	}
	var refs int64
	h.DB.Model(&models.PhieuNhapKho{}).Where("ma_ncc = ?", id).Count(&refs)
	if refs > 0 {
		h.DB.Model(&s).Update("trang_thai", false)
		services.GhiLog(h.DB, c.GetString("user_id"), "ngung_ncc", "NhaCungCap", id)
		c.JSON(http.StatusOK, gin.H{"message": "NCC đã có phiếu nhập, chuyển sang NGỪNG hoạt động", "locked": true})
		return
	}
	if err := h.DB.Delete(&s).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "xoa_ncc", "NhaCungCap", id)
	c.JSON(http.StatusOK, gin.H{"message": "đã xóa nhà cung cấp", "locked": false})
}
