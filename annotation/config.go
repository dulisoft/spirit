package annotation

var (
	APP_NAME             = "bear"
	PORT                 = "8888"                                                    //端口
	ROUTER_PATH          = "D:/workspace/go/pro/personal/spirit/annotation/examples" //路由文件的名称
	TEMPLATE_DIR         = "D:/workspace/go/pro/personal/spirit/annotation/templates"
	ROUTER_TEMPLATE_NAME = "router.tmpl"
)

const (
	RouterAnnotation     = "@Router"
	ControllerAnnotation = "@Controller"
)

func getRouterTemplatePath() string {
	return TEMPLATE_DIR + "/" + ROUTER_TEMPLATE_NAME
}

func getMainCodePath() string {
	return ROUTER_PATH + "/" + APP_NAME + "/main.go"
}
