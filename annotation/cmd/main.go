package main

import (
	"fmt"

	"github.com/dulisoft/spirit/annotation"
)

func main() {
	serviceRootPath := "D:/workspace/go/pro/personal/spirit/annotation/examples/bear"
	app := annotation.NewApp(serviceRootPath)
	if err := app.Work(); err != nil {
		fmt.Printf("app work error %v", err)
	}
}
