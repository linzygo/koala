# koala
go语言微服务框架，代码生成工具，解决一些共性问题  

client - rpc客户端，封装一些公共的处理，封装每个rpc的扩展处理，方便扩展  
config - 服务端与客户端的共用配置, 配置使用的是toml  
errno - 一些错误定义  
loadbalance - 负载均衡算法，客户端调用服务时选择服务节点  
logger - 日志组件  
meta - 中间件元数据定义，通过Context在各个中间件传递  
middleware - 中间件  
registry - 注册插件，用于服务注册&发现  
server - rpc服务端，封装一些公共的处理，方便扩展  
tools  
  -codegen - 代码生成工具  
util - 辅助功能  
