package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"warehouse-backend/internal/middleware"
	"warehouse-backend/internal/models"
)

// Register khai báo toàn bộ route theo ma trận phân quyền (Bảng 3.2):
//   - Tài khoản: chỉ QTV.
//   - Kho, duyệt phiếu, báo cáo: chỉ QL kho.
//   - Danh mục NCC / hàng hóa, lập / sửa / xóa phiếu, kiểm kê: NV kho + QL kho.
func Register(r *gin.Engine, gdb *gorm.DB, secret string) {
	h := &H{DB: gdb, Secret: secret}

	api := r.Group("/api")
	api.POST("/auth/login", h.Login)

	auth := api.Group("", middleware.Auth(secret))
	auth.POST("/auth/logout", h.Logout)

	qtv := auth.Group("", middleware.RequireRoles(models.RoleQTV))
	{
		qtv.GET("/users", h.ListUsers)
		qtv.POST("/users", h.CreateUser)
		qtv.PUT("/users/:id", h.UpdateUser)
		qtv.DELETE("/users/:id", h.DeleteUser)
		qtv.GET("/logs", h.ListLogs)
	}

	// danh mục NCC / hàng hóa: NV kho + QL kho
	staff := auth.Group("", middleware.RequireRoles(models.RoleNVKho, models.RoleQLKho))
	{
		staff.GET("/suppliers", h.ListSuppliers)
		staff.POST("/suppliers", h.CreateSupplier)
		staff.PUT("/suppliers/:id", h.UpdateSupplier)
		staff.DELETE("/suppliers/:id", h.DeleteSupplier)

		staff.GET("/products", h.ListProducts)
		staff.POST("/products", h.CreateProduct)
		staff.PUT("/products/:id", h.UpdateProduct)
		staff.DELETE("/products/:id", h.DeleteProduct)

		staff.GET("/groups", h.ListGroups)
		staff.GET("/units", h.ListUnits)
		staff.GET("/stocks", h.ListStocks)

		// phiếu nhập UC23-27
		staff.GET("/imports", h.ListImports)
		staff.GET("/imports/:id", h.GetImport)
		staff.POST("/imports", h.CreateImport)
		staff.PUT("/imports/:id", h.UpdateImport)
		staff.DELETE("/imports/:id", h.DeleteImport)

		// phiếu xuất UC28-32
		staff.GET("/exports", h.ListExports)
		staff.GET("/exports/:id", h.GetExport)
		staff.POST("/exports", h.CreateExport)
		staff.PUT("/exports/:id", h.UpdateExport)
		staff.DELETE("/exports/:id", h.DeleteExport)

		// kiểm kê UC33
		staff.GET("/stocktakes", h.ListStocktakes)
		staff.GET("/stocktakes/:id", h.GetStocktake)
		staff.POST("/stocktakes", h.CreateStocktake)
		staff.PUT("/stocktakes/:id", h.UpdateStocktake)
	}

	// kho, duyệt phiếu, báo cáo: chỉ QL kho
	manager := auth.Group("", middleware.RequireRoles(models.RoleQLKho))
	{
		manager.GET("/warehouses", h.ListWarehouses)
		manager.POST("/warehouses", h.CreateWarehouse)
		manager.PUT("/warehouses/:id", h.UpdateWarehouse)
		manager.DELETE("/warehouses/:id", h.DeleteWarehouse)

		manager.POST("/imports/:id/approve", h.ApproveImport)
		manager.POST("/imports/:id/reject", h.RejectImport)
		manager.POST("/exports/:id/approve", h.ApproveExport)
		manager.POST("/exports/:id/reject", h.RejectExport)
		manager.POST("/stocktakes/:id/approve", h.ApproveStocktake)
		manager.POST("/stocktakes/:id/reject", h.RejectStocktake)

		manager.GET("/reports/nxt", h.ReportNXT)
	}

	// NV kho cần xem danh sách kho để chọn trên form phiếu (chỉ đọc)
	auth.GET("/warehouses/options", middleware.RequireRoles(models.RoleNVKho, models.RoleQLKho), h.ListWarehouses)
}
