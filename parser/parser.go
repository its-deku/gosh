package parser

import (
	"errors"
	"strings"

	"mod.org/shellit/logger"
)

func getDelim(toks map[string]bool, s string) int {
	ind := -1
	// delim := ""
	for i, v := range s {
		if toks[string(v)] {
			// delim = string(v)
			ind = i
			break
		}
	}
	return ind
}

func parseTokAndCmd(toks map[string]bool, comb string, tok string) logger.Job {
	logger.Log(comb)
	comb = strings.TrimSpace(comb)
	ind := getDelim(toks, comb)

	cmdOrOpt := ""
	ptr := 0
	args := []string{}

	// seperates command and args before the operator
	// if ind == -1 {
	for str := range strings.SplitSeq(sanitize(comb), " ") {
		if str != "" && str != " " {
			if ptr == 0 && !toks[str] {
				cmdOrOpt = str
				ptr += 1
				continue
			}
			args = append(args, str)
		}
	}
	// }
	// else {
	// 	for str := range strings.SplitSeq(sanitize(comb[ind:]), " ") {
	// 		if str != "" && str != " " {
	// 			args = append(args, str)
	// 		}
	// 	}
	// }

	opt := ""
	switch tok {
	case ">", "$":
		if ind != 0 {
			tok = ""
		} else {
			opt = cmdOrOpt
			args = []string{}
		}
		cmdOrOpt = ""
	case "|":
		args = []string{}
	default:
		tok = ""
	}

	job := logger.Job{
		Cmd:      cmdOrOpt,
		Args:     args,
		Operator: tok,
		Opt:      opt,
	}

	logger.PrintJob(job)
	return job
}

func Parse(cmds map[string]logger.Cmd, stream string) ([]logger.Job, error) {
	tokens := map[string]bool{
		">": true,
		"$": true, // substitute for >>
		"|": true,
	}

	stream = sanitize(stream) // removes whitespace and replaces >> with $

	if stream == "" {
		return []logger.Job{}, errors.New(stream)
	}

	var job logger.Job

	// commands without operators
	if !strings.Contains(stream, ">") && !strings.Contains(stream, "$") && !strings.Contains(stream, "|") {
		arr := []string{}
		for str := range strings.SplitSeq(stream, " ") {
			if str != "" && str != " " {
				arr = append(arr, str)
			}
		}

		job.Cmd = arr[0]
		job.Args = arr[1:]

		return []logger.Job{job}, nil
	}

	// check if the stream starts/ends with an operator
	if tokens[string(stream[0])] || tokens[string(stream[len(stream)-1])] {
		return []logger.Job{}, errors.New("not a super valid command")
	}

	// store the individual commands seperately
	parts := []string{}
	pInd := 0
	tok := ""
	for i, s := range stream {
		if tokens[string(s)] {
			// if len(tok) == 0 {
			// 	tok = string(s)
			// } else {
			// 	parts = append(parts, stream[pInd:i])
			// 	pInd = i
			// }
			tok = string(s)
			parts = append(parts, stream[pInd:i])
			parseTokAndCmd(tokens, parts[len(parts)-1], string(stream[pInd]))
			pInd = i
		}
	}
	parts = append(parts, stream[pInd:])
	parseTokAndCmd(tokens, parts[len(parts)-1], tok)

	// delim := []string{}
	// out := [][]string{}
	// for _, v := range parts {
	// 	arr, del := parseTokAndCmd(tokens, v)
	// 	out = append(out, arr)
	// 	delim = append(delim, del)
	// }
	// return out, delim, nil
	return []logger.Job{}, nil
}

func sanitize(str string) string {
	str = strings.TrimSpace(str)

	// check if input is empty
	if len(str) == 0 {
		return str
	}

	return strings.ReplaceAll(str, ">>", "$")
}
