package response

type JhinihCode struct {
	Code   int
	Jhinih string
}

//五位业务状态码

var (
	/* 成功 */
	SUCCESS = JhinihCode{Code: 20000, Jhinih: "成功"}

	/* 默认失败 */
	COMMON_FAIL = JhinihCode{-43960, "失败"}

	/* 请求错误 <0 */
	TOKEN_IS_EXPIRED   = JhinihCode{-20000, "token已过期"}
	TOKEN_IS_BLANK     = JhinihCode{-20001, "token为空"}
	TOKEN_NOT_VALID    = JhinihCode{-20002, "token无效"}
	TOKEN_TYPE_ERROR   = JhinihCode{-20003, "token类型错误"}
	TOKEN_FORMAT_ERROR = JhinihCode{-20004, "token格式错误"}
	HAVE_NOT_BEARER    = JhinihCode{-20005, "请求头中需要有Bearer字段"}
	RTOKEN_IS_EXPIRED  = JhinihCode{-20006, "rtoken已过期"}
	REQUEST_FREQUENTLY = JhinihCode{-20007, "请求过于频繁"}
	PERMISSION_DENIED  = JhinihCode{-20008, "权限不足"}

	/* 内部错误 60000 ~ 69999 */
	INTERNAL_ERROR              = JhinihCode{60001, "内部错误, check log"}
	INTERNAL_FILE_UPLOAD_ERROR  = JhinihCode{60002, "文件上传失败"}
	SNOWFLAKE_ID_GENERATE_ERROR = JhinihCode{60003, "snowflake id生成失败"}
	DATABASE_ERROR              = JhinihCode{60004, "数据库错误"}
	REDIS_ERROR                 = JhinihCode{60005, "redis错误"}
	EMAIL_SEND_ERROR            = JhinihCode{60006, "邮件发送失败"}

	/* 参数错误：10000 ~ 19999 */
	PARAM_NOT_VALID         = JhinihCode{10001, "参数无效"}
	PARAM_IS_BLANK          = JhinihCode{10002, "参数为空"}
	PARAM_TYPE_ERROR        = JhinihCode{10003, "参数类型错误"}
	PARAM_NOT_COMPLETE      = JhinihCode{10004, "参数缺失"}
	MEMBER_NOT_EXIST        = JhinihCode{10005, "用户不存在"}
	MESSAGE_NOT_EXIST       = JhinihCode{10006, "消息不存在"}
	EMAIL_NOT_VALID         = JhinihCode{10007, "邮箱格式错误"}
	VERIFY_CODE_VALID       = JhinihCode{10008, "验证码无效"}
	EMAIL_OR_PASSWORD_ERROR = JhinihCode{10009, "账号或密码错误"}
	USER_ALREADY_EXIST      = JhinihCode{10010, "用户已存在"}
	USER_NOT_EXIST          = JhinihCode{10011, "用户不存在"}
	ME_AND_ME               = JhinihCode{10012, "不能加自己"}
	FRIEND_YES_FRIEN        = JhinihCode{10013, "不能重复添加"}
	FRIEND_NOT_EXIT         = JhinihCode{10014, "好友ID为空"}
	COMMUNITY_IS_BLANK      = JhinihCode{10015, "群名称为空"}
	FACK_FACK_FACK          = JhinihCode{10016, "开发者不允许"}

	/* 用户错误 20000 ~ 29999 */
	USER_NOT_LOGIN = JhinihCode{20001, "用户未登录"}

	/*
	 USER_ACCOUNT_DISABLE(20005, "账号不可用"),
	 USER_ACCOUNT_LOCKED(20006, "账号被锁定"),
	 USER_ACCOUNT_NOT_EXIST(20007, "账号不存在"),
	 USER_ACCOUNT_USE_BY_OTHERS(20009, "账号下线"),
	 USER_ACCOUNT_EXPIRED(20010, "账号已过期"),
	*/
)
