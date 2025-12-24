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
			jobs = append(jobs, subParser(parts[i], true))
		}
		jobs = append(jobs, subParser(parts[len(parts)-1], false))
	} else {
		jobs = append(jobs, subParser(stream, false))
	}

	return jobs
}

func subParser(stream string, bg bool) logger.Job2 {
	pInd := 0
	commands := []string{}
	for i, v := range stream {
		if string(v) == "|" {
			commands = append(commands, stream[pInd:i])
			pInd = i + 1
		}
	}
	commands = append(commands, stream[pInd:])

	return logger.Job2{
		Background: bg,
		Execs:      commands,
	}
}
