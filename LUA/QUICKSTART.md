# TikTok 全链路压测 - 快速开始指南

## 🚀 5分钟快速上手

### 1. 环境准备

#### 安装 wrk 工具

**Windows:**
```bash
# 下载预编译版本
# 访问: https://github.com/wg/wrk/releases
# 下载 wrk.exe 并添加到 PATH 环境变量
```

**Linux/macOS:**
```bash
# Ubuntu/Debian
sudo apt-get install wrk

# CentOS/RHEL  
sudo yum install wrk

# macOS
brew install wrk
```

#### 启动 TikTok 后端服务
```bash
cd /path/to/tiktok-backend
go run main.go
```

### 2. 执行压测

#### Windows 用户
```cmd
# 进入 LUA 目录
cd LUA

# 执行所有压测场景
run_tests.bat all

# 执行单个场景
run_tests.bat full-chain
run_tests.bat hot-key
run_tests.bat write-intensive
run_tests.bat mixed-scenario
```

#### Linux/macOS 用户
```bash
# 进入 LUA 目录
cd LUA

# 给脚本执行权限
chmod +x run_tests.sh

# 执行所有压测场景
./run_tests.sh all

# 执行单个场景
./run_tests.sh full-chain
./run_tests.sh hot-key
./run_tests.sh write-intensive
./run_tests.sh mixed-scenario
```

### 3. 查看结果

压测完成后，结果会保存在 `results/` 目录下：
- `*.txt` - 详细的压测结果
- `test_report_*.md` - 汇总报告

## 📊 压测场景说明

| 场景 | 描述 | 适用情况 |
|------|------|----------|
| **full-chain** | 全链路业务流程 | 综合性能评估 |
| **hot-key** | 热点数据访问 | 缓存性能测试 |
| **write-intensive** | 写入密集操作 | 数据库写入性能 |
| **mixed-scenario** | 混合业务场景 | 生产环境模拟 |
| **read-only** | 纯读取操作 | 读取性能基准 |

## 🎯 关键指标解读

### QPS (每秒请求数)
```
Requests/sec: 1234.56
```
- **优秀**: >1000 QPS
- **良好**: 500-1000 QPS  
- **需优化**: <500 QPS

### 延迟 (Latency)
```
Latency     50.00%   45.67ms
            75.00%   78.90ms
            90.00%  123.45ms
            99.00%  234.56ms
```
- **平均延迟**: <100ms 优秀
- **99%延迟**: <500ms 可接受
- **99.9%延迟**: <1000ms 临界

### 成功率
```
Non-2xx or 3xx responses: 12
```
- **目标**: 错误数 <1% 总请求数
- **告警**: 错误率 >5%

## 🔧 常见问题

### Q: 连接被拒绝
```
connect() failed: Connection refused
```
**解决方案:**
1. 确认后端服务已启动
2. 检查端口是否正确 (默认8080)
3. 检查防火墙设置

### Q: 延迟过高
```
99.00%  5000.00ms
```
**解决方案:**
1. 检查数据库性能
2. 查看缓存命中率
3. 分析慢查询日志
4. 调整连接池大小

### Q: 内存不足
```
wrk: unable to create thread
```
**解决方案:**
1. 减少并发连接数 (-c 参数)
2. 减少线程数 (-t 参数)
3. 增加系统内存

## ⚙️ 自定义配置

### 修改目标地址
编辑脚本文件，修改 `base_url`:
```lua
local config = {
    base_url = "http://your-server:port"
}
```

### 调整压测参数
```bash
# 自定义参数执行
wrk -t10 -c100 -d60s -s full_chain_test.lua http://localhost:8080

# 参数说明:
# -t10: 10个线程
# -c100: 100个并发连接
# -d60s: 持续60秒
```

### 修改业务权重
编辑对应的 `.lua` 文件，调整权重配置:
```lua
local weights = {
    login = 30,           -- 30% 登录请求
    get_videos = 25,      -- 25% 获取视频
    -- 其他权重...
}
```

## 📈 性能优化建议

### 1. 数据库优化
- 添加适当索引
- 优化慢查询
- 调整连接池大小
- 考虑读写分离

### 2. 缓存优化
- 实现热点数据缓存
- 设置合理的过期时间
- 使用缓存预热
- 监控缓存命中率

### 3. 应用优化
- 异步处理非关键操作
- 实现接口限流
- 优化序列化性能
- 减少不必要的数据库查询

### 4. 系统优化
- 调整系统参数
- 优化网络配置
- 监控系统资源
- 实现负载均衡

## 📋 压测检查清单

### 压测前
- [ ] 后端服务正常运行
- [ ] 数据库连接正常
- [ ] 缓存服务可用
- [ ] 监控系统就绪
- [ ] 备份重要数据

### 压测中
- [ ] 监控系统资源使用
- [ ] 观察错误日志
- [ ] 记录异常情况
- [ ] 调整压测参数

### 压测后
- [ ] 分析测试结果
- [ ] 识别性能瓶颈
- [ ] 制定优化方案
- [ ] 验证优化效果

## 🆘 获取帮助

如遇到问题，请：
1. 查看 `README.md` 详细文档
2. 检查 `results/` 目录下的错误日志
3. 联系开发团队或提交 Issue

---

**提示**: 建议先在测试环境执行压测，确认无误后再在生产环境使用。