package annotation

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

// GenPosition  代码生成助手，不同的位置生成有不同的实现
type GenPosition interface {
	Write(code []byte) error
}

// Annotation 注解信息，可以方便包含所有的信息，方便生成代码
type Annotation struct {
	pkg   string        //文件的包名称
	path  string        //文件路径
	f     *ast.File     //注解文件
	fdel  *ast.FuncDecl //被注解的方法
	Props []string      //键值信息
}

func PareserAnnotations(filePath string) ([]Router, error) {
	anns, err := ParseGoFileFunc(filePath)
	if err != nil {
		fmt.Printf("error :%v", err)
		return nil, err
	}
	routers := make([]Router, 0)
	for _, anno := range anns {
		router := Router{
			Receiver: getFuncReceiver(anno.fdel.Recv),
			Method:   strings.ToUpper(strings.TrimRight(strings.TrimLeft(anno.Props[2], "["), "]")),
			Path:     anno.Props[1],
			Name:     anno.fdel.Name.Name,
		}
		routers = append(routers, router)
	}
	return routers, nil
}

func ParseGoFileFunc(filePath string) ([]*Annotation, error) {
	file, err := ParseFile(filePath)
	if err != nil {
		return nil, err
	}

	annotations := make([]*Annotation, 0)
	for _, decl := range file.Decls {
		//解析方法
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Doc == nil {
				continue
			}
			comment := f.Doc.Text()
			if !strings.Contains(comment, "Router") {
				continue
			}

			posStr := fmt.Sprintf("%v", f.Body.Rbrace)
			pos, _ := strconv.Atoi(posStr)

			if err := WriteAtNextLine(filePath, pos, "//this is c comment"); err != nil {
				fmt.Printf("error: %v", err)
			}
			annotations = append(annotations, &Annotation{
				pkg:   file.Name.Name,
				path:  filePath,
				f:     file,
				fdel:  f,
				Props: parserComment(f.Doc.List),
			})
		}
	}
	return annotations, nil
}

func parserComment(comments []*ast.Comment) []string {
	for _, comment := range comments {
		cmt := comment.Text
		if !strings.Contains(cmt, "Router") {
			continue
		}
		return parserController(cmt)
	}
	return []string{}
}

func parserController(cmt string) []string {
	cmt = strings.TrimLeft(cmt, "//")
	ans := strings.Split(cmt, " ")

	j := 0
	for i := 0; i < len(ans); i++ {
		if ans[i] != "" {
			ans[j] = ans[i]
			j++
		}
	}
	ans = ans[:j+1]
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
