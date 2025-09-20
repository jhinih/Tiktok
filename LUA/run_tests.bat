@echo off
REM TikTok 全链路压测执行脚本 (Windows版本)
REM 使用 wrk 工具执行不同场景的压测

setlocal enabledelayedexpansion

REM 配置参数
set BASE_URL=http://localhost:8080
set DURATION=60s
set THREADS=10
set CONNECTIONS=100

REM 颜色定义 (Windows CMD)
set RED=[91m
set GREEN=[92m
set YELLOW=[93m
set BLUE=[94m
set NC=[0m

REM 检查 wrk 是否安装
:check_wrk
where wrk >nul 2>&1
if %errorlevel% neq 0 (
    echo %RED%错误: wrk 工具未安装%NC%
    echo 请先安装 wrk: https://github.com/wg/wrk
    echo 或者下载 Windows 版本: https://github.com/wg/wrk/releases
    pause
    exit /b 1
)

REM 检查服务器是否可用
:check_server
echo %BLUE%检查服务器连接...%NC%
curl -s --connect-timeout 5 "%BASE_URL%/api/videos" >nul 2>&1
if %errorlevel% neq 0 (
    echo %RED%错误: 无法连接到服务器 %BASE_URL%%NC%
    echo 请确保 TikTok 后端服务正在运行
    pause
    exit /b 1
)
echo %GREEN%服务器连接正常%NC%
goto :eof

REM 打印测试信息
:print_test_info
set test_name=%1
set script_file=%2
echo.
echo %YELLOW%========================================%NC%
echo %YELLOW%开始执行: %test_name%%NC%
echo %YELLOW%========================================%NC%
echo 脚本文件: %script_file%
echo 目标地址: %BASE_URL%
echo 持续时间: %DURATION%
echo 线程数: %THREADS%
echo 连接数: %CONNECTIONS%
echo %YELLOW%----------------------------------------%NC%
goto :eof

REM 执行压测
:run_test
set test_name=%1
set script_file=%2
set custom_duration=%3
set custom_threads=%4
set custom_connections=%5

if "%custom_duration%"=="" set custom_duration=%DURATION%
if "%custom_threads%"=="" set custom_threads=%THREADS%
if "%custom_connections%"=="" set custom_connections=%CONNECTIONS%

call :print_test_info "%test_name%" "%script_file%"

REM 创建结果目录
if not exist results mkdir results

REM 生成时间戳
for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set "YY=%dt:~2,2%" & set "YYYY=%dt:~0,4%" & set "MM=%dt:~4,2%" & set "DD=%dt:~6,2%"
set "HH=%dt:~8,2%" & set "Min=%dt:~10,2%" & set "Sec=%dt:~12,2%"
set "timestamp=%YYYY%%MM%%DD%_%HH%%Min%%Sec%"

set result_file=results\%test_name%_%timestamp%.txt

REM 执行压测
echo 执行命令: wrk -t%custom_threads% -c%custom_connections% -d%custom_duration% -s %script_file% %BASE_URL%
wrk -t%custom_threads% -c%custom_connections% -d%custom_duration% -s %script_file% %BASE_URL% > %result_file% 2>&1

if %errorlevel% equ 0 (
    echo %GREEN%测试完成，结果已保存到: %result_file%%NC%
    type %result_file%
) else (
    echo %RED%测试执行失败%NC%
)
goto :eof

REM 全链路压测
:run_full_chain_test
echo.
echo %BLUE%=== 全链路压测 ===%NC%
call :run_test "full_chain" "full_chain_test.lua" "120s" "12" "150"
goto :eof

REM 热点数据压测
:run_hot_key_test
echo.
echo %BLUE%=== 热点数据压测 ===%NC%
call :run_test "hot_key" "hot_key_test.lua" "90s" "15" "200"
goto :eof

REM 写入密集型压测
:run_write_intensive_test
echo.
echo %BLUE%=== 写入密集型压测 ===%NC%
call :run_test "write_intensive" "write_intensive_test.lua" "90s" "8" "80"
goto :eof

REM 混合场景压测
:run_mixed_scenario_test
echo.
echo %BLUE%=== 混合场景压测 ===%NC%
call :run_test "mixed_scenario" "mixed_scenario_test.lua" "180s" "10" "120"
goto :eof

REM 原有的读取压测
:run_read_test
echo.
echo %BLUE%=== 读取压测 (原有) ===%NC%
call :run_test "read_only" "read.lua" "60s" "20" "300"
goto :eof

REM 压力递增测试
:run_stress_test
echo.
echo %BLUE%=== 压力递增测试 ===%NC%

set connections_list=50 100 200 400 800
set threads_list=5 10 15 20 25
set i=0

for %%c in (%connections_list%) do (
    set /a i+=1
    for /f "tokens=!i!" %%t in ("%threads_list%") do (
        echo.
        echo %YELLOW%压力级别 !i!: %%t 线程, %%c 连接%NC%
        call :run_test "stress_level_!i!" "mixed_scenario_test.lua" "60s" "%%t" "%%c"
        
        echo %BLUE%等待系统恢复...%NC%
        timeout /t 10 /nobreak >nul
    )
)
goto :eof

REM 生成测试报告
:generate_report
echo.
echo %BLUE%=== 生成测试报告 ===%NC%

for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set "YY=%dt:~2,2%" & set "YYYY=%dt:~0,4%" & set "MM=%dt:~4,2%" & set "DD=%dt:~6,2%"
set "HH=%dt:~8,2%" & set "Min=%dt:~10,2%" & set "Sec=%dt:~12,2%"
set "timestamp=%YYYY%%MM%%DD%_%HH%%Min%%Sec%"

set report_file=results\test_report_%timestamp%.md

echo # TikTok 全链路压测报告 > %report_file%
echo. >> %report_file%
echo ## 测试环境 >> %report_file%
echo - 目标服务器: %BASE_URL% >> %report_file%
echo - 测试时间: %date% %time% >> %report_file%
echo - 操作系统: Windows >> %report_file%
echo. >> %report_file%
echo ## 测试场景 >> %report_file%
echo. >> %report_file%
echo ### 1. 全链路压测 >> %report_file%
echo - **目标**: 模拟真实用户完整业务流程 >> %report_file%
echo - **包含**: 注册、登录、浏览视频、点赞、评论、社交互动 >> %report_file%
echo - **权重**: 登录30%%、获取视频25%%、点赞查询20%%、互动25%% >> %report_file%
echo. >> %report_file%
echo ### 2. 热点数据压测 >> %report_file%
echo - **目标**: 测试缓存和热key处理能力 >> %report_file%
echo - **场景**: 80%%流量访问1%%热点数据 >> %report_file%
echo - **包含**: 热点视频查询、大V用户信息、热门评论 >> %report_file%
echo. >> %report_file%
echo ### 3. 写入密集型压测 >> %report_file%
echo - **目标**: 测试数据库写入性能 >> %report_file%
echo - **包含**: 视频创建、点赞写入、评论发布、用户资料更新 >> %report_file%
echo - **比例**: 创建30%%、点赞25%%、评论20%%、其他25%% >> %report_file%
echo. >> %report_file%
echo ### 4. 混合场景压测 >> %report_file%
echo - **目标**: 模拟生产环境复杂场景 >> %report_file%
echo - **用户模式**: 活跃用户80%%流量、普通用户15%%、潜水用户5%% >> %report_file%
echo - **时间权重**: 根据当前时间调整流量分布 >> %report_file%
echo. >> %report_file%
echo ## 测试结果 >> %report_file%
echo. >> %report_file%

REM 添加最新的测试结果到报告
if exist results\*.txt (
    echo ### 最新测试结果 >> %report_file%
    echo. >> %report_file%
    
    for %%f in (results\*.txt) do (
        echo #### %%~nf >> %report_file%
        echo ``` >> %report_file%
        REM 获取文件最后20行 (Windows没有tail命令，使用PowerShell)
        powershell -command "Get-Content '%%f' | Select-Object -Last 20" >> %report_file%
        echo ``` >> %report_file%
        echo. >> %report_file%
    )
)

echo %GREEN%测试报告已生成: %report_file%%NC%
goto :eof

REM 清理旧结果
:clean_results
echo %YELLOW%清理旧的测试结果...%NC%
if exist results rmdir /s /q results
echo %GREEN%清理完成%NC%
goto :eof

REM 显示帮助信息
:show_help
echo %BLUE%TikTok 全链路压测工具 (Windows版本)%NC%
echo.
echo 用法: %0 [选项]
echo.
echo 选项:
echo   all              执行所有压测场景
echo   full-chain       全链路压测
echo   hot-key          热点数据压测
echo   write-intensive  写入密集型压测
echo   mixed-scenario   混合场景压测
echo   read-only        原有读取压测
echo   stress           压力递增测试
echo   report           生成测试报告
echo   clean            清理测试结果
echo   help             显示此帮助信息
echo.
echo 示例:
echo   %0 all                    # 执行所有测试
echo   %0 full-chain            # 只执行全链路测试
echo   %0 stress                # 执行压力递增测试
echo.
pause
goto :eof

REM 主函数
:main
set command=%1
if "%command%"=="" set command=help

if "%command%"=="all" (
    call :check_server
    call :run_full_chain_test
    call :run_hot_key_test
    call :run_write_intensive_test
    call :run_mixed_scenario_test
    call :run_read_test
    call :generate_report
) else if "%command%"=="full-chain" (
    call :check_server
    call :run_full_chain_test
) else if "%command%"=="hot-key" (
    call :check_server
    call :run_hot_key_test
) else if "%command%"=="write-intensive" (
    call :check_server
    call :run_write_intensive_test
) else if "%command%"=="mixed-scenario" (
    call :check_server
    call :run_mixed_scenario_test
) else if "%command%"=="read-only" (
    call :check_server
    call :run_read_test
) else if "%command%"=="stress" (
    call :check_server
    call :run_stress_test
) else if "%command%"=="report" (
    call :generate_report
) else if "%command%"=="clean" (
    call :clean_results
) else (
    call :show_help
)

pause
goto :eof

REM 执行主函数
call :main %1