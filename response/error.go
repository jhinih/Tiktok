package response

type RespError struct {
	JsonJhinihResult
}

// NewResponseError 包装响应错误类型，简化返回信息流程。
func ErrResponse(err error, result JhinihCode) error {
	respError := &RespError{}
	respError.Code = result.Code
	respError.Message = result.Jhinih
	respError.Data = err
	return respError
}

func (r RespError) Error() string {
	return r.Message
}
