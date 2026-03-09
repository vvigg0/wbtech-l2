package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	mycut "github.com/vvigg0/wbtech-l2/13/internal/cutanalog"
)

func main() {
	config := mycut.Config{}
	var fields string
	var delim string
	var fileName string

	flag.StringVar(&fileName, "file", "", "Название файла, который нужно обрабатывать(по умолчанию Stdin)")

	flag.StringVar(&fields, "f", "", "Номера полей(колонок), которые нужно вывести")
	flag.StringVar(&delim, "d", "\t", "Разделитель(один символ)")
	flag.BoolVar(&config.Sep, "s", false, "Игнорировать/не игнорировать строки без разделителя")
	flag.Parse()

	err := mycut.ParseFields(fields, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ParseFields error: %v\n", err)
		os.Exit(1)
	}

	if delim == `\t` {
		delim = "\t"
	}

	runes := []rune(delim)
	if len(runes) != 1 {
		fmt.Fprintln(os.Stderr, "Разделитель должен быть одним символом")
		os.Exit(1)
	}

	config.Delim = runes[0]

	match := mycut.BuildMatcher(&config)

	var r io.Reader = os.Stdin
	if fileName != "" {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка при открытии файла: ", err)
			os.Exit(1)
		}
		defer file.Close()
		r = file
	}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		match(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Scanner error: %v", err)
		os.Exit(1)
	}
}
