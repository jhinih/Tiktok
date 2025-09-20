-- full_chain_test.lua
-- TikTok 全链路压测脚本
-- 模拟真实用户行为：注册->登录->浏览视频->点赞->评论->社交互动

local json = require("json")
local math = require("math")

-- 全局配置
local config = {
    base_url = "http://localhost:8080",
    -- 热点视频ID池 (模拟热key场景)
    hot_video_ids = {},
    -- 用户池
    users = {},
    -- 测试邮箱域名
    email_domains = {"test.com", "example.com", "demo.com"},
    -- 随机用户名前缀
    username_prefixes = {"user", "test", "demo", "tiktok"},
}

-- 初始化热点视频ID
for i = 1, 1000 do
    config.hot_video_ids[i] = "video_hot_" .. i
end

-- 初始化测试用户池
for i = 1, 100 do
    local user = {
        email = "testuser" .. i .. "@" .. config.email_domains[math.random(1, #config.email_domains)],
        username = config.username_prefixes[math.random(1, #config.username_prefixes)] .. "_" .. i,
        password = "test123456",
        user_id = nil,
        token = nil
    }
    config.users[i] = user
end

-- 权重配置 (模拟真实业务比例)
local weights = {
    login = 30,           -- 30% 登录请求
    get_videos = 25,      -- 25% 获取视频列表
    get_video_likes = 20, -- 20% 获取视频点赞数 (热点查询)
    like_video = 10,      -- 10% 点赞视频
    comment_video = 8,    -- 8% 评论视频
    get_user_info = 5,    -- 5% 获取用户信息
    upload_file = 2       -- 2% 文件上传
}

-- 累计权重计算
local cumulative_weights = {}
local total = 0
for action, weight in pairs(weights) do
    total = total + weight
    cumulative_weights[#cumulative_weights + 1] = {action = action, weight = total}
end

-- 随机选择用户
function get_random_user()
    return config.users[math.random(1, #config.users)]
end

-- 随机选择热点视频
function get_random_hot_video()
    return config.hot_video_ids[math.random(1, #config.hot_video_ids)]
end

-- 根据权重随机选择操作
function get_weighted_action()
    local rand = math.random(1, total)
    for _, item in ipairs(cumulative_weights) do
        if rand <= item.weight then
            return item.action
        end
    end
    return "get_videos" -- 默认操作
end

-- 生成随机验证码
function generate_code()
    return string.format("%06d", math.random(100000, 999999))
end

-- 构建请求头
function build_headers(token)
    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }
    if token then
        headers["Authorization"] = "Bearer " .. token
    end
    return headers
end

-- 用户注册请求
function register_request()
    local user = get_random_user()
    local body = json.encode({
        email = user.email,
        code = generate_code(),
        password = user.password,
        username = user.username,
        avatar = "https://example.com/avatar" .. math.random(1, 100) .. ".jpg"
    })
    
    return wrk.format("POST", "/api/login/register", build_headers(), body)
end

-- 用户登录请求
function login_request()
    local user = get_random_user()
    local body = json.encode({
        email = user.email,
        password = user.password
    })
    
    return wrk.format("POST", "/api/login/login", build_headers(), body)
end

-- 获取视频列表请求
function get_videos_request()
    local user = get_random_user()
    local params = string.format("?page=%d&page_size=%d&order_by=created_at", 
                                math.random(1, 10), 
                                math.random(10, 50))
    
    return wrk.format("GET", "/api/videos" .. params, build_headers(user.token))
end

-- 获取视频点赞数请求 (热点查询)
function get_video_likes_request()
    local video_id = get_random_hot_video()
    local path = "/api/videos/get-video-likes?video_id=" .. video_id
    local user = get_random_user()
    
    return wrk.format("GET", path, build_headers(user.token))
end

-- 点赞视频请求
function like_video_request()
    local user = get_random_user()
    local video_id = get_random_hot_video()
    local body = json.encode({
        video_id = video_id,
        owner_id = "owner_" .. math.random(1, 100),
        user_id = user.user_id or "user_" .. math.random(1, 1000)
    })
    
    return wrk.format("POST", "/api/videos/like-video", build_headers(user.token), body)
end

-- 评论视频请求
function comment_video_request()
    local user = get_random_user()
    local video_id = get_random_hot_video()
    local comments = {
        "太棒了！", "很有趣", "喜欢这个视频", "不错哦", "继续加油",
        "Amazing!", "Love it!", "So cool!", "Awesome!", "Great content!"
    }
    
    local body = json.encode({
        video_id = video_id,
        content = comments[math.random(1, #comments)],
        owner_id = "owner_" .. math.random(1, 100),
        user_id = user.user_id or "user_" .. math.random(1, 1000)
    })
    
    return wrk.format("POST", "/api/videos/comment-video", build_headers(user.token), body)
end

-- 获取用户信息请求
function get_user_info_request()
    local user = get_random_user()
    local target_user_id = "user_" .. math.random(1, 1000)
    local path = "/api/user/get-user-info?id=" .. target_user_id
    
    return wrk.format("GET", path, build_headers(user.token))
end

-- 文件上传请求 (模拟)
function upload_file_request()
    local user = get_random_user()
    -- 注意：这里只是模拟请求结构，实际文件上传需要multipart/form-data
    local headers = build_headers(user.token)
    headers["Content-Type"] = "multipart/form-data"
    
    return wrk.format("POST", "/api/file/upload", headers, "mock_file_data")
end

-- 主请求函数
function request()
    local action = get_weighted_action()
    
    if action == "login" then
        return login_request()
    elseif action == "get_videos" then
        return get_videos_request()
    elseif action == "get_video_likes" then
        return get_video_likes_request()
    elseif action == "like_video" then
        return like_video_request()
    elseif action == "comment_video" then
        return comment_video_request()
    elseif action == "get_user_info" then
        return get_user_info_request()
    elseif action == "upload_file" then
        return upload_file_request()
    else
        return get_videos_request()
    end
end

-- 响应处理函数
function response(status, headers, body)
    -- 统计不同状态码的响应
    if status == 200 then
        -- 成功响应，可以解析token等信息
        if body and string.find(body, "atoken") then
            -- 登录成功，提取token (简化处理)
            -- 实际项目中可以用JSON解析库
        end
    elseif status >= 400 then
        -- 错误响应统计
        print("Error response: " .. status)
    end
end

-- 初始化函数
function init(args)
    math.randomseed(os.time())
    print("TikTok 全链路压测初始化完成")
    print("用户池大小: " .. #config.users)
    print("热点视频数量: " .. #config.hot_video_ids)
end

-- 完成函数
function done(summary, latency, requests)
    print("\n=== 压测结果统计 ===")
    print("总请求数: " .. summary.requests)
    print("总错误数: " .. summary.errors.connect + summary.errors.read + summary.errors.write + summary.errors.status + summary.errors.timeout)
    print("平均延迟: " .. latency.mean / 1000 .. "ms")
    print("99%延迟: " .. latency:percentile(99) / 1000 .. "ms")
    print("QPS: " .. summary.requests / (summary.duration / 1000000))
end