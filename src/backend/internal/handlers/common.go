package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// H gom DB + secret cho mọi handler.
type H struct {
	DB     *gorm.DB
	Secret string
}

func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "không tìm thấy bản ghi"})
}

func serverErr(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// paginate đọc ?page=&page_size= (mặc định 1 / 20, tối đa 200).
func paginate(c *gin.Context) (offset, limit, page int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 200 {
		size = 20
	}
	return (page - 1) * size, size, page
}

func listResponse(c *gin.Context, items any, total int64, page int) {
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": page})
}

// nextCode sinh mã phiếu dạng <prefix><số thứ tự 5 chữ số> trong transaction.
func nextCode(tx *gorm.DB, table, column, prefix string) (string, error) {
	var maxCode *string
	row := tx.Table(table).Select("MAX(" + column + ")").Row()
	if err := row.Scan(&maxCode); err != nil {
		return "", err
	}
	n := 0
	if maxCode != nil && len(*maxCode) > len(prefix) {
		fmt.Sscanf((*maxCode)[len(prefix):], "%d", &n)
	}
	return fmt.Sprintf("%s%05d", prefix, n+1), nil
}
