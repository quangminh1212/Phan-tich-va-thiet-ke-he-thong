#!/usr/bin/env bash
# Sinh file DOCX từ báo cáo Markdown bằng md-to-docx (cần Node.js / npx).
# Chạy từ bất kỳ thư mục nào; script tự trỏ về thư mục gốc dự án.
set -euo pipefail

cd "$(dirname "$0")/.."

echo "Đang chuyển BaoCao_QuanLyKhoHang.md sang DOCX..."
npx -y -p @luytbq43/md-to-docx md-to-docx BaoCao_QuanLyKhoHang.md -o BaoCao_QuanLyKhoHang.docx

echo "Xong: $(pwd)/BaoCao_QuanLyKhoHang.docx"
