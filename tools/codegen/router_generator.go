package main

import (
	"path/filepath"
)

var routerDir = "router"

func init() {
	gen := &RouterGenerator{}
	RegisterServerGenerator("router", gen)
	AddServerDir(routerDir)
}

// RouterGenerator 路由
type RouterGenerator struct {
}

// Run 实现Generator接口
func (gen *RouterGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	fpath := filepath.Join(GetGeneratorRootDir(), routerDir, "router.go")
	err = SaveCodeToFile(fpath, "router", routerTemplate, metaData)
	return
}
