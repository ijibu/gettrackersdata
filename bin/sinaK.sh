#!/bin/bash
#######################
##抓取新浪分时K线图脚本
#######################

ymd=`date +%y%m%d`
execdir=/root/go/src/ijibu/trackers/sina
cd $execdir

./sinaK -d=$ymd -n=942 -s=sh -t=k -k=min
./sinaK -d=$ymd -n=1590 -s=sz -t=k -k=min
