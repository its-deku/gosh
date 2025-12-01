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

func parseTokAndCmd(toks map[string]bool, comb string) []string {
	comb = strings.TrimSpace(comb)
	ind, delim := getDelim(toks, comb)

	command := []string{}
	if ind != 0 {
		command = append(command, strings.TrimSpace(comb[:ind]))
	}
	command = append(command, " "+delim+" ")
	command = append(command, strings.TrimSpace(comb[ind+1:]))

	return command
}

func Parse(cmds map[string]logger.Cmd, stream string) ([][]string, error) {
	tokens := map[string]bool{
		">": true,
		"$": true, // substitute for >>
		"|": true,
	}

	stream = sanitize(stream)

	// check if the stream starts/ends with an operator
	if tokens[string(stream[0])] || tokens[string(stream[len(stream)-1])] {
		return nil, errors.New("not a super valid command")
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

	out := [][]string{}
	for _, v := range parts {
		out = append(out, parseTokAndCmd(tokens, v))
	}
	return out, nil
}

func sanitize(str string) string {
	i, out := 0, ""
	str = strings.TrimSpace(str)
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
