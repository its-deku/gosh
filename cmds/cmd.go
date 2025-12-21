package cmds

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"mod.org/shellit/logger"
)

func Cd(args []string) string {
	err := os.Chdir(args[0])
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

func Sleep(args []string) string {
	duration, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		logger.Log(err)
	}

	time.Sleep(time.Duration(duration) * time.Second)
	return "slept for " + fmt.Sprint(duration) + " seconds"
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
	// cmd := exec.Command("tput", "clear")
	// rows, err := cmd.Output()
	// if err != nil {
	// 	return err.Error()
	// }
	// for row := 0; row <= int(rows[0]); row++ {
	// 	logger.Log("\n")
	// }
	logger.Log("\033[2J\033[H")
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
		logger.Log("while appending :()")
		return err.Error()
	}
	f.Close()
	return ""
}

func Redirect(cmd map[string]logger.Cmd, data string, delim string, fileName []string, isPipe bool) string {

	// check if it is a pipe
	if isPipe {
		logger.Log(fileName[0])
		return cmd[fileName[0]](strings.Split(data, " "))
	}

	if len(fileName) > 2 {
		return "Only one file can be specified!"
	}

	files := strings.Split(Ls([]string{"."}), " ")

	isPresent := slices.Contains(files, fileName[0])

	var fileCreate string
	if !isPresent {
		fileCreate = Touch(fileName)
	}

	// write data to the file
	if len(fileCreate) > 1 {
		return "Issue with file creation!"
	}

	if delim == ">" {
		return overWrite(fileName[0], []byte(data+"\n"))
	}
	return appendToFile(fileName[0], []byte(data+"\n"))
}
