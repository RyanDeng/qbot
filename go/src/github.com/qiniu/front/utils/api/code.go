package api

type Code int

func (c Code) IsEmpty() bool {
	return c == 0
}

func (c Code) Humanize() string {
	return CodeHumanize[c]
}

// 通用错误
const (
	OK              Code = 200 // ok
	InvalidArgs     Code = 400 // 请求参数错误，或者数据未通过验证
	Unauthorized    Code = 401 // 提供的授权数据未通过（登录已过期，或者不正确）
	Forbidden       Code = 403 // 不允许使用此接口
	NotFound        Code = 404 // 资源不存在
	TooManyRequests Code = 429 // 访问频率超过限制
	ResultError     Code = 500 // 请求结果发生错误
	CSRFDetected    Code = 599 // 检查到 CSRF
)

// 特殊错误
const (
	// just a example
	ErrorCodeExample Code = 5000 // 特殊错误代码以 5000 起始

	SigninWrongInfo Code = 5100 // 账户或密码错误
	SigninFailed    Code = 5101 // 登录失败，可能服务器错误
	SigninBlocked   Code = 5102 // 超过5次，被Block，等待5分钟
	InvalidToken    Code = 5103 // token, refresh_token 过期或错误

	BalanceFailed Code = 5200 // 用户余额获取失败

	BucketNotFound Code = 5300 // Bucket 不存在
)

var CodeHumanize = map[Code]string{
	OK:              "ok",
	InvalidArgs:     "invalid args",
	Unauthorized:    "unauthorized",
	Forbidden:       "forbidden",
	NotFound:        "not found",
	TooManyRequests: "too many requests",
	ResultError:     "response result error",
	CSRFDetected:    "csrf attack detected",

	SigninWrongInfo: "signin wrong info, username or password is wrong",
	SigninFailed:    "signin failed, server return error",
	SigninBlocked:   "user is blocked 5 minutes",
	InvalidToken:    "signout failed, expired or invalid token",

	BalanceFailed: "get balance failed",

	BucketNotFound: "bucket not found",
}
