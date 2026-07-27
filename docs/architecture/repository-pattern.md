# Repository Pattern Rules

## Nguyên tắc bắt buộc
- Interface repository khai báo trong package `domain/<entity>/repository.go`. Implementation đặt trong `repository/<postgres|redis|...>/`. Service import domain, không import trực tiếp implementation.
- Chỉ khai báo method thực sự được dùng — không thêm "cho chắc".
- Repository không trả lỗi driver-specific ra ngoài (`sql.ErrNoRows`, `pgx.ErrNoRows`...) — luôn convert sang lỗi domain (`apperror.ErrNotFound`...).
- Repository không log lỗi — chỉ trả lỗi kèm ngữ cảnh (`%w`), tầng gọi nó quyết định log ở đâu.
- Repository không chứa business logic (default value, tính toán nghiệp vụ) — chỉ thực thi lưu trữ/truy vấn dữ liệu đã được chuẩn bị sẵn từ Service.
- Method nhận `context.Context` làm tham số đầu tiên, lấy executor (DB/Tx) từ context để hỗ trợ transaction xuyên suốt.

## Scope của 1 repository
- 1 repository tương ứng 1 aggregate/entity chính — không tạo `GenericRepository` dùng chung cho mọi bảng.
- Query phục vụ read-case đặc thù (join nhiều bảng cho dashboard...) đặt trong repository liên quan nhất hoặc tách riêng Query Object, không nhét vào repository không rõ ranh giới.

## Pagination/filter — bắt buộc
- Danh sách có phân trang phải trả kèm `Total` trong cùng 1 lần gọi (dùng struct kết quả riêng như `ListXxxResult{Items, Total}`), không để tầng service tự query đếm riêng.

## Checklist khi generate code repository
- [ ] Interface ở domain, implementation ở repository/<impl>.
- [ ] Không có business rule nào trong file repository.
- [ ] Lỗi trả ra là lỗi domain, không phải lỗi driver gốc.
- [ ] Có `context.Context` làm tham số đầu, hỗ trợ transaction qua context.
- [ ] Danh sách phân trang trả kèm Total.