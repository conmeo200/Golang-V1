# ---------- Stage 1: Build ----------
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Cài git để tải module nếu cần
RUN apk add --no-cache git

# Tối ưu cache layer
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/main.go


# ---------- Stage 2: Runtime ----------
FROM alpine:3.19

WORKDIR /app

# Tạo user không phải root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary từ stage build
COPY --from=builder /app/app .

# Đổi quyền
RUN chown appuser:appgroup app

USER appuser

EXPOSE 8080

CMD ["./app"]
