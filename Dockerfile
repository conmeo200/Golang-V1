# Sử dụng image golang alpine cho môi trường Development
FROM golang:1.25-alpine

# Thiết lập thư mục làm việc
WORKDIR /app

# Cài đặt git để Air có thể theo dõi thay đổi
RUN apk add --no-cache git

# 1. Cài đặt Air (Hot Reload) để tự động build khi code thay đổi
# Rất hữu ích khi dev, giúp bạn không phải restart docker thủ công
RUN go install github.com/air-verse/air@latest

# 2. Tối ưu Cache cho dependencies
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy source code
COPY . .

# Port mặc định của App
EXPOSE 8080

# Lệnh mặc định chạy Air để tự động reload
CMD ["air"]