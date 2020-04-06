# centos7 protoc工具安装

    1、下载https://github.com/protocolbuffers/protobuf/archive/v3.11.4.tar.gz
    cd /usr/local/src
    sudo wget https://github.com/protocolbuffers/protobuf/archive/v3.11.4.tar.gz
    2、开始安装
    sudo mv v3.11.4.tar.gz protobuf-3.11.4.tar.gz
    sudo tar zxvf protobuf-3.11.4.tar.gz
    cd protobuf-3.11.4
    sudo yum install gcc-c++ cmake libtool
    ]$ sudo mkdir /usr/local/protobuf

    需要编译, 在新版的 PB 源码中，是不包含 .configure 文件的，需要生成，此时先执行 ./autogen.sh 
    脚本说明如下：
    # Run this script to generate the configure script and other files that will
    # be included in the distribution. These files are not checked in because they
    # are automatically generated.

    此时生成了 .configure 文件，可以开始编译了
    sudo ./configure --prefix=/usr/local/protobuf
    sudo make && make install

    安装完成后,查看版本:
    $ cd /usr/local/protobuf/bin
    $ ./protoc --version
    libprotoc 3.11.4
     
    建立软链接
    
    $ sudo ln -s /usr/local/protobuf/bin/protoc /usr/bin/protoc

# go protoc工具安装

    go get -u github.com/golang/protobuf/proto
    
    go get -u github.com/golang/protobuf/protoc-gen-go

    go get -u google.golang.org/grpc

    cd $GOPATH/pkg/mod/github.com/golang/protobuf@v1.3.5/protoc-gen-go
    go install

    cd $GOPATH/pkg/mod/github.com/golang/protobuf@v1.3.5/proto
    go install

