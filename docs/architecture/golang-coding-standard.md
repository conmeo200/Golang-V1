# Golang Coding Standard Rules

## Công cụ bắt buộc chạy trước khi coi code là hoàn chỉnh
- `gofmt`/`goimports` để format.
- `golangci-lint` (gồm govet, errcheck, staticcheck, unused, gosimple) — code phải pass lint không lỗi.
- `go vet`.

## Naming — bắt buộc tuân theo
- Package: chữ thường, không gạch dưới, không viết hoa.
- Interface: không thêm tiền tố `I` (`UserRepository`, không phải `IUserRepository`). Interface 1 method có thể đặt tên theo hành vi + hậu tố `er`.
- Struct implement interface: hậu tố mô tả nguồn (`PostgresUserRepository`).
- Biến: `camelCase`. Hằng export: `PascalCase`. Không dùng `MAX_RETRY` kiểu hằng C — dùng `MaxRetry`.
- Viết tắt giữ nguyên hoa: `userID`, `HTTPClient` — không viết `userId`, `HttpClient`.
- Tên hàm không cần hậu tố kiểu `WithError` — Go convention là trả `(T, error)`.

## Error handling — bắt buộc
- Luôn check lỗi ngay sau khi gọi; không bỏ qua bằng `_` trừ khi có comment giải thích rõ lý do.
- Không dùng `panic` cho lỗi nghiệp vụ có thể lường trước (validation fail, not found). Chỉ dùng cho lỗi lập trình không thể phục hồi, luôn có `recover` ở middleware ngoài cùng.
- Wrap lỗi bằng `%w` để giữ được `errors.Is`/`errors.As` qua nhiều tầng.
- Chi tiết đầy đủ xem `error-handling.md`.

## Context — bắt buộc
- Hàm thực hiện I/O (DB, HTTP, Redis) nhận `context.Context` làm tham số đầu tiên.
- Không lưu `context.Context` vào struct field.
- Không dùng `context.Background()` trong handler/service — context phải truyền từ request gốc.
- Context chỉ dùng để truyền request-scoped value (request ID, user ID sau auth), không dùng để truyền business data.

## Interface & Dependency — bắt buộc
- Interface khai báo ở nơi **sử dụng** (domain/service), không khai báo ở nơi implement.
- Interface nhỏ, chỉ chứa method thực sự cần dùng — không tạo interface CRUD khổng lồ nếu không dùng hết.
- Dependency truyền qua constructor, không dùng biến global hay singleton ẩn.

## Concurrency
- Không share mutable state giữa goroutine mà không có mutex/channel bảo vệ.
- Mọi goroutine phải có cách dừng (context cancel/done channel) — không để leak.
- Dùng `errgroup.Group` khi chạy nhiều goroutine song song cần gom lỗi.

## Comment/Doc — bắt buộc
- Mọi identifier export đều có comment bắt đầu bằng chính tên đó (chuẩn godoc).
- Comment giải thích **tại sao**, không lặp lại **cái gì** code đã tự rõ.

## Struct & method
- Nhất quán pointer receiver hay value receiver cho cùng 1 type, không trộn lẫn.
- 1 file cho 1 struct/concept chính.

## Import — thứ tự bắt buộc
Standard library → third-party → internal package, cách nhau dòng trống (`goimports` tự xử lý).

## Cấm tuyệt đối
- Biến global cho DB connection, config, logger.
- Dùng `interface{}`/`any` khi đã biết rõ kiểu — ưu tiên generics.
- Trả `nil, nil` mập mờ khi hàm có kiểu trả `(*T, error)` — không tìm thấy phải trả lỗi rõ ràng (`ErrNotFound`), không để caller đoán.