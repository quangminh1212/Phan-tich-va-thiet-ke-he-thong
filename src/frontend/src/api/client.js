// Fetch wrapper: gắn Bearer token, tự đăng xuất khi 401.
export function getToken() {
  return localStorage.getItem('token') || '';
}

export async function api(path, { method = 'GET', body } = {}) {
  const headers = { 'Content-Type': 'application/json' };
  const token = getToken();
  if (token) headers.Authorization = `Bearer ${token}`;
  const res = await fetch(`/api${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });
  if (res.status === 401) {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
    throw new Error('phiên hết hạn, đăng nhập lại');
  }
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || `lỗi ${res.status}`);
  return data;
}
