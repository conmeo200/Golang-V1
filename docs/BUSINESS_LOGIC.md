# Bối cảnh & Quy tắc Nghiệp vụ Kinh doanh (Business Logic)

**[CRITICAL INSTRUCTION FOR ALL AI ASSISTANTS]**
Tất cả các AI khi phân tích, thiết kế Database, thiết kế API, hoặc xây dựng UI liên quan đến tài chính/doanh thu cho dự án này **BẮT BUỘC** phải tuân thủ các quy tắc trong tài liệu này.

## 1. Doanh thu & Dòng tiền (Revenue & Cashflow)

Trong dự án thương mại điện tử này,    doanh thu của người bán (Seller) không phải là một con số đơn giản. Phải phân biệt rõ 2 loại doanh thu: 
-     **Doanh thu Gộp (Gross Revenue): ** Tổng giá trị đơn hàng mà khách hàng đã thanh toán (bao gồm cả phí ship, trừ khi có cấu hình khác).
-     **Doanh thu Thuần (Net Revenue): ** Số tiền thực nhận của Seller sau khi đã trừ đi tất cả các loại phí và thuế (Xem phần 2).

> **Lưu ý về phí ship: ** Phí ship thu từ khách thường là tiền thu hộ để trả cho đơn vị vận chuyển (pass-through), không phải doanh thu thực của Seller. Nếu tính phí ship vào `Gross Revenue`, thì khi tính `Net Revenue` cũng phải trừ ra khoản "chi phí ship trả hộ" tương ứng — nếu không, doanh thu của Seller sẽ bị đội ảo lên.

### Cấu trúc 4 Ví (Wallets)
Bất kỳ cấu trúc dữ liệu nào liên quan đến số dư của người bán phải chia làm 4 loại: 
1.  `Pending Balance` (Chờ đối soát)                                              : Tiền của các đơn hàng đang giao hoặc đang trong thời gian cho phép khiếu nại đổi trả. Số tiền này CHƯA thể rút.
2.  `On-hold / Frozen Balance` (Đang tranh chấp)                                  : Tiền của các đơn hàng đang bị khiếu nại/tranh chấp (dispute) đã mở nhưng chưa xử lý xong. Đây là trạng thái **tách riêng** khỏi `Pending` thông thường, vì nó không tự động chuyển sang `Available` theo thời gian mà phải chờ kết quả xử lý khiếu nại.
3.  `Available Balance` (Có thể rút)                                              : Tiền từ các đơn hàng đã hoàn thành (Khách đã ấn "Đã nhận được hàng", hết thời gian khiếu nại, hoặc khiếu nại đã xử lý xong theo hướng có lợi cho Seller).
4.  `Withdrawn Amount` (Đã rút)                                                   : Tổng số tiền Seller đã rút về ngân hàng thành công.

### Xử lý số dư âm (Negative Balance)
Cần có quy tắc rõ ràng cho trường hợp: đơn hàng đã chuyển tiền sang `Available` (hoặc Seller đã rút), sau đó khách khiếu nại và được hoàn tiền thành công: 
- Nếu `Available Balance` đủ để trừ ngược → trừ trực tiếp.
- Nếu không đủ (tiền đã nằm ở `Withdrawn`) → hệ thống phải cho phép `Available Balance` **âm**, và số âm này sẽ được cấn trừ dần vào các đơn hàng tiếp theo trước khi Seller được rút tiền mới.
- Phải có cảnh báo/khóa tính năng rút tiền khi số dư đang âm vượt ngưỡng cấu hình.

## 2. Các khoản Khấu trừ (Deductions / Fees)
Khi thiết kế API trả về Doanh thu hoặc thiết kế UI hiển thị hóa đơn đối soát, phải bao gồm các loại phí sau (các con số dưới đây là cấu hình chuẩn hiện tại, có thể lưu trong biến môi trường hoặc bảng `Config`): 
1.  **Phí thanh toán (Payment Gateway Fee)                                  : ** ~2% trên tổng giá trị thanh toán của khách.
2.  **Phí nền tảng/Phí sàn (Platform Commission)                            : ** ~5% (dao động tùy ngành hàng). Phí này thu trên giá trị sản phẩm.
3.  **Thuế thu nhập (Taxes - TNCN/GTGT)                                     : ** Sàn bắt buộc thu hộ 1.5% đối với các hộ kinh doanh trên sàn thương mại điện tử.
4.  **Chi phí Marketing (Nội bộ Sàn)                                        : ** Ví dụ như Freeship Extra,                                                                                                         Hoàn Xu Extra, hoặc chạy Quảng cáo (Ads).
5.  **Chiết khấu/Voucher của Shop                                           : ** Nếu khách dùng Voucher do Shop tạo,                                                                                               Shop chịu 100% chi phí này.

### Cơ sở tính phí (Fee Base) & Thứ tự áp dụng
**Quan trọng: ** Các loại phí trên **không cùng một cơ sở tính (base)**. Nếu áp dụng công thức cộng trừ đơn giản trên cùng một `Gross Revenue` sẽ ra kết quả sai. Phải tính tuần tự theo bảng sau: 

| Bước | Loại phí | Tính trên (Base) | Ghi chú |
|---|---|---|---|
| 1 | Chiết khấu/Voucher của Shop | Giá trị sản phẩm gốc | Trừ trước để ra "giá trị sản phẩm sau voucher" |
| 2 | Platform Commission (~5%) | Giá trị sản phẩm **sau** voucher Shop | Không tính trên phí ship |
| 3 | Payment Gateway Fee (~2%) | Tổng giá trị thanh toán thực tế của khách (gồm ship, nếu có) | Đây là phí cổng thanh toán nên tính trên dòng tiền thực chảy qua cổng |
| 4 | Thuế (TNCN/GTGT ~1.5%) | Gross Revenue (hoặc Gross sau voucher, tùy cấu hình pháp lý) — **cần xác nhận với kế toán** | Không được tự suy diễn base thuế, phải lấy từ cấu hình |
| 5 | Chi phí Marketing nội bộ (Ads, Freeship Extra, Hoàn Xu Extra) | Trừ trực tiếp bằng số tiền thực chi | Ghi nhận theo từng chiến dịch/khoảng thời gian |

**Công thức chuẩn (tuần tự, không phải cộng dồn 1 lần):**
```
Step1 = Giá trị sản phẩm gốc - Voucher Shop
Step2 = Step1 - Platform Commission (5% x Step1)
Step3 = Step2 + Phí ship (nếu Gross có bao gồm ship) - Payment Gateway Fee (2% x Tổng thanh toán)
Step4 = Step3 - Thuế (1.5% x Base_thuế_theo_cấu_hình)
Net Revenue = Step4 - Chi phí Marketing nội bộ
```
> AI khi code phần này **KHÔNG được gộp tất cả phí trừ 1 lần trên Gross Revenue**. Phải lưu từng bước tính (mỗi loại phí = 1 dòng trong bảng `Fee_Breakdown` hoặc `Transaction_Detail`) để phục vụ đối soát và audit.

### Hoàn tiền/Hủy đơn một phần (Partial Refund)
Khi đơn hàng bị hoàn **một phần** (ví dụ mua 3 món, trả lại 1 món):
- Phải tính lại `Platform Commission` và `Thuế` theo đúng tỷ lệ giá trị của (các) sản phẩm bị hoàn, không tính lại trên toàn bộ đơn.
- `Payment Gateway Fee` đã thu thường **không hoàn lại** (tùy chính sách cổng thanh toán) — cần cấu hình rõ có hoàn hay không.
- Phải tạo transaction điều chỉnh riêng (loại `Partial Refund Adjustment`), không sửa trực tiếp lên transaction gốc.

**Công thức chung (tổng quát, tham chiếu):**
`Net Revenue = Gross Revenue - (Payment Fee + Platform Commission + Taxes + Shop's Marketing Costs)`
*(Công thức này chỉ mang tính minh họa tổng quát — khi triển khai thực tế phải tuân theo thứ tự và base ở bảng trên.)*

## 3. Quy trình Rút tiền (Withdrawal Flow)
- Seller chỉ được rút từ `Available Balance` (không được rút từ `Pending` hoặc `On-hold`).
- Mỗi giao dịch rút tiền phải tạo ra một Transaction record với trạng thái: `Pending` -> `Processing` -> `Success`/`Failed`.
- Khi trạng thái là `Failed`: bắt buộc phải hoàn tiền lại vào `Available Balance`. Việc hoàn tiền này phải **idempotent** (có cơ chế kiểm tra transaction đã hoàn hay chưa) để tránh cộng tiền 2 lần nếu webhook/callback bị gọi lặp.
- Cần có các ràng buộc cấu hình (lưu trong bảng `Config`):
  - Số tiền rút tối thiểu / tối đa mỗi lần.
  - Phí rút tiền (nếu có) và ai chịu phí.
  - Giới hạn số lần rút trong ngày/tháng.
- Mọi giao dịch thay đổi số dư (Bán được hàng, Hoàn tiền, Rút tiền, Trừ phí Ads) đều phải được ghi log dưới dạng Báo cáo để phục vụ cho việc xuất file Excel Khai Thuế.

## 4. Nguyên tắc thiết kế dữ liệu (Data Integrity Principles)
- **Transaction là bất biến (immutable):** Không được `UPDATE` hay `DELETE` một transaction đã ghi nhận. Mọi điều chỉnh (hoàn tiền, sửa lỗi, khiếu nại) phải tạo ra một bản ghi điều chỉnh mới (correction/adjustment entry) tham chiếu ngược lại transaction gốc.
- **Mô hình sổ cái (Ledger-style):** Khuyến nghị mỗi thay đổi số dư là 1 dòng ghi có `balance_before` và `balance_after`, thay vì chỉ `UPDATE` một cột `balance` duy nhất — nhằm tránh race condition khi nhiều đơn hàng được xử lý đồng thời và đảm bảo có thể truy vết (audit trail) đầy đủ.
- **Quy tắc làm tròn (Rounding):** Vì các phần trăm phí sẽ ra số tiền lẻ, phải quy định rõ đơn vị làm tròn (ví dụ: làm tròn xuống đến đồng, hoặc đến hàng trăm VNĐ) và áp dụng nhất quán ở tất cả các bước tính phí.

##  Lời nhắc cho các hệ thống AI                                                        : 
Khi bạn được yêu cầu "tạo bảng database cho doanh thu" hoặc "viết API báo cáo doanh thu", bạn **KHÔNG ĐƯỢC** chỉ tạo một cột `Revenue`. Bạn phải tạo bảng `Transactions` (bất biến, ledger-style), các bảng cấu hình `Fees`/`Config`, bảng `Fee_Breakdown` chi tiết từng loại phí theo đúng base và thứ tự ở Mục 2, và xử lý đầy đủ các case: hoàn tiền một phần, tranh chấp/khiếu nại (`On-hold Balance`), và số dư âm.