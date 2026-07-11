import { useCallback, useEffect, useState } from 'react';
import { api } from '../api/client';
import Modal from './Modal';

/**
 * Trang danh mục dùng chung cho Tài khoản / NCC / Hàng hóa / Kho:
 * xem danh sách (UC xem DS), tìm kiếm (UC tìm), thêm / sửa / xóa (modal form).
 *
 * config = {
 *   title, endpoint, idKey,
 *   columns: [{ key, label, render? }],
 *   fields:  [{ key, label, type: 'text'|'number'|'password'|'select'|'checkbox',
 *               options?: [{value,label}], required?, createOnly? }],
 * }
 */
export default function CatalogPage({ title, endpoint, idKey, columns, fields }) {
  const [items, setItems] = useState([]);
  const [total, setTotal] = useState(0);
  const [q, setQ] = useState('');
  const [editing, setEditing] = useState(null); // null | {} (thêm) | row (sửa)
  const [form, setForm] = useState({});
  const [msg, setMsg] = useState('');
  const [err, setErr] = useState('');

  const load = useCallback(async (keyword = '') => {
    try {
      const data = await api(`${endpoint}?q=${encodeURIComponent(keyword)}&page_size=100`);
      setItems(data.items || []);
      setTotal(data.total || 0);
    } catch (e) {
      setErr(e.message);
    }
  }, [endpoint]);

  useEffect(() => {
    load();
  }, [load]);

  function openCreate() {
    const init = {};
    fields.forEach((f) => {
      init[f.key] = f.type === 'checkbox' ? true : '';
    });
    setForm(init);
    setEditing({});
    setErr('');
  }

  function openEdit(row) {
    setForm({ ...row });
    setEditing(row);
    setErr('');
  }

  async function save() {
    setErr('');
    try {
      const body = { ...form };
      fields.forEach((f) => {
        if (f.type === 'number') body[f.key] = Number(body[f.key] || 0);
      });
      const isCreate = !editing[idKey];
      if (isCreate) {
        await api(endpoint, { method: 'POST', body });
      } else {
        await api(`${endpoint}/${editing[idKey]}`, { method: 'PUT', body });
      }
      setEditing(null);
      setMsg(isCreate ? 'Đã thêm mới' : 'Đã cập nhật');
      load(q);
    } catch (e) {
      setErr(e.message);
    }
  }

  async function remove(row) {
    if (!window.confirm(`Xóa bản ghi ${row[idKey]}?`)) return;
    try {
      const res = await api(`${endpoint}/${row[idKey]}`, { method: 'DELETE' });
      setMsg(res.message || 'Đã xóa');
      load(q);
    } catch (e) {
      setErr(e.message);
    }
  }

  return (
    <div>
      <h2>{title}</h2>
      <div className="toolbar">
        <input
          placeholder="Từ khóa tìm kiếm..."
          value={q}
          onChange={(e) => setQ(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && load(q)}
        />
        <button className="secondary" onClick={() => load(q)}>Tìm kiếm</button>
        <button onClick={openCreate}>Thêm mới</button>
        <span className="muted">{total} bản ghi</span>
      </div>
      {msg && <div className="info">{msg}</div>}
      {err && !editing && <div className="error">{err}</div>}
      <table>
        <thead>
          <tr>
            {columns.map((c) => <th key={c.key}>{c.label}</th>)}
            <th style={{ width: 130 }}></th>
          </tr>
        </thead>
        <tbody>
          {items.map((row) => (
            <tr key={row[idKey]}>
              {columns.map((c) => (
                <td key={c.key}>{c.render ? c.render(row) : String(row[c.key] ?? '')}</td>
              ))}
              <td>
                <button className="secondary" onClick={() => openEdit(row)}>Sửa</button>{' '}
                <button className="danger" onClick={() => remove(row)}>Xóa</button>
              </td>
            </tr>
          ))}
          {items.length === 0 && (
            <tr><td colSpan={columns.length + 1} className="muted">Chưa có dữ liệu - bấm "Thêm mới"</td></tr>
          )}
        </tbody>
      </table>

      {editing !== null && (
        <Modal title={editing[idKey] ? `Sửa ${editing[idKey]}` : 'Thêm mới'} onClose={() => setEditing(null)}>
          <div className="form-grid">
            {fields.map((f) => {
              if (f.createOnly && editing[idKey]) return null;
              return (
                <FieldRow
                  key={f.key}
                  field={f}
                  value={form[f.key]}
                  onChange={(v) => setForm({ ...form, [f.key]: v })}
                />
              );
            })}
          </div>
          {err && <div className="error">{err}</div>}
          <div className="actions">
            <button className="secondary" onClick={() => setEditing(null)}>Hủy</button>
            <button onClick={save}>Lưu</button>
          </div>
        </Modal>
      )}
    </div>
  );
}

function FieldRow({ field, value, onChange }) {
  const { key, label, type = 'text', options = [] } = field;
  return (
    <>
      <label htmlFor={key}>{label}{field.required ? ' *' : ''}</label>
      {type === 'select' ? (
        <select id={key} value={value ?? ''} onChange={(e) => onChange(e.target.value)}>
          <option value="">-- chọn --</option>
          {options.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
        </select>
      ) : type === 'checkbox' ? (
        <input id={key} type="checkbox" checked={!!value} onChange={(e) => onChange(e.target.checked)} />
      ) : (
        <input id={key} type={type} value={value ?? ''} onChange={(e) => onChange(e.target.value)} />
      )}
    </>
  );
}
