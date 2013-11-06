#1/bin/bash
if test "$1" = ""; then
	ymd=`date -d yesterday +%Y%m%d`
else
	ymd=$1
fi
cd /root/go/src/ijibu/gettarckers
./163 -d=$ymd -n=1 -s=sh -t=cjmx
./163 -d=$ymd -n=1 -s=sh -t=chddata
./163 -d=$ymd -n=1 -s=sh -t=lszjlx
./163 -d=$ymd -n=1 -s=sz -t=cjmx
./163 -d=$ymd -n=1 -s=sz -t=chddata
./163 -d=$ymd -n=1 -s=sz -t=lszjlx
exit 0
