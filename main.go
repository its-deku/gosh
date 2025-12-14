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

func runCmd(cmd map[string]logger.Cmd, cargs []string, delim string) string {
	// commands and args map
	cmdInfo := map[string]int{
		"ls":    3,
		"pwd":   0,
		"cd":    1,
		"touch": 3,
		"echo":  3,
	}

	_, found := cmd[cargs[0]]
	if !found {
		logger.Log(cargs[0])
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
		return cmd[cargs[0]]([]string{"."})
	}

	// check for operators
	if delim != ">" && delim != "$" && delim != "|" {
		return cmd[cargs[0]](cargs[1:])
	}

	// handle commands with redirect operators
	opInd := logger.FindStr(delim, cargs)
	preCmd := cargs[:opInd]
	output := cmd[preCmd[0]](preCmd[1:])

	isPipe := false
	if delim == "|" {
		isPipe = true
	}

	return cmds.Redirect(cmd, output, delim, cargs[opInd+1:], isPipe)
}

func parseInput(cmd map[string]logger.Cmd, s string) any {

	out, delim, err := parser.Parse(cmd, s)

	if err != nil {
		return err.Error()
	}
	logger.Log(out)

	// check if the input contains operators > or |
	if strings.Contains(s, ">") || strings.Contains(s, "$") || strings.Contains(s, "|") {
		prev := "" // stores output of previous command, useful for chaining
		for i, curr := range out {
			if curr[0] == ">" || curr[0] == "$" || curr[0] == "|" {
				arr := []string{}
				arr = append(arr, prev)

				if curr[0] == "|" {
					prev = cmd[curr[1]](arr)
				} else {
					prev = cmds.Redirect(cmd, prev, curr[0], []string{curr[1]}, false)
				}
			} else {
				prev = runCmd(cmd, curr, delim[i])
			}
		}
		return prev
		// return runCmd(cmds, out[0], delim[0])
	}

	return runCmd(cmd, out[0], "N")
}
