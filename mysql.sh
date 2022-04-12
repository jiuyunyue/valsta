#!/bin/bash

# bash for macos
docker pull mysql:oracle
docker run --name docker-mysql --privileged=true -d -p 3306:3306  -e MYSQL_ROOT_PASSWORD=123456  mysql:oracle
