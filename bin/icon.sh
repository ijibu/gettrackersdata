#!/bin/bash

##批量替换文件中的单引号为空。直接替换原文件。
sed -i "s/'//g" `find ./ -type f`;

##批量替换文件中的None为0。直接替换原文件。
sed -i "s/None/0/g" `find ./ -type f`;

##批量转换文件编码，因为下载下来的CSV和XLS文件ANSI编码的，必须进行转码。
##在window下查看csv文件的编码为ANSI。在linux下面查看编码，VIM显示为lant1,file命令查看编码为ISO-8008,我擦
##不晓得是啥原因。
##批量进行编码转换，并且覆盖原文件,需要手工创建tmp目录
find ./ -name '*.csv' -exec iconv -f GBK -t UTF8 {} -o ./tmp/{} \;