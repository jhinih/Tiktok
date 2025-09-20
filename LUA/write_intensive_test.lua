-- write_intensive_test.lua
-- 写入密集型压测脚本 - 测试数据库写入性能和事务处理能力

local json = require("json")

-- 写入配置
local write_config = {
    users = {},
    video_titles = {},
    comments = {},
    batch_size = 10
}

-- 初始化测试数据
function init_write_data()
    -- 用户池
    for i = 1, 1000 do
        write_config.users[i] = {
            id = "user_" .. i,
            username = "testuser" .. i,
            email = "test" .. i .. "@example.com"
        }
    end
    
    -- 视频标题池
    local titles = {
        "搞笑视频", "美食分享", "旅行日记", "音乐MV", "舞蹈表演",
        "技能教学", "生活记录", "宠物日常", "运动健身", "时尚穿搭",
        "Funny Video", "Food Share", "Travel Vlog", "Music Video", "Dance Performance"
    }
    for i = 1, #titles do
        for j = 1, 20 do
            write_config.video_titles[#write_config.video_titles + 1] = titles[i] .. " " .. j
        end
    end
    
    -- 评论内容池
    write_config.comments = {
        "太棒了！", "很有趣", "喜欢", "不错", "继续加油", "支持", "赞一个",
        "Amazing!", "Love it!", "So cool!", "Awesome!", "Great!", "Nice!", "Perfect!",
        "哈哈哈", "太搞笑了", "笑死我了", "有才", "创意十足", "很用心", "质量很高",
        "学到了", "很实用", "感谢分享", "收藏了", "转发", "关注了", "期待更新"
    }
end

-- 获取随机用户
function get_random_user()
    return write_config.users[math.random(1, #write_config.users)]
end

-- 创建视频请求 (30% 的写入请求)
function create_video_request()
    local user = get_random_user()
    local title = write_config.video_titles[math.random(1, #write_config.video_titles)]
    local descriptions = {
        "分享生活中的美好时刻", "希望大家喜欢", "用心制作的内容",
        "记录生活点滴", "与大家分享", "原创内容，请多支持"
    }
    
    local body = json.encode({
        video_path = "https://example.com/video/" .. math.random(1, 10000) .. ".mp4",
        cover_path = "https://example.com/cover/" .. math.random(1, 10000) .. ".jpg",
        title = title,
        description = descriptions[math.random(1, #descriptions)],
        is_private = math.random() > 0.8, -- 20% 概率设为私密
        type = "video",
        user_id = user.id
    })
    
    return wrk.format("POST", "/api/videos/create-video",
                    {["Content-Type"] = "application/json"}, body)
end

-- 批量点赞请求 (25% 的写入请求)
function batch_like_request()
    local user = get_random_user()
    local video_id = "video_" .. math.random(1, 50000)
    local owner_id = "user_" .. math.random(1, 1000)
    
    local body = json.encode({
        video_id = video_id,
        owner_id = owner_id,
        user_id = user.id
    })
    
    return wrk.format("POST", "/api/videos/like-video",
                    {["Content-Type"] = "application/json"}, body)
end

-- 评论写入请求 (20% 的写入请求)
function comment_write_request()
    local user = get_random_user()
    local video_id = "video_" .. math.random(1, 50000)
    local owner_id = "user_" .. math.random(1, 1000)
    local comment = write_config.comments[math.random(1, #write_config.comments)]
    
    local body = json.encode({
        video_id = video_id,
        content = comment,
        owner_id = owner_id,
        user_id = user.id
    })
    
    return wrk.format("POST", "/api/videos/comment-video",
                    {["Content-Type"] = "application/json"}, body)
end

-- 用户资料更新请求 (10% 的写入请求)
function update_profile_request()
    local user = get_random_user()
    local bios = {
        "热爱生活，分享快乐", "记录美好时光", "用心创作内容",
        "Life is beautiful", "Share happiness", "Create with passion",
        "音乐爱好者", "旅行达人", "美食探索者", "健身达人", "时尚博主"
    }
    
    local body = json.encode({
        id = user.id,
        username = user.username,
        avatar = "https://example.com/avatar/" .. math.random(1, 1000) .. ".jpg",
        bio = bios[math.random(1, #bios)],
        email = user.email
    })
    
    return wrk.format("POST", "/api/user/set-profile",
                    {["Content-Type"] = "application/json"}, body)
end

-- 添加好友请求 (10% 的写入请求)
function add_friend_request()
    local target_user = get_random_user()
    
    local body = json.encode({
        username = target_user.username
    })
    
    return wrk.format("POST", "/api/contact/add-friend",
                    {["Content-Type"] = "application/json"}, body)
end

-- 创建群聊请求 (3% 的写入请求)
function create_community_request()
    local user = get_random_user()
    local group_names = {
        "朋友圈", "兴趣小组", "学习交流", "工作群", "游戏群",
        "Friends", "Interest Group", "Study Group", "Work Team", "Gaming"
    }
    
    local body = json.encode({
        owner_name = user.username,
        name = group_names[math.random(1, #group_names)] .. "_" .. math.random(1, 1000),
        icon = "https://example.com/group/" .. math.random(1, 100) .. ".jpg",
        desc = "欢迎大家加入讨论"
    })
    
    return wrk.format("POST", "/api/contact/create-community",
                    {["Content-Type"] = "application/json"}, body)
end

-- 用户注册请求 (2% 的写入请求)
function register_request()
    local user_id = math.random(10000, 99999)
    local body = json.encode({
        email = "newuser" .. user_id .. "@test.com",
        code = string.format("%06d", math.random(100000, 999999)),
        password = "test123456",
        username = "newuser" .. user_id,
        avatar = "https://example.com/avatar/default.jpg"
    })
    
    return wrk.format("POST", "/api/login/register",
                    {["Content-Type"] = "application/json"}, body)
end

-- 主请求函数
function request()
    local rand = math.random(1, 100)
    
    if rand <= 30 then
        return create_video_request()
    elseif rand <= 55 then
        return batch_like_request()
    elseif rand <= 75 then
        return comment_write_request()
    elseif rand <= 85 then
        return update_profile_request()
    elseif rand <= 95 then
        return add_friend_request()
    elseif rand <= 98 then
        return create_community_request()
    else
        return register_request()
    end
end

-- 初始化
function init(args)
    math.randomseed(os.time())
    init_write_data()
    print("写入密集型压测初始化完成")
    print("用户池大小: " .. #write_config.users)
    print("视频标题池: " .. #write_config.video_titles)
    print("评论池: " .. #write_config.comments)
end

-- 响应处理
function response(status, headers, body)
    if status >= 500 then
        print("Server error: " .. status)
    elseif status == 429 then
        print("Rate limited - write too fast")
    elseif status >= 400 then
        print("Client error: " .. status)
    end
end

-- 完成统计
function done(summary, latency, requests)
    print("\n=== 写入密集型压测结果 ===")
    print("总写入请求: " .. summary.requests)
    print("写入错误数: " .. summary.errors.connect + summary.errors.read + summary.errors.write + summary.errors.status + summary.errors.timeout)
    print("平均写入延迟: " .. latency.mean / 1000 .. "ms")
    print("99%写入延迟: " .. latency:percentile(99) / 1000 .. "ms")
    print("写入QPS: " .. summary.requests / (summary.duration / 1000000))
    
    -- 计算写入性能指标
    local write_success_rate = (summary.requests - summary.errors.status) / summary.requests * 100
    print("写入成功率: " .. string.format("%.2f", write_success_rate) .. "%")
end