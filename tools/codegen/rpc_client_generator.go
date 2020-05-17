package main

import (
	"path/filepath"
	"unicode"
)

const (
	rpcClientDir     = "client"
	rpcClientGenName = "rpc_client"
)

func init() {
	gen := &RPCClientGenerator{}
	RegisterClientGenerator(rpcClientGenName, gen)
	AddClientDir(rpcClientDir)
}

// RPCClientMeta 封装数据传递给模板
type RPCClientMeta struct {
	*ServiceMetaData
	ClientName string
}

// RPCClientGenerator rpc客户端代码生成器
type RPCClientGenerator struct {
}

// Run 实现Generator接口
func (gen *RPCClientGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	clientName := []rune(metaData.PackageName)
	clientName[0] = unicode.ToUpper(clientName[0])
	meta := &RPCClientMeta{
		ServiceMetaData: metaData,
		ClientName:      string(clientName),
	}
	fpath := filepath.Join(GetGeneratorRootDir(), rpcClientDir, "client.go")
	err = SaveCodeToFile(fpath, rpcClientGenName, rpcClientTemplate, meta)
	return
}
