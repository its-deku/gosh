package logger

import "strings"

type Job struct {
	cmd      string
	args     []string
	operator string
	opt      string
}

func FindStr(x string, s []string) int {
	for i, v := range s {
		if strings.TrimSpace(v) == x {
			return i
		}
	}
	return -1
}
