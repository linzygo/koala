package main

import (
	"fmt"
	"koala/util"
	"path/filepath"
)

var controllerDir = "controller"

func init() {
	gen := &ControllerGenerator{}
	RegisterServerGenerator("controller", gen)
	AddServerDir(controllerDir)
}

// ControllerGenerator 路由代码生成器
type ControllerGenerator struct {
}

// RPCMeta 封装数据传给template
type RPCMeta struct {
	RPCMetaData
	PackagePrefix string
	PackageName   string
}

// Run 实现Generator接口
func (gen *ControllerGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	rpcMeta := &RPCMeta{
		PackagePrefix: metaData.PackagePrefix,
		PackageName:   metaData.PackageName,
	}

	for _, RPC := range metaData.RPCs {
		rpcMeta.RPCMetaData = RPC
		err = gen.generateFile(rpcMeta)
		if err != nil {
			return
		}
	}

	return
}

// generateFile 生成代码文件
// 参数
//   metaData: RPC信息
// 返回值
//   err: error
func (gen *ControllerGenerator) generateFile(metaData *RPCMeta) (err error) {
	fpath := filepath.Join(GetGeneratorRootDir(), controllerDir, fmt.Sprintf("%s.go", ToUnderscoreString(metaData.RPCMetaData.Name)))
	if util.PathIsExist(fpath) {
		return
	}
	err = SaveCodeToFile(fpath, "controller", controllerTemplate, metaData)
	return
}
