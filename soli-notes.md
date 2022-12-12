> 看官方github的issue区bug数，评估稳定性
## 过期特性
## github issues
[Why we couldn't set Key-Value larger than 1/1024 of cache size?](https://github.com/coocood/freecache/issues/28)
## pprof
 go tool pprof -seconds 30 -http 127.0.0.1:6061 http://127.0.0.1:6060/debug/pprof/profile
 go tool pprof -seconds 30 -http 127.0.0.1:6061 http://127.0.0.1:6060/debug/pprof/goroutine

 lsof -i:18000|head -2

指定端口连接数：netstat -nat | grep -i "8080" | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'