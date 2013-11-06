#!/bin/bash
if test "$1" = ""; then
	dateYear=`date +%Y`
	dateMonth=`date +%Y`
	dateDay=`date +%d`
	dateMonth=`expr  $dateMonth - 1`
	dateDay=`expr  $dateDay - 1`
	startDate=`echo "${dateYear}-${dateMonth}-${dateDay}"`	##计算时间，主要是yahoo接口的特殊性。	
else
	startDate=$1
fi
if test "$2" = ""; then
	dateYear1=`date +%Y`
	dateMonth1=`date +%Y`
	dateDay1=`date +%d`
	dateMonth1=`expr  $dateMonth1 - 1`
	dateDay1=`expr  $dateDay1 - 1`
	endDate=`echo "${dateYear1}-${dateMonth1}-${dateDay1}"`	##计算时间，主要是yahoo接口的特殊性。	
else
	endDate=$2
fi

cd /root/go/src/ijibu/gettrackersdata
./yahoo -d=$startDate -e=$endDate -n=941 -s=sh
./yahoo -d=$startDate -e=$endDate -n=1589 -s=sz
exit 0
