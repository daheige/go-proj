package conf

import (
	"errors"

	"github.com/daheige/thinkgo/rediscache"
	"github.com/daheige/thinkgo/yamlconf"

	"github.com/gomodule/redigo/redis"
)

var conf *yamlconf.ConfigEngine

// InitConf 初始化配置文件
func InitConf(path string) {
	conf = yamlconf.NewConf()
	err := conf.LoadConf(path + "/app.yaml")
	if err != nil {
		panic(err)
	}

	AppEnv = conf.GetString("AppEnv", "production")
	switch AppEnv {
	case "local", "testing", "staging":
		AppDebug = true
	default:
		AppDebug = false
	}
}

// InitRedis 初始化redis
func InitRedis() {
	//初始化redis
	redisConf := &rediscache.RedisConf{}
	conf.GetStruct("RedisCommon", redisConf)

	// log.Println(redisConf)
	redisConf.SetRedisPool("default")
}

// GetRedisObj 从连接池中获取redis client
//用完就需要调用redisObj.Close()释放连接，防止过多的连接导致redis连接过多
// 导致当前请求而陷入长久等待，从而redis崩溃
func GetRedisObj(name string) (redis.Conn, error) {
	conn := rediscache.GetRedisClient(name)
	if conn == nil {
		return nil, errors.New("get redis client error")
	}

	return conn, nil
}
