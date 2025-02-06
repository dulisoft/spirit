package annotation

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"log"
	"os"
	"text/template"
)

// CommonWriter 基础写文件组件
type CommonWriter struct {
	Pkg      string
	FilePath string
	Position string
	GenDels  []*ast.GenDecl
	FuncDels []*ast.FuncDecl
}

func (c *CommonWriter) Write() error {
	// 创建一个新的文件集合，用于存储我们的Go代码
	fset := token.NewFileSet()

	// 创建一个ast.File, 一个包含包声明和文件级声明的单元
	file := &ast.File{
		Name:  ast.NewIdent(c.Pkg),
		Decls: []ast.Decl{},
	}

	// 将变量声明添加到文件的声明列表中
	for i := range c.GenDels {
		file.Decls = append(file.Decls, c.GenDels[i])
	}
	for i := range c.FuncDels {
		file.Decls = append(file.Decls, c.FuncDels[i])
	}

	// 使用go/format包格式化AST
	buf := bytes.NewBuffer([]byte(""))
	var err error
	if err = format.Node(buf, fset, file); err != nil {
		fmt.Println("Error formatting AST:", err)
		return err
	}

	// 将格式化后的代码写入文件
	if err := os.WriteFile(c.FilePath, buf.Bytes(), 0666); err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}

func WirteTemplate(tmpl *template.Template, app any, filePath string) error {
	destFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Open Files Error:", err)
		return err
	}
	return tmpl.Execute(destFile, app)
}

func WriteAtNextLine(filePath string, pos int, content string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	sb := bytes.NewBuffer(nil)
	sb.Write(fileContent[:pos])
	sb.WriteByte(byte('\n'))
	sb.Write([]byte(content))
	sb.Write(fileContent[pos:])

	newContent := sb.Bytes()
	if _, err := file.WriteAt(newContent, 0); err != nil {
		return err
	}
	return nil
}
