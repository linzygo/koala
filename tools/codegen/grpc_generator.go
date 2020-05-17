package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	gen := &GRPCGenerator{}
	RegisterCommonGenerator("grpc", gen)
}

// GRPCGenerator grpc代码生成器
type GRPCGenerator struct {
}

// Run 实现Generator接口
func (gen *GRPCGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	dir := filepath.Join(GetGeneratorRootDir(), metaData.PackageName)
	err = os.Mkdir(dir, os.ModeDir)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("创建rpc代码目录[%s]失败, err=%v\n", dir, err)
		return
	}
	err = nil

	outputParams := fmt.Sprintf("plugins=grpc:%s", dir)
	params := []string{
		"--go_out",
		outputParams,
		opt.ProtoFile,
		"-I",
		filepath.Dir(opt.ProtoFile),
	}
	fmt.Printf("params=%#v\n", params)
	cmd := exec.Command("protoc", params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("生成rpc代码失败, err=%v\n", err)
	}
	return
}
