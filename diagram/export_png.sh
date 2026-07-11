#!/usr/bin/env bash
# Xuất PNG từ drawio rồi CROP sát phần hình vẽ (bỏ khoảng trắng thừa, chừa viền 10px).
# Cần: drawio CLI (brew install --cask drawio) + ImageMagick (brew install imagemagick).
#
# Dùng:
#   ./export_png.sh              # xuất toàn bộ drawio/*.drawio
#   ./export_png.sh 5.34 4.1     # chỉ xuất file có tên bắt đầu bằng các tiền tố này
set -euo pipefail
cd "$(dirname "$0")"

MAGICK=$(command -v magick || command -v convert)

for f in drawio/*.drawio; do
  b=$(basename "${f%.drawio}")
  if [ $# -gt 0 ]; then
    match=0
    for p in "$@"; do [[ "$b" == "$p"* ]] && match=1; done
    [ $match -eq 1 ] || continue
  fi
  # -b 40: nới canvas quanh hình để chữ tràn ra ngoài khung (tên actor...) không bị cắt,
  # sau đó -trim cắt sát phần có nét vẽ thật, chừa viền 10px.
  # QUAN TRỌNG: ép nền TRẮNG ĐỤC và bỏ kênh alpha - pixel trong suốt khi bị Word/PDF
  # resample sẽ tạo viền xám mờ quanh ảnh.
  drawio -x -f png -s 2 -b 40 -o "png/$b.png" "$f" >/dev/null 2>&1
  "$MAGICK" "png/$b.png" -background white -alpha remove -alpha off \
    -fuzz 1% -trim +repage -bordercolor white -border 10 "png/$b.png"
  echo "exported + cropped: png/$b.png"
done
