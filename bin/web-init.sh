#!/usr/bin/env bash
#初始化目录权限
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

appName=$1
if [ -z $appName ];then
    echo "请指定构建的appName"
    exit 0
fi

if [ ! -d $root_dir/cmd/$appName ];then
    echo "你输入的appName不在$roo_dir/cmd中"
    exit 0
fi

mkdir -p $root_dir/logs/$appName
chmod 777 -R $root_dir/logs/$appName
echo "初始化成功！"

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

if [ -d $root_dir/bin/go-gen/$appName ];then
    rm -rf $root_dir/bin/go-gen/$appName
fi

mkdir -p $root_dir/bin/go-gen/$appName

echo "开始构建web二进制文件"
cd $root_dir/cmd/$appName
go build -o $root_dir/bin/go-gen/$appName/$appName

#清除cmd/下面由于go build生成的二进制文件
cd $root_dir/cmd/$appName
go clean .

echo "构建$appName成功！"

exit 0
