# go-proj 项目

    基于golang gin框架和grpc框架封装而成。
    涉及到的包：gin,grpc,protobuf,daheige/thinkgo

# 目录结构

    .
    ├── app                             应用目录
    │   ├── job                         job/task作业层
    │   ├── logic                       公共逻辑层，上下文采用标准上下文ctx
    │   │   ├── BaseLogic.go
    │   │   ├── HomeLogic.go
    │   │   └── readme.md
    │   ├── rpc                         grpc service层
    │   │   └── service
    │   └── web                         web/api
    │       ├── controller
    │       ├── middleware
    │       └── routes
    ├── bin                             存放golang生成的二进制文件和shell脚本
    │   ├── go-gen                      golang生成的二进制文件
    │   │   ├── rpc
    │   │   └── web
    │   ├── nodejs-generate.sh
    │   ├── pb-generate.sh              golang pb和php pb代码生成脚本
    │   ├── php7.2_install.sh
    │   ├── pprof-check-version.sh      pprof性能监控生成自动版本号
    │   ├── web-check-version.sh        gin框架应用性能监控自动生成版本号
    │   └── web-init.sh                 golang rpc,web,job自动化构建脚本
    ├── conf                            项目配置文件目录
    ├── clients                         golang,php,nodejs客户端生成的代码
    │   ├── go
    │   │   └── client.go
    │   └── php
    │       ├── App                     自动生成的php代码
    │       ├── composer.json           composer文件，可以指定App命名空间自动加载
    │       ├── composer.lock
    │       ├── hello_client.php
    │       ├── readme.md
    │       └── vendor
    ├── cmd                             各个应用的main.go文件和配置文件app.yaml,线上可以放在别的目录
    │   ├── job
    │   ├── rpc
    │   │   ├── app.yaml                开发模式下的配置文件
    │   │   ├── logs
    │   │   └── main.go
    │   └── web
    │       ├── app.yaml
    │       ├── logs
    │       └── main.go
    ├── go.mod
    ├── go.sum
    ├── HealthCheck                     健康检查自动生成的代码
    │   ├── ginCheck
    │   │   └── checkversion.go
    │   ├── pprofCheck
    │   │   └── checkversion.go
    │   └── readme.md
    ├── library                         公共库主要是第三方库，logger,gin metrics监控等
    │   ├── helper                      助手函数库
    │   ├── ginMonitor                  gin web/api打点监控
    │   │   └── monitor.go
    │   └── Logger                      日志服务
    │       ├── log.go
    │       └── readme.md
    ├── LICENSE
    ├── logs                            运行日志目录，线上可放在别的目录,开发模式goland日志放在logs中
    │   ├── rpc
    │   └── web
    ├── pb                              根据pb协议，自动生成的golang pb代码
    │   └── hello.pb.go
    ├── protos                          pb协议文件
    │   └── hello.proto
    └── readme.md

# go-grpc 和 php grpc 工具安装

    参考https://github.com/daheige/hg-grpc

# 设置 golang 环境变量和 go mod 代理

    vim ~/.bashrc
    export GOROOT=/usr/local/go
    export GOOS=linux
    export GOPATH=/mygo
    export GOSRC=$GOPATH/src
    export GOBIN=$GOPATH/bin
    export GOPKG=$GOPATH/pkg

    #开启go mod机制
    export GO111MODULE=auto

    #禁用cgo模块
    export CGO_ENABLED=0

    # 阿里云代理
    export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

    # 也可以使用下面这个代理
    #export GOPROXY=https://goproxy.cn,direct

    #下面一行请根据实际情况修改
    export PATH=$GOROOT/bin:$GOBIN:$PATH

    保存退出:wq 使配置文件生效 source ~/.bashrc

# grpc 运行

    1、生成pb代码
        sh bin/pb-generate.sh
    2、启动服务端
    $ cp app.exam.yaml app.yaml
    $ sh bin/app-start.sh rpc
    2019/07/14 11:25:26 server pprof run on:  51051
    2019/07/14 11:25:26 go-proj grpc run on: 50051

    3、运行客户端
    $ go run clients/go/client.go
    2019/07/14 11:26:36 name:hello,golang grpc,message:call ok

    php客户端
    $ php clients/php/hello_client.php
    检测App\Grpc\GPBMetadata\Hello\HelloReq是否存在
    bool(true)
    status code: 0
    name:hello,world
    call ok

# woker job/task 运行

    开发环境下运行job/task
    $ sh bin/pprof-check-version.sh
    $ cp app.exam.yaml cmd/worker/app.yaml
    $ go run cmd/worker/worker.go
    2019/07/17 21:29:37 ===worker service start===
    2019/07/17 21:29:37 server pprof run on:  30031
    2019/07/17 21:29:38 hello world
    2019/07/17 21:29:39 current id:  heige
    2019/07/17 21:29:40 hello world
    2019/07/17 21:29:42 current id:  heige

# 项目工程化构建

    构建web
    $ sh bin/web-init.sh web
    初始化成功！
    生成自动版本号
    HealthCheck/pprofCheck/checkversion.go
    生成checkVersion.go成功
    HealthCheck/ginCheck/checkversion.go
    生成checkVersion.go成功
    开始构建web二进制文件
    构建web成功！

    构建rpc
    $ sh bin/web-init.sh rpc
    初始化成功！
    生成自动版本号
    HealthCheck/pprofCheck/checkversion.go
    生成checkVersion.go成功

    Generating codes...

    generating golang stubs...
    generating golang code success
    generating php stubs...
    generating php stubs from: /web/go/go-proj/protos/hello.proto
            [DONE]


    Generate codes successfully!

    开始构建web二进制文件
    构建rpc成功！

# 开发模式启动

    可以把项目中的app.exam.yaml复制到cmd对应的应用中，然后go run main.go启动

# 关于项目部署

    建议将web,grpc,job分开单独部署，可采用不同的app.yaml配置文件启动

# 项目上线说明

    1、可将bin下面的对应cmd下面的main.go生成的二进制文件，分发到线上部署，配置文件参考cmd/web/app.yaml
    2、上线二进制文件，需要指定app.yaml目录和logs目录

# grpc 中间件

    chain.go
    定义多个中间件（拦截器）
    // 注册interceptor和中间件
    opts = append(opts, grpc.UnaryInterceptor(
    	middleware.ChainUnaryServer(
    		middleware.RequestInterceptor,
    		middleware.Limit(&middleware.MockPassLimiter{}),
    	)))

    server := grpc.NewServer(opts...)
    具体demo参考cmd/rpc/main.go

    grpc中间件参考： https://github.com/grpc-ecosystem/go-grpc-middleware

# go mod 编译方式

    方式1:
                如果采用go mod tidy 拉取依赖，会在$GOPATH/pkg/mod缓存go.mod中的包
                cd cmd/web 然后执行 go build进行编译就可以，更新go.mod包，会自动进行包依赖拉取

    方式2:
                如果采用vendor机制，执行  go mod vendor  后会按需把项目用到的包，放在go-proj根目录的vendor下面
                cd cmd/web  然后执行  go build -mod=vendor 拉取当前目录下的vendor作为编译依赖包，进行编译

    以上两种方式推荐使用方式1，这样只需要升级go.mod中的包，就不需要担心包依赖问题

# 版权

    MIT
