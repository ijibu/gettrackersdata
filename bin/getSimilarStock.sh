#!/bin/bash
#	获取相似股票曲线脚本
#		该脚本首先把每只股票和其它股票的相似性计算出来，保存在一个文件中，
#		然后把该文件进行排序，找到股票相似性最高的股票。
for line in $(cat /root/go/src/github.com/ijibu/gettrackersdata/ini/shang_new.ini)
do
	echo ${line}
	cd /root/go/src/github.com/ijibu/gettrackersdata
	./getSimilarStock -n=${line} >> ijibu.log
done
