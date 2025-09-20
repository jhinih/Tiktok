#!/bin/bash

# TikTok 全链路压测执行脚本
# 使用 wrk 工具执行不同场景的压测

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置参数
BASE_URL="http://localhost:8080"
DURATION="60s"
THREADS=10
CONNECTIONS=100

# 检查 wrk 是否安装
check_wrk() {
    if ! command -v wrk &> /dev/null; then
        echo -e "${RED}错误: wrk 工具未安装${NC}"
        echo "请先安装 wrk: https://github.com/wg/wrk"
        exit 1
    fi
}

# 检查服务器是否可用
check_server() {
    echo -e "${BLUE}检查服务器连接...${NC}"
    if ! curl -s --connect-timeout 5 "$BASE_URL/api/videos" > /dev/null; then
        echo -e "${RED}错误: 无法连接到服务器 $BASE_URL${NC}"
        echo "请确保 TikTok 后端服务正在运行"
        exit 1
    fi
    echo -e "${GREEN}服务器连接正常${NC}"
}

# 打印测试信息
print_test_info() {
    local test_name=$1
    local script_file=$2
    
    echo -e "\n${YELLOW}========================================${NC}"
    echo -e "${YELLOW}开始执行: $test_name${NC}"
    echo -e "${YELLOW}========================================${NC}"
    echo -e "脚本文件: $script_file"
    echo -e "目标地址: $BASE_URL"
    echo -e "持续时间: $DURATION"
    echo -e "线程数: $THREADS"
    echo -e "连接数: $CONNECTIONS"
    echo -e "${YELLOW}----------------------------------------${NC}"
}

# 执行压测
run_test() {
    local test_name=$1
    local script_file=$2
    local custom_duration=${3:-$DURATION}
    local custom_threads=${4:-$THREADS}
    local custom_connections=${5:-$CONNECTIONS}
    
    print_test_info "$test_name" "$script_file"
    
    # 创建结果目录
    mkdir -p results
    local result_file="results/${test_name}_$(date +%Y%m%d_%H%M%S).txt"
    
    # 执行压测
    wrk -t$custom_threads -c$custom_connections -d$custom_duration \
        -s "$script_file" \
        "$BASE_URL" | tee "$result_file"
    
    echo -e "${GREEN}测试完成，结果已保存到: $result_file${NC}"
}

# 全链路压测
run_full_chain_test() {
    echo -e "\n${BLUE}=== 全链路压测 ===${NC}"
    run_test "full_chain" "full_chain_test.lua" "120s" 12 150
}

# 热点数据压测
run_hot_key_test() {
    echo -e "\n${BLUE}=== 热点数据压测 ===${NC}"
    run_test "hot_key" "hot_key_test.lua" "90s" 15 200
}

# 写入密集型压测
run_write_intensive_test() {
    echo -e "\n${BLUE}=== 写入密集型压测 ===${NC}"
    run_test "write_intensive" "write_intensive_test.lua" "90s" 8 80
}

# 混合场景压测
run_mixed_scenario_test() {
    echo -e "\n${BLUE}=== 混合场景压测 ===${NC}"
    run_test "mixed_scenario" "mixed_scenario_test.lua" "180s" 10 120
}

# 原有的读取压测
run_read_test() {
    echo -e "\n${BLUE}=== 读取压测 (原有) ===${NC}"
    run_test "read_only" "read.lua" "60s" 20 300
}

# 压力递增测试
run_stress_test() {
    echo -e "\n${BLUE}=== 压力递增测试 ===${NC}"
    
    local connections_list=(50 100 200 400 800)
    local threads_list=(5 10 15 20 25)
    
    for i in "${!connections_list[@]}"; do
        local conn=${connections_list[$i]}
        local thread=${threads_list[$i]}
        
        echo -e "\n${YELLOW}压力级别 $((i+1)): $thread 线程, $conn 连接${NC}"
        run_test "stress_level_$((i+1))" "mixed_scenario_test.lua" "60s" $thread $conn
        
        # 等待系统恢复
        echo -e "${BLUE}等待系统恢复...${NC}"
        sleep 10
    done
}

# 生成测试报告
generate_report() {
    echo -e "\n${BLUE}=== 生成测试报告 ===${NC}"
    
    local report_file="results/test_report_$(date +%Y%m%d_%H%M%S).md"
    
    cat > "$report_file" << EOF
# TikTok 全链路压测报告

## 测试环境
- 目标服务器: $BASE_URL
- 测试时间: $(date)
- 操作系统: $(uname -s)

## 测试场景

### 1. 全链路压测
- **目标**: 模拟真实用户完整业务流程
- **包含**: 注册、登录、浏览视频、点赞、评论、社交互动
- **权重**: 登录30%、获取视频25%、点赞查询20%、互动25%

### 2. 热点数据压测
- **目标**: 测试缓存和热key处理能力
- **场景**: 80%流量访问1%热点数据
- **包含**: 热点视频查询、大V用户信息、热门评论

### 3. 写入密集型压测
- **目标**: 测试数据库写入性能
- **包含**: 视频创建、点赞写入、评论发布、用户资料更新
- **比例**: 创建30%、点赞25%、评论20%、其他25%

### 4. 混合场景压测
- **目标**: 模拟生产环境复杂场景
- **用户模式**: 活跃用户80%流量、普通用户15%、潜水用户5%
- **时间权重**: 根据当前时间调整流量分布

## 测试结果

EOF

    # 添加最新的测试结果到报告
    if [ -d "results" ]; then
        echo "### 最新测试结果" >> "$report_file"
        echo "" >> "$report_file"
        
        for file in results/*.txt; do
            if [ -f "$file" ]; then
                echo "#### $(basename "$file" .txt)" >> "$report_file"
                echo '```' >> "$report_file"
                tail -20 "$file" >> "$report_file"
                echo '```' >> "$report_file"
                echo "" >> "$report_file"
            fi
        done
    fi
    
    echo -e "${GREEN}测试报告已生成: $report_file${NC}"
}

# 清理旧结果
clean_results() {
    echo -e "${YELLOW}清理旧的测试结果...${NC}"
    rm -rf results/*
    echo -e "${GREEN}清理完成${NC}"
}

# 显示帮助信息
show_help() {
    echo -e "${BLUE}TikTok 全链路压测工具${NC}"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  all              执行所有压测场景"
    echo "  full-chain       全链路压测"
    echo "  hot-key          热点数据压测"
    echo "  write-intensive  写入密集型压测"
    echo "  mixed-scenario   混合场景压测"
    echo "  read-only        原有读取压测"
    echo "  stress           压力递增测试"
    echo "  report           生成测试报告"
    echo "  clean            清理测试结果"
    echo "  help             显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 all                    # 执行所有测试"
    echo "  $0 full-chain            # 只执行全链路测试"
    echo "  $0 stress                # 执行压力递增测试"
}

# 主函数
main() {
    local command=${1:-help}
    
    # 检查依赖
    check_wrk
    
    case $command in
        "all")
            check_server
            run_full_chain_test
            run_hot_key_test
            run_write_intensive_test
            run_mixed_scenario_test
            run_read_test
            generate_report
            ;;
        "full-chain")
            check_server
            run_full_chain_test
            ;;
        "hot-key")
            check_server
            run_hot_key_test
            ;;
        "write-intensive")
            check_server
            run_write_intensive_test
            ;;
        "mixed-scenario")
            check_server
            run_mixed_scenario_test
            ;;
        "read-only")
            check_server
            run_read_test
            ;;
        "stress")
            check_server
            run_stress_test
            ;;
        "report")
            generate_report
            ;;
        "clean")
            clean_results
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# 执行主函数
main "$@"