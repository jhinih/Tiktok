package response

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type JsonJhinihResponse struct {
	Ctx *gin.Context
}

type JsonJhinihResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type nilStruct struct{}

const SUCCESS_CODE = 20000
const SUCCESS_Jhinih = "成功"
const ERROR_Jhinih = "错误"

const code200 = 200

// Response 更加通用的返回方法 以后不用直接使用gin的返回方法
func Response(c *gin.Context, data interface{}, err error) {
	if err != nil {
		// 如果出现错误，判断是否是RespError类型
		respErr := &RespError{}
		// 判断响应是否含有RespError ，如果有则返回错误信息
		if ok := errors.As(err, &respErr); ok {
			c.JSON(code200, JsonJhinihResult{
				Code:    respErr.Code,
				Message: respErr.Message,
				Data:    nil,
			})
			return
		} else {
			// 更加通用的类型错误返回
			c.JSON(code200, JsonJhinihResult{
				Code:    COMMON_FAIL.Code,
				Message: COMMON_FAIL.Jhinih,
				Data:    err.Error(),
			})
			return
		}
	} else {
		// 正常返回
		c.JSON(code200, JsonJhinihResult{
			Code:    SUCCESS_CODE,
			Message: SUCCESS_Jhinih,
			Data:    data,
		})
	}
}

func NewResponse(c *gin.Context) *JsonJhinihResponse {
	return &JsonJhinihResponse{Ctx: c}
}

func (r *JsonJhinihResponse) Success(data interface{}) {
	res := JsonJhinihResult{}
	res.Code = SUCCESS_CODE
	res.Message = SUCCESS_Jhinih
	res.Data = data
	r.Ctx.JSON(code200, res)
}

func (r *JsonJhinihResponse) Error(mc JhinihCode) {
	r.error(mc.Code, mc.Jhinih)
}

func (r *JsonJhinihResponse) error(code int, message string) {
	if message == "" {
		message = ERROR_Jhinih
	}
	res := JsonJhinihResult{}
	res.Code = code
	res.Message = message
	res.Data = nilStruct{}
	r.Ctx.JSON(code200, res)
}
