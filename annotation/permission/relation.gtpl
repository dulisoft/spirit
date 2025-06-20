package {{.Pkg}}

import (
	"github.com/gin-gonic/gin"
    "github.com/dulisoft/spirit/annotation"
    {{- range .Imports}}
    "{{- .}}"
    {{- end}}
)

{{- /* a comment */}}

func main() {
    InitObjectPool()
	r := gin.Default()
    {{- range .Routers }}
        r.Handle("{{- .Method}}", "{{.Path}}", annotation.GetInstance[{{.HandlerTypePkg}}.{{.HandlerType}}]().{{.Handler}})
    {{- end}}
	if err := r.Run("0.0.0.0:{{.Port}}"); err != nil {
		panic(err)
	}
}