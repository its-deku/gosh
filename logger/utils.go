package logger

import "strings"

func FindStr(x string, s []string) int {
	for i, v := range s {
		if strings.TrimSpace(v) == x {
			return i
		}
	}
	return -1
}
