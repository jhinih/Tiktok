-- mixed_scenario_test.lua
-- 混合场景压测脚本 - 模拟真实生产环境的复杂业务场景

local json = require("json")

-- 场景配置
local scenario_config = {
    -- 用户行为模式
    user_patterns = {
        -- 活跃用户 (20% 用户产生 80% 流量)
        active_users = {},
        -- 普通用户 (60% 用户产生 15% 流量)
        normal_users = {},
        -- 潜水用户 (20% 用户产生 5% 流量)
        lurker_users = {}
    },
    
    -- 时间段权重 (模拟不同时间段的流量分布)
    time_weights = {
        peak_hours = {weight = 40, start_hour = 19, end_hour = 22},    -- 晚高峰
        normal_hours = {weight = 35, start_hour = 9, end_hour = 18},   -- 工作时间
        off_peak = {weight = 25, start_hour = 0, end_hour = 8}         -- 低峰期
    },
    
    -- 业务场景权重
    business_scenarios = {
        content_consumption = 45,  -- 内容消费 (浏览、观看)
        social_interaction = 25,   -- 社交互动 (点赞、评论、关注)
        content_creation = 15,     -- 内容创作 (发布视频、上传)
        user_management = 10,      -- 用户管理 (资料更新、设置)
        system_operation = 5       -- 系统操作 (搜索、推荐)
    }
}

-- 初始化用户模式
function init_user_patterns()
    -- 活跃用户 (高频操作)
    for i = 1, 200 do
        scenario_config.user_patterns.active_users[i] = {
            id = "active_user_" .. i,
            token = "token_active_" .. i,
            behavior_frequency = "high"
        }
    end
    
    -- 普通用户 (中频操作)
    for i = 1, 600 do
        scenario_config.user_patterns.normal_users[i] = {
            id = "normal_user_" .. i,
            token = "token_normal_" .. i,
            behavior_frequency = "medium"
        }
    end
    
    -- 潜水用户 (低频操作)
    for i = 1, 200 do
        scenario_config.user_patterns.lurker_users[i] = {
            id = "lurker_user_" .. i,
            token = "token_lurker_" .. i,
            behavior_frequency = "low"
        }
    end
end

-- 根据流量分布选择用户
function select_user_by_pattern()
    local rand = math.random(1, 100)
    
    if rand <= 80 then
        -- 80% 流量来自活跃用户
        return scenario_config.user_patterns.active_users[math.random(1, #scenario_config.user_patterns.active_users)]
    elseif rand <= 95 then
        -- 15% 流量来自普通用户
        return scenario_config.user_patterns.normal_users[math.random(1, #scenario_config.user_patterns.normal_users)]
    else
        -- 5% 流量来自潜水用户
        return scenario_config.user_patterns.lurker_users[math.random(1, #scenario_config.user_patterns.lurker_users)]
    end
end

-- 获取当前时间段权重
function get_current_time_weight()
    local current_hour = tonumber(os.date("%H"))
    
    if current_hour >= 19 and current_hour <= 22 then
        return scenario_config.time_weights.peak_hours.weight
    elseif current_hour >= 9 and current_hour <= 18 then
        return scenario_config.time_weights.normal_hours.weight
    else
        return scenario_config.time_weights.off_peak.weight
    end
end

-- 内容消费场景 (45%)
function content_consumption_scenario()
    local user = select_user_by_pattern()
    local scenarios = {
        -- 获取推荐视频列表
        function()
            local params = "?page=" .. math.random(1, 5) .. "&page_size=" .. math.random(10, 30) .. "&order_by=recommend"
            return wrk.format("GET", "/api/videos" .. params, 
                            {["Authorization"] = "Bearer " .. user.token})
        end,
        
        -- 获取热门视频点赞数
        function()
            local video_id = "hot_video_" .. math.random(1, 100)
            return wrk.format("GET", "/api/videos/get-video-likes?video_id=" .. video_id,
                            {["Authorization"] = "Bearer " .. user.token})
        end,
        
        -- 获取视频评论
        function()
            local video_id = "video_" .. math.random(1, 1000)
            return wrk.format("GET", "/api/videos/get-comments?id=" .. video_id .. "&is_video=true&page=1&page_size=20",
                            {["Authorization"] = "Bearer " .. user.token})
        end,
        
        -- 获取用户资料
        function()
            local target_user = "user_" .. math.random(1, 1000)
            return wrk.format("GET", "/api/user/get-profile?id=" .. target_user,
                            {["Authorization"] = "Bearer " .. user.token})
        end
    }
    
    return scenarios[math.random(1, #scenarios)]()
end

-- 社交互动场景 (25%)
function social_interaction_scenario()
    local user = select_user_by_pattern()
    local scenarios = {
        -- 点赞视频
        function()
            local body = json.encode({
                video_id = "video_" .. math.random(1, 1000),
                owner_id = "user_" .. math.random(1, 1000),
                user_id = user.id
            })
            return wrk.format("POST", "/api/videos/like-video",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end,
        
        -- 评论视频
        function()
            local comments = {
                "太棒了！", "很有趣", "喜欢这个", "不错哦", "继续加油",
                "Amazing!", "Love it!", "So cool!", "Awesome!", "Great!"
            }
            local body = json.encode({
                video_id = "video_" .. math.random(1, 1000),
                content = comments[math.random(1, #comments)],
                owner_id = "user_" .. math.random(1, 1000),
                user_id = user.id
            })
            return wrk.format("POST", "/api/videos/comment-video",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end,
        
        -- 添加好友
        function()
            local body = json.encode({
                username = "user_" .. math.random(1, 1000)
            })
            return wrk.format("POST", "/api/contact/add-friend",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end,
        
        -- 点赞评论
        function()
            local body = json.encode({
                user_id = user.id,
                comment_id = "comment_" .. math.random(1, 5000)
            })
            return wrk.format("POST", "/api/videos/like-comments",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end
    }
    
    return scenarios[math.random(1, #scenarios)]()
end

-- 内容创作场景 (15%)
function content_creation_scenario()
    local user = select_user_by_pattern()
    local scenarios = {
        -- 创建视频
        function()
            local titles = {"生活记录", "搞笑视频", "美食分享", "旅行日记", "技能教学"}
            local body = json.encode({
                video_path = "https://example.com/video/" .. math.random(1, 10000) .. ".mp4",
                cover_path = "https://example.com/cover/" .. math.random(1, 10000) .. ".jpg",
                title = titles[math.random(1, #titles)] .. " " .. math.random(1, 1000),
                description = "分享生活中的美好时刻",
                is_private = math.random() > 0.9,
                type = "video",
                user_id = user.id
            })
            return wrk.format("POST", "/api/videos/create-video",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end,
        
        -- 文件上传
        function()
            return wrk.format("POST", "/api/file/upload",
                            {["Content-Type"] = "multipart/form-data", ["Authorization"] = "Bearer " .. user.token}, 
                            "mock_file_data")
        end,
        
        -- 创建群聊
        function()
            local body = json.encode({
                owner_name = user.id,
                name = "群聊_" .. math.random(1, 1000),
                icon = "https://example.com/group.jpg",
                desc = "欢迎大家加入"
            })
            return wrk.format("POST", "/api/contact/create-community",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end
    }
    
    return scenarios[math.random(1, #scenarios)]()
end

-- 用户管理场景 (10%)
function user_management_scenario()
    local user = select_user_by_pattern()
    local scenarios = {
        -- 获取我的信息
        function()
            return wrk.format("GET", "/api/user/get-my-info",
                            {["Authorization"] = "Bearer " .. user.token})
        end,
        
        -- 更新用户资料
        function()
            local body = json.encode({
                id = user.id,
                username = user.id,
                avatar = "https://example.com/avatar/" .. math.random(1, 100) .. ".jpg",
                bio = "更新的个人简介 " .. math.random(1, 1000),
                email = user.id .. "@example.com"
            })
            return wrk.format("POST", "/api/user/set-profile",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end,
        
        -- 获取好友列表
        function()
            local body = json.encode({
                user_id = user.id,
                page = "1",
                page_size = "20",
                order_by = "created_at"
            })
            return wrk.format("POST", "/api/contact/get-friend-list",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end,
        
        -- 刷新Token
        function()
            local body = json.encode({
                rtoken = "refresh_" .. user.token
            })
            return wrk.format("POST", "/api/login/refresh-token",
                            {["Content-Type"] = "application/json"}, body)
        end
    }
    
    return scenarios[math.random(1, #scenarios)]()
end

-- 系统操作场景 (5%)
function system_operation_scenario()
    local user = select_user_by_pattern()
    local scenarios = {
        -- 获取在线用户列表
        function()
            return wrk.format("GET", "/api/contact/get-user-list-online",
                            {["Authorization"] = "Bearer " .. user.token})
        end,
        
        -- 加载群聊列表
        function()
            local body = json.encode({
                owner_id = user.id
            })
            return wrk.format("POST", "/api/contact/load-community",
                            {["Content-Type"] = "application/json", ["Authorization"] = "Bearer " .. user.token}, body)
        end,
        
        -- 获取WebSocket票据
        function()
            return wrk.format("GET", "/api/chat/ws-ticket",
                            {["Authorization"] = "Bearer " .. user.token})
        end
    }
    
    return scenarios[math.random(1, #scenarios)]()
end

-- 主请求函数
function request()
    local time_weight = get_current_time_weight()
    local adjusted_rand = math.random(1, 100) * (time_weight / 40) -- 根据时间段调整权重
    
    if adjusted_rand <= scenario_config.business_scenarios.content_consumption then
        return content_consumption_scenario()
    elseif adjusted_rand <= scenario_config.business_scenarios.content_consumption + scenario_config.business_scenarios.social_interaction then
        return social_interaction_scenario()
    elseif adjusted_rand <= scenario_config.business_scenarios.content_consumption + scenario_config.business_scenarios.social_interaction + scenario_config.business_scenarios.content_creation then
        return content_creation_scenario()
    elseif adjusted_rand <= 90 then
        return user_management_scenario()
    else
        return system_operation_scenario()
    end
end

-- 初始化
function init(args)
    math.randomseed(os.time())
    init_user_patterns()
    print("混合场景压测初始化完成")
    print("活跃用户: " .. #scenario_config.user_patterns.active_users)
    print("普通用户: " .. #scenario_config.user_patterns.normal_users)
    print("潜水用户: " .. #scenario_config.user_patterns.lurker_users)
    print("当前时间段权重: " .. get_current_time_weight())
end

-- 响应处理
function response(status, headers, body)
    if status >= 500 then
        print("Server error: " .. status)
    elseif status == 429 then
        print("Rate limited")
    elseif status >= 400 then
        print("Client error: " .. status)
    end
end

-- 完成统计
function done(summary, latency, requests)
    print("\n=== 混合场景压测结果 ===")
    print("总请求数: " .. summary.requests)
    print("总错误数: " .. summary.errors.connect + summary.errors.read + summary.errors.write + summary.errors.status + summary.errors.timeout)
    print("平均延迟: " .. latency.mean / 1000 .. "ms")
    print("50%延迟: " .. latency:percentile(50) / 1000 .. "ms")
    print("90%延迟: " .. latency:percentile(90) / 1000 .. "ms")
    print("99%延迟: " .. latency:percentile(99) / 1000 .. "ms")
    print("QPS: " .. summary.requests / (summary.duration / 1000000))
    
    -- 计算成功率
    local success_rate = (summary.requests - summary.errors.status) / summary.requests * 100
    print("成功率: " .. string.format("%.2f", success_rate) .. "%")
end