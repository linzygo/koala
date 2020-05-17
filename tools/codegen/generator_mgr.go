package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
)

var (
	genMgr = &GeneratorMgr{
		genCommonMap: make(map[string]Generator),
		genServerMap: make(map[string]Generator),
		genClientMap: make(map[string]Generator),
	}
)

// GeneratorMgr 管理所有生成代码的工具
type GeneratorMgr struct {
	// 服务端、客户端共用的代码生成器
	genCommonMap map[string]Generator
	// 服务端代码生成器
	genServerMap map[string]Generator
	// 客户端代码生成器
	genClientMap map[string]Generator
	// proto文件元信息
	meta ServiceMetaData
	// 根目录
	rootDir string
	// 公共代码目录
	commonDirs []string
	// 服务端代码目录
	serverDirs []string
	// 客户端代码目录
	clientDirs []string
}

// Run 生成各个代码
// 参数
//   opt: 命令行参数
// 返回值
//   err: error
func (mgr *GeneratorMgr) Run(opt *Option) (err error) {
	err = mgr.parseProto(opt.ProtoFile)
	if err != nil {
		return
	}

	// 如果没有指定project名称，那么使用解析得到的package名称
	if len(opt.ProjectName) == 0 {
		opt.ProjectName = mgr.meta.PackageName
	}
	mgr.meta.PackagePrefix = path.Join(opt.PackagePrefix, opt.ProjectPrefix, opt.ProjectName)
	mgr.meta.PackagePrefix = strings.ReplaceAll(mgr.meta.PackagePrefix, "\\", "/")
	fmt.Printf("meta.PackagePrefix=%s\n", mgr.meta.PackagePrefix)

	err = mgr.initDir(opt)
	if err != nil {
		return
	}

	err = mgr.runGenerator(opt, mgr.genCommonMap)
	if err != nil {
		return
	}
	if opt.ServerCode {
		err = mgr.runGenerator(opt, mgr.genServerMap)
		if err != nil {
			return
		}
	}
	if opt.ClientCode {
		err = mgr.runGenerator(opt, mgr.genClientMap)
		if err != nil {
			return
		}
	}
	return
}

// runGenerator 执行各个map的代码生成器
// 参数
//   opt: 命令行参数
//   genMap: 对应的map
// 返回值
//   err: error
func (mgr *GeneratorMgr) runGenerator(opt *Option, genMap map[string]Generator) (err error) {
	for _, gen := range genMap {
		err = gen.Run(opt, &mgr.meta)
		if err != nil {
			return
		}
	}
	return
}

// parseProto 解析文件
// 参数
//   path: 文件路径
// 返回值
//   err: error
func (mgr *GeneratorMgr) parseProto(path string) (err error) {
	reader, err := os.Open(path)
	if err != nil {
		fmt.Printf("打开proto文件[%s]失败, err=%v\n", path, err)
		return
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	protoData, err := parser.Parse()
	if err != nil {
		fmt.Printf("解析proto文件失败, err=%v\n", err)
		return
	}

	proto.Walk(protoData, proto.WithService(mgr.handleService), proto.WithRPC(mgr.handleRPC), proto.WithPackage(mgr.handlePackage))

	return
}

// handleService 用于处理解析proto文件得到的service信息
func (mgr *GeneratorMgr) handleService(service *proto.Service) {
	mgr.meta.ServiceName = service.Name
}

// handleRPC 用于处理解析proto文件得到的RPC信息
// 参数
//   RPC: 解析出来的RPC信息
func (mgr *GeneratorMgr) handleRPC(RPC *proto.RPC) {
	data := RPCMetaData{
		Name:        RPC.Name,
		RequestType: RPC.RequestType,
		ReturnsType: RPC.ReturnsType,
	}
	mgr.meta.RPCs = append(mgr.meta.RPCs, data)
}

// handlePackage 用于处理解析proto文件得到的Package信息
// 参数
//   Package: 解析出来的Package信息
func (mgr *GeneratorMgr) handlePackage(Package *proto.Package) {
	mgr.meta.PackageName = Package.Name
}

// initDir 初始化代码目录
// 参数
//   opt: 命令行参数
// 返回值
//   err: error
func (mgr *GeneratorMgr) initDir(opt *Option) (err error) {
	dirpath := path.Join(opt.OutputDir, opt.ProjectPrefix, opt.ProjectName)
	dirpath = filepath.FromSlash(dirpath)
	err = os.MkdirAll(dirpath, os.ModeDir)

	if err != nil && !os.IsExist(err) {
		fmt.Printf("创建项目目录[%s]失败, err=%v\n", dirpath, err)
		return
	}

	err = nil
	mgr.rootDir = dirpath

	err = mgr.createDirs(mgr.commonDirs...)
	if err != nil {
		return
	}

	if opt.ServerCode {
		err = mgr.createDirs(mgr.serverDirs...)
		if err != nil {
			return
		}
	}

	if opt.ClientCode {
		err = mgr.createDirs(mgr.clientDirs...)
		if err != nil {
			return
		}
	}

	return
}

// createDirs 创建目录
// 参数
//   dirs: 要创建的目录不定参(多个同级目录)
// 返回值
//   err: error
func (mgr *GeneratorMgr) createDirs(dirs ...string) (err error) {
	if len(dirs) == 0 {
		return
	}
	for _, dir := range dirs {
		fullpath := filepath.Join(mgr.rootDir, dir)
		err = os.Mkdir(fullpath, os.ModeDir)
		if err != nil {
			if !os.IsExist(err) {
				fmt.Printf("创建目录[%s]失败, err=%v\n", fullpath, err)
				return
			}
			err = nil
		}
	}
	return
}

// RegisterCommonGenerator 公共代码生成器注册
// 参数
//   name: 代码生成器名称
//   gen: 代码生成器
// 返回值
//   err: error
func RegisterCommonGenerator(name string, gen Generator) (err error) {
	if _, ok := genMgr.genCommonMap[name]; ok {
		err = fmt.Errorf("公共代码生成器[%s]已经存在", name)
		return
	}
	genMgr.genCommonMap[name] = gen
	return
}

// RegisterServerGenerator 服务端代码生成器注册
// 参数
//   name: 代码生成器名称
//   gen: 代码生成器
// 返回值
//   err: error
func RegisterServerGenerator(name string, gen Generator) (err error) {
	if _, ok := genMgr.genServerMap[name]; ok {
		err = fmt.Errorf("服务端代码生成器[%s]已经存在", name)
		return
	}
	genMgr.genServerMap[name] = gen
	return
}

// RegisterClientGenerator 客户端代码生成器注册
// 参数
//   name: 代码生成器名称
//   gen: 代码生成器
// 返回值
//   err: error
func RegisterClientGenerator(name string, gen Generator) (err error) {
	if _, ok := genMgr.genClientMap[name]; ok {
		err = fmt.Errorf("客户端代码生成器[%s]已经存在", name)
		return
	}
	genMgr.genClientMap[name] = gen
	return
}

// GetGeneratorRootDir 获取生成的代码根目录
// 返回值
//   rootDir, 根目录
func GetGeneratorRootDir() string {
	return genMgr.rootDir
}

// AddCommonDir 添加要生成公共代码的目录
// 参数
//   dir: 目录名称
func AddCommonDir(dir string) {
	if len(dir) != 0 {
		genMgr.commonDirs = append(genMgr.commonDirs, dir)
	}
}

// AddServerDir 添加服务端代码生成目录
// 参数
//   dir: 目录名称
func AddServerDir(dir string) {
	if len(dir) != 0 {
		genMgr.serverDirs = append(genMgr.serverDirs, dir)
	}
}

// AddClientDir 添加客户端代码生成目录
// 参数
//   dir: 目录名称
func AddClientDir(dir string) {
	if len(dir) != 0 {
		genMgr.clientDirs = append(genMgr.clientDirs, dir)
	}
}
