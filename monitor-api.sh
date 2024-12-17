#!/bin/bash

start() {
 cd /opt/monitor_api
 ./counter --start > /dev/null 2> /dev/null &
}

stop() {
 kill -s SIGKILL `ps -ef | grep -i "./counter --start" | awk '{print $2;}'`
}

case $1 in
 start|stop) "$1" ;;
esac
echo ""
