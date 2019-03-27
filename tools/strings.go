package tools

import "strings"

//NotEmptyAll is not empty for all
func NotEmptyAll(str ...string) bool {
	if len(str) == 0 {
		return false
	}
	for i := 0; i < len(str); i++ {
		if str[i] == "" {
			return false
		}
	}
	return true
}

//MapToString is used for k8s select string from map convert
func MapToString(labels map[string]string) string {
	var returnStringArr []string
	for index, value := range labels {
		singleLabel := index + "=" + value
		returnStringArr = append(returnStringArr, singleLabel)
	}
	return strings.Join(returnStringArr, ",")
}
