#!/bin/sh

set -e

echo "run DB migration"
# 路徑都是 docker image 內的路徑
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# 作用是使用传递给脚本的参数来替换当前的shell进程。这个命令常见于启动脚本的末尾，用于启动一个程序，并且让这个程序接管当前的进程。
exec "$@"