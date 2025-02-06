package annotation

var (
	APP_NAME             = "bear"
	PORT                 = "8888"                                                    //端口
	ROUTER_PATH          = "D:/workspace/go/pro/personal/spirit/annotation/examples" //路由文件的名称
	TEMPLATE_DIR         = "templates"
	ROUTER_TEMPLATE_NAME = "router.tmpl"
)

const (
	ANNOTATION_PREFIX = "@"
	ANNOTATION_ROUTER = "Router"
)

const (
	STATIC    = iota //静态变量
	VARIABLE         //变量
	STRUCT           //结构体
	METHOD           //方法
	FUNCTIONS        //函数
)
