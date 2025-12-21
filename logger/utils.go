package logger

import (
	"fmt"
	"strings"
)

type Job struct {
	Cmd      string
	Args     []string
	Operator string
	Opt      string
}

func PrintJob(j Job) {
	fmt.Println("Cmd: "+j.Cmd, "Args: ", j.Args, "Operator: "+j.Operator, "Opt: "+j.Opt)
}

func FindStr(x string, s []string) int {
	for i, v := range s {
		if strings.TrimSpace(v) == x {
			return i
		}
	}
	return -1
}
