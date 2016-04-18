#!/bin/sh
#
# Please see:
# http://www.cnblogs.com/getong/archive/2013/04/01/2993139.html
#

go test -test.benchtime=2s -test.bench=. -test.benchmem > redis-go-driver-benchmark.log

for i in AlphazeroRedis GaryburdRedigo GosexyRedis FishRedis Simonz05Godis 
do
  grep $i redis-go-driver-benchmark.log | awk '{print $3}' > $i.tmp
  grep $i redis-go-driver-benchmark.log | awk '{print $5}' > $i.mem.tmp
  grep $i redis-go-driver-benchmark.log | awk '{print $7}' > $i.alloc.tmp
done

R --no-save < go-redis-getongs-data.R > /dev/null
R --no-save < go-redis-mem-data.R > /dev/null
R --no-save < go-redis-alloc-data.R > /dev/null
rm -f *.tmp
