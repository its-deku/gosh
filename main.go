package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mod.org/shellit/cmds"
	"mod.org/shellit/logger"
	"mod.org/shellit/parser"
)

func main() {
	commands := map[string]logger.Cmd{
		"cd":    cmds.Cd,
		"pwd":   cmds.Pwd,
		"ls":    cmds.Ls,
		"touch": cmds.Touch,
		"clear": cmds.Clear,
		"echo":  cmds.Echo,
	}
	// repl loop
	for {
		fmt.Print("$ ")
		inp, err := getInput()
		inp = strings.Trim(inp, "\n")
		if err != nil || inp == "exit" {
			return
		}
		logger.Log(parseInput(commands, inp))
	}
}

func getInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	inp, err := reader.ReadString('\n')

	if err != nil {
		logger.Log("Error while parsing input command !")
		return "", err
	}

	return inp, nil
}

func runCmd(cmds map[string]logger.Cmd, cargs []string) string {
	// commands and args map
	cmdInfo := map[string]int{
		"ls":    3,
		"pwd":   0,
		"cd":    1,
		"touch": 3,
		"echo":  3,
	}

	_, found := cmds[cargs[0]]
	if !found {
		logger.Log("Command not found!")
		return ""
	}

	ln := len(cargs)
	if cmdInfo[cargs[0]] != 3 {
		if cmdInfo[cargs[0]] != ln-1 {
			return "too many arguments..."
		}
	}

	if ln == 1 {
		return cmds[cargs[0]]([]string{"."})
	}

	return cmds[cargs[0]](cargs[1:])
}

func parseInput(cmds map[string]logger.Cmd, s string) any {

	out, delim, err := parser.Parse(cmds, s)

	// check if the input contains operators > or |
	if strings.Contains(s, ">") || strings.Contains(s, ">>") || strings.Contains(s, "|") {
		if err != nil {
			return err.Error()
		}

		prev := ""
		for i := range out {
			if i != 0 {
				for v := range strings.SplitSeq(prev, " ") {
					out[i] = append(out[i], v)
				}
			}
			prev = runCmd(cmds, out[i])
		}

		logger.Log(delim)
	}

	return runCmd(cmds, out[0])
}
