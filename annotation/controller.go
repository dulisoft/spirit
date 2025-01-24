package annotation

// Controller  集合多个Router的struct
type Controller struct {
	Imports []string
	Pkg     string
	Name    string
	Routers []Router
}

// Router  控制器，表示该方法是一个api路由处理器
// 将路由注册写到指定模版中
type Router struct {
	Receiver string
	Method   string
	Path     string
	Name     string
}

func NewController(routers []Router) *Controller {
	return &Controller{
		Imports: []string{},
		Pkg:     "main",
		Name:    "",
		Routers: routers,
	}
}
