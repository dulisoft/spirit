package main

import (
	"github.com/samber/lo"
	"go/ast"
	"os"
	"strings"
)

const (
	PermissionTag = "@Permission"
	RouterTag     = "@Router"
)

type ProjectParser struct {
	FilesPath string
	file      *os.File
}

func NewProjectParser(path string) (*ProjectParser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &ProjectParser{
		FilesPath: path,
		file:      file,
	}, nil
}

type PermissionAnnotation struct {
	RouterCmt     string `json:"router_cmt"`
	PermissionCmt string `json:"permission_cmt"`
}

func (p *PermissionAnnotation) Parse() *PermissionResource {
	//routerTag
	rs := strings.SplitAfter(p.RouterCmt, RouterTag)
	rts := strings.Split(rs[1], " ")
	rts = lo.Filter(rts, func(item string, index int) bool { return item != "" })
	//permissionTag
	ps := strings.SplitAfter(p.PermissionCmt, PermissionTag)
	pts := strings.Split(ps[1], " ")
	pts = lo.Filter(pts, func(item string, index int) bool { return item != "" })

	return &PermissionResource{
		Path:           rts[0],
		Method:         rts[1],
		Action:         pts[0],
		PermissionName: pts[1],
		Scope:          pts[2],
	}
}

type PermissionResource struct {
	Path           string `json:"path"`
	Method         string `json:"method"`
	Action         string `json:"action"`
	PermissionName string `json:"permission_name"`
	Scope          string `json:"scope"`
}

func (p *ProjectParser) ReadAnnotationSlice() []*PermissionResource {
	astFile, err := ParseFile(p.FilesPath)
	if err != nil {
		return nil
	}
	annos := make([]*PermissionResource, 0)
	for _, decl := range astFile.Decls {
		switch f := decl.(type) {
		case *ast.FuncDecl:
			routerCmts := parseFuncCmt(f, RouterTag)
			permissionCmts := parseFuncCmt(f, PermissionTag)
			if len(routerCmts) > 0 && len(permissionCmts) > 0 {
				pa := &PermissionAnnotation{
					RouterCmt:     routerCmts[0],
					PermissionCmt: permissionCmts[0],
				}
				annos = append(annos, pa.Parse())
			}
		}
	}
	return annos
}
