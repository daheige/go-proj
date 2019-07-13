# go-proj项目
    基于golang gin框架和grpc框架封装而成。
    涉及到的包：gin,gorm,redisgo,grpc,protobuf,viper

# 目录结构
    .
    ├── app             应用目录
    │   ├── conf        配置文件golang结构体定义
    │   ├── extensions  拓展层主要是errCode,logger定义
    │   ├── helper      助手函数，主要是redis,gorm,context相关函数
    │   ├── job         job任务/task任务
    │   ├── logic       逻辑层，服务于rpc,web,job三个层
    │   ├── rpc         grpc service服务层
    │   └── web         web/api层
    │       ├── api     api接口
    │       └── home    web首页，可以没有，根据项目来
    ├── bin             生成的golang二进制目录和shell脚步目录
    ├── logs            日志目录
    ├── cmd             各个应用对应main.go启动文件和对应
    │   ├── job
    │   ├── rpc
    │   └── web
    └── readme.md
# 关于项目部署
    建议将web,grpc,job分开单独部署，可采用不同的app.yaml配置文件启动

# 版权
    MIT
