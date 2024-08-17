package utils

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`[^A-zA-я]+`)

func CleanString(data string) string {
	return re.ReplaceAllString(strings.ToLower(data), "")
}

func GetFirstValidString(params ...string) string {
	for i := range params {
		if strings.TrimSpace(params[i]) != "" {
			return params[i]
		}
	}
	return ""
}
