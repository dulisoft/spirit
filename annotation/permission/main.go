package main

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
)

var args *CmdArg

type CmdArg struct {
	Path string `json:"path"`
}

// ExampleController
// @Permission 读取   管理逻辑视图   全部&本部门
// @Router  /api/data-view/v1/form_view   GET
func ExampleController() {
}

func main() {
	parser, err := NewProjectParser("main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	ans := parser.ReadAnnotationSlice()
	fmt.Printf("%s\n", lo.T2(json.Marshal(ans)).A)
}
