FROM golang:1.16.3

ENV PATH="/protoc/bin:${PATH}"

# build golang environment
ENV PROTOC_VER 3.15.6

RUN apt-get update && \
    apt-get -y install autoconf automake cmake make libtool curl g++ zip unzip vim git;

RUN mkdir -p /tmp/temp && wget -q -P /tmp/temp/ https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip && \
    cd /usr && unzip /tmp/temp/protoc-${PROTOC_VER}-linux-x86_64.zip && rm -rf /tmp/temp/*

RUN go env -w GO111MODULE=on; go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct;go env -w CGO_ENABLED=0 &&  \
    go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u google.golang.org/grpc  && \
    go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc  && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway  && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    go get -u github.com/go-playground/validator/v10 && \
    go get  -u  github.com/golang/mock/mockgen  && \
    go get -u github.com/favadi/protoc-go-inject-tag
