package templates

import (
	_ "embed"
	"text/template"
)

//go:embed  pool.gtpl
var ObjectPoolTemplateText string

//go:embed  router.gtpl
var RouterTemplateText string

var (
	ObjectPoolTemplate = template.Must(template.New("pool.gtpl").Parse(ObjectPoolTemplateText))
	RouterTemplate     = template.Must(template.New("router.gtpl").Parse(RouterTemplateText))
)
