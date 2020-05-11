#!/bin/bash

function newContainer() {
  docker kill mysql_test
  docker container rm -f mysql_test
  docker run -d --name mysql_test -p 3305:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test_database mysql
  docker cp tables.sql mysql_test:/docker-entrypoint-initdb.d/init.sql
  echo "completed"
}

newContainer