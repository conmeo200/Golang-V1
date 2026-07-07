# Sử dụng image golang alpine cho môi trường Development
FROM golang:1.25-alpine

# Thiết lập thư mục làm việc
WORKDIR /app

# Cài đặt git và curl để tải Air binary
RUN apk add --no-cache git curl

# 1. Tải Air binary dựng sẵn (Nhanh hơn rất nhiều so với go install)
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

  # 2. Tối ưu Cache cho dependencies
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy source code
COPY . .

# Port mặc định của App
EXPOSE 8080

# Lệnh mặc định chạy Air để tự động reload
CMD ["air"]