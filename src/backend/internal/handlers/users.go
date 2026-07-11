package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

// UC03-UC07: quản lý tài khoản (chỉ QTV).

// ListUsers - UC07 xem danh sách + UC06 tìm kiếm (?q=, ?vai_tro=).
func (h *H) ListUsers(c *gin.Context) {
	offset, limit, page := paginate(c)
	q := h.DB.Model(&models.NguoiDung{})
	if kw := c.Query("q"); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("ten_dang_nhap ILIKE ? OR ho_ten ILIKE ?", like, like)
	}
	if r := c.Query("vai_tro"); r != "" {
		q = q.Where("vai_tro = ?", r)
	}
	var total int64
	q.Count(&total)
	var items []models.NguoiDung
	if err := q.Order("ma_nguoi_dung").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		serverErr(c, err)
		return
	}
	listResponse(c, items, total, page)
}

type userReq struct {
	MaNguoiDung string `json:"ma_nguoi_dung"`
	TenDangNhap string `json:"ten_dang_nhap" binding:"required"`
	MatKhau     string `json:"mat_khau"`
	HoTen       string `json:"ho_ten" binding:"required"`
	Email       string `json:"email"`
	VaiTro      string `json:"vai_tro" binding:"required"`
	TrangThai   *bool  `json:"trang_thai"`
}

// CreateUser - UC03 thêm tài khoản: kiểm tra trùng tên đăng nhập, mã hóa mật khẩu.
func (h *H) CreateUser(c *gin.Context) {
	var req userReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "thiếu trường bắt buộc: "+err.Error())
		return
	}
	if req.MatKhau == "" {
		badRequest(c, "cần nhập mật khẩu")
		return
	}
	var cnt int64
	h.DB.Model(&models.NguoiDung{}).Where("ten_dang_nhap = ?", req.TenDangNhap).Count(&cnt)
	if cnt > 0 {
		badRequest(c, "tên đăng nhập đã tồn tại")
		return
	}
	if req.MaNguoiDung == "" {
		code, err := nextCode(h.DB, "nguoi_dung", "ma_nguoi_dung", "ND")
		if err != nil {
			serverErr(c, err)
			return
		}
		req.MaNguoiDung = code
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.MatKhau), bcrypt.DefaultCost)
	u := models.NguoiDung{
		MaNguoiDung: req.MaNguoiDung, TenDangNhap: req.TenDangNhap, MatKhau: string(hash),
		HoTen: req.HoTen, Email: req.Email, VaiTro: req.VaiTro, TrangThai: true,
	}
	if err := h.DB.Create(&u).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "them_tai_khoan", "NguoiDung", u.MaNguoiDung)
	c.JSON(http.StatusCreated, u)
}

// UpdateUser - UC04 sửa tài khoản (tên đăng nhập không đổi; đổi mật khẩu nếu gửi kèm).
func (h *H) UpdateUser(c *gin.Context) {
	var u models.NguoiDung
	if err := h.DB.First(&u, "ma_nguoi_dung = ?", c.Param("id")).Error; err != nil {
		notFound(c)
		return
	}
	var req userReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "thiếu trường bắt buộc: "+err.Error())
		return
	}
	u.HoTen, u.Email, u.VaiTro = req.HoTen, req.Email, req.VaiTro
	if req.TrangThai != nil {
		u.TrangThai = *req.TrangThai
	}
	if req.MatKhau != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.MatKhau), bcrypt.DefaultCost)
		u.MatKhau = string(hash)
	}
	if err := h.DB.Save(&u).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "sua_tai_khoan", "NguoiDung", u.MaNguoiDung)
	c.JSON(http.StatusOK, u)
}

// DeleteUser - UC05 xóa tài khoản: đã phát sinh phiếu / log thì chỉ KHÓA, không xóa cứng;
// không cho tự xóa chính mình.
func (h *H) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == c.GetString("user_id") {
		badRequest(c, "không thể xóa tài khoản đang đăng nhập")
		return
	}
	var u models.NguoiDung
	if err := h.DB.First(&u, "ma_nguoi_dung = ?", id).Error; err != nil {
		notFound(c)
		return
	}
	var refs int64
	h.DB.Model(&models.PhieuNhapKho{}).Where("ma_nguoi_dung = ?", id).Count(&refs)
	var n2, n3, n4 int64
	h.DB.Model(&models.PhieuXuatKho{}).Where("ma_nguoi_dung = ?", id).Count(&n2)
	h.DB.Model(&models.PhieuKiemKe{}).Where("ma_nguoi_dung = ?", id).Count(&n3)
	h.DB.Model(&models.LichSuThaoTac{}).Where("ma_nguoi_dung = ?", id).Count(&n4)
	refs += n2 + n3 + n4
	if refs > 0 {
		h.DB.Model(&u).Update("trang_thai", false)
		services.GhiLog(h.DB, c.GetString("user_id"), "khoa_tai_khoan", "NguoiDung", id)
		c.JSON(http.StatusOK, gin.H{"message": "tài khoản đã phát sinh dữ liệu, chuyển sang KHÓA", "locked": true})
		return
	}
	if err := h.DB.Delete(&u).Error; err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, c.GetString("user_id"), "xoa_tai_khoan", "NguoiDung", id)
	c.JSON(http.StatusOK, gin.H{"message": "đã xóa tài khoản", "locked": false})
}
