# Golang Project Structure Rules

Áp dụng cho backend Go phục vụ đồng thời Web (server-rendered, html/template) và REST API cho mobile app, dùng chung business logic.

## Cấu trúc thư mục bắt buộc
- `cmd/server/main.go` — entrypoint duy nhất, chỉ làm nhiệm vụ wiring dependency (khởi tạo DB, repository, service, router). Không chứa business logic, không chứa route handler cụ thể.
- `internal/auth/` — issue/verify JWT, định nghĩa claims.
- `internal/middleware/` — auth, permission, CSRF, logging, recover.
- `internal/rbac/` — định nghĩa role, permission, logic check quyền.
- `internal/config/` — load & validate config từ env/yaml.
- `internal/domain/<entity>/` — entity, interface repository, service (business logic thuần, không phụ thuộc framework HTTP/DB cụ thể).
- `internal/repository/<postgres|redis|...>/` — implement interface khai báo ở `domain/*/repository.go`.
- `internal/handler/web/` — HTTP handler trả HTML, dùng cookie.
- `internal/handler/api/v1/` — HTTP handler trả JSON, dùng Bearer token.
- `internal/router/` — khai báo route group + gắn middleware, không chứa logic.
- `internal/validator/` — custom validation rule dùng chung.
- `pkg/logger/`, `pkg/apperror/`, `pkg/response/` — code tái sử dụng được ở project khác.
- `web/templates/`, `web/static/` — asset cho SSR.
- `migrations/` — SQL migration.
- `test/integration/`, `test/testdata/`.

## Nguyên tắc phân tầng bắt buộc
`handler → service → repository (interface) → repository/<impl>`

- `domain/` không import `net/http`, không import framework HTTP (Gin/Echo/Fiber).
- `repository/<impl>` implement interface khai báo trong `domain/*/repository.go`; service chỉ phụ thuộc interface, không phụ thuộc implementation cụ thể.
- `handler/` là nơi duy nhất biết HTTP request/response/cookie/header. Handler gọi service, không tự viết business logic.
- `router/` chỉ định nghĩa route + middleware, không chứa logic xử lý.

## Vì sao tách handler/web và handler/api
- Response format khác nhau hoàn toàn (HTML vs JSON).
- Middleware khác nhau (cookie+CSRF vs Bearer token).
- Cả hai luôn gọi chung 1 service — không được viết business logic riêng cho từng platform.

## Quy tắc đặt package — bắt buộc
- Không dùng `common`, `util`, `helper` làm tên package chung chung.
- Không tạo package `models` chứa toàn bộ entity — tách theo domain.
- 1 file không vượt quá ~300–400 dòng; vượt thì tách theo trách nhiệm, không gộp bừa.
- Package name: chữ thường, không gạch dưới.

## Khi generate code mới, luôn:
1. Xác định entity thuộc domain nào → đặt đúng thư mục `internal/domain/<entity>/`.
2. Không viết SQL/query trực tiếp trong handler hoặc service.
3. Không tạo route mới mà không gắn qua `router/`.
4. Không thêm business logic vào `main.go`.