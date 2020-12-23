package stringer

import "sort"

// RemoveDuplicateAndBlankString 删除重复和空字符串
func RemoveDuplicateAndBlankString(strs []string) []string {
	sort.Strings(strs)
	cleanStrs := []string{}
	pre := ""
	for _, v := range strs {
		if v != pre && v != "" {
			cleanStrs = append(cleanStrs, v)
		}
		pre = v
	}
	return cleanStrs
}
