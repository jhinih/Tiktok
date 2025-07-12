package ulearning

// LoginResp 定义了登录响应的数据结构
type LoginResp struct {
	Authorization string `json:"AUTHORIZATION"` // 授权令牌
	UserID        int    `json:"userId"`        // 用户ID
	RoleID        int    `json:"roleId"`        // 角色ID
}

// GetAllCoursesResp 定义了获取所有课程列表的响应数据结构
type GetAllCoursesResp struct {
	CourseList []struct {
		ID      int    `json:"id"`      // 课程ID
		Name    string `json:"name"`    // 课程名称
		ClassID int    `json:"classId"` // 班级ID
	} `json:"courseList"` // 课程列表
}

// GetCourseActivitiesResp 定义了获取课程活动列表的响应数据结构
type GetCourseActivitiesResp struct {
	OtherActivityDTOList []struct {
		RelationID   int    `json:"relationId"`   // 关联ID
		RelationType int    `json:"relationType"` // 关联类型
		Title        string `json:"title"`        // 活动标题
		PersonStatus int    `json:"personStatus"` // 个人状态
		Status       int    `json:"status"`       // 活动状态
	} `json:"otherActivityDTOList"` // 其他活动列表
}

// GetActivityDetailResp 定义了获取活动详情的响应数据结构
type GetActivityDetailResp struct {
	AbsenceNum     int    `json:"absenceNum"`     // 缺勤人数
	NotAbsenceNum  int    `json:"notAbsenceNum"`  // 非缺勤人数
	Finish         string `json:"finish"`         // 完成状态
	Location       string `json:"location"`       // 位置信息
	AttendanceCode string `json:"attendanceCode"` // 签到码
}

// SigninUser 定义了签到用户的数据结构
type SigninUser struct {
	UserID int `json:"userID"` // 用户ID
	Status int `json:"status"` // 签到状态
}

// SigninTeacherOPReq 定义了教师签到操作的请求数据结构
type SigninTeacherOPReq struct {
	AttendanceID int          `json:"attendanceID"` // 签到ID
	Users        []SigninUser `json:"users"`        // 签到用户列表
}

// SigninTeacherOPResp 定义了教师签到操作的响应数据结构
type SigninTeacherOPResp struct {
	Msg    string `json:"msg"`    // 响应消息
	Status string `json:"status"` // 响应状态
}

// SigninOperationReq 定义了签到操作的请求数据结构
type SigninOperationReq struct {
	AttendanceID   int    `json:"attendanceID"`   // 签到ID
	ClassID        int    `json:"classID"`        // 班级ID
	UserID         int    `json:"userID"`         // 用户ID
	Location       string `json:"location"`       // 位置信息
	Address        string `json:"address"`        // 地址信息
	EnterWay       int    `json:"enterWay"`       // 签到方式
	AttendanceCode string `json:"attendanceCode"` // 签到码
}

// SigninOperationResp 定义了签到操作的响应数据结构
type SigninOperationResp struct {
	Msg       string `json:"msg"`       // 响应消息
	NewStatus int    `json:"newStatus"` // 新状态
	Status    int    `json:"status"`    // 响应状态
}
