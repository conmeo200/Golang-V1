# Transaction Rules

## Nguyên tắc bắt buộc
- Transaction chỉ được mở ở tầng **Service**. Không mở ở Repository, không mở ở Handler.
- Repository không tự quyết định commit/rollback — chỉ thực thi query bằng connection/transaction được truyền vào qua `context.Context`.
- Dùng pattern Unit of Work: `TxManager.WithTx(ctx, fn)` — mở transaction, chạy `fn`, tự commit nếu `fn` trả `nil`, tự rollback nếu `fn` trả lỗi. Service không tự gọi `Commit()`/`Rollback()` thủ công.
- Repository lấy executor (DB hoặc Tx) từ context — không có tham số `tx` riêng truyền qua từng method.

## Khi nào BẮT BUỘC dùng transaction
- Nhiều thao tác ghi trên nhiều bảng phải cùng thành công hoặc cùng thất bại.
- Read-then-write cần tránh race condition (dùng thêm `SELECT ... FOR UPDATE` hoặc version column).
- Thao tác thay đổi role/permission kèm ghi audit log — phải atomic.

## Khi KHÔNG cần transaction
- Chỉ 1 câu lệnh ghi duy nhất.
- Thao tác đọc thuần (trừ khi cần isolation đặc biệt).

## Isolation & khoá
- Mặc định `READ COMMITTED`.
- Ưu tiên `SELECT ... FOR UPDATE` để lock row thay vì nâng isolation level toàn transaction.
- Transaction phải ngắn nhất có thể — cấm gọi HTTP/service bên thứ 3 bên trong transaction block.

## Timeout — bắt buộc
- Mọi transaction phải được bọc bởi context có timeout hợp lý (5–10s tuỳ độ phức tạp).

## Checklist khi generate code có transaction
- [ ] Transaction mở ở Service, không ở Repository/Handler.
- [ ] Dùng `TxManager.WithTx`, không tự quản lý `Begin/Commit/Rollback` rải rác.
- [ ] Không có lời gọi HTTP/gRPC bên trong transaction.
- [ ] Context truyền vào có timeout.
- [ ] Thao tác update số dư/tồn kho có `FOR UPDATE` hoặc version column.