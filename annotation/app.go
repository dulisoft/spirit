package annotation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dulisoft/spirit/annotation/templates"
)

type App struct {
	Service     string
	Port        string
	servicePath string
	genPath     string //对象池初始化Go文件在项目中的默认路径
	tmplPath    string //对象文件模版文件路径
	Annotations []*Annotation
	workers     map[string][]AnnoWork
}

func NewApp(path string) *App {
	return &App{
		servicePath: path,
		Service:     APP_NAME,
		Port:        PORT,
		genPath:     path,
		tmplPath:    TEMPLATE_DIR,
		Annotations: []*Annotation{},
		workers:     make(map[string][]AnnoWork),
	}
}

func (a *App) Work() error {
	//解析项目的注解
	if err := a.parese(); err != nil {
		return err
	}
	//初始化对象池
	if err := a.WritePool(); err != nil {
		return err
	}
	//写main
	if err := a.WriteMain(); err != nil {
		return err
	}
	//执行标注
	for _, anno := range a.Annotations {
		workers := a.workers[anno.Name]
		for _, worker := range workers {
			if err := worker(anno); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *App) RegisterAnnotation(anno string, worker AnnoWork) {
	a.workers[anno] = append(a.workers[anno], worker)
}

func (a *App) parese() error {
	//读取go.mod，获得该项目的path
	servicePkg, err := a.serviceMode()
	if err != nil {
		return err
	}
	//遍历所有的go文件,获取标注信息
	filepath.Walk(a.servicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		//处理文件
		annos, _, err := ParseGoFileDecls(path)
		if err != nil {
			return err
		}
		//补充一些信息
		for _, anno := range annos {
			anno.obj.PkgPath = a.pkgPath(servicePkg, path)
		}
		a.Annotations = append(a.Annotations, annos...)
		return nil
	})
	return nil
}

func (a *App) serviceMode() (string, error) {
	modPath := a.goModePath()
	file, err := os.Open(modPath)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "module") {
			line = strings.TrimLeft(line, "module")
			line = strings.TrimSpace(line)
			return line, nil
		}
	}
	return "", fmt.Errorf("invalid service")
}

func (a *App) WritePool() error {
	tmplData := &templates.PoolTmplData{
		Pkg: "main",
	}

	objUniqueDict := make(map[string]any)
	for _, anno := range a.Annotations {
		if anno.Kind != METHOD {
			continue
		}
		tmplData.Imports = append(tmplData.Imports, anno.obj.Import())
		obj := templates.PoolObject{
			Type:    anno.obj.FuncSign.Reveiver.Name,
			TypePkg: anno.obj.PkgName,
		}
		key := fmt.Sprintf("%v-%v", obj.Type, obj.TypePkg)
		if _, has := objUniqueDict[key]; !has {
			tmplData.Objects = append(tmplData.Objects, obj)
			objUniqueDict[key] = 1
		}
	}
	//去重
	tmplData.Imports = UniqueSlice(tmplData.Imports)
	return WirteTemplate(templates.ObjectPoolTemplate, tmplData, a.poolGoPath())
}

func (a *App) WriteMain() error {
	tmplData := &templates.RouterTmplData{
		Pkg:  "main",
		Port: a.Port,
	}
	for _, anno := range a.Annotations {
		//收集路由注解
		if anno.Name != ANNOTATION_ROUTER {
			continue
		}
		tmplData.Imports = append(tmplData.Imports, anno.obj.Import())
		if len(anno.Props) < 2 {
			panic(fmt.Sprintf("invalid router annotation %v", anno.obj.PkgPath))
		}
		tmplData.Routers = append(tmplData.Routers, templates.SingleRouter{
			Method:         parseMethod(anno.Props[1]),
			Path:           anno.Props[0],
			Handler:        anno.obj.Name,
			HandlerType:    anno.obj.FuncSign.Reveiver.Name,
			HandlerTypePkg: anno.obj.PkgName,
		})
	}
	return WirteTemplate(templates.RouterTemplate, tmplData, a.mainGoPath())
}

// pkgPath  根据当前的包名称，项目包路径，文件路径，推断完整的包路径
func (a *App) pkgPath(pkgPath string, filePath string) string {
	filePath = strings.Replace(filePath, "\\", "/", -1)
	servicePath := strings.Replace(a.servicePath, "\\", "/", -1)
	fileServicePath, _ := strings.CutPrefix(filePath, servicePath)
	pathes := strings.Split(fileServicePath, "/")
	fileServicePath = strings.Join(pathes[:len(pathes)-1], "/")
	return strings.TrimLeft(strings.Replace(pkgPath+"/"+fileServicePath, "//", "/", -1), "/")
}

func (a *App) goModePath() string {
	return a.servicePath + "/go.mod"
}

func (a *App) poolGoPath() string {
	return a.genPath + "/object_pool.go"
}

func (a *App) mainGoPath() string {
	return a.genPath + "/main.go"
}
