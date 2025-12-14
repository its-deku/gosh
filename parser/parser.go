package parser

import (
	"errors"
	"strings"

	"mod.org/shellit/logger"
)

func getDelim(toks map[string]bool, s string) (int, string) {
	ind := -1
	delim := ""
	for i, v := range s {
		if toks[string(v)] {
			delim = string(v)
			ind = i
			break
		}
	}
	return ind, delim
}

func parseTokAndCmd(toks map[string]bool, comb string) ([]string, string) {
	comb = strings.TrimSpace(comb)
	ind, delim := getDelim(toks, comb)

	command := []string{}
	// seperates command and args before the operator
	if ind != 0 {
		for str := range strings.SplitSeq(sanitize(comb[:ind]), " ") {
			if str != "" && str != " " {
				command = append(command, str)
			}
		}
	}

	command = append(command, delim)
	command = append(command, strings.TrimSpace(comb[ind+1:]))

	return command, delim
}

func Parse(cmds map[string]logger.Cmd, stream string) ([][]string, []string, error) {
	tokens := map[string]bool{
		">": true,
		"$": true, // substitute for >>
		"|": true,
	}

	stream = sanitize(stream)

	if stream == "" {
		return nil, nil, errors.New(stream)
	}

	if !strings.Contains(stream, ">") && !strings.Contains(stream, "$") && !strings.Contains(stream, "|") {
		arr := []string{}
		for str := range strings.SplitSeq(stream, " ") {
			if str != "" && str != " " {
				arr = append(arr, str)
			}
		}

		return [][]string{arr}, nil, nil
	}

	// check if the stream starts/ends with an operator
	if tokens[string(stream[0])] || tokens[string(stream[len(stream)-1])] {
		return nil, nil, errors.New("not a super valid command")
	}

	// store the individual commands seperately
	parts := []string{}
	pInd := 0
	tok := ""
	for i, s := range stream {
		if tokens[string(s)] {
			if len(tok) == 0 {
				tok = string(s)
			} else {
				parts = append(parts, stream[pInd:i])
				pInd = i
			}
		}
	}
	parts = append(parts, stream[pInd:])

	delim := []string{}
	out := [][]string{}
	for _, v := range parts {
		arr, del := parseTokAndCmd(tokens, v)
		out = append(out, arr)
		delim = append(delim, del)
	}
	return out, delim, nil
}

func sanitize(str string) string {
	i, out := 0, ""
	str = strings.TrimSpace(str)

	// check if input is empty
	if len(str) == 0 {
		return str
	}

	for {
		if i >= len(str)-1 {
			break
		}
		if string(str[i]) == ">" && string(str[i+1]) == ">" {
			out += "$"
			i += 1
		} else {
			out += string(str[i])
		}
		i += 1
	}
	return out + string(str[len(str)-1])
}
