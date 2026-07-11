import { useEffect, useState } from 'react';
import { api } from '../api/client';

// UC34: báo cáo nhập - xuất - tồn theo kỳ / kho / mặt hàng (chỉ QL kho).
export default function ReportPage() {
  const firstDay = new Date();
  firstDay.setDate(1);
  const [from, setFrom] = useState(firstDay.toISOString().slice(0, 10));
  const [to, setTo] = useState(new Date().toISOString().slice(0, 10));
  const [maKho, setMaKho] = useState('');
  const [warehouses, setWarehouses] = useState([]);
  const [rows, setRows] = useState([]);
  const [err, setErr] = useState('');

  useEffect(() => {
    api('/warehouses/options?page_size=200').then((d) => setWarehouses(d.items || [])).catch(() => {});
  }, []);

  async function run() {
    setErr('');
    try {
      const data = await api(`/reports/nxt?from=${from}&to=${to}&ma_kho=${maKho}`);
      setRows(data.items || []);
    } catch (e) { setErr(e.message); }
  }

  function exportCSV() {
    const header = 'Ma hang,Ten hang,Ton dau,Nhap,Xuat,Dieu chinh,Ton cuoi';
    const lines = rows.map((r) =>
      [r.ma_hang_hoa, `"${r.ten_hang_hoa}"`, r.ton_dau, r.nhap, r.xuat, r.dieu_chinh, r.ton_cuoi].join(','));
    const blob = new Blob(['﻿' + [header, ...lines].join('\n')], { type: 'text/csv;charset=utf-8' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = `bao-cao-nxt_${from}_${to}.csv`;
    a.click();
  }

  return (
    <div>
      <h2>Báo cáo nhập - xuất - tồn</h2>
      <div className="toolbar">
        <label>Từ</label>
        <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} />
        <label>đến</label>
        <input type="date" value={to} onChange={(e) => setTo(e.target.value)} />
        <select value={maKho} onChange={(e) => setMaKho(e.target.value)}>
          <option value="">Tất cả kho</option>
          {warehouses.map((k) => <option key={k.ma_kho} value={k.ma_kho}>{k.ten_kho}</option>)}
        </select>
        <button onClick={run}>Xem báo cáo</button>
        <button className="secondary" onClick={exportCSV} disabled={rows.length === 0}>Xuất CSV</button>
      </div>
      {err && <div className="error">{err}</div>}
      <table>
        <thead>
          <tr>
            <th>Mã hàng</th><th>Tên hàng</th><th>Tồn đầu</th><th>Nhập</th>
            <th>Xuất</th><th>Điều chỉnh KK</th><th>Tồn cuối</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((r) => (
            <tr key={r.ma_hang_hoa}>
              <td>{r.ma_hang_hoa}</td><td>{r.ten_hang_hoa}</td><td>{r.ton_dau}</td>
              <td>{r.nhap}</td><td>{r.xuat}</td><td>{r.dieu_chinh}</td><td><b>{r.ton_cuoi}</b></td>
            </tr>
          ))}
          {rows.length === 0 && <tr><td colSpan={7} className="muted">Chọn kỳ báo cáo rồi bấm "Xem báo cáo"</td></tr>}
        </tbody>
      </table>
    </div>
  );
}
