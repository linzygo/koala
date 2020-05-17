package main

import (
	"fmt"
	"path/filepath"
)

var confDir = "conf"

func init() {
	gen := &ConfigGenerator{}
	RegisterServerGenerator("config", gen)
	AddServerDir(confDir)
}

// ConfigGenerator 配置文件生成器
type ConfigGenerator struct {
}

// Run 实现Generator接口
func (gen *ConfigGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	fpath := filepath.Join(GetGeneratorRootDir(), confDir, fmt.Sprintf("%s.toml", metaData.PackageName))
	err = SaveCodeToFile(fpath, "config", configTemplate, metaData)
	return
}
