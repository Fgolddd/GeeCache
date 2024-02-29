# 分布式缓存
## 项目介绍
基于HTTP的分布式缓存系统

## 功能介绍
* 支持多种缓存淘汰策略：lru (default), fifo (已拓展)
* 支持并发读写
* 支持一致性哈希
* 分布式的节点访问
* singleflight 防止缓存击穿
* protobuf 通信

## 拓展方向（todo）
1. 将 http 通信改为 rpc 通信提高网络通信效率
2. 加入 etcd 进行分布式节点的监测实现动态管理
3. 加入缓存过期机制
4. 拓展缓存淘汰策略lfu、arc

## 项目结构
``` txt
.
|-- README.md
|-- geecache
|   |-- byteview.go
|   |-- cache.go                     //支持缓存淘汰策略的底层缓存
|   |-- consistenthash               //一致性哈希（负载均衡）
|   |   |-- consistenthash.go
|   |   `-- consistenthash_test.go
|   |-- geecache.go                  //对底层缓存的封装
|   |-- geecache_test.go
|   |-- geecachepb                   //切换protobuf
|   |   |-- geecachepb.pb.go
|   |   `-- geecachepb.proto
|   |-- go.mod
|   |-- go.sum
|   |-- http.go                      //HTTP通信
|   |-- peers.go   
|   |-- singleflight                 //并发请求处理优化（协程编排）
|   |   |-- singleflight.go
|   |   `-- singleflight_test.go
|   `-- tactics                      //缓存淘汰策略
|       |-- fifo                     //FIFO淘汰策略
|       |   |-- fifo.go
|       |   `-- fifo_test.go
|       `-- lru                      //LRU淘汰策略
|           |-- lru.go
|           `-- lru_test.go
|-- go.mod
|-- go.sum
|-- main.go
`-- run.sh                           //测试脚本
```

