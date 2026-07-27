# Logging Rules

## Thư viện
`log/slog` (built-in Go 1.21+). Nếu project đã dùng `zap`/`zerolog`, áp dụng nguyên tắc tương tự trên thư viện đó.

## Nguyên tắc bắt buộc
- Luôn dùng structured logging (key-value) — cấm dùng `Printf`/nối chuỗi để log.
- Log level dùng đúng mục đích:
  - `Debug`: chi tiết kỹ thuật, tắt ở production.
  - `Info`: sự kiện nghiệp vụ (login, order created) — bật ở production, coi như audit trail.
  - `Warn`: bất thường nhưng tự phục hồi được.
  - `Error`: request thất bại, cần chú ý.
  - Panic/Fatal-level: chỉ dùng ở `main.go` lúc khởi động, không dùng trong luồng xử lý request.

## Request ID — bắt buộc cho cả web lẫn API
- Mỗi request (web và mobile API) phải có `request_id` duy nhất, gắn kèm mọi dòng log trong quá trình xử lý request đó, đọc từ header `X-Request-ID` nếu có, tự generate UUID nếu không.
- Logger gắn `request_id` được lưu vào `context.Context`, lấy ra bằng `logger.FromContext(ctx)` ở bất kỳ đâu — không truyền logger qua từng tham số hàm.
- Gắn thêm `client_type` (`web`/`mobile`) vào log liên quan auth để phân tích sự cố theo platform.

## Cấm log
- Mật khẩu (kể cả đã hash), access token/refresh token đầy đủ, số thẻ ngân hàng, OTP, câu SQL kèm giá trị nhạy cảm.
- Nếu cần log token để debug: chỉ log vài ký tự đầu/cuối, không log full.

## Log ở đâu — bắt buộc
- Repository: không log, chỉ trả lỗi.
- Service: log `Info` cho sự kiện nghiệp vụ quan trọng; không log `Error` nếu lỗi sẽ được log lại ở handler.
- Handler: log `Error` khi trả lỗi cho client, kèm đủ ngữ cảnh (user_id, request_id, lỗi gốc).
- Middleware: log tổng quan mỗi request (method, path, status, duration).

## Output format
- Production: JSON.
- Local dev: text handler, level Debug.
- Chọn theo biến môi trường `ENV`.

## Checklist khi generate code có logging
- [ ] Dùng key-value structured logging.
- [ ] Log trong cùng 1 request có chung `request_id`.
- [ ] Không log password/token đầy đủ/thông tin thanh toán.
- [ ] Repository không log; log tập trung ở service (sự kiện) và handler/middleware (lỗi, lifecycle).
- [ ] Debug level không bật ở production.