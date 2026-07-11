import { useCallback, useEffect, useState } from 'react';
import { api } from '../api/client';
import { useAuth } from '../auth/AuthContext';
import Modal from '../components/Modal';
import { StatusBadge } from '../components/TicketPage';

// UC33: lập phiếu kiểm kê - tạo theo kho (nạp SL sổ), nhập SL thực tế,
// QL kho duyệt điều chỉnh thì tồn = số thực tế.
export default function StocktakePage() {
  const { user } = useAuth();
  const isManager = user?.vai_tro === 'QLKho';

  const [items, setItems] = useState([]);
  const [warehouses, setWarehouses] = useState([]);
  const [maKho, setMaKho] = useState('');
  const [detail, setDetail] = useState(null);
  const [msg, setMsg] = useState('');
  const [err, setErr] = useState('');

  const load = useCallback(async () => {
    try {
      const data = await api('/stocktakes?page_size=100');
      setItems(data.items || []);
    } catch (e) { setErr(e.message); }
  }, []);

  useEffect(() => {
    load();
    api('/warehouses/options?page_size=200').then((d) => setWarehouses(d.items || [])).catch(() => {});
  }, [load]);

  async function create() {
    if (!maKho) { setErr('Chọn kho cần kiểm kê'); return; }
    setErr('');
    try {
      const p = await api('/stocktakes', { method: 'POST', body: { ma_kho: maKho } });
      setMsg(`Đã tạo phiếu ${p.ma_phieu_kiem_ke} (SL sổ sách đã nạp từ tồn kho)`);
      setDetail(p);
      load();
    } catch (e) { setErr(e.message); }
  }

  async function openDetail(row) {
    try { setDetail(await api(`/stocktakes/${row.ma_phieu_kiem_ke}`)); } catch (e) { setErr(e.message); }
  }

  function setActual(i, v) {
    const ct = [...detail.chi_tiet];
    const sl = Number(v);
    ct[i] = { ...ct[i], sl_thuc_te: sl, chenh_lech: sl - ct[i].sl_so_sach };
    setDetail({ ...detail, chi_tiet: ct });
  }

  async function saveActual() {
    try {
      await api(`/stocktakes/${detail.ma_phieu_kiem_ke}`, {
        method: 'PUT',
        body: { chi_tiet: detail.chi_tiet.map((l) => ({ ma_hang_hoa: l.ma_hang_hoa, sl_thuc_te: Number(l.sl_thuc_te) })) },
      });
      setMsg('Đã lưu số lượng thực tế');
    } catch (e) { setErr(e.message); }
  }

  async function decide(action) {
    const lyDo = action === 'reject' ? window.prompt('Lý do yêu cầu kiểm lại?') || '' : '';
    try {
      await api(`/stocktakes/${detail.ma_phieu_kiem_ke}/${action}`, { method: 'POST', body: { ly_do: lyDo } });
      setMsg(action === 'approve' ? 'Đã duyệt điều chỉnh - tồn kho gán theo số thực tế' : 'Đã từ chối');
      setDetail(null);
      load();
    } catch (e) { setErr(e.message); }
  }

  return (
    <div>
      <h2>Kiểm kê tồn kho</h2>
      <div className="toolbar">
        <select value={maKho} onChange={(e) => setMaKho(e.target.value)}>
          <option value="">-- chọn kho --</option>
          {warehouses.filter((k) => k.trang_thai).map((k) => (
            <option key={k.ma_kho} value={k.ma_kho}>{k.ten_kho}</option>
          ))}
        </select>
        <button onClick={create}>Lập phiếu kiểm kê</button>
      </div>
      {msg && <div className="info">{msg}</div>}
      {err && !detail && <div className="error">{err}</div>}
      <table>
        <thead><tr><th>Mã phiếu</th><th>Kho</th><th>Ngày kiểm</th><th>Người lập</th><th>Trạng thái</th><th></th></tr></thead>
        <tbody>
          {items.map((row) => (
            <tr key={row.ma_phieu_kiem_ke}>
              <td>{row.ma_phieu_kiem_ke}</td>
              <td>{row.ma_kho}</td>
              <td>{(row.ngay_kiem || '').slice(0, 10)}</td>
              <td>{row.ma_nguoi_dung}</td>
              <td><StatusBadge value={row.trang_thai} /></td>
              <td><button className="secondary" onClick={() => openDetail(row)}>Mở</button></td>
            </tr>
          ))}
          {items.length === 0 && <tr><td colSpan={6} className="muted">Chưa có phiếu kiểm kê</td></tr>}
        </tbody>
      </table>

      {detail && (
        <Modal title={`Phiếu kiểm kê ${detail.ma_phieu_kiem_ke} - kho ${detail.ma_kho}`} onClose={() => setDetail(null)}>
          <p><StatusBadge value={detail.trang_thai} /> <span className="muted">Tồn chỉ thay đổi sau khi quản lý duyệt điều chỉnh.</span></p>
          <table className="lines-table">
            <thead><tr><th>Hàng hóa</th><th>SL sổ sách</th><th>SL thực tế</th><th>Chênh lệch</th></tr></thead>
            <tbody>
              {(detail.chi_tiet || []).map((l, i) => (
                <tr key={l.ma_hang_hoa}>
                  <td>{l.ma_hang_hoa}</td>
                  <td>{l.sl_so_sach}</td>
                  <td>
                    {detail.trang_thai === 'cho_duyet'
                      ? <input type="number" min="0" value={l.sl_thuc_te} onChange={(e) => setActual(i, e.target.value)} />
                      : l.sl_thuc_te}
                  </td>
                  <td>{l.chenh_lech}</td>
                </tr>
              ))}
            </tbody>
          </table>
          {err && <div className="error">{err}</div>}
          <div className="actions">
            {detail.trang_thai === 'cho_duyet' && (
              <button className="secondary" onClick={saveActual}>Lưu SL thực tế</button>
            )}
            {isManager && detail.trang_thai === 'cho_duyet' && (
              <>
                <button className="danger" onClick={() => decide('reject')}>Yêu cầu kiểm lại</button>
                <button className="success" onClick={() => decide('approve')}>Duyệt điều chỉnh</button>
              </>
            )}
            <button className="secondary" onClick={() => setDetail(null)}>Đóng</button>
          </div>
        </Modal>
      )}
    </div>
  );
}
