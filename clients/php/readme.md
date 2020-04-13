# php grpc运行

    composer install
    php hello_client.php

# 关于是否需要使用protobuf.so

    对于php7.0+，protoc3可以安装php protobuf拓展
    vim php.ini
    ;不一定要安装,一般建议用protobuf拓展比较好
    extension=protobuf.so
    extension=grpc.so

    这个时候，可以去掉composer2.json中的
    "google/protobuf": "^3.8"
    mv composer2.json composer.json
    然后composer update

    对于不支持protobuf的php版本，可以用用composer2.json 替换为composer.json

# composer镜像设置

    采用composer config -g repo.packagist composer https://mirrors.aliyun.com/composer/
