package annotation

import (
	"strings"
)

// TypePack 反射出类型的具体包名称
func TypePack[T any]() string {
	return ""
}

func UniqueSlice(ds []string) []string {
	results := make([]string, 0)
	uniqueDict := make(map[string]any)
	for i := range ds {
		if _, has := uniqueDict[ds[i]]; !has {
			results = append(results, ds[i])
			uniqueDict[ds[i]] = 1
		}
	}
	return results
}

func ParseImportName(path string) string {
	path = strings.Trim(path, `"`)
	ims := strings.Split(path, "/")
	return ims[len(ims)-1]
}

func parseMethod(m string) string {
	return strings.ToUpper(strings.Trim(m, "[]"))
}
