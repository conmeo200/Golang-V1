#!/bin/bash

# Configuration
DOCKER_USER=${DOCKER_USER:-"nick2677"} # Default to nick2677 if not set
PROJECT_NAME="golang-v1"

# Tự động lấy mã commit ngắn (7 ký tự) làm version, nếu không có git thì dùng 'latest'
VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "v1.0.0")
TAG="latest"

echo "🚀 Bắt đầu quy trình Build và Push Docker images ($VERSION) cho dự án: $PROJECT_NAME"

# 1. Build & Push App Server
echo "📦 Building App Server Image..."
docker build -t $DOCKER_USER/$PROJECT_NAME-app:$VERSION -f Dockerfile.pro .
docker tag $DOCKER_USER/$PROJECT_NAME-app:$VERSION $DOCKER_USER/$PROJECT_NAME-app:$TAG
echo "📤 Pushing App Server to Docker Hub..."
docker push $DOCKER_USER/$PROJECT_NAME-app:$VERSION
docker push $DOCKER_USER/$PROJECT_NAME-app:$TAG

# 2. Build & Push Kong Gateway
echo "📦 Building Kong Gateway Image..."
docker build -t $DOCKER_USER/$PROJECT_NAME-kong:$VERSION ./kong
docker tag $DOCKER_USER/$PROJECT_NAME-kong:$VERSION $DOCKER_USER/$PROJECT_NAME-kong:$TAG
echo "📤 Pushing Kong to Docker Hub..."
docker push $DOCKER_USER/$PROJECT_NAME-kong:$VERSION
docker push $DOCKER_USER/$PROJECT_NAME-kong:$TAG

# 3. Build & Push Nginx Proxy
echo "📦 Building Nginx Proxy Image..."
docker build -t $DOCKER_USER/$PROJECT_NAME-nginx:$VERSION ./nginx
docker tag $DOCKER_USER/$PROJECT_NAME-nginx:$VERSION $DOCKER_USER/$PROJECT_NAME-nginx:$TAG
echo "📤 Pushing Nginx to Docker Hub..."
docker push $DOCKER_USER/$PROJECT_NAME-nginx:$VERSION
docker push $DOCKER_USER/$PROJECT_NAME-nginx:$TAG


echo "✅ Hoàn tất! Tất cả images đã được đẩy lên Docker Hub tại: https://hub.docker.com/u/$DOCKER_USER"
