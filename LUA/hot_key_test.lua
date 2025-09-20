-- hot_key_test.lua
-- 热点数据压测脚本 - 专门测试缓存和热key处理能力

local json = require("json")

-- 热点配置
local hot_config = {
    -- 超级热点视频 (1% 的视频承载 80% 的流量)
    super_hot_videos = {},
    -- 普通热点视频 (9% 的视频承载 15% 的流量)  
    normal_hot_videos = {},
    -- 长尾视频 (90% 的视频承载 5% 的流量)
    long_tail_videos = {},
    
    -- 热点用户 (大V用户)
    hot_users = {},
    
    -- 流量分布权重
    super_hot_weight = 80,
    normal_hot_weight = 15,
    long_tail_weight = 5
}

-- 初始化热点数据
function init_hot_data()
    -- 10个超级热点视频
    for i = 1, 10 do
        hot_config.super_hot_videos[i] = "super_hot_" .. i
    end
    
    -- 90个普通热点视频
    for i = 1, 90 do
        hot_config.normal_hot_videos[i] = "normal_hot_" .. i
    end
    
    -- 900个长尾视频
    for i = 1, 900 do
        hot_config.long_tail_videos[i] = "long_tail_" .. i
    end
    
    -- 50个热点用户
    for i = 1, 50 do
        hot_config.hot_users[i] = "hot_user_" .. i
    end
end

-- 根据热点分布选择视频
function select_video_by_hotness()
    local rand = math.random(1, 100)
    
    if rand <= hot_config.super_hot_weight then
        -- 80% 概率选择超级热点视频
        return hot_config.super_hot_videos[math.random(1, #hot_config.super_hot_videos)]
    elseif rand <= hot_config.super_hot_weight + hot_config.normal_hot_weight then
        -- 15% 概率选择普通热点视频
        return hot_config.normal_hot_videos[math.random(1, #hot_config.normal_hot_videos)]
    else
        -- 5% 概率选择长尾视频
        return hot_config.long_tail_videos[math.random(1, #hot_config.long_tail_videos)]
    end
end

-- 选择热点用户
function select_hot_user()
    return hot_config.hot_users[math.random(1, #hot_config.hot_users)]
end

-- 热点读取场景 (70% 的请求)
function hot_read_request()
    local scenarios = {
        -- 获取视频点赞数 (最热点的查询)
        function()
            local video_id = select_video_by_hotness()
            return wrk.format("GET", "/api/videos/get-video-likes?video_id=" .. video_id)
        end,
        
        -- 获取用户信息 (热点用户查询)
        function()
            local user_id = select_hot_user()
            return wrk.format("GET", "/api/user/get-user-info?id=" .. user_id)
        end,
        
        -- 获取视频评论 (热点视频评论)
        function()
            local video_id = select_video_by_hotness()
            return wrk.format("GET", "/api/videos/get-comments?id=" .. video_id .. "&is_video=true")
        end,
        
        -- 获取用户资料 (热点用户资料)
        function()
            local user_id = select_hot_user()
            return wrk.format("GET", "/api/user/get-profile?id=" .. user_id)
        end
    }
    
    local scenario = scenarios[math.random(1, #scenarios)]
    return scenario()
end

-- 热点写入场景 (20% 的请求)
function hot_write_request()
    local scenarios = {
        -- 点赞热点视频
        function()
            local video_id = select_video_by_hotness()
            local user_id = "user_" .. math.random(1, 10000)
            local body = json.encode({
                video_id = video_id,
                owner_id = select_hot_user(),
                user_id = user_id
            })
            return wrk.format("POST", "/api/videos/like-video", 
                            {["Content-Type"] = "application/json"}, body)
        end,
        
        -- 评论热点视频
        function()
            local video_id = select_video_by_hotness()
            local user_id = "user_" .. math.random(1, 10000)
            local comments = {
                "太棒了！", "第一！", "沙发", "前排", "火钳刘明",
                "Amazing!", "First!", "Love it!", "So good!", "Awesome!"
            }
            local body = json.encode({
                video_id = video_id,
                content = comments[math.random(1, #comments)],
                owner_id = select_hot_user(),
                user_id = user_id
            })
            return wrk.format("POST", "/api/videos/comment-video",
                            {["Content-Type"] = "application/json"}, body)
        end,
        
        -- 关注热点用户
        function()
            local hot_user = select_hot_user()
            local body = json.encode({
                username = hot_user
            })
            return wrk.format("POST", "/api/contact/add-friend",
                            {["Content-Type"] = "application/json"}, body)
        end
    }
    
    local scenario = scenarios[math.random(1, #scenarios)]
    return scenario()
end

-- 冷数据访问场景 (10% 的请求)
function cold_data_request()
    local scenarios = {
        -- 访问长尾视频
        function()
            local video_id = hot_config.long_tail_videos[math.random(1, #hot_config.long_tail_videos)]
            return wrk.format("GET", "/api/videos/get-video-likes?video_id=" .. video_id)
        end,
        
        -- 获取普通用户信息
        function()
            local user_id = "normal_user_" .. math.random(1, 100000)
            return wrk.format("GET", "/api/user/get-user-info?id=" .. user_id)
        end,
        
        -- 获取视频列表 (分页查询)
        function()
            local page = math.random(10, 100) -- 访问较后面的页面
            return wrk.format("GET", "/api/videos?page=" .. page .. "&page_size=20")
        end
    }
    
    local scenario = scenarios[math.random(1, #scenarios)]
    return scenario()
end

-- 主请求函数
function request()
    local rand = math.random(1, 100)
    
    if rand <= 70 then
        -- 70% 热点读取
        return hot_read_request()
    elseif rand <= 90 then
        -- 20% 热点写入
        return hot_write_request()
    else
        -- 10% 冷数据访问
        return cold_data_request()
    end
end

-- 初始化
function init(args)
    math.randomseed(os.time())
    init_hot_data()
    print("热点数据压测初始化完成")
    print("超级热点视频: " .. #hot_config.super_hot_videos)
    print("普通热点视频: " .. #hot_config.normal_hot_videos)
    print("长尾视频: " .. #hot_config.long_tail_videos)
    print("热点用户: " .. #hot_config.hot_users)
end

-- 响应处理
function response(status, headers, body)
    if status >= 500 then
        print("Server error: " .. status)
    elseif status == 429 then
        print("Rate limited")
    end
end