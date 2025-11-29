package logger

import "fmt"

type Cmd func([]string) string

func Log(a any) {
	fmt.Println(a)
}
