#!/bin/bash

# bash for macos
docker pull mysql:oracle

# shellcheck disable=SC2046
if [ $(docker ps -aqf name=docker-mysql ) ];then
echo "start "
docker start $(docker ps -aqf name=docker-mysql)
else
echo "run "
docker run --name docker-mysql --privileged=true -d -p 3306:3306  -e MYSQL_ROOT_PASSWORD=123456  mysql:oracle
fi

mkdir cache
sudo chmod 777 cache

