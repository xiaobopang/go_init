#!/bin/bash
# @Author: pangxiaobo
# @Date:   2018-11-06 14:34:30
# @Last Modified by:   pangxiaobo
# @Last Modified time: 2018-11-06 14:36:30


case $1 in 
	start)
		nohup ./go_init 2>&1 >> info.log 2>&1 /dev/null &
		echo "服务已启动..."
		sleep 1
	;;
	stop)
		killall go_init
		echo "服务已停止..."
		sleep 1
	;;
	restart)
		killall go_init
		sleep 1
		nohup ./go_init 2>&1 >> info.log 2>&1 /dev/null &
		echo "服务已重启..."
		sleep 1
	;;
	*) 
		echo "$0 {start|stop|restart}"
		exit 4
	;;
esac