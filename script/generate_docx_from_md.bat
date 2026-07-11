@echo off
REM Sinh file DOCX tu bao cao Markdown bang md-to-docx (can Node.js / npx).
REM Chay tu bat ky thu muc nao; script tu tro ve thu muc v2.

cd /d "%~dp0.."

echo Dang chuyen BaoCao_QuanLyKhoHang.md sang DOCX...
call npx -y -p @luytbq43/md-to-docx md-to-docx BaoCao_QuanLyKhoHang.md -o BaoCao_QuanLyKhoHang.docx

if %errorlevel% neq 0 (
    echo LOI: chuyen doi that bai. Kiem tra Node.js da cai va cac canh bao o tren.
    exit /b 1
)

echo Xong: %cd%\BaoCao_QuanLyKhoHang.docx
