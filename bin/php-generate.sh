#!/usr/bin/env bash
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

protoExec=$(which "protoc")
if [ -z $protoExec ]; then
    echo 'Please install protoc!'
    echo "Please look readme.md to install proto3"
    echo "if you use centos7,please look docs/centos7-protoc-install.md"
    exit 0
fi

# grpc php 工具
grpc_php_plugin=$(which "grpc_php_plugin")
if [ -z $grpc_php_plugin ]; then
    echo 'Please install grpc_php_plugin!'
    echo "Please look docs/php_grpc.md to install grpc_php_plugin tool"
    #exit
fi

# protos协议文件目录
protos_dir=$root_dir/protos

#如果想放在其他目录，可以更改
php_client_dir=$root_dir/clients/php

mkdir -p $php_client_dir

#删除老的pb生成的php文件
rm -rf $php_client_dir/App

echo "\n\033[0;32mGenerating codes...\033[39;49;0m\n"

echo "generating php stubs..."

#生成php pb文件
cd $protos_dir

for file in $protos_dir/*.proto; do
    echo "generating php stubs from: $file"
    $protoExec --proto_path=$protos_dir --php_out=$root_dir/clients/php/ --grpc_out=$root_dir/clients/php/ --plugin=protoc-gen-grpc=$grpc_php_plugin $file
    echo "\t\033[0;32m[DONE]\033[39;49;0m\n"
done

#将php grpc生成的文件归档处理
mv $root_dir/clients/php/GPBMetadata $php_client_dir/App/Grpc/

#为了避免composer install grpc protobuf的命名空间GPBMetadata冲突
#这里需要把GPBMetadata加上命名空间
#在GPBMetadata前面加上命名空间App\Grpc
sed -i 's/GPBMetadata/App\\\Grpc\\\GPBMetadata/g' `grep GPBMetadata -rl $php_client_dir/App/Grpc/`

echo "\n\033[0;32mGenerate codes successfully!\033[39;49;0m\n"

exit 0
