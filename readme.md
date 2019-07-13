# go-proj项目
    基于golang gin框架和grpc框架封装而成。
    涉及到的包：gin,gorm,redisgo,grpc,protobuf

# 目录结构
    .
    ├── app                             应用目录
    │   ├── conf                        项目配置文件golang定义
    │   │   ├── base.go
    │   │   └── bootstrap.go
    │   ├── extensions                  拓展层，一般是logger,eCode,第三方拓展
    │   │   └── Logger
    │   ├── helper                      公共库函数
    │   │   ├── context.go
    │   │   └── utils.go
    │   ├── job                         job/task作业层
    │   ├── logic                       公共逻辑层，上下文采用标准上下文ctx
    │   │   ├── BaseLogic.go
    │   │   ├── HomeLogic.go
    │   ├── rpc                         grpc服务层
    │   │   └── service
    │   └── web                         web/api层
    │       ├── controller              控制器
    │       ├── middleware              中间件
    │       └── routes                  路由设置
    ├── bin                             存放golang生成的二进制文件和shell脚本
    │   ├── web
    │   ├── web-check-version.sh        自动版本shell
    │   └── web-init.sh                 cmd各个应用构建脚本
    ├── cmd                             各个应用入口文件main.go
    │   ├── job
    │   ├── rpc
    │   └── web
    │       ├── app.yaml                项目上线可根据实际情况放置
    │       ├── logs
    │       ├── main.go                 web入口文件
    │       └── web                     go build生成的二进制文件
    ├── go.mod
    ├── go.sum
    ├── LICENSE
    ├── logs                            各个应用的日志文件
    │   └── web
    └── readme.md
# 关于项目部署
    建议将web,grpc,job分开单独部署，可采用不同的app.yaml配置文件启动

# 项目上线说明
    1、可将bin下面的对应cmd下面的main.go生成的二进制文件，分发到线上部署，配置文件参考cmd/web/app.yaml
    2、上线二进制文件，需要指定app.yaml目录和logs目录

# 版权
    MIT
