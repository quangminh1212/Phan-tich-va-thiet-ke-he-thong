package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	MaNguoiDung string `json:"ma_nguoi_dung"`
	HoTen       string `json:"ho_ten"`
	VaiTro      string `json:"vai_tro"`
	jwt.RegisteredClaims
}

// Auth xác thực JWT Bearer (UC01) và gắn thông tin người dùng vào context.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "thiếu token"})
			return
		}
		claims := &Claims{}
		tok, err := jwt.ParseWithClaims(strings.TrimPrefix(h, "Bearer "), claims,
			func(t *jwt.Token) (any, error) { return []byte(secret), nil })
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token không hợp lệ"})
			return
		}
		c.Set("user_id", claims.MaNguoiDung)
		c.Set("user_name", claims.HoTen)
		c.Set("user_role", claims.VaiTro)
		c.Next()
	}
}

// RequireRoles chặn theo ma trận phân quyền (Bảng 3.2 của báo cáo).
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		if !allowed[c.GetString("user_role")] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "không có quyền thực hiện chức năng này"})
			return
		}
		c.Next()
	}
}
