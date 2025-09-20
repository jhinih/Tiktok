-- config.lua
-- 全链路压测配置文件

local config = {}

-- 服务器配置
config.server = {
    base_url = "http://localhost:8080",
    timeout = 30,  -- 超时时间(秒)
    retry_count = 3  -- 重试次数
}

-- 压测参数配置
config.load_test = {
    -- 默认压测参数
    default = {
        threads = 10,
        connections = 100,
        duration = "60s"
    },
    
    -- 全链路压测
    full_chain = {
        threads = 12,
        connections = 150,
        duration = "120s"
    },
    
    -- 热点数据压测
    hot_key = {
        threads = 15,
        connections = 200,
        duration = "90s"
    },
    
    -- 写入密集型压测
    write_intensive = {
        threads = 8,
        connections = 80,
        duration = "90s"
    },
    
    -- 混合场景压测
    mixed_scenario = {
        threads = 10,
        connections = 120,
        duration = "180s"
    },
    
    -- 压力递增测试
    stress_levels = {
        {threads = 5, connections = 50, duration = "60s"},
        {threads = 10, connections = 100, duration = "60s"},
        {threads = 15, connections = 200, duration = "60s"},
        {threads = 20, connections = 400, duration = "60s"},
        {threads = 25, connections = 800, duration = "60s"}
    }
}

-- 业务数据配置
config.business = {
    -- 用户配置
    users = {
        pool_size = 1000,  -- 用户池大小
        active_ratio = 0.2,  -- 活跃用户比例
        normal_ratio = 0.6,  -- 普通用户比例
        lurker_ratio = 0.2   -- 潜水用户比例
    },
    
    -- 视频配置
    videos = {
        hot_video_count = 1000,  -- 热点视频数量
        super_hot_count = 10,    -- 超级热点视频数量
        normal_hot_count = 90,   -- 普通热点视频数量
        long_tail_count = 900    -- 长尾视频数量
    },
    
    -- 流量分布权重
    traffic_weights = {
        super_hot = 80,  -- 超级热点流量权重
        normal_hot = 15, -- 普通热点流量权重
        long_tail = 5    -- 长尾流量权重
    },
    
    -- 业务场景权重
    scenario_weights = {
        -- 全链路场景
        full_chain = {
            login = 30,
            get_videos = 25,
            get_video_likes = 20,
            like_video = 10,
            comment_video = 8,
            get_user_info = 5,
            upload_file = 2
        },
        
        -- 热点数据场景
        hot_key = {
            hot_read = 70,
            hot_write = 20,
            cold_data = 10
        },
        
        -- 写入密集型场景
        write_intensive = {
            create_video = 30,
            batch_like = 25,
            comment_write = 20,
            update_profile = 10,
            add_friend = 10,
            create_community = 3,
            register = 2
        },
        
        -- 混合场景
        mixed_scenario = {
            content_consumption = 45,
            social_interaction = 25,
            content_creation = 15,
            user_management = 10,
            system_operation = 5
        }
    }
}

-- 时间段配置
config.time_patterns = {
    peak_hours = {
        start_hour = 19,
        end_hour = 22,
        weight_multiplier = 1.5  -- 高峰期权重倍数
    },
    normal_hours = {
        start_hour = 9,
        end_hour = 18,
        weight_multiplier = 1.0  -- 正常时间权重
    },
    off_peak = {
        start_hour = 0,
        end_hour = 8,
        weight_multiplier = 0.6  -- 低峰期权重倍数
    }
}

-- API 端点配置
config.endpoints = {
    -- 登录相关
    auth = {
        send_code = "/api/login/send-code",
        register = "/api/login/register",
        login = "/api/login/login",
        refresh_token = "/api/login/refresh-token"
    },
    
    -- 用户相关
    user = {
        get_user_info = "/api/user/get-user-info",
        get_my_info = "/api/user/get-my-info",
        get_profile = "/api/user/get-profile",
        set_profile = "/api/user/set-profile",
        set_role = "/api/user/set-role"
    },
    
    -- 视频相关
    videos = {
        list = "/api/videos",
        create = "/api/videos/create-video",
        like = "/api/videos/like-video",
        get_likes = "/api/videos/get-video-likes",
        comment = "/api/videos/comment-video",
        comment_comments = "/api/videos/comment-comments",
        like_comments = "/api/videos/like-comments",
        get_comments = "/api/videos/get-comments",
        get_comment_all = "/api/videos/get-comment-all",
        get_comment_member = "/api/videos/get-comment-member"
    },
    
    -- 文件相关
    file = {
        upload = "/api/file/upload"
    },
    
    -- 社交相关
    contact = {
        add_friend = "/api/contact/add-friend",
        get_friend_list = "/api/contact/get-friend-list",
        get_user_list_online = "/api/contact/get-user-list-online",
        get_group_users = "/api/contact/get-group-users",
        get_group_list = "/api/contact/get-group-list",
        create_community = "/api/contact/create-community",
        join_community = "/api/contact/join-community",
        load_community = "/api/contact/load-community"
    },
    
    -- 聊天相关
    chat = {
        ws_ticket = "/api/chat/ws-ticket",
        ws = "/api/chat/ws"
    }
}

-- 测试数据模板
config.test_data = {
    -- 邮箱域名
    email_domains = {
        "test.com", "example.com", "demo.com", "mock.com", "load.test"
    },
    
    -- 用户名前缀
    username_prefixes = {
        "user", "test", "demo", "tiktok", "load", "perf", "stress"
    },
    
    -- 视频标题
    video_titles = {
        "搞笑视频", "美食分享", "旅行日记", "音乐MV", "舞蹈表演",
        "技能教学", "生活记录", "宠物日常", "运动健身", "时尚穿搭",
        "Funny Video", "Food Share", "Travel Vlog", "Music Video", "Dance Performance"
    },
    
    -- 评论内容
    comments = {
        chinese = {
            "太棒了！", "很有趣", "喜欢这个", "不错哦", "继续加油",
            "支持", "赞一个", "哈哈哈", "太搞笑了", "笑死我了",
            "有才", "创意十足", "很用心", "质量很高", "学到了",
            "很实用", "感谢分享", "收藏了", "转发", "关注了"
        },
        english = {
            "Amazing!", "Love it!", "So cool!", "Awesome!", "Great!",
            "Nice!", "Perfect!", "Fantastic!", "Incredible!", "Outstanding!",
            "Brilliant!", "Excellent!", "Wonderful!", "Superb!", "Marvelous!"
        }
    },
    
    -- 用户简介
    user_bios = {
        "热爱生活，分享快乐", "记录美好时光", "用心创作内容",
        "Life is beautiful", "Share happiness", "Create with passion",
        "音乐爱好者", "旅行达人", "美食探索者", "健身达人", "时尚博主"
    },
    
    -- 群组名称
    group_names = {
        "朋友圈", "兴趣小组", "学习交流", "工作群", "游戏群",
        "Friends", "Interest Group", "Study Group", "Work Team", "Gaming"
    }
}

-- 性能基准配置
config.benchmarks = {
    -- QPS 目标
    qps_targets = {
        full_chain = 800,
        hot_key = 1500,
        write_intensive = 300,
        mixed_scenario = 600
    },
    
    -- 延迟目标 (毫秒)
    latency_targets = {
        avg_latency = 100,
        p99_latency = 500,
        p999_latency = 1000
    },
    
    -- 成功率目标 (百分比)
    success_rate_targets = {
        full_chain = 99.5,
        hot_key = 99.8,
        write_intensive = 99.0,
        mixed_scenario = 99.5
    }
}

-- 监控配置
config.monitoring = {
    -- 需要监控的状态码
    monitor_status_codes = {200, 201, 400, 401, 403, 404, 429, 500, 502, 503, 504},
    
    -- 告警阈值
    alert_thresholds = {
        error_rate = 1.0,      -- 错误率超过1%告警
        avg_latency = 200,     -- 平均延迟超过200ms告警
        p99_latency = 1000,    -- 99%延迟超过1000ms告警
        qps_drop = 0.2         -- QPS下降超过20%告警
    }
}

-- 输出配置
config.output = {
    -- 结果目录
    results_dir = "results",
    
    -- 报告格式
    report_formats = {"txt", "json", "csv"},
    
    -- 是否生成详细日志
    verbose_logging = true,
    
    -- 是否保存原始数据
    save_raw_data = false
}

return config