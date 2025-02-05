package annotation

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

const ANNOTATION_PREFIX = "@"
const (
	STATIC   = iota //静态变量
	VARIABLE        //变量
	STRUCT          //结构体
	METHOD          //方法
	FUNC            //函数
)

var anoRegexp, _ = regexp.Compile("@\\w+")

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
	Pkg      string    //包名
	Path     string    //路径
	FuncSign *FuncSign //方法签名
}

type FuncSign struct {
	Reveiver string   //方法的接受者
	Args     []string //参数类型数组
	Resps    []string //返回值
}

func ParseGoFileFunc(filePath string) ([]*Annotation, error) {
	astFile, err := ParseFile(filePath)
	if err != nil {
		return nil, err
	}
	annotations := make([]*Annotation, 0)
	for _, decl := range astFile.Decls {
		switch f := decl.(type) {
		case *ast.FuncDecl:
			annotations = append(annotations, parseFuncDecl(f)...)
		case *ast.GenDecl:
			annotations = append(annotations, parseGenDecl(f)...)
		default:
			panic(fmt.Sprintf("sytax error,position: %v", f.Pos()))
		}
	}
	return annotations, nil
}

// parseFuncDecl 解析方法
func parseFuncDecl(funcDecl *ast.FuncDecl) (ans []*Annotation) {
	if funcDecl.Doc == nil {
		return
	}
	comment := funcDecl.Doc.Text()
	if !strings.Contains(comment, ANNOTATION_PREFIX) {
		return
	}
	for _, cmt := range funcDecl.Doc.List {
		cmtLine := cmt.Text
		if strings.Contains(cmtLine, ANNOTATION_PREFIX) {
			continue
		}
		cmts := parserCmt(cmtLine)
		if len(cmts) <= 0 {
			continue
		}
		ano := &Annotation{
			Name:  cmts[0],
			Kind:  0,
			Props: cmts[1:],
			obj: Object{
				Name: funcDecl.Name.Name,
				Pkg:  "",
				Path: "",
			},
			pos: Position{
				DocBegin: int(funcDecl.Doc.Pos()),
				DeclEnd:  int(funcDecl.Body.End()),
			},
		}
		ans = append(ans, ano)
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
			Name:  cmts[0],
			Kind:  0,
			Props: cmts[1:],
			obj: Object{
				Name: genDecl.Tok.String(),
				Pkg:  "",
				Path: "",
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

func parserCmt(cmt string) []string {
	cmt = strings.TrimLeft(cmt, "//")
	ans := strings.Split(cmt, " ")

	j := 0
	for i := 0; i < len(ans); i++ {
		if ans[i] != "" {
			ans[j] = ans[i]
			j++
		}
	}
	return ans
}

func getFuncReceiver(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) <= 0 {
		return ""
	}
	startExpr, ok := recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return ""
	}
	ident, ok := startExpr.X.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func ParseFile(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, filePath, nil, parser.ParseComments)
}
