package tools

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
	var returnVar string
	for index, value := range labels {
		singleLabel := index + "=" + value + ","
		returnVar += singleLabel
	}
	return returnVar[:len(returnVar)-1]
}
