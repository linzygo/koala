package main

import (
	"koala/util"
	"path/filepath"
)

var mainDir = "main"

func init() {
	gen := &ServerMainGenerator{}
	RegisterServerGenerator("main", gen)
	AddServerDir(mainDir)
}

// ServerMainGenerator 服务端main函数代码生成器
type ServerMainGenerator struct {
}

// Run 实现Generator接口
func (gen *ServerMainGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	fpath := filepath.Join(GetGeneratorRootDir(), mainDir, "main.go")
	if util.PathIsExist(fpath) {
		return
	}
	err = SaveCodeToFile(fpath, "main", mainTemplate, metaData)
	return
}
