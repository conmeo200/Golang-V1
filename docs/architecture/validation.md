# Validation Rules

## Phân tầng validation — bắt buộc
1. **Format/structural validation** (tầng Handler): đúng kiểu, đúng định dạng, field bắt buộc — dùng struct tag (`github.com/go-playground/validator/v10`). Không cần biết business rule.
2. **Business validation** (tầng Service): hợp lệ theo nghiệp vụ (email đã tồn tại, đủ quyền, đủ số dư) — cần query DB hoặc biết context nghiệp vụ.

Không trộn 2 loại này: Handler không tự query DB để check trùng; Service không check định dạng field.

## Dùng chung cho Web form và API JSON — bắt buộc
- Cùng 1 struct request + cùng tag validator dùng cho cả `ShouldBindJSON` (API) và `ShouldBind` (web form). Không định nghĩa 2 struct validate riêng cho 2 platform.
- Có 1 hàm dùng chung để format lỗi validator thành `map[string]string` theo field — handler web render vào template, handler api trả JSON `details`. Không viết 2 bộ format lỗi riêng.

## Custom validation rule
- Rule nghiệp vụ đặc thù (độ mạnh mật khẩu...) đăng ký qua `validator.RegisterValidation`, không viết if/else rời rạc trong handler.

## Bảo mật — bắt buộc
- Role/permission không bao giờ nhận trực tiếp từ input client — luôn lấy từ token đã verify hoặc từ DB.
- Action nhạy cảm (đổi mật khẩu, đổi email, xoá tài khoản) validate thêm: xác nhận mật khẩu hiện tại hoặc yêu cầu re-auth gần đây.

## Cấm
- Validate trùng lặp cùng 1 điều kiện ở cả handler và service.
- Trả lỗi validate chung chung không rõ field nào sai.

## Checklist khi generate code validation
- [ ] Format validation (struct tag) tách biệt khỏi business validation (Service).
- [ ] Cùng struct request + tag validator dùng cho cả web form và API JSON.
- [ ] Lỗi trả về kèm chi tiết theo field.
- [ ] Role/permission không nhận trực tiếp từ input client.
- [ ] Custom rule dùng `RegisterValidation`, không viết rời rạc trong handler.