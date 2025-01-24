package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseFile(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, filePath, nil, parser.ParseComments)
}

func ExtractFunctionDocs(file *ast.File) map[string]string {
	docs := make(map[string]string)
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Doc != nil {
				fmt.Println()
				docs[f.Name.Name] = f.Doc.Text()
				if f.Name.Name == "greet" {
					gen(f)
				}
			}
		}
	}
	return docs
}

func main() {
	Comnment()
}

func Comnment() {
	file, err := ParseFile("main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	docs := ExtractFunctionDocs(file)
	fmt.Println("Documentation for function 'greet':", docs["greet"])
}

// @Controller GET /hello
func greet(name string) string {
	return "Hello," + name
}
