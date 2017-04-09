#!/bin/sh

sysctl -w fs.file-max=128000

l1=`cat /proc/sys/fs/file-max`
l2=`ulimit -Sn`
l3=`ulimit -Hn`
echo "file limits: kernel=${l1}, soft ulimit=${l2}, hard ulimit=${l3}"

./dbheroapp_linux
