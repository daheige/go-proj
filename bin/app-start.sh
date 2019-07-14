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

#构建工程
sh $root_dir/bin/web-init.sh $appName

#开始运行
$root_dir/bin/go-gen/$appName -config_dir=$root_dir -log_dir=$root_dir/logs/$appName
