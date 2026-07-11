import { useCallback, useEffect, useState } from 'react';
import { api } from '../api/client';
import { useAuth } from '../auth/AuthContext';
import Modal from './Modal';

export function StatusBadge({ value }) {
  const label = { cho_duyet: 'Chờ duyệt', da_duyet: 'Đã duyệt', tu_choi: 'Từ chối' }[value] || value;
  return <span className={`badge ${value}`}>{label}</span>;
}

/**
 * Trang phiếu dùng chung cho phiếu nhập (UC23-27) và phiếu xuất (UC28-32):
 * danh sách + tìm kiếm + lập / sửa / xóa + duyệt / từ chối (chỉ QL kho).
 * kind = 'import' | 'export'
 */
export default function TicketPage({ kind }) {
  const isImport = kind === 'import';
  const endpoint = isImport ? '/imports' : '/exports';
  const idKey = isImport ? 'ma_phieu_nhap' : 'ma_phieu_xuat';
  const title = isImport ? 'Phiếu nhập kho' : 'Phiếu xuất kho';

  const { user } = useAuth();
  const isManager = user?.vai_tro === 'QLKho';

  const [items, setItems] = useState([]);
  const [q, setQ] = useState('');
  const [status, setStatus] = useState('');
  const [detail, setDetail] = useState(null);
  const [editing, setEditing] = useState(null); // phiếu đang lập / sửa
  const [msg, setMsg] = useState('');
  const [err, setErr] = useState('');

  const [warehouses, setWarehouses] = useState([]);
  const [suppliers, setSuppliers] = useState([]);
  const [products, setProducts] = useState([]);

  const load = useCallback(async () => {
    try {
      const data = await api(`${endpoint}?q=${encodeURIComponent(q)}&trang_thai=${status}&page_size=100`);
      setItems(data.items || []);
    } catch (e) {
      setErr(e.message);
    }
  }, [endpoint, q, status]);

  useEffect(() => { load(); }, [load]);

  useEffect(() => {
    api('/warehouses/options?page_size=200').then((d) => setWarehouses(d.items || [])).catch(() => {});
    api('/products?page_size=200').then((d) => setProducts(d.items || [])).catch(() => {});
    if (isImport) api('/suppliers?page_size=200').then((d) => setSuppliers(d.items || [])).catch(() => {});
  }, [isImport]);

  function emptyForm() {
    return {
      ma_ncc: '', ma_kho: '',
      ngay: new Date().toISOString().slice(0, 10),
      ly_do_xuat: '', ghi_chu: '',
      chi_tiet: [{ ma_hang_hoa: '', so_luong: 1, don_gia: 0 }],
    };
  }

  async function openEdit(row) {
    try {
      const p = await api(`${endpoint}/${row[idKey]}`);
      setEditing({
        id: row[idKey],
        ma_ncc: p.ma_ncc || '', ma_kho: p.ma_kho,
        ngay: (isImport ? p.ngay_nhap : p.ngay_xuat || '').slice(0, 10),
        ly_do_xuat: p.ly_do_xuat || '', ghi_chu: p.ghi_chu || '',
        chi_tiet: (p.chi_tiet || []).map((l) => ({
          ma_hang_hoa: l.ma_hang_hoa, so_luong: l.so_luong, don_gia: l.don_gia,
        })),
      });
      setErr('');
    } catch (e) { setErr(e.message); }
  }

  async function save() {
    setErr('');
    const f = editing;
    const body = isImport
      ? { ma_ncc: f.ma_ncc, ma_kho: f.ma_kho, ngay_nhap: f.ngay, chi_tiet: [] }
      : { ma_kho: f.ma_kho, ngay_xuat: f.ngay, ly_do_xuat: f.ly_do_xuat, ghi_chu: f.ghi_chu, chi_tiet: [] };
    body.chi_tiet = f.chi_tiet
      .filter((l) => l.ma_hang_hoa)
      .map((l) => ({ ma_hang_hoa: l.ma_hang_hoa, so_luong: Number(l.so_luong), don_gia: Number(l.don_gia) }));
    if (body.chi_tiet.length === 0) { setErr('Cần ít nhất một dòng hàng'); return; }
    try {
      if (f.id) await api(`${endpoint}/${f.id}`, { method: 'PUT', body });
      else await api(endpoint, { method: 'POST', body });
      setEditing(null);
      setMsg(f.id ? 'Đã cập nhật phiếu' : 'Đã lập phiếu, trạng thái Chờ duyệt');
      load();
    } catch (e) { setErr(e.message); }
  }

  async function remove(row) {
    if (!window.confirm(`Xóa phiếu ${row[idKey]}?`)) return;
    try {
      await api(`${endpoint}/${row[idKey]}`, { method: 'DELETE' });
      setMsg('Đã xóa phiếu');
      load();
    } catch (e) { setErr(e.message); }
  }

  async function decide(row, action) {
    const lyDo = action === 'reject' ? window.prompt('Lý do từ chối?') || '' : '';
    try {
      await api(`${endpoint}/${row[idKey]}/${action}`, { method: 'POST', body: { ly_do: lyDo } });
      setMsg(action === 'approve' ? 'Đã duyệt - tồn kho đã cập nhật' : 'Đã từ chối - tồn kho không đổi');
      setDetail(null);
      load();
    } catch (e) { setErr(e.message); }
  }

  async function openDetail(row) {
    try { setDetail(await api(`${endpoint}/${row[idKey]}`)); } catch (e) { setErr(e.message); }
  }

  const productName = (ma) => products.find((p) => p.ma_hang_hoa === ma)?.ten_hang_hoa || ma;

  return (
    <div>
      <h2>{title}</h2>
      <div className="toolbar">
        <input placeholder="Mã phiếu..." value={q} onChange={(e) => setQ(e.target.value)} />
        <select value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">Tất cả trạng thái</option>
          <option value="cho_duyet">Chờ duyệt</option>
          <option value="da_duyet">Đã duyệt</option>
          <option value="tu_choi">Từ chối</option>
        </select>
        <button className="secondary" onClick={load}>Tìm kiếm</button>
        <button onClick={() => { setEditing(emptyForm()); setErr(''); }}>Lập phiếu</button>
      </div>
      {msg && <div className="info">{msg}</div>}
      {err && !editing && <div className="error">{err}</div>}
      <table>
        <thead>
          <tr>
            <th>Mã phiếu</th>{isImport && <th>NCC</th>}<th>Kho</th><th>Ngày</th>
            <th>Tổng tiền</th><th>Trạng thái</th><th style={{ width: 220 }}></th>
          </tr>
        </thead>
        <tbody>
          {items.map((row) => (
            <tr key={row[idKey]}>
              <td>{row[idKey]}</td>
              {isImport && <td>{row.ma_ncc}</td>}
              <td>{row.ma_kho}</td>
              <td>{(isImport ? row.ngay_nhap : row.ngay_xuat || '').slice(0, 10)}</td>
              <td>{Number(row.tong_tien).toLocaleString('vi-VN')}</td>
              <td><StatusBadge value={row.trang_thai} /></td>
              <td>
                <button className="secondary" onClick={() => openDetail(row)}>Xem</button>{' '}
                {row.trang_thai === 'cho_duyet' && (
                  <>
                    <button className="secondary" onClick={() => openEdit(row)}>Sửa</button>{' '}
                    <button className="danger" onClick={() => remove(row)}>Xóa</button>
                  </>
                )}
              </td>
            </tr>
          ))}
          {items.length === 0 && <tr><td colSpan={7} className="muted">Không có phiếu nào</td></tr>}
        </tbody>
      </table>

      {detail && (
        <Modal title={`Chi tiết phiếu ${detail[idKey]}`} onClose={() => setDetail(null)}>
          <p>
            {isImport && <>NCC: <b>{detail.ma_ncc}</b> - </>}
            Kho: <b>{detail.ma_kho}</b> - Người lập: <b>{detail.ma_nguoi_dung}</b> -{' '}
            <StatusBadge value={detail.trang_thai} />
          </p>
          {!isImport && <p className="muted">Lý do xuất: {detail.ly_do_xuat} {detail.ghi_chu && `- ${detail.ghi_chu}`}</p>}
          <table>
            <thead><tr><th>Hàng hóa</th><th>SL</th><th>Đơn giá</th><th>Thành tiền</th></tr></thead>
            <tbody>
              {(detail.chi_tiet || []).map((l) => (
                <tr key={l.ma_hang_hoa}>
                  <td>{productName(l.ma_hang_hoa)}</td><td>{l.so_luong}</td>
                  <td>{Number(l.don_gia).toLocaleString('vi-VN')}</td>
                  <td>{Number(l.thanh_tien).toLocaleString('vi-VN')}</td>
                </tr>
              ))}
            </tbody>
          </table>
          {err && <div className="error">{err}</div>}
          <div className="actions">
            {isManager && detail.trang_thai === 'cho_duyet' && (
              <>
                <button className="danger" onClick={() => decide(detail, 'reject')}>Từ chối</button>
                <button className="success" onClick={() => decide(detail, 'approve')}>Duyệt phiếu</button>
              </>
            )}
            <button className="secondary" onClick={() => setDetail(null)}>Đóng</button>
          </div>
        </Modal>
      )}

      {editing && (
        <Modal title={editing.id ? `Sửa phiếu ${editing.id}` : `Lập ${title.toLowerCase()}`} onClose={() => setEditing(null)}>
          <div className="form-grid">
            {isImport && (
              <>
                <label>Nhà cung cấp *</label>
                <select value={editing.ma_ncc} onChange={(e) => setEditing({ ...editing, ma_ncc: e.target.value })}>
                  <option value="">-- chọn --</option>
                  {suppliers.filter((s) => s.trang_thai).map((s) => (
                    <option key={s.ma_ncc} value={s.ma_ncc}>{s.ten_ncc}</option>
                  ))}
                </select>
              </>
            )}
            <label>Kho *</label>
            <select value={editing.ma_kho} onChange={(e) => setEditing({ ...editing, ma_kho: e.target.value })}>
              <option value="">-- chọn --</option>
              {warehouses.filter((k) => k.trang_thai).map((k) => (
                <option key={k.ma_kho} value={k.ma_kho}>{k.ten_kho}</option>
              ))}
            </select>
            <label>Ngày</label>
            <input type="date" value={editing.ngay} onChange={(e) => setEditing({ ...editing, ngay: e.target.value })} />
            {!isImport && (
              <>
                <label>Lý do xuất *</label>
                <input value={editing.ly_do_xuat} onChange={(e) => setEditing({ ...editing, ly_do_xuat: e.target.value })} />
                <label>Ghi chú người nhận</label>
                <input value={editing.ghi_chu} onChange={(e) => setEditing({ ...editing, ghi_chu: e.target.value })} />
              </>
            )}
          </div>

          <h4>Dòng hàng</h4>
          <table className="lines-table">
            <thead><tr><th>Hàng hóa</th><th>SL</th><th>Đơn giá</th><th></th></tr></thead>
            <tbody>
              {editing.chi_tiet.map((l, i) => (
                <tr key={i}>
                  <td>
                    <select
                      value={l.ma_hang_hoa}
                      onChange={(e) => {
                        const ct = [...editing.chi_tiet];
                        const p = products.find((x) => x.ma_hang_hoa === e.target.value);
                        ct[i] = {
                          ...ct[i],
                          ma_hang_hoa: e.target.value,
                          don_gia: p ? (isImport ? p.gia_nhap : p.gia_xuat) : 0,
                        };
                        setEditing({ ...editing, chi_tiet: ct });
                      }}
                    >
                      <option value="">-- chọn hàng --</option>
                      {products.filter((p) => p.trang_thai).map((p) => (
                        <option key={p.ma_hang_hoa} value={p.ma_hang_hoa}>
                          {p.ten_hang_hoa}{!isImport ? ` (tồn ${p.tong_ton})` : ''}
                        </option>
                      ))}
                    </select>
                  </td>
                  <td>
                    <input type="number" min="1" value={l.so_luong}
                      onChange={(e) => {
                        const ct = [...editing.chi_tiet];
                        ct[i] = { ...ct[i], so_luong: e.target.value };
                        setEditing({ ...editing, chi_tiet: ct });
                      }} />
                  </td>
                  <td>
                    <input type="number" min="0" value={l.don_gia}
                      onChange={(e) => {
                        const ct = [...editing.chi_tiet];
                        ct[i] = { ...ct[i], don_gia: e.target.value };
                        setEditing({ ...editing, chi_tiet: ct });
                      }} />
                  </td>
                  <td>
                    <button className="danger" onClick={() => setEditing({
                      ...editing, chi_tiet: editing.chi_tiet.filter((_, j) => j !== i),
                    })}>x</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <button className="secondary" style={{ marginTop: 8 }}
            onClick={() => setEditing({ ...editing, chi_tiet: [...editing.chi_tiet, { ma_hang_hoa: '', so_luong: 1, don_gia: 0 }] })}>
            + Thêm dòng
          </button>
          {err && <div className="error">{err}</div>}
          <div className="actions">
            <button className="secondary" onClick={() => setEditing(null)}>Hủy</button>
            <button onClick={save}>{editing.id ? 'Lưu' : 'Lập phiếu (gửi duyệt)'}</button>
          </div>
        </Modal>
      )}
    </div>
  );
}
