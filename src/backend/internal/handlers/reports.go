package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UC34: báo cáo nhập - xuất - tồn (chỉ QL kho, chỉ đọc từ phiếu ĐÃ DUYỆT và TonKho).
// GET /api/reports/nxt?from=2026-07-01&to=2026-07-31&ma_kho=&ma_hang_hoa=
//
// Với mỗi mặt hàng: nhap / xuat = tổng SL trên phiếu đã duyệt trong kỳ;
// dieu_chinh = tổng chênh lệch kiểm kê đã duyệt trong kỳ; ton_cuoi = tồn hiện tại;
// ton_dau suy ra: ton_cuoi - nhap + xuat - dieu_chinh.
func (h *H) ReportNXT(c *gin.Context) {
	from := c.DefaultQuery("from", "1970-01-01")
	to := c.DefaultQuery("to", time.Now().Format("2006-01-02"))
	maKho := c.Query("ma_kho")
	maHang := c.Query("ma_hang_hoa")

	type agg struct {
		MaHangHoa string
		Total     int
	}
	sum := func(rows []agg) map[string]int {
		m := map[string]int{}
		for _, r := range rows {
			m[r.MaHangHoa] = r.Total
		}
		return m
	}

	// nhập trong kỳ (phiếu đã duyệt)
	var nhapRows []agg
	qn := h.DB.Table("chi_tiet_phieu_nhap ct").
		Select("ct.ma_hang_hoa, COALESCE(SUM(ct.so_luong),0) AS total").
		Joins("JOIN phieu_nhap_kho p ON p.ma_phieu_nhap = ct.ma_phieu_nhap").
		Where("p.trang_thai = 'da_duyet' AND p.ngay_nhap::date BETWEEN ? AND ?", from, to).
		Group("ct.ma_hang_hoa")
	if maKho != "" {
		qn = qn.Where("p.ma_kho = ?", maKho)
	}
	qn.Scan(&nhapRows)

	// xuất trong kỳ (phiếu đã duyệt)
	var xuatRows []agg
	qx := h.DB.Table("chi_tiet_phieu_xuat ct").
		Select("ct.ma_hang_hoa, COALESCE(SUM(ct.so_luong),0) AS total").
		Joins("JOIN phieu_xuat_kho p ON p.ma_phieu_xuat = ct.ma_phieu_xuat").
		Where("p.trang_thai = 'da_duyet' AND p.ngay_xuat::date BETWEEN ? AND ?", from, to).
		Group("ct.ma_hang_hoa")
	if maKho != "" {
		qx = qx.Where("p.ma_kho = ?", maKho)
	}
	qx.Scan(&xuatRows)

	// điều chỉnh kiểm kê trong kỳ (đã duyệt)
	var dcRows []agg
	qd := h.DB.Table("chi_tiet_kiem_ke ct").
		Select("ct.ma_hang_hoa, COALESCE(SUM(ct.chenh_lech),0) AS total").
		Joins("JOIN phieu_kiem_ke p ON p.ma_phieu_kiem_ke = ct.ma_phieu_kiem_ke").
		Where("p.trang_thai = 'da_duyet' AND p.ngay_kiem::date BETWEEN ? AND ?", from, to).
		Group("ct.ma_hang_hoa")
	if maKho != "" {
		qd = qd.Where("p.ma_kho = ?", maKho)
	}
	qd.Scan(&dcRows)

	// tồn hiện tại
	var tonRows []agg
	qt := h.DB.Table("ton_kho").
		Select("ma_hang_hoa, COALESCE(SUM(so_luong_ton),0) AS total").
		Group("ma_hang_hoa")
	if maKho != "" {
		qt = qt.Where("ma_kho = ?", maKho)
	}
	qt.Scan(&tonRows)

	nhap, xuat, dc, ton := sum(nhapRows), sum(xuatRows), sum(dcRows), sum(tonRows)

	type prod struct {
		MaHangHoa  string
		TenHangHoa string
	}
	var prods []prod
	qp := h.DB.Table("hang_hoa").Select("ma_hang_hoa, ten_hang_hoa").Order("ma_hang_hoa")
	if maHang != "" {
		qp = qp.Where("ma_hang_hoa = ?", maHang)
	}
	qp.Scan(&prods)

	out := []gin.H{}
	for _, p := range prods {
		n, x, d, t := nhap[p.MaHangHoa], xuat[p.MaHangHoa], dc[p.MaHangHoa], ton[p.MaHangHoa]
		if n == 0 && x == 0 && d == 0 && t == 0 {
			continue // bỏ dòng không phát sinh và không tồn
		}
		out = append(out, gin.H{
			"ma_hang_hoa":  p.MaHangHoa,
			"ten_hang_hoa": p.TenHangHoa,
			"ton_dau":      t - n + x - d,
			"nhap":         n,
			"xuat":         x,
			"dieu_chinh":   d,
			"ton_cuoi":     t,
		})
	}
	c.JSON(http.StatusOK, gin.H{"from": from, "to": to, "ma_kho": maKho, "items": out})
}
