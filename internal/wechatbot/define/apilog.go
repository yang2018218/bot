package define

type ApiLogType string

const (
	ApiLogTypeRequest ApiLogType = "request" // 记录请求和返回
	ApiLogTypeReceive ApiLogType = "receive" // 只记录请求
)

type ApiLogLevel int

const (
	ApiLogLevel1 ApiLogLevel = 1 // 全记录
	ApiLogLevel2 ApiLogLevel = 2 // 只记录入参
)
