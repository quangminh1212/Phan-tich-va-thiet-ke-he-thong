# Template đặc tả & Quy tắc đồng bộ biểu đồ

## 1. Template đặc tả Use Case (kịch bản)

Dùng cho mọi use case trọng yếu. Mỗi UC một bảng như sau:

| Mục | Nội dung |
|---|---|
| **Mã / Tên UC** | UC-03 / Mượn sách |
| **Tác nhân chính** | Thủ thư |
| **Tác nhân phụ** | (hệ thống thẻ, nếu có) |
| **Mô tả** | Thủ thư lập phiếu mượn cho độc giả |
| **Tiền điều kiện** | Thủ thư đã đăng nhập; độc giả có thẻ hợp lệ |
| **Kích hoạt (trigger)** | Độc giả yêu cầu mượn sách |
| **Luồng sự kiện chính** | 1. Thủ thư tra cứu độc giả. 2. Hệ thống hiển thị thông tin & số sách đang mượn. 3. Thủ thư quét mã sách. 4. Hệ thống kiểm tra tình trạng sách. 5. Hệ thống tạo phiếu mượn. 6. Hệ thống cập nhật tình trạng sách = "Đang mượn". |
| **Luồng phụ / ngoại lệ** | 4a. Sách đã được mượn → báo lỗi, dừng. 2a. Độc giả vượt hạn mức → từ chối. |
| **Hậu điều kiện** | Phiếu mượn được lưu; tồn kho sách giảm |
| **Yêu cầu liên quan** | FR-05 |

> **Quy tắc**: số bước trong "Luồng chính" phải ánh xạ 1–1 với các message trong biểu đồ tuần tự của UC đó.

## 2. Quy tắc đồng bộ chéo (checklist khi rà soát)

- [ ] **Actor**: mọi tác nhân ở mô tả nghiệp vụ (Ch.1) đều có trong biểu đồ use case tổng quát, và ngược lại.
- [ ] **Use case ↔ Yêu cầu**: mọi yêu cầu chức năng (FR) được phủ bởi ≥1 use case; mọi use case truy về ≥1 FR.
- [ ] **Use case ↔ Sequence**: mọi UC trọng yếu có ≥1 biểu đồ tuần tự; tên UC khớp.
- [ ] **Sequence ↔ Class**: mọi đối tượng (lifeline) trong sequence là một lớp trong biểu đồ lớp; mọi message gọi tới một đối tượng phải tương ứng một **phương thức** của lớp đó (xuất hiện trong biểu đồ lớp thiết kế chi tiết Ch.5).
- [ ] **Lifeline phải là TÊN LỚP CỤ THỂ** — `DieuKhienDatVe`, `BenhNhanService`, `HoaDon`... KHÔNG đặt tên vai trò chung chung như "Điều khiển", "Giao diện", "Hệ thống". Đây là nguyên nhân gốc khiến message không khớp lớp nào.
- [ ] **Thông báo/kết quả trả về dùng return `-->>` (nét đứt)**, KHÔNG viết `->>` kèm `()` như gọi phương thức. Vd đúng: `DK-->>GD: thongBaoThanhCong` (nhãn, không ngoặc); SAI: `DK->>GD: thongBaoThanhCong()`.
- [ ] **Class ↔ CSDL**: mọi lớp thực thể (entity) ánh xạ tới ≥1 bảng; thuộc tính khớp cột; quan hệ + bội số khớp khóa ngoại.
- [ ] **Tên nhất quán**: 1 khái niệm = 1 tên (vd luôn "DocGia", không lẫn "BanDoc"/"KhachHang").
- [ ] **Bội số (multiplicity)** hợp lý và nhất quán giữa class diagram và lược đồ CSDL.
- [ ] **include/extend** dùng đúng nghĩa; không vẽ association giữa 2 use case.

## 3. Mermaid templates

> **Ba mức biểu đồ lớp** (theo cấu trúc 6 chương):
> - **Mô hình lĩnh vực / Domain (Ch3.1)** — chỉ lớp thực thể + quan hệ, KHÔNG thuộc tính chi tiết/phương thức/lớp biên-điều khiển.
> - **Lớp phân tích / Analysis (Ch3.2)** — thêm lớp Biên–Điều khiển–Thực thể (BCE) tham gia UC tiêu biểu.
> - **Lớp thiết kế chi tiết (Ch5.1)** — đầy đủ thuộc tính + kiểu, phương thức, tầm vực (+/-/#). Đây là bản `check-sync.py` đối chiếu với sequence.

### Use case (mermaid không có UC diagram gốc → dùng graph)
```mermaid
graph LR
  actor1((Thủ thư))
  actor2((Độc giả))
  subgraph "Hệ thống Quản lý Thư viện"
    uc1([Đăng nhập])
    uc2([Tra cứu sách])
    uc3([Mượn sách])
    uc4([Trả sách])
  end
  actor1 --- uc1
  actor1 --- uc3
  actor1 --- uc4
  actor2 --- uc2
  uc3 -. include .-> uc1
```
> Nếu xuất sang draw.io để đẹp hơn, dùng skill `mermaid-to-drawio`.

### Class diagram
```mermaid
classDiagram
  class DocGia {
    -maDocGia: String
    -hoTen: String
    +dangKy()
    +traCuu()
  }
  class PhieuMuon {
    -maPhieu: String
    -ngayMuon: Date
    -ngayHenTra: Date
    +taoPhieu()
  }
  class Sach {
    -maSach: String
    -tenSach: String
    -tinhTrang: String
  }
  DocGia "1" --> "0..*" PhieuMuon : thực hiện
  PhieuMuon "1" --> "1..*" Sach : gồm
```

### Sequence diagram (luồng chính, message khớp phương thức)
```mermaid
sequenceDiagram
  actor TT as Thủ thư
  participant GD as GiaoDienMuon
  participant DK as DieuKhienMuon
  participant PM as PhieuMuon
  participant S as Sach
  TT->>GD: chonMuonSach(maDG, maSach)
  GD->>DK: muonSach(maDG, maSach)
  DK->>S: kiemTraTinhTrang(maSach)
  S-->>DK: tinhTrang
  DK->>PM: taoPhieu(maDG, maSach)
  PM-->>DK: maPhieu
  DK->>S: capNhatTinhTrang("Đang mượn")
  DK-->>GD: thongBaoThanhCong
```

### State machine (vòng đời đối tượng)
```mermaid
stateDiagram-v2
  [*] --> SanSang
  SanSang --> DangMuon : muonSach()
  DangMuon --> SanSang : traSach()
  DangMuon --> QuaHan : quaNgayHenTra
  QuaHan --> SanSang : traSach()+phat()
```

### Activity diagram
```mermaid
flowchart TD
  A([Bắt đầu]) --> B[Tra cứu độc giả]
  B --> C{Độc giả hợp lệ?}
  C -- Không --> Z([Kết thúc - từ chối])
  C -- Có --> D[Quét mã sách]
  D --> E{Sách sẵn sàng?}
  E -- Không --> Z
  E -- Có --> F[Tạo phiếu mượn]
  F --> G[Cập nhật tình trạng sách]
  G --> H([Kết thúc])
```

### Package diagram (Ch4.2 — phân chia module / hệ thống con, kiến trúc BCE / 3 lớp)
```mermaid
flowchart TB
  subgraph Boundary["«layer» Giao diện"]
    GD[GiaoDienMuon]
  end
  subgraph Control["«layer» Điều khiển"]
    DK[DieuKhienMuon]
  end
  subgraph Entity["«layer» Thực thể"]
    E[DocGia, Sach, PhieuMuon]
  end
  Boundary --> Control --> Entity
```

### Component diagram (Ch4.3 — thành phần vật lý của phần mềm)
```mermaid
flowchart LR
  UI["«component» Web UI"] --> API["«component» REST API"]
  API --> SVC["«component» Service Layer"]
  SVC --> DAO["«component» Data Access"]
  DAO --> DB[("«database» CSDL")]
```

### Deployment diagram (Ch4.4 — bố trí trên phần cứng & mạng)
```mermaid
flowchart TB
  subgraph Client["«device» Máy người dùng"]
    Browser["Trình duyệt"]
  end
  subgraph AppServer["«server» App Server"]
    App["«artifact» app.war"]
  end
  subgraph DBServer["«server» DB Server"]
    DBMS["«artifact» MySQL"]
  end
  Browser -->|HTTPS| App
  App -->|JDBC| DBMS
```
