package main

// Option 命令行参数
type Option struct {
	OutputDir     string // 输出目录
	ProtoFile     string // protobuf的idl文件
	ProjectName   string // 项目名称, 如果用户没有指定，则使用proto文件的package名称
	ProjectPrefix string // 项目前缀, 生成的代码放在"ProjectPrefix/ProjectName/"下
	PackagePrefix string // 包前缀，生成的代码import自己生成的包，import "PackagePrefix/ProjectPrefix/ProjectName/xxx"
	ServerCode    bool   // 服务端代码
	ClientCode    bool   // 客户端代码
}
