package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/term"
)

// Shell представляет интерактивную командную оболочку.
type Shell struct {
	in        io.Reader
	out, errw io.Writer
}

// New создаёт новый экземпляр Shell с заданными потоками ввода, вывода и ошибок.
func New(in io.Reader, out, errw io.Writer) *Shell {
	return &Shell{in, out, errw}
}

// Run запускает основной цикл чтения и выполнения команд.
// Возвращает итоговый код завершения Shell.
func (s *Shell) Run() int {
	fd := int(os.Stdin.Fd())
	interactive := term.IsTerminal(fd)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)
	sc := bufio.NewScanner(s.in)
	code := 0
	for {
		if interactive {
			fmt.Fprint(s.out, "$ ")
		}
		if !sc.Scan() {
			return code
		}
		line := sc.Text()
		c, exit := s.runLine(line)
		code = c
		if exit != nil {
			return exit.Code
		}
	}
}

func (s *Shell) runLine(line string) (int, *exitRequest) {
	ts, err := tokenize(line)
	if err != nil {
		fmt.Fprintln(s.errw, err)
		return 1, nil
	}
	if len(ts) == 0 {
		return 0, nil
	}

	expr, err := parse(ts)
	if err != nil {
		fmt.Fprintln(s.errw, err)
		return 1, nil
	}
	expand(&expr)

	code, exit := runExpression(expr, s.out, s.errw)
	return code, exit
}
