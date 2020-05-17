package main

// Generator 代码生成器接口
type Generator interface {
	Run(opt *Option, metaData *ServiceMetaData) (err error)
}

// RPCMetaData proto服务的rpc信息
// 字段
//   Name: rpc接口的名称
//   RequestType: rpc接口的参数类型
//   ReturnsType: rpc接口的返回值类型
type RPCMetaData struct {
	Name        string
	RequestType string
	ReturnsType string
}

// ServiceMetaData proto元信息
// 字段
//   PackageName: 对应proto文件的package
//   PackagePrifix: 用于模板import package时能正确指定路径, 命令行参数传入
//   ServiceName: 对应proto文件的service名称
//   RPC: 包含proto文件里面的service的每个rpc信息
type ServiceMetaData struct {
	PackageName   string
	PackagePrefix string
	ServiceName   string
	RPCs          []RPCMetaData
}
