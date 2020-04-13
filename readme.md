# go-proj 项目

    基于golang gin框架和grpc框架封装而成。
    涉及到的包：gin,grpc,protobuf,redigo,daheige/thinkgo

# go version选择

    推荐使用go v1.14.1+版本
    如果是用的go v1.14以下版本，请使用 go-proj/v1分支

# 目录结构

    .
    ├── app                             应用目录
    │   ├── job                         job/task作业层
    │   ├── logic                       公共逻辑层，上下文采用标准上下文ctx
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
    ├── library                         公共库主要是第三方库，logger,gin metrics监控等
    │   ├── helper                      助手函数库
    │   ├── ginMonitor                  gin web/api打点监控
    │   └── logger                      日志服务
    ├── LICENSE
    ├── logs                            运行日志目录，线上可放在别的目录,开发模式goland日志放在logs中
    │   ├── rpc
    │   └── web
    ├── pb                              根据pb协议，自动生成的golang pb代码
    │   └── hello.pb.go
    ├── protos                          pb协议文件
    │   └── hello.proto
    └── readme.md

# 关于web层

    基于gin1.5.0+框架封装而成
    
# 关于gin validate参数校验

    gin1.6.2+ 基于gopkg.in/go-playground/validator.v10封装之后
    将validator库的validate tag改成了binding方便gin使用
    
    参考手册：
        https://github.com/go-playground/validator/tree/v9
        https://godoc.org/github.com/go-playground/validator
        https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
        
# gin使用手册
    
    参考 https://github.com/gin-gonic/gin
    中文翻译: https://github.com/daheige/gin-doc-cn 如果有更新，以官网为准
    
# go-grpc 和 php grpc 工具安装

    ubuntu系统
        参考 https://github.com/daheige/hg-grpc
    centos系统
        参考 docs/centos7-protoc-install.md
        
# golang 环境安装

    golang下载地址:
       https://golang.google.cn/dl/

    以go最新版本go1.14版本为例
    https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
    1、linux环境，下载
        cd /usr/local/
        sudo wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
        sudo tar zxvf go1.14.1.linux-amd64.tar.gz
        创建golang需要的目录
        sudo mkdir /mygo
        sudo mkdir /mygo/bin
        sudo mkdir /mygo/src
        sudo mkdir /mygo/pkg

    2、设置环境变量vim ~/.bashrc 或者sudo vim /etc/profile
        export GOROOT=/usr/local/go
        export GOOS=linux
        export GOPATH=/mygo
        export GOSRC=$GOPATH/src
        export GOBIN=$GOPATH/bin
        export GOPKG=$GOPATH/pkg
        #开启go mod机制
        export GO111MODULE=on

        #禁用cgo模块
        export CGO_ENABLED=0

        export PATH=$GOROOT/bin:$GOBIN:$PATH

    3、source ~/.bashrc 生效配置

# 设置 goproxy 代理

    go version >= 1.13
    设置goproxy代理
    vim ~/.bashrc添加如下内容:
    export GOPROXY=https://goproxy.io,direct
    或者
    export GOPROXY=https://goproxy.cn,direct
    或者
    export GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct

    让bashrc生效
    source ~/.bashrc

    go version < 1.13
    vim ~/.bashrc添加如下内容：
    export GOPROXY=https://goproxy.io
    或者使用 export GOPROXY=https://athens.azurefd.net
    或者使用 export GOPROXY=https://mirrors.aliyun.com/goproxy/
    让bashrc生效
    source ~/.bashrc

# php grpc工具和拓展安装

    参考 docs/centos7-protoc-install.md

# grpc 运行

    1、生成pb代码 （生成代码之前，请先安装好go-grpc 和 php grpc 工具）
        sh bin/go-generate.sh
        sh bin/php-generate.sh

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
    开始构建web二进制文件
    构建web成功！

    构建rpc
    $ sh bin/web-init.sh rpc
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

# nginx grpc_pass

    启动两个实例
    $ cd cmd/rpc
    $ go run main.go --port=50051
    2019/11/09 18:55:00 server pprof run on:  51051
    2019/11/09 18:55:00 go-proj grpc run on: 50051

    新开一个终端
    $ cd cmd/rpc
    $ go run main.go --port=50052 --log_dir=/web/wwwlogs
    2019/11/09 18:56:32 server pprof run on:  51052
    2019/11/09 18:56:32 go-proj grpc run on: 50052

    配置nginx grpc负载均衡
    参考:
    https://github.com/daheige/hg-grpc/blob/master/readme.md#%E9%85%8D%E7%BD%AEnginx-grpc%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1
    为了go grpc服务高可用，需要对grpc服务做负载均衡处理，这里借助nginx grpc模块实现，配置如下:

    #nginx gprc负载均衡配置，要求nginx1.13.0+以上版本
    #nginx gprc负载均衡配置
    # 多个ip:port实例
    upstream go_grpc {
        server 127.0.0.1:50051 weight=5 max_fails=3 fail_timeout=10;
        server 127.0.0.1:50052 weight=1 max_fails=3 fail_timeout=10;
    }

    server {
        listen 50050 http2;
        server_name localhost;

        access_log /web/wwwlogs/go-grpc-access.log;

        location / {
            grpc_pass grpc://go_grpc;
        }
    }

    重启nginx
    sudo service nginx restart

    运行grpc client
    $ go run clients/go/client_rpc.go
    2019/11/09 19:01:18 name:hello,golang grpc,message:call ok
    $ go run clients/go/client_rpc.go
    2019/11/09 19:01:20 name:hello,golang grpc,message:call ok
    $ go run clients/go/client_rpc.go
    2019/11/09 19:01:23 name:hello,golang grpc,message:call ok

    请求过程中可以看到两个实例，会产生响应日志
    2019/11/09 19:07:46 req method:  /App.Grpc.Hello.GreeterService/SayHello
    2019/11/09 19:07:46 req data:  name:"golang grpc"

    查看请求日志
    $ tail -f cmd/rpc/logs/go-grpc.log

    $ tail -f /web/wwwlogs/go-grpc.log

    通过查看日志，可以看到请求到50051这个实例的grpc请求相对多一点，50052这个实例相对少一点
    因为nginx grpc的权重不一样

# go mod 编译方式

    方式1:
                如果采用go mod tidy 拉取依赖，会在$GOPATH/pkg/mod缓存go.mod中的包
                cd cmd/web 然后执行 go build进行编译就可以，更新go.mod包，会自动进行包依赖拉取

    方式2:
                如果采用vendor机制，执行  go mod vendor  后会按需把项目用到的包，放在go-proj根目录的vendor下面
                cd cmd/web  然后执行  go build -mod=vendor 拉取当前目录下的vendor作为编译依赖包，进行编译

    以上两种方式推荐使用方式1，这样只需要升级go.mod中的包，就不需要担心包依赖问题

# docker 构建镜像

    web层:
        $ docker build -t go-proj-web:v1 -f web-Dockerfile .
        运行：
        $ docker run -it --name=go-proj-web -v /data/logs:/data/logs -v /data/www/go-proj:/data/conf -p 1338:1338 -p 2338:2338 go-proj-web:v1

    job层:
        $ docker build -it go-proj-job:v1 -f job-Dockerfile .
        运行:
        $ docker run -it --name=go-proj-job -v /data/logs:/data/logs -v /data/www/go-proj:/data/conf -p 30031:30031 go-proj-job:v1

    rpc层:
        $ docker build -t go-proj-job:v1 -f job-Dockerfile .
        运行:
             $ docker run -it --name=go-proj-rpc -v /data/logs:/data/logs -v /data/www/go-proj:/data/conf -p 50051:50051 -p 51051:51051 go-proj-rpc:v1

    如果要在后台运行，docker run 加一个 -d参数
    
    采用脚步构建镜像 sh bin/docker-build.sh web

# 版权

    MIT
