# Service Pattern Rules

## Nguyên tắc bắt buộc
- Service chứa toàn bộ business logic: gọi repository, validate nghiệp vụ, quyết định khi nào cần transaction, trả lỗi domain.
- Service KHÔNG được import `net/http`, không nhận `*gin.Context` (hay context tương đương của framework), không biết gì về cookie/header/status code.
- Service nhận Input struct riêng (không dùng thẳng entity từ HTTP request) — tránh trường hợp client gửi field không mong muốn (vd: `role: "admin"`) ghi đè business rule.
- 1 Service ứng với 1 domain — không tạo "GodService" xử lý nhiều domain không liên quan.
- Use case cần phối hợp nhiều domain (vd: checkout cần trừ ví + trừ kho) → tạo Application Service/Orchestrator riêng ở tầng cao hơn, inject các domain service vào. Domain service không được gọi chéo lẫn nhau hay gọi ngược lên orchestrator.
- Service trả lỗi domain (sentinel error hoặc `AppError`), không trả HTTP status code — việc map sang status code là của handler.

## Nguyên tắc dùng chung cho Web + Mobile — bắt buộc
- Cùng 1 Service method phải phục vụ được cả handler web và handler api. Không viết 2 bản logic nghiệp vụ riêng cho 2 platform.
- Sự khác biệt duy nhất giữa web và mobile nằm ở handler (web set cookie/render HTML, api trả JSON) — không đẩy sự khác biệt này xuống service.

## Checklist khi generate code service
- [ ] Không import package framework HTTP nào.
- [ ] Nhận Input struct riêng, không bind thẳng entity từ request.
- [ ] Business rule (default value, điều kiện hợp lệ) nằm ở đây, không ở handler/repository.
- [ ] Trả lỗi domain, không trả status code.
- [ ] Use case đa domain có Orchestrator riêng, domain service không gọi chéo nhau.
- [ ] Method này gọi được từ cả handler web lẫn handler api mà không cần sửa.