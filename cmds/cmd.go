package cmds

import (
	"os"
	"os/exec"
	"slices"
	"strings"

	"mod.org/shellit/logger"
)

func Cd(dir []string) string {
	err := os.Chdir(dir[0])
	if err != nil {
		return err.Error()
	}
	// os.Stdout.Write([]byte(s))
	return ""
}

func Pwd(args []string) string {
	cur, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	return cur
}

func Echo(args []string) string {
	echo := ""
	for i, v := range args {
		if i == 0 {
			echo += v
		} else {
			echo += " " + v
		}
	}

	return echo
}

func Ls(args []string) string {
	output := []string{}
	for _, dir := range args {
		files, err := os.ReadDir(dir)
		if err != nil {
			if err.Error() == "open "+dir+": not a directory" {
				logger.Log(dir)
			} else {
				logger.Log(err)
			}
			continue
		}

		for _, file := range files {
			output = append(output, file.Name())
		}
	}

	var strOut string

	for i, str := range output {
		if i > 0 {
			strOut += " " + str
		} else {
			strOut += str
		}
	}

	return strOut
}

func Touch(args []string) string {
	for _, file := range args {
		_, err := os.Create(file)
		if err != nil {
			return err.Error()
		}
	}
	return ""
}

func Clear(args []string) string {
	cmd := exec.Command("clear", args...)
	err := cmd.Run()
	if err != nil {
		return err.Error()
	}
	return ""
}

func overWrite(file string, data []byte) string {
	err := os.WriteFile(file, data, 0666)
	if err != nil {
		return err.Error()
	}

	return ""
}

func appendToFile(file string, data []byte) string {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err.Error()
	}

	_, err = f.Write(data)
	if err != nil {
		logger.Log("while appendign :()")
		return err.Error()
	}
	f.Close()
	return ""
}

func Redirect(commands map[string]logger.Cmd, s string, delim string) string {
	ind := strings.Index(s, delim)
	cmd := strings.Split(s[:ind], " ")
	file := strings.Split(s[ind+1:], " ")[1:]

	if len(file) > 2 {
		return "Only one file can be specified!"
	}

	cmdOut := commands[cmd[0]](cmd[1 : len(cmd)-1])
	files := strings.Split(Ls([]string{"."}), " ")

	isPresent := slices.Contains(files, file[0])

	var fileCreate string
	if !isPresent {
		fileCreate = Touch(file)
	}

	// write data to the file
	if len(fileCreate) > 1 {
		return "Issue with file creation!"
	}

	if delim == ">" {
		return overWrite(file[0], []byte(cmdOut+"\n"))
	}
	return appendToFile(file[0], []byte(cmdOut+"\n"))
}
