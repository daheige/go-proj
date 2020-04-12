package conf

// 配置文件对应的变量
var (
	AppEnv   string
	AppDebug bool

	// web 服务中是否存在gRpc服务，主要是启动的时候是否要加载grpcconf中的gRpc client初始化
	WebHasGRPCService bool
)
