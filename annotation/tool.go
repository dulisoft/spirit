package annotation

//TypePack 反射出类型的具体包名称
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
