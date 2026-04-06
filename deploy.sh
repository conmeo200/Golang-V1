#!/bin/bash

# Configuration
DOCKER_USER=${DOCKER_USER:-"nick2677"} # Default to nick2677 if not set
PROJECT_NAME="golang-v1"

TAG="latest"

echo "🚀 Bắt đầu quy trình Build và Push Docker images cho dự án: $PROJECT_NAME"

# 1. Build & Push App Server
echo "📦 Building App Server Image..."
docker build -t $DOCKER_USER/$PROJECT_NAME-app:$TAG -f Dockerfile.pro .
echo "📤 Pushing App Server to Docker Hub..."
docker push $DOCKER_USER/$PROJECT_NAME-app:$TAG

# 2. Build & Push Kong Gateway
echo "📦 Building Kong Gateway Image..."
docker build -t $DOCKER_USER/$PROJECT_NAME-kong:$TAG ./kong
echo "📤 Pushing Kong to Docker Hub..."
docker push $DOCKER_USER/$PROJECT_NAME-kong:$TAG

# 3. Build & Push Nginx Proxy
echo "📦 Building Nginx Proxy Image..."
docker build -t $DOCKER_USER/$PROJECT_NAME-nginx:$TAG ./nginx
echo "📤 Pushing Nginx to Docker Hub..."
docker push $DOCKER_USER/$PROJECT_NAME-nginx:$TAG

echo "✅ Hoàn tất! Tất cả images đã được đẩy lên Docker Hub tại: https://hub.docker.com/u/$DOCKER_USER"
