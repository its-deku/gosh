package parser

import (
	"strings"

	"mod.org/shellit/logger"
)

func ShellParse(stream string) []logger.Job2 {
	stream = sanitize(stream)

	jobs := []logger.Job2{}

	if strings.Contains(stream, "&") {
		parts := strings.Split(stream, "&")
		for i := 0; i < len(parts)-1; i++ {
			parts[i] = sanitize(parts[i])
			jobs = append(jobs, subParser(parts[i], true))
		}
		if string(stream[len(stream)-1]) != "&" {
			parts[len(parts)-1] = sanitize(parts[len(parts)-1])
			jobs = append(jobs, subParser(parts[len(parts)-1], false))
		}
	} else {
		jobs = append(jobs, subParser(stream, false))
	}

	return jobs
}

func parseRedirect(cmd string) (string, string) {
	var op string

	if strings.Contains(cmd, ">") {
		op = ">"
	} else {
		op = "$"
	}

	out := []string{}

	for str := range strings.SplitSeq(cmd, op) {
		out = append(out, str)
	}

	return strings.Join(out[:len(out)-1], " "), strings.Join(out[len(out)-1:], " ")
}

func subParser(stream string, bg bool) logger.Job2 {
	pInd := 0
	var arg, exec string
	commands := []string{}
	opt := []string{}
	stream += " |"

	for i, v := range stream {
		arg = ""
		if string(v) == "|" {
			cmd := sanitize(stream[pInd:i])

			// check for redirect operators in command
			if strings.Contains(cmd, ">") || strings.Contains(cmd, "$") {
				exec, arg = parseRedirect(cmd)
				commands = append(commands, exec)
			} else {
				commands = append(commands, cmd)
			}
			opt = append(opt, arg)

			pInd = i + 1
		}
	}

	// commands = append(commands, sanitize(stream[pInd:]))

	return logger.Job2{
		Background: bg,
		Execs:      commands,
		Opt:        opt,
	}
}
