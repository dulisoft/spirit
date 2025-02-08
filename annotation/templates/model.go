package templates

// PoolTmplData 模版文件对象
type PoolTmplData struct {
	Pkg     string
	Imports []string
	Objects []PoolObject
}

type PoolObject struct {
	Type    string //对象类型
	TypePkg string //对象的包名称
}

type RouterTmplData struct {
	Pkg     string
	Imports []string       //import
	Port    string         //监听端口
	Routers []SingleRouter //路由信息
}

type SingleRouter struct {
	Method         string //方法
	Path           string //路径
	Handler        string //处理方法
	HandlerType    string //处理方法的类
	HandlerTypePkg string //处理方法类所在的包
}
