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

#自动版本号生成
if [ -f $root_dir/bin/pprof-check-version.sh ];then
    echo "生成自动版本号"
    sh $root_dir/bin/pprof-check-version.sh
fi

if [ "$appName" = "web" ];then
    sh $root_dir/bin/web-check-version.sh
fi

#rpc pb协议自动生成golang pb代码和php代码
if [ "$appName" = "rpc" ];then
    sh $root_dir/bin/pb-generate.sh
fi

# 开始构建docker镜像
cd $root_dir

docker build -t go-proj-${appName}:v1 -f $root_dir/${appName}-Dockerfile $root_dir
