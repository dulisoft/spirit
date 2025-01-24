package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

func gen(fdel *ast.FuncDecl) {
	// 创建一个新的文件集合，用于存储我们的Go代码
	fset := token.NewFileSet()

	// 创建一个ast.File, 一个包含包声明和文件级声明的单元
	file := &ast.File{
		Name:  ast.NewIdent("gcode"),
		Decls: []ast.Decl{},
	}

	// 创建一个变量声明
	varDecl := &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names:  []*ast.Ident{ast.NewIdent("x")},
				Values: []ast.Expr{ast.NewIdent("42")},
			},
		},
	}

	// 将变量声明添加到文件的声明列表中
	file.Decls = append(file.Decls, varDecl, fdel)

	// 使用go/format包格式化AST
	buf := bytes.NewBuffer([]byte(""))
	var err error
	if err = format.Node(buf, fset, file); err != nil {
		fmt.Println("Error formatting AST:", err)
		return
	}

	// 将格式化后的代码写入文件
	if err := os.WriteFile("gcode/output.go", buf.Bytes(), 0666); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Code generated successfully!")
}
