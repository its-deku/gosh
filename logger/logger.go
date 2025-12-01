package logger

import "fmt"

type Cmd func([]string) string

func Log(a any) {
	fmt.Println(a)
}

func LogSeq(arr []string) {
	for i, v := range arr {
		if i == 0 {
			fmt.Print("[" + v + ", ")
		} else if i == len(arr)-1 {
			fmt.Print(v + "]")
		} else {
			fmt.Print(v + ", ")
		}
	}
	fmt.Println()
}
