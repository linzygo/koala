package main

var mainTemplate = `
package main

import (
	"koala/server"
	"{{.PackagePrefix}}/{{.PackageName}}"
	"{{.PackagePrefix}}/router"
)

func main() {
	err := server.InitServer("{{.PackageName}}")
	if err != nil {
		return
	}

	{{.PackageName}}.Register{{.ServiceName}}Server(server.GetServer(), &router.ServerRouter{})
	server.Run()
}
`
