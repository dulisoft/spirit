package {{.Pkg}}
import (
    "github.com/dulisoft/spirit/annotation"
    {{- range .Imports}}
        "{{.}}"
    {{- end}}
)

func InitObjectPool(){
    {{ range .Objects }}
        annotation.Put[{{.Type}}](new({{.Obj}}))
    {{- end}}
}