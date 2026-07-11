package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"warehouse-backend/internal/middleware"
	"warehouse-backend/internal/models"
	"warehouse-backend/internal/services"
)

// Login - UC01 Đăng nhập: xác thực tên đăng nhập / mật khẩu, kiểm tra trạng thái,
// tạo "phiên" dưới dạng JWT và ghi log.
func (h *H) Login(c *gin.Context) {
	var req struct {
		TenDangNhap string `json:"ten_dang_nhap" binding:"required"`
		MatKhau     string `json:"mat_khau" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "cần nhập đủ tên đăng nhập và mật khẩu")
		return
	}
	var u models.NguoiDung
	if err := h.DB.Where("ten_dang_nhap = ?", req.TenDangNhap).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sai tên đăng nhập hoặc mật khẩu"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.MatKhau), []byte(req.MatKhau)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sai tên đăng nhập hoặc mật khẩu"})
		return
	}
	if !u.TrangThai {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tài khoản đã bị khóa, liên hệ quản trị viên"})
		return
	}
	claims := middleware.Claims{
		MaNguoiDung: u.MaNguoiDung,
		HoTen:       u.HoTen,
		VaiTro:      u.VaiTro,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.MaNguoiDung,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.Secret))
	if err != nil {
		serverErr(c, err)
		return
	}
	services.GhiLog(h.DB, u.MaNguoiDung, "dang_nhap", "NguoiDung", "Đăng nhập thành công")
	c.JSON(http.StatusOK, gin.H{"token": tok, "user": u})
}

// Logout - UC02 Đăng xuất: JWT stateless nên server chỉ ghi log; client xóa token.
func (h *H) Logout(c *gin.Context) {
	services.GhiLog(h.DB, c.GetString("user_id"), "dang_xuat", "NguoiDung", "Đăng xuất")
	c.JSON(http.StatusOK, gin.H{"message": "đã đăng xuất"})
}
