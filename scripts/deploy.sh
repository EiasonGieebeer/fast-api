#!/bin/bash
# ============================================================
# Fast API 手动部署脚本
# 用法:
#   sh scripts/deploy.sh                    # 构建并部署到服务器
#   sh scripts/deploy.sh build-only         # 只构建，不上传
# ============================================================
set -e

# ---- 配置 ----
SERVER="root@139.129.36.232"
CONTAINER="fast-api"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BUILD_DIR="/tmp/fast-api-deploy"

mkdir -p "$BUILD_DIR"

echo ""
echo "=== 1. 构建前端 ==="
cd "$PROJECT_DIR/web"
bun run build

echo ""
echo "=== 2. 编译 Go 二进制（Linux amd64）==="
cd "$PROJECT_DIR"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -ldflags "-s -w -X 'github.com/QuantumNous/new-api/common.Version=$(cat VERSION)'" \
  -o "$BUILD_DIR/fast-api" .

echo ""
echo "=== 3. 检查构建产物 ==="
ls -lh "$BUILD_DIR/fast-api"

if [ "$1" = "build-only" ]; then
    echo ""
    echo "=== 构建完成（仅本地）==="
    echo "二进制: $BUILD_DIR/fast-api"
    exit 0
fi

echo ""
echo "=== 4. 上传到服务器 ==="
MSYS_NO_PATHCONV=1 scp "$BUILD_DIR/fast-api" "${SERVER}:/tmp/fast-api-new"

echo ""
echo "=== 5. 替换容器内二进制并重启 ==="
ssh "$SERVER" <<'REMOTE'
    set -e
    chmod +x /tmp/fast-api-new
    docker stop fast-api 2>/dev/null || true
    docker cp /tmp/fast-api-new fast-api:/fast-api
    docker start fast-api
    sleep 4
    echo "--- 容器日志 ---"
    docker logs --tail 5 fast-api
    echo "--- 健康检查 ---"
    curl -s -o /dev/null -w "首页: %{http_code}\n" -m 5 https://www.fastapi.cool/
    curl -s -o /dev/null -w "About: %{http_code}\n" -m 5 https://www.fastapi.cool/about
REMOTE

echo ""
echo "=== 部署完成 ==="
