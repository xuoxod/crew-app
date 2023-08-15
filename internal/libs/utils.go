package libs

import (
	"strings"
)

func StringNoSpaces(arg string) bool {
	return !strings.Contains(arg, " ")
}

func Cap(arg string) string {
	var capped string

	for i, c := range strings.Split(arg, "") {
		if i == 0 {
			capped += strings.ToUpper(c)
		} else {
			capped += c
		}
	}

	return capped
}
