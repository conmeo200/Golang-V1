# ---------- Giai đoạn 1: Build (Giai đoạn xây dựng) ----------
# Sử dụng image golang phiên bản mới nhất trên nền alpine cực nhẹ
FROM golang:1.24-alpine AS builder

# Thiết lập thư mục làm việc bên trong container
WORKDIR /app

# Cài đặt các công cụ cần thiết (git, build-base...)
RUN apk add --no-cache git build-base

# 1. Tối ưu hóa Cache cho Dependency
# Copy các file quản lý module trước để Docker có thể cache layer này
# Nếu go.mod/sum không đổi, Docker sẽ bỏ qua bước download ở lần build sau
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy toàn bộ mã nguồn vào container
COPY . .

# 3. Build tất cả các binary (Đa dịch vụ)
# -s -w: Giảm dung lượng file binary bằng cách xóa bảng ký hiệu và thông tin debug
# CGO_ENABLED=0: Tạo static binary để chạy được trên alpine không cần glibc
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/server ./cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/worker ./cmd/worker && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/migrate ./cmd/migrate && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/seed ./cmd/seed


# ---------- Giai đoạn 2: Runtime (Giai đoạn chạy) ----------
# Sử dụng image alpine phiên bản nhỏ nhất cho môi trường Production
FROM alpine:3.19

# Cài đặt thư viện cần thiết (nếu code có dùng SSL/TLS)
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

# Tạo User không phải Root để tăng cường bảo mật (Security Best Practice)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Chỉ copy những file binary cần thiết từ stage builder sang
# Điều này giúp image Production cực kỳ nhỏ gọn (~20-30MB)
COPY --from=builder /app/bin /app/bin
COPY --from=builder /app/web /app/web


# Đổi quyền sở hữu cho appuser
RUN chown -R appuser:appgroup /app

# Sử dụng User bảo mật
USER appuser

# Expose port (Dành cho Server)
EXPOSE 8080

# Mặc định là chạy Server, nhưng có thể ghi đè command để chạy Worker/Migrate
CMD ["./bin/server"]
