package main

var configTemplate = `
[service]
name="{{.PackageName}}"
port=8081

[log]
type="file" #日志类型, [console|file|console,file]
level="info" #日志级别, [debug|trace|info|access|warn|error]
path="./log" #日志文件根目录, 只有type为file或者console,file的时候才起作用
file_size=10485760 #日志文件最大10M
chan_size=50000 #日志队列大小,队列太小,打印日志太多可能导致部分丢失

[register]
switch_on=true
register_name="etcd"
register_addrs=["0.0.0.0:2379"]
register_path="/linzygo.com/koala/"
timeout=1
heartbeat=10

[prometheus]
port=8080
switch_on=true

[limit]
qps=50000
switch_on=true

[trace]
switch_on=true
report_addr="http://localhost:9411/api/v1/spans"
sample_type="const"
sample_rate=1
`
