package schema

// StatusText 定义状态文本
type StatusText string

func (t StatusText) String() string {
	return string(t)
}

// 定义HTTP状态文本常量
const (
	OKStatus    StatusText = "OK"
	ErrorStatus StatusText = "ERROR"
	FailStatus  StatusText = "FAIL"
)

// StatusResult 响应状态
type StatusResult struct {
	Status StatusText `json:"status"` // 状态(OK)
}

// ErrorResult 响应错误
type ErrorResult struct {
	Error ErrorItem `json:"error"` // 错误项
}

// ErrorItem 响应错误项
type ErrorItem struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

// NewIDResult 创建响应唯一标识实例
func NewIDResult(id string) *IDResult {
	return &IDResult{
		ID: id,
	}
}

// IDResult 响应唯一标识
type IDResult struct {
	ID string `json:"id"`
}
