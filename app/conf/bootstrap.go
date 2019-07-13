package conf

import (
	"errors"

	"github.com/daheige/thinkgo/redisCache"
	"github.com/daheige/thinkgo/yamlConf"

	"github.com/gomodule/redigo/redis"
)

var conf *yamlConf.ConfigEngine

func InitConf(path string) {
	conf = yamlConf.NewConf()
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

func InitRedis() {
	//初始化redis
	redisConf := &redisCache.RedisConf{}
	conf.GetStruct("RedisCommon", redisConf)

	// log.Println(redisConf)
	redisConf.SetRedisPool("default")
}

//从连接池中获取redis client
//用完就需要调用redisObj.Close()释放连接，防止过多的连接导致redis连接过多
// 导致当前请求而陷入长久等待，从而redis崩溃
func GetRedisObj(name string) (redis.Conn, error) {
	conn := redisCache.GetRedisClient(name)
	if conn == nil {
		return nil, errors.New("get redis client error")
	}

	return conn, nil
}
