package types

type BizCode int

const (
	Success       BizCode = 0
	Failed        BizCode = 1
	NotAuthorized BizCode = 401

	ErrorMsg    = "系统开小差了"
	InvalidArgs = "非法参数或参数解析失败"
)

type BizVo struct {
	Code    BizCode `json:"code"`
	Message string  `json:"message,omitempty"`
	Data    any     `json:"data,omitempty"`
}