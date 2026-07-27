# Testing Rules

## Tỉ lệ các loại test
- Unit test (mock hết dependency): nhiều nhất (~70%).
- Integration test (Service + DB thật qua testcontainers): vừa phải (~20%).
- E2E test (qua HTTP thật): ít nhất (~10%), chỉ cover critical path (login, checkout).

## Nguyên tắc bắt buộc
- Test dùng table-driven pattern cho các case tương tự nhau; assertion dùng `testify` (`require`/`assert`).
- Unit test cho Service: mock repository qua interface đã khai báo ở domain — không chạm DB thật. Ưu tiên `mockgen` để generate mock khi interface nhiều method.
- Integration test cho Repository: dùng DB thật qua `testcontainers-go`, không mock SQL — đảm bảo query đúng cú pháp và đúng transaction behavior. Tách job riêng trong CI (`-short` cho unit, job riêng cho integration).
- Test Handler dùng `httptest`, không cần chạy server thật.
- Mỗi test tự tạo và dọn dẹp data của mình — không phụ thuộc thứ tự chạy hay data để lại từ test khác.

## Case bắt buộc phải cover cho JWT/Auth
- Token hợp lệ → verify thành công, claims đúng.
- Token hết hạn → lỗi `ErrTokenExpired`.
- Token sai chữ ký (giả mạo) → verify thất bại.
- Token đúng chữ ký nhưng `aud`/`iss` sai → verify thất bại.
- Refresh token đã bị revoke → không cấp lại access token.
- Refresh token bị dùng lại sau khi rotate (reuse detection) → revoke toàn bộ family token.
- Role trong token khác role hiện tại trong DB (permission_version) → request bị từ chối dù token còn hạn.

## Case bắt buộc phải cover cho Middleware web/API
- Middleware đọc token từ cookie (route web) hoạt động đúng khi có cookie, không có header.
- Middleware đọc token từ header (route api) hoạt động đúng khi có `Authorization: Bearer`, không có cookie.
- CSRF: request POST không có CSRF token hợp lệ ở route web bị từ chối; route API (Bearer) không bị áp CSRF.

## Coverage
- Không đặt mục tiêu % coverage cứng nhắc — ưu tiên cover đầy đủ business rule, error path, edge case của auth/permission hơn là cover 100% code đơn giản (getter/setter).

## Checklist khi generate test
- [ ] Unit test Service dùng mock, không chạm DB thật.
- [ ] Integration test Repository dùng DB thật, tách CI job riêng.
- [ ] Dùng table-driven pattern khi có nhiều case tương tự.
- [ ] Có test cho các case lỗi JWT (hết hạn, sai chữ ký, bị revoke, reuse detection).
- [ ] Middleware test cover cả nhánh web (cookie) và api (Bearer).
- [ ] Test không phụ thuộc thứ tự chạy hay side-effect từ test khác.