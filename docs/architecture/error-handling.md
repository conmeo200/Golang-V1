# Error Handling Rules

## Nguyên tắc dòng chảy lỗi — bắt buộc
- Lỗi tạo ra ở tầng thấp (repository) → wrap thêm ngữ cảnh khi đi qua service (dùng `%w`) → map thành response ở handler (tầng cao nhất).
- Chỉ log lỗi 1 lần, ở handler hoặc middleware recover — không log lặp ở repository/service.
- Không bao giờ trả lỗi kỹ thuật thô (stack trace, câu SQL, đường dẫn file) ra response cho client.

## Định nghĩa lỗi — bắt buộc
- Lỗi domain định nghĩa dưới dạng sentinel error trong package `apperror` (`ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidCredentials`, `ErrPermissionDenied`, `ErrValidation`, `ErrTokenExpired`, `ErrTokenRevoked`...).
- Lỗi cần kèm metadata cho response (field lỗi, status code) dùng struct `AppError{Code, Message, Status, Details, Err}`, implement `Error()` và `Unwrap()`.
- Không so sánh lỗi bằng chuỗi (`err.Error() == "..."`) — luôn dùng `errors.Is`/`errors.As`.

## Map lỗi → response — bắt buộc dùng chung 1 hàm
- Có đúng 1 hàm `MapError(err error) *AppError` dùng chung cho cả handler web và handler api — không viết switch-case lỗi riêng cho từng platform.
- Handler web dùng kết quả `MapError` để render HTML lỗi; handler api dùng để trả JSON. Logic map lỗi → status code/message giống hệt nhau.

## Panic/Recover
- `panic` chỉ dùng cho lỗi lập trình không nên xảy ra (vi phạm invariant), không dùng cho lỗi nghiệp vụ.
- Middleware recover đặt ngoài cùng route chain, log đầy đủ (kèm stack trace), trả response 500 chung chung, không lộ chi tiết kỹ thuật cho client.

## Checklist khi generate code liên quan lỗi
- [ ] Lỗi domain là sentinel error hoặc AppError, không phải string rời rạc.
- [ ] Mọi wrap lỗi dùng `%w`.
- [ ] Chỉ log lỗi đúng 1 lần trong toàn luồng xử lý request.
- [ ] Response ra client không chứa thông tin kỹ thuật nhạy cảm.
- [ ] Dùng chung `MapError` cho cả web và api, không tạo bản riêng.