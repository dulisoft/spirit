package main

import (
	_ "embed"
	"text/template"
)

//go:embed  relation.gtpl
var ObjectPoolTemplateText string

var permissionResrouceTemplate = template.Must(template.New("relation.gtpl").Parse(ObjectPoolTemplateText))
