package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	var opt Option
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "out",
			Usage:       "输出路径, 会在此路径下创建相应项目存放生成的代码。",
			Destination: &opt.OutputDir,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "f",
			Usage:       "protobuf的proto文件, 用于生成rpc代码。",
			Destination: &opt.ProtoFile,
			Required:    true,
		},
		&cli.StringFlag{
			Name: "p",
			Usage: `项目名称,
					如果未指定, 则项目名称默认为f选项指定的proto文件里面的package名称。`,
			Destination: &opt.ProjectName,
		},
		&cli.StringFlag{
			Name: "pp",
			Usage: `项目前缀,
					如果指定此选项, 则生成的代码放在"out路径/项目前缀/项目名/"目录下,
					生成的代码import生成的package时将是import "项目前缀/项目名/xxx";
					如果未指定此选项, 则生成的代码放在"out路径/项目名/"目录下,
					生成的代码import生成的package时将是import "项目名/xxx"。
					注意import package收pkgp参数影响，具体看pkgp参数。
				   `,
			Destination: &opt.ProjectPrefix,
		},
		&cli.StringFlag{
			Name: "pkgp",
			Usage: `包前缀,
					如果指定此选项, 则生成的代码import生成的package时将是
					"包前缀/项目前缀/项目名/xxx"或"包前缀/项目名/xxx";
					如果没指定此选项, 则生成的代码import生成的package时将是
					"项目前缀/项目名/xxx"或"项目名/xxx"。
					`,
			Destination: &opt.PackagePrefix,
		},
		&cli.BoolFlag{
			Name:        "s",
			Usage:       "如果指定此选项，则生成服务端代码，如果没有指定此选项也没有指定c选项，则默认使用此选项。",
			Destination: &opt.ServerCode,
		},
		&cli.BoolFlag{
			Name:        "c",
			Usage:       "如果指定此选项，则生成客户端代码，必须显式指定此选项，否则不会生成客户端代码。",
			Destination: &opt.ClientCode,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		if !opt.ServerCode && !opt.ClientCode {
			opt.ServerCode = true
		}

		genMgr.Run(&opt)
		fmt.Println("代码生成完成")

		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
