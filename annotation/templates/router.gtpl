package {{.Pkg}}
import (
	"net/http"
	"github.com/gin-gonic/gin"
    {{- range .Imports}}
        "{{.}}"
    {{- end}}
)

{{- /* a comment */}}

func main() {
    InitObjectPool()
	r := gin.Default()
    {{ range .Routers }}
        r.Handle("{{.Method}}", "{{.Path}}", {{.Handler}})
    {{- end}}
	if err := r.Run("0.0.0.0:{{.Port}}"); err != nil {
		panic(err)
	}
}