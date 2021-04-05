FROM golang:1.16.3

ENV PATH="/protoc/bin:${PATH}"

# build golang environment
ENV PROTOC_VER 3.15.6

RUN apt-get update && apt-get install -y apt-utils zip unzip curl  vim cmake make;

RUN wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip \  
    && mkdir -p /protoc && unzip protoc.zip -d /protoc

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct  && \
    go get -u google.golang.org/grpc \
    go get -u github.com/golang/protobuf/proto \
    go get -u github.com/golang/protobuf/protoc-gen-go \
    go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    go get  -u  github.com/golang/mock/mockgen \
    go get -u   github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    go get  -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    go get -u  github.com/mwitkow/go-proto-validators/protoc-gen-govalidators 
