package annotation

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

const ANNOTATION_PREFIX = "//"

var (
	METHOD    = 1
	FUNCTIONS = 2
)

var anoRegexp, _ = regexp.Compile("@\\w+")

// IdentType 类型描述
type IdentType struct {
	Path string
	Name string
}

// AnnoWork 注解工作方法
type AnnoWork func(anno *Annotation) error

// Annotation 注解信息，可以方便包含所有的信息，方便生成代码
type Annotation struct {
	Name  string   //注解的名称
	Kind  int      //被注解对象类型
	Props []string //注解属性
	obj   Object   //被注解对象信息
	pos   Position //位置信息
}

type Position struct {
	DocBegin int //注释之上的位置
	DeclEnd  int //声明后面的位置
}

// Object 被注解对象
type Object struct {
	Name     string    //名称
	PkgPath  string    //包路径
	PkgName  string    //包名称
	FuncSign *FuncSign //方法签名
}

func (obj Object) Import() string {
	if strings.HasSuffix(obj.PkgPath, obj.PkgName) {
		return obj.PkgPath
	}
	return obj.PkgName + " " + obj.PkgPath
}

func ParseGoFileDecls(filePath string) ([]*Annotation, map[string]*IdentType, error) {
	astFile, err := ParseFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	importDict := make(map[string]*IdentType)
	annotations := make([]*Annotation, 0)
	for _, decl := range astFile.Decls {
		switch f := decl.(type) {
		case *ast.FuncDecl:
			annotations = append(annotations, parseFuncDecl(f, astFile.Name.Name)...)
		case *ast.GenDecl:
			annotations = append(annotations, parseGenDecl(f)...)
			if f.Tok == token.IMPORT {
				//import的理论上只有一个地方
				importDict, err = parserImports(f)
				if err != nil {
					fmt.Printf("parse import error %v", err.Error())
				}
			}
		default:
			panic(fmt.Sprintf("sytax error,position: %v", f.Pos()))
		}
	}
	return annotations, importDict, nil
}

// parseFuncDecl 解析方法
func parseFuncDecl(funcDecl *ast.FuncDecl, pkgName string) (ans []*Annotation) {
	if funcDecl.Doc == nil {
		return
	}
	comment := funcDecl.Doc.Text()
	if !strings.Contains(comment, ANNOTATION_PREFIX) {
		return
	}
	for _, cmt := range funcDecl.Doc.List {
		cmtLine := cmt.Text
		if !strings.Contains(cmtLine, ANNOTATION_PREFIX) {
			continue
		}
		cmts := parserCmt(cmtLine)
		if len(cmts) <= 0 {
			continue
		}
		anno := &Annotation{
			Name:  strings.TrimPrefix(cmts[0], ANNOTATION_PREFIX),
			Kind:  METHOD,
			Props: cmts[1:],
			obj: Object{
				Name:     funcDecl.Name.Name,
				PkgName:  pkgName,
				FuncSign: parseFuncSign(funcDecl),
			},
			pos: Position{
				DocBegin: int(funcDecl.Doc.Pos()),
				DeclEnd:  int(funcDecl.Body.End()),
			},
		}
		//判断是方法还是函数
		if funcDecl.Recv == nil {
			anno.Kind = FUNCTIONS
		}
		ans = append(ans, anno)
	}
	return ans
}

// parseGenDecl 解析变量结构体等
func parseGenDecl(genDecl *ast.GenDecl) (ans []*Annotation) {
	if genDecl.Doc == nil {
		return
	}
	comment := genDecl.Doc.Text()
	if !strings.Contains(comment, ANNOTATION_PREFIX) {
		return
	}
	for _, cmt := range genDecl.Doc.List {
		cmtLine := cmt.Text
		if strings.Contains(cmtLine, ANNOTATION_PREFIX) {
			continue
		}
		cmts := parserCmt(cmtLine)
		if len(cmts) <= 0 {
			continue
		}
		ano := &Annotation{
			Name:  strings.TrimPrefix(cmts[0], ANNOTATION_PREFIX),
			Kind:  0,
			Props: cmts[1:],
			obj: Object{
				Name:    genDecl.Tok.String(),
				PkgPath: "",
				PkgName: "",
			},
			pos: Position{
				DocBegin: int(genDecl.Doc.Pos()),
				DeclEnd:  int(genDecl.TokPos),
			},
		}
		ans = append(ans, ano)
	}
	return ans
}

// parserImports  解析import
func parserImports(importDecl *ast.GenDecl) (map[string]*IdentType, error) {
	if importDecl.Tok != token.IMPORT {
		return nil, fmt.Errorf("not import token")
	}
	importDict := make(map[string]*IdentType)
	for i := range importDecl.Specs {
		spec, ok := importDecl.Specs[i].(*ast.ImportSpec)
		if !ok {
			continue
		}
		key := ParseImportName(spec.Path.Value)
		if spec.Name != nil {
			key = spec.Name.Name
		}
		importDict[key] = &IdentType{
			Name: key,
			Path: strings.Trim(spec.Path.Value, `"`),
		}
	}
	return importDict, nil
}

// parseFuncSign  解析方法签名
func parseFuncSign(funcDecl *ast.FuncDecl) *FuncSign {
	funcSign := &FuncSign{}
	if funcDecl.Recv != nil {
		starIdent := parserStartExpr(funcDecl.Recv.List[0].Type)
		if starIdent != nil {
			funcSign.Reveiver = IdentType{
				Path: "",
				Name: starIdent.Name,
			}
		}
	}
	if funcDecl.Type.Params != nil {
		for _, field := range funcDecl.Type.Params.List {
			xIdent, selIdent := parserSelectExpr(field.Type)
			arg := new(IdentType)
			if xIdent != nil && selIdent != nil {
				arg = &IdentType{
					Path: xIdent.Name,
					Name: selIdent.Name,
				}
			}
			if xIdent == nil && selIdent != nil {
				arg = &IdentType{
					Path: "",
					Name: selIdent.Name,
				}
			}
			if arg != new(IdentType) {
				funcSign.Args = append(funcSign.Args, arg)
			}
		}
	}
	if funcDecl.Type.Results != nil {
		for _, field := range funcDecl.Type.Results.List {
			xIdent, selIdent := parserSelectExpr(field.Type)
			re := new(IdentType)
			if xIdent != nil && selIdent != nil {
				re = &IdentType{
					Path: xIdent.Name,
					Name: selIdent.Name,
				}
			}
			if xIdent == nil && selIdent != nil {
				re = &IdentType{
					Path: "",
					Name: selIdent.Name,
				}
			}
			if re != new(IdentType) {
				funcSign.Returns = append(funcSign.Returns, re)
			}
		}
	}
	return funcSign
}

func parserStartExpr(expr ast.Expr) *ast.Ident {
	startExpr, ok := expr.(*ast.StarExpr)
	if ok {
		starIdent, ok := startExpr.X.(*ast.Ident)
		if ok {
			return starIdent
		}
	}
	return nil
}

func parserSelectExpr(expr ast.Expr) (*ast.Ident, *ast.Ident) {
	selectExpr, ok := expr.(*ast.StarExpr)
	if ok {
		selectExpr, ok := selectExpr.X.(*ast.SelectorExpr)
		if !ok {
			return nil, nil
		}
		xIdent, ok := selectExpr.X.(*ast.Ident)
		if ok {
			return xIdent, selectExpr.Sel
		}
		return nil, selectExpr.Sel
	}
	return nil, nil
}

func parserCmt(cmt string) []string {
	cmt = strings.TrimLeft(cmt, "//")
	ans := strings.Split(cmt, " ")

	j := 0
	for i := 0; i < len(ans); i++ {
		if ans[i] != "" {
			ans[j] = ans[i]
			ans[i] = ""
			j++
		}
	}
	return ans[:j]
}

func ParseFile(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, filePath, nil, parser.ParseComments)
}
