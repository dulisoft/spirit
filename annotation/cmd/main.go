package main

import (
	"fmt"

	"github.com/dulisoft/spirit/annotation"
)

func main() {
	goPath := "D:/workspace/go/pro/personal/spirit/annotation/examples/bear/controller/example/router.go"
	routers, err := annotation.PareserAnnotations(goPath)
	if err != nil {
		fmt.Printf("error :%v", err)
	}
	controller := annotation.NewController(routers)
	app := annotation.NewApp([]annotation.Controller{*controller})
	if err := app.RegisterRouters(); err != nil {
		fmt.Printf("error :%v", err)
	}
}
