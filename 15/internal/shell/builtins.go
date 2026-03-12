package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var errExit = errors.New("exit requested")

type exitRequest struct {
	Code int
}

func runBuiltin(cmd command, out, errw io.Writer) (handled bool, code int, exit *exitRequest) {
	if len(cmd.argv) == 0 {
		return false, 1, nil
	}
	if cmd.argv[0] == "exit" {
		c := 0
		if len(cmd.argv) > 1 {
			v, e := strconv.Atoi(cmd.argv[1])
			if e != nil {
				fmt.Fprintln(errw, "exit: invalid code")
				return true, 1, nil
			}
			c = v
		}
		return true, 0, &exitRequest{c}
	}
	switch cmd.argv[0] {
	case "cd":
		code = builtinCd(cmd.argv, errw)
	case "pwd":
		code = builtinPwd(out, errw)
	case "echo":
		code = builtinEcho(cmd.argv, out)
	case "kill":
		code = builtinKill(cmd.argv, errw)
	default:
		return false, 1, nil
	}
	return true, code, nil
}
func builtinCd(argv []string, errw io.Writer) int {
	path := ""
	if len(argv) < 2 {
		h, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(errw, "cd:", err)
			return 1
		}
		path = h
	} else {
		path = argv[1]
	}
	if err := os.Chdir(path); err != nil {
		fmt.Fprintln(errw, "cd:", err)
		return 1
	}
	return 0
}

func builtinPwd(out, errw io.Writer) int {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(errw, "pwd:", err)
	}
	fmt.Fprintln(out, wd)
	return 0
}

func builtinEcho(argv []string, out io.Writer) int {
	if len(argv) <= 1 {
		fmt.Fprintln(out)
		return 0
	}
	fmt.Fprintln(out, strings.Join(argv[1:], " "))
	return 0
}

func builtinKill(argv []string, errw io.Writer) int {
	if len(argv) < 2 {
		fmt.Fprintln(errw, "kill: pid required")
		return 1
	}
	pid, err := strconv.Atoi(argv[1])
	if err != nil || pid <= 0 {
		fmt.Fprintln(errw, "kill: invalid pid")
		return 1
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintln(errw, "kill:", err)
		return 1
	}
	err = process.Kill()
	if err != nil {
		fmt.Fprintln(errw, "kill:", err)
		return 1
	}
	return 0
}
