#!/bin/bash

# API 集成测试执行脚本
# 用法: ./scripts/run_api_tests.sh [测试选项]
# 示例:
#   ./scripts/run_api_tests.sh                    # 运行所有测试
#   ./scripts/run_api_tests.sh -v                 # 运行所有测试，显示详细信息
#   ./scripts/run_api_tests.sh -run TestAuth      # 只运行认证相关测试

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
TEST_DIR="/root/code/bluebell/tests/api"
CONTAINER_NAME="bluebell-ai"
API_BASE_URL="http://localhost:8084"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Bluebell API 集成测试执行脚本       ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查 Docker 容器是否运行
echo -e "${YELLOW}[1/5] 检查 Docker 容器状态...${NC}"
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}错误: Docker 容器 $CONTAINER_NAME 未运行${NC}"
    echo "请先启动容器: docker start $CONTAINER_NAME"
    exit 1
fi
echo -e "${GREEN}✓ Docker 容器已运行${NC}"
echo ""

# 检查服务是否可访问
echo -e "${YELLOW}[2/5] 检查 API 服务可用性...${NC}"
if ! curl -sf "$API_BASE_URL/ping" > /dev/null 2>&1; then
    echo -e "${RED}错误: API 服务无法访问 ($API_BASE_URL)${NC}"
    echo "请确保服务已启动在 $API_BASE_URL"
    exit 1
fi
echo -e "${GREEN}✓ API 服务可访问${NC}"
echo ""

# 进入容器执行测试
echo -e "${YELLOW}[3/5] 执行 API 测试...${NC}"
echo ""

# 构建测试命令
TEST_ARGS="$@"
if [ -z "$TEST_ARGS" ]; then
    # 默认运行所有测试
    TEST_CMD="cd /root/code/bluebell && go test -v ./tests/api/... -count=1 2>&1"
else
    # 使用用户提供的参数
    TEST_CMD="cd /root/code/bluebell && go test ./tests/api/... $TEST_ARGS -count=1 2>&1"
fi

# 在容器中执行测试
docker exec -i "$CONTAINER_NAME" bash -c "$TEST_CMD" | tee /tmp/test_output.log || true

echo ""

# 分析测试结果
echo -e "${YELLOW}[4/5] 分析测试结果...${NC}"
if grep -q "FAIL" /tmp/test_output.log; then
    echo -e "${RED}✗ 测试失败${NC}"
    FAIL_COUNT=$(grep -c "FAIL" /tmp/test_output.log || echo "0")
    echo -e "${RED}失败用例数: $FAIL_COUNT${NC}"
    EXIT_CODE=1
else
    echo -e "${GREEN}✓ 所有测试通过${NC}"
    EXIT_CODE=0
fi

# 统计测试用例数
PASS_COUNT=$(grep -c "^--- PASS:" /tmp/test_output.log || echo "0")
SKIP_COUNT=$(grep -c "^--- SKIP:" /tmp/test_output.log || echo "0")
TOTAL_COUNT=$((PASS_COUNT + SKIP_COUNT))

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}           测试统计报告                ${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "通过: ${GREEN}$PASS_COUNT${NC}"
echo -e "跳过: ${YELLOW}$SKIP_COUNT${NC}"
if [ "$EXIT_CODE" -eq 1 ]; then
    echo -e "失败: ${RED}$FAIL_COUNT${NC}"
fi
echo -e "总计: ${BLUE}$TOTAL_COUNT${NC}"
echo -e "${BLUE}========================================${NC}"

# 如果有失败的测试，显示详细信息
if [ "$EXIT_CODE" -eq 1 ]; then
    echo ""
    echo -e "${YELLOW}失败用例详情:${NC}"
    grep "^--- FAIL:" /tmp/output.log | head -20 || echo "无法获取失败详情"
fi

exit $EXIT_CODE
