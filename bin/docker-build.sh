#!/usr/bin/env bash
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

appName=$1

if [ -z $appName ];then
    echo "请指定运行的appName"
    exit 0
fi

if [ ! -d $root_dir/cmd/$appName ];then
    echo "你输入的appName不在$roo_dir/cmd中"
    exit 0
fi

#rpc pb协议自动生成golang pb代码和php代码
if [ "$appName" = "rpc" ];then
    sh $root_dir/bin/go-generate.sh
    sh $roo_dir/bin/php-generate.sh
fi

# 开始构建docker镜像
cd $root_dir

docker build -t go-proj-${appName}:v1 -f $root_dir/${appName}-Dockerfile $root_dir
