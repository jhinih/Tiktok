package ulearning

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"tgwp/global"
	"tgwp/log/zlog"

	"github.com/levigross/grequests"
)

// User 定义了用户结构体
type User struct {
	Token  string // 用户令牌
	UserID int    // 用户ID
}

// teacherUser 全局教师用户变量
var teacherUser *User

// NewUser 创建新的用户实例
func NewUser() *User {
	return &User{}
}

// TeacherLogin 教师登录方法，使用配置中的教师账号和密码
func (l *User) TeacherLogin() (err error) {
	err = l.Login(global.Config.Ulearning.Teacher, global.Config.Ulearning.Password)
	teacherUser = l
	return
}

// RefreshTeacherToken 刷新教师令牌
func RefreshTeacherToken() (err error) {
	teacherUser = NewUser()
	err = teacherUser.TeacherLogin()
	return
}

// Login 基础登录方法，处理用户登录认证
func (l *User) Login(username, password string) (err error) {
	// 创建HTTP客户端，禁止自动重定向
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 构建登录表单数据
	data := url.Values{
		"loginName": {username},
		"password":  {password},
	}
	var req *http.Request
	req, err = http.NewRequest("POST", global.Config.Ulearning.Login, strings.NewReader(data.Encode()))
	if err != nil {
		zlog.Errorf("请求登录失败: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求并处理响应
	resp, err := client.Do(req)
	if err != nil {
		zlog.Errorf("请求登录失败: %v", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusFound {
		body, _ := ioutil.ReadAll(resp.Body)
		zlog.Errorf("请求登录失败: %s", string(body))
		return
	}

	// 从Cookie中获取认证信息
	authToken := ""
	userInfo := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "AUTHORIZATION" {
			authToken = cookie.Value
		} else if cookie.Name == "USERINFO" {
			userInfo, err = url.QueryUnescape(cookie.Value)
			if err != nil {
				zlog.Errorf("解析用户信息失败: %v", err)
				return
			}
		}
	}

	// 验证认证令牌
	if authToken == "" {
		zlog.Errorf("登录失败: 未获取到 AUTHORIZATION 值")
		return
	}

	// 解析用户信息
	var userDetails LoginResp
	if err = json.Unmarshal([]byte(userInfo), &userDetails); err != nil {
		zlog.Errorf("解析用户信息失败: %v", err)
		return
	}

	// 保存用户信息
	l.Token = authToken
	l.UserID = userDetails.UserID
	zlog.Debugf("Token: %s, UserID: %d", l.Token, l.UserID)
	return
}

// GetAllCourses 获取所有课程列表
func (l *User) GetAllCourses() (resp GetAllCoursesResp, err error) {
	// 构建请求参数
	Url := global.Config.Ulearning.GetAllCourses
	geq := &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded", // 设置内容类型
			"Authorization": l.Token,                             // 设置认证令牌
		},
	}
	var response *grequests.Response
	response, err = grequests.Get(Url, geq)
	if err != nil {
		zlog.Errorf("请求课程列表失败: %v", err)
		return
	}
	defer response.Close()

	// 解析响应数据
	if err = response.JSON(&resp); err != nil {
		zlog.Errorf("解析响应失败: %v", err)
		return
	}
	return
}

// GetCourseActivities 获取课程活动列表
func (l *User) GetCourseActivities(courseID int) (resp GetCourseActivitiesResp, err error) {
	// 构建请求URL和参数
	Url := fmt.Sprintf(global.Config.Ulearning.GetCourseActivities, courseID)
	geq := &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": l.Token,
		},
	}
	var response *grequests.Response
	response, err = grequests.Get(Url, geq)
	if err != nil {
		zlog.Errorf("请求课程活动失败: %v", err)
		return
	}
	defer response.Close()

	// 解析响应数据
	if err = response.JSON(&resp); err != nil {
		zlog.Errorf("解析响应失败: %v", err)
		return
	}
	return
}

// GetActivityDetail 获取活动详情
func (l *User) GetActivityDetail(relationID int) (resp GetActivityDetailResp, err error) {
	// 构建请求URL和参数
	Url := fmt.Sprintf(global.Config.Ulearning.GetActivityDetail, relationID)
	geq := &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": l.Token,
		},
	}
	var response *grequests.Response
	response, err = grequests.Get(Url, geq)
	if err != nil {
		zlog.Errorf("请求活动详情失败: %v", err)
		return
	}
	defer response.Close()

	// 解析响应数据
	if err = response.JSON(&resp); err != nil {
		zlog.Errorf("解析响应失败: %v", err)
		return
	}
	zlog.Debugf("活动详情: %v", resp)
	return
}

// SigninByTeacher 教师为学生进行签到操作
func (l *User) SigninByTeacher(relationID int, userID int) (err error) {
	// 构建签到用户数据
	user := SigninUser{
		UserID: userID,
		Status: 1,
	}
	users := []SigninUser{user}
	var reqData = SigninTeacherOPReq{
		AttendanceID: relationID,
		Users:        users,
	}

	// 构建请求参数
	Url := global.Config.Ulearning.SigninTeacherOperation
	geq := &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": l.Token,
		},
		JSON: reqData,
	}

	// 发送请求
	var response *grequests.Response
	response, err = grequests.Post(Url, geq)
	if err != nil {
		zlog.Errorf("请求签到失败: %v", err)
		return
	}
	defer response.Close()

	// 解析响应数据
	var resp SigninTeacherOPResp
	if err = response.JSON(&resp); err != nil {
		zlog.Errorf("解析响应失败: %v", err)
		return
	}

	// 检查签到状态
	if resp.Status != "success" {
		zlog.Errorf("签到失败: %s", resp.Msg)
		return errors.New(resp.Msg)
	}
	zlog.Infof("签到成功")
	return
}

// SigninByStudent 学生进行签到操作
func (l *User) SigninByStudent(relationID int, classID int) (err error) {
	// 获取签到活动详情
	data, err := teacherUser.GetActivityDetail(relationID)

	// 构建签到请求数据
	Url := global.Config.Ulearning.SigninOperation
	var reqData = SigninOperationReq{
		AttendanceID:   relationID,
		ClassID:        classID,
		UserID:         l.UserID,
		Location:       data.Location,
		Address:        "",
		EnterWay:       1,
		AttendanceCode: data.AttendanceCode,
	}

	// 构建请求参数
	geq := &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": l.Token,
		},
		JSON: reqData,
	}

	// 发送请求
	var response *grequests.Response
	response, err = grequests.Post(Url, geq)
	if err != nil {
		zlog.Errorf("请求签到失败: %v", err)
		return
	}
	defer response.Close()

	// 解析响应数据
	var resp SigninOperationResp
	if err = response.JSON(&resp); err != nil {
		zlog.Errorf("解析响应失败: %v", err)
		return
	}

	// 检查签到状态
	if resp.NewStatus != 1 {
		zlog.Errorf("签到失败: %s", resp.Msg)
		return errors.New(resp.Msg)
	}
	zlog.Infof("签到成功")
	return
}
