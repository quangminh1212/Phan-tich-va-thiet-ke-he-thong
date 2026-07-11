package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

// UC13-UC17: quản lý hàng hóa (NV kho, QL kho).

// ListProducts - UC17 xem danh sách + UC16 tìm kiếm (?q=, ?ma_nhom_hang=), trả kèm tổng tồn.
func (h *H) ListProducts(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.HangHoa{})
	if kw := c.Query("q"); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("ma_hang_hoa ILIKE ? OR ten_hang_hoa ILIKE ?", like, like)
	}
	if g := c.Query("ma_nhom_hang"); g != "" {
		q = q.Where("ma_nhom_hang = ?", g)
	}
	var total int64
	q.Count(&total)
	var items []models.HangHoa
	if err := q.Order("ma_hang_hoa").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	// đọc kèm tồn (TonKho.layTon) cho các mã trong trang
	type row struct {
		MaHangHoa string
		Ton       int
	}
	tons := map[string]int{}
	if len(items) > 0 {
		codes := make([]string, len(items))
		for i, it := range items {
			codes[i] = it.MaHangHoa
		}
		var rows []row
		h.DB.Model(&models.TonKho{}).Select("ma_hang_hoa, COALESCE(SUM(so_luong_ton),0) AS ton").
			Where("ma_hang_hoa IN ?", codes).Group("ma_hang_hoa").Scan(&rows)
		for _, r := range rows {
			tons[r.MaHangHoa] = r.Ton
		}
	}
	out := make([]gin.H, len(items))
	for i, it := range items {
		out[i] = gin.H{
			"ma_hang_hoa": it.MaHangHoa, "ten_hang_hoa": it.TenHangHoa,
			"ma_nhom_hang": it.MaNhomHang, "ma_don_vi_tinh": it.MaDonViTinh,
			"gia_nhap": it.GiaNhap, "gia_xuat": it.GiaXuat, "mo_ta": it.MoTa,
			"trang_thai": it.TrangThai, "tong_ton": tons[it.MaHangHoa],
		}
	}
	listResponse(c, out, total, page)
}

// CreateProduct - UC13 thêm hàng hóa: kiểm tra trùng mã, lưu, KHỞI TẠO TỒN = 0 tại mọi kho.
func (h *H) CreateProduct(c *gin.Context) {
	var p models.HangHoa
	if err := c.ShouldBindJSON(&p); err != nil || p.TenHangHoa == "" {
		badRequest(c, "cần nhập tối thiểu tên hàng hóa")
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if p.MaHangHoa == "" {
			code, err := nextCode(tx, "hang_hoa", "ma_hang_hoa", "HH")
			if err != nil {
				return err
			}
			p.MaHangHoa = code
		} else {
			var cnt int64
			tx.Model(&models.HangHoa{}).Where("ma_hang_hoa = ?", p.MaHangHoa).Count(&cnt)
			if cnt > 0 {
				return errDuplicate
			}
		}
		p.TrangThai = true
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		var khoList []models.Kho
		tx.Where("trang_thai = true").Find(&khoList)
		for _, k := range khoList {
			if err := services.KhoiTaoTon(tx, k.MaKho, p.MaHangHoa); err != nil {
				return err
			}
		}
		services.GhiLog(tx, c.GetString("user_id"), "them_hang_hoa", "HangHoa", p.MaHangHoa)
		return nil
	})
	if err == errDuplicate {
		badRequest(c, "mã hàng hóa đã tồn tại")
		return
	}
	if err != nil {
		serverErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

// UpdateProduct - UC14 sửa hàng hóa (mã không đổi).
func (h *H) UpdateProduct(c *gin.Context) {
	var p models.HangHoa
	if err := h.DB.First(&p, "ma_hang_hoa = ?", c.Param("id")).Error; err != nil {
		notFound(c)
		return
	}
	var req models.HangHoa
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}
	p.TenHangHoa, p.MaNhomHang, p.MaDonViTinh = req.TenHangHoa, req.MaNhomHang, req.MaDonViTinh
	p.GiaNhap, p.GiaXuat, p.MoTa, p.TrangThai = req.GiaNhap, req.GiaXuat, req.MoTa, req.TrangThai
	if err := h.DB.Save(&p).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "sua_hang_hoa", "HangHoa", p.MaHangHoa)
	c.JSON(http.StatusOK, p)
}

// DeleteProduct - UC15: hàng đã có trên phiếu nhập / xuất thì chỉ NGỪNG sử dụng.
func (h *H) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var p models.HangHoa
	if err := h.DB.First(&p, "ma_hang_hoa = ?", id).Error; err != nil {
		notFound(c)
		return
	}
	var refs, n2 int64
	h.DB.Model(&models.ChiTietPhieuNhap{}).Where("ma_hang_hoa = ?", id).Count(&refs)
	h.DB.Model(&models.ChiTietPhieuXuat{}).Where("ma_hang_hoa = ?", id).Count(&n2)
	refs += n2
	if refs > 0 {
		h.DB.Model(&p).Update("trang_thai", false)
		services.GhiLog(h.DB, c.GetString("user_id"), "ngung_hang_hoa", "HangHoa", id)
		c.JSON(http.StatusOK, gin.H{"message": "hàng đã phát sinh phiếu, chuyển sang NGỪNG sử dụng", "locked": true})
		return
	}
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("ma_hang_hoa = ?", id).Delete(&models.TonKho{}).Error; err != nil {
			return err
		}
		return tx.Delete(&p).Error
	}); err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "xoa_hang_hoa", "HangHoa", id)
	c.JSON(http.StatusOK, gin.H{"message": "đã xóa hàng hóa", "locked": false})
}
